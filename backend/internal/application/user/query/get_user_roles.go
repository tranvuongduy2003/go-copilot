package userquery

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	roledto "github.com/tranvuongduy2003/go-copilot/internal/application/role/dto"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/role"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/user"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

type GetUserRolesQuery struct {
	UserID uuid.UUID
}

type GetUserRolesHandler struct {
	userRepository user.Repository
	roleRepository role.Repository
	logger         logger.Logger
}

func NewGetUserRolesHandler(
	userRepository user.Repository,
	roleRepository role.Repository,
	logger logger.Logger,
) *GetUserRolesHandler {
	return &GetUserRolesHandler{
		userRepository: userRepository,
		roleRepository: roleRepository,
		logger:         logger,
	}
}

func (handler *GetUserRolesHandler) Handle(context context.Context, query GetUserRolesQuery) ([]*roledto.RoleDTO, error) {
	existingUser, err := handler.userRepository.FindByID(context, query.UserID)
	if err != nil {
		return nil, err
	}

	roleIDs := existingUser.RoleIDs()
	if len(roleIDs) == 0 {
		return []*roledto.RoleDTO{}, nil
	}

	roles, err := handler.roleRepository.FindByIDs(context, roleIDs)
	if err != nil {
		return nil, fmt.Errorf("get user roles: %w", err)
	}

	return roledto.RolesFromDomain(roles), nil
}
