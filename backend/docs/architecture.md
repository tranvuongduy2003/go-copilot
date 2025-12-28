# Authentication & RBAC Architecture

## Overview

This document describes the authentication and Role-Based Access Control (RBAC) architecture for the Go Copilot API.

## RBAC Model

### Entity Relationships

```
┌─────────────┐       ┌─────────────┐       ┌─────────────────┐
│    User     │──────<│  UserRole   │>──────│      Role       │
│             │       │  (junction) │       │                 │
│ - id        │       │ - user_id   │       │ - id            │
│ - email     │       │ - role_id   │       │ - name          │
│ - password  │       │ - assigned_at│      │ - display_name  │
│ - full_name │       └─────────────┘       │ - is_system     │
│ - status    │                             │ - is_default    │
└─────────────┘                             │ - priority      │
                                            └────────┬────────┘
                                                     │
                                            ┌────────┴────────┐
                                            │ RolePermission  │
                                            │   (junction)    │
                                            │ - role_id       │
                                            │ - permission_id │
                                            └────────┬────────┘
                                                     │
                                            ┌────────┴────────┐
                                            │   Permission    │
                                            │                 │
                                            │ - id            │
                                            │ - resource      │
                                            │ - action        │
                                            │ - is_system     │
                                            └─────────────────┘
```

### Permission Code Format

Permissions follow the format: `{resource}:{action}`

Examples:
- `users:read` - Read user information
- `users:create` - Create new users
- `roles:assign` - Assign roles to users
- `system:admin` - Super admin permission (grants all access)

### System Roles

| Role | Description | Key Permissions |
|------|-------------|-----------------|
| `super_admin` | Full system access | `system:admin` (all permissions) |
| `admin` | User and role management | `users:*`, `roles:*`, `permissions:read` |
| `manager` | Read access with limited updates | `users:read`, `users:list` |
| `user` | Default role for new users | Basic read permissions |

## Authentication Flow

### Login Flow

```
┌──────────┐     ┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│  Client  │     │ Auth Handler│     │ Login Handler│     │   Services  │
└────┬─────┘     └──────┬──────┘     └──────┬───────┘     └──────┬──────┘
     │                  │                   │                    │
     │ POST /auth/login │                   │                    │
     │ {email, password}│                   │                    │
     │─────────────────>│                   │                    │
     │                  │                   │                    │
     │                  │ Check rate limit  │                    │
     │                  │──────────────────>│                    │
     │                  │                   │                    │
     │                  │                   │ Check account lockout
     │                  │                   │───────────────────>│
     │                  │                   │                    │
     │                  │                   │ Find user by email │
     │                  │                   │───────────────────>│
     │                  │                   │                    │
     │                  │                   │ Verify password    │
     │                  │                   │───────────────────>│
     │                  │                   │                    │
     │                  │                   │ Load roles &       │
     │                  │                   │ permissions        │
     │                  │                   │───────────────────>│
     │                  │                   │                    │
     │                  │                   │ Generate tokens    │
     │                  │                   │───────────────────>│
     │                  │                   │                    │
     │                  │                   │ Store refresh token│
     │                  │                   │───────────────────>│
     │                  │                   │                    │
     │                  │                   │ Publish login event│
     │                  │                   │───────────────────>│
     │                  │                   │                    │
     │ 200 OK           │                   │                    │
     │ {access_token,   │                   │                    │
     │  refresh_token}  │                   │                    │
     │<─────────────────│                   │                    │
     │                  │                   │                    │
```

### Token Refresh Flow

```
┌──────────┐     ┌─────────────┐     ┌────────────────┐     ┌─────────────┐
│  Client  │     │ Auth Handler│     │ Refresh Handler│     │   Services  │
└────┬─────┘     └──────┬──────┘     └───────┬────────┘     └──────┬──────┘
     │                  │                    │                     │
     │ POST /auth/refresh                    │                     │
     │ {refresh_token}  │                    │                     │
     │─────────────────>│                    │                     │
     │                  │                    │                     │
     │                  │ Hash refresh token │                     │
     │                  │───────────────────>│                     │
     │                  │                    │                     │
     │                  │                    │ Find token by hash  │
     │                  │                    │────────────────────>│
     │                  │                    │                     │
     │                  │                    │ Validate not expired│
     │                  │                    │ Validate not revoked│
     │                  │                    │────────────────────>│
     │                  │                    │                     │
     │                  │                    │ Load user           │
     │                  │                    │────────────────────>│
     │                  │                    │                     │
     │                  │                    │ Load current roles  │
     │                  │                    │ & permissions       │
     │                  │                    │────────────────────>│
     │                  │                    │                     │
     │                  │                    │ Generate new tokens │
     │                  │                    │────────────────────>│
     │                  │                    │                     │
     │                  │                    │ Rotate refresh token│
     │                  │                    │ (revoke old, store  │
     │                  │                    │  new)               │
     │                  │                    │────────────────────>│
     │                  │                    │                     │
     │ 200 OK           │                    │                     │
     │ {new tokens}     │                    │                     │
     │<─────────────────│                    │                     │
```

## Authorization Middleware Chain

