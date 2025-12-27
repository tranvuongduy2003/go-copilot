package testutil

import (
	"context"

	"github.com/google/uuid"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/shared"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/user"
	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

type MockUserRepository struct {
	Users       map[uuid.UUID]*user.User
	EmailIndex  map[string]*user.User
	CreateError error
	UpdateError error
	DeleteError error
	FindError   error
	ListError   error
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		Users:      make(map[uuid.UUID]*user.User),
		EmailIndex: make(map[string]*user.User),
	}
}

func (m *MockUserRepository) Create(ctx context.Context, u *user.User) error {
	if m.CreateError != nil {
		return m.CreateError
	}
	if _, exists := m.EmailIndex[u.Email().String()]; exists {
		return user.NewEmailAlreadyExistsError(u.Email().String())
	}
	m.Users[u.ID()] = u
	m.EmailIndex[u.Email().String()] = u
	return nil
}

func (m *MockUserRepository) Update(ctx context.Context, u *user.User) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	if _, exists := m.Users[u.ID()]; !exists {
		return user.NewUserNotFoundError(u.ID().String())
	}
	m.Users[u.ID()] = u
	m.EmailIndex[u.Email().String()] = u
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	u, exists := m.Users[id]
	if !exists {
		return user.NewUserNotFoundError(id.String())
	}
	delete(m.Users, id)
	delete(m.EmailIndex, u.Email().String())
	return nil
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	if m.FindError != nil {
		return nil, m.FindError
	}
	u, exists := m.Users[id]
	if !exists {
		return nil, user.NewUserNotFoundError(id.String())
	}
	return u, nil
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	if m.FindError != nil {
		return nil, m.FindError
	}
	u, exists := m.EmailIndex[email]
	if !exists {
		return nil, user.NewUserNotFoundError(email)
	}
	return u, nil
}

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	if m.FindError != nil {
		return false, m.FindError
	}
	_, exists := m.EmailIndex[email]
	return exists, nil
}

func (m *MockUserRepository) List(ctx context.Context, filter user.Filter, pagination shared.Pagination) ([]*user.User, int64, error) {
	if m.ListError != nil {
		return nil, 0, m.ListError
	}

	result := make([]*user.User, 0)
	for _, u := range m.Users {
		if filter.Status != nil && u.Status() != *filter.Status {
			continue
		}
		result = append(result, u)
	}

	total := int64(len(result))

	offset := pagination.Offset()
	limit := pagination.Limit()

	if offset >= int(total) {
		return []*user.User{}, total, nil
	}

	end := offset + limit
	if end > int(total) {
		end = int(total)
	}

	return result[offset:end], total, nil
}

func (m *MockUserRepository) AddUser(u *user.User) {
	m.Users[u.ID()] = u
	m.EmailIndex[u.Email().String()] = u
}

type MockPasswordHasher struct {
	HashResult    string
	HashError     error
	CompareError  error
}

func NewMockPasswordHasher() *MockPasswordHasher {
	return &MockPasswordHasher{
		HashResult: "$2a$10$hashedpassword",
	}
}

func (m *MockPasswordHasher) Hash(password string) (string, error) {
	if m.HashError != nil {
		return "", m.HashError
	}
	return m.HashResult, nil
}

func (m *MockPasswordHasher) Compare(hashedPassword, password string) error {
	return m.CompareError
}

type MockEventBus struct {
	PublishedEvents []shared.DomainEvent
	PublishError    error
	Handlers        map[string][]shared.EventHandler
}

func NewMockEventBus() *MockEventBus {
	return &MockEventBus{
		PublishedEvents: make([]shared.DomainEvent, 0),
		Handlers:        make(map[string][]shared.EventHandler),
	}
}

func (m *MockEventBus) Publish(ctx context.Context, events ...shared.DomainEvent) error {
	if m.PublishError != nil {
		return m.PublishError
	}
	m.PublishedEvents = append(m.PublishedEvents, events...)
	return nil
}

func (m *MockEventBus) Subscribe(eventType string, handler shared.EventHandler) {
	m.Handlers[eventType] = append(m.Handlers[eventType], handler)
}

func (m *MockEventBus) Unsubscribe(eventType string, handler shared.EventHandler) {
}

type NoopLogger struct{}

func NewNoopLogger() *NoopLogger {
	return &NoopLogger{}
}

func (l *NoopLogger) Debug(msg string, fields ...logger.Field) {}
func (l *NoopLogger) Info(msg string, fields ...logger.Field)  {}
func (l *NoopLogger) Warn(msg string, fields ...logger.Field)  {}
func (l *NoopLogger) Error(msg string, fields ...logger.Field) {}
func (l *NoopLogger) Fatal(msg string, fields ...logger.Field) {}
func (l *NoopLogger) With(fields ...logger.Field) logger.Logger {
	return l
}
func (l *NoopLogger) Sync() error {
	return nil
}
