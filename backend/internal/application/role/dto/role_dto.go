package roledto

import (
	"time"

	"github.com/google/uuid"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/role"
)

type RoleDTO struct {
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

func RoleFromDomain(domainRole *role.Role) *RoleDTO {
	if domainRole == nil {
		return nil
	}
	return &RoleDTO{
		ID:            domainRole.ID(),
		Name:          domainRole.Name(),
		DisplayName:   domainRole.DisplayName(),
		Description:   domainRole.Description(),
		PermissionIDs: domainRole.PermissionIDs(),
		IsSystem:      domainRole.IsSystem(),
		IsDefault:     domainRole.IsDefault(),
		Priority:      domainRole.Priority(),
		CreatedAt:     domainRole.CreatedAt(),
		UpdatedAt:     domainRole.UpdatedAt(),
	}
}

func RolesFromDomain(domainRoles []*role.Role) []*RoleDTO {
	dtos := make([]*RoleDTO, len(domainRoles))
	for i, domainRole := range domainRoles {
		dtos[i] = RoleFromDomain(domainRole)
	}
	return dtos
}

type RoleWithPermissionsDTO struct {
	ID            uuid.UUID     `json:"id"`
	Name          string        `json:"name"`
	DisplayName   string        `json:"display_name"`
	Description   string        `json:"description"`
	Permissions   []string      `json:"permissions"`
	PermissionIDs []uuid.UUID   `json:"permission_ids"`
	IsSystem      bool          `json:"is_system"`
	IsDefault     bool          `json:"is_default"`
	Priority      int           `json:"priority"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

func RoleWithPermissionsFromDomain(domainRole *role.Role, permissionCodes []string) *RoleWithPermissionsDTO {
	if domainRole == nil {
		return nil
	}
	return &RoleWithPermissionsDTO{
		ID:            domainRole.ID(),
		Name:          domainRole.Name(),
		DisplayName:   domainRole.DisplayName(),
		Description:   domainRole.Description(),
		Permissions:   permissionCodes,
		PermissionIDs: domainRole.PermissionIDs(),
		IsSystem:      domainRole.IsSystem(),
		IsDefault:     domainRole.IsDefault(),
		Priority:      domainRole.Priority(),
		CreatedAt:     domainRole.CreatedAt(),
		UpdatedAt:     domainRole.UpdatedAt(),
	}
}
