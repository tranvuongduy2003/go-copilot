---
applyTo: "backend/migrations/**/*.sql,backend/internal/infrastructure/persistence/**/*.go"
---

# Database Development Instructions

These instructions cover database migrations using golang-migrate CLI and repository implementations.

## golang-migrate CLI Setup

### Installation

```bash
# Install golang-migrate CLI
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Verify installation
migrate -version
```

### Environment Configuration

Set the DATABASE_URL environment variable:

```bash
export DATABASE_URL="postgres://user:password@localhost:5432/dbname?sslmode=disable"
```

Or use a `.env` file:

```env
DATABASE_URL=postgres://user:password@localhost:5432/dbname?sslmode=disable
```

## Migration File Structure

golang-migrate uses separate `.up.sql` and `.down.sql` files:

```
backend/
└── migrations/
    ├── 000001_create_users_table.up.sql
    ├── 000001_create_users_table.down.sql
    ├── 000002_add_users_indexes.up.sql
    ├── 000002_add_users_indexes.down.sql
    ├── 000003_create_posts_table.up.sql
    └── 000003_create_posts_table.down.sql
```

## golang-migrate CLI Commands

### Create New Migration

```bash
# Create a new migration (creates .up.sql and .down.sql files)
migrate create -ext sql -dir backend/migrations -seq create_users_table

# Creates:
# - backend/migrations/000001_create_users_table.up.sql
# - backend/migrations/000001_create_users_table.down.sql
```

### Run Migrations

```bash
# Apply all pending migrations
migrate -path backend/migrations -database "$DATABASE_URL" up

# Apply only N migrations
migrate -path backend/migrations -database "$DATABASE_URL" up 2
```

### Rollback Migrations

```bash
# Rollback the last migration
migrate -path backend/migrations -database "$DATABASE_URL" down 1

# Rollback all migrations (DANGER!)
migrate -path backend/migrations -database "$DATABASE_URL" down
```

### Check Status

```bash
# Show current version
migrate -path backend/migrations -database "$DATABASE_URL" version

# Force set version (for fixing dirty state)
migrate -path backend/migrations -database "$DATABASE_URL" force 1
```

### Other Commands

```bash
# Go to a specific version
migrate -path backend/migrations -database "$DATABASE_URL" goto 3

# Drop everything (DANGER!)
migrate -path backend/migrations -database "$DATABASE_URL" drop -f
```

## Migration Template

### Standard Migration (Up)

```sql
-- backend/migrations/000001_create_users_table.up.sql

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
```

### Standard Migration (Down)

```sql
-- backend/migrations/000001_create_users_table.down.sql

DROP TABLE IF EXISTS users;
```

### Concurrent Index Migration

Some statements cannot run inside a transaction (e.g., `CREATE INDEX CONCURRENTLY`). Use separate migration files:

```sql
-- backend/migrations/000002_add_search_index.up.sql
CREATE INDEX CONCURRENTLY idx_posts_title_search ON posts USING gin(to_tsvector('english', title));
```

```sql
-- backend/migrations/000002_add_search_index.down.sql
DROP INDEX CONCURRENTLY IF EXISTS idx_posts_title_search;
```

### Data Migration

```sql
-- backend/migrations/000003_add_user_status.up.sql

-- Add new column
ALTER TABLE users ADD COLUMN status VARCHAR(50);

-- Backfill data
UPDATE users SET status = 'active' WHERE deleted_at IS NULL;
UPDATE users SET status = 'deleted' WHERE deleted_at IS NOT NULL;

-- Make column NOT NULL after backfill
ALTER TABLE users ALTER COLUMN status SET NOT NULL;
ALTER TABLE users ALTER COLUMN status SET DEFAULT 'active';
```

```sql
-- backend/migrations/000003_add_user_status.down.sql

ALTER TABLE users DROP COLUMN status;
```

## Migration Best Practices

### 1. Keep Migrations Small and Focused

```sql
-- GOOD: Single responsibility
-- 000001_create_users_table.up.sql - Only creates users table
-- 000002_create_posts_table.up.sql - Only creates posts table
-- 000003_add_posts_user_fk.up.sql - Only adds foreign key

-- BAD: Too many changes in one migration
-- 000001_initial_schema.up.sql - Creates 10 tables, all indexes, all constraints
```

### 2. Always Provide Reversible Down Migrations

Each `.up.sql` file needs a corresponding `.down.sql` file:

