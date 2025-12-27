package testutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/shared"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/user"
)

func AssertNotFoundError(t *testing.T, err error) {
	t.Helper()
	require.Error(t, err)
	assert.True(t, shared.IsNotFoundError(err), "expected NotFoundError, got %T: %v", err, err)
}

func AssertValidationError(t *testing.T, err error) {
	t.Helper()
	require.Error(t, err)
	assert.True(t, shared.IsValidationError(err), "expected ValidationError, got %T: %v", err, err)
}

func AssertConflictError(t *testing.T, err error) {
	t.Helper()
	require.Error(t, err)
	assert.True(t, shared.IsConflictError(err), "expected ConflictError, got %T: %v", err, err)
}

func AssertAuthorizationError(t *testing.T, err error) {
	t.Helper()
	require.Error(t, err)
	assert.True(t, shared.IsAuthorizationError(err), "expected AuthorizationError, got %T: %v", err, err)
}

func AssertBusinessRuleError(t *testing.T, err error) {
	t.Helper()
	require.Error(t, err)
	assert.True(t, shared.IsBusinessRuleViolationError(err), "expected BusinessRuleViolationError, got %T: %v", err, err)
}

func AssertInvalidStatusTransitionError(t *testing.T, err error) {
	t.Helper()
	require.Error(t, err)
	assert.True(t, shared.IsInvalidStatusTransitionError(err), "expected InvalidStatusTransitionError, got %T: %v", err, err)
}

func AssertUserNotFoundError(t *testing.T, err error) {
	t.Helper()
	require.Error(t, err)
	var notFoundError *shared.NotFoundError
	if assert.ErrorAs(t, err, &notFoundError, "expected NotFoundError, got %T: %v", err, err) {
		assert.Equal(t, "User", notFoundError.EntityType, "expected User entity type")
	}
}

func AssertEmailAlreadyExistsError(t *testing.T, err error) {
	t.Helper()
	require.Error(t, err)
	var conflictError *shared.ConflictError
	if assert.ErrorAs(t, err, &conflictError, "expected ConflictError, got %T: %v", err, err) {
		assert.Equal(t, "User", conflictError.EntityType, "expected User entity type")
		assert.Equal(t, "email", conflictError.Field, "expected email field")
	}
}

func AssertUserStatus(t *testing.T, u *user.User, expected user.Status) {
	t.Helper()
	assert.Equal(t, expected, u.Status(), "expected status %s, got %s", expected, u.Status())
}

func AssertUserEmail(t *testing.T, u *user.User, expected string) {
	t.Helper()
	assert.Equal(t, expected, u.Email().String(), "expected email %s, got %s", expected, u.Email().String())
}

func AssertDomainEventPublished(t *testing.T, eventBus *MockEventBus, eventType string) {
	t.Helper()
	found := false
	for _, event := range eventBus.PublishedEvents {
		if event.EventType() == eventType {
			found = true
			break
		}
	}
	assert.True(t, found, "expected event %s to be published", eventType)
}

func AssertDomainEventCount(t *testing.T, eventBus *MockEventBus, expectedCount int) {
	t.Helper()
	assert.Len(t, eventBus.PublishedEvents, expectedCount, "expected %d events, got %d", expectedCount, len(eventBus.PublishedEvents))
}

func AssertNoError(t *testing.T, err error) {
	t.Helper()
	require.NoError(t, err)
}

func AssertError(t *testing.T, err error) {
	t.Helper()
	require.Error(t, err)
}
