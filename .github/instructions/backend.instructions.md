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
│   │   ├── auth/                      # Auth aggregate
│   │   ├── role/                      # Role aggregate
│   │   ├── permission/                # Permission aggregate
│   │   └── shared/                    # Shared domain concepts
│   │       ├── errors.go              # Domain errors
│   │       └── valueobjects.go        # Shared value objects
│   │
│   ├── application/                   # Application Layer (CQRS, domain-aligned)
│   │   ├── cqrs/                      # Base CQRS interfaces
│   │   ├── user/                      # User bounded context
│   │   │   ├── command/               # usercommand package
│   │   │   │   ├── create_user.go
│   │   │   │   └── update_user.go
│   │   │   ├── query/                 # userquery package
│   │   │   │   ├── get_user.go
│   │   │   │   └── list_users.go
│   │   │   └── dto/                   # userdto package
│   │   │       └── user_dto.go
│   │   └── auth/                      # Auth bounded context
│   │       ├── command/               # authcommand package
│   │       ├── query/                 # authquery package
│   │       └── dto/                   # authdto package
│   │
│   ├── infrastructure/                # Infrastructure Layer (adapters)
│   │   ├── persistence/               # Database implementations
│   │   │   ├── postgres/              # Database utilities
│   │   │   │   └── ...
│   │   │   └── repository/            # Repository implementations
│   │   │       └── user_repository.go # Implements domain.UserRepository
│   │   ├── cache/                     # Cache implementations
│   │   │   └── redis/
│   │   ├── messaging/                 # Event bus implementations
│   │   │   └── memory/
│   │   │       └── event_bus.go
│   │   └── audit/                     # Audit logging
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
└── migrations/                        # golang-migrate migrations
    ├── 000001_create_users_table.up.sql
    └── 000001_create_users_table.down.sql
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

## Application Layer (CQRS, Domain-Aligned)

The application layer is organized by domain/bounded context. Each domain has its own `command/`, `query/`, and `dto/` packages with unique package names.

### Package Naming Convention

```go
import (
    usercommand "yourapp/internal/application/user/command"
    userquery "yourapp/internal/application/user/query"
    userdto "yourapp/internal/application/user/dto"

    authcommand "yourapp/internal/application/auth/command"
    authquery "yourapp/internal/application/auth/query"
    authdto "yourapp/internal/application/auth/dto"
)
```

### Commands (Write Operations)

```go
// internal/application/user/command/create_user.go
package usercommand

import (
    "context"
    "fmt"

    userdto "github.com/yourorg/app/internal/application/user/dto"
    "github.com/yourorg/app/internal/domain/shared"
    "github.com/yourorg/app/internal/domain/user"
)

type CreateUserCommand struct {
    Email string
    Name  string
}

type CreateUserHandler struct {
    repository user.Repository
    logger     Logger
}

func NewCreateUserHandler(repository user.Repository, logger Logger) *CreateUserHandler {
    return &CreateUserHandler{repository: repository, logger: logger}
}

func (handler *CreateUserHandler) Handle(ctx context.Context, command CreateUserCommand) (*userdto.UserDTO, error) {
    email, err := user.NewEmail(command.Email)
    if err != nil {
        return nil, fmt.Errorf("invalid email: %w", err)
    }

    existing, err := handler.repository.FindByEmail(ctx, email)
    if err != nil && err != user.ErrUserNotFound {
        return nil, fmt.Errorf("check existing user: %w", err)
    }
    if existing != nil {
        return nil, shared.ErrConflict
    }

    newUser, err := user.NewUser(email, command.Name)
    if err != nil {
        return nil, fmt.Errorf("create user: %w", err)
    }

    if err := handler.repository.Save(ctx, newUser); err != nil {
        return nil, fmt.Errorf("save user: %w", err)
    }

    handler.logger.Info("user created", "id", newUser.ID(), "email", email.String())

    return userdto.UserFromDomain(newUser), nil
}
```

### Queries (Read Operations)

