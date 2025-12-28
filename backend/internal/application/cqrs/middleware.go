package cqrs

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
	"github.com/tranvuongduy2003/go-copilot/pkg/validator"
)

type CommandMiddleware func(next CommandDispatcher) CommandDispatcher

type QueryMiddleware func(next QueryDispatcher) QueryDispatcher

type CommandDispatcher func(ctx context.Context, command Command) (interface{}, error)

type QueryDispatcher func(ctx context.Context, query Query) (interface{}, error)

func LoggingCommandMiddleware(log logger.Logger) CommandMiddleware {
	return func(next CommandDispatcher) CommandDispatcher {
		return func(ctx context.Context, command Command) (interface{}, error) {
			commandType := reflect.TypeOf(command).String()
			startTime := time.Now()

			log.Info("executing command",
				logger.String("command_type", commandType),
			)

			result, err := next(ctx, command)
			duration := time.Since(startTime)

			if err != nil {
				log.Error("command failed",
					logger.String("command_type", commandType),
					logger.Duration("duration", duration),
					logger.Err(err),
				)
			} else {
				log.Info("command completed",
					logger.String("command_type", commandType),
					logger.Duration("duration", duration),
				)
			}

			return result, err
		}
	}
}

func LoggingQueryMiddleware(log logger.Logger) QueryMiddleware {
	return func(next QueryDispatcher) QueryDispatcher {
		return func(ctx context.Context, query Query) (interface{}, error) {
			queryType := reflect.TypeOf(query).String()
			startTime := time.Now()

			log.Debug("executing query",
				logger.String("query_type", queryType),
			)

			result, err := next(ctx, query)
			duration := time.Since(startTime)

			if err != nil {
				log.Error("query failed",
					logger.String("query_type", queryType),
					logger.Duration("duration", duration),
					logger.Err(err),
				)
			} else {
				log.Debug("query completed",
					logger.String("query_type", queryType),
					logger.Duration("duration", duration),
				)
			}

			return result, err
		}
	}
}

func ValidationCommandMiddleware(validate *validator.Validator) CommandMiddleware {
	return func(next CommandDispatcher) CommandDispatcher {
		return func(ctx context.Context, command Command) (interface{}, error) {
			if err := validate.Validate(command); err != nil {
				return nil, fmt.Errorf("command validation failed: %w", err)
			}
			return next(ctx, command)
		}
	}
}

func ValidationQueryMiddleware(validate *validator.Validator) QueryMiddleware {
	return func(next QueryDispatcher) QueryDispatcher {
		return func(ctx context.Context, query Query) (interface{}, error) {
			if err := validate.Validate(query); err != nil {
				return nil, fmt.Errorf("query validation failed: %w", err)
			}
			return next(ctx, query)
		}
	}
}

type TransactionManager interface {
	Begin(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

func TransactionCommandMiddleware(transactionManager TransactionManager) CommandMiddleware {
	return func(next CommandDispatcher) CommandDispatcher {
		return func(ctx context.Context, command Command) (result interface{}, err error) {
			transactionContext, beginError := transactionManager.Begin(ctx)
			if beginError != nil {
				return nil, fmt.Errorf("failed to begin transaction: %w", beginError)
			}

			defer func() {
				if panicValue := recover(); panicValue != nil {
					_ = transactionManager.Rollback(transactionContext)
					panic(panicValue)
				}
			}()

			result, err = next(transactionContext, command)
			if err != nil {
				if rollbackError := transactionManager.Rollback(transactionContext); rollbackError != nil {
					return nil, fmt.Errorf("command failed: %w, rollback error: %v", err, rollbackError)
				}
				return nil, err
			}

			if commitError := transactionManager.Commit(transactionContext); commitError != nil {
				return nil, fmt.Errorf("failed to commit transaction: %w", commitError)
			}

			return result, nil
		}
	}
}

type MetricsRecorder interface {
	RecordCommandDuration(commandType string, duration time.Duration, success bool)
	RecordQueryDuration(queryType string, duration time.Duration, success bool)
}

func MetricsCommandMiddleware(recorder MetricsRecorder) CommandMiddleware {
	return func(next CommandDispatcher) CommandDispatcher {
		return func(ctx context.Context, command Command) (interface{}, error) {
			commandType := reflect.TypeOf(command).String()
			startTime := time.Now()

			result, err := next(ctx, command)

			duration := time.Since(startTime)
			recorder.RecordCommandDuration(commandType, duration, err == nil)

			return result, err
		}
	}
}

func MetricsQueryMiddleware(recorder MetricsRecorder) QueryMiddleware {
	return func(next QueryDispatcher) QueryDispatcher {
		return func(ctx context.Context, query Query) (interface{}, error) {
			queryType := reflect.TypeOf(query).String()
			startTime := time.Now()

			result, err := next(ctx, query)

			duration := time.Since(startTime)
			recorder.RecordQueryDuration(queryType, duration, err == nil)

			return result, err
		}
	}
}

func ChainCommandMiddleware(middlewares ...CommandMiddleware) CommandMiddleware {
	return func(final CommandDispatcher) CommandDispatcher {
		for i := len(middlewares) - 1; i >= 0; i-- {
			final = middlewares[i](final)
		}
		return final
	}
}

func ChainQueryMiddleware(middlewares ...QueryMiddleware) QueryMiddleware {
	return func(final QueryDispatcher) QueryDispatcher {
		for i := len(middlewares) - 1; i >= 0; i-- {
			final = middlewares[i](final)
		}
		return final
	}
}
