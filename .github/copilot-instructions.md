# Copilot Global Instructions

This document provides global instructions that are always loaded by GitHub Copilot for this project.

## Project Overview

This is a full-stack application built with:

- **Backend**: Go 1.25 with clean architecture, REST APIs, PostgreSQL
- **Frontend**: React 19 + TypeScript + Tailwind CSS v4 + shadcn/ui (new-york style)
- **Infrastructure**: Docker, Docker Compose, GitHub Actions

## Tech Stack Details

### Backend (Go 1.25)

- **Framework**: Standard library `net/http` with Chi router
- **Database**: PostgreSQL with `pgx` driver
- **Migrations**: `golang-migrate`
- **Validation**: `go-playground/validator`
- **Logging**: `slog` (structured logging)
- **Configuration**: Environment variables with `envconfig`
- **Testing**: Standard `testing` package with `testify`

### Frontend (React 19)

- **Build Tool**: Vite
- **Language**: TypeScript 5.x (strict mode)
- **Styling**: Tailwind CSS v4 with CSS-first configuration
- **Components**: shadcn/ui (new-york style)
- **State Management**: TanStack Query for server state, Zustand for client state
- **Forms**: React Hook Form with Zod validation
- **Routing**: TanStack Router (type-safe)
- **Testing**: Vitest + React Testing Library

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
| `--secondary` | `oklch(0.75 0.15 200)` | Secondary actions, accents |
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

## File Naming Conventions

### Backend (Go)

```
backend/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── domain/
│   │   └── user.go           # Domain models
│   ├── handlers/
│   │   └── user_handler.go   # HTTP handlers
│   ├── middleware/
│   │   └── auth.go
│   ├── repository/
│   │   └── user_repository.go
│   └── service/
│       └── user_service.go
├── migrations/
│   └── 001_create_users.up.sql
└── pkg/
    └── response/
        └── response.go
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
