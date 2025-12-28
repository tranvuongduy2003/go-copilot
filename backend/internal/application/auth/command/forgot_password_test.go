package authcommand

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/user"
	"github.com/tranvuongduy2003/go-copilot/pkg/testutil"
)

func TestForgotPasswordHandler_Handle(t *testing.T) {
	ctx := context.Background()

	createTestUser := func(email string, status user.Status) *user.User {
		now := time.Now().UTC()
		testUser, _ := user.ReconstructUser(user.ReconstructUserParams{
			ID:           uuid.New(),
			Email:        email,
			PasswordHash: "$2a$10$hashedpassword",
			FullName:     "Test User",
			Status:       status,
			CreatedAt:    now,
			UpdatedAt:    now,
		})
		return testUser
	}

	tests := []struct {
		name        string
		setupMocks  func(*testutil.MockUserRepository, *testutil.MockPasswordResetTokenStore)
		command     ForgotPasswordCommand
		wantResult  bool
		wantErr     bool
		errContains string
	}{
		{
			name: "successfully generate reset token for active user",
			setupMocks: func(userRepo *testutil.MockUserRepository, tokenStore *testutil.MockPasswordResetTokenStore) {
				testUser := createTestUser("test@example.com", user.StatusActive)
				userRepo.AddUser(testUser)
			},
			command: ForgotPasswordCommand{
				Email: "test@example.com",
			},
			wantResult: true,
			wantErr:    false,
		},
		{
			name: "return nil for non-existent email (prevent enumeration)",
			setupMocks: func(userRepo *testutil.MockUserRepository, tokenStore *testutil.MockPasswordResetTokenStore) {
			},
			command: ForgotPasswordCommand{
				Email: "nonexistent@example.com",
			},
			wantResult: false,
			wantErr:    false,
		},
		{
			name: "return nil for inactive user (prevent enumeration)",
			setupMocks: func(userRepo *testutil.MockUserRepository, tokenStore *testutil.MockPasswordResetTokenStore) {
				testUser := createTestUser("inactive@example.com", user.StatusInactive)
				userRepo.AddUser(testUser)
			},
			command: ForgotPasswordCommand{
				Email: "inactive@example.com",
			},
			wantResult: false,
			wantErr:    false,
		},
		{
			name: "return nil for banned user (prevent enumeration)",
			setupMocks: func(userRepo *testutil.MockUserRepository, tokenStore *testutil.MockPasswordResetTokenStore) {
				testUser := createTestUser("banned@example.com", user.StatusBanned)
				userRepo.AddUser(testUser)
			},
			command: ForgotPasswordCommand{
				Email: "banned@example.com",
			},
			wantResult: false,
			wantErr:    false,
		},
		{
			name: "fail when token store returns error",
			setupMocks: func(userRepo *testutil.MockUserRepository, tokenStore *testutil.MockPasswordResetTokenStore) {
				testUser := createTestUser("test@example.com", user.StatusActive)
				userRepo.AddUser(testUser)
				tokenStore.StoreError = assert.AnError
			},
			command: ForgotPasswordCommand{
				Email: "test@example.com",
			},
			wantResult: false,
			wantErr:    true,
			errContains: "store reset token",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			userRepo := testutil.NewMockUserRepository()
			tokenStore := testutil.NewMockPasswordResetTokenStore()
			tokenGenerator := testutil.NewMockTokenGenerator()
			eventBus := testutil.NewMockEventBus()
			logger := testutil.NewNoopLogger()

			testCase.setupMocks(userRepo, tokenStore)

			handler := NewForgotPasswordHandler(ForgotPasswordHandlerParams{
				UserRepository:          userRepo,
				TokenGenerator:          tokenGenerator,
				PasswordResetTokenStore: tokenStore,
				EventBus:                eventBus,
				ResetTokenTTL:           time.Hour,
				Logger:                  logger,
			})

			result, err := handler.Handle(ctx, testCase.command)

			if testCase.wantErr {
				require.Error(t, err)
				if testCase.errContains != "" {
					assert.Contains(t, err.Error(), testCase.errContains)
				}
			} else {
				require.NoError(t, err)
			}

			if testCase.wantResult {
				require.NotNil(t, result)
				assert.NotEmpty(t, result.ResetToken)
				assert.True(t, result.ExpiresAt.After(time.Now()))
			} else {
				assert.Nil(t, result)
			}
		})
	}
}