```sql
-- 000004_add_avatar_url.up.sql
ALTER TABLE users ADD COLUMN avatar_url TEXT;
```

```sql
-- 000004_add_avatar_url.down.sql
ALTER TABLE users DROP COLUMN avatar_url;
```

### 3. Transactions

golang-migrate runs each migration in a transaction by default. For statements that can't run in a transaction (like `CREATE INDEX CONCURRENTLY`), they should be in separate migration files.

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
.PHONY: migrate-up migrate-down migrate-version migrate-create

migrate-up:
	migrate -path backend/migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path backend/migrations -database "$(DATABASE_URL)" down 1

migrate-version:
	migrate -path backend/migrations -database "$(DATABASE_URL)" version

migrate-create:
	@read -p "Migration name: " name; \
	migrate create -ext sql -dir backend/migrations -seq $$name

migrate-force:
	@read -p "Version: " version; \
	migrate -path backend/migrations -database "$(DATABASE_URL)" force $$version
```

Usage:

```bash
make migrate-up
make migrate-down
make migrate-version
make migrate-create  # Prompts for migration name
```

## Persistence Layer Structure

The persistence layer is organized into two packages:

```
backend/internal/infrastructure/persistence/
├── postgres/              # Database utilities
│   ├── connection.go      # Connection pool management
│   ├── unit_of_work.go    # Transaction support
│   ├── query_builder.go   # SQL query helpers
│   └── errors.go          # Database error types
└── repository/            # Repository implementations
    └── user_repository.go # Implements domain.UserRepository
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
// internal/infrastructure/persistence/repository/user_repository.go
package repository

import (
    "context"
    "errors"
    "fmt"
    "time"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5"
    "github.com/yourorg/app/internal/domain/user"
    "github.com/yourorg/app/internal/infrastructure/persistence/postgres"
)

type UserRepository struct {
    pool postgres.ConnectionPool
}

func NewUserRepository(pool postgres.ConnectionPool) *UserRepository {
    return &UserRepository{pool: pool}
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
    query := `
        SELECT id, email, name, role, created_at, updated_at
        FROM users
        WHERE id = $1 AND deleted_at IS NULL
    `

    querier := postgres.GetQuerier(ctx, r.pool)

    var (
        dbID      uuid.UUID
        email     string
        name      string
        role      string
        createdAt time.Time
        updatedAt time.Time
    )

    err := querier.QueryRow(ctx, query, id).Scan(
        &dbID, &email, &name, &role, &createdAt, &updatedAt,
    )

    if errors.Is(err, pgx.ErrNoRows) {
        return nil, user.ErrUserNotFound
    }
    if err != nil {
        return nil, postgres.NewDBError("find user by id", err)
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

    querier := postgres.GetQuerier(ctx, r.pool)

    var (
        dbID      uuid.UUID
        dbEmail   string
        name      string
        role      string
        createdAt time.Time
        updatedAt time.Time
    )

    err := querier.QueryRow(ctx, query, email.String()).Scan(
        &dbID, &dbEmail, &name, &role, &createdAt, &updatedAt,
    )

    if errors.Is(err, pgx.ErrNoRows) {
        return nil, user.ErrUserNotFound
    }
    if err != nil {
        return nil, postgres.NewDBError("find user by email", err)
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

    querier := postgres.GetQuerier(ctx, r.pool)

    _, err := querier.Exec(ctx, query,
        u.ID(),
        u.Email().String(),
        u.Name(),
        string(u.Role()),
        u.CreatedAt(),
        u.UpdatedAt(),
    )

    if err != nil {
        dbErr := postgres.NewDBError("save user", err)
        if dbErr.IsUniqueViolation() {
            return user.ErrEmailAlreadyExists
        }
        return dbErr
    }

    return nil
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
    query := `
        UPDATE users
        SET deleted_at = $2, updated_at = $2
        WHERE id = $1 AND deleted_at IS NULL
    `

    querier := postgres.GetQuerier(ctx, r.pool)
    now := time.Now().UTC()

    result, err := querier.Exec(ctx, query, id, now)
    if err != nil {
        return postgres.NewDBError("delete user", err)
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

- [golang-migrate Documentation](https://github.com/golang-migrate/migrate)
- [golang-migrate CLI Usage](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
- [PostgreSQL Best Practices](https://wiki.postgresql.org/wiki/Don't_Do_This)
