package authcommand

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/auth"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/user"
	"github.com/tranvuongduy2003/go-copilot/pkg/testutil"
)

func TestResetPasswordHandler_Handle(t *testing.T) {
	ctx := context.Background()

	createTestUser := func(email string) *user.User {
		now := time.Now().UTC()
		testUser, _ := user.ReconstructUser(user.ReconstructUserParams{
			ID:           uuid.New(),
			Email:        email,
			PasswordHash: "$2a$10$oldhashedpassword",
			FullName:     "Test User",
			Status:       user.StatusActive,
			CreatedAt:    now,
			UpdatedAt:    now,
		})
		return testUser
	}

	tests := []struct {
		name        string
		setupMocks  func(*testutil.MockUserRepository, *testutil.MockPasswordResetTokenStore, *testutil.MockRefreshTokenRepository)
		command     ResetPasswordCommand
		wantErr     bool
		errContains string
		checkResult func(*testing.T, *testutil.MockUserRepository, *testutil.MockPasswordResetTokenStore, *testutil.MockRefreshTokenRepository)
	}{
		{
			name: "successfully reset password",
			setupMocks: func(userRepo *testutil.MockUserRepository, tokenStore *testutil.MockPasswordResetTokenStore, refreshRepo *testutil.MockRefreshTokenRepository) {
				testUser := createTestUser("test@example.com")
				userRepo.AddUser(testUser)
				tokenStore.Tokens["mock_hash"] = "test@example.com"
			},
			command: ResetPasswordCommand{
				ResetToken:  "valid_token",
				NewPassword: "NewPassword123!",
			},
			wantErr: false,
			checkResult: func(t *testing.T, userRepo *testutil.MockUserRepository, tokenStore *testutil.MockPasswordResetTokenStore, refreshRepo *testutil.MockRefreshTokenRepository) {
				assert.Empty(t, tokenStore.Tokens)
			},
		},
		{
			name: "fail when password is too weak",
			setupMocks: func(userRepo *testutil.MockUserRepository, tokenStore *testutil.MockPasswordResetTokenStore, refreshRepo *testutil.MockRefreshTokenRepository) {
				testUser := createTestUser("test@example.com")
				userRepo.AddUser(testUser)
				tokenStore.Tokens["mock_hash"] = "test@example.com"
			},
			command: ResetPasswordCommand{
				ResetToken:  "valid_token",
				NewPassword: "weak",
			},
			wantErr:     true,
			errContains: "password",
		},
		{
			name: "fail when token store is nil",
			setupMocks: func(userRepo *testutil.MockUserRepository, tokenStore *testutil.MockPasswordResetTokenStore, refreshRepo *testutil.MockRefreshTokenRepository) {
			},
			command: ResetPasswordCommand{
				ResetToken:  "valid_token",
				NewPassword: "NewPassword123!",
			},
			wantErr:     true,
			errContains: "invalid",
		},
		{
			name: "fail when token not found",
			setupMocks: func(userRepo *testutil.MockUserRepository, tokenStore *testutil.MockPasswordResetTokenStore, refreshRepo *testutil.MockRefreshTokenRepository) {
			},
			command: ResetPasswordCommand{
				ResetToken:  "invalid_token",
				NewPassword: "NewPassword123!",
			},
			wantErr:     true,
			errContains: "invalid",
		},
		{
			name: "fail when user not found by email",
			setupMocks: func(userRepo *testutil.MockUserRepository, tokenStore *testutil.MockPasswordResetTokenStore, refreshRepo *testutil.MockRefreshTokenRepository) {
				tokenStore.Tokens["mock_hash"] = "nonexistent@example.com"
			},
			command: ResetPasswordCommand{
				ResetToken:  "valid_token",
				NewPassword: "NewPassword123!",
			},
			wantErr:     true,
			errContains: "invalid",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			userRepo := testutil.NewMockUserRepository()
			tokenStore := testutil.NewMockPasswordResetTokenStore()
			refreshRepo := testutil.NewMockRefreshTokenRepository()
			tokenGenerator := testutil.NewMockTokenGenerator()
			passwordHasher := testutil.NewMockPasswordHasher()
			eventBus := testutil.NewMockEventBus()
			logger := testutil.NewNoopLogger()

			testCase.setupMocks(userRepo, tokenStore, refreshRepo)

			var store PasswordResetTokenStore = tokenStore
			if testCase.name == "fail when token store is nil" {
				store = nil
			}

			handler := NewResetPasswordHandler(ResetPasswordHandlerParams{
				UserRepository:          userRepo,
				RefreshTokenRepository:  refreshRepo,
				TokenGenerator:          tokenGenerator,
				PasswordHasher:          passwordHasher,
				PasswordResetTokenStore: store,
				EventBus:                eventBus,
				Logger:                  logger,
			})

			err := handler.Handle(ctx, testCase.command)

			if testCase.wantErr {
				require.Error(t, err)
				if testCase.errContains != "" {
					assert.Contains(t, err.Error(), testCase.errContains)
				}
			} else {
				require.NoError(t, err)
				if testCase.checkResult != nil {
					testCase.checkResult(t, userRepo, tokenStore, refreshRepo)
				}
			}
		})
	}
}

