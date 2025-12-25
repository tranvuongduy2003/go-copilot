# System Architecture

This document provides a comprehensive overview of our full-stack application architecture, covering both the Go backend and React frontend.

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                         Client Layer                            │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │              React 19 + TypeScript Frontend              │   │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐    │   │
│  │  │ Pages   │  │Components│  │ Hooks   │  │ Stores  │    │   │
│  │  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘    │   │
│  │       │           │           │           │            │   │
│  │  ┌────┴───────────┴───────────┴───────────┴────┐       │   │
│  │  │           TanStack Query + Zustand           │       │   │
│  │  └──────────────────────┬──────────────────────┘       │   │
│  └─────────────────────────┼───────────────────────────────┘   │
└────────────────────────────┼────────────────────────────────────┘
                             │ HTTP/REST
┌────────────────────────────┼────────────────────────────────────┐
│                         API Layer                               │
│  ┌─────────────────────────┼───────────────────────────────┐   │
│  │              Go 1.25 Backend (Clean Architecture)        │   │
│  │  ┌─────────────────────┴─────────────────────────┐      │   │
│  │  │              HTTP Handlers (chi router)        │      │   │
│  │  └─────────────────────┬─────────────────────────┘      │   │
│  │  ┌─────────────────────┴─────────────────────────┐      │   │
│  │  │               Service Layer                    │      │   │
│  │  └─────────────────────┬─────────────────────────┘      │   │
│  │  ┌─────────────────────┴─────────────────────────┐      │   │
│  │  │             Repository Layer                   │      │   │
│  │  └─────────────────────┬─────────────────────────┘      │   │
│  └─────────────────────────┼───────────────────────────────┘   │
└────────────────────────────┼────────────────────────────────────┘
                             │
┌────────────────────────────┼────────────────────────────────────┐
│                       Data Layer                                │
│  ┌─────────────┐  ┌────────┴────────┐  ┌─────────────┐         │
│  │  PostgreSQL │  │      Redis      │  │   S3/Blob   │         │
│  │  (Primary)  │  │    (Cache)      │  │  (Storage)  │         │
│  └─────────────┘  └─────────────────┘  └─────────────┘         │
└─────────────────────────────────────────────────────────────────┘
```

## Backend Architecture (Go)

### Clean Architecture Layers

```
backend/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── domain/                  # Enterprise business rules
│   │   ├── entity/              # Core entities
│   │   ├── valueobject/         # Value objects
│   │   └── errors/              # Domain errors
│   ├── application/             # Application business rules
│   │   ├── service/             # Use cases/services
│   │   ├── dto/                 # Data transfer objects
│   │   └── port/                # Interface definitions
│   ├── infrastructure/          # Frameworks & drivers
│   │   ├── repository/          # Database implementations
│   │   ├── cache/               # Caching implementations
│   │   └── external/            # External service clients
│   └── interface/               # Interface adapters
│       ├── http/                # HTTP handlers
│       │   ├── handler/         # Request handlers
│       │   ├── middleware/      # HTTP middleware
│       │   └── router/          # Route definitions
│       └── dto/                 # API request/response DTOs
├── pkg/                         # Shared utilities
│   ├── config/                  # Configuration management
│   ├── logger/                  # Logging utilities
│   └── validator/               # Validation helpers
└── migrations/                  # Database migrations
```

### Dependency Flow

```
         ┌─────────────────┐
         │    Handlers     │  Interface Layer
         └────────┬────────┘
                  │ uses
         ┌────────▼────────┐
         │    Services     │  Application Layer
         └────────┬────────┘
                  │ uses
         ┌────────▼────────┐
         │  Repositories   │  Infrastructure Layer
         └────────┬────────┘
                  │
         ┌────────▼────────┐
         │    Database     │  External Systems
         └─────────────────┘

Dependencies point inward. Domain layer has no external dependencies.
```

### Domain Layer

The innermost layer containing enterprise-wide business rules.

```go
// internal/domain/entity/user.go
package entity

import (
    "time"

    "github.com/google/uuid"
)

