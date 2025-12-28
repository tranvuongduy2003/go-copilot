package role

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRole(t *testing.T) {
	tests := []struct {
		name        string
		params      NewRoleParams
		wantErr     bool
		errContains string
	}{
		{
			name: "valid role",
			params: NewRoleParams{
				Name:        "admin",
				DisplayName: "Administrator",
				Description: "Admin role",
				IsSystem:    false,
				IsDefault:   false,
				Priority:    10,
			},
			wantErr: false,
		},
		{
			name: "valid system role",
			params: NewRoleParams{
				Name:        "super_admin",
				DisplayName: "Super Administrator",
				IsSystem:    true,
				IsDefault:   false,
				Priority:    100,
			},
			wantErr: false,
		},
		{
			name: "valid role with permissions",
			params: NewRoleParams{
				Name:          "editor",
				DisplayName:   "Editor",
				PermissionIDs: []uuid.UUID{uuid.New(), uuid.New()},
			},
			wantErr: false,
		},
		{
			name: "deduplicates permissions",
			params: NewRoleParams{
				Name:          "viewer",
				DisplayName:   "Viewer",
				PermissionIDs: []uuid.UUID{uuid.MustParse("00000000-0000-0000-0000-000000000001"), uuid.MustParse("00000000-0000-0000-0000-000000000001")},
			},
			wantErr: false,
		},
		{
			name: "empty role name",
			params: NewRoleParams{
				Name:        "",
				DisplayName: "Test",
			},
			wantErr:     true,
			errContains: "role name cannot be empty",
		},
		{
			name: "empty display name",
			params: NewRoleParams{
				Name:        "test",
				DisplayName: "",
			},
			wantErr:     true,
			errContains: "display name cannot be empty",
		},
		{
			name: "invalid role name format",
			params: NewRoleParams{
				Name:        "Admin Role",
				DisplayName: "Admin Role",
			},
			wantErr:     true,
			errContains: "must be lowercase alphanumeric",
		},
		{
			name: "role name starts with number",
			params: NewRoleParams{
				Name:        "1admin",
				DisplayName: "Admin",
			},
			wantErr:     true,
			errContains: "must be lowercase alphanumeric",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			role, err := NewRole(tt.params)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				assert.Nil(t, role)
			} else {
				require.NoError(t, err)
				require.NotNil(t, role)
				assert.NotEqual(t, uuid.Nil, role.ID())
				assert.Equal(t, tt.params.DisplayName, role.DisplayName())
				assert.Equal(t, tt.params.IsSystem, role.IsSystem())
				assert.Equal(t, tt.params.IsDefault, role.IsDefault())
				assert.Equal(t, tt.params.Priority, role.Priority())
				assert.Len(t, role.DomainEvents(), 1)
			}
		})
	}
}

func TestRole_DeduplicatesPermissions(t *testing.T) {
	permID := uuid.New()
	role, err := NewRole(NewRoleParams{
		Name:          "test",
		DisplayName:   "Test",
		PermissionIDs: []uuid.UUID{permID, permID, permID},
	})

	require.NoError(t, err)
	assert.Len(t, role.PermissionIDs(), 1)
	assert.Equal(t, permID, role.PermissionIDs()[0])
}

func TestReconstructRole(t *testing.T) {
	id := uuid.New()
	permIDs := []uuid.UUID{uuid.New(), uuid.New()}
	now := time.Now().UTC()

	role, err := ReconstructRole(ReconstructRoleParams{
		ID:            id,
		Name:          "admin",
		DisplayName:   "Administrator",
		Description:   "Admin role",
		PermissionIDs: permIDs,
		IsSystem:      true,
		IsDefault:     false,
		Priority:      100,
		CreatedAt:     now,
		UpdatedAt:     now,
	})

	require.NoError(t, err)
	require.NotNil(t, role)
	assert.Equal(t, id, role.ID())
	assert.Equal(t, "admin", role.Name())
	assert.Equal(t, "Administrator", role.DisplayName())
	assert.True(t, role.IsSystem())
	assert.Len(t, role.PermissionIDs(), 2)
	assert.Empty(t, role.DomainEvents())
}

func TestRole_AddPermission(t *testing.T) {
	role, _ := NewRole(NewRoleParams{
		Name:        "editor",
		DisplayName: "Editor",
	})
	role.ClearDomainEvents()

	permID := uuid.New()
	err := role.AddPermission(permID)
	require.NoError(t, err)
	assert.True(t, role.HasPermission(permID))
	assert.Len(t, role.DomainEvents(), 1)

	err = role.AddPermission(permID)
	assert.ErrorIs(t, err, ErrPermissionAlreadyAssigned)
}

