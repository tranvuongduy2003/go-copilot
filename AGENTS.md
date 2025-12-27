# AI Agent Instructions

This document provides global instructions for AI coding assistants (GitHub Copilot, Claude, Cursor) working on this codebase.

## Project Overview

Full-stack application following **Clean Architecture + Domain-Driven Design (DDD) + CQRS** patterns.

| Layer | Technology | Version | Notes |
|-------|------------|---------|-------|
| **Backend** | Go | 1.25+ | Clean Architecture + DDD + CQRS |
| **Router** | Chi | v5 | Lightweight, idiomatic |
| **Database** | PostgreSQL | 16+ | With pgx v5 driver |
| **Migrations** | golang-migrate | v4 | CLI-based SQL migrations |
| **Frontend** | React | 19 | With TypeScript 5.x strict |
| **Styling** | Tailwind CSS | v4 | CSS-first configuration |
| **Components** | shadcn/ui | Latest | new-york style |
| **State** | TanStack Query + Zustand | Latest | Server + client state |
| **Testing** | Vitest + testify | Latest | Frontend + Backend |

---

## Quick Commands

### Backend

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

# Rollback all migrations
migrate -path backend/migrations -database "$DATABASE_URL" down

# Check current version
migrate -path backend/migrations -database "$DATABASE_URL" version

# Force set version (for fixing dirty state)
migrate -path backend/migrations -database "$DATABASE_URL" force <version>
```

### Frontend

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

### Docker

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

### Naming Conventions (CRITICAL)

**Use meaningful, descriptive names instead of comments. No abbreviations allowed.**

```go
// BAD - Abbreviations
func (r *repo) FindAll(ctx context.Context, opts ListOpts) ([]*User, int, error)
var usrRepo UserRepository
var cfg *Config

// GOOD - Full words
func (repository *userRepository) FindAll(ctx context.Context, options ListOptions) ([]*User, int, error)
var userRepository UserRepository
var configuration *Configuration
```

```tsx
// BAD - Abbreviations
const [usr, setUsr] = useState<User | null>(null);
const handleBtnClick = () => { ... };

