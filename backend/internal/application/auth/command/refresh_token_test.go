package authcommand

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/auth"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/user"
	"github.com/tranvuongduy2003/go-copilot/pkg/testutil"
)

func TestRefreshTokenHandler_Handle(t *testing.T) {
	ctx := context.Background()

	createActiveUser := func() *user.User {
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
		return testUser
	}

	createInactiveUser := func() *user.User {
		now := time.Now().UTC()
		testUser, _ := user.ReconstructUser(user.ReconstructUserParams{
			ID:           uuid.New(),
			Email:        "inactive@example.com",
			PasswordHash: "$2a$10$hashedpassword",
			FullName:     "Inactive User",
			Status:       user.StatusInactive,
			CreatedAt:    now,
			UpdatedAt:    now,
		})
		return testUser
	}

	createValidRefreshToken := func(userID uuid.UUID, hash string) *auth.RefreshToken {
		token, _ := auth.NewRefreshToken(auth.NewRefreshTokenParams{
			UserID:    userID,
			TokenHash: hash,
			ExpiresAt: time.Now().UTC().Add(24 * time.Hour),
		})
		return token
	}

	createExpiredRefreshToken := func(userID uuid.UUID, hash string) *auth.RefreshToken {
		now := time.Now().UTC()
		return auth.ReconstructRefreshToken(auth.ReconstructRefreshTokenParams{
			ID:        uuid.New(),
			UserID:    userID,
			TokenHash: hash,
			ExpiresAt: now.Add(-1 * time.Hour),
			CreatedAt: now.Add(-2 * time.Hour),
			IsRevoked: false,
		})
	}

	tests := []struct {
		name        string
		setupMocks  func(*testutil.MockUserRepository, *testutil.MockRoleRepository, *testutil.MockPermissionRepository, *testutil.MockRefreshTokenRepository, *testutil.MockTokenGenerator, *testutil.MockTokenBlacklist) string
		command     func(string) RefreshTokenCommand
		wantErr     bool
		errContains string
		checkResult func(*testing.T, *testutil.MockRefreshTokenRepository)
	}{
		{
			name: "successfully refresh token",
			setupMocks: func(userRepo *testutil.MockUserRepository, roleRepo *testutil.MockRoleRepository, permissionRepo *testutil.MockPermissionRepository, tokenRepo *testutil.MockRefreshTokenRepository, tokenGen *testutil.MockTokenGenerator, blacklist *testutil.MockTokenBlacklist) string {
				testUser := createActiveUser()
				userRepo.AddUser(testUser)

				tokenGen.RefreshTokenHash = "mock_hash"
				token := createValidRefreshToken(testUser.ID(), "mock_hash")
				tokenRepo.Tokens[token.ID()] = token
				tokenRepo.HashIndex["mock_hash"] = token
				tokenRepo.UserTokens[testUser.ID()] = append(tokenRepo.UserTokens[testUser.ID()], token)

				return "valid_refresh_token"
			},
			command: func(refreshToken string) RefreshTokenCommand {
				return RefreshTokenCommand{
					RefreshToken: refreshToken,
					IPAddress:    net.ParseIP("192.168.1.1"),
					UserAgent:    "Mozilla/5.0",
				}
			},
			wantErr: false,
			checkResult: func(t *testing.T, tokenRepo *testutil.MockRefreshTokenRepository) {
				assert.GreaterOrEqual(t, len(tokenRepo.Tokens), 1)
			},
		},
		{
			name: "fail when refresh token not found",
			setupMocks: func(userRepo *testutil.MockUserRepository, roleRepo *testutil.MockRoleRepository, permissionRepo *testutil.MockPermissionRepository, tokenRepo *testutil.MockRefreshTokenRepository, tokenGen *testutil.MockTokenGenerator, blacklist *testutil.MockTokenBlacklist) string {
				tokenGen.RefreshTokenHash = "nonexistent_hash"
				return "invalid_refresh_token"
			},
			command: func(refreshToken string) RefreshTokenCommand {
				return RefreshTokenCommand{
					RefreshToken: refreshToken,
					IPAddress:    net.ParseIP("192.168.1.1"),
					UserAgent:    "Mozilla/5.0",
				}
			},
			wantErr:     true,
			errContains: "invalid",
		},
		{
			name: "fail when refresh token is expired",
			setupMocks: func(userRepo *testutil.MockUserRepository, roleRepo *testutil.MockRoleRepository, permissionRepo *testutil.MockPermissionRepository, tokenRepo *testutil.MockRefreshTokenRepository, tokenGen *testutil.MockTokenGenerator, blacklist *testutil.MockTokenBlacklist) string {
				testUser := createActiveUser()
				userRepo.AddUser(testUser)

				tokenGen.RefreshTokenHash = "expired_hash"
				token := createExpiredRefreshToken(testUser.ID(), "expired_hash")
				tokenRepo.Tokens[token.ID()] = token
				tokenRepo.HashIndex["expired_hash"] = token
				tokenRepo.UserTokens[testUser.ID()] = append(tokenRepo.UserTokens[testUser.ID()], token)

				return "expired_refresh_token"
			},
			command: func(refreshToken string) RefreshTokenCommand {
				return RefreshTokenCommand{
					RefreshToken: refreshToken,
					IPAddress:    net.ParseIP("192.168.1.1"),
					UserAgent:    "Mozilla/5.0",
				}
			},
			wantErr:     true,
			errContains: "invalid",
		},
		{
			name: "fail when user is inactive",
			setupMocks: func(userRepo *testutil.MockUserRepository, roleRepo *testutil.MockRoleRepository, permissionRepo *testutil.MockPermissionRepository, tokenRepo *testutil.MockRefreshTokenRepository, tokenGen *testutil.MockTokenGenerator, blacklist *testutil.MockTokenBlacklist) string {
				inactiveUser := createInactiveUser()
				userRepo.AddUser(inactiveUser)

				tokenGen.RefreshTokenHash = "inactive_user_hash"
				token := createValidRefreshToken(inactiveUser.ID(), "inactive_user_hash")
				tokenRepo.Tokens[token.ID()] = token
				tokenRepo.HashIndex["inactive_user_hash"] = token
				tokenRepo.UserTokens[inactiveUser.ID()] = append(tokenRepo.UserTokens[inactiveUser.ID()], token)

				return "inactive_user_token"
			},
			command: func(refreshToken string) RefreshTokenCommand {
				return RefreshTokenCommand{
					RefreshToken: refreshToken,
					IPAddress:    net.ParseIP("192.168.1.1"),
					UserAgent:    "Mozilla/5.0",
				}
			},
			wantErr:     true,
			errContains: "not active",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := testutil.NewMockUserRepository()
			roleRepo := testutil.NewMockRoleRepository()
			permissionRepo := testutil.NewMockPermissionRepository()
			tokenRepo := testutil.NewMockRefreshTokenRepository()
			tokenGen := testutil.NewMockTokenGenerator()
			blacklist := testutil.NewMockTokenBlacklist()
			eventBus := testutil.NewMockEventBus()
			logger := testutil.NewNoopLogger()

			refreshToken := tt.setupMocks(userRepo, roleRepo, permissionRepo, tokenRepo, tokenGen, blacklist)

			handler := NewRefreshTokenHandler(RefreshTokenHandlerParams{
				UserRepository:         userRepo,
				RoleRepository:         roleRepo,
				PermissionRepository:   permissionRepo,
				RefreshTokenRepository: tokenRepo,
				TokenGenerator:         tokenGen,
				TokenBlacklist:         blacklist,
				EventBus:               eventBus,
				RefreshTokenTTL:        24 * time.Hour,
				Logger:                 logger,
			})

			result, err := handler.Handle(ctx, tt.command(refreshToken))

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.NotEmpty(t, result.AccessToken)
				assert.NotEmpty(t, result.RefreshToken)
				if tt.checkResult != nil {
					tt.checkResult(t, tokenRepo)
				}
			}
		})
	}
}

