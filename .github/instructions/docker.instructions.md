---
applyTo: "docker/**/*,Dockerfile*"
---

# Docker Development Instructions

These instructions apply to all Docker-related files including Dockerfiles and docker-compose configurations.

## Tech Stack

- **Backend**: Go 1.25+ (compiled binary, ~10-20MB image)
- **Frontend**: React 19 with pnpm (Next.js standalone output)
- **Package Manager**: pnpm (frontend)
- **Base Images**: Alpine Linux for minimal size

## Project Docker Structure

```
docker/
├── Dockerfile.frontend          # React/Next.js multi-stage build (pnpm)
├── Dockerfile.backend           # Go multi-stage build
├── Dockerfile.backend.go        # Go backend (legacy name)
├── Dockerfile.nginx             # Nginx reverse proxy
├── docker-compose.yml           # Development environment
├── docker-compose.prod.yml      # Production environment
├── docker-compose.test.yml      # Testing environment
└── nginx/
    ├── nginx.conf               # Main nginx config
    └── conf.d/
        └── default.conf         # Server blocks
```

## Dockerfile Best Practices

### 1. Go Backend Multi-stage Build

```dockerfile
# Stage 1: Builder
FROM golang:1.25-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=${VERSION}" \
    -o /app/server ./cmd/server

# Stage 2: Runner (minimal ~10-20MB)
FROM alpine:3.20 AS runner
WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata wget
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder /app/server /app/server
COPY --from=builder /app/migrations ./migrations 2>/dev/null || true

USER appuser
EXPOSE 8080
ENTRYPOINT ["/app/server"]
```

### 2. Run as Non-root User

```dockerfile
# Alpine-style user creation (for Go/Alpine images)
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Node-style user creation (for frontend images)
RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs

# Set ownership
COPY --chown=appuser:appgroup --from=builder /app/server ./server

# Switch to non-root user
USER appuser
```

### 3. Add Health Checks

```dockerfile
HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1
```

### 4. Use Specific Base Image Tags

```dockerfile
# Good - specific version
FROM node:20.10-alpine3.19

# Avoid - floating tags
# FROM node:latest
# FROM node:20
```

### 5. Optimize Layer Caching

```dockerfile
# Copy package files first (changes less frequently)
COPY package*.json ./
RUN npm ci

# Then copy source (changes more frequently)
COPY . .
RUN npm run build
```

### 6. Minimize Image Size

```dockerfile
# Use alpine variants
FROM node:20-alpine

# Remove unnecessary files
RUN npm ci --only=production && \
    npm cache clean --force

# Use .dockerignore
```

## Docker Compose Patterns

### Development Environment

```yaml
version: "3.9"

services:
  app:
    build:
      context: ..
      dockerfile: docker/Dockerfile.backend
      target: deps  # Stop at deps stage for dev
    volumes:
      - ../backend:/app
      - /app/node_modules  # Preserve node_modules
    ports:
      - "8080:8080"
    environment:
      - NODE_ENV=development
    command: npm run dev
    depends_on:
      postgres:
        condition: service_healthy

  postgres:
    image: postgres:16-alpine
    environment:
      - POSTGRES_PASSWORD=postgres
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5
```

### Production Environment

```yaml
version: "3.9"

services:
  app:
    image: ${REGISTRY}/${IMAGE}:${TAG}
    restart: always
    deploy:
      replicas: 2
      resources:
        limits:
          cpus: '1'
          memory: 512M
        reservations:
          cpus: '0.25'
          memory: 256M
    healthcheck:
      test: ["CMD", "wget", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

## Service-specific Patterns

### Backend (Go) - Primary

```dockerfile
# =============================================================================
# Go Backend - Multi-stage Production Build
# =============================================================================
FROM golang:1.25-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}" \
    -o /app/server ./cmd/server

FROM alpine:3.20 AS runner
WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata wget && \
    update-ca-certificates

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder /app/server /app/server
COPY --from=builder /app/migrations ./migrations 2>/dev/null || true

RUN chown -R appuser:appgroup /app

USER appuser

ENV PORT=8080

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

ENTRYPOINT ["/app/server"]
```

### Frontend (React 19 with pnpm)

```dockerfile
# =============================================================================
# React 19 Frontend - Multi-stage Production Build
# =============================================================================
FROM node:20-alpine AS deps
WORKDIR /app

RUN apk add --no-cache libc6-compat
RUN corepack enable pnpm

COPY package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

FROM node:20-alpine AS builder
WORKDIR /app

RUN corepack enable pnpm

COPY --from=deps /app/node_modules ./node_modules
COPY . .

ARG NEXT_PUBLIC_API_URL
ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL
ENV NEXT_TELEMETRY_DISABLED=1
ENV NODE_ENV=production

RUN pnpm build

FROM node:20-alpine AS runner
WORKDIR /app

ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1
ENV PORT=3000

RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs

COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

EXPOSE 3000

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:3000/api/health || exit 1

CMD ["node", "server.js"]
```

## .dockerignore Template

```
# Dependencies
node_modules/
vendor/

# Build outputs
dist/
build/
.next/

# Development
.git/
.env*
!.env.example
*.log

# IDE
.idea/
.vscode/
*.swp

# Testing
coverage/
.nyc_output/

# Docker
Dockerfile*
docker-compose*
.docker/
```

## Common Commands

```bash
# Build images
docker build -f docker/Dockerfile.frontend -t frontend:latest .
docker build -f docker/Dockerfile.backend -t backend:latest .

# Development
docker compose -f docker/docker-compose.yml up -d
docker compose -f docker/docker-compose.yml logs -f
docker compose -f docker/docker-compose.yml down

# Production
docker compose -f docker/docker-compose.prod.yml up -d

# Testing
docker compose -f docker/docker-compose.test.yml up --abort-on-container-exit

# Cleanup
docker system prune -af
docker volume prune -f
```

## Security Checklist

- [ ] Using specific base image versions
- [ ] Running as non-root user
- [ ] No secrets in images or layers
- [ ] Health checks defined
- [ ] Minimal attack surface (alpine images)
- [ ] .dockerignore configured
- [ ] Images scanned for vulnerabilities
- [ ] Read-only root filesystem where possible
