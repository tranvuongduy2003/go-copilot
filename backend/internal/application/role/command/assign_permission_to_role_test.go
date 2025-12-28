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

func TestAssignPermissionToRoleHandler_Handle(t *testing.T) {
	ctx := context.Background()

	createTestRole := func(isSystem bool) *role.Role {
		now := time.Now().UTC()
		testRole, _ := role.ReconstructRole(role.ReconstructRoleParams{
			ID:          uuid.New(),
			Name:        "editor",
			DisplayName: "Editor",
			Description: "Test role",
			IsSystem:    isSystem,
			IsDefault:   false,
			Priority:    0,
			CreatedAt:   now,
			UpdatedAt:   now,
		})
		return testRole
	}

	createTestPermission := func() *permission.Permission {
		now := time.Now().UTC()
		perm, _ := permission.ReconstructPermission(permission.ReconstructPermissionParams{
			ID:          uuid.New(),
			Resource:    "articles",
			Action:      "create",
			Description: "Create articles",
			IsSystem:    false,
			CreatedAt:   now,
			UpdatedAt:   now,
		})
		return perm
	}

	tests := []struct {
		name        string
		setupMocks  func(*testutil.MockRoleRepository, *testutil.MockPermissionRepository) (*role.Role, *permission.Permission)
		command     func(*role.Role, *permission.Permission) AssignPermissionToRoleCommand
		wantErr     bool
		errContains string
		checkResult func(*testing.T, *role.Role, *permission.Permission, *testutil.MockRoleRepository)
	}{
		{
			name: "successfully assign permission to role",
			setupMocks: func(roleRepo *testutil.MockRoleRepository, permissionRepo *testutil.MockPermissionRepository) (*role.Role, *permission.Permission) {
				testRole := createTestRole(false)
				testPermission := createTestPermission()
				roleRepo.AddRole(testRole)
				permissionRepo.AddPermission(testPermission)
				return testRole, testPermission
			},
			command: func(r *role.Role, p *permission.Permission) AssignPermissionToRoleCommand {
				return AssignPermissionToRoleCommand{
					RoleID:       r.ID(),
					PermissionID: p.ID(),
				}
			},
			wantErr: false,
			checkResult: func(t *testing.T, r *role.Role, p *permission.Permission, roleRepo *testutil.MockRoleRepository) {
				updatedRole, _ := roleRepo.FindByID(context.Background(), r.ID())
				assert.True(t, updatedRole.HasPermission(p.ID()))
			},
		},
		{
			name: "fail when role is system role",
			setupMocks: func(roleRepo *testutil.MockRoleRepository, permissionRepo *testutil.MockPermissionRepository) (*role.Role, *permission.Permission) {
				systemRole := createTestRole(true)
				testPermission := createTestPermission()
				roleRepo.AddRole(systemRole)
				permissionRepo.AddPermission(testPermission)
				return systemRole, testPermission
			},
			command: func(r *role.Role, p *permission.Permission) AssignPermissionToRoleCommand {
				return AssignPermissionToRoleCommand{
					RoleID:       r.ID(),
					PermissionID: p.ID(),
				}
			},
			wantErr:     true,
			errContains: "system role",
		},
		{
			name: "fail when role not found",
			setupMocks: func(roleRepo *testutil.MockRoleRepository, permissionRepo *testutil.MockPermissionRepository) (*role.Role, *permission.Permission) {
				testRole := createTestRole(false)
				testPermission := createTestPermission()
				roleRepo.AddRole(testRole)
				permissionRepo.AddPermission(testPermission)
				return testRole, testPermission
			},
			command: func(r *role.Role, p *permission.Permission) AssignPermissionToRoleCommand {
				return AssignPermissionToRoleCommand{
					RoleID:       uuid.New(),
					PermissionID: p.ID(),
				}
			},
			wantErr:     true,
			errContains: "not found",
		},
		{
			name: "fail when permission not found",
			setupMocks: func(roleRepo *testutil.MockRoleRepository, permissionRepo *testutil.MockPermissionRepository) (*role.Role, *permission.Permission) {
				testRole := createTestRole(false)
				testPermission := createTestPermission()
				roleRepo.AddRole(testRole)
				permissionRepo.AddPermission(testPermission)
				return testRole, testPermission
			},
			command: func(r *role.Role, p *permission.Permission) AssignPermissionToRoleCommand {
				return AssignPermissionToRoleCommand{
					RoleID:       r.ID(),
					PermissionID: uuid.New(),
				}
			},
			wantErr:     true,
			errContains: "not found",
		},
		{
			name: "fail when permission already assigned",
			setupMocks: func(roleRepo *testutil.MockRoleRepository, permissionRepo *testutil.MockPermissionRepository) (*role.Role, *permission.Permission) {
				testPermission := createTestPermission()
				now := time.Now().UTC()
				testRole, _ := role.ReconstructRole(role.ReconstructRoleParams{
					ID:            uuid.New(),
					Name:          "editor",
					DisplayName:   "Editor",
					Description:   "Test role",
					PermissionIDs: []uuid.UUID{testPermission.ID()},
					IsSystem:      false,
					IsDefault:     false,
					Priority:      0,
					CreatedAt:     now,
					UpdatedAt:     now,
				})
				roleRepo.AddRole(testRole)
				permissionRepo.AddPermission(testPermission)
				return testRole, testPermission
			},
			command: func(r *role.Role, p *permission.Permission) AssignPermissionToRoleCommand {
				return AssignPermissionToRoleCommand{
					RoleID:       r.ID(),
					PermissionID: p.ID(),
				}
			},
			wantErr:     true,
			errContains: "already assigned",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roleRepo := testutil.NewMockRoleRepository()
			permissionRepo := testutil.NewMockPermissionRepository()
			eventBus := testutil.NewMockEventBus()
			logger := testutil.NewNoopLogger()

			testRole, testPermission := tt.setupMocks(roleRepo, permissionRepo)

			handler := NewAssignPermissionToRoleHandler(roleRepo, permissionRepo, eventBus, logger)
			cmd := tt.command(testRole, testPermission)

			result, err := handler.Handle(ctx, cmd)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				if tt.checkResult != nil {
					tt.checkResult(t, testRole, testPermission, roleRepo)
				}
			}
		})
	}
}
