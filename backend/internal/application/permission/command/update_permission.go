package permissioncommand

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	permissiondto "github.com/tranvuongduy2003/go-copilot/internal/application/permission/dto"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/permission"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

type UpdatePermissionCommand struct {
	PermissionID uuid.UUID
	Description  string
}

type UpdatePermissionHandler struct {
	permissionRepository permission.Repository
	logger               logger.Logger
}

func NewUpdatePermissionHandler(
	permissionRepository permission.Repository,
	logger logger.Logger,
) *UpdatePermissionHandler {
	return &UpdatePermissionHandler{
		permissionRepository: permissionRepository,
		logger:               logger,
	}
}

func (handler *UpdatePermissionHandler) Handle(context context.Context, command UpdatePermissionCommand) (*permissiondto.PermissionDTO, error) {
	existingPermission, err := handler.permissionRepository.FindByID(context, command.PermissionID)
	if err != nil {
		return nil, err
	}

	if existingPermission.IsSystem() {
		return nil, permission.ErrSystemPermissionCannotBeModified
	}

	existingPermission.UpdateDescription(command.Description)

	if err := handler.permissionRepository.Update(context, existingPermission); err != nil {
		return nil, fmt.Errorf("update permission: %w", err)
	}

	handler.logger.Info("permission updated successfully",
		logger.String("permission_id", existingPermission.ID().String()),
		logger.String("code", existingPermission.CodeString()),
	)

	return permissiondto.PermissionFromDomain(existingPermission), nil
}
