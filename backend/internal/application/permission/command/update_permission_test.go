package permissioncommand

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/permission"
	"github.com/tranvuongduy2003/go-copilot/pkg/testutil"
)

func TestUpdatePermissionHandler_Handle(t *testing.T) {
	ctx := context.Background()

	createTestPermission := func(isSystem bool) *permission.Permission {
		now := time.Now().UTC()
		perm, _ := permission.ReconstructPermission(permission.ReconstructPermissionParams{
			ID:          uuid.New(),
			Resource:    "articles",
			Action:      "create",
			Description: "Original description",
			IsSystem:    isSystem,
			CreatedAt:   now,
			UpdatedAt:   now,
		})
		return perm
	}

	tests := []struct {
		name        string
		setupMocks  func(*testutil.MockPermissionRepository) *permission.Permission
		command     func(*permission.Permission) UpdatePermissionCommand
		wantErr     bool
		errContains string
		checkResult func(*testing.T, *permission.Permission, *testutil.MockPermissionRepository)
	}{
		{
			name: "successfully update permission description",
			setupMocks: func(permissionRepo *testutil.MockPermissionRepository) *permission.Permission {
				testPermission := createTestPermission(false)
				permissionRepo.AddPermission(testPermission)
				return testPermission
			},
			command: func(perm *permission.Permission) UpdatePermissionCommand {
				return UpdatePermissionCommand{
					PermissionID: perm.ID(),
					Description:  "Updated description",
				}
			},
			wantErr: false,
			checkResult: func(t *testing.T, perm *permission.Permission, permissionRepo *testutil.MockPermissionRepository) {
				updatedPermission, _ := permissionRepo.FindByID(context.Background(), perm.ID())
				assert.Equal(t, "Updated description", updatedPermission.Description())
			},
		},
		{
			name: "fail when updating system permission",
			setupMocks: func(permissionRepo *testutil.MockPermissionRepository) *permission.Permission {
				systemPermission := createTestPermission(true)
				permissionRepo.AddPermission(systemPermission)
				return systemPermission
			},
			command: func(perm *permission.Permission) UpdatePermissionCommand {
				return UpdatePermissionCommand{
					PermissionID: perm.ID(),
					Description:  "Trying to update system permission",
				}
			},
			wantErr:     true,
			errContains: "system permission",
		},
		{
			name: "fail when permission not found",
			setupMocks: func(permissionRepo *testutil.MockPermissionRepository) *permission.Permission {
				testPermission := createTestPermission(false)
				permissionRepo.AddPermission(testPermission)
				return testPermission
			},
			command: func(perm *permission.Permission) UpdatePermissionCommand {
				return UpdatePermissionCommand{
					PermissionID: uuid.New(),
					Description:  "Updating non-existent permission",
				}
			},
			wantErr:     true,
			errContains: "not found",
		},
		{
			name: "successfully update with same description (no-op)",
			setupMocks: func(permissionRepo *testutil.MockPermissionRepository) *permission.Permission {
				testPermission := createTestPermission(false)
				permissionRepo.AddPermission(testPermission)
				return testPermission
			},
			command: func(perm *permission.Permission) UpdatePermissionCommand {
				return UpdatePermissionCommand{
					PermissionID: perm.ID(),
					Description:  "Original description",
				}
			},
			wantErr: false,
			checkResult: func(t *testing.T, perm *permission.Permission, permissionRepo *testutil.MockPermissionRepository) {
				updatedPermission, _ := permissionRepo.FindByID(context.Background(), perm.ID())
				assert.Equal(t, "Original description", updatedPermission.Description())
			},
		},
		{
			name: "successfully update with empty description",
			setupMocks: func(permissionRepo *testutil.MockPermissionRepository) *permission.Permission {
				testPermission := createTestPermission(false)
				permissionRepo.AddPermission(testPermission)
				return testPermission
			},
			command: func(perm *permission.Permission) UpdatePermissionCommand {
				return UpdatePermissionCommand{
					PermissionID: perm.ID(),
					Description:  "",
				}
			},
			wantErr: false,
			checkResult: func(t *testing.T, perm *permission.Permission, permissionRepo *testutil.MockPermissionRepository) {
				updatedPermission, _ := permissionRepo.FindByID(context.Background(), perm.ID())
				assert.Equal(t, "", updatedPermission.Description())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			permissionRepo := testutil.NewMockPermissionRepository()
			logger := testutil.NewNoopLogger()

			testPermission := tt.setupMocks(permissionRepo)

			handler := NewUpdatePermissionHandler(permissionRepo, logger)
			cmd := tt.command(testPermission)

			result, err := handler.Handle(ctx, cmd)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, cmd.Description, result.Description)
				if tt.checkResult != nil {
					tt.checkResult(t, testPermission, permissionRepo)
				}
			}
		})
	}
}

func TestUpdatePermissionHandler_Handle_PreservesOriginalFields(t *testing.T) {
	ctx := context.Background()
	permissionRepo := testutil.NewMockPermissionRepository()
	logger := testutil.NewNoopLogger()

	now := time.Now().UTC()
	originalPermission, _ := permission.ReconstructPermission(permission.ReconstructPermissionParams{
		ID:          uuid.New(),
		Resource:    "users",
		Action:      "manage",
		Description: "Original description",
		IsSystem:    false,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	permissionRepo.AddPermission(originalPermission)

	handler := NewUpdatePermissionHandler(permissionRepo, logger)

	result, err := handler.Handle(ctx, UpdatePermissionCommand{
		PermissionID: originalPermission.ID(),
		Description:  "New description",
	})

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, originalPermission.ID(), result.ID)
	assert.Equal(t, "users", result.Resource)
	assert.Equal(t, "manage", result.Action)
	assert.Equal(t, "New description", result.Description)
	assert.False(t, result.IsSystem)
}
