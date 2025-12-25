---
applyTo: "backend/internal/repository/**/*.go,backend/migrations/**/*.sql"
---

# Database Development Instructions

These instructions apply to database migrations and repository implementations.

## Migration Guidelines

### File Naming

```
migrations/
├── 001_create_users.up.sql
├── 001_create_users.down.sql
├── 002_add_user_roles.up.sql
├── 002_add_user_roles.down.sql
├── 003_create_posts.up.sql
└── 003_create_posts.down.sql
```

### Migration Template

```sql
-- migrations/001_create_users.up.sql

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'user',
    email_verified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT users_email_unique UNIQUE (email)
);

-- Create indexes
CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_role ON users(role) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_created_at ON users(created_at);

-- Add comment
COMMENT ON TABLE users IS 'Application users';
COMMENT ON COLUMN users.deleted_at IS 'Soft delete timestamp';
```

```sql
-- migrations/001_create_users.down.sql

DROP TABLE IF EXISTS users;
```

### Migration Best Practices

1. **Always provide down migrations**
   ```sql
   -- up: Add column
   ALTER TABLE users ADD COLUMN avatar_url TEXT;

   -- down: Remove column
   ALTER TABLE users DROP COLUMN avatar_url;
   ```

2. **Use transactions for safety**
   ```sql
   BEGIN;

   ALTER TABLE users ADD COLUMN status VARCHAR(50);
   UPDATE users SET status = 'active' WHERE deleted_at IS NULL;
   UPDATE users SET status = 'deleted' WHERE deleted_at IS NOT NULL;
   ALTER TABLE users ALTER COLUMN status SET NOT NULL;

   COMMIT;
   ```

3. **Add indexes for foreign keys and commonly queried columns**
   ```sql
   CREATE INDEX idx_posts_user_id ON posts(user_id);
   CREATE INDEX idx_posts_status ON posts(status) WHERE deleted_at IS NULL;
   CREATE INDEX idx_posts_created_at ON posts(created_at DESC);
   ```

4. **Use partial indexes for soft deletes**
   ```sql
   -- Only index non-deleted records
   CREATE UNIQUE INDEX idx_users_email_active
   ON users(email)
   WHERE deleted_at IS NULL;
   ```

5. **Consider adding constraints**
   ```sql
   ALTER TABLE posts
   ADD CONSTRAINT posts_status_check
   CHECK (status IN ('draft', 'published', 'archived'));
   ```

## Repository Pattern

### Repository Interface

```go
package repository

import (
    "context"

    "github.com/yourorg/app/internal/domain"
)

// ListOptions contains pagination and filtering options.
type ListOptions struct {
    Page    int
    PerPage int
    Sort    string
    Order   string // "asc" or "desc"
    Search  string
    Status  string
}

// UserRepository defines the interface for user data operations.
type UserRepository interface {
    // FindByID retrieves a user by ID.
    // Returns ErrNotFound if user doesn't exist.
    FindByID(ctx context.Context, id string) (*domain.User, error)

    // FindByEmail retrieves a user by email (case-insensitive).
    // Returns ErrNotFound if user doesn't exist.
    FindByEmail(ctx context.Context, email string) (*domain.User, error)

    // FindAll retrieves users with pagination and filtering.
    // Returns the users, total count, and any error.
    FindAll(ctx context.Context, opts ListOptions) ([]*domain.User, int, error)

    // Create persists a new user.
    // Returns ErrConflict if email already exists.
    Create(ctx context.Context, user *domain.User) error

    // Update persists changes to an existing user.
    // Returns ErrNotFound if user doesn't exist.
    Update(ctx context.Context, user *domain.User) error

    // Delete soft-deletes a user.
    // Returns ErrNotFound if user doesn't exist.
    Delete(ctx context.Context, id string) error
}
```

### PostgreSQL Implementation

