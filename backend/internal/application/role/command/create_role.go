package rolecommand

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	roledto "github.com/tranvuongduy2003/go-copilot/internal/application/role/dto"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/permission"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/role"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/shared"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

type CreateRoleCommand struct {
	Name          string
	DisplayName   string
	Description   string
	PermissionIDs []uuid.UUID
}

type CreateRoleHandler struct {
	roleRepository       role.Repository
	permissionRepository permission.Repository
	eventBus             shared.EventBus
	logger               logger.Logger
}

func NewCreateRoleHandler(
	roleRepository role.Repository,
	permissionRepository permission.Repository,
	eventBus shared.EventBus,
	logger logger.Logger,
) *CreateRoleHandler {
	return &CreateRoleHandler{
		roleRepository:       roleRepository,
		permissionRepository: permissionRepository,
		eventBus:             eventBus,
		logger:               logger,
	}
}

func (handler *CreateRoleHandler) Handle(context context.Context, command CreateRoleCommand) (*roledto.RoleDTO, error) {
	exists, err := handler.roleRepository.ExistsByName(context, command.Name)
	if err != nil {
		return nil, fmt.Errorf("check role name exists: %w", err)
	}
	if exists {
		return nil, role.ErrRoleNameExists
	}

	if len(command.PermissionIDs) > 0 {
		permissions, err := handler.permissionRepository.FindByIDs(context, command.PermissionIDs)
		if err != nil {
			return nil, fmt.Errorf("validate permissions: %w", err)
		}
		if len(permissions) != len(command.PermissionIDs) {
			return nil, permission.ErrPermissionNotFound
		}
	}

	newRole, err := role.NewRole(role.NewRoleParams{
		Name:          command.Name,
		DisplayName:   command.DisplayName,
		Description:   command.Description,
		PermissionIDs: command.PermissionIDs,
		IsSystem:      false,
		IsDefault:     false,
		Priority:      0,
	})
	if err != nil {
		return nil, fmt.Errorf("create role: %w", err)
	}

	if err := handler.roleRepository.Create(context, newRole); err != nil {
		return nil, fmt.Errorf("save role: %w", err)
	}

	if handler.eventBus != nil {
		if err := handler.eventBus.Publish(context, newRole.DomainEvents()...); err != nil {
			handler.logger.Error("failed to publish domain events",
				logger.String("role_id", newRole.ID().String()),
				logger.Err(err),
			)
		}
		newRole.ClearDomainEvents()
	}

	handler.logger.Info("role created successfully",
		logger.String("role_id", newRole.ID().String()),
		logger.String("name", newRole.Name()),
	)

	return roledto.RoleFromDomain(newRole), nil
}
