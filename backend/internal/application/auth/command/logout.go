package authcommand

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/auth"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/shared"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

type LogoutCommand struct {
	UserID    uuid.UUID
	TokenID   string
	ExpiresAt int64
	LogoutAll bool
}

type LogoutHandler struct {
	refreshTokenRepository auth.RefreshTokenRepository
	tokenBlacklist         auth.TokenBlacklist
	eventBus               shared.EventBus
	logger                 logger.Logger
}

type LogoutHandlerParams struct {
	RefreshTokenRepository auth.RefreshTokenRepository
	TokenBlacklist         auth.TokenBlacklist
	EventBus               shared.EventBus
	Logger                 logger.Logger
}

func NewLogoutHandler(params LogoutHandlerParams) *LogoutHandler {
	return &LogoutHandler{
		refreshTokenRepository: params.RefreshTokenRepository,
		tokenBlacklist:         params.TokenBlacklist,
		eventBus:               params.EventBus,
		logger:                 params.Logger,
	}
}

func (handler *LogoutHandler) Handle(ctx context.Context, command LogoutCommand) error {
	if err := handler.tokenBlacklist.Add(ctx, command.TokenID, command.ExpiresAt); err != nil {
		return fmt.Errorf("add token to blacklist: %w", err)
	}

	if command.LogoutAll {
		if err := handler.refreshTokenRepository.RevokeAllByUserID(ctx, command.UserID); err != nil {
			return fmt.Errorf("revoke all refresh tokens: %w", err)
		}
	}

	if handler.eventBus != nil {
		event := auth.NewUserLoggedOutEvent(command.UserID, command.LogoutAll)
		if err := handler.eventBus.Publish(ctx, event); err != nil {
			handler.logger.Error("failed to publish logout event",
				logger.String("user_id", command.UserID.String()),
				logger.Err(err),
			)
		}
	}

	handler.logger.Info("user logged out successfully",
		logger.String("user_id", command.UserID.String()),
		logger.Bool("logout_all", command.LogoutAll),
	)

	return nil
}
