package rolecommand

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	roledto "github.com/tranvuongduy2003/go-copilot/internal/application/role/dto"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/permission"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/role"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/shared"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

type AssignPermissionToRoleCommand struct {
	RoleID       uuid.UUID
	PermissionID uuid.UUID
}

type AssignPermissionToRoleHandler struct {
	roleRepository       role.Repository
	permissionRepository permission.Repository
	eventBus             shared.EventBus
	logger               logger.Logger
}

func NewAssignPermissionToRoleHandler(
	roleRepository role.Repository,
	permissionRepository permission.Repository,
	eventBus shared.EventBus,
	logger logger.Logger,
) *AssignPermissionToRoleHandler {
	return &AssignPermissionToRoleHandler{
		roleRepository:       roleRepository,
		permissionRepository: permissionRepository,
		eventBus:             eventBus,
		logger:               logger,
	}
}

func (handler *AssignPermissionToRoleHandler) Handle(context context.Context, command AssignPermissionToRoleCommand) (*roledto.RoleDTO, error) {
	existingRole, err := handler.roleRepository.FindByID(context, command.RoleID)
	if err != nil {
		return nil, err
	}

	if !existingRole.CanBeModified() {
		return nil, role.ErrSystemRoleCannotBeModified
	}

	_, err = handler.permissionRepository.FindByID(context, command.PermissionID)
	if err != nil {
		return nil, err
	}

	if err := existingRole.AddPermission(command.PermissionID); err != nil {
		return nil, err
	}

	if err := handler.roleRepository.Update(context, existingRole); err != nil {
		return nil, fmt.Errorf("update role: %w", err)
	}

	if handler.eventBus != nil {
		if err := handler.eventBus.Publish(context, existingRole.DomainEvents()...); err != nil {
			handler.logger.Error("failed to publish domain events",
				logger.String("role_id", existingRole.ID().String()),
				logger.Err(err),
			)
		}
		existingRole.ClearDomainEvents()
	}

	handler.logger.Info("permission assigned to role",
		logger.String("role_id", existingRole.ID().String()),
		logger.String("permission_id", command.PermissionID.String()),
	)

	return roledto.RoleFromDomain(existingRole), nil
}