func TestRefreshTokenHandler_Handle_RotatesToken(t *testing.T) {
	ctx := context.Background()
	userRepo := testutil.NewMockUserRepository()
	roleRepo := testutil.NewMockRoleRepository()
	permissionRepo := testutil.NewMockPermissionRepository()
	tokenRepo := testutil.NewMockRefreshTokenRepository()
	tokenGen := testutil.NewMockTokenGenerator()
	blacklist := testutil.NewMockTokenBlacklist()
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

	tokenGen.RefreshTokenHash = "original_hash"
	originalToken, _ := auth.NewRefreshToken(auth.NewRefreshTokenParams{
		UserID:    testUser.ID(),
		TokenHash: "original_hash",
		ExpiresAt: time.Now().UTC().Add(24 * time.Hour),
	})
	tokenRepo.Tokens[originalToken.ID()] = originalToken
	tokenRepo.HashIndex["original_hash"] = originalToken
	tokenRepo.UserTokens[testUser.ID()] = []*auth.RefreshToken{originalToken}

	handler := NewRefreshTokenHandler(RefreshTokenHandlerParams{
		UserRepository:         userRepo,
		RoleRepository:         roleRepo,
		PermissionRepository:   permissionRepo,
		RefreshTokenRepository: tokenRepo,
		TokenGenerator:         tokenGen,
		TokenBlacklist:         blacklist,
		EventBus:               eventBus,
		RefreshTokenTTL:        24 * time.Hour,
		Logger:                 logger,
	})

	result, err := handler.Handle(ctx, RefreshTokenCommand{
		RefreshToken: "original_token",
		IPAddress:    net.ParseIP("192.168.1.1"),
		UserAgent:    "Mozilla/5.0",
	})

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.True(t, originalToken.IsRevoked())
}

