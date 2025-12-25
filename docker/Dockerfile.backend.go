# =============================================================================
# Go Backend - Multi-stage Production Build (Legacy Name)
# =============================================================================
# NOTE: This file is kept for backwards compatibility.
# Primary backend Dockerfile is now Dockerfile.backend
# =============================================================================

# -----------------------------------------------------------------------------
# Stage 1: Builder
# -----------------------------------------------------------------------------
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files
COPY backend/go.mod backend/go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY backend/ .

# Build arguments
ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT

# Build with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}" \
    -o /app/server ./cmd/server

# -----------------------------------------------------------------------------
# Stage 2: Production Runner
# -----------------------------------------------------------------------------
FROM alpine:3.20 AS runner
WORKDIR /app

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata wget && \
    update-ca-certificates

# Security: Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy binary
COPY --from=builder /app/server /app/server

# Copy migrations if present (optional - uncomment if migrations folder exists)
# COPY --from=builder /app/migrations ./migrations/

# Set ownership
RUN chown -R appuser:appgroup /app

USER appuser

ENV PORT=8080

EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

ENTRYPOINT ["/app/server"]
