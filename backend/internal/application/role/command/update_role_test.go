package rolecommand

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

func TestUpdateRoleHandler_Handle(t *testing.T) {
	ctx := context.Background()

	createTestRole := func(isSystem bool) *role.Role {
		now := time.Now().UTC()
		testRole, _ := role.ReconstructRole(role.ReconstructRoleParams{
			ID:          uuid.New(),
			Name:        "editor",
			DisplayName: "Editor",
			Description: "Original description",
			IsSystem:    isSystem,
			IsDefault:   false,
			Priority:    0,
			CreatedAt:   now,
			UpdatedAt:   now,
		})
		return testRole
	}

	tests := []struct {
		name        string
		setupMocks  func(*testutil.MockRoleRepository) *role.Role
		command     func(*role.Role) UpdateRoleCommand
		wantErr     bool
		errContains string
		checkResult func(*testing.T, *role.Role, *testutil.MockRoleRepository)
	}{
		{
			name: "successfully update role display name and description",
			setupMocks: func(roleRepo *testutil.MockRoleRepository) *role.Role {
				testRole := createTestRole(false)
				roleRepo.AddRole(testRole)
				return testRole
			},
			command: func(r *role.Role) UpdateRoleCommand {
				return UpdateRoleCommand{
					RoleID:      r.ID(),
					DisplayName: "Updated Editor",
					Description: "Updated description",
				}
			},
			wantErr: false,
			checkResult: func(t *testing.T, r *role.Role, roleRepo *testutil.MockRoleRepository) {
				updatedRole, _ := roleRepo.FindByID(context.Background(), r.ID())
				assert.Equal(t, "Updated Editor", updatedRole.DisplayName())
				assert.Equal(t, "Updated description", updatedRole.Description())
			},
		},
		{
			name: "fail when updating system role",
			setupMocks: func(roleRepo *testutil.MockRoleRepository) *role.Role {
				systemRole := createTestRole(true)
				roleRepo.AddRole(systemRole)
				return systemRole
			},
			command: func(r *role.Role) UpdateRoleCommand {
				return UpdateRoleCommand{
					RoleID:      r.ID(),
					DisplayName: "Trying to update system role",
					Description: "New description",
				}
			},
			wantErr:     true,
			errContains: "system role",
		},
		{
			name: "fail when role not found",
			setupMocks: func(roleRepo *testutil.MockRoleRepository) *role.Role {
				testRole := createTestRole(false)
				roleRepo.AddRole(testRole)
				return testRole
			},
			command: func(r *role.Role) UpdateRoleCommand {
				return UpdateRoleCommand{
					RoleID:      uuid.New(),
					DisplayName: "Updated Editor",
					Description: "Updated description",
				}
			},
			wantErr:     true,
			errContains: "not found",
		},
		{
			name: "successfully update only description",
			setupMocks: func(roleRepo *testutil.MockRoleRepository) *role.Role {
				testRole := createTestRole(false)
				roleRepo.AddRole(testRole)
				return testRole
			},
			command: func(r *role.Role) UpdateRoleCommand {
				return UpdateRoleCommand{
					RoleID:      r.ID(),
					DisplayName: "",
					Description: "Only description changed",
				}
			},
			wantErr: false,
			checkResult: func(t *testing.T, r *role.Role, roleRepo *testutil.MockRoleRepository) {
				updatedRole, _ := roleRepo.FindByID(context.Background(), r.ID())
				assert.Equal(t, "Editor", updatedRole.DisplayName())
				assert.Equal(t, "Only description changed", updatedRole.Description())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roleRepo := testutil.NewMockRoleRepository()
			eventBus := testutil.NewMockEventBus()
			logger := testutil.NewNoopLogger()

			testRole := tt.setupMocks(roleRepo)

			handler := NewUpdateRoleHandler(roleRepo, eventBus, logger)
			cmd := tt.command(testRole)

			result, err := handler.Handle(ctx, cmd)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				if tt.checkResult != nil {
					tt.checkResult(t, testRole, roleRepo)
				}
			}
		})
	}
}
