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

func TestCreatePermissionHandler_Handle(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		setupMocks  func(*testutil.MockPermissionRepository)
		command     CreatePermissionCommand
		wantErr     bool
		errContains string
		checkResult func(*testing.T, *testutil.MockPermissionRepository)
	}{
		{
			name:       "successfully create permission",
			setupMocks: func(permissionRepo *testutil.MockPermissionRepository) {},
			command: CreatePermissionCommand{
				Resource:    "articles",
				Action:      "create",
				Description: "Create articles",
			},
			wantErr: false,
			checkResult: func(t *testing.T, permissionRepo *testutil.MockPermissionRepository) {
				assert.Len(t, permissionRepo.Permissions, 1)
				for _, perm := range permissionRepo.Permissions {
					assert.Equal(t, "articles", perm.Resource().String())
					assert.Equal(t, "create", perm.Action().String())
					assert.Equal(t, "Create articles", perm.Description())
					assert.False(t, perm.IsSystem())
				}
			},
		},
		{
			name: "fail when permission code already exists",
			setupMocks: func(permissionRepo *testutil.MockPermissionRepository) {
				existingPermission, _ := permission.NewPermission(permission.NewPermissionParams{
					Resource:    "articles",
					Action:      "create",
					Description: "Existing permission",
					IsSystem:    false,
				})
				permissionRepo.AddPermission(existingPermission)
			},
			command: CreatePermissionCommand{
				Resource:    "articles",
				Action:      "create",
				Description: "New permission with same code",
			},
			wantErr:     true,
			errContains: "already exists",
		},
		{
			name:       "fail when resource is invalid",
			setupMocks: func(permissionRepo *testutil.MockPermissionRepository) {},
			command: CreatePermissionCommand{
				Resource:    "",
				Action:      "create",
				Description: "Invalid resource",
			},
			wantErr:     true,
			errContains: "resource",
		},
		{
			name:       "fail when action is invalid",
			setupMocks: func(permissionRepo *testutil.MockPermissionRepository) {},
			command: CreatePermissionCommand{
				Resource:    "articles",
				Action:      "",
				Description: "Invalid action",
			},
			wantErr:     true,
			errContains: "action",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			permissionRepo := testutil.NewMockPermissionRepository()
			logger := testutil.NewNoopLogger()

			tt.setupMocks(permissionRepo)

			handler := NewCreatePermissionHandler(permissionRepo, logger)

			result, err := handler.Handle(ctx, tt.command)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tt.command.Resource, result.Resource)
				assert.Equal(t, tt.command.Action, result.Action)
				assert.Equal(t, tt.command.Description, result.Description)
				if tt.checkResult != nil {
					tt.checkResult(t, permissionRepo)
				}
			}
		})
	}
}

func TestCreatePermissionHandler_Handle_WithDifferentResources(t *testing.T) {
	ctx := context.Background()

	resources := []string{"users", "roles", "permissions", "orders", "products"}
	actions := []string{"create", "read", "update", "delete", "list", "manage"}

	permissionRepo := testutil.NewMockPermissionRepository()
	logger := testutil.NewNoopLogger()
	handler := NewCreatePermissionHandler(permissionRepo, logger)

	for _, resource := range resources {
		for _, action := range actions {
			t.Run(resource+":"+action, func(t *testing.T) {
				command := CreatePermissionCommand{
					Resource:    resource,
					Action:      action,
					Description: "Test permission for " + resource + ":" + action,
				}

				result, err := handler.Handle(ctx, command)

				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, resource, result.Resource)
				assert.Equal(t, action, result.Action)
			})
		}
	}
}

func TestCreatePermissionHandler_Handle_GeneratesUniqueID(t *testing.T) {
	ctx := context.Background()
	permissionRepo := testutil.NewMockPermissionRepository()
	logger := testutil.NewNoopLogger()
	handler := NewCreatePermissionHandler(permissionRepo, logger)

	result1, err := handler.Handle(ctx, CreatePermissionCommand{
		Resource:    "articles",
		Action:      "create",
		Description: "First permission",
	})
	require.NoError(t, err)

	result2, err := handler.Handle(ctx, CreatePermissionCommand{
		Resource:    "articles",
		Action:      "read",
		Description: "Second permission",
	})
	require.NoError(t, err)

	assert.NotEqual(t, result1.ID, result2.ID)
	_, err = uuid.Parse(result1.ID.String())
	assert.NoError(t, err)
	_, err = uuid.Parse(result2.ID.String())
	assert.NoError(t, err)
}

func TestCreatePermissionHandler_Handle_SetsCorrectTimestamps(t *testing.T) {
	ctx := context.Background()
	permissionRepo := testutil.NewMockPermissionRepository()
	logger := testutil.NewNoopLogger()
	handler := NewCreatePermissionHandler(permissionRepo, logger)

	beforeCreate := time.Now().UTC()

	result, err := handler.Handle(ctx, CreatePermissionCommand{
		Resource:    "articles",
		Action:      "create",
		Description: "Test permission",
	})

	afterCreate := time.Now().UTC()

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.True(t, result.CreatedAt.After(beforeCreate) || result.CreatedAt.Equal(beforeCreate))
	assert.True(t, result.CreatedAt.Before(afterCreate) || result.CreatedAt.Equal(afterCreate))
	assert.True(t, result.UpdatedAt.After(beforeCreate) || result.UpdatedAt.Equal(beforeCreate))
	assert.True(t, result.UpdatedAt.Before(afterCreate) || result.UpdatedAt.Equal(afterCreate))
}
