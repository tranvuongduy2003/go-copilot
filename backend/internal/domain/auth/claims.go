package auth

import (
	"time"

	"github.com/google/uuid"
)

type Claims struct {
	UserID      uuid.UUID
	Email       string
	Roles       []string
	Permissions []string
	TokenID     string
	IssuedAt    time.Time
	ExpiresAt   time.Time
	Issuer      string
	Audience    string
}

func NewClaims(
	userID uuid.UUID,
	email string,
	roles []string,
	permissions []string,
	tokenID string,
	issuedAt time.Time,
	expiresAt time.Time,
	issuer string,
	audience string,
) Claims {
	return Claims{
		UserID:      userID,
		Email:       email,
		Roles:       roles,
		Permissions: permissions,
		TokenID:     tokenID,
		IssuedAt:    issuedAt,
		ExpiresAt:   expiresAt,
		Issuer:      issuer,
		Audience:    audience,
	}
}

func (c Claims) IsExpired() bool {
	return time.Now().UTC().After(c.ExpiresAt)
}

func (c Claims) HasRole(role string) bool {
	for _, r := range c.Roles {
		if r == role {
			return true
		}
	}
	return false
}

func (c Claims) HasPermission(permission string) bool {
	for _, p := range c.Permissions {
		if p == permission {
			return true
		}
	}
	return false
}

func (c Claims) HasAnyPermission(permissions ...string) bool {
	permissionSet := make(map[string]bool)
	for _, p := range c.Permissions {
		permissionSet[p] = true
	}
	for _, p := range permissions {
		if permissionSet[p] {
			return true
		}
	}
	return false
}

func (c Claims) HasAllPermissions(permissions ...string) bool {
	permissionSet := make(map[string]bool)
	for _, p := range c.Permissions {
		permissionSet[p] = true
	}
	for _, p := range permissions {
		if !permissionSet[p] {
			return false
		}
	}
	return true
}

func (c Claims) HasAnyRole(roles ...string) bool {
	roleSet := make(map[string]bool)
	for _, r := range c.Roles {
		roleSet[r] = true
	}
	for _, r := range roles {
		if roleSet[r] {
			return true
		}
	}
	return false
}
