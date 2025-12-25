---
applyTo: "backend/internal/handlers/**/*.go"
---

# API Handler Development Instructions

These instructions apply to all HTTP handlers in the backend.

## Handler Structure

### Basic Handler Template

```go
package handlers

import (
    "encoding/json"
    "errors"
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/yourorg/app/internal/domain"
    "github.com/yourorg/app/internal/service"
    "github.com/yourorg/app/pkg/response"
)

// UserHandler handles HTTP requests for user operations.
type UserHandler struct {
    svc    service.UserService
    logger *slog.Logger
}

// NewUserHandler creates a new UserHandler with the given dependencies.
func NewUserHandler(svc service.UserService, logger *slog.Logger) *UserHandler {
    return &UserHandler{
        svc:    svc,
        logger: logger,
    }
}

// RegisterRoutes registers all user routes on the given router.
func (h *UserHandler) RegisterRoutes(r chi.Router) {
    r.Route("/users", func(r chi.Router) {
        r.Get("/", h.List)
        r.Post("/", h.Create)
        r.Get("/{id}", h.Get)
        r.Put("/{id}", h.Update)
        r.Delete("/{id}", h.Delete)
    })
}
```

## Request Handling

### Parse and Validate Input

```go
import "github.com/go-playground/validator/v10"

var validate = validator.New()

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    // 1. Parse request body
    var input domain.CreateUserInput
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        response.BadRequest(w, "Invalid request body")
        return
    }

    // 2. Validate input
    if err := validate.Struct(input); err != nil {
        validationErrors := formatValidationErrors(err.(validator.ValidationErrors))
        response.ValidationError(w, validationErrors)
        return
    }

    // 3. Call service
    user, err := h.svc.CreateUser(ctx, input)
    if err != nil {
        h.handleError(w, err)
        return
    }

    // 4. Return response
    response.JSON(w, http.StatusCreated, user)
}
```

### Extract URL Parameters

```go
func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    // Get path parameter
    id := chi.URLParam(r, "id")
    if id == "" {
        response.BadRequest(w, "Missing user ID")
        return
    }

    // Validate UUID format if needed
    if _, err := uuid.Parse(id); err != nil {
        response.BadRequest(w, "Invalid user ID format")
        return
    }

    user, err := h.svc.GetUser(ctx, id)
    if err != nil {
        h.handleError(w, err)
        return
    }

    response.JSON(w, http.StatusOK, user)
}
```

### Parse Query Parameters

```go
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    // Parse pagination
    page, _ := strconv.Atoi(r.URL.Query().Get("page"))
    if page < 1 {
        page = 1
    }

    perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
    if perPage < 1 || perPage > 100 {
        perPage = 20
    }

    // Parse filters
    status := r.URL.Query().Get("status")
    search := r.URL.Query().Get("search")

    opts := repository.ListOptions{
        Page:    page,
        PerPage: perPage,
        Status:  status,
        Search:  search,
    }

    users, total, err := h.svc.ListUsers(ctx, opts)
    if err != nil {
        h.handleError(w, err)
        return
    }

    response.JSONWithMeta(w, http.StatusOK, users, &response.Meta{
        Page:    page,
        PerPage: perPage,
        Total:   total,
    })
}
```

## Response Formatting

### Standard Response Structure

```go
package response

// Success response
type Response struct {
    Data interface{} `json:"data,omitempty"`
    Meta *Meta       `json:"meta,omitempty"`
}

type Meta struct {
    Page    int `json:"page"`
    PerPage int `json:"per_page"`
    Total   int `json:"total"`
}

// Error response
type ErrorResponse struct {
    Error ErrorBody `json:"error"`
}

type ErrorBody struct {
    Code    string      `json:"code"`
    Message string      `json:"message"`
    Details interface{} `json:"details,omitempty"`
}
```

### Response Helper Functions

```go
package response

func JSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(Response{Data: data})
}

func JSONWithMeta(w http.ResponseWriter, status int, data interface{}, meta *Meta) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(Response{Data: data, Meta: meta})
}

func Error(w http.ResponseWriter, status int, code, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(ErrorResponse{
        Error: ErrorBody{Code: code, Message: message},
    })
}

func BadRequest(w http.ResponseWriter, message string) {
    Error(w, http.StatusBadRequest, "BAD_REQUEST", message)
}

func NotFound(w http.ResponseWriter, message string) {
    Error(w, http.StatusNotFound, "NOT_FOUND", message)
}

func Conflict(w http.ResponseWriter, message string) {
    Error(w, http.StatusConflict, "CONFLICT", message)
}

func Unauthorized(w http.ResponseWriter, message string) {
    Error(w, http.StatusUnauthorized, "UNAUTHORIZED", message)
}

func Forbidden(w http.ResponseWriter, message string) {
    Error(w, http.StatusForbidden, "FORBIDDEN", message)
}

func ValidationError(w http.ResponseWriter, errors []ValidationError) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(ErrorResponse{
        Error: ErrorBody{
            Code:    "VALIDATION_ERROR",
            Message: "Validation failed",
            Details: errors,
        },
    })
}

func InternalError(w http.ResponseWriter, err error) {
    slog.Error("internal error", "error", err)
    Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "An unexpected error occurred")
}
```

## Error Handling

### Centralized Error Handler

