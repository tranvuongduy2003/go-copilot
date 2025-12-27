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

- [x] Install golang-migrate CLI: `go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`
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

- [ ] Define Command marker interface (can be empty)
- [ ] Define CommandHandler interface with Handle(ctx, cmd) error method
- [ ] Define CommandBus interface with Dispatch(ctx, cmd) error method
- [ ] Implement InMemoryCommandBus
- [ ] Store handlers by command type
- [ ] Implement Register() method to register handlers
- [ ] Implement Dispatch() method to find and execute handler
- [ ] Return error if handler not found for command type
- [ ] Add logging for command dispatch

### 5.1.2 Query Infrastructure

- [ ] Define Query marker interface (can be empty)
- [ ] Define QueryHandler interface with Handle(ctx, query) (result, error) method
- [ ] Use generics for type-safe query results if Go 1.18+
- [ ] Define QueryBus interface with Dispatch(ctx, query) (result, error) method
- [ ] Implement InMemoryQueryBus
- [ ] Store handlers by query type
- [ ] Implement Register() method
- [ ] Implement Dispatch() method
- [ ] Return error if handler not found

### 5.1.3 Middleware Support (Optional)

- [ ] Define CommandMiddleware function type
- [ ] Define QueryMiddleware function type
- [ ] Implement logging middleware for commands
- [ ] Implement logging middleware for queries
- [ ] Implement validation middleware
- [ ] Implement transaction middleware for commands
- [ ] Add middleware chain execution in buses

## 5.2 DTOs (`internal/application/dto/`)

### 5.2.1 User DTOs

- [ ] Define UserDTO struct for responses
- [ ] Include ID as string (UUID string representation)
- [ ] Include Email as string
- [ ] Include FullName as string
- [ ] Include Status as string
- [ ] Include CreatedAt as time.Time or string
- [ ] Include UpdatedAt as time.Time or string
- [ ] Exclude sensitive fields (password hash)
- [ ] Add JSON tags for serialization
- [ ] Implement ToDTO(user \*domain.User) UserDTO mapper function
- [ ] Implement ToDTOs(users []\*domain.User) []UserDTO mapper function

### 5.2.2 Pagination DTOs

- [ ] Define PaginationRequest struct with Page, Limit fields
- [ ] Add validation tags for minimum and maximum values
- [ ] Implement Offset() method to calculate SQL offset
- [ ] Define PaginatedResponse[T] generic struct
- [ ] Include Items []T field
- [ ] Include Total int64 field
- [ ] Include Page int field
- [ ] Include Limit int field
- [ ] Include TotalPages int field (calculated)
- [ ] Include HasNext bool field (calculated)
- [ ] Include HasPrev bool field (calculated)

## 5.3 User Commands (`internal/application/command/`)

### 5.3.1 CreateUser Command

- [ ] Define CreateUserCommand struct
- [ ] Include Email field with validation tags
- [ ] Include Password field with validation tags
- [ ] Include FullName field with validation tags
- [ ] Define CreateUserHandler struct
- [ ] Inject UserRepository dependency
- [ ] Inject PasswordHasher dependency (interface for bcrypt)
- [ ] Inject EventBus dependency
- [ ] Inject Logger dependency
- [ ] Implement Handle(ctx, cmd) (string, error) method returning created user ID
- [ ] Validate command input using validator
- [ ] Check if email already exists
- [ ] Return ErrEmailAlreadyExists if duplicate
- [ ] Hash password using injected hasher
- [ ] Create new User domain entity
- [ ] Save user using repository
- [ ] Publish domain events from aggregate
- [ ] Return created user ID
- [ ] Log successful creation
- [ ] Write unit tests with mocked repository
- [ ] Test successful creation flow
- [ ] Test duplicate email rejection
- [ ] Test validation failures
- [ ] Test password hashing integration

### 5.3.2 UpdateUser Command

