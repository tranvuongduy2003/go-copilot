package user

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestUser(t *testing.T) *User {
	t.Helper()
	user, err := NewUser(NewUserParams{
		Email:        "test@example.com",
		PasswordHash: "$2a$10$hashedpassword",
		FullName:     "Test User",
	})
	require.NoError(t, err)
	return user
}

func createTestUserWithRoles(t *testing.T, roleIDs []uuid.UUID) *User {
	t.Helper()
	now := time.Now().UTC()
	user, err := ReconstructUser(ReconstructUserParams{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "$2a$10$hashedpassword",
		FullName:     "Test User",
		Status:       StatusActive,
		RoleIDs:      roleIDs,
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	require.NoError(t, err)
	return user
}

func TestUser_RoleIDs(t *testing.T) {
	roleID1 := uuid.New()
	roleID2 := uuid.New()
	roleIDs := []uuid.UUID{roleID1, roleID2}

	user := createTestUserWithRoles(t, roleIDs)

	result := user.RoleIDs()

	assert.Len(t, result, 2)
	assert.Contains(t, result, roleID1)
	assert.Contains(t, result, roleID2)
}

func TestUser_RoleIDs_ReturnsCopy(t *testing.T) {
	roleID := uuid.New()
	user := createTestUserWithRoles(t, []uuid.UUID{roleID})

	result := user.RoleIDs()
	result[0] = uuid.New()

	assert.Equal(t, roleID, user.RoleIDs()[0])
}

func TestUser_RoleIDs_EmptyWhenNoRoles(t *testing.T) {
	user := createTestUser(t)

	result := user.RoleIDs()

	assert.Empty(t, result)
}

func TestUser_HasRole(t *testing.T) {
	roleID1 := uuid.New()
	roleID2 := uuid.New()
	unassignedRoleID := uuid.New()

	user := createTestUserWithRoles(t, []uuid.UUID{roleID1, roleID2})

	assert.True(t, user.HasRole(roleID1))
	assert.True(t, user.HasRole(roleID2))
	assert.False(t, user.HasRole(unassignedRoleID))
}

func TestUser_HasRole_EmptyRoles(t *testing.T) {
	user := createTestUser(t)

	assert.False(t, user.HasRole(uuid.New()))
}

func TestUser_AssignRole(t *testing.T) {
	tests := []struct {
		name        string
		initialRole []uuid.UUID
		newRoleID   uuid.UUID
		wantErr     bool
		errType     error
	}{
		{
			name:        "successfully assign role to user without roles",
			initialRole: nil,
			newRoleID:   uuid.New(),
			wantErr:     false,
		},
		{
			name:        "successfully assign role to user with existing roles",
			initialRole: []uuid.UUID{uuid.New()},
			newRoleID:   uuid.New(),
			wantErr:     false,
		},
		{
			name: "fail when role already assigned",
			initialRole: func() []uuid.UUID {
				id := uuid.New()
				return []uuid.UUID{id}
			}(),
			newRoleID: uuid.Nil,
			wantErr:   true,
			errType:   ErrRoleAlreadyAssigned,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var user *User
			if tt.initialRole != nil {
				user = createTestUserWithRoles(t, tt.initialRole)
			} else {
				user = createTestUser(t)
			}

			roleID := tt.newRoleID
			if tt.name == "fail when role already assigned" {
				roleID = tt.initialRole[0]
			}

			err := user.AssignRole(roleID)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.errType, err)
			} else {
				require.NoError(t, err)
				assert.True(t, user.HasRole(roleID))
			}
		})
	}
}

func TestUser_AssignRole_PublishesDomainEvent(t *testing.T) {
	user := createTestUser(t)
	user.ClearDomainEvents()
	roleID := uuid.New()

	err := user.AssignRole(roleID)

	require.NoError(t, err)
	events := user.DomainEvents()
	assert.Len(t, events, 1)

	event, ok := events[0].(UserRoleAssignedEvent)
	assert.True(t, ok)
	assert.Equal(t, user.ID(), event.AggregateID())
	assert.Equal(t, roleID, event.RoleID)
}

func TestUser_AssignRole_UpdatesTimestamp(t *testing.T) {
	user := createTestUser(t)
	originalUpdatedAt := user.UpdatedAt()

	time.Sleep(1 * time.Millisecond)
	err := user.AssignRole(uuid.New())

	require.NoError(t, err)
	assert.True(t, user.UpdatedAt().After(originalUpdatedAt))
}

