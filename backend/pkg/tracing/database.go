package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	databaseTracerName = "database"
)

type DatabaseSpanConfig struct {
	DBSystem    string
	DBName      string
	DBUser      string
	ServerAddr  string
	ServerPort  int
}

func StartDatabaseSpan(ctx context.Context, operation string, table string, config DatabaseSpanConfig) (context.Context, trace.Span) {
	tracer := otel.Tracer(databaseTracerName)

	attrs := []attribute.KeyValue{
		semconv.DBSystemPostgreSQL,
		semconv.DBCollectionName(table),
		attribute.String("db.operation", operation),
	}

	if config.DBName != "" {
		attrs = append(attrs, semconv.DBNamespace(config.DBName))
	}

	if config.DBUser != "" {
		attrs = append(attrs, attribute.String("db.user", config.DBUser))
	}

	if config.ServerAddr != "" {
		attrs = append(attrs, semconv.ServerAddress(config.ServerAddr))
	}

	if config.ServerPort > 0 {
		attrs = append(attrs, semconv.ServerPort(config.ServerPort))
	}

	spanName := operation + " " + table

	return tracer.Start(ctx, spanName,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attrs...),
	)
}

func RecordDatabaseError(span trace.Span, err error) {
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
}

func AddDatabaseStatement(span trace.Span, statement string) {
	span.SetAttributes(semconv.DBQueryText(statement))
}

func AddDatabaseRowsAffected(span trace.Span, rowsAffected int64) {
	span.SetAttributes(attribute.Int64("db.rows_affected", rowsAffected))
}

type TracedOperation func(ctx context.Context) error

func TraceDatabase(ctx context.Context, operation string, table string, config DatabaseSpanConfig, fn TracedOperation) error {
	ctx, span := StartDatabaseSpan(ctx, operation, table, config)
	defer span.End()

	err := fn(ctx)
	if err != nil {
		RecordDatabaseError(span, err)
	}
	return err
}

type TracedOperationWithResult[T any] func(ctx context.Context) (T, error)

func TraceDatabaseWithResult[T any](ctx context.Context, operation string, table string, config DatabaseSpanConfig, fn TracedOperationWithResult[T]) (T, error) {
	ctx, span := StartDatabaseSpan(ctx, operation, table, config)
	defer span.End()

	result, err := fn(ctx)
	if err != nil {
		RecordDatabaseError(span, err)
	}
	return result, err
}
