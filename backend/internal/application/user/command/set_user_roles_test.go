package usercommand

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

func TestSetUserRolesHandler_Handle(t *testing.T) {
	ctx := context.Background()

	editorRole, _ := role.NewRole(role.NewRoleParams{
		Name:        "editor",
		DisplayName: "Editor",
	})

	viewerRole, _ := role.NewRole(role.NewRoleParams{
		Name:        "viewer",
		DisplayName: "Viewer",
	})

	adminRole, _ := role.NewRole(role.NewRoleParams{
		Name:        "admin",
		DisplayName: "Administrator",
	})

	createTestUser := func() *user.User {
		testUser, _ := user.NewUser(user.NewUserParams{
			Email:        "test@example.com",
			PasswordHash: "$2a$10$hashedpassword",
			FullName:     "Test User",
		})
		return testUser
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
		setupMocks  func(*testutil.MockUserRepository, *testutil.MockRoleRepository, *testutil.MockEventBus) *user.User
		command     func(*user.User) SetUserRolesCommand
		wantErr     bool
		errContains string
		checkResult func(*testing.T, *user.User)
	}{
		{
			name: "successfully set roles for user",
			setupMocks: func(userRepo *testutil.MockUserRepository, roleRepo *testutil.MockRoleRepository, eventBus *testutil.MockEventBus) *user.User {
				roleRepo.AddRole(editorRole)
				roleRepo.AddRole(viewerRole)
				testUser := createTestUser()
				userRepo.AddUser(testUser)
				return testUser
			},
			command: func(u *user.User) SetUserRolesCommand {
				return SetUserRolesCommand{
					UserID:  u.ID(),
					RoleIDs: []uuid.UUID{editorRole.ID(), viewerRole.ID()},
				}
			},
			wantErr: false,
			checkResult: func(t *testing.T, u *user.User) {
				assert.True(t, u.HasRole(editorRole.ID()))
				assert.True(t, u.HasRole(viewerRole.ID()))
				assert.Len(t, u.RoleIDs(), 2)
			},
		},
		{
			name: "successfully replace existing roles",
			setupMocks: func(userRepo *testutil.MockUserRepository, roleRepo *testutil.MockRoleRepository, eventBus *testutil.MockEventBus) *user.User {
				roleRepo.AddRole(editorRole)
				roleRepo.AddRole(viewerRole)
				roleRepo.AddRole(adminRole)
				testUser := createTestUserWithRole(editorRole.ID())
				userRepo.AddUser(testUser)
				return testUser
			},
			command: func(u *user.User) SetUserRolesCommand {
				return SetUserRolesCommand{
					UserID:  u.ID(),
					RoleIDs: []uuid.UUID{viewerRole.ID(), adminRole.ID()},
				}
			},
			wantErr: false,
			checkResult: func(t *testing.T, u *user.User) {
				assert.False(t, u.HasRole(editorRole.ID()))
				assert.True(t, u.HasRole(viewerRole.ID()))
				assert.True(t, u.HasRole(adminRole.ID()))
				assert.Len(t, u.RoleIDs(), 2)
			},
		},
		{
			name: "fail when user not found",
			setupMocks: func(userRepo *testutil.MockUserRepository, roleRepo *testutil.MockRoleRepository, eventBus *testutil.MockEventBus) *user.User {
				roleRepo.AddRole(editorRole)
				testUser := createTestUser()
				userRepo.AddUser(testUser)
				return testUser
			},
			command: func(u *user.User) SetUserRolesCommand {
				return SetUserRolesCommand{
					UserID:  uuid.New(),
					RoleIDs: []uuid.UUID{editorRole.ID()},
				}
			},
			wantErr:     true,
			errContains: "not found",
		},
		{
			name: "fail when role not found",
			setupMocks: func(userRepo *testutil.MockUserRepository, roleRepo *testutil.MockRoleRepository, eventBus *testutil.MockEventBus) *user.User {
				testUser := createTestUser()
				userRepo.AddUser(testUser)
				return testUser
			},
			command: func(u *user.User) SetUserRolesCommand {
				return SetUserRolesCommand{
					UserID:  u.ID(),
					RoleIDs: []uuid.UUID{uuid.New()},
				}
			},
			wantErr:     true,
			errContains: "not found",
		},
		{
			name: "successfully set empty roles",
			setupMocks: func(userRepo *testutil.MockUserRepository, roleRepo *testutil.MockRoleRepository, eventBus *testutil.MockEventBus) *user.User {
				roleRepo.AddRole(editorRole)
				testUser := createTestUserWithRole(editorRole.ID())
				userRepo.AddUser(testUser)
				return testUser
			},
			command: func(u *user.User) SetUserRolesCommand {
				return SetUserRolesCommand{
					UserID:  u.ID(),
					RoleIDs: []uuid.UUID{},
				}
			},
			wantErr: false,
			checkResult: func(t *testing.T, u *user.User) {
				assert.False(t, u.HasRole(editorRole.ID()))
				assert.Len(t, u.RoleIDs(), 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := testutil.NewMockUserRepository()
			roleRepo := testutil.NewMockRoleRepository()
			eventBus := testutil.NewMockEventBus()
			logger := testutil.NewNoopLogger()

			testUser := tt.setupMocks(userRepo, roleRepo, eventBus)

			handler := NewSetUserRolesHandler(userRepo, roleRepo, eventBus, logger)
			cmd := tt.command(testUser)

			result, err := handler.Handle(ctx, cmd)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				if tt.checkResult != nil {
					updatedUser, _ := userRepo.FindByID(ctx, testUser.ID())
					tt.checkResult(t, updatedUser)
				}
			}
		})
	}
}
