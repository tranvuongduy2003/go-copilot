---
name: database-migration
description: Create database migrations safely. Use when modifying database schema.
---

# Database Migration Skill

This skill guides you through creating safe, reversible database migrations for PostgreSQL.

## When to Use This Skill

- Creating new tables
- Adding or modifying columns
- Creating indexes
- Adding constraints
- Modifying relationships

## Migration File Naming

```
migrations/
├── 001_create_users.up.sql
├── 001_create_users.down.sql
├── 002_add_user_roles.up.sql
├── 002_add_user_roles.down.sql
```

Format: `{sequence}_{description}.{up|down}.sql`

## Migration Templates

### Template 1: Create Table

```sql
-- migrations/001_create_users.up.sql

-- Users table stores application user accounts
CREATE TABLE IF NOT EXISTS users (
    -- Primary key
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Core fields
    email VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,

    -- Status and role
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    email_verified_at TIMESTAMPTZ,

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    -- Constraints
    CONSTRAINT users_email_unique UNIQUE (email),
    CONSTRAINT users_role_check CHECK (role IN ('user', 'admin', 'moderator')),
    CONSTRAINT users_status_check CHECK (status IN ('active', 'inactive', 'suspended'))
);

-- Indexes for common queries
CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_role ON users(role) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_status ON users(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_created_at ON users(created_at DESC);

-- Table and column comments
COMMENT ON TABLE users IS 'Application user accounts';
COMMENT ON COLUMN users.password_hash IS 'bcrypt hashed password';
COMMENT ON COLUMN users.deleted_at IS 'Soft delete timestamp - NULL means active';
```

```sql
-- migrations/001_create_users.down.sql

DROP TABLE IF EXISTS users;
```

### Template 2: Create Table with Foreign Key

```sql
-- migrations/002_create_posts.up.sql

CREATE TABLE IF NOT EXISTS posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Foreign key to users
    user_id UUID NOT NULL,

    -- Content
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    content TEXT,
    excerpt VARCHAR(500),

    -- Metadata
    status VARCHAR(50) NOT NULL DEFAULT 'draft',
    published_at TIMESTAMPTZ,

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    -- Foreign key constraint
    CONSTRAINT posts_user_id_fk
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    -- Check constraints
    CONSTRAINT posts_status_check
        CHECK (status IN ('draft', 'published', 'archived')),

    -- Unique constraint (partial - only non-deleted)
    CONSTRAINT posts_slug_unique
        UNIQUE (slug)
);

-- Indexes
CREATE INDEX idx_posts_user_id ON posts(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_posts_status ON posts(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_posts_slug ON posts(slug) WHERE deleted_at IS NULL;
CREATE INDEX idx_posts_published_at ON posts(published_at DESC) WHERE status = 'published';
CREATE INDEX idx_posts_created_at ON posts(created_at DESC);

-- Full text search index
CREATE INDEX idx_posts_search ON posts USING gin(to_tsvector('english', title || ' ' || COALESCE(content, '')));
```

```sql
-- migrations/002_create_posts.down.sql

DROP TABLE IF EXISTS posts;
```

### Template 3: Add Column

```sql
-- migrations/003_add_user_avatar.up.sql

-- Add avatar URL column to users
ALTER TABLE users
ADD COLUMN avatar_url TEXT;

-- Add profile columns
ALTER TABLE users
ADD COLUMN bio VARCHAR(500),
ADD COLUMN website VARCHAR(255),
ADD COLUMN location VARCHAR(100);

-- Comment on new columns
COMMENT ON COLUMN users.avatar_url IS 'URL to user avatar image';
COMMENT ON COLUMN users.bio IS 'User biography/description';
```

```sql
-- migrations/003_add_user_avatar.down.sql

ALTER TABLE users
DROP COLUMN IF EXISTS avatar_url,
DROP COLUMN IF EXISTS bio,
DROP COLUMN IF EXISTS website,
DROP COLUMN IF EXISTS location;
```

### Template 4: Add Column with Default (Safe for Large Tables)

```sql
-- migrations/004_add_user_settings.up.sql

-- For large tables, add column without default first
ALTER TABLE users
ADD COLUMN settings JSONB;

-- Then update existing rows in batches (do this in application code for very large tables)
UPDATE users SET settings = '{}' WHERE settings IS NULL;

-- Finally add the default for new rows
ALTER TABLE users
ALTER COLUMN settings SET DEFAULT '{}',
ALTER COLUMN settings SET NOT NULL;
```

```sql
-- migrations/004_add_user_settings.down.sql

ALTER TABLE users
DROP COLUMN IF EXISTS settings;
```

### Template 5: Create Junction Table (Many-to-Many)

