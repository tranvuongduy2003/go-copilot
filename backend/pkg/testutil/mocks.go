package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/tranvuongduy2003/go-copilot/internal/domain/auth"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/permission"
	"github.com/tranvuongduy2003/go-copilot/internal/domain/role"
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

func (m *MockUserRepository) FindByRole(ctx context.Context, roleID uuid.UUID) ([]*user.User, error) {
	if m.FindError != nil {
		return nil, m.FindError
	}
	result := make([]*user.User, 0)
	for _, u := range m.Users {
		for _, rid := range u.RoleIDs() {
			if rid == roleID {
				result = append(result, u)
				break
			}
		}
	}
	return result, nil
}

func (m *MockUserRepository) AddUser(u *user.User) {
	m.Users[u.ID()] = u
	m.EmailIndex[u.Email().String()] = u
}

type MockPasswordHasher struct {
	HashResult   string
	HashError    error
	VerifyResult bool
	VerifyError  error
}

func NewMockPasswordHasher() *MockPasswordHasher {
	return &MockPasswordHasher{
		HashResult:   "$2a$10$hashedpassword",
		VerifyResult: true,
	}
}

func (m *MockPasswordHasher) Hash(password string) (string, error) {
	if m.HashError != nil {
		return "", m.HashError
	}
	return m.HashResult, nil
}

func (m *MockPasswordHasher) Compare(hashedPassword, password string) error {
	if m.VerifyError != nil {
		return m.VerifyError
	}
	if !m.VerifyResult {
		return fmt.Errorf("password mismatch")
	}
	return nil
}

func (m *MockPasswordHasher) Verify(hashedPassword, plainPassword string) (bool, error) {
	if m.VerifyError != nil {
		return false, m.VerifyError
	}
	return m.VerifyResult, nil
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

type MockRoleRepository struct {
	Roles       map[uuid.UUID]*role.Role
	NameIndex   map[string]*role.Role
	CreateError error
	UpdateError error
	DeleteError error
	FindError   error
}

func NewMockRoleRepository() *MockRoleRepository {
	return &MockRoleRepository{
		Roles:     make(map[uuid.UUID]*role.Role),
		NameIndex: make(map[string]*role.Role),
	}
}

func (m *MockRoleRepository) Create(ctx context.Context, r *role.Role) error {
	if m.CreateError != nil {
		return m.CreateError
	}
	m.Roles[r.ID()] = r
	m.NameIndex[r.Name()] = r
	return nil
}

func (m *MockRoleRepository) Update(ctx context.Context, r *role.Role) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	m.Roles[r.ID()] = r
	return nil
}

func (m *MockRoleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	r, exists := m.Roles[id]
	if !exists {
		return role.ErrRoleNotFound
	}
	delete(m.Roles, id)
	delete(m.NameIndex, r.Name())
	return nil
}

func (m *MockRoleRepository) FindByID(ctx context.Context, id uuid.UUID) (*role.Role, error) {
	if m.FindError != nil {
		return nil, m.FindError
	}
	r, exists := m.Roles[id]
	if !exists {
		return nil, role.ErrRoleNotFound
	}
	return r, nil
}

func (m *MockRoleRepository) FindByName(ctx context.Context, name string) (*role.Role, error) {
	if m.FindError != nil {
		return nil, m.FindError
	}
	r, exists := m.NameIndex[name]
	if !exists {
		return nil, role.ErrRoleNotFound
	}
	return r, nil
}

func (m *MockRoleRepository) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]*role.Role, error) {
	if m.FindError != nil {
		return nil, m.FindError
	}
	result := make([]*role.Role, 0)
	for _, id := range ids {
		if r, exists := m.Roles[id]; exists {
			result = append(result, r)
		}
	}
	return result, nil
}

func (m *MockRoleRepository) FindDefault(ctx context.Context) (*role.Role, error) {
	if m.FindError != nil {
		return nil, m.FindError
	}
	for _, r := range m.Roles {
		if r.IsDefault() {
			return r, nil
		}
	}
	return nil, role.ErrRoleNotFound
}

