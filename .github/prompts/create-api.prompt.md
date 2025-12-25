---
description: Create a new REST API endpoint with handler, service, repository, and tests
agent: "Backend Engineer"
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

### 2. Domain Model

Create `backend/internal/domain/{{resourceName}}.go`:
- Define the main struct with JSON tags
- Define CreateInput struct with validation tags
- Define UpdateInput struct with pointer fields
- Add any domain-specific errors

### 3. Repository

Create `backend/internal/repository/postgres/{{resourceName}}_repository.go`:
- Implement FindByID
- Implement FindAll with pagination
- Implement Create
- Implement Update
- Implement Delete (soft delete)
- Handle errors properly (ErrNotFound, etc.)

### 4. Service

Create `backend/internal/service/{{resourceName}}_service.go`:
- Implement business logic
- Add validation
- Handle authorization if needed
- Add logging

### 5. Handler

Create `backend/internal/handlers/{{resourceName}}_handler.go`:
- Implement HTTP handlers for each operation
- Parse and validate request input
- Handle errors with appropriate status codes
- Format responses consistently
- Register routes

### 6. Tests

Create tests:
- `backend/internal/service/{{resourceName}}_service_test.go`
- `backend/internal/handlers/{{resourceName}}_handler_test.go`

## Code Templates

Use these patterns from the project:

```go
// Domain model template
type {{ResourceName}} struct {
    ID          string    `json:"id"`
    // Add fields
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// Handler error handling template
if errors.Is(err, domain.ErrNotFound) {
    response.NotFound(w, "{{ResourceName}} not found")
    return
}
```

## Validation

After implementation:
1. Run migrations: `migrate -path migrations -database "$DATABASE_URL" up`
2. Run tests: `go test ./...`
3. Run linter: `golangci-lint run`
4. Test manually with curl or API client

## Output

Provide:
1. List of files created
2. API endpoint documentation
3. Example curl commands for testing
4. Any additional configuration needed
