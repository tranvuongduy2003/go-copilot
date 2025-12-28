package repository

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/auth"
	"github.com/tranvuongduy2003/go-copilot/internal/infrastructure/persistence/postgres"
)

const (
	queryInsertRefreshToken = `
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, created_at, last_used_at, is_revoked, device_info, ip_address)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	queryUpdateRefreshToken = `
		UPDATE refresh_tokens
		SET last_used_at = $2, is_revoked = $3
		WHERE id = $1`

	queryFindRefreshTokenByID = `
		SELECT id, user_id, token_hash, expires_at, created_at, last_used_at, is_revoked, device_info, ip_address
		FROM refresh_tokens
		WHERE id = $1`

	queryFindRefreshTokenByHash = `
		SELECT id, user_id, token_hash, expires_at, created_at, last_used_at, is_revoked, device_info, ip_address
		FROM refresh_tokens
		WHERE token_hash = $1`

	queryFindRefreshTokensByUserID = `
		SELECT id, user_id, token_hash, expires_at, created_at, last_used_at, is_revoked, device_info, ip_address
		FROM refresh_tokens
		WHERE user_id = $1
		ORDER BY created_at DESC`

	queryFindActiveRefreshTokensByUserID = `
		SELECT id, user_id, token_hash, expires_at, created_at, last_used_at, is_revoked, device_info, ip_address
		FROM refresh_tokens
		WHERE user_id = $1 AND is_revoked = FALSE AND expires_at > NOW()
		ORDER BY created_at DESC`

	queryRevokeRefreshToken = `
		UPDATE refresh_tokens SET is_revoked = TRUE WHERE id = $1`

	queryRevokeAllRefreshTokensByUserID = `
		UPDATE refresh_tokens SET is_revoked = TRUE WHERE user_id = $1 AND is_revoked = FALSE`

	queryDeleteExpiredRefreshTokens = `
		DELETE FROM refresh_tokens WHERE expires_at < NOW() AND is_revoked = TRUE`

	queryCountActiveRefreshTokensByUserID = `
		SELECT COUNT(*) FROM refresh_tokens WHERE user_id = $1 AND is_revoked = FALSE AND expires_at > NOW()`
)

type refreshTokenRow struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	TokenHash  string
	ExpiresAt  time.Time
	CreatedAt  time.Time
	LastUsedAt *time.Time
	IsRevoked  bool
	DeviceInfo []byte
	IPAddress  net.IP
}

func (r *refreshTokenRow) toDomain() *auth.RefreshToken {
	var deviceInfo *auth.DeviceInfo
	if r.DeviceInfo != nil {
		di, err := auth.DeviceInfoFromJSON(r.DeviceInfo)
		if err == nil {
			deviceInfo = &di
		}
	}

	return auth.ReconstructRefreshToken(auth.ReconstructRefreshTokenParams{
		ID:         r.ID,
		UserID:     r.UserID,
		TokenHash:  r.TokenHash,
		ExpiresAt:  r.ExpiresAt,
		CreatedAt:  r.CreatedAt,
		LastUsedAt: r.LastUsedAt,
		IsRevoked:  r.IsRevoked,
		DeviceInfo: deviceInfo,
		IPAddress:  r.IPAddress,
	})
}

type RefreshTokenRepository struct {
	pool *pgxpool.Pool
}

func NewRefreshTokenRepository(pool *pgxpool.Pool) *RefreshTokenRepository {
	return &RefreshTokenRepository{pool: pool}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, token *auth.RefreshToken) error {
	querier := postgres.GetQuerier(ctx, r.pool)

	var deviceInfoBytes []byte
	if token.DeviceInfo() != nil {
		var err error
		deviceInfoBytes, err = token.DeviceInfo().ToJSON()
		if err != nil {
			return postgres.NewDBError("serialize device info", err)
		}
	}

	var ipAddress interface{}
	if token.IPAddress() != nil {
		ipAddress = token.IPAddress().String()
	}

	_, err := querier.Exec(ctx, queryInsertRefreshToken,
		token.ID(),
		token.UserID(),
		token.TokenHash(),
		token.ExpiresAt(),
		token.CreatedAt(),
		token.LastUsedAt(),
		token.IsRevoked(),
		deviceInfoBytes,
		ipAddress,
	)
	if err != nil {
		return postgres.NewDBError("create refresh token", err)
	}

	return nil
}

func (r *RefreshTokenRepository) FindByID(ctx context.Context, id uuid.UUID) (*auth.RefreshToken, error) {
	querier := postgres.GetQuerier(ctx, r.pool)

	row := &refreshTokenRow{}
	err := querier.QueryRow(ctx, queryFindRefreshTokenByID, id).Scan(
		&row.ID,
		&row.UserID,
		&row.TokenHash,
		&row.ExpiresAt,
		&row.CreatedAt,
		&row.LastUsedAt,
		&row.IsRevoked,
		&row.DeviceInfo,
		&row.IPAddress,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, auth.ErrSessionNotFound
		}
		return nil, postgres.NewDBError("find refresh token by id", err)
	}

	return row.toDomain(), nil
}

func (r *RefreshTokenRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*auth.RefreshToken, error) {
	querier := postgres.GetQuerier(ctx, r.pool)

	row := &refreshTokenRow{}
	err := querier.QueryRow(ctx, queryFindRefreshTokenByHash, tokenHash).Scan(
		&row.ID,
		&row.UserID,
		&row.TokenHash,
		&row.ExpiresAt,
		&row.CreatedAt,
		&row.LastUsedAt,
		&row.IsRevoked,
		&row.DeviceInfo,
		&row.IPAddress,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, auth.NewRefreshTokenNotFoundError(tokenHash)
		}
		return nil, postgres.NewDBError("find refresh token by hash", err)
	}

	return row.toDomain(), nil
}

func (r *RefreshTokenRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*auth.RefreshToken, error) {
	querier := postgres.GetQuerier(ctx, r.pool)

	rows, err := querier.Query(ctx, queryFindRefreshTokensByUserID, userID)
	if err != nil {
		return nil, postgres.NewDBError("find refresh tokens by user id", err)
	}
	defer rows.Close()

	tokens := make([]*auth.RefreshToken, 0)
	for rows.Next() {
		row := &refreshTokenRow{}
		err := rows.Scan(
			&row.ID,
			&row.UserID,
			&row.TokenHash,
			&row.ExpiresAt,
			&row.CreatedAt,
			&row.LastUsedAt,
			&row.IsRevoked,
			&row.DeviceInfo,
			&row.IPAddress,
		)
		if err != nil {
			return nil, postgres.NewDBError("scan refresh token row", err)
		}
		tokens = append(tokens, row.toDomain())
	}

	if err := rows.Err(); err != nil {
		return nil, postgres.NewDBError("iterate refresh token rows", err)
	}

	return tokens, nil
}

func (r *RefreshTokenRepository) FindActiveByUserID(ctx context.Context, userID uuid.UUID) ([]*auth.RefreshToken, error) {
	querier := postgres.GetQuerier(ctx, r.pool)

	rows, err := querier.Query(ctx, queryFindActiveRefreshTokensByUserID, userID)
	if err != nil {
		return nil, postgres.NewDBError("find active refresh tokens by user id", err)
	}
	defer rows.Close()

	tokens := make([]*auth.RefreshToken, 0)
	for rows.Next() {
		row := &refreshTokenRow{}
		err := rows.Scan(
			&row.ID,
			&row.UserID,
			&row.TokenHash,
			&row.ExpiresAt,
			&row.CreatedAt,
			&row.LastUsedAt,
			&row.IsRevoked,
			&row.DeviceInfo,
			&row.IPAddress,
		)
		if err != nil {
			return nil, postgres.NewDBError("scan refresh token row", err)
		}
		tokens = append(tokens, row.toDomain())
	}

	if err := rows.Err(); err != nil {
		return nil, postgres.NewDBError("iterate refresh token rows", err)
	}

	return tokens, nil
}

func (r *RefreshTokenRepository) Update(ctx context.Context, token *auth.RefreshToken) error {
	querier := postgres.GetQuerier(ctx, r.pool)

	cmdTag, err := querier.Exec(ctx, queryUpdateRefreshToken,
		token.ID(),
		token.LastUsedAt(),
		token.IsRevoked(),
	)
	if err != nil {
		return postgres.NewDBError("update refresh token", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return auth.NewRefreshTokenNotFoundError(token.ID().String())
	}

	return nil
}

func (r *RefreshTokenRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	querier := postgres.GetQuerier(ctx, r.pool)

	cmdTag, err := querier.Exec(ctx, queryRevokeRefreshToken, id)
	if err != nil {
		return postgres.NewDBError("revoke refresh token", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return auth.NewRefreshTokenNotFoundError(id.String())
	}

	return nil
}

func (r *RefreshTokenRepository) RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error {
	querier := postgres.GetQuerier(ctx, r.pool)

	_, err := querier.Exec(ctx, queryRevokeAllRefreshTokensByUserID, userID)
	if err != nil {
		return postgres.NewDBError("revoke all refresh tokens by user id", err)
	}

	return nil
}

func (r *RefreshTokenRepository) DeleteExpired(ctx context.Context) (int64, error) {
	querier := postgres.GetQuerier(ctx, r.pool)

	cmdTag, err := querier.Exec(ctx, queryDeleteExpiredRefreshTokens)
	if err != nil {
		return 0, postgres.NewDBError("delete expired refresh tokens", err)
	}

	return cmdTag.RowsAffected(), nil
}

func (r *RefreshTokenRepository) CountActiveByUserID(ctx context.Context, userID uuid.UUID) (int, error) {
	querier := postgres.GetQuerier(ctx, r.pool)

	var count int
	err := querier.QueryRow(ctx, queryCountActiveRefreshTokensByUserID, userID).Scan(&count)
	if err != nil {
		return 0, postgres.NewDBError("count active refresh tokens by user id", err)
	}

	return count, nil
}