- [ ] Define UpdateUserCommand struct
- [ ] Include UserID field (required)
- [ ] Include FullName field (optional, pointer or wrapper type)
- [ ] Include other updatable fields as optional
- [ ] Define UpdateUserHandler struct
- [ ] Inject dependencies
- [ ] Implement Handle(ctx, cmd) error method
- [ ] Validate command input
- [ ] Fetch existing user by ID
- [ ] Return ErrUserNotFound if not exists
- [ ] Apply updates to domain entity using business methods
- [ ] Save updated user
- [ ] Publish domain events
- [ ] Write unit tests for update scenarios

### 5.3.3 DeleteUser Command

- [ ] Define DeleteUserCommand struct with UserID field
- [ ] Define DeleteUserHandler struct
- [ ] Inject dependencies
- [ ] Implement Handle(ctx, cmd) error method
- [ ] Validate command input
- [ ] Fetch existing user by ID
- [ ] Return ErrUserNotFound if not exists
- [ ] Call Delete() method on domain entity
- [ ] Save updated user (soft delete)
- [ ] Publish domain events
- [ ] Write unit tests for delete scenarios

### 5.3.4 ChangePassword Command

- [ ] Define ChangePasswordCommand struct
- [ ] Include UserID field
- [ ] Include CurrentPassword field
- [ ] Include NewPassword field with validation
- [ ] Define ChangePasswordHandler struct
- [ ] Inject dependencies including PasswordHasher
- [ ] Implement Handle(ctx, cmd) error method
- [ ] Fetch existing user by ID
- [ ] Verify current password matches stored hash
- [ ] Return ErrInvalidPassword if mismatch
- [ ] Validate new password strength
- [ ] Hash new password
- [ ] Call ChangePassword() on domain entity
- [ ] Save updated user
- [ ] Publish domain events
- [ ] Write unit tests including password verification

### 5.3.5 ActivateUser Command

- [ ] Define ActivateUserCommand struct with UserID field
- [ ] Define ActivateUserHandler struct
- [ ] Inject dependencies
- [ ] Implement Handle(ctx, cmd) error method
- [ ] Fetch existing user by ID
- [ ] Call Activate() on domain entity
- [ ] Handle InvalidStatusTransition error from domain
- [ ] Save updated user
- [ ] Publish domain events
- [ ] Write unit tests

### 5.3.6 DeactivateUser Command

- [ ] Define DeactivateUserCommand struct with UserID field
- [ ] Define DeactivateUserHandler struct
- [ ] Implement Handle() method similar to ActivateUser
- [ ] Write unit tests

## 5.4 User Queries (`internal/application/query/`)

### 5.4.1 GetUser Query

- [ ] Define GetUserQuery struct with UserID field
- [ ] Define GetUserHandler struct
- [ ] Inject UserRepository dependency (can use same repo or dedicated read repo)
- [ ] Inject Logger dependency
- [ ] Implement Handle(ctx, query) (\*UserDTO, error) method
- [ ] Validate query input
- [ ] Fetch user from repository
- [ ] Return ErrUserNotFound if not exists
- [ ] Map domain entity to DTO
- [ ] Return DTO
- [ ] Write unit tests with mocked repository

### 5.4.2 ListUsers Query

- [ ] Define ListUsersQuery struct
- [ ] Include pagination fields (Page, Limit)
- [ ] Include filter fields (Status, Search)
- [ ] Include sort fields (SortBy, SortOrder)
- [ ] Define ListUsersHandler struct
- [ ] Inject dependencies
- [ ] Implement Handle(ctx, query) (\*PaginatedResponse[UserDTO], error) method
- [ ] Validate query input
- [ ] Apply default pagination if not specified
- [ ] Build filter from query fields
- [ ] Fetch users from repository with filter and pagination
- [ ] Map domain entities to DTOs
- [ ] Build paginated response with metadata
- [ ] Return response
- [ ] Write unit tests for various filter combinations

### 5.4.3 GetUserByEmail Query

- [ ] Define GetUserByEmailQuery struct with Email field
- [ ] Define GetUserByEmailHandler struct
- [ ] Implement Handle() method
- [ ] Write unit tests

### 5.4.4 CheckEmailExists Query

