# Operations Guide

This guide covers common administrative operations for the authentication and RBAC system.

## Table of Contents

1. [Health Checks & Monitoring](#health-checks--monitoring)
2. [Metrics](#metrics)
3. [Creating Permissions](#creating-permissions)
4. [Creating Roles](#creating-roles)
5. [Assigning Roles to Users](#assigning-roles-to-users)
6. [Handling Locked Accounts](#handling-locked-accounts)
7. [Token Cleanup](#token-cleanup)
8. [Audit Log Monitoring](#audit-log-monitoring)
9. [Incident Response](#incident-response)

---

## Health Checks & Monitoring

### Endpoints

| Endpoint | Purpose | Response Code |
|----------|---------|---------------|
| `GET /health` | Readiness check (all dependencies) | 200 OK / 503 Service Unavailable |
| `GET /health/live` | Liveness probe (process alive) | 200 OK |
| `GET /health/ready` | Alias for `/health` | 200 OK / 503 Service Unavailable |

### Liveness Check

The liveness endpoint checks if the application process is running. Use this for Kubernetes liveness probes.

```bash
curl http://localhost:8080/health/live
```

Response:
```json
{
  "status": "ok"
}
```

### Readiness Check

The readiness endpoint verifies all dependencies are healthy before accepting traffic.

```bash
curl http://localhost:8080/health/ready
```

Response (healthy):
```json
{
  "status": "healthy",
  "checks": {
    "database": {
      "status": "healthy"
    },
    "redis": {
      "status": "healthy"
    }
  }
}
```

Response (unhealthy):
```json
{
  "status": "unhealthy",
  "checks": {
    "database": {
      "status": "healthy"
    },
    "redis": {
      "status": "unhealthy",
      "message": "connection refused"
    }
  }
}
```

### Kubernetes Configuration

```yaml
# Deployment probe configuration
livenessProbe:
  httpGet:
    path: /health/live
    port: 8080
  initialDelaySeconds: 10
  periodSeconds: 15
  timeoutSeconds: 5
  failureThreshold: 3

readinessProbe:
  httpGet:
    path: /health/ready
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 10
  timeoutSeconds: 5
  failureThreshold: 3
```

---

## Metrics

### Prometheus Endpoint

Metrics are exposed at `GET /metrics` in Prometheus format.

```bash
curl http://localhost:8080/metrics
```

### Available Metrics

#### HTTP Request Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `http_requests_total` | Counter | Total HTTP requests by method, path, status |
| `http_request_duration_seconds` | Histogram | Request duration in seconds |
| `http_request_size_bytes` | Histogram | Request body size |
| `http_response_size_bytes` | Histogram | Response body size |

#### Database Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `db_pool_connections_total` | Gauge | Total connections in pool |
| `db_pool_connections_idle` | Gauge | Idle connections in pool |
| `db_pool_connections_in_use` | Gauge | Active connections in pool |
| `db_query_duration_seconds` | Histogram | Query execution time |

#### Business Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `user_registrations_total` | Counter | Total user registrations |
| `user_logins_total` | Counter | Total login attempts (success/failure) |
| `active_sessions_total` | Gauge | Current active sessions |

### Prometheus Configuration

```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'go-copilot'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: /metrics
    scrape_interval: 15s
```

### Grafana Dashboard

Import the provided dashboard or create panels for:

1. **Request Rate**: `rate(http_requests_total[5m])`
2. **Error Rate**: `rate(http_requests_total{status=~"5.."}[5m])`
3. **Latency P95**: `histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))`
4. **DB Connection Usage**: `db_pool_connections_in_use / db_pool_connections_total`

### Alerting Rules

```yaml
# alerting_rules.yml
groups:
  - name: go-copilot
    rules:
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.1
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: High 5xx error rate

      - alert: HighLatency
        expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: P95 latency above 1 second

      - alert: DatabaseConnectionPoolExhausted
        expr: db_pool_connections_in_use / db_pool_connections_total > 0.9
        for: 2m
        labels:
          severity: warning
        annotations:
          summary: Database connection pool near capacity
```

---

## Creating Permissions

### Via API

```bash
# Create a new permission
curl -X POST http://localhost:8080/api/v1/permissions \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "resource": "articles",
    "action": "publish",
    "description": "Allows publishing articles"
  }'
```

### Via Database Migration

For system permissions, add them in a migration:

```sql
-- migrations/000011_add_article_permissions.up.sql
INSERT INTO permissions (id, resource, action, description, is_system, created_at, updated_at)
VALUES
    (gen_random_uuid(), 'articles', 'create', 'Create articles', true, NOW(), NOW()),
    (gen_random_uuid(), 'articles', 'read', 'Read articles', true, NOW(), NOW()),
    (gen_random_uuid(), 'articles', 'update', 'Update articles', true, NOW(), NOW()),
    (gen_random_uuid(), 'articles', 'delete', 'Delete articles', true, NOW(), NOW()),
    (gen_random_uuid(), 'articles', 'publish', 'Publish articles', true, NOW(), NOW())
ON CONFLICT (resource, action) DO NOTHING;
```

### Permission Naming Conventions

| Resource | Actions | Examples |
|----------|---------|----------|
| Singular noun, lowercase | Standard CRUD + custom | `users:create`, `users:read` |
| Use underscores for multi-word | | `audit_logs:read` |
| Special permissions | | `system:admin`, `users:manage` |

---

## Creating Roles

### Via API

```bash
# Create a new role with permissions
curl -X POST http://localhost:8080/api/v1/roles \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "content_editor",
    "display_name": "Content Editor",
    "description": "Can create and edit articles",
    "permission_ids": [
      "uuid-of-articles-create",
      "uuid-of-articles-read",
      "uuid-of-articles-update"
    ]
  }'
```

### Via Database Migration

For system roles:

```sql
-- migrations/000012_add_content_editor_role.up.sql
DO $$
DECLARE
    role_id UUID := gen_random_uuid();
    perm_id UUID;
BEGIN
    -- Create the role
    INSERT INTO roles (id, name, display_name, description, is_system, is_default, priority, created_at, updated_at)
    VALUES (role_id, 'content_editor', 'Content Editor', 'Can create and edit articles', true, false, 50, NOW(), NOW())
    ON CONFLICT (name) DO NOTHING;

    -- Assign permissions
    FOR perm_id IN
        SELECT id FROM permissions
        WHERE resource = 'articles' AND action IN ('create', 'read', 'update')
    LOOP
        INSERT INTO role_permissions (role_id, permission_id, created_at)
        VALUES (role_id, perm_id, NOW())
        ON CONFLICT DO NOTHING;
    END LOOP;
END $$;
```

### Modifying Role Permissions

```bash
# Set all permissions for a role (replaces existing)
curl -X PUT http://localhost:8080/api/v1/roles/<role_id>/permissions \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "permission_ids": ["perm-uuid-1", "perm-uuid-2", "perm-uuid-3"]
  }'

# Add a single permission
curl -X POST http://localhost:8080/api/v1/roles/<role_id>/permissions/<permission_id> \
  -H "Authorization: Bearer <admin_token>"

# Remove a single permission
curl -X DELETE http://localhost:8080/api/v1/roles/<role_id>/permissions/<permission_id> \
  -H "Authorization: Bearer <admin_token>"
```

---

## Assigning Roles to Users

### Via API

```bash
# Assign a single role to a user
curl -X POST http://localhost:8080/api/v1/users/<user_id>/roles/<role_id> \
  -H "Authorization: Bearer <admin_token>"

# Set all roles for a user (replaces existing)
curl -X PUT http://localhost:8080/api/v1/users/<user_id>/roles \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "role_ids": ["role-uuid-1", "role-uuid-2"]
  }'

# Revoke a role from a user
curl -X DELETE http://localhost:8080/api/v1/users/<user_id>/roles/<role_id> \
  -H "Authorization: Bearer <admin_token>"
```

### Via Database (Emergency)

```sql
-- Assign role to user
INSERT INTO user_roles (user_id, role_id, assigned_at)
VALUES ('user-uuid', 'role-uuid', NOW())
ON CONFLICT DO NOTHING;

-- Revoke role from user
DELETE FROM user_roles
WHERE user_id = 'user-uuid' AND role_id = 'role-uuid';

-- Assign super_admin role to a user (emergency admin access)
INSERT INTO user_roles (user_id, role_id, assigned_at)
SELECT 'user-uuid', id, NOW()
FROM roles WHERE name = 'super_admin';
```

---

## Handling Locked Accounts

### Check Account Lock Status

```bash
# Check via Redis CLI
redis-cli GET "lockout:user@example.com"
redis-cli TTL "lockout:user@example.com"
```

### Unlock an Account

```bash
# Via Redis CLI
redis-cli DEL "lockout:user@example.com"
redis-cli DEL "login_attempts:user@example.com"
```

### Unlock via API (if implemented)

```bash
# Admin endpoint to unlock account
curl -X POST http://localhost:8080/api/v1/admin/users/<user_id>/unlock \
  -H "Authorization: Bearer <admin_token>"
```

### View Failed Attempt Count

```bash
redis-cli GET "login_attempts:user@example.com"
```

### Lockout Configuration

Environment variables:
```bash
LOCKOUT_MAX_ATTEMPTS=5          # Lock after 5 failures
LOCKOUT_DURATION=15m            # Lock for 15 minutes
LOCKOUT_ATTEMPT_WINDOW=15m      # Reset counter after 15 minutes of no failures
```

---

## Token Cleanup

### Expired Refresh Tokens

Set up a cron job to clean expired tokens:

```bash
# Run daily at 3 AM
0 3 * * * /usr/local/bin/cleanup-tokens
```

Cleanup script:
```bash
#!/bin/bash
# cleanup-tokens.sh

# Delete expired and revoked refresh tokens
psql $DATABASE_URL -c "
    DELETE FROM refresh_tokens
    WHERE expires_at < NOW() OR is_revoked = true;
"

echo "Cleaned up expired refresh tokens at $(date)"
```

### Via Application

The repository provides a cleanup method:

```go
// Call periodically (e.g., via cron job or background worker)
deletedCount, err := refreshTokenRepository.DeleteExpired(ctx)
log.Info("cleaned up expired tokens", "count", deletedCount)
```

### Token Blacklist (Redis)

Token blacklist entries automatically expire based on their TTL. No manual cleanup needed.

To manually clear (use with caution):
```bash
redis-cli KEYS "token:blacklist:*" | xargs redis-cli DEL
```

---

## Audit Log Monitoring

### Query Audit Logs

```sql
-- Recent login failures
SELECT * FROM audit_logs
WHERE event_type = 'LoginFailed'
ORDER BY created_at DESC
LIMIT 100;

-- Permission denied events
SELECT * FROM audit_logs
WHERE event_type = 'PermissionDenied'
ORDER BY created_at DESC
LIMIT 100;

-- All events for a specific user
SELECT * FROM audit_logs
WHERE user_id = 'user-uuid'
ORDER BY created_at DESC;

-- Events by IP address
SELECT * FROM audit_logs
WHERE metadata->>'ip_address' = '192.168.1.100'
ORDER BY created_at DESC;

-- Account lockouts in last 24 hours
SELECT * FROM audit_logs
WHERE event_type = 'AccountLocked'
AND created_at > NOW() - INTERVAL '24 hours';
```

### Monitoring Alerts

Set up alerts for:

1. **High login failure rate**
   ```sql
   SELECT COUNT(*) FROM audit_logs
   WHERE event_type = 'LoginFailed'
   AND created_at > NOW() - INTERVAL '5 minutes';
   -- Alert if > 100
   ```

2. **Multiple account lockouts**
   ```sql
   SELECT COUNT(*) FROM audit_logs
   WHERE event_type = 'AccountLocked'
   AND created_at > NOW() - INTERVAL '1 hour';
   -- Alert if > 10
   ```

3. **Permission denied spikes**
   ```sql
   SELECT COUNT(*) FROM audit_logs
   WHERE event_type = 'PermissionDenied'
   AND created_at > NOW() - INTERVAL '5 minutes';
   -- Alert if > 50
   ```

### Log Retention

Implement retention policy:

```sql
-- Delete logs older than 90 days
DELETE FROM audit_logs
WHERE created_at < NOW() - INTERVAL '90 days';
```

Schedule as daily cron job.

---

## Incident Response

### Suspected Account Compromise

1. **Immediately revoke all sessions**
   ```bash
   curl -X POST http://localhost:8080/api/v1/auth/logout-all \
     -H "Authorization: Bearer <user_or_admin_token>"
   ```

   Or via database:
   ```sql
   UPDATE refresh_tokens SET is_revoked = true
   WHERE user_id = 'compromised-user-uuid';
   ```

2. **Lock the account**
   ```sql
   UPDATE users SET status = 'inactive'
   WHERE id = 'compromised-user-uuid';
   ```

3. **Review audit logs**
   ```sql
   SELECT * FROM audit_logs
   WHERE user_id = 'compromised-user-uuid'
   ORDER BY created_at DESC;
   ```

4. **Force password reset** (when user recovers)

### Suspected API Key/Token Leak

1. **Identify affected tokens**
   ```sql
   SELECT * FROM refresh_tokens
   WHERE created_at > 'leak-suspected-time';
   ```

2. **Revoke all potentially affected tokens**
   ```sql
   UPDATE refresh_tokens SET is_revoked = true
   WHERE created_at > 'leak-suspected-time';
   ```

3. **If access token leaked, add to blacklist** (requires token ID)

### Brute Force Attack Detection

1. **Identify attacking IPs**
   ```sql
   SELECT metadata->>'ip_address', COUNT(*) as attempts
   FROM audit_logs
   WHERE event_type = 'LoginFailed'
   AND created_at > NOW() - INTERVAL '1 hour'
   GROUP BY metadata->>'ip_address'
   ORDER BY attempts DESC;
   ```

2. **Block at firewall/WAF level**

3. **Consider temporary rate limit reduction**

### Privilege Escalation Attempt

1. **Review permission denied logs**
   ```sql
   SELECT user_id, metadata->>'permission', COUNT(*)
   FROM audit_logs
   WHERE event_type = 'PermissionDenied'
   AND created_at > NOW() - INTERVAL '1 hour'
   GROUP BY user_id, metadata->>'permission'
   ORDER BY COUNT(*) DESC;
   ```

2. **Verify user's assigned roles are correct**
   ```sql
   SELECT u.email, r.name
   FROM users u
   JOIN user_roles ur ON u.id = ur.user_id
   JOIN roles r ON ur.role_id = r.id
   WHERE u.id = 'suspicious-user-uuid';
   ```

3. **Review recent role assignments**
   ```sql
   SELECT * FROM audit_logs
   WHERE event_type IN ('UserRoleAssigned', 'UserRoleRevoked')
   AND created_at > NOW() - INTERVAL '24 hours';
   ```

---

## Environment Variables Reference

```bash
# JWT Configuration
JWT_SECRET_KEY=your-secret-key
JWT_ACCESS_TOKEN_TTL=15m
JWT_REFRESH_TOKEN_TTL=168h  # 7 days
JWT_ISSUER=your-app-name
JWT_AUDIENCE=your-app-name

# Account Lockout
LOCKOUT_MAX_ATTEMPTS=5
LOCKOUT_DURATION=15m
LOCKOUT_ATTEMPT_WINDOW=15m

# Rate Limiting
RATE_LIMIT_LOGIN_RPS=1
RATE_LIMIT_LOGIN_BURST=5
RATE_LIMIT_REGISTER_RPS=1
RATE_LIMIT_REGISTER_BURST=3

# Redis (for lockout and token blacklist)
REDIS_URL=redis://localhost:6379

# Database
DATABASE_URL=postgres://user:pass@localhost:5432/dbname
```
