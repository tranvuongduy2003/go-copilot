package permissionquery

import (
	"context"
	"fmt"

	permissiondto "github.com/tranvuongduy2003/go-copilot/internal/application/permission/dto"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/permission"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

type ListPermissionsQuery struct {
	Resource *string
}

type ListPermissionsHandler struct {
	permissionRepository permission.Repository
	logger               logger.Logger
}

func NewListPermissionsHandler(
	permissionRepository permission.Repository,
	logger logger.Logger,
) *ListPermissionsHandler {
	return &ListPermissionsHandler{
		permissionRepository: permissionRepository,
		logger:               logger,
	}
}

func (handler *ListPermissionsHandler) Handle(context context.Context, query ListPermissionsQuery) ([]*permissiondto.PermissionDTO, error) {
	var permissions []*permission.Permission
	var err error

	if query.Resource != nil && *query.Resource != "" {
		resource, resourceErr := permission.NewResource(*query.Resource)
		if resourceErr != nil {
			return nil, resourceErr
		}
		permissions, err = handler.permissionRepository.FindByResource(context, resource)
	} else {
		permissions, err = handler.permissionRepository.FindAll(context)
	}

	if err != nil {
		return nil, fmt.Errorf("list permissions: %w", err)
	}

	return permissiondto.PermissionsFromDomain(permissions), nil
}
