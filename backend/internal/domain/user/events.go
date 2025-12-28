package user

import (
	"time"

	"github.com/google/uuid"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/shared"
)

const (
	EventTypeUserCreated       = "user.created"
	EventTypeUserActivated     = "user.activated"
	EventTypeUserDeactivated   = "user.deactivated"
	EventTypeUserBanned        = "user.banned"
	EventTypePasswordChanged   = "user.password_changed"
	EventTypeProfileUpdated    = "user.profile_updated"
	EventTypeUserDeleted       = "user.deleted"
	EventTypeUserRoleAssigned  = "user.role.assigned"
	EventTypeUserRoleRevoked   = "user.role.revoked"
	EventTypeUserRolesUpdated  = "user.roles.updated"
)

type UserCreatedEvent struct {
	shared.BaseDomainEvent
	Email     string
	FullName  string
}

func NewUserCreatedEvent(userID uuid.UUID, email, fullName string) UserCreatedEvent {
	return UserCreatedEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent(userID, EventTypeUserCreated),
		Email:           email,
		FullName:        fullName,
	}
}

type UserActivatedEvent struct {
	shared.BaseDomainEvent
}

func NewUserActivatedEvent(userID uuid.UUID) UserActivatedEvent {
	return UserActivatedEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent(userID, EventTypeUserActivated),
	}
}

type UserDeactivatedEvent struct {
	shared.BaseDomainEvent
}

func NewUserDeactivatedEvent(userID uuid.UUID) UserDeactivatedEvent {
	return UserDeactivatedEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent(userID, EventTypeUserDeactivated),
	}
}

type UserBannedEvent struct {
	shared.BaseDomainEvent
	Reason string
}

func NewUserBannedEvent(userID uuid.UUID, reason string) UserBannedEvent {
	return UserBannedEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent(userID, EventTypeUserBanned),
		Reason:          reason,
	}
}

type PasswordChangedEvent struct {
	shared.BaseDomainEvent
}

func NewPasswordChangedEvent(userID uuid.UUID) PasswordChangedEvent {
	return PasswordChangedEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent(userID, EventTypePasswordChanged),
	}
}

type ProfileUpdatedEvent struct {
	shared.BaseDomainEvent
	ChangedFields []string
}

func NewProfileUpdatedEvent(userID uuid.UUID, changedFields []string) ProfileUpdatedEvent {
	return ProfileUpdatedEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent(userID, EventTypeProfileUpdated),
		ChangedFields:   changedFields,
	}
}

type UserDeletedEvent struct {
	shared.BaseDomainEvent
	DeletedAt time.Time
}

func NewUserDeletedEvent(userID uuid.UUID, deletedAt time.Time) UserDeletedEvent {
	return UserDeletedEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent(userID, EventTypeUserDeleted),
		DeletedAt:       deletedAt,
	}
}

type UserRoleAssignedEvent struct {
	shared.BaseDomainEvent
	RoleID uuid.UUID
}

func NewUserRoleAssignedEvent(userID, roleID uuid.UUID) UserRoleAssignedEvent {
	return UserRoleAssignedEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent(userID, EventTypeUserRoleAssigned),
		RoleID:          roleID,
	}
}

type UserRoleRevokedEvent struct {
	shared.BaseDomainEvent
	RoleID uuid.UUID
}

func NewUserRoleRevokedEvent(userID, roleID uuid.UUID) UserRoleRevokedEvent {
	return UserRoleRevokedEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent(userID, EventTypeUserRoleRevoked),
		RoleID:          roleID,
	}
}

type UserRolesUpdatedEvent struct {
	shared.BaseDomainEvent
	OldRoleIDs []uuid.UUID
	NewRoleIDs []uuid.UUID
}

func NewUserRolesUpdatedEvent(userID uuid.UUID, oldRoleIDs, newRoleIDs []uuid.UUID) UserRolesUpdatedEvent {
	return UserRolesUpdatedEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent(userID, EventTypeUserRolesUpdated),
		OldRoleIDs:      oldRoleIDs,
		NewRoleIDs:      newRoleIDs,
	}
}
