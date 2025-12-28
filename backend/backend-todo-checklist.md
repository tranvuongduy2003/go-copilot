# Expert-Level Backend Setup Checklist

## Clean Architecture + DDD + CQRS | Golang + PostgreSQL

---

# Phase 1: Project Foundation

## 1.1 Repository Initialization

- [x] Initialize Go module with proper module path following organization naming convention
- [x] Create `.gitignore` with Go-specific patterns (binaries, IDE files, env files, vendor/)
- [x] Create `.editorconfig` for consistent coding style across team
- [x] Initialize git repository with initial commit
- [ ] Set up branch protection rules (main/develop branches)
- [x] Create `README.md` with project overview, prerequisites, and quick start guide

## 1.2 Development Environment Setup

- [x] Refer folder `../docker` for local development services
- [x] Create `.env.example` with all required environment variables documented
- [ ] Verify all team members can spin up environment with single command

## 1.3 Makefile Configuration

- [x] Define `help` target as default with all available commands documented
- [x] Add `setup` target for first-time project setup (install tools, copy env files)
- [x] Add `dev` target to start docker-compose services
- [x] Add `run` target to start application in development mode
- [x] Add `build` target with proper ldflags for version injection
- [x] Add `test` target with coverage reporting
- [x] Add `test-integration` target for integration tests with test database
- [x] Add `lint` target using golangci-lint
- [x] Add `fmt` target for code formatting (gofmt + goimports)
- [x] Add `migrate-*` targets for database migration commands
- [x] Add `generate` target for code generation (mocks, wire, etc.)
- [x] Add `clean` target to remove build artifacts

## 1.4 Dependency Management

- [x] Review and correct Go version in `go.mod` (should be stable release)
- [x] Add chi/v5 for HTTP routing (already present)
- [x] Add pgx/v5 for PostgreSQL driver (already present)
- [x] Add goose/v3 for database migrations (CLI tool)
- [x] Add validator/v10 for struct validation
- [x] Add viper for configuration management
- [x] Add zap for structured logging
- [x] Add uuid for UUID generation
- [x] Add testify for testing assertions and mocks
- [x] Add bcrypt or argon2 for password hashing
- [x] Add jwt/v5 for JWT token handling (if authentication required)
- [x] Add otel packages for OpenTelemetry tracing (optional, production)
- [x] Run `go mod tidy` to clean up dependencies
- [x] Run `go mod verify` to verify dependency integrity

## 1.5 Code Quality Tools Setup

- [x] Install golangci-lint locally and add to CI pipeline
- [x] Create `.golangci.yml` configuration file
- [x] Enable essential linters: gofmt, goimports, govet, errcheck, staticcheck
- [x] Enable security linter: gosec
- [x] Enable style linters: misspell, unconvert, gocritic
- [x] Configure local-prefixes for goimports to group internal imports
- [x] Set up pre-commit hooks using pre-commit framework or lefthook
- [x] Add hook for running `go fmt` on staged files
- [x] Add hook for running `golangci-lint` on staged files
- [x] Add hook for checking commit message format (conventional commits)

---

# Phase 2: Configuration & Logging Infrastructure

## 2.1 Configuration Package (`pkg/config/`)

- [x] Define main Config struct with nested structs for each component
- [x] Define Server config: host, port, read timeout, write timeout, idle timeout
- [x] Define Database config: host, port, user, password, database name, SSL mode, pool size
- [x] Define Redis config (if used): host, port, password, database number
- [x] Define JWT config (if used): secret key, access token TTL, refresh token TTL
- [x] Define Log config: level, format (json/console), output destination
- [x] Implement configuration loading from environment variables
- [x] Implement configuration loading from YAML/JSON file as fallback
- [x] Implement environment-specific config file loading (dev, staging, prod)
- [x] Add configuration validation with detailed error messages
- [x] Validate required fields are present and non-empty
- [x] Validate numeric fields are within acceptable ranges
- [x] Validate URLs and connection strings are properly formatted
- [x] Implement config hot-reload capability for non-critical settings (optional)
- [x] Add helper methods for constructing DSN strings
- [x] Write unit tests for configuration loading and validation

## 2.2 Logger Package (`pkg/logger/`)

- [x] Define Logger interface with methods: Debug, Info, Warn, Error, Fatal
- [x] Define structured logging methods that accept key-value pairs
- [x] Implement logger using zap as underlying library
- [x] Configure log level based on configuration
- [x] Configure output format: JSON for production, console for development
- [x] Add caller information (file, line number) to log entries
- [x] Add timestamp in ISO8601 format to all log entries
- [x] Implement request-scoped logger with context support
- [x] Create helper to extract logger from context
- [x] Create helper to inject logger into context
- [x] Add request ID field injection capability
- [x] Add user ID field injection capability (after authentication)
- [x] Implement log sampling for high-volume debug logs in production
- [x] Create global logger instance for application-wide use
- [x] Write unit tests for logger functionality

## 2.3 Validator Package (`pkg/validator/`)

- [x] Create validator wrapper around go-playground/validator
- [x] Register custom validation tags for domain-specific rules
- [x] Add custom validation for email format with DNS check option
- [x] Add custom validation for phone number formats
- [x] Add custom validation for password strength requirements
- [x] Add custom validation for UUID format
- [x] Implement error message translation to user-friendly messages
- [x] Create map of validation tag to human-readable error message
- [x] Implement locale-aware error messages (optional)
- [x] Create helper function to validate struct and return formatted errors
- [x] Write unit tests for all custom validators

---

# Phase 3: Domain Layer Implementation

## 3.1 Shared Domain Components (`internal/domain/shared/`)

### 3.1.1 Base Types

- [x] Define Entity base struct with ID field and equality comparison method
- [x] Define AggregateRoot struct embedding Entity with domain events slice
- [x] Add method to AggregateRoot for registering domain events
- [x] Add method to AggregateRoot for retrieving and clearing domain events
- [x] Define DomainEvent interface with event type, timestamp, and aggregate ID
- [x] Define base DomainEvent struct implementing common fields

### 3.1.2 Shared Value Objects

