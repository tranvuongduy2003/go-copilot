package query

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/tranvuongduy2003/go-copilot/internal/application/dto"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/user"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

type GetUserQuery struct {
	UserID uuid.UUID
}

type GetUserHandler struct {
	userRepository user.Repository
	logger         logger.Logger
}

func NewGetUserHandler(
	userRepository user.Repository,
	logger logger.Logger,
) *GetUserHandler {
	return &GetUserHandler{
		userRepository: userRepository,
		logger:         logger,
	}
}

func (handler *GetUserHandler) Handle(context context.Context, query GetUserQuery) (*dto.UserDTO, error) {
	foundUser, err := handler.userRepository.FindByID(context, query.UserID)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}

	return dto.UserFromDomain(foundUser), nil
}

type GetUserByEmailQuery struct {
	Email string
}

type GetUserByEmailHandler struct {
	userRepository user.Repository
	logger         logger.Logger
}

func NewGetUserByEmailHandler(
	userRepository user.Repository,
	logger logger.Logger,
) *GetUserByEmailHandler {
	return &GetUserByEmailHandler{
		userRepository: userRepository,
		logger:         logger,
	}
}

func (handler *GetUserByEmailHandler) Handle(context context.Context, query GetUserByEmailQuery) (*dto.UserDTO, error) {
	foundUser, err := handler.userRepository.FindByEmail(context, query.Email)
	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}

	return dto.UserFromDomain(foundUser), nil
}

type CheckEmailExistsQuery struct {
	Email string
}

type CheckEmailExistsHandler struct {
	userRepository user.Repository
	logger         logger.Logger
}

func NewCheckEmailExistsHandler(
	userRepository user.Repository,
	logger logger.Logger,
) *CheckEmailExistsHandler {
	return &CheckEmailExistsHandler{
		userRepository: userRepository,
		logger:         logger,
	}
}

func (handler *CheckEmailExistsHandler) Handle(context context.Context, query CheckEmailExistsQuery) (bool, error) {
	exists, err := handler.userRepository.ExistsByEmail(context, query.Email)
	if err != nil {
		return false, fmt.Errorf("check email exists: %w", err)
	}

	return exists, nil
}
