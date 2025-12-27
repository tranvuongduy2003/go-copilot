# Dockerfile Builder Skill

Generate optimized Dockerfiles for Go and React applications.

## Usage

```
/project:skill:dockerfile <type>
```

Types: `go-api`, `react-spa`, `fullstack`

## Go API Dockerfile

**`docker/Dockerfile.api`**
```dockerfile
# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.version=${VERSION:-dev}" \
    -o /app/api \
    ./cmd/api

# Runtime stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/api .

# Copy migrations if needed
COPY --from=builder /app/migrations ./migrations

# Create non-root user
RUN addgroup -g 1000 -S appgroup && \
    adduser -u 1000 -S appuser -G appgroup

# Set ownership
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run binary
ENTRYPOINT ["./api"]
```

## React SPA Dockerfile

**`docker/Dockerfile.web`**
```dockerfile
# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# Install dependencies
COPY package.json package-lock.json ./
RUN npm ci

# Copy source code
COPY . .

# Build application
ARG VITE_API_URL
ENV VITE_API_URL=${VITE_API_URL}

RUN npm run build

# Runtime stage with nginx
FROM nginx:alpine

# Copy custom nginx config
COPY docker/nginx.conf /etc/nginx/nginx.conf

# Copy built files
COPY --from=builder /app/dist /usr/share/nginx/html

# Create non-root user
RUN addgroup -g 1000 -S appgroup && \
    adduser -u 1000 -S appuser -G appgroup && \
    chown -R appuser:appgroup /var/cache/nginx && \
    chown -R appuser:appgroup /var/log/nginx && \
    touch /var/run/nginx.pid && \
    chown -R appuser:appgroup /var/run/nginx.pid

USER appuser

EXPOSE 80

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:80 || exit 1

CMD ["nginx", "-g", "daemon off;"]
```

## Nginx Configuration

**`docker/nginx.conf`**
```nginx
user nginx;
worker_processes auto;
error_log /var/log/nginx/error.log warn;
pid /var/run/nginx.pid;

events {
    worker_connections 1024;
}

http {
    include /etc/nginx/mime.types;
    default_type application/octet-stream;

    log_format main '$remote_addr - $remote_user [$time_local] "$request" '
                    '$status $body_bytes_sent "$http_referer" '
                    '"$http_user_agent" "$http_x_forwarded_for"';

    access_log /var/log/nginx/access.log main;

    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    keepalive_timeout 65;
    types_hash_max_size 2048;

    # Gzip compression
    gzip on;
    gzip_vary on;
    gzip_proxied any;
    gzip_comp_level 6;
    gzip_types text/plain text/css text/xml application/json application/javascript
               application/xml application/xml+rss text/javascript application/wasm;

    server {
        listen 80;
        server_name localhost;
        root /usr/share/nginx/html;
        index index.html;

        # Security headers
        add_header X-Frame-Options "SAMEORIGIN" always;
        add_header X-Content-Type-Options "nosniff" always;
        add_header X-XSS-Protection "1; mode=block" always;
        add_header Referrer-Policy "strict-origin-when-cross-origin" always;

        # SPA routing - try files, fallback to index.html
        location / {
            try_files $uri $uri/ /index.html;
        }

        # Cache static assets
        location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2)$ {
            expires 1y;
            add_header Cache-Control "public, immutable";
        }

        # Don't cache HTML
        location ~* \.html$ {
            expires -1;
            add_header Cache-Control "no-store, no-cache, must-revalidate";
        }

        # Health check endpoint
        location /health {
            access_log off;
            return 200 "healthy\n";
            add_header Content-Type text/plain;
        }
    }
}
```

## Docker Compose

**`docker-compose.yml`**
```yaml
version: '3.8'

services:
  api:
    build:
      context: ./backend
      dockerfile: ../docker/Dockerfile.api
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://postgres:postgres@database:5432/app?sslmode=disable
      - REDIS_URL=redis://cache:6379
      - JWT_SECRET=${JWT_SECRET}
    depends_on:
      database:
        condition: service_healthy
      cache:
        condition: service_started
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 10s
    restart: unless-stopped

  web:
    build:
      context: ./frontend
      dockerfile: ../docker/Dockerfile.web
      args:
        - VITE_API_URL=http://localhost:8080
    ports:
      - "3000:80"
    depends_on:
      - api
    restart: unless-stopped

  database:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: app
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped

  cache:
    image: redis:7-alpine
    command: redis-server --appendonly yes
    volumes:
      - redis_data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped

  migrate:
    build:
      context: ./backend
      dockerfile: ../docker/Dockerfile.api
    command: ["migrate", "-path", "/app/migrations", "-database", "$DATABASE_URL", "up"]
    environment:
      - DATABASE_URL=postgres://postgres:postgres@database:5432/app?sslmode=disable
    depends_on:
      database:
        condition: service_healthy
    restart: "no"

volumes:
  postgres_data:
  redis_data:
```

## Development Docker Compose

**`docker-compose.dev.yml`**
```yaml
version: '3.8'

services:
  api:
    build:
      context: ./backend
      dockerfile: ../docker/Dockerfile.api.dev
    ports:
      - "8080:8080"
    volumes:
      - ./backend:/app
    environment:
      - DATABASE_URL=postgres://postgres:postgres@database:5432/app?sslmode=disable
      - REDIS_URL=redis://cache:6379
    depends_on:
      - database
      - cache

  database:
    image: postgres:16-alpine
    ports:
      - "5432:5432"
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: app
    volumes:
      - postgres_dev_data:/var/lib/postgresql/data

  cache:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_dev_data:/data

volumes:
  postgres_dev_data:
  redis_dev_data:
```

## Commands

```bash
# Build images
docker-compose build

# Start all services
docker-compose up -d

# View logs
docker-compose logs -f api

# Stop all services
docker-compose down

# Remove volumes
docker-compose down -v

# Run migrations
docker-compose run --rm migrate

# Access database
docker-compose exec database psql -U postgres -d app
```

## Checklist

- [ ] Multi-stage build for smaller images
- [ ] Non-root user for security
- [ ] Health checks configured
- [ ] Proper layer caching (COPY order)
- [ ] Security headers in nginx
- [ ] Gzip compression enabled
- [ ] Static asset caching
- [ ] Environment variables for configuration
- [ ] Volume mounts for persistence
- [ ] Service dependencies defined
