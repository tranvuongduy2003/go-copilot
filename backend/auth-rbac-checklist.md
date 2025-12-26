# Authentication & Authorization Checklist
## RBAC (Role-Based Access Control) | Expert-Level Implementation

---

# Phase 1: RBAC Domain Model Design

## 1.1 Permission Aggregate (`internal/domain/permission/`)

### 1.1.1 Permission Entity Design
- [ ] Define Permission entity struct
- [ ] Include ID as UUID
- [ ] Include Resource field (e.g., "users", "orders", "products")
- [ ] Include Action field (e.g., "create", "read", "update", "delete", "list")
- [ ] Include Description field for human-readable explanation
- [ ] Include IsSystem field to mark built-in permissions (non-deletable)
- [ ] Include CreatedAt timestamp
- [ ] Include UpdatedAt timestamp
- [ ] Implement NewPermission constructor with validation
- [ ] Validate Resource is not empty and follows naming convention
- [ ] Validate Action is from allowed action set
- [ ] Generate unique permission code combining resource:action
- [ ] Implement permission code generation method (e.g., "users:create")
- [ ] Implement equality comparison based on resource and action

### 1.1.2 Permission Value Objects
- [ ] Define Resource value object
- [ ] Validate format (lowercase, alphanumeric with underscores)
- [ ] Define allowed resources enum or registry
- [ ] Define Action value object
- [ ] Define standard actions: create, read, update, delete, list, manage
- [ ] Allow custom actions for specific resources
- [ ] Define PermissionCode value object
- [ ] Format: "{resource}:{action}"
- [ ] Implement parsing from string
- [ ] Implement validation

### 1.1.3 Permission Repository Interface
- [ ] Define PermissionRepository interface
- [ ] Define Create(ctx, permission) error
- [ ] Define Update(ctx, permission) error
- [ ] Define Delete(ctx, id) error
- [ ] Define FindByID(ctx, id) (*Permission, error)
- [ ] Define FindByCode(ctx, code) (*Permission, error)
- [ ] Define FindByResource(ctx, resource) ([]*Permission, error)
- [ ] Define FindAll(ctx) ([]*Permission, error)
- [ ] Define ExistsByCode(ctx, code) (bool, error)

### 1.1.4 Permission Domain Errors
- [ ] Define ErrPermissionNotFound
- [ ] Define ErrPermissionCodeExists
- [ ] Define ErrSystemPermissionCannotBeDeleted
- [ ] Define ErrInvalidResource
- [ ] Define ErrInvalidAction

## 1.2 Role Aggregate (`internal/domain/role/`)

### 1.2.1 Role Entity Design
- [ ] Define Role entity struct
- [ ] Include ID as UUID
- [ ] Include Name field (unique, e.g., "admin", "manager", "user")
- [ ] Include DisplayName field for UI
- [ ] Include Description field
- [ ] Include Permissions as slice of permission IDs or permission codes
- [ ] Include IsSystem field to mark built-in roles (non-deletable)
- [ ] Include IsDefault field to mark role assigned to new users
- [ ] Include Priority field for role hierarchy (higher = more privileged)
- [ ] Include CreatedAt timestamp
- [ ] Include UpdatedAt timestamp
- [ ] Embed AggregateRoot for domain events

### 1.2.2 Role Business Methods
- [ ] Implement NewRole constructor with validation
- [ ] Validate Name follows naming convention (lowercase, no spaces)
- [ ] Validate DisplayName is not empty
- [ ] Implement AddPermission(permissionID) method
- [ ] Validate permission is not already assigned
- [ ] Register RolePermissionAdded domain event
- [ ] Implement RemovePermission(permissionID) method
- [ ] Validate permission is currently assigned
- [ ] Register RolePermissionRemoved domain event
- [ ] Implement SetPermissions(permissionIDs) method for bulk update
- [ ] Clear existing permissions
- [ ] Add all new permissions
- [ ] Register RolePermissionsUpdated domain event
- [ ] Implement HasPermission(permissionID) method
- [ ] Implement UpdateDetails(displayName, description) method
- [ ] Implement CanBeDeleted() method checking IsSystem flag

### 1.2.3 Role Repository Interface
- [ ] Define RoleRepository interface
- [ ] Define Create(ctx, role) error
- [ ] Define Update(ctx, role) error
- [ ] Define Delete(ctx, id) error
- [ ] Define FindByID(ctx, id) (*Role, error)
- [ ] Define FindByName(ctx, name) (*Role, error)
- [ ] Define FindByIDs(ctx, ids) ([]*Role, error)
- [ ] Define FindAll(ctx) ([]*Role, error)
- [ ] Define FindDefault(ctx) (*Role, error)
- [ ] Define ExistsByName(ctx, name) (bool, error)
- [ ] Define FindByPermission(ctx, permissionID) ([]*Role, error)

### 1.2.4 Role Domain Events
- [ ] Define RoleCreatedEvent
- [ ] Define RoleUpdatedEvent
- [ ] Define RoleDeletedEvent
- [ ] Define RolePermissionAddedEvent with role ID and permission ID
- [ ] Define RolePermissionRemovedEvent
- [ ] Define RolePermissionsUpdatedEvent with old and new permission lists

### 1.2.5 Role Domain Errors
- [ ] Define ErrRoleNotFound
- [ ] Define ErrRoleNameExists
- [ ] Define ErrSystemRoleCannotBeDeleted
- [ ] Define ErrSystemRoleCannotBeModified
- [ ] Define ErrPermissionAlreadyAssigned
- [ ] Define ErrPermissionNotAssigned
- [ ] Define ErrDefaultRoleCannotBeDeleted

## 1.3 User-Role Association

