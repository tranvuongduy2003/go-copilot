package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/permission"
	"github.com/tranvuongduy2003/go-copilot/internal/infrastructure/persistence/postgres"
)

const (
	queryInsertPermission = `
		INSERT INTO permissions (id, resource, action, description, is_system, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	queryUpdatePermission = `
		UPDATE permissions
		SET description = $2, updated_at = $3
		WHERE id = $1`

	queryDeletePermission = `
		DELETE FROM permissions WHERE id = $1`

	queryFindPermissionByID = `
		SELECT id, resource, action, description, is_system, created_at, updated_at
		FROM permissions
		WHERE id = $1`

	queryFindPermissionByCode = `
		SELECT id, resource, action, description, is_system, created_at, updated_at
		FROM permissions
		WHERE resource = $1 AND action = $2`

	queryFindPermissionByCodeString = `
		SELECT id, resource, action, description, is_system, created_at, updated_at
		FROM permissions
		WHERE CONCAT(resource, ':', action) = $1`

	queryFindPermissionsByResource = `
		SELECT id, resource, action, description, is_system, created_at, updated_at
		FROM permissions
		WHERE resource = $1
		ORDER BY action`

	queryFindAllPermissions = `
		SELECT id, resource, action, description, is_system, created_at, updated_at
		FROM permissions
		ORDER BY resource, action`

	queryFindPermissionsByIDs = `
		SELECT id, resource, action, description, is_system, created_at, updated_at
		FROM permissions
		WHERE id = ANY($1)`

	queryExistsPermissionByCode = `
		SELECT EXISTS(SELECT 1 FROM permissions WHERE resource = $1 AND action = $2)`
)

type permissionRow struct {
	ID          uuid.UUID
	Resource    string
	Action      string
	Description *string
	IsSystem    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (r *permissionRow) toDomain() (*permission.Permission, error) {
	description := ""
	if r.Description != nil {
		description = *r.Description
	}
	return permission.ReconstructPermission(permission.ReconstructPermissionParams{
		ID:          r.ID,
		Resource:    r.Resource,
		Action:      r.Action,
		Description: description,
		IsSystem:    r.IsSystem,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	})
}

func permissionToRow(p *permission.Permission) *permissionRow {
	desc := p.Description()
	return &permissionRow{
		ID:          p.ID(),
		Resource:    p.Resource().String(),
		Action:      p.Action().String(),
		Description: &desc,
		IsSystem:    p.IsSystem(),
		CreatedAt:   p.CreatedAt(),
		UpdatedAt:   p.UpdatedAt(),
	}
}

type PermissionRepository struct {
	pool *pgxpool.Pool
}

func NewPermissionRepository(pool *pgxpool.Pool) *PermissionRepository {
	return &PermissionRepository{pool: pool}
}

func (r *PermissionRepository) Create(ctx context.Context, p *permission.Permission) error {
	querier := postgres.GetQuerier(ctx, r.pool)
	row := permissionToRow(p)

	_, err := querier.Exec(ctx, queryInsertPermission,
		row.ID,
		row.Resource,
		row.Action,
		row.Description,
		row.IsSystem,
		row.CreatedAt,
		row.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolationCode {
			return permission.NewPermissionCodeExistsError(p.CodeString())
		}
		return postgres.NewDBError("create permission", err)
	}

	return nil
}

func (r *PermissionRepository) Update(ctx context.Context, p *permission.Permission) error {
	querier := postgres.GetQuerier(ctx, r.pool)

	cmdTag, err := querier.Exec(ctx, queryUpdatePermission,
		p.ID(),
		p.Description(),
		p.UpdatedAt(),
	)
	if err != nil {
		return postgres.NewDBError("update permission", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return permission.NewPermissionNotFoundError(p.ID().String())
	}

	return nil
}

func (r *PermissionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	querier := postgres.GetQuerier(ctx, r.pool)

	cmdTag, err := querier.Exec(ctx, queryDeletePermission, id)
	if err != nil {
		return postgres.NewDBError("delete permission", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return permission.NewPermissionNotFoundError(id.String())
	}

	return nil
}

func (r *PermissionRepository) FindByID(ctx context.Context, id uuid.UUID) (*permission.Permission, error) {
	querier := postgres.GetQuerier(ctx, r.pool)

	row := &permissionRow{}
	err := querier.QueryRow(ctx, queryFindPermissionByID, id).Scan(
		&row.ID,
		&row.Resource,
		&row.Action,
		&row.Description,
		&row.IsSystem,
		&row.CreatedAt,
		&row.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, permission.NewPermissionNotFoundError(id.String())
		}
		return nil, postgres.NewDBError("find permission by id", err)
	}

	return row.toDomain()
}

func (r *PermissionRepository) FindByCode(ctx context.Context, code permission.PermissionCode) (*permission.Permission, error) {
	querier := postgres.GetQuerier(ctx, r.pool)

	row := &permissionRow{}
	err := querier.QueryRow(ctx, queryFindPermissionByCode, code.Resource().String(), code.Action().String()).Scan(
		&row.ID,
		&row.Resource,
		&row.Action,
		&row.Description,
		&row.IsSystem,
		&row.CreatedAt,
		&row.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, permission.NewPermissionNotFoundError(code.String())
		}
		return nil, postgres.NewDBError("find permission by code", err)
	}

	return row.toDomain()
}

func (r *PermissionRepository) FindByCodeString(ctx context.Context, code string) (*permission.Permission, error) {
	querier := postgres.GetQuerier(ctx, r.pool)

	row := &permissionRow{}
	err := querier.QueryRow(ctx, queryFindPermissionByCodeString, code).Scan(
		&row.ID,
		&row.Resource,
		&row.Action,
		&row.Description,
		&row.IsSystem,
		&row.CreatedAt,
		&row.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, permission.NewPermissionNotFoundError(code)
		}
		return nil, postgres.NewDBError("find permission by code string", err)
	}

	return row.toDomain()
}

func (r *PermissionRepository) FindByResource(ctx context.Context, resource permission.Resource) ([]*permission.Permission, error) {
	querier := postgres.GetQuerier(ctx, r.pool)

	rows, err := querier.Query(ctx, queryFindPermissionsByResource, resource.String())
	if err != nil {
		return nil, postgres.NewDBError("find permissions by resource", err)
	}
	defer rows.Close()

	permissions := make([]*permission.Permission, 0)
	for rows.Next() {
		row := &permissionRow{}
		err := rows.Scan(
			&row.ID,
			&row.Resource,
			&row.Action,
			&row.Description,
			&row.IsSystem,
			&row.CreatedAt,
			&row.UpdatedAt,
		)
		if err != nil {
			return nil, postgres.NewDBError("scan permission row", err)
		}

		p, err := row.toDomain()
		if err != nil {
			return nil, postgres.NewDBError("convert permission row to domain", err)
		}
		permissions = append(permissions, p)
	}

	if err := rows.Err(); err != nil {
		return nil, postgres.NewDBError("iterate permission rows", err)
	}

	return permissions, nil
}

func (r *PermissionRepository) FindAll(ctx context.Context) ([]*permission.Permission, error) {
	querier := postgres.GetQuerier(ctx, r.pool)

	rows, err := querier.Query(ctx, queryFindAllPermissions)
	if err != nil {
		return nil, postgres.NewDBError("find all permissions", err)
	}
	defer rows.Close()

	permissions := make([]*permission.Permission, 0)
	for rows.Next() {
		row := &permissionRow{}
		err := rows.Scan(
			&row.ID,
			&row.Resource,
			&row.Action,
			&row.Description,
			&row.IsSystem,
			&row.CreatedAt,
			&row.UpdatedAt,
		)
		if err != nil {
			return nil, postgres.NewDBError("scan permission row", err)
		}

		p, err := row.toDomain()
		if err != nil {
			return nil, postgres.NewDBError("convert permission row to domain", err)
		}
		permissions = append(permissions, p)
	}

	if err := rows.Err(); err != nil {
		return nil, postgres.NewDBError("iterate permission rows", err)
	}

	return permissions, nil
}

func (r *PermissionRepository) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]*permission.Permission, error) {
	if len(ids) == 0 {
		return []*permission.Permission{}, nil
	}

	querier := postgres.GetQuerier(ctx, r.pool)

	rows, err := querier.Query(ctx, queryFindPermissionsByIDs, ids)
	if err != nil {
		return nil, postgres.NewDBError("find permissions by ids", err)
	}
	defer rows.Close()

	permissions := make([]*permission.Permission, 0)
	for rows.Next() {
		row := &permissionRow{}
		err := rows.Scan(
			&row.ID,
			&row.Resource,
			&row.Action,
			&row.Description,
			&row.IsSystem,
			&row.CreatedAt,
			&row.UpdatedAt,
		)
		if err != nil {
			return nil, postgres.NewDBError("scan permission row", err)
		}

		p, err := row.toDomain()
		if err != nil {
			return nil, postgres.NewDBError("convert permission row to domain", err)
		}
		permissions = append(permissions, p)
	}

	if err := rows.Err(); err != nil {
		return nil, postgres.NewDBError("iterate permission rows", err)
	}

	return permissions, nil
}

func (r *PermissionRepository) ExistsByCode(ctx context.Context, code permission.PermissionCode) (bool, error) {
	querier := postgres.GetQuerier(ctx, r.pool)

	var exists bool
	err := querier.QueryRow(ctx, queryExistsPermissionByCode, code.Resource().String(), code.Action().String()).Scan(&exists)
	if err != nil {
		return false, postgres.NewDBError("check permission exists by code", err)
	}

	return exists, nil
}
