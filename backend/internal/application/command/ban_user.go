package command

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/tranvuongduy2003/go-copilot/internal/application/dto"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/shared"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/user"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

type BanUserCommand struct {
	UserID uuid.UUID
	Reason string
}

type BanUserHandler struct {
	userRepository user.Repository
	eventBus       shared.EventBus
	logger         logger.Logger
}

func NewBanUserHandler(
	userRepository user.Repository,
	eventBus shared.EventBus,
	logger logger.Logger,
) *BanUserHandler {
	return &BanUserHandler{
		userRepository: userRepository,
		eventBus:       eventBus,
		logger:         logger,
	}
}

func (handler *BanUserHandler) Handle(context context.Context, command BanUserCommand) (*dto.UserDTO, error) {
	existingUser, err := handler.userRepository.FindByID(context, command.UserID)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}

	if err := existingUser.Ban(command.Reason); err != nil {
		return nil, fmt.Errorf("ban user: %w", err)
	}

	if err := handler.userRepository.Update(context, existingUser); err != nil {
		return nil, fmt.Errorf("save user: %w", err)
	}

	if handler.eventBus != nil {
		if err := handler.eventBus.Publish(context, existingUser.DomainEvents()...); err != nil {
			handler.logger.Error("failed to publish domain events",
				logger.String("user_id", existingUser.ID().String()),
				logger.Err(err),
			)
		}
		existingUser.ClearDomainEvents()
	}

	handler.logger.Info("user banned successfully",
		logger.String("user_id", existingUser.ID().String()),
		logger.String("reason", command.Reason),
	)

	return dto.UserFromDomain(existingUser), nil
}
