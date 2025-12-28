package security

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const passwordResetKeyPrefix = "password_reset:"

type RedisPasswordResetTokenStore struct {
	client *redis.Client
}

func NewRedisPasswordResetTokenStore(client *redis.Client) *RedisPasswordResetTokenStore {
	return &RedisPasswordResetTokenStore{client: client}
}

func (store *RedisPasswordResetTokenStore) Store(ctx context.Context, email string, tokenHash string, expiresAt time.Time) error {
	key := store.buildKey(tokenHash)
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return fmt.Errorf("expiration time must be in the future")
	}

	err := store.client.Set(ctx, key, email, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to store password reset token: %w", err)
	}

	return nil
}

func (store *RedisPasswordResetTokenStore) Get(ctx context.Context, tokenHash string) (string, error) {
	key := store.buildKey(tokenHash)

	email, err := store.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", fmt.Errorf("failed to get password reset token: %w", err)
	}

	return email, nil
}

func (store *RedisPasswordResetTokenStore) Delete(ctx context.Context, tokenHash string) error {
	key := store.buildKey(tokenHash)

	err := store.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete password reset token: %w", err)
	}

	return nil
}

func (store *RedisPasswordResetTokenStore) buildKey(tokenHash string) string {
	return passwordResetKeyPrefix + tokenHash
}
