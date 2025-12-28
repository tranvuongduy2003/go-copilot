package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/role"
	"github.com/tranvuongduy2003/go-copilot/internal/infrastructure/persistence/postgres"
)

const (
	queryInsertRole = `
		INSERT INTO roles (id, name, display_name, description, is_system, is_default, priority, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	queryUpdateRole = `
		UPDATE roles
		SET display_name = $2, description = $3, is_default = $4, priority = $5, updated_at = $6
		WHERE id = $1`

	queryDeleteRole = `
		DELETE FROM roles WHERE id = $1`

	queryFindRoleByID = `
		SELECT id, name, display_name, description, is_system, is_default, priority, created_at, updated_at
		FROM roles
		WHERE id = $1`

	queryFindRoleByName = `
		SELECT id, name, display_name, description, is_system, is_default, priority, created_at, updated_at
		FROM roles
		WHERE name = $1`

	queryFindRolesByIDs = `
		SELECT id, name, display_name, description, is_system, is_default, priority, created_at, updated_at
		FROM roles
		WHERE id = ANY($1)`

	queryFindAllRoles = `
		SELECT id, name, display_name, description, is_system, is_default, priority, created_at, updated_at
		FROM roles
		ORDER BY priority DESC, name`

	queryFindDefaultRole = `
		SELECT id, name, display_name, description, is_system, is_default, priority, created_at, updated_at
		FROM roles
		WHERE is_default = TRUE
		LIMIT 1`

	queryExistsRoleByName = `
		SELECT EXISTS(SELECT 1 FROM roles WHERE name = $1)`

	queryFindRolesByPermission = `
		SELECT r.id, r.name, r.display_name, r.description, r.is_system, r.is_default, r.priority, r.created_at, r.updated_at
		FROM roles r
		INNER JOIN role_permissions rp ON r.id = rp.role_id
		WHERE rp.permission_id = $1
		ORDER BY r.priority DESC, r.name`

	queryFindRolePermissions = `
		SELECT permission_id FROM role_permissions WHERE role_id = $1`

	queryDeleteRolePermissions = `
		DELETE FROM role_permissions WHERE role_id = $1`

	queryInsertRolePermission = `
		INSERT INTO role_permissions (role_id, permission_id, created_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (role_id, permission_id) DO NOTHING`
)

type roleRow struct {
	ID          uuid.UUID
	Name        string
	DisplayName string
	Description *string
	IsSystem    bool
	IsDefault   bool
	Priority    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (r *roleRow) toDomain(permissionIDs []uuid.UUID) (*role.Role, error) {
	description := ""
	if r.Description != nil {
		description = *r.Description
	}
	return role.ReconstructRole(role.ReconstructRoleParams{
		ID:            r.ID,
		Name:          r.Name,
		DisplayName:   r.DisplayName,
		Description:   description,
		PermissionIDs: permissionIDs,
		IsSystem:      r.IsSystem,
		IsDefault:     r.IsDefault,
		Priority:      r.Priority,
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
	})
}

func roleToRow(r *role.Role) *roleRow {
	desc := r.Description()
	return &roleRow{
		ID:          r.ID(),
		Name:        r.Name(),
		DisplayName: r.DisplayName(),
		Description: &desc,
		IsSystem:    r.IsSystem(),
		IsDefault:   r.IsDefault(),
		Priority:    r.Priority(),
		CreatedAt:   r.CreatedAt(),
		UpdatedAt:   r.UpdatedAt(),
	}
}

type RoleRepository struct {
	pool *pgxpool.Pool
}

func NewRoleRepository(pool *pgxpool.Pool) *RoleRepository {
	return &RoleRepository{pool: pool}
}

func (r *RoleRepository) Create(ctx context.Context, rl *role.Role) error {
	querier := postgres.GetQuerier(ctx, r.pool)
	row := roleToRow(rl)

	_, err := querier.Exec(ctx, queryInsertRole,
		row.ID,
		row.Name,
		row.DisplayName,
		row.Description,
		row.IsSystem,
		row.IsDefault,
		row.Priority,
		row.CreatedAt,
		row.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolationCode {
			return role.NewRoleNameExistsError(row.Name)
		}
		return postgres.NewDBError("create role", err)
	}

	if err := r.syncPermissions(ctx, querier, rl.ID(), rl.PermissionIDs()); err != nil {
		return err
	}

	return nil
}

func (r *RoleRepository) Update(ctx context.Context, rl *role.Role) error {
	querier := postgres.GetQuerier(ctx, r.pool)
	row := roleToRow(rl)

	cmdTag, err := querier.Exec(ctx, queryUpdateRole,
		row.ID,
		row.DisplayName,
		row.Description,
		row.IsDefault,
		row.Priority,
		row.UpdatedAt,
	)
	if err != nil {
		return postgres.NewDBError("update role", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return role.NewRoleNotFoundError(row.ID.String())
	}

	if err := r.syncPermissions(ctx, querier, rl.ID(), rl.PermissionIDs()); err != nil {
		return err
	}

	return nil
}

func (r *RoleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	querier := postgres.GetQuerier(ctx, r.pool)

	cmdTag, err := querier.Exec(ctx, queryDeleteRole, id)
	if err != nil {
		return postgres.NewDBError("delete role", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return role.NewRoleNotFoundError(id.String())
	}

	return nil
}

func (r *RoleRepository) FindByID(ctx context.Context, id uuid.UUID) (*role.Role, error) {
	querier := postgres.GetQuerier(ctx, r.pool)

	row := &roleRow{}
	err := querier.QueryRow(ctx, queryFindRoleByID, id).Scan(
		&row.ID,
		&row.Name,
		&row.DisplayName,
		&row.Description,
		&row.IsSystem,
		&row.IsDefault,
		&row.Priority,
		&row.CreatedAt,
		&row.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, role.NewRoleNotFoundError(id.String())
		}
		return nil, postgres.NewDBError("find role by id", err)
	}

	permissionIDs, err := r.loadPermissionIDs(ctx, querier, id)
	if err != nil {
		return nil, err
	}

	return row.toDomain(permissionIDs)
}

func (r *RoleRepository) FindByName(ctx context.Context, name string) (*role.Role, error) {
	querier := postgres.GetQuerier(ctx, r.pool)

	row := &roleRow{}
	err := querier.QueryRow(ctx, queryFindRoleByName, name).Scan(
		&row.ID,
		&row.Name,
		&row.DisplayName,
		&row.Description,
		&row.IsSystem,
		&row.IsDefault,
		&row.Priority,
		&row.CreatedAt,
		&row.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, role.NewRoleNotFoundError(name)
		}
		return nil, postgres.NewDBError("find role by name", err)
	}

	permissionIDs, err := r.loadPermissionIDs(ctx, querier, row.ID)
	if err != nil {
		return nil, err
	}

	return row.toDomain(permissionIDs)
}

func (r *RoleRepository) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]*role.Role, error) {
	if len(ids) == 0 {
		return []*role.Role{}, nil
	}

	querier := postgres.GetQuerier(ctx, r.pool)

	rows, err := querier.Query(ctx, queryFindRolesByIDs, ids)
	if err != nil {
		return nil, postgres.NewDBError("find roles by ids", err)
	}
	defer rows.Close()

	roles := make([]*role.Role, 0)
	for rows.Next() {
		row := &roleRow{}
		err := rows.Scan(
			&row.ID,
			&row.Name,
			&row.DisplayName,
			&row.Description,
			&row.IsSystem,
			&row.IsDefault,
			&row.Priority,
			&row.CreatedAt,
			&row.UpdatedAt,
		)
		if err != nil {
			return nil, postgres.NewDBError("scan role row", err)
		}

		permissionIDs, err := r.loadPermissionIDs(ctx, querier, row.ID)
		if err != nil {
			return nil, err
		}

		rl, err := row.toDomain(permissionIDs)
		if err != nil {
			return nil, postgres.NewDBError("convert role row to domain", err)
		}
		roles = append(roles, rl)
	}

	if err := rows.Err(); err != nil {
		return nil, postgres.NewDBError("iterate role rows", err)
	}

	return roles, nil
}

func (r *RoleRepository) FindAll(ctx context.Context) ([]*role.Role, error) {
	querier := postgres.GetQuerier(ctx, r.pool)

	rows, err := querier.Query(ctx, queryFindAllRoles)
	if err != nil {
		return nil, postgres.NewDBError("find all roles", err)
	}
	defer rows.Close()

	roles := make([]*role.Role, 0)
	for rows.Next() {
		row := &roleRow{}
		err := rows.Scan(
			&row.ID,
			&row.Name,
			&row.DisplayName,
			&row.Description,
			&row.IsSystem,
			&row.IsDefault,
			&row.Priority,
			&row.CreatedAt,
			&row.UpdatedAt,
		)
		if err != nil {
			return nil, postgres.NewDBError("scan role row", err)
		}

		permissionIDs, err := r.loadPermissionIDs(ctx, querier, row.ID)
		if err != nil {
			return nil, err
		}

		rl, err := row.toDomain(permissionIDs)
		if err != nil {
			return nil, postgres.NewDBError("convert role row to domain", err)
		}
		roles = append(roles, rl)
	}

	if err := rows.Err(); err != nil {
		return nil, postgres.NewDBError("iterate role rows", err)
	}

	return roles, nil
}

func (r *RoleRepository) FindDefault(ctx context.Context) (*role.Role, error) {
	querier := postgres.GetQuerier(ctx, r.pool)

	row := &roleRow{}
	err := querier.QueryRow(ctx, queryFindDefaultRole).Scan(
		&row.ID,
		&row.Name,
		&row.DisplayName,
		&row.Description,
		&row.IsSystem,
		&row.IsDefault,
		&row.Priority,
		&row.CreatedAt,
		&row.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, role.ErrNoDefaultRole
		}
		return nil, postgres.NewDBError("find default role", err)
	}

	permissionIDs, err := r.loadPermissionIDs(ctx, querier, row.ID)
	if err != nil {
		return nil, err
	}

	return row.toDomain(permissionIDs)
}

func (r *RoleRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	querier := postgres.GetQuerier(ctx, r.pool)

	var exists bool
	err := querier.QueryRow(ctx, queryExistsRoleByName, name).Scan(&exists)
	if err != nil {
		return false, postgres.NewDBError("check role exists by name", err)
	}

	return exists, nil
}

func (r *RoleRepository) FindByPermission(ctx context.Context, permissionID uuid.UUID) ([]*role.Role, error) {
	querier := postgres.GetQuerier(ctx, r.pool)

	rows, err := querier.Query(ctx, queryFindRolesByPermission, permissionID)
	if err != nil {
		return nil, postgres.NewDBError("find roles by permission", err)
	}
	defer rows.Close()

	roles := make([]*role.Role, 0)
	for rows.Next() {
		row := &roleRow{}
		err := rows.Scan(
			&row.ID,
			&row.Name,
			&row.DisplayName,
			&row.Description,
			&row.IsSystem,
			&row.IsDefault,
			&row.Priority,
			&row.CreatedAt,
			&row.UpdatedAt,
		)
		if err != nil {
			return nil, postgres.NewDBError("scan role row", err)
		}

		permissionIDs, err := r.loadPermissionIDs(ctx, querier, row.ID)
		if err != nil {
			return nil, err
		}

		rl, err := row.toDomain(permissionIDs)
		if err != nil {
			return nil, postgres.NewDBError("convert role row to domain", err)
		}
		roles = append(roles, rl)
	}

	if err := rows.Err(); err != nil {
		return nil, postgres.NewDBError("iterate role rows", err)
	}

	return roles, nil
}

func (r *RoleRepository) loadPermissionIDs(ctx context.Context, querier postgres.Querier, roleID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := querier.Query(ctx, queryFindRolePermissions, roleID)
	if err != nil {
		return nil, postgres.NewDBError("load role permissions", err)
	}
	defer rows.Close()

	permissionIDs := make([]uuid.UUID, 0)
	for rows.Next() {
		var permissionID uuid.UUID
		if err := rows.Scan(&permissionID); err != nil {
			return nil, postgres.NewDBError("scan permission id", err)
		}
		permissionIDs = append(permissionIDs, permissionID)
	}

	if err := rows.Err(); err != nil {
		return nil, postgres.NewDBError("iterate permission ids", err)
	}

	return permissionIDs, nil
}

func (r *RoleRepository) syncPermissions(ctx context.Context, querier postgres.Querier, roleID uuid.UUID, permissionIDs []uuid.UUID) error {
	_, err := querier.Exec(ctx, queryDeleteRolePermissions, roleID)
	if err != nil {
		return postgres.NewDBError("delete role permissions", err)
	}

	now := time.Now().UTC()
	for _, permissionID := range permissionIDs {
		_, err := querier.Exec(ctx, queryInsertRolePermission, roleID, permissionID, now)
		if err != nil {
			return postgres.NewDBError("insert role permission", err)
		}
	}

	return nil
}
