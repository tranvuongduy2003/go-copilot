package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestHasPermission(t *testing.T) {
	tests := []struct {
		name            string
		userPermissions []string
		permission      string
		expected        bool
	}{
		{
			name:            "user has exact permission",
			userPermissions: []string{"users:read", "users:create"},
			permission:      "users:read",
			expected:        true,
		},
		{
			name:            "user does not have permission",
			userPermissions: []string{"users:read"},
			permission:      "users:create",
			expected:        false,
		},
		{
			name:            "system:admin grants any permission",
			userPermissions: []string{"system:admin"},
			permission:      "any:permission",
			expected:        true,
		},
		{
			name:            "empty permissions list",
			userPermissions: []string{},
			permission:      "users:read",
			expected:        false,
		},
		{
			name:            "nil permissions list",
			userPermissions: nil,
			permission:      "users:read",
			expected:        false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			result := hasPermission(testCase.userPermissions, testCase.permission)
			assert.Equal(t, testCase.expected, result)
		})
	}
}

func TestHasRole(t *testing.T) {
	tests := []struct {
		name      string
		userRoles []string
		role      string
		expected  bool
	}{
		{
			name:      "user has exact role",
			userRoles: []string{"admin", "user"},
			role:      "admin",
			expected:  true,
		},
		{
			name:      "user does not have role",
			userRoles: []string{"user"},
			role:      "admin",
			expected:  false,
		},
		{
			name:      "empty roles list",
			userRoles: []string{},
			role:      "admin",
			expected:  false,
		},
		{
			name:      "nil roles list",
			userRoles: nil,
			role:      "admin",
			expected:  false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			result := hasRole(testCase.userRoles, testCase.role)
			assert.Equal(t, testCase.expected, result)
		})
	}
}

func TestRequirePermission(t *testing.T) {
	tests := []struct {
		name               string
		setupContext       func() context.Context
		requiredPermission string
		expectedStatus     int
	}{
		{
			name: "allow when user has required permission",
			setupContext: func() context.Context {
				authContext := &AuthContext{
					UserID:      uuid.New(),
					Permissions: []string{"users:read", "users:create"},
				}
				return context.WithValue(context.Background(), authContextKey{}, authContext)
			},
			requiredPermission: "users:read",
			expectedStatus:     http.StatusOK,
		},
		{
			name: "allow when user has system:admin permission",
			setupContext: func() context.Context {
				authContext := &AuthContext{
					UserID:      uuid.New(),
					Permissions: []string{"system:admin"},
				}
				return context.WithValue(context.Background(), authContextKey{}, authContext)
			},
			requiredPermission: "any:permission",
			expectedStatus:     http.StatusOK,
		},
		{
			name: "deny when user lacks required permission",
			setupContext: func() context.Context {
				authContext := &AuthContext{
					UserID:      uuid.New(),
					Permissions: []string{"users:read"},
				}
				return context.WithValue(context.Background(), authContextKey{}, authContext)
			},
			requiredPermission: "users:create",
			expectedStatus:     http.StatusForbidden,
		},
		{
			name: "deny when no auth context",
			setupContext: func() context.Context {
				return context.Background()
			},
			requiredPermission: "users:read",
			expectedStatus:     http.StatusUnauthorized,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			nextHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(http.StatusOK)
			})

			middleware := RequirePermission(testCase.requiredPermission)
			handler := middleware(nextHandler)

			request := httptest.NewRequest(http.MethodGet, "/test", nil)
			request = request.WithContext(testCase.setupContext())
			recorder := httptest.NewRecorder()

			handler.ServeHTTP(recorder, request)

			assert.Equal(t, testCase.expectedStatus, recorder.Code)
		})
	}
}

func TestRequireAnyPermission(t *testing.T) {
	tests := []struct {
		name                string
		setupContext        func() context.Context
		requiredPermissions []string
		expectedStatus      int
	}{
		{
			name: "allow when user has one of required permissions",
			setupContext: func() context.Context {
				authContext := &AuthContext{
					UserID:      uuid.New(),
					Permissions: []string{"users:read"},
				}
				return context.WithValue(context.Background(), authContextKey{}, authContext)
			},
			requiredPermissions: []string{"users:read", "users:create"},
			expectedStatus:      http.StatusOK,
		},
		{
			name: "allow when user has multiple required permissions",
			setupContext: func() context.Context {
				authContext := &AuthContext{
					UserID:      uuid.New(),
					Permissions: []string{"users:read", "users:create"},
				}
				return context.WithValue(context.Background(), authContextKey{}, authContext)
			},
			requiredPermissions: []string{"users:read", "users:create"},
			expectedStatus:      http.StatusOK,
		},
		{
			name: "deny when user has none of required permissions",
			setupContext: func() context.Context {
				authContext := &AuthContext{
					UserID:      uuid.New(),
					Permissions: []string{"users:delete"},
				}
				return context.WithValue(context.Background(), authContextKey{}, authContext)
			},
			requiredPermissions: []string{"users:read", "users:create"},
			expectedStatus:      http.StatusForbidden,
		},
		{
			name: "deny when no auth context",
			setupContext: func() context.Context {
				return context.Background()
			},
			requiredPermissions: []string{"users:read"},
			expectedStatus:      http.StatusUnauthorized,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			nextHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(http.StatusOK)
			})

			middleware := RequireAnyPermission(testCase.requiredPermissions...)
			handler := middleware(nextHandler)

			request := httptest.NewRequest(http.MethodGet, "/test", nil)
			request = request.WithContext(testCase.setupContext())
			recorder := httptest.NewRecorder()

			handler.ServeHTTP(recorder, request)

			assert.Equal(t, testCase.expectedStatus, recorder.Code)
		})
	}
}

