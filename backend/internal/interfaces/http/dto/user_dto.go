package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/user"
)

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
	FullName string `json:"full_name" validate:"required,min=2,max=255"`
}

type UpdateUserRequest struct {
	FullName *string `json:"full_name,omitempty" validate:"omitempty,min=2,max=255"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,password"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

type BanUserRequest struct {
	Reason string `json:"reason" validate:"required,min=1,max=500"`
}

type ListUsersRequest struct {
	Page      int     `json:"page" validate:"omitempty,gte=1"`
	Limit     int     `json:"limit" validate:"omitempty,gte=1,lte=100"`
	Status    *string `json:"status,omitempty" validate:"omitempty,oneof=pending active inactive banned"`
	Search    *string `json:"search,omitempty" validate:"omitempty,max=255"`
	SortBy    *string `json:"sort_by,omitempty" validate:"omitempty,oneof=created_at updated_at email full_name"`
	SortOrder *string `json:"sort_order,omitempty" validate:"omitempty,oneof=asc desc"`
	DateFrom  *string `json:"date_from,omitempty"`
	DateTo    *string `json:"date_to,omitempty"`
}

type UserResponse struct {
	ID        uuid.UUID  `json:"id"`
	Email     string     `json:"email"`
	FullName  string     `json:"full_name"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func UserResponseFromDomain(domainUser *user.User) UserResponse {
	return UserResponse{
		ID:        domainUser.ID(),
		Email:     domainUser.Email().String(),
		FullName:  domainUser.FullName().String(),
		Status:    domainUser.Status().String(),
		CreatedAt: domainUser.CreatedAt(),
		UpdatedAt: domainUser.UpdatedAt(),
		DeletedAt: domainUser.DeletedAt(),
	}
}

func UserResponsesFromDomain(domainUsers []*user.User) []UserResponse {
	responses := make([]UserResponse, len(domainUsers))
	for i, domainUser := range domainUsers {
		responses[i] = UserResponseFromDomain(domainUser)
	}
	return responses
}

type PaginatedUsersResponse struct {
	Items      []UserResponse `json:"items"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
	HasNext    bool           `json:"has_next"`
	HasPrev    bool           `json:"has_prev"`
}