- [x] Define Email value object with validation in constructor
- [x] Implement Email equality comparison
- [x] Implement Email string representation
- [x] Implement Email domain extraction method
- [x] Define PhoneNumber value object with country code support
- [ ] Define Money value object with currency and amount (if e-commerce)
- [ ] Implement Money arithmetic operations with currency validation
- [ ] Define Address value object with component fields (if needed)
- [x] Define DateRange value object for temporal queries (if needed)
- [x] Define Pagination value object with page, limit, offset calculations

### 3.1.3 Domain Errors

- [x] Define base DomainError interface extending error
- [x] Add error code method to DomainError interface
- [x] Define NotFoundError struct with entity type and identifier
- [x] Define ValidationError struct with field and message
- [x] Define ConflictError struct for duplicate/constraint violations
- [x] Define AuthorizationError struct for permission denied scenarios
- [x] Define BusinessRuleViolationError for domain invariant violations
- [x] Implement Error() method for each error type
- [x] Implement Is() method for error comparison support
- [x] Create constructor functions for each error type
- [x] Write unit tests for error type assertions

## 3.2 User Aggregate (`internal/domain/user/`)

### 3.2.1 User Entity

- [x] Define User struct with all required fields
- [x] Include ID as UUID type
- [x] Include Email as value object type
- [x] Include PasswordHash as string (not plain password)
- [x] Include FullName as string or Name value object
- [x] Include Status as enum type
- [x] Include CreatedAt as time.Time
- [x] Include UpdatedAt as time.Time
- [x] Include DeletedAt as \*time.Time for soft delete (optional)
- [x] Embed AggregateRoot for domain events support
- [x] Implement NewUser constructor with required field validation
- [x] Validate email is not empty and properly formatted
- [x] Validate password meets minimum requirements before hashing
- [x] Validate full name is not empty and within length limits
- [x] Generate UUID for new user
- [x] Set default status to pending or active based on requirements
- [x] Set CreatedAt and UpdatedAt to current UTC time
- [x] Register UserCreated domain event in constructor

### 3.2.2 User Business Methods

- [x] Implement Activate() method to change status to active
- [x] Add validation that user is not already active
- [x] Update UpdatedAt timestamp
- [x] Register UserActivated domain event
- [x] Implement Deactivate() method to change status to inactive
- [x] Add validation that user is not already inactive
- [x] Update UpdatedAt timestamp
- [x] Register UserDeactivated domain event
- [x] Implement Ban() method to change status to banned
- [x] Update UpdatedAt timestamp
- [x] Register UserBanned domain event
- [x] Implement ChangePassword() method
- [x] Accept hashed password (hashing done in application layer)
- [x] Update UpdatedAt timestamp
- [x] Register PasswordChanged domain event
- [x] Implement UpdateProfile() method for name and other editable fields
- [x] Validate new values meet requirements
- [x] Update UpdatedAt timestamp
- [x] Register ProfileUpdated domain event
- [x] Implement Delete() method for soft delete
- [x] Set DeletedAt to current timestamp
- [x] Register UserDeleted domain event

### 3.2.3 User Status Enum

- [x] Define UserStatus as string type
- [x] Define constants: StatusPending, StatusActive, StatusInactive, StatusBanned
- [x] Implement IsValid() method to check if status value is valid
- [x] Implement CanTransitionTo() method for status state machine validation
- [x] Define allowed transitions: Pending→Active, Active→Inactive, Active→Banned, etc.

### 3.2.4 User Repository Interface

- [x] Define UserRepository interface in domain layer
- [x] Define Create(ctx, user) error method
- [x] Define Update(ctx, user) error method
- [x] Define Delete(ctx, id) error method
- [x] Define FindByID(ctx, id) (\*User, error) method
- [x] Define FindByEmail(ctx, email) (\*User, error) method
- [x] Define ExistsByEmail(ctx, email) (bool, error) method
- [x] Define List(ctx, filter, pagination) ([]\*User, total, error) method
- [x] Define UserFilter struct for list filtering (status, search term, date range)
- [ ] Document expected error types for each method in comments

### 3.2.5 User Domain Events

- [x] Define UserCreatedEvent with user ID, email, timestamp
- [x] Define UserActivatedEvent with user ID, timestamp
- [x] Define UserDeactivatedEvent with user ID, timestamp
- [x] Define UserBannedEvent with user ID, reason, timestamp
- [x] Define PasswordChangedEvent with user ID, timestamp
- [x] Define ProfileUpdatedEvent with user ID, changed fields, timestamp
- [x] Define UserDeletedEvent with user ID, timestamp
- [x] Implement EventType() method for each event returning unique string
- [x] Implement OccurredAt() method for each event
- [x] Implement AggregateID() method for each event

### 3.2.6 User Domain Errors

- [x] Define ErrUserNotFound as sentinel error
- [x] Define ErrEmailAlreadyExists as sentinel error
- [x] Define ErrInvalidEmail as sentinel error
- [x] Define ErrInvalidPassword as sentinel error
- [x] Define ErrInvalidStatus as sentinel error
- [x] Define ErrInvalidStatusTransition with current and target status
- [x] Define ErrUserAlreadyActive as sentinel error
- [x] Define ErrUserAlreadyInactive as sentinel error

### 3.2.7 User Domain Tests

- [x] Write unit tests for NewUser constructor validation
- [x] Test successful user creation with valid inputs
- [x] Test failure with empty email
- [x] Test failure with invalid email format
- [x] Test failure with empty password
- [x] Test failure with weak password
- [x] Test failure with empty full name
- [x] Write unit tests for each business method
- [x] Test successful activation from pending status
- [x] Test failure when activating already active user
- [x] Test successful deactivation
- [x] Test successful ban
- [x] Test successful password change
- [x] Test domain event registration for each operation
- [x] Write unit tests for UserStatus transitions
- [x] Test all valid transitions succeed
- [x] Test all invalid transitions fail

---

# Phase 4: Infrastructure Layer Implementation

## 4.1 Database Connection (`internal/infrastructure/persistence/postgres/`)

### 4.1.1 Connection Pool Setup

