---
description: Set up a complete CI/CD pipeline with GitHub Actions for testing, building, and deploying
---

# Setup CI/CD Pipeline

Create a comprehensive CI/CD pipeline using GitHub Actions for automated testing, building, and deployment.

## Tech Stack

This project uses:
- **Backend**: Go 1.25+ with Chi router (Clean Architecture + DDD + CQRS)
- **Frontend**: React 19 with pnpm (Next.js standalone output)
- **Database**: PostgreSQL 16
- **Package Manager**: pnpm (frontend)
- **Container Registry**: ghcr.io

## Pipeline Requirements

**Project Type**: {{projectType}}
- [ ] Full-stack (Go Backend + React 19 Frontend)
- [ ] Frontend only (React 19)
- [ ] Backend only (Go)

**Deployment Target**: {{deploymentTarget}}
- [ ] Kubernetes (EKS/GKE/AKS) with Kustomize
- [ ] Docker Compose
- [ ] AWS ECS
- [ ] Serverless (Lambda)

**Environments**:
- [ ] Staging (auto-deploy from develop/main)
- [ ] Production (manual approval required)

## Pipeline Components

### 1. CI Pipeline (ci.yml)

Features:
- [ ] Change detection (only build what changed)
- [ ] Parallel jobs for frontend/backend
- [ ] Linting and type checking
- [ ] Unit tests with coverage
- [ ] Integration tests
- [ ] Docker image build test
- [ ] Security scanning

### 2. CD Pipeline (deploy.yml)

Features:
- [ ] Semantic versioning
- [ ] Docker image build and push
- [ ] Staging auto-deployment
- [ ] Production deployment with approval
- [ ] Rollback capability
- [ ] Slack notifications

### 3. Security Pipeline (security.yml)

Features:
- [ ] Dependency vulnerability scanning (Trivy)
- [ ] Secret detection (Gitleaks)
- [ ] Static code analysis (CodeQL)
- [ ] Container image scanning
- [ ] Infrastructure scanning (Checkov)

## Implementation

### CI Pipeline

Location: `.github/workflows/ci.yml`

```yaml
name: CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main, develop]

concurrency:
  group: ${{ github.workflow }}-${{ github.ref }}
  cancel-in-progress: true

jobs:
  changes:
    runs-on: ubuntu-latest
    outputs:
      frontend: ${{ steps.changes.outputs.frontend }}
      backend: ${{ steps.changes.outputs.backend }}
    steps:
      - uses: actions/checkout@v4
      - uses: dorny/paths-filter@v3
        id: changes
        with:
          filters: |
            frontend:
              - 'frontend/**'
            backend:
              - 'backend/**'

  frontend:
    needs: changes
    if: needs.changes.outputs.frontend == 'true'
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: frontend
    steps:
      - uses: actions/checkout@v4
      - uses: pnpm/action-setup@v4
        with:
          version: 9
      - uses: actions/setup-node@v4
        with:
          node-version: '20'
          cache: 'pnpm'
          cache-dependency-path: frontend/pnpm-lock.yaml
      - name: Install dependencies
        run: pnpm install --frozen-lockfile
      - name: Lint
        run: pnpm lint
      - name: Type check
        run: pnpm type-check
      - name: Run tests
        run: pnpm test:ci
      - name: Build
        run: pnpm build

  backend:
    needs: changes
    if: needs.changes.outputs.backend == 'true'
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:16-alpine
        env:
          POSTGRES_USER: test
          POSTGRES_PASSWORD: test
          POSTGRES_DB: test
        options: --health-cmd pg_isready --health-interval 10s --health-timeout 5s --health-retries 5
        ports:
          - 5432:5432
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.25'
          cache-dependency-path: backend/go.sum
      - name: Install dependencies
        run: go mod download
        working-directory: backend
      - name: Verify dependencies
        run: go mod verify
        working-directory: backend
      - name: Run go vet
        run: go vet ./...
        working-directory: backend
      - name: Run golangci-lint
        uses: golangci/golangci-lint-action@v4
        with:
          working-directory: backend
      - name: Run tests
        run: go test -race -coverprofile=coverage.out -covermode=atomic ./...
        working-directory: backend
        env:
          DATABASE_URL: postgresql://test:test@localhost:5432/test
      - name: Build
        run: go build -o server ./cmd/server
        working-directory: backend

  security:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: aquasecurity/trivy-action@master
        with:
          scan-type: 'fs'
          severity: 'CRITICAL,HIGH'
      - uses: gitleaks/gitleaks-action@v2
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### Deploy Pipeline

Location: `.github/workflows/deploy.yml`

```yaml
name: Deploy

on:
  push:
    branches: [main]
    tags: ['v*']
  workflow_dispatch:
    inputs:
      environment:
        required: true
        type: choice
        options: [staging, production]

jobs:
  build:
    runs-on: ubuntu-latest
    outputs:
      version: ${{ steps.version.outputs.version }}
    steps:
      - uses: actions/checkout@v4
      - id: version
        run: echo "version=$(git describe --tags --always)" >> $GITHUB_OUTPUT
      - uses: docker/setup-buildx-action@v3
      - uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
      - uses: docker/build-push-action@v5
        with:
          context: .
          file: docker/Dockerfile.backend
          push: true
          tags: ghcr.io/${{ github.repository }}/backend:${{ steps.version.outputs.version }}
          cache-from: type=gha
          cache-to: type=gha,mode=max

  deploy-staging:
    needs: build
    runs-on: ubuntu-latest
    environment:
      name: staging
      url: https://staging.example.com
    steps:
      - uses: actions/checkout@v4
      - uses: azure/setup-kubectl@v4
      - run: |
          kubectl apply -k k8s/overlays/staging
          kubectl rollout status deployment/backend -n staging

  deploy-production:
    needs: [build, deploy-staging]
    if: startsWith(github.ref, 'refs/tags/v')
    runs-on: ubuntu-latest
    environment:
      name: production
      url: https://example.com
    steps:
      - uses: actions/checkout@v4
      - uses: azure/setup-kubectl@v4
      - run: |
          kubectl apply -k k8s/overlays/production
          kubectl rollout status deployment/backend -n production
```

## Required Secrets

Configure in GitHub repository settings:

```
# Required
GITHUB_TOKEN              # Auto-provided

# Code coverage
CODECOV_TOKEN             # From codecov.io

# Kubernetes
KUBE_CONFIG_STAGING       # base64 encoded kubeconfig
KUBE_CONFIG_PRODUCTION    # base64 encoded kubeconfig

# Notifications
SLACK_WEBHOOK             # Slack incoming webhook URL
```

## Validation

After implementation:

1. Create a PR to trigger CI
2. Merge to main to trigger staging deploy
3. Create a tag to trigger production deploy

```bash
# Test locally with act
act push -j frontend
act push -j backend

# Create release tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

## Output

Provide:
1. CI workflow file
2. Deploy workflow file
3. Security workflow file
4. Required secrets list
5. Instructions for setting up environments
