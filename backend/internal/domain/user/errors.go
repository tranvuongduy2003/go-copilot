package user

import (
	"github.com/tranvuongduy2003/go-copilot/internal/domain/shared"
)

var (
	ErrUserNotFound = shared.NewNotFoundError("User", "")
	ErrEmailAlreadyExists = shared.NewConflictError("User", "email", "")
	ErrInvalidEmail = shared.NewValidationError("email", "invalid email format")
	ErrInvalidPassword = shared.NewValidationError("password", "invalid password")
	ErrInvalidStatus = shared.NewValidationError("status", "invalid status")
	ErrUserAlreadyActive = shared.NewBusinessRuleViolationError("user_already_active", "user is already active")
	ErrUserAlreadyInactive = shared.NewBusinessRuleViolationError("user_already_inactive", "user is already inactive")
	ErrUserIsBanned = shared.NewBusinessRuleViolationError("user_is_banned", "user is banned and cannot perform this action")
)

func NewUserNotFoundError(identifier string) *shared.NotFoundError {
	return shared.NewNotFoundError("User", identifier)
}

func NewEmailAlreadyExistsError(email string) *shared.ConflictError {
	return shared.NewConflictError("User", "email", email)
}

func NewInvalidStatusTransitionError(current, target Status) *shared.InvalidStatusTransitionError {
	return shared.NewInvalidStatusTransitionError(current.String(), target.String())
}