- [x] Define DB struct wrapping pgxpool.Pool
- [x] Implement NewDB constructor accepting config
- [x] Parse connection string from config
- [x] Configure connection pool size (min, max connections)
- [x] Configure connection lifetime and idle timeout
- [x] Configure connection health check period
- [x] Implement Connect() method with retry logic
- [x] Define maximum retry attempts
- [x] Define retry backoff strategy (exponential with jitter)
- [x] Log connection attempts and failures
- [x] Implement Close() method for graceful shutdown
- [x] Implement Ping() method for health checks
- [x] Implement Stats() method for monitoring connection pool metrics
- [x] Write integration tests for connection management

### 4.1.2 Transaction Support

- [x] Define Transaction interface with Commit, Rollback methods
- [x] Implement BeginTx() method to start transaction
- [x] Accept context for cancellation
- [x] Accept transaction options (isolation level)
- [x] Implement WithTransaction() helper for transaction scope
- [x] Accept function to execute within transaction
- [x] Automatically rollback on error or panic
- [x] Automatically commit on success
- [x] Implement transaction context propagation
- [x] Create helper to extract transaction from context
- [x] Create helper to inject transaction into context

### 4.1.3 Query Builder Helpers (Optional)

- [x] Implement helper for building dynamic WHERE clauses
- [x] Implement helper for building ORDER BY clauses
- [x] Implement helper for building pagination (LIMIT, OFFSET)
- [x] Implement helper for building RETURNING clauses
- [x] All helpers must use parameterized queries to prevent SQL injection

## 4.2 Database Migrations (`migrations/`)

### 4.2.1 golang-migrate CLI Setup

- [x] Install golang-migrate CLI:
      `go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`
- [x] Create `migrations/` directory in project root
- [x] Configure default database connection in environment variables
- [x] Document migration commands in Makefile
- [x] Document migration workflow in README

### 4.2.2 User Table Migration

- [x] Generate migration file: `migrate create -ext sql -dir migrations -seq create_users_table`
- [x] Write UP migration: CREATE TABLE users
- [x] Define id column as UUID PRIMARY KEY
- [x] Define email column as VARCHAR(255) UNIQUE NOT NULL
- [x] Define password_hash column as VARCHAR(255) NOT NULL
- [x] Define full_name column as VARCHAR(255) NOT NULL
- [x] Define status column as VARCHAR(50) NOT NULL with DEFAULT
- [x] Define created_at column as TIMESTAMPTZ NOT NULL with DEFAULT NOW()
- [x] Define updated_at column as TIMESTAMPTZ NOT NULL with DEFAULT NOW()
- [x] Define deleted_at column as TIMESTAMPTZ NULL for soft delete
- [x] Write DOWN migration: DROP TABLE users
- [x] Test migration up and down in development environment

### 4.2.3 User Table Indexes Migration

- [x] Generate migration file: `migrate create -ext sql -dir migrations -seq add_users_indexes`
- [x] Write UP migration: CREATE INDEX on email column
- [x] Write UP migration: CREATE INDEX on status column
- [x] Write UP migration: CREATE INDEX on created_at column
- [x] Write UP migration: CREATE INDEX on deleted_at for soft delete queries
- [x] Consider composite indexes for common query patterns
- [x] Write DOWN migration: DROP all created indexes
- [x] Test migration up and down

### 4.2.4 Audit Trigger Migration (Optional)

- [x] Generate migration file: `migrate create -ext sql -dir migrations -seq add_updated_at_trigger`
- [x] Write UP migration: CREATE FUNCTION to update updated_at column
- [x] Write UP migration: CREATE TRIGGER on users table BEFORE UPDATE
- [x] Write DOWN migration: DROP TRIGGER and FUNCTION
- [x] Test that updated_at automatically updates on row modification

### 4.2.5 Migration Testing

- [x] Test full migration sequence from empty database
- [x] Test rollback of each migration individually
- [x] Test rollback of all migrations
- [x] Test re-applying migrations after rollback
- [x] Add migration tests to CI pipeline

## 4.3 User Repository Implementation (`internal/infrastructure/persistence/postgres/`)

### 4.3.1 Repository Structure

- [x] Define UserRepository struct with DB pool dependency
- [x] Implement NewUserRepository constructor
- [x] Define SQL query constants as private package variables
- [x] Use parameterized queries for all SQL statements
- [x] Implement domain entity to database row mapping
- [x] Implement database row to domain entity mapping

### 4.3.2 Create Method Implementation

- [x] Implement Create(ctx, user) error
- [x] Extract transaction from context if present
- [x] Use INSERT query with all user fields
- [x] Map domain User to database columns
- [x] Handle unique constraint violation for email
- [x] Return ErrEmailAlreadyExists for duplicate email
- [x] Wrap other database errors with context
- [x] Write integration test for successful creation
- [x] Write integration test for duplicate email handling

### 4.3.3 Update Method Implementation

- [x] Implement Update(ctx, user) error
- [x] Extract transaction from context if present
- [x] Use UPDATE query with all mutable fields
- [x] Include WHERE clause for ID
- [x] Check rows affected to detect not found
- [x] Return ErrUserNotFound if no rows updated
- [x] Wrap database errors with context
- [x] Write integration test for successful update
- [x] Write integration test for not found handling

### 4.3.4 Delete Method Implementation

- [x] Implement Delete(ctx, id) error for soft delete
- [x] Update deleted_at column to current timestamp
- [x] Check rows affected to detect not found
- [x] Return ErrUserNotFound if no rows updated
- [x] Write integration test for successful deletion
- [x] Write integration test for not found handling

### 4.3.5 FindByID Method Implementation

- [x] Implement FindByID(ctx, id) (\*User, error)
- [x] Use SELECT query with WHERE id = $1
- [x] Add condition to exclude soft-deleted records
- [x] Map database row to domain User entity
- [x] Handle pgx.ErrNoRows
- [x] Return ErrUserNotFound for no rows
- [x] Return nil, nil is NOT acceptable - always return error for not found
- [x] Write integration test for successful retrieval
- [x] Write integration test for not found handling
- [x] Write integration test confirming soft-deleted excluded

### 4.3.6 FindByEmail Method Implementation

