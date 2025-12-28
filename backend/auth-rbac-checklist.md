# Authentication & Authorization Checklist
## RBAC (Role-Based Access Control) | Expert-Level Implementation

---

# Phase 1: RBAC Domain Model Design

## 1.1 Permission Aggregate (`internal/domain/permission/`)

### 1.1.1 Permission Entity Design
- [x] Define Permission entity struct
- [x] Include ID as UUID
- [x] Include Resource field (e.g., "users", "orders", "products")
- [x] Include Action field (e.g., "create", "read", "update", "delete", "list")
- [x] Include Description field for human-readable explanation
- [x] Include IsSystem field to mark built-in permissions (non-deletable)
- [x] Include CreatedAt timestamp
- [x] Include UpdatedAt timestamp
- [x] Implement NewPermission constructor with validation
- [x] Validate Resource is not empty and follows naming convention
- [x] Validate Action is from allowed action set
- [x] Generate unique permission code combining resource:action
- [x] Implement permission code generation method (e.g., "users:create")
- [x] Implement equality comparison based on resource and action

### 1.1.2 Permission Value Objects
- [x] Define Resource value object
- [x] Validate format (lowercase, alphanumeric with underscores)
- [x] Define allowed resources enum or registry
- [x] Define Action value object
- [x] Define standard actions: create, read, update, delete, list, manage
- [x] Allow custom actions for specific resources
- [x] Define PermissionCode value object
- [x] Format: "{resource}:{action}"
- [x] Implement parsing from string
- [x] Implement validation

### 1.1.3 Permission Repository Interface
- [x] Define PermissionRepository interface
- [x] Define Create(ctx, permission) error
- [x] Define Update(ctx, permission) error
- [x] Define Delete(ctx, id) error
- [x] Define FindByID(ctx, id) (*Permission, error)
- [x] Define FindByCode(ctx, code) (*Permission, error)
- [x] Define FindByResource(ctx, resource) ([]*Permission, error)
- [x] Define FindAll(ctx) ([]*Permission, error)
- [x] Define ExistsByCode(ctx, code) (bool, error)

### 1.1.4 Permission Domain Errors
- [x] Define ErrPermissionNotFound
- [x] Define ErrPermissionCodeExists
- [x] Define ErrSystemPermissionCannotBeDeleted
- [x] Define ErrInvalidResource
- [x] Define ErrInvalidAction

## 1.2 Role Aggregate (`internal/domain/role/`)

### 1.2.1 Role Entity Design
- [x] Define Role entity struct
- [x] Include ID as UUID
- [x] Include Name field (unique, e.g., "admin", "manager", "user")
- [x] Include DisplayName field for UI
- [x] Include Description field
- [x] Include Permissions as slice of permission IDs or permission codes
- [x] Include IsSystem field to mark built-in roles (non-deletable)
- [x] Include IsDefault field to mark role assigned to new users
- [x] Include Priority field for role hierarchy (higher = more privileged)
- [x] Include CreatedAt timestamp
- [x] Include UpdatedAt timestamp
- [x] Embed AggregateRoot for domain events

### 1.2.2 Role Business Methods
- [x] Implement NewRole constructor with validation
- [x] Validate Name follows naming convention (lowercase, no spaces)
- [x] Validate DisplayName is not empty
- [x] Implement AddPermission(permissionID) method
- [x] Validate permission is not already assigned
- [x] Register RolePermissionAdded domain event
- [x] Implement RemovePermission(permissionID) method
- [x] Validate permission is currently assigned
- [x] Register RolePermissionRemoved domain event
- [x] Implement SetPermissions(permissionIDs) method for bulk update
- [x] Clear existing permissions
- [x] Add all new permissions
- [x] Register RolePermissionsUpdated domain event
- [x] Implement HasPermission(permissionID) method
- [x] Implement UpdateDetails(displayName, description) method
- [x] Implement CanBeDeleted() method checking IsSystem flag

### 1.2.3 Role Repository Interface
- [x] Define RoleRepository interface
- [x] Define Create(ctx, role) error
- [x] Define Update(ctx, role) error
- [x] Define Delete(ctx, id) error
- [x] Define FindByID(ctx, id) (*Role, error)
- [x] Define FindByName(ctx, name) (*Role, error)
- [x] Define FindByIDs(ctx, ids) ([]*Role, error)
- [x] Define FindAll(ctx) ([]*Role, error)
- [x] Define FindDefault(ctx) (*Role, error)
- [x] Define ExistsByName(ctx, name) (bool, error)
- [x] Define FindByPermission(ctx, permissionID) ([]*Role, error)

### 1.2.4 Role Domain Events
- [x] Define RoleCreatedEvent
- [x] Define RoleUpdatedEvent
- [x] Define RoleDeletedEvent
- [x] Define RolePermissionAddedEvent with role ID and permission ID
- [x] Define RolePermissionRemovedEvent
- [x] Define RolePermissionsUpdatedEvent with old and new permission lists

### 1.2.5 Role Domain Errors
- [x] Define ErrRoleNotFound
- [x] Define ErrRoleNameExists
- [x] Define ErrSystemRoleCannotBeDeleted
- [x] Define ErrSystemRoleCannotBeModified
- [x] Define ErrPermissionAlreadyAssigned
- [x] Define ErrPermissionNotAssigned
- [x] Define ErrDefaultRoleCannotBeDeleted

## 1.3 User-Role Association

### 1.3.1 Extend User Aggregate
- [x] Add Roles field to User entity as slice of role IDs
- [x] Implement AssignRole(roleID) method
- [x] Validate role is not already assigned
- [x] Register UserRoleAssigned domain event
- [x] Implement RevokeRole(roleID) method
- [x] Validate role is currently assigned
- [x] Validate user has at least one role after revocation (optional)
- [x] Register UserRoleRevoked domain event
- [x] Implement SetRoles(roleIDs) method for bulk update
- [x] Implement HasRole(roleID) method
- [x] Implement GetRoleIDs() method