// GOOD - Full words
const [user, setUser] = useState<User | null>(null);
const handleButtonClick = () => { ... };
```

### Comments Policy (CRITICAL)

**Do NOT write comments unless absolutely necessary. Code should be self-documenting.**

Only add comments for:
- Complex algorithms that cannot be simplified
- Legal/license requirements
- TODO/FIXME with ticket references

---

## Architecture Boundaries

### Always Do

**Backend (DDD + CQRS)**
- Follow dependency rule: domain <- application <- infrastructure <- interfaces
- Use CQRS pattern: separate Command handlers (writes) from Query handlers (reads)
- Define repository interfaces (ports) in domain layer
- Implement repositories in infrastructure layer
- Pass `context.Context` as first parameter to all functions
- Wrap errors with context: `fmt.Errorf("failed to do X: %w", err)`
- Use value objects for domain concepts (Email, Role, etc.)
- **Use full, descriptive names** - no abbreviations (`repository` not `repo`)
- **Write self-documenting code** - meaningful names instead of comments

**Frontend (React + Design System)**
- Use design system colors: `bg-primary`, `text-foreground`, `bg-destructive`
- Use standard spacing: `p-4`, `gap-6`, `mt-8` (4px base unit)
- Use shadcn/ui components from `@/components/ui/`
- Create TypeScript types that match Go domain models
- Use React Query for server state, Zustand for client state
- Handle loading, error, and empty states in all components
- **Use full, descriptive names** - no abbreviations (`button` not `btn`)
- **Write self-documenting code** - meaningful names instead of comments

### Ask First

- Before creating new database tables or migrations
- Before adding new external dependencies
- Before making breaking API changes
- Before modifying authentication/authorization flows
- Before creating new aggregates or domain entities
- When multiple valid implementation approaches exist

### Never Do

**Naming & Comments**
- Never use abbreviations in names (`usr`, `repo`, `cfg`, `opts`, `btn`, `msg`)
- Never write comments that explain "what" - make code self-documenting
- Never add comments to obvious code - delete unnecessary comments

**Backend**
- Never put business logic in HTTP handlers (use application layer)
- Never import infrastructure packages in domain layer
- Never ignore errors - always handle explicitly
- Never use string concatenation for SQL queries (use parameterized)
- Never log sensitive data (passwords, tokens, PII)

**Frontend**
- Never use arbitrary colors: `bg-[#7c3aed]`, `bg-purple-500`
- Never use arbitrary spacing: `p-[13px]`, `mt-[7px]`
- Never use `any` type - provide proper TypeScript types
- Never store sensitive data in localStorage
- Never use `dangerouslySetInnerHTML` without sanitization

**Security**
- Never hardcode secrets or credentials
- Never commit .env files with real values
- Never approve code with SQL injection vulnerabilities

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
│   │   │   │   ├── connection.go      # Connection pool management
│   │   │   │   ├── unit_of_work.go    # Transaction support
│   │   │   │   ├── query_builder.go   # SQL query helpers
│   │   │   │   └── errors.go          # Database error types
│   │   │   └── repository/            # Repository implementations
│   │   │       └── user_repository.go # Implements domain.UserRepository
│   │   ├── messaging/
│   │   │   └── memory/
│   │   │       └── event_bus.go       # In-memory event bus
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
├── migrations/                            # golang-migrate migrations
│   ├── 000001_create_users_table.up.sql
│   ├── 000001_create_users_table.down.sql
│   ├── 000002_add_users_indexes.up.sql
│   └── 000002_add_users_indexes.down.sql
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

## Design System

### Color Palette (OKLCH)

| Token | Usage |
|-------|-------|
| `primary` | Primary actions, links, focus states (violet hue 290) |
| `primary-dark` | Primary hover states, gradients |
| `secondary` | Secondary actions, accents (cyan hue 220) |
| `destructive` | Error states, destructive actions |
| `muted` | Subtle backgrounds, secondary text |
| `success` | Success states, confirmations |
| `warning` | Warning states, cautions |

```tsx
// CORRECT: Using design system colors
<button className="bg-primary hover:bg-primary/90 text-primary-foreground">

// WRONG: Using arbitrary colors
<button className="bg-purple-500 hover:bg-purple-600 text-white">
```

### Spacing Scale (4px base)

| Token | Value | Usage |
|-------|-------|-------|
| `1` | 4px | Tight padding |
| `2` | 8px | Inline elements |
| `3` | 12px | Medium-small |
| `4` | 16px | Default component padding |
| `6` | 24px | Card spacing |
| `8` | 32px | Section padding |
| `12` | 48px | Page sections |

```tsx
// CORRECT: Using spacing scale
<div className="p-4 gap-6 mt-8">

// WRONG: Using arbitrary spacing
<div className="p-[13px] gap-[7px] mt-[15px]">
```

### Border Radius

| Token | Value | Usage |
|-------|-------|-------|
| `sm` | 4px | Badges, chips |
| `md` | 8px | Buttons, inputs (default) |
| `lg` | 12px | Cards, dialogs |
| `xl` | 16px | Large cards, modals |
| `full` | 9999px | Pills, avatars |

---

## Code Patterns

### Go Error Handling

```go
result, err := service.DoSomething(ctx, input)
if err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}
```

### Go Domain Entity

```go
// Entity with private fields and getters
type User struct {
    id        uuid.UUID
    email     Email     // Value Object
    name      string
    createdAt time.Time
}

func (u *User) ID() uuid.UUID { return u.id }
func (u *User) Email() Email  { return u.email }
```

### React Component with Props

```tsx
interface UserCardProps {
  user: User;
  onEdit?: () => void;
}

export function UserCard({ user, onEdit }: UserCardProps) {
  return (
    <Card className="p-4">
      <CardHeader>
        <CardTitle>{user.name}</CardTitle>
      </CardHeader>
      <CardContent>
        <p className="text-muted-foreground">{user.email}</p>
      </CardContent>
    </Card>
  );
}
```

### React Query Hook

```tsx
export function useUsers() {
  return useQuery({
    queryKey: ['users'],
    queryFn: () => api.get<User[]>('/users'),
  });
}

export function useCreateUser() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (input: CreateUserInput) => api.post<User>('/users', input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
    },
  });
}
```

---

## Available Agents

| Agent | Description |
|-------|-------------|
| `backend-engineer` | Go backend with Clean Architecture + DDD + CQRS |
| `frontend-engineer` | React 19 + Tailwind CSS v4 + shadcn/ui |
| `fullstack-engineer` | End-to-end feature development |
| `test-agent` | Comprehensive test suites (Go + React) |
| `code-reviewer` | Code review for quality, security, design system |
| `security-auditor` | OWASP Top 10 vulnerability auditing |
| `documentation-writer` | Technical documentation |
| `technical-planner` | Feature planning and technical designs |

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
- `feat(api): add user registration endpoint`
- `fix(ui): resolve button focus state in dark mode`
- `docs(readme): update installation instructions`

---

## References

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [DDD + CQRS in Go](https://threedots.tech/post/ddd-cqrs-clean-architecture-combined/)
- [W3C Design Tokens](https://tr.designtokens.org/format/)
- [shadcn/ui](https://ui.shadcn.com/)
- [Tailwind CSS v4](https://tailwindcss.com/docs/v4-beta)