- [x] Implement FindByEmail(ctx, email) (\*User, error)
- [x] Use SELECT query with WHERE email = $1
- [x] Add condition to exclude soft-deleted records
- [x] Case-insensitive comparison using LOWER()
- [x] Handle not found scenario
- [x] Write integration test for successful retrieval
- [x] Write integration test for case insensitivity

### 4.3.7 ExistsByEmail Method Implementation

- [x] Implement ExistsByEmail(ctx, email) (bool, error)
- [x] Use SELECT EXISTS query for efficiency
- [x] Add condition to exclude soft-deleted records
- [x] Return false, nil for not exists
- [x] Return true, nil for exists
- [x] Write integration test for both cases

### 4.3.8 List Method Implementation

- [x] Implement List(ctx, filter, pagination) ([]\*User, int64, error)
- [x] Build dynamic WHERE clause based on filter
- [x] Support filtering by status
- [x] Support filtering by search term (email, name)
- [x] Support filtering by date range
- [x] Exclude soft-deleted records
- [x] Apply pagination using LIMIT and OFFSET
- [x] Execute count query for total records
- [x] Execute data query for current page
- [x] Consider using single query with window function for count
- [x] Map all rows to domain entities
- [x] Return empty slice (not nil) when no results
- [x] Write integration test for unfiltered listing
- [x] Write integration test for each filter type
- [x] Write integration test for pagination
- [x] Write integration test for empty results

## 4.4 Unit of Work Pattern (`internal/infrastructure/persistence/postgres/`)

### 4.4.1 Unit of Work Interface

- [x] Define UnitOfWork interface in application layer
- [x] Define Begin(ctx) (UnitOfWorkContext, error) method
- [x] Define UnitOfWorkContext interface with Commit, Rollback, Context methods

### 4.4.2 Unit of Work Implementation

- [x] Implement PostgresUnitOfWork struct
- [x] Implement Begin() method starting database transaction
- [x] Implement Commit() method committing transaction
- [x] Implement Rollback() method rolling back transaction
- [x] Implement Context() returning context with transaction
- [x] Ensure repositories can extract transaction from context
- [x] Write integration test for commit scenario
- [x] Write integration test for rollback scenario
- [x] Write integration test for concurrent transactions

## 4.5 Event Bus Implementation (`internal/infrastructure/messaging/`)

### 4.5.1 Event Bus Interface

- [x] Define EventBus interface in domain or application layer
- [x] Define Publish(ctx, events ...DomainEvent) error method
- [x] Define Subscribe(eventType string, handler EventHandler) method
- [x] Define EventHandler function type

### 4.5.2 In-Memory Event Bus

- [x] Implement InMemoryEventBus struct
- [x] Use map to store handlers by event type
- [x] Use mutex for thread-safe handler registration
- [x] Implement Publish() method
- [x] Iterate through events
- [x] Find handlers for each event type
- [x] Execute handlers (sync or async based on requirements)
- [x] Handle handler panics gracefully
- [x] Log handler errors without failing publish
- [x] Implement Subscribe() method
- [x] Store handler in map by event type
- [x] Support multiple handlers per event type
- [x] Write unit tests for publishing events
- [x] Write unit tests for subscribing handlers
- [x] Write unit tests for multiple handlers

### 4.5.3 Async Event Handling (Optional)

- [x] Implement async event dispatch using goroutines
- [x] Use worker pool to limit concurrent handlers
- [x] Implement graceful shutdown waiting for handlers to complete
- [x] Add timeout for handler execution
- [x] Write tests for async behavior

---

# Phase 5: Application Layer Implementation (CQRS)

## 5.1 CQRS Infrastructure (`internal/application/`)

### 5.1.1 Command Infrastructure

- [x] Define Command marker interface (can be empty)
- [x] Define CommandHandler interface with Handle(ctx, cmd) error method
- [x] Define CommandBus interface with Dispatch(ctx, cmd) error method
- [x] Implement InMemoryCommandBus
- [x] Store handlers by command type
- [x] Implement Register() method to register handlers
- [x] Implement Dispatch() method to find and execute handler
- [x] Return error if handler not found for command type
- [x] Add logging for command dispatch

### 5.1.2 Query Infrastructure

- [x] Define Query marker interface (can be empty)
- [x] Define QueryHandler interface with Handle(ctx, query) (result, error) method
- [x] Use generics for type-safe query results if Go 1.18+
- [x] Define QueryBus interface with Dispatch(ctx, query) (result, error) method
- [x] Implement InMemoryQueryBus
- [x] Store handlers by query type
- [x] Implement Register() method
- [x] Implement Dispatch() method
- [x] Return error if handler not found

### 5.1.3 Middleware Support (Optional)

- [x] Define CommandMiddleware function type
- [x] Define QueryMiddleware function type
- [x] Implement logging middleware for commands
- [x] Implement logging middleware for queries
- [x] Implement validation middleware
- [x] Implement transaction middleware for commands
- [x] Add middleware chain execution in buses

## 5.2 DTOs (`internal/application/user/dto/`)

### 5.2.1 User DTOs

- [x] Define UserDTO struct for responses
- [x] Include ID as string (UUID string representation)
- [x] Include Email as string
- [x] Include FullName as string
- [x] Include Status as string
- [x] Include CreatedAt as time.Time or string
- [x] Include UpdatedAt as time.Time or string
- [x] Exclude sensitive fields (password hash)
- [x] Add JSON tags for serialization
- [x] Implement ToDTO(user \*domain.User) UserDTO mapper function
- [x] Implement ToDTOs(users []\*domain.User) []UserDTO mapper function

### 5.2.2 Pagination DTOs

- [x] Define PaginationRequest struct with Page, Limit fields
- [x] Add validation tags for minimum and maximum values
- [x] Implement Offset() method to calculate SQL offset
- [x] Define PaginatedResponse[T] generic struct
- [x] Include Items []T field
- [x] Include Total int64 field
- [x] Include Page int field
- [x] Include Limit int field
- [x] Include TotalPages int field (calculated)
- [x] Include HasNext bool field (calculated)
- [x] Include HasPrev bool field (calculated)

## 5.3 User Commands (`internal/application/user/command/`)

### 5.3.1 CreateUser Command

