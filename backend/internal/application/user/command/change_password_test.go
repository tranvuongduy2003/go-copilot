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

func TestChangePasswordHandler_Handle(t *testing.T) {
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
		setupMocks  func(*testutil.MockUserRepository, *testutil.MockPasswordHasher, *testutil.MockEventBus)
		command     func(*user.User) ChangePasswordCommand
		wantErr     bool
		errContains string
		checkResult func(*testing.T, *testutil.MockEventBus)
	}{
		{
			name: "successfully change password",
			setupMocks: func(userRepo *testutil.MockUserRepository, hasher *testutil.MockPasswordHasher, eventBus *testutil.MockEventBus) {
				hasher.VerifyResult = true
			},
			command: func(u *user.User) ChangePasswordCommand {
				return ChangePasswordCommand{
					UserID:          u.ID(),
					CurrentPassword: "OldPassword123!",
					NewPassword:     "NewSecurePass123!",
				}
			},
			wantErr: false,
			checkResult: func(t *testing.T, eventBus *testutil.MockEventBus) {
				assert.Len(t, eventBus.PublishedEvents, 1)
				assert.Equal(t, "user.password_changed", eventBus.PublishedEvents[0].EventType())
			},
		},
		{
			name: "fail when user not found",
			setupMocks: func(userRepo *testutil.MockUserRepository, hasher *testutil.MockPasswordHasher, eventBus *testutil.MockEventBus) {
			},
			command: func(u *user.User) ChangePasswordCommand {
				return ChangePasswordCommand{
					UserID:          uuid.New(),
					CurrentPassword: "OldPassword123!",
					NewPassword:     "NewSecurePass123!",
				}
			},
			wantErr:     true,
			errContains: "not found",
		},
		{
			name: "fail when current password is incorrect",
			setupMocks: func(userRepo *testutil.MockUserRepository, hasher *testutil.MockPasswordHasher, eventBus *testutil.MockEventBus) {
				hasher.VerifyResult = false
				hasher.VerifyError = errors.New("password mismatch")
			},
			command: func(u *user.User) ChangePasswordCommand {
				return ChangePasswordCommand{
					UserID:          u.ID(),
					CurrentPassword: "WrongPassword123!",
					NewPassword:     "NewSecurePass123!",
				}
			},
			wantErr:     true,
			errContains: "invalid password",
		},
		{
			name: "fail when new password is too weak",
			setupMocks: func(userRepo *testutil.MockUserRepository, hasher *testutil.MockPasswordHasher, eventBus *testutil.MockEventBus) {
				hasher.VerifyResult = true
			},
			command: func(u *user.User) ChangePasswordCommand {
				return ChangePasswordCommand{
					UserID:          u.ID(),
					CurrentPassword: "OldPassword123!",
					NewPassword:     "weak",
				}
			},
			wantErr:     true,
			errContains: "password",
		},
		{
			name: "fail when password hasher returns error",
			setupMocks: func(userRepo *testutil.MockUserRepository, hasher *testutil.MockPasswordHasher, eventBus *testutil.MockEventBus) {
				hasher.VerifyResult = true
				hasher.HashError = errors.New("hashing failed")
			},
			command: func(u *user.User) ChangePasswordCommand {
				return ChangePasswordCommand{
					UserID:          u.ID(),
					CurrentPassword: "OldPassword123!",
					NewPassword:     "NewSecurePass123!",
				}
			},
			wantErr:     true,
			errContains: "hash password",
		},
		{
			name: "fail when repository update returns error",
			setupMocks: func(userRepo *testutil.MockUserRepository, hasher *testutil.MockPasswordHasher, eventBus *testutil.MockEventBus) {
				hasher.VerifyResult = true
				userRepo.UpdateError = errors.New("database error")
			},
			command: func(u *user.User) ChangePasswordCommand {
				return ChangePasswordCommand{
					UserID:          u.ID(),
					CurrentPassword: "OldPassword123!",
					NewPassword:     "NewSecurePass123!",
				}
			},
			wantErr:     true,
			errContains: "save user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := testutil.NewMockUserRepository()
			hasher := testutil.NewMockPasswordHasher()
			eventBus := testutil.NewMockEventBus()
			logger := testutil.NewNoopLogger()

			testUser := createTestUser()
			userRepo.AddUser(testUser)

			tt.setupMocks(userRepo, hasher, eventBus)

			handler := NewChangePasswordHandler(userRepo, hasher, eventBus, logger)
			cmd := tt.command(testUser)

			err := handler.Handle(ctx, cmd)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				require.NoError(t, err)
				if tt.checkResult != nil {
					tt.checkResult(t, eventBus)
				}
			}
		})
	}
}
