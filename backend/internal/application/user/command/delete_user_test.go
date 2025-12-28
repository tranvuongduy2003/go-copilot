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

func TestDeleteUserHandler_Handle(t *testing.T) {
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
		setupMocks  func(*testutil.MockUserRepository, *testutil.MockEventBus)
		command     func(*user.User) DeleteUserCommand
		wantErr     bool
		errContains string
		checkResult func(*testing.T, *testutil.MockUserRepository, *testutil.MockEventBus, *user.User)
	}{
		{
			name: "successfully delete user",
			setupMocks: func(userRepo *testutil.MockUserRepository, eventBus *testutil.MockEventBus) {
			},
			command: func(u *user.User) DeleteUserCommand {
				return DeleteUserCommand{
					UserID: u.ID(),
				}
			},
			wantErr: false,
			checkResult: func(t *testing.T, userRepo *testutil.MockUserRepository, eventBus *testutil.MockEventBus, u *user.User) {
				assert.Len(t, eventBus.PublishedEvents, 1)
				assert.Equal(t, "user.deleted", eventBus.PublishedEvents[0].EventType())
				updatedUser, _ := userRepo.FindByID(ctx, u.ID())
				assert.NotNil(t, updatedUser.DeletedAt())
			},
		},
		{
			name: "fail when user not found",
			setupMocks: func(userRepo *testutil.MockUserRepository, eventBus *testutil.MockEventBus) {
			},
			command: func(u *user.User) DeleteUserCommand {
				return DeleteUserCommand{
					UserID: uuid.New(),
				}
			},
			wantErr:     true,
			errContains: "not found",
		},
		{
			name: "fail when repository update returns error",
			setupMocks: func(userRepo *testutil.MockUserRepository, eventBus *testutil.MockEventBus) {
				userRepo.UpdateError = errors.New("database error")
			},
			command: func(u *user.User) DeleteUserCommand {
				return DeleteUserCommand{
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

			tt.setupMocks(userRepo, eventBus)

			handler := NewDeleteUserHandler(userRepo, eventBus, logger)
			cmd := tt.command(testUser)

			err := handler.Handle(ctx, cmd)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				require.NoError(t, err)
				if tt.checkResult != nil {
					tt.checkResult(t, userRepo, eventBus, testUser)
				}
			}
		})
	}
}
