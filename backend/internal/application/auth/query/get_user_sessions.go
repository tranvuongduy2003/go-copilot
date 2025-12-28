package authquery

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	authdto "github.com/tranvuongduy2003/go-copilot/internal/application/auth/dto"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/auth"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

type GetUserSessionsQuery struct {
	UserID         uuid.UUID
	CurrentTokenID uuid.UUID
}

type GetUserSessionsHandler struct {
	refreshTokenRepository auth.RefreshTokenRepository
	logger                 logger.Logger
}

type GetUserSessionsHandlerParams struct {
	RefreshTokenRepository auth.RefreshTokenRepository
	Logger                 logger.Logger
}

func NewGetUserSessionsHandler(params GetUserSessionsHandlerParams) *GetUserSessionsHandler {
	return &GetUserSessionsHandler{
		refreshTokenRepository: params.RefreshTokenRepository,
		logger:                 params.Logger,
	}
}

func (handler *GetUserSessionsHandler) Handle(ctx context.Context, query GetUserSessionsQuery) ([]*authdto.SessionDTO, error) {
	tokens, err := handler.refreshTokenRepository.FindActiveByUserID(ctx, query.UserID)
	if err != nil {
		return nil, fmt.Errorf("find active refresh tokens: %w", err)
	}

	sessions := authdto.SessionsFromRefreshTokens(tokens, query.CurrentTokenID)

	handler.logger.Debug("user sessions retrieved",
		logger.String("user_id", query.UserID.String()),
		logger.Int("session_count", len(sessions)),
	)

	return sessions, nil
}
