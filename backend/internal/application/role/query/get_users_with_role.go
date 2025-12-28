package rolequery

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	userdto "github.com/tranvuongduy2003/go-copilot/internal/application/user/dto"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/role"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/user"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

type GetUsersWithRoleQuery struct {
	RoleID uuid.UUID
}

type GetUsersWithRoleHandler struct {
	roleRepository role.Repository
	userRepository user.Repository
	logger         logger.Logger
}

func NewGetUsersWithRoleHandler(
	roleRepository role.Repository,
	userRepository user.Repository,
	logger logger.Logger,
) *GetUsersWithRoleHandler {
	return &GetUsersWithRoleHandler{
		roleRepository: roleRepository,
		userRepository: userRepository,
		logger:         logger,
	}
}

func (handler *GetUsersWithRoleHandler) Handle(context context.Context, query GetUsersWithRoleQuery) ([]*userdto.UserDTO, error) {
	_, err := handler.roleRepository.FindByID(context, query.RoleID)
	if err != nil {
		return nil, err
	}

	users, err := handler.userRepository.FindByRole(context, query.RoleID)
	if err != nil {
		return nil, fmt.Errorf("get users with role: %w", err)
	}

	return userdto.UsersFromDomain(users), nil
}
