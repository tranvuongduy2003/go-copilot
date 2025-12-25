---
name: github-actions
description: Create GitHub Actions workflows for CI/CD, security scanning, and automation
---

# GitHub Actions Skill

This skill guides you through creating GitHub Actions workflows for CI/CD pipelines, security scanning, and automation.

## Tech Stack

This project uses:
- **Backend**: Go 1.25+ with Chi router (Clean Architecture + DDD + CQRS)
- **Frontend**: React 19 with pnpm (Next.js standalone output)
- **Database**: PostgreSQL 16
- **Package Manager**: pnpm (frontend)
- **Container Registry**: ghcr.io

## When to Use This Skill

- Setting up CI pipelines for testing and building
- Creating CD pipelines for deployment
- Implementing security scanning workflows
- Automating release processes
- Creating reusable workflows

## Workflow Structure

```
.github/
└── workflows/
    ├── ci.yml                   # Continuous Integration
    ├── deploy.yml               # Continuous Deployment
    ├── security.yml             # Security Scanning
    ├── release.yml              # Release Automation
    └── reusable-build.yml       # Reusable Workflows
```

## Templates

### Template 1: CI Pipeline

```yaml
# .github/workflows/ci.yml
name: CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main, develop]
  workflow_dispatch:

concurrency:
  group: ${{ github.workflow }}-${{ github.ref }}
  cancel-in-progress: true

env:
  NODE_VERSION: '20'
  PNPM_VERSION: '8'

jobs:
  # Detect what changed
  changes:
    name: Detect Changes
    runs-on: ubuntu-latest
    outputs:
      frontend: ${{ steps.changes.outputs.frontend }}
      backend: ${{ steps.changes.outputs.backend }}
      docker: ${{ steps.changes.outputs.docker }}
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
            docker:
              - 'docker/**'

  # Frontend CI
  frontend:
    name: Frontend
    needs: changes
    if: needs.changes.outputs.frontend == 'true'
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: frontend
    steps:
      - uses: actions/checkout@v4

      - uses: pnpm/action-setup@v2
        with:
          version: ${{ env.PNPM_VERSION }}

      - uses: actions/setup-node@v4
        with:
          node-version: ${{ env.NODE_VERSION }}
          cache: 'pnpm'
          cache-dependency-path: frontend/pnpm-lock.yaml

      - name: Install dependencies
        run: pnpm install --frozen-lockfile

      - name: Lint
        run: pnpm lint

      - name: Type check
        run: pnpm type-check

      - name: Test
        run: pnpm test:ci

      - name: Build
        run: pnpm build

      - uses: codecov/codecov-action@v4
        with:
          token: ${{ secrets.CODECOV_TOKEN }}
          flags: frontend

  # Backend CI (Go)
  backend:
    name: Backend
    needs: changes
    if: needs.changes.outputs.backend == 'true'
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: backend
    services:
      postgres:
        image: postgres:16-alpine
        env:
          POSTGRES_USER: test
          POSTGRES_PASSWORD: test
          POSTGRES_DB: test
        ports:
          - 5432:5432
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
      redis:
        image: redis:7-alpine
        ports:
          - 6379:6379
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '1.25'
          cache-dependency-path: backend/go.sum

      - name: Install dependencies
        run: go mod download

      - name: Verify dependencies
        run: go mod verify

      - name: Run go vet
        run: go vet ./...

      - name: Run golangci-lint
        uses: golangci/golangci-lint-action@v4
        with:
          working-directory: backend

      - name: Test
        run: go test -race -coverprofile=coverage.out -covermode=atomic ./...
        env:
          DATABASE_URL: postgresql://test:test@localhost:5432/test
          REDIS_URL: redis://localhost:6379

      - name: Build
        run: go build -o server ./cmd/server

  # Docker Build Test
  docker:
    name: Docker Build
    needs: [changes]
    if: needs.changes.outputs.docker == 'true' || needs.changes.outputs.frontend == 'true' || needs.changes.outputs.backend == 'true'
    runs-on: ubuntu-latest
    strategy:
      matrix:
        include:
          - image: frontend
            dockerfile: docker/Dockerfile.frontend
          - image: backend
            dockerfile: docker/Dockerfile.backend
    steps:
      - uses: actions/checkout@v4

      - uses: docker/setup-buildx-action@v3

      - name: Build
        uses: docker/build-push-action@v5
        with:
          context: .
          file: ${{ matrix.dockerfile }}
          push: false
          tags: ${{ matrix.image }}:test
          cache-from: type=gha
          cache-to: type=gha,mode=max

  # Security Scan
  security:
    name: Security
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

  # CI Summary
  ci-success:
    name: CI Success
    needs: [frontend, backend, docker, security]
    if: always()
    runs-on: ubuntu-latest
    steps:
      - name: Check results
        run: |
          if [[ "${{ needs.frontend.result }}" == "failure" ]] || \
             [[ "${{ needs.backend.result }}" == "failure" ]] || \
             [[ "${{ needs.security.result }}" == "failure" ]]; then
            exit 1
          fi
```

### Template 2: Deploy Pipeline

