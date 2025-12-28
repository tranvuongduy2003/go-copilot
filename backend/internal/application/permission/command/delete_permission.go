package permissioncommand

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/permission"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/role"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

type DeletePermissionCommand struct {
	PermissionID uuid.UUID
}

type DeletePermissionHandler struct {
	permissionRepository permission.Repository
	roleRepository       role.Repository
	logger               logger.Logger
}

func NewDeletePermissionHandler(
	permissionRepository permission.Repository,
	roleRepository role.Repository,
	logger logger.Logger,
) *DeletePermissionHandler {
	return &DeletePermissionHandler{
		permissionRepository: permissionRepository,
		roleRepository:       roleRepository,
		logger:               logger,
	}
}

func (handler *DeletePermissionHandler) Handle(context context.Context, command DeletePermissionCommand) error {
	existingPermission, err := handler.permissionRepository.FindByID(context, command.PermissionID)
	if err != nil {
		return err
	}

	if !existingPermission.CanBeDeleted() {
		return permission.ErrSystemPermissionCannotBeDeleted
	}

	rolesWithPermission, err := handler.roleRepository.FindByPermission(context, command.PermissionID)
	if err != nil {
		return fmt.Errorf("check permission usage: %w", err)
	}
	if len(rolesWithPermission) > 0 {
		return permission.ErrPermissionInUse
	}

	if err := handler.permissionRepository.Delete(context, command.PermissionID); err != nil {
		return fmt.Errorf("delete permission: %w", err)
	}

	handler.logger.Info("permission deleted successfully",
		logger.String("permission_id", existingPermission.ID().String()),
		logger.String("code", existingPermission.CodeString()),
	)

	return nil
}
