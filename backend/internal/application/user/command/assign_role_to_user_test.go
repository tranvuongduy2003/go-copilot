package usercommand

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/role"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/user"
	"github.com/tranvuongduy2003/go-copilot/pkg/testutil"
)

func TestAssignRoleToUserHandler_Handle(t *testing.T) {
	ctx := context.Background()

	testRole, _ := role.NewRole(role.NewRoleParams{
		Name:        "editor",
		DisplayName: "Editor",
	})

	createTestUser := func() *user.User {
		testUser, _ := user.NewUser(user.NewUserParams{
			Email:        "test@example.com",
			PasswordHash: "$2a$10$hashedpassword",
			FullName:     "Test User",
		})
		return testUser
	}

	tests := []struct {
		name          string
		setupMocks    func(*testutil.MockUserRepository, *testutil.MockRoleRepository, *testutil.MockEventBus)
		command       func(*user.User) AssignRoleToUserCommand
		wantErr       bool
		errContains   string
		checkResult   func(*testing.T, *user.User)
	}{
		{
			name: "successfully assign role to user",
			setupMocks: func(userRepo *testutil.MockUserRepository, roleRepo *testutil.MockRoleRepository, eventBus *testutil.MockEventBus) {
				roleRepo.AddRole(testRole)
			},
			command: func(u *user.User) AssignRoleToUserCommand {
				return AssignRoleToUserCommand{
					UserID: u.ID(),
					RoleID: testRole.ID(),
				}
			},
			wantErr: false,
			checkResult: func(t *testing.T, u *user.User) {
				assert.True(t, u.HasRole(testRole.ID()))
			},
		},
		{
			name: "fail when user not found",
			setupMocks: func(userRepo *testutil.MockUserRepository, roleRepo *testutil.MockRoleRepository, eventBus *testutil.MockEventBus) {
				roleRepo.AddRole(testRole)
			},
			command: func(u *user.User) AssignRoleToUserCommand {
				return AssignRoleToUserCommand{
					UserID: uuid.New(),
					RoleID: testRole.ID(),
				}
			},
			wantErr:     true,
			errContains: "not found",
		},
		{
			name: "fail when role not found",
			setupMocks: func(userRepo *testutil.MockUserRepository, roleRepo *testutil.MockRoleRepository, eventBus *testutil.MockEventBus) {
			},
			command: func(u *user.User) AssignRoleToUserCommand {
				return AssignRoleToUserCommand{
					UserID: u.ID(),
					RoleID: uuid.New(),
				}
			},
			wantErr:     true,
			errContains: "not found",
		},
		{
			name: "fail when role already assigned",
			setupMocks: func(userRepo *testutil.MockUserRepository, roleRepo *testutil.MockRoleRepository, eventBus *testutil.MockEventBus) {
				roleRepo.AddRole(testRole)
			},
			command: func(u *user.User) AssignRoleToUserCommand {
				u.AssignRole(testRole.ID())
				return AssignRoleToUserCommand{
					UserID: u.ID(),
					RoleID: testRole.ID(),
				}
			},
			wantErr:     true,
			errContains: "already assigned",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := testutil.NewMockUserRepository()
			roleRepo := testutil.NewMockRoleRepository()
			eventBus := testutil.NewMockEventBus()
			logger := testutil.NewNoopLogger()

			testUser := createTestUser()
			userRepo.AddUser(testUser)

			tt.setupMocks(userRepo, roleRepo, eventBus)

			handler := NewAssignRoleToUserHandler(userRepo, roleRepo, eventBus, logger)
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
