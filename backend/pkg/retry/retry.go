package retry

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"time"
)

var ErrMaxRetriesExceeded = errors.New("maximum retries exceeded")

type Config struct {
	MaxRetries      int
	InitialInterval time.Duration
	MaxInterval     time.Duration
	Multiplier      float64
	Jitter          float64
}

func DefaultConfig() Config {
	return Config{
		MaxRetries:      3,
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     10 * time.Second,
		Multiplier:      2.0,
		Jitter:          0.1,
	}
}

type RetryableFunc func(ctx context.Context) error

type RetryableFuncWithResult[T any] func(ctx context.Context) (T, error)

func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	var retryableError *RetryableError
	if errors.As(err, &retryableError) {
		return retryableError.Retryable
	}

	return true
}

type RetryableError struct {
	Err       error
	Retryable bool
}

func (e *RetryableError) Error() string {
	return e.Err.Error()
}

func (e *RetryableError) Unwrap() error {
	return e.Err
}

func NewRetryableError(err error, retryable bool) *RetryableError {
	return &RetryableError{
		Err:       err,
		Retryable: retryable,
	}
}

func NonRetryable(err error) *RetryableError {
	return NewRetryableError(err, false)
}

func Retryable(err error) *RetryableError {
	return NewRetryableError(err, true)
}

func Do(ctx context.Context, config Config, operation RetryableFunc) error {
	var lastErr error

	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		lastErr = operation(ctx)
		if lastErr == nil {
			return nil
		}

		if !IsRetryable(lastErr) {
			return lastErr
		}

		if attempt == config.MaxRetries {
			break
		}

		sleepDuration := calculateBackoff(config, attempt)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(sleepDuration):
		}
	}

	return errors.Join(ErrMaxRetriesExceeded, lastErr)
}

func DoWithResult[T any](ctx context.Context, config Config, operation RetryableFuncWithResult[T]) (T, error) {
	var result T
	var lastErr error

	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		if err := ctx.Err(); err != nil {
			return result, err
		}

		result, lastErr = operation(ctx)
		if lastErr == nil {
			return result, nil
		}

		if !IsRetryable(lastErr) {
			return result, lastErr
		}

		if attempt == config.MaxRetries {
			break
		}

		sleepDuration := calculateBackoff(config, attempt)

		select {
		case <-ctx.Done():
			return result, ctx.Err()
		case <-time.After(sleepDuration):
		}
	}

	return result, errors.Join(ErrMaxRetriesExceeded, lastErr)
}

func calculateBackoff(config Config, attempt int) time.Duration {
	backoff := float64(config.InitialInterval) * math.Pow(config.Multiplier, float64(attempt))

	if backoff > float64(config.MaxInterval) {
		backoff = float64(config.MaxInterval)
	}

	if config.Jitter > 0 {
		jitter := backoff * config.Jitter * (2*rand.Float64() - 1)
		backoff += jitter
	}

	if backoff < 0 {
		backoff = 0
	}

	return time.Duration(backoff)
}

type Retryer struct {
	config Config
}

func New(config Config) *Retryer {
	return &Retryer{config: config}
}

func NewDefault() *Retryer {
	return New(DefaultConfig())
}

func (retryer *Retryer) Do(ctx context.Context, operation RetryableFunc) error {
	return Do(ctx, retryer.config, operation)
}

func (retryer *Retryer) DoWithResult(ctx context.Context, operation RetryableFuncWithResult[any]) (any, error) {
	return DoWithResult(ctx, retryer.config, operation)
}
