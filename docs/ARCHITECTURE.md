# System Architecture

This document provides a comprehensive overview of our full-stack application architecture, covering both the Go backend (Clean Architecture + DDD + CQRS) and React frontend.

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                         Client Layer                            │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │              React 19 + TypeScript Frontend              │   │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐    │   │
│  │  │ Pages   │  │Components│  │ Hooks   │  │ Stores  │    │   │
│  │  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘    │   │
│  │       │           │           │           │            │   │
│  │  ┌────┴───────────┴───────────┴───────────┴────┐       │   │
│  │  │           TanStack Query + Zustand           │       │   │
│  │  └──────────────────────┬──────────────────────┘       │   │
│  └─────────────────────────┼───────────────────────────────┘   │
└────────────────────────────┼────────────────────────────────────┘
                             │ HTTP/REST
┌────────────────────────────┼────────────────────────────────────┐
│                 Go Backend (DDD + CQRS)                         │
│  ┌─────────────────────────┼───────────────────────────────┐   │
│  │              Interfaces Layer (HTTP Handlers)            │   │
│  │  ┌─────────────────────┴─────────────────────────┐      │   │
│  │  │         HTTP Handlers (chi router)             │      │   │
│  │  └──────────────┬────────────────┬───────────────┘      │   │
│  │                 │                │                       │   │
│  │  ┌──────────────▼────┐  ┌───────▼────────────┐          │   │
│  │  │  Command Handlers │  │  Query Handlers    │  CQRS    │   │
│  │  │  (Write Ops)      │  │  (Read Ops)        │          │   │
│  │  └──────────────┬────┘  └───────┬────────────┘          │   │
│  │                 │               │                        │   │
│  │  ┌──────────────▼───────────────▼────────────┐          │   │
│  │  │        Domain Layer (Aggregates)           │  DDD    │   │
│  │  │  Entities │ Value Objects │ Domain Events  │          │   │
│  │  │  Repository Interfaces (Ports)             │          │   │
│  │  └──────────────────────┬────────────────────┘          │   │
│  │                         │                                │   │
│  │  ┌──────────────────────▼────────────────────┐          │   │
│  │  │      Infrastructure (Implementations)      │          │   │
│  │  │  PostgreSQL Repos │ Redis Cache │ Events   │          │   │
│  │  └───────────────────────────────────────────┘          │   │
│  └─────────────────────────────────────────────────────────┘   │
└────────────────────────────┬────────────────────────────────────┘
                             │
