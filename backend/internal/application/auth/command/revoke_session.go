package authcommand

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/auth"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/shared"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

type RevokeSessionCommand struct {
	UserID    uuid.UUID
	SessionID uuid.UUID
}

type RevokeSessionHandler struct {
	refreshTokenRepository auth.RefreshTokenRepository
	eventBus               shared.EventBus
	logger                 logger.Logger
}

type RevokeSessionHandlerParams struct {
	RefreshTokenRepository auth.RefreshTokenRepository
	EventBus               shared.EventBus
	Logger                 logger.Logger
}

func NewRevokeSessionHandler(params RevokeSessionHandlerParams) *RevokeSessionHandler {
	return &RevokeSessionHandler{
		refreshTokenRepository: params.RefreshTokenRepository,
		eventBus:               params.EventBus,
		logger:                 params.Logger,
	}
}

func (handler *RevokeSessionHandler) Handle(context context.Context, command RevokeSessionCommand) error {
	session, err := handler.refreshTokenRepository.FindByID(context, command.SessionID)
	if err != nil {
		return err
	}

	if session.UserID() != command.UserID {
		return auth.ErrSessionNotFound
	}

	if session.IsRevoked() {
		return nil
	}

	if err := handler.refreshTokenRepository.Revoke(context, command.SessionID); err != nil {
		return fmt.Errorf("revoke session: %w", err)
	}

	if handler.eventBus != nil {
		event := auth.NewSessionRevokedEvent(command.UserID, command.SessionID)
		if err := handler.eventBus.Publish(context, event); err != nil {
			handler.logger.Error("failed to publish session revoked event",
				logger.String("user_id", command.UserID.String()),
				logger.String("session_id", command.SessionID.String()),
				logger.Err(err),
			)
		}
	}

	handler.logger.Info("session revoked successfully",
		logger.String("user_id", command.UserID.String()),
		logger.String("session_id", command.SessionID.String()),
	)

	return nil
}