### 1.3.1 Extend User Aggregate
- [ ] Add Roles field to User entity as slice of role IDs
- [ ] Implement AssignRole(roleID) method
- [ ] Validate role is not already assigned
- [ ] Register UserRoleAssigned domain event
- [ ] Implement RevokeRole(roleID) method
- [ ] Validate role is currently assigned
- [ ] Validate user has at least one role after revocation (optional)
- [ ] Register UserRoleRevoked domain event
- [ ] Implement SetRoles(roleIDs) method for bulk update
- [ ] Implement HasRole(roleID) method
- [ ] Implement GetRoleIDs() method

### 1.3.2 User Domain Events for Roles
- [ ] Define UserRoleAssignedEvent with user ID and role ID
- [ ] Define UserRoleRevokedEvent with user ID and role ID
- [ ] Define UserRolesUpdatedEvent with user ID and role list

### 1.3.3 User Repository Updates
- [ ] Update UserRepository interface if needed
- [ ] Add FindByRole(ctx, roleID) ([]*User, error) method
- [ ] Add method to load user with roles eagerly

## 1.4 Default Roles & Permissions Seed

### 1.4.1 Define System Permissions
- [ ] Define all CRUD permissions for each resource
- [ ] users:create, users:read, users:update, users:delete, users:list
- [ ] roles:create, roles:read, roles:update, roles:delete, roles:list
- [ ] permissions:read, permissions:list
- [ ] Define special permissions
- [ ] users:manage (super permission for user management)
- [ ] roles:assign (permission to assign roles to users)
- [ ] system:admin (super admin permission)

### 1.4.2 Define System Roles
- [ ] Define SuperAdmin role
- [ ] Has all permissions
- [ ] IsSystem = true
- [ ] Highest priority
- [ ] Define Admin role
- [ ] Has user and role management permissions
- [ ] IsSystem = true
- [ ] Define Manager role (if needed)
- [ ] Has read and limited update permissions
- [ ] IsSystem = true
- [ ] Define User role
- [ ] Has basic read permissions for own data
- [ ] IsSystem = true
- [ ] IsDefault = true

### 1.4.3 Seeder Implementation
- [ ] Create database seeder for permissions
- [ ] Check if permission exists before creating
- [ ] Mark as IsSystem = true
- [ ] Create database seeder for roles
- [ ] Check if role exists before creating
- [ ] Assign appropriate permissions to each role
- [ ] Mark as IsSystem = true
- [ ] Create seeder for initial SuperAdmin user (optional)
- [ ] Integrate seeder with application startup or migration

---

# Phase 2: Authentication System

## 2.1 Authentication Domain (`internal/domain/auth/`)

### 2.1.1 Credential Value Objects
- [ ] Define Password value object
- [ ] Implement minimum length validation (e.g., 8 characters)
- [ ] Implement complexity validation (uppercase, lowercase, number, special char)
- [ ] Implement maximum length validation
- [ ] Implement common password check (optional, use dictionary)
- [ ] Define HashedPassword value object
- [ ] Store algorithm identifier with hash
- [ ] Implement Verify(plainPassword) method
- [ ] Define Email value object (if not already in shared)
- [ ] Validate format
- [ ] Normalize (lowercase, trim)

### 2.1.2 Token Value Objects
- [ ] Define AccessToken value object
- [ ] Include token string
- [ ] Include expiration time
- [ ] Include token type (Bearer)
- [ ] Implement IsExpired() method
- [ ] Define RefreshToken value object
- [ ] Include token string (opaque or JWT)
- [ ] Include expiration time
- [ ] Include user ID association
- [ ] Include device/session identifier (optional)
- [ ] Implement IsExpired() method
- [ ] Define TokenPair struct containing both tokens

### 2.1.3 Session/RefreshToken Entity (if storing refresh tokens)
- [ ] Define RefreshTokenEntity struct
- [ ] Include ID as UUID
- [ ] Include UserID
- [ ] Include TokenHash (store hash, not plain token)
- [ ] Include ExpiresAt
- [ ] Include CreatedAt
- [ ] Include LastUsedAt
- [ ] Include DeviceInfo (user agent, IP - optional)
- [ ] Include IsRevoked flag
- [ ] Implement IsValid() method checking expiration and revocation
- [ ] Implement Revoke() method
- [ ] Implement UpdateLastUsed() method

### 2.1.4 RefreshToken Repository Interface
- [ ] Define RefreshTokenRepository interface
- [ ] Define Create(ctx, token) error
- [ ] Define FindByTokenHash(ctx, hash) (*RefreshToken, error)
- [ ] Define FindByUserID(ctx, userID) ([]*RefreshToken, error)
- [ ] Define Revoke(ctx, id) error
- [ ] Define RevokeAllByUserID(ctx, userID) error
- [ ] Define DeleteExpired(ctx) (int64, error) for cleanup
- [ ] Define CountActiveByUserID(ctx, userID) (int, error) for session limit

### 2.1.5 Authentication Domain Services
- [ ] Define PasswordHasher interface
- [ ] Define Hash(password string) (HashedPassword, error)
- [ ] Define Verify(hashed HashedPassword, plain string) (bool, error)
- [ ] Define TokenGenerator interface
- [ ] Define GenerateAccessToken(user, roles, permissions) (AccessToken, error)
- [ ] Define GenerateRefreshToken() (string, error)
- [ ] Define ParseAccessToken(token string) (*Claims, error)
- [ ] Define Claims struct for token payload
- [ ] Include UserID
- [ ] Include Email
- [ ] Include Roles (role names or IDs)
- [ ] Include Permissions (permission codes)
- [ ] Include IssuedAt
- [ ] Include ExpiresAt
- [ ] Include TokenID (jti) for revocation

### 2.1.6 Authentication Domain Errors
- [ ] Define ErrInvalidCredentials
- [ ] Define ErrAccountLocked
- [ ] Define ErrAccountInactive
- [ ] Define ErrTokenExpired
- [ ] Define ErrTokenInvalid
- [ ] Define ErrTokenRevoked
- [ ] Define ErrRefreshTokenNotFound
- [ ] Define ErrRefreshTokenExpired
- [ ] Define ErrSessionLimitExceeded
- [ ] Define ErrPasswordTooWeak

