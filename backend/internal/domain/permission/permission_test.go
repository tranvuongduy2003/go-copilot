package permission

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPermission(t *testing.T) {
	tests := []struct {
		name        string
		params      NewPermissionParams
		wantErr     bool
		errContains string
	}{
		{
			name: "valid permission",
			params: NewPermissionParams{
				Resource:    "users",
				Action:      "create",
				Description: "Create users",
				IsSystem:    false,
			},
			wantErr: false,
		},
		{
			name: "valid system permission",
			params: NewPermissionParams{
				Resource:    "system",
				Action:      "admin",
				Description: "System admin",
				IsSystem:    true,
			},
			wantErr: false,
		},
		{
			name: "empty resource",
			params: NewPermissionParams{
				Resource: "",
				Action:   "create",
			},
			wantErr:     true,
			errContains: "resource cannot be empty",
		},
		{
			name: "empty action",
			params: NewPermissionParams{
				Resource: "users",
				Action:   "",
			},
			wantErr:     true,
			errContains: "action cannot be empty",
		},
		{
			name: "invalid resource format",
			params: NewPermissionParams{
				Resource: "Users@123",
				Action:   "create",
			},
			wantErr:     true,
			errContains: "must be lowercase",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			perm, err := NewPermission(tt.params)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				assert.Nil(t, perm)
			} else {
				require.NoError(t, err)
				require.NotNil(t, perm)
				assert.Equal(t, tt.params.Resource, perm.Resource().String())
				assert.Equal(t, tt.params.Action, perm.Action().String())
				assert.Equal(t, tt.params.Description, perm.Description())
				assert.Equal(t, tt.params.IsSystem, perm.IsSystem())
				assert.NotEqual(t, uuid.Nil, perm.ID())
			}
		})
	}
}

func TestReconstructPermission(t *testing.T) {
	id := uuid.New()
	now := time.Now().UTC()

	perm, err := ReconstructPermission(ReconstructPermissionParams{
		ID:          id,
		Resource:    "users",
		Action:      "read",
		Description: "Read users",
		IsSystem:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	})

	require.NoError(t, err)
	require.NotNil(t, perm)
	assert.Equal(t, id, perm.ID())
	assert.Equal(t, "users", perm.Resource().String())
	assert.Equal(t, "read", perm.Action().String())
	assert.True(t, perm.IsSystem())
}

func TestPermission_Code(t *testing.T) {
	perm, err := NewPermission(NewPermissionParams{
		Resource: "users",
		Action:   "create",
	})
	require.NoError(t, err)

	assert.Equal(t, "users:create", perm.CodeString())
	assert.Equal(t, "users:create", perm.Code().String())
}

func TestPermission_UpdateDescription(t *testing.T) {
	perm, err := NewPermission(NewPermissionParams{
		Resource:    "users",
		Action:      "create",
		Description: "Original description",
	})
	require.NoError(t, err)

	originalUpdatedAt := perm.UpdatedAt()
	time.Sleep(time.Millisecond)

	perm.UpdateDescription("New description")
	assert.Equal(t, "New description", perm.Description())
	assert.True(t, perm.UpdatedAt().After(originalUpdatedAt))

	updatedAt := perm.UpdatedAt()
	perm.UpdateDescription("New description")
	assert.Equal(t, updatedAt, perm.UpdatedAt())
}

func TestPermission_CanBeDeleted(t *testing.T) {
	systemPerm, _ := NewPermission(NewPermissionParams{
		Resource: "system",
		Action:   "admin",
		IsSystem: true,
	})
	assert.False(t, systemPerm.CanBeDeleted())

	regularPerm, _ := NewPermission(NewPermissionParams{
		Resource: "posts",
		Action:   "create",
		IsSystem: false,
	})
	assert.True(t, regularPerm.CanBeDeleted())
}

func TestPermission_Equals(t *testing.T) {
	perm1, _ := NewPermission(NewPermissionParams{
		Resource: "users",
		Action:   "create",
	})

	perm2, _ := NewPermission(NewPermissionParams{
		Resource: "users",
		Action:   "create",
	})

	perm3, _ := NewPermission(NewPermissionParams{
		Resource: "users",
		Action:   "delete",
	})

	assert.True(t, perm1.Equals(perm2))
	assert.False(t, perm1.Equals(perm3))
	assert.False(t, perm1.Equals(nil))
}