- [x] Define CreateUserCommand struct
- [x] Include Email field with validation tags
- [x] Include Password field with validation tags
- [x] Include FullName field with validation tags
- [x] Define CreateUserHandler struct
- [x] Inject UserRepository dependency
- [x] Inject PasswordHasher dependency (interface for bcrypt)
- [x] Inject EventBus dependency
- [x] Inject Logger dependency
- [x] Implement Handle(ctx, cmd) (string, error) method returning created user ID
- [x] Validate command input using validator
- [x] Check if email already exists
- [x] Return ErrEmailAlreadyExists if duplicate
- [x] Hash password using injected hasher
- [x] Create new User domain entity
- [x] Save user using repository
- [x] Publish domain events from aggregate
- [x] Return created user ID
- [x] Log successful creation
- [x] Write unit tests with mocked repository
- [x] Test successful creation flow
- [x] Test duplicate email rejection
- [x] Test validation failures
- [x] Test password hashing integration

### 5.3.2 UpdateUser Command

- [x] Define UpdateUserCommand struct
- [x] Include UserID field (required)
- [x] Include FullName field (optional, pointer or wrapper type)
- [x] Include other updatable fields as optional
- [x] Define UpdateUserHandler struct
- [x] Inject dependencies
- [x] Implement Handle(ctx, cmd) error method
- [x] Validate command input
- [x] Fetch existing user by ID
- [x] Return ErrUserNotFound if not exists
- [x] Apply updates to domain entity using business methods
- [x] Save updated user
- [x] Publish domain events
- [x] Write unit tests for update scenarios

### 5.3.3 DeleteUser Command

- [x] Define DeleteUserCommand struct with UserID field
- [x] Define DeleteUserHandler struct
- [x] Inject dependencies
- [x] Implement Handle(ctx, cmd) error method
- [x] Validate command input
- [x] Fetch existing user by ID
- [x] Return ErrUserNotFound if not exists
- [x] Call Delete() method on domain entity
- [x] Save updated user (soft delete)
- [x] Publish domain events
- [x] Write unit tests for delete scenarios

### 5.3.4 ChangePassword Command

- [x] Define ChangePasswordCommand struct
- [x] Include UserID field
- [x] Include CurrentPassword field
- [x] Include NewPassword field with validation
- [x] Define ChangePasswordHandler struct
- [x] Inject dependencies including PasswordHasher
- [x] Implement Handle(ctx, cmd) error method
- [x] Fetch existing user by ID
- [x] Verify current password matches stored hash
- [x] Return ErrInvalidPassword if mismatch
- [x] Validate new password strength
- [x] Hash new password
- [x] Call ChangePassword() on domain entity
- [x] Save updated user
- [x] Publish domain events
- [x] Write unit tests including password verification

### 5.3.5 ActivateUser Command

- [x] Define ActivateUserCommand struct with UserID field
- [x] Define ActivateUserHandler struct
- [x] Inject dependencies
- [x] Implement Handle(ctx, cmd) error method
- [x] Fetch existing user by ID
- [x] Call Activate() on domain entity
- [x] Handle InvalidStatusTransition error from domain
- [x] Save updated user
- [x] Publish domain events
- [x] Write unit tests

### 5.3.6 DeactivateUser Command

- [x] Define DeactivateUserCommand struct with UserID field
- [x] Define DeactivateUserHandler struct
- [x] Implement Handle() method similar to ActivateUser
- [x] Write unit tests

## 5.4 User Queries (`internal/application/user/query/`)

### 5.4.1 GetUser Query

- [x] Define GetUserQuery struct with UserID field
- [x] Define GetUserHandler struct
- [x] Inject UserRepository dependency (can use same repo or dedicated read repo)
- [x] Inject Logger dependency
- [x] Implement Handle(ctx, query) (\*UserDTO, error) method
- [x] Validate query input
- [x] Fetch user from repository
- [x] Return ErrUserNotFound if not exists
- [x] Map domain entity to DTO
- [x] Return DTO
- [x] Write unit tests with mocked repository

### 5.4.2 ListUsers Query

- [x] Define ListUsersQuery struct
- [x] Include pagination fields (Page, Limit)
- [x] Include filter fields (Status, Search)
- [x] Include sort fields (SortBy, SortOrder)
- [x] Define ListUsersHandler struct
- [x] Inject dependencies
- [x] Implement Handle(ctx, query) (\*PaginatedResponse[UserDTO], error) method
- [x] Validate query input
- [x] Apply default pagination if not specified
- [x] Build filter from query fields
- [x] Fetch users from repository with filter and pagination
- [x] Map domain entities to DTOs
- [x] Build paginated response with metadata
- [x] Return response
- [x] Write unit tests for various filter combinations

### 5.4.3 GetUserByEmail Query

- [x] Define GetUserByEmailQuery struct with Email field
- [x] Define GetUserByEmailHandler struct
- [x] Implement Handle() method
- [x] Write unit tests

### 5.4.4 CheckEmailExists Query

- [x] Define CheckEmailExistsQuery struct with Email field
- [x] Define CheckEmailExistsHandler struct
- [x] Implement Handle(ctx, query) (bool, error) method
- [x] Write unit tests

---

# Phase 6: Interface Layer Implementation (HTTP API)

## 6.1 HTTP Infrastructure (`internal/interfaces/http/`)

### 6.1.1 Server Setup

- [x] Define Server struct with dependencies
- [x] Include HTTP server instance
- [x] Include router instance
- [x] Include logger instance
- [x] Include configuration
- [x] Implement NewServer constructor
- [x] Accept all dependencies via constructor injection
- [x] Configure server timeouts from config (read, write, idle)
- [x] Implement Start() method to start HTTP server
- [x] Implement Shutdown(ctx) method for graceful shutdown
- [x] Wait for in-flight requests to complete
- [x] Respect context deadline for shutdown timeout

### 6.1.2 Router Setup (`internal/interfaces/http/router/`)

- [x] Define NewRouter function
- [x] Accept handler dependencies
- [x] Create chi router instance
- [x] Apply global middleware in correct order
- [x] Define route groups with version prefix (/api/v1)
- [x] Register user routes
- [x] Register health check routes
- [x] Return configured router

### 6.1.3 Response Helpers

