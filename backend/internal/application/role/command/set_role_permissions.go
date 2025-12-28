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

type SetRolePermissionsCommand struct {
	RoleID        uuid.UUID
	PermissionIDs []uuid.UUID
}

type SetRolePermissionsHandler struct {
	roleRepository       role.Repository
	permissionRepository permission.Repository
	eventBus             shared.EventBus
	logger               logger.Logger
}

func NewSetRolePermissionsHandler(
	roleRepository role.Repository,
	permissionRepository permission.Repository,
	eventBus shared.EventBus,
	logger logger.Logger,
) *SetRolePermissionsHandler {
	return &SetRolePermissionsHandler{
		roleRepository:       roleRepository,
		permissionRepository: permissionRepository,
		eventBus:             eventBus,
		logger:               logger,
	}
}

func (handler *SetRolePermissionsHandler) Handle(context context.Context, command SetRolePermissionsCommand) (*roledto.RoleDTO, error) {
	existingRole, err := handler.roleRepository.FindByID(context, command.RoleID)
	if err != nil {
		return nil, err
	}

	if !existingRole.CanBeModified() {
		return nil, role.ErrSystemRoleCannotBeModified
	}

	if len(command.PermissionIDs) > 0 {
		permissions, err := handler.permissionRepository.FindByIDs(context, command.PermissionIDs)
		if err != nil {
			return nil, fmt.Errorf("validate permissions: %w", err)
		}
		if len(permissions) != len(command.PermissionIDs) {
			return nil, permission.ErrPermissionNotFound
		}
	}

	existingRole.SetPermissions(command.PermissionIDs)

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

	handler.logger.Info("role permissions updated",
		logger.String("role_id", existingRole.ID().String()),
		logger.Int("permission_count", len(command.PermissionIDs)),
	)

	return roledto.RoleFromDomain(existingRole), nil
}
