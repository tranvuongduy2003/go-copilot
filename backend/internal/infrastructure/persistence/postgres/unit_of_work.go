package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type txKey struct{}

type Querier interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type Transaction interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type UnitOfWork interface {
	Begin(ctx context.Context) (UnitOfWorkContext, error)
	BeginTx(ctx context.Context, opts pgx.TxOptions) (UnitOfWorkContext, error)
}

type UnitOfWorkContext interface {
	Transaction
	Context() context.Context
	Querier() Querier
}

type pgxTransaction struct {
	tx  pgx.Tx
	ctx context.Context
}

func (t *pgxTransaction) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

func (t *pgxTransaction) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}

func (t *pgxTransaction) Context() context.Context {
	return t.ctx
}

func (t *pgxTransaction) Querier() Querier {
	return t.tx
}

type PostgresUnitOfWork struct {
	pool *pgxpool.Pool
}

func NewUnitOfWork(pool *pgxpool.Pool) *PostgresUnitOfWork {
	return &PostgresUnitOfWork{pool: pool}
}

func (uow *PostgresUnitOfWork) Begin(ctx context.Context) (UnitOfWorkContext, error) {
	return uow.BeginTx(ctx, pgx.TxOptions{})
}

func (uow *PostgresUnitOfWork) BeginTx(ctx context.Context, opts pgx.TxOptions) (UnitOfWorkContext, error) {
	tx, err := uow.pool.BeginTx(ctx, opts)
	if err != nil {
		return nil, newDBError("begin transaction", errors.Join(ErrBeginTransaction, err))
	}

	txCtx := context.WithValue(ctx, txKey{}, tx)
	return &pgxTransaction{
		tx:  tx,
		ctx: txCtx,
	}, nil
}

func WithTransaction(ctx context.Context, pool *pgxpool.Pool, fn func(ctx context.Context) error) error {
	return WithTransactionOptions(ctx, pool, pgx.TxOptions{}, fn)
}

func WithTransactionOptions(ctx context.Context, pool *pgxpool.Pool, opts pgx.TxOptions, fn func(ctx context.Context) error) error {
	tx, err := pool.BeginTx(ctx, opts)
	if err != nil {
		return newDBError("begin transaction", errors.Join(ErrBeginTransaction, err))
	}

	txCtx := context.WithValue(ctx, txKey{}, tx)

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback(ctx)
			panic(r)
		}
	}()

	if err := fn(txCtx); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return newDBError("rollback transaction", errors.Join(ErrRollbackTransaction, rbErr, err))
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return newDBError("commit transaction", errors.Join(ErrCommitTransaction, err))
	}

	return nil
}

func TxFromContext(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	return tx, ok
}

func InjectTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func GetQuerier(ctx context.Context, pool *pgxpool.Pool) Querier {
	if tx, ok := TxFromContext(ctx); ok {
		return tx
	}
	return pool
}
