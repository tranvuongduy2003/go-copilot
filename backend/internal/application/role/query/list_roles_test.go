package rolequery

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/role"
	"github.com/tranvuongduy2003/go-copilot/pkg/testutil"
)

func TestListRolesHandler_Handle(t *testing.T) {
	ctx := context.Background()

	createTestRole := func(name string) *role.Role {
		now := time.Now().UTC()
		testRole, _ := role.ReconstructRole(role.ReconstructRoleParams{
			ID:          uuid.New(),
			Name:        name,
			DisplayName: name,
			Description: "Test role",
			IsSystem:    false,
			IsDefault:   false,
			Priority:    0,
			CreatedAt:   now,
			UpdatedAt:   now,
		})
		return testRole
	}

	tests := []struct {
		name          string
		setupMocks    func(*testutil.MockRoleRepository)
		wantErr       bool
		errContains   string
		expectedCount int
	}{
		{
			name: "successfully list all roles",
			setupMocks: func(roleRepo *testutil.MockRoleRepository) {
				roleRepo.AddRole(createTestRole("admin"))
				roleRepo.AddRole(createTestRole("user"))
				roleRepo.AddRole(createTestRole("manager"))
			},
			wantErr:       false,
			expectedCount: 3,
		},
		{
			name: "return empty list when no roles",
			setupMocks: func(roleRepo *testutil.MockRoleRepository) {
			},
			wantErr:       false,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roleRepo := testutil.NewMockRoleRepository()
			logger := testutil.NewNoopLogger()

			tt.setupMocks(roleRepo)

			handler := NewListRolesHandler(roleRepo, logger)

			result, err := handler.Handle(ctx, ListRolesQuery{})

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