- [ ] Define CheckEmailExistsQuery struct with Email field
- [ ] Define CheckEmailExistsHandler struct
- [ ] Implement Handle(ctx, query) (bool, error) method
- [ ] Write unit tests

---

# Phase 6: Interface Layer Implementation (HTTP API)

## 6.1 HTTP Infrastructure (`internal/interfaces/http/`)

### 6.1.1 Server Setup

- [ ] Define Server struct with dependencies
- [ ] Include HTTP server instance
- [ ] Include router instance
- [ ] Include logger instance
- [ ] Include configuration
- [ ] Implement NewServer constructor
- [ ] Accept all dependencies via constructor injection
- [ ] Configure server timeouts from config (read, write, idle)
- [ ] Implement Start() method to start HTTP server
- [ ] Implement Shutdown(ctx) method for graceful shutdown
- [ ] Wait for in-flight requests to complete
- [ ] Respect context deadline for shutdown timeout

### 6.1.2 Router Setup (`internal/interfaces/http/router/`)

- [ ] Define NewRouter function
- [ ] Accept handler dependencies
- [ ] Create chi router instance
- [ ] Apply global middleware in correct order
- [ ] Define route groups with version prefix (/api/v1)
- [ ] Register user routes
- [ ] Register health check routes
- [ ] Return configured router

### 6.1.3 Response Helpers

- [ ] Define standard success response structure
- [ ] Include data field
- [ ] Include optional message field
- [ ] Include optional metadata field
- [ ] Define standard error response structure
- [ ] Include error code field
- [ ] Include message field
- [ ] Include details field (for validation errors)
- [ ] Include request ID field
- [ ] Implement JSON response helper function
- [ ] Set Content-Type header
- [ ] Set status code
- [ ] Encode response body
- [ ] Implement error response helper function
- [ ] Map domain errors to HTTP status codes
- [ ] NotFoundError → 404
- [ ] ValidationError → 400
- [ ] ConflictError → 409
- [ ] AuthorizationError → 403
- [ ] Unknown errors → 500
- [ ] Log errors appropriately (5xx with stack, 4xx without)

## 6.2 Middleware (`internal/interfaces/http/middleware/`)

### 6.2.1 Request ID Middleware

- [ ] Implement middleware function
- [ ] Extract X-Request-ID header if present
- [ ] Generate new UUID if header not present
- [ ] Inject request ID into request context
- [ ] Set X-Request-ID response header
- [ ] Write tests for middleware

### 6.2.2 Logging Middleware

- [ ] Implement middleware function
- [ ] Extract logger from context or use global
- [ ] Extract request ID from context
- [ ] Log request start with method, path, request ID
- [ ] Wrap response writer to capture status code
- [ ] Log request completion with duration and status
- [ ] Include additional fields: remote addr, user agent
- [ ] Skip logging for health check endpoints (optional)
- [ ] Write tests for middleware

### 6.2.3 Recovery Middleware

- [ ] Implement middleware function
- [ ] Use defer to catch panics
- [ ] Log panic with stack trace
- [ ] Return 500 error response to client
- [ ] Include request ID in error response
- [ ] Do not expose stack trace to client in production
- [ ] Write tests for panic recovery

### 6.2.4 Timeout Middleware

- [ ] Implement middleware function
- [ ] Accept timeout duration as parameter
- [ ] Create context with timeout
- [ ] Pass timeout context to next handler
- [ ] Handle context deadline exceeded
- [ ] Return 504 Gateway Timeout on deadline
- [ ] Write tests for timeout behavior

### 6.2.5 CORS Middleware

- [ ] Configure allowed origins from config
- [ ] Configure allowed methods
- [ ] Configure allowed headers
- [ ] Configure exposed headers
- [ ] Configure credentials support
- [ ] Configure max age for preflight cache
- [ ] Handle preflight OPTIONS requests
- [ ] Write tests for CORS behavior

### 6.2.6 Authentication Middleware (if required)

- [ ] Implement middleware function
- [ ] Extract Authorization header
- [ ] Validate Bearer token format
- [ ] Parse and validate JWT token
- [ ] Extract user claims from token
- [ ] Inject user info into request context
- [ ] Return 401 for missing or invalid token
- [ ] Write tests for auth scenarios