┌────────────────────────────┼────────────────────────────────────┐
│                       Data Layer                                │
│  ┌─────────────┐  ┌────────┴────────┐  ┌─────────────┐         │
│  │  PostgreSQL │  │      Redis      │  │   S3/Blob   │         │
│  │  (Primary)  │  │    (Cache)      │  │  (Storage)  │         │
│  └─────────────┘  └─────────────────┘  └─────────────┘         │
└─────────────────────────────────────────────────────────────────┘
```

## Backend Architecture (Go)

### Clean Architecture + DDD + CQRS

Our backend follows Clean Architecture principles combined with Domain-Driven Design (DDD) and Command Query Responsibility Segregation (CQRS).

#### The Dependency Rule

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

### Project Structure

```
backend/
├── cmd/
│   └── api/
│       └── main.go                    # Entry point, dependency wiring
├── internal/
│   ├── domain/                        # Domain Layer (innermost, pure business logic)
│   │   ├── user/                      # User aggregate
│   │   │   ├── user.go                # Entity + Value Objects
│   │   │   ├── repository.go          # Repository interface (port)
│   │   │   ├── errors.go              # Domain-specific errors
│   │   │   └── events.go              # Domain events
│   │   ├── order/                     # Order aggregate
│   │   │   └── ...
│   │   └── shared/                    # Shared domain concepts
│   │       ├── errors.go              # Common domain errors
│   │       └── valueobjects.go        # Shared value objects
│   │
│   ├── application/                   # Application Layer (CQRS)
│   │   ├── command/                   # Commands (write operations)
│   │   │   ├── create_user.go
│   │   │   ├── update_user.go
│   │   │   └── delete_user.go
│   │   ├── query/                     # Queries (read operations)
│   │   │   ├── get_user.go
│   │   │   ├── list_users.go
│   │   │   └── search_users.go
│   │   └── dto/                       # Data Transfer Objects
│   │       └── user_dto.go
│   │
│   ├── infrastructure/                # Infrastructure Layer (adapters)
│   │   ├── persistence/               # Database implementations
│   │   │   ├── postgres/              # Database utilities
│   │   │   │   ├── connection.go      # Connection pool management
│   │   │   │   ├── unit_of_work.go    # Transaction management
│   │   │   │   ├── query_builder.go   # SQL query helpers
│   │   │   │   └── errors.go          # Database error types
│   │   │   └── repository/            # Repository implementations
│   │   │       └── user_repository.go # Implements domain.UserRepository
│   │   ├── cache/                     # Cache implementations
│   │   │   └── redis/
│   │   │       └── cache.go
│   │   └── messaging/                 # Event bus implementations
│   │       └── memory/
│   │           └── event_bus.go
│   │
│   └── interfaces/                    # Interface Adapters Layer
│       └── http/
│           ├── handler/               # HTTP handlers
│           │   └── user_handler.go
│           ├── middleware/            # HTTP middleware
│           │   ├── auth.go
│           │   ├── logging.go
│           │   └── recovery.go
│           ├── router/                # Route definitions
│           │   └── router.go
│           └── dto/                   # HTTP request/response DTOs
│               └── user_dto.go
│
├── pkg/                               # Shared packages (can be imported)
│   ├── config/                        # Configuration management
│   ├── logger/                        # Logging utilities
│   └── validator/                     # Validation helpers
│
└── migrations/                        # golang-migrate database migrations
    ├── 000001_create_users_table.up.sql
    ├── 000001_create_users_table.down.sql
    ├── 000002_create_orders_table.up.sql
    └── 000002_create_orders_table.down.sql
```

## Domain Layer (DDD)

The domain layer is the heart of the application. It contains enterprise-wide business rules and has NO external dependencies.

### Entities and Aggregates

Entities have identity and lifecycle. An Aggregate is a cluster of entities treated as a single unit, with an Aggregate Root as the entry point.

```go
// internal/domain/user/user.go
package user

import (
    "errors"
    "time"

    "github.com/google/uuid"
)

// User is the aggregate root for user-related operations.
// Fields are private to enforce invariants through methods.
type User struct {
    id        uuid.UUID
    email     Email      // Value Object
    name      string
    role      Role       // Value Object
    createdAt time.Time
    updatedAt time.Time
}

// NewUser creates a new User entity with validation.
// This is the only way to create a valid User.
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

// Reconstitute creates a User from persisted data (bypasses validation).
// Used only by repositories when loading from database.
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

// Business methods enforce invariants
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
func (u *User) ID() uuid.UUID        { return u.id }
func (u *User) Email() Email         { return u.email }
func (u *User) Name() string         { return u.name }
func (u *User) Role() Role           { return u.role }
func (u *User) CreatedAt() time.Time { return u.createdAt }
func (u *User) UpdatedAt() time.Time { return u.updatedAt }
```

### Value Objects

Value Objects are immutable and defined by their attributes, not identity.

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

Repository interfaces are defined in the domain layer but implemented in infrastructure. This is the "Port" in Ports & Adapters (Hexagonal Architecture).

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

// ReadRepository for CQRS query side (optional, for complex read models)
type ReadRepository interface {
    List(ctx context.Context, opts ListOptions) ([]*User, int, error)
    Search(ctx context.Context, query string) ([]*User, error)
}
```

### Domain Errors

```go
// internal/domain/user/errors.go
package user

import "errors"

var (
    ErrInvalidEmail      = errors.New("invalid email address")
    ErrInvalidName       = errors.New("name cannot be empty")
    ErrUserNotFound      = errors.New("user not found")
    ErrEmailAlreadyExists = errors.New("email already exists")
)

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
```