func (m *MockRoleRepository) FindAll(ctx context.Context) ([]*role.Role, error) {
	if m.FindError != nil {
		return nil, m.FindError
	}
	result := make([]*role.Role, 0)
	for _, r := range m.Roles {
		result = append(result, r)
	}
	return result, nil
}

func (m *MockRoleRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	if m.FindError != nil {
		return false, m.FindError
	}
	_, exists := m.NameIndex[name]
	return exists, nil
}

func (m *MockRoleRepository) FindByPermission(ctx context.Context, permissionID uuid.UUID) ([]*role.Role, error) {
	if m.FindError != nil {
		return nil, m.FindError
	}
	result := make([]*role.Role, 0)
	for _, r := range m.Roles {
		if r.HasPermission(permissionID) {
			result = append(result, r)
		}
	}
	return result, nil
}

func (m *MockRoleRepository) AddRole(r *role.Role) {
	m.Roles[r.ID()] = r
	m.NameIndex[r.Name()] = r
}

type MockPermissionRepository struct {
	Permissions map[uuid.UUID]*permission.Permission
	CodeIndex   map[string]*permission.Permission
	FindError   error
}

func NewMockPermissionRepository() *MockPermissionRepository {
	return &MockPermissionRepository{
		Permissions: make(map[uuid.UUID]*permission.Permission),
		CodeIndex:   make(map[string]*permission.Permission),
	}
}

func (m *MockPermissionRepository) Create(ctx context.Context, p *permission.Permission) error {
	m.Permissions[p.ID()] = p
	m.CodeIndex[p.CodeString()] = p
	return nil
}

func (m *MockPermissionRepository) Update(ctx context.Context, p *permission.Permission) error {
	m.Permissions[p.ID()] = p
	return nil
}

func (m *MockPermissionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	p, exists := m.Permissions[id]
	if !exists {
		return permission.ErrPermissionNotFound
	}
	delete(m.Permissions, id)
	delete(m.CodeIndex, p.CodeString())
	return nil
}

func (m *MockPermissionRepository) FindByID(ctx context.Context, id uuid.UUID) (*permission.Permission, error) {
	if m.FindError != nil {
		return nil, m.FindError
	}
	p, exists := m.Permissions[id]
	if !exists {
		return nil, permission.ErrPermissionNotFound
	}
	return p, nil
}

func (m *MockPermissionRepository) FindByCode(ctx context.Context, code permission.PermissionCode) (*permission.Permission, error) {
	if m.FindError != nil {
		return nil, m.FindError
	}
	p, exists := m.CodeIndex[code.String()]
	if !exists {
		return nil, permission.ErrPermissionNotFound
	}
	return p, nil
}

func (m *MockPermissionRepository) FindByCodeString(ctx context.Context, code string) (*permission.Permission, error) {
	if m.FindError != nil {
		return nil, m.FindError
	}
	p, exists := m.CodeIndex[code]
	if !exists {
		return nil, permission.ErrPermissionNotFound
	}
	return p, nil
}

func (m *MockPermissionRepository) FindByResource(ctx context.Context, resource permission.Resource) ([]*permission.Permission, error) {
	if m.FindError != nil {
		return nil, m.FindError
	}
	result := make([]*permission.Permission, 0)
	for _, p := range m.Permissions {
		if p.Resource() == resource {
			result = append(result, p)
		}
	}
	return result, nil
}

func (m *MockPermissionRepository) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]*permission.Permission, error) {
	if m.FindError != nil {
		return nil, m.FindError
	}
	result := make([]*permission.Permission, 0)
	for _, id := range ids {
		if p, exists := m.Permissions[id]; exists {
			result = append(result, p)
		}
	}
	return result, nil
}

func (m *MockPermissionRepository) FindAll(ctx context.Context) ([]*permission.Permission, error) {
	if m.FindError != nil {
		return nil, m.FindError
	}
	result := make([]*permission.Permission, 0)
	for _, p := range m.Permissions {
		result = append(result, p)
	}
	return result, nil
}

