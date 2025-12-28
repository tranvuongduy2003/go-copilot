package audit

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresAuditLogger struct {
	pool *pgxpool.Pool
}

func NewPostgresAuditLogger(pool *pgxpool.Pool) *PostgresAuditLogger {
	return &PostgresAuditLogger{
		pool: pool,
	}
}

func (auditLogger *PostgresAuditLogger) Log(ctx context.Context, entry AuditEntry) error {
	metadata, err := json.Marshal(entry.Metadata)
	if err != nil {
		metadata = []byte("{}")
	}

	var ipAddress interface{}
	if entry.IPAddress != "" {
		ipAddress = entry.IPAddress
	}

	query := `
		INSERT INTO audit_logs (
			id, timestamp, event_type, user_id, action, resource_type,
			resource_id, ip_address, user_agent, success, failure_reason, metadata
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		)
	`

	_, err = auditLogger.pool.Exec(ctx, query,
		entry.ID,
		entry.Timestamp,
		entry.EventType,
		entry.UserID,
		entry.Action,
		entry.ResourceType,
		entry.ResourceID,
		ipAddress,
		entry.UserAgent,
		entry.Success,
		entry.FailureReason,
		metadata,
	)

	if err != nil {
		return fmt.Errorf("failed to insert audit log: %w", err)
	}

	return nil
}
