package security

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/auth"
)

const tokenBlacklistKeyPrefix = "token:blacklist:"

type redisTokenBlacklist struct {
	client *redis.Client
}

func NewRedisTokenBlacklist(client *redis.Client) auth.TokenBlacklist {
	return &redisTokenBlacklist{client: client}
}

func (blacklist *redisTokenBlacklist) Add(ctx context.Context, tokenID string, expiresAt int64) error {
	key := blacklist.buildKey(tokenID)

	ttl := time.Until(time.Unix(expiresAt, 0))
	if ttl <= 0 {
		return nil
	}

	err := blacklist.client.Set(ctx, key, "1", ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to add token to blacklist: %w", err)
	}

	return nil
}

func (blacklist *redisTokenBlacklist) IsBlacklisted(ctx context.Context, tokenID string) (bool, error) {
	key := blacklist.buildKey(tokenID)

	exists, err := blacklist.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check token blacklist: %w", err)
	}

	return exists > 0, nil
}

func (blacklist *redisTokenBlacklist) buildKey(tokenID string) string {
	return tokenBlacklistKeyPrefix + tokenID
}