## Application Layer (CQRS)

The application layer implements use cases using CQRS pattern: Commands for writes, Queries for reads.

### Commands (Write Operations)

Commands change state. They validate input, orchestrate domain operations, and persist changes.

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

Queries retrieve data without changing state. They can use optimized read models.

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

// internal/application/query/list_users.go
package query

type ListUsersQuery struct {
    Page     int
    PageSize int
    Role     string
}

type ListUsersHandler struct {
    repo user.ReadRepository
}

func NewListUsersHandler(repo user.ReadRepository) *ListUsersHandler {
    return &ListUsersHandler{repo: repo}
}

func (h *ListUsersHandler) Handle(ctx context.Context, q ListUsersQuery) (*dto.PaginatedResponse[dto.UserResponse], error) {
    opts := user.ListOptions{
        Offset: (q.Page - 1) * q.PageSize,
        Limit:  q.PageSize,
        Role:   q.Role,
    }

    users, total, err := h.repo.List(ctx, opts)
    if err != nil {
        return nil, err
    }

    items := make([]dto.UserResponse, len(users))
    for i, u := range users {
        items[i] = *dto.UserResponseFromDomain(u)
    }

    return &dto.PaginatedResponse[dto.UserResponse]{
        Items: items,
        Total: total,
        Page:  q.Page,
    }, nil
}
```

### Data Transfer Objects (DTOs)

DTOs transform domain objects for external communication.

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

type PaginatedResponse[T any] struct {
    Items []T `json:"items"`
    Total int `json:"total"`
    Page  int `json:"page"`
}
```

## Infrastructure Layer

Implements interfaces defined in the domain layer (Adapters).

### PostgreSQL Repository Implementation

```go
// internal/infrastructure/persistence/postgres/user_repository.go
package postgres

import (
    "context"
    "errors"
    "fmt"
    "time"

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
        dbID      uuid.UUID
        email     string
        name      string
        role      string
        createdAt time.Time
        updatedAt time.Time
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
        if isUniqueViolation(err) {
            return user.ErrEmailAlreadyExists
        }
        return fmt.Errorf("save user: %w", err)
    }

    return nil
}
```

### Unit of Work (Transaction Management)

```go
// internal/infrastructure/persistence/postgres/unit_of_work.go
package postgres

import (
    "context"
    "fmt"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

type UnitOfWork struct {
    db *pgxpool.Pool
}

func NewUnitOfWork(db *pgxpool.Pool) *UnitOfWork {
    return &UnitOfWork{db: db}
}

func (uow *UnitOfWork) Execute(ctx context.Context, fn func(ctx context.Context, tx pgx.Tx) error) error {
    tx, err := uow.db.Begin(ctx)
    if err != nil {
        return fmt.Errorf("begin transaction: %w", err)
    }
    defer tx.Rollback(ctx)

    if err := fn(ctx, tx); err != nil {
        return err
    }

    if err := tx.Commit(ctx); err != nil {
        return fmt.Errorf("commit transaction: %w", err)
    }

    return nil
}
```

## Interfaces Layer (HTTP)

The outermost layer adapts HTTP requests to application commands/queries.

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
    updateUser *command.UpdateUserHandler
    getUser    *query.GetUserHandler
    listUsers  *query.ListUsersHandler
}

func NewUserHandler(
    createUser *command.CreateUserHandler,
    updateUser *command.UpdateUserHandler,
    getUser *query.GetUserHandler,
    listUsers *query.ListUsersHandler,
) *UserHandler {
    return &UserHandler{
        createUser: createUser,
        updateUser: updateUser,
        getUser:    getUser,
        listUsers:  listUsers,
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
    case errors.Is(err, shared.ErrConflict), errors.Is(err, user.ErrEmailAlreadyExists):
        respondError(w, http.StatusConflict, "resource already exists")
    case errors.Is(err, user.ErrInvalidEmail), errors.Is(err, user.ErrInvalidName):
        respondError(w, http.StatusBadRequest, err.Error())
    default:
        respondError(w, http.StatusInternalServerError, "internal server error")
    }
}
```

## Dependency Injection (Wiring)

```go
// cmd/api/main.go
package main