type User struct {
    ID        uuid.UUID
    Email     string
    Name      string
    Role      Role
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Role string

const (
    RoleAdmin  Role = "admin"
    RoleMember Role = "member"
)

func NewUser(email, name string) (*User, error) {
    if email == "" {
        return nil, ErrInvalidEmail
    }
    return &User{
        ID:        uuid.New(),
        Email:     email,
        Name:      name,
        Role:      RoleMember,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }, nil
}
```

### Application Layer

Contains application-specific business rules and orchestrates domain entities.

```go
// internal/application/port/repository.go
package port

import (
    "context"

    "github.com/google/uuid"
    "project/internal/domain/entity"
)

type UserRepository interface {
    Create(ctx context.Context, user *entity.User) error
    GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
    GetByEmail(ctx context.Context, email string) (*entity.User, error)
    Update(ctx context.Context, user *entity.User) error
    Delete(ctx context.Context, id uuid.UUID) error
    List(ctx context.Context, opts ListOptions) ([]*entity.User, error)
}

// internal/application/service/user_service.go
package service

type UserService struct {
    userRepo port.UserRepository
    cache    port.Cache
    logger   *slog.Logger
}

func NewUserService(
    userRepo port.UserRepository,
    cache port.Cache,
    logger *slog.Logger,
) *UserService {
    return &UserService{
        userRepo: userRepo,
        cache:    cache,
        logger:   logger,
    }
}

func (s *UserService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
    user, err := entity.NewUser(req.Email, req.Name)
    if err != nil {
        return nil, fmt.Errorf("create user entity: %w", err)
    }

    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, fmt.Errorf("save user: %w", err)
    }

    return dto.UserResponseFromEntity(user), nil
}
```

### Infrastructure Layer

Implements interfaces defined in the application layer.

```go
// internal/infrastructure/repository/postgres/user_repository.go
package postgres

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
    query := `
        INSERT INTO users (id, email, name, role, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
    _, err := r.db.ExecContext(ctx, query,
        user.ID, user.Email, user.Name, user.Role, user.CreatedAt, user.UpdatedAt,
    )
    return err
}
```

### Interface Layer (HTTP)

```go
// internal/interface/http/handler/user_handler.go
package handler

type UserHandler struct {
    userService *service.UserService
    logger      *slog.Logger
}

func NewUserHandler(userService *service.UserService, logger *slog.Logger) *UserHandler {
    return &UserHandler{
        userService: userService,
        logger:      logger,
    }
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req dto.CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, http.StatusBadRequest, "invalid request body")
        return
    }

    if err := validate.Struct(req); err != nil {
        respondError(w, http.StatusBadRequest, formatValidationError(err))
        return
    }

    user, err := h.userService.CreateUser(r.Context(), req)
    if err != nil {
        h.logger.Error("failed to create user", "error", err)
        respondError(w, http.StatusInternalServerError, "failed to create user")
        return
    }

    respondJSON(w, http.StatusCreated, user)
}
```

## Frontend Architecture (React)

### Directory Structure

```
frontend/
├── public/                      # Static assets
├── src/
│   ├── app/                     # App router pages
│   │   ├── (auth)/              # Auth route group
│   │   │   ├── login/
│   │   │   └── register/
│   │   ├── (dashboard)/         # Dashboard route group
│   │   │   ├── layout.tsx
│   │   │   └── page.tsx
│   │   ├── layout.tsx           # Root layout
│   │   └── page.tsx             # Home page
│   ├── components/
│   │   ├── ui/                  # shadcn/ui components
│   │   ├── features/            # Feature components
│   │   ├── layouts/             # Layout components
│   │   └── shared/              # Shared components
│   ├── hooks/                   # Custom hooks
│   │   ├── use-auth.ts
│   │   └── use-media-query.ts
│   ├── lib/                     # Utilities
│   │   ├── api/                 # API client
│   │   ├── utils.ts             # Helper functions
│   │   └── validations/         # Zod schemas
│   ├── stores/                  # Zustand stores
│   │   ├── auth-store.ts
│   │   └── ui-store.ts
│   ├── types/                   # TypeScript types
│   │   ├── api.ts
│   │   └── models.ts
│   └── styles/
│       └── globals.css          # Global styles
├── tailwind.config.ts
└── tsconfig.json
```

### State Management

#### Server State (TanStack Query)

```typescript
// src/lib/api/users.ts
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from './client';

export const userKeys = {
  all: ['users'] as const,
  lists: () => [...userKeys.all, 'list'] as const,
  list: (filters: UserFilters) => [...userKeys.lists(), filters] as const,
  details: () => [...userKeys.all, 'detail'] as const,
  detail: (id: string) => [...userKeys.details(), id] as const,
};

export function useUsers(filters: UserFilters) {
  return useQuery({
    queryKey: userKeys.list(filters),
    queryFn: () => api.get<User[]>('/users', { params: filters }),
  });
}

export function useUser(id: string) {
  return useQuery({
    queryKey: userKeys.detail(id),
    queryFn: () => api.get<User>(`/users/${id}`),
    enabled: !!id,
  });
}

export function useCreateUser() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateUserInput) => api.post<User>('/users', data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: userKeys.lists() });
    },
  });
}
```

#### Client State (Zustand)

```typescript
// src/stores/auth-store.ts
import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  login: (user: User, token: string) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      token: null,
      isAuthenticated: false,
      login: (user, token) =>
        set({ user, token, isAuthenticated: true }),
      logout: () =>
        set({ user: null, token: null, isAuthenticated: false }),
    }),
    {
      name: 'auth-storage',
    }
  )
);
```

