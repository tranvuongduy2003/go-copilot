# Claude Code Project Instructions

> This document provides global instructions for Claude Code (Anthropic's AI coding assistant) working on this codebase. These instructions complement the existing GitHub Copilot configuration in `.github/`.

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

# Build binary
cd backend && go build -o bin/api cmd/api/main.go
```

### Database Migrations (golang-migrate CLI)

```bash
# Create new migration
migrate create -ext sql -dir backend/migrations -seq <name>

# Apply all migrations
migrate -path backend/migrations -database "$DATABASE_URL" up

# Rollback last migration
migrate -path backend/migrations -database "$DATABASE_URL" down 1

# Check version
migrate -path backend/migrations -database "$DATABASE_URL" version
```

### Frontend

```bash
# Start dev server
cd frontend && npm run dev

# Run tests
cd frontend && npm test

# Build for production
cd frontend && npm run build
```

### Docker

```bash
# Start all services
docker-compose up -d

# Stop all services
docker-compose down
```

---

## Coding Standards

### Naming Conventions

**CRITICAL: Use meaningful, descriptive names instead of comments.**

- **No abbreviations** in function names, variable names, parameters, or type names
- Names should be self-documenting and reveal intent
- If you need a comment to explain what code does, rename it instead

```go
// BAD - Abbreviations
func GetUsrByID(ctx context.Context, id uuid.UUID) (*Usr, error)
func (r *repo) FindAll(ctx context.Context, opts ListOpts) ([]*User, int, error)
var usrRepo UserRepository
var cfg *Config

// GOOD - Full words
func GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
func (repository *userRepository) FindAll(ctx context.Context, options ListOptions) ([]*User, int, error)
var userRepository UserRepository
var configuration *Configuration
```

```tsx
// BAD - Abbreviations
const [usr, setUsr] = useState<User | null>(null);
const handleBtnClick = () => { ... };
interface Props { ... }
function Btn({ ...props }) { ... }

// GOOD - Full words
const [user, setUser] = useState<User | null>(null);
const handleButtonClick = () => { ... };
interface UserCardProps { ... }
function Button({ ...props }) { ... }
```

### Comments Policy

**CRITICAL: Do NOT write comments unless absolutely necessary.**

- Code should be self-explanatory through meaningful names
- Only add comments for:
  - Complex algorithms that cannot be simplified
  - Legal/license requirements
  - TODO/FIXME with ticket references
  - Public API documentation (godoc, JSDoc) when required
- Delete comments that explain "what" - the code should show that
- If code needs explanation, refactor it to be clearer

```go
// BAD - Unnecessary comments
// GetUserByID retrieves a user by their ID
func GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
    // Find the user in the repository
    user, err := repository.FindByID(ctx, id)
    if err != nil {
        return nil, err // Return error if not found
    }
    return user, nil // Return the user
}

// GOOD - Self-documenting code, no comments needed
func GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
    user, err := repository.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    return user, nil
}
```

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

**Frontend (React + Design System)**
- Use design system colors: `bg-primary`, `text-foreground`, `bg-destructive`
- Use standard spacing: `p-4`, `gap-6`, `mt-8` (4px base unit)
- Use shadcn/ui components from `@/components/ui/`
- Create TypeScript types that match Go domain models
- Use React Query for server state, Zustand for client state
- Handle loading, error, and empty states in all components

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
├── cmd/api/main.go                    # Entry point, dependency wiring
├── internal/
│   ├── domain/                        # Domain Layer (pure business logic)
│   │   ├── user/                      # User aggregate
│   │   │   ├── user.go                # Entity with private fields + getters
│   │   │   ├── repository.go          # Repository interface (port)
│   │   │   ├── errors.go              # Domain-specific errors
│   │   │   └── events.go              # Domain events
│   │   └── shared/                    # Shared domain concepts
│   │       ├── entity.go
│   │       ├── errors.go
│   │       ├── event_bus.go
│   │       └── valueobjects.go
│   │
│   ├── application/                   # Application Layer (CQRS)
│   │   ├── command/                   # Commands (write operations)
│   │   ├── query/                     # Queries (read operations)
│   │   └── dto/                       # Data Transfer Objects
│   │
│   ├── infrastructure/                # Infrastructure Layer (adapters)
│   │   ├── persistence/
│   │   │   ├── postgres/
│   │   │   └── repository/
│   │   ├── messaging/
│   │   └── cache/
│   │
│   └── interfaces/http/               # Interface Adapters Layer
│       ├── handler/
│       ├── middleware/
│       └── router/
│
├── migrations/                        # golang-migrate migrations
└── pkg/                               # Shared packages
```