```sql
-- migrations/005_create_post_tags.up.sql

-- Tags table
CREATE TABLE IF NOT EXISTS tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT tags_name_unique UNIQUE (name),
    CONSTRAINT tags_slug_unique UNIQUE (slug)
);

-- Junction table for posts and tags
CREATE TABLE IF NOT EXISTS post_tags (
    post_id UUID NOT NULL,
    tag_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (post_id, tag_id),

    CONSTRAINT post_tags_post_fk
        FOREIGN KEY (post_id)
        REFERENCES posts(id)
        ON DELETE CASCADE,

    CONSTRAINT post_tags_tag_fk
        FOREIGN KEY (tag_id)
        REFERENCES tags(id)
        ON DELETE CASCADE
);

-- Indexes for junction table
CREATE INDEX idx_post_tags_post_id ON post_tags(post_id);
CREATE INDEX idx_post_tags_tag_id ON post_tags(tag_id);
```

```sql
-- migrations/005_create_post_tags.down.sql

DROP TABLE IF EXISTS post_tags;
DROP TABLE IF EXISTS tags;
```

### Template 6: Add Index

```sql
-- migrations/006_add_performance_indexes.up.sql

-- Composite index for common query pattern
CREATE INDEX CONCURRENTLY idx_posts_user_status
ON posts(user_id, status)
WHERE deleted_at IS NULL;

-- Partial index for active users
CREATE INDEX CONCURRENTLY idx_users_active_email
ON users(email)
WHERE status = 'active' AND deleted_at IS NULL;

-- Index for sorting
CREATE INDEX CONCURRENTLY idx_posts_hot
ON posts(published_at DESC, id)
WHERE status = 'published';
```

```sql
-- migrations/006_add_performance_indexes.down.sql

DROP INDEX CONCURRENTLY IF EXISTS idx_posts_user_status;
DROP INDEX CONCURRENTLY IF EXISTS idx_users_active_email;
DROP INDEX CONCURRENTLY IF EXISTS idx_posts_hot;
```

### Template 7: Modify Column (Safe Rename)

```sql
-- migrations/007_rename_column.up.sql

-- Rename column (PostgreSQL 9.4+)
ALTER TABLE users
RENAME COLUMN name TO full_name;

-- Update any related indexes or constraints if needed
```

```sql
-- migrations/007_rename_column.down.sql

ALTER TABLE users
RENAME COLUMN full_name TO name;
```

### Template 8: Add Enum Type

```sql
-- migrations/008_add_priority_enum.up.sql

-- Create enum type
CREATE TYPE priority_level AS ENUM ('low', 'medium', 'high', 'urgent');

-- Add column using enum
ALTER TABLE posts
ADD COLUMN priority priority_level DEFAULT 'medium';

CREATE INDEX idx_posts_priority ON posts(priority) WHERE deleted_at IS NULL;
```

```sql
-- migrations/008_add_priority_enum.down.sql

ALTER TABLE posts
DROP COLUMN IF EXISTS priority;

DROP TYPE IF EXISTS priority_level;
```

## Migration Commands

```bash
# Create new migration
migrate create -ext sql -dir migrations -seq create_orders

# Run all pending migrations
migrate -path migrations -database "$DATABASE_URL" up

# Run specific number of migrations
migrate -path migrations -database "$DATABASE_URL" up 1

# Rollback last migration
migrate -path migrations -database "$DATABASE_URL" down 1

# Rollback all migrations
migrate -path migrations -database "$DATABASE_URL" down

# Check current version
migrate -path migrations -database "$DATABASE_URL" version

# Force set version (use carefully!)
migrate -path migrations -database "$DATABASE_URL" force 5
```

## Best Practices

### 1. Always Provide Reversible Migrations
```sql
-- Every .up.sql needs a corresponding .down.sql
```

### 2. Use Transactions
```sql
BEGIN;
-- Your migration SQL here
COMMIT;
```

### 3. Use CONCURRENTLY for Indexes on Large Tables
```sql
-- Won't lock the table
CREATE INDEX CONCURRENTLY idx_name ON table(column);
```

### 4. Add NOT NULL Safely
```sql
-- Step 1: Add column nullable
ALTER TABLE users ADD COLUMN new_col TEXT;

-- Step 2: Backfill data
UPDATE users SET new_col = 'default' WHERE new_col IS NULL;

-- Step 3: Add NOT NULL constraint
ALTER TABLE users ALTER COLUMN new_col SET NOT NULL;
```

### 5. Use Partial Indexes for Soft Deletes
```sql
CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
```

### 6. Document Your Schema
```sql
COMMENT ON TABLE users IS 'Application user accounts';
COMMENT ON COLUMN users.deleted_at IS 'Soft delete timestamp';
```

## Migration Checklist

- [ ] Migration file follows naming convention
- [ ] Up migration creates/modifies schema
- [ ] Down migration reverses changes completely
- [ ] Indexes added for foreign keys
- [ ] Indexes added for commonly queried columns
- [ ] Constraints have meaningful names
- [ ] Comments added for complex logic
- [ ] Tested locally with up and down
- [ ] Considered performance for large tables
