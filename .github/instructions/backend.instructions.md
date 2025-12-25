---
applyTo: "backend/**/*.go"
---

# Go Backend Development Instructions

These instructions apply to all Go files in the backend directory, following Clean Architecture, Domain-Driven Design (DDD), and CQRS patterns.

## Project Structure (Clean Architecture + DDD + CQRS)

```
backend/
├── cmd/
│   └── api/
│       └── main.go                    # Entry point, wiring
├── internal/
│   ├── domain/                        # Domain Layer (innermost)
│   │   ├── user/                      # User aggregate
│   │   │   ├── user.go                # Entity + Value Objects
│   │   │   ├── repository.go          # Repository interface (port)
│   │   │   └── events.go              # Domain events
│   │   └── shared/                    # Shared domain concepts
│   │       ├── errors.go              # Domain errors
│   │       └── valueobjects.go        # Shared value objects
│   │
│   ├── application/                   # Application Layer (use cases)
│   │   ├── command/                   # CQRS Commands (write operations)
│   │   │   ├── create_user.go
│   │   │   └── update_user.go
│   │   ├── query/                     # CQRS Queries (read operations)
│   │   │   ├── get_user.go
│   │   │   └── list_users.go
│   │   └── dto/                       # Data Transfer Objects
│   │       └── user_dto.go
│   │
│   ├── infrastructure/                # Infrastructure Layer (adapters)
│   │   ├── persistence/               # Database implementations
│   │   │   └── postgres/
│   │   │       ├── user_repository.go
│   │   │       └── queries/           # SQLC generated (optional)
│   │   ├── cache/                     # Cache implementations
│   │   │   └── redis/
│   │   └── messaging/                 # Event bus implementations
│   │       └── memory/
│   │
│   └── interfaces/                    # Interface Adapters Layer
│       └── http/
│           ├── handler/               # HTTP handlers
│           │   └── user_handler.go
│           ├── middleware/            # HTTP middleware
│           ├── router/                # Route definitions
│           └── dto/                   # HTTP request/response DTOs
│
├── pkg/                               # Shared packages (can be imported)
│   ├── config/
│   ├── logger/
│   └── validator/
│
└── migrations/                        # Goose migrations
    └── sql/
```

## Dependency Rule

Dependencies MUST point inward. Inner layers cannot know about outer layers.

```
┌─────────────────────────────────────────────┐
│              Interfaces (HTTP)              │ ← Outermost
├─────────────────────────────────────────────┤
│      Infrastructure (Postgres, Redis)       │
├─────────────────────────────────────────────┤
│    Application (Commands, Queries, DTOs)    │
├─────────────────────────────────────────────┤
│          Domain (Entities, Ports)           │ ← Innermost
└─────────────────────────────────────────────┘
```

## Domain Layer (DDD)

### Entities and Aggregates

```go
// internal/domain/user/user.go
package user

import (
    "errors"
    "time"

    "github.com/google/uuid"
)

// User is the aggregate root for user-related operations.
type User struct {
    id        uuid.UUID
    email     Email      // Value Object
    name      string
    role      Role       // Value Object
    createdAt time.Time
    updatedAt time.Time
}

// NewUser creates a new User entity with validation.
func NewUser(email Email, name string) (*User, error) {
    if name == "" {
        return nil, ErrInvalidName
    }

    now := time.Now().UTC()
    return &User{
        id:        uuid.New(),
        email:     email,
        name:      name,
        role:      RoleMember,
        createdAt: now,
        updatedAt: now,
    }, nil
}

// Reconstitute creates a User from persisted data (no validation).
func Reconstitute(id uuid.UUID, email Email, name string, role Role, createdAt, updatedAt time.Time) *User {
    return &User{
        id:        id,
        email:     email,
        name:      name,
        role:      role,
        createdAt: createdAt,
        updatedAt: updatedAt,
    }
}

// Business methods on the entity
func (u *User) ChangeName(name string) error {
    if name == "" {
        return ErrInvalidName
    }
    u.name = name
    u.updatedAt = time.Now().UTC()
    return nil
}

func (u *User) PromoteToAdmin() {
    u.role = RoleAdmin
    u.updatedAt = time.Now().UTC()
}

// Getters (no setters - enforce invariants through methods)
func (u *User) ID() uuid.UUID      { return u.id }
func (u *User) Email() Email       { return u.email }
func (u *User) Name() string       { return u.name }
func (u *User) Role() Role         { return u.role }
func (u *User) CreatedAt() time.Time { return u.createdAt }
func (u *User) UpdatedAt() time.Time { return u.updatedAt }
```

### Value Objects

