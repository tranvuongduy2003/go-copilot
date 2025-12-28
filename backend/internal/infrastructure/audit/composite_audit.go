package audit

import (
	"context"
)

type CompositeAuditLogger struct {
	loggers []AuditLogger
}

func NewCompositeAuditLogger(loggers ...AuditLogger) *CompositeAuditLogger {
	return &CompositeAuditLogger{
		loggers: loggers,
	}
}

func (composite *CompositeAuditLogger) Log(ctx context.Context, entry AuditEntry) error {
	var lastError error
	for _, auditLogger := range composite.loggers {
		if err := auditLogger.Log(ctx, entry); err != nil {
			lastError = err
		}
	}
	return lastError
}

func (composite *CompositeAuditLogger) AddLogger(auditLogger AuditLogger) {
	composite.loggers = append(composite.loggers, auditLogger)
}