func TestRole_RemovePermission(t *testing.T) {
	permID := uuid.New()
	role, _ := NewRole(NewRoleParams{
		Name:          "editor",
		DisplayName:   "Editor",
		PermissionIDs: []uuid.UUID{permID},
	})
	role.ClearDomainEvents()

	err := role.RemovePermission(permID)
	require.NoError(t, err)
	assert.False(t, role.HasPermission(permID))
	assert.Len(t, role.DomainEvents(), 1)

	err = role.RemovePermission(permID)
	assert.ErrorIs(t, err, ErrPermissionNotAssigned)
}

func TestRole_SetPermissions(t *testing.T) {
	oldPermIDs := []uuid.UUID{uuid.New(), uuid.New()}
	role, _ := NewRole(NewRoleParams{
		Name:          "editor",
		DisplayName:   "Editor",
		PermissionIDs: oldPermIDs,
	})
	role.ClearDomainEvents()

	newPermIDs := []uuid.UUID{uuid.New(), uuid.New(), uuid.New()}
	role.SetPermissions(newPermIDs)

	assert.Len(t, role.PermissionIDs(), 3)
	for _, id := range newPermIDs {
		assert.True(t, role.HasPermission(id))
	}
	for _, id := range oldPermIDs {
		assert.False(t, role.HasPermission(id))
	}
	assert.Len(t, role.DomainEvents(), 1)
}

func TestRole_UpdateDetails(t *testing.T) {
	role, _ := NewRole(NewRoleParams{
		Name:        "editor",
		DisplayName: "Editor",
		Description: "Original description",
	})
	role.ClearDomainEvents()
	originalUpdatedAt := role.UpdatedAt()

	time.Sleep(time.Millisecond)
	err := role.UpdateDetails("New Display Name", "New description")
	require.NoError(t, err)
	assert.Equal(t, "New Display Name", role.DisplayName())
	assert.Equal(t, "New description", role.Description())
	assert.True(t, role.UpdatedAt().After(originalUpdatedAt))
	assert.Len(t, role.DomainEvents(), 1)
}

func TestRole_UpdateDetails_NoChange(t *testing.T) {
	role, _ := NewRole(NewRoleParams{
		Name:        "editor",
		DisplayName: "Editor",
		Description: "Description",
	})
	role.ClearDomainEvents()
	updatedAt := role.UpdatedAt()

	time.Sleep(time.Millisecond)
	err := role.UpdateDetails("", "Description")
	require.NoError(t, err)
	assert.Equal(t, updatedAt, role.UpdatedAt())
	assert.Empty(t, role.DomainEvents())
}

func TestRole_CanBeDeleted(t *testing.T) {
	systemRole, _ := NewRole(NewRoleParams{
		Name:        "system",
		DisplayName: "System",
		IsSystem:    true,
	})
	assert.False(t, systemRole.CanBeDeleted())

	defaultRole, _ := NewRole(NewRoleParams{
		Name:        "user",
		DisplayName: "User",
		IsDefault:   true,
	})
	assert.False(t, defaultRole.CanBeDeleted())

	regularRole, _ := NewRole(NewRoleParams{
		Name:        "custom",
		DisplayName: "Custom",
		IsSystem:    false,
		IsDefault:   false,
	})
	assert.True(t, regularRole.CanBeDeleted())
}

func TestRole_CanBeModified(t *testing.T) {
	systemRole, _ := NewRole(NewRoleParams{
		Name:        "system",
		DisplayName: "System",
		IsSystem:    true,
	})
	assert.False(t, systemRole.CanBeModified())

	regularRole, _ := NewRole(NewRoleParams{
		Name:        "custom",
		DisplayName: "Custom",
		IsSystem:    false,
	})
	assert.True(t, regularRole.CanBeModified())
}

func TestRole_PermissionIDs_ReturnsDefensiveCopy(t *testing.T) {
	permID := uuid.New()
	role, _ := NewRole(NewRoleParams{
		Name:          "editor",
		DisplayName:   "Editor",
		PermissionIDs: []uuid.UUID{permID},
	})

	permIDs := role.PermissionIDs()
	permIDs[0] = uuid.New()

	assert.Equal(t, permID, role.PermissionIDs()[0])
}
