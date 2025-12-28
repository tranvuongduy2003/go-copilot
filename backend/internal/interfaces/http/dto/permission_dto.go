package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreatePermissionRequest struct {
	Resource    string `json:"resource" validate:"required,min=1,max=100"`
	Action      string `json:"action" validate:"required,min=1,max=100"`
	Description string `json:"description" validate:"omitempty,max=500"`
}

type UpdatePermissionRequest struct {
	Description string `json:"description" validate:"omitempty,max=500"`
}

type PermissionResponse struct {
	ID          uuid.UUID `json:"id"`
	Resource    string    `json:"resource"`
	Action      string    `json:"action"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	IsSystem    bool      `json:"is_system"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
