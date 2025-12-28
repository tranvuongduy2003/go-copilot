# Go API Builder Skill

Generate complete REST API endpoints following DDD + CQRS patterns.

## Usage

```
/project:skill:go-api <entity-name>
```

## Generated Files

For an entity named `Product`, this skill generates:

### 1. Domain Layer

**`internal/domain/product/product.go`**
```go
package product

import (
    "time"

    "github.com/google/uuid"
)

type Product struct {
    id          uuid.UUID
    name        string
    description string
    price       int64
    status      Status
    createdAt   time.Time
    updatedAt   time.Time
}

func NewProduct(name, description string, price int64) (*Product, error) {
    if name == "" {
        return nil, ErrInvalidName
    }
    if price < 0 {
        return nil, ErrInvalidPrice
    }

    now := time.Now()
    return &Product{
        id:          uuid.New(),
        name:        name,
        description: description,
        price:       price,
        status:      StatusActive,
        createdAt:   now,
        updatedAt:   now,
    }, nil
}

func (product *Product) ID() uuid.UUID          { return product.id }
func (product *Product) Name() string           { return product.name }
func (product *Product) Description() string    { return product.description }
func (product *Product) Price() int64           { return product.price }
func (product *Product) Status() Status         { return product.status }
func (product *Product) CreatedAt() time.Time   { return product.createdAt }
func (product *Product) UpdatedAt() time.Time   { return product.updatedAt }

func (product *Product) UpdateName(name string) error {
    if name == "" {
        return ErrInvalidName
    }
    product.name = name
    product.updatedAt = time.Now()
    return nil
}

func Reconstitute(
    id uuid.UUID,
    name, description string,
    price int64,
    status Status,
    createdAt, updatedAt time.Time,
) *Product {
    return &Product{
        id:          id,
        name:        name,
        description: description,
        price:       price,
        status:      status,
        createdAt:   createdAt,
        updatedAt:   updatedAt,
    }
}
```

**`internal/domain/product/repository.go`**
```go
package product

import (
    "context"

    "github.com/google/uuid"
)

type Repository interface {
    FindByID(context context.Context, id uuid.UUID) (*Product, error)
    FindAll(context context.Context, options ListOptions) ([]*Product, int, error)
    Save(context context.Context, product *Product) error
    Update(context context.Context, product *Product) error
    Delete(context context.Context, id uuid.UUID) error
}

type ListOptions struct {
    Page    int
    PerPage int
    SortBy  string
    Order   string
}
```

**`internal/domain/product/errors.go`**
```go
package product

import "errors"

var (
    ErrNotFound     = errors.New("product not found")
    ErrInvalidName  = errors.New("invalid product name")
    ErrInvalidPrice = errors.New("invalid product price")
)
```

**`internal/domain/product/status.go`**
```go
package product

type Status string

const (
    StatusActive   Status = "active"
    StatusInactive Status = "inactive"
    StatusDeleted  Status = "deleted"
)

func (status Status) String() string {
    return string(status)
}

func (status Status) IsValid() bool {
    switch status {
    case StatusActive, StatusInactive, StatusDeleted:
        return true
    }
    return false
}
```

### 2. Application Layer (CQRS, Domain-Aligned)

The application layer is organized by domain/bounded context. Each domain has its own `command/`, `query/`, and `dto/` packages.

**`internal/application/product/command/create_product.go`**
```go
package productcommand

import (
    "context"
    "fmt"

    productdto "yourapp/internal/application/product/dto"
    "yourapp/internal/domain/product"
)

type CreateProductCommand struct {
    Name        string
    Description string
    Price       int64
}

type CreateProductHandler struct {
    repository product.Repository
}

func NewCreateProductHandler(repository product.Repository) *CreateProductHandler {
    return &CreateProductHandler{repository: repository}
}

func (handler *CreateProductHandler) Handle(context context.Context, command CreateProductCommand) (*productdto.ProductDTO, error) {
    newProduct, err := product.NewProduct(command.Name, command.Description, command.Price)
    if err != nil {
        return nil, fmt.Errorf("invalid product: %w", err)
    }

    if err := handler.repository.Save(context, newProduct); err != nil {
        return nil, fmt.Errorf("save product: %w", err)
    }

    return productdto.ProductFromDomain(newProduct), nil
}
```

**`internal/application/product/query/get_product.go`**
```go
package productquery

import (
    "context"

    "github.com/google/uuid"
    productdto "yourapp/internal/application/product/dto"
    "yourapp/internal/domain/product"
)

type GetProductQuery struct {
    ID uuid.UUID
}

type GetProductHandler struct {
    repository product.Repository
}

func NewGetProductHandler(repository product.Repository) *GetProductHandler {
    return &GetProductHandler{repository: repository}
}

func (handler *GetProductHandler) Handle(context context.Context, query GetProductQuery) (*productdto.ProductDTO, error) {
    foundProduct, err := handler.repository.FindByID(context, query.ID)
    if err != nil {
        return nil, err
    }
    return productdto.ProductFromDomain(foundProduct), nil
}
```

**`internal/application/product/dto/product_dto.go`**
```go
package productdto

import (
    "time"

    "github.com/google/uuid"
    "yourapp/internal/domain/product"
)

type ProductDTO struct {
    ID          uuid.UUID `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Price       int64     `json:"price"`
    Status      string    `json:"status"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}

func ProductFromDomain(domainProduct *product.Product) *ProductDTO {
    return &ProductDTO{
        ID:          domainProduct.ID(),
        Name:        domainProduct.Name(),
        Description: domainProduct.Description(),
        Price:       domainProduct.Price(),
        Status:      domainProduct.Status().String(),
        CreatedAt:   domainProduct.CreatedAt(),
        UpdatedAt:   domainProduct.UpdatedAt(),
    }
}
```

