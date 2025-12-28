package authcommand

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/auth"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/shared"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/user"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

type ForgotPasswordCommand struct {
	Email string
}

type ForgotPasswordResult struct {
	ResetToken string
	ExpiresAt  time.Time
}

type PasswordResetTokenStore interface {
	Store(ctx context.Context, email string, tokenHash string, expiresAt time.Time) error
	Get(ctx context.Context, tokenHash string) (string, error)
	Delete(ctx context.Context, tokenHash string) error
}

type ForgotPasswordHandler struct {
	userRepository          user.Repository
	tokenGenerator          auth.TokenGenerator
	passwordResetTokenStore PasswordResetTokenStore
	eventBus                shared.EventBus
	resetTokenTTL           time.Duration
	logger                  logger.Logger
}

type ForgotPasswordHandlerParams struct {
	UserRepository          user.Repository
	TokenGenerator          auth.TokenGenerator
	PasswordResetTokenStore PasswordResetTokenStore
	EventBus                shared.EventBus
	ResetTokenTTL           time.Duration
	Logger                  logger.Logger
}

func NewForgotPasswordHandler(params ForgotPasswordHandlerParams) *ForgotPasswordHandler {
	return &ForgotPasswordHandler{
		userRepository:          params.UserRepository,
		tokenGenerator:          params.TokenGenerator,
		passwordResetTokenStore: params.PasswordResetTokenStore,
		eventBus:                params.EventBus,
		resetTokenTTL:           params.ResetTokenTTL,
		logger:                  params.Logger,
	}
}

func (handler *ForgotPasswordHandler) Handle(ctx context.Context, command ForgotPasswordCommand) (*ForgotPasswordResult, error) {
	existingUser, err := handler.userRepository.FindByEmail(ctx, command.Email)
	if err != nil {
		handler.logger.Info("password reset requested for non-existent email",
			logger.String("email", command.Email),
		)
		return nil, nil
	}

	if !existingUser.Status().IsActive() {
		handler.logger.Info("password reset requested for inactive account",
			logger.String("email", command.Email),
		)
		return nil, nil
	}

	resetToken, err := handler.generateResetToken()
	if err != nil {
		return nil, fmt.Errorf("generate reset token: %w", err)
	}

	tokenHash := handler.tokenGenerator.HashRefreshToken(resetToken)
	expiresAt := time.Now().UTC().Add(handler.resetTokenTTL)

	if handler.passwordResetTokenStore != nil {
		if err := handler.passwordResetTokenStore.Store(ctx, command.Email, tokenHash, expiresAt); err != nil {
			return nil, fmt.Errorf("store reset token: %w", err)
		}
	}

	if handler.eventBus != nil {
		event := auth.NewPasswordResetRequestedEvent(existingUser.ID(), command.Email)
		if err := handler.eventBus.Publish(ctx, event); err != nil {
			handler.logger.Error("failed to publish password reset requested event",
				logger.String("user_id", existingUser.ID().String()),
				logger.Err(err),
			)
		}
	}

	handler.logger.Info("password reset token generated",
		logger.String("email", command.Email),
	)

	return &ForgotPasswordResult{
		ResetToken: resetToken,
		ExpiresAt:  expiresAt,
	}, nil
}

func (handler *ForgotPasswordHandler) generateResetToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate random bytes: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
