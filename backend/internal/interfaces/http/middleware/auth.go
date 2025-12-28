package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/auth"
	"github.com/tranvuongduy2003/go-copilot/internal/interfaces/http/response"
)

type authContextKey struct{}

type AuthContext struct {
	UserID      uuid.UUID
	Email       string
	Roles       []string
	Permissions []string
	TokenID     string
	ExpiresAt   time.Time
}

type AuthMiddleware struct {
	tokenGenerator auth.TokenGenerator
	tokenBlacklist auth.TokenBlacklist
}

func NewAuthMiddleware(tokenGenerator auth.TokenGenerator, tokenBlacklist auth.TokenBlacklist) *AuthMiddleware {
	return &AuthMiddleware{
		tokenGenerator: tokenGenerator,
		tokenBlacklist: tokenBlacklist,
	}
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		token, err := m.extractToken(request)
		if err != nil {
			response.Unauthorized(writer, request, err.Error())
			return
		}

		claims, err := m.tokenGenerator.ParseAccessToken(token)
		if err != nil {
			response.Unauthorized(writer, request, "invalid token")
			return
		}

		if claims.IsExpired() {
			response.Unauthorized(writer, request, "token expired")
			return
		}

		if m.tokenBlacklist != nil {
			isBlacklisted, err := m.tokenBlacklist.IsBlacklisted(request.Context(), claims.TokenID)
			if err == nil && isBlacklisted {
				response.Unauthorized(writer, request, "token has been revoked")
				return
			}
		}

		authContext := &AuthContext{
			UserID:      claims.UserID,
			Email:       claims.Email,
			Roles:       claims.Roles,
			Permissions: claims.Permissions,
			TokenID:     claims.TokenID,
			ExpiresAt:   claims.ExpiresAt,
		}

		ctx := context.WithValue(request.Context(), authContextKey{}, authContext)
		request = request.WithContext(ctx)

		next.ServeHTTP(writer, request)
	})
}

func (m *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		token, err := m.extractToken(request)
		if err != nil {
			next.ServeHTTP(writer, request)
			return
		}

		claims, err := m.tokenGenerator.ParseAccessToken(token)
		if err != nil {
			next.ServeHTTP(writer, request)
			return
		}

		if claims.IsExpired() {
			next.ServeHTTP(writer, request)
			return
		}

		if m.tokenBlacklist != nil {
			isBlacklisted, _ := m.tokenBlacklist.IsBlacklisted(request.Context(), claims.TokenID)
			if isBlacklisted {
				next.ServeHTTP(writer, request)
				return
			}
		}

		authContext := &AuthContext{
			UserID:      claims.UserID,
			Email:       claims.Email,
			Roles:       claims.Roles,
			Permissions: claims.Permissions,
			TokenID:     claims.TokenID,
			ExpiresAt:   claims.ExpiresAt,
		}

		ctx := context.WithValue(request.Context(), authContextKey{}, authContext)
		request = request.WithContext(ctx)

		next.ServeHTTP(writer, request)
	})
}

func (m *AuthMiddleware) extractToken(request *http.Request) (string, error) {
	authHeader := request.Header.Get("Authorization")
	if authHeader == "" {
		return "", auth.ErrMissingToken
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", auth.ErrInvalidTokenFormat
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", auth.ErrMissingToken
	}

	return token, nil
}

func GetAuthContext(ctx context.Context) (*AuthContext, bool) {
	if ctx == nil {
		return nil, false
	}
	authContext, ok := ctx.Value(authContextKey{}).(*AuthContext)
	return authContext, ok
}

func RequireAuthContext(ctx context.Context) *AuthContext {
	authContext, ok := GetAuthContext(ctx)
	if !ok {
		return nil
	}
	return authContext
}

type userContextKey struct{}

type AuthenticatedUser struct {
	ID    uuid.UUID
	Email string
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		authHeader := request.Header.Get("Authorization")
		if authHeader == "" {
			response.Unauthorized(writer, request, "missing authorization header")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Unauthorized(writer, request, "invalid authorization header format")
			return
		}

		token := parts[1]
		if token == "" {
			response.Unauthorized(writer, request, "missing token")
			return
		}

		authenticatedUser, err := validateToken(token)
		if err != nil {
			response.Unauthorized(writer, request, "invalid token")
			return
		}

		contextWithUser := context.WithValue(request.Context(), userContextKey{}, authenticatedUser)
		request = request.WithContext(contextWithUser)

		next.ServeHTTP(writer, request)
	})
}

func validateToken(token string) (*AuthenticatedUser, error) {
	_ = token
	return nil, nil
}

func GetAuthenticatedUser(ctx context.Context) *AuthenticatedUser {
	if ctx == nil {
		return nil
	}
	if authenticatedUser, ok := ctx.Value(userContextKey{}).(*AuthenticatedUser); ok {
		return authenticatedUser
	}
	return nil
}

func OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		authHeader := request.Header.Get("Authorization")
		if authHeader == "" {
			next.ServeHTTP(writer, request)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			next.ServeHTTP(writer, request)
			return
		}

		token := parts[1]
		if token == "" {
			next.ServeHTTP(writer, request)
			return
		}

		authenticatedUser, err := validateToken(token)
		if err == nil && authenticatedUser != nil {
			contextWithUser := context.WithValue(request.Context(), userContextKey{}, authenticatedUser)
			request = request.WithContext(contextWithUser)
		}

		next.ServeHTTP(writer, request)
	})
}