```go
func (h *UserHandler) handleError(w http.ResponseWriter, err error) {
    switch {
    case errors.Is(err, domain.ErrNotFound):
        response.NotFound(w, "Resource not found")
    case errors.Is(err, domain.ErrConflict):
        response.Conflict(w, "Resource already exists")
    case errors.Is(err, domain.ErrInvalidInput):
        response.BadRequest(w, err.Error())
    case errors.Is(err, domain.ErrUnauthorized):
        response.Unauthorized(w, "Authentication required")
    case errors.Is(err, domain.ErrForbidden):
        response.Forbidden(w, "Access denied")
    default:
        response.InternalError(w, err)
    }
}
```

### Validation Error Formatting

```go
type ValidationError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}

func formatValidationErrors(errs validator.ValidationErrors) []ValidationError {
    var result []ValidationError
    for _, err := range errs {
        result = append(result, ValidationError{
            Field:   toSnakeCase(err.Field()),
            Message: getErrorMessage(err),
        })
    }
    return result
}

func getErrorMessage(err validator.FieldError) string {
    switch err.Tag() {
    case "required":
        return "This field is required"
    case "email":
        return "Invalid email format"
    case "min":
        return fmt.Sprintf("Must be at least %s characters", err.Param())
    case "max":
        return fmt.Sprintf("Must be at most %s characters", err.Param())
    default:
        return fmt.Sprintf("Failed validation: %s", err.Tag())
    }
}
```

## Middleware Integration

### Using Middleware in Handlers

```go
func (h *UserHandler) RegisterRoutes(r chi.Router) {
    r.Route("/users", func(r chi.Router) {
        // Public routes
        r.Post("/", h.Create)
        r.Post("/login", h.Login)

        // Protected routes
        r.Group(func(r chi.Router) {
            r.Use(middleware.RequireAuth)
            r.Get("/me", h.GetCurrentUser)
            r.Put("/me", h.UpdateCurrentUser)
        })

        // Admin-only routes
        r.Group(func(r chi.Router) {
            r.Use(middleware.RequireAuth)
            r.Use(middleware.RequireRole("admin"))
            r.Get("/", h.List)
            r.Get("/{id}", h.Get)
            r.Delete("/{id}", h.Delete)
        })
    })
}
```

### Extracting User from Context

```go
func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    // Get user ID from context (set by auth middleware)
    userID := middleware.GetUserID(ctx)
    if userID == "" {
        response.Unauthorized(w, "Not authenticated")
        return
    }

    user, err := h.svc.GetUser(ctx, userID)
    if err != nil {
        h.handleError(w, err)
        return
    }

    response.JSON(w, http.StatusOK, user)
}
```

## API Versioning

### URL-Based Versioning

```go
func SetupRoutes(r chi.Router) {
    // API v1
    r.Route("/api/v1", func(r chi.Router) {
        r.Use(middleware.Logger)
        r.Use(middleware.Recovery)
        r.Use(middleware.CORS)

        userHandler.RegisterRoutes(r)
        authHandler.RegisterRoutes(r)
    })

    // API v2 (future)
    r.Route("/api/v2", func(r chi.Router) {
        // New handlers with breaking changes
    })
}
```

## Request Logging

```go
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    h.logger.Info("create user request",
        "method", r.Method,
        "path", r.URL.Path,
        "ip", r.RemoteAddr,
    )

    // ... handler logic ...

    h.logger.Info("user created",
        "user_id", user.ID,
        "email", user.Email,
    )
}
```

## Testing Handlers

```go
func TestUserHandler_Create(t *testing.T) {
    tests := []struct {
        name       string
        body       interface{}
        setup      func(*MockService)
        wantStatus int
        wantCode   string
    }{
        {
            name: "success",
            body: map[string]string{
                "email":    "test@example.com",
                "name":     "Test User",
                "password": "password123",
            },
            setup: func(m *MockService) {
                m.On("CreateUser", mock.Anything, mock.AnythingOfType("CreateUserInput")).
                    Return(&domain.User{ID: "123", Email: "test@example.com"}, nil)
            },
            wantStatus: http.StatusCreated,
        },
        {
            name: "validation error - invalid email",
            body: map[string]string{
                "email":    "invalid",
                "name":     "Test",
                "password": "password123",
            },
            setup:      func(m *MockService) {},
            wantStatus: http.StatusBadRequest,
            wantCode:   "VALIDATION_ERROR",
        },
        {
            name: "conflict - email exists",
            body: map[string]string{
                "email":    "exists@example.com",
                "name":     "Test",
                "password": "password123",
            },
            setup: func(m *MockService) {
                m.On("CreateUser", mock.Anything, mock.AnythingOfType("CreateUserInput")).
                    Return(nil, domain.ErrConflict)
            },
            wantStatus: http.StatusConflict,
            wantCode:   "CONFLICT",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockSvc := new(MockService)
            tt.setup(mockSvc)

            handler := NewUserHandler(mockSvc, slog.Default())
            router := chi.NewRouter()
            handler.RegisterRoutes(router)

            body, _ := json.Marshal(tt.body)
            req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
            req.Header.Set("Content-Type", "application/json")
            rec := httptest.NewRecorder()

            router.ServeHTTP(rec, req)

            assert.Equal(t, tt.wantStatus, rec.Code)

            if tt.wantCode != "" {
                var resp response.ErrorResponse
                json.Unmarshal(rec.Body.Bytes(), &resp)
                assert.Equal(t, tt.wantCode, resp.Error.Code)
            }
        })
    }
}
```
