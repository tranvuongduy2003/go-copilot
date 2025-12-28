package permission

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(context context.Context, permission *Permission) error
	Update(context context.Context, permission *Permission) error
	Delete(context context.Context, id uuid.UUID) error
	FindByID(context context.Context, id uuid.UUID) (*Permission, error)
	FindByCode(context context.Context, code PermissionCode) (*Permission, error)
	FindByCodeString(context context.Context, code string) (*Permission, error)
	FindByResource(context context.Context, resource Resource) ([]*Permission, error)
	FindAll(context context.Context) ([]*Permission, error)
	FindByIDs(context context.Context, ids []uuid.UUID) ([]*Permission, error)
	ExistsByCode(context context.Context, code PermissionCode) (bool, error)
}
