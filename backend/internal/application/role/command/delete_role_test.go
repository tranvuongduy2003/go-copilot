package rolecommand

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/role"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/user"
	"github.com/tranvuongduy2003/go-copilot/pkg/testutil"
)

func TestDeleteRoleHandler_Handle(t *testing.T) {
	ctx := context.Background()

	createTestRole := func(isSystem, isDefault bool) *role.Role {
		now := time.Now().UTC()
		testRole, _ := role.ReconstructRole(role.ReconstructRoleParams{
			ID:          uuid.New(),
			Name:        "editor",
			DisplayName: "Editor",
			Description: "Test role",
			IsSystem:    isSystem,
			IsDefault:   isDefault,
			Priority:    0,
			CreatedAt:   now,
			UpdatedAt:   now,
		})
		return testRole
	}

	createTestUserWithRole := func(roleID uuid.UUID) *user.User {
		now := time.Now().UTC()
		testUser, _ := user.ReconstructUser(user.ReconstructUserParams{
			ID:           uuid.New(),
			Email:        "test@example.com",
			PasswordHash: "$2a$10$hashedpassword",
			FullName:     "Test User",
			Status:       user.StatusActive,
			RoleIDs:      []uuid.UUID{roleID},
			CreatedAt:    now,
			UpdatedAt:    now,
		})
		return testUser
	}

	tests := []struct {
		name        string
		setupMocks  func(*testutil.MockRoleRepository, *testutil.MockUserRepository) *role.Role
		command     func(*role.Role) DeleteRoleCommand
		wantErr     bool
		errContains string
		checkResult func(*testing.T, *role.Role, *testutil.MockRoleRepository)
	}{
		{
			name: "successfully delete role",
			setupMocks: func(roleRepo *testutil.MockRoleRepository, userRepo *testutil.MockUserRepository) *role.Role {
				testRole := createTestRole(false, false)
				roleRepo.AddRole(testRole)
				return testRole
			},
			command: func(r *role.Role) DeleteRoleCommand {
				return DeleteRoleCommand{
					RoleID: r.ID(),
				}
			},
			wantErr: false,
			checkResult: func(t *testing.T, r *role.Role, roleRepo *testutil.MockRoleRepository) {
				_, err := roleRepo.FindByID(context.Background(), r.ID())
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "not found")
			},
		},
		{
			name: "fail when deleting system role",
			setupMocks: func(roleRepo *testutil.MockRoleRepository, userRepo *testutil.MockUserRepository) *role.Role {
				systemRole := createTestRole(true, false)
				roleRepo.AddRole(systemRole)
				return systemRole
			},
			command: func(r *role.Role) DeleteRoleCommand {
				return DeleteRoleCommand{
					RoleID: r.ID(),
				}
			},
			wantErr:     true,
			errContains: "system role",
		},
		{
			name: "fail when deleting default role",
			setupMocks: func(roleRepo *testutil.MockRoleRepository, userRepo *testutil.MockUserRepository) *role.Role {
				defaultRole := createTestRole(false, true)
				roleRepo.AddRole(defaultRole)
				return defaultRole
			},
			command: func(r *role.Role) DeleteRoleCommand {
				return DeleteRoleCommand{
					RoleID: r.ID(),
				}
			},
			wantErr:     true,
			errContains: "default role",
		},
		{
			name: "fail when role not found",
			setupMocks: func(roleRepo *testutil.MockRoleRepository, userRepo *testutil.MockUserRepository) *role.Role {
				testRole := createTestRole(false, false)
				roleRepo.AddRole(testRole)
				return testRole
			},
			command: func(r *role.Role) DeleteRoleCommand {
				return DeleteRoleCommand{
					RoleID: uuid.New(),
				}
			},
			wantErr:     true,
			errContains: "not found",
		},
		{
			name: "fail when role is assigned to users",
			setupMocks: func(roleRepo *testutil.MockRoleRepository, userRepo *testutil.MockUserRepository) *role.Role {
				testRole := createTestRole(false, false)
				roleRepo.AddRole(testRole)

				testUser := createTestUserWithRole(testRole.ID())
				userRepo.AddUser(testUser)

				return testRole
			},
			command: func(r *role.Role) DeleteRoleCommand {
				return DeleteRoleCommand{
					RoleID: r.ID(),
				}
			},
			wantErr:     true,
			errContains: "assigned to users",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roleRepo := testutil.NewMockRoleRepository()
			userRepo := testutil.NewMockUserRepository()
			eventBus := testutil.NewMockEventBus()
			logger := testutil.NewNoopLogger()

			testRole := tt.setupMocks(roleRepo, userRepo)

			handler := NewDeleteRoleHandler(roleRepo, userRepo, eventBus, logger)
			cmd := tt.command(testRole)

			err := handler.Handle(ctx, cmd)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				require.NoError(t, err)
				if tt.checkResult != nil {
					tt.checkResult(t, testRole, roleRepo)
				}
			}
		})
	}
}

