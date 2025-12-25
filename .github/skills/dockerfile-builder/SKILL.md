---
name: dockerfile-builder
description: Create production-ready Dockerfiles with multi-stage builds and security best practices
---

# Dockerfile Builder Skill

This skill guides you through creating production-ready Dockerfiles following security, performance, and maintainability best practices.

## Tech Stack

This project uses:
- **Backend**: Go 1.25+ with Chi router (Clean Architecture + DDD + CQRS)
- **Frontend**: React 19 with pnpm (Next.js standalone output)
- **Package Manager**: pnpm (frontend)

## When to Use This Skill

- Creating new Dockerfiles for services
- Optimizing existing Dockerfiles
- Converting development Dockerfiles to production
- Adding multi-stage builds

## Dockerfile Naming Convention

```
docker/
├── Dockerfile.frontend          # React 19 frontend (pnpm)
├── Dockerfile.backend           # Go backend API (primary)
├── Dockerfile.backend.go        # Go backend (legacy name)
├── Dockerfile.nginx             # Nginx reverse proxy
├── Dockerfile.worker            # Background workers
└── Dockerfile.migration         # Database migrations (Goose)
```

## Templates

### Template 1: Go Backend (Primary)

```dockerfile
# =============================================================================
# Go Backend - Multi-stage Production Build
# =============================================================================

# -----------------------------------------------------------------------------
# Stage 1: Builder
# -----------------------------------------------------------------------------
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build arguments
ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT

# Build with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}" \
    -o /app/server ./cmd/server

# -----------------------------------------------------------------------------
# Stage 2: Runner
# -----------------------------------------------------------------------------
FROM alpine:3.20 AS runner
WORKDIR /app

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata wget && \
    update-ca-certificates

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy binary
COPY --from=builder /app/server /app/server

# Copy migrations if needed
COPY --from=builder /app/migrations ./migrations 2>/dev/null || true

RUN chown -R appuser:appgroup /app

USER appuser

ENV PORT=8080

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

ENTRYPOINT ["/app/server"]
```

### Template 2: React 19 Frontend (pnpm)

```dockerfile
# =============================================================================
# React 19 Frontend - Multi-stage Production Build
# =============================================================================

# -----------------------------------------------------------------------------
# Stage 1: Dependencies
# -----------------------------------------------------------------------------
FROM node:20-alpine AS deps
WORKDIR /app

RUN apk add --no-cache libc6-compat
RUN corepack enable pnpm

COPY package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

# -----------------------------------------------------------------------------
# Stage 2: Builder
# -----------------------------------------------------------------------------
FROM node:20-alpine AS builder
WORKDIR /app

RUN corepack enable pnpm

COPY --from=deps /app/node_modules ./node_modules
COPY . .

# Build arguments for environment
ARG NEXT_PUBLIC_API_URL
ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL
ENV NEXT_TELEMETRY_DISABLED=1
ENV NODE_ENV=production

RUN pnpm build

# -----------------------------------------------------------------------------
# Stage 3: Runner
# -----------------------------------------------------------------------------
FROM node:20-alpine AS runner
WORKDIR /app

ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1
ENV PORT=3000

RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs

# Copy built application (standalone output)
COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

EXPOSE 3000

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:3000/api/health || exit 1

CMD ["node", "server.js"]
```

### Template 3: Goose Database Migrations

```dockerfile
# =============================================================================
# Goose Database Migrations
# =============================================================================

FROM golang:1.25-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git

# Install goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM alpine:3.20 AS runner
WORKDIR /app

RUN apk --no-cache add ca-certificates

# Copy goose binary
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Copy migration files
COPY migrations ./migrations

ENV GOOSE_DRIVER=postgres
ENV GOOSE_MIGRATION_DIR=/app/migrations

ENTRYPOINT ["goose"]
CMD ["up"]
```

### Template 4: Python FastAPI

```dockerfile
# =============================================================================
# Python FastAPI - Multi-stage Production Build
# =============================================================================

# -----------------------------------------------------------------------------
# Stage 1: Builder
# -----------------------------------------------------------------------------
FROM python:3.12-slim AS builder
WORKDIR /app

# Install build dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

# Create virtual environment
RUN python -m venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"

# Install dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# -----------------------------------------------------------------------------
# Stage 2: Runner
# -----------------------------------------------------------------------------
FROM python:3.12-slim AS runner
WORKDIR /app

# Create non-root user
RUN groupadd --gid 1001 appgroup && \
    useradd --uid 1001 --gid appgroup --shell /bin/bash appuser

# Copy virtual environment
COPY --from=builder /opt/venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"

# Copy application
COPY --chown=appuser:appgroup . .

USER appuser

ENV PORT=8080
ENV PYTHONUNBUFFERED=1
ENV PYTHONDONTWRITEBYTECODE=1

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
  CMD python -c "import urllib.request; urllib.request.urlopen('http://localhost:8080/health')" || exit 1

CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8080"]
```

### Template 5: Nginx Reverse Proxy

```dockerfile
# =============================================================================
# Nginx Reverse Proxy
# =============================================================================

FROM nginx:1.25-alpine

# Remove default config
RUN rm /etc/nginx/conf.d/default.conf

# Copy custom configuration
COPY docker/nginx/nginx.conf /etc/nginx/nginx.conf
COPY docker/nginx/conf.d/ /etc/nginx/conf.d/

# Create cache directories
RUN mkdir -p /var/cache/nginx/client_temp && \
    chown -R nginx:nginx /var/cache/nginx /var/log/nginx /etc/nginx/conf.d

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost/health || exit 1

EXPOSE 80 443

CMD ["nginx", "-g", "daemon off;"]
```

## .dockerignore Template

```
# Dependencies
node_modules/
vendor/
__pycache__/
*.pyc
.venv/

# Build outputs
dist/
build/
.next/
*.exe

# Development
.git/
.gitignore
.env*
!.env.example
*.log
*.md
!README.md

# IDE
.idea/
.vscode/
*.swp

# Testing
coverage/
.nyc_output/
htmlcov/
.pytest_cache/

# Docker
Dockerfile*
docker-compose*
.docker/
```

## Build Commands

```bash
# Build with default tag
docker build -f docker/Dockerfile.backend -t myapp-backend:latest .

# Build with version tag
docker build -f docker/Dockerfile.backend -t myapp-backend:v1.0.0 .

# Build with build args
docker build -f docker/Dockerfile.backend \
  --build-arg VERSION=v1.0.0 \
  --build-arg BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
  --build-arg GIT_COMMIT=$(git rev-parse HEAD) \
  -t myapp-backend:v1.0.0 .

# Build with cache
docker build -f docker/Dockerfile.backend \
  --cache-from myapp-backend:latest \
  -t myapp-backend:latest .

# Multi-platform build
docker buildx build -f docker/Dockerfile.backend \
  --platform linux/amd64,linux/arm64 \
  -t myapp-backend:latest \
  --push .
```

## Security Checklist

- [ ] Using specific base image version
- [ ] Running as non-root user
- [ ] No secrets in image
- [ ] Health check defined
- [ ] .dockerignore configured
- [ ] Minimal image size
- [ ] Security updates installed
- [ ] Read-only filesystem where possible

## Validation

```bash
# Check image size
docker images myapp-backend:latest

# Verify non-root user
docker run --rm myapp-backend:latest whoami

# Scan for vulnerabilities
trivy image myapp-backend:latest

# Lint Dockerfile
hadolint docker/Dockerfile.backend
```
