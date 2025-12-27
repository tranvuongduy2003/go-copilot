package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var (
	ErrCacheMiss    = errors.New("cache miss")
	ErrCacheExpired = errors.New("cache expired")
	ErrInvalidData  = errors.New("invalid cache data")
)

type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, key string) (bool, error)
	TTL(ctx context.Context, key string) (time.Duration, error)
	Clear(ctx context.Context, pattern string) error
}

type TypedCache[T any] struct {
	cache  Cache
	prefix string
}

func NewTypedCache[T any](cache Cache, prefix string) *TypedCache[T] {
	return &TypedCache[T]{
		cache:  cache,
		prefix: prefix,
	}
}

func (c *TypedCache[T]) key(key string) string {
	if c.prefix == "" {
		return key
	}
	return c.prefix + ":" + key
}

func (c *TypedCache[T]) Get(ctx context.Context, key string) (T, error) {
	var zero T
	data, err := c.cache.Get(ctx, c.key(key))
	if err != nil {
		return zero, err
	}

	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		return zero, fmt.Errorf("%w: %v", ErrInvalidData, err)
	}

	return result, nil
}

func (c *TypedCache[T]) Set(ctx context.Context, key string, value T, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal cache value: %w", err)
	}

	return c.cache.Set(ctx, c.key(key), data, ttl)
}

func (c *TypedCache[T]) Delete(ctx context.Context, keys ...string) error {
	prefixedKeys := make([]string, len(keys))
	for i, k := range keys {
		prefixedKeys[i] = c.key(k)
	}
	return c.cache.Delete(ctx, prefixedKeys...)
}

func (c *TypedCache[T]) Exists(ctx context.Context, key string) (bool, error) {
	return c.cache.Exists(ctx, c.key(key))
}

func (c *TypedCache[T]) GetOrSet(ctx context.Context, key string, ttl time.Duration, loader func(ctx context.Context) (T, error)) (T, error) {
	value, err := c.Get(ctx, key)
	if err == nil {
		return value, nil
	}

	if !errors.Is(err, ErrCacheMiss) {
		return value, err
	}

	value, err = loader(ctx)
	if err != nil {
		return value, err
	}

	if setErr := c.Set(ctx, key, value, ttl); setErr != nil {
		return value, nil
	}

	return value, nil
}

func (c *TypedCache[T]) Clear(ctx context.Context) error {
	return c.cache.Clear(ctx, c.prefix+":*")
}

type CacheStats struct {
	Hits       int64 `json:"hits"`
	Misses     int64 `json:"misses"`
	HitRate    float64 `json:"hit_rate"`
	TotalKeys  int64 `json:"total_keys"`
	MemoryUsed int64 `json:"memory_used"`
}

type TTLConfig struct {
	Default time.Duration
	Short   time.Duration
	Medium  time.Duration
	Long    time.Duration
}

func DefaultTTLConfig() TTLConfig {
	return TTLConfig{
		Default: 5 * time.Minute,
		Short:   1 * time.Minute,
		Medium:  15 * time.Minute,
		Long:    1 * time.Hour,
	}
}