### 6.2.7 Rate Limiting Middleware (optional)

- [ ] Choose rate limiting strategy (token bucket, sliding window)
- [ ] Configure limits from config
- [ ] Implement per-IP rate limiting
- [ ] Implement per-user rate limiting (if authenticated)
- [ ] Return 429 Too Many Requests when limited
- [ ] Include Retry-After header
- [ ] Write tests for rate limiting

## 6.3 HTTP DTOs (`internal/interfaces/http/dto/`)

### 6.3.1 Request DTOs

- [ ] Define CreateUserRequest struct
- [ ] Include Email field with json tag and validation tags
- [ ] Include Password field with validation tags
- [ ] Include FullName field with validation tags
- [ ] Define UpdateUserRequest struct
- [ ] Include optional fields with pointer types or omitempty
- [ ] Define ChangePasswordRequest struct
- [ ] Include CurrentPassword field
- [ ] Include NewPassword field with validation
- [ ] Include ConfirmPassword field (optional, for UI validation)
- [ ] Define ListUsersRequest struct for query parameters
- [ ] Include page, limit as query params
- [ ] Include status filter as query param
- [ ] Include search as query param
- [ ] Include sort_by, sort_order as query params

### 6.3.2 Response DTOs

- [ ] Define UserResponse struct mirroring application DTO
- [ ] Add JSON tags with appropriate naming (snake_case or camelCase)
- [ ] Define PaginatedResponse struct for list endpoints
- [ ] Define ErrorResponse struct matching standard error format

### 6.3.3 DTO Validation

- [ ] Add validation tags to all request DTOs
- [ ] Document validation rules in comments or OpenAPI spec
- [ ] Create reusable validation error formatter

## 6.4 User Handlers (`internal/interfaces/http/handler/`)

### 6.4.1 Handler Structure

- [ ] Define UserHandler struct
- [ ] Inject CommandBus dependency
- [ ] Inject QueryBus dependency
- [ ] Inject Validator dependency
- [ ] Inject Logger dependency
- [ ] Implement NewUserHandler constructor

### 6.4.2 Create User Endpoint

- [ ] Implement CreateUser handler for POST /users
- [ ] Parse request body into CreateUserRequest DTO
- [ ] Handle JSON parsing errors with 400 response
- [ ] Validate request using validator
- [ ] Return 400 with validation details on failure
- [ ] Map HTTP DTO to CreateUserCommand
- [ ] Dispatch command via CommandBus
- [ ] Handle domain errors and map to HTTP responses
- [ ] Return 201 Created with user ID on success
- [ ] Include Location header with resource URL
- [ ] Write integration tests for endpoint

### 6.4.3 Get User Endpoint

- [ ] Implement GetUser handler for GET /users/{id}
- [ ] Extract user ID from URL path parameter
- [ ] Validate UUID format
- [ ] Return 400 for invalid UUID
- [ ] Create GetUserQuery
- [ ] Dispatch query via QueryBus
- [ ] Handle not found error with 404 response
- [ ] Return 200 with user data on success
- [ ] Write integration tests for endpoint

### 6.4.4 List Users Endpoint

- [ ] Implement ListUsers handler for GET /users
- [ ] Parse query parameters into ListUsersRequest
- [ ] Apply default pagination if not specified
- [ ] Validate query parameters
- [ ] Create ListUsersQuery
- [ ] Dispatch query via QueryBus
- [ ] Return 200 with paginated response
- [ ] Write integration tests for endpoint

### 6.4.5 Update User Endpoint

- [ ] Implement UpdateUser handler for PUT /users/{id}
- [ ] Extract user ID from URL path
- [ ] Parse request body
- [ ] Validate request
- [ ] Create UpdateUserCommand
- [ ] Dispatch command
- [ ] Handle not found and other errors
- [ ] Return 200 on success
- [ ] Write integration tests

### 6.4.6 Delete User Endpoint

