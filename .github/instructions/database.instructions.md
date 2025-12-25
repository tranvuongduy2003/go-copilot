---
applyTo: "backend/migrations/**/*.sql,backend/internal/infrastructure/persistence/**/*.go"
---

# Database Development Instructions

These instructions cover database migrations using Goose CLI and repository implementations.

## Goose CLI Setup

### Installation

```bash
# Install goose CLI
go install github.com/pressly/goose/v3/cmd/goose@latest

# Verify installation
goose --version
```

### Environment Configuration

Create a `.env` file or set environment variables:

```bash
# Environment variables for goose
export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING="postgres://user:password@localhost:5432/dbname?sslmode=disable"
export GOOSE_MIGRATION_DIR=./migrations/sql
```

Or use a `.env` file:

```env
GOOSE_DRIVER=postgres
GOOSE_DBSTRING=postgres://user:password@localhost:5432/dbname?sslmode=disable
GOOSE_MIGRATION_DIR=./migrations/sql
```

## Migration File Structure

```
backend/
└── migrations/
    └── sql/
        ├── 00001_create_users.sql
        ├── 00002_add_user_roles.sql
        ├── 00003_create_posts.sql
        └── 00004_add_posts_indexes.sql
```

## Goose CLI Commands

### Create New Migration

```bash
# Create a new SQL migration
goose -dir migrations/sql create create_users sql

# Creates: migrations/sql/20241225120000_create_users.sql
```

### Run Migrations

```bash
# Apply all pending migrations
goose -dir migrations/sql postgres "postgres://..." up

# With environment variables
goose up

# Apply only one migration
goose up-by-one

# Apply up to a specific version
goose up-to 20241225120000
```

### Rollback Migrations

```bash
# Rollback the last migration
goose down

# Rollback to a specific version
goose down-to 20241225100000

# Rollback all migrations (DANGER!)
goose reset
```

### Check Status

```bash
# Show migration status
goose status

# Output:
#     Applied At                  Migration
#     =======================================
#     Mon Jan 01 00:00:00 2024    00001_create_users.sql
#     Mon Jan 01 00:00:01 2024    00002_add_user_roles.sql
#     Pending                     00003_create_posts.sql
```

### Other Commands

```bash
# Show current version
goose version

# Redo the last migration (down then up)
goose redo

# Validate migration files
goose validate
```

## Migration Template

### Standard Migration

```sql
-- +goose Up
-- +goose StatementBegin
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

-- Create indexes
CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_role ON users(role) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_created_at ON users(created_at DESC);

-- Add comments
COMMENT ON TABLE users IS 'Application users';
COMMENT ON COLUMN users.deleted_at IS 'Soft delete timestamp';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
```

### No Transaction Migration

Some statements cannot run inside a transaction (e.g., `CREATE INDEX CONCURRENTLY`):

```sql
-- +goose NO TRANSACTION

-- +goose Up
CREATE INDEX CONCURRENTLY idx_posts_title_search ON posts USING gin(to_tsvector('english', title));

-- +goose Down
DROP INDEX CONCURRENTLY IF EXISTS idx_posts_title_search;
```

### Data Migration

```sql
-- +goose Up
-- +goose StatementBegin
-- Add new column
ALTER TABLE users ADD COLUMN status VARCHAR(50);

-- Backfill data
UPDATE users SET status = 'active' WHERE deleted_at IS NULL;
UPDATE users SET status = 'deleted' WHERE deleted_at IS NOT NULL;

-- Make column NOT NULL after backfill
ALTER TABLE users ALTER COLUMN status SET NOT NULL;
ALTER TABLE users ALTER COLUMN status SET DEFAULT 'active';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN status;
-- +goose StatementEnd
```

## Migration Best Practices

### 1. Keep Migrations Small and Focused

```sql
-- GOOD: Single responsibility
-- 00001_create_users.sql - Only creates users table
-- 00002_create_posts.sql - Only creates posts table
-- 00003_add_posts_user_fk.sql - Only adds foreign key

-- BAD: Too many changes in one migration
-- 00001_initial_schema.sql - Creates 10 tables, all indexes, all constraints
```

### 2. Always Provide Reversible Down Migrations

```sql
-- +goose Up
ALTER TABLE users ADD COLUMN avatar_url TEXT;

-- +goose Down
ALTER TABLE users DROP COLUMN avatar_url;
```

### 3. Use Transactions (Default Behavior)

Goose runs each migration in a transaction by default. This ensures atomicity.

```sql
-- +goose Up
-- +goose StatementBegin
-- All statements here run in a single transaction
CREATE TABLE posts (...);
CREATE INDEX idx_posts_user ON posts(user_id);
-- +goose StatementEnd
```

### 4. Use Partial Indexes for Soft Deletes

```sql
-- Only index non-deleted records
CREATE UNIQUE INDEX idx_users_email_active
ON users(email)
WHERE deleted_at IS NULL;
```

### 5. Add Foreign Key Indexes

```sql
-- Always index foreign keys
CREATE TABLE posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    -- ...
);

CREATE INDEX idx_posts_user_id ON posts(user_id);
```

### 6. Use Descriptive Constraint Names

```sql
ALTER TABLE posts
ADD CONSTRAINT posts_user_id_fkey
FOREIGN KEY (user_id) REFERENCES users(id);

ALTER TABLE posts
ADD CONSTRAINT posts_status_check
CHECK (status IN ('draft', 'published', 'archived'));
```

## Makefile Integration

Add these targets to your Makefile:

