package query

import (
	"context"
	"fmt"

	"github.com/tranvuongduy2003/go-copilot/internal/application/dto"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/shared"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/user"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

type ListUsersQuery struct {
	Page      int
	Limit     int
	Status    *string
	Search    *string
	SortBy    *string
	SortOrder *string
	DateFrom  *string
	DateTo    *string
}

type ListUsersHandler struct {
	userRepository user.Repository
	logger         logger.Logger
}

func NewListUsersHandler(
	userRepository user.Repository,
	logger logger.Logger,
) *ListUsersHandler {
	return &ListUsersHandler{
		userRepository: userRepository,
		logger:         logger,
	}
}

func (handler *ListUsersHandler) Handle(context context.Context, query ListUsersQuery) (*dto.PaginatedUsersDTO, error) {
	pagination := shared.NewPagination(query.Page, query.Limit)

	filter := user.Filter{
		Search:    query.Search,
		DateRange: shared.NewDateRange(query.DateFrom, query.DateTo),
	}

	if query.Status != nil {
		status, valid := user.ParseStatus(*query.Status)
		if valid {
			filter.Status = &status
		}
	}

	users, total, err := handler.userRepository.List(context, filter, pagination)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	return dto.NewPaginatedUsersDTO(users, total, pagination), nil
}