```
Request
   │
   ▼
┌─────────────────────────────────────────────────────────────┐
│                     Global Middleware                        │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌────────┐│
│  │RequestID│→│ Logging │→│Recovery │→│  CORS   │→│Timeout ││
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘ └────────┘│
└─────────────────────────────────────────────────────────────┘
   │
   ▼
┌─────────────────────────────────────────────────────────────┐
│                   Route-Specific Middleware                  │
│                                                              │
│  Public Routes (e.g., /auth/login):                         │
│  ┌───────────┐                                              │
│  │Rate Limit │ → Handler                                    │
│  └───────────┘                                              │
│                                                              │
│  Protected Routes (e.g., /users):                           │
│  ┌───────────┐   ┌──────────────┐   ┌──────────────────┐   │
│  │Rate Limit │ → │ RequireAuth  │ → │RequirePermission │   │
│  └───────────┘   │(parse JWT,   │   │(check permission │   │
│                  │ set context) │   │ in claims)       │   │
│                  └──────────────┘   └──────────────────┘   │
│                           │                   │              │
│                           ▼                   ▼              │
│                      AuthContext         Handler             │
│                      {UserID,                                │
│                       Roles,                                 │
│                       Permissions}                           │
└─────────────────────────────────────────────────────────────┘
```

### Middleware Components

| Middleware | Purpose |
|------------|---------|
| `RequestID` | Adds unique request ID for tracing |
| `Logging` | Logs request/response details |
| `Recovery` | Catches panics, returns 500 |
| `CORS` | Handles cross-origin requests |
| `Timeout` | Enforces request timeout |
| `RateLimit` | Limits requests per IP |
| `RequireAuth` | Parses JWT, sets AuthContext |
| `RequirePermission` | Checks specific permission |
| `RequireAnyPermission` | Checks any of listed permissions |
| `RequireAllPermissions` | Checks all listed permissions |
| `RequireRole` | Checks specific role |
| `ResourceOwner` | Checks resource ownership or admin |

## Permission Checking Strategy

### JWT-Based Approach (Current Implementation)

Permissions are embedded in JWT access tokens at login/refresh time.

**Pros:**
- No database lookup on each request
- Fast authorization checks
- Stateless verification

**Cons:**
- Permissions not immediately revoked (wait for token expiry)
- Larger token size with many permissions

```go
// Permissions stored in JWT claims
type Claims struct {
    UserID      uuid.UUID `json:"user_id"`
    Email       string    `json:"email"`
    Roles       []string  `json:"roles"`
    Permissions []string  `json:"permissions"`
    // ...
}

// Permission check in middleware
func hasPermission(userPermissions []string, required string) bool {
    for _, p := range userPermissions {
        if p == required || p == "system:admin" {
            return true
        }
    }
    return false
}
```

### Permission Resolution

When a user logs in:

1. Load user's assigned roles
2. For each role, load associated permissions
3. Aggregate all unique permission codes
4. Include in JWT claims

```
User → [Role1, Role2] → [Perm1, Perm2, Perm3, Perm2, Perm4] → [Perm1, Perm2, Perm3, Perm4]
                                 (deduplicated)
```

## Token Security

### Access Token
- Short-lived (15 minutes default)
- Contains user info, roles, permissions
- Signed with HS256/RS256
- Validated on each request

### Refresh Token
- Longer-lived (7 days default)
- Opaque token (random string)
- Stored as hash in database
- Rotated on each use
- Can be revoked

### Token Blacklist (Redis)
- Stores revoked access token IDs
- TTL matches token expiration
- Checked during authentication

## Account Security Features

### Rate Limiting

| Endpoint | Limit | Burst |
|----------|-------|-------|
| Login | 1/sec | 5 |
| Register | 1/sec | 3 |
| Password Reset | 1/sec | 3 |
| Token Refresh | 1/sec | 30 |
| Default | 10/sec | 20 |

### Account Lockout

- 5 failed login attempts → 15 minute lockout
- Counter resets on successful login
- Stored in Redis with TTL

## Database Schema

### Core Tables

```sql
-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Permissions table
CREATE TABLE permissions (
    id UUID PRIMARY KEY,
    resource VARCHAR(100) NOT NULL,
    action VARCHAR(100) NOT NULL,
    description TEXT,
    is_system BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(resource, action)
);

-- Roles table
CREATE TABLE roles (
    id UUID PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    description TEXT,
    is_system BOOLEAN DEFAULT FALSE,
    is_default BOOLEAN DEFAULT FALSE,
    priority INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Junction tables
CREATE TABLE user_roles (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID REFERENCES roles(id) ON DELETE CASCADE,
    assigned_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, role_id)
);

CREATE TABLE role_permissions (
    role_id UUID REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

-- Refresh tokens table
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    is_revoked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    last_used_at TIMESTAMPTZ
);
```

## Audit Logging

All security-relevant events are logged:

- `UserLoggedIn` - Successful login with IP, device
- `UserLoggedOut` - Logout event
- `LoginFailed` - Failed login attempt
- `AccountLocked` - Account lockout triggered
- `PasswordChanged` - Password change
- `PasswordReset` - Password reset completion
- `UserRoleAssigned` - Role assignment
- `UserRoleRevoked` - Role revocation
- `RoleCreated/Updated/Deleted` - Role changes
- `PermissionDenied` - Authorization failures

Events are stored in the `audit_logs` table with:
- Event type
- User ID (if applicable)
- IP address
- User agent
- Timestamp
- Additional metadata (JSON)