func (m *MockPermissionRepository) ExistsByCode(ctx context.Context, code permission.PermissionCode) (bool, error) {
	if m.FindError != nil {
		return false, m.FindError
	}
	_, exists := m.CodeIndex[code.String()]
	return exists, nil
}

func (m *MockPermissionRepository) AddPermission(p *permission.Permission) {
	m.Permissions[p.ID()] = p
	m.CodeIndex[p.CodeString()] = p
}

type MockRefreshTokenRepository struct {
	Tokens        map[uuid.UUID]*auth.RefreshToken
	UserTokens    map[uuid.UUID][]*auth.RefreshToken
	HashIndex     map[string]*auth.RefreshToken
	CreateError   error
	UpdateError   error
	DeleteError   error
	FindError     error
}

func NewMockRefreshTokenRepository() *MockRefreshTokenRepository {
	return &MockRefreshTokenRepository{
		Tokens:     make(map[uuid.UUID]*auth.RefreshToken),
		UserTokens: make(map[uuid.UUID][]*auth.RefreshToken),
		HashIndex:  make(map[string]*auth.RefreshToken),
	}
}

func (m *MockRefreshTokenRepository) Create(ctx context.Context, token *auth.RefreshToken) error {
	if m.CreateError != nil {
		return m.CreateError
	}
	m.Tokens[token.ID()] = token
	m.HashIndex[token.TokenHash()] = token
	m.UserTokens[token.UserID()] = append(m.UserTokens[token.UserID()], token)
	return nil
}

func (m *MockRefreshTokenRepository) Update(ctx context.Context, token *auth.RefreshToken) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	m.Tokens[token.ID()] = token
	return nil
}

func (m *MockRefreshTokenRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	token, exists := m.Tokens[id]
	if !exists {
		return auth.ErrRefreshTokenNotFound
	}
	token.Revoke()
	return nil
}

func (m *MockRefreshTokenRepository) FindByID(ctx context.Context, id uuid.UUID) (*auth.RefreshToken, error) {
	if m.FindError != nil {
		return nil, m.FindError
	}
	token, exists := m.Tokens[id]
	if !exists {
		return nil, auth.ErrSessionNotFound
	}
	return token, nil
}

func (m *MockRefreshTokenRepository) FindByTokenHash(ctx context.Context, hash string) (*auth.RefreshToken, error) {
	if m.FindError != nil {
		return nil, m.FindError
	}
	token, exists := m.HashIndex[hash]
	if !exists {
		return nil, auth.ErrRefreshTokenNotFound
	}
	return token, nil
}

func (m *MockRefreshTokenRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*auth.RefreshToken, error) {
	if m.FindError != nil {
		return nil, m.FindError
	}
	return m.UserTokens[userID], nil
}

func (m *MockRefreshTokenRepository) FindActiveByUserID(ctx context.Context, userID uuid.UUID) ([]*auth.RefreshToken, error) {
	if m.FindError != nil {
		return nil, m.FindError
	}
	result := make([]*auth.RefreshToken, 0)
	for _, token := range m.UserTokens[userID] {
		if token.IsValid() {
			result = append(result, token)
		}
	}
	return result, nil
}

func (m *MockRefreshTokenRepository) RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	for _, token := range m.UserTokens[userID] {
		token.Revoke()
	}
	return nil
}

func (m *MockRefreshTokenRepository) DeleteExpired(ctx context.Context) (int64, error) {
	return 0, nil
}

func (m *MockRefreshTokenRepository) CountActiveByUserID(ctx context.Context, userID uuid.UUID) (int, error) {
	if m.FindError != nil {
		return 0, m.FindError
	}
	count := 0
	for _, token := range m.UserTokens[userID] {
		if token.IsValid() {
			count++
		}
	}
	return count, nil
}

type MockTokenGenerator struct {
	RefreshToken     string
	RefreshTokenHash string
	GenerateError    error
	ParseError       error
	ParsedClaims     *auth.Claims
}

