# Database Migration Skill

Generate database migrations following golang-migrate patterns.

## Usage

```
/project:skill:migration <operation> <table-name>
```

## Migration Templates

### Create Table

**`migrations/XXXXXX_create_<table>_table.up.sql`**
```sql
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    email_verified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT users_email_check CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
    CONSTRAINT users_role_check CHECK (role IN ('admin', 'user', 'guest')),
    CONSTRAINT users_status_check CHECK (status IN ('active', 'inactive', 'suspended'))
);

CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_status ON users(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_role ON users(role) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_created_at ON users(created_at DESC);
```

**`migrations/XXXXXX_create_<table>_table.down.sql`**
```sql
DROP TABLE IF EXISTS users;
```

### Add Column

**`migrations/XXXXXX_add_<column>_to_<table>.up.sql`**
```sql
ALTER TABLE users ADD COLUMN phone VARCHAR(20);

CREATE INDEX idx_users_phone ON users(phone) WHERE phone IS NOT NULL AND deleted_at IS NULL;
```

**`migrations/XXXXXX_add_<column>_to_<table>.down.sql`**
```sql
DROP INDEX IF EXISTS idx_users_phone;
ALTER TABLE users DROP COLUMN IF EXISTS phone;
```

### Create Junction Table (Many-to-Many)

**`migrations/XXXXXX_create_user_roles_table.up.sql`**
```sql
CREATE TABLE IF NOT EXISTS user_roles (
    user_id UUID NOT NULL,
    role_id UUID NOT NULL,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    assigned_by UUID,

    PRIMARY KEY (user_id, role_id),
    CONSTRAINT fk_user_roles_user FOREIGN KEY (user_id)
        REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_roles_role FOREIGN KEY (role_id)
        REFERENCES roles(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_roles_assigned_by FOREIGN KEY (assigned_by)
        REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX idx_user_roles_user ON user_roles(user_id);
CREATE INDEX idx_user_roles_role ON user_roles(role_id);
```

**`migrations/XXXXXX_create_user_roles_table.down.sql`**
```sql
DROP TABLE IF EXISTS user_roles;
```

### Add Foreign Key

**`migrations/XXXXXX_add_<column>_fk_to_<table>.up.sql`**
```sql
ALTER TABLE orders ADD COLUMN user_id UUID;

ALTER TABLE orders
ADD CONSTRAINT fk_orders_user
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL;

CREATE INDEX idx_orders_user ON orders(user_id) WHERE user_id IS NOT NULL;
```

**`migrations/XXXXXX_add_<column>_fk_to_<table>.down.sql`**
```sql
ALTER TABLE orders DROP CONSTRAINT IF EXISTS fk_orders_user;
DROP INDEX IF EXISTS idx_orders_user;
ALTER TABLE orders DROP COLUMN IF EXISTS user_id;
```

### Create Enum Type

**`migrations/XXXXXX_create_<enum>_type.up.sql`**
```sql
DO $$ BEGIN
    CREATE TYPE order_status AS ENUM (
        'pending',
        'confirmed',
        'processing',
        'shipped',
        'delivered',
        'cancelled',
        'refunded'
    );
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

ALTER TABLE orders ADD COLUMN status order_status NOT NULL DEFAULT 'pending';

CREATE INDEX idx_orders_status ON orders(status);
```

**`migrations/XXXXXX_create_<enum>_type.down.sql`**
```sql
ALTER TABLE orders DROP COLUMN IF EXISTS status;
DROP TYPE IF EXISTS order_status;
```

### Add Trigger for Updated At

**`migrations/XXXXXX_add_updated_at_trigger.up.sql`**
```sql
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_orders_updated_at
    BEFORE UPDATE ON orders
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
```

**`migrations/XXXXXX_add_updated_at_trigger.down.sql`**
```sql
DROP TRIGGER IF EXISTS update_orders_updated_at ON orders;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP FUNCTION IF EXISTS update_updated_at_column();
```

