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

type RevokeRoleFromUserCommand struct {
	UserID uuid.UUID
	RoleID uuid.UUID
}

type RevokeRoleFromUserHandler struct {
	userRepository user.Repository
	eventBus       shared.EventBus
	logger         logger.Logger
}

func NewRevokeRoleFromUserHandler(
	userRepository user.Repository,
	eventBus shared.EventBus,
	logger logger.Logger,
) *RevokeRoleFromUserHandler {
	return &RevokeRoleFromUserHandler{
		userRepository: userRepository,
		eventBus:       eventBus,
		logger:         logger,
	}
}

func (handler *RevokeRoleFromUserHandler) Handle(context context.Context, command RevokeRoleFromUserCommand) (*userdto.UserDTO, error) {
	existingUser, err := handler.userRepository.FindByID(context, command.UserID)
	if err != nil {
		return nil, err
	}

	if err := existingUser.RevokeRole(command.RoleID); err != nil {
		return nil, err
	}

	if err := handler.userRepository.Update(context, existingUser); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
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

	handler.logger.Info("role revoked from user",
		logger.String("user_id", existingUser.ID().String()),
		logger.String("role_id", command.RoleID.String()),
	)

	return userdto.UserFromDomain(existingUser), nil
}