func TestUser_RevokeRole(t *testing.T) {
	tests := []struct {
		name        string
		setupUser   func(t *testing.T) (*User, uuid.UUID)
		wantErr     bool
		errType     error
	}{
		{
			name: "successfully revoke assigned role",
			setupUser: func(t *testing.T) (*User, uuid.UUID) {
				roleID := uuid.New()
				user := createTestUserWithRoles(t, []uuid.UUID{roleID})
				return user, roleID
			},
			wantErr: false,
		},
		{
			name: "successfully revoke one of multiple roles",
			setupUser: func(t *testing.T) (*User, uuid.UUID) {
				roleID1 := uuid.New()
				roleID2 := uuid.New()
				user := createTestUserWithRoles(t, []uuid.UUID{roleID1, roleID2})
				return user, roleID1
			},
			wantErr: false,
		},
		{
			name: "fail when role not assigned",
			setupUser: func(t *testing.T) (*User, uuid.UUID) {
				user := createTestUser(t)
				return user, uuid.New()
			},
			wantErr: true,
			errType: ErrRoleNotAssigned,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, roleID := tt.setupUser(t)

			err := user.RevokeRole(roleID)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.errType, err)
			} else {
				require.NoError(t, err)
				assert.False(t, user.HasRole(roleID))
			}
		})
	}
}

func TestUser_RevokeRole_PublishesDomainEvent(t *testing.T) {
	roleID := uuid.New()
	user := createTestUserWithRoles(t, []uuid.UUID{roleID})
	user.ClearDomainEvents()

	err := user.RevokeRole(roleID)

	require.NoError(t, err)
	events := user.DomainEvents()
	assert.Len(t, events, 1)

	event, ok := events[0].(UserRoleRevokedEvent)
	assert.True(t, ok)
	assert.Equal(t, user.ID(), event.AggregateID())
	assert.Equal(t, roleID, event.RoleID)
}

func TestUser_RevokeRole_UpdatesTimestamp(t *testing.T) {
	roleID := uuid.New()
	user := createTestUserWithRoles(t, []uuid.UUID{roleID})
	originalUpdatedAt := user.UpdatedAt()

	time.Sleep(1 * time.Millisecond)
	err := user.RevokeRole(roleID)

	require.NoError(t, err)
	assert.True(t, user.UpdatedAt().After(originalUpdatedAt))
}

func TestUser_SetRoles(t *testing.T) {
	tests := []struct {
		name           string
		initialRoleIDs []uuid.UUID
		newRoleIDs     []uuid.UUID
		expectedCount  int
	}{
		{
			name:           "set roles on user without roles",
			initialRoleIDs: nil,
			newRoleIDs:     []uuid.UUID{uuid.New(), uuid.New()},
			expectedCount:  2,
		},
		{
			name:           "replace existing roles",
			initialRoleIDs: []uuid.UUID{uuid.New()},
			newRoleIDs:     []uuid.UUID{uuid.New(), uuid.New()},
			expectedCount:  2,
		},
		{
			name:           "clear all roles with empty slice",
			initialRoleIDs: []uuid.UUID{uuid.New(), uuid.New()},
			newRoleIDs:     []uuid.UUID{},
			expectedCount:  0,
		},
		{
			name:           "deduplicate roles",
			initialRoleIDs: nil,
			newRoleIDs: func() []uuid.UUID {
				id := uuid.New()
				return []uuid.UUID{id, id, uuid.New()}
			}(),
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var user *User
			if tt.initialRoleIDs != nil {
				user = createTestUserWithRoles(t, tt.initialRoleIDs)
			} else {
				user = createTestUser(t)
			}

			user.SetRoles(tt.newRoleIDs)

			assert.Len(t, user.RoleIDs(), tt.expectedCount)
		})
	}
}

func TestUser_SetRoles_PublishesDomainEvent(t *testing.T) {
	oldRoleID := uuid.New()
	user := createTestUserWithRoles(t, []uuid.UUID{oldRoleID})
	user.ClearDomainEvents()

	newRoleID := uuid.New()
	user.SetRoles([]uuid.UUID{newRoleID})

	events := user.DomainEvents()
	assert.Len(t, events, 1)

	event, ok := events[0].(UserRolesUpdatedEvent)
	assert.True(t, ok)
	assert.Equal(t, user.ID(), event.AggregateID())
	assert.Contains(t, event.OldRoleIDs, oldRoleID)
	assert.Contains(t, event.NewRoleIDs, newRoleID)
}

func TestUser_SetRoles_UpdatesTimestamp(t *testing.T) {
	user := createTestUser(t)
	originalUpdatedAt := user.UpdatedAt()

	time.Sleep(1 * time.Millisecond)
	user.SetRoles([]uuid.UUID{uuid.New()})

	assert.True(t, user.UpdatedAt().After(originalUpdatedAt))
}

func TestUser_NewUser_StartsWithNoRoles(t *testing.T) {
	user := createTestUser(t)

	assert.Empty(t, user.RoleIDs())
}

func TestUser_ReconstructUser_PreservesRoles(t *testing.T) {
	roleID1 := uuid.New()
	roleID2 := uuid.New()

	user := createTestUserWithRoles(t, []uuid.UUID{roleID1, roleID2})

	assert.Len(t, user.RoleIDs(), 2)
	assert.True(t, user.HasRole(roleID1))
	assert.True(t, user.HasRole(roleID2))
}
