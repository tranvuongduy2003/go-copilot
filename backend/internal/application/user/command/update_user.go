package usercommand

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	userdto "github.com/tranvuongduy2003/go-copilot/internal/application/user/dto"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/shared"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/user"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

type UpdateUserCommand struct {
	UserID   uuid.UUID
	FullName *string
}

type UpdateUserHandler struct {
	userRepository user.Repository
	eventBus       shared.EventBus
	logger         logger.Logger
}

func NewUpdateUserHandler(
	userRepository user.Repository,
	eventBus shared.EventBus,
	logger logger.Logger,
) *UpdateUserHandler {
	return &UpdateUserHandler{
		userRepository: userRepository,
		eventBus:       eventBus,
		logger:         logger,
	}
}

func (handler *UpdateUserHandler) Handle(context context.Context, command UpdateUserCommand) (*userdto.UserDTO, error) {
	existingUser, err := handler.userRepository.FindByID(context, command.UserID)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}

	if command.FullName != nil {
		if err := existingUser.UpdateProfile(*command.FullName); err != nil {
			return nil, fmt.Errorf("update profile: %w", err)
		}
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

	handler.logger.Info("user updated successfully",
		logger.String("user_id", existingUser.ID().String()),
	)

	return userdto.UserFromDomain(existingUser), nil
}
