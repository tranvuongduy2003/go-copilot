---
applyTo: "backend/**/*.go"
---

# Go Backend Development Instructions

These instructions apply to all Go files in the backend directory.

## Project Structure

Follow clean architecture with these layers:

```
backend/
├── cmd/api/main.go           # Entry point
├── internal/
│   ├── config/               # Configuration management
│   ├── domain/               # Domain models and errors
│   ├── handlers/             # HTTP handlers
│   ├── middleware/           # HTTP middleware
│   ├── repository/           # Data access layer
│   │   └── postgres/         # PostgreSQL implementation
│   └── service/              # Business logic
├── migrations/               # Database migrations
└── pkg/                      # Shared packages
    ├── response/             # HTTP response helpers
    └── validator/            # Input validation
```

## Coding Standards

### Package Naming
- Use short, lowercase names: `user`, `auth`, `config`
- No underscores or mixed caps
- Package name should describe purpose

### File Naming
- Use `snake_case.go`: `user_service.go`, `auth_middleware.go`
- Test files: `user_service_test.go`
- Keep files focused; split large files

### Function and Variable Naming
```go
// Exported (public) - PascalCase
func CreateUser(ctx context.Context, input CreateUserInput) (*User, error)

// Unexported (private) - camelCase
func hashPassword(password string) (string, error)

// Constants - PascalCase for exported, camelCase for unexported
const MaxRetries = 3
const defaultTimeout = 30 * time.Second
```

## Error Handling

### Define Domain Errors
```go
package domain

import "errors"

var (
    ErrNotFound     = errors.New("resource not found")
    ErrConflict     = errors.New("resource already exists")
    ErrInvalidInput = errors.New("invalid input")
    ErrUnauthorized = errors.New("unauthorized")
    ErrForbidden    = errors.New("forbidden")
)
```

### Wrap Errors with Context
```go
// Always wrap errors with context
result, err := s.repo.FindByID(ctx, id)
if err != nil {
    if errors.Is(err, repository.ErrNotFound) {
        return nil, ErrNotFound
    }
    return nil, fmt.Errorf("failed to find user by id %s: %w", id, err)
}

// Never ignore errors
result, _ := doSomething() // BAD
```

### Handle Errors at Boundaries
```go
func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
    user, err := h.svc.GetUser(ctx, id)
    if err != nil {
        switch {
        case errors.Is(err, domain.ErrNotFound):
            response.NotFound(w, "User not found")
        case errors.Is(err, domain.ErrForbidden):
            response.Forbidden(w, "Access denied")
        default:
            response.InternalError(w, err)
        }
        return
    }
    response.JSON(w, http.StatusOK, user)
}
```

## Context Usage

### Always Pass Context
```go
// Context as first parameter
func (s *Service) GetUser(ctx context.Context, id string) (*User, error)

// Pass context to all downstream calls
func (s *Service) CreateUser(ctx context.Context, input CreateUserInput) (*User, error) {
    // Check context for cancellation
    if err := ctx.Err(); err != nil {
        return nil, err
    }

    // Pass to repository
    return s.repo.Create(ctx, user)
}
```

### Context Values for Request-Scoped Data
```go
type contextKey string

const userIDKey contextKey = "user_id"

func SetUserID(ctx context.Context, id string) context.Context {
    return context.WithValue(ctx, userIDKey, id)
}

func GetUserID(ctx context.Context) string {
    id, _ := ctx.Value(userIDKey).(string)
    return id
}
```

## Repository Pattern

### Define Interface at Point of Use
```go
// In service package, not repository package
package service

type UserRepository interface {
    FindByID(ctx context.Context, id string) (*domain.User, error)
    FindByEmail(ctx context.Context, email string) (*domain.User, error)
    Create(ctx context.Context, user *domain.User) error
    Update(ctx context.Context, user *domain.User) error
    Delete(ctx context.Context, id string) error
}
```

### Implement with Specific Database
```go
package postgres

type userRepository struct {
    db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *userRepository {
    return &userRepository{db: db}
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
    query := `
        SELECT id, email, name, created_at, updated_at
        FROM users
        WHERE id = $1 AND deleted_at IS NULL
    `

    var user domain.User
    err := r.db.QueryRow(ctx, query, id).Scan(
        &user.ID,
        &user.Email,
        &user.Name,
        &user.CreatedAt,
        &user.UpdatedAt,
    )

    if errors.Is(err, pgx.ErrNoRows) {
        return nil, domain.ErrNotFound
    }
    if err != nil {
        return nil, fmt.Errorf("query user by id: %w", err)
    }

    return &user, nil
}
```

## Service Layer