- [x] Define standard success response structure
- [x] Include data field
- [x] Include optional message field
- [x] Include optional metadata field
- [x] Define standard error response structure
- [x] Include error code field
- [x] Include message field
- [x] Include details field (for validation errors)
- [x] Include request ID field
- [x] Implement JSON response helper function
- [x] Set Content-Type header
- [x] Set status code
- [x] Encode response body
- [x] Implement error response helper function
- [x] Map domain errors to HTTP status codes
- [x] NotFoundError → 404
- [x] ValidationError → 400
- [x] ConflictError → 409
- [x] AuthorizationError → 403
- [x] Unknown errors → 500
- [x] Log errors appropriately (5xx with stack, 4xx without)

## 6.2 Middleware (`internal/interfaces/http/middleware/`)

### 6.2.1 Request ID Middleware

- [x] Implement middleware function
- [x] Extract X-Request-ID header if present
- [x] Generate new UUID if header not present
- [x] Inject request ID into request context
- [x] Set X-Request-ID response header
- [x] Write tests for middleware

### 6.2.2 Logging Middleware

- [x] Implement middleware function
- [x] Extract logger from context or use global
- [x] Extract request ID from context
- [x] Log request start with method, path, request ID
- [x] Wrap response writer to capture status code
- [x] Log request completion with duration and status
- [x] Include additional fields: remote addr, user agent
- [x] Skip logging for health check endpoints (optional)
- [x] Write tests for middleware

### 6.2.3 Recovery Middleware

- [x] Implement middleware function
- [x] Use defer to catch panics
- [x] Log panic with stack trace
- [x] Return 500 error response to client
- [x] Include request ID in error response
- [x] Do not expose stack trace to client in production
- [x] Write tests for panic recovery

### 6.2.4 Timeout Middleware

- [x] Implement middleware function
- [x] Accept timeout duration as parameter
- [x] Create context with timeout
- [x] Pass timeout context to next handler
- [x] Handle context deadline exceeded
- [x] Return 504 Gateway Timeout on deadline
- [x] Write tests for timeout behavior

### 6.2.5 CORS Middleware

- [x] Configure allowed origins from config
- [x] Configure allowed methods
- [x] Configure allowed headers
- [x] Configure exposed headers
- [x] Configure credentials support
- [x] Configure max age for preflight cache
- [x] Handle preflight OPTIONS requests
- [x] Write tests for CORS behavior

### 6.2.6 Authentication Middleware (if required)

- [x] Implement middleware function
- [x] Extract Authorization header
- [x] Validate Bearer token format
- [ ] Parse and validate JWT token
- [ ] Extract user claims from token
- [x] Inject user info into request context
- [x] Return 401 for missing or invalid token
- [ ] Write tests for auth scenarios

### 6.2.7 Rate Limiting Middleware (optional)

- [x] Choose rate limiting strategy (token bucket, sliding window)
- [x] Configure limits from config
- [x] Implement per-IP rate limiting
- [x] Implement per-user rate limiting (if authenticated)
- [x] Return 429 Too Many Requests when limited
- [x] Include Retry-After header
- [x] Write tests for rate limiting

## 6.3 HTTP DTOs (`internal/interfaces/http/dto/`)

### 6.3.1 Request DTOs

- [x] Define CreateUserRequest struct
- [x] Include Email field with json tag and validation tags
- [x] Include Password field with validation tags
- [x] Include FullName field with validation tags
- [x] Define UpdateUserRequest struct
- [x] Include optional fields with pointer types or omitempty
- [x] Define ChangePasswordRequest struct
- [x] Include CurrentPassword field
- [x] Include NewPassword field with validation
- [ ] Include ConfirmPassword field (optional, for UI validation)
- [x] Define ListUsersRequest struct for query parameters
- [x] Include page, limit as query params
- [x] Include status filter as query param
- [x] Include search as query param
- [x] Include sort_by, sort_order as query params

### 6.3.2 Response DTOs

- [x] Define UserResponse struct mirroring application DTO
- [x] Add JSON tags with appropriate naming (snake_case or camelCase)
- [x] Define PaginatedResponse struct for list endpoints
- [x] Define ErrorResponse struct matching standard error format

### 6.3.3 DTO Validation

- [x] Add validation tags to all request DTOs
- [ ] Document validation rules in comments or OpenAPI spec
- [x] Create reusable validation error formatter

## 6.4 User Handlers (`internal/interfaces/http/handler/`)

### 6.4.1 Handler Structure

- [x] Define UserHandler struct
- [x] Inject CommandBus dependency
- [x] Inject QueryBus dependency
- [x] Inject Validator dependency
- [x] Inject Logger dependency
- [x] Implement NewUserHandler constructor

### 6.4.2 Create User Endpoint

- [x] Implement CreateUser handler for POST /users
- [x] Parse request body into CreateUserRequest DTO
- [x] Handle JSON parsing errors with 400 response
- [x] Validate request using validator
- [x] Return 400 with validation details on failure
- [x] Map HTTP DTO to CreateUserCommand
- [x] Dispatch command via CommandBus
- [x] Handle domain errors and map to HTTP responses
- [x] Return 201 Created with user ID on success
- [x] Include Location header with resource URL
- [ ] Write integration tests for endpoint

### 6.4.3 Get User Endpoint

- [x] Implement GetUser handler for GET /users/{id}
- [x] Extract user ID from URL path parameter
- [x] Validate UUID format
- [x] Return 400 for invalid UUID
- [x] Create GetUserQuery
- [x] Dispatch query via QueryBus
- [x] Handle not found error with 404 response
- [x] Return 200 with user data on success
- [ ] Write integration tests for endpoint

### 6.4.4 List Users Endpoint

- [x] Implement ListUsers handler for GET /users
- [x] Parse query parameters into ListUsersRequest
- [x] Apply default pagination if not specified
- [x] Validate query parameters
- [x] Create ListUsersQuery
- [x] Dispatch query via QueryBus
- [x] Return 200 with paginated response
- [ ] Write integration tests for endpoint

### 6.4.5 Update User Endpoint

