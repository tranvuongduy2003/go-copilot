package rolecommand

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	roledto "github.com/tranvuongduy2003/go-copilot/internal/application/role/dto"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/role"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/shared"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

type RemovePermissionFromRoleCommand struct {
	RoleID       uuid.UUID
	PermissionID uuid.UUID
}

type RemovePermissionFromRoleHandler struct {
	roleRepository role.Repository
	eventBus       shared.EventBus
	logger         logger.Logger
}

func NewRemovePermissionFromRoleHandler(
	roleRepository role.Repository,
	eventBus shared.EventBus,
	logger logger.Logger,
) *RemovePermissionFromRoleHandler {
	return &RemovePermissionFromRoleHandler{
		roleRepository: roleRepository,
		eventBus:       eventBus,
		logger:         logger,
	}
}

func (handler *RemovePermissionFromRoleHandler) Handle(context context.Context, command RemovePermissionFromRoleCommand) (*roledto.RoleDTO, error) {
	existingRole, err := handler.roleRepository.FindByID(context, command.RoleID)
	if err != nil {
		return nil, err
	}

	if !existingRole.CanBeModified() {
		return nil, role.ErrSystemRoleCannotBeModified
	}

	if err := existingRole.RemovePermission(command.PermissionID); err != nil {
		return nil, err
	}

	if err := handler.roleRepository.Update(context, existingRole); err != nil {
		return nil, fmt.Errorf("update role: %w", err)
	}

	if handler.eventBus != nil {
		if err := handler.eventBus.Publish(context, existingRole.DomainEvents()...); err != nil {
			handler.logger.Error("failed to publish domain events",
				logger.String("role_id", existingRole.ID().String()),
				logger.Err(err),
			)
		}
		existingRole.ClearDomainEvents()
	}

	handler.logger.Info("permission removed from role",
		logger.String("role_id", existingRole.ID().String()),
		logger.String("permission_id", command.PermissionID.String()),
	)

	return roledto.RoleFromDomain(existingRole), nil
}