### 3. Infrastructure Layer

**`internal/infrastructure/persistence/repository/product_repository.go`**
```go
package repository

import (
    "context"
    "fmt"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
    "yourapp/internal/domain/product"
)

type productRepository struct {
    pool *pgxpool.Pool
}

func NewProductRepository(pool *pgxpool.Pool) product.Repository {
    return &productRepository{pool: pool}
}

func (repository *productRepository) FindByID(context context.Context, id uuid.UUID) (*product.Product, error) {
    query := `
        SELECT id, name, description, price, status, created_at, updated_at
        FROM products
        WHERE id = $1 AND deleted_at IS NULL
    `

    var (
        productID                uuid.UUID
        name, description        string
        price                    int64
        status                   string
        createdAt, updatedAt     time.Time
    )

    err := repository.pool.QueryRow(context, query, id).Scan(
        &productID, &name, &description, &price, &status, &createdAt, &updatedAt,
    )
    if err == pgx.ErrNoRows {
        return nil, product.ErrNotFound
    }
    if err != nil {
        return nil, fmt.Errorf("query product: %w", err)
    }

    return product.Reconstitute(
        productID, name, description, price,
        product.Status(status), createdAt, updatedAt,
    ), nil
}

func (repository *productRepository) Save(context context.Context, productEntity *product.Product) error {
    query := `
        INSERT INTO products (id, name, description, price, status, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `

    _, err := repository.pool.Exec(context, query,
        productEntity.ID(),
        productEntity.Name(),
        productEntity.Description(),
        productEntity.Price(),
        productEntity.Status().String(),
        productEntity.CreatedAt(),
        productEntity.UpdatedAt(),
    )
    if err != nil {
        return fmt.Errorf("insert product: %w", err)
    }

    return nil
}
```

### 4. Interface Layer

**`internal/interfaces/http/handler/product_handler.go`**
```go
package handler

import (
    "encoding/json"
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/google/uuid"
    productcommand "yourapp/internal/application/product/command"
    productquery "yourapp/internal/application/product/query"
    "yourapp/internal/domain/product"
    "yourapp/internal/interfaces/http/response"
)

type ProductHandler struct {
    createHandler *productcommand.CreateProductHandler
    getHandler    *productquery.GetProductHandler
    listHandler   *productquery.ListProductsHandler
}

func NewProductHandler(
    createHandler *productcommand.CreateProductHandler,
    getHandler *productquery.GetProductHandler,
    listHandler *productquery.ListProductsHandler,
) *ProductHandler {
    return &ProductHandler{
        createHandler: createHandler,
        getHandler:    getHandler,
        listHandler:   listHandler,
    }
}

func (handler *ProductHandler) RegisterRoutes(router chi.Router) {
    router.Route("/products", func(router chi.Router) {
        router.Get("/", handler.List)
        router.Post("/", handler.Create)
        router.Get("/{id}", handler.Get)
        router.Put("/{id}", handler.Update)
        router.Delete("/{id}", handler.Delete)
    })
}

type CreateProductRequest struct {
    Name        string `json:"name"`
    Description string `json:"description"`
    Price       int64  `json:"price"`
}

func (handler *ProductHandler) Create(writer http.ResponseWriter, request *http.Request) {
    var createRequest CreateProductRequest
    if err := json.NewDecoder(request.Body).Decode(&createRequest); err != nil {
        response.BadRequest(writer, "Invalid request body")
        return
    }

    result, err := handler.createHandler.Handle(request.Context(), productcommand.CreateProductCommand{
        Name:        createRequest.Name,
        Description: createRequest.Description,
        Price:       createRequest.Price,
    })
    if err != nil {
        response.HandleError(writer, err)
        return
    }

    response.JSON(writer, http.StatusCreated, result)
}

func (handler *ProductHandler) Get(writer http.ResponseWriter, request *http.Request) {
    idParam := chi.URLParam(request, "id")
    id, err := uuid.Parse(idParam)
    if err != nil {
        response.BadRequest(writer, "Invalid product ID")
        return
    }

    result, err := handler.getHandler.Handle(request.Context(), productquery.GetProductQuery{ID: id})
    if err != nil {
        if err == product.ErrNotFound {
            response.NotFound(writer, "Product not found")
            return
        }
        response.HandleError(writer, err)
        return
    }

    response.JSON(writer, http.StatusOK, result)
}
```

### 5. Migration

**`migrations/XXXXXX_create_products_table.up.sql`**
```sql
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price BIGINT NOT NULL CHECK (price >= 0),
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_products_status ON products(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_products_created_at ON products(created_at DESC);
```

**`migrations/XXXXXX_create_products_table.down.sql`**
```sql
DROP TABLE IF EXISTS products;
```

## Checklist

- [ ] Migration files created
- [ ] Domain entity with private fields + getters
- [ ] Repository interface in domain layer
- [ ] Domain errors defined
- [ ] Status enum defined
- [ ] Command handlers (Create, Update, Delete)
- [ ] Query handlers (Get, List)
- [ ] DTOs for API responses
- [ ] Repository implementation with Reconstitute
- [ ] HTTP handler with routes
- [ ] Request validation
- [ ] Tests