```go
// internal/application/user/query/get_user.go
package userquery

import (
    "context"

    "github.com/google/uuid"
    userdto "github.com/yourorg/app/internal/application/user/dto"
    "github.com/yourorg/app/internal/domain/user"
)

type GetUserQuery struct {
    ID uuid.UUID
}

type GetUserHandler struct {
    repository user.Repository
}

func NewGetUserHandler(repository user.Repository) *GetUserHandler {
    return &GetUserHandler{repository: repository}
}

func (handler *GetUserHandler) Handle(ctx context.Context, query GetUserQuery) (*userdto.UserDTO, error) {
    foundUser, err := handler.repository.FindByID(ctx, query.ID)
    if err != nil {
        return nil, err
    }

    return userdto.UserFromDomain(foundUser), nil
}
```

### DTOs (Data Transfer Objects)

```go
// internal/application/user/dto/user_dto.go
package userdto

import (
    "time"

    "github.com/google/uuid"
    "github.com/yourorg/app/internal/domain/user"
)

type UserDTO struct {
    ID        uuid.UUID `json:"id"`
    Email     string    `json:"email"`
    Name      string    `json:"name"`
    Role      string    `json:"role"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func UserFromDomain(domainUser *user.User) *UserDTO {
    return &UserDTO{
        ID:        domainUser.ID(),
        Email:     domainUser.Email().String(),
        Name:      domainUser.Name(),
        Role:      string(domainUser.Role()),
        CreatedAt: domainUser.CreatedAt(),
        UpdatedAt: domainUser.UpdatedAt(),
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
    usercommand "github.com/yourorg/app/internal/application/user/command"
    userquery "github.com/yourorg/app/internal/application/user/query"
    "github.com/yourorg/app/internal/domain/shared"
    "github.com/yourorg/app/internal/domain/user"
)

type UserHandler struct {
    createUserHandler *usercommand.CreateUserHandler
    getUserHandler    *userquery.GetUserHandler
}

func NewUserHandler(
    createUserHandler *usercommand.CreateUserHandler,
    getUserHandler *userquery.GetUserHandler,
) *UserHandler {
    return &UserHandler{
        createUserHandler: createUserHandler,
        getUserHandler:    getUserHandler,
    }
}

type CreateUserRequest struct {
    Email string `json:"email" validate:"required,email"`
    Name  string `json:"name" validate:"required,min=2,max=100"`
}

func (handler *UserHandler) Create(writer http.ResponseWriter, request *http.Request) {
    var requestBody CreateUserRequest
    if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
        respondError(writer, http.StatusBadRequest, "invalid request body")
        return
    }

    command := usercommand.CreateUserCommand{
        Email: requestBody.Email,
        Name:  requestBody.Name,
    }

    result, err := handler.createUserHandler.Handle(request.Context(), command)
    if err != nil {
        handleDomainError(writer, err)
        return
    }

    respondJSON(writer, http.StatusCreated, result)
}

func (handler *UserHandler) Get(writer http.ResponseWriter, request *http.Request) {
    idParam := chi.URLParam(request, "id")
    id, err := uuid.Parse(idParam)
    if err != nil {
        respondError(writer, http.StatusBadRequest, "invalid user id")
        return
    }

    query := userquery.GetUserQuery{ID: id}
    result, err := handler.getUserHandler.Handle(request.Context(), query)
    if err != nil {
        handleDomainError(writer, err)
        return
    }

    respondJSON(writer, http.StatusOK, result)
}

func handleDomainError(writer http.ResponseWriter, err error) {
    switch {
    case errors.Is(err, user.ErrUserNotFound), errors.Is(err, shared.ErrNotFound):
        respondError(writer, http.StatusNotFound, "resource not found")
    case errors.Is(err, shared.ErrConflict):
        respondError(writer, http.StatusConflict, "resource already exists")
    case errors.Is(err, user.ErrInvalidEmail), errors.Is(err, user.ErrInvalidName):
        respondError(writer, http.StatusBadRequest, err.Error())
    default:
        respondError(writer, http.StatusInternalServerError, "internal server error")
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
