---
name: go-api-builder
description: Build REST API endpoints in Go with handlers, services, repositories, and tests. Use when creating new API endpoints.
---

# Go API Builder Skill

This skill guides you through building complete REST API endpoints in Go following clean architecture patterns.

## When to Use This Skill

- Creating a new API endpoint (GET, POST, PUT, DELETE)
- Adding a new resource to the API
- Implementing CRUD operations for a domain entity

## Step-by-Step Process

### Step 1: Define the Domain Model

Create the domain model and input/output types in `internal/domain/`.

```go
// internal/domain/product.go
package domain

import "time"

// Product represents a product in the catalog.
type Product struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Price       int       `json:"price"` // Price in cents
    Category    string    `json:"category"`
    Stock       int       `json:"stock"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// CreateProductInput represents input for creating a product.
type CreateProductInput struct {
    Name        string `json:"name" validate:"required,min=1,max=255"`
    Description string `json:"description" validate:"max=5000"`
    Price       int    `json:"price" validate:"required,min=1"`
    Category    string `json:"category" validate:"required,max=100"`
    Stock       int    `json:"stock" validate:"min=0"`
}

// UpdateProductInput represents input for updating a product.
type UpdateProductInput struct {
    Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
    Description *string `json:"description,omitempty" validate:"omitempty,max=5000"`
    Price       *int    `json:"price,omitempty" validate:"omitempty,min=1"`
    Category    *string `json:"category,omitempty" validate:"omitempty,max=100"`
    Stock       *int    `json:"stock,omitempty" validate:"omitempty,min=0"`
}
```

### Step 2: Create the Database Migration

Create migration files in `migrations/`.

```sql
-- migrations/XXX_create_products.up.sql
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

```sql
-- migrations/XXX_create_products.down.sql
DROP TABLE IF EXISTS products;
```

### Step 3: Define Repository Interface (Domain Layer)

Create the repository interface in the domain layer (port):

```go
// internal/domain/product/repository.go
package product

import (
    "context"

    "github.com/yourorg/app/internal/domain"
)

type ProductRepository interface {
    FindByID(ctx context.Context, id string) (*domain.Product, error)
    FindAll(ctx context.Context, opts ListOptions) ([]*domain.Product, int, error)
    Create(ctx context.Context, product *domain.Product) error
    Update(ctx context.Context, product *domain.Product) error
    Delete(ctx context.Context, id string) error
}
```

### Step 4: Implement Repository (Infrastructure Layer)

