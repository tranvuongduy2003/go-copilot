package usercommand

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	userdto "github.com/tranvuongduy2003/go-copilot/internal/application/user/dto"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/role"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/shared"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/user"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

type SetUserRolesCommand struct {
	UserID  uuid.UUID
	RoleIDs []uuid.UUID
}

type SetUserRolesHandler struct {
	userRepository user.Repository
	roleRepository role.Repository
	eventBus       shared.EventBus
	logger         logger.Logger
}

func NewSetUserRolesHandler(
	userRepository user.Repository,
	roleRepository role.Repository,
	eventBus shared.EventBus,
	logger logger.Logger,
) *SetUserRolesHandler {
	return &SetUserRolesHandler{
		userRepository: userRepository,
		roleRepository: roleRepository,
		eventBus:       eventBus,
		logger:         logger,
	}
}

func (handler *SetUserRolesHandler) Handle(context context.Context, command SetUserRolesCommand) (*userdto.UserDTO, error) {
	existingUser, err := handler.userRepository.FindByID(context, command.UserID)
	if err != nil {
		return nil, err
	}

	if len(command.RoleIDs) > 0 {
		roles, err := handler.roleRepository.FindByIDs(context, command.RoleIDs)
		if err != nil {
			return nil, fmt.Errorf("validate roles: %w", err)
		}
		if len(roles) != len(command.RoleIDs) {
			return nil, role.ErrRoleNotFound
		}
	}

	existingUser.SetRoles(command.RoleIDs)

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

	handler.logger.Info("user roles updated",
		logger.String("user_id", existingUser.ID().String()),
		logger.Int("role_count", len(command.RoleIDs)),
	)

	return userdto.UserFromDomain(existingUser), nil
}
