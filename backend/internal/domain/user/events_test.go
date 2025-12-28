package user

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewUserCreatedEvent(t *testing.T) {
	userID := uuid.New()
	email := "test@example.com"
	fullName := "Test User"

	event := NewUserCreatedEvent(userID, email, fullName)

	assert.Equal(t, userID, event.AggregateID())
	assert.Equal(t, EventTypeUserCreated, event.EventType())
	assert.Equal(t, email, event.Email)
	assert.Equal(t, fullName, event.FullName)
	assert.NotZero(t, event.OccurredAt())
}

func TestNewUserActivatedEvent(t *testing.T) {
	userID := uuid.New()

	event := NewUserActivatedEvent(userID)

	assert.Equal(t, userID, event.AggregateID())
	assert.Equal(t, EventTypeUserActivated, event.EventType())
}

func TestNewUserDeactivatedEvent(t *testing.T) {
	userID := uuid.New()

	event := NewUserDeactivatedEvent(userID)

	assert.Equal(t, userID, event.AggregateID())
	assert.Equal(t, EventTypeUserDeactivated, event.EventType())
}

func TestNewUserBannedEvent(t *testing.T) {
	userID := uuid.New()
	reason := "Violation of terms"

	event := NewUserBannedEvent(userID, reason)

	assert.Equal(t, userID, event.AggregateID())
	assert.Equal(t, EventTypeUserBanned, event.EventType())
	assert.Equal(t, reason, event.Reason)
}

func TestNewPasswordChangedEvent(t *testing.T) {
	userID := uuid.New()

	event := NewPasswordChangedEvent(userID)

	assert.Equal(t, userID, event.AggregateID())
	assert.Equal(t, EventTypePasswordChanged, event.EventType())
}

func TestNewProfileUpdatedEvent(t *testing.T) {
	userID := uuid.New()
	changedFields := []string{"full_name", "avatar"}

	event := NewProfileUpdatedEvent(userID, changedFields)

	assert.Equal(t, userID, event.AggregateID())
	assert.Equal(t, EventTypeProfileUpdated, event.EventType())
	assert.Equal(t, changedFields, event.ChangedFields)
}

func TestNewUserDeletedEvent(t *testing.T) {
	userID := uuid.New()
	deletedAt := time.Now().UTC()

	event := NewUserDeletedEvent(userID, deletedAt)

	assert.Equal(t, userID, event.AggregateID())
	assert.Equal(t, EventTypeUserDeleted, event.EventType())
	assert.Equal(t, deletedAt, event.DeletedAt)
}

func TestNewUserRoleAssignedEvent(t *testing.T) {
	userID := uuid.New()
	roleID := uuid.New()

	event := NewUserRoleAssignedEvent(userID, roleID)

	assert.Equal(t, userID, event.AggregateID())
	assert.Equal(t, EventTypeUserRoleAssigned, event.EventType())
	assert.Equal(t, roleID, event.RoleID)
}

func TestNewUserRoleRevokedEvent(t *testing.T) {
	userID := uuid.New()
	roleID := uuid.New()

	event := NewUserRoleRevokedEvent(userID, roleID)

	assert.Equal(t, userID, event.AggregateID())
	assert.Equal(t, EventTypeUserRoleRevoked, event.EventType())
	assert.Equal(t, roleID, event.RoleID)
}

func TestNewUserRolesUpdatedEvent(t *testing.T) {
	userID := uuid.New()
	oldRoleIDs := []uuid.UUID{uuid.New()}
	newRoleIDs := []uuid.UUID{uuid.New(), uuid.New()}

	event := NewUserRolesUpdatedEvent(userID, oldRoleIDs, newRoleIDs)

	assert.Equal(t, userID, event.AggregateID())
	assert.Equal(t, EventTypeUserRolesUpdated, event.EventType())
	assert.Equal(t, oldRoleIDs, event.OldRoleIDs)
	assert.Equal(t, newRoleIDs, event.NewRoleIDs)
}
