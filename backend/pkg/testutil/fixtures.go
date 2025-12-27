package testutil

import (
	"time"

	"github.com/google/uuid"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/user"
)

type UserBuilder struct {
	id           uuid.UUID
	email        string
	passwordHash string
	fullName     string
	status       user.Status
	createdAt    time.Time
	updatedAt    time.Time
	deletedAt    *time.Time
}

func NewUserBuilder() *UserBuilder {
	now := time.Now().UTC()
	return &UserBuilder{
		id:           uuid.New(),
		email:        RandomEmail(),
		passwordHash: "$2a$10$hashedpassword",
		fullName:     RandomFullName(),
		status:       user.StatusPending,
		createdAt:    now,
		updatedAt:    now,
		deletedAt:    nil,
	}
}

func (b *UserBuilder) WithID(id uuid.UUID) *UserBuilder {
	b.id = id
	return b
}

func (b *UserBuilder) WithEmail(email string) *UserBuilder {
	b.email = email
	return b
}

func (b *UserBuilder) WithPasswordHash(hash string) *UserBuilder {
	b.passwordHash = hash
	return b
}

func (b *UserBuilder) WithFullName(name string) *UserBuilder {
	b.fullName = name
	return b
}

func (b *UserBuilder) WithStatus(status user.Status) *UserBuilder {
	b.status = status
	return b
}

func (b *UserBuilder) WithCreatedAt(t time.Time) *UserBuilder {
	b.createdAt = t
	return b
}

func (b *UserBuilder) WithUpdatedAt(t time.Time) *UserBuilder {
	b.updatedAt = t
	return b
}

func (b *UserBuilder) WithDeletedAt(t *time.Time) *UserBuilder {
	b.deletedAt = t
	return b
}

func (b *UserBuilder) Active() *UserBuilder {
	b.status = user.StatusActive
	return b
}

func (b *UserBuilder) Inactive() *UserBuilder {
	b.status = user.StatusInactive
	return b
}

func (b *UserBuilder) Banned() *UserBuilder {
	b.status = user.StatusBanned
	return b
}

func (b *UserBuilder) Pending() *UserBuilder {
	b.status = user.StatusPending
	return b
}

func (b *UserBuilder) Deleted() *UserBuilder {
	now := time.Now().UTC()
	b.deletedAt = &now
	return b
}

func (b *UserBuilder) Build() (*user.User, error) {
	return user.ReconstructUser(user.ReconstructUserParams{
		ID:           b.id,
		Email:        b.email,
		PasswordHash: b.passwordHash,
		FullName:     b.fullName,
		Status:       b.status,
		CreatedAt:    b.createdAt,
		UpdatedAt:    b.updatedAt,
		DeletedAt:    b.deletedAt,
	})
}

func (b *UserBuilder) MustBuild() *user.User {
	u, err := b.Build()
	if err != nil {
		panic(err)
	}
	return u
}

func CreateTestUser() *user.User {
	return NewUserBuilder().MustBuild()
}

func CreateActiveUser() *user.User {
	return NewUserBuilder().Active().MustBuild()
}

func CreatePendingUser() *user.User {
	return NewUserBuilder().Pending().MustBuild()
}

func CreateInactiveUser() *user.User {
	return NewUserBuilder().Inactive().MustBuild()
}

func CreateBannedUser() *user.User {
	return NewUserBuilder().Banned().MustBuild()
}

func CreateDeletedUser() *user.User {
	return NewUserBuilder().Deleted().MustBuild()
}

func CreateUsers(count int) []*user.User {
	users := make([]*user.User, count)
	for i := 0; i < count; i++ {
		users[i] = CreateTestUser()
	}
	return users
}
