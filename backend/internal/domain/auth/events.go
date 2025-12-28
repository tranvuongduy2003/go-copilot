package auth

import (
	"github.com/google/uuid"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/shared"
)

const (
	EventTypeUserLoggedIn             = "auth.user.logged_in"
	EventTypeUserLoggedOut            = "auth.user.logged_out"
	EventTypeUserRegistered           = "auth.user.registered"
	EventTypePasswordResetRequested   = "auth.password_reset.requested"
	EventTypePasswordReset            = "auth.password_reset.completed"
	EventTypeRefreshTokenRotated      = "auth.refresh_token.rotated"
	EventTypeLoginFailed              = "auth.login.failed"
	EventTypeAccountLocked            = "auth.account.locked"
	EventTypeSessionRevoked           = "auth.session.revoked"
)

type UserLoggedInEvent struct {
	shared.BaseDomainEvent
	Email     string
	IPAddress string
	UserAgent string
}

func NewUserLoggedInEvent(userID uuid.UUID, email, ipAddress, userAgent string) UserLoggedInEvent {
	return UserLoggedInEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent(userID, EventTypeUserLoggedIn),
		Email:           email,
		IPAddress:       ipAddress,
		UserAgent:       userAgent,
	}
}

type UserLoggedOutEvent struct {
	shared.BaseDomainEvent
	LogoutAll bool
}

func NewUserLoggedOutEvent(userID uuid.UUID, logoutAll bool) UserLoggedOutEvent {
	return UserLoggedOutEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent(userID, EventTypeUserLoggedOut),
		LogoutAll:       logoutAll,
	}
}

type UserRegisteredEvent struct {
	shared.BaseDomainEvent
	Email    string
	FullName string
}

func NewUserRegisteredEvent(userID uuid.UUID, email, fullName string) UserRegisteredEvent {
	return UserRegisteredEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent(userID, EventTypeUserRegistered),
		Email:           email,
		FullName:        fullName,
	}
}

type PasswordResetRequestedEvent struct {
	shared.BaseDomainEvent
	Email string
}

func NewPasswordResetRequestedEvent(userID uuid.UUID, email string) PasswordResetRequestedEvent {
	return PasswordResetRequestedEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent(userID, EventTypePasswordResetRequested),
		Email:           email,
	}
}

type PasswordResetEvent struct {
	shared.BaseDomainEvent
}

func NewPasswordResetEvent(userID uuid.UUID) PasswordResetEvent {
	return PasswordResetEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent(userID, EventTypePasswordReset),
	}
}

type RefreshTokenRotatedEvent struct {
	shared.BaseDomainEvent
	OldTokenID uuid.UUID
	NewTokenID uuid.UUID
}

func NewRefreshTokenRotatedEvent(userID, oldTokenID, newTokenID uuid.UUID) RefreshTokenRotatedEvent {
	return RefreshTokenRotatedEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent(userID, EventTypeRefreshTokenRotated),
		OldTokenID:      oldTokenID,
		NewTokenID:      newTokenID,
	}
}

type LoginFailedEvent struct {
	shared.BaseDomainEvent
	Email         string
	IPAddress     string
	FailureReason string
	AttemptCount  int
}

func NewLoginFailedEvent(userID uuid.UUID, email, ipAddress, reason string, attemptCount int) LoginFailedEvent {
	return LoginFailedEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent(userID, EventTypeLoginFailed),
		Email:           email,
		IPAddress:       ipAddress,
		FailureReason:   reason,
		AttemptCount:    attemptCount,
	}
}

type AccountLockedEvent struct {
	shared.BaseDomainEvent
	Email         string
	LockDuration  int
	FailedAttempts int
}

func NewAccountLockedEvent(userID uuid.UUID, email string, lockDuration, failedAttempts int) AccountLockedEvent {
	return AccountLockedEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent(userID, EventTypeAccountLocked),
		Email:           email,
		LockDuration:    lockDuration,
		FailedAttempts:  failedAttempts,
	}
}

type SessionRevokedEvent struct {
	shared.BaseDomainEvent
	SessionID uuid.UUID
}

func NewSessionRevokedEvent(userID, sessionID uuid.UUID) SessionRevokedEvent {
	return SessionRevokedEvent{
		BaseDomainEvent: shared.NewBaseDomainEvent(userID, EventTypeSessionRevoked),
		SessionID:       sessionID,
	}
}
