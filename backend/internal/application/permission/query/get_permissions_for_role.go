package permissionquery

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	permissiondto "github.com/tranvuongduy2003/go-copilot/internal/application/permission/dto"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/permission"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/role"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

type GetPermissionsForRoleQuery struct {
	RoleID uuid.UUID
}

type GetPermissionsForRoleHandler struct {
	roleRepository       role.Repository
	permissionRepository permission.Repository
	logger               logger.Logger
}

func NewGetPermissionsForRoleHandler(
	roleRepository role.Repository,
	permissionRepository permission.Repository,
	logger logger.Logger,
) *GetPermissionsForRoleHandler {
	return &GetPermissionsForRoleHandler{
		roleRepository:       roleRepository,
		permissionRepository: permissionRepository,
		logger:               logger,
	}
}

func (handler *GetPermissionsForRoleHandler) Handle(context context.Context, query GetPermissionsForRoleQuery) ([]*permissiondto.PermissionDTO, error) {
	foundRole, err := handler.roleRepository.FindByID(context, query.RoleID)
	if err != nil {
		return nil, err
	}

	permissionIDs := foundRole.PermissionIDs()
	if len(permissionIDs) == 0 {
		return []*permissiondto.PermissionDTO{}, nil
	}

	permissions, err := handler.permissionRepository.FindByIDs(context, permissionIDs)
	if err != nil {
		return nil, fmt.Errorf("get permissions: %w", err)
	}

	return permissiondto.PermissionsFromDomain(permissions), nil
}
