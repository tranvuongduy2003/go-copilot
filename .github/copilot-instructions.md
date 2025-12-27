# Copilot Global Instructions

This document provides global instructions that are always loaded by GitHub Copilot for this project.

## Project Overview

This is a production-ready full-stack application for building modern web services. The backend follows **Clean Architecture + Domain-Driven Design (DDD) + CQRS** patterns, while the frontend uses **React 19 with Tailwind CSS v4 and shadcn/ui** components. The application demonstrates best practices for scalable, maintainable enterprise software.

## Tech Stack

| Layer | Technology | Version | Notes |
|-------|------------|---------|-------|
| **Backend** | Go | 1.25+ | Clean Architecture + DDD + CQRS |
| **Router** | Chi | v5 | Lightweight, idiomatic |
| **Database** | PostgreSQL | 16+ | With pgx v5 driver |
| **Migrations** | golang-migrate | v4 | CLI-based migrations |
| **Frontend** | React | 19 | With TypeScript 5.x strict |
| **Styling** | Tailwind CSS | v4 | CSS-first configuration |
| **Components** | shadcn/ui | Latest | new-york style |
| **State** | TanStack Query + Zustand | Latest | Server + client state |
| **Forms** | React Hook Form + Zod | Latest | Validation |
| **Testing** | Vitest + testify | Latest | Frontend + Backend |

## Available Scripts & Commands

### Backend Commands

```bash
# Run the API server
cd backend && go run cmd/api/main.go

# Run tests
cd backend && go test ./...

# Run tests with coverage
cd backend && go test -cover ./...

# Run linter
cd backend && golangci-lint run

# Build binary
cd backend && go build -o bin/api cmd/api/main.go
```

### Database Migrations (golang-migrate CLI)

```bash
# Install golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Create new migration (creates .up.sql and .down.sql files)
migrate create -ext sql -dir backend/migrations -seq <name>

# Apply all migrations
migrate -path backend/migrations -database "$DATABASE_URL" up

# Rollback last migration
migrate -path backend/migrations -database "$DATABASE_URL" down 1

# Check version
migrate -path backend/migrations -database "$DATABASE_URL" version
```

### Frontend Commands

```bash
# Install dependencies
cd frontend && npm install

# Start dev server
cd frontend && npm run dev

# Run tests
cd frontend && npm test

# Build for production
cd frontend && npm run build

# Run linter
cd frontend && npm run lint
```

### Docker Commands

```bash
# Start all services
docker-compose up -d

# Stop all services
docker-compose down

# View logs
docker-compose logs -f api
```

---

## Coding Standards

### Go Coding Standards

1. **Package Naming**: Use short, lowercase names without underscores
2. **File Naming**: Use `snake_case.go` for file names
3. **Error Handling**: Always handle errors explicitly; never ignore them
4. **Context**: Pass `context.Context` as the first parameter
5. **Interfaces**: Define interfaces where they are used, not where they are implemented
6. **Structs**: Use pointer receivers for methods that modify state
7. **Comments**: Write godoc comments for all exported types, functions, and methods

```go
// Good: Error handling
result, err := service.DoSomething(ctx, input)
if err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}

// Good: Context as first parameter
func (s *Service) GetUser(ctx context.Context, id string) (*User, error)
```

### TypeScript/React Coding Standards

1. **File Naming**: Use `kebab-case.tsx` for components, `camelCase.ts` for utilities
2. **Component Naming**: Use PascalCase for component names
3. **Props**: Define explicit interfaces for all component props
4. **Hooks**: Prefix custom hooks with `use`
5. **Types**: Prefer `interface` over `type` for object shapes
6. **Imports**: Use absolute imports from `@/` alias
7. **Exports**: Use named exports for components

```tsx
// Good: Component with typed props
interface ButtonProps {
  variant?: 'primary' | 'secondary' | 'ghost';
  size?: 'sm' | 'md' | 'lg';
  children: React.ReactNode;
  onClick?: () => void;
}

export function Button({ variant = 'primary', size = 'md', children, onClick }: ButtonProps) {
  return (
    <button
      className={cn(buttonVariants({ variant, size }))}
      onClick={onClick}
    >
      {children}
    </button>
  );
}
```

---

## Design System Rules

### CRITICAL: UI Consistency

