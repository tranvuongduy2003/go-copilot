package postgres

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tranvuongduy2003/go-copilot/pkg/config"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

const (
	defaultMaxRetries     = 5
	defaultBaseDelay      = 1 * time.Second
	defaultMaxDelay       = 30 * time.Second
	defaultHealthCheckPeriod = 30 * time.Second
)

type DB struct {
	pool   *pgxpool.Pool
	config *config.DatabaseConfig
	logger logger.Logger
}

type DBStats struct {
	TotalConns       int32
	AcquiredConns    int32
	IdleConns        int32
	ConstructingConns int32
	MaxConns         int32
	AcquireCount     int64
	AcquireDuration  time.Duration
	CanceledAcquires int64
	EmptyAcquires    int64
}

type DBOption func(*dbOptions)

type dbOptions struct {
	maxRetries      int
	baseDelay       time.Duration
	maxDelay        time.Duration
	healthCheckPeriod time.Duration
}

func WithMaxRetries(retries int) DBOption {
	return func(o *dbOptions) {
		o.maxRetries = retries
	}
}

func WithBaseDelay(delay time.Duration) DBOption {
	return func(o *dbOptions) {
		o.baseDelay = delay
	}
}

func WithMaxDelay(delay time.Duration) DBOption {
	return func(o *dbOptions) {
		o.maxDelay = delay
	}
}

func WithHealthCheckPeriod(period time.Duration) DBOption {
	return func(o *dbOptions) {
		o.healthCheckPeriod = period
	}
}

func NewDB(cfg *config.DatabaseConfig, log logger.Logger, opts ...DBOption) *DB {
	return &DB{
		config: cfg,
		logger: log,
	}
}

func (db *DB) Connect(ctx context.Context, opts ...DBOption) error {
	options := &dbOptions{
		maxRetries:      defaultMaxRetries,
		baseDelay:       defaultBaseDelay,
		maxDelay:        defaultMaxDelay,
		healthCheckPeriod: defaultHealthCheckPeriod,
	}

	for _, opt := range opts {
		opt(options)
	}

	poolConfig, err := pgxpool.ParseConfig(db.config.PgxConnString())
	if err != nil {
		return newDBError("parse connection string", errors.Join(ErrParseConfig, err))
	}

	poolConfig.MaxConns = int32(db.config.MaxOpenConns)
	poolConfig.MinConns = int32(db.config.MaxIdleConns)
	poolConfig.MaxConnLifetime = db.config.ConnMaxLifetime
	poolConfig.MaxConnIdleTime = db.config.ConnMaxLifetime / 2
	poolConfig.HealthCheckPeriod = options.healthCheckPeriod

	var pool *pgxpool.Pool
	var lastErr error

	for attempt := 0; attempt <= options.maxRetries; attempt++ {
		if attempt > 0 {
			delay := calculateBackoff(attempt, options.baseDelay, options.maxDelay)
			db.logger.Warn("retrying database connection",
				logger.Int("attempt", attempt),
				logger.Duration("delay", delay),
				logger.Err(lastErr),
			)

			select {
			case <-ctx.Done():
				return newDBError("connect", errors.Join(ErrConnectionCancelled, ctx.Err()))
			case <-time.After(delay):
			}
		}

		db.logger.Info("attempting database connection",
			logger.Int("attempt", attempt+1),
			logger.Int("max_attempts", options.maxRetries+1),
		)

		pool, err = pgxpool.NewWithConfig(ctx, poolConfig)
		if err != nil {
			lastErr = err
			continue
		}

		if err = pool.Ping(ctx); err != nil {
			pool.Close()
			lastErr = err
			continue
		}

		db.pool = pool
		db.logger.Info("database connection established",
			logger.String("host", db.config.Host),
			logger.Int("port", db.config.Port),
			logger.String("database", db.config.Name),
		)
		return nil
	}

	return newDBError("connect", errors.Join(ErrConnectionFailed, lastErr))
}

func (db *DB) Close() {
	if db.pool != nil {
		db.pool.Close()
		db.logger.Info("database connection closed")
	}
}

func (db *DB) Ping(ctx context.Context) error {
	if db.pool == nil {
		return ErrPoolNotInitialized
	}
	return db.pool.Ping(ctx)
}

func (db *DB) Stats() DBStats {
	if db.pool == nil {
		return DBStats{}
	}

	stat := db.pool.Stat()
	return DBStats{
		TotalConns:       stat.TotalConns(),
		AcquiredConns:    stat.AcquiredConns(),
		IdleConns:        stat.IdleConns(),
		ConstructingConns: stat.ConstructingConns(),
		MaxConns:         stat.MaxConns(),
		AcquireCount:     stat.AcquireCount(),
		AcquireDuration:  stat.AcquireDuration(),
		CanceledAcquires: stat.CanceledAcquireCount(),
		EmptyAcquires:    stat.EmptyAcquireCount(),
	}
}

func (db *DB) Pool() *pgxpool.Pool {
	return db.pool
}

func calculateBackoff(attempt int, baseDelay, maxDelay time.Duration) time.Duration {
	delay := baseDelay * time.Duration(1<<uint(attempt-1))
	if delay > maxDelay {
		delay = maxDelay
	}

	jitter := time.Duration(rand.Int63n(int64(delay) / 2))
	return delay + jitter
}