### 1.3.2 User Domain Events for Roles
- [x] Define UserRoleAssignedEvent with user ID and role ID
- [x] Define UserRoleRevokedEvent with user ID and role ID
- [x] Define UserRolesUpdatedEvent with user ID and role list

### 1.3.3 User Repository Updates
- [x] Update UserRepository interface if needed
- [x] Add FindByRole(ctx, roleID) ([]*User, error) method
- [x] Add method to load user with roles eagerly

## 1.4 Default Roles & Permissions Seed

### 1.4.1 Define System Permissions
- [x] Define all CRUD permissions for each resource
- [x] users:create, users:read, users:update, users:delete, users:list
- [x] roles:create, roles:read, roles:update, roles:delete, roles:list
- [x] permissions:read, permissions:list
- [x] Define special permissions
- [x] users:manage (super permission for user management)
- [x] roles:assign (permission to assign roles to users)
- [x] system:admin (super admin permission)

### 1.4.2 Define System Roles
- [x] Define SuperAdmin role
- [x] Has all permissions
- [x] IsSystem = true
- [x] Highest priority
- [x] Define Admin role
- [x] Has user and role management permissions
- [x] IsSystem = true
- [x] Define Manager role (if needed)
- [x] Has read and limited update permissions
- [x] IsSystem = true
- [x] Define User role
- [x] Has basic read permissions for own data
- [x] IsSystem = true
- [x] IsDefault = true

### 1.4.3 Seeder Implementation
- [x] Create database seeder for permissions
- [x] Check if permission exists before creating
- [x] Mark as IsSystem = true
- [x] Create database seeder for roles
- [x] Check if role exists before creating
- [x] Assign appropriate permissions to each role
- [x] Mark as IsSystem = true
- [ ] Create seeder for initial SuperAdmin user (optional)
- [x] Integrate seeder with application startup or migration

---

# Phase 2: Authentication System

## 2.1 Authentication Domain (`internal/domain/auth/`)

### 2.1.1 Credential Value Objects
- [x] Define Password value object
- [x] Implement minimum length validation (e.g., 8 characters)
- [x] Implement complexity validation (uppercase, lowercase, number, special char)
- [x] Implement maximum length validation
- [ ] Implement common password check (optional, use dictionary)
- [x] Define HashedPassword value object
- [x] Store algorithm identifier with hash
- [x] Implement Verify(plainPassword) method
- [x] Define Email value object (if not already in shared)
- [x] Validate format
- [x] Normalize (lowercase, trim)

### 2.1.2 Token Value Objects
- [x] Define AccessToken value object
- [x] Include token string
- [x] Include expiration time
- [x] Include token type (Bearer)
- [x] Implement IsExpired() method
- [x] Define RefreshToken value object
- [x] Include token string (opaque or JWT)
- [x] Include expiration time
- [x] Include user ID association
- [x] Include device/session identifier (optional)
- [x] Implement IsExpired() method
- [x] Define TokenPair struct containing both tokens

### 2.1.3 Session/RefreshToken Entity (if storing refresh tokens)
- [x] Define RefreshTokenEntity struct
- [x] Include ID as UUID
- [x] Include UserID
- [x] Include TokenHash (store hash, not plain token)
- [x] Include ExpiresAt
- [x] Include CreatedAt
- [x] Include LastUsedAt
- [x] Include DeviceInfo (user agent, IP - optional)
- [x] Include IsRevoked flag
- [x] Implement IsValid() method checking expiration and revocation
- [x] Implement Revoke() method
- [x] Implement UpdateLastUsed() method

### 2.1.4 RefreshToken Repository Interface
- [x] Define RefreshTokenRepository interface
- [x] Define Create(ctx, token) error
- [x] Define FindByTokenHash(ctx, hash) (*RefreshToken, error)
- [x] Define FindByUserID(ctx, userID) ([]*RefreshToken, error)
- [x] Define Revoke(ctx, id) error
- [x] Define RevokeAllByUserID(ctx, userID) error
- [x] Define DeleteExpired(ctx) (int64, error) for cleanup
- [x] Define CountActiveByUserID(ctx, userID) (int, error) for session limit

### 2.1.5 Authentication Domain Services
- [x] Define PasswordHasher interface
- [x] Define Hash(password string) (HashedPassword, error)
- [x] Define Verify(hashed HashedPassword, plain string) (bool, error)
- [x] Define TokenGenerator interface
- [x] Define GenerateAccessToken(user, roles, permissions) (AccessToken, error)
- [x] Define GenerateRefreshToken() (string, error)
- [x] Define ParseAccessToken(token string) (*Claims, error)
- [x] Define Claims struct for token payload
- [x] Include UserID
- [x] Include Email
- [x] Include Roles (role names or IDs)
- [x] Include Permissions (permission codes)
- [x] Include IssuedAt
- [x] Include ExpiresAt
- [x] Include TokenID (jti) for revocation

### 2.1.6 Authentication Domain Errors
- [x] Define ErrInvalidCredentials
- [x] Define ErrAccountLocked
- [x] Define ErrAccountInactive
- [x] Define ErrTokenExpired
- [x] Define ErrTokenInvalid
- [x] Define ErrTokenRevoked
- [x] Define ErrRefreshTokenNotFound
- [x] Define ErrRefreshTokenExpired
- [x] Define ErrSessionLimitExceeded
- [x] Define ErrPasswordTooWeak

## 2.2 Authentication Infrastructure