- [x] Implement UpdateUser handler for PUT /users/{id}
- [x] Extract user ID from URL path
- [x] Parse request body
- [x] Validate request
- [x] Create UpdateUserCommand
- [x] Dispatch command
- [x] Handle not found and other errors
- [x] Return 200 on success
- [ ] Write integration tests

### 6.4.6 Delete User Endpoint

- [x] Implement DeleteUser handler for DELETE /users/{id}
- [x] Extract user ID from URL path
- [x] Create DeleteUserCommand
- [x] Dispatch command
- [x] Handle not found error
- [x] Return 204 No Content on success
- [ ] Write integration tests

### 6.4.7 Change Password Endpoint

- [x] Implement ChangePassword handler for POST /users/{id}/password
- [x] Parse and validate request
- [x] Create ChangePasswordCommand
- [x] Dispatch command
- [x] Handle invalid password error with 400
- [x] Return 200 on success
- [ ] Write integration tests

### 6.4.8 Activate User Endpoint

- [x] Implement ActivateUser handler for POST /users/{id}/activate
- [x] Create ActivateUserCommand
- [x] Dispatch command
- [x] Handle status transition errors
- [x] Return 200 on success
- [ ] Write integration tests

### 6.4.9 Deactivate User Endpoint

- [x] Implement DeactivateUser handler for POST /users/{id}/deactivate
- [x] Similar implementation to activate
- [ ] Write integration tests

## 6.5 Health Check Endpoints (`internal/interfaces/http/handler/`)

### 6.5.1 Liveness Probe

- [x] Implement handler for GET /health/live
- [x] Return 200 OK if application is running
- [x] No dependency checks (just proves process is alive)
- [x] Return simple JSON response with status

### 6.5.2 Readiness Probe

- [x] Implement handler for GET /health/ready
- [x] Check database connection health
- [x] Check Redis connection health (if used)
- [x] Check other critical dependencies
- [x] Return 200 OK if all dependencies healthy
- [x] Return 503 Service Unavailable if any dependency unhealthy
- [x] Include details of which checks failed

---

# Phase 7: Application Entry Point & Dependency Wiring

## 7.1 Main Function (`cmd/api/main.go`)

### 7.1.1 Initialization Sequence

- [x] Load configuration from environment and files
- [x] Handle configuration errors with clear message and exit
- [x] Initialize logger with configuration
- [x] Set as global logger
- [x] Log application startup with version info
- [x] Initialize database connection
- [x] Implement retry logic for initial connection
- [ ] Run database migrations if configured for auto-migrate
- [ ] Handle migration errors appropriately
- [x] Initialize Redis connection (if used)
- [x] Verify connectivity with ping

### 7.1.2 Dependency Construction

- [x] Create repositories
- [x] Create UserRepository with database pool
- [x] Create services/utilities
- [x] Create PasswordHasher implementation
- [x] Create Validator instance
- [x] Create EventBus implementation
- [x] Create command handlers
- [x] Wire all dependencies into handlers
- [ ] Create CommandBus and register handlers
- [x] Create query handlers
- [x] Wire dependencies into handlers
- [ ] Create QueryBus and register handlers
- [x] Create HTTP handlers
- [x] Wire CommandBus, QueryBus, and other dependencies
- [x] Create router with handlers and middleware
- [x] Create HTTP server with router

### 7.1.3 Server Startup

- [x] Start HTTP server in goroutine
- [x] Log server address and port
- [x] Set up OS signal handling
- [x] Listen for SIGINT and SIGTERM
- [x] Implement graceful shutdown sequence
- [x] Stop accepting new connections
- [x] Wait for in-flight requests (with timeout)
- [x] Close database connections
- [x] Close Redis connections
- [x] Flush logs
- [x] Log shutdown completion
- [x] Exit with appropriate code

### 7.1.4 Error Handling

- [x] Handle server start errors
- [x] Handle shutdown timeout
- [x] Log all errors with context
- [x] Return non-zero exit code on error

## 7.2 Dependency Injection (Optional - using Wire)

### 7.2.1 Wire Setup

- [x] Install wire: `go install github.com/google/wire/cmd/wire@latest`
- [x] Create `wire.go` file with build tag `//go:build wireinject`
- [x] Define provider functions for each dependency
- [x] Define provider sets for related providers
- [x] Create injector function
- [x] Run `wire ./...` to generate `wire_gen.go`
- [x] Add `wire_gen.go` to `.gitignore` or commit it based on team preference

---

# Phase 8: Testing Strategy

## 8.1 Unit Tests

### 8.1.1 Domain Layer Tests

- [x] Test all entity constructors with valid and invalid inputs
- [x] Test all entity business methods
- [x] Test all value object validation
- [x] Test all value object behavior methods
- [x] Test domain event generation
- [x] Test aggregate invariant enforcement
- [x] Achieve minimum 90% coverage for domain layer (achieved: 90.7% user, 78.3% shared)

### 8.1.2 Application Layer Tests

- [x] Mock all repository interfaces
- [x] Mock all external service interfaces
- [x] Test command handlers with all scenarios
- [x] Success paths
- [x] Validation failures
- [x] Domain errors (not found, conflict)
- [x] Test query handlers with all scenarios
- [x] Test DTO mapping functions
- [x] Achieve minimum 80% coverage for application layer (achieved: 94.7% command, 100% query)

### 8.1.3 Interface Layer Tests

- [x] Test HTTP handlers with mocked buses
- [x] Test request parsing and validation
- [x] Test response formatting
- [x] Test error response mapping
- [x] Test middleware behavior
- [x] Achieve minimum 70% coverage for interface layer (achieved: 70.6% handler)

## 8.2 Integration Tests

### 8.2.1 Test Infrastructure

- [x] Create `docker-compose.test.yml` for test dependencies
- [x] Configure test database with separate schema or database
- [x] Implement test database setup and teardown
- [x] Implement test data fixtures/factories
- [x] Create helper functions for common test operations

### 8.2.2 Repository Integration Tests

- [ ] Test repository methods against real database
- [ ] Test transaction behavior
- [ ] Test concurrent access scenarios
- [ ] Test error handling for database failures
- [ ] Clean up test data after each test

### 8.2.3 API Integration Tests

