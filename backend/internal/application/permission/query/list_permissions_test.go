package permissionquery

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

func TestListPermissionsHandler_Handle(t *testing.T) {
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

	tests := []struct {
		name          string
		setupMocks    func(*testutil.MockPermissionRepository)
		query         ListPermissionsQuery
		wantErr       bool
		errContains   string
		expectedCount int
	}{
		{
			name: "successfully list all permissions",
			setupMocks: func(permRepo *testutil.MockPermissionRepository) {
				permRepo.AddPermission(createTestPermission("users", "read"))
				permRepo.AddPermission(createTestPermission("users", "create"))
				permRepo.AddPermission(createTestPermission("roles", "read"))
			},
			query:         ListPermissionsQuery{},
			wantErr:       false,
			expectedCount: 3,
		},
		{
			name: "successfully filter by resource",
			setupMocks: func(permRepo *testutil.MockPermissionRepository) {
				permRepo.AddPermission(createTestPermission("users", "read"))
				permRepo.AddPermission(createTestPermission("users", "create"))
				permRepo.AddPermission(createTestPermission("roles", "read"))
			},
			query: ListPermissionsQuery{
				Resource: stringPtr("users"),
			},
			wantErr:       false,
			expectedCount: 2,
		},
		{
			name: "return empty list when no permissions",
			setupMocks: func(permRepo *testutil.MockPermissionRepository) {
			},
			query:         ListPermissionsQuery{},
			wantErr:       false,
			expectedCount: 0,
		},
		{
			name: "return empty list when resource has no permissions",
			setupMocks: func(permRepo *testutil.MockPermissionRepository) {
				permRepo.AddPermission(createTestPermission("users", "read"))
			},
			query: ListPermissionsQuery{
				Resource: stringPtr("roles"),
			},
			wantErr:       false,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			permRepo := testutil.NewMockPermissionRepository()
			logger := testutil.NewNoopLogger()

			tt.setupMocks(permRepo)

			handler := NewListPermissionsHandler(permRepo, logger)

			result, err := handler.Handle(ctx, tt.query)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				require.NoError(t, err)
				assert.Len(t, result, tt.expectedCount)
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
