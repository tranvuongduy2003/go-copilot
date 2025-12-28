package authcommand

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/auth"
	"github.com/tranvuongduy2003/go-copilot/pkg/testutil"
)

func TestLogoutHandler_Handle(t *testing.T) {
	ctx := context.Background()

	createRefreshToken := func(userID uuid.UUID) *auth.RefreshToken {
		token, _ := auth.NewRefreshToken(auth.NewRefreshTokenParams{
			UserID:    userID,
			TokenHash: "test_hash",
			ExpiresAt: time.Now().UTC().Add(24 * time.Hour),
		})
		return token
	}

	tests := []struct {
		name        string
		setupMocks  func(*testutil.MockRefreshTokenRepository, *testutil.MockTokenBlacklist) uuid.UUID
		command     func(uuid.UUID) LogoutCommand
		wantErr     bool
		errContains string
		checkResult func(*testing.T, *testutil.MockRefreshTokenRepository, *testutil.MockTokenBlacklist)
	}{
		{
			name: "successfully logout single session",
			setupMocks: func(tokenRepo *testutil.MockRefreshTokenRepository, blacklist *testutil.MockTokenBlacklist) uuid.UUID {
				userID := uuid.New()
				token := createRefreshToken(userID)
				tokenRepo.Tokens[token.ID()] = token
				tokenRepo.UserTokens[userID] = append(tokenRepo.UserTokens[userID], token)
				return userID
			},
			command: func(userID uuid.UUID) LogoutCommand {
				return LogoutCommand{
					UserID:    userID,
					TokenID:   "test_token_id",
					ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
					LogoutAll: false,
				}
			},
			wantErr: false,
			checkResult: func(t *testing.T, tokenRepo *testutil.MockRefreshTokenRepository, blacklist *testutil.MockTokenBlacklist) {
				assert.True(t, blacklist.BlacklistedTokens["test_token_id"])
			},
		},
		{
			name: "successfully logout all sessions",
			setupMocks: func(tokenRepo *testutil.MockRefreshTokenRepository, blacklist *testutil.MockTokenBlacklist) uuid.UUID {
				userID := uuid.New()
				token1 := createRefreshToken(userID)
				token2 := createRefreshToken(userID)
				tokenRepo.Tokens[token1.ID()] = token1
				tokenRepo.Tokens[token2.ID()] = token2
				tokenRepo.UserTokens[userID] = append(tokenRepo.UserTokens[userID], token1, token2)
				return userID
			},
			command: func(userID uuid.UUID) LogoutCommand {
				return LogoutCommand{
					UserID:    userID,
					TokenID:   "test_token_id",
					ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
					LogoutAll: true,
				}
			},
			wantErr: false,
			checkResult: func(t *testing.T, tokenRepo *testutil.MockRefreshTokenRepository, blacklist *testutil.MockTokenBlacklist) {
				assert.True(t, blacklist.BlacklistedTokens["test_token_id"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenRepo := testutil.NewMockRefreshTokenRepository()
			blacklist := testutil.NewMockTokenBlacklist()
			eventBus := testutil.NewMockEventBus()
			logger := testutil.NewNoopLogger()

			userID := tt.setupMocks(tokenRepo, blacklist)

			handler := NewLogoutHandler(LogoutHandlerParams{
				RefreshTokenRepository: tokenRepo,
				TokenBlacklist:         blacklist,
				EventBus:               eventBus,
				Logger:                 logger,
			})

			err := handler.Handle(ctx, tt.command(userID))

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				require.NoError(t, err)
				if tt.checkResult != nil {
					tt.checkResult(t, tokenRepo, blacklist)
				}
			}
		})
	}
}

func TestLogoutHandler_Handle_PublishesEvent(t *testing.T) {
	ctx := context.Background()
	tokenRepo := testutil.NewMockRefreshTokenRepository()
	blacklist := testutil.NewMockTokenBlacklist()
	eventBus := testutil.NewMockEventBus()
	logger := testutil.NewNoopLogger()

	userID := uuid.New()

	handler := NewLogoutHandler(LogoutHandlerParams{
		RefreshTokenRepository: tokenRepo,
		TokenBlacklist:         blacklist,
		EventBus:               eventBus,
		Logger:                 logger,
	})

	err := handler.Handle(ctx, LogoutCommand{
		UserID:    userID,
		TokenID:   "test_token_id",
		ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
		LogoutAll: false,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, eventBus.PublishedEvents)
}

func TestLogoutHandler_Handle_LogoutAllRevokesAllTokens(t *testing.T) {
	ctx := context.Background()
	tokenRepo := testutil.NewMockRefreshTokenRepository()
	blacklist := testutil.NewMockTokenBlacklist()
	eventBus := testutil.NewMockEventBus()
	logger := testutil.NewNoopLogger()

	userID := uuid.New()
	token1, _ := auth.NewRefreshToken(auth.NewRefreshTokenParams{
		UserID:    userID,
		TokenHash: "hash1",
		ExpiresAt: time.Now().UTC().Add(24 * time.Hour),
	})
	token2, _ := auth.NewRefreshToken(auth.NewRefreshTokenParams{
		UserID:    userID,
		TokenHash: "hash2",
		ExpiresAt: time.Now().UTC().Add(24 * time.Hour),
	})

	tokenRepo.Tokens[token1.ID()] = token1
	tokenRepo.Tokens[token2.ID()] = token2
	tokenRepo.UserTokens[userID] = []*auth.RefreshToken{token1, token2}

	handler := NewLogoutHandler(LogoutHandlerParams{
		RefreshTokenRepository: tokenRepo,
		TokenBlacklist:         blacklist,
		EventBus:               eventBus,
		Logger:                 logger,
	})

	err := handler.Handle(ctx, LogoutCommand{
		UserID:    userID,
		TokenID:   "test_token_id",
		ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
		LogoutAll: true,
	})

	require.NoError(t, err)

	for _, token := range tokenRepo.UserTokens[userID] {
		assert.True(t, token.IsRevoked())
	}
}
