package command

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/tranvuongduy2003/go-copilot/internal/application/dto"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/shared"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/user"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

type DeactivateUserCommand struct {
	UserID uuid.UUID
}

type DeactivateUserHandler struct {
	userRepository user.Repository
	eventBus       shared.EventBus
	logger         logger.Logger
}

func NewDeactivateUserHandler(
	userRepository user.Repository,
	eventBus shared.EventBus,
	logger logger.Logger,
) *DeactivateUserHandler {
	return &DeactivateUserHandler{
		userRepository: userRepository,
		eventBus:       eventBus,
		logger:         logger,
	}
}

func (handler *DeactivateUserHandler) Handle(context context.Context, command DeactivateUserCommand) (*dto.UserDTO, error) {
	existingUser, err := handler.userRepository.FindByID(context, command.UserID)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}

	if err := existingUser.Deactivate(); err != nil {
		return nil, fmt.Errorf("deactivate user: %w", err)
	}

	if err := handler.userRepository.Update(context, existingUser); err != nil {
		return nil, fmt.Errorf("save user: %w", err)
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

	handler.logger.Info("user deactivated successfully",
		logger.String("user_id", existingUser.ID().String()),
	)

	return dto.UserFromDomain(existingUser), nil
}