### Frontend (React)

```
frontend/src/
├── components/
│   ├── ui/                # shadcn/ui components
│   ├── features/          # Feature-specific components
│   └── layout/            # Layout components
├── hooks/                 # Custom hooks
├── lib/                   # Utilities
├── pages/                 # Page components
├── stores/                # Zustand stores
└── types/                 # TypeScript types
```

---

## Design System

### Color Palette (OKLCH)

| Token | Usage |
|-------|-------|
| `primary` | Primary actions, links, focus states (violet hue 290) |
| `secondary` | Secondary actions, accents (cyan hue 220) |
| `destructive` | Error states, destructive actions |
| `muted` | Subtle backgrounds, secondary text |
| `success` | Success states, confirmations |
| `warning` | Warning states, cautions |

### Spacing Scale (4px base)

| Token | Value | Usage |
|-------|-------|-------|
| `1` | 4px | Tight padding |
| `2` | 8px | Inline elements |
| `4` | 16px | Default component padding |
| `6` | 24px | Card spacing |
| `8` | 32px | Section padding |

---

## Code Patterns

### Go Domain Entity

```go
type User struct {
    id        uuid.UUID
    email     Email     // Value Object
    name      string
    createdAt time.Time
}

func (u *User) ID() uuid.UUID { return u.id }
func (u *User) Email() Email  { return u.email }
```

### Go Error Handling

```go
result, err := service.DoSomething(ctx, input)
if err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}
```

### React Component

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
```

---

## Custom Commands (Agents)

Claude Code custom commands are available in `.claude/commands/`. These act as specialized agents for different tasks:

### Core Development

| Command | Description |
|---------|-------------|
| `/project:backend` | Backend engineering with Go + DDD + CQRS |
| `/project:frontend` | Frontend engineering with React + shadcn/ui |
| `/project:fullstack` | End-to-end feature development |
| `/project:api` | Create new API endpoints |
| `/project:migrate` | Database migration operations |
| `/project:test` | Testing patterns and generation |

### Operations & Quality

| Command | Description |
|---------|-------------|
| `/project:review` | Code review checklist |
| `/project:plan` | Technical planning and design |
| `/project:devops` | Infrastructure and CI/CD |
| `/project:security` | Security auditing (OWASP Top 10) |
| `/project:sre` | Site Reliability Engineering |
| `/project:docs` | Documentation generation |

---

## Skills

Skills are reusable templates in `.claude/skills/` that generate code following project patterns:

### Code Generation

| Skill | Usage | Description |
|-------|-------|-------------|
| **Go API Builder** | `/project:skill:go-api <resource>` | Generate complete REST API endpoints with DDD + CQRS |
| **React Component** | `/project:skill:react-component <name>` | Generate React components with design system |
| **Testing Patterns** | `/project:skill:testing <file>` | Generate test suites for Go and React |

### Infrastructure

| Skill | Usage | Description |
|-------|-------|-------------|
| **Database Migration** | `/project:skill:migration <operation>` | Generate golang-migrate migrations |
| **Dockerfile Builder** | `/project:skill:dockerfile <type>` | Generate optimized Dockerfiles |
| **Kubernetes Manifest** | `/project:skill:k8s <resource>` | Generate Kubernetes manifests |

### Skill Files Location

```
.claude/skills/
├── go-api-builder.md          # REST API generation
├── react-component-builder.md # React component generation
├── testing-patterns.md        # Test suite generation
├── database-migration.md      # Migration templates
├── dockerfile-builder.md      # Docker configurations
└── kubernetes-manifest.md     # K8s manifests
```

---

## Commit Message Format

Follow Conventional Commits:

```
<type>(<scope>): <description>
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`, `perf`

Examples:
- `feat(api): add user registration endpoint`
- `fix(ui): resolve button focus state in dark mode`

---

## Related Documentation

### Claude Code

- **Commands (Agents)**: `.claude/commands/` - Specialized command configurations
- **Skills**: `.claude/skills/` - Reusable code generation templates
- **Settings**: `.claude/settings.local.json` - Local permissions and configuration

### GitHub Copilot

- **Agents**: `.github/agents/` - Specialized agent configurations
- **Instructions**: `.github/instructions/` - Context-aware instructions
- **Prompts**: `.github/prompts/` - Reusable prompt templates
- **Skills**: `.github/skills/` - Custom AI skills

---

## References

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [DDD + CQRS in Go](https://threedots.tech/post/ddd-cqrs-clean-architecture-combined/)
- [shadcn/ui](https://ui.shadcn.com/)
- [Tailwind CSS v4](https://tailwindcss.com/docs/v4-beta)
