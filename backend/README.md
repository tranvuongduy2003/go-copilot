# Go Copilot Backend

A production-ready backend API built with Go, following Clean Architecture, DDD, and CQRS patterns.

## Tech Stack

- **Language**: Go 1.23+
- **HTTP Router**: chi/v5
- **Database**: PostgreSQL with pgx/v5
- **Cache**: Redis
- **Migrations**: golang-migrate/v4
- **Logging**: zap
- **Validation**: go-playground/validator
- **Configuration**: viper
- **Authentication**: JWT

## Prerequisites

- Go 1.23 or higher
- Docker and Docker Compose
- Make
- golangci-lint (for linting)

## Project Structure

```
backend/
├── cmd/
│   └── api/              # Application entrypoints
├── internal/
│   ├── domain/           # Domain layer (entities, value objects, events)
│   ├── application/      # Application layer (commands, queries, DTOs)
│   ├── infrastructure/   # Infrastructure layer (repositories, external services)
│   └── interfaces/       # Interface layer (HTTP handlers, middleware)
├── pkg/                  # Shared packages (config, logger, validator)
├── migrations/           # Database migrations
└── ...
```

## Quick Start

### 1. Clone and Setup

```bash
# Clone the repository
git clone <repository-url>
cd backend

# Copy environment file
cp .env.example .env

# Install dependencies and tools
make setup
```

### 2. Start Development Environment

```bash
# Start PostgreSQL, Redis via Docker
make dev

# Run the application
make run
```

### 3. Run Tests

```bash
# Run unit tests
make test

# Run integration tests
make test-integration

# Run tests with coverage
make test-coverage
```

## Available Make Commands

Run `make help` to see all available commands:

| Command             | Description                              |
|---------------------|------------------------------------------|
| `make help`         | Show available commands                  |
| `make setup`        | First-time project setup                 |
| `make dev`          | Start docker-compose services            |
| `make run`          | Run application in development mode      |
| `make build`        | Build production binary                  |
| `make test`         | Run unit tests with coverage             |
| `make test-integration` | Run integration tests                |
| `make lint`         | Run golangci-lint                        |
| `make fmt`          | Format code                              |
| `make migrate-up`   | Run database migrations                  |
| `make migrate-down` | Rollback last migration                  |
| `make generate`     | Run code generation                      |
| `make clean`        | Remove build artifacts                   |

## Configuration

Configuration is loaded from environment variables. See `.env.example` for all available options.

### Required Environment Variables

| Variable      | Description               | Default     |
|---------------|---------------------------|-------------|
| `SERVER_PORT` | HTTP server port          | `8080`      |
| `DB_HOST`     | PostgreSQL host           | `localhost` |
| `DB_PORT`     | PostgreSQL port           | `5432`      |
| `DB_USER`     | PostgreSQL user           | `postgres`  |
| `DB_PASSWORD` | PostgreSQL password       | `postgres`  |
| `DB_NAME`     | PostgreSQL database name  | `app_dev`   |
| `JWT_SECRET`  | JWT signing secret        | -           |

## API Documentation

API documentation is available at `/api/docs` when running in development mode.

## Database Migrations

This project uses [golang-migrate](https://github.com/golang-migrate/migrate) for database migrations.

### Migration Commands

```bash
# Create a new migration
make migrate-create name=create_users_table

# Run all pending migrations
make migrate-up

# Run next migration only
make migrate-up-one

# Rollback last migration
make migrate-down

# Rollback all migrations
make migrate-down-all

# Show current migration version
make migrate-version

# Force set migration version (for fixing dirty state)
make migrate-force version=1

# Reset database (rollback all and re-run)
make migrate-reset
```

### Migration Workflow

1. **Creating a Migration**
   ```bash
   make migrate-create name=add_users_indexes
   ```
   This creates two files: `{version}_{name}.up.sql` and `{version}_{name}.down.sql`

2. **Writing Migrations**
   ```sql
   -- 000001_create_users.up.sql
   CREATE TABLE users (...);

   -- 000001_create_users.down.sql
   DROP TABLE users;
   ```

3. **Testing Migrations**
   ```bash
   # Apply migration
   make migrate-up

   # Verify it works, then rollback
   make migrate-down

   # Re-apply to confirm both directions work
   make migrate-up
   ```

4. **Environment-Specific DSN**
   ```bash
   # Override database connection for different environments
   DB_DSN="postgres://user:pass@host:5432/dbname?sslmode=disable" make migrate-up
   ```

### Current Migrations

| Version | Name | Description |
|---------|------|-------------|
| 000001 | create_users_table | Creates users table with all fields |
| 000002 | add_users_indexes | Adds indexes for performance |
| 000003 | add_updated_at_trigger | Auto-updates updated_at on row changes |

## Code Quality

```bash
# Run linter
make lint

# Format code
make fmt

# Run all checks before commit
make check
```

## Architecture

This project follows Clean Architecture principles with the following layers:

1. **Domain Layer** (`internal/domain/`)
   - Entities and Aggregate Roots
   - Value Objects
   - Domain Events
   - Repository Interfaces

2. **Application Layer** (`internal/application/`)
   - Commands (CQRS write operations)
   - Queries (CQRS read operations)
   - DTOs

3. **Infrastructure Layer** (`internal/infrastructure/`)
   - Repository Implementations
   - External Service Adapters
   - Database Connections

4. **Interface Layer** (`internal/interfaces/`)
   - HTTP Handlers
   - Middleware
   - Request/Response DTOs

## License

[MIT License](LICENSE)
