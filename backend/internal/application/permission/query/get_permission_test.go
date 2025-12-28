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

func TestGetPermissionHandler_Handle(t *testing.T) {
	ctx := context.Background()

	createTestPermission := func() *permission.Permission {
		now := time.Now().UTC()
		perm, _ := permission.ReconstructPermission(permission.ReconstructPermissionParams{
			ID:          uuid.New(),
			Resource:    "users",
			Action:      "read",
			Description: "Read users",
			IsSystem:    false,
			CreatedAt:   now,
			UpdatedAt:   now,
		})
		return perm
	}

	tests := []struct {
		name        string
		setupMocks  func(*testutil.MockPermissionRepository) uuid.UUID
		wantErr     bool
		errContains string
	}{
		{
			name: "successfully get permission",
			setupMocks: func(permRepo *testutil.MockPermissionRepository) uuid.UUID {
				perm := createTestPermission()
				permRepo.AddPermission(perm)
				return perm.ID()
			},
			wantErr: false,
		},
		{
			name: "fail when permission not found",
			setupMocks: func(permRepo *testutil.MockPermissionRepository) uuid.UUID {
				return uuid.New()
			},
			wantErr:     true,
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			permRepo := testutil.NewMockPermissionRepository()
			logger := testutil.NewNoopLogger()

			permissionID := tt.setupMocks(permRepo)

			handler := NewGetPermissionHandler(permRepo, logger)

			result, err := handler.Handle(ctx, GetPermissionQuery{PermissionID: permissionID})

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, permissionID, result.ID)
			}
		})
	}
}
