package user

import (
	"context"

	"github.com/google/uuid"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/shared"
)

type Filter struct {
	Status     *Status
	Search     *string
	DateRange  shared.DateRange
}

type Repository interface {
	Create(ctx context.Context, user *User) error

	Update(ctx context.Context, user *User) error

	Delete(ctx context.Context, id uuid.UUID) error

	FindByID(ctx context.Context, id uuid.UUID) (*User, error)

	FindByEmail(ctx context.Context, email string) (*User, error)

	ExistsByEmail(ctx context.Context, email string) (bool, error)

	List(ctx context.Context, filter Filter, pagination shared.Pagination) ([]*User, int64, error)
}
