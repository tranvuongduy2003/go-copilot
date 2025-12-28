package auth

import (
	"context"

	"github.com/google/uuid"
)

type RefreshTokenRepository interface {
	Create(context context.Context, token *RefreshToken) error
	FindByID(context context.Context, id uuid.UUID) (*RefreshToken, error)
	FindByTokenHash(context context.Context, tokenHash string) (*RefreshToken, error)
	FindByUserID(context context.Context, userID uuid.UUID) ([]*RefreshToken, error)
	FindActiveByUserID(context context.Context, userID uuid.UUID) ([]*RefreshToken, error)
	Update(context context.Context, token *RefreshToken) error
	Revoke(context context.Context, id uuid.UUID) error
	RevokeAllByUserID(context context.Context, userID uuid.UUID) error
	DeleteExpired(context context.Context) (int64, error)
	CountActiveByUserID(context context.Context, userID uuid.UUID) (int, error)
}
