# DevOps Documentation

This document provides comprehensive documentation for the DevOps infrastructure, CI/CD pipelines, deployment processes, and operational procedures for the full-stack application.

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Docker](#docker)
- [CI/CD Pipelines](#cicd-pipelines)
- [Kubernetes](#kubernetes)
- [Terraform Infrastructure](#terraform-infrastructure)
- [Monitoring & Observability](#monitoring--observability)
- [Security](#security)
- [Scripts & Automation](#scripts--automation)
- [Environment Configuration](#environment-configuration)
- [Troubleshooting](#troubleshooting)

---

## Overview

The DevOps infrastructure follows modern best practices including:

- **Infrastructure as Code (IaC)** - All infrastructure defined in Terraform
- **GitOps** - Kubernetes manifests managed via Kustomize
- **CI/CD** - Automated pipelines with GitHub Actions
- **Observability** - Full-stack monitoring with Prometheus, Grafana, Loki, and Tempo
- **Security-first** - Automated scanning, secret management, and compliance checks

### Technology Stack

| Component | Technology |
|-----------|------------|
| Container Runtime | Docker |
| Orchestration | Kubernetes (EKS) |
| CI/CD | GitHub Actions |
| IaC | Terraform |
| Monitoring | Prometheus + Grafana |
| Logging | Loki + Promtail |
| Tracing | Tempo (OpenTelemetry) |
| Alerting | Alertmanager |
| Secret Management | Sealed Secrets / AWS Secrets Manager |

---

## Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              CloudFront CDN                                  │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Application Load Balancer                            │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                    ┌─────────────────┼─────────────────┐
                    ▼                 ▼                 ▼
            ┌──────────────┐  ┌──────────────┐  ┌──────────────┐
            │   Frontend   │  │   Backend    │  │    Nginx     │
            │   (Next.js)  │  │  (Node/Go)   │  │   (Proxy)    │
            └──────────────┘  └──────────────┘  └──────────────┘
                                      │
                    ┌─────────────────┼─────────────────┐
                    ▼                 ▼                 ▼
            ┌──────────────┐  ┌──────────────┐  ┌──────────────┐
            │  PostgreSQL  │  │    Redis     │  │      S3      │
            │    (RDS)     │  │(ElastiCache) │  │   (Assets)   │
            └──────────────┘  └──────────────┘  └──────────────┘
```

### Network Architecture

- **VPC**: 10.0.0.0/16 CIDR block
- **Public Subnets**: ALB, NAT Gateway
- **Private Subnets**: Application workloads, databases
- **Intra Subnets**: Internal services only

---

## Docker

### Directory Structure

```
docker/
├── Dockerfile.frontend      # Next.js multi-stage build
├── Dockerfile.backend       # Node.js backend
├── Dockerfile.backend.go    # Go backend alternative
├── Dockerfile.nginx         # Reverse proxy
├── docker-compose.yml       # Development environment
├── docker-compose.prod.yml  # Production environment
├── docker-compose.test.yml  # Testing environment
└── nginx/
    ├── nginx.conf           # Main nginx configuration
    └── conf.d/
        └── default.conf     # Server blocks
```

### Dockerfile Best Practices

All Dockerfiles follow these principles:

1. **Multi-stage builds** - Separate build and runtime stages
2. **Non-root user** - Applications run as unprivileged users
3. **Security hardening** - Minimal base images, no unnecessary packages
4. **Health checks** - Built-in container health verification
5. **Layer optimization** - Efficient caching with proper ordering

### Building Images

```bash
# Build all images
make docker-build

# Build specific image
docker build -f docker/Dockerfile.frontend -t myapp-frontend:latest .
docker build -f docker/Dockerfile.backend -t myapp-backend:latest .

# Build with custom tag
docker build -f docker/Dockerfile.frontend -t myapp-frontend:v1.2.3 .
```

### Docker Compose Environments

#### Development (`docker-compose.yml`)

```bash
# Start all services
docker compose -f docker/docker-compose.yml up -d

# Start specific services
docker compose -f docker/docker-compose.yml up -d postgres redis

# View logs
docker compose -f docker/docker-compose.yml logs -f

# Stop services
docker compose -f docker/docker-compose.yml down
```

#### Production (`docker-compose.prod.yml`)

```bash
# Set environment variables
export REGISTRY=ghcr.io
export IMAGE_NAME=myorg/myapp
export TAG=v1.0.0

# Deploy
docker compose -f docker/docker-compose.prod.yml up -d
```

#### Testing (`docker-compose.test.yml`)

```bash
# Run tests
docker compose -f docker/docker-compose.test.yml up --abort-on-container-exit

# Run with specific profile
docker compose -f docker/docker-compose.test.yml --profile e2e up
```

### Image Registry

Images are pushed to GitHub Container Registry (ghcr.io):

```bash
# Login to registry
echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin

# Push images
make docker-push
```

---

## CI/CD Pipelines

### Pipeline Overview

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│    Commit    │────▶│      CI      │────▶│    Build     │────▶│    Deploy    │
│              │     │  (Test/Lint) │     │   (Docker)   │     │   (K8s)      │
└──────────────┘     └──────────────┘     └──────────────┘     └──────────────┘
        │                   │                    │                    │
        ▼                   ▼                    ▼                    ▼
   PR Created          Tests Pass          Images Built        Deployed to
   Push to main        Security Scan       Push to Registry    Staging/Prod
```

### Workflow Files

| Workflow | Trigger | Purpose |
|----------|---------|---------|
| `ci.yml` | PR, Push to main/develop | Lint, test, build, security scan |
| `deploy.yml` | Push to main, Tags | Deploy to staging/production |
| `security.yml` | Daily, PR | Security scanning (SAST, DAST, dependencies) |
| `release.yml` | Push to main | Semantic versioning, changelog generation |

### CI Pipeline (`ci.yml`)

The CI pipeline runs on every pull request and push:

```yaml
Jobs:
  1. changes        # Detect which files changed
  2. frontend       # Lint, type-check, test, build frontend
  3. backend-node   # Lint, type-check, test backend (Node.js)
  4. backend-go     # Lint, test backend (Go) - if applicable
  5. security       # Trivy scan, Gitleaks
  6. docker-build   # Test Docker builds
  7. e2e            # End-to-end tests (Playwright)
  8. ci-success     # Summary gate
```

#### Running Locally

```bash
# Install act for local GitHub Actions testing
brew install act

# Run CI workflow locally
act push -j frontend
act push -j backend-node
```

### Deploy Pipeline (`deploy.yml`)

Deployment follows blue-green strategy:

```yaml
Stages:
  1. setup          # Determine environment, generate version
  2. build          # Build and push Docker images
  3. deploy-staging # Deploy to staging (auto on develop)
  4. deploy-prod    # Deploy to production (auto on main/tags)
  5. rollback       # Automatic rollback on failure
```

#### Manual Deployment

```bash
# Trigger deployment via GitHub CLI
gh workflow run deploy.yml -f environment=staging

# Deploy specific version
gh workflow run deploy.yml -f environment=production
```

### Security Pipeline (`security.yml`)

Comprehensive security scanning:

| Scan Type | Tool | Purpose |
|-----------|------|---------|
| Dependency | Trivy | Vulnerability scanning |
| Container | Trivy | Image vulnerability scanning |
| SAST | CodeQL | Static code analysis |
| Secrets | Gitleaks, TruffleHog | Secret detection |
| IaC | Checkov, KICS | Infrastructure security |
| DAST | OWASP ZAP | Dynamic security testing |
| License | license-checker | License compliance |

### Release Pipeline (`release.yml`)

Automated semantic versioning:

```bash
# Commit message format
feat: add new feature      # → Minor version bump (1.0.0 → 1.1.0)
fix: resolve bug           # → Patch version bump (1.0.0 → 1.0.1)
feat!: breaking change     # → Major version bump (1.0.0 → 2.0.0)
```

### Required Secrets

Configure these secrets in GitHub repository settings:

```
GITHUB_TOKEN          # Auto-provided
CODECOV_TOKEN         # Code coverage
KUBE_CONFIG_STAGING   # Kubernetes config (base64)
KUBE_CONFIG_PRODUCTION
SLACK_WEBHOOK         # Notifications
PAGERDUTY_SERVICE_KEY # Incident alerts
GITLEAKS_LICENSE      # Gitleaks (optional)
```

---

## Kubernetes

### Directory Structure

```
k8s/
├── base/                           # Base manifests
│   ├── kustomization.yaml
│   ├── namespace.yaml
│   ├── ingress.yaml
│   ├── frontend/
│   │   ├── deployment.yaml
│   │   └── service.yaml
│   └── backend/
│       ├── deployment.yaml
│       └── service.yaml
└── overlays/
    ├── staging/                    # Staging environment
    │   ├── kustomization.yaml
    │   └── namespace.yaml
    └── production/                 # Production environment
        ├── kustomization.yaml
        ├── namespace.yaml
        └── sealed-secrets.yaml
```

### Kustomize Usage

```bash
# Preview manifests
kubectl kustomize k8s/overlays/staging

# Apply to cluster
kubectl apply -k k8s/overlays/staging

# Delete deployment
kubectl delete -k k8s/overlays/staging
```

### Resource Specifications

#### Staging Environment

| Resource | Replicas | CPU Request | Memory Request |
|----------|----------|-------------|----------------|
| Frontend | 1 | 50m | 128Mi |
| Backend | 1 | 100m | 256Mi |

#### Production Environment

| Resource | Replicas | CPU Request | Memory Request |
|----------|----------|-------------|----------------|
| Frontend | 2-10 | 200m | 512Mi |
| Backend | 3-20 | 500m | 1Gi |

### Autoscaling

Horizontal Pod Autoscaler (HPA) is configured for both frontend and backend:

```yaml
Metrics:
  - CPU utilization: 70%
  - Memory utilization: 80%

Scale behavior:
  - Scale up: Immediate (100% or 4 pods per 15s)
  - Scale down: Gradual (10% per 60s, 5min stabilization)
```

### Networking

#### Ingress

```yaml
Hosts:
  - example.com          → frontend
  - www.example.com      → frontend
  - api.example.com      → backend

Features:
  - TLS termination (Let's Encrypt)
  - Rate limiting (100 req/s)
  - Security headers
```

#### Network Policies

```yaml
Policies:
  - default-deny-all     # Deny all traffic by default
  - allow-frontend       # Frontend can access backend
  - allow-backend        # Backend can access DB, Redis
```

### Secrets Management

Using Sealed Secrets for GitOps-compatible secret management:

```bash
# Install kubeseal
brew install kubeseal

# Create sealed secret
kubectl create secret generic backend-secrets \
  --from-literal=DATABASE_URL='...' \
  --from-literal=JWT_SECRET='...' \
  --dry-run=client -o yaml | \
  kubeseal --format yaml > k8s/overlays/production/sealed-secrets.yaml
```

### Useful Commands

```bash
# Get cluster status
make k8s-status

# View logs
make k8s-logs

# Port forward for local access
make k8s-port-forward

# Rollback deployment
kubectl rollout undo deployment/backend -n production

# Scale deployment
kubectl scale deployment/backend --replicas=5 -n production
```

---

## Terraform Infrastructure

### Directory Structure

```
terraform/
├── main.tf                 # Main configuration
├── variables.tf            # Input variables
├── outputs.tf              # Output values
└── environments/
    ├── staging.tfvars      # Staging variables
    └── production.tfvars   # Production variables
```

### Resources Created

| Resource | Service | Purpose |
|----------|---------|---------|
| VPC | AWS VPC | Network isolation |
| EKS | AWS EKS | Kubernetes cluster |
| RDS | AWS RDS | PostgreSQL database |
| ElastiCache | AWS ElastiCache | Redis cache |
| S3 | AWS S3 | Asset storage |
| CloudFront | AWS CloudFront | CDN |
| ALB | AWS ALB | Load balancing |
| Secrets Manager | AWS Secrets Manager | Secret storage |

### Workflow

```bash
# Initialize Terraform
make tf-init

# Plan changes (staging)
make tf-plan-staging

# Plan changes (production)
make tf-plan-prod

# Apply changes
make tf-apply

# View outputs
make tf-output

# Destroy infrastructure (staging only)
make tf-destroy-staging
```

### State Management

Terraform state is stored in S3 with DynamoDB locking:

```hcl
backend "s3" {
  bucket         = "terraform-state-bucket"
  key            = "fullstack-app/terraform.tfstate"
  region         = "us-east-1"
  encrypt        = true
  dynamodb_table = "terraform-locks"
}
```

### Environment Differences

| Setting | Staging | Production |
|---------|---------|------------|
| EKS Nodes | 1-3 (t3.medium) | 3-20 (t3.large) |
| RDS | db.t3.micro, Single-AZ | db.t3.medium, Multi-AZ |
| Redis | cache.t3.micro | cache.t3.medium |
| Backups | 7 days | 30 days |
| Monitoring | Basic | Enhanced |

### Cost Estimation

```bash
# Install Infracost
brew install infracost

# Estimate costs
infracost breakdown --path terraform/ --terraform-var-file=environments/staging.tfvars
```

---

## Monitoring & Observability

### Architecture

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│  Application │────▶│  Prometheus  │────▶│   Grafana    │
│   Metrics    │     │              │     │  Dashboards  │
└──────────────┘     └──────────────┘     └──────────────┘
                            │
                            ▼
                     ┌──────────────┐
                     │ Alertmanager │
                     └──────────────┘
                            │
              ┌─────────────┼─────────────┐
              ▼             ▼             ▼
         ┌────────┐   ┌──────────┐   ┌─────────┐
         │ Slack  │   │PagerDuty │   │  Email  │
         └────────┘   └──────────┘   └─────────┘

┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│  Application │────▶│    Loki      │────▶│   Grafana    │
│     Logs     │     │              │     │   Explore    │
└──────────────┘     └──────────────┘     └──────────────┘
        │
        └───────────▶┌──────────────┐
                     │   Promtail   │
                     └──────────────┘

┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│  Application │────▶│    Tempo     │────▶│   Grafana    │
│    Traces    │     │              │     │   Explore    │
└──────────────┘     └──────────────┘     └──────────────┘
```

### Starting the Stack

```bash
# Start all monitoring services
make monitoring-up

# View logs
make monitoring-logs

# Stop services
make monitoring-down
```

### Access URLs

| Service | URL | Default Credentials |
|---------|-----|---------------------|
| Grafana | http://localhost:3001 | admin / admin |
| Prometheus | http://localhost:9090 | - |
| Alertmanager | http://localhost:9093 | - |
| Loki | http://localhost:3100 | - |
| Tempo | http://localhost:3200 | - |

### Metrics

Key metrics collected:

```
# Application Metrics
http_requests_total              # Request count
http_request_duration_seconds    # Request latency
http_requests_in_flight          # Concurrent requests

# System Metrics
node_cpu_seconds_total           # CPU usage
node_memory_MemAvailable_bytes   # Available memory
node_filesystem_avail_bytes      # Disk space

# Container Metrics
container_cpu_usage_seconds_total
container_memory_usage_bytes
container_network_receive_bytes_total
```

### Alerting Rules

| Alert | Severity | Condition |
|-------|----------|-----------|
| HighErrorRate | Critical | >5% error rate for 5m |
| HighResponseTime | Warning | P95 latency >2s for 5m |
| ServiceDown | Critical | Service unreachable for 1m |
| HighMemoryUsage | Warning | >85% memory for 5m |
| HighCPUUsage | Warning | >85% CPU for 5m |
| DiskSpaceLow | Warning | <15% disk space |
| DiskSpaceCritical | Critical | <5% disk space |

### SLO Monitoring

Service Level Objectives:

| SLO | Target | Window |
|-----|--------|--------|
| Availability | 99.9% | 30 days |
| P99 Latency | <500ms | 30 days |

### Logging

Log format (JSON):

```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "level": "info",
  "message": "Request processed",
  "traceId": "abc123",
  "spanId": "def456",
  "service": "backend",
  "requestId": "req-789"
}
```

Log queries in Grafana:

```logql
# Error logs
{service="backend"} |= "error"

# Slow requests
{service="backend"} | json | request_time > 1s

# Specific trace
{service="backend"} |= "traceId=abc123"
```

### Tracing

OpenTelemetry configuration:

```bash
# Environment variables
OTEL_EXPORTER_OTLP_ENDPOINT=http://tempo:4318
OTEL_SERVICE_NAME=backend
```

---

## Security

### Security Scanning

| Tool | Purpose | Frequency |
|------|---------|-----------|
| Trivy | Dependency & container scanning | Every PR, Daily |
| CodeQL | Static code analysis (SAST) | Every PR |
| Gitleaks | Secret detection | Every commit |
| TruffleHog | Advanced secret detection | Every PR |
| Checkov | IaC security | Every PR |
| KICS | Kubernetes & Docker security | Every PR |
| OWASP ZAP | Dynamic security testing (DAST) | Daily |

### Running Security Scans Locally

```bash
# Full security scan
make security-scan

# Dependency audit
make security-audit

# Trivy scan
trivy fs --severity HIGH,CRITICAL .

# Secret detection
gitleaks detect --source . --verbose

# IaC scan
checkov -d terraform/
```

### Secret Management

**Development:**
- Use `.env` files (never commit)
- Generate secrets: `openssl rand -hex 32`

**Staging/Production:**
- Sealed Secrets for Kubernetes
- AWS Secrets Manager for Terraform
- GitHub Secrets for CI/CD

### Dependabot

Automated dependency updates configured in `.github/dependabot.yml`:

- **Frontend**: Weekly npm updates
- **Backend**: Weekly npm/Go updates
- **GitHub Actions**: Weekly updates
- **Docker**: Weekly base image updates
- **Terraform**: Weekly provider updates

### CODEOWNERS

Code review requirements defined in `.github/CODEOWNERS`:

- Infrastructure changes require DevOps team review
- Security configurations require Security team review
- Database migrations require DBA review

---

## Scripts & Automation

### Available Make Commands

```bash
make help              # Show all commands

# Development
make install           # Install dependencies
make dev               # Start development
make dev-docker        # Start with Docker
make dev-stop          # Stop development

# Build
make build             # Build all
make build-frontend    # Build frontend
make build-backend     # Build backend

# Testing
make test              # Run all tests
make test-e2e          # Run E2E tests
make test-coverage     # Run with coverage

# Code Quality
make lint              # Lint all code
make lint-fix          # Fix lint issues
make format            # Format code
make typecheck         # Type checking

# Docker
make docker-build      # Build images
make docker-push       # Push to registry
make docker-up         # Start production
make docker-down       # Stop containers
make docker-clean      # Clean resources

# Database
make db-migrate        # Run migrations
make db-seed           # Seed data
make db-reset          # Reset database
make db-studio         # Open studio

# Kubernetes
make k8s-apply-staging # Deploy staging
make k8s-apply-prod    # Deploy production
make k8s-status        # Check status
make k8s-logs          # View logs

# Terraform
make tf-init           # Initialize
make tf-plan-staging   # Plan staging
make tf-plan-prod      # Plan production
make tf-apply          # Apply changes

# Monitoring
make monitoring-up     # Start monitoring
make monitoring-down   # Stop monitoring

# Security
make security-scan     # Run scans
make security-audit    # Audit deps
```

### Setup Script

```bash
# Development setup
./scripts/setup.sh --dev

# Production setup
./scripts/setup.sh --prod

# CI setup
./scripts/setup.sh --ci
```

### Deployment Script

```bash
# Deploy to staging
./scripts/deploy.sh staging

# Deploy to production
./scripts/deploy.sh production

# Dry run
./scripts/deploy.sh staging --dry-run
```

### Backup Script

```bash
# Create backup
./scripts/backup.sh

# Restore from backup
./scripts/backup.sh --restore backups/backup_20240115_103000_postgres.sql.gz

# List backups
./scripts/backup.sh --list

# Cleanup old backups
./scripts/backup.sh --cleanup
```

---

## Environment Configuration

### Environment Files

| File | Purpose | Committed |
|------|---------|-----------|
| `.env.example` | Template with all variables | Yes |
| `.env` | Root environment variables | No |
| `frontend/.env.local` | Frontend variables | No |
| `backend/.env` | Backend variables | No |

### Required Variables

**Application:**
```bash
NODE_ENV=production
APP_NAME=fullstack-app
LOG_LEVEL=info
```

**Database:**
```bash
DATABASE_URL=postgresql://user:pass@host:5432/db
DB_POOL_MIN=2
DB_POOL_MAX=10
```

**Redis:**
```bash
REDIS_URL=redis://host:6379
```

**Authentication:**
```bash
JWT_SECRET=<64-char-hex>
JWT_EXPIRES_IN=7d
```

**Monitoring:**
```bash
SENTRY_DSN=https://...
OTEL_EXPORTER_OTLP_ENDPOINT=http://tempo:4318
```

### Generating Secrets

```bash
# JWT Secret
openssl rand -hex 64

# Session Secret
openssl rand -hex 32

# Database Password
openssl rand -base64 32
```

---

## Troubleshooting

### Common Issues

#### Docker Build Fails

```bash
# Clear Docker cache
docker builder prune -af

# Rebuild without cache
docker build --no-cache -f docker/Dockerfile.frontend .
```

#### Kubernetes Pods Not Starting

```bash
# Check pod status
kubectl describe pod <pod-name> -n <namespace>

# Check logs
kubectl logs <pod-name> -n <namespace>

# Check events
kubectl get events -n <namespace> --sort-by='.lastTimestamp'
```

#### Database Connection Issues

```bash
# Test connection
psql $DATABASE_URL -c "SELECT 1"

# Check from pod
kubectl exec -it <pod-name> -- nc -zv postgres 5432
```

#### CI Pipeline Failures

```bash
# Run locally with act
act push -j frontend --verbose

# Check GitHub Actions logs
gh run view <run-id> --log
```

### Health Checks

```bash
# Frontend health
curl http://localhost:3000/api/health

# Backend health
curl http://localhost:8080/health

# Kubernetes readiness
kubectl get pods -o wide
```

### Log Locations

| Component | Location |
|-----------|----------|
| Frontend | stdout / Loki |
| Backend | stdout / Loki |
| Nginx | /var/log/nginx/ |
| PostgreSQL | CloudWatch (RDS) |
| Kubernetes | kubectl logs |

### Support Contacts

- **Infrastructure Issues**: DevOps Team (#devops-support)
- **Security Issues**: Security Team (#security)
- **On-Call**: PagerDuty escalation

---

## References

- [Docker Documentation](https://docs.docker.com/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Terraform AWS Provider](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)
- [GitHub Actions](https://docs.github.com/en/actions)
- [Prometheus](https://prometheus.io/docs/)
- [Grafana](https://grafana.com/docs/)
