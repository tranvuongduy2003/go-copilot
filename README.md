# Full-Stack Application

Production-ready, AI-driven full-stack development environment following **Clean Architecture + Domain-Driven Design (DDD) + CQRS** patterns.

## Tech Stack

| Layer | Technology | Version |
|-------|------------|---------|
| **Backend** | Go | 1.25+ |
| **Router** | Chi | v5 |
| **Database** | PostgreSQL | 16+ |
| **Migrations** | Goose | v3 |
| **Cache** | Redis | 7+ |
| **Frontend** | React | 19 |
| **Styling** | Tailwind CSS | v4 |
| **Components** | shadcn/ui | Latest |
| **State** | TanStack Query + Zustand | Latest |
| **Testing** | Vitest + testify | Latest |

## Project Structure

```
.
├── backend/                    # Go backend (Clean Architecture + DDD + CQRS)
├── frontend/                   # React 19 + TypeScript + Tailwind CSS v4
├── design-system/              # Shared design tokens and components
├── docker/                     # Docker configurations
│   ├── docker-compose.yml      # Development environment
│   ├── docker-compose.prod.yml # Production environment
│   └── docker-compose.test.yml # Testing environment
├── k8s/                        # Kubernetes manifests (Kustomize)
│   ├── base/                   # Base configurations
│   └── overlays/               # Environment-specific overlays
├── terraform/                  # Infrastructure as Code (AWS)
│   ├── main.tf
│   ├── variables.tf
│   ├── outputs.tf
│   └── environments/           # Environment tfvars
├── monitoring/                 # Observability stack
│   ├── prometheus/
│   ├── grafana/
│   ├── loki/
│   ├── tempo/
│   └── alertmanager/
├── scripts/                    # Development scripts
├── docs/                       # Documentation
│   ├── API.md
│   └── CONTRIBUTING.md
└── .github/
    ├── workflows/              # CI/CD pipelines
    ├── agents/                 # AI agent configurations
    ├── instructions/           # Context-aware instructions
    ├── prompts/                # Reusable prompt templates
    └── skills/                 # Custom AI skills
```

## Quick Start

### Prerequisites

- Node.js 20+
- Go 1.25+
- pnpm 9+
- Docker & Docker Compose
- Make

### Environment Setup

```bash
# Check environment
make env-check

# Copy environment template
cp .env.example .env

# Install dependencies
make install
```

### Development

```bash
# Start development environment (DB + Redis + services)
make dev

# Start only frontend
make dev-frontend

# Start only backend
make dev-backend

# Start full Docker environment
make dev-docker

# Stop development environment
make dev-stop
```

### Build

```bash
# Build all services
make build

# Build frontend only
make build-frontend

# Build backend only
make build-backend
```

### Testing

```bash
# Run all tests
make test

# Run frontend tests
make test-frontend

# Run backend tests
make test-backend

# Run E2E tests
make test-e2e

# Run tests with coverage
make test-coverage
```

### Code Quality

```bash
# Lint all code
make lint

# Fix linting issues
make lint-fix

# Format code
make format

# Type checking
make typecheck
```

## Docker

```bash
# Build Docker images
make docker-build

# Push to registry
make docker-push

# Start production environment
make docker-up

# Stop all containers
make docker-down

# Run tests in Docker
make docker-test

# Clean Docker resources
make docker-clean
```

## Database

```bash
# Run migrations
make db-migrate

# Create new migration
make db-migrate-create

# Seed database
make db-seed

# Reset database
make db-reset

# Open database studio
make db-studio
```

### Goose CLI (Direct)

```bash
# Install goose
go install github.com/pressly/goose/v3/cmd/goose@latest

# Create migration
goose -dir backend/migrations/sql create <name> sql

# Apply migrations
goose -dir backend/migrations/sql postgres "$DATABASE_URL" up

# Rollback
goose -dir backend/migrations/sql postgres "$DATABASE_URL" down

# Status
goose -dir backend/migrations/sql postgres "$DATABASE_URL" status
```

## Kubernetes

```bash
# Deploy to staging
make k8s-apply-staging

# Deploy to production
make k8s-apply-prod

# Check status
make k8s-status

# View logs
make k8s-logs

# Port forward services
make k8s-port-forward

# Delete staging
make k8s-delete-staging
```

## Terraform

```bash
# Initialize
make tf-init

# Plan staging
make tf-plan-staging

# Plan production
make tf-plan-prod

# Apply changes
make tf-apply

# Show outputs
make tf-output

# Destroy staging
make tf-destroy-staging
```

## Monitoring

```bash
# Start monitoring stack
make monitoring-up

# Stop monitoring
make monitoring-down

# View logs
make monitoring-logs
```

### Endpoints (when running)

| Service | URL |
|---------|-----|
| Prometheus | http://localhost:9090 |
| Grafana | http://localhost:3001 |
| Alertmanager | http://localhost:9093 |

## Security

```bash
# Run security scans (Trivy + Gitleaks)
make security-scan

# Run dependency audit
make security-audit
```

## AI-Driven Development

This project is optimized for AI coding assistants with comprehensive context and instructions.

### Available Agents

