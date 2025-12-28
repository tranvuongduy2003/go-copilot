package permissioncommand

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

func TestDeletePermissionHandler_Handle(t *testing.T) {
	ctx := context.Background()

	createTestPermission := func(isSystem bool) *permission.Permission {
		now := time.Now().UTC()
		perm, _ := permission.ReconstructPermission(permission.ReconstructPermissionParams{
			ID:          uuid.New(),
			Resource:    "articles",
			Action:      "delete",
			Description: "Delete articles",
			IsSystem:    isSystem,
			CreatedAt:   now,
			UpdatedAt:   now,
		})
		return perm
	}

	createTestRole := func(permissionID uuid.UUID) *role.Role {
		testRole, _ := role.NewRole(role.NewRoleParams{
			Name:        "editor",
			DisplayName: "Editor",
		})
		testRole.AddPermission(permissionID)
		return testRole
	}

	tests := []struct {
		name        string
		setupMocks  func(*testutil.MockPermissionRepository, *testutil.MockRoleRepository) *permission.Permission
		command     func(*permission.Permission) DeletePermissionCommand
		wantErr     bool
		errContains string
		checkResult func(*testing.T, *permission.Permission, *testutil.MockPermissionRepository)
	}{
		{
			name: "successfully delete permission",
			setupMocks: func(permissionRepo *testutil.MockPermissionRepository, roleRepo *testutil.MockRoleRepository) *permission.Permission {
				testPermission := createTestPermission(false)
				permissionRepo.AddPermission(testPermission)
				return testPermission
			},
			command: func(perm *permission.Permission) DeletePermissionCommand {
				return DeletePermissionCommand{
					PermissionID: perm.ID(),
				}
			},
			wantErr: false,
			checkResult: func(t *testing.T, perm *permission.Permission, permissionRepo *testutil.MockPermissionRepository) {
				_, err := permissionRepo.FindByID(context.Background(), perm.ID())
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "not found")
			},
		},
		{
			name: "fail when deleting system permission",
			setupMocks: func(permissionRepo *testutil.MockPermissionRepository, roleRepo *testutil.MockRoleRepository) *permission.Permission {
				systemPermission := createTestPermission(true)
				permissionRepo.AddPermission(systemPermission)
				return systemPermission
			},
			command: func(perm *permission.Permission) DeletePermissionCommand {
				return DeletePermissionCommand{
					PermissionID: perm.ID(),
				}
			},
			wantErr:     true,
			errContains: "system permission",
		},
		{
			name: "fail when permission not found",
			setupMocks: func(permissionRepo *testutil.MockPermissionRepository, roleRepo *testutil.MockRoleRepository) *permission.Permission {
				testPermission := createTestPermission(false)
				permissionRepo.AddPermission(testPermission)
				return testPermission
			},
			command: func(perm *permission.Permission) DeletePermissionCommand {
				return DeletePermissionCommand{
					PermissionID: uuid.New(),
				}
			},
			wantErr:     true,
			errContains: "not found",
		},
		{
			name: "fail when permission is assigned to a role",
			setupMocks: func(permissionRepo *testutil.MockPermissionRepository, roleRepo *testutil.MockRoleRepository) *permission.Permission {
				testPermission := createTestPermission(false)
				permissionRepo.AddPermission(testPermission)

				testRole := createTestRole(testPermission.ID())
				roleRepo.AddRole(testRole)

				return testPermission
			},
			command: func(perm *permission.Permission) DeletePermissionCommand {
				return DeletePermissionCommand{
					PermissionID: perm.ID(),
				}
			},
			wantErr:     true,
			errContains: "assigned to one or more roles",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			permissionRepo := testutil.NewMockPermissionRepository()
			roleRepo := testutil.NewMockRoleRepository()
			logger := testutil.NewNoopLogger()

			testPermission := tt.setupMocks(permissionRepo, roleRepo)

			handler := NewDeletePermissionHandler(permissionRepo, roleRepo, logger)
			cmd := tt.command(testPermission)

			err := handler.Handle(ctx, cmd)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				require.NoError(t, err)
				if tt.checkResult != nil {
					tt.checkResult(t, testPermission, permissionRepo)
				}
			}
		})
	}
}

func TestDeletePermissionHandler_Handle_MultipleRolesWithPermission(t *testing.T) {
	ctx := context.Background()
	permissionRepo := testutil.NewMockPermissionRepository()
	roleRepo := testutil.NewMockRoleRepository()
	logger := testutil.NewNoopLogger()

	now := time.Now().UTC()
	testPermission, _ := permission.ReconstructPermission(permission.ReconstructPermissionParams{
		ID:          uuid.New(),
		Resource:    "articles",
		Action:      "delete",
		Description: "Delete articles",
		IsSystem:    false,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	permissionRepo.AddPermission(testPermission)

	role1, _ := role.NewRole(role.NewRoleParams{
		Name:        "editor",
		DisplayName: "Editor",
	})
	role1.AddPermission(testPermission.ID())
	roleRepo.AddRole(role1)

	role2, _ := role.NewRole(role.NewRoleParams{
		Name:        "admin",
		DisplayName: "Administrator",
	})
	role2.AddPermission(testPermission.ID())
	roleRepo.AddRole(role2)

	handler := NewDeletePermissionHandler(permissionRepo, roleRepo, logger)

	err := handler.Handle(ctx, DeletePermissionCommand{
		PermissionID: testPermission.ID(),
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "assigned to one or more roles")

	_, findErr := permissionRepo.FindByID(ctx, testPermission.ID())
	require.NoError(t, findErr)
}

func TestDeletePermissionHandler_Handle_PermissionNotAssignedToAnyRole(t *testing.T) {
	ctx := context.Background()
	permissionRepo := testutil.NewMockPermissionRepository()
	roleRepo := testutil.NewMockRoleRepository()
	logger := testutil.NewNoopLogger()

	now := time.Now().UTC()
	testPermission, _ := permission.ReconstructPermission(permission.ReconstructPermissionParams{
		ID:          uuid.New(),
		Resource:    "articles",
		Action:      "archive",
		Description: "Archive articles",
		IsSystem:    false,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	permissionRepo.AddPermission(testPermission)

	role1, _ := role.NewRole(role.NewRoleParams{
		Name:        "editor",
		DisplayName: "Editor",
	})
	roleRepo.AddRole(role1)

	handler := NewDeletePermissionHandler(permissionRepo, roleRepo, logger)

	err := handler.Handle(ctx, DeletePermissionCommand{
		PermissionID: testPermission.ID(),
	})

	require.NoError(t, err)

	_, findErr := permissionRepo.FindByID(ctx, testPermission.ID())
	require.Error(t, findErr)
	assert.Contains(t, findErr.Error(), "not found")
}