```go
// internal/infrastructure/persistence/postgres/product_repository.go
package postgres

import (
    "context"
    "errors"
    "fmt"
    "time"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/yourorg/app/internal/domain"
    "github.com/yourorg/app/internal/repository"
)

type productRepository struct {
    db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) repository.ProductRepository {
    return &productRepository{db: db}
}

func (r *productRepository) FindByID(ctx context.Context, id string) (*domain.Product, error) {
    query := `
        SELECT id, name, description, price, category, stock, created_at, updated_at
        FROM products
        WHERE id = $1 AND deleted_at IS NULL
    `

    var p domain.Product
    err := r.db.QueryRow(ctx, query, id).Scan(
        &p.ID, &p.Name, &p.Description, &p.Price,
        &p.Category, &p.Stock, &p.CreatedAt, &p.UpdatedAt,
    )

    if errors.Is(err, pgx.ErrNoRows) {
        return nil, domain.ErrNotFound
    }
    if err != nil {
        return nil, fmt.Errorf("query product: %w", err)
    }

    return &p, nil
}

func (r *productRepository) FindAll(ctx context.Context, opts repository.ListOptions) ([]*domain.Product, int, error) {
    // Count total
    var total int
    countQuery := `SELECT COUNT(*) FROM products WHERE deleted_at IS NULL`
    if err := r.db.QueryRow(ctx, countQuery).Scan(&total); err != nil {
        return nil, 0, fmt.Errorf("count products: %w", err)
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

    rows, err := r.db.Query(ctx, query, opts.PerPage, offset)
    if err != nil {
        return nil, 0, fmt.Errorf("query products: %w", err)
    }
    defer rows.Close()

    var products []*domain.Product
    for rows.Next() {
        var p domain.Product
        if err := rows.Scan(
            &p.ID, &p.Name, &p.Description, &p.Price,
            &p.Category, &p.Stock, &p.CreatedAt, &p.UpdatedAt,
        ); err != nil {
            return nil, 0, fmt.Errorf("scan product: %w", err)
        }
        products = append(products, &p)
    }

    return products, total, nil
}

func (r *productRepository) Create(ctx context.Context, product *domain.Product) error {
    query := `
        INSERT INTO products (id, name, description, price, category, stock, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `

    _, err := r.db.Exec(ctx, query,
        product.ID, product.Name, product.Description, product.Price,
        product.Category, product.Stock, product.CreatedAt, product.UpdatedAt,
    )
    if err != nil {
        return fmt.Errorf("insert product: %w", err)
    }

    return nil
}

func (r *productRepository) Update(ctx context.Context, product *domain.Product) error {
    query := `
        UPDATE products
        SET name = $2, description = $3, price = $4, category = $5, stock = $6, updated_at = $7
        WHERE id = $1 AND deleted_at IS NULL
    `

    result, err := r.db.Exec(ctx, query,
        product.ID, product.Name, product.Description, product.Price,
        product.Category, product.Stock, time.Now().UTC(),
    )
    if err != nil {
        return fmt.Errorf("update product: %w", err)
    }

    if result.RowsAffected() == 0 {
        return domain.ErrNotFound
    }

    return nil
}

func (r *productRepository) Delete(ctx context.Context, id string) error {
    query := `UPDATE products SET deleted_at = $2 WHERE id = $1 AND deleted_at IS NULL`

    result, err := r.db.Exec(ctx, query, id, time.Now().UTC())
    if err != nil {
        return fmt.Errorf("delete product: %w", err)
    }

    if result.RowsAffected() == 0 {
        return domain.ErrNotFound
    }

    return nil
}
```

### Step 5: Implement Command/Query Handlers (Application Layer)

Create the CQRS handlers with business logic.

```go
// internal/application/command/create_product.go
package command

import (
    "context"
    "fmt"
    "time"

    "github.com/google/uuid"
    "github.com/yourorg/app/internal/domain"
    "github.com/yourorg/app/internal/repository"
    "log/slog"
)

type ProductService struct {
    repo   repository.ProductRepository
    logger *slog.Logger
}

func NewProductService(repo repository.ProductRepository, logger *slog.Logger) *ProductService {
    return &ProductService{repo: repo, logger: logger}
}

func (s *ProductService) GetProduct(ctx context.Context, id string) (*domain.Product, error) {
    return s.repo.FindByID(ctx, id)
}

func (s *ProductService) ListProducts(ctx context.Context, opts repository.ListOptions) ([]*domain.Product, int, error) {
    return s.repo.FindAll(ctx, opts)
}

func (s *ProductService) CreateProduct(ctx context.Context, input domain.CreateProductInput) (*domain.Product, error) {
    now := time.Now().UTC()
    product := &domain.Product{
        ID:          uuid.New().String(),
        Name:        input.Name,
        Description: input.Description,
        Price:       input.Price,
        Category:    input.Category,
        Stock:       input.Stock,
        CreatedAt:   now,
        UpdatedAt:   now,
    }

    if err := s.repo.Create(ctx, product); err != nil {
        return nil, fmt.Errorf("create product: %w", err)
    }

    s.logger.Info("product created", "id", product.ID, "name", product.Name)

    return product, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, id string, input domain.UpdateProductInput) (*domain.Product, error) {
    product, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }

    // Apply updates
    if input.Name != nil {
        product.Name = *input.Name
    }
    if input.Description != nil {
        product.Description = *input.Description
    }
    if input.Price != nil {
        product.Price = *input.Price
    }
    if input.Category != nil {
        product.Category = *input.Category
    }
    if input.Stock != nil {
        product.Stock = *input.Stock
    }
    product.UpdatedAt = time.Now().UTC()

    if err := s.repo.Update(ctx, product); err != nil {
        return nil, fmt.Errorf("update product: %w", err)
    }

    return product, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id string) error {
    return s.repo.Delete(ctx, id)
}
```

### Step 6: Implement the HTTP Handler (Interface Layer)

Create the HTTP handler with validation.

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
    "github.com/yourorg/app/internal/domain"
    "github.com/yourorg/app/internal/repository"
    "github.com/yourorg/app/internal/service"
    "github.com/yourorg/app/pkg/response"
)