All UI components MUST follow the defined design system. Never use arbitrary colors, spacing, or typography.

### Color Palette (OKLCH)

Use these colors exclusively. Reference them via CSS custom properties or Tailwind classes.

| Token | OKLCH Value | Usage |
|-------|-------------|-------|
| `--primary` | `oklch(0.7 0.15 290)` | Primary actions, links, focus states |
| `--primary-dark` | `oklch(0.6 0.2 280)` | Primary hover states, gradients |
| `--secondary` | `oklch(0.75 0.15 220)` | Secondary actions, accents |
| `--background-dark` | `oklch(0.1 0.01 260)` | Dark mode background |
| `--background-light` | `oklch(0.98 0.01 260)` | Light mode background |
| `--success` | `oklch(0.7 0.17 160)` | Success states, confirmations |
| `--warning` | `oklch(0.8 0.15 85)` | Warning states, cautions |
| `--error` | `oklch(0.65 0.2 15)` | Error states, destructive actions |

```tsx
// CORRECT: Using design system colors
<button className="bg-primary hover:bg-primary-dark text-primary-foreground">

// WRONG: Using arbitrary colors
<button className="bg-purple-500 hover:bg-purple-600 text-white">
```

### Typography

- **Font Family**: Inter (sans-serif), JetBrains Mono (monospace)
- **Scale**: 12, 14, 16, 18, 20, 24, 30, 36, 48, 60, 72px
- **Weights**: 400 (normal), 500 (medium), 600 (semibold), 700 (bold)

```tsx
// Use Tailwind typography classes
<h1 className="text-4xl font-bold tracking-tight">
<p className="text-base text-muted-foreground">
<code className="font-mono text-sm">
```

### Spacing Scale

Use the defined spacing scale: 4, 8, 12, 16, 20, 24, 32, 40, 48, 64, 80, 96, 128px

```tsx
// CORRECT: Using spacing scale
<div className="p-4 gap-6 mt-8">

// WRONG: Using arbitrary spacing
<div className="p-[13px] gap-[7px] mt-[15px]">
```

### Border Radius

| Token | Value | Usage |
|-------|-------|-------|
| `sm` | 4px | Small elements, badges |
| `md` | 8px | Buttons, inputs |
| `lg` | 12px | Cards, dialogs |
| `xl` | 16px | Large cards, modals |
| `2xl` | 24px | Feature sections |
| `full` | 9999px | Pills, avatars |

### Shadows

Use violet-tinted shadows for elevated elements:

```css
--shadow-sm: 0 1px 2px oklch(0.3 0.05 290 / 0.05);
--shadow-md: 0 4px 6px oklch(0.3 0.05 290 / 0.1);
--shadow-lg: 0 10px 15px oklch(0.3 0.05 290 / 0.15);
```

---

## Project Structure

### Backend (DDD + CQRS Architecture)

```
backend/
├── cmd/
│   └── api/
│       └── main.go                    # Entry point, dependency wiring
├── internal/
│   ├── domain/                        # Domain Layer (innermost, pure business logic)
│   │   ├── user/                      # User aggregate
│   │   │   ├── user.go                # Entity with private fields + getters
│   │   │   ├── repository.go          # Repository interface (port)
│   │   │   ├── errors.go              # Domain-specific errors
│   │   │   └── events.go              # Domain events
│   │   └── shared/                    # Shared domain concepts
│   │       ├── errors.go
│   │       └── valueobjects.go
│   │
│   ├── application/                   # Application Layer (CQRS)
│   │   ├── command/                   # Commands (write operations)
│   │   │   ├── create_user.go
│   │   │   └── update_user.go
│   │   ├── query/                     # Queries (read operations)
│   │   │   ├── get_user.go
│   │   │   └── list_users.go
│   │   └── dto/                       # Data Transfer Objects
│   │       └── user_dto.go
│   │
│   ├── infrastructure/                # Infrastructure Layer (adapters)
│   │   ├── persistence/
│   │   │   ├── postgres/              # Database utilities
│   │   │   │   ├── connection.go
│   │   │   │   ├── unit_of_work.go
│   │   │   │   ├── query_builder.go
│   │   │   │   └── errors.go
│   │   │   └── repository/            # Repository implementations
│   │   │       └── user_repository.go # Implements domain.UserRepository
│   │   └── cache/
│   │       └── redis/
│   │
│   └── interfaces/                    # Interface Adapters Layer
│       └── http/
│           ├── handler/
│           │   └── user_handler.go
│           ├── middleware/
│           └── router/
│
├── migrations/                        # golang-migrate migrations
│   ├── 000001_create_users_table.up.sql
│   └── 000001_create_users_table.down.sql
└── pkg/
    ├── config/
    ├── logger/
    └── validator/
```