func main() {
    // Load configuration
    cfg := config.Load()

    // Infrastructure
    db := postgres.NewConnection(cfg.DatabaseURL)
    userRepo := postgres.NewUserRepository(db)
    userReadRepo := postgres.NewUserReadRepository(db)
    uow := postgres.NewUnitOfWork(db)
    logger := logger.New(cfg.LogLevel)

    // Application - Command Handlers
    createUserHandler := command.NewCreateUserHandler(userRepo, logger)
    updateUserHandler := command.NewUpdateUserHandler(userRepo, uow, logger)
    deleteUserHandler := command.NewDeleteUserHandler(userRepo, logger)

    // Application - Query Handlers
    getUserHandler := query.NewGetUserHandler(userRepo)
    listUsersHandler := query.NewListUsersHandler(userReadRepo)

    // Interfaces - HTTP Handlers
    userHTTPHandler := handler.NewUserHandler(
        createUserHandler,
        updateUserHandler,
        getUserHandler,
        listUsersHandler,
    )

    // Router
    r := chi.NewRouter()
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)

    r.Route("/api/v1", func(r chi.Router) {
        r.Route("/users", func(r chi.Router) {
            r.Post("/", userHTTPHandler.Create)
            r.Get("/", userHTTPHandler.List)
            r.Get("/{id}", userHTTPHandler.Get)
            r.Put("/{id}", userHTTPHandler.Update)
            r.Delete("/{id}", userHTTPHandler.Delete)
        })
    })

    // Start server
    server := &http.Server{
        Addr:    cfg.ServerAddress,
        Handler: r,
    }

    logger.Info("server starting", "addr", cfg.ServerAddress)
    if err := server.ListenAndServe(); err != nil {
        logger.Error("server error", "error", err)
    }
}
```

## Database Migrations (golang-migrate)

We use golang-migrate CLI for database migrations.

### Migration Commands

```bash
# Install golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Create a new migration (creates .up.sql and .down.sql files)
migrate create -ext sql -dir backend/migrations -seq create_users_table

# Apply all pending migrations
migrate -path backend/migrations -database "$DATABASE_URL" up

# Rollback the last migration
migrate -path backend/migrations -database "$DATABASE_URL" down 1

# Check migration version
migrate -path backend/migrations -database "$DATABASE_URL" version

# Force set version (for fixing dirty state)
migrate -path backend/migrations -database "$DATABASE_URL" force 1
```

### Migration Example

**Up Migration:**
```sql
-- migrations/000001_create_users_table.up.sql

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'member',
    email_verified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT users_email_unique UNIQUE (email)
);

CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_role ON users(role) WHERE deleted_at IS NULL;
```

**Down Migration:**
```sql
-- migrations/000001_create_users_table.down.sql

DROP TABLE IF EXISTS users;
```

## Frontend Architecture (React)

### Directory Structure

```
frontend/
├── public/                      # Static assets
├── src/
│   ├── app/                     # App router pages
│   │   ├── (auth)/              # Auth route group
│   │   │   ├── login/
│   │   │   └── register/
│   │   ├── (dashboard)/         # Dashboard route group
│   │   │   ├── layout.tsx
│   │   │   └── page.tsx
│   │   ├── layout.tsx           # Root layout
│   │   └── page.tsx             # Home page
│   ├── components/
│   │   ├── ui/                  # shadcn/ui components
│   │   ├── features/            # Feature components
│   │   ├── layouts/             # Layout components
│   │   └── shared/              # Shared components
│   ├── hooks/                   # Custom hooks
│   │   ├── use-auth.ts
│   │   └── use-media-query.ts
│   ├── lib/                     # Utilities
│   │   ├── api/                 # API client
│   │   ├── utils.ts             # Helper functions
│   │   └── validations/         # Zod schemas
│   ├── stores/                  # Zustand stores
│   │   ├── auth-store.ts
│   │   └── ui-store.ts
│   ├── types/                   # TypeScript types
│   │   ├── api.ts
│   │   └── models.ts
│   └── styles/
│       └── globals.css          # Global styles (Tailwind v4)
├── tailwind.config.ts
└── tsconfig.json
```

### State Management

#### Server State (TanStack Query)

```typescript
// src/lib/api/users.ts
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from './client';