- [ ] Implement DeleteUser handler for DELETE /users/{id}
- [ ] Extract user ID from URL path
- [ ] Create DeleteUserCommand
- [ ] Dispatch command
- [ ] Handle not found error
- [ ] Return 204 No Content on success
- [ ] Write integration tests

### 6.4.7 Change Password Endpoint

- [ ] Implement ChangePassword handler for POST /users/{id}/password
- [ ] Parse and validate request
- [ ] Create ChangePasswordCommand
- [ ] Dispatch command
- [ ] Handle invalid password error with 400
- [ ] Return 200 on success
- [ ] Write integration tests

### 6.4.8 Activate User Endpoint

- [ ] Implement ActivateUser handler for POST /users/{id}/activate
- [ ] Create ActivateUserCommand
- [ ] Dispatch command
- [ ] Handle status transition errors
- [ ] Return 200 on success
- [ ] Write integration tests

### 6.4.9 Deactivate User Endpoint

- [ ] Implement DeactivateUser handler for POST /users/{id}/deactivate
- [ ] Similar implementation to activate
- [ ] Write integration tests

## 6.5 Health Check Endpoints (`internal/interfaces/http/handler/`)

### 6.5.1 Liveness Probe

- [ ] Implement handler for GET /health/live
- [ ] Return 200 OK if application is running
- [ ] No dependency checks (just proves process is alive)
- [ ] Return simple JSON response with status

### 6.5.2 Readiness Probe

- [ ] Implement handler for GET /health/ready
- [ ] Check database connection health
- [ ] Check Redis connection health (if used)
- [ ] Check other critical dependencies
- [ ] Return 200 OK if all dependencies healthy
- [ ] Return 503 Service Unavailable if any dependency unhealthy
- [ ] Include details of which checks failed

---

# Phase 7: Application Entry Point & Dependency Wiring

## 7.1 Main Function (`cmd/api/main.go`)

### 7.1.1 Initialization Sequence

- [ ] Load configuration from environment and files
- [ ] Handle configuration errors with clear message and exit
- [ ] Initialize logger with configuration
- [ ] Set as global logger
- [ ] Log application startup with version info
- [ ] Initialize database connection
- [ ] Implement retry logic for initial connection
- [ ] Run database migrations if configured for auto-migrate
- [ ] Handle migration errors appropriately
- [ ] Initialize Redis connection (if used)
- [ ] Verify connectivity with ping

### 7.1.2 Dependency Construction

- [ ] Create repositories
- [ ] Create UserRepository with database pool
- [ ] Create services/utilities
- [ ] Create PasswordHasher implementation
- [ ] Create Validator instance
- [ ] Create EventBus implementation
- [ ] Create command handlers
- [ ] Wire all dependencies into handlers
- [ ] Create CommandBus and register handlers
- [ ] Create query handlers
- [ ] Wire dependencies into handlers
- [ ] Create QueryBus and register handlers
- [ ] Create HTTP handlers
- [ ] Wire CommandBus, QueryBus, and other dependencies
- [ ] Create router with handlers and middleware
- [ ] Create HTTP server with router

### 7.1.3 Server Startup

- [ ] Start HTTP server in goroutine
- [ ] Log server address and port
- [ ] Set up OS signal handling
- [ ] Listen for SIGINT and SIGTERM
- [ ] Implement graceful shutdown sequence
- [ ] Stop accepting new connections
- [ ] Wait for in-flight requests (with timeout)
- [ ] Close database connections
- [ ] Close Redis connections
- [ ] Flush logs
- [ ] Log shutdown completion
- [ ] Exit with appropriate code

### 7.1.4 Error Handling

- [ ] Handle server start errors
- [ ] Handle shutdown timeout
- [ ] Log all errors with context
- [ ] Return non-zero exit code on error

## 7.2 Dependency Injection (Optional - using Wire)

### 7.2.1 Wire Setup

- [ ] Install wire: `go install github.com/google/wire/cmd/wire@latest`
- [ ] Create `wire.go` file with build tag `//go:build wireinject`
- [ ] Define provider functions for each dependency
- [ ] Define provider sets for related providers
- [ ] Create injector function
- [ ] Run `wire ./...` to generate `wire_gen.go`
- [ ] Add `wire_gen.go` to `.gitignore` or commit it based on team preference

