package role

import (
	"github.com/google/uuid"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/shared"
)

const (
	EventTypeRoleCreated            = "role.created"
	EventTypeRoleUpdated            = "role.updated"
	EventTypeRoleDeleted            = "role.deleted"
	EventTypeRolePermissionAdded    = "role.permission.added"
	EventTypeRolePermissionRemoved  = "role.permission.removed"
	EventTypeRolePermissionsUpdated = "role.permissions.updated"
)

type RoleCreatedEvent struct {
	shared.BaseDomainEvent
	RoleName    string
	DisplayName string
}

func NewRoleCreatedEvent(roleID uuid.UUID, name, displayName string) RoleCreatedEvent {
	return RoleCreatedEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent(roleID, EventTypeRoleCreated),
		RoleName:        name,
		DisplayName:     displayName,
	}
}

type RoleUpdatedEvent struct {
	shared.BaseDomainEvent
	RoleName string
}

func NewRoleUpdatedEvent(roleID uuid.UUID, name string) RoleUpdatedEvent {
	return RoleUpdatedEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent(roleID, EventTypeRoleUpdated),
		RoleName:        name,
	}
}

type RoleDeletedEvent struct {
	shared.BaseDomainEvent
	RoleName string
}

func NewRoleDeletedEvent(roleID uuid.UUID, name string) RoleDeletedEvent {
	return RoleDeletedEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent(roleID, EventTypeRoleDeleted),
		RoleName:        name,
	}
}

type RolePermissionAddedEvent struct {
	shared.BaseDomainEvent
	PermissionID uuid.UUID
}

func NewRolePermissionAddedEvent(roleID, permissionID uuid.UUID) RolePermissionAddedEvent {
	return RolePermissionAddedEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent(roleID, EventTypeRolePermissionAdded),
		PermissionID:    permissionID,
	}
}

type RolePermissionRemovedEvent struct {
	shared.BaseDomainEvent
	PermissionID uuid.UUID
}

func NewRolePermissionRemovedEvent(roleID, permissionID uuid.UUID) RolePermissionRemovedEvent {
	return RolePermissionRemovedEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent(roleID, EventTypeRolePermissionRemoved),
		PermissionID:    permissionID,
	}
}

type RolePermissionsUpdatedEvent struct {
	shared.BaseDomainEvent
	OldPermissionIDs []uuid.UUID
	NewPermissionIDs []uuid.UUID
}

func NewRolePermissionsUpdatedEvent(roleID uuid.UUID, oldIDs, newIDs []uuid.UUID) RolePermissionsUpdatedEvent {
	return RolePermissionsUpdatedEvent{
		BaseDomainEvent:  shared.NewBaseDomainEvent(roleID, EventTypeRolePermissionsUpdated),
		OldPermissionIDs: oldIDs,
		NewPermissionIDs: newIDs,
	}
}
