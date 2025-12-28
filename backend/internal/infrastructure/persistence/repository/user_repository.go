package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/shared"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/user"
	"github.com/tranvuongduy2003/go-copilot/internal/infrastructure/persistence/postgres"
)

const (
	usersTable = "users"

	queryInsertUser = `
		INSERT INTO users (id, email, password_hash, full_name, status, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	queryUpdateUser = `
		UPDATE users
		SET email = $2, password_hash = $3, full_name = $4, status = $5, updated_at = $6, deleted_at = $7
		WHERE id = $1 AND deleted_at IS NULL`

	querySoftDeleteUser = `
		UPDATE users
		SET deleted_at = $2, updated_at = $2
		WHERE id = $1 AND deleted_at IS NULL`

	queryFindUserByID = `
		SELECT id, email, password_hash, full_name, status, created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL`

	queryFindUserByEmail = `
		SELECT id, email, password_hash, full_name, status, created_at, updated_at, deleted_at
		FROM users
		WHERE LOWER(email) = LOWER($1) AND deleted_at IS NULL`

	queryExistsByEmail = `
		SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(email) = LOWER($1) AND deleted_at IS NULL)`

	queryCountUsers = `SELECT COUNT(*) FROM users`

	querySelectUsers = `
		SELECT id, email, password_hash, full_name, status, created_at, updated_at, deleted_at
		FROM users`

	queryFindUserRoles = `
		SELECT role_id FROM user_roles WHERE user_id = $1`

	queryDeleteUserRoles = `
		DELETE FROM user_roles WHERE user_id = $1`

	queryInsertUserRole = `
		INSERT INTO user_roles (user_id, role_id, assigned_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, role_id) DO NOTHING`

	queryFindUsersByRole = `
		SELECT u.id, u.email, u.password_hash, u.full_name, u.status, u.created_at, u.updated_at, u.deleted_at
		FROM users u
		INNER JOIN user_roles ur ON u.id = ur.user_id
		WHERE ur.role_id = $1 AND u.deleted_at IS NULL
		ORDER BY u.created_at DESC`
)

const pgUniqueViolationCode = "23505"

type userRow struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	FullName     string
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

func (r *userRow) toDomain(roleIDs []uuid.UUID) (*user.User, error) {
	return user.ReconstructUser(user.ReconstructUserParams{
		ID:           r.ID,
		Email:        r.Email,
		PasswordHash: r.PasswordHash,
		FullName:     r.FullName,
		Status:       user.Status(r.Status),
		RoleIDs:      roleIDs,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
		DeletedAt:    r.DeletedAt,
	})
}

func userToRow(u *user.User) *userRow {
	return &userRow{
		ID:           u.ID(),
		Email:        u.Email().String(),
		PasswordHash: u.PasswordHash().String(),
		FullName:     u.FullName().String(),
		Status:       u.Status().String(),
		CreatedAt:    u.CreatedAt(),
		UpdatedAt:    u.UpdatedAt(),
		DeletedAt:    u.DeletedAt(),
	}
}

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	querier := postgres.GetQuerier(ctx, r.pool)
	row := userToRow(u)

	_, err := querier.Exec(ctx, queryInsertUser,
		row.ID,
		row.Email,
		row.PasswordHash,
		row.FullName,
		row.Status,
		row.CreatedAt,
		row.UpdatedAt,
		row.DeletedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolationCode {
			return user.NewEmailAlreadyExistsError(row.Email)
		}
		return postgres.NewDBError("create user", err)
	}

	if err := r.syncRoles(ctx, querier, u.ID(), u.RoleIDs()); err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	querier := postgres.GetQuerier(ctx, r.pool)
	row := userToRow(u)

	cmdTag, err := querier.Exec(ctx, queryUpdateUser,
		row.ID,
		row.Email,
		row.PasswordHash,
		row.FullName,
		row.Status,
		row.UpdatedAt,
		row.DeletedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolationCode {
			return user.NewEmailAlreadyExistsError(row.Email)
		}
		return postgres.NewDBError("update user", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return user.NewUserNotFoundError(row.ID.String())
	}

	if err := r.syncRoles(ctx, querier, u.ID(), u.RoleIDs()); err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	querier := postgres.GetQuerier(ctx, r.pool)
	now := time.Now().UTC()

	cmdTag, err := querier.Exec(ctx, querySoftDeleteUser, id, now)
	if err != nil {
		return postgres.NewDBError("delete user", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return user.NewUserNotFoundError(id.String())
	}

	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	querier := postgres.GetQuerier(ctx, r.pool)

	row := &userRow{}
	err := querier.QueryRow(ctx, queryFindUserByID, id).Scan(
		&row.ID,
		&row.Email,
		&row.PasswordHash,
		&row.FullName,
		&row.Status,
		&row.CreatedAt,
		&row.UpdatedAt,
		&row.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user.NewUserNotFoundError(id.String())
		}
		return nil, postgres.NewDBError("find user by id", err)
	}

	roleIDs, err := r.loadRoleIDs(ctx, querier, id)
	if err != nil {
		return nil, err
	}

	return row.toDomain(roleIDs)
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	querier := postgres.GetQuerier(ctx, r.pool)

	row := &userRow{}
	err := querier.QueryRow(ctx, queryFindUserByEmail, email).Scan(
		&row.ID,
		&row.Email,
		&row.PasswordHash,
		&row.FullName,
		&row.Status,
		&row.CreatedAt,
		&row.UpdatedAt,
		&row.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user.NewUserNotFoundError(email)
		}
		return nil, postgres.NewDBError("find user by email", err)
	}

	roleIDs, err := r.loadRoleIDs(ctx, querier, row.ID)
	if err != nil {
		return nil, err
	}

	return row.toDomain(roleIDs)
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	querier := postgres.GetQuerier(ctx, r.pool)

	var exists bool
	err := querier.QueryRow(ctx, queryExistsByEmail, email).Scan(&exists)
	if err != nil {
		return false, postgres.NewDBError("check email exists", err)
	}

	return exists, nil
}

func (r *UserRepository) List(ctx context.Context, filter user.Filter, pagination shared.Pagination) ([]*user.User, int64, error) {
	querier := postgres.GetQuerier(ctx, r.pool)

	where := postgres.NewWhereClause()
	where.IsNull("deleted_at")

	if filter.Status != nil {
		where.Eq("status", filter.Status.String())
	}

	if filter.Search != nil && *filter.Search != "" {
		searchPattern := "%" + *filter.Search + "%"
		where.AddCondition("(LOWER(email) LIKE LOWER($%d) OR LOWER(full_name) LIKE LOWER($%d))", searchPattern, searchPattern)
	}

	if filter.DateRange.HasFrom() {
		where.Gte("created_at", *filter.DateRange.From())
	}

	if filter.DateRange.HasTo() {
		where.Lte("created_at", *filter.DateRange.To())
	}

	whereClause, args := where.Build()

	countQuery := queryCountUsers
	if whereClause != "" {
		countQuery = countQuery + " " + whereClause
	}

	var total int64
	err := querier.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, postgres.NewDBError("count users", err)
	}

	if total == 0 {
		return []*user.User{}, 0, nil
	}

	orderBy := postgres.NewOrderByClause().Desc("created_at")
	paginationClause := postgres.NewPaginationClauseFromOffset(pagination.Limit(), pagination.Offset())

	dataQuery := querySelectUsers
	if whereClause != "" {
		dataQuery = dataQuery + " " + whereClause
	}
	dataQuery = dataQuery + " " + orderBy.Build() + " " + paginationClause.Build()

	rows, err := querier.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, postgres.NewDBError("list users", err)
	}
	defer rows.Close()

	users := make([]*user.User, 0)
	for rows.Next() {
		row := &userRow{}
		err := rows.Scan(
			&row.ID,
			&row.Email,
			&row.PasswordHash,
			&row.FullName,
			&row.Status,
			&row.CreatedAt,
			&row.UpdatedAt,
			&row.DeletedAt,
		)
		if err != nil {
			return nil, 0, postgres.NewDBError("scan user row", err)
		}

		roleIDs, err := r.loadRoleIDs(ctx, querier, row.ID)
		if err != nil {
			return nil, 0, err
		}

		u, err := row.toDomain(roleIDs)
		if err != nil {
			return nil, 0, postgres.NewDBError("convert user row to domain", err)
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, postgres.NewDBError("iterate user rows", err)
	}

	return users, total, nil
}

func (r *UserRepository) FindByRole(ctx context.Context, roleID uuid.UUID) ([]*user.User, error) {
	querier := postgres.GetQuerier(ctx, r.pool)

	rows, err := querier.Query(ctx, queryFindUsersByRole, roleID)
	if err != nil {
		return nil, postgres.NewDBError("find users by role", err)
	}
	defer rows.Close()

	users := make([]*user.User, 0)
	for rows.Next() {
		row := &userRow{}
		err := rows.Scan(
			&row.ID,
			&row.Email,
			&row.PasswordHash,
			&row.FullName,
			&row.Status,
			&row.CreatedAt,
			&row.UpdatedAt,
			&row.DeletedAt,
		)
		if err != nil {
			return nil, postgres.NewDBError("scan user row", err)
		}

		roleIDs, err := r.loadRoleIDs(ctx, querier, row.ID)
		if err != nil {
			return nil, err
		}

		u, err := row.toDomain(roleIDs)
		if err != nil {
			return nil, postgres.NewDBError("convert user row to domain", err)
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, postgres.NewDBError("iterate user rows", err)
	}

	return users, nil
}

func (r *UserRepository) loadRoleIDs(ctx context.Context, querier postgres.Querier, userID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := querier.Query(ctx, queryFindUserRoles, userID)
	if err != nil {
		return nil, postgres.NewDBError("load user roles", err)
	}
	defer rows.Close()

	roleIDs := make([]uuid.UUID, 0)
	for rows.Next() {
		var roleID uuid.UUID
		if err := rows.Scan(&roleID); err != nil {
			return nil, postgres.NewDBError("scan role id", err)
		}
		roleIDs = append(roleIDs, roleID)
	}

	if err := rows.Err(); err != nil {
		return nil, postgres.NewDBError("iterate role ids", err)
	}

	return roleIDs, nil
}

func (r *UserRepository) syncRoles(ctx context.Context, querier postgres.Querier, userID uuid.UUID, roleIDs []uuid.UUID) error {
	_, err := querier.Exec(ctx, queryDeleteUserRoles, userID)
	if err != nil {
		return postgres.NewDBError("delete user roles", err)
	}

	now := time.Now().UTC()
	for _, roleID := range roleIDs {
		_, err := querier.Exec(ctx, queryInsertUserRole, userID, roleID, now)
		if err != nil {
			return postgres.NewDBError("insert user role", err)
		}
	}

	return nil
}
