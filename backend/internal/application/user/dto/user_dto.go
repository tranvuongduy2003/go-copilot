package userdto

import (
	"time"

	"github.com/google/uuid"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/shared"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/user"
)

type UserDTO struct {
	ID        uuid.UUID  `json:"id"`
	Email     string     `json:"email"`
	FullName  string     `json:"full_name"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func UserFromDomain(domainUser *user.User) *UserDTO {
	if domainUser == nil {
		return nil
	}
	return &UserDTO{
		ID:        domainUser.ID(),
		Email:     domainUser.Email().String(),
		FullName:  domainUser.FullName().String(),
		Status:    domainUser.Status().String(),
		CreatedAt: domainUser.CreatedAt(),
		UpdatedAt: domainUser.UpdatedAt(),
		DeletedAt: domainUser.DeletedAt(),
	}
}

func UsersFromDomain(domainUsers []*user.User) []*UserDTO {
	dtos := make([]*UserDTO, len(domainUsers))
	for i, domainUser := range domainUsers {
		dtos[i] = UserFromDomain(domainUser)
	}
	return dtos
}

type PaginatedUsersDTO struct {
	Items      []*UserDTO `json:"items"`
	Total      int64      `json:"total"`
	Page       int        `json:"page"`
	Limit      int        `json:"limit"`
	TotalPages int        `json:"total_pages"`
	HasNext    bool       `json:"has_next"`
	HasPrev    bool       `json:"has_prev"`
}

func NewPaginatedUsersDTO(users []*user.User, total int64, pagination shared.Pagination) *PaginatedUsersDTO {
	return &PaginatedUsersDTO{
		Items:      UsersFromDomain(users),
		Total:      total,
		Page:       pagination.Page(),
		Limit:      pagination.Limit(),
		TotalPages: pagination.TotalPages(total),
		HasNext:    pagination.HasNext(total),
		HasPrev:    pagination.HasPrev(),
	}
}
