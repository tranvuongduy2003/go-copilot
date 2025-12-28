package usercommand

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/user"
	"github.com/tranvuongduy2003/go-copilot/pkg/testutil"
)

func TestActivateUserHandler_Handle(t *testing.T) {
	ctx := context.Background()

	createTestUser := func() *user.User {
		testUser, _ := user.NewUser(user.NewUserParams{
			Email:        "test@example.com",
			PasswordHash: "$2a$10$hashedpassword",
			FullName:     "Test User",
		})
		testUser.ClearDomainEvents()
		return testUser
	}

	tests := []struct {
		name        string
		setupMocks  func(*testutil.MockUserRepository, *testutil.MockEventBus, *user.User)
		command     func(*user.User) ActivateUserCommand
		wantErr     bool
		errContains string
		checkResult func(*testing.T, *testutil.MockUserRepository, *testutil.MockEventBus, *user.User)
	}{
		{
			name: "successfully activate pending user",
			setupMocks: func(userRepo *testutil.MockUserRepository, eventBus *testutil.MockEventBus, u *user.User) {
			},
			command: func(u *user.User) ActivateUserCommand {
				return ActivateUserCommand{
					UserID: u.ID(),
				}
			},
			wantErr: false,
			checkResult: func(t *testing.T, userRepo *testutil.MockUserRepository, eventBus *testutil.MockEventBus, u *user.User) {
				assert.Len(t, eventBus.PublishedEvents, 1)
				assert.Equal(t, "user.activated", eventBus.PublishedEvents[0].EventType())
				updatedUser, _ := userRepo.FindByID(ctx, u.ID())
				assert.Equal(t, user.StatusActive, updatedUser.Status())
			},
		},
		{
			name: "fail when user not found",
			setupMocks: func(userRepo *testutil.MockUserRepository, eventBus *testutil.MockEventBus, u *user.User) {
			},
			command: func(u *user.User) ActivateUserCommand {
				return ActivateUserCommand{
					UserID: uuid.New(),
				}
			},
			wantErr:     true,
			errContains: "not found",
		},
		{
			name: "fail when user is already active",
			setupMocks: func(userRepo *testutil.MockUserRepository, eventBus *testutil.MockEventBus, u *user.User) {
				u.Activate()
				u.ClearDomainEvents()
			},
			command: func(u *user.User) ActivateUserCommand {
				return ActivateUserCommand{
					UserID: u.ID(),
				}
			},
			wantErr:     true,
			errContains: "activate user",
		},
		{
			name: "fail when repository update returns error",
			setupMocks: func(userRepo *testutil.MockUserRepository, eventBus *testutil.MockEventBus, u *user.User) {
				userRepo.UpdateError = errors.New("database error")
			},
			command: func(u *user.User) ActivateUserCommand {
				return ActivateUserCommand{
					UserID: u.ID(),
				}
			},
			wantErr:     true,
			errContains: "save user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := testutil.NewMockUserRepository()
			eventBus := testutil.NewMockEventBus()
			logger := testutil.NewNoopLogger()

			testUser := createTestUser()
			userRepo.AddUser(testUser)

			tt.setupMocks(userRepo, eventBus, testUser)

			handler := NewActivateUserHandler(userRepo, eventBus, logger)
			cmd := tt.command(testUser)

			result, err := handler.Handle(ctx, cmd)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, "active", result.Status)
				if tt.checkResult != nil {
					tt.checkResult(t, userRepo, eventBus, testUser)
				}
			}
		})
	}
}