```yaml
# .github/workflows/deploy.yml
name: Deploy

on:
  push:
    branches: [main]
    tags: ['v*']
  workflow_dispatch:
    inputs:
      environment:
        description: 'Target environment'
        required: true
        type: choice
        options:
          - staging
          - production

concurrency:
  group: deploy-${{ github.ref }}
  cancel-in-progress: false

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  # Setup
  setup:
    name: Setup
    runs-on: ubuntu-latest
    outputs:
      environment: ${{ steps.env.outputs.environment }}
      version: ${{ steps.version.outputs.version }}
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Determine environment
        id: env
        run: |
          if [[ "${{ github.event_name }}" == "workflow_dispatch" ]]; then
            echo "environment=${{ inputs.environment }}" >> $GITHUB_OUTPUT
          elif [[ "${{ github.ref }}" == refs/tags/v* ]]; then
            echo "environment=production" >> $GITHUB_OUTPUT
          else
            echo "environment=staging" >> $GITHUB_OUTPUT
          fi

      - name: Generate version
        id: version
        run: |
          if [[ "${{ github.ref }}" == refs/tags/v* ]]; then
            echo "version=${GITHUB_REF#refs/tags/v}" >> $GITHUB_OUTPUT
          else
            echo "version=$(git describe --tags --always)-$(date +%Y%m%d%H%M%S)" >> $GITHUB_OUTPUT
          fi

  # Build Images
  build:
    name: Build ${{ matrix.image }}
    needs: setup
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write
    strategy:
      matrix:
        include:
          - image: frontend
            dockerfile: docker/Dockerfile.frontend
          - image: backend
            dockerfile: docker/Dockerfile.backend
    steps:
      - uses: actions/checkout@v4

      - uses: docker/setup-buildx-action@v3

      - uses: docker/login-action@v3
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - uses: docker/metadata-action@v5
        id: meta
        with:
          images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}/${{ matrix.image }}
          tags: |
            type=raw,value=${{ needs.setup.outputs.version }}
            type=raw,value=latest,enable=${{ github.ref == 'refs/heads/main' }}

      - uses: docker/build-push-action@v5
        with:
          context: .
          file: ${{ matrix.dockerfile }}
          push: true
          tags: ${{ steps.meta.outputs.tags }}
          cache-from: type=gha
          cache-to: type=gha,mode=max
          build-args: |
            VERSION=${{ needs.setup.outputs.version }}

  # Deploy Staging
  deploy-staging:
    name: Deploy Staging
    needs: [setup, build]
    if: needs.setup.outputs.environment == 'staging'
    runs-on: ubuntu-latest
    environment:
      name: staging
      url: https://staging.example.com
    steps:
      - uses: actions/checkout@v4

      - uses: azure/setup-kubectl@v4

      - name: Deploy
        run: |
          echo "${{ secrets.KUBE_CONFIG_STAGING }}" | base64 -d > kubeconfig
          export KUBECONFIG=kubeconfig
          kubectl apply -k k8s/overlays/staging
          kubectl rollout status deployment/backend -n staging --timeout=300s

  # Deploy Production
  deploy-production:
    name: Deploy Production
    needs: [setup, build]
    if: needs.setup.outputs.environment == 'production'
    runs-on: ubuntu-latest
    environment:
      name: production
      url: https://example.com
    steps:
      - uses: actions/checkout@v4

      - uses: azure/setup-kubectl@v4

      - name: Deploy
        run: |
          echo "${{ secrets.KUBE_CONFIG_PRODUCTION }}" | base64 -d > kubeconfig
          export KUBECONFIG=kubeconfig
          kubectl apply -k k8s/overlays/production
          kubectl rollout status deployment/backend -n production --timeout=600s

      - name: Notify
        if: always()
        uses: 8398a7/action-slack@v3
        with:
          status: ${{ job.status }}
          fields: repo,message,commit,author
        env:
          SLACK_WEBHOOK_URL: ${{ secrets.SLACK_WEBHOOK }}
```

### Template 3: Security Scanning

```yaml
# .github/workflows/security.yml
name: Security

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]
  schedule:
    - cron: '0 2 * * *'

permissions:
  contents: read
  security-events: write

jobs:
  trivy:
    name: Trivy Scan
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: aquasecurity/trivy-action@master
        with:
          scan-type: 'fs'
          format: 'sarif'
          output: 'trivy-results.sarif'
          severity: 'CRITICAL,HIGH,MEDIUM'

      - uses: github/codeql-action/upload-sarif@v3
        with:
          sarif_file: 'trivy-results.sarif'

  gitleaks:
    name: Secret Detection
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - uses: gitleaks/gitleaks-action@v2
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

  codeql:
    name: CodeQL Analysis
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: github/codeql-action/init@v3
        with:
          languages: javascript

      - uses: github/codeql-action/autobuild@v3

      - uses: github/codeql-action/analyze@v3
```

### Template 4: Release Automation

```yaml
# .github/workflows/release.yml
name: Release

on:
  push:
    branches: [main]

permissions:
  contents: write
  packages: write
  pull-requests: write

jobs:
  release:
    name: Release
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - uses: actions/setup-node@v4
        with:
          node-version: '20'

      - name: Install semantic-release
        run: npm install -g semantic-release @semantic-release/changelog @semantic-release/git

      - name: Release
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        run: npx semantic-release
```

## Best Practices

1. **Use concurrency control** to prevent duplicate runs
2. **Cache dependencies** for faster builds
3. **Use path filters** to skip unnecessary jobs
4. **Pin action versions** to specific SHA or version
5. **Use environments** for deployment approvals
6. **Store secrets securely** in GitHub Secrets
7. **Use matrix builds** for multiple configurations
8. **Add status checks** as branch protection rules

## Checklist

- [ ] Concurrency control configured
- [ ] Dependencies cached
- [ ] Path filters used
- [ ] Action versions pinned
- [ ] Environments configured
- [ ] Secrets not hardcoded
- [ ] Security scanning enabled
- [ ] Notifications configured