func NewMockTokenGenerator() *MockTokenGenerator {
	return &MockTokenGenerator{
		RefreshToken:     "mock_refresh_token",
		RefreshTokenHash: "mock_hash",
	}
}

func (m *MockTokenGenerator) GenerateAccessToken(userID uuid.UUID, email string, roles, permissions []string) (auth.AccessToken, error) {
	if m.GenerateError != nil {
		return auth.AccessToken{}, m.GenerateError
	}
	return auth.NewAccessToken("mock_access_token", time.Now().Add(15*time.Minute)), nil
}

func (m *MockTokenGenerator) GenerateRefreshToken() (string, error) {
	if m.GenerateError != nil {
		return "", m.GenerateError
	}
	return m.RefreshToken, nil
}

func (m *MockTokenGenerator) HashRefreshToken(token string) string {
	return m.RefreshTokenHash
}

func (m *MockTokenGenerator) ParseAccessToken(token string) (*auth.Claims, error) {
	if m.ParseError != nil {
		return nil, m.ParseError
	}
	return m.ParsedClaims, nil
}

type MockTokenBlacklist struct {
	BlacklistedTokens map[string]bool
	AddError          error
	IsBlacklistedErr  error
}

func NewMockTokenBlacklist() *MockTokenBlacklist {
	return &MockTokenBlacklist{
		BlacklistedTokens: make(map[string]bool),
	}
}

func (m *MockTokenBlacklist) Add(ctx context.Context, tokenID string, expiresAt int64) error {
	if m.AddError != nil {
		return m.AddError
	}
	m.BlacklistedTokens[tokenID] = true
	return nil
}

func (m *MockTokenBlacklist) IsBlacklisted(ctx context.Context, tokenID string) (bool, error) {
	if m.IsBlacklistedErr != nil {
		return false, m.IsBlacklistedErr
	}
	return m.BlacklistedTokens[tokenID], nil
}

type MockAccountLockout struct {
	Locked          bool
	RemainingTime   time.Duration
	AttemptCount    int
	IsLockedErr     error
	RecordErr       error
}

func NewMockAccountLockout() *MockAccountLockout {
	return &MockAccountLockout{}
}

func (m *MockAccountLockout) IsLocked(ctx context.Context, identifier string) (bool, time.Duration, error) {
	if m.IsLockedErr != nil {
		return false, 0, m.IsLockedErr
	}
	return m.Locked, m.RemainingTime, nil
}

func (m *MockAccountLockout) RecordFailedAttempt(ctx context.Context, identifier string) (int, error) {
	if m.RecordErr != nil {
		return 0, m.RecordErr
	}
	m.AttemptCount++
	return m.AttemptCount, nil
}

func (m *MockAccountLockout) ResetAttempts(ctx context.Context, identifier string) error {
	m.AttemptCount = 0
	return nil
}

func (m *MockAccountLockout) GetAttemptCount(ctx context.Context, identifier string) (int, error) {
	return m.AttemptCount, nil
}

type MockPasswordResetTokenStore struct {
	Tokens      map[string]string
	StoreError  error
	GetError    error
	DeleteError error
}

func NewMockPasswordResetTokenStore() *MockPasswordResetTokenStore {
	return &MockPasswordResetTokenStore{
		Tokens: make(map[string]string),
	}
}

func (m *MockPasswordResetTokenStore) Store(ctx context.Context, email string, tokenHash string, expiresAt time.Time) error {
	if m.StoreError != nil {
		return m.StoreError
	}
	m.Tokens[tokenHash] = email
	return nil
}

func (m *MockPasswordResetTokenStore) Get(ctx context.Context, tokenHash string) (string, error) {
	if m.GetError != nil {
		return "", m.GetError
	}
	email, exists := m.Tokens[tokenHash]
	if !exists {
		return "", nil
	}
	return email, nil
}

func (m *MockPasswordResetTokenStore) Delete(ctx context.Context, tokenHash string) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	delete(m.Tokens, tokenHash)
	return nil
}
