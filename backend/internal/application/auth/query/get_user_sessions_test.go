package authquery

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

func TestGetUserSessionsHandler_Handle(t *testing.T) {
	ctx := context.Background()

	createRefreshToken := func(userID uuid.UUID) *auth.RefreshToken {
		now := time.Now().UTC()
		return auth.ReconstructRefreshToken(auth.ReconstructRefreshTokenParams{
			ID:        uuid.New(),
			UserID:    userID,
			TokenHash: "hash_" + uuid.New().String(),
			ExpiresAt: now.Add(24 * time.Hour),
			CreatedAt: now,
			IsRevoked: false,
		})
	}

	tests := []struct {
		name           string
		setupMocks     func(*testutil.MockRefreshTokenRepository) (uuid.UUID, uuid.UUID)
		wantErr        bool
		errContains    string
		expectedCount  int
	}{
		{
			name: "successfully get user sessions",
			setupMocks: func(tokenRepo *testutil.MockRefreshTokenRepository) (uuid.UUID, uuid.UUID) {
				userID := uuid.New()
				token1 := createRefreshToken(userID)
				token2 := createRefreshToken(userID)
				tokenRepo.Tokens[token1.ID()] = token1
				tokenRepo.Tokens[token2.ID()] = token2
				tokenRepo.UserTokens[userID] = []*auth.RefreshToken{token1, token2}
				return userID, token1.ID()
			},
			wantErr:       false,
			expectedCount: 2,
		},
		{
			name: "return empty list when no sessions",
			setupMocks: func(tokenRepo *testutil.MockRefreshTokenRepository) (uuid.UUID, uuid.UUID) {
				userID := uuid.New()
				return userID, uuid.Nil
			},
			wantErr:       false,
			expectedCount: 0,
		},
		{
			name: "successfully identify current session",
			setupMocks: func(tokenRepo *testutil.MockRefreshTokenRepository) (uuid.UUID, uuid.UUID) {
				userID := uuid.New()
				token1 := createRefreshToken(userID)
				token2 := createRefreshToken(userID)
				tokenRepo.Tokens[token1.ID()] = token1
				tokenRepo.Tokens[token2.ID()] = token2
				tokenRepo.UserTokens[userID] = []*auth.RefreshToken{token1, token2}
				return userID, token1.ID()
			},
			wantErr:       false,
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenRepo := testutil.NewMockRefreshTokenRepository()
			logger := testutil.NewNoopLogger()

			userID, currentTokenID := tt.setupMocks(tokenRepo)

			handler := NewGetUserSessionsHandler(GetUserSessionsHandlerParams{
				RefreshTokenRepository: tokenRepo,
				Logger:                 logger,
			})

			result, err := handler.Handle(ctx, GetUserSessionsQuery{
				UserID:         userID,
				CurrentTokenID: currentTokenID,
			})

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				require.NoError(t, err)
				assert.Len(t, result, tt.expectedCount)
			}
		})
	}
}

func TestGetUserSessionsHandler_Handle_MarkCurrentSession(t *testing.T) {
	ctx := context.Background()
	tokenRepo := testutil.NewMockRefreshTokenRepository()
	logger := testutil.NewNoopLogger()

	userID := uuid.New()
	now := time.Now().UTC()

	token1 := auth.ReconstructRefreshToken(auth.ReconstructRefreshTokenParams{
		ID:        uuid.New(),
		UserID:    userID,
		TokenHash: "hash1",
		ExpiresAt: now.Add(24 * time.Hour),
		CreatedAt: now,
		IsRevoked: false,
	})
	token2 := auth.ReconstructRefreshToken(auth.ReconstructRefreshTokenParams{
		ID:        uuid.New(),
		UserID:    userID,
		TokenHash: "hash2",
		ExpiresAt: now.Add(24 * time.Hour),
		CreatedAt: now.Add(-1 * time.Hour),
		IsRevoked: false,
	})

	tokenRepo.Tokens[token1.ID()] = token1
	tokenRepo.Tokens[token2.ID()] = token2
	tokenRepo.UserTokens[userID] = []*auth.RefreshToken{token1, token2}

	handler := NewGetUserSessionsHandler(GetUserSessionsHandlerParams{
		RefreshTokenRepository: tokenRepo,
		Logger:                 logger,
	})

	result, err := handler.Handle(ctx, GetUserSessionsQuery{
		UserID:         userID,
		CurrentTokenID: token1.ID(),
	})

	require.NoError(t, err)
	require.Len(t, result, 2)

	var currentFound bool
	for _, session := range result {
		if session.IsCurrent {
			currentFound = true
			assert.Equal(t, token1.ID(), session.ID)
		}
	}
	assert.True(t, currentFound, "Current session should be marked")
}