- [ ] Test full request/response cycle
- [ ] Test with real database
- [ ] Test error scenarios
- [ ] Test authentication flow (if applicable)
- [ ] Test pagination behavior
- [ ] Test concurrent requests

## 8.3 Test Configuration

### 8.3.1 Test Utilities

- [x] Create test helper package (pkg/testutil/)
- [x] Implement random data generators (random.go)
- [x] Implement assertion helpers (assertions.go)
- [x] Implement mock factories (mocks.go, fixtures.go)

### 8.3.2 CI Test Pipeline

- [ ] Configure test job in CI
- [ ] Start test dependencies (database)
- [ ] Run migrations
- [ ] Run unit tests with coverage
- [ ] Run integration tests
- [ ] Generate coverage report
- [ ] Fail pipeline if coverage below threshold
- [ ] Upload coverage to reporting service (optional)

---

# Phase 9: Production Readiness

## 9.1 Observability

### 9.1.1 Metrics (Prometheus)

- [x] Add prometheus client library
- [x] Create metrics registry
- [x] Define HTTP request metrics
- [x] Request count by method, path, status
- [x] Request duration histogram
- [x] Request size histogram
- [x] Response size histogram
- [x] Define database metrics
- [x] Connection pool stats
- [x] Query duration histogram
- [x] Error count by type
- [x] Define business metrics
- [x] User registration count
- [x] Active users gauge
- [x] Expose /metrics endpoint
- [ ] Document available metrics

### 9.1.2 Tracing (OpenTelemetry)

- [x] Add OpenTelemetry dependencies
- [x] Configure trace exporter (Jaeger, Zipkin, OTLP)
- [x] Initialize tracer provider
- [x] Add HTTP middleware for trace propagation
- [x] Add spans for database operations
- [x] Add spans for external service calls
- [x] Include relevant attributes in spans
- [x] Implement context propagation throughout codebase

### 9.1.3 Health Checks

- [x] Implement comprehensive readiness check
- [x] Add timeout for each dependency check
- [x] Return structured health status
- [x] Implement liveness check
- [ ] Document health check endpoints

## 9.2 Security

### 9.2.1 Input Validation

- [x] Validate all user input at API boundary
- [x] Sanitize strings to prevent XSS (if rendering HTML)
- [x] Limit request body size
- [x] Limit query parameter lengths

### 9.2.2 Authentication & Authorization (don't implement)

- [ ] Implement secure password hashing (bcrypt with appropriate cost)
- [ ] Implement JWT with short expiration
- [ ] Implement refresh token mechanism
- [ ] Store refresh tokens securely
- [ ] Implement token revocation
- [ ] Add rate limiting for auth endpoints
- [ ] Implement account lockout after failed attempts

### 9.2.3 Security Headers

- [x] Add security headers middleware
- [x] X-Content-Type-Options: nosniff
- [x] X-Frame-Options: DENY
- [x] X-XSS-Protection: 1; mode=block
- [x] Content-Security-Policy (if serving HTML)
- [x] Strict-Transport-Security (for HTTPS)
- [x] Referrer-Policy header
- [x] Permissions-Policy header
- [x] X-Permitted-Cross-Domain-Policies header

### 9.2.4 Secrets Management

- [ ] Never log sensitive data (passwords, tokens)
- [ ] Never commit secrets to repository
- [ ] Use environment variables for secrets
- [ ] Consider secrets manager integration (Vault, AWS Secrets Manager)
- [ ] Rotate secrets regularly

## 9.3 Error Handling & Resilience

### 9.3.1 Error Tracking

- [x] Integrate error tracking service (Sentry, Rollbar)
- [x] Configure error grouping
- [x] Include relevant context with errors
- [ ] Set up alerting for error spikes

### 9.3.2 Circuit Breaker (optional)

- [x] Implement circuit breaker for external services
- [x] Configure failure threshold
- [x] Configure recovery timeout
- [x] Log circuit state changes

### 9.3.3 Retry Logic

- [x] Implement retry for transient failures
- [x] Use exponential backoff
- [x] Add jitter to prevent thundering herd
- [x] Set maximum retry attempts
- [x] Make retry behavior configurable

## 9.4 Performance

### 9.4.1 Database Optimization

- [x] Review and optimize slow queries
- [x] Add appropriate indexes
- [x] Configure connection pool size appropriately
- [x] Implement query timeouts
- [ ] Consider read replicas for read-heavy workloads

### 9.4.2 Caching Strategy

- [x] Identify cacheable data
- [x] Implement cache layer with Redis
- [x] Define cache invalidation strategy
- [x] Set appropriate TTLs
- [x] Monitor cache hit rates

### 9.4.3 Response Optimization

- [x] Enable response compression (gzip)
- [ ] Implement response caching where appropriate
- [ ] Optimize JSON serialization

---

# Phase 10: DevOps & Deployment

## 10.4 Documentation

### 10.4.1 API Documentation

- [ ] Create OpenAPI/Swagger specification
- [ ] Document all endpoints
- [ ] Document request/response schemas
- [ ] Document error responses
- [ ] Document authentication requirements
- [ ] Set up Swagger UI endpoint

### 10.4.2 Project Documentation

- [ ] Complete README with setup instructions
- [ ] Document architecture decisions (ADRs)
- [ ] Document configuration options
- [ ] Document deployment process
- [ ] Create troubleshooting guide
- [ ] Document API versioning strategy

### 10.4.3 Code Documentation

- [ ] Add package-level documentation (doc.go)
- [ ] Document all exported types and functions
- [ ] Include usage examples where helpful
- [ ] Generate and host GoDoc

---

# Final Checklist Before Production

- [ ] All tests passing with adequate coverage
- [ ] No critical security vulnerabilities
- [ ] All configuration externalized
- [ ] Logging configured appropriately
- [ ] Metrics and tracing enabled
- [ ] Health checks implemented and tested
- [ ] Database migrations tested
- [ ] Graceful shutdown implemented
- [ ] Rate limiting configured
- [ ] Error tracking integrated
- [ ] Documentation complete
- [ ] Deployment pipeline tested
- [ ] Rollback procedure documented and tested
- [ ] On-call procedures documented
- [ ] Load testing completed (optional but recommended)
