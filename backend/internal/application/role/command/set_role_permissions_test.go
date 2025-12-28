package rolecommand

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

func TestSetRolePermissionsHandler_Handle(t *testing.T) {
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

	createTestRole := func(permissionIDs []uuid.UUID, isSystem bool) *role.Role {
		now := time.Now().UTC()
		testRole, _ := role.ReconstructRole(role.ReconstructRoleParams{
			ID:            uuid.New(),
			Name:          "editor",
			DisplayName:   "Editor",
			Description:   "Test role",
			PermissionIDs: permissionIDs,
			IsSystem:      isSystem,
			IsDefault:     false,
			Priority:      0,
			CreatedAt:     now,
			UpdatedAt:     now,
		})
		return testRole
	}

	tests := []struct {
		name        string
		setupMocks  func(*testutil.MockRoleRepository, *testutil.MockPermissionRepository) (*role.Role, []uuid.UUID)
		command     func(*role.Role, []uuid.UUID) SetRolePermissionsCommand
		wantErr     bool
		errContains string
		checkResult func(*testing.T, *role.Role, []uuid.UUID, *testutil.MockRoleRepository)
	}{
		{
			name: "successfully set permissions for role",
			setupMocks: func(roleRepo *testutil.MockRoleRepository, permissionRepo *testutil.MockPermissionRepository) (*role.Role, []uuid.UUID) {
				perm1 := createTestPermission("articles", "create")
				perm2 := createTestPermission("articles", "read")
				permissionRepo.AddPermission(perm1)
				permissionRepo.AddPermission(perm2)
				testRole := createTestRole([]uuid.UUID{}, false)
				roleRepo.AddRole(testRole)
				return testRole, []uuid.UUID{perm1.ID(), perm2.ID()}
			},
			command: func(r *role.Role, permIDs []uuid.UUID) SetRolePermissionsCommand {
				return SetRolePermissionsCommand{
					RoleID:        r.ID(),
					PermissionIDs: permIDs,
				}
			},
			wantErr: false,
			checkResult: func(t *testing.T, r *role.Role, permIDs []uuid.UUID, roleRepo *testutil.MockRoleRepository) {
				updatedRole, _ := roleRepo.FindByID(context.Background(), r.ID())
				assert.Len(t, updatedRole.PermissionIDs(), 2)
				for _, permID := range permIDs {
					assert.True(t, updatedRole.HasPermission(permID))
				}
			},
		},
		{
			name: "successfully replace existing permissions",
			setupMocks: func(roleRepo *testutil.MockRoleRepository, permissionRepo *testutil.MockPermissionRepository) (*role.Role, []uuid.UUID) {
				oldPerm := createTestPermission("articles", "delete")
				newPerm1 := createTestPermission("articles", "create")
				newPerm2 := createTestPermission("articles", "read")
				permissionRepo.AddPermission(oldPerm)
				permissionRepo.AddPermission(newPerm1)
				permissionRepo.AddPermission(newPerm2)
				testRole := createTestRole([]uuid.UUID{oldPerm.ID()}, false)
				roleRepo.AddRole(testRole)
				return testRole, []uuid.UUID{newPerm1.ID(), newPerm2.ID()}
			},
			command: func(r *role.Role, permIDs []uuid.UUID) SetRolePermissionsCommand {
				return SetRolePermissionsCommand{
					RoleID:        r.ID(),
					PermissionIDs: permIDs,
				}
			},
			wantErr: false,
			checkResult: func(t *testing.T, r *role.Role, permIDs []uuid.UUID, roleRepo *testutil.MockRoleRepository) {
				updatedRole, _ := roleRepo.FindByID(context.Background(), r.ID())
				assert.Len(t, updatedRole.PermissionIDs(), 2)
				for _, permID := range permIDs {
					assert.True(t, updatedRole.HasPermission(permID))
				}
			},
		},
		{
			name: "successfully set empty permissions",
			setupMocks: func(roleRepo *testutil.MockRoleRepository, permissionRepo *testutil.MockPermissionRepository) (*role.Role, []uuid.UUID) {
				oldPerm := createTestPermission("articles", "delete")
				permissionRepo.AddPermission(oldPerm)
				testRole := createTestRole([]uuid.UUID{oldPerm.ID()}, false)
				roleRepo.AddRole(testRole)
				return testRole, []uuid.UUID{}
			},
			command: func(r *role.Role, permIDs []uuid.UUID) SetRolePermissionsCommand {
				return SetRolePermissionsCommand{
					RoleID:        r.ID(),
					PermissionIDs: []uuid.UUID{},
				}
			},
			wantErr: false,
			checkResult: func(t *testing.T, r *role.Role, permIDs []uuid.UUID, roleRepo *testutil.MockRoleRepository) {
				updatedRole, _ := roleRepo.FindByID(context.Background(), r.ID())
				assert.Len(t, updatedRole.PermissionIDs(), 0)
			},
		},
		{
			name: "fail when role is system role",
			setupMocks: func(roleRepo *testutil.MockRoleRepository, permissionRepo *testutil.MockPermissionRepository) (*role.Role, []uuid.UUID) {
				perm := createTestPermission("articles", "create")
				permissionRepo.AddPermission(perm)
				systemRole := createTestRole([]uuid.UUID{}, true)
				roleRepo.AddRole(systemRole)
				return systemRole, []uuid.UUID{perm.ID()}
			},
			command: func(r *role.Role, permIDs []uuid.UUID) SetRolePermissionsCommand {
				return SetRolePermissionsCommand{
					RoleID:        r.ID(),
					PermissionIDs: permIDs,
				}
			},
			wantErr:     true,
			errContains: "system role",
		},
		{
			name: "fail when role not found",
			setupMocks: func(roleRepo *testutil.MockRoleRepository, permissionRepo *testutil.MockPermissionRepository) (*role.Role, []uuid.UUID) {
				perm := createTestPermission("articles", "create")
				permissionRepo.AddPermission(perm)
				testRole := createTestRole([]uuid.UUID{}, false)
				roleRepo.AddRole(testRole)
				return testRole, []uuid.UUID{perm.ID()}
			},
			command: func(r *role.Role, permIDs []uuid.UUID) SetRolePermissionsCommand {
				return SetRolePermissionsCommand{
					RoleID:        uuid.New(),
					PermissionIDs: permIDs,
				}
			},
			wantErr:     true,
			errContains: "not found",
		},
		{
			name: "fail when some permissions not found",
			setupMocks: func(roleRepo *testutil.MockRoleRepository, permissionRepo *testutil.MockPermissionRepository) (*role.Role, []uuid.UUID) {
				perm := createTestPermission("articles", "create")
				permissionRepo.AddPermission(perm)
				testRole := createTestRole([]uuid.UUID{}, false)
				roleRepo.AddRole(testRole)
				return testRole, []uuid.UUID{perm.ID(), uuid.New()}
			},
			command: func(r *role.Role, permIDs []uuid.UUID) SetRolePermissionsCommand {
				return SetRolePermissionsCommand{
					RoleID:        r.ID(),
					PermissionIDs: permIDs,
				}
			},
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roleRepo := testutil.NewMockRoleRepository()
			permissionRepo := testutil.NewMockPermissionRepository()
			eventBus := testutil.NewMockEventBus()
			logger := testutil.NewNoopLogger()

			testRole, permIDs := tt.setupMocks(roleRepo, permissionRepo)

			handler := NewSetRolePermissionsHandler(roleRepo, permissionRepo, eventBus, logger)
			cmd := tt.command(testRole, permIDs)

			result, err := handler.Handle(ctx, cmd)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				if tt.checkResult != nil {
					tt.checkResult(t, testRole, permIDs, roleRepo)
				}
			}
		})
	}
}
