package authdto

import (
	"time"

	"github.com/google/uuid"

	userdto "github.com/tranvuongduy2003/go-copilot/internal/application/user/dto"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/auth"
)

type AuthResponseDTO struct {
	User         *userdto.UserDTO `json:"user"`
	AccessToken  string           `json:"access_token"`
	RefreshToken string           `json:"refresh_token"`
	ExpiresAt    time.Time        `json:"expires_at"`
}

type TokenPairDTO struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type SessionDTO struct {
	ID         uuid.UUID  `json:"id"`
	DeviceInfo DeviceDTO  `json:"device_info"`
	IPAddress  string     `json:"ip_address"`
	CreatedAt  time.Time  `json:"created_at"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	IsCurrent  bool       `json:"is_current"`
}

type DeviceDTO struct {
	UserAgent string `json:"user_agent,omitempty"`
	Platform  string `json:"platform,omitempty"`
	Browser   string `json:"browser,omitempty"`
}

func SessionFromRefreshToken(token *auth.RefreshToken, currentTokenID uuid.UUID) *SessionDTO {
	if token == nil {
		return nil
	}

	var deviceDTO DeviceDTO
	if deviceInfo := token.DeviceInfo(); deviceInfo != nil {
		deviceDTO = DeviceDTO{
			UserAgent: deviceInfo.UserAgent,
			Platform:  deviceInfo.Platform,
			Browser:   deviceInfo.Browser,
		}
	}

	var ipAddress string
	if token.IPAddress() != nil {
		ipAddress = token.IPAddress().String()
	}

	return &SessionDTO{
		ID:         token.ID(),
		DeviceInfo: deviceDTO,
		IPAddress:  ipAddress,
		CreatedAt:  token.CreatedAt(),
		LastUsedAt: token.LastUsedAt(),
		IsCurrent:  token.ID() == currentTokenID,
	}
}

func SessionsFromRefreshTokens(tokens []*auth.RefreshToken, currentTokenID uuid.UUID) []*SessionDTO {
	sessions := make([]*SessionDTO, 0, len(tokens))
	for _, token := range tokens {
		if token.IsValid() {
			sessions = append(sessions, SessionFromRefreshToken(token, currentTokenID))
		}
	}
	return sessions
}

type AuthUserDTO struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	FullName    string    `json:"full_name"`
	Status      string    `json:"status"`
	Roles       []string  `json:"roles"`
	Permissions []string  `json:"permissions"`
}

type ClaimsDTO struct {
	UserID      uuid.UUID `json:"user_id"`
	Email       string    `json:"email"`
	Roles       []string  `json:"roles"`
	Permissions []string  `json:"permissions"`
	TokenID     string    `json:"token_id"`
	IssuedAt    time.Time `json:"issued_at"`
	ExpiresAt   time.Time `json:"expires_at"`
}

func ClaimsFromDomain(claims *auth.Claims) *ClaimsDTO {
	if claims == nil {
		return nil
	}
	return &ClaimsDTO{
		UserID:      claims.UserID,
		Email:       claims.Email,
		Roles:       claims.Roles,
		Permissions: claims.Permissions,
		TokenID:     claims.TokenID,
		IssuedAt:    claims.IssuedAt,
		ExpiresAt:   claims.ExpiresAt,
	}
}
