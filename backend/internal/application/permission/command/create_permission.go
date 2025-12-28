package permissioncommand

import (
	"context"
	"fmt"

	permissiondto "github.com/tranvuongduy2003/go-copilot/internal/application/permission/dto"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/permission"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

type CreatePermissionCommand struct {
	Resource    string
	Action      string
	Description string
}

type CreatePermissionHandler struct {
	permissionRepository permission.Repository
	logger               logger.Logger
}

func NewCreatePermissionHandler(
	permissionRepository permission.Repository,
	logger logger.Logger,
) *CreatePermissionHandler {
	return &CreatePermissionHandler{
		permissionRepository: permissionRepository,
		logger:               logger,
	}
}

func (handler *CreatePermissionHandler) Handle(context context.Context, command CreatePermissionCommand) (*permissiondto.PermissionDTO, error) {
	newPermission, err := permission.NewPermission(permission.NewPermissionParams{
		Resource:    command.Resource,
		Action:      command.Action,
		Description: command.Description,
		IsSystem:    false,
	})
	if err != nil {
		return nil, fmt.Errorf("create permission: %w", err)
	}

	exists, err := handler.permissionRepository.ExistsByCode(context, newPermission.Code())
	if err != nil {
		return nil, fmt.Errorf("check permission exists: %w", err)
	}
	if exists {
		return nil, permission.ErrPermissionCodeExists
	}

	if err := handler.permissionRepository.Create(context, newPermission); err != nil {
		return nil, fmt.Errorf("save permission: %w", err)
	}

	handler.logger.Info("permission created successfully",
		logger.String("permission_id", newPermission.ID().String()),
		logger.String("code", newPermission.CodeString()),
	)

	return permissiondto.PermissionFromDomain(newPermission), nil
}