```go
package postgres

import (
    "context"
    "errors"
    "fmt"
    "strings"
    "time"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/yourorg/app/internal/domain"
    "github.com/yourorg/app/internal/repository"
)

type userRepository struct {
    db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) repository.UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
    query := `
        SELECT id, email, name, role, created_at, updated_at
        FROM users
        WHERE id = $1 AND deleted_at IS NULL
    `

    var user domain.User
    err := r.db.QueryRow(ctx, query, id).Scan(
        &user.ID,
        &user.Email,
        &user.Name,
        &user.Role,
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

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
    query := `
        SELECT id, email, name, password_hash, role, created_at, updated_at
        FROM users
        WHERE LOWER(email) = LOWER($1) AND deleted_at IS NULL
    `

    var user domain.User
    err := r.db.QueryRow(ctx, query, email).Scan(
        &user.ID,
        &user.Email,
        &user.Name,
        &user.PasswordHash,
        &user.Role,
        &user.CreatedAt,
        &user.UpdatedAt,
    )

    if errors.Is(err, pgx.ErrNoRows) {
        return nil, domain.ErrNotFound
    }
    if err != nil {
        return nil, fmt.Errorf("query user by email: %w", err)
    }

    return &user, nil
}

func (r *userRepository) FindAll(ctx context.Context, opts repository.ListOptions) ([]*domain.User, int, error) {
    // Build query with filters
    baseQuery := `FROM users WHERE deleted_at IS NULL`
    args := []interface{}{}
    argIndex := 1

    if opts.Search != "" {
        baseQuery += fmt.Sprintf(` AND (name ILIKE $%d OR email ILIKE $%d)`, argIndex, argIndex)
        args = append(args, "%"+opts.Search+"%")
        argIndex++
    }

    if opts.Status != "" {
        baseQuery += fmt.Sprintf(` AND role = $%d`, argIndex)
        args = append(args, opts.Status)
        argIndex++
    }

    // Count total
    var total int
    countQuery := `SELECT COUNT(*) ` + baseQuery
    if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
        return nil, 0, fmt.Errorf("count users: %w", err)
    }

    // Build select with pagination
    sortColumn := "created_at"
    if opts.Sort != "" && isValidSortColumn(opts.Sort) {
        sortColumn = opts.Sort
    }

    sortOrder := "DESC"
    if opts.Order == "asc" {
        sortOrder = "ASC"
    }

    offset := (opts.Page - 1) * opts.PerPage

    selectQuery := fmt.Sprintf(`
        SELECT id, email, name, role, created_at, updated_at
        %s
        ORDER BY %s %s
        LIMIT $%d OFFSET $%d
    `, baseQuery, sortColumn, sortOrder, argIndex, argIndex+1)

    args = append(args, opts.PerPage, offset)

    rows, err := r.db.Query(ctx, selectQuery, args...)
    if err != nil {
        return nil, 0, fmt.Errorf("query users: %w", err)
    }
    defer rows.Close()

    var users []*domain.User
    for rows.Next() {
        var user domain.User
        if err := rows.Scan(
            &user.ID,
            &user.Email,
            &user.Name,
            &user.Role,
            &user.CreatedAt,
            &user.UpdatedAt,
        ); err != nil {
            return nil, 0, fmt.Errorf("scan user: %w", err)
        }
        users = append(users, &user)
    }

    return users, total, nil
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
    query := `
        INSERT INTO users (id, email, name, password_hash, role, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `

    _, err := r.db.Exec(ctx, query,
        user.ID,
        user.Email,
        user.Name,
        user.PasswordHash,
        user.Role,
        user.CreatedAt,
        user.UpdatedAt,
    )

    if err != nil {
        if isUniqueViolation(err) {
            return domain.ErrConflict
        }
        return fmt.Errorf("insert user: %w", err)
    }

    return nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
    query := `
        UPDATE users
        SET name = $2, role = $3, updated_at = $4
        WHERE id = $1 AND deleted_at IS NULL
    `

    result, err := r.db.Exec(ctx, query,
        user.ID,
        user.Name,
        user.Role,
        time.Now().UTC(),
    )

    if err != nil {
        return fmt.Errorf("update user: %w", err)
    }

    if result.RowsAffected() == 0 {
        return domain.ErrNotFound
    }

    return nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
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
        return domain.ErrNotFound
    }

    return nil
}

// Helper functions

func isUniqueViolation(err error) bool {
    var pgErr *pgconn.PgError
    if errors.As(err, &pgErr) {
        return pgErr.Code == "23505" // unique_violation
    }
    return false
}

func isValidSortColumn(col string) bool {
    validColumns := map[string]bool{
        "created_at": true,
        "updated_at": true,
        "name":       true,
        "email":      true,
    }
    return validColumns[col]
}
```

## Transaction Handling

```go
func (r *Repository) TransferFunds(ctx context.Context, fromID, toID string, amount int) error {
    tx, err := r.db.Begin(ctx)
    if err != nil {
        return fmt.Errorf("begin transaction: %w", err)
    }
    defer tx.Rollback(ctx)

    // Check source balance
    var balance int
    err = tx.QueryRow(ctx,
        "SELECT balance FROM accounts WHERE id = $1 FOR UPDATE",
        fromID,
    ).Scan(&balance)
    if err != nil {
        return fmt.Errorf("get source balance: %w", err)
    }

    if balance < amount {
        return domain.ErrInsufficientFunds
    }

    // Debit source
    _, err = tx.Exec(ctx,
        "UPDATE accounts SET balance = balance - $1 WHERE id = $2",
        amount, fromID,
    )
    if err != nil {
        return fmt.Errorf("debit account: %w", err)
    }

    // Credit destination
    _, err = tx.Exec(ctx,
        "UPDATE accounts SET balance = balance + $1 WHERE id = $2",
        amount, toID,
    )
    if err != nil {
        return fmt.Errorf("credit account: %w", err)
    }

    if err := tx.Commit(ctx); err != nil {
        return fmt.Errorf("commit transaction: %w", err)
    }

    return nil
}
```

## Query Optimization

### Use Prepared Statements for Repeated Queries

```go
type userRepository struct {
    db              *pgxpool.Pool
    findByIDStmt    string
    findByEmailStmt string
}

func NewUserRepository(db *pgxpool.Pool) *userRepository {
    return &userRepository{
        db: db,
        findByIDStmt: `
            SELECT id, email, name, role, created_at, updated_at
            FROM users
            WHERE id = $1 AND deleted_at IS NULL
        `,
        findByEmailStmt: `
            SELECT id, email, name, password_hash, role, created_at, updated_at
            FROM users
            WHERE LOWER(email) = LOWER($1) AND deleted_at IS NULL
        `,
    }
}
```

### Batch Operations

```go
func (r *Repository) CreateMany(ctx context.Context, users []*domain.User) error {
    batch := &pgx.Batch{}

    for _, user := range users {
        batch.Queue(
            `INSERT INTO users (id, email, name, created_at, updated_at)
             VALUES ($1, $2, $3, $4, $5)`,
            user.ID, user.Email, user.Name, user.CreatedAt, user.UpdatedAt,
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

## Security Considerations

1. **Always use parameterized queries**
   ```go
   // CORRECT
   query := "SELECT * FROM users WHERE email = $1"
   row := db.QueryRow(ctx, query, email)

   // NEVER DO THIS
   query := "SELECT * FROM users WHERE email = '" + email + "'"
   ```

2. **Validate input before queries**
   ```go
   func isValidSortColumn(col string) bool {
       validColumns := map[string]bool{"created_at": true, "name": true}
       return validColumns[col]
   }
   ```

3. **Use row-level security when appropriate**
   ```sql
   ALTER TABLE documents ENABLE ROW LEVEL SECURITY;

   CREATE POLICY documents_user_policy ON documents
       FOR ALL
       USING (user_id = current_user_id());
   ```
