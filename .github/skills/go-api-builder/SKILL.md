---
name: go-api-builder
description: Build REST API endpoints in Go following DDD + CQRS patterns. Use when creating new API endpoints with domain entities, command/query handlers, and HTTP handlers.
---

# Go API Builder Skill

This skill guides you through building complete REST API endpoints in Go following **DDD + CQRS** patterns with Clean Architecture.

## When to Use This Skill

- Creating a new API endpoint (GET, POST, PUT, DELETE)
- Adding a new resource/aggregate to the API
- Implementing CRUD operations using command/query handlers

## Architecture Overview

```
internal/
├── domain/product/           # Domain Layer (entities, repository interface)
├── application/
│   ├── command/              # Write operations (Create, Update, Delete)
│   ├── query/                # Read operations (Get, List)
│   └── dto/                  # Data Transfer Objects
├── infrastructure/
│   └── persistence/
│       ├── postgres/         # Database utilities (connection, unit_of_work, errors)
│       └── repository/       # Repository implementations
└── interfaces/http/handler/  # HTTP handlers
```

## Step-by-Step Process

### Step 1: Define the Domain Entity (DDD Pattern)

Create the domain entity with **private fields and getter methods** in `internal/domain/product/`.

```go
// internal/domain/product/product.go
package product

import (
    "errors"
    "time"

    "github.com/google/uuid"
)

// Domain errors
var (
    ErrNotFound     = errors.New("product not found")
    ErrInvalidPrice = errors.New("price must be positive")
    ErrInvalidName  = errors.New("name cannot be empty")
)

// Product is the aggregate root for product operations.
// Uses private fields with getter methods (DDD pattern).
type Product struct {
    id          uuid.UUID
    name        string
    description string
    price       int // Price in cents
    category    string
    stock       int
    createdAt   time.Time
    updatedAt   time.Time
}

// NewProduct creates a new Product with validation.
func NewProduct(name, description, category string, price, stock int) (*Product, error) {
    if name == "" {
        return nil, ErrInvalidName
    }
    if price <= 0 {
        return nil, ErrInvalidPrice
    }
    now := time.Now()
    return &Product{
        id:          uuid.New(),
        name:        name,
        description: description,
        price:       price,
        category:    category,
        stock:       stock,
        createdAt:   now,
        updatedAt:   now,
    }, nil
}

// Getter methods
func (p *Product) ID() uuid.UUID        { return p.id }
func (p *Product) Name() string         { return p.name }
func (p *Product) Description() string  { return p.description }
func (p *Product) Price() int           { return p.price }
func (p *Product) Category() string     { return p.category }
func (p *Product) Stock() int           { return p.stock }
func (p *Product) CreatedAt() time.Time { return p.createdAt }
func (p *Product) UpdatedAt() time.Time { return p.updatedAt }

// UpdatePrice updates the product price with validation.
func (p *Product) UpdatePrice(price int) error {
    if price <= 0 {
        return ErrInvalidPrice
    }
    p.price = price
    p.updatedAt = time.Now()
    return nil
}

// Reconstitute recreates a Product from persistence (used by repository).
func Reconstitute(id uuid.UUID, name, description, category string, price, stock int, createdAt, updatedAt time.Time) *Product {
    return &Product{
        id:          id,
        name:        name,
        description: description,
        price:       price,
        category:    category,
        stock:       stock,
        createdAt:   createdAt,
        updatedAt:   updatedAt,
    }
}
```

### Step 2: Create the Database Migration (golang-migrate)

Create migration files in `backend/migrations/`.

**Up Migration:**
```sql
-- backend/migrations/000001_create_products_table.up.sql

CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price INTEGER NOT NULL CHECK (price > 0),
    category VARCHAR(100) NOT NULL,
    stock INTEGER NOT NULL DEFAULT 0 CHECK (stock >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_products_category ON products(category) WHERE deleted_at IS NULL;
CREATE INDEX idx_products_created_at ON products(created_at DESC);
```