### Business Logic Lives Here
```go
package service

type UserService struct {
    repo   UserRepository
    hasher PasswordHasher
    logger *slog.Logger
}

func NewUserService(repo UserRepository, hasher PasswordHasher, logger *slog.Logger) *UserService {
    return &UserService{
        repo:   repo,
        hasher: hasher,
        logger: logger,
    }
}

func (s *UserService) CreateUser(ctx context.Context, input domain.CreateUserInput) (*domain.User, error) {
    // Business rule: Check for duplicate email
    existing, err := s.repo.FindByEmail(ctx, input.Email)
    if err != nil && !errors.Is(err, domain.ErrNotFound) {
        return nil, fmt.Errorf("check existing email: %w", err)
    }
    if existing != nil {
        return nil, domain.ErrConflict
    }

    // Business logic: Hash password
    hash, err := s.hasher.Hash(input.Password)
    if err != nil {
        return nil, fmt.Errorf("hash password: %w", err)
    }

    now := time.Now().UTC()
    user := &domain.User{
        ID:           uuid.New().String(),
        Email:        input.Email,
        Name:         input.Name,
        PasswordHash: hash,
        CreatedAt:    now,
        UpdatedAt:    now,
    }

    if err := s.repo.Create(ctx, user); err != nil {
        return nil, fmt.Errorf("create user: %w", err)
    }

    s.logger.Info("user created", "id", user.ID, "email", user.Email)

    return user, nil
}
```

## Logging

### Use Structured Logging (slog)
```go
import "log/slog"

// Create logger
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelInfo,
}))

// Log with context
logger.Info("user created",
    "user_id", user.ID,
    "email", user.Email,
)

// Log errors
logger.Error("failed to create user",
    "error", err,
    "email", input.Email,
)

// Never log sensitive data
logger.Info("login attempt", "email", email) // OK
logger.Info("login attempt", "password", password) // NEVER
```

## Configuration

### Use Environment Variables
```go
package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
    Port         int    `envconfig:"PORT" default:"8080"`
    Environment  string `envconfig:"ENVIRONMENT" default:"development"`
    DatabaseURL  string `envconfig:"DATABASE_URL" required:"true"`
    JWTSecret    string `envconfig:"JWT_SECRET" required:"true"`
    LogLevel     string `envconfig:"LOG_LEVEL" default:"info"`
}

func Load() (*Config, error) {
    var cfg Config
    if err := envconfig.Process("", &cfg); err != nil {
        return nil, fmt.Errorf("load config: %w", err)
    }
    return &cfg, nil
}
```

## Input Validation

### Validate at Handler Level
```go
import "github.com/go-playground/validator/v10"

var validate = validator.New()

type CreateUserInput struct {
    Email    string `json:"email" validate:"required,email,max=255"`
    Name     string `json:"name" validate:"required,min=2,max=100"`
    Password string `json:"password" validate:"required,min=8,max=128"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
    var input CreateUserInput
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        response.BadRequest(w, "Invalid request body")
        return
    }

    if err := validate.Struct(input); err != nil {
        response.ValidationError(w, formatValidationErrors(err))
        return
    }

    // Input is valid, proceed
}
```

## Database Queries

### Use Parameterized Queries
```go
// CORRECT: Parameterized
query := "SELECT * FROM users WHERE email = $1"
row := db.QueryRow(ctx, query, email)

// WRONG: String concatenation (SQL injection risk)
query := "SELECT * FROM users WHERE email = '" + email + "'"
```

### Handle Transactions
```go
func (r *Repository) Transfer(ctx context.Context, fromID, toID string, amount int) error {
    tx, err := r.db.Begin(ctx)
    if err != nil {
        return fmt.Errorf("begin transaction: %w", err)
    }
    defer tx.Rollback(ctx) // Rollback if not committed

    // Debit from account
    if _, err := tx.Exec(ctx,
        "UPDATE accounts SET balance = balance - $1 WHERE id = $2",
        amount, fromID,
    ); err != nil {
        return fmt.Errorf("debit account: %w", err)
    }

    // Credit to account
    if _, err := tx.Exec(ctx,
        "UPDATE accounts SET balance = balance + $1 WHERE id = $2",
        amount, toID,
    ); err != nil {
        return fmt.Errorf("credit account: %w", err)
    }

    if err := tx.Commit(ctx); err != nil {
        return fmt.Errorf("commit transaction: %w", err)
    }

    return nil
}
```

## Testing

### Table-Driven Tests
```go
func TestUserService_CreateUser(t *testing.T) {
    tests := []struct {
        name    string
        input   CreateUserInput
        setup   func(*MockRepo)
        wantErr error
    }{
        {
            name:  "success",
            input: CreateUserInput{Email: "test@example.com", Name: "Test"},
            setup: func(m *MockRepo) {
                m.On("FindByEmail", mock.Anything, "test@example.com").Return(nil, ErrNotFound)
                m.On("Create", mock.Anything, mock.Anything).Return(nil)
            },
            wantErr: nil,
        },
        {
            name:  "duplicate email",
            input: CreateUserInput{Email: "exists@example.com", Name: "Test"},
            setup: func(m *MockRepo) {
                m.On("FindByEmail", mock.Anything, "exists@example.com").Return(&User{}, nil)
            },
            wantErr: ErrConflict,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            repo := new(MockRepo)
            tt.setup(repo)
            svc := NewUserService(repo, hasher, logger)

            _, err := svc.CreateUser(context.Background(), tt.input)

            if tt.wantErr != nil {
                assert.ErrorIs(t, err, tt.wantErr)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```