### API Client

```typescript
// src/lib/api/client.ts
import ky from 'ky';
import { useAuthStore } from '@/stores/auth-store';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export const api = ky.create({
  prefixUrl: `${API_BASE_URL}/api/v1`,
  timeout: 30000,
  hooks: {
    beforeRequest: [
      (request) => {
        const token = useAuthStore.getState().token;
        if (token) {
          request.headers.set('Authorization', `Bearer ${token}`);
        }
      },
    ],
    afterResponse: [
      async (_request, _options, response) => {
        if (response.status === 401) {
          useAuthStore.getState().logout();
          window.location.href = '/login';
        }
        return response;
      },
    ],
  },
});

// Type-safe wrapper
export async function apiGet<T>(path: string, options?: Parameters<typeof api.get>[1]): Promise<T> {
  return api.get(path, options).json<T>();
}

export async function apiPost<T>(path: string, data?: unknown): Promise<T> {
  return api.post(path, { json: data }).json<T>();
}

export async function apiPut<T>(path: string, data?: unknown): Promise<T> {
  return api.put(path, { json: data }).json<T>();
}

export async function apiDelete<T>(path: string): Promise<T> {
  return api.delete(path).json<T>();
}
```

## Data Flow

### Request Flow

```
User Action
    │
    ▼
React Component
    │
    ▼ (TanStack Query mutation/query)
API Client (ky)
    │
    ▼ (HTTP Request)
Go HTTP Handler
    │
    ▼ (Validate & Parse)
Service Layer
    │
    ▼ (Business Logic)
Repository Layer
    │
    ▼ (SQL Query)
PostgreSQL
    │
    ▼ (Response)
[Reverse path back to UI]
```

### Authentication Flow

```
1. User submits credentials
    │
    ▼
2. POST /api/v1/auth/login
    │
    ▼
3. Validate credentials against DB
    │
    ▼
4. Generate JWT token
    │
    ▼
5. Return token + user data
    │
    ▼
6. Store in Zustand (persisted)
    │
    ▼
7. Subsequent requests include Bearer token
```

## Database Schema

### Core Tables

```sql
-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'member',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);

-- Audit log
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(100) NOT NULL,
    entity_id UUID,
    changes JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at DESC);
```

### Migration Strategy

```go
// migrations/001_initial_schema.up.sql
// migrations/001_initial_schema.down.sql

// Run migrations
goose -dir migrations postgres "postgres://..." up
```

## Security Considerations

### Authentication

- JWT tokens with short expiry (15 minutes)
- Refresh token rotation
- Secure HTTP-only cookies for web clients
- Rate limiting on auth endpoints

### API Security

- Input validation at handler level
- SQL parameterized queries (no string concatenation)
- CORS configuration for allowed origins
- Request size limits
- Timeout middleware

### Data Protection

- Passwords hashed with bcrypt (cost 12)
- Sensitive data encrypted at rest
- TLS for all connections
- Audit logging for sensitive operations

## Performance Optimizations

### Backend

- Connection pooling for database
- Redis caching for frequently accessed data
- Pagination for list endpoints
- Database indexing on query patterns
- Graceful shutdown handling

### Frontend

- Code splitting with dynamic imports
- Image optimization
- TanStack Query caching and deduplication
- Virtualized lists for large datasets
- Optimistic updates for better UX

## Monitoring & Observability

### Logging

```go
// Structured logging with slog
logger.Info("user created",
    "user_id", user.ID,
    "email", user.Email,
    "duration_ms", duration.Milliseconds(),
)
```

### Metrics

- Request duration histograms
- Error rate counters
- Database connection pool stats
- Cache hit/miss rates

### Tracing

- OpenTelemetry integration
- Distributed trace context propagation
- Span creation for database queries

## Deployment

### Docker Compose (Development)

```yaml
version: '3.8'
services:
  api:
    build: ./backend
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://user:pass@db:5432/app
      - REDIS_URL=redis://cache:6379
    depends_on:
      - db
      - cache

  frontend:
    build: ./frontend
    ports:
      - "3000:3000"
    environment:
      - NEXT_PUBLIC_API_URL=http://localhost:8080

  db:
    image: postgres:16
    volumes:
      - postgres_data:/var/lib/postgresql/data
    environment:
      - POSTGRES_USER=user
      - POSTGRES_PASSWORD=pass
      - POSTGRES_DB=app

  cache:
    image: redis:7-alpine

volumes:
  postgres_data:
```

### Production Considerations

- Kubernetes deployment with horizontal scaling
- Load balancer with health checks
- Database read replicas
- CDN for static assets
- Secrets management (Vault/AWS Secrets Manager)
