package permissiondto

import (
	"time"

	"github.com/google/uuid"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/permission"
)

type PermissionDTO struct {
	ID          uuid.UUID `json:"id"`
	Resource    string    `json:"resource"`
	Action      string    `json:"action"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	IsSystem    bool      `json:"is_system"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func PermissionFromDomain(domainPermission *permission.Permission) *PermissionDTO {
	if domainPermission == nil {
		return nil
	}
	return &PermissionDTO{
		ID:          domainPermission.ID(),
		Resource:    domainPermission.Resource().String(),
		Action:      domainPermission.Action().String(),
		Code:        domainPermission.CodeString(),
		Description: domainPermission.Description(),
		IsSystem:    domainPermission.IsSystem(),
		CreatedAt:   domainPermission.CreatedAt(),
		UpdatedAt:   domainPermission.UpdatedAt(),
	}
}

func PermissionsFromDomain(domainPermissions []*permission.Permission) []*PermissionDTO {
	dtos := make([]*PermissionDTO, len(domainPermissions))
	for i, domainPermission := range domainPermissions {
		dtos[i] = PermissionFromDomain(domainPermission)
	}
	return dtos
}
