package usercommand

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/shared"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/user"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

type DeleteUserCommand struct {
	UserID uuid.UUID
}

type DeleteUserHandler struct {
	userRepository user.Repository
	eventBus       shared.EventBus
	logger         logger.Logger
}

func NewDeleteUserHandler(
	userRepository user.Repository,
	eventBus shared.EventBus,
	logger logger.Logger,
) *DeleteUserHandler {
	return &DeleteUserHandler{
		userRepository: userRepository,
		eventBus:       eventBus,
		logger:         logger,
	}
}

func (handler *DeleteUserHandler) Handle(context context.Context, command DeleteUserCommand) error {
	existingUser, err := handler.userRepository.FindByID(context, command.UserID)
	if err != nil {
		return fmt.Errorf("find user: %w", err)
	}

	if err := existingUser.Delete(); err != nil {
		return fmt.Errorf("delete user: %w", err)
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

	handler.logger.Info("user deleted successfully",
		logger.String("user_id", existingUser.ID().String()),
	)

	return nil
}
