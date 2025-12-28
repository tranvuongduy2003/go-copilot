DROP TRIGGER IF EXISTS trigger_single_default_role ON roles;
DROP FUNCTION IF EXISTS check_single_default_role();
DROP INDEX IF EXISTS idx_roles_priority;
DROP INDEX IF EXISTS idx_roles_is_default;
DROP INDEX IF EXISTS idx_roles_name;
DROP TABLE IF EXISTS roles;
