package role

import (
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/shared"
)

var roleNameRegex = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

type Role struct {
	shared.AggregateRoot
	name          string
	displayName   string
	description   string
	permissionIDs []uuid.UUID
	isSystem      bool
	isDefault     bool
	priority      int
	createdAt     time.Time
	updatedAt     time.Time
}

type NewRoleParams struct {
	Name          string
	DisplayName   string
	Description   string
	PermissionIDs []uuid.UUID
	IsSystem      bool
	IsDefault     bool
	Priority      int
}

func NewRole(params NewRoleParams) (*Role, error) {
	name := strings.ToLower(strings.TrimSpace(params.Name))
	if name == "" {
		return nil, shared.NewValidationError("name", "role name cannot be empty")
	}
	if len(name) > 100 {
		return nil, shared.NewValidationError("name", "role name cannot exceed 100 characters")
	}
	if !roleNameRegex.MatchString(name) {
		return nil, shared.NewValidationError("name", "role name must be lowercase alphanumeric with underscores, starting with a letter")
	}

	displayName := strings.TrimSpace(params.DisplayName)
	if displayName == "" {
		return nil, shared.NewValidationError("display_name", "display name cannot be empty")
	}
	if len(displayName) > 255 {
		return nil, shared.NewValidationError("display_name", "display name cannot exceed 255 characters")
	}

	permissionIDs := make([]uuid.UUID, 0)
	if params.PermissionIDs != nil {
		seen := make(map[uuid.UUID]bool)
		for _, id := range params.PermissionIDs {
			if !seen[id] {
				permissionIDs = append(permissionIDs, id)
				seen[id] = true
			}
		}
	}

	now := time.Now().UTC()
	role := &Role{
		AggregateRoot: shared.NewAggregateRoot(),
		name:          name,
		displayName:   displayName,
		description:   params.Description,
		permissionIDs: permissionIDs,
		isSystem:      params.IsSystem,
		isDefault:     params.IsDefault,
		priority:      params.Priority,
		createdAt:     now,
		updatedAt:     now,
	}

	role.AddDomainEvent(NewRoleCreatedEvent(role.ID(), name, displayName))

	return role, nil
}

type ReconstructRoleParams struct {
	ID            uuid.UUID
	Name          string
	DisplayName   string
	Description   string
	PermissionIDs []uuid.UUID
	IsSystem      bool
	IsDefault     bool
	Priority      int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func ReconstructRole(params ReconstructRoleParams) (*Role, error) {
	if params.Name == "" {
		return nil, shared.NewValidationError("name", "role name cannot be empty")
	}

	permissionIDs := make([]uuid.UUID, 0)
	if params.PermissionIDs != nil {
		permissionIDs = append(permissionIDs, params.PermissionIDs...)
	}

	return &Role{
		AggregateRoot: shared.NewAggregateRootWithID(params.ID),
		name:          params.Name,
		displayName:   params.DisplayName,
		description:   params.Description,
		permissionIDs: permissionIDs,
		isSystem:      params.IsSystem,
		isDefault:     params.IsDefault,
		priority:      params.Priority,
		createdAt:     params.CreatedAt,
		updatedAt:     params.UpdatedAt,
	}, nil
}

func (r *Role) Name() string {
	return r.name
}

func (r *Role) DisplayName() string {
	return r.displayName
}

func (r *Role) Description() string {
	return r.description
}

func (r *Role) PermissionIDs() []uuid.UUID {
	result := make([]uuid.UUID, len(r.permissionIDs))
	copy(result, r.permissionIDs)
	return result
}

func (r *Role) IsSystem() bool {
	return r.isSystem
}

func (r *Role) IsDefault() bool {
	return r.isDefault
}

func (r *Role) Priority() int {
	return r.priority
}

func (r *Role) CreatedAt() time.Time {
	return r.createdAt
}

func (r *Role) UpdatedAt() time.Time {
	return r.updatedAt
}

func (r *Role) HasPermission(permissionID uuid.UUID) bool {
	for _, id := range r.permissionIDs {
		if id == permissionID {
			return true
		}
	}
	return false
}

func (r *Role) AddPermission(permissionID uuid.UUID) error {
	if r.HasPermission(permissionID) {
		return ErrPermissionAlreadyAssigned
	}

	r.permissionIDs = append(r.permissionIDs, permissionID)
	r.updatedAt = time.Now().UTC()
	r.AddDomainEvent(NewRolePermissionAddedEvent(r.ID(), permissionID))

	return nil
}

func (r *Role) RemovePermission(permissionID uuid.UUID) error {
	index := -1
	for i, id := range r.permissionIDs {
		if id == permissionID {
			index = i
			break
		}
	}

	if index == -1 {
		return ErrPermissionNotAssigned
	}

	r.permissionIDs = append(r.permissionIDs[:index], r.permissionIDs[index+1:]...)
	r.updatedAt = time.Now().UTC()
	r.AddDomainEvent(NewRolePermissionRemovedEvent(r.ID(), permissionID))

	return nil
}

func (r *Role) SetPermissions(permissionIDs []uuid.UUID) {
	oldPermissionIDs := r.permissionIDs

	seen := make(map[uuid.UUID]bool)
	newPermissionIDs := make([]uuid.UUID, 0)
	for _, id := range permissionIDs {
		if !seen[id] {
			newPermissionIDs = append(newPermissionIDs, id)
			seen[id] = true
		}
	}

	r.permissionIDs = newPermissionIDs
	r.updatedAt = time.Now().UTC()
	r.AddDomainEvent(NewRolePermissionsUpdatedEvent(r.ID(), oldPermissionIDs, newPermissionIDs))
}

func (r *Role) UpdateDetails(displayName, description string) error {
	changed := false

	if displayName != "" && displayName != r.displayName {
		if len(displayName) > 255 {
			return shared.NewValidationError("display_name", "display name cannot exceed 255 characters")
		}
		r.displayName = displayName
		changed = true
	}

	if description != r.description {
		r.description = description
		changed = true
	}

	if changed {
		r.updatedAt = time.Now().UTC()
		r.AddDomainEvent(NewRoleUpdatedEvent(r.ID(), r.name))
	}

	return nil
}

func (r *Role) CanBeDeleted() bool {
	return !r.isSystem && !r.isDefault
}

func (r *Role) CanBeModified() bool {
	return !r.isSystem
}