### 2.2.1 Password Hasher Implementation (`internal/infrastructure/security/`)
- [x] Implement BcryptPasswordHasher
- [x] Configure cost factor (recommend 12+ for production)
- [x] Make cost factor configurable
- [x] Implement Hash method using bcrypt.GenerateFromPassword
- [x] Implement Verify method using bcrypt.CompareHashAndPassword
- [x] Handle timing attacks (constant time comparison built into bcrypt)
- [ ] Consider implementing Argon2 hasher as alternative
- [ ] Configure memory, iterations, parallelism
- [x] Write unit tests for hasher

### 2.2.2 JWT Token Implementation (`internal/infrastructure/security/`)
- [x] Implement JWTTokenGenerator
- [x] Configure signing method (RS256 recommended for production, HS256 for simplicity)
- [x] Configure secret key or key pair from config
- [x] Configure access token expiration (e.g., 15 minutes)
- [x] Configure refresh token expiration (e.g., 7 days)
- [x] Configure issuer claim
- [x] Configure audience claim
- [x] Implement GenerateAccessToken method
- [x] Create claims with user info, roles, permissions
- [x] Include standard claims (iss, aud, exp, iat, jti)
- [x] Sign token with configured method
- [x] Return AccessToken value object
- [x] Implement GenerateRefreshToken method
- [x] Generate cryptographically secure random string
- [ ] Or generate JWT with minimal claims
- [x] Implement ParseAccessToken method
- [x] Parse and validate token signature
- [x] Validate expiration
- [x] Validate issuer and audience
- [x] Return Claims struct
- [x] Handle all error cases with appropriate error types
- [x] Write unit tests for token generation and parsing

### 2.2.3 RefreshToken Repository Implementation
- [x] Create refresh_tokens table migration
- [x] Include id, user_id, token_hash, expires_at, created_at, last_used_at, is_revoked
- [x] Add index on token_hash for lookup
- [x] Add index on user_id for user's sessions
- [x] Add index on expires_at for cleanup job
- [x] Implement PostgresRefreshTokenRepository
- [x] Implement all interface methods
- [x] Hash token before storing (use SHA256)
- [ ] Write integration tests

### 2.2.4 Token Blacklist (Optional - for access token revocation)
- [x] Design token blacklist strategy
- [x] Option 1: Store revoked token IDs in Redis with TTL
- [ ] Option 2: Store in database with cleanup job
- [ ] Option 3: Use short-lived tokens and rely on refresh token revocation
- [x] Implement TokenBlacklist interface
- [x] Define Add(tokenID, expiresAt) error
- [x] Define IsBlacklisted(tokenID) (bool, error)
- [x] Implement Redis-based blacklist
- [x] Use token ID as key
- [x] Set TTL to match token expiration
- [x] Write tests for blacklist

## 2.3 Authentication Application Layer

### 2.3.1 Auth Commands (`internal/application/command/auth/`)

#### Register Command
- [x] Define RegisterCommand struct
- [x] Include Email
- [x] Include Password
- [x] Include FullName
- [x] Include other registration fields
- [x] Define RegisterHandler struct
- [x] Inject UserRepository
- [x] Inject RoleRepository
- [x] Inject PasswordHasher
- [x] Inject EventBus
- [x] Inject Logger
- [x] Implement Handle method
- [x] Validate input
- [x] Check email uniqueness
- [x] Hash password
- [x] Create User entity
- [x] Assign default role to user
- [x] Save user
- [x] Publish domain events
- [x] Return created user ID
- [x] Write unit tests

