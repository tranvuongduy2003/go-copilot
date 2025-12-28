package auth

import (
	"context"

	"github.com/google/uuid"
)

type TokenGeneratorConfig struct {
	SecretKey             string
	AccessTokenExpiry     int
	RefreshTokenExpiry    int
	Issuer                string
	Audience              string
}

type TokenGenerator interface {
	GenerateAccessToken(userID uuid.UUID, email string, roles []string, permissions []string) (AccessToken, error)
	GenerateRefreshToken() (string, error)
	ParseAccessToken(token string) (*Claims, error)
	HashRefreshToken(token string) string
}

type TokenBlacklist interface {
	Add(context context.Context, tokenID string, expiresAt int64) error
	IsBlacklisted(context context.Context, tokenID string) (bool, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(hashedPassword, plainPassword string) (bool, error)
}
