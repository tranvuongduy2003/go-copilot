package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/tranvuongduy2003/go-copilot/internal/interfaces/http/response"
)

func RequirePermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			authContext, ok := GetAuthContext(request.Context())
			if !ok {
				response.Unauthorized(writer, request, "unauthorized")
				return
			}

			if !hasPermission(authContext.Permissions, permission) {
				response.Forbidden(writer, request, "insufficient permissions")
				return
			}

			next.ServeHTTP(writer, request)
		})
	}
}

func RequireAnyPermission(permissions ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			authContext, ok := GetAuthContext(request.Context())
			if !ok {
				response.Unauthorized(writer, request, "unauthorized")
				return
			}

			for _, permission := range permissions {
				if hasPermission(authContext.Permissions, permission) {
					next.ServeHTTP(writer, request)
					return
				}
			}

			response.Forbidden(writer, request, "insufficient permissions")
		})
	}
}

func RequireAllPermissions(permissions ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			authContext, ok := GetAuthContext(request.Context())
			if !ok {
				response.Unauthorized(writer, request, "unauthorized")
				return
			}

			for _, permission := range permissions {
				if !hasPermission(authContext.Permissions, permission) {
					response.Forbidden(writer, request, "insufficient permissions")
					return
				}
			}

			next.ServeHTTP(writer, request)
		})
	}
}

func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			authContext, ok := GetAuthContext(request.Context())
			if !ok {
				response.Unauthorized(writer, request, "unauthorized")
				return
			}

			if !hasRole(authContext.Roles, role) {
				response.Forbidden(writer, request, "insufficient role")
				return
			}

			next.ServeHTTP(writer, request)
		})
	}
}

func RequireAnyRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			authContext, ok := GetAuthContext(request.Context())
			if !ok {
				response.Unauthorized(writer, request, "unauthorized")
				return
			}

			for _, role := range roles {
				if hasRole(authContext.Roles, role) {
					next.ServeHTTP(writer, request)
					return
				}
			}

			response.Forbidden(writer, request, "insufficient role")
		})
	}
}

func ResourceOwner(paramName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			authContext, ok := GetAuthContext(request.Context())
			if !ok {
				response.Unauthorized(writer, request, "unauthorized")
				return
			}

			resourceID := chi.URLParam(request, paramName)
			if resourceID == "" {
				response.BadRequest(writer, request, "resource id is required")
				return
			}

			if resourceID != authContext.UserID.String() {
				if hasRole(authContext.Roles, "super_admin") || hasRole(authContext.Roles, "admin") {
					next.ServeHTTP(writer, request)
					return
				}

				response.Forbidden(writer, request, "access denied")
				return
			}

			next.ServeHTTP(writer, request)
		})
	}
}

func hasPermission(userPermissions []string, permission string) bool {
	for _, userPermission := range userPermissions {
		if userPermission == permission || userPermission == "system:admin" {
			return true
		}
	}
	return false
}

func hasRole(userRoles []string, role string) bool {
	for _, userRole := range userRoles {
		if userRole == role {
			return true
		}
	}
	return false
}