func TestRequireAllPermissions(t *testing.T) {
	tests := []struct {
		name                string
		setupContext        func() context.Context
		requiredPermissions []string
		expectedStatus      int
	}{
		{
			name: "allow when user has all required permissions",
			setupContext: func() context.Context {
				authContext := &AuthContext{
					UserID:      uuid.New(),
					Permissions: []string{"users:read", "users:create", "users:delete"},
				}
				return context.WithValue(context.Background(), authContextKey{}, authContext)
			},
			requiredPermissions: []string{"users:read", "users:create"},
			expectedStatus:      http.StatusOK,
		},
		{
			name: "deny when user has only some required permissions",
			setupContext: func() context.Context {
				authContext := &AuthContext{
					UserID:      uuid.New(),
					Permissions: []string{"users:read"},
				}
				return context.WithValue(context.Background(), authContextKey{}, authContext)
			},
			requiredPermissions: []string{"users:read", "users:create"},
			expectedStatus:      http.StatusForbidden,
		},
		{
			name: "deny when user has none of required permissions",
			setupContext: func() context.Context {
				authContext := &AuthContext{
					UserID:      uuid.New(),
					Permissions: []string{"users:delete"},
				}
				return context.WithValue(context.Background(), authContextKey{}, authContext)
			},
			requiredPermissions: []string{"users:read", "users:create"},
			expectedStatus:      http.StatusForbidden,
		},
		{
			name: "deny when no auth context",
			setupContext: func() context.Context {
				return context.Background()
			},
			requiredPermissions: []string{"users:read"},
			expectedStatus:      http.StatusUnauthorized,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			nextHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(http.StatusOK)
			})

			middleware := RequireAllPermissions(testCase.requiredPermissions...)
			handler := middleware(nextHandler)

			request := httptest.NewRequest(http.MethodGet, "/test", nil)
			request = request.WithContext(testCase.setupContext())
			recorder := httptest.NewRecorder()

			handler.ServeHTTP(recorder, request)

			assert.Equal(t, testCase.expectedStatus, recorder.Code)
		})
	}
}

func TestRequireRole(t *testing.T) {
	tests := []struct {
		name           string
		setupContext   func() context.Context
		requiredRole   string
		expectedStatus int
	}{
		{
			name: "allow when user has required role",
			setupContext: func() context.Context {
				authContext := &AuthContext{
					UserID: uuid.New(),
					Roles:  []string{"admin", "user"},
				}
				return context.WithValue(context.Background(), authContextKey{}, authContext)
			},
			requiredRole:   "admin",
			expectedStatus: http.StatusOK,
		},
		{
			name: "deny when user lacks required role",
			setupContext: func() context.Context {
				authContext := &AuthContext{
					UserID: uuid.New(),
					Roles:  []string{"user"},
				}
				return context.WithValue(context.Background(), authContextKey{}, authContext)
			},
			requiredRole:   "admin",
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "deny when no auth context",
			setupContext: func() context.Context {
				return context.Background()
			},
			requiredRole:   "admin",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			nextHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(http.StatusOK)
			})

			middleware := RequireRole(testCase.requiredRole)
			handler := middleware(nextHandler)

			request := httptest.NewRequest(http.MethodGet, "/test", nil)
			request = request.WithContext(testCase.setupContext())
			recorder := httptest.NewRecorder()

			handler.ServeHTTP(recorder, request)

			assert.Equal(t, testCase.expectedStatus, recorder.Code)
		})
	}
}

