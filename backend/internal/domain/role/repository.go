package role

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(context context.Context, role *Role) error
	Update(context context.Context, role *Role) error
	Delete(context context.Context, id uuid.UUID) error
	FindByID(context context.Context, id uuid.UUID) (*Role, error)
	FindByName(context context.Context, name string) (*Role, error)
	FindByIDs(context context.Context, ids []uuid.UUID) ([]*Role, error)
	FindAll(context context.Context) ([]*Role, error)
	FindDefault(context context.Context) (*Role, error)
	ExistsByName(context context.Context, name string) (bool, error)
	FindByPermission(context context.Context, permissionID uuid.UUID) ([]*Role, error)
}
