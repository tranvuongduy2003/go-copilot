package permission

import (
	"github.com/tranvuongduy2003/go-copilot/internal/domain/shared"
)

var (
	ErrPermissionNotFound = shared.NewNotFoundError("Permission", "")

	ErrPermissionCodeExists = shared.NewConflictError("Permission", "code", "")

	ErrSystemPermissionCannotBeDeleted = shared.NewBusinessRuleViolationError(
		"system_permission_cannot_be_deleted",
		"system permissions cannot be deleted",
	)

	ErrSystemPermissionCannotBeModified = shared.NewBusinessRuleViolationError(
		"system_permission_cannot_be_modified",
		"system permissions cannot be modified",
	)

	ErrPermissionInUse = shared.NewBusinessRuleViolationError(
		"permission_in_use",
		"permission is currently assigned to one or more roles",
	)
)

func NewPermissionNotFoundError(identifier string) *shared.NotFoundError {
	return shared.NewNotFoundError("Permission", identifier)
}

func NewPermissionCodeExistsError(code string) *shared.ConflictError {
	return shared.NewConflictError("Permission", "code", code)
}
