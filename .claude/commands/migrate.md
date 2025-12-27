# Database Migration Command

Create and manage database migrations using **golang-migrate** CLI.

## Task: $ARGUMENTS

## Quick Commands

```bash
# Install golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Create new migration (creates .up.sql and .down.sql files)
migrate create -ext sql -dir backend/migrations -seq <name>

# Apply all migrations
migrate -path backend/migrations -database "$DATABASE_URL" up

# Rollback last migration
migrate -path backend/migrations -database "$DATABASE_URL" down 1

# Rollback all migrations
migrate -path backend/migrations -database "$DATABASE_URL" down

# Check current version
migrate -path backend/migrations -database "$DATABASE_URL" version

# Force set version (for fixing dirty state)
migrate -path backend/migrations -database "$DATABASE_URL" force <version>
```

## Migration File Structure

```
backend/migrations/
├── 000001_create_users_table.up.sql
├── 000001_create_users_table.down.sql
├── 000002_add_users_indexes.up.sql
├── 000002_add_users_indexes.down.sql
└── ...
```

## Migration Naming Convention

```
XXXXXX_<action>_<table>_<details>.up.sql
XXXXXX_<action>_<table>_<details>.down.sql

Examples:
- 000001_create_users_table.up.sql
- 000002_add_users_email_index.up.sql
- 000003_add_users_role_column.up.sql
- 000004_create_products_table.up.sql
```

## Migration Templates

### Create Table

```sql
-- up migration
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Indexes for common queries
CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_status ON users(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_created_at ON users(created_at DESC);

-- down migration
DROP TABLE IF EXISTS users;
```

### Add Column

```sql
-- up migration
ALTER TABLE users ADD COLUMN phone VARCHAR(20);

-- down migration
ALTER TABLE users DROP COLUMN IF EXISTS phone;
```

### Add Index

```sql
-- up migration
CREATE INDEX idx_users_phone ON users(phone) WHERE phone IS NOT NULL AND deleted_at IS NULL;

-- down migration
DROP INDEX IF EXISTS idx_users_phone;
```

### Add Foreign Key

```sql
-- up migration
ALTER TABLE orders
ADD CONSTRAINT fk_orders_user_id
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- down migration
ALTER TABLE orders DROP CONSTRAINT IF EXISTS fk_orders_user_id;
```

### Create Junction Table (Many-to-Many)

```sql
-- up migration
CREATE TABLE IF NOT EXISTS user_roles (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, role_id)
);

CREATE INDEX idx_user_roles_user ON user_roles(user_id);
CREATE INDEX idx_user_roles_role ON user_roles(role_id);

-- down migration
DROP TABLE IF EXISTS user_roles;
```

### Add Enum Type

```sql
-- up migration
CREATE TYPE order_status AS ENUM ('pending', 'confirmed', 'shipped', 'delivered', 'cancelled');

ALTER TABLE orders ADD COLUMN status order_status NOT NULL DEFAULT 'pending';

-- down migration
ALTER TABLE orders DROP COLUMN IF EXISTS status;
DROP TYPE IF EXISTS order_status;
```

## Best Practices

### Always Do

1. **Make migrations reversible** - Always write both up and down migrations
2. **Use IF EXISTS/IF NOT EXISTS** - Prevent errors on re-runs
3. **Add indexes for foreign keys** - Improves JOIN performance
4. **Use soft deletes** - Add `deleted_at TIMESTAMPTZ` column
5. **Add timestamps** - Include `created_at` and `updated_at`
6. **Use UUID primary keys** - Better for distributed systems
7. **Add appropriate indexes** - Index columns used in WHERE/JOIN
8. **Use CHECK constraints** - Enforce data integrity at DB level

### Ask First

- Before modifying existing columns
- Before dropping columns or tables
- Before changing column types
- Before adding constraints to existing data

### Never Do

- Never modify an existing migration that has been applied
- Never delete migration files
- Never manually modify the schema_migrations table
- Never use `CASCADE` without understanding implications

## Handling Common Scenarios

### Fixing Dirty State

If a migration fails partway through:

```bash
# Check current state
migrate -path backend/migrations -database "$DATABASE_URL" version

# Force to known good version
migrate -path backend/migrations -database "$DATABASE_URL" force <version>

# Then manually fix the database if needed
```

### Data Migration

For migrations that require data transformation:

```sql
-- up migration
-- 1. Add new column
ALTER TABLE users ADD COLUMN full_name VARCHAR(255);

-- 2. Migrate data
UPDATE users SET full_name = CONCAT(first_name, ' ', last_name);

-- 3. Make non-nullable if needed
ALTER TABLE users ALTER COLUMN full_name SET NOT NULL;

-- 4. Drop old columns (optional, do in separate migration)
-- ALTER TABLE users DROP COLUMN first_name, DROP COLUMN last_name;
```

### Zero-Downtime Migrations

For production deployments:

1. **Add columns** - New columns should be nullable or have defaults
2. **Backfill data** - Run data migration separately
3. **Update code** - Deploy code that writes to both old and new
4. **Switch reads** - Update code to read from new column
5. **Remove old** - Drop old column in future migration
