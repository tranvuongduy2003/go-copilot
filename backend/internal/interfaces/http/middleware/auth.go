package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/tranvuongduy2003/go-copilot/internal/interfaces/http/response"
)

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
