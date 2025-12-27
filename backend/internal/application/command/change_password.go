package command

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/shared"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/user"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
	"github.com/tranvuongduy2003/go-copilot/pkg/security"
)

type ChangePasswordCommand struct {
	UserID          uuid.UUID
	CurrentPassword string
	NewPassword     string
}

type ChangePasswordHandler struct {
	userRepository user.Repository
	passwordHasher security.PasswordHasher
	eventBus       shared.EventBus
	logger         logger.Logger
}

func NewChangePasswordHandler(
	userRepository user.Repository,
	passwordHasher security.PasswordHasher,
	eventBus shared.EventBus,
	logger logger.Logger,
) *ChangePasswordHandler {
	return &ChangePasswordHandler{
		userRepository: userRepository,
		passwordHasher: passwordHasher,
		eventBus:       eventBus,
		logger:         logger,
	}
}

func (handler *ChangePasswordHandler) Handle(context context.Context, command ChangePasswordCommand) error {
	existingUser, err := handler.userRepository.FindByID(context, command.UserID)
	if err != nil {
		return fmt.Errorf("find user: %w", err)
	}

	if err := handler.passwordHasher.Compare(existingUser.PasswordHash().String(), command.CurrentPassword); err != nil {
		return user.ErrInvalidPassword
	}

	if err := shared.ValidatePassword(command.NewPassword); err != nil {
		return err
	}

	newHashedPassword, err := handler.passwordHasher.Hash(command.NewPassword)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	if err := existingUser.ChangePassword(newHashedPassword); err != nil {
		return fmt.Errorf("change password: %w", err)
	}

	if err := handler.userRepository.Update(context, existingUser); err != nil {
		return fmt.Errorf("save user: %w", err)
	}

	if handler.eventBus != nil {
		if err := handler.eventBus.Publish(context, existingUser.DomainEvents()...); err != nil {
			handler.logger.Error("failed to publish domain events",
				logger.String("user_id", existingUser.ID().String()),
				logger.Err(err),
			)
		}
		existingUser.ClearDomainEvents()
	}

	handler.logger.Info("password changed successfully",
		logger.String("user_id", existingUser.ID().String()),
	)

	return nil
}
