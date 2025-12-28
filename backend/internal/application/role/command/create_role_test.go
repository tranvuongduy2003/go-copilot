package rolecommand

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/permission"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/role"
	"github.com/tranvuongduy2003/go-copilot/pkg/testutil"
)

func TestCreateRoleHandler_Handle(t *testing.T) {
	ctx := context.Background()

	createTestPermission := func() *permission.Permission {
		perm, _ := permission.NewPermission(permission.NewPermissionParams{
			Resource:    "articles",
			Action:      "create",
			Description: "Create articles",
			IsSystem:    false,
		})
		return perm
	}

	tests := []struct {
		name        string
		setupMocks  func(*testutil.MockRoleRepository, *testutil.MockPermissionRepository) []uuid.UUID
		command     func([]uuid.UUID) CreateRoleCommand
		wantErr     bool
		errContains string
		checkResult func(*testing.T, *testutil.MockRoleRepository)
	}{
		{
			name: "successfully create role without permissions",
			setupMocks: func(roleRepo *testutil.MockRoleRepository, permissionRepo *testutil.MockPermissionRepository) []uuid.UUID {
				return nil
			},
			command: func(permIDs []uuid.UUID) CreateRoleCommand {
				return CreateRoleCommand{
					Name:          "editor",
					DisplayName:   "Editor",
					Description:   "Can edit content",
					PermissionIDs: []uuid.UUID{},
				}
			},
			wantErr: false,
			checkResult: func(t *testing.T, roleRepo *testutil.MockRoleRepository) {
				assert.Len(t, roleRepo.Roles, 1)
				for _, r := range roleRepo.Roles {
					assert.Equal(t, "editor", r.Name())
					assert.Equal(t, "Editor", r.DisplayName())
					assert.Equal(t, "Can edit content", r.Description())
					assert.False(t, r.IsSystem())
					assert.False(t, r.IsDefault())
				}
			},
		},
		{
			name: "successfully create role with permissions",
			setupMocks: func(roleRepo *testutil.MockRoleRepository, permissionRepo *testutil.MockPermissionRepository) []uuid.UUID {
				perm1 := createTestPermission()
				perm2, _ := permission.NewPermission(permission.NewPermissionParams{
					Resource:    "articles",
					Action:      "read",
					Description: "Read articles",
					IsSystem:    false,
				})
				permissionRepo.AddPermission(perm1)
				permissionRepo.AddPermission(perm2)
				return []uuid.UUID{perm1.ID(), perm2.ID()}
			},
			command: func(permIDs []uuid.UUID) CreateRoleCommand {
				return CreateRoleCommand{
					Name:          "editor",
					DisplayName:   "Editor",
					Description:   "Can edit content",
					PermissionIDs: permIDs,
				}
			},
			wantErr: false,
			checkResult: func(t *testing.T, roleRepo *testutil.MockRoleRepository) {
				assert.Len(t, roleRepo.Roles, 1)
				for _, r := range roleRepo.Roles {
					assert.Len(t, r.PermissionIDs(), 2)
				}
			},
		},
		{
			name: "fail when role name already exists",
			setupMocks: func(roleRepo *testutil.MockRoleRepository, permissionRepo *testutil.MockPermissionRepository) []uuid.UUID {
				existingRole, _ := role.NewRole(role.NewRoleParams{
					Name:        "editor",
					DisplayName: "Existing Editor",
				})
				roleRepo.AddRole(existingRole)
				return nil
			},
			command: func(permIDs []uuid.UUID) CreateRoleCommand {
				return CreateRoleCommand{
					Name:          "editor",
					DisplayName:   "New Editor",
					Description:   "New description",
					PermissionIDs: []uuid.UUID{},
				}
			},
			wantErr:     true,
			errContains: "already exists",
		},
		{
			name: "fail when permission not found",
			setupMocks: func(roleRepo *testutil.MockRoleRepository, permissionRepo *testutil.MockPermissionRepository) []uuid.UUID {
				return []uuid.UUID{uuid.New()}
			},
			command: func(permIDs []uuid.UUID) CreateRoleCommand {
				return CreateRoleCommand{
					Name:          "editor",
					DisplayName:   "Editor",
					Description:   "Can edit content",
					PermissionIDs: permIDs,
				}
			},
			wantErr:     true,
			errContains: "not found",
		},
		{
			name: "fail when role name is empty",
			setupMocks: func(roleRepo *testutil.MockRoleRepository, permissionRepo *testutil.MockPermissionRepository) []uuid.UUID {
				return nil
			},
			command: func(permIDs []uuid.UUID) CreateRoleCommand {
				return CreateRoleCommand{
					Name:          "",
					DisplayName:   "Editor",
					Description:   "Can edit content",
					PermissionIDs: []uuid.UUID{},
				}
			},
			wantErr:     true,
			errContains: "name",
		},
		{
			name: "fail when display name is empty",
			setupMocks: func(roleRepo *testutil.MockRoleRepository, permissionRepo *testutil.MockPermissionRepository) []uuid.UUID {
				return nil
			},
			command: func(permIDs []uuid.UUID) CreateRoleCommand {
				return CreateRoleCommand{
					Name:          "editor",
					DisplayName:   "",
					Description:   "Can edit content",
					PermissionIDs: []uuid.UUID{},
				}
			},
			wantErr:     true,
			errContains: "display_name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roleRepo := testutil.NewMockRoleRepository()
			permissionRepo := testutil.NewMockPermissionRepository()
			eventBus := testutil.NewMockEventBus()
			logger := testutil.NewNoopLogger()

			permIDs := tt.setupMocks(roleRepo, permissionRepo)

			handler := NewCreateRoleHandler(roleRepo, permissionRepo, eventBus, logger)
			cmd := tt.command(permIDs)

			result, err := handler.Handle(ctx, cmd)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, cmd.Name, result.Name)
				assert.Equal(t, cmd.DisplayName, result.DisplayName)
				if tt.checkResult != nil {
					tt.checkResult(t, roleRepo)
				}
			}
		})
	}
}
