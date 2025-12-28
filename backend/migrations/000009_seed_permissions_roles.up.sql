-- Seed System Permissions
INSERT INTO permissions (id, resource, action, description, is_system) VALUES
    -- User permissions
    ('a0000000-0000-0000-0000-000000000001', 'users', 'create', 'Create new users', TRUE),
    ('a0000000-0000-0000-0000-000000000002', 'users', 'read', 'View user details', TRUE),
    ('a0000000-0000-0000-0000-000000000003', 'users', 'update', 'Update user information', TRUE),
    ('a0000000-0000-0000-0000-000000000004', 'users', 'delete', 'Delete users', TRUE),
    ('a0000000-0000-0000-0000-000000000005', 'users', 'list', 'List all users', TRUE),
    ('a0000000-0000-0000-0000-000000000006', 'users', 'manage', 'Full user management access', TRUE),
    -- Role permissions
    ('a0000000-0000-0000-0000-000000000007', 'roles', 'create', 'Create new roles', TRUE),
    ('a0000000-0000-0000-0000-000000000008', 'roles', 'read', 'View role details', TRUE),
    ('a0000000-0000-0000-0000-000000000009', 'roles', 'update', 'Update role information', TRUE),
    ('a0000000-0000-0000-0000-000000000010', 'roles', 'delete', 'Delete roles', TRUE),
    ('a0000000-0000-0000-0000-000000000011', 'roles', 'list', 'List all roles', TRUE),
    ('a0000000-0000-0000-0000-000000000012', 'roles', 'assign', 'Assign roles to users', TRUE),
    -- Permission permissions
    ('a0000000-0000-0000-0000-000000000013', 'permissions', 'read', 'View permission details', TRUE),
    ('a0000000-0000-0000-0000-000000000014', 'permissions', 'list', 'List all permissions', TRUE),
    -- System permissions
    ('a0000000-0000-0000-0000-000000000015', 'system', 'admin', 'Full system administration access', TRUE)
ON CONFLICT (resource, action) DO NOTHING;

-- Seed System Roles
INSERT INTO roles (id, name, display_name, description, is_system, is_default, priority) VALUES
    ('b0000000-0000-0000-0000-000000000001', 'super_admin', 'Super Administrator', 'Full system access with all permissions', TRUE, FALSE, 100),
    ('b0000000-0000-0000-0000-000000000002', 'admin', 'Administrator', 'Administrative access for user and role management', TRUE, FALSE, 80),
    ('b0000000-0000-0000-0000-000000000003', 'manager', 'Manager', 'Management access with read and limited update permissions', TRUE, FALSE, 50),
    ('b0000000-0000-0000-0000-000000000004', 'user', 'User', 'Basic user access for own data', TRUE, TRUE, 10)
ON CONFLICT (name) DO NOTHING;

-- Assign all permissions to super_admin
INSERT INTO role_permissions (role_id, permission_id)
SELECT 'b0000000-0000-0000-0000-000000000001', id FROM permissions
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Assign user and role management permissions to admin
INSERT INTO role_permissions (role_id, permission_id) VALUES
    ('b0000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000001'), -- users:create
    ('b0000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000002'), -- users:read
    ('b0000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000003'), -- users:update
    ('b0000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000004'), -- users:delete
    ('b0000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000005'), -- users:list
    ('b0000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000006'), -- users:manage
    ('b0000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000008'), -- roles:read
    ('b0000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000011'), -- roles:list
    ('b0000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000012'), -- roles:assign
    ('b0000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000013'), -- permissions:read
    ('b0000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000014')  -- permissions:list
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Assign read and limited update permissions to manager
INSERT INTO role_permissions (role_id, permission_id) VALUES
    ('b0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000002'), -- users:read
    ('b0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000003'), -- users:update
    ('b0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000005'), -- users:list
    ('b0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000008'), -- roles:read
    ('b0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000011'), -- roles:list
    ('b0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000013'), -- permissions:read
    ('b0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000014')  -- permissions:list
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Assign basic read permission to user role (own data only - enforced by application)
INSERT INTO role_permissions (role_id, permission_id) VALUES
    ('b0000000-0000-0000-0000-000000000004', 'a0000000-0000-0000-0000-000000000002') -- users:read (own)
ON CONFLICT (role_id, permission_id) DO NOTHING;