```go
// internal/domain/user/email.go
package user

import (
    "regexp"
    "strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// Email is a value object representing a validated email address.
type Email struct {
    value string
}

func NewEmail(value string) (Email, error) {
    normalized := strings.ToLower(strings.TrimSpace(value))
    if !emailRegex.MatchString(normalized) {
        return Email{}, ErrInvalidEmail
    }
    return Email{value: normalized}, nil
}

func (e Email) String() string { return e.value }
func (e Email) IsZero() bool   { return e.value == "" }

// Role is a value object representing user roles.
type Role string

const (
    RoleMember Role = "member"
    RoleAdmin  Role = "admin"
)

func (r Role) IsValid() bool {
    return r == RoleMember || r == RoleAdmin
}
```

### Repository Interface (Port)

```go
// internal/domain/user/repository.go
package user

import (
    "context"

    "github.com/google/uuid"
)

// Repository defines the port for user persistence.
// This interface is defined in the domain layer and implemented in infrastructure.
type Repository interface {
    FindByID(ctx context.Context, id uuid.UUID) (*User, error)
    FindByEmail(ctx context.Context, email Email) (*User, error)
    Save(ctx context.Context, user *User) error
    Delete(ctx context.Context, id uuid.UUID) error
}
```

### Domain Errors

```go
// internal/domain/shared/errors.go
package shared

import "errors"

var (
    ErrNotFound     = errors.New("resource not found")
    ErrConflict     = errors.New("resource already exists")
    ErrInvalidInput = errors.New("invalid input")
    ErrUnauthorized = errors.New("unauthorized")
    ErrForbidden    = errors.New("forbidden")
)

// internal/domain/user/errors.go
package user

import "errors"

var (
    ErrInvalidEmail = errors.New("invalid email address")
    ErrInvalidName  = errors.New("name cannot be empty")
    ErrUserNotFound = errors.New("user not found")
)
```

## Application Layer (CQRS)

### Commands (Write Operations)

```go
// internal/application/command/create_user.go
package command

import (
    "context"
    "fmt"

    "github.com/yourorg/app/internal/domain/shared"
    "github.com/yourorg/app/internal/domain/user"
)

type CreateUserCommand struct {
    Email string
    Name  string
}

type CreateUserHandler struct {
    repo   user.Repository
    logger Logger
}

func NewCreateUserHandler(repo user.Repository, logger Logger) *CreateUserHandler {
    return &CreateUserHandler{repo: repo, logger: logger}
}

func (h *CreateUserHandler) Handle(ctx context.Context, cmd CreateUserCommand) (*user.User, error) {
    // Validate and create value objects
    email, err := user.NewEmail(cmd.Email)
    if err != nil {
        return nil, fmt.Errorf("invalid email: %w", err)
    }

    // Check for existing user (business rule)
    existing, err := h.repo.FindByEmail(ctx, email)
    if err != nil && err != user.ErrUserNotFound {
        return nil, fmt.Errorf("check existing user: %w", err)
    }
    if existing != nil {
        return nil, shared.ErrConflict
    }

    // Create domain entity
    u, err := user.NewUser(email, cmd.Name)
    if err != nil {
        return nil, fmt.Errorf("create user: %w", err)
    }

    // Persist
    if err := h.repo.Save(ctx, u); err != nil {
        return nil, fmt.Errorf("save user: %w", err)
    }

    h.logger.Info("user created", "id", u.ID(), "email", email.String())

    return u, nil
}
```

### Queries (Read Operations)

```go
// internal/application/query/get_user.go
package query

import (
    "context"

    "github.com/google/uuid"
    "github.com/yourorg/app/internal/application/dto"
    "github.com/yourorg/app/internal/domain/user"
)

type GetUserQuery struct {
    ID uuid.UUID
}

type GetUserHandler struct {
    repo user.Repository
}

func NewGetUserHandler(repo user.Repository) *GetUserHandler {
    return &GetUserHandler{repo: repo}
}

func (h *GetUserHandler) Handle(ctx context.Context, q GetUserQuery) (*dto.UserResponse, error) {
    u, err := h.repo.FindByID(ctx, q.ID)
    if err != nil {
        return nil, err
    }

    return dto.UserResponseFromDomain(u), nil
}
```

### DTOs (Data Transfer Objects)

```go
// internal/application/dto/user_dto.go
package dto

import (
    "time"

    "github.com/google/uuid"
    "github.com/yourorg/app/internal/domain/user"
)

type UserResponse struct {
    ID        uuid.UUID `json:"id"`
    Email     string    `json:"email"`
    Name      string    `json:"name"`
    Role      string    `json:"role"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func UserResponseFromDomain(u *user.User) *UserResponse {
    return &UserResponse{
        ID:        u.ID(),
        Email:     u.Email().String(),
        Name:      u.Name(),
        Role:      string(u.Role()),
        CreatedAt: u.CreatedAt(),
        UpdatedAt: u.UpdatedAt(),
    }
}
```

## Infrastructure Layer

### Repository Implementation

```go
// internal/infrastructure/persistence/postgres/user_repository.go
package postgres

import (
    "context"
    "errors"
    "fmt"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/yourorg/app/internal/domain/user"
)

