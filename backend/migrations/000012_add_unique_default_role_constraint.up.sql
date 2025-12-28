CREATE UNIQUE INDEX idx_roles_single_default ON roles (is_default) WHERE is_default = true;
