package auth

import (
	"time"
)

type AccessToken struct {
	token     string
	expiresAt time.Time
	tokenType string
}

func NewAccessToken(token string, expiresAt time.Time) AccessToken {
	return AccessToken{
		token:     token,
		expiresAt: expiresAt,
		tokenType: "Bearer",
	}
}

func (t AccessToken) Token() string {
	return t.token
}

func (t AccessToken) ExpiresAt() time.Time {
	return t.expiresAt
}

func (t AccessToken) TokenType() string {
	return t.tokenType
}

func (t AccessToken) IsExpired() bool {
	return time.Now().UTC().After(t.expiresAt)
}

func (t AccessToken) ExpiresIn() int64 {
	return int64(time.Until(t.expiresAt).Seconds())
}

type RefreshTokenValue struct {
	token     string
	expiresAt time.Time
}

func NewRefreshTokenValue(token string, expiresAt time.Time) RefreshTokenValue {
	return RefreshTokenValue{
		token:     token,
		expiresAt: expiresAt,
	}
}

func (t RefreshTokenValue) Token() string {
	return t.token
}

func (t RefreshTokenValue) ExpiresAt() time.Time {
	return t.expiresAt
}

func (t RefreshTokenValue) IsExpired() bool {
	return time.Now().UTC().After(t.expiresAt)
}

type TokenPair struct {
	accessToken  AccessToken
	refreshToken RefreshTokenValue
}

func NewTokenPair(accessToken AccessToken, refreshToken RefreshTokenValue) TokenPair {
	return TokenPair{
		accessToken:  accessToken,
		refreshToken: refreshToken,
	}
}

func (tp TokenPair) AccessToken() AccessToken {
	return tp.accessToken
}

func (tp TokenPair) RefreshToken() RefreshTokenValue {
	return tp.refreshToken
}