func TestRefreshTokenHandler_Handle_PublishesEvent(t *testing.T) {
	ctx := context.Background()
	userRepo := testutil.NewMockUserRepository()
	roleRepo := testutil.NewMockRoleRepository()
	permissionRepo := testutil.NewMockPermissionRepository()
	tokenRepo := testutil.NewMockRefreshTokenRepository()
	tokenGen := testutil.NewMockTokenGenerator()
	blacklist := testutil.NewMockTokenBlacklist()
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

	tokenGen.RefreshTokenHash = "test_hash"
	token, _ := auth.NewRefreshToken(auth.NewRefreshTokenParams{
		UserID:    testUser.ID(),
		TokenHash: "test_hash",
		ExpiresAt: time.Now().UTC().Add(24 * time.Hour),
	})
	tokenRepo.Tokens[token.ID()] = token
	tokenRepo.HashIndex["test_hash"] = token
	tokenRepo.UserTokens[testUser.ID()] = []*auth.RefreshToken{token}

	handler := NewRefreshTokenHandler(RefreshTokenHandlerParams{
		UserRepository:         userRepo,
		RoleRepository:         roleRepo,
		PermissionRepository:   permissionRepo,
		RefreshTokenRepository: tokenRepo,
		TokenGenerator:         tokenGen,
		TokenBlacklist:         blacklist,
		EventBus:               eventBus,
		RefreshTokenTTL:        24 * time.Hour,
		Logger:                 logger,
	})

	_, err := handler.Handle(ctx, RefreshTokenCommand{
		RefreshToken: "test_token",
		IPAddress:    net.ParseIP("192.168.1.1"),
		UserAgent:    "Mozilla/5.0",
	})

	require.NoError(t, err)
	assert.NotEmpty(t, eventBus.PublishedEvents)
}
