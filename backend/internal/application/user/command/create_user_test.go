package usercommand

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/user"
	"github.com/tranvuongduy2003/go-copilot/pkg/testutil"
)

func TestCreateUserHandler_Handle(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		command     CreateUserCommand
		setupMocks  func(*testutil.MockUserRepository, *testutil.MockPasswordHasher, *testutil.MockEventBus)
		wantErr     bool
		errContains string
		checkResult func(*testing.T, *testutil.MockUserRepository, *testutil.MockEventBus)
	}{
		{
			name: "successfully create user",
			command: CreateUserCommand{
				Email:    "newuser@example.com",
				Password: "SecurePass123!",
				FullName: "New User",
			},
			setupMocks: func(userRepo *testutil.MockUserRepository, hasher *testutil.MockPasswordHasher, eventBus *testutil.MockEventBus) {
			},
			wantErr: false,
			checkResult: func(t *testing.T, userRepo *testutil.MockUserRepository, eventBus *testutil.MockEventBus) {
				assert.Len(t, userRepo.Users, 1)
				assert.Len(t, eventBus.PublishedEvents, 1)
				assert.Equal(t, "user.created", eventBus.PublishedEvents[0].EventType())
			},
		},
		{
			name: "fail when email already exists",
			command: CreateUserCommand{
				Email:    "existing@example.com",
				Password: "SecurePass123!",
				FullName: "New User",
			},
			setupMocks: func(userRepo *testutil.MockUserRepository, hasher *testutil.MockPasswordHasher, eventBus *testutil.MockEventBus) {
				existingUser, _ := user.NewUser(user.NewUserParams{
					Email:        "existing@example.com",
					PasswordHash: "$2a$10$hashedpassword",
					FullName:     "Existing User",
				})
				userRepo.AddUser(existingUser)
			},
			wantErr:     true,
			errContains: "already exists",
		},
		{
			name: "fail when password is too weak",
			command: CreateUserCommand{
				Email:    "newuser@example.com",
				Password: "weak",
				FullName: "New User",
			},
			setupMocks: func(userRepo *testutil.MockUserRepository, hasher *testutil.MockPasswordHasher, eventBus *testutil.MockEventBus) {
			},
			wantErr:     true,
			errContains: "password",
		},
		{
			name: "fail when password hasher returns error",
			command: CreateUserCommand{
				Email:    "newuser@example.com",
				Password: "SecurePass123!",
				FullName: "New User",
			},
			setupMocks: func(userRepo *testutil.MockUserRepository, hasher *testutil.MockPasswordHasher, eventBus *testutil.MockEventBus) {
				hasher.HashError = errors.New("hashing failed")
			},
			wantErr:     true,
			errContains: "hash password",
		},
		{
			name: "fail when repository create returns error",
			command: CreateUserCommand{
				Email:    "newuser@example.com",
				Password: "SecurePass123!",
				FullName: "New User",
			},
			setupMocks: func(userRepo *testutil.MockUserRepository, hasher *testutil.MockPasswordHasher, eventBus *testutil.MockEventBus) {
				userRepo.CreateError = errors.New("database error")
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

			tt.setupMocks(userRepo, hasher, eventBus)

			handler := NewCreateUserHandler(userRepo, hasher, eventBus, logger)
			result, err := handler.Handle(ctx, tt.command)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tt.command.Email, result.Email)
				assert.Equal(t, tt.command.FullName, result.FullName)
				if tt.checkResult != nil {
					tt.checkResult(t, userRepo, eventBus)
				}
			}
		})
	}
}