func TestResetPasswordHandler_Handle_RevokesAllRefreshTokens(t *testing.T) {
	ctx := context.Background()
	userRepo := testutil.NewMockUserRepository()
	tokenStore := testutil.NewMockPasswordResetTokenStore()
	refreshRepo := testutil.NewMockRefreshTokenRepository()
	tokenGenerator := testutil.NewMockTokenGenerator()
	passwordHasher := testutil.NewMockPasswordHasher()
	eventBus := testutil.NewMockEventBus()
	logger := testutil.NewNoopLogger()

	now := time.Now().UTC()
	testUser, _ := user.ReconstructUser(user.ReconstructUserParams{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "$2a$10$oldhashedpassword",
		FullName:     "Test User",
		Status:       user.StatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	userRepo.AddUser(testUser)
	tokenStore.Tokens["mock_hash"] = "test@example.com"

	refreshToken1, _ := auth.NewRefreshToken(auth.NewRefreshTokenParams{
		UserID:    testUser.ID(),
		TokenHash: "hash1",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})
	refreshToken2, _ := auth.NewRefreshToken(auth.NewRefreshTokenParams{
		UserID:    testUser.ID(),
		TokenHash: "hash2",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})
	refreshRepo.Create(ctx, refreshToken1)
	refreshRepo.Create(ctx, refreshToken2)

	handler := NewResetPasswordHandler(ResetPasswordHandlerParams{
		UserRepository:          userRepo,
		RefreshTokenRepository:  refreshRepo,
		TokenGenerator:          tokenGenerator,
		PasswordHasher:          passwordHasher,
		PasswordResetTokenStore: tokenStore,
		EventBus:                eventBus,
		Logger:                  logger,
	})

	err := handler.Handle(ctx, ResetPasswordCommand{
		ResetToken:  "valid_token",
		NewPassword: "NewPassword123!",
	})

	require.NoError(t, err)

	activeTokens, _ := refreshRepo.FindActiveByUserID(ctx, testUser.ID())
	assert.Empty(t, activeTokens)
}

func TestResetPasswordHandler_Handle_PublishesEvents(t *testing.T) {
	ctx := context.Background()
	userRepo := testutil.NewMockUserRepository()
	tokenStore := testutil.NewMockPasswordResetTokenStore()
	refreshRepo := testutil.NewMockRefreshTokenRepository()
	tokenGenerator := testutil.NewMockTokenGenerator()
	passwordHasher := testutil.NewMockPasswordHasher()
	eventBus := testutil.NewMockEventBus()
	logger := testutil.NewNoopLogger()

	now := time.Now().UTC()
	testUser, _ := user.ReconstructUser(user.ReconstructUserParams{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "$2a$10$oldhashedpassword",
		FullName:     "Test User",
		Status:       user.StatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	userRepo.AddUser(testUser)
	tokenStore.Tokens["mock_hash"] = "test@example.com"

	handler := NewResetPasswordHandler(ResetPasswordHandlerParams{
		UserRepository:          userRepo,
		RefreshTokenRepository:  refreshRepo,
		TokenGenerator:          tokenGenerator,
		PasswordHasher:          passwordHasher,
		PasswordResetTokenStore: tokenStore,
		EventBus:                eventBus,
		Logger:                  logger,
	})

	err := handler.Handle(ctx, ResetPasswordCommand{
		ResetToken:  "valid_token",
		NewPassword: "NewPassword123!",
	})

	require.NoError(t, err)
	assert.NotEmpty(t, eventBus.PublishedEvents)
}

func TestResetPasswordHandler_Handle_UpdatesUserPassword(t *testing.T) {
	ctx := context.Background()
	userRepo := testutil.NewMockUserRepository()
	tokenStore := testutil.NewMockPasswordResetTokenStore()
	refreshRepo := testutil.NewMockRefreshTokenRepository()
	tokenGenerator := testutil.NewMockTokenGenerator()
	passwordHasher := testutil.NewMockPasswordHasher()
	eventBus := testutil.NewMockEventBus()
	logger := testutil.NewNoopLogger()

	now := time.Now().UTC()
	testUser, _ := user.ReconstructUser(user.ReconstructUserParams{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "$2a$10$oldhashedpassword",
		FullName:     "Test User",
		Status:       user.StatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	originalPasswordHash := testUser.PasswordHash()
	userRepo.AddUser(testUser)
	tokenStore.Tokens["mock_hash"] = "test@example.com"

	handler := NewResetPasswordHandler(ResetPasswordHandlerParams{
		UserRepository:          userRepo,
		RefreshTokenRepository:  refreshRepo,
		TokenGenerator:          tokenGenerator,
		PasswordHasher:          passwordHasher,
		PasswordResetTokenStore: tokenStore,
		EventBus:                eventBus,
		Logger:                  logger,
	})

	err := handler.Handle(ctx, ResetPasswordCommand{
		ResetToken:  "valid_token",
		NewPassword: "NewPassword123!",
	})

	require.NoError(t, err)

	updatedUser, _ := userRepo.FindByEmail(ctx, "test@example.com")
	assert.NotEqual(t, originalPasswordHash, updatedUser.PasswordHash())
}

func TestResetPasswordHandler_Handle_DeletesResetToken(t *testing.T) {
	ctx := context.Background()
	userRepo := testutil.NewMockUserRepository()
	tokenStore := testutil.NewMockPasswordResetTokenStore()
	refreshRepo := testutil.NewMockRefreshTokenRepository()
	tokenGenerator := testutil.NewMockTokenGenerator()
	passwordHasher := testutil.NewMockPasswordHasher()
	eventBus := testutil.NewMockEventBus()
	logger := testutil.NewNoopLogger()

	now := time.Now().UTC()
	testUser, _ := user.ReconstructUser(user.ReconstructUserParams{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "$2a$10$oldhashedpassword",
		FullName:     "Test User",
		Status:       user.StatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	userRepo.AddUser(testUser)
	tokenStore.Tokens["mock_hash"] = "test@example.com"

	handler := NewResetPasswordHandler(ResetPasswordHandlerParams{
		UserRepository:          userRepo,
		RefreshTokenRepository:  refreshRepo,
		TokenGenerator:          tokenGenerator,
		PasswordHasher:          passwordHasher,
		PasswordResetTokenStore: tokenStore,
		EventBus:                eventBus,
		Logger:                  logger,
	})

	assert.Len(t, tokenStore.Tokens, 1)

	err := handler.Handle(ctx, ResetPasswordCommand{
		ResetToken:  "valid_token",
		NewPassword: "NewPassword123!",
	})

	require.NoError(t, err)
	assert.Empty(t, tokenStore.Tokens)
}
