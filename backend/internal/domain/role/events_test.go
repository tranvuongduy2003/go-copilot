package role

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewRoleCreatedEvent(t *testing.T) {
	roleID := uuid.New()
	name := "admin"
	displayName := "Administrator"

	event := NewRoleCreatedEvent(roleID, name, displayName)

	assert.Equal(t, roleID, event.AggregateID())
	assert.Equal(t, EventTypeRoleCreated, event.EventType())
	assert.Equal(t, name, event.RoleName)
	assert.Equal(t, displayName, event.DisplayName)
	assert.NotZero(t, event.OccurredAt())
}

func TestNewRoleUpdatedEvent(t *testing.T) {
	roleID := uuid.New()
	name := "admin"

	event := NewRoleUpdatedEvent(roleID, name)

	assert.Equal(t, roleID, event.AggregateID())
	assert.Equal(t, EventTypeRoleUpdated, event.EventType())
	assert.Equal(t, name, event.RoleName)
}

func TestNewRoleDeletedEvent(t *testing.T) {
	roleID := uuid.New()
	name := "admin"

	event := NewRoleDeletedEvent(roleID, name)

	assert.Equal(t, roleID, event.AggregateID())
	assert.Equal(t, EventTypeRoleDeleted, event.EventType())
	assert.Equal(t, name, event.RoleName)
}

func TestNewRolePermissionAddedEvent(t *testing.T) {
	roleID := uuid.New()
	permissionID := uuid.New()

	event := NewRolePermissionAddedEvent(roleID, permissionID)

	assert.Equal(t, roleID, event.AggregateID())
	assert.Equal(t, EventTypeRolePermissionAdded, event.EventType())
	assert.Equal(t, permissionID, event.PermissionID)
}

func TestNewRolePermissionRemovedEvent(t *testing.T) {
	roleID := uuid.New()
	permissionID := uuid.New()

	event := NewRolePermissionRemovedEvent(roleID, permissionID)

	assert.Equal(t, roleID, event.AggregateID())
	assert.Equal(t, EventTypeRolePermissionRemoved, event.EventType())
	assert.Equal(t, permissionID, event.PermissionID)
}

func TestNewRolePermissionsUpdatedEvent(t *testing.T) {
	roleID := uuid.New()
	oldIDs := []uuid.UUID{uuid.New(), uuid.New()}
	newIDs := []uuid.UUID{uuid.New(), uuid.New(), uuid.New()}

	event := NewRolePermissionsUpdatedEvent(roleID, oldIDs, newIDs)

	assert.Equal(t, roleID, event.AggregateID())
	assert.Equal(t, EventTypeRolePermissionsUpdated, event.EventType())
	assert.Equal(t, oldIDs, event.OldPermissionIDs)
	assert.Equal(t, newIDs, event.NewPermissionIDs)
}
