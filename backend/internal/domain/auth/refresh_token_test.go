package auth

import (
	"net"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRefreshToken(t *testing.T) {
	validUserID := uuid.New()
	validExpiresAt := time.Now().UTC().Add(24 * time.Hour)

	tests := []struct {
		name        string
		params      NewRefreshTokenParams
		wantErr     bool
		errContains string
	}{
		{
			name: "valid refresh token",
			params: NewRefreshTokenParams{
				UserID:    validUserID,
				TokenHash: "hashed_token_value",
				ExpiresAt: validExpiresAt,
				DeviceInfo: &DeviceInfo{
					UserAgent: "Mozilla/5.0",
					Platform:  "Windows",
					Browser:   "Chrome",
				},
				IPAddress: net.ParseIP("192.168.1.1"),
			},
			wantErr: false,
		},
		{
			name: "valid without optional fields",
			params: NewRefreshTokenParams{
				UserID:    validUserID,
				TokenHash: "hashed_token_value",
				ExpiresAt: validExpiresAt,
			},
			wantErr: false,
		},
		{
			name: "missing user ID",
			params: NewRefreshTokenParams{
				UserID:    uuid.Nil,
				TokenHash: "hashed_token",
				ExpiresAt: validExpiresAt,
			},
			wantErr:     true,
			errContains: "user ID is required",
		},
		{
			name: "missing token hash",
			params: NewRefreshTokenParams{
				UserID:    validUserID,
				TokenHash: "",
				ExpiresAt: validExpiresAt,
			},
			wantErr:     true,
			errContains: "token hash is required",
		},
		{
			name: "expired token",
			params: NewRefreshTokenParams{
				UserID:    validUserID,
				TokenHash: "hashed_token",
				ExpiresAt: time.Now().UTC().Add(-1 * time.Hour),
			},
			wantErr:     true,
			errContains: "expiration time must be in the future",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := NewRefreshToken(tt.params)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				assert.Nil(t, token)
			} else {
				require.NoError(t, err)
				require.NotNil(t, token)
				assert.NotEqual(t, uuid.Nil, token.ID())
				assert.Equal(t, tt.params.UserID, token.UserID())
				assert.Equal(t, tt.params.TokenHash, token.TokenHash())
				assert.False(t, token.IsRevoked())
				assert.Nil(t, token.LastUsedAt())
			}
		})
	}
}

func TestReconstructRefreshToken(t *testing.T) {
	id := uuid.New()
	userID := uuid.New()
	now := time.Now().UTC()
	lastUsed := now.Add(-1 * time.Hour)

	token := ReconstructRefreshToken(ReconstructRefreshTokenParams{
		ID:         id,
		UserID:     userID,
		TokenHash:  "hashed_token",
		ExpiresAt:  now.Add(24 * time.Hour),
		CreatedAt:  now,
		LastUsedAt: &lastUsed,
		IsRevoked:  false,
		DeviceInfo: &DeviceInfo{UserAgent: "Test"},
		IPAddress:  net.ParseIP("10.0.0.1"),
	})

	assert.Equal(t, id, token.ID())
	assert.Equal(t, userID, token.UserID())
	assert.Equal(t, "hashed_token", token.TokenHash())
	assert.NotNil(t, token.LastUsedAt())
	assert.Equal(t, lastUsed, *token.LastUsedAt())
	assert.False(t, token.IsRevoked())
}

func TestRefreshToken_IsExpired(t *testing.T) {
	userID := uuid.New()

	futureToken, _ := NewRefreshToken(NewRefreshTokenParams{
		UserID:    userID,
		TokenHash: "hash1",
		ExpiresAt: time.Now().UTC().Add(24 * time.Hour),
	})
	assert.False(t, futureToken.IsExpired())

	expiredToken := ReconstructRefreshToken(ReconstructRefreshTokenParams{
		ID:        uuid.New(),
		UserID:    userID,
		TokenHash: "hash2",
		ExpiresAt: time.Now().UTC().Add(-1 * time.Hour),
		CreatedAt: time.Now().UTC().Add(-25 * time.Hour),
	})
	assert.True(t, expiredToken.IsExpired())
}

func TestRefreshToken_IsValid(t *testing.T) {
	userID := uuid.New()

	validToken, _ := NewRefreshToken(NewRefreshTokenParams{
		UserID:    userID,
		TokenHash: "hash",
		ExpiresAt: time.Now().UTC().Add(24 * time.Hour),
	})
	assert.True(t, validToken.IsValid())

	revokedToken := ReconstructRefreshToken(ReconstructRefreshTokenParams{
		ID:        uuid.New(),
		UserID:    userID,
		TokenHash: "hash",
		ExpiresAt: time.Now().UTC().Add(24 * time.Hour),
		CreatedAt: time.Now().UTC(),
		IsRevoked: true,
	})
	assert.False(t, revokedToken.IsValid())

	expiredToken := ReconstructRefreshToken(ReconstructRefreshTokenParams{
		ID:        uuid.New(),
		UserID:    userID,
		TokenHash: "hash",
		ExpiresAt: time.Now().UTC().Add(-1 * time.Hour),
		CreatedAt: time.Now().UTC().Add(-25 * time.Hour),
	})
	assert.False(t, expiredToken.IsValid())
}

func TestRefreshToken_Revoke(t *testing.T) {
	token, _ := NewRefreshToken(NewRefreshTokenParams{
		UserID:    uuid.New(),
		TokenHash: "hash",
		ExpiresAt: time.Now().UTC().Add(24 * time.Hour),
	})

	assert.False(t, token.IsRevoked())
	assert.True(t, token.IsValid())

	token.Revoke()

	assert.True(t, token.IsRevoked())
	assert.False(t, token.IsValid())
}

func TestRefreshToken_UpdateLastUsed(t *testing.T) {
	token, _ := NewRefreshToken(NewRefreshTokenParams{
		UserID:    uuid.New(),
		TokenHash: "hash",
		ExpiresAt: time.Now().UTC().Add(24 * time.Hour),
	})

	assert.Nil(t, token.LastUsedAt())

	token.UpdateLastUsed()

	require.NotNil(t, token.LastUsedAt())
	assert.WithinDuration(t, time.Now().UTC(), *token.LastUsedAt(), time.Second)
}

func TestDeviceInfo_ToJSON(t *testing.T) {
	info := DeviceInfo{
		UserAgent: "Mozilla/5.0",
		Platform:  "Windows",
		Browser:   "Chrome",
	}

	data, err := info.ToJSON()
	require.NoError(t, err)
	assert.Contains(t, string(data), "Mozilla/5.0")
	assert.Contains(t, string(data), "Windows")
	assert.Contains(t, string(data), "Chrome")
}

func TestDeviceInfoFromJSON(t *testing.T) {
	jsonData := []byte(`{"user_agent":"Mozilla/5.0","platform":"MacOS","browser":"Safari"}`)

	info, err := DeviceInfoFromJSON(jsonData)
	require.NoError(t, err)
	assert.Equal(t, "Mozilla/5.0", info.UserAgent)
	assert.Equal(t, "MacOS", info.Platform)
	assert.Equal(t, "Safari", info.Browser)

	emptyInfo, err := DeviceInfoFromJSON([]byte{})
	require.NoError(t, err)
	assert.Empty(t, emptyInfo.UserAgent)
}