---

# Phase 8: Testing Strategy

## 8.1 Unit Tests

### 8.1.1 Domain Layer Tests

- [ ] Test all entity constructors with valid and invalid inputs
- [ ] Test all entity business methods
- [ ] Test all value object validation
- [ ] Test all value object behavior methods
- [ ] Test domain event generation
- [ ] Test aggregate invariant enforcement
- [ ] Achieve minimum 90% coverage for domain layer

### 8.1.2 Application Layer Tests

- [ ] Mock all repository interfaces
- [ ] Mock all external service interfaces
- [ ] Test command handlers with all scenarios
- [ ] Success paths
- [ ] Validation failures
- [ ] Domain errors (not found, conflict)
- [ ] Test query handlers with all scenarios
- [ ] Test DTO mapping functions
- [ ] Achieve minimum 80% coverage for application layer

### 8.1.3 Interface Layer Tests

- [ ] Test HTTP handlers with mocked buses
- [ ] Test request parsing and validation
- [ ] Test response formatting
- [ ] Test error response mapping
- [ ] Test middleware behavior
- [ ] Achieve minimum 70% coverage for interface layer

## 8.2 Integration Tests

### 8.2.1 Test Infrastructure

- [ ] Create `docker-compose.test.yml` for test dependencies
- [ ] Configure test database with separate schema or database
- [ ] Implement test database setup and teardown
- [ ] Implement test data fixtures/factories
- [ ] Create helper functions for common test operations

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

- [ ] Create test helper package
- [ ] Implement random data generators
- [ ] Implement assertion helpers
- [ ] Implement mock factories

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

- [ ] Add prometheus client library
- [ ] Create metrics registry
- [ ] Define HTTP request metrics
- [ ] Request count by method, path, status
- [ ] Request duration histogram
- [ ] Request size histogram
- [ ] Response size histogram
- [ ] Define database metrics
- [ ] Connection pool stats
- [ ] Query duration histogram
- [ ] Error count by type
- [ ] Define business metrics
- [ ] User registration count
- [ ] Active users gauge
- [ ] Expose /metrics endpoint
- [ ] Document available metrics

### 9.1.2 Tracing (OpenTelemetry)

- [ ] Add OpenTelemetry dependencies
- [ ] Configure trace exporter (Jaeger, Zipkin, OTLP)
- [ ] Initialize tracer provider
- [ ] Add HTTP middleware for trace propagation
- [ ] Add spans for database operations
- [ ] Add spans for external service calls
- [ ] Include relevant attributes in spans
- [ ] Implement context propagation throughout codebase

### 9.1.3 Health Checks

- [ ] Implement comprehensive readiness check
- [ ] Add timeout for each dependency check
- [ ] Return structured health status
- [ ] Implement liveness check
- [ ] Document health check endpoints

## 9.2 Security

### 9.2.1 Input Validation

- [ ] Validate all user input at API boundary
- [ ] Sanitize strings to prevent XSS (if rendering HTML)
- [ ] Limit request body size
- [ ] Limit query parameter lengths

### 9.2.2 Authentication & Authorization

- [ ] Implement secure password hashing (bcrypt with appropriate cost)
- [ ] Implement JWT with short expiration
- [ ] Implement refresh token mechanism
- [ ] Store refresh tokens securely
- [ ] Implement token revocation
- [ ] Add rate limiting for auth endpoints
- [ ] Implement account lockout after failed attempts

### 9.2.3 Security Headers

- [ ] Add security headers middleware
- [ ] X-Content-Type-Options: nosniff
- [ ] X-Frame-Options: DENY
- [ ] X-XSS-Protection: 1; mode=block
- [ ] Content-Security-Policy (if serving HTML)
- [ ] Strict-Transport-Security (for HTTPS)

### 9.2.4 Secrets Management

