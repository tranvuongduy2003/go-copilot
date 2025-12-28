package role

import (
	"github.com/tranvuongduy2003/go-copilot/internal/domain/shared"
)

var (
	ErrRoleNotFound = shared.NewNotFoundError("Role", "")

	ErrRoleNameExists = shared.NewConflictError("Role", "name", "")

	ErrSystemRoleCannotBeDeleted = shared.NewBusinessRuleViolationError(
		"system_role_cannot_be_deleted",
		"system roles cannot be deleted",
	)

	ErrSystemRoleCannotBeModified = shared.NewBusinessRuleViolationError(
		"system_role_cannot_be_modified",
		"system roles cannot be modified",
	)

	ErrDefaultRoleCannotBeDeleted = shared.NewBusinessRuleViolationError(
		"default_role_cannot_be_deleted",
		"the default role cannot be deleted",
	)

	ErrPermissionAlreadyAssigned = shared.NewBusinessRuleViolationError(
		"permission_already_assigned",
		"permission is already assigned to this role",
	)

	ErrPermissionNotAssigned = shared.NewBusinessRuleViolationError(
		"permission_not_assigned",
		"permission is not assigned to this role",
	)

	ErrRoleInUse = shared.NewBusinessRuleViolationError(
		"role_in_use",
		"role is assigned to users and cannot be deleted",
	)

	ErrNoDefaultRole = shared.NewNotFoundError("Role", "default")
)

func NewRoleNotFoundError(identifier string) *shared.NotFoundError {
	return shared.NewNotFoundError("Role", identifier)
}

func NewRoleNameExistsError(name string) *shared.ConflictError {
	return shared.NewConflictError("Role", "name", name)
}
