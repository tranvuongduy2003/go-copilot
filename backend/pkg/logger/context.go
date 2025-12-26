package logger

import (
	"context"
)

type contextKey string

const (
	loggerKey    contextKey = "logger"
	requestIDKey contextKey = "request_id"
	userIDKey    contextKey = "user_id"
)

func WithContext(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

func FromContext(ctx context.Context) Logger {
	if l, ok := ctx.Value(loggerKey).(Logger); ok {
		return l
	}
	return L()
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

func RequestIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}
	return ""
}

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func UserIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(userIDKey).(string); ok {
		return id
	}
	return ""
}

func Ctx(ctx context.Context) Logger {
	l := FromContext(ctx)

	var fields []Field

	if requestID := RequestIDFromContext(ctx); requestID != "" {
		fields = append(fields, String("request_id", requestID))
	}

	if userID := UserIDFromContext(ctx); userID != "" {
		fields = append(fields, String("user_id", userID))
	}

	if len(fields) > 0 {
		return l.With(fields...)
	}

	return l
}
