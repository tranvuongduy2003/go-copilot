package rolecommand

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/role"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/shared"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/user"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

type DeleteRoleCommand struct {
	RoleID uuid.UUID
}

type DeleteRoleHandler struct {
	roleRepository role.Repository
	userRepository user.Repository
	eventBus       shared.EventBus
	logger         logger.Logger
}

func NewDeleteRoleHandler(
	roleRepository role.Repository,
	userRepository user.Repository,
	eventBus shared.EventBus,
	logger logger.Logger,
) *DeleteRoleHandler {
	return &DeleteRoleHandler{
		roleRepository: roleRepository,
		userRepository: userRepository,
		eventBus:       eventBus,
		logger:         logger,
	}
}

func (handler *DeleteRoleHandler) Handle(context context.Context, command DeleteRoleCommand) error {
	existingRole, err := handler.roleRepository.FindByID(context, command.RoleID)
	if err != nil {
		return err
	}

	if !existingRole.CanBeDeleted() {
		if existingRole.IsSystem() {
			return role.ErrSystemRoleCannotBeDeleted
		}
		if existingRole.IsDefault() {
			return role.ErrDefaultRoleCannotBeDeleted
		}
	}

	usersWithRole, err := handler.userRepository.FindByRole(context, command.RoleID)
	if err != nil {
		return fmt.Errorf("check role usage: %w", err)
	}
	if len(usersWithRole) > 0 {
		return role.ErrRoleInUse
	}

	if err := handler.roleRepository.Delete(context, command.RoleID); err != nil {
		return fmt.Errorf("delete role: %w", err)
	}

	if handler.eventBus != nil {
		event := role.NewRoleDeletedEvent(existingRole.ID(), existingRole.Name())
		if err := handler.eventBus.Publish(context, event); err != nil {
			handler.logger.Error("failed to publish role deleted event",
				logger.String("role_id", existingRole.ID().String()),
				logger.Err(err),
			)
		}
	}

	handler.logger.Info("role deleted successfully",
		logger.String("role_id", existingRole.ID().String()),
		logger.String("name", existingRole.Name()),
	)

	return nil
}
