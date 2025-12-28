package authcommand

import (
	"context"
	"fmt"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/auth"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/shared"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/user"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
	"github.com/tranvuongduy2003/go-copilot/pkg/security"
)

type ResetPasswordCommand struct {
	ResetToken  string
	NewPassword string
}

type ResetPasswordHandler struct {
	userRepository          user.Repository
	refreshTokenRepository  auth.RefreshTokenRepository
	tokenGenerator          auth.TokenGenerator
	passwordHasher          security.PasswordHasher
	passwordResetTokenStore PasswordResetTokenStore
	eventBus                shared.EventBus
	logger                  logger.Logger
}

type ResetPasswordHandlerParams struct {
	UserRepository          user.Repository
	RefreshTokenRepository  auth.RefreshTokenRepository
	TokenGenerator          auth.TokenGenerator
	PasswordHasher          security.PasswordHasher
	PasswordResetTokenStore PasswordResetTokenStore
	EventBus                shared.EventBus
	Logger                  logger.Logger
}

func NewResetPasswordHandler(params ResetPasswordHandlerParams) *ResetPasswordHandler {
	return &ResetPasswordHandler{
		userRepository:          params.UserRepository,
		refreshTokenRepository:  params.RefreshTokenRepository,
		tokenGenerator:          params.TokenGenerator,
		passwordHasher:          params.PasswordHasher,
		passwordResetTokenStore: params.PasswordResetTokenStore,
		eventBus:                params.EventBus,
		logger:                  params.Logger,
	}
}

func (handler *ResetPasswordHandler) Handle(ctx context.Context, command ResetPasswordCommand) error {
	if err := shared.ValidatePassword(command.NewPassword); err != nil {
		return err
	}

	tokenHash := handler.tokenGenerator.HashRefreshToken(command.ResetToken)

	if handler.passwordResetTokenStore == nil {
		return auth.ErrInvalidResetToken
	}

	email, err := handler.passwordResetTokenStore.Get(ctx, tokenHash)
	if err != nil || email == "" {
		return auth.ErrInvalidResetToken
	}

	existingUser, err := handler.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return auth.ErrInvalidResetToken
	}

	hashedPassword, err := handler.passwordHasher.Hash(command.NewPassword)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	if err := existingUser.ChangePassword(hashedPassword); err != nil {
		return fmt.Errorf("change password: %w", err)
	}

	if err := handler.userRepository.Update(ctx, existingUser); err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	if err := handler.passwordResetTokenStore.Delete(ctx, tokenHash); err != nil {
		handler.logger.Error("failed to delete reset token",
			logger.Err(err),
		)
	}

	if err := handler.refreshTokenRepository.RevokeAllByUserID(ctx, existingUser.ID()); err != nil {
		handler.logger.Error("failed to revoke refresh tokens after password reset",
			logger.String("user_id", existingUser.ID().String()),
			logger.Err(err),
		)
	}

	if handler.eventBus != nil {
		events := existingUser.DomainEvents()
		events = append(events, auth.NewPasswordResetEvent(existingUser.ID()))

		if err := handler.eventBus.Publish(ctx, events...); err != nil {
			handler.logger.Error("failed to publish password reset event",
				logger.String("user_id", existingUser.ID().String()),
				logger.Err(err),
			)
		}
		existingUser.ClearDomainEvents()
	}

	handler.logger.Info("password reset successfully",
		logger.String("user_id", existingUser.ID().String()),
	)

	return nil
}
