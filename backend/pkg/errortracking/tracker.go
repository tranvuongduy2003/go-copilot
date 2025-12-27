package errortracking

import (
	"context"
	"runtime"
	"time"

	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

type Severity string

const (
	SeverityDebug   Severity = "debug"
	SeverityInfo    Severity = "info"
	SeverityWarning Severity = "warning"
	SeverityError   Severity = "error"
	SeverityFatal   Severity = "fatal"
)

type ErrorContext struct {
	UserID      string
	RequestID   string
	Path        string
	Method      string
	Tags        map[string]string
	Extra       map[string]interface{}
	Fingerprint []string
}

type ErrorTracker interface {
	CaptureError(ctx context.Context, err error, errorContext *ErrorContext)
	CaptureMessage(ctx context.Context, message string, severity Severity, errorContext *ErrorContext)
	AddBreadcrumb(ctx context.Context, category, message string, data map[string]interface{})
	Flush(timeout time.Duration)
}

type LogErrorTracker struct {
	logger logger.Logger
}

func NewLogErrorTracker(logger logger.Logger) *LogErrorTracker {
	return &LogErrorTracker{
		logger: logger,
	}
}

func (tracker *LogErrorTracker) CaptureError(ctx context.Context, err error, errorContext *ErrorContext) {
	fields := tracker.buildFields(errorContext)
	fields = append(fields, logger.Err(err))
	fields = append(fields, logger.String("stack", getStackTrace()))

	tracker.logger.Error("error captured", fields...)
}

func (tracker *LogErrorTracker) CaptureMessage(ctx context.Context, message string, severity Severity, errorContext *ErrorContext) {
	fields := tracker.buildFields(errorContext)
	fields = append(fields, logger.String("severity", string(severity)))

	switch severity {
	case SeverityDebug:
		tracker.logger.Debug(message, fields...)
	case SeverityInfo:
		tracker.logger.Info(message, fields...)
	case SeverityWarning:
		tracker.logger.Warn(message, fields...)
	case SeverityError, SeverityFatal:
		tracker.logger.Error(message, fields...)
	default:
		tracker.logger.Info(message, fields...)
	}
}

func (tracker *LogErrorTracker) AddBreadcrumb(ctx context.Context, category, message string, data map[string]interface{}) {
	fields := []logger.Field{
		logger.String("category", category),
		logger.String("breadcrumb", message),
	}
	for key, value := range data {
		fields = append(fields, logger.Any(key, value))
	}
	tracker.logger.Debug("breadcrumb", fields...)
}

func (tracker *LogErrorTracker) Flush(timeout time.Duration) {
	if err := tracker.logger.Sync(); err != nil {
		tracker.logger.Warn("failed to flush logger", logger.Err(err))
	}
}

func (tracker *LogErrorTracker) buildFields(errorContext *ErrorContext) []logger.Field {
	var fields []logger.Field

	if errorContext == nil {
		return fields
	}

	if errorContext.UserID != "" {
		fields = append(fields, logger.String("user_id", errorContext.UserID))
	}
	if errorContext.RequestID != "" {
		fields = append(fields, logger.String("request_id", errorContext.RequestID))
	}
	if errorContext.Path != "" {
		fields = append(fields, logger.String("path", errorContext.Path))
	}
	if errorContext.Method != "" {
		fields = append(fields, logger.String("method", errorContext.Method))
	}

	for key, value := range errorContext.Tags {
		fields = append(fields, logger.String("tag."+key, value))
	}

	for key, value := range errorContext.Extra {
		fields = append(fields, logger.Any("extra."+key, value))
	}

	return fields
}

func getStackTrace() string {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

type NoopErrorTracker struct{}

func NewNoopErrorTracker() *NoopErrorTracker {
	return &NoopErrorTracker{}
}

func (tracker *NoopErrorTracker) CaptureError(ctx context.Context, err error, errorContext *ErrorContext) {
}

func (tracker *NoopErrorTracker) CaptureMessage(ctx context.Context, message string, severity Severity, errorContext *ErrorContext) {
}

func (tracker *NoopErrorTracker) AddBreadcrumb(ctx context.Context, category, message string, data map[string]interface{}) {
}

func (tracker *NoopErrorTracker) Flush(timeout time.Duration) {}
