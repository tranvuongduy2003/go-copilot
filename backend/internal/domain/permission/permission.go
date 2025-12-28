package permission

import (
	"time"

	"github.com/google/uuid"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/shared"
)

type Permission struct {
	shared.Entity
	resource    Resource
	action      Action
	description string
	isSystem    bool
	createdAt   time.Time
	updatedAt   time.Time
}

type NewPermissionParams struct {
	Resource    string
	Action      string
	Description string
	IsSystem    bool
}

func NewPermission(params NewPermissionParams) (*Permission, error) {
	resource, err := NewResource(params.Resource)
	if err != nil {
		return nil, err
	}

	action, err := NewAction(params.Action)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	return &Permission{
		Entity:      shared.NewEntity(),
		resource:    resource,
		action:      action,
		description: params.Description,
		isSystem:    params.IsSystem,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

type ReconstructPermissionParams struct {
	ID          uuid.UUID
	Resource    string
	Action      string
	Description string
	IsSystem    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func ReconstructPermission(params ReconstructPermissionParams) (*Permission, error) {
	resource, err := NewResource(params.Resource)
	if err != nil {
		return nil, err
	}

	action, err := NewAction(params.Action)
	if err != nil {
		return nil, err
	}

	return &Permission{
		Entity:      shared.NewEntityWithID(params.ID),
		resource:    resource,
		action:      action,
		description: params.Description,
		isSystem:    params.IsSystem,
		createdAt:   params.CreatedAt,
		updatedAt:   params.UpdatedAt,
	}, nil
}

func (p *Permission) Resource() Resource {
	return p.resource
}

func (p *Permission) Action() Action {
	return p.action
}

func (p *Permission) Description() string {
	return p.description
}

func (p *Permission) IsSystem() bool {
	return p.isSystem
}

func (p *Permission) CreatedAt() time.Time {
	return p.createdAt
}

func (p *Permission) UpdatedAt() time.Time {
	return p.updatedAt
}

func (p *Permission) Code() PermissionCode {
	return NewPermissionCode(p.resource, p.action)
}

func (p *Permission) CodeString() string {
	return p.Code().String()
}

func (p *Permission) UpdateDescription(description string) {
	if description != p.description {
		p.description = description
		p.updatedAt = time.Now().UTC()
	}
}

func (p *Permission) CanBeDeleted() bool {
	return !p.isSystem
}

func (p *Permission) Equals(other *Permission) bool {
	if other == nil {
		return false
	}
	return p.resource.Equals(other.resource) && p.action == other.action
}
