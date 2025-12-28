package audit

import (
	"context"
	"encoding/json"

	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

type LoggerAuditLogger struct {
	logger logger.Logger
}

func NewLoggerAuditLogger(log logger.Logger) *LoggerAuditLogger {
	return &LoggerAuditLogger{
		logger: log,
	}
}

func (auditLogger *LoggerAuditLogger) Log(ctx context.Context, entry AuditEntry) error {
	metadata, _ := json.Marshal(entry.Metadata)

	auditLogger.logger.Info("AUDIT",
		logger.String("audit_id", entry.ID.String()),
		logger.String("event_type", entry.EventType),
		logger.String("user_id", entry.UserID.String()),
		logger.String("action", entry.Action),
		logger.String("resource_type", entry.ResourceType),
		logger.String("resource_id", entry.ResourceID),
		logger.String("ip_address", entry.IPAddress),
		logger.String("user_agent", entry.UserAgent),
		logger.Bool("success", entry.Success),
		logger.String("failure_reason", entry.FailureReason),
		logger.String("metadata", string(metadata)),
		logger.Time("timestamp", entry.Timestamp),
	)

	return nil
}
