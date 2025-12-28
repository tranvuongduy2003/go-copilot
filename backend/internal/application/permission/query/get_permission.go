package permissionquery

import (
	"context"

	"github.com/google/uuid"

	permissiondto "github.com/tranvuongduy2003/go-copilot/internal/application/permission/dto"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/permission"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

type GetPermissionQuery struct {
	PermissionID uuid.UUID
}

type GetPermissionHandler struct {
	permissionRepository permission.Repository
	logger               logger.Logger
}

func NewGetPermissionHandler(
	permissionRepository permission.Repository,
	logger logger.Logger,
) *GetPermissionHandler {
	return &GetPermissionHandler{
		permissionRepository: permissionRepository,
		logger:               logger,
	}
}

func (handler *GetPermissionHandler) Handle(context context.Context, query GetPermissionQuery) (*permissiondto.PermissionDTO, error) {
	foundPermission, err := handler.permissionRepository.FindByID(context, query.PermissionID)
	if err != nil {
		return nil, err
	}

	return permissiondto.PermissionFromDomain(foundPermission), nil
}
