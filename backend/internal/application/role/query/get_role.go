package rolequery

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	roledto "github.com/tranvuongduy2003/go-copilot/internal/application/role/dto"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/permission"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/role"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

type GetRoleQuery struct {
	RoleID uuid.UUID
}

type GetRoleHandler struct {
	roleRepository       role.Repository
	permissionRepository permission.Repository
	logger               logger.Logger
}

func NewGetRoleHandler(
	roleRepository role.Repository,
	permissionRepository permission.Repository,
	logger logger.Logger,
) *GetRoleHandler {
	return &GetRoleHandler{
		roleRepository:       roleRepository,
		permissionRepository: permissionRepository,
		logger:               logger,
	}
}

func (handler *GetRoleHandler) Handle(context context.Context, query GetRoleQuery) (*roledto.RoleWithPermissionsDTO, error) {
	foundRole, err := handler.roleRepository.FindByID(context, query.RoleID)
	if err != nil {
		return nil, err
	}

	permissionIDs := foundRole.PermissionIDs()
	var permissionCodes []string

	if len(permissionIDs) > 0 {
		permissions, err := handler.permissionRepository.FindByIDs(context, permissionIDs)
		if err != nil {
			return nil, fmt.Errorf("get permissions: %w", err)
		}
		permissionCodes = make([]string, len(permissions))
		for i, permission := range permissions {
			permissionCodes[i] = permission.CodeString()
		}
	} else {
		permissionCodes = []string{}
	}

	return roledto.RoleWithPermissionsFromDomain(foundRole, permissionCodes), nil
}