export const userKeys = {
  all: ['users'] as const,
  lists: () => [...userKeys.all, 'list'] as const,
  list: (filters: UserFilters) => [...userKeys.lists(), filters] as const,
  details: () => [...userKeys.all, 'detail'] as const,
  detail: (id: string) => [...userKeys.details(), id] as const,
};

export function useUsers(filters: UserFilters) {
  return useQuery({
    queryKey: userKeys.list(filters),
    queryFn: () => api.get<User[]>('/users', { params: filters }),
  });
}

export function useCreateUser() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateUserInput) => api.post<User>('/users', data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: userKeys.lists() });
    },
  });
}
```

#### Client State (Zustand)

```typescript
// src/stores/auth-store.ts
import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  login: (user: User, token: string) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      token: null,
      isAuthenticated: false,
      login: (user, token) =>
        set({ user, token, isAuthenticated: true }),
      logout: () =>
        set({ user: null, token: null, isAuthenticated: false }),
    }),
    {
      name: 'auth-storage',
    }
  )
);
```

## Data Flow

### CQRS Request Flow

```
┌─────────────────────────────────────────────────────────────┐
│                     WRITE (Commands)                        │
├─────────────────────────────────────────────────────────────┤
│  User Action                                                │
│       │                                                     │
│       ▼                                                     │
│  HTTP POST/PUT/DELETE                                       │
│       │                                                     │
│       ▼                                                     │
│  HTTP Handler → Command Handler → Domain Entity → Repository│
│                      │                                      │
│                      ▼                                      │
│              Validates business rules                       │
│              Mutates aggregate state                        │
│              Persists to database                           │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                      READ (Queries)                         │
├─────────────────────────────────────────────────────────────┤
│  User Action                                                │
│       │                                                     │
│       ▼                                                     │
│  HTTP GET                                                   │
│       │                                                     │
│       ▼                                                     │
│  HTTP Handler → Query Handler → Read Repository → DTO       │
│                      │                                      │
│                      ▼                                      │
│              Can use optimized read models                  │
│              No domain logic validation                     │
│              Returns DTOs directly                          │
└─────────────────────────────────────────────────────────────┘
```

## Security Considerations

### Authentication
- JWT tokens with short expiry (15 minutes)
- Refresh token rotation
- Secure HTTP-only cookies for web clients
- Rate limiting on auth endpoints

### API Security
- Input validation at handler level
- SQL parameterized queries (no string concatenation)
- CORS configuration for allowed origins
- Request size limits
- Timeout middleware

### Data Protection
- Passwords hashed with bcrypt (cost 12)
- Sensitive data encrypted at rest
- TLS for all connections
- Audit logging for sensitive operations

## Performance Optimizations

### Backend
- Connection pooling for database (pgxpool)
- Redis caching for frequently accessed data
- Pagination for list endpoints
- Database indexing on query patterns
- Graceful shutdown handling

### Frontend
- Code splitting with dynamic imports
- Image optimization
- TanStack Query caching and deduplication
- Virtualized lists for large datasets
- Optimistic updates for better UX

## References

- [Clean Architecture by Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Three Dots Labs - DDD, CQRS, Clean Architecture](https://threedots.tech/post/ddd-cqrs-clean-architecture-combined/)
- [Wild Workouts Go DDD Example](https://github.com/ThreeDotsLabs/wild-workouts-go-ddd-example)
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [shadcn/ui](https://ui.shadcn.com/)
- [Tailwind CSS v4](https://tailwindcss.com/docs/v4-beta)
