package authquery

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	authdto "github.com/tranvuongduy2003/go-copilot/internal/application/auth/dto"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/permission"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/role"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/user"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

type GetCurrentUserQuery struct {
	UserID uuid.UUID
}

type GetCurrentUserHandler struct {
	userRepository       user.Repository
	roleRepository       role.Repository
	permissionRepository permission.Repository
	logger               logger.Logger
}

type GetCurrentUserHandlerParams struct {
	UserRepository       user.Repository
	RoleRepository       role.Repository
	PermissionRepository permission.Repository
	Logger               logger.Logger
}

func NewGetCurrentUserHandler(params GetCurrentUserHandlerParams) *GetCurrentUserHandler {
	return &GetCurrentUserHandler{
		userRepository:       params.UserRepository,
		roleRepository:       params.RoleRepository,
		permissionRepository: params.PermissionRepository,
		logger:               params.Logger,
	}
}

func (handler *GetCurrentUserHandler) Handle(ctx context.Context, query GetCurrentUserQuery) (*authdto.AuthUserDTO, error) {
	existingUser, err := handler.userRepository.FindByID(ctx, query.UserID)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}

	roleNames, permissions := handler.loadUserRolesAndPermissions(ctx, existingUser)

	return &authdto.AuthUserDTO{
		ID:          existingUser.ID(),
		Email:       existingUser.Email().String(),
		FullName:    existingUser.FullName().String(),
		Status:      existingUser.Status().String(),
		Roles:       roleNames,
		Permissions: permissions,
	}, nil
}

func (handler *GetCurrentUserHandler) loadUserRolesAndPermissions(ctx context.Context, domainUser *user.User) ([]string, []string) {
	roleIDs := domainUser.RoleIDs()
	if len(roleIDs) == 0 {
		return []string{}, []string{}
	}

	roles, err := handler.roleRepository.FindByIDs(ctx, roleIDs)
	if err != nil {
		handler.logger.Error("failed to load user roles",
			logger.String("user_id", domainUser.ID().String()),
			logger.Err(err),
		)
		return []string{}, []string{}
	}

	roleNames := make([]string, 0, len(roles))
	permissionIDSet := make(map[uuid.UUID]bool)

	for _, roleEntity := range roles {
		roleNames = append(roleNames, roleEntity.Name())
		for _, permissionID := range roleEntity.PermissionIDs() {
			permissionIDSet[permissionID] = true
		}
	}

	if len(permissionIDSet) == 0 {
		return roleNames, []string{}
	}

	permissionIDs := make([]uuid.UUID, 0, len(permissionIDSet))
	for permissionID := range permissionIDSet {
		permissionIDs = append(permissionIDs, permissionID)
	}

	permissions, err := handler.permissionRepository.FindByIDs(ctx, permissionIDs)
	if err != nil {
		handler.logger.Error("failed to load permissions",
			logger.String("user_id", domainUser.ID().String()),
			logger.Err(err),
		)
		return roleNames, []string{}
	}

	permissionCodes := make([]string, 0, len(permissions))
	for _, permissionEntity := range permissions {
		permissionCodes = append(permissionCodes, permissionEntity.Code().String())
	}

	return roleNames, permissionCodes
}