func TestForgotPasswordHandler_Handle_TokenStored(t *testing.T) {
	ctx := context.Background()
	userRepo := testutil.NewMockUserRepository()
	tokenStore := testutil.NewMockPasswordResetTokenStore()
	tokenGenerator := testutil.NewMockTokenGenerator()
	eventBus := testutil.NewMockEventBus()
	logger := testutil.NewNoopLogger()

	now := time.Now().UTC()
	testUser, _ := user.ReconstructUser(user.ReconstructUserParams{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "$2a$10$hashedpassword",
		FullName:     "Test User",
		Status:       user.StatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	userRepo.AddUser(testUser)

	handler := NewForgotPasswordHandler(ForgotPasswordHandlerParams{
		UserRepository:          userRepo,
		TokenGenerator:          tokenGenerator,
		PasswordResetTokenStore: tokenStore,
		EventBus:                eventBus,
		ResetTokenTTL:           time.Hour,
		Logger:                  logger,
	})

	result, err := handler.Handle(ctx, ForgotPasswordCommand{
		Email: "test@example.com",
	})

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Len(t, tokenStore.Tokens, 1)
	storedEmail, exists := tokenStore.Tokens[tokenGenerator.RefreshTokenHash]
	assert.True(t, exists)
	assert.Equal(t, "test@example.com", storedEmail)
}

func TestForgotPasswordHandler_Handle_EventPublished(t *testing.T) {
	ctx := context.Background()
	userRepo := testutil.NewMockUserRepository()
	tokenStore := testutil.NewMockPasswordResetTokenStore()
	tokenGenerator := testutil.NewMockTokenGenerator()
	eventBus := testutil.NewMockEventBus()
	logger := testutil.NewNoopLogger()

	now := time.Now().UTC()
	testUser, _ := user.ReconstructUser(user.ReconstructUserParams{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "$2a$10$hashedpassword",
		FullName:     "Test User",
		Status:       user.StatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	userRepo.AddUser(testUser)

	handler := NewForgotPasswordHandler(ForgotPasswordHandlerParams{
		UserRepository:          userRepo,
		TokenGenerator:          tokenGenerator,
		PasswordResetTokenStore: tokenStore,
		EventBus:                eventBus,
		ResetTokenTTL:           time.Hour,
		Logger:                  logger,
	})

	_, err := handler.Handle(ctx, ForgotPasswordCommand{
		Email: "test@example.com",
	})

	require.NoError(t, err)
	assert.Len(t, eventBus.PublishedEvents, 1)
}

func TestForgotPasswordHandler_Handle_NilTokenStore(t *testing.T) {
	ctx := context.Background()
	userRepo := testutil.NewMockUserRepository()
	tokenGenerator := testutil.NewMockTokenGenerator()
	eventBus := testutil.NewMockEventBus()
	logger := testutil.NewNoopLogger()

	now := time.Now().UTC()
	testUser, _ := user.ReconstructUser(user.ReconstructUserParams{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "$2a$10$hashedpassword",
		FullName:     "Test User",
		Status:       user.StatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	userRepo.AddUser(testUser)

	handler := NewForgotPasswordHandler(ForgotPasswordHandlerParams{
		UserRepository:          userRepo,
		TokenGenerator:          tokenGenerator,
		PasswordResetTokenStore: nil,
		EventBus:                eventBus,
		ResetTokenTTL:           time.Hour,
		Logger:                  logger,
	})

	result, err := handler.Handle(ctx, ForgotPasswordCommand{
		Email: "test@example.com",
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.ResetToken)
}
