package rolequery

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/permission"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/role"
	"github.com/tranvuongduy2003/go-copilot/pkg/testutil"
)

func TestGetRoleHandler_Handle(t *testing.T) {
	ctx := context.Background()

	createTestPermission := func(resource, action string) *permission.Permission {
		now := time.Now().UTC()
		perm, _ := permission.ReconstructPermission(permission.ReconstructPermissionParams{
			ID:          uuid.New(),
			Resource:    resource,
			Action:      action,
			Description: resource + ":" + action,
			IsSystem:    false,
			CreatedAt:   now,
			UpdatedAt:   now,
		})
		return perm
	}

	createTestRole := func(name string, permissionIDs []uuid.UUID) *role.Role {
		now := time.Now().UTC()
		testRole, _ := role.ReconstructRole(role.ReconstructRoleParams{
			ID:            uuid.New(),
			Name:          name,
			DisplayName:   name,
			Description:   "Test role",
			PermissionIDs: permissionIDs,
			IsSystem:      false,
			IsDefault:     false,
			Priority:      0,
			CreatedAt:     now,
			UpdatedAt:     now,
		})
		return testRole
	}

	tests := []struct {
		name        string
		setupMocks  func(*testutil.MockRoleRepository, *testutil.MockPermissionRepository) uuid.UUID
		wantErr     bool
		errContains string
	}{
		{
			name: "successfully get role with permissions",
			setupMocks: func(roleRepo *testutil.MockRoleRepository, permRepo *testutil.MockPermissionRepository) uuid.UUID {
				perm1 := createTestPermission("users", "read")
				perm2 := createTestPermission("users", "create")
				permRepo.AddPermission(perm1)
				permRepo.AddPermission(perm2)

				testRole := createTestRole("admin", []uuid.UUID{perm1.ID(), perm2.ID()})
				roleRepo.AddRole(testRole)
				return testRole.ID()
			},
			wantErr: false,
		},
		{
			name: "successfully get role without permissions",
			setupMocks: func(roleRepo *testutil.MockRoleRepository, permRepo *testutil.MockPermissionRepository) uuid.UUID {
				testRole := createTestRole("user", []uuid.UUID{})
				roleRepo.AddRole(testRole)
				return testRole.ID()
			},
			wantErr: false,
		},
		{
			name: "fail when role not found",
			setupMocks: func(roleRepo *testutil.MockRoleRepository, permRepo *testutil.MockPermissionRepository) uuid.UUID {
				return uuid.New()
			},
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roleRepo := testutil.NewMockRoleRepository()
			permRepo := testutil.NewMockPermissionRepository()
			logger := testutil.NewNoopLogger()

			roleID := tt.setupMocks(roleRepo, permRepo)

			handler := NewGetRoleHandler(roleRepo, permRepo, logger)

			result, err := handler.Handle(ctx, GetRoleQuery{RoleID: roleID})

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, roleID, result.ID)
			}
		})
	}
}

func TestGetRoleHandler_Handle_LoadsPermissionCodes(t *testing.T) {
	ctx := context.Background()
	roleRepo := testutil.NewMockRoleRepository()
	permRepo := testutil.NewMockPermissionRepository()
	logger := testutil.NewNoopLogger()

	now := time.Now().UTC()
	perm1, _ := permission.ReconstructPermission(permission.ReconstructPermissionParams{
		ID:          uuid.New(),
		Resource:    "users",
		Action:      "read",
		Description: "Read users",
		IsSystem:    false,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	perm2, _ := permission.ReconstructPermission(permission.ReconstructPermissionParams{
		ID:          uuid.New(),
		Resource:    "users",
		Action:      "create",
		Description: "Create users",
		IsSystem:    false,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	permRepo.AddPermission(perm1)
	permRepo.AddPermission(perm2)

	testRole, _ := role.ReconstructRole(role.ReconstructRoleParams{
		ID:            uuid.New(),
		Name:          "admin",
		DisplayName:   "Admin",
		Description:   "Administrator role",
		PermissionIDs: []uuid.UUID{perm1.ID(), perm2.ID()},
		IsSystem:      false,
		IsDefault:     false,
		Priority:      0,
		CreatedAt:     now,
		UpdatedAt:     now,
	})
	roleRepo.AddRole(testRole)

	handler := NewGetRoleHandler(roleRepo, permRepo, logger)

	result, err := handler.Handle(ctx, GetRoleQuery{RoleID: testRole.ID()})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result.Permissions, 2)
	assert.Contains(t, result.Permissions, "users:read")
	assert.Contains(t, result.Permissions, "users:create")
}