## 2.2 Authentication Infrastructure

### 2.2.1 Password Hasher Implementation (`internal/infrastructure/security/`)
- [ ] Implement BcryptPasswordHasher
- [ ] Configure cost factor (recommend 12+ for production)
- [ ] Make cost factor configurable
- [ ] Implement Hash method using bcrypt.GenerateFromPassword
- [ ] Implement Verify method using bcrypt.CompareHashAndPassword
- [ ] Handle timing attacks (constant time comparison built into bcrypt)
- [ ] Consider implementing Argon2 hasher as alternative
- [ ] Configure memory, iterations, parallelism
- [ ] Write unit tests for hasher

### 2.2.2 JWT Token Implementation (`internal/infrastructure/security/`)
- [ ] Implement JWTTokenGenerator
- [ ] Configure signing method (RS256 recommended for production, HS256 for simplicity)
- [ ] Configure secret key or key pair from config
- [ ] Configure access token expiration (e.g., 15 minutes)
- [ ] Configure refresh token expiration (e.g., 7 days)
- [ ] Configure issuer claim
- [ ] Configure audience claim
- [ ] Implement GenerateAccessToken method
- [ ] Create claims with user info, roles, permissions
- [ ] Include standard claims (iss, aud, exp, iat, jti)
- [ ] Sign token with configured method
- [ ] Return AccessToken value object
- [ ] Implement GenerateRefreshToken method
- [ ] Generate cryptographically secure random string
- [ ] Or generate JWT with minimal claims
- [ ] Implement ParseAccessToken method
- [ ] Parse and validate token signature
- [ ] Validate expiration
- [ ] Validate issuer and audience
- [ ] Return Claims struct
- [ ] Handle all error cases with appropriate error types
- [ ] Write unit tests for token generation and parsing

### 2.2.3 RefreshToken Repository Implementation
- [ ] Create refresh_tokens table migration
- [ ] Include id, user_id, token_hash, expires_at, created_at, last_used_at, is_revoked
- [ ] Add index on token_hash for lookup
- [ ] Add index on user_id for user's sessions
- [ ] Add index on expires_at for cleanup job
- [ ] Implement PostgresRefreshTokenRepository
- [ ] Implement all interface methods
- [ ] Hash token before storing (use SHA256)
- [ ] Write integration tests

### 2.2.4 Token Blacklist (Optional - for access token revocation)
- [ ] Design token blacklist strategy
- [ ] Option 1: Store revoked token IDs in Redis with TTL
- [ ] Option 2: Store in database with cleanup job
- [ ] Option 3: Use short-lived tokens and rely on refresh token revocation
- [ ] Implement TokenBlacklist interface
- [ ] Define Add(tokenID, expiresAt) error
- [ ] Define IsBlacklisted(tokenID) (bool, error)
- [ ] Implement Redis-based blacklist
- [ ] Use token ID as key
- [ ] Set TTL to match token expiration
- [ ] Write tests for blacklist

## 2.3 Authentication Application Layer

### 2.3.1 Auth Commands (`internal/application/command/auth/`)

#### Register Command
- [ ] Define RegisterCommand struct
- [ ] Include Email
- [ ] Include Password
- [ ] Include FullName
- [ ] Include other registration fields
- [ ] Define RegisterHandler struct
- [ ] Inject UserRepository
- [ ] Inject RoleRepository
- [ ] Inject PasswordHasher
- [ ] Inject EventBus
- [ ] Inject Logger
- [ ] Implement Handle method
- [ ] Validate input
- [ ] Check email uniqueness
- [ ] Hash password
- [ ] Create User entity
- [ ] Assign default role to user
- [ ] Save user
- [ ] Publish domain events
- [ ] Return created user ID
- [ ] Write unit tests

