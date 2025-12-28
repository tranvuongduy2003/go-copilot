package audit

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/auth"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/shared"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

type AuditEntry struct {
	ID            uuid.UUID
	Timestamp     time.Time
	EventType     string
	UserID        uuid.UUID
	Action        string
	ResourceType  string
	ResourceID    string
	IPAddress     string
	UserAgent     string
	Success       bool
	FailureReason string
	Metadata      map[string]interface{}
}

type AuditLogger interface {
	Log(ctx context.Context, entry AuditEntry) error
}

type AuthAuditHandler struct {
	auditLogger AuditLogger
	logger      logger.Logger
}

func NewAuthAuditHandler(auditLogger AuditLogger, log logger.Logger) *AuthAuditHandler {
	return &AuthAuditHandler{
		auditLogger: auditLogger,
		logger:      log,
	}
}

func (handler *AuthAuditHandler) HandleEvent(ctx context.Context, event shared.DomainEvent) error {
	var entry AuditEntry

	switch e := event.(type) {
	case auth.UserLoggedInEvent:
		entry = AuditEntry{
			ID:           uuid.New(),
			Timestamp:    e.OccurredAt(),
			EventType:    e.EventType(),
			UserID:       e.AggregateID(),
			Action:       "login",
			ResourceType: "session",
			IPAddress:    e.IPAddress,
			UserAgent:    e.UserAgent,
			Success:      true,
			Metadata: map[string]interface{}{
				"email": e.Email,
			},
		}

	case auth.UserLoggedOutEvent:
		entry = AuditEntry{
			ID:           uuid.New(),
			Timestamp:    e.OccurredAt(),
			EventType:    e.EventType(),
			UserID:       e.AggregateID(),
			Action:       "logout",
			ResourceType: "session",
			Success:      true,
			Metadata: map[string]interface{}{
				"logout_all": e.LogoutAll,
			},
		}

	case auth.UserRegisteredEvent:
		entry = AuditEntry{
			ID:           uuid.New(),
			Timestamp:    e.OccurredAt(),
			EventType:    e.EventType(),
			UserID:       e.AggregateID(),
			Action:       "register",
			ResourceType: "user",
			ResourceID:   e.AggregateID().String(),
			Success:      true,
			Metadata: map[string]interface{}{
				"email":     e.Email,
				"full_name": e.FullName,
			},
		}

	case auth.LoginFailedEvent:
		entry = AuditEntry{
			ID:            uuid.New(),
			Timestamp:     e.OccurredAt(),
			EventType:     e.EventType(),
			UserID:        e.AggregateID(),
			Action:        "login",
			ResourceType:  "session",
			IPAddress:     e.IPAddress,
			Success:       false,
			FailureReason: e.FailureReason,
			Metadata: map[string]interface{}{
				"email":         e.Email,
				"attempt_count": e.AttemptCount,
			},
		}

	case auth.AccountLockedEvent:
		entry = AuditEntry{
			ID:           uuid.New(),
			Timestamp:    e.OccurredAt(),
			EventType:    e.EventType(),
			UserID:       e.AggregateID(),
			Action:       "account_locked",
			ResourceType: "user",
			ResourceID:   e.AggregateID().String(),
			Success:      true,
			Metadata: map[string]interface{}{
				"email":           e.Email,
				"lock_duration":   e.LockDuration,
				"failed_attempts": e.FailedAttempts,
			},
		}

	case auth.PasswordResetRequestedEvent:
		entry = AuditEntry{
			ID:           uuid.New(),
			Timestamp:    e.OccurredAt(),
			EventType:    e.EventType(),
			UserID:       e.AggregateID(),
			Action:       "password_reset_requested",
			ResourceType: "user",
			ResourceID:   e.AggregateID().String(),
			Success:      true,
			Metadata: map[string]interface{}{
				"email": e.Email,
			},
		}

	case auth.PasswordResetEvent:
		entry = AuditEntry{
			ID:           uuid.New(),
			Timestamp:    e.OccurredAt(),
			EventType:    e.EventType(),
			UserID:       e.AggregateID(),
			Action:       "password_reset_completed",
			ResourceType: "user",
			ResourceID:   e.AggregateID().String(),
			Success:      true,
		}

	case auth.RefreshTokenRotatedEvent:
		entry = AuditEntry{
			ID:           uuid.New(),
			Timestamp:    e.OccurredAt(),
			EventType:    e.EventType(),
			UserID:       e.AggregateID(),
			Action:       "token_refresh",
			ResourceType: "refresh_token",
			Success:      true,
			Metadata: map[string]interface{}{
				"old_token_id": e.OldTokenID.String(),
				"new_token_id": e.NewTokenID.String(),
			},
		}

	default:
		return nil
	}

	if handler.auditLogger != nil {
		if err := handler.auditLogger.Log(ctx, entry); err != nil {
			handler.logger.Error("failed to write audit log",
				logger.String("event_type", entry.EventType),
				logger.String("user_id", entry.UserID.String()),
				logger.Err(err),
			)
		}
	}

	handler.logger.Info("audit event",
		logger.String("event_type", entry.EventType),
		logger.String("user_id", entry.UserID.String()),
		logger.String("action", entry.Action),
		logger.Bool("success", entry.Success),
	)

	return nil
}

func (handler *AuthAuditHandler) SubscribedEventTypes() []string {
	return []string{
		auth.EventTypeUserLoggedIn,
		auth.EventTypeUserLoggedOut,
		auth.EventTypeUserRegistered,
		auth.EventTypePasswordResetRequested,
		auth.EventTypePasswordReset,
		auth.EventTypeRefreshTokenRotated,
		auth.EventTypeLoginFailed,
		auth.EventTypeAccountLocked,
	}
}