func TestDeleteRoleHandler_Handle_MultipleUsersWithRole(t *testing.T) {
	ctx := context.Background()
	roleRepo := testutil.NewMockRoleRepository()
	userRepo := testutil.NewMockUserRepository()
	eventBus := testutil.NewMockEventBus()
	logger := testutil.NewNoopLogger()

	now := time.Now().UTC()
	testRole, _ := role.ReconstructRole(role.ReconstructRoleParams{
		ID:          uuid.New(),
		Name:        "editor",
		DisplayName: "Editor",
		Description: "Test role",
		IsSystem:    false,
		IsDefault:   false,
		Priority:    0,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	roleRepo.AddRole(testRole)

	user1, _ := user.ReconstructUser(user.ReconstructUserParams{
		ID:           uuid.New(),
		Email:        "user1@example.com",
		PasswordHash: "$2a$10$hashedpassword",
		FullName:     "User One",
		Status:       user.StatusActive,
		RoleIDs:      []uuid.UUID{testRole.ID()},
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	userRepo.AddUser(user1)

	user2, _ := user.ReconstructUser(user.ReconstructUserParams{
		ID:           uuid.New(),
		Email:        "user2@example.com",
		PasswordHash: "$2a$10$hashedpassword",
		FullName:     "User Two",
		Status:       user.StatusActive,
		RoleIDs:      []uuid.UUID{testRole.ID()},
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	userRepo.AddUser(user2)

	handler := NewDeleteRoleHandler(roleRepo, userRepo, eventBus, logger)

	err := handler.Handle(ctx, DeleteRoleCommand{
		RoleID: testRole.ID(),
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "assigned to users")

	_, findErr := roleRepo.FindByID(ctx, testRole.ID())
	require.NoError(t, findErr)
}

func TestDeleteRoleHandler_Handle_PublishesRoleDeletedEvent(t *testing.T) {
	ctx := context.Background()
	roleRepo := testutil.NewMockRoleRepository()
	userRepo := testutil.NewMockUserRepository()
	eventBus := testutil.NewMockEventBus()
	logger := testutil.NewNoopLogger()

	now := time.Now().UTC()
	testRole, _ := role.ReconstructRole(role.ReconstructRoleParams{
		ID:          uuid.New(),
		Name:        "editor",
		DisplayName: "Editor",
		Description: "Test role",
		IsSystem:    false,
		IsDefault:   false,
		Priority:    0,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	roleRepo.AddRole(testRole)

	handler := NewDeleteRoleHandler(roleRepo, userRepo, eventBus, logger)

	err := handler.Handle(ctx, DeleteRoleCommand{
		RoleID: testRole.ID(),
	})

	require.NoError(t, err)
	assert.Len(t, eventBus.PublishedEvents, 1)
	assert.Equal(t, role.EventTypeRoleDeleted, eventBus.PublishedEvents[0].EventType())
}
