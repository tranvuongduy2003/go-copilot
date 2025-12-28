package userquery

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

func TestGetUserHandler_Handle(t *testing.T) {
	ctx := context.Background()

	createTestUser := func() *user.User {
		testUser, _ := user.NewUser(user.NewUserParams{
			Email:        "test@example.com",
			PasswordHash: "$2a$10$hashedpassword",
			FullName:     "Test User",
		})
		return testUser
	}

	tests := []struct {
		name        string
		setupMocks  func(*testutil.MockUserRepository)
		query       func(*user.User) GetUserQuery
		wantErr     bool
		errContains string
		checkResult func(*testing.T, *user.User)
	}{
		{
			name: "successfully get user by ID",
			setupMocks: func(userRepo *testutil.MockUserRepository) {
			},
			query: func(u *user.User) GetUserQuery {
				return GetUserQuery{
					UserID: u.ID(),
				}
			},
			wantErr: false,
			checkResult: func(t *testing.T, u *user.User) {
			},
		},
		{
			name: "fail when user not found",
			setupMocks: func(userRepo *testutil.MockUserRepository) {
			},
			query: func(u *user.User) GetUserQuery {
				return GetUserQuery{
					UserID: uuid.New(),
				}
			},
			wantErr:     true,
			errContains: "not found",
		},
		{
			name: "fail when repository returns error",
			setupMocks: func(userRepo *testutil.MockUserRepository) {
				userRepo.FindError = errors.New("database error")
			},
			query: func(u *user.User) GetUserQuery {
				return GetUserQuery{
					UserID: u.ID(),
				}
			},
			wantErr:     true,
			errContains: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := testutil.NewMockUserRepository()
			logger := testutil.NewNoopLogger()

			testUser := createTestUser()
			userRepo.AddUser(testUser)

			tt.setupMocks(userRepo)

			handler := NewGetUserHandler(userRepo, logger)
			q := tt.query(testUser)

			result, err := handler.Handle(ctx, q)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, testUser.ID(), result.ID)
				assert.Equal(t, testUser.Email().String(), result.Email)
				assert.Equal(t, testUser.FullName().String(), result.FullName)
				if tt.checkResult != nil {
					tt.checkResult(t, testUser)
				}
			}
		})
	}
}

func TestGetUserByEmailHandler_Handle(t *testing.T) {
	ctx := context.Background()

	createTestUser := func() *user.User {
		testUser, _ := user.NewUser(user.NewUserParams{
			Email:        "test@example.com",
			PasswordHash: "$2a$10$hashedpassword",
			FullName:     "Test User",
		})
		return testUser
	}

	tests := []struct {
		name        string
		setupMocks  func(*testutil.MockUserRepository)
		query       GetUserByEmailQuery
		wantErr     bool
		errContains string
	}{
		{
			name: "successfully get user by email",
			setupMocks: func(userRepo *testutil.MockUserRepository) {
			},
			query: GetUserByEmailQuery{
				Email: "test@example.com",
			},
			wantErr: false,
		},
		{
			name: "fail when user not found",
			setupMocks: func(userRepo *testutil.MockUserRepository) {
			},
			query: GetUserByEmailQuery{
				Email: "nonexistent@example.com",
			},
			wantErr:     true,
			errContains: "not found",
		},
		{
			name: "fail when repository returns error",
			setupMocks: func(userRepo *testutil.MockUserRepository) {
				userRepo.FindError = errors.New("database error")
			},
			query: GetUserByEmailQuery{
				Email: "test@example.com",
			},
			wantErr:     true,
			errContains: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := testutil.NewMockUserRepository()
			logger := testutil.NewNoopLogger()

			testUser := createTestUser()
			userRepo.AddUser(testUser)

			tt.setupMocks(userRepo)

			handler := NewGetUserByEmailHandler(userRepo, logger)
			result, err := handler.Handle(ctx, tt.query)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tt.query.Email, result.Email)
			}
		})
	}
}

func TestCheckEmailExistsHandler_Handle(t *testing.T) {
	ctx := context.Background()

	createTestUser := func() *user.User {
		testUser, _ := user.NewUser(user.NewUserParams{
			Email:        "existing@example.com",
			PasswordHash: "$2a$10$hashedpassword",
			FullName:     "Test User",
		})
		return testUser
	}

	tests := []struct {
		name        string
		setupMocks  func(*testutil.MockUserRepository)
		query       CheckEmailExistsQuery
		wantResult  bool
		wantErr     bool
		errContains string
	}{
		{
			name: "returns true when email exists",
			setupMocks: func(userRepo *testutil.MockUserRepository) {
			},
			query: CheckEmailExistsQuery{
				Email: "existing@example.com",
			},
			wantResult: true,
			wantErr:    false,
		},
		{
			name: "returns false when email does not exist",
			setupMocks: func(userRepo *testutil.MockUserRepository) {
			},
			query: CheckEmailExistsQuery{
				Email: "nonexistent@example.com",
			},
			wantResult: false,
			wantErr:    false,
		},
		{
			name: "fail when repository returns error",
			setupMocks: func(userRepo *testutil.MockUserRepository) {
				userRepo.FindError = errors.New("database error")
			},
			query: CheckEmailExistsQuery{
				Email: "existing@example.com",
			},
			wantErr:     true,
			errContains: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := testutil.NewMockUserRepository()
			logger := testutil.NewNoopLogger()

			testUser := createTestUser()
			userRepo.AddUser(testUser)

			tt.setupMocks(userRepo)

			handler := NewCheckEmailExistsHandler(userRepo, logger)
			result, err := handler.Handle(ctx, tt.query)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantResult, result)
			}
		})
	}
}
