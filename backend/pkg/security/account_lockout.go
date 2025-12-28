package security

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	lockoutKeyPrefix       = "lockout:"
	attemptsKeyPrefix      = "login_attempts:"
	DefaultMaxAttempts     = 5
	DefaultLockoutDuration = 15 * time.Minute
	DefaultAttemptWindow   = 15 * time.Minute
)

type AccountLockoutConfig struct {
	MaxAttempts     int
	LockoutDuration time.Duration
	AttemptWindow   time.Duration
}

func DefaultAccountLockoutConfig() AccountLockoutConfig {
	return AccountLockoutConfig{
		MaxAttempts:     DefaultMaxAttempts,
		LockoutDuration: DefaultLockoutDuration,
		AttemptWindow:   DefaultAttemptWindow,
	}
}

type AccountLockout interface {
	IsLocked(ctx context.Context, identifier string) (bool, time.Duration, error)
	RecordFailedAttempt(ctx context.Context, identifier string) (int, error)
	ResetAttempts(ctx context.Context, identifier string) error
	GetAttemptCount(ctx context.Context, identifier string) (int, error)
}

type RedisAccountLockout struct {
	client *redis.Client
	config AccountLockoutConfig
}

func NewRedisAccountLockout(client *redis.Client, config AccountLockoutConfig) *RedisAccountLockout {
	if config.MaxAttempts <= 0 {
		config.MaxAttempts = DefaultMaxAttempts
	}
	if config.LockoutDuration <= 0 {
		config.LockoutDuration = DefaultLockoutDuration
	}
	if config.AttemptWindow <= 0 {
		config.AttemptWindow = DefaultAttemptWindow
	}

	return &RedisAccountLockout{
		client: client,
		config: config,
	}
}

func (lockout *RedisAccountLockout) IsLocked(ctx context.Context, identifier string) (bool, time.Duration, error) {
	lockKey := lockout.buildLockKey(identifier)

	ttl, err := lockout.client.TTL(ctx, lockKey).Result()
	if err != nil {
		return false, 0, fmt.Errorf("failed to check lockout status: %w", err)
	}

	if ttl > 0 {
		return true, ttl, nil
	}

	return false, 0, nil
}

func (lockout *RedisAccountLockout) RecordFailedAttempt(ctx context.Context, identifier string) (int, error) {
	attemptsKey := lockout.buildAttemptsKey(identifier)

	pipe := lockout.client.Pipeline()
	incrCmd := pipe.Incr(ctx, attemptsKey)
	pipe.Expire(ctx, attemptsKey, lockout.config.AttemptWindow)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to record failed attempt: %w", err)
	}

	attempts := int(incrCmd.Val())

	if attempts >= lockout.config.MaxAttempts {
		lockKey := lockout.buildLockKey(identifier)
		err = lockout.client.Set(ctx, lockKey, "locked", lockout.config.LockoutDuration).Err()
		if err != nil {
			return attempts, fmt.Errorf("failed to set lockout: %w", err)
		}

		lockout.client.Del(ctx, attemptsKey)
	}

	return attempts, nil
}

func (lockout *RedisAccountLockout) ResetAttempts(ctx context.Context, identifier string) error {
	attemptsKey := lockout.buildAttemptsKey(identifier)
	lockKey := lockout.buildLockKey(identifier)

	pipe := lockout.client.Pipeline()
	pipe.Del(ctx, attemptsKey)
	pipe.Del(ctx, lockKey)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to reset attempts: %w", err)
	}

	return nil
}

func (lockout *RedisAccountLockout) GetAttemptCount(ctx context.Context, identifier string) (int, error) {
	attemptsKey := lockout.buildAttemptsKey(identifier)

	count, err := lockout.client.Get(ctx, attemptsKey).Int()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get attempt count: %w", err)
	}

	return count, nil
}

func (lockout *RedisAccountLockout) buildLockKey(identifier string) string {
	return lockoutKeyPrefix + identifier
}

func (lockout *RedisAccountLockout) buildAttemptsKey(identifier string) string {
	return attemptsKeyPrefix + identifier
}
