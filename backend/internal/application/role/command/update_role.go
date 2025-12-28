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

type UpdateRoleCommand struct {
	RoleID      uuid.UUID
	DisplayName string
	Description string
}

type UpdateRoleHandler struct {
	roleRepository role.Repository
	eventBus       shared.EventBus
	logger         logger.Logger
}

func NewUpdateRoleHandler(
	roleRepository role.Repository,
	eventBus shared.EventBus,
	logger logger.Logger,
) *UpdateRoleHandler {
	return &UpdateRoleHandler{
		roleRepository: roleRepository,
		eventBus:       eventBus,
		logger:         logger,
	}
}

func (handler *UpdateRoleHandler) Handle(context context.Context, command UpdateRoleCommand) (*roledto.RoleDTO, error) {
	existingRole, err := handler.roleRepository.FindByID(context, command.RoleID)
	if err != nil {
		return nil, err
	}

	if !existingRole.CanBeModified() {
		return nil, role.ErrSystemRoleCannotBeModified
	}

	if err := existingRole.UpdateDetails(command.DisplayName, command.Description); err != nil {
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

	handler.logger.Info("role updated successfully",
		logger.String("role_id", existingRole.ID().String()),
		logger.String("name", existingRole.Name()),
	)

	return roledto.RoleFromDomain(existingRole), nil
}
