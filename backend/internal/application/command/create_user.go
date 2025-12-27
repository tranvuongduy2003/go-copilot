package command

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/tranvuongduy2003/go-copilot/internal/application/dto"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/shared"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/user"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
	"github.com/tranvuongduy2003/go-copilot/pkg/security"
)

type CreateUserCommand struct {
	Email    string
	Password string
	FullName string
}

type CreateUserHandler struct {
	userRepository user.Repository
	passwordHasher security.PasswordHasher
	eventBus       shared.EventBus
	logger         logger.Logger
}

func NewCreateUserHandler(
	userRepository user.Repository,
	passwordHasher security.PasswordHasher,
	eventBus shared.EventBus,
	logger logger.Logger,
) *CreateUserHandler {
	return &CreateUserHandler{
		userRepository: userRepository,
		passwordHasher: passwordHasher,
		eventBus:       eventBus,
		logger:         logger,
	}
}

func (handler *CreateUserHandler) Handle(context context.Context, command CreateUserCommand) (*dto.UserDTO, error) {
	if err := shared.ValidatePassword(command.Password); err != nil {
		return nil, err
	}

	exists, err := handler.userRepository.ExistsByEmail(context, command.Email)
	if err != nil {
		return nil, fmt.Errorf("check email exists: %w", err)
	}
	if exists {
		return nil, user.NewEmailAlreadyExistsError(command.Email)
	}

	hashedPassword, err := handler.passwordHasher.Hash(command.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	newUser, err := user.NewUser(user.NewUserParams{
		Email:        command.Email,
		PasswordHash: hashedPassword,
		FullName:     command.FullName,
	})
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	if err := handler.userRepository.Create(context, newUser); err != nil {
		return nil, fmt.Errorf("save user: %w", err)
	}

	if handler.eventBus != nil {
		if err := handler.eventBus.Publish(context, newUser.DomainEvents()...); err != nil {
			handler.logger.Error("failed to publish domain events",
				logger.String("user_id", newUser.ID().String()),
				logger.Err(err),
			)
		}
		newUser.ClearDomainEvents()
	}

	handler.logger.Info("user created successfully",
		logger.String("user_id", newUser.ID().String()),
		logger.String("email", newUser.Email().String()),
	)

	return dto.UserFromDomain(newUser), nil
}

type CreateUserResult struct {
	UserID uuid.UUID
	User   *dto.UserDTO
}
