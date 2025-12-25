---
applyTo: "docker/**/*,k8s/**/*,terraform/**/*,monitoring/**/*,.github/workflows/**/*,Makefile,scripts/**/*"
---

# DevOps Development Instructions

These instructions apply to all DevOps-related files including Docker, Kubernetes, Terraform, monitoring configurations, CI/CD workflows, and automation scripts.

## Tech Stack

- **Backend**: Go 1.25+ with Chi router, PostgreSQL 16, Goose migrations
- **Frontend**: React 19, TypeScript 5, Tailwind CSS v4, shadcn/ui
- **Architecture**: Clean Architecture + DDD + CQRS
- **Package Manager**: pnpm (frontend)
- **Container Registry**: ghcr.io
- **Orchestration**: Kubernetes with Kustomize
- **IaC**: Terraform (AWS)
- **Monitoring**: Prometheus, Grafana, Loki, Tempo

## Project DevOps Structure

```
project/
├── backend/                         # Go backend service
│   ├── cmd/server/                  # Application entrypoint
│   ├── internal/                    # Private packages (DDD layers)
│   │   ├── domain/                  # Domain entities, value objects
│   │   ├── application/             # Use cases, CQRS handlers
│   │   ├── infrastructure/          # Database, external services
│   │   └── interfaces/              # HTTP handlers, middleware
│   ├── migrations/                  # Goose SQL migrations
│   ├── go.mod
│   └── go.sum
│
├── frontend/                        # React 19 frontend
│   └── src/
│       ├── components/              # React components
│       ├── features/                # Feature modules
│       ├── hooks/                   # Custom hooks
│       └── lib/                     # Utilities
│
├── docker/                          # Container configurations
│   ├── Dockerfile.frontend          # React multi-stage build (pnpm)
│   ├── Dockerfile.backend           # Go multi-stage build
│   ├── Dockerfile.backend.go        # Go backend (legacy name)
│   ├── Dockerfile.nginx             # Nginx reverse proxy
│   ├── docker-compose.yml           # Development environment
│   ├── docker-compose.prod.yml      # Production environment
│   ├── docker-compose.test.yml      # Testing environment
│   └── nginx/                       # Nginx configurations
│
├── k8s/                             # Kubernetes manifests (Kustomize)
│   ├── base/                        # Base configurations
│   │   ├── frontend/
│   │   ├── backend/
│   │   ├── ingress.yaml
│   │   └── kustomization.yaml
│   └── overlays/                    # Environment-specific overlays
│       ├── staging/
│       └── production/
│
├── terraform/                       # Infrastructure as Code (AWS)
│   ├── main.tf                      # VPC, EKS, RDS, ElastiCache
│   ├── variables.tf
│   ├── outputs.tf
│   └── environments/
│       ├── staging.tfvars
│       └── production.tfvars
│
├── monitoring/                      # Observability stack
│   ├── prometheus/                  # Metrics collection
│   ├── grafana/                     # Visualization
│   ├── alertmanager/                # Alert routing
│   ├── loki/                        # Log aggregation
│   └── tempo/                       # Distributed tracing
│
├── .github/workflows/               # CI/CD pipelines
│   ├── ci.yml
│   ├── deploy.yml
│   ├── security.yml
│   └── release.yml
│
├── scripts/                         # Automation scripts
│   ├── setup.sh
│   ├── deploy.sh
│   └── backup.sh
│
├── design-system/                   # Design tokens (colors, spacing)
│   └── tokens/
│
├── docs/                            # Documentation
│   ├── ARCHITECTURE.md
│   └── DEVOPS.md
│
└── Makefile                         # 80+ command shortcuts
```

## Core Principles

### 1. Infrastructure as Code

- All infrastructure must be defined in code
- Changes go through version control
- Use modules for reusability
- Document all resources

### 2. GitOps

- Git is the source of truth
- Changes deployed via pull requests
- Automated synchronization
- Audit trail for all changes

### 3. Security First

- Scan for vulnerabilities in CI
- Use non-root containers
- Implement least privilege
- Encrypt secrets at rest

### 4. Observability

- Implement the three pillars: metrics, logs, traces
- Create actionable alerts
- Build meaningful dashboards
- Use structured logging

## Naming Conventions

### Docker

```
Dockerfile.<service>           # e.g., Dockerfile.frontend
docker-compose.<env>.yml       # e.g., docker-compose.prod.yml
```

### Kubernetes

```
k8s/base/<resource>.yaml       # Base manifests
k8s/overlays/<env>/            # Environment overlays
```

### Terraform

```
terraform/main.tf              # Main configuration
terraform/variables.tf         # Input variables
terraform/outputs.tf           # Output values
terraform/modules/<name>/      # Reusable modules
```

### GitHub Actions

```
.github/workflows/ci.yml       # Continuous Integration
.github/workflows/deploy.yml   # Continuous Deployment
.github/workflows/security.yml # Security scanning
```

## Version Tagging

Use semantic versioning for releases:

```
v1.0.0        # Major.Minor.Patch
v1.0.0-rc.1   # Release candidate
v1.0.0-beta.1 # Beta release
```

## Environment Variables

### Naming Pattern

```bash
# Application
APP_NAME=fullstack-app
APP_ENV=production

# Database
DATABASE_URL=postgresql://...
DB_HOST=localhost
DB_PORT=5432

# Cache
REDIS_URL=redis://...
REDIS_HOST=localhost

# Secrets (use external secret managers in production)
JWT_SECRET=...
API_KEY=...
```

### Never Commit

- `.env` files with real values
- Passwords or API keys
- TLS certificates and private keys
- Kubeconfig files

## Common Commands

### Development

```bash
make dev              # Start development environment
make dev-stop         # Stop development environment
make dev-logs         # View logs
```

### Build

```bash
make build            # Build all services
make docker-build     # Build Docker images
make docker-push      # Push to registry
```

### Testing

```bash
make test             # Run all tests
make test-ci          # Run tests in CI mode
make docker-test      # Run tests in Docker
```

### Deployment

```bash
make k8s-apply-staging    # Deploy to staging
make k8s-apply-prod       # Deploy to production
make k8s-status           # Check deployment status
```

### Monitoring

```bash
make monitoring-up    # Start monitoring stack
make monitoring-down  # Stop monitoring stack
```

### Security

```bash
make security-scan    # Run security scans
make security-audit   # Audit dependencies
```

## Best Practices Checklist

### Docker

- [ ] Multi-stage builds for smaller images
- [ ] Non-root user in final stage
- [ ] Health checks defined
- [ ] .dockerignore configured
- [ ] No secrets in images

### Kubernetes

- [ ] Resource requests and limits set
- [ ] Liveness and readiness probes
- [ ] Pod disruption budgets
- [ ] Network policies defined
- [ ] RBAC configured

### Terraform

- [ ] Remote state with locking
- [ ] Separate environments
- [ ] Modules for reusability
- [ ] Variables validated
- [ ] Outputs documented

### CI/CD

- [ ] Tests run before deploy
- [ ] Security scanning enabled
- [ ] Deployment requires approval
- [ ] Rollback capability
- [ ] Notifications configured

### Monitoring

- [ ] Key metrics collected
- [ ] Alerts have runbooks
- [ ] Dashboards created
- [ ] Log aggregation enabled
- [ ] Tracing implemented