```makefile
# Database migration commands
.PHONY: migrate-up migrate-down migrate-status migrate-create

GOOSE_FLAGS = -dir migrations/sql

migrate-up:
	goose $(GOOSE_FLAGS) postgres "$(DATABASE_URL)" up

migrate-down:
	goose $(GOOSE_FLAGS) postgres "$(DATABASE_URL)" down

migrate-status:
	goose $(GOOSE_FLAGS) postgres "$(DATABASE_URL)" status

migrate-create:
	@read -p "Migration name: " name; \
	goose $(GOOSE_FLAGS) create $$name sql

migrate-reset:
	goose $(GOOSE_FLAGS) postgres "$(DATABASE_URL)" reset

migrate-redo:
	goose $(GOOSE_FLAGS) postgres "$(DATABASE_URL)" redo
```

Usage:

```bash
make migrate-up
make migrate-down
make migrate-status
make migrate-create  # Prompts for migration name
```

## Programmatic Migration (Embedded)

For running migrations from Go code:

```go
package main

import (
    "database/sql"
    "embed"
    "log"

    "github.com/pressly/goose/v3"
    _ "github.com/lib/pq"
)

//go:embed migrations/sql/*.sql
var embedMigrations embed.FS

func main() {
    db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    goose.SetBaseFS(embedMigrations)

    if err := goose.SetDialect("postgres"); err != nil {
        log.Fatal(err)
    }

    if err := goose.Up(db, "migrations/sql"); err != nil {
        log.Fatal(err)
    }
}
```

## Repository Implementation (DDD Style)

### Domain Repository Interface

```go
// internal/domain/user/repository.go
package user

import (
    "context"
    "github.com/google/uuid"
)

// Repository is the port for user persistence.
// Defined in domain, implemented in infrastructure.
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

### PostgreSQL Implementation

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
        return nil, fmt.Errorf("find user by id: %w", err)
    }

    // Reconstitute domain entity
    emailVO, _ := user.NewEmail(email)
    return user.Reconstitute(dbID, emailVO, name, user.Role(role), createdAt, updatedAt), nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email user.Email) (*user.User, error) {
    query := `
        SELECT id, email, name, role, created_at, updated_at
        FROM users
        WHERE LOWER(email) = LOWER($1) AND deleted_at IS NULL
    `

    var (
        dbID      uuid.UUID
        dbEmail   string
        name      string
        role      string
        createdAt time.Time
        updatedAt time.Time
    )

    err := r.db.QueryRow(ctx, query, email.String()).Scan(
        &dbID, &dbEmail, &name, &role, &createdAt, &updatedAt,
    )

    if errors.Is(err, pgx.ErrNoRows) {
        return nil, user.ErrUserNotFound
    }
    if err != nil {
        return nil, fmt.Errorf("find user by email: %w", err)
    }

    emailVO, _ := user.NewEmail(dbEmail)
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

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
    query := `
        UPDATE users
        SET deleted_at = $2, updated_at = $2
        WHERE id = $1 AND deleted_at IS NULL
    `

    now := time.Now().UTC()
    result, err := r.db.Exec(ctx, query, id, now)
    if err != nil {
        return fmt.Errorf("delete user: %w", err)
    }

    if result.RowsAffected() == 0 {
        return user.ErrUserNotFound
    }

    return nil
}
```

## Transaction Handling

### Unit of Work Pattern

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

### Usage in Application Layer

```go
// internal/application/command/transfer_funds.go
func (h *TransferFundsHandler) Handle(ctx context.Context, cmd TransferFundsCommand) error {
    return h.uow.Execute(ctx, func(ctx context.Context, tx pgx.Tx) error {
        // All operations in this function run in a transaction

        from, err := h.accountRepo.FindByIDForUpdate(ctx, tx, cmd.FromAccountID)
        if err != nil {
            return err
        }

        to, err := h.accountRepo.FindByIDForUpdate(ctx, tx, cmd.ToAccountID)
        if err != nil {
            return err
        }

        if err := from.Debit(cmd.Amount); err != nil {
            return err
        }

        to.Credit(cmd.Amount)

        if err := h.accountRepo.SaveWithTx(ctx, tx, from); err != nil {
            return err
        }

        if err := h.accountRepo.SaveWithTx(ctx, tx, to); err != nil {
            return err
        }

        return nil
    })
}
```

## Query Optimization

### Use Row Locking for Updates

```go
func (r *Repository) FindByIDForUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*Account, error) {
    query := `
        SELECT id, balance, updated_at
        FROM accounts
        WHERE id = $1
        FOR UPDATE
    `

    // FOR UPDATE locks the row until transaction commits
    var account Account
    err := tx.QueryRow(ctx, query, id).Scan(&account.ID, &account.Balance, &account.UpdatedAt)
    // ...
}
```

### Batch Operations

```go
func (r *Repository) CreateMany(ctx context.Context, users []*user.User) error {
    batch := &pgx.Batch{}

    for _, u := range users {
        batch.Queue(
            `INSERT INTO users (id, email, name, role, created_at, updated_at)
             VALUES ($1, $2, $3, $4, $5, $6)`,
            u.ID(), u.Email().String(), u.Name(), string(u.Role()), u.CreatedAt(), u.UpdatedAt(),
        )
    }

    results := r.db.SendBatch(ctx, batch)
    defer results.Close()

    for range users {
        if _, err := results.Exec(); err != nil {
            return fmt.Errorf("batch insert: %w", err)
        }
    }

    return nil
}
```

## References

- [Goose Documentation](https://github.com/pressly/goose)
- [Goose Go Package](https://pkg.go.dev/github.com/pressly/goose/v3)
- [PostgreSQL Best Practices](https://wiki.postgresql.org/wiki/Don't_Do_This)
