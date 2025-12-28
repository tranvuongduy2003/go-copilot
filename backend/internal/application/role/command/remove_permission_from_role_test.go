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

func TestRemovePermissionFromRoleHandler_Handle(t *testing.T) {
	ctx := context.Background()

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

	createTestRoleWithPermission := func(permissionID uuid.UUID, isSystem bool) *role.Role {
		now := time.Now().UTC()
		testRole, _ := role.ReconstructRole(role.ReconstructRoleParams{
			ID:            uuid.New(),
			Name:          "editor",
			DisplayName:   "Editor",
			Description:   "Test role",
			PermissionIDs: []uuid.UUID{permissionID},
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
		setupMocks  func(*testutil.MockRoleRepository) (*role.Role, *permission.Permission)
		command     func(*role.Role, *permission.Permission) RemovePermissionFromRoleCommand
		wantErr     bool
		errContains string
		checkResult func(*testing.T, *role.Role, *permission.Permission, *testutil.MockRoleRepository)
	}{
		{
			name: "successfully remove permission from role",
			setupMocks: func(roleRepo *testutil.MockRoleRepository) (*role.Role, *permission.Permission) {
				testPermission := createTestPermission()
				testRole := createTestRoleWithPermission(testPermission.ID(), false)
				roleRepo.AddRole(testRole)
				return testRole, testPermission
			},
			command: func(r *role.Role, p *permission.Permission) RemovePermissionFromRoleCommand {
				return RemovePermissionFromRoleCommand{
					RoleID:       r.ID(),
					PermissionID: p.ID(),
				}
			},
			wantErr: false,
			checkResult: func(t *testing.T, r *role.Role, p *permission.Permission, roleRepo *testutil.MockRoleRepository) {
				updatedRole, _ := roleRepo.FindByID(context.Background(), r.ID())
				assert.False(t, updatedRole.HasPermission(p.ID()))
			},
		},
		{
			name: "fail when role is system role",
			setupMocks: func(roleRepo *testutil.MockRoleRepository) (*role.Role, *permission.Permission) {
				testPermission := createTestPermission()
				systemRole := createTestRoleWithPermission(testPermission.ID(), true)
				roleRepo.AddRole(systemRole)
				return systemRole, testPermission
			},
			command: func(r *role.Role, p *permission.Permission) RemovePermissionFromRoleCommand {
				return RemovePermissionFromRoleCommand{
					RoleID:       r.ID(),
					PermissionID: p.ID(),
				}
			},
			wantErr:     true,
			errContains: "system role",
		},
		{
			name: "fail when role not found",
			setupMocks: func(roleRepo *testutil.MockRoleRepository) (*role.Role, *permission.Permission) {
				testPermission := createTestPermission()
				testRole := createTestRoleWithPermission(testPermission.ID(), false)
				roleRepo.AddRole(testRole)
				return testRole, testPermission
			},
			command: func(r *role.Role, p *permission.Permission) RemovePermissionFromRoleCommand {
				return RemovePermissionFromRoleCommand{
					RoleID:       uuid.New(),
					PermissionID: p.ID(),
				}
			},
			wantErr:     true,
			errContains: "not found",
		},
		{
			name: "fail when permission not assigned to role",
			setupMocks: func(roleRepo *testutil.MockRoleRepository) (*role.Role, *permission.Permission) {
				testPermission := createTestPermission()
				testRole := createTestRoleWithPermission(testPermission.ID(), false)
				roleRepo.AddRole(testRole)
				return testRole, testPermission
			},
			command: func(r *role.Role, p *permission.Permission) RemovePermissionFromRoleCommand {
				return RemovePermissionFromRoleCommand{
					RoleID:       r.ID(),
					PermissionID: uuid.New(),
				}
			},
			wantErr:     true,
			errContains: "not assigned",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roleRepo := testutil.NewMockRoleRepository()
			eventBus := testutil.NewMockEventBus()
			logger := testutil.NewNoopLogger()

			testRole, testPermission := tt.setupMocks(roleRepo)

			handler := NewRemovePermissionFromRoleHandler(roleRepo, eventBus, logger)
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

func TestRemovePermissionFromRoleHandler_Handle_RemoveMultiplePermissions(t *testing.T) {
	ctx := context.Background()
	roleRepo := testutil.NewMockRoleRepository()
	eventBus := testutil.NewMockEventBus()
	logger := testutil.NewNoopLogger()

	perm1ID := uuid.New()
	perm2ID := uuid.New()
	perm3ID := uuid.New()

	now := time.Now().UTC()
	testRole, _ := role.ReconstructRole(role.ReconstructRoleParams{
		ID:            uuid.New(),
		Name:          "editor",
		DisplayName:   "Editor",
		Description:   "Test role",
		PermissionIDs: []uuid.UUID{perm1ID, perm2ID, perm3ID},
		IsSystem:      false,
		IsDefault:     false,
		Priority:      0,
		CreatedAt:     now,
		UpdatedAt:     now,
	})
	roleRepo.AddRole(testRole)

	handler := NewRemovePermissionFromRoleHandler(roleRepo, eventBus, logger)

	_, err := handler.Handle(ctx, RemovePermissionFromRoleCommand{
		RoleID:       testRole.ID(),
		PermissionID: perm1ID,
	})
	require.NoError(t, err)

	_, err = handler.Handle(ctx, RemovePermissionFromRoleCommand{
		RoleID:       testRole.ID(),
		PermissionID: perm2ID,
	})
	require.NoError(t, err)

	updatedRole, _ := roleRepo.FindByID(ctx, testRole.ID())
	assert.False(t, updatedRole.HasPermission(perm1ID))
	assert.False(t, updatedRole.HasPermission(perm2ID))
	assert.True(t, updatedRole.HasPermission(perm3ID))
	assert.Len(t, updatedRole.PermissionIDs(), 1)
}
