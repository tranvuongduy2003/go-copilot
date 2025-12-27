# API Builder Command

Build REST API endpoints in Go following **DDD + CQRS** patterns with Clean Architecture.

## Task: $ARGUMENTS

## Architecture Overview

```
internal/
├── domain/<entity>/           # Domain Layer
│   ├── <entity>.go            # Entity with private fields + getters
│   ├── repository.go          # Repository interface (port)
│   └── errors.go              # Domain errors
├── application/
│   ├── command/               # Write operations (Create, Update, Delete)
│   ├── query/                 # Read operations (Get, List)
│   └── dto/                   # Data Transfer Objects
├── infrastructure/
│   └── persistence/
│       ├── postgres/          # Database utilities
│       └── repository/        # Repository implementations
└── interfaces/http/handler/   # HTTP handlers
```

## Step-by-Step Process

### Step 1: Create Migration (golang-migrate)

```sql
-- backend/migrations/000001_create_<entity>_table.up.sql
CREATE TABLE <entities> (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_<entities>_created_at ON <entities>(created_at DESC);
```

```sql
-- backend/migrations/000001_create_<entity>_table.down.sql
DROP TABLE IF EXISTS <entities>;
```

### Step 2: Domain Entity (DDD)

```go
// internal/domain/<entity>/<entity>.go
package <entity>

type <Entity> struct {
    id        uuid.UUID
    name      string
    createdAt time.Time
    updatedAt time.Time
}

func New<Entity>(name string) (*<Entity>, error) {
    if name == "" {
        return nil, ErrInvalidName
    }
    now := time.Now()
    return &<Entity>{
        id:        uuid.New(),
        name:      name,
        createdAt: now,
        updatedAt: now,
    }, nil
}

func (e *<Entity>) ID() uuid.UUID        { return e.id }
func (e *<Entity>) Name() string         { return e.name }
func (e *<Entity>) CreatedAt() time.Time { return e.createdAt }

func Reconstitute(id uuid.UUID, name string, createdAt, updatedAt time.Time) *<Entity> {
    return &<Entity>{id: id, name: name, createdAt: createdAt, updatedAt: updatedAt}
}
```

### Step 3: Repository Interface (Port)

```go
// internal/domain/<entity>/repository.go
type Repository interface {
    FindByID(ctx context.Context, id uuid.UUID) (*<Entity>, error)
    FindAll(ctx context.Context, opts ListOptions) ([]*<Entity>, int, error)
    Save(ctx context.Context, entity *<Entity>) error
    Delete(ctx context.Context, id uuid.UUID) error
}
```

### Step 4: Command Handlers (CQRS - Writes)

```go
// internal/application/command/create_<entity>.go
type Create<Entity>Command struct {
    Name string
}

type Create<Entity>Handler struct {
    repo <entity>.Repository
}

func (h *Create<Entity>Handler) Handle(ctx context.Context, cmd Create<Entity>Command) (*dto.<Entity>DTO, error) {
    e, err := <entity>.New<Entity>(cmd.Name)
    if err != nil {
        return nil, fmt.Errorf("invalid <entity>: %w", err)
    }
    if err := h.repo.Save(ctx, e); err != nil {
        return nil, fmt.Errorf("save <entity>: %w", err)
    }
    return dto.<Entity>FromDomain(e), nil
}
```

### Step 5: Query Handlers (CQRS - Reads)

```go
// internal/application/query/get_<entity>.go
type Get<Entity>Query struct {
    ID uuid.UUID
}

type Get<Entity>Handler struct {
    repo <entity>.Repository
}

func (h *Get<Entity>Handler) Handle(ctx context.Context, q Get<Entity>Query) (*dto.<Entity>DTO, error) {
    e, err := h.repo.FindByID(ctx, q.ID)
    if err != nil {
        return nil, err
    }
    return dto.<Entity>FromDomain(e), nil
}
```

### Step 6: DTO

```go
// internal/application/dto/<entity>_dto.go
type <Entity>DTO struct {
    ID        uuid.UUID `json:"id"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func <Entity>FromDomain(e *<entity>.<Entity>) *<Entity>DTO {
    return &<Entity>DTO{
        ID:        e.ID(),
        Name:      e.Name(),
        CreatedAt: e.CreatedAt(),
        UpdatedAt: e.UpdatedAt(),
    }
}
```

### Step 7: Repository Implementation

```go
// internal/infrastructure/persistence/repository/<entity>_repository.go
func (r *<entity>Repository) FindByID(ctx context.Context, id uuid.UUID) (*<entity>.<Entity>, error) {
    query := `SELECT id, name, created_at, updated_at FROM <entities> WHERE id = $1 AND deleted_at IS NULL`

    var dbID uuid.UUID
    var name string
    var createdAt, updatedAt time.Time

    err := r.pool.QueryRow(ctx, query, id).Scan(&dbID, &name, &createdAt, &updatedAt)
    if errors.Is(err, pgx.ErrNoRows) {
        return nil, <entity>.ErrNotFound
    }
    if err != nil {
        return nil, fmt.Errorf("query <entity>: %w", err)
    }

    return <entity>.Reconstitute(dbID, name, createdAt, updatedAt), nil
}
```

### Step 8: HTTP Handler

```go
// internal/interfaces/http/handler/<entity>_handler.go
func (h *<Entity>Handler) RegisterRoutes(r chi.Router) {
    r.Route("/<entities>", func(r chi.Router) {
        r.Get("/", h.List)
        r.Post("/", h.Create)
        r.Get("/{id}", h.Get)
        r.Put("/{id}", h.Update)
        r.Delete("/{id}", h.Delete)
    })
}

func (h *<Entity>Handler) Get(w http.ResponseWriter, r *http.Request) {
    id, err := uuid.Parse(chi.URLParam(r, "id"))
    if err != nil {
        response.BadRequest(w, "Invalid ID")
        return
    }

    result, err := h.getHandler.Handle(r.Context(), query.Get<Entity>Query{ID: id})
    if err != nil {
        if errors.Is(err, <entity>.ErrNotFound) {
            response.NotFound(w, "<Entity> not found")
            return
        }
        response.InternalError(w, err)
        return
    }

    response.JSON(w, http.StatusOK, result)
}
```

## API Response Format

```json
// Success Response
{
  "data": { ... }
}

// Success with Pagination
{
  "data": [...],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 100
  }
}

// Error Response
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Resource not found"
  }
}
```

## Checklist

- [ ] Migration files (.up.sql and .down.sql)
- [ ] Domain entity with private fields + getters
- [ ] Repository interface (port) in domain layer
- [ ] Command handlers (Create, Update, Delete)
- [ ] Query handlers (Get, List)
- [ ] DTOs for API responses
- [ ] Repository implementation using Reconstitute
- [ ] HTTP handler with CQRS handlers
- [ ] Request validation
- [ ] Routes registered
- [ ] Unit tests
