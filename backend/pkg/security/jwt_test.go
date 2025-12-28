package security

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/auth"
)

func TestJWTTokenGenerator_GenerateAccessToken(t *testing.T) {
	config := JWTConfig{
		SecretKey:       "test-secret-key-that-is-long-enough",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
		Issuer:          "test-issuer",
		Audience:        "test-audience",
	}

	generator := NewJWTTokenGenerator(config)

	userID := uuid.New()
	email := "test@example.com"
	roles := []string{"admin", "user"}
	permissions := []string{"users:read", "users:create"}

	accessToken, err := generator.GenerateAccessToken(userID, email, roles, permissions)

	require.NoError(t, err)
	assert.NotEmpty(t, accessToken.Token())
	assert.False(t, accessToken.IsExpired())
}

func TestJWTTokenGenerator_ParseAccessToken(t *testing.T) {
	config := JWTConfig{
		SecretKey:       "test-secret-key-that-is-long-enough",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
		Issuer:          "test-issuer",
		Audience:        "test-audience",
	}

	generator := NewJWTTokenGenerator(config)

	userID := uuid.New()
	email := "test@example.com"
	roles := []string{"admin", "user"}
	permissions := []string{"users:read", "users:create"}

	accessToken, err := generator.GenerateAccessToken(userID, email, roles, permissions)
	require.NoError(t, err)

	claims, err := generator.ParseAccessToken(accessToken.Token())

	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, roles, claims.Roles)
	assert.Equal(t, permissions, claims.Permissions)
	assert.Equal(t, "test-issuer", claims.Issuer)
	assert.Equal(t, "test-audience", claims.Audience)
	assert.NotEmpty(t, claims.TokenID)
}

func TestJWTTokenGenerator_ExpiredToken(t *testing.T) {
	config := JWTConfig{
		SecretKey:       "test-secret-key-that-is-long-enough",
		AccessTokenTTL:  -1 * time.Hour,
		RefreshTokenTTL: 7 * 24 * time.Hour,
		Issuer:          "test-issuer",
		Audience:        "test-audience",
	}

	generator := NewJWTTokenGenerator(config)

	userID := uuid.New()
	accessToken, err := generator.GenerateAccessToken(userID, "test@example.com", nil, nil)
	require.NoError(t, err)

	_, err = generator.ParseAccessToken(accessToken.Token())

	assert.ErrorIs(t, err, auth.ErrTokenExpired)
}

func TestJWTTokenGenerator_InvalidToken(t *testing.T) {
	config := JWTConfig{
		SecretKey:       "test-secret-key-that-is-long-enough",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
		Issuer:          "test-issuer",
		Audience:        "test-audience",
	}

	generator := NewJWTTokenGenerator(config)

	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "empty token",
			token: "",
		},
		{
			name:  "malformed token",
			token: "not.a.valid.token",
		},
		{
			name:  "random string",
			token: "randomstring",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := generator.ParseAccessToken(testCase.token)
			assert.ErrorIs(t, err, auth.ErrTokenInvalid)
		})
	}
}

func TestJWTTokenGenerator_WrongSecret(t *testing.T) {
	config1 := JWTConfig{
		SecretKey:       "secret-key-1",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
		Issuer:          "test-issuer",
		Audience:        "test-audience",
	}

	config2 := JWTConfig{
		SecretKey:       "secret-key-2",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
		Issuer:          "test-issuer",
		Audience:        "test-audience",
	}

	generator1 := NewJWTTokenGenerator(config1)
	generator2 := NewJWTTokenGenerator(config2)

	userID := uuid.New()
	accessToken, err := generator1.GenerateAccessToken(userID, "test@example.com", nil, nil)
	require.NoError(t, err)

	_, err = generator2.ParseAccessToken(accessToken.Token())

	assert.ErrorIs(t, err, auth.ErrTokenInvalid)
}

func TestJWTTokenGenerator_GenerateRefreshToken(t *testing.T) {
	config := JWTConfig{
		SecretKey:       "test-secret-key-that-is-long-enough",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
		Issuer:          "test-issuer",
		Audience:        "test-audience",
	}

	generator := NewJWTTokenGenerator(config)

	token1, err := generator.GenerateRefreshToken()
	require.NoError(t, err)
	assert.NotEmpty(t, token1)

	token2, err := generator.GenerateRefreshToken()
	require.NoError(t, err)
	assert.NotEmpty(t, token2)

	assert.NotEqual(t, token1, token2)
}

func TestJWTTokenGenerator_HashRefreshToken(t *testing.T) {
	config := JWTConfig{
		SecretKey:       "test-secret-key-that-is-long-enough",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
		Issuer:          "test-issuer",
		Audience:        "test-audience",
	}

	generator := NewJWTTokenGenerator(config).(*jwtTokenGenerator)

	token := "test-refresh-token"
	hash1 := generator.HashRefreshToken(token)
	hash2 := generator.HashRefreshToken(token)

	assert.Equal(t, hash1, hash2)

	differentToken := "different-token"
	differentHash := generator.HashRefreshToken(differentToken)
	assert.NotEqual(t, hash1, differentHash)
}

func TestJWTTokenGenerator_TokenContainsUniqueID(t *testing.T) {
	config := JWTConfig{
		SecretKey:       "test-secret-key-that-is-long-enough",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
		Issuer:          "test-issuer",
		Audience:        "test-audience",
	}

	generator := NewJWTTokenGenerator(config)

	userID := uuid.New()

	token1, _ := generator.GenerateAccessToken(userID, "test@example.com", nil, nil)
	token2, _ := generator.GenerateAccessToken(userID, "test@example.com", nil, nil)

	claims1, _ := generator.ParseAccessToken(token1.Token())
	claims2, _ := generator.ParseAccessToken(token2.Token())

	assert.NotEqual(t, claims1.TokenID, claims2.TokenID)
}

func TestJWTTokenGenerator_ClaimsExpiration(t *testing.T) {
	config := JWTConfig{
		SecretKey:       "test-secret-key-that-is-long-enough",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
		Issuer:          "test-issuer",
		Audience:        "test-audience",
	}

	generator := NewJWTTokenGenerator(config)

	userID := uuid.New()
	accessToken, _ := generator.GenerateAccessToken(userID, "test@example.com", nil, nil)
	claims, _ := generator.ParseAccessToken(accessToken.Token())

	assert.False(t, claims.IsExpired())

	expectedExpiration := time.Now().UTC().Add(15 * time.Minute)
	assert.WithinDuration(t, expectedExpiration, claims.ExpiresAt, 2*time.Second)
}
