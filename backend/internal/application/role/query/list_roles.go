package rolequery

import (
	"context"
	"fmt"

	roledto "github.com/tranvuongduy2003/go-copilot/internal/application/role/dto"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/role"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

type ListRolesQuery struct{}

type ListRolesHandler struct {
	roleRepository role.Repository
	logger         logger.Logger
}

func NewListRolesHandler(
	roleRepository role.Repository,
	logger logger.Logger,
) *ListRolesHandler {
	return &ListRolesHandler{
		roleRepository: roleRepository,
		logger:         logger,
	}
}

func (handler *ListRolesHandler) Handle(context context.Context, query ListRolesQuery) ([]*roledto.RoleDTO, error) {
	roles, err := handler.roleRepository.FindAll(context)
	if err != nil {
		return nil, fmt.Errorf("list roles: %w", err)
	}

	return roledto.RolesFromDomain(roles), nil
}
