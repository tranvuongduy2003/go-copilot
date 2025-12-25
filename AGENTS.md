# Copilot Coding Agent Instructions

This document provides root-level instructions for the GitHub Copilot Coding Agent working on this project.

## Project Overview

Full-stack application with:
- **Backend**: Go 1.25 with clean architecture
- **Frontend**: React 19 + Tailwind CSS v4 + shadcn/ui
- **Database**: PostgreSQL 16
- **Infrastructure**: Docker, Docker Compose

## Quick Start

### Prerequisites

- Go 1.25+
- Node.js 20+
- Docker and Docker Compose
- PostgreSQL 16 (or use Docker)

### Development Setup

```bash
# Clone and setup
git clone <repository-url>
cd copilot

# Start infrastructure
docker compose up -d postgres redis

# Backend setup
cd backend
cp .env.example .env
go mod download
go run cmd/api/main.go

# Frontend setup (new terminal)
cd frontend
npm install
npm run dev
```

### Environment Variables

Backend (`backend/.env`):
```env
DATABASE_URL=postgres://postgres:postgres@localhost:5432/app?sslmode=disable
REDIS_URL=redis://localhost:6379
JWT_SECRET=your-secret-key
PORT=8080
ENVIRONMENT=development
```

Frontend (`frontend/.env`):
```env
VITE_API_URL=http://localhost:8080/api/v1
```

---

## Build Commands

### Backend

```bash
# Build
cd backend && go build -o bin/api cmd/api/main.go

# Run
./bin/api

# Run with hot reload (using air)
air

# Build Docker image
docker build -t app-backend -f docker/Dockerfile.backend .
```

### Frontend

```bash
# Development
cd frontend && npm run dev

# Build for production
npm run build

# Preview production build
npm run preview

# Build Docker image
docker build -t app-frontend -f docker/Dockerfile.frontend .
```

### Full Stack

```bash
# Start all services
docker compose up -d

# View logs
docker compose logs -f

# Stop all services
docker compose down

# Rebuild and restart
docker compose up -d --build
```

---

## Test Commands

### Backend Tests

```bash
cd backend

# Run all tests
go test ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package tests
go test ./internal/service/...

# Run with verbose output
go test -v ./...

# Run integration tests
go test -tags=integration ./...
```

### Frontend Tests

```bash
cd frontend

# Run all tests
npm test

# Run with coverage
npm run test:coverage

# Run in watch mode
npm run test:watch

# Run e2e tests
npm run test:e2e
```

---

## Development Workflow

### Starting New Feature

1. Create a new branch from `main`
2. Implement the feature following architecture patterns
3. Write tests (aim for 80%+ coverage)
4. Update documentation if needed
5. Create a pull request
6. Address code review feedback
7. Merge after approval

### Database Migrations

```bash
cd backend

# Create new migration
migrate create -ext sql -dir migrations -seq <migration_name>

# Run migrations
migrate -path migrations -database "$DATABASE_URL" up

# Rollback last migration
migrate -path migrations -database "$DATABASE_URL" down 1

# Check migration status
migrate -path migrations -database "$DATABASE_URL" version
```

### Code Generation

```bash
# Generate mocks (backend)
go generate ./...

# Generate API types from OpenAPI spec
npm run generate:api

# Generate shadcn/ui component
npx shadcn@latest add <component-name>
```

---

## Branch Naming Conventions

Use the following branch naming format:

```
<type>/<ticket-id>-<short-description>
```

### Types

| Type | Usage |
|------|-------|
| `feature/` | New features |
| `fix/` | Bug fixes |
| `hotfix/` | Critical production fixes |
| `refactor/` | Code refactoring |
| `docs/` | Documentation changes |
| `test/` | Test additions or fixes |
| `chore/` | Maintenance tasks |

### Examples

- `feature/PROJ-123-user-authentication`
- `fix/PROJ-456-login-redirect-loop`
- `refactor/PROJ-789-simplify-error-handling`
- `docs/update-api-documentation`

---

## Commit Message Format

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, semicolons, etc.)
- `refactor`: Code refactoring (no feature or fix)
- `perf`: Performance improvements
- `test`: Adding or updating tests
- `chore`: Maintenance tasks
- `ci`: CI/CD changes
- `build`: Build system changes

### Scopes

- `api`: Backend API changes
- `ui`: Frontend UI changes
- `auth`: Authentication related
- `db`: Database related
- `config`: Configuration changes
- `deps`: Dependency updates

### Examples

```
feat(api): add user registration endpoint

Implement POST /api/v1/users endpoint with:
- Email validation
- Password hashing
- Duplicate email check

Closes #123
```

```
fix(ui): resolve button focus state in dark mode

The primary button was not showing focus ring in dark mode
due to incorrect color contrast calculation.
```

---

## Pull Request Requirements

### Before Creating PR

- [ ] Code follows project coding standards
- [ ] All tests pass locally
- [ ] New code has appropriate test coverage (80%+)
- [ ] Documentation updated if needed
- [ ] No linting errors or warnings
- [ ] Commit messages follow conventional format
- [ ] Branch is up to date with main

### PR Description Template

```markdown
## Summary
Brief description of changes

## Changes
- Change 1
- Change 2

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
Describe testing performed

## Screenshots (if applicable)
Add screenshots for UI changes

## Checklist
- [ ] Self-review completed
- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] No breaking changes (or documented)
```

### Review Process

1. Automated checks must pass (CI, tests, linting)
2. At least 1 approval required
3. All comments must be resolved
4. Branch must be up to date with main
5. Squash merge to main

---

## Architecture Guidelines

### Backend Structure (Clean Architecture)

```
backend/
├── cmd/api/           # Application entry point
├── internal/          # Private application code
│   ├── config/        # Configuration
│   ├── domain/        # Domain models and interfaces
│   ├── handlers/      # HTTP handlers
│   ├── middleware/    # HTTP middleware
│   ├── repository/    # Data access layer
│   └── service/       # Business logic
├── migrations/        # Database migrations
└── pkg/              # Public packages
```

### Frontend Structure

```
frontend/src/
├── components/        # UI components
│   ├── ui/           # shadcn/ui components
│   └── features/     # Feature-specific components
├── hooks/            # Custom React hooks
├── lib/              # Utilities and helpers
├── pages/            # Page components
├── stores/           # State management
├── styles/           # Global styles
└── types/            # TypeScript types
```

---

## Troubleshooting

### Common Issues

**Database connection failed**
```bash
# Check if PostgreSQL is running
docker compose ps

# Check connection
psql "$DATABASE_URL" -c "SELECT 1"
```

**Frontend build errors**
```bash
# Clear cache and reinstall
rm -rf node_modules .vite
npm install
```

**Go module issues**
```bash
# Clear module cache
go clean -modcache
go mod download
```

**Port already in use**
```bash
# Find and kill process
lsof -i :8080
kill -9 <PID>
```

---

## Useful Commands

```bash
# Format Go code
gofmt -w .

# Lint Go code
golangci-lint run

# Format frontend code
npm run format

# Lint frontend code
npm run lint

# Check types
npm run typecheck

# Database shell
docker compose exec postgres psql -U postgres -d app
```
