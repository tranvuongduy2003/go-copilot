package user

import (
	"time"

	"github.com/google/uuid"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/shared"
)

type User struct {
	shared.AggregateRoot
	email        shared.Email
	passwordHash shared.PasswordHash
	fullName     shared.FullName
	status       Status
	createdAt    time.Time
	updatedAt    time.Time
	deletedAt    *time.Time
}

type NewUserParams struct {
	Email        string
	PasswordHash string
	FullName     string
}

func NewUser(params NewUserParams) (*User, error) {
	email, err := shared.NewEmail(params.Email)
	if err != nil {
		return nil, err
	}

	passwordHash, err := shared.NewPasswordHash(params.PasswordHash)
	if err != nil {
		return nil, err
	}

	fullName, err := shared.NewFullName(params.FullName)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	user := &User{
		AggregateRoot: shared.NewAggregateRoot(),
		email:         email,
		passwordHash:  passwordHash,
		fullName:      fullName,
		status:        StatusPending,
		createdAt:     now,
		updatedAt:     now,
		deletedAt:     nil,
	}

	user.AddDomainEvent(NewUserCreatedEvent(user.ID(), email.String(), fullName.String()))

	return user, nil
}

type ReconstructUserParams struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	FullName     string
	Status       Status
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

func ReconstructUser(params ReconstructUserParams) (*User, error) {
	email, err := shared.NewEmail(params.Email)
	if err != nil {
		return nil, err
	}

	passwordHash, err := shared.NewPasswordHash(params.PasswordHash)
	if err != nil {
		return nil, err
	}

	fullName, err := shared.NewFullName(params.FullName)
	if err != nil {
		return nil, err
	}

	if !params.Status.IsValid() {
		return nil, ErrInvalidStatus
	}

	return &User{
		AggregateRoot: shared.NewAggregateRootWithID(params.ID),
		email:         email,
		passwordHash:  passwordHash,
		fullName:      fullName,
		status:        params.Status,
		createdAt:     params.CreatedAt,
		updatedAt:     params.UpdatedAt,
		deletedAt:     params.DeletedAt,
	}, nil
}

func (u *User) Email() shared.Email {
	return u.email
}

func (u *User) PasswordHash() shared.PasswordHash {
	return u.passwordHash
}

func (u *User) FullName() shared.FullName {
	return u.fullName
}

func (u *User) Status() Status {
	return u.status
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

func (u *User) DeletedAt() *time.Time {
	return u.deletedAt
}

func (u *User) IsDeleted() bool {
	return u.deletedAt != nil
}

func (u *User) Activate() error {
	if u.status.IsActive() {
		return ErrUserAlreadyActive
	}

	if !u.status.CanTransitionTo(StatusActive) {
		return NewInvalidStatusTransitionError(u.status, StatusActive)
	}

	u.status = StatusActive
	u.updatedAt = time.Now().UTC()
	u.AddDomainEvent(NewUserActivatedEvent(u.ID()))

	return nil
}

func (u *User) Deactivate() error {
	if u.status.IsInactive() {
		return ErrUserAlreadyInactive
	}

	if !u.status.CanTransitionTo(StatusInactive) {
		return NewInvalidStatusTransitionError(u.status, StatusInactive)
	}

	u.status = StatusInactive
	u.updatedAt = time.Now().UTC()
	u.AddDomainEvent(NewUserDeactivatedEvent(u.ID()))

	return nil
}

func (u *User) Ban(reason string) error {
	if u.status.IsBanned() {
		return ErrUserIsBanned
	}

	if !u.status.CanTransitionTo(StatusBanned) {
		return NewInvalidStatusTransitionError(u.status, StatusBanned)
	}

	u.status = StatusBanned
	u.updatedAt = time.Now().UTC()
	u.AddDomainEvent(NewUserBannedEvent(u.ID(), reason))

	return nil
}

func (u *User) ChangePassword(newPasswordHash string) error {
	passwordHash, err := shared.NewPasswordHash(newPasswordHash)
	if err != nil {
		return err
	}

	u.passwordHash = passwordHash
	u.updatedAt = time.Now().UTC()
	u.AddDomainEvent(NewPasswordChangedEvent(u.ID()))

	return nil
}

func (u *User) UpdateProfile(fullName string) error {
	var changedFields []string

	if fullName != "" && fullName != u.fullName.String() {
		newFullName, err := shared.NewFullName(fullName)
		if err != nil {
			return err
		}
		u.fullName = newFullName
		changedFields = append(changedFields, "full_name")
	}

	if len(changedFields) > 0 {
		u.updatedAt = time.Now().UTC()
		u.AddDomainEvent(NewProfileUpdatedEvent(u.ID(), changedFields))
	}

	return nil
}

func (u *User) Delete() error {
	if u.IsDeleted() {
		return shared.NewBusinessRuleViolationError("user_already_deleted", "user is already deleted")
	}

	now := time.Now().UTC()
	u.deletedAt = &now
	u.updatedAt = now
	u.AddDomainEvent(NewUserDeletedEvent(u.ID(), now))

	return nil
}
