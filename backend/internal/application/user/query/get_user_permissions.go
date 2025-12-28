package userquery

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	permissiondto "github.com/tranvuongduy2003/go-copilot/internal/application/permission/dto"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/permission"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/role"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/user"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

type GetUserPermissionsQuery struct {
	UserID uuid.UUID
}

type GetUserPermissionsHandler struct {
	userRepository       user.Repository
	roleRepository       role.Repository
	permissionRepository permission.Repository
	logger               logger.Logger
}

func NewGetUserPermissionsHandler(
	userRepository user.Repository,
	roleRepository role.Repository,
	permissionRepository permission.Repository,
	logger logger.Logger,
) *GetUserPermissionsHandler {
	return &GetUserPermissionsHandler{
		userRepository:       userRepository,
		roleRepository:       roleRepository,
		permissionRepository: permissionRepository,
		logger:               logger,
	}
}

func (handler *GetUserPermissionsHandler) Handle(context context.Context, query GetUserPermissionsQuery) ([]*permissiondto.PermissionDTO, error) {
	existingUser, err := handler.userRepository.FindByID(context, query.UserID)
	if err != nil {
		return nil, err
	}

	roleIDs := existingUser.RoleIDs()
	if len(roleIDs) == 0 {
		return []*permissiondto.PermissionDTO{}, nil
	}

	roles, err := handler.roleRepository.FindByIDs(context, roleIDs)
	if err != nil {
		return nil, fmt.Errorf("get user roles: %w", err)
	}

	permissionIDSet := make(map[uuid.UUID]bool)
	for _, roleEntity := range roles {
		for _, permissionID := range roleEntity.PermissionIDs() {
			permissionIDSet[permissionID] = true
		}
	}

	if len(permissionIDSet) == 0 {
		return []*permissiondto.PermissionDTO{}, nil
	}

	permissionIDs := make([]uuid.UUID, 0, len(permissionIDSet))
	for permissionID := range permissionIDSet {
		permissionIDs = append(permissionIDs, permissionID)
	}

	permissions, err := handler.permissionRepository.FindByIDs(context, permissionIDs)
	if err != nil {
		return nil, fmt.Errorf("get permissions: %w", err)
	}

	return permissiondto.PermissionsFromDomain(permissions), nil
}
