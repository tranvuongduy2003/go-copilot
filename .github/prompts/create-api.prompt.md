---
description: Create a new REST API endpoint with handler, service, repository, and tests
---

# Create REST API Endpoint

Create a new REST API endpoint following the project's clean architecture.

## Endpoint Details

**Resource Name**: {{resourceName}}

**Description**: {{description}}

**Operations Needed**:
- [ ] List (GET /)
- [ ] Get by ID (GET /:id)
- [ ] Create (POST /)
- [ ] Update (PUT /:id)
- [ ] Delete (DELETE /:id)

## Implementation Steps

### 1. Database Migration

Create migration files:
- `backend/migrations/XXX_create_{{resourceName}}.up.sql`
- `backend/migrations/XXX_create_{{resourceName}}.down.sql`

Include:
- Primary key (UUID)
- Required fields
- Timestamps (created_at, updated_at, deleted_at)
- Appropriate indexes
- Foreign key constraints if needed

### 2. Domain Layer

Create `backend/internal/domain/{{resourceName}}/`:
- `{{resourceName}}.go` - Entity with private fields + getter methods
- `repository.go` - Repository interface (port)
- `errors.go` - Domain-specific errors
- Define value objects if needed

### 3. Application Layer (CQRS)

Create command handlers in `backend/internal/application/command/`:
- `create_{{resourceName}}.go` - Create command handler
- `update_{{resourceName}}.go` - Update command handler
- `delete_{{resourceName}}.go` - Delete command handler

Create query handlers in `backend/internal/application/query/`:
- `get_{{resourceName}}.go` - Get by ID query handler
- `list_{{resourceName}}s.go` - List query handler with pagination

### 4. Infrastructure Layer

Create `backend/internal/infrastructure/persistence/postgres/{{resourceName}}_repository.go`:
- Implement the domain repository interface
- Use parameterized queries
- Handle errors properly (ErrNotFound, etc.)

### 5. Handler

Create `backend/internal/interfaces/http/handler/{{resourceName}}_handler.go`:
- Implement HTTP handlers for each operation
- Parse and validate request input
- Handle errors with appropriate status codes
- Format responses consistently
- Register routes

### 6. Tests

Create tests:
- `backend/internal/application/command/create_{{resourceName}}_test.go`
- `backend/internal/application/query/get_{{resourceName}}_test.go`
- `backend/internal/interfaces/http/handler/{{resourceName}}_handler_test.go`

## Code Templates

Use these patterns from the project:

```go
// Domain entity template (private fields + getters)
// internal/domain/{{resourceName}}/{{resourceName}}.go
package {{resourceName}}

type {{ResourceName}} struct {
    id        uuid.UUID
    name      string
    createdAt time.Time
    updatedAt time.Time
}

func (e *{{ResourceName}}) ID() uuid.UUID      { return e.id }
func (e *{{ResourceName}}) Name() string       { return e.name }
func (e *{{ResourceName}}) CreatedAt() time.Time { return e.createdAt }

// Repository interface (port)
// internal/domain/{{resourceName}}/repository.go
type Repository interface {
    FindByID(ctx context.Context, id uuid.UUID) (*{{ResourceName}}, error)
    Save(ctx context.Context, entity *{{ResourceName}}) error
}

// Handler error handling template
if errors.Is(err, domain.ErrNotFound) {
    response.NotFound(w, "{{ResourceName}} not found")
    return
}
```

## Validation

After implementation:
1. Run migrations: `goose -dir backend/migrations/sql postgres "$DATABASE_URL" up`
2. Run tests: `cd backend && go test ./...`
3. Run linter: `cd backend && golangci-lint run`
4. Test manually with curl or API client

## Output

Provide:
1. List of files created
2. API endpoint documentation
3. Example curl commands for testing
4. Any additional configuration needed