**Down Migration:**
```sql
-- backend/migrations/000001_create_products_table.down.sql

DROP TABLE IF EXISTS products;
```

Run migration: `migrate -path backend/migrations -database "$DATABASE_URL" up`

### Step 3: Define Repository Interface (Domain Layer)

Create the repository interface in the domain layer (port). The interface works with domain entities.

```go
// internal/domain/product/repository.go
package product

import (
    "context"

    "github.com/google/uuid"
)

// Repository defines the port for product persistence (DDD pattern).
type Repository interface {
    FindByID(ctx context.Context, id uuid.UUID) (*Product, error)
    FindAll(ctx context.Context, opts ListOptions) ([]*Product, int, error)
    Save(ctx context.Context, product *Product) error
    Delete(ctx context.Context, id uuid.UUID) error
}

// ListOptions defines pagination and filtering options.
type ListOptions struct {
    Page     int
    PerPage  int
    Category string
}
```

### Step 4: Implement Repository (Infrastructure Layer)

The repository adapter uses `Reconstitute` to create domain entities from database rows and getter methods to extract values for persistence.

```go
// internal/infrastructure/persistence/repository/product_repository.go
package repository

import (
    "context"
    "errors"
    "time"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5"
    "github.com/yourorg/app/internal/domain/product"
    "github.com/yourorg/app/internal/infrastructure/persistence/postgres"
)

type productRepository struct {
    pool postgres.ConnectionPool
}

// NewProductRepository creates a new PostgreSQL product repository (adapter).
func NewProductRepository(pool postgres.ConnectionPool) product.Repository {
    return &productRepository{pool: pool}
}

func (r *productRepository) FindByID(ctx context.Context, id uuid.UUID) (*product.Product, error) {
    query := `
        SELECT id, name, description, price, category, stock, created_at, updated_at
        FROM products
        WHERE id = $1 AND deleted_at IS NULL
    `

    querier := postgres.GetQuerier(ctx, r.pool)

    var (
        dbID                      uuid.UUID
        name, description, category string
        price, stock              int
        createdAt, updatedAt      time.Time
    )

    err := querier.QueryRow(ctx, query, id).Scan(
        &dbID, &name, &description, &price, &category, &stock, &createdAt, &updatedAt,
    )
    if errors.Is(err, pgx.ErrNoRows) {
        return nil, product.ErrNotFound
    }
    if err != nil {
        return nil, postgres.NewDBError("query product", err)
    }

    // Use Reconstitute to create domain entity from DB values
    return product.Reconstitute(dbID, name, description, category, price, stock, createdAt, updatedAt), nil
}

func (r *productRepository) FindAll(ctx context.Context, opts product.ListOptions) ([]*product.Product, int, error) {
    querier := postgres.GetQuerier(ctx, r.pool)

    // Count total
    var total int
    countQuery := `SELECT COUNT(*) FROM products WHERE deleted_at IS NULL`
    if err := querier.QueryRow(ctx, countQuery).Scan(&total); err != nil {
        return nil, 0, postgres.NewDBError("count products", err)
    }

    // Fetch products
    offset := (opts.Page - 1) * opts.PerPage
    query := `
        SELECT id, name, description, price, category, stock, created_at, updated_at
        FROM products
        WHERE deleted_at IS NULL
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `

    rows, err := querier.Query(ctx, query, opts.PerPage, offset)
    if err != nil {
        return nil, 0, postgres.NewDBError("query products", err)
    }
    defer rows.Close()

    var products []*product.Product
    for rows.Next() {
        var (
            id                        uuid.UUID
            name, description, category string
            price, stock              int
            createdAt, updatedAt      time.Time
        )
        if err := rows.Scan(&id, &name, &description, &price, &category, &stock, &createdAt, &updatedAt); err != nil {
            return nil, 0, postgres.NewDBError("scan product", err)
        }
        products = append(products, product.Reconstitute(id, name, description, category, price, stock, createdAt, updatedAt))
    }

    return products, total, nil
}

func (r *productRepository) Save(ctx context.Context, p *product.Product) error {
    query := `
        INSERT INTO products (id, name, description, price, category, stock, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        ON CONFLICT (id) DO UPDATE SET
            name = EXCLUDED.name,
            description = EXCLUDED.description,
            price = EXCLUDED.price,
            category = EXCLUDED.category,
            stock = EXCLUDED.stock,
            updated_at = EXCLUDED.updated_at
    `

    querier := postgres.GetQuerier(ctx, r.pool)

    // Use getter methods to extract values from domain entity
    _, err := querier.Exec(ctx, query,
        p.ID(), p.Name(), p.Description(), p.Price(),
        p.Category(), p.Stock(), p.CreatedAt(), p.UpdatedAt(),
    )
    if err != nil {
        return postgres.NewDBError("save product", err)
    }

    return nil
}

func (r *productRepository) Delete(ctx context.Context, id uuid.UUID) error {
    query := `UPDATE products SET deleted_at = $2 WHERE id = $1 AND deleted_at IS NULL`

    querier := postgres.GetQuerier(ctx, r.pool)

    result, err := querier.Exec(ctx, query, id, time.Now().UTC())
    if err != nil {
        return postgres.NewDBError("delete product", err)
    }

    if result.RowsAffected() == 0 {
        return product.ErrNotFound
    }

    return nil
}
```

### Step 5: Implement Command/Query Handlers (Application Layer - CQRS)

Create separate **command handlers** (writes) and **query handlers** (reads).

#### Command: Create Product

```go
// internal/application/command/create_product.go
package command

import (
    "context"
    "fmt"
    "log/slog"

    "github.com/yourorg/app/internal/application/dto"
    "github.com/yourorg/app/internal/domain/product"
)

// CreateProductCommand contains the data needed to create a product.
type CreateProductCommand struct {
    Name        string
    Description string
    Price       int
    Category    string
    Stock       int
}

// CreateProductHandler handles the CreateProductCommand.
type CreateProductHandler struct {
    repo   product.Repository
    logger *slog.Logger
}

func NewCreateProductHandler(repo product.Repository, logger *slog.Logger) *CreateProductHandler {
    return &CreateProductHandler{repo: repo, logger: logger}
}

// Handle executes the create product command.
func (h *CreateProductHandler) Handle(ctx context.Context, cmd CreateProductCommand) (*dto.ProductDTO, error) {
    // Create domain entity with validation
    p, err := product.NewProduct(cmd.Name, cmd.Description, cmd.Category, cmd.Price, cmd.Stock)
    if err != nil {
        return nil, fmt.Errorf("invalid product: %w", err)
    }

    // Persist via repository
    if err := h.repo.Save(ctx, p); err != nil {
        return nil, fmt.Errorf("save product: %w", err)
    }

    h.logger.Info("product created", "id", p.ID(), "name", p.Name())

    return dto.ProductFromDomain(p), nil
}
```

#### Command: Update Product

```go
// internal/application/command/update_product.go
package command

import (
    "context"
    "fmt"

    "github.com/google/uuid"
    "github.com/yourorg/app/internal/application/dto"
    "github.com/yourorg/app/internal/domain/product"
)

// UpdateProductCommand contains the data needed to update a product.
type UpdateProductCommand struct {
    ID          uuid.UUID
    Name        *string
    Description *string
    Price       *int
    Category    *string
    Stock       *int
}

// UpdateProductHandler handles the UpdateProductCommand.
type UpdateProductHandler struct {
    repo product.Repository
}

func NewUpdateProductHandler(repo product.Repository) *UpdateProductHandler {
    return &UpdateProductHandler{repo: repo}
}

// Handle executes the update product command.
func (h *UpdateProductHandler) Handle(ctx context.Context, cmd UpdateProductCommand) (*dto.ProductDTO, error) {
    // Fetch existing product
    p, err := h.repo.FindByID(ctx, cmd.ID)
    if err != nil {
        return nil, err
    }

    // Apply updates via domain methods
    if cmd.Price != nil {
        if err := p.UpdatePrice(*cmd.Price); err != nil {
            return nil, err
        }
    }
    // ... apply other updates via domain methods

    // Persist changes
    if err := h.repo.Save(ctx, p); err != nil {
        return nil, fmt.Errorf("save product: %w", err)
    }

    return dto.ProductFromDomain(p), nil
}
```

#### Command: Delete Product

```go
// internal/application/command/delete_product.go
package command

import (
    "context"

    "github.com/google/uuid"
    "github.com/yourorg/app/internal/domain/product"
)

// DeleteProductCommand contains the ID of the product to delete.
type DeleteProductCommand struct {
    ID uuid.UUID
}

// DeleteProductHandler handles the DeleteProductCommand.
type DeleteProductHandler struct {
    repo product.Repository
}

func NewDeleteProductHandler(repo product.Repository) *DeleteProductHandler {
    return &DeleteProductHandler{repo: repo}
}

// Handle executes the delete product command.
func (h *DeleteProductHandler) Handle(ctx context.Context, cmd DeleteProductCommand) error {
    return h.repo.Delete(ctx, cmd.ID)
}
```

#### Query: Get Product

```go
// internal/application/query/get_product.go
package query

import (
    "context"

    "github.com/google/uuid"
    "github.com/yourorg/app/internal/application/dto"
    "github.com/yourorg/app/internal/domain/product"
)

// GetProductQuery contains the ID of the product to retrieve.
type GetProductQuery struct {
    ID uuid.UUID
}

// GetProductHandler handles the GetProductQuery.
type GetProductHandler struct {
    repo product.Repository
}

func NewGetProductHandler(repo product.Repository) *GetProductHandler {
    return &GetProductHandler{repo: repo}
}

// Handle executes the get product query.
func (h *GetProductHandler) Handle(ctx context.Context, q GetProductQuery) (*dto.ProductDTO, error) {
    p, err := h.repo.FindByID(ctx, q.ID)
    if err != nil {
        return nil, err
    }
    return dto.ProductFromDomain(p), nil
}
```

#### Query: List Products

```go
// internal/application/query/list_products.go
package query

import (
    "context"

    "github.com/yourorg/app/internal/application/dto"
    "github.com/yourorg/app/internal/domain/product"
)

// ListProductsQuery contains pagination options.
type ListProductsQuery struct {
    Page     int
    PerPage  int
    Category string
}

// ListProductsResult contains the paginated results.
type ListProductsResult struct {
    Products []*dto.ProductDTO
    Total    int
}

// ListProductsHandler handles the ListProductsQuery.
type ListProductsHandler struct {
    repo product.Repository
}

func NewListProductsHandler(repo product.Repository) *ListProductsHandler {
    return &ListProductsHandler{repo: repo}
}

// Handle executes the list products query.
func (h *ListProductsHandler) Handle(ctx context.Context, q ListProductsQuery) (*ListProductsResult, error) {
    products, total, err := h.repo.FindAll(ctx, product.ListOptions{
        Page:     q.Page,
        PerPage:  q.PerPage,
        Category: q.Category,
    })
    if err != nil {
        return nil, err
    }

    dtos := make([]*dto.ProductDTO, len(products))
    for i, p := range products {
        dtos[i] = dto.ProductFromDomain(p)
    }

    return &ListProductsResult{Products: dtos, Total: total}, nil
}
```

#### DTO: Product Data Transfer Object

```go
// internal/application/dto/product_dto.go
package dto

import (
    "time"

    "github.com/google/uuid"
    "github.com/yourorg/app/internal/domain/product"
)

// ProductDTO is the data transfer object for API responses.
type ProductDTO struct {
    ID          uuid.UUID `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Price       int       `json:"price"`
    Category    string    `json:"category"`
    Stock       int       `json:"stock"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// ProductFromDomain converts a domain entity to DTO using getter methods.
func ProductFromDomain(p *product.Product) *ProductDTO {
    return &ProductDTO{
        ID:          p.ID(),
        Name:        p.Name(),
        Description: p.Description(),
        Price:       p.Price(),
        Category:    p.Category(),
        Stock:       p.Stock(),
        CreatedAt:   p.CreatedAt(),
        UpdatedAt:   p.UpdatedAt(),
    }
}
```

### Step 6: Implement the HTTP Handler (Interface Layer)

Create the HTTP handler that uses command/query handlers.

```go
// internal/interfaces/http/handler/product_handler.go
package handler

import (
    "encoding/json"
    "errors"
    "net/http"
    "strconv"

    "github.com/go-chi/chi/v5"
    "github.com/go-playground/validator/v10"
    "github.com/google/uuid"
    "github.com/yourorg/app/internal/application/command"
    "github.com/yourorg/app/internal/application/query"
    "github.com/yourorg/app/internal/domain/product"
    "github.com/yourorg/app/pkg/response"
)

var validate = validator.New()

// ProductHandler handles HTTP requests for products.
type ProductHandler struct {
    createHandler *command.CreateProductHandler
    updateHandler *command.UpdateProductHandler
    deleteHandler *command.DeleteProductHandler
    getHandler    *query.GetProductHandler
    listHandler   *query.ListProductsHandler
}

// NewProductHandler creates a new product handler with CQRS handlers.
func NewProductHandler(
    createHandler *command.CreateProductHandler,
    updateHandler *command.UpdateProductHandler,
    deleteHandler *command.DeleteProductHandler,
    getHandler *query.GetProductHandler,
    listHandler *query.ListProductsHandler,
) *ProductHandler {
    return &ProductHandler{
        createHandler: createHandler,
        updateHandler: updateHandler,
        deleteHandler: deleteHandler,
        getHandler:    getHandler,
        listHandler:   listHandler,
    }
}

func (h *ProductHandler) RegisterRoutes(r chi.Router) {
    r.Route("/products", func(r chi.Router) {
        r.Get("/", h.List)
        r.Post("/", h.Create)
        r.Get("/{id}", h.Get)
        r.Put("/{id}", h.Update)
        r.Delete("/{id}", h.Delete)
    })
}

// List handles GET /products
func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    page, _ := strconv.Atoi(r.URL.Query().Get("page"))
    if page < 1 {
        page = 1
    }
    perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
    if perPage < 1 || perPage > 100 {
        perPage = 20
    }

    result, err := h.listHandler.Handle(ctx, query.ListProductsQuery{
        Page:    page,
        PerPage: perPage,
    })
    if err != nil {
        response.InternalError(w, err)
        return
    }

    response.JSONWithMeta(w, http.StatusOK, result.Products, &response.Meta{
        Page:    page,
        PerPage: perPage,
        Total:   result.Total,
    })
}

// Get handles GET /products/{id}
func (h *ProductHandler) Get(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    id, err := uuid.Parse(chi.URLParam(r, "id"))
    if err != nil {
        response.BadRequest(w, "Invalid product ID")
        return
    }

    result, err := h.getHandler.Handle(ctx, query.GetProductQuery{ID: id})
    if err != nil {
        if errors.Is(err, product.ErrNotFound) {
            response.NotFound(w, "Product not found")
            return
        }
        response.InternalError(w, err)
        return
    }

    response.JSON(w, http.StatusOK, result)
}

// CreateProductRequest is the request body for creating a product.
type CreateProductRequest struct {
    Name        string `json:"name" validate:"required,min=1,max=255"`
    Description string `json:"description" validate:"max=1000"`
    Price       int    `json:"price" validate:"required,gt=0"`
    Category    string `json:"category" validate:"required,min=1,max=100"`
    Stock       int    `json:"stock" validate:"gte=0"`
}

// Create handles POST /products
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    var req CreateProductRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        response.BadRequest(w, "Invalid request body")
        return
    }

    if err := validate.Struct(req); err != nil {
        response.ValidationError(w, formatValidationErrors(err.(validator.ValidationErrors)))
        return
    }

    // Convert request to command
    result, err := h.createHandler.Handle(ctx, command.CreateProductCommand{
        Name:        req.Name,
        Description: req.Description,
        Price:       req.Price,
        Category:    req.Category,
        Stock:       req.Stock,
    })
    if err != nil {
        response.InternalError(w, err)
        return
    }

    response.JSON(w, http.StatusCreated, result)
}

// UpdateProductRequest is the request body for updating a product.
type UpdateProductRequest struct {
    Name        *string `json:"name" validate:"omitempty,min=1,max=255"`
    Description *string `json:"description" validate:"omitempty,max=1000"`
    Price       *int    `json:"price" validate:"omitempty,gt=0"`
    Category    *string `json:"category" validate:"omitempty,min=1,max=100"`
    Stock       *int    `json:"stock" validate:"omitempty,gte=0"`
}

// Update handles PUT /products/{id}
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    id, err := uuid.Parse(chi.URLParam(r, "id"))
    if err != nil {
        response.BadRequest(w, "Invalid product ID")
        return
    }

    var req UpdateProductRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        response.BadRequest(w, "Invalid request body")
        return
    }

    if err := validate.Struct(req); err != nil {
        response.ValidationError(w, formatValidationErrors(err.(validator.ValidationErrors)))
        return
    }

    // Convert request to command
    result, err := h.updateHandler.Handle(ctx, command.UpdateProductCommand{
        ID:          id,
        Name:        req.Name,
        Description: req.Description,
        Price:       req.Price,
        Category:    req.Category,
        Stock:       req.Stock,
    })
    if err != nil {
        if errors.Is(err, product.ErrNotFound) {
            response.NotFound(w, "Product not found")
            return
        }
        response.InternalError(w, err)
        return
    }

    response.JSON(w, http.StatusOK, result)
}

// Delete handles DELETE /products/{id}
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    id, err := uuid.Parse(chi.URLParam(r, "id"))
    if err != nil {
        response.BadRequest(w, "Invalid product ID")
        return
    }

    if err := h.deleteHandler.Handle(ctx, command.DeleteProductCommand{ID: id}); err != nil {
        if errors.Is(err, product.ErrNotFound) {
            response.NotFound(w, "Product not found")
            return
        }
        response.InternalError(w, err)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
```

### Step 7: Write Tests

Create comprehensive tests for the command/query handlers.

#### Command Handler Test

```go
// internal/application/command/create_product_test.go
package command_test

import (
    "context"
    "log/slog"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/yourorg/app/internal/application/command"
    "github.com/yourorg/app/internal/domain/product"
)

// MockRepository is a mock implementation of product.Repository
type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) FindByID(ctx context.Context, id uuid.UUID) (*product.Product, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*product.Product), args.Error(1)
}

func (m *MockRepository) Save(ctx context.Context, p *product.Product) error {
    args := m.Called(ctx, p)
    return args.Error(0)
}

// ... other mock methods

func TestCreateProductHandler_Handle(t *testing.T) {
    tests := []struct {
        name    string
        cmd     command.CreateProductCommand
        setup   func(*MockRepository)
        wantErr bool
    }{
        {
            name: "success",
            cmd: command.CreateProductCommand{
                Name:     "Test Product",
                Price:    1000,
                Category: "Electronics",
                Stock:    10,
            },
            setup: func(m *MockRepository) {
                m.On("Save", mock.Anything, mock.AnythingOfType("*product.Product")).Return(nil)
            },
            wantErr: false,
        },
        {
            name: "validation error - empty name",
            cmd: command.CreateProductCommand{
                Name:     "",
                Price:    1000,
                Category: "Electronics",
            },
            setup:   func(m *MockRepository) {},
            wantErr: true,
        },
        {
            name: "validation error - invalid price",
            cmd: command.CreateProductCommand{
                Name:     "Test Product",
                Price:    -100,
                Category: "Electronics",
            },
            setup:   func(m *MockRepository) {},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            repo := new(MockRepository)
            tt.setup(repo)
            handler := command.NewCreateProductHandler(repo, slog.Default())

            result, err := handler.Handle(context.Background(), tt.cmd)

            if tt.wantErr {
                assert.Error(t, err)
                assert.Nil(t, result)
                return
            }

            assert.NoError(t, err)
            assert.NotNil(t, result)
            assert.NotEqual(t, uuid.Nil, result.ID)
            assert.Equal(t, tt.cmd.Name, result.Name)
            assert.Equal(t, tt.cmd.Price, result.Price)
            repo.AssertExpectations(t)
        })
    }
}
```

#### Query Handler Test

```go
// internal/application/query/get_product_test.go
package query_test

import (
    "context"
    "testing"

    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "github.com/yourorg/app/internal/application/query"
    "github.com/yourorg/app/internal/domain/product"
)

func TestGetProductHandler_Handle(t *testing.T) {
    testID := uuid.New()
    existingProduct := product.Reconstitute(
        testID, "Test Product", "Description", "Electronics",
        1000, 10, time.Now(), time.Now(),
    )

    tests := []struct {
        name    string
        query   query.GetProductQuery
        setup   func(*MockRepository)
        wantErr error
    }{
        {
            name:  "success",
            query: query.GetProductQuery{ID: testID},
            setup: func(m *MockRepository) {
                m.On("FindByID", mock.Anything, testID).Return(existingProduct, nil)
            },
            wantErr: nil,
        },
        {
            name:  "not found",
            query: query.GetProductQuery{ID: uuid.New()},
            setup: func(m *MockRepository) {
                m.On("FindByID", mock.Anything, mock.Anything).Return(nil, product.ErrNotFound)
            },
            wantErr: product.ErrNotFound,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            repo := new(MockRepository)
            tt.setup(repo)
            handler := query.NewGetProductHandler(repo)

            result, err := handler.Handle(context.Background(), tt.query)

            if tt.wantErr != nil {
                assert.ErrorIs(t, err, tt.wantErr)
                return
            }

            assert.NoError(t, err)
            assert.NotNil(t, result)
            assert.Equal(t, existingProduct.ID(), result.ID)
            repo.AssertExpectations(t)
        })
    }
}
```

#### HTTP Handler Test

```go
// internal/interfaces/http/handler/product_handler_test.go
package handler_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/go-chi/chi/v5"
    "github.com/stretchr/testify/assert"
    "github.com/yourorg/app/internal/interfaces/http/handler"
)

func TestProductHandler_Create(t *testing.T) {
    // Setup mock handlers
    createHandler := setupMockCreateHandler(t)
    h := handler.NewProductHandler(createHandler, nil, nil, nil, nil)

    router := chi.NewRouter()
    h.RegisterRoutes(router)

    body := `{"name": "Test Product", "price": 1000, "category": "Electronics", "stock": 10}`
    req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString(body))
    req.Header.Set("Content-Type", "application/json")
    rec := httptest.NewRecorder()

    router.ServeHTTP(rec, req)

    assert.Equal(t, http.StatusCreated, rec.Code)

    var response map[string]interface{}
    err := json.NewDecoder(rec.Body).Decode(&response)
    assert.NoError(t, err)
    assert.Equal(t, "Test Product", response["name"])
}
```

## Checklist

- [ ] Domain entity with private fields + getters (DDD)
- [ ] Migration created with golang-migrate (.up.sql and .down.sql files)
- [ ] Repository interface (port) in domain layer
- [ ] Repository implementation (adapter) in repository/ package with Reconstitute
- [ ] Command handlers (Create, Update, Delete)
- [ ] Query handlers (Get, List)
- [ ] DTOs for API responses
- [ ] HTTP handler using command/query handlers
- [ ] Request validation structs
- [ ] Routes registered
- [ ] Unit tests for command handlers
- [ ] Unit tests for query handlers
- [ ] HTTP handler tests
