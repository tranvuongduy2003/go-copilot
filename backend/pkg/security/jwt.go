package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/auth"
)

type JWTConfig struct {
	SecretKey       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	Issuer          string
	Audience        string
}

type jwtTokenGenerator struct {
	config JWTConfig
}

func NewJWTTokenGenerator(config JWTConfig) auth.TokenGenerator {
	return &jwtTokenGenerator{config: config}
}

type jwtClaims struct {
	jwt.RegisteredClaims
	Email       string   `json:"email"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}

func (generator *jwtTokenGenerator) GenerateAccessToken(userID uuid.UUID, email string, roles []string, permissions []string) (auth.AccessToken, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(generator.config.AccessTokenTTL)
	tokenID := uuid.New().String()

	claims := jwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID,
			Subject:   userID.String(),
			Issuer:    generator.config.Issuer,
			Audience:  jwt.ClaimStrings{generator.config.Audience},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
		},
		Email:       email,
		Roles:       roles,
		Permissions: permissions,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(generator.config.SecretKey))
	if err != nil {
		return auth.AccessToken{}, fmt.Errorf("failed to sign access token: %w", err)
	}

	return auth.NewAccessToken(tokenString, expiresAt), nil
}

func (generator *jwtTokenGenerator) GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func (generator *jwtTokenGenerator) ParseAccessToken(tokenString string) (*auth.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(generator.config.SecretKey), nil
	})
	if err != nil {
		return nil, auth.ErrTokenInvalid
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok || !token.Valid {
		return nil, auth.ErrTokenInvalid
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, auth.ErrTokenInvalid
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now().UTC()) {
		return nil, auth.ErrTokenExpired
	}

	issuedAt := time.Time{}
	if claims.IssuedAt != nil {
		issuedAt = claims.IssuedAt.Time
	}

	expiresAt := time.Time{}
	if claims.ExpiresAt != nil {
		expiresAt = claims.ExpiresAt.Time
	}

	audience := ""
	if len(claims.Audience) > 0 {
		audience = claims.Audience[0]
	}

	return &auth.Claims{
		UserID:      userID,
		Email:       claims.Email,
		Roles:       claims.Roles,
		Permissions: claims.Permissions,
		TokenID:     claims.ID,
		IssuedAt:    issuedAt,
		ExpiresAt:   expiresAt,
		Issuer:      claims.Issuer,
		Audience:    audience,
	}, nil
}

func (generator *jwtTokenGenerator) HashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
