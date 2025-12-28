package auth

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewUserLoggedInEvent(t *testing.T) {
	userID := uuid.New()
	email := "test@example.com"
	ipAddress := "192.168.1.1"
	userAgent := "Mozilla/5.0"

	event := NewUserLoggedInEvent(userID, email, ipAddress, userAgent)

	assert.Equal(t, userID, event.AggregateID())
	assert.Equal(t, EventTypeUserLoggedIn, event.EventType())
	assert.Equal(t, email, event.Email)
	assert.Equal(t, ipAddress, event.IPAddress)
	assert.Equal(t, userAgent, event.UserAgent)
	assert.NotZero(t, event.OccurredAt())
}

func TestNewUserLoggedOutEvent(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name      string
		logoutAll bool
	}{
		{"single logout", false},
		{"logout all", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := NewUserLoggedOutEvent(userID, tt.logoutAll)

			assert.Equal(t, userID, event.AggregateID())
			assert.Equal(t, EventTypeUserLoggedOut, event.EventType())
			assert.Equal(t, tt.logoutAll, event.LogoutAll)
		})
	}
}

func TestNewUserRegisteredEvent(t *testing.T) {
	userID := uuid.New()
	email := "test@example.com"
	fullName := "Test User"

	event := NewUserRegisteredEvent(userID, email, fullName)

	assert.Equal(t, userID, event.AggregateID())
	assert.Equal(t, EventTypeUserRegistered, event.EventType())
	assert.Equal(t, email, event.Email)
	assert.Equal(t, fullName, event.FullName)
}

func TestNewPasswordResetRequestedEvent(t *testing.T) {
	userID := uuid.New()
	email := "test@example.com"

	event := NewPasswordResetRequestedEvent(userID, email)

	assert.Equal(t, userID, event.AggregateID())
	assert.Equal(t, EventTypePasswordResetRequested, event.EventType())
	assert.Equal(t, email, event.Email)
}

func TestNewPasswordResetEvent(t *testing.T) {
	userID := uuid.New()

	event := NewPasswordResetEvent(userID)

	assert.Equal(t, userID, event.AggregateID())
	assert.Equal(t, EventTypePasswordReset, event.EventType())
}

func TestNewRefreshTokenRotatedEvent(t *testing.T) {
	userID := uuid.New()
	oldTokenID := uuid.New()
	newTokenID := uuid.New()

	event := NewRefreshTokenRotatedEvent(userID, oldTokenID, newTokenID)

	assert.Equal(t, userID, event.AggregateID())
	assert.Equal(t, EventTypeRefreshTokenRotated, event.EventType())
	assert.Equal(t, oldTokenID, event.OldTokenID)
	assert.Equal(t, newTokenID, event.NewTokenID)
}

func TestNewLoginFailedEvent(t *testing.T) {
	userID := uuid.New()
	email := "test@example.com"
	ipAddress := "192.168.1.1"
	reason := "invalid_password"
	attemptCount := 3

	event := NewLoginFailedEvent(userID, email, ipAddress, reason, attemptCount)

	assert.Equal(t, userID, event.AggregateID())
	assert.Equal(t, EventTypeLoginFailed, event.EventType())
	assert.Equal(t, email, event.Email)
	assert.Equal(t, ipAddress, event.IPAddress)
	assert.Equal(t, reason, event.FailureReason)
	assert.Equal(t, attemptCount, event.AttemptCount)
}

func TestNewAccountLockedEvent(t *testing.T) {
	userID := uuid.New()
	email := "test@example.com"
	lockDuration := 15
	failedAttempts := 5

	event := NewAccountLockedEvent(userID, email, lockDuration, failedAttempts)

	assert.Equal(t, userID, event.AggregateID())
	assert.Equal(t, EventTypeAccountLocked, event.EventType())
	assert.Equal(t, email, event.Email)
	assert.Equal(t, lockDuration, event.LockDuration)
	assert.Equal(t, failedAttempts, event.FailedAttempts)
}

func TestNewSessionRevokedEvent(t *testing.T) {
	userID := uuid.New()
	sessionID := uuid.New()

	event := NewSessionRevokedEvent(userID, sessionID)

	assert.Equal(t, userID, event.AggregateID())
	assert.Equal(t, EventTypeSessionRevoked, event.EventType())
	assert.Equal(t, sessionID, event.SessionID)
}