| Agent | File | Description |
|-------|------|-------------|
| Backend Engineer | `.github/agents/backend-engineer.agent.md` | Go + Clean Architecture + DDD + CQRS |
| Frontend Engineer | `.github/agents/frontend-engineer.agent.md` | React 19 + Tailwind CSS v4 + shadcn/ui |
| Fullstack Engineer | `.github/agents/fullstack-engineer.agent.md` | End-to-end feature development |
| Testing Agent | `.github/agents/testing-agent.agent.md` | Comprehensive test suites |
| Code Reviewer | `.github/agents/code-reviewer.agent.md` | Quality + security + design system |
| Security Auditor | `.github/agents/security-auditor.agent.md` | OWASP Top 10 auditing |
| DevOps Engineer | `.github/agents/devops-engineer.agent.md` | Infrastructure + CI/CD |
| SRE Engineer | `.github/agents/sre-engineer.agent.md` | Reliability + observability |
| Planner | `.github/agents/planner.agent.md` | Technical planning + design |
| Documentation | `.github/agents/documentation.agent.md` | Technical writing |

### Instructions (Context-Aware)

Located in `.github/instructions/`:

- `frontend.instructions.md` - React + TypeScript patterns
- `backend.instructions.md` - Go + DDD patterns
- `components.instructions.md` - UI component guidelines
- `testing.instructions.md` - Testing strategies
- `api.instructions.md` - API design
- `database.instructions.md` - Database patterns
- `docker.instructions.md` - Container guidelines
- `kubernetes.instructions.md` - K8s best practices
- `terraform.instructions.md` - IaC patterns
- `cicd.instructions.md` - Pipeline configuration
- `devops.instructions.md` - DevOps practices

### Prompts

Located in `.github/prompts/`:

- `create-component.prompt.md` - Create React components
- `create-api.prompt.md` - Create API endpoints
- `add-tests.prompt.md` - Add test coverage
- `fix-bug.prompt.md` - Debug and fix issues
- `refactor.prompt.md` - Refactor code
- `code-review.prompt.md` - Review changes
- `new-feature.prompt.md` - Plan new features
- `create-dockerfile.prompt.md` - Create Dockerfiles
- `deploy-kubernetes.prompt.md` - K8s deployments
- `setup-cicd.prompt.md` - CI/CD pipelines
- `setup-monitoring.prompt.md` - Observability setup
- `infrastructure-review.prompt.md` - IaC review

### Skills

Located in `.github/skills/`:

- `react-component-builder/` - Component generation
- `shadcn-theming/` - Theme configuration
- `testing-patterns/` - Test patterns

## CI/CD Workflows

| Workflow | Trigger | Description |
|----------|---------|-------------|
| `ci.yml` | Push/PR | Lint, test, build |
| `deploy.yml` | Tag/Manual | Deploy to environments |
| `release.yml` | Main branch | Semantic versioning |
| `security.yml` | Schedule/PR | Security scanning |

## Architecture

### Backend (Clean Architecture + DDD + CQRS)

```
backend/
├── cmd/api/                    # Entry point
├── internal/
│   ├── domain/                 # Domain layer (entities, value objects, repositories interfaces)
│   ├── application/            # Application layer (commands, queries, DTOs)
│   ├── infrastructure/         # Infrastructure layer (DB, cache implementations)
│   └── interfaces/             # Interface adapters (HTTP handlers, middleware)
├── migrations/sql/             # Goose SQL migrations
└── pkg/                        # Shared packages
```

### Frontend

```
frontend/
├── src/
│   ├── components/
│   │   ├── ui/                 # shadcn/ui components
│   │   └── features/           # Feature components
│   ├── hooks/                  # Custom hooks
│   ├── lib/                    # Utilities
│   ├── pages/                  # Page components
│   ├── stores/                 # Zustand stores
│   └── types/                  # TypeScript types
```

## Design System

### Colors (OKLCH)

| Token | Usage |
|-------|-------|
| `primary` | Primary actions, links |
| `secondary` | Secondary actions |
| `destructive` | Error states |
| `muted` | Subtle backgrounds |
| `success` | Success states |
| `warning` | Warning states |

### Spacing (4px base)

| Token | Value | Usage |
|-------|-------|-------|
| `1` | 4px | Tight padding |
| `2` | 8px | Inline elements |
| `4` | 16px | Default padding |
| `6` | 24px | Card spacing |
| `8` | 32px | Section padding |

## Commit Convention

Follow [Conventional Commits](https://conventionalcommits.org/):

```
<type>(<scope>): <description>

[body]

[footer]
```

**Types:** `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`, `perf`

**Examples:**
- `feat(api): add user registration endpoint`
- `fix(ui): resolve button focus state`
- `docs(readme): update installation guide`

## Environment Variables

See [.env.example](.env.example) for all configuration options:

- Application settings
- Database configuration
- Redis configuration
- Authentication & Security
- External services (AWS, OAuth)
- Email configuration
- Monitoring & Observability
- Feature flags
- Third-party APIs

## Useful Commands

```bash
# Show all available commands
make help

# Show current version
make version

# Generate TypeScript types from backend
make generate-types

# Generate documentation
make docs

# Clean all build artifacts
make clean

# Clean everything including Docker
make clean-all
```

## Documentation

- [API Documentation](docs/API.md)
- [Contributing Guide](docs/CONTRIBUTING.md)
- [AI Agent Instructions](AGENTS.md)

## License

MIT
