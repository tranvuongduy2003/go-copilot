package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateRoleRequest struct {
	Name          string      `json:"name" validate:"required,min=1,max=100"`
	DisplayName   string      `json:"display_name" validate:"required,min=1,max=255"`
	Description   string      `json:"description" validate:"omitempty,max=500"`
	PermissionIDs []uuid.UUID `json:"permission_ids" validate:"omitempty,dive,uuid4"`
}

type UpdateRoleRequest struct {
	DisplayName string `json:"display_name" validate:"omitempty,min=1,max=255"`
	Description string `json:"description" validate:"omitempty,max=500"`
}

type SetRolePermissionsRequest struct {
	PermissionIDs []uuid.UUID `json:"permission_ids" validate:"required,dive,uuid4"`
}

type SetUserRolesRequest struct {
	RoleIDs []uuid.UUID `json:"role_ids" validate:"required,dive,uuid4"`
}

type RoleResponse struct {
	ID            uuid.UUID   `json:"id"`
	Name          string      `json:"name"`
	DisplayName   string      `json:"display_name"`
	Description   string      `json:"description"`
	PermissionIDs []uuid.UUID `json:"permission_ids"`
	IsSystem      bool        `json:"is_system"`
	IsDefault     bool        `json:"is_default"`
	Priority      int         `json:"priority"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

type RoleWithPermissionsResponse struct {
	ID            uuid.UUID   `json:"id"`
	Name          string      `json:"name"`
	DisplayName   string      `json:"display_name"`
	Description   string      `json:"description"`
	Permissions   []string    `json:"permissions"`
	PermissionIDs []uuid.UUID `json:"permission_ids"`
	IsSystem      bool        `json:"is_system"`
	IsDefault     bool        `json:"is_default"`
	Priority      int         `json:"priority"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}