#### Login Command
- [x] Define LoginCommand struct
- [x] Include Email
- [x] Include Password
- [x] Include DeviceInfo (optional)
- [x] Include IP address (optional)
- [x] Define LoginHandler struct
- [x] Inject UserRepository
- [x] Inject RoleRepository
- [x] Inject PermissionRepository
- [x] Inject RefreshTokenRepository
- [x] Inject PasswordHasher
- [x] Inject TokenGenerator
- [x] Inject EventBus
- [x] Inject Logger
- [x] Implement Handle method
- [x] Validate input
- [x] Find user by email
- [x] Return ErrInvalidCredentials if not found (don't reveal user exists)
- [x] Check user status (active, not banned)
- [x] Verify password
- [x] Implement account lockout on repeated failures (optional)
- [x] Load user's roles
- [x] Load permissions for all roles
- [x] Aggregate unique permissions
- [x] Generate access token with roles and permissions
- [x] Generate refresh token
- [x] Store refresh token entity (hashed)
- [ ] Check session limit before creating new session (optional)
- [x] Publish UserLoggedIn event
- [x] Return TokenPair
- [x] Write unit tests for all scenarios

#### Refresh Token Command
- [x] Define RefreshTokenCommand struct
- [x] Include RefreshToken string
- [x] Define RefreshTokenHandler struct
- [x] Inject RefreshTokenRepository
- [x] Inject UserRepository
- [x] Inject RoleRepository
- [x] Inject PermissionRepository
- [x] Inject TokenGenerator
- [x] Implement Handle method
- [x] Hash incoming refresh token
- [x] Find refresh token entity by hash
- [x] Return ErrRefreshTokenNotFound if not found
- [x] Validate refresh token is not expired
- [x] Validate refresh token is not revoked
- [x] Load user by ID from refresh token
- [x] Validate user is still active
- [x] Load current roles and permissions (may have changed)
- [x] Generate new access token
- [x] Optionally rotate refresh token (generate new, revoke old)
- [x] Update last_used_at on refresh token
- [x] Return new TokenPair or just AccessToken
- [x] Write unit tests

#### Logout Command
- [x] Define LogoutCommand struct
- [x] Include RefreshToken string (or AccessToken to extract jti)
- [x] Include LogoutAll flag (optional)
- [x] Define LogoutHandler struct
- [x] Inject RefreshTokenRepository
- [x] Inject TokenBlacklist (if using)
- [x] Implement Handle method
- [x] If LogoutAll: revoke all user's refresh tokens
- [x] If single logout: revoke specific refresh token
- [x] Optionally add access token to blacklist
- [x] Publish UserLoggedOut event
- [x] Write unit tests

#### Change Password Command (Auth context)
- [x] Ensure existing ChangePasswordCommand revokes all refresh tokens
- [x] Add step to revoke all sessions after password change
- [x] Publish PasswordChanged event for audit

#### Forgot Password Command
- [x] Define ForgotPasswordCommand struct
- [x] Include Email
- [x] Define ForgotPasswordHandler struct
- [x] Inject UserRepository
- [x] Inject PasswordResetTokenRepository (or use same as refresh token)
- [ ] Inject EmailService (interface)
- [x] Inject TokenGenerator
- [x] Implement Handle method
- [x] Find user by email
- [x] If not found, still return success (prevent email enumeration)
- [x] Generate password reset token (short-lived, e.g., 1 hour)
- [x] Store reset token with user ID
- [x] Send email with reset link (async via event)
- [x] Publish PasswordResetRequested event
- [x] Write unit tests

#### Reset Password Command
- [x] Define ResetPasswordCommand struct
- [x] Include ResetToken
- [x] Include NewPassword
- [x] Define ResetPasswordHandler struct
- [x] Inject UserRepository
- [x] Inject PasswordResetTokenRepository
- [x] Inject PasswordHasher
- [x] Inject RefreshTokenRepository
- [x] Implement Handle method
- [x] Find and validate reset token
- [x] Load user
- [x] Validate new password strength
- [x] Hash new password
- [x] Update user password
- [x] Invalidate reset token
- [x] Revoke all existing refresh tokens
- [x] Publish PasswordReset event
- [x] Write unit tests

### 2.3.2 Auth Queries (`internal/application/query/auth/`)

#### Get Current User Query
- [x] Define GetCurrentUserQuery struct
- [x] Include UserID (extracted from token)
- [x] Define GetCurrentUserHandler struct
- [x] Implement Handle method returning UserDTO with roles

#### Get User Sessions Query
- [x] Define GetUserSessionsQuery struct
- [x] Include UserID
- [x] Define GetUserSessionsHandler struct
- [x] Implement Handle method returning list of active sessions
- [x] Include device info, created at, last used at

### 2.3.3 Auth DTOs (`internal/application/dto/`)
- [x] Define TokenPairDTO
- [x] Include AccessToken
- [x] Include RefreshToken
- [x] Include ExpiresIn (seconds)
- [x] Include TokenType ("Bearer")
- [x] Define AuthUserDTO
- [x] Include user info
- [x] Include roles
- [x] Include permissions
- [x] Define SessionDTO for session listing

## 2.4 Authentication Interface Layer

### 2.4.1 Auth HTTP Handlers (`internal/interfaces/http/handler/`)

#### Auth Handler Structure
- [x] Define AuthHandler struct
- [x] Inject CommandBus
- [x] Inject QueryBus
- [x] Inject Validator
- [x] Inject Logger

#### Register Endpoint
- [x] Implement POST /auth/register
- [x] Parse RegisterRequest DTO
- [x] Validate input
- [x] Dispatch RegisterCommand
- [x] Return 201 with user ID or 200 with token pair (auto-login)

#### Login Endpoint
- [x] Implement POST /auth/login
- [x] Parse LoginRequest DTO (email, password)
- [x] Extract device info from headers (User-Agent)
- [x] Extract IP from request
- [x] Dispatch LoginCommand
- [x] Return 200 with TokenPairDTO
- [x] Return 401 for invalid credentials
- [x] Return 403 for locked/inactive account

#### Refresh Token Endpoint
- [x] Implement POST /auth/refresh
- [x] Parse RefreshTokenRequest DTO
- [ ] Or extract from HTTP-only cookie
- [x] Dispatch RefreshTokenCommand
- [x] Return 200 with new tokens
- [x] Return 401 for invalid/expired refresh token

#### Logout Endpoint
- [x] Implement POST /auth/logout
- [x] Extract refresh token from request or cookie
- [x] Dispatch LogoutCommand
- [ ] Clear HTTP-only cookie if using cookies
- [x] Return 204 No Content

#### Logout All Sessions Endpoint
- [x] Implement POST /auth/logout-all
- [x] Require authentication
- [x] Extract user ID from token
- [x] Dispatch LogoutCommand with LogoutAll flag
- [x] Return 204 No Content

#### Forgot Password Endpoint
- [x] Implement POST /auth/forgot-password
- [x] Parse ForgotPasswordRequest (email only)
- [x] Dispatch ForgotPasswordCommand
- [x] Always return 200 (prevent email enumeration)
- [x] Include message about checking email

#### Reset Password Endpoint
- [x] Implement POST /auth/reset-password
- [x] Parse ResetPasswordRequest (token, new password)
- [x] Dispatch ResetPasswordCommand
- [x] Return 200 on success
- [x] Return 400 for invalid/expired token

#### Get Current User Endpoint
- [x] Implement GET /auth/me
- [x] Require authentication
- [x] Extract user ID from token claims
- [x] Dispatch GetCurrentUserQuery
- [x] Return 200 with AuthUserDTO

#### Get Sessions Endpoint
- [x] Implement GET /auth/sessions
- [x] Require authentication
- [x] Dispatch GetUserSessionsQuery
- [x] Return 200 with list of sessions

#### Revoke Session Endpoint
- [x] Implement DELETE /auth/sessions/{sessionId}
- [x] Require authentication
- [x] Verify session belongs to current user
- [x] Dispatch command to revoke specific session
- [x] Return 204 No Content

### 2.4.2 Auth HTTP DTOs (`internal/interfaces/http/dto/`)
- [x] Define RegisterRequest
- [x] Define LoginRequest
- [x] Define RefreshTokenRequest
- [x] Define ForgotPasswordRequest
- [x] Define ResetPasswordRequest
- [x] Define TokenResponse matching TokenPairDTO
- [x] Add appropriate validation tags to all request DTOs

### 2.4.3 Auth Routes
- [x] Create auth route group /auth
- [x] Register all auth endpoints
- [x] Apply rate limiting to sensitive endpoints (login, register, forgot-password)
- [x] Apply authentication middleware to protected endpoints (me, sessions, logout)

---

# Phase 3: Authorization System (RBAC)

## 3.1 Authorization Infrastructure

### 3.1.1 Permission Checker Service
- [x] Define PermissionChecker interface
- [x] Define HasPermission(ctx, userID, permission) (bool, error)
- [x] Define HasAnyPermission(ctx, userID, permissions) (bool, error)
- [x] Define HasAllPermissions(ctx, userID, permissions) (bool, error)
- [x] Define GetUserPermissions(ctx, userID) ([]string, error)
- [x] Implement PermissionChecker
- [ ] Option 1: Load from database on each check
- [x] Option 2: Load from JWT claims (already in token)
- [ ] Option 3: Load from cache with invalidation
- [x] For JWT-based approach:
- [x] Extract permissions from token claims in context
- [x] Check if required permission exists in claims
- [ ] For database approach:
- [ ] Load user's roles
- [ ] Load permissions for each role
- [ ] Aggregate and check
- [ ] Implement caching layer for database approach
- [ ] Cache user permissions with TTL
- [ ] Invalidate on role assignment change
- [ ] Invalidate on role permission change

### 3.1.2 Role Checker Service
- [x] Define RoleChecker interface
- [x] Define HasRole(ctx, userID, role) (bool, error)
- [x] Define HasAnyRole(ctx, userID, roles) (bool, error)
- [x] Define GetUserRoles(ctx, userID) ([]string, error)
- [x] Implement similar to PermissionChecker

### 3.1.3 Authorization Context
- [x] Define AuthContext struct to hold authenticated user info
- [x] Include UserID
- [x] Include Email
- [x] Include Roles (slice of role names)
- [x] Include Permissions (slice of permission codes)
- [x] Define context key for AuthContext
- [x] Implement helper to get AuthContext from context
- [x] Implement helper to set AuthContext in context
- [x] Return error or nil if not authenticated

## 3.2 Authorization Middleware

### 3.2.1 Authentication Middleware Enhancement
- [x] Extract and parse JWT from Authorization header
- [x] Validate token signature and expiration
- [x] Check token blacklist (if implemented)
- [x] Create AuthContext from token claims
- [x] Inject AuthContext into request context
- [x] Continue to next handler if valid
- [x] Return 401 if token missing or invalid
- [x] Log authentication attempts

### 3.2.2 Require Auth Middleware
- [x] Implement middleware that requires authentication
- [x] Check if AuthContext exists in context
- [x] Return 401 if not authenticated
- [x] Continue if authenticated
- [x] Apply to routes that require any authenticated user

### 3.2.3 Require Permission Middleware
- [x] Implement RequirePermission(permission string) middleware factory
- [x] Extract AuthContext from context
- [x] Return 401 if not authenticated
- [x] Check if user has required permission
- [x] Return 403 Forbidden if permission denied
- [x] Log authorization failures
- [x] Continue if authorized
- [x] Implement RequireAnyPermission(permissions ...string) variant
- [x] Implement RequireAllPermissions(permissions ...string) variant

### 3.2.4 Require Role Middleware
- [x] Implement RequireRole(role string) middleware factory
- [x] Similar logic to RequirePermission
- [x] Return 403 if role not present
- [x] Implement RequireAnyRole(roles ...string) variant

### 3.2.5 Resource Owner Middleware
- [x] Implement middleware for resource ownership check
- [x] Extract resource ID from URL parameter
- [x] Extract user ID from AuthContext
- [x] Compare resource owner with current user
- [x] Allow if user is owner OR has admin permission
- [x] Return 403 if not owner and not admin
- [x] Make configurable for different resources

## 3.3 Permission Application Layer

### 3.3.1 Permission Commands (`internal/application/command/permission/`)

#### Create Permission Command
- [x] Define CreatePermissionCommand
- [x] Include Resource
- [x] Include Action
- [x] Include Description
- [x] Define CreatePermissionHandler
- [x] Validate input
- [x] Check if permission code already exists
- [x] Create Permission entity
- [x] Save to repository
- [x] Return permission ID
- [x] Write unit tests

#### Update Permission Command
- [x] Define UpdatePermissionCommand
- [x] Include PermissionID
- [x] Include Description (only description should be updatable)
- [x] Define UpdatePermissionHandler
- [x] Validate system permissions cannot be modified
- [x] Update and save
- [x] Write unit tests

#### Delete Permission Command
- [x] Define DeletePermissionCommand
- [x] Include PermissionID
- [x] Define DeletePermissionHandler
- [x] Validate system permissions cannot be deleted
- [x] Check if permission is assigned to any role
- [x] Option: Prevent deletion if in use
- [ ] Option: Cascade remove from roles
- [x] Delete permission
- [x] Write unit tests

### 3.3.2 Permission Queries (`internal/application/query/permission/`)
- [x] Implement ListPermissions query
- [x] Implement GetPermission query
- [x] Implement ListPermissionsByResource query
- [x] Implement GetPermissionsForRole query

## 3.4 Role Application Layer

### 3.4.1 Role Commands (`internal/application/command/role/`)

#### Create Role Command
- [x] Define CreateRoleCommand
- [x] Include Name
- [x] Include DisplayName
- [x] Include Description
- [x] Include PermissionIDs (initial permissions)
- [x] Define CreateRoleHandler
- [x] Validate input
- [x] Check name uniqueness
- [x] Validate all permission IDs exist
- [x] Create Role entity with permissions
- [x] Save to repository
- [x] Publish RoleCreated event
- [x] Return role ID
- [x] Write unit tests

#### Update Role Command
- [x] Define UpdateRoleCommand
- [x] Include RoleID
- [x] Include DisplayName
- [x] Include Description
- [x] Define UpdateRoleHandler
- [x] Validate system roles have limited modification
- [x] Update allowed fields
- [x] Save to repository
- [x] Publish RoleUpdated event
- [x] Write unit tests

#### Delete Role Command
- [x] Define DeleteRoleCommand
- [x] Include RoleID
- [x] Define DeleteRoleHandler
- [x] Validate not a system role
- [x] Validate not the default role
- [x] Check if role is assigned to any users
- [x] Option: Prevent deletion if in use
- [ ] Option: Reassign users to default role
- [x] Delete role
- [x] Publish RoleDeleted event
- [x] Write unit tests

#### Assign Permission to Role Command
- [x] Define AssignPermissionToRoleCommand
- [x] Include RoleID
- [x] Include PermissionID
- [x] Define AssignPermissionToRoleHandler
- [x] Validate role exists
- [x] Validate permission exists
- [x] Load role
- [x] Add permission using domain method
- [x] Save role
- [ ] Invalidate permission cache for affected users
- [x] Publish domain events
- [x] Write unit tests

#### Remove Permission from Role Command
- [x] Define RemovePermissionFromRoleCommand
- [x] Similar to assign but removes
- [x] Validate system role restrictions
- [x] Write unit tests

#### Set Role Permissions Command (bulk update)
- [x] Define SetRolePermissionsCommand
- [x] Include RoleID
- [x] Include PermissionIDs (complete list)
- [x] Define SetRolePermissionsHandler
- [x] Validate all permissions exist
- [x] Load role
- [x] Set permissions using domain method
- [x] Save role
- [ ] Invalidate caches
- [x] Publish events
- [x] Write unit tests

### 3.4.2 Role Queries (`internal/application/query/role/`)
- [x] Implement ListRoles query
- [x] Implement GetRole query with permissions
- [x] Implement GetRolePermissions query
- [x] Implement GetUsersWithRole query

## 3.5 User Role Management

### 3.5.1 User Role Commands (`internal/application/command/user/`)

#### Assign Role to User Command
- [x] Define AssignRoleToUserCommand
- [x] Include UserID
- [x] Include RoleID
- [x] Define AssignRoleToUserHandler
- [x] Validate user exists
- [x] Validate role exists
- [x] Load user
- [x] Assign role using domain method
- [x] Save user
- [ ] Invalidate user's permission cache
- [x] Publish UserRoleAssigned event
- [x] Write unit tests

#### Revoke Role from User Command
- [x] Define RevokeRoleFromUserCommand
- [x] Include UserID
- [x] Include RoleID
- [x] Define RevokeRoleFromUserHandler
- [x] Validate user has at least one role remaining (optional)
- [x] Load user
- [x] Revoke role using domain method
- [x] Save user
- [ ] Invalidate cache
- [x] Publish UserRoleRevoked event
- [x] Write unit tests

#### Set User Roles Command (bulk update)
- [x] Define SetUserRolesCommand
- [x] Include UserID
- [x] Include RoleIDs
- [x] Define SetUserRolesHandler
- [x] Validate at least one role (optional)
- [x] Validate all roles exist
- [x] Load user
- [x] Set roles
- [x] Save user
- [ ] Invalidate cache
- [x] Publish event
- [x] Write unit tests

### 3.5.2 User Role Queries
- [x] Implement GetUserRoles query
- [x] Implement GetUserPermissions query (aggregated from all roles)

## 3.6 Role & Permission HTTP Handlers

### 3.6.1 Permission Handler
- [x] Implement GET /permissions - List all permissions
- [x] Require permissions:list permission
- [x] Implement GET /permissions/{id} - Get permission details
- [x] Require permissions:read permission
- [x] Implement POST /permissions - Create permission (if allowing custom)
- [x] Require permissions:create permission
- [x] Implement PUT /permissions/{id} - Update permission
- [x] Require permissions:update permission
- [x] Implement DELETE /permissions/{id} - Delete permission
- [x] Require permissions:delete permission

### 3.6.2 Role Handler
- [x] Implement GET /roles - List all roles
- [x] Require roles:list permission
- [x] Implement GET /roles/{id} - Get role with permissions
- [x] Require roles:read permission
- [x] Implement POST /roles - Create role
- [x] Require roles:create permission
- [x] Implement PUT /roles/{id} - Update role details
- [x] Require roles:update permission
- [x] Implement DELETE /roles/{id} - Delete role
- [x] Require roles:delete permission
- [x] Implement PUT /roles/{id}/permissions - Set role permissions
- [x] Require roles:update permission
- [x] Implement POST /roles/{id}/permissions/{permissionId} - Add permission
- [x] Require roles:update permission
- [x] Implement DELETE /roles/{id}/permissions/{permissionId} - Remove permission
- [x] Require roles:update permission

### 3.6.3 User Role Management Endpoints
- [x] Implement GET /users/{id}/roles - Get user's roles
- [x] Require users:read or self
- [x] Implement PUT /users/{id}/roles - Set user's roles
- [x] Require roles:assign permission
- [x] Implement POST /users/{id}/roles/{roleId} - Assign role
- [x] Require roles:assign permission
- [x] Implement DELETE /users/{id}/roles/{roleId} - Revoke role
- [x] Require roles:assign permission
- [x] Implement GET /users/{id}/permissions - Get user's effective permissions
- [x] Require users:read or self

---

# Phase 4: Database Schema for RBAC

## 4.1 Permissions Table Migration
- [x] Generate migration: `goose create create_permissions_table sql`
- [x] Define permissions table
- [x] id UUID PRIMARY KEY
- [x] resource VARCHAR(100) NOT NULL
- [x] action VARCHAR(100) NOT NULL
- [x] description TEXT
- [x] is_system BOOLEAN DEFAULT FALSE
- [x] created_at TIMESTAMPTZ DEFAULT NOW()
- [x] updated_at TIMESTAMPTZ DEFAULT NOW()
- [x] Add unique constraint on (resource, action)
- [x] Add index on resource
- [x] Add index on is_system
- [x] Write down migration

## 4.2 Roles Table Migration
- [x] Generate migration: `goose create create_roles_table sql`
- [x] Define roles table
- [x] id UUID PRIMARY KEY
- [x] name VARCHAR(100) UNIQUE NOT NULL
- [x] display_name VARCHAR(255) NOT NULL
- [x] description TEXT
- [x] is_system BOOLEAN DEFAULT FALSE
- [x] is_default BOOLEAN DEFAULT FALSE
- [x] priority INTEGER DEFAULT 0
- [x] created_at TIMESTAMPTZ DEFAULT NOW()
- [x] updated_at TIMESTAMPTZ DEFAULT NOW()
- [x] Add index on name
- [x] Add index on is_default
- [x] Add constraint: only one default role
- [x] Write down migration

## 4.3 Role Permissions Junction Table Migration
- [x] Generate migration: `goose create create_role_permissions_table sql`
- [x] Define role_permissions table
- [x] role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE
- [x] permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE
- [x] created_at TIMESTAMPTZ DEFAULT NOW()
- [x] Add PRIMARY KEY (role_id, permission_id)
- [x] Add index on permission_id for reverse lookup
- [x] Write down migration

## 4.4 User Roles Junction Table Migration
- [x] Generate migration: `goose create create_user_roles_table sql`
- [x] Define user_roles table
- [x] user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
- [x] role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE
- [x] assigned_at TIMESTAMPTZ DEFAULT NOW()
- [x] assigned_by UUID REFERENCES users(id) (optional, for audit)
- [x] Add PRIMARY KEY (user_id, role_id)
- [x] Add index on role_id for reverse lookup
- [x] Write down migration

## 4.5 Refresh Tokens Table Migration
- [x] Generate migration: `goose create create_refresh_tokens_table sql`
- [x] Define refresh_tokens table
- [x] id UUID PRIMARY KEY
- [x] user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
- [x] token_hash VARCHAR(255) UNIQUE NOT NULL
- [x] expires_at TIMESTAMPTZ NOT NULL
- [x] created_at TIMESTAMPTZ DEFAULT NOW()
- [x] last_used_at TIMESTAMPTZ
- [x] is_revoked BOOLEAN DEFAULT FALSE
- [x] device_info JSONB (optional)
- [x] ip_address INET (optional)
- [x] Add index on user_id
- [x] Add index on expires_at for cleanup
- [x] Add index on is_revoked
- [x] Write down migration

## 4.6 Password Reset Tokens Table Migration (Optional)
- [x] Generate migration: `goose create create_password_reset_tokens_table sql`
- [x] Define password_reset_tokens table
- [x] id UUID PRIMARY KEY
- [x] user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
- [x] token_hash VARCHAR(255) UNIQUE NOT NULL
- [x] expires_at TIMESTAMPTZ NOT NULL
- [x] created_at TIMESTAMPTZ DEFAULT NOW()
- [x] used_at TIMESTAMPTZ
- [x] Add index on user_id
- [x] Add index on expires_at
- [x] Write down migration

## 4.7 Seed Data Migration
- [x] Generate migration: `goose create seed_permissions_and_roles go`
- [x] Use Go migration for complex seeding logic
- [x] Seed all system permissions
- [x] Seed all system roles
- [x] Assign permissions to roles
- [ ] Create initial super admin user (optional)
- [x] Make idempotent (check before insert)

---

# Phase 5: Repository Implementations

## 5.1 Permission Repository Implementation
- [x] Implement PostgresPermissionRepository
- [x] Implement Create with unique constraint handling
- [x] Implement Update
- [x] Implement Delete with system permission check
- [x] Implement FindByID
- [x] Implement FindByCode with (resource, action) lookup
- [x] Implement FindByResource
- [x] Implement FindAll with optional filtering
- [x] Implement ExistsByCode
- [ ] Write integration tests for all methods

## 5.2 Role Repository Implementation
- [x] Implement PostgresRoleRepository
- [x] Implement Create with permissions insertion
- [x] Implement Update
- [x] Implement Delete with cascade consideration
- [x] Implement FindByID with permissions loading
- [x] Implement FindByName
- [x] Implement FindByIDs for batch loading
- [x] Implement FindAll
- [x] Implement FindDefault
- [x] Implement ExistsByName
- [x] Implement FindByPermission
- [ ] Write integration tests

## 5.3 User Repository Updates
- [x] Update user repository to handle roles
- [x] Load roles when finding user (eager or lazy)
- [x] Update Create to assign default role
- [x] Implement FindByRole
- [ ] Update integration tests

## 5.4 RefreshToken Repository Implementation
- [x] Implement PostgresRefreshTokenRepository
- [x] Implement Create with token hashing
- [x] Implement FindByTokenHash
- [x] Implement FindByUserID
- [x] Implement Revoke
- [x] Implement RevokeAllByUserID
- [x] Implement DeleteExpired for cleanup job
- [x] Implement CountActiveByUserID
- [ ] Write integration tests

---

# Phase 6: Security Hardening

## 6.1 Rate Limiting
- [x] Implement rate limiting for login endpoint
- [x] Limit by IP address
- [x] Stricter limit for failed attempts
- [x] Implement rate limiting for registration
- [x] Implement rate limiting for password reset request
- [x] Implement rate limiting for token refresh
- [x] Configure rate limits via environment
- [x] Return 429 with Retry-After header
- [x] Log rate limit violations

## 6.2 Account Security
- [x] Implement account lockout after N failed login attempts
- [x] Track failed attempts in database or cache
- [x] Lock duration configurable (e.g., 15 minutes)
- [x] Reset counter on successful login
- [ ] Implement CAPTCHA integration for repeated failures (optional)
- [ ] Implement login notification email (optional)
- [ ] Detect suspicious login (new device, new location)
- [ ] Send email notification

## 6.3 Token Security
- [x] Use short-lived access tokens (15 minutes recommended)
- [x] Use longer-lived refresh tokens with rotation
- [x] Implement refresh token rotation
- [x] Issue new refresh token on each refresh
- [x] Revoke old refresh token
- [ ] Detect refresh token reuse (potential theft)
- [ ] Consider using HTTP-only cookies for tokens
- [ ] Set Secure flag for HTTPS only
- [ ] Set SameSite attribute
- [ ] Implement proper CSRF protection if using cookies

## 6.4 Password Security
- [x] Enforce minimum password length (12+ recommended)
- [x] Enforce complexity requirements
- [ ] Check against common password list
- [ ] Consider using zxcvbn for strength estimation
- [ ] Implement password history (prevent reuse of last N passwords)
- [ ] Implement password expiration policy (optional, controversial)
- [x] Hash passwords with appropriate cost factor
- [ ] Consider upgrading hash on login if algorithm changes

## 6.5 Session Management
- [ ] Implement maximum concurrent sessions per user
- [x] Allow users to view active sessions
- [x] Allow users to revoke specific sessions
- [x] Allow users to revoke all other sessions
- [ ] Implement session timeout for inactivity (optional)
- [ ] Implement absolute session timeout

## 6.6 Audit Logging
- [x] Log all authentication events
- [x] Successful login with IP, device info
- [x] Failed login with IP, device info
- [x] Logout
- [x] Password change
- [x] Password reset request and completion
- [x] Log all authorization events
- [x] Permission denied attempts
- [x] Role assignments and revocations
- [x] Store audit logs securely
- [ ] Implement log retention policy
- [x] Consider separate audit log storage

---

# Phase 7: Testing

## 7.1 Unit Tests

### 7.1.1 Domain Tests
- [x] Test Permission entity creation and validation
- [x] Test Role entity and all business methods
- [x] Test permission assignment/removal
- [x] Test status checks
- [x] Test User role methods
- [x] Test all value objects
- [x] Test all domain events

### 7.1.2 Application Tests
- [x] Test all auth command handlers
- [x] Mock repositories and services
- [x] Test success and error paths
- [x] Test all role/permission command handlers
- [x] Test all query handlers
- [x] Test authorization logic in handlers

### 7.1.3 Infrastructure Tests
- [x] Test password hasher
- [x] Test JWT token generator
- [x] Token generation and parsing
- [x] Expiration handling
- [x] Invalid token handling
- [x] Test permission checker service

## 7.2 Integration Tests

### 7.2.1 Repository Tests
- [ ] Test all permission repository methods with database
- [ ] Test all role repository methods with database
- [ ] Test role-permission relationships
- [ ] Test refresh token repository methods
- [ ] Test user-role relationships

### 7.2.2 API Tests
- [ ] Test complete registration flow
- [ ] Test complete login flow
- [ ] Test token refresh flow
- [ ] Test logout flows (single and all)
- [ ] Test password reset flow
- [ ] Test protected endpoint access with valid token
- [ ] Test protected endpoint rejection without token
- [ ] Test protected endpoint rejection with expired token
- [ ] Test authorization checks
- [ ] Access granted with correct permission
- [ ] Access denied without permission
- [ ] Test role management endpoints
- [ ] Test permission management endpoints
- [ ] Test user role assignment endpoints

## 7.3 Security Tests
- [x] Test rate limiting is enforced
- [x] Test account lockout works
- [x] Test password validation rejects weak passwords
- [x] Test token expiration is enforced
- [x] Test refresh token rotation works
- [x] Test revoked tokens are rejected
- [x] Test system roles/permissions cannot be deleted

---

# Phase 8: Documentation

## 8.1 API Documentation
- [x] Document all auth endpoints in OpenAPI spec
- [x] Include request/response schemas
- [x] Include error responses
- [x] Document authentication mechanism
- [x] Document all role/permission endpoints
- [x] Document required permissions for each endpoint
- [x] Include examples for common flows

## 8.2 Architecture Documentation
- [x] Document RBAC model and relationships
- [x] Create diagram showing User-Role-Permission relationships
- [x] Document authentication flow with sequence diagram
- [x] Document token refresh flow
- [x] Document authorization middleware chain
- [x] Document permission checking strategies
- [x] Document caching strategy for permissions

## 8.3 Operations Documentation
- [x] Document how to create new permissions
- [x] Document how to create new roles
- [x] Document how to assign roles to users
- [x] Document how to handle locked accounts
- [x] Document token cleanup job setup
- [x] Document audit log monitoring
- [x] Document incident response for security events

---

# Final Checklist

- [x] All RBAC entities implemented and tested
- [x] All authentication flows working
- [x] All authorization checks in place
- [x] Rate limiting configured
- [x] Account security measures implemented
- [x] Token security best practices followed
- [x] Audit logging enabled
- [x] All tests passing
- [x] API documentation complete
- [ ] Security review completed
- [ ] Performance testing for auth endpoints done
- [ ] Monitoring and alerting configured for auth failures