#### Login Command
- [ ] Define LoginCommand struct
- [ ] Include Email
- [ ] Include Password
- [ ] Include DeviceInfo (optional)
- [ ] Include IP address (optional)
- [ ] Define LoginHandler struct
- [ ] Inject UserRepository
- [ ] Inject RoleRepository
- [ ] Inject PermissionRepository
- [ ] Inject RefreshTokenRepository
- [ ] Inject PasswordHasher
- [ ] Inject TokenGenerator
- [ ] Inject EventBus
- [ ] Inject Logger
- [ ] Implement Handle method
- [ ] Validate input
- [ ] Find user by email
- [ ] Return ErrInvalidCredentials if not found (don't reveal user exists)
- [ ] Check user status (active, not banned)
- [ ] Verify password
- [ ] Implement account lockout on repeated failures (optional)
- [ ] Load user's roles
- [ ] Load permissions for all roles
- [ ] Aggregate unique permissions
- [ ] Generate access token with roles and permissions
- [ ] Generate refresh token
- [ ] Store refresh token entity (hashed)
- [ ] Check session limit before creating new session (optional)
- [ ] Publish UserLoggedIn event
- [ ] Return TokenPair
- [ ] Write unit tests for all scenarios

#### Refresh Token Command
- [ ] Define RefreshTokenCommand struct
- [ ] Include RefreshToken string
- [ ] Define RefreshTokenHandler struct
- [ ] Inject RefreshTokenRepository
- [ ] Inject UserRepository
- [ ] Inject RoleRepository
- [ ] Inject PermissionRepository
- [ ] Inject TokenGenerator
- [ ] Implement Handle method
- [ ] Hash incoming refresh token
- [ ] Find refresh token entity by hash
- [ ] Return ErrRefreshTokenNotFound if not found
- [ ] Validate refresh token is not expired
- [ ] Validate refresh token is not revoked
- [ ] Load user by ID from refresh token
- [ ] Validate user is still active
- [ ] Load current roles and permissions (may have changed)
- [ ] Generate new access token
- [ ] Optionally rotate refresh token (generate new, revoke old)
- [ ] Update last_used_at on refresh token
- [ ] Return new TokenPair or just AccessToken
- [ ] Write unit tests

#### Logout Command
- [ ] Define LogoutCommand struct
- [ ] Include RefreshToken string (or AccessToken to extract jti)
- [ ] Include LogoutAll flag (optional)
- [ ] Define LogoutHandler struct
- [ ] Inject RefreshTokenRepository
- [ ] Inject TokenBlacklist (if using)
- [ ] Implement Handle method
- [ ] If LogoutAll: revoke all user's refresh tokens
- [ ] If single logout: revoke specific refresh token
- [ ] Optionally add access token to blacklist
- [ ] Publish UserLoggedOut event
- [ ] Write unit tests

#### Change Password Command (Auth context)
- [ ] Ensure existing ChangePasswordCommand revokes all refresh tokens
- [ ] Add step to revoke all sessions after password change
- [ ] Publish PasswordChanged event for audit

#### Forgot Password Command
- [ ] Define ForgotPasswordCommand struct
- [ ] Include Email
- [ ] Define ForgotPasswordHandler struct
- [ ] Inject UserRepository
- [ ] Inject PasswordResetTokenRepository (or use same as refresh token)
- [ ] Inject EmailService (interface)
- [ ] Inject TokenGenerator
- [ ] Implement Handle method
- [ ] Find user by email
- [ ] If not found, still return success (prevent email enumeration)
- [ ] Generate password reset token (short-lived, e.g., 1 hour)
- [ ] Store reset token with user ID
- [ ] Send email with reset link (async via event)
- [ ] Publish PasswordResetRequested event
- [ ] Write unit tests

#### Reset Password Command
- [ ] Define ResetPasswordCommand struct
- [ ] Include ResetToken
- [ ] Include NewPassword
- [ ] Define ResetPasswordHandler struct
- [ ] Inject UserRepository
- [ ] Inject PasswordResetTokenRepository
- [ ] Inject PasswordHasher
- [ ] Inject RefreshTokenRepository
- [ ] Implement Handle method
- [ ] Find and validate reset token
- [ ] Load user
- [ ] Validate new password strength
- [ ] Hash new password
- [ ] Update user password
- [ ] Invalidate reset token
- [ ] Revoke all existing refresh tokens
- [ ] Publish PasswordReset event
- [ ] Write unit tests

### 2.3.2 Auth Queries (`internal/application/query/auth/`)

#### Get Current User Query
- [ ] Define GetCurrentUserQuery struct
- [ ] Include UserID (extracted from token)
- [ ] Define GetCurrentUserHandler struct
- [ ] Implement Handle method returning UserDTO with roles

#### Get User Sessions Query
- [ ] Define GetUserSessionsQuery struct
- [ ] Include UserID
- [ ] Define GetUserSessionsHandler struct
- [ ] Implement Handle method returning list of active sessions
- [ ] Include device info, created at, last used at

### 2.3.3 Auth DTOs (`internal/application/dto/`)
- [ ] Define TokenPairDTO
- [ ] Include AccessToken
- [ ] Include RefreshToken
- [ ] Include ExpiresIn (seconds)
- [ ] Include TokenType ("Bearer")
- [ ] Define AuthUserDTO
- [ ] Include user info
- [ ] Include roles
- [ ] Include permissions
- [ ] Define SessionDTO for session listing

## 2.4 Authentication Interface Layer

### 2.4.1 Auth HTTP Handlers (`internal/interfaces/http/handler/`)

#### Auth Handler Structure
- [ ] Define AuthHandler struct
- [ ] Inject CommandBus
- [ ] Inject QueryBus
- [ ] Inject Validator
- [ ] Inject Logger

#### Register Endpoint
- [ ] Implement POST /auth/register
- [ ] Parse RegisterRequest DTO
- [ ] Validate input
- [ ] Dispatch RegisterCommand
- [ ] Return 201 with user ID or 200 with token pair (auto-login)

#### Login Endpoint
- [ ] Implement POST /auth/login
- [ ] Parse LoginRequest DTO (email, password)
- [ ] Extract device info from headers (User-Agent)
- [ ] Extract IP from request
- [ ] Dispatch LoginCommand
- [ ] Return 200 with TokenPairDTO
- [ ] Return 401 for invalid credentials
- [ ] Return 403 for locked/inactive account

#### Refresh Token Endpoint
- [ ] Implement POST /auth/refresh
- [ ] Parse RefreshTokenRequest DTO
- [ ] Or extract from HTTP-only cookie
- [ ] Dispatch RefreshTokenCommand
- [ ] Return 200 with new tokens
- [ ] Return 401 for invalid/expired refresh token

#### Logout Endpoint
- [ ] Implement POST /auth/logout
- [ ] Extract refresh token from request or cookie
- [ ] Dispatch LogoutCommand
- [ ] Clear HTTP-only cookie if using cookies
- [ ] Return 204 No Content

#### Logout All Sessions Endpoint
- [ ] Implement POST /auth/logout-all
- [ ] Require authentication
- [ ] Extract user ID from token
- [ ] Dispatch LogoutCommand with LogoutAll flag
- [ ] Return 204 No Content

#### Forgot Password Endpoint
- [ ] Implement POST /auth/forgot-password
- [ ] Parse ForgotPasswordRequest (email only)
- [ ] Dispatch ForgotPasswordCommand
- [ ] Always return 200 (prevent email enumeration)
- [ ] Include message about checking email

#### Reset Password Endpoint
- [ ] Implement POST /auth/reset-password
- [ ] Parse ResetPasswordRequest (token, new password)
- [ ] Dispatch ResetPasswordCommand
- [ ] Return 200 on success
- [ ] Return 400 for invalid/expired token

#### Get Current User Endpoint
- [ ] Implement GET /auth/me
- [ ] Require authentication
- [ ] Extract user ID from token claims
- [ ] Dispatch GetCurrentUserQuery
- [ ] Return 200 with AuthUserDTO

#### Get Sessions Endpoint
- [ ] Implement GET /auth/sessions
- [ ] Require authentication
- [ ] Dispatch GetUserSessionsQuery
- [ ] Return 200 with list of sessions

#### Revoke Session Endpoint
- [ ] Implement DELETE /auth/sessions/{sessionId}
- [ ] Require authentication
- [ ] Verify session belongs to current user
- [ ] Dispatch command to revoke specific session
- [ ] Return 204 No Content

### 2.4.2 Auth HTTP DTOs (`internal/interfaces/http/dto/`)
- [ ] Define RegisterRequest
- [ ] Define LoginRequest
- [ ] Define RefreshTokenRequest
- [ ] Define ForgotPasswordRequest
- [ ] Define ResetPasswordRequest
- [ ] Define TokenResponse matching TokenPairDTO
- [ ] Add appropriate validation tags to all request DTOs

### 2.4.3 Auth Routes
- [ ] Create auth route group /auth
- [ ] Register all auth endpoints
- [ ] Apply rate limiting to sensitive endpoints (login, register, forgot-password)
- [ ] Apply authentication middleware to protected endpoints (me, sessions, logout)

---

# Phase 3: Authorization System (RBAC)

## 3.1 Authorization Infrastructure

### 3.1.1 Permission Checker Service
- [ ] Define PermissionChecker interface
- [ ] Define HasPermission(ctx, userID, permission) (bool, error)
- [ ] Define HasAnyPermission(ctx, userID, permissions) (bool, error)
- [ ] Define HasAllPermissions(ctx, userID, permissions) (bool, error)
- [ ] Define GetUserPermissions(ctx, userID) ([]string, error)
- [ ] Implement PermissionChecker
- [ ] Option 1: Load from database on each check
- [ ] Option 2: Load from JWT claims (already in token)
- [ ] Option 3: Load from cache with invalidation
- [ ] For JWT-based approach:
- [ ] Extract permissions from token claims in context
- [ ] Check if required permission exists in claims
- [ ] For database approach:
- [ ] Load user's roles
- [ ] Load permissions for each role
- [ ] Aggregate and check
- [ ] Implement caching layer for database approach
- [ ] Cache user permissions with TTL
- [ ] Invalidate on role assignment change
- [ ] Invalidate on role permission change

### 3.1.2 Role Checker Service
- [ ] Define RoleChecker interface
- [ ] Define HasRole(ctx, userID, role) (bool, error)
- [ ] Define HasAnyRole(ctx, userID, roles) (bool, error)
- [ ] Define GetUserRoles(ctx, userID) ([]string, error)
- [ ] Implement similar to PermissionChecker

### 3.1.3 Authorization Context
- [ ] Define AuthContext struct to hold authenticated user info
- [ ] Include UserID
- [ ] Include Email
- [ ] Include Roles (slice of role names)
- [ ] Include Permissions (slice of permission codes)
- [ ] Define context key for AuthContext
- [ ] Implement helper to get AuthContext from context
- [ ] Implement helper to set AuthContext in context
- [ ] Return error or nil if not authenticated

## 3.2 Authorization Middleware

### 3.2.1 Authentication Middleware Enhancement
- [ ] Extract and parse JWT from Authorization header
- [ ] Validate token signature and expiration
- [ ] Check token blacklist (if implemented)
- [ ] Create AuthContext from token claims
- [ ] Inject AuthContext into request context
- [ ] Continue to next handler if valid
- [ ] Return 401 if token missing or invalid
- [ ] Log authentication attempts

### 3.2.2 Require Auth Middleware
- [ ] Implement middleware that requires authentication
- [ ] Check if AuthContext exists in context
- [ ] Return 401 if not authenticated
- [ ] Continue if authenticated
- [ ] Apply to routes that require any authenticated user

### 3.2.3 Require Permission Middleware
- [ ] Implement RequirePermission(permission string) middleware factory
- [ ] Extract AuthContext from context
- [ ] Return 401 if not authenticated
- [ ] Check if user has required permission
- [ ] Return 403 Forbidden if permission denied
- [ ] Log authorization failures
- [ ] Continue if authorized
- [ ] Implement RequireAnyPermission(permissions ...string) variant
- [ ] Implement RequireAllPermissions(permissions ...string) variant

### 3.2.4 Require Role Middleware
- [ ] Implement RequireRole(role string) middleware factory
- [ ] Similar logic to RequirePermission
- [ ] Return 403 if role not present
- [ ] Implement RequireAnyRole(roles ...string) variant

### 3.2.5 Resource Owner Middleware
- [ ] Implement middleware for resource ownership check
- [ ] Extract resource ID from URL parameter
- [ ] Extract user ID from AuthContext
- [ ] Compare resource owner with current user
- [ ] Allow if user is owner OR has admin permission
- [ ] Return 403 if not owner and not admin
- [ ] Make configurable for different resources

## 3.3 Permission Application Layer

### 3.3.1 Permission Commands (`internal/application/command/permission/`)

#### Create Permission Command
- [ ] Define CreatePermissionCommand
- [ ] Include Resource
- [ ] Include Action
- [ ] Include Description
- [ ] Define CreatePermissionHandler
- [ ] Validate input
- [ ] Check if permission code already exists
- [ ] Create Permission entity
- [ ] Save to repository
- [ ] Return permission ID
- [ ] Write unit tests

#### Update Permission Command
- [ ] Define UpdatePermissionCommand
- [ ] Include PermissionID
- [ ] Include Description (only description should be updatable)
- [ ] Define UpdatePermissionHandler
- [ ] Validate system permissions cannot be modified
- [ ] Update and save
- [ ] Write unit tests

#### Delete Permission Command
- [ ] Define DeletePermissionCommand
- [ ] Include PermissionID
- [ ] Define DeletePermissionHandler
- [ ] Validate system permissions cannot be deleted
- [ ] Check if permission is assigned to any role
- [ ] Option: Prevent deletion if in use
- [ ] Option: Cascade remove from roles
- [ ] Delete permission
- [ ] Write unit tests

### 3.3.2 Permission Queries (`internal/application/query/permission/`)
- [ ] Implement ListPermissions query
- [ ] Implement GetPermission query
- [ ] Implement ListPermissionsByResource query
- [ ] Implement GetPermissionsForRole query

## 3.4 Role Application Layer

### 3.4.1 Role Commands (`internal/application/command/role/`)

#### Create Role Command
- [ ] Define CreateRoleCommand
- [ ] Include Name
- [ ] Include DisplayName
- [ ] Include Description
- [ ] Include PermissionIDs (initial permissions)
- [ ] Define CreateRoleHandler
- [ ] Validate input
- [ ] Check name uniqueness
- [ ] Validate all permission IDs exist
- [ ] Create Role entity with permissions
- [ ] Save to repository
- [ ] Publish RoleCreated event
- [ ] Return role ID
- [ ] Write unit tests

#### Update Role Command
- [ ] Define UpdateRoleCommand
- [ ] Include RoleID
- [ ] Include DisplayName
- [ ] Include Description
- [ ] Define UpdateRoleHandler
- [ ] Validate system roles have limited modification
- [ ] Update allowed fields
- [ ] Save to repository
- [ ] Publish RoleUpdated event
- [ ] Write unit tests

#### Delete Role Command
- [ ] Define DeleteRoleCommand
- [ ] Include RoleID
- [ ] Define DeleteRoleHandler
- [ ] Validate not a system role
- [ ] Validate not the default role
- [ ] Check if role is assigned to any users
- [ ] Option: Prevent deletion if in use
- [ ] Option: Reassign users to default role
- [ ] Delete role
- [ ] Publish RoleDeleted event
- [ ] Write unit tests

#### Assign Permission to Role Command
- [ ] Define AssignPermissionToRoleCommand
- [ ] Include RoleID
- [ ] Include PermissionID
- [ ] Define AssignPermissionToRoleHandler
- [ ] Validate role exists
- [ ] Validate permission exists
- [ ] Load role
- [ ] Add permission using domain method
- [ ] Save role
- [ ] Invalidate permission cache for affected users
- [ ] Publish domain events
- [ ] Write unit tests

#### Remove Permission from Role Command
- [ ] Define RemovePermissionFromRoleCommand
- [ ] Similar to assign but removes
- [ ] Validate system role restrictions
- [ ] Write unit tests

#### Set Role Permissions Command (bulk update)
- [ ] Define SetRolePermissionsCommand
- [ ] Include RoleID
- [ ] Include PermissionIDs (complete list)
- [ ] Define SetRolePermissionsHandler
- [ ] Validate all permissions exist
- [ ] Load role
- [ ] Set permissions using domain method
- [ ] Save role
- [ ] Invalidate caches
- [ ] Publish events
- [ ] Write unit tests

### 3.4.2 Role Queries (`internal/application/query/role/`)
- [ ] Implement ListRoles query
- [ ] Implement GetRole query with permissions
- [ ] Implement GetRolePermissions query
- [ ] Implement GetUsersWithRole query

## 3.5 User Role Management

### 3.5.1 User Role Commands (`internal/application/command/user/`)

#### Assign Role to User Command
- [ ] Define AssignRoleToUserCommand
- [ ] Include UserID
- [ ] Include RoleID
- [ ] Define AssignRoleToUserHandler
- [ ] Validate user exists
- [ ] Validate role exists
- [ ] Load user
- [ ] Assign role using domain method
- [ ] Save user
- [ ] Invalidate user's permission cache
- [ ] Publish UserRoleAssigned event
- [ ] Write unit tests

#### Revoke Role from User Command
- [ ] Define RevokeRoleFromUserCommand
- [ ] Include UserID
- [ ] Include RoleID
- [ ] Define RevokeRoleFromUserHandler
- [ ] Validate user has at least one role remaining (optional)
- [ ] Load user
- [ ] Revoke role using domain method
- [ ] Save user
- [ ] Invalidate cache
- [ ] Publish UserRoleRevoked event
- [ ] Write unit tests

#### Set User Roles Command (bulk update)
- [ ] Define SetUserRolesCommand
- [ ] Include UserID
- [ ] Include RoleIDs
- [ ] Define SetUserRolesHandler
- [ ] Validate at least one role (optional)
- [ ] Validate all roles exist
- [ ] Load user
- [ ] Set roles
- [ ] Save user
- [ ] Invalidate cache
- [ ] Publish event
- [ ] Write unit tests

### 3.5.2 User Role Queries
- [ ] Implement GetUserRoles query
- [ ] Implement GetUserPermissions query (aggregated from all roles)

## 3.6 Role & Permission HTTP Handlers

### 3.6.1 Permission Handler
- [ ] Implement GET /permissions - List all permissions
- [ ] Require permissions:list permission
- [ ] Implement GET /permissions/{id} - Get permission details
- [ ] Require permissions:read permission
- [ ] Implement POST /permissions - Create permission (if allowing custom)
- [ ] Require permissions:create permission
- [ ] Implement PUT /permissions/{id} - Update permission
- [ ] Require permissions:update permission
- [ ] Implement DELETE /permissions/{id} - Delete permission
- [ ] Require permissions:delete permission

### 3.6.2 Role Handler
- [ ] Implement GET /roles - List all roles
- [ ] Require roles:list permission
- [ ] Implement GET /roles/{id} - Get role with permissions
- [ ] Require roles:read permission
- [ ] Implement POST /roles - Create role
- [ ] Require roles:create permission
- [ ] Implement PUT /roles/{id} - Update role details
- [ ] Require roles:update permission
- [ ] Implement DELETE /roles/{id} - Delete role
- [ ] Require roles:delete permission
- [ ] Implement PUT /roles/{id}/permissions - Set role permissions
- [ ] Require roles:update permission
- [ ] Implement POST /roles/{id}/permissions/{permissionId} - Add permission
- [ ] Require roles:update permission
- [ ] Implement DELETE /roles/{id}/permissions/{permissionId} - Remove permission
- [ ] Require roles:update permission

### 3.6.3 User Role Management Endpoints
- [ ] Implement GET /users/{id}/roles - Get user's roles
- [ ] Require users:read or self
- [ ] Implement PUT /users/{id}/roles - Set user's roles
- [ ] Require roles:assign permission
- [ ] Implement POST /users/{id}/roles/{roleId} - Assign role
- [ ] Require roles:assign permission
- [ ] Implement DELETE /users/{id}/roles/{roleId} - Revoke role
- [ ] Require roles:assign permission
- [ ] Implement GET /users/{id}/permissions - Get user's effective permissions
- [ ] Require users:read or self

---

# Phase 4: Database Schema for RBAC

## 4.1 Permissions Table Migration
- [ ] Generate migration: `goose create create_permissions_table sql`
- [ ] Define permissions table
- [ ] id UUID PRIMARY KEY
- [ ] resource VARCHAR(100) NOT NULL
- [ ] action VARCHAR(100) NOT NULL
- [ ] description TEXT
- [ ] is_system BOOLEAN DEFAULT FALSE
- [ ] created_at TIMESTAMPTZ DEFAULT NOW()
- [ ] updated_at TIMESTAMPTZ DEFAULT NOW()
- [ ] Add unique constraint on (resource, action)
- [ ] Add index on resource
- [ ] Add index on is_system
- [ ] Write down migration

## 4.2 Roles Table Migration
- [ ] Generate migration: `goose create create_roles_table sql`
- [ ] Define roles table
- [ ] id UUID PRIMARY KEY
- [ ] name VARCHAR(100) UNIQUE NOT NULL
- [ ] display_name VARCHAR(255) NOT NULL
- [ ] description TEXT
- [ ] is_system BOOLEAN DEFAULT FALSE
- [ ] is_default BOOLEAN DEFAULT FALSE
- [ ] priority INTEGER DEFAULT 0
- [ ] created_at TIMESTAMPTZ DEFAULT NOW()
- [ ] updated_at TIMESTAMPTZ DEFAULT NOW()
- [ ] Add index on name
- [ ] Add index on is_default
- [ ] Add constraint: only one default role
- [ ] Write down migration

## 4.3 Role Permissions Junction Table Migration
- [ ] Generate migration: `goose create create_role_permissions_table sql`
- [ ] Define role_permissions table
- [ ] role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE
- [ ] permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE
- [ ] created_at TIMESTAMPTZ DEFAULT NOW()
- [ ] Add PRIMARY KEY (role_id, permission_id)
- [ ] Add index on permission_id for reverse lookup
- [ ] Write down migration

## 4.4 User Roles Junction Table Migration
- [ ] Generate migration: `goose create create_user_roles_table sql`
- [ ] Define user_roles table
- [ ] user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
- [ ] role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE
- [ ] assigned_at TIMESTAMPTZ DEFAULT NOW()
- [ ] assigned_by UUID REFERENCES users(id) (optional, for audit)
- [ ] Add PRIMARY KEY (user_id, role_id)
- [ ] Add index on role_id for reverse lookup
- [ ] Write down migration

## 4.5 Refresh Tokens Table Migration
- [ ] Generate migration: `goose create create_refresh_tokens_table sql`
- [ ] Define refresh_tokens table
- [ ] id UUID PRIMARY KEY
- [ ] user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
- [ ] token_hash VARCHAR(255) UNIQUE NOT NULL
- [ ] expires_at TIMESTAMPTZ NOT NULL
- [ ] created_at TIMESTAMPTZ DEFAULT NOW()
- [ ] last_used_at TIMESTAMPTZ
- [ ] is_revoked BOOLEAN DEFAULT FALSE
- [ ] device_info JSONB (optional)
- [ ] ip_address INET (optional)
- [ ] Add index on user_id
- [ ] Add index on expires_at for cleanup
- [ ] Add index on is_revoked
- [ ] Write down migration

## 4.6 Password Reset Tokens Table Migration (Optional)
- [ ] Generate migration: `goose create create_password_reset_tokens_table sql`
- [ ] Define password_reset_tokens table
- [ ] id UUID PRIMARY KEY
- [ ] user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
- [ ] token_hash VARCHAR(255) UNIQUE NOT NULL
- [ ] expires_at TIMESTAMPTZ NOT NULL
- [ ] created_at TIMESTAMPTZ DEFAULT NOW()
- [ ] used_at TIMESTAMPTZ
- [ ] Add index on user_id
- [ ] Add index on expires_at
- [ ] Write down migration

## 4.7 Seed Data Migration
- [ ] Generate migration: `goose create seed_permissions_and_roles go`
- [ ] Use Go migration for complex seeding logic
- [ ] Seed all system permissions
- [ ] Seed all system roles
- [ ] Assign permissions to roles
- [ ] Create initial super admin user (optional)
- [ ] Make idempotent (check before insert)

---

# Phase 5: Repository Implementations

## 5.1 Permission Repository Implementation
- [ ] Implement PostgresPermissionRepository
- [ ] Implement Create with unique constraint handling
- [ ] Implement Update
- [ ] Implement Delete with system permission check
- [ ] Implement FindByID
- [ ] Implement FindByCode with (resource, action) lookup
- [ ] Implement FindByResource
- [ ] Implement FindAll with optional filtering
- [ ] Implement ExistsByCode
- [ ] Write integration tests for all methods

## 5.2 Role Repository Implementation
- [ ] Implement PostgresRoleRepository
- [ ] Implement Create with permissions insertion
- [ ] Implement Update
- [ ] Implement Delete with cascade consideration
- [ ] Implement FindByID with permissions loading
- [ ] Implement FindByName
- [ ] Implement FindByIDs for batch loading
- [ ] Implement FindAll
- [ ] Implement FindDefault
- [ ] Implement ExistsByName
- [ ] Implement FindByPermission
- [ ] Write integration tests

## 5.3 User Repository Updates
- [ ] Update user repository to handle roles
- [ ] Load roles when finding user (eager or lazy)
- [ ] Update Create to assign default role
- [ ] Implement FindByRole
- [ ] Update integration tests

## 5.4 RefreshToken Repository Implementation
- [ ] Implement PostgresRefreshTokenRepository
- [ ] Implement Create with token hashing
- [ ] Implement FindByTokenHash
- [ ] Implement FindByUserID
- [ ] Implement Revoke
- [ ] Implement RevokeAllByUserID
- [ ] Implement DeleteExpired for cleanup job
- [ ] Implement CountActiveByUserID
- [ ] Write integration tests

---

# Phase 6: Security Hardening

## 6.1 Rate Limiting
- [ ] Implement rate limiting for login endpoint
- [ ] Limit by IP address
- [ ] Stricter limit for failed attempts
- [ ] Implement rate limiting for registration
- [ ] Implement rate limiting for password reset request
- [ ] Implement rate limiting for token refresh
- [ ] Configure rate limits via environment
- [ ] Return 429 with Retry-After header
- [ ] Log rate limit violations

## 6.2 Account Security
- [ ] Implement account lockout after N failed login attempts
- [ ] Track failed attempts in database or cache
- [ ] Lock duration configurable (e.g., 15 minutes)
- [ ] Reset counter on successful login
- [ ] Implement CAPTCHA integration for repeated failures (optional)
- [ ] Implement login notification email (optional)
- [ ] Detect suspicious login (new device, new location)
- [ ] Send email notification

## 6.3 Token Security
- [ ] Use short-lived access tokens (15 minutes recommended)
- [ ] Use longer-lived refresh tokens with rotation
- [ ] Implement refresh token rotation
- [ ] Issue new refresh token on each refresh
- [ ] Revoke old refresh token
- [ ] Detect refresh token reuse (potential theft)
- [ ] Consider using HTTP-only cookies for tokens
- [ ] Set Secure flag for HTTPS only
- [ ] Set SameSite attribute
- [ ] Implement proper CSRF protection if using cookies

## 6.4 Password Security
- [ ] Enforce minimum password length (12+ recommended)
- [ ] Enforce complexity requirements
- [ ] Check against common password list
- [ ] Consider using zxcvbn for strength estimation
- [ ] Implement password history (prevent reuse of last N passwords)
- [ ] Implement password expiration policy (optional, controversial)
- [ ] Hash passwords with appropriate cost factor
- [ ] Consider upgrading hash on login if algorithm changes

## 6.5 Session Management
- [ ] Implement maximum concurrent sessions per user
- [ ] Allow users to view active sessions
- [ ] Allow users to revoke specific sessions
- [ ] Allow users to revoke all other sessions
- [ ] Implement session timeout for inactivity (optional)
- [ ] Implement absolute session timeout

## 6.6 Audit Logging
- [ ] Log all authentication events
- [ ] Successful login with IP, device info
- [ ] Failed login with IP, device info
- [ ] Logout
- [ ] Password change
- [ ] Password reset request and completion
- [ ] Log all authorization events
- [ ] Permission denied attempts
- [ ] Role assignments and revocations
- [ ] Store audit logs securely
- [ ] Implement log retention policy
- [ ] Consider separate audit log storage

---

# Phase 7: Testing

## 7.1 Unit Tests

### 7.1.1 Domain Tests
- [ ] Test Permission entity creation and validation
- [ ] Test Role entity and all business methods
- [ ] Test permission assignment/removal
- [ ] Test status checks
- [ ] Test User role methods
- [ ] Test all value objects
- [ ] Test all domain events

### 7.1.2 Application Tests
- [ ] Test all auth command handlers
- [ ] Mock repositories and services
- [ ] Test success and error paths
- [ ] Test all role/permission command handlers
- [ ] Test all query handlers
- [ ] Test authorization logic in handlers

### 7.1.3 Infrastructure Tests
- [ ] Test password hasher
- [ ] Test JWT token generator
- [ ] Token generation and parsing
- [ ] Expiration handling
- [ ] Invalid token handling
- [ ] Test permission checker service

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
- [ ] Test rate limiting is enforced
- [ ] Test account lockout works
- [ ] Test password validation rejects weak passwords
- [ ] Test token expiration is enforced
- [ ] Test refresh token rotation works
- [ ] Test revoked tokens are rejected
- [ ] Test system roles/permissions cannot be deleted

---

# Phase 8: Documentation

## 8.1 API Documentation
- [ ] Document all auth endpoints in OpenAPI spec
- [ ] Include request/response schemas
- [ ] Include error responses
- [ ] Document authentication mechanism
- [ ] Document all role/permission endpoints
- [ ] Document required permissions for each endpoint
- [ ] Include examples for common flows

## 8.2 Architecture Documentation
- [ ] Document RBAC model and relationships
- [ ] Create diagram showing User-Role-Permission relationships
- [ ] Document authentication flow with sequence diagram
- [ ] Document token refresh flow
- [ ] Document authorization middleware chain
- [ ] Document permission checking strategies
- [ ] Document caching strategy for permissions

## 8.3 Operations Documentation
- [ ] Document how to create new permissions
- [ ] Document how to create new roles
- [ ] Document how to assign roles to users
- [ ] Document how to handle locked accounts
- [ ] Document token cleanup job setup
- [ ] Document audit log monitoring
- [ ] Document incident response for security events

---

# Final Checklist

- [ ] All RBAC entities implemented and tested
- [ ] All authentication flows working
- [ ] All authorization checks in place
- [ ] Rate limiting configured
- [ ] Account security measures implemented
- [ ] Token security best practices followed
- [ ] Audit logging enabled
- [ ] All tests passing
- [ ] API documentation complete
- [ ] Security review completed
- [ ] Performance testing for auth endpoints done
- [ ] Monitoring and alerting configured for auth failures
