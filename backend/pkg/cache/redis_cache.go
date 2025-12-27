package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
	"github.com/tranvuongduy2003/go-copilot/pkg/metrics"
)

type RedisCache struct {
	client *redis.Client
	logger logger.Logger
}

func NewRedisCache(client *redis.Client, logger logger.Logger) *RedisCache {
	return &RedisCache{
		client: client,
		logger: logger,
	}
}

func (c *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	start := time.Now()

	result, err := c.client.Get(ctx, key).Bytes()
	duration := time.Since(start).Seconds()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			metrics.RecordCacheOperation("get", false, duration)
			return nil, ErrCacheMiss
		}
		metrics.RecordCacheOperation("get", false, duration)
		c.logger.Error("redis get error",
			logger.String("key", key),
			logger.Err(err),
		)
		return nil, err
	}

	metrics.RecordCacheOperation("get", true, duration)
	return result, nil
}

func (c *RedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	start := time.Now()

	err := c.client.Set(ctx, key, value, ttl).Err()
	duration := time.Since(start).Seconds()

	if err != nil {
		metrics.RecordCacheOperation("set", false, duration)
		c.logger.Error("redis set error",
			logger.String("key", key),
			logger.Err(err),
		)
		return err
	}

	metrics.RecordCacheOperation("set", true, duration)
	return nil
}

func (c *RedisCache) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	start := time.Now()
	err := c.client.Del(ctx, keys...).Err()
	duration := time.Since(start).Seconds()

	if err != nil {
		metrics.RecordCacheOperation("delete", false, duration)
		c.logger.Error("redis delete error",
			logger.Strings("keys", keys),
			logger.Err(err),
		)
		return err
	}

	metrics.RecordCacheOperation("delete", true, duration)
	return nil
}

func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	start := time.Now()
	result, err := c.client.Exists(ctx, key).Result()
	duration := time.Since(start).Seconds()

	if err != nil {
		metrics.RecordCacheOperation("exists", false, duration)
		return false, err
	}

	metrics.RecordCacheOperation("exists", result > 0, duration)
	return result > 0, nil
}

func (c *RedisCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	return c.client.TTL(ctx, key).Result()
}

func (c *RedisCache) Clear(ctx context.Context, pattern string) error {
	start := time.Now()

	var cursor uint64
	var keysToDelete []string

	for {
		var keys []string
		var err error
		keys, cursor, err = c.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			metrics.RecordCacheOperation("clear", false, time.Since(start).Seconds())
			return err
		}

		keysToDelete = append(keysToDelete, keys...)

		if cursor == 0 {
			break
		}
	}

	if len(keysToDelete) > 0 {
		if err := c.client.Del(ctx, keysToDelete...).Err(); err != nil {
			metrics.RecordCacheOperation("clear", false, time.Since(start).Seconds())
			return err
		}
	}

	metrics.RecordCacheOperation("clear", true, time.Since(start).Seconds())
	c.logger.Info("cache cleared",
		logger.String("pattern", pattern),
		logger.Int("keys_deleted", len(keysToDelete)),
	)

	return nil
}

func (c *RedisCache) GetMulti(ctx context.Context, keys ...string) (map[string][]byte, error) {
	if len(keys) == 0 {
		return make(map[string][]byte), nil
	}

	start := time.Now()
	results, err := c.client.MGet(ctx, keys...).Result()
	duration := time.Since(start).Seconds()

	if err != nil {
		metrics.RecordCacheOperation("mget", false, duration)
		return nil, err
	}

	output := make(map[string][]byte, len(keys))
	for i, result := range results {
		if result != nil {
			if str, ok := result.(string); ok {
				output[keys[i]] = []byte(str)
			}
		}
	}

	metrics.RecordCacheOperation("mget", true, duration)
	return output, nil
}

func (c *RedisCache) SetMulti(ctx context.Context, items map[string][]byte, ttl time.Duration) error {
	if len(items) == 0 {
		return nil
	}

	start := time.Now()
	pipe := c.client.Pipeline()

	for key, value := range items {
		pipe.Set(ctx, key, value, ttl)
	}

	_, err := pipe.Exec(ctx)
	duration := time.Since(start).Seconds()

	if err != nil {
		metrics.RecordCacheOperation("mset", false, duration)
		return err
	}

	metrics.RecordCacheOperation("mset", true, duration)
	return nil
}

func (c *RedisCache) Stats(ctx context.Context) (*CacheStats, error) {
	info, err := c.client.Info(ctx, "stats", "memory", "keyspace").Result()
	if err != nil {
		return nil, err
	}

	c.logger.Debug("redis info", logger.String("info", info))

	return &CacheStats{}, nil
}