- [ ] Never log sensitive data (passwords, tokens)
- [ ] Never commit secrets to repository
- [ ] Use environment variables for secrets
- [ ] Consider secrets manager integration (Vault, AWS Secrets Manager)
- [ ] Rotate secrets regularly

## 9.3 Error Handling & Resilience

### 9.3.1 Error Tracking

- [ ] Integrate error tracking service (Sentry, Rollbar)
- [ ] Configure error grouping
- [ ] Include relevant context with errors
- [ ] Set up alerting for error spikes

### 9.3.2 Circuit Breaker (optional)

- [ ] Implement circuit breaker for external services
- [ ] Configure failure threshold
- [ ] Configure recovery timeout
- [ ] Log circuit state changes

### 9.3.3 Retry Logic

- [ ] Implement retry for transient failures
- [ ] Use exponential backoff
- [ ] Add jitter to prevent thundering herd
- [ ] Set maximum retry attempts
- [ ] Make retry behavior configurable

## 9.4 Performance

### 9.4.1 Database Optimization

- [ ] Review and optimize slow queries
- [ ] Add appropriate indexes
- [ ] Configure connection pool size appropriately
- [ ] Implement query timeouts
- [ ] Consider read replicas for read-heavy workloads

### 9.4.2 Caching Strategy

- [ ] Identify cacheable data
- [ ] Implement cache layer with Redis
- [ ] Define cache invalidation strategy
- [ ] Set appropriate TTLs
- [ ] Monitor cache hit rates

### 9.4.3 Response Optimization

- [ ] Enable response compression (gzip)
- [ ] Implement response caching where appropriate
- [ ] Optimize JSON serialization

---

# Phase 10: DevOps & Deployment

## 10.1 Docker

### 10.1.1 Dockerfile

- [ ] Use multi-stage build
- [ ] Stage 1: Build binary with Go image
- [ ] Stage 2: Runtime with minimal base image (distroless or alpine)
- [ ] Set appropriate labels
- [ ] Create non-root user for runtime
- [ ] Copy only necessary files
- [ ] Set appropriate EXPOSE port
- [ ] Define ENTRYPOINT and CMD

### 10.1.2 Docker Compose

- [ ] Define application service
- [ ] Define database service with health check
- [ ] Define Redis service (if used)
- [ ] Configure networking between services
- [ ] Define volumes for data persistence
- [ ] Create environment-specific compose files

### 10.1.3 Docker Optimization

- [ ] Use .dockerignore to reduce context size
- [ ] Order Dockerfile instructions for optimal caching
- [ ] Pin base image versions
- [ ] Scan images for vulnerabilities

## 10.2 CI/CD Pipeline

### 10.2.1 CI Pipeline

- [ ] Trigger on push and pull request
- [ ] Lint stage
- [ ] Run golangci-lint
- [ ] Fail on lint errors
- [ ] Test stage
- [ ] Run unit tests with coverage
- [ ] Run integration tests
- [ ] Upload coverage report
- [ ] Security stage
- [ ] Run gosec for security issues
- [ ] Run dependency vulnerability scan
- [ ] Build stage
- [ ] Build Docker image
- [ ] Tag with commit SHA and branch
- [ ] Push to container registry

### 10.2.2 CD Pipeline

- [ ] Deploy to staging on merge to develop
- [ ] Run smoke tests against staging
- [ ] Deploy to production on merge to main
- [ ] Implement deployment strategy (rolling, blue-green, canary)
- [ ] Implement rollback mechanism
- [ ] Send deployment notifications

## 10.3 Kubernetes (Optional)

### 10.3.1 Manifests

- [ ] Create Deployment manifest
- [ ] Define resource requests and limits
- [ ] Configure liveness and readiness probes
- [ ] Set appropriate replica count
- [ ] Create Service manifest
- [ ] Create ConfigMap for non-sensitive config
- [ ] Create Secret for sensitive config
- [ ] Create Ingress for external access
- [ ] Create HorizontalPodAutoscaler

### 10.3.2 Helm Chart (Optional)

- [ ] Create Helm chart structure
- [ ] Parameterize environment-specific values
- [ ] Create values files for each environment
- [ ] Document chart usage

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