type UserRepository struct {
    db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
    query := `
        SELECT id, email, name, role, created_at, updated_at
        FROM users
        WHERE id = $1 AND deleted_at IS NULL
    `

    var (
        dbID        uuid.UUID
        email       string
        name        string
        role        string
        createdAt   time.Time
        updatedAt   time.Time
    )

    err := r.db.QueryRow(ctx, query, id).Scan(
        &dbID, &email, &name, &role, &createdAt, &updatedAt,
    )

    if errors.Is(err, pgx.ErrNoRows) {
        return nil, user.ErrUserNotFound
    }
    if err != nil {
        return nil, fmt.Errorf("query user: %w", err)
    }

    // Reconstitute domain entity from persistence
    emailVO, _ := user.NewEmail(email) // Already validated in DB
    return user.Reconstitute(dbID, emailVO, name, user.Role(role), createdAt, updatedAt), nil
}

func (r *UserRepository) Save(ctx context.Context, u *user.User) error {
    query := `
        INSERT INTO users (id, email, name, role, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        ON CONFLICT (id) DO UPDATE SET
            name = EXCLUDED.name,
            role = EXCLUDED.role,
            updated_at = EXCLUDED.updated_at
    `

    _, err := r.db.Exec(ctx, query,
        u.ID(),
        u.Email().String(),
        u.Name(),
        string(u.Role()),
        u.CreatedAt(),
        u.UpdatedAt(),
    )

    if err != nil {
        return fmt.Errorf("save user: %w", err)
    }

    return nil
}
```

## Interface Layer (HTTP Handlers)

```go
// internal/interfaces/http/handler/user_handler.go
package handler

import (
    "encoding/json"
    "errors"
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/google/uuid"
    "github.com/yourorg/app/internal/application/command"
    "github.com/yourorg/app/internal/application/query"
    "github.com/yourorg/app/internal/domain/shared"
    "github.com/yourorg/app/internal/domain/user"
)

type UserHandler struct {
    createUser *command.CreateUserHandler
    getUser    *query.GetUserHandler
}

func NewUserHandler(
    createUser *command.CreateUserHandler,
    getUser *query.GetUserHandler,
) *UserHandler {
    return &UserHandler{
        createUser: createUser,
        getUser:    getUser,
    }
}

type CreateUserRequest struct {
    Email string `json:"email" validate:"required,email"`
    Name  string `json:"name" validate:"required,min=2,max=100"`
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, http.StatusBadRequest, "invalid request body")
        return
    }

    cmd := command.CreateUserCommand{
        Email: req.Email,
        Name:  req.Name,
    }

    u, err := h.createUser.Handle(r.Context(), cmd)
    if err != nil {
        handleDomainError(w, err)
        return
    }

    respondJSON(w, http.StatusCreated, dto.UserResponseFromDomain(u))
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        respondError(w, http.StatusBadRequest, "invalid user id")
        return
    }

    q := query.GetUserQuery{ID: id}
    result, err := h.getUser.Handle(r.Context(), q)
    if err != nil {
        handleDomainError(w, err)
        return
    }

    respondJSON(w, http.StatusOK, result)
}

func handleDomainError(w http.ResponseWriter, err error) {
    switch {
    case errors.Is(err, user.ErrUserNotFound), errors.Is(err, shared.ErrNotFound):
        respondError(w, http.StatusNotFound, "resource not found")
    case errors.Is(err, shared.ErrConflict):
        respondError(w, http.StatusConflict, "resource already exists")
    case errors.Is(err, user.ErrInvalidEmail), errors.Is(err, user.ErrInvalidName):
        respondError(w, http.StatusBadRequest, err.Error())
    default:
        respondError(w, http.StatusInternalServerError, "internal server error")
    }
}
```

## Best Practices

### 1. Keep Domain Pure

- No framework dependencies in domain layer
- No infrastructure concerns (database, HTTP, etc.)
- Business logic only

### 2. Use Dependency Injection

```go
// cmd/api/main.go - Wire everything together
func main() {
    // Infrastructure
    db := postgres.NewConnection(cfg.DatabaseURL)
    userRepo := postgres.NewUserRepository(db)

    // Application
    createUserHandler := command.NewCreateUserHandler(userRepo, logger)
    getUserHandler := query.NewGetUserHandler(userRepo)

    // Interfaces
    userHTTPHandler := handler.NewUserHandler(createUserHandler, getUserHandler)

    // Router
    r := chi.NewRouter()
    r.Post("/users", userHTTPHandler.Create)
    r.Get("/users/{id}", userHTTPHandler.Get)
}
```

### 3. Separate Read and Write Models (CQRS)

- Commands: Change state, return minimal data
- Queries: Read state, can use optimized read models

### 4. Error Handling

- Define domain errors in domain layer
- Map to HTTP status codes at interface layer
- Always wrap errors with context

## References

- [Three Dots Labs - DDD, CQRS, Clean Architecture](https://threedots.tech/post/ddd-cqrs-clean-architecture-combined/)
- [Wild Workouts Go DDD Example](https://github.com/ThreeDotsLabs/wild-workouts-go-ddd-example)
