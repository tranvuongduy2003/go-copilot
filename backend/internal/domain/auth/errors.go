package auth

import (
	"github.com/tranvuongduy2003/go-copilot/internal/domain/shared"
)

var (
	ErrInvalidCredentials = shared.NewAuthorizationError("authenticate", "user")

	ErrAccountLocked = shared.NewBusinessRuleViolationError(
		"account_locked",
		"account is temporarily locked due to too many failed login attempts",
	)

	ErrAccountInactive = shared.NewBusinessRuleViolationError(
		"account_inactive",
		"account is not active",
	)

	ErrAccountBanned = shared.NewBusinessRuleViolationError(
		"account_banned",
		"account has been banned",
	)

	ErrMissingToken = shared.NewAuthorizationError("authenticate", "missing token")

	ErrInvalidTokenFormat = shared.NewAuthorizationError("authenticate", "invalid token format")

	ErrTokenExpired = shared.NewAuthorizationError("validate", "token (expired)")

	ErrTokenInvalid = shared.NewAuthorizationError("validate", "token (invalid)")

	ErrTokenRevoked = shared.NewAuthorizationError("validate", "token (revoked)")

	ErrRefreshTokenNotFound = shared.NewNotFoundError("RefreshToken", "")

	ErrRefreshTokenExpired = shared.NewBusinessRuleViolationError(
		"refresh_token_expired",
		"refresh token has expired",
	)

	ErrRefreshTokenRevoked = shared.NewBusinessRuleViolationError(
		"refresh_token_revoked",
		"refresh token has been revoked",
	)

	ErrRefreshTokenInvalid = shared.NewAuthorizationError("validate", "refresh token (invalid)")

	ErrInvalidResetToken = shared.NewAuthorizationError("validate", "password reset token (invalid)")

	ErrSessionLimitExceeded = shared.NewBusinessRuleViolationError(
		"session_limit_exceeded",
		"maximum number of active sessions exceeded",
	)

	ErrPasswordTooWeak = shared.NewValidationError(
		"password",
		"password does not meet security requirements",
	)

	ErrPasswordResetTokenInvalid = shared.NewAuthorizationError("validate", "password reset token")

	ErrPasswordResetTokenExpired = shared.NewBusinessRuleViolationError(
		"password_reset_token_expired",
		"password reset token has expired",
	)

	ErrEmailAlreadyVerified = shared.NewBusinessRuleViolationError(
		"email_already_verified",
		"email is already verified",
	)

	ErrSessionNotFound = shared.NewNotFoundError("Session", "")
)

func NewRefreshTokenNotFoundError(identifier string) *shared.NotFoundError {
	return shared.NewNotFoundError("RefreshToken", identifier)
}
