package auth

import (
	"encoding/json"
	"net"
	"time"

	"github.com/google/uuid"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/shared"
)

type DeviceInfo struct {
	UserAgent string `json:"user_agent,omitempty"`
	Platform  string `json:"platform,omitempty"`
	Browser   string `json:"browser,omitempty"`
}

func (d DeviceInfo) ToJSON() ([]byte, error) {
	return json.Marshal(d)
}

func DeviceInfoFromJSON(data []byte) (DeviceInfo, error) {
	var d DeviceInfo
	if len(data) == 0 {
		return d, nil
	}
	err := json.Unmarshal(data, &d)
	return d, err
}

type RefreshToken struct {
	shared.Entity
	userID     uuid.UUID
	tokenHash  string
	expiresAt  time.Time
	createdAt  time.Time
	lastUsedAt *time.Time
	isRevoked  bool
	deviceInfo *DeviceInfo
	ipAddress  net.IP
}

type NewRefreshTokenParams struct {
	UserID     uuid.UUID
	TokenHash  string
	ExpiresAt  time.Time
	DeviceInfo *DeviceInfo
	IPAddress  net.IP
}

func NewRefreshToken(params NewRefreshTokenParams) (*RefreshToken, error) {
	if params.UserID == uuid.Nil {
		return nil, shared.NewValidationError("user_id", "user ID is required")
	}
	if params.TokenHash == "" {
		return nil, shared.NewValidationError("token_hash", "token hash is required")
	}
	if params.ExpiresAt.Before(time.Now().UTC()) {
		return nil, shared.NewValidationError("expires_at", "expiration time must be in the future")
	}

	now := time.Now().UTC()
	return &RefreshToken{
		Entity:     shared.NewEntity(),
		userID:     params.UserID,
		tokenHash:  params.TokenHash,
		expiresAt:  params.ExpiresAt,
		createdAt:  now,
		lastUsedAt: nil,
		isRevoked:  false,
		deviceInfo: params.DeviceInfo,
		ipAddress:  params.IPAddress,
	}, nil
}

type ReconstructRefreshTokenParams struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	TokenHash  string
	ExpiresAt  time.Time
	CreatedAt  time.Time
	LastUsedAt *time.Time
	IsRevoked  bool
	DeviceInfo *DeviceInfo
	IPAddress  net.IP
}

func ReconstructRefreshToken(params ReconstructRefreshTokenParams) *RefreshToken {
	return &RefreshToken{
		Entity:     shared.NewEntityWithID(params.ID),
		userID:     params.UserID,
		tokenHash:  params.TokenHash,
		expiresAt:  params.ExpiresAt,
		createdAt:  params.CreatedAt,
		lastUsedAt: params.LastUsedAt,
		isRevoked:  params.IsRevoked,
		deviceInfo: params.DeviceInfo,
		ipAddress:  params.IPAddress,
	}
}

func (rt *RefreshToken) UserID() uuid.UUID {
	return rt.userID
}

func (rt *RefreshToken) TokenHash() string {
	return rt.tokenHash
}

func (rt *RefreshToken) ExpiresAt() time.Time {
	return rt.expiresAt
}

func (rt *RefreshToken) CreatedAt() time.Time {
	return rt.createdAt
}

func (rt *RefreshToken) LastUsedAt() *time.Time {
	return rt.lastUsedAt
}

func (rt *RefreshToken) IsRevoked() bool {
	return rt.isRevoked
}

func (rt *RefreshToken) DeviceInfo() *DeviceInfo {
	return rt.deviceInfo
}

func (rt *RefreshToken) IPAddress() net.IP {
	return rt.ipAddress
}

func (rt *RefreshToken) IsExpired() bool {
	return time.Now().UTC().After(rt.expiresAt)
}

func (rt *RefreshToken) IsValid() bool {
	return !rt.isRevoked && !rt.IsExpired()
}

func (rt *RefreshToken) Revoke() {
	rt.isRevoked = true
}

func (rt *RefreshToken) UpdateLastUsed() {
	now := time.Now().UTC()
	rt.lastUsedAt = &now
}