var validate = validator.New()

type ProductHandler struct {
    svc *service.ProductService
}

func NewProductHandler(svc *service.ProductService) *ProductHandler {
    return &ProductHandler{svc: svc}
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

    products, total, err := h.svc.ListProducts(ctx, repository.ListOptions{
        Page:    page,
        PerPage: perPage,
    })
    if err != nil {
        response.InternalError(w, err)
        return
    }

    response.JSONWithMeta(w, http.StatusOK, products, &response.Meta{
        Page:    page,
        PerPage: perPage,
        Total:   total,
    })
}

func (h *ProductHandler) Get(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    id := chi.URLParam(r, "id")

    product, err := h.svc.GetProduct(ctx, id)
    if err != nil {
        if errors.Is(err, domain.ErrNotFound) {
            response.NotFound(w, "Product not found")
            return
        }
        response.InternalError(w, err)
        return
    }

    response.JSON(w, http.StatusOK, product)
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    var input domain.CreateProductInput
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        response.BadRequest(w, "Invalid request body")
        return
    }

    if err := validate.Struct(input); err != nil {
        response.ValidationError(w, formatValidationErrors(err.(validator.ValidationErrors)))
        return
    }

    product, err := h.svc.CreateProduct(ctx, input)
    if err != nil {
        response.InternalError(w, err)
        return
    }

    response.JSON(w, http.StatusCreated, product)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    id := chi.URLParam(r, "id")

    var input domain.UpdateProductInput
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        response.BadRequest(w, "Invalid request body")
        return
    }

    if err := validate.Struct(input); err != nil {
        response.ValidationError(w, formatValidationErrors(err.(validator.ValidationErrors)))
        return
    }

    product, err := h.svc.UpdateProduct(ctx, id, input)
    if err != nil {
        if errors.Is(err, domain.ErrNotFound) {
            response.NotFound(w, "Product not found")
            return
        }
        response.InternalError(w, err)
        return
    }

    response.JSON(w, http.StatusOK, product)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    id := chi.URLParam(r, "id")

    if err := h.svc.DeleteProduct(ctx, id); err != nil {
        if errors.Is(err, domain.ErrNotFound) {
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

```go
// internal/application/command/create_product_test.go
package command_test

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/yourorg/app/internal/domain"
    "github.com/yourorg/app/internal/service"
)

func TestProductService_CreateProduct(t *testing.T) {
    tests := []struct {
        name    string
        input   domain.CreateProductInput
        setup   func(*MockProductRepository)
        wantErr bool
    }{
        {
            name: "success",
            input: domain.CreateProductInput{
                Name:     "Test Product",
                Price:    1000,
                Category: "Electronics",
            },
            setup: func(m *MockProductRepository) {
                m.On("Create", mock.Anything, mock.AnythingOfType("*domain.Product")).Return(nil)
            },
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            repo := new(MockProductRepository)
            tt.setup(repo)
            svc := service.NewProductService(repo, slog.Default())

            product, err := svc.CreateProduct(context.Background(), tt.input)

            if tt.wantErr {
                assert.Error(t, err)
                return
            }

            assert.NoError(t, err)
            assert.NotEmpty(t, product.ID)
            assert.Equal(t, tt.input.Name, product.Name)
        })
    }
}
```

## Checklist

- [ ] Domain model defined with validation tags
- [ ] Migration created (up and down)
- [ ] Repository interface defined
- [ ] PostgreSQL repository implemented
- [ ] Service layer with business logic
- [ ] HTTP handler with validation
- [ ] Routes registered
- [ ] Unit tests for service
- [ ] Handler tests
- [ ] Integration tests (optional)