### Frontend (React)

```
frontend/
├── src/
│   ├── components/
│   │   ├── ui/               # shadcn/ui components
│   │   │   ├── button.tsx
│   │   │   └── card.tsx
│   │   └── features/         # Feature-specific components
│   │       └── user-profile/
│   │           ├── user-profile.tsx
│   │           └── user-avatar.tsx
│   ├── hooks/
│   │   └── use-user.ts
│   ├── lib/
│   │   ├── api.ts
│   │   └── utils.ts
│   ├── pages/
│   │   └── dashboard.tsx
│   ├── stores/
│   │   └── auth-store.ts
│   ├── styles/
│   │   └── globals.css
│   └── types/
│       └── user.ts
```

---

## Error Handling Patterns

### Go Error Handling

```go
// Define domain errors
var (
    ErrUserNotFound = errors.New("user not found")
    ErrInvalidInput = errors.New("invalid input")
)

// Wrap errors with context
func (s *UserService) GetUser(ctx context.Context, id string) (*User, error) {
    user, err := s.repo.FindByID(ctx, id)
    if err != nil {
        if errors.Is(err, repository.ErrNotFound) {
            return nil, ErrUserNotFound
        }
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    return user, nil
}
```

### React Error Handling

```tsx
// Use error boundaries for component errors
import { ErrorBoundary } from '@/components/error-boundary';

function App() {
  return (
    <ErrorBoundary fallback={<ErrorFallback />}>
      <Routes />
    </ErrorBoundary>
  );
}

// Use TanStack Query for API errors
const { data, error, isLoading } = useQuery({
  queryKey: ['user', userId],
  queryFn: () => fetchUser(userId),
});

if (error) {
  return <ErrorMessage error={error} />;
}
```

---

## API Design Principles

### RESTful Endpoints

```
GET    /api/v1/users          # List users
POST   /api/v1/users          # Create user
GET    /api/v1/users/:id      # Get user by ID
PUT    /api/v1/users/:id      # Update user
DELETE /api/v1/users/:id      # Delete user
```

### Response Format

```json
{
  "data": { ... },
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 100
  }
}
```

### Error Response Format

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input data",
    "details": [
      { "field": "email", "message": "Invalid email format" }
    ]
  }
}
```

---

## Security Best Practices

### Backend Security

1. **Input Validation**: Validate all input at the handler level
2. **SQL Injection**: Use parameterized queries exclusively
3. **Authentication**: Use JWT with short expiration, refresh tokens
4. **Authorization**: Implement RBAC at the service layer
5. **Rate Limiting**: Apply rate limiting middleware
6. **CORS**: Configure strict CORS policies
7. **Secrets**: Never log or expose secrets

### Frontend Security

1. **XSS Prevention**: Use React's built-in escaping, avoid `dangerouslySetInnerHTML`
2. **CSRF**: Include CSRF tokens in state-changing requests
3. **Sensitive Data**: Never store sensitive data in localStorage
4. **API Keys**: Never expose API keys in client-side code
5. **Content Security Policy**: Configure strict CSP headers

---

## Performance Guidelines

### Backend Performance

1. Use connection pooling for database connections
2. Implement caching for frequently accessed data
3. Use pagination for list endpoints
4. Optimize database queries with proper indexes
5. Profile and benchmark critical paths

### Frontend Performance

1. Use React.lazy for code splitting
2. Implement virtual scrolling for long lists
3. Memoize expensive computations
4. Optimize images with next-gen formats
5. Use TanStack Query for efficient data fetching and caching

---

## Commit Message Format

Follow Conventional Commits:

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`, `perf`

Examples:
- `feat(auth): add OAuth2 login support`
- `fix(api): handle null values in user response`
- `docs(readme): update installation instructions`
