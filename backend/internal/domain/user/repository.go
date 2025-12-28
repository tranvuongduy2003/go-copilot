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
	// Create persists a new user to the data store.
	// Returns ErrEmailAlreadyExists if the email is already registered.
	// Returns wrapped database errors for other failures.
	Create(ctx context.Context, user *User) error

	// Update modifies an existing user in the data store.
	// Returns ErrUserNotFound if the user does not exist.
	// Returns wrapped database errors for other failures.
	Update(ctx context.Context, user *User) error

	// Delete performs a soft delete on the user by setting deleted_at.
	// Returns ErrUserNotFound if the user does not exist.
	// Returns wrapped database errors for other failures.
	Delete(ctx context.Context, id uuid.UUID) error

	// FindByID retrieves a user by their unique identifier.
	// Returns ErrUserNotFound if the user does not exist or is soft-deleted.
	// Never returns (nil, nil) - always returns an error for not found.
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)

	// FindByEmail retrieves a user by their email address (case-insensitive).
	// Returns ErrUserNotFound if the user does not exist or is soft-deleted.
	// Never returns (nil, nil) - always returns an error for not found.
	FindByEmail(ctx context.Context, email string) (*User, error)

	// ExistsByEmail checks if a user with the given email exists.
	// Returns (true, nil) if exists, (false, nil) if not exists.
	// Only returns error for database failures.
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// FindByRole retrieves all users assigned to a specific role.
	// Returns empty slice (not nil) if no users have the role.
	// Returns wrapped database errors for failures.
	FindByRole(ctx context.Context, roleID uuid.UUID) ([]*User, error)

	// List retrieves users matching the filter with pagination.
	// Returns (users, totalCount, nil) on success.
	// Returns empty slice (not nil) when no results match.
	// Returns wrapped database errors for failures.
	List(ctx context.Context, filter Filter, pagination shared.Pagination) ([]*User, int64, error)
}
