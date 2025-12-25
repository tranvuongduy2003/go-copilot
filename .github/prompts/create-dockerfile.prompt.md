---
description: Create a production-ready Dockerfile with multi-stage builds and security best practices
---

# Create Dockerfile

Create a production-ready Dockerfile following best practices for security, performance, and maintainability.

## Tech Stack

This project uses:
- **Backend**: Go 1.25+ with Chi router (Clean Architecture + DDD + CQRS)
- **Frontend**: React 19 with pnpm (Next.js standalone output)
- **Package Manager**: pnpm (frontend)

## Service Details

**Service Name**: {{serviceName}}

**Service Type**: {{serviceType}}
- [ ] Go (API server) - **Primary backend**
- [ ] React 19 (Next.js frontend) - **Primary frontend**
- [ ] Python (FastAPI/Django)
- [ ] Static files (Nginx)

**Port**: {{port}}

**Build Command**: {{buildCommand}}

**Start Command**: {{startCommand}}

## Requirements

### Multi-stage Build

Create a multi-stage Dockerfile with:
1. **builder** stage - Compile/build application
2. **runner** stage - Minimal production image (~10-20MB for Go)

### Security

- [ ] Use specific base image version (no `latest`)
- [ ] Create non-root user
- [ ] Run as non-root user
- [ ] Use read-only filesystem where possible
- [ ] Drop all capabilities
- [ ] No secrets in image

### Health Check

- [ ] Add HEALTHCHECK instruction
- [ ] Use appropriate health endpoint
- [ ] Configure interval, timeout, retries

### Optimization

- [ ] Use alpine-based images
- [ ] Order layers for caching
- [ ] Remove unnecessary files
- [ ] Use .dockerignore

## Implementation Steps

### 1. Create Dockerfile

Location: `docker/Dockerfile.{{serviceName}}`

#### For Go Backend:

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

ENV PORT={{port}}

EXPOSE {{port}}

HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:{{port}}/health || exit 1

ENTRYPOINT ["/app/server"]
```

#### For React 19 Frontend:

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

ENV NEXT_TELEMETRY_DISABLED=1
ENV NODE_ENV=production

RUN pnpm build

FROM node:20-alpine AS runner
WORKDIR /app

ENV NODE_ENV=production
ENV PORT={{port}}

RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs

COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

EXPOSE {{port}}

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:{{port}}/api/health || exit 1

CMD ["node", "server.js"]
```

### 2. Create .dockerignore

Location: `docker/.dockerignore` or root `.dockerignore`

```
node_modules/
dist/
.git/
.env*
!.env.example
*.log
coverage/
.nyc_output/
```

### 3. Add to docker-compose.yml

```yaml
services:
  {{serviceName}}:
    build:
      context: ..
      dockerfile: docker/Dockerfile.{{serviceName}}
    ports:
      - "{{port}}:{{port}}"
    environment:
      - NODE_ENV=production
    healthcheck:
      test: ["CMD", "wget", "--spider", "http://localhost:{{port}}/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

## Validation

After implementation:

```bash
# Build image
docker build -f docker/Dockerfile.{{serviceName}} -t {{serviceName}}:test .

# Check image size
docker images {{serviceName}}:test

# Run container
docker run -d --name {{serviceName}}-test -p {{port}}:{{port}} {{serviceName}}:test

# Check health
curl http://localhost:{{port}}/health

# Check non-root user
docker exec {{serviceName}}-test whoami

# Scan for vulnerabilities
trivy image {{serviceName}}:test

# Cleanup
docker rm -f {{serviceName}}-test
```

## Output

Provide:
1. Complete Dockerfile
2. Updated .dockerignore
3. Docker Compose service configuration
4. Build and run commands
5. Any additional configuration needed
