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

func TestRevokeRoleFromUserHandler_Handle(t *testing.T) {
	ctx := context.Background()

	testRole, _ := role.NewRole(role.NewRoleParams{
		Name:        "editor",
		DisplayName: "Editor",
	})

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

	createTestUserWithoutRole := func() *user.User {
		testUser, _ := user.NewUser(user.NewUserParams{
			Email:        "test@example.com",
			PasswordHash: "$2a$10$hashedpassword",
			FullName:     "Test User",
		})
		return testUser
	}

	tests := []struct {
		name        string
		setupMocks  func(*testutil.MockUserRepository, *testutil.MockRoleRepository, *testutil.MockEventBus) *user.User
		command     func(*user.User) RevokeRoleFromUserCommand
		wantErr     bool
		errContains string
		checkResult func(*testing.T, *user.User)
	}{
		{
			name: "successfully revoke role from user",
			setupMocks: func(userRepo *testutil.MockUserRepository, roleRepo *testutil.MockRoleRepository, eventBus *testutil.MockEventBus) *user.User {
				roleRepo.AddRole(testRole)
				testUser := createTestUserWithRole(testRole.ID())
				userRepo.AddUser(testUser)
				return testUser
			},
			command: func(u *user.User) RevokeRoleFromUserCommand {
				return RevokeRoleFromUserCommand{
					UserID: u.ID(),
					RoleID: testRole.ID(),
				}
			},
			wantErr: false,
			checkResult: func(t *testing.T, u *user.User) {
				assert.False(t, u.HasRole(testRole.ID()))
			},
		},
		{
			name: "fail when user not found",
			setupMocks: func(userRepo *testutil.MockUserRepository, roleRepo *testutil.MockRoleRepository, eventBus *testutil.MockEventBus) *user.User {
				roleRepo.AddRole(testRole)
				testUser := createTestUserWithRole(testRole.ID())
				userRepo.AddUser(testUser)
				return testUser
			},
			command: func(u *user.User) RevokeRoleFromUserCommand {
				return RevokeRoleFromUserCommand{
					UserID: uuid.New(),
					RoleID: testRole.ID(),
				}
			},
			wantErr:     true,
			errContains: "not found",
		},
		{
			name: "fail when role not assigned to user",
			setupMocks: func(userRepo *testutil.MockUserRepository, roleRepo *testutil.MockRoleRepository, eventBus *testutil.MockEventBus) *user.User {
				testUser := createTestUserWithoutRole()
				userRepo.AddUser(testUser)
				return testUser
			},
			command: func(u *user.User) RevokeRoleFromUserCommand {
				return RevokeRoleFromUserCommand{
					UserID: u.ID(),
					RoleID: testRole.ID(),
				}
			},
			wantErr:     true,
			errContains: "not assigned",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := testutil.NewMockUserRepository()
			roleRepo := testutil.NewMockRoleRepository()
			eventBus := testutil.NewMockEventBus()
			logger := testutil.NewNoopLogger()

			testUser := tt.setupMocks(userRepo, roleRepo, eventBus)

			handler := NewRevokeRoleFromUserHandler(userRepo, eventBus, logger)
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
