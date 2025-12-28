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

func TestUpdateUserHandler_Handle(t *testing.T) {
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
		command     func(*user.User) UpdateUserCommand
		wantErr     bool
		errContains string
		checkResult func(*testing.T, *testutil.MockUserRepository, *testutil.MockEventBus)
	}{
		{
			name: "successfully update user full name",
			setupMocks: func(userRepo *testutil.MockUserRepository, eventBus *testutil.MockEventBus) {
			},
			command: func(u *user.User) UpdateUserCommand {
				newName := "Updated Name"
				return UpdateUserCommand{
					UserID:   u.ID(),
					FullName: &newName,
				}
			},
			wantErr: false,
			checkResult: func(t *testing.T, userRepo *testutil.MockUserRepository, eventBus *testutil.MockEventBus) {
				assert.Len(t, eventBus.PublishedEvents, 1)
				assert.Equal(t, "user.profile_updated", eventBus.PublishedEvents[0].EventType())
			},
		},
		{
			name: "successfully update user with no changes",
			setupMocks: func(userRepo *testutil.MockUserRepository, eventBus *testutil.MockEventBus) {
			},
			command: func(u *user.User) UpdateUserCommand {
				return UpdateUserCommand{
					UserID:   u.ID(),
					FullName: nil,
				}
			},
			wantErr: false,
			checkResult: func(t *testing.T, userRepo *testutil.MockUserRepository, eventBus *testutil.MockEventBus) {
				assert.Len(t, eventBus.PublishedEvents, 0)
			},
		},
		{
			name: "fail when user not found",
			setupMocks: func(userRepo *testutil.MockUserRepository, eventBus *testutil.MockEventBus) {
			},
			command: func(u *user.User) UpdateUserCommand {
				newName := "Updated Name"
				return UpdateUserCommand{
					UserID:   uuid.New(),
					FullName: &newName,
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
			command: func(u *user.User) UpdateUserCommand {
				newName := "Updated Name"
				return UpdateUserCommand{
					UserID:   u.ID(),
					FullName: &newName,
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

			handler := NewUpdateUserHandler(userRepo, eventBus, logger)
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
					tt.checkResult(t, userRepo, eventBus)
				}
			}
		})
	}
}