### Create Full-Text Search Index

**`migrations/XXXXXX_add_search_index_to_<table>.up.sql`**
```sql
ALTER TABLE products ADD COLUMN search_vector tsvector
    GENERATED ALWAYS AS (
        setweight(to_tsvector('english', coalesce(name, '')), 'A') ||
        setweight(to_tsvector('english', coalesce(description, '')), 'B')
    ) STORED;

CREATE INDEX idx_products_search ON products USING GIN(search_vector);
```

**`migrations/XXXXXX_add_search_index_to_<table>.down.sql`**
```sql
DROP INDEX IF EXISTS idx_products_search;
ALTER TABLE products DROP COLUMN IF EXISTS search_vector;
```

### Add Soft Delete Support

**`migrations/XXXXXX_add_soft_delete_to_<table>.up.sql`**
```sql
ALTER TABLE products ADD COLUMN deleted_at TIMESTAMPTZ;
ALTER TABLE products ADD COLUMN deleted_by UUID REFERENCES users(id);

CREATE INDEX idx_products_deleted_at ON products(deleted_at) WHERE deleted_at IS NOT NULL;

DROP INDEX IF EXISTS idx_products_status;
CREATE INDEX idx_products_status ON products(status) WHERE deleted_at IS NULL;
```

**`migrations/XXXXXX_add_soft_delete_to_<table>.down.sql`**
```sql
DROP INDEX IF EXISTS idx_products_deleted_at;
DROP INDEX IF EXISTS idx_products_status;
ALTER TABLE products DROP COLUMN IF EXISTS deleted_by;
ALTER TABLE products DROP COLUMN IF EXISTS deleted_at;
CREATE INDEX idx_products_status ON products(status);
```

## Data Migration Example

**`migrations/XXXXXX_migrate_user_names.up.sql`**
```sql
ALTER TABLE users ADD COLUMN IF NOT EXISTS full_name VARCHAR(255);

UPDATE users
SET full_name = TRIM(CONCAT(first_name, ' ', last_name))
WHERE full_name IS NULL;

ALTER TABLE users ALTER COLUMN full_name SET NOT NULL;
```

**`migrations/XXXXXX_migrate_user_names.down.sql`**
```sql
ALTER TABLE users DROP COLUMN IF EXISTS full_name;
```

## Common Commands

```bash
# Create new migration
migrate create -ext sql -dir backend/migrations -seq <name>

# Apply all pending migrations
migrate -path backend/migrations -database "$DATABASE_URL" up

# Apply N migrations
migrate -path backend/migrations -database "$DATABASE_URL" up N

# Rollback last migration
migrate -path backend/migrations -database "$DATABASE_URL" down 1

# Rollback all migrations
migrate -path backend/migrations -database "$DATABASE_URL" down

# Check current version
migrate -path backend/migrations -database "$DATABASE_URL" version

# Force set version (fix dirty state)
migrate -path backend/migrations -database "$DATABASE_URL" force VERSION
```

## Best Practices

1. **Always write reversible migrations** - up and down
2. **Use IF EXISTS/IF NOT EXISTS** - prevent errors on re-runs
3. **Add indexes for foreign keys** - improve JOIN performance
4. **Use soft deletes** - add `deleted_at` column
5. **Add timestamps** - `created_at` and `updated_at`
6. **Use UUID primary keys** - better for distributed systems
7. **Add CHECK constraints** - enforce data integrity
8. **Create indexes for common queries** - WHERE/JOIN columns

## Checklist

- [ ] Migration file naming follows pattern: `XXXXXX_<action>_<table>_<details>`
- [ ] Both up.sql and down.sql created
- [ ] Down migration reverses up migration exactly
- [ ] Uses IF EXISTS/IF NOT EXISTS
- [ ] Appropriate indexes created
- [ ] Foreign key indexes added
- [ ] CHECK constraints for enums/validations
- [ ] Tested in development environment