func TestRequireAnyRole(t *testing.T) {
	tests := []struct {
		name           string
		setupContext   func() context.Context
		requiredRoles  []string
		expectedStatus int
	}{
		{
			name: "allow when user has one of required roles",
			setupContext: func() context.Context {
				authContext := &AuthContext{
					UserID: uuid.New(),
					Roles:  []string{"user"},
				}
				return context.WithValue(context.Background(), authContextKey{}, authContext)
			},
			requiredRoles:  []string{"admin", "user"},
			expectedStatus: http.StatusOK,
		},
		{
			name: "allow when user has multiple required roles",
			setupContext: func() context.Context {
				authContext := &AuthContext{
					UserID: uuid.New(),
					Roles:  []string{"admin", "user"},
				}
				return context.WithValue(context.Background(), authContextKey{}, authContext)
			},
			requiredRoles:  []string{"admin", "user"},
			expectedStatus: http.StatusOK,
		},
		{
			name: "deny when user has none of required roles",
			setupContext: func() context.Context {
				authContext := &AuthContext{
					UserID: uuid.New(),
					Roles:  []string{"guest"},
				}
				return context.WithValue(context.Background(), authContextKey{}, authContext)
			},
			requiredRoles:  []string{"admin", "user"},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "deny when no auth context",
			setupContext: func() context.Context {
				return context.Background()
			},
			requiredRoles:  []string{"admin"},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			nextHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(http.StatusOK)
			})

			middleware := RequireAnyRole(testCase.requiredRoles...)
			handler := middleware(nextHandler)

			request := httptest.NewRequest(http.MethodGet, "/test", nil)
			request = request.WithContext(testCase.setupContext())
			recorder := httptest.NewRecorder()

			handler.ServeHTTP(recorder, request)

			assert.Equal(t, testCase.expectedStatus, recorder.Code)
		})
	}
}

func TestResourceOwner(t *testing.T) {
	ownerID := uuid.New()
	otherID := uuid.New()

	tests := []struct {
		name           string
		setupContext   func() context.Context
		resourceID     string
		expectedStatus int
	}{
		{
			name: "allow when user is resource owner",
			setupContext: func() context.Context {
				authContext := &AuthContext{
					UserID: ownerID,
					Roles:  []string{"user"},
				}
				return context.WithValue(context.Background(), authContextKey{}, authContext)
			},
			resourceID:     ownerID.String(),
			expectedStatus: http.StatusOK,
		},
		{
			name: "allow when user is super_admin even if not owner",
			setupContext: func() context.Context {
				authContext := &AuthContext{
					UserID: otherID,
					Roles:  []string{"super_admin"},
				}
				return context.WithValue(context.Background(), authContextKey{}, authContext)
			},
			resourceID:     ownerID.String(),
			expectedStatus: http.StatusOK,
		},
		{
			name: "allow when user is admin even if not owner",
			setupContext: func() context.Context {
				authContext := &AuthContext{
					UserID: otherID,
					Roles:  []string{"admin"},
				}
				return context.WithValue(context.Background(), authContextKey{}, authContext)
			},
			resourceID:     ownerID.String(),
			expectedStatus: http.StatusOK,
		},
		{
			name: "deny when user is not owner and not admin",
			setupContext: func() context.Context {
				authContext := &AuthContext{
					UserID: otherID,
					Roles:  []string{"user"},
				}
				return context.WithValue(context.Background(), authContextKey{}, authContext)
			},
			resourceID:     ownerID.String(),
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "deny when no auth context",
			setupContext: func() context.Context {
				return context.Background()
			},
			resourceID:     ownerID.String(),
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			nextHandlerCalled := false
			nextHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				nextHandlerCalled = true
				writer.WriteHeader(http.StatusOK)
			})

			middleware := ResourceOwner("id")
			handler := middleware(nextHandler)

			request := httptest.NewRequest(http.MethodGet, "/users/"+testCase.resourceID, nil)

			routeContext := chi.NewRouteContext()
			routeContext.URLParams.Add("id", testCase.resourceID)
			ctx := context.WithValue(testCase.setupContext(), chi.RouteCtxKey, routeContext)
			request = request.WithContext(ctx)

			recorder := httptest.NewRecorder()

			handler.ServeHTTP(recorder, request)

			assert.Equal(t, testCase.expectedStatus, recorder.Code)
			if testCase.expectedStatus == http.StatusOK {
				assert.True(t, nextHandlerCalled)
			}
		})
	}
}

func TestResourceOwner_MissingResourceID(t *testing.T) {
	nextHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})

	authContext := &AuthContext{
		UserID: uuid.New(),
		Roles:  []string{"user"},
	}

	middleware := ResourceOwner("id")
	handler := middleware(nextHandler)

	request := httptest.NewRequest(http.MethodGet, "/users/", nil)

	routeContext := chi.NewRouteContext()
	ctx := context.WithValue(context.Background(), authContextKey{}, authContext)
	ctx = context.WithValue(ctx, chi.RouteCtxKey, routeContext)
	request = request.WithContext(ctx)

	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}
