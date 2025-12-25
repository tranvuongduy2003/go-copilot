---
name: testing-patterns
description: Write comprehensive tests for Go and React. Use when adding tests.
---

# Testing Patterns Skill

This skill provides testing patterns and templates for both Go backend and React frontend code.

## Go Testing Patterns

### Pattern 1: Table-Driven Unit Tests

```go
func TestCalculateDiscount(t *testing.T) {
    tests := []struct {
        name     string
        price    int
        quantity int
        want     int
        wantErr  error
    }{
        {
            name:     "no discount for small order",
            price:    100,
            quantity: 1,
            want:     100,
            wantErr:  nil,
        },
        {
            name:     "5% discount for 10+ items",
            price:    100,
            quantity: 10,
            want:     950, // 1000 - 5%
            wantErr:  nil,
        },
        {
            name:     "10% discount for 50+ items",
            price:    100,
            quantity: 50,
            want:     4500, // 5000 - 10%
            wantErr:  nil,
        },
        {
            name:     "error for zero price",
            price:    0,
            quantity: 1,
            want:     0,
            wantErr:  ErrInvalidPrice,
        },
        {
            name:     "error for negative quantity",
            price:    100,
            quantity: -1,
            want:     0,
            wantErr:  ErrInvalidQuantity,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := CalculateDiscount(tt.price, tt.quantity)

            if tt.wantErr != nil {
                require.Error(t, err)
                assert.ErrorIs(t, err, tt.wantErr)
                return
            }

            require.NoError(t, err)
            assert.Equal(t, tt.want, got)
        })
    }
}
```

### Pattern 2: Command Handler Tests (CQRS)

Test command handlers that orchestrate domain operations.

```go
// internal/application/command/create_user_test.go
package command_test

import (
    "context"
    "log/slog"
    "testing"

    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"
    "github.com/yourorg/app/internal/application/command"
    "github.com/yourorg/app/internal/domain/user"
)

// MockUserRepository implements user.Repository for testing
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
    args := m.Called(ctx, email)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepository) Save(ctx context.Context, u *user.User) error {
    args := m.Called(ctx, u)
    return args.Error(0)
}

func TestCreateUserHandler_Handle(t *testing.T) {
    tests := []struct {
        name    string
        cmd     command.CreateUserCommand
        setup   func(*MockUserRepository, *MockPasswordHasher)
        wantErr error
    }{
        {
            name: "success - creates user",
            cmd: command.CreateUserCommand{
                Email:    "test@example.com",
                Name:     "Test User",
                Password: "password123",
            },
            setup: func(repo *MockUserRepository, hasher *MockPasswordHasher) {
                repo.On("FindByEmail", mock.Anything, "test@example.com").
                    Return(nil, user.ErrNotFound)
                hasher.On("Hash", "password123").
                    Return("hashed_password", nil)
                repo.On("Save", mock.Anything, mock.AnythingOfType("*user.User")).
                    Return(nil)
            },
            wantErr: nil,
        },
        {
            name: "error - email already exists",
            cmd: command.CreateUserCommand{
                Email:    "exists@example.com",
                Name:     "Test User",
                Password: "password123",
            },
            setup: func(repo *MockUserRepository, hasher *MockPasswordHasher) {
                existingUser := user.Reconstitute(uuid.New(), "exists@example.com", "Existing", "", time.Now(), time.Now())
                repo.On("FindByEmail", mock.Anything, "exists@example.com").
                    Return(existingUser, nil)
            },
            wantErr: user.ErrEmailExists,
        },
        {
            name: "error - invalid email format",
            cmd: command.CreateUserCommand{
                Email:    "invalid-email",
                Name:     "Test User",
                Password: "password123",
            },
            setup:   func(repo *MockUserRepository, hasher *MockPasswordHasher) {},
            wantErr: user.ErrInvalidEmail,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            repo := new(MockUserRepository)
            hasher := new(MockPasswordHasher)
            tt.setup(repo, hasher)

            handler := command.NewCreateUserHandler(repo, hasher, slog.Default())

            // Act
            result, err := handler.Handle(context.Background(), tt.cmd)

            // Assert
            if tt.wantErr != nil {
                require.Error(t, err)
                assert.ErrorIs(t, err, tt.wantErr)
                assert.Nil(t, result)
                return
            }

            require.NoError(t, err)
            assert.NotNil(t, result)
            assert.Equal(t, tt.cmd.Email, result.Email)
            assert.NotEqual(t, uuid.Nil, result.ID)

            // Verify mock expectations
            repo.AssertExpectations(t)
            hasher.AssertExpectations(t)
        })
    }
}
```

### Pattern 3: HTTP Handler Tests (with CQRS Handlers)

Test HTTP handlers that use command/query handlers.

```go
// internal/interfaces/http/handler/user_handler_test.go
package handler_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/go-chi/chi/v5"
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"
    "github.com/yourorg/app/internal/application/command"
    "github.com/yourorg/app/internal/application/dto"
    "github.com/yourorg/app/internal/application/query"
    "github.com/yourorg/app/internal/domain/user"
    "github.com/yourorg/app/internal/interfaces/http/handler"
)

// MockCreateUserHandler mocks the create user command handler
type MockCreateUserHandler struct {
    mock.Mock
}

func (m *MockCreateUserHandler) Handle(ctx context.Context, cmd command.CreateUserCommand) (*dto.UserDTO, error) {
    args := m.Called(ctx, cmd)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.UserDTO), args.Error(1)
}

// MockGetUserHandler mocks the get user query handler
type MockGetUserHandler struct {
    mock.Mock
}

func (m *MockGetUserHandler) Handle(ctx context.Context, q query.GetUserQuery) (*dto.UserDTO, error) {
    args := m.Called(ctx, q)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*dto.UserDTO), args.Error(1)
}

func TestUserHandler_Create(t *testing.T) {
    tests := []struct {
        name       string
        body       interface{}
        setup      func(*MockCreateUserHandler)
        wantStatus int
        wantBody   func(*testing.T, []byte)
    }{
        {
            name: "success - returns created user",
            body: map[string]string{
                "email":    "test@example.com",
                "name":     "Test User",
                "password": "password123",
            },
            setup: func(m *MockCreateUserHandler) {
                m.On("Handle", mock.Anything, mock.MatchedBy(func(cmd command.CreateUserCommand) bool {
                    return cmd.Email == "test@example.com" && cmd.Name == "Test User"
                })).Return(&dto.UserDTO{
                    ID:    uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
                    Email: "test@example.com",
                    Name:  "Test User",
                }, nil)
            },
            wantStatus: http.StatusCreated,
            wantBody: func(t *testing.T, body []byte) {
                var resp map[string]interface{}
                require.NoError(t, json.Unmarshal(body, &resp))

                data := resp["data"].(map[string]interface{})
                assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", data["id"])
                assert.Equal(t, "test@example.com", data["email"])
            },
        },
        {
            name: "validation error - invalid email format",
            body: map[string]string{
                "email":    "invalid-email",
                "name":     "Test User",
                "password": "password123",
            },
            setup:      func(m *MockCreateUserHandler) {},
            wantStatus: http.StatusBadRequest,
            wantBody: func(t *testing.T, body []byte) {
                var resp map[string]interface{}
                require.NoError(t, json.Unmarshal(body, &resp))
                assert.Contains(t, resp["error"].(map[string]interface{})["code"], "VALIDATION")
            },
        },
        {
            name: "conflict error - email already exists",
            body: map[string]string{
                "email":    "exists@example.com",
                "name":     "Test User",
                "password": "password123",
            },
            setup: func(m *MockCreateUserHandler) {
                m.On("Handle", mock.Anything, mock.Anything).
                    Return(nil, user.ErrEmailExists)
            },
            wantStatus: http.StatusConflict,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            createHandler := new(MockCreateUserHandler)
            tt.setup(createHandler)

            h := handler.NewUserHandler(createHandler, nil, nil, nil, nil)
            router := chi.NewRouter()
            h.RegisterRoutes(router)

            body, _ := json.Marshal(tt.body)
            req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
            req.Header.Set("Content-Type", "application/json")
            rec := httptest.NewRecorder()

            // Act
            router.ServeHTTP(rec, req)

            // Assert
            assert.Equal(t, tt.wantStatus, rec.Code)
            if tt.wantBody != nil {
                tt.wantBody(t, rec.Body.Bytes())
            }
            createHandler.AssertExpectations(t)
        })
    }
}

func TestUserHandler_Get(t *testing.T) {
    userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

    tests := []struct {
        name       string
        id         string
        setup      func(*MockGetUserHandler)
        wantStatus int
    }{
        {
            name: "success - returns user",
            id:   userID.String(),
            setup: func(m *MockGetUserHandler) {
                m.On("Handle", mock.Anything, query.GetUserQuery{ID: userID}).
                    Return(&dto.UserDTO{ID: userID, Email: "test@example.com"}, nil)
            },
            wantStatus: http.StatusOK,
        },
        {
            name: "not found",
            id:   uuid.New().String(),
            setup: func(m *MockGetUserHandler) {
                m.On("Handle", mock.Anything, mock.Anything).
                    Return(nil, user.ErrNotFound)
            },
            wantStatus: http.StatusNotFound,
        },
        {
            name:       "invalid uuid",
            id:         "invalid-uuid",
            setup:      func(m *MockGetUserHandler) {},
            wantStatus: http.StatusBadRequest,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            getHandler := new(MockGetUserHandler)
            tt.setup(getHandler)

            h := handler.NewUserHandler(nil, nil, nil, getHandler, nil)
            router := chi.NewRouter()
            h.RegisterRoutes(router)

            req := httptest.NewRequest(http.MethodGet, "/users/"+tt.id, nil)
            rec := httptest.NewRecorder()

            router.ServeHTTP(rec, req)

            assert.Equal(t, tt.wantStatus, rec.Code)
        })
    }
}
```

### Pattern 4: Integration Tests with Test Containers

Test repository adapters with a real database using DDD entities.

```go
//go:build integration

package integration_test

import (
    "context"
    "testing"
    "time"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/pressly/goose/v3"
    "github.com/stretchr/testify/suite"
    "github.com/testcontainers/testcontainers-go/modules/postgres"
    "github.com/yourorg/app/internal/domain/user"
    pgrepo "github.com/yourorg/app/internal/infrastructure/persistence/postgres"
)

type IntegrationSuite struct {
    suite.Suite
    container *postgres.PostgresContainer
    db        *pgxpool.Pool
}

func (s *IntegrationSuite) SetupSuite() {
    ctx := context.Background()

    container, err := postgres.Run(ctx,
        "postgres:16-alpine",
        postgres.WithDatabase("testdb"),
        postgres.WithUsername("test"),
        postgres.WithPassword("test"),
    )
    s.Require().NoError(err)
    s.container = container

    connStr, err := container.ConnectionString(ctx, "sslmode=disable")
    s.Require().NoError(err)

    s.db, err = pgxpool.New(ctx, connStr)
    s.Require().NoError(err)

    // Run Goose migrations
    s.runMigrations(connStr)
}

func (s *IntegrationSuite) runMigrations(connStr string) {
    db, err := goose.OpenDBWithDriver("postgres", connStr)
    s.Require().NoError(err)
    defer db.Close()

    err = goose.Up(db, "../../migrations/sql")
    s.Require().NoError(err)
}

func (s *IntegrationSuite) TearDownSuite() {
    s.db.Close()
    s.container.Terminate(context.Background())
}

func (s *IntegrationSuite) SetupTest() {
    // Clean tables between tests
    s.db.Exec(context.Background(), "TRUNCATE users, posts CASCADE")
}

func (s *IntegrationSuite) TestUserRepository_SaveAndFind() {
    ctx := context.Background()
    repo := pgrepo.NewUserRepository(s.db)

    // Create a domain entity using NewUser factory
    newUser, err := user.NewUser("test@example.com", "Test User", "hashed_password")
    s.Require().NoError(err)

    // Save via repository
    err = repo.Save(ctx, newUser)
    s.NoError(err)

    // Fetch and verify using getter methods
    found, err := repo.FindByID(ctx, newUser.ID())
    s.NoError(err)
    s.Equal(newUser.Email(), found.Email())
    s.Equal(newUser.Name(), found.Name())
}

func (s *IntegrationSuite) TestUserRepository_FindByEmail() {
    ctx := context.Background()
    repo := pgrepo.NewUserRepository(s.db)

    // Create user
    newUser, _ := user.NewUser("findme@example.com", "Find Me", "hashed_password")
    err := repo.Save(ctx, newUser)
    s.Require().NoError(err)

    // Find by email
    found, err := repo.FindByEmail(ctx, "findme@example.com")
    s.NoError(err)
    s.Equal(newUser.ID(), found.ID())
}

func (s *IntegrationSuite) TestUserRepository_NotFound() {
    ctx := context.Background()
    repo := pgrepo.NewUserRepository(s.db)

    // Try to find non-existent user
    _, err := repo.FindByID(ctx, uuid.New())
    s.ErrorIs(err, user.ErrNotFound)
}

func (s *IntegrationSuite) TestUserRepository_UpdateViaReconstitute() {
    ctx := context.Background()
    repo := pgrepo.NewUserRepository(s.db)

    // Create initial user
    newUser, _ := user.NewUser("update@example.com", "Original Name", "hashed_password")
    err := repo.Save(ctx, newUser)
    s.Require().NoError(err)

    // Fetch, modify via domain method, and save
    found, err := repo.FindByID(ctx, newUser.ID())
    s.Require().NoError(err)

    // Use domain method to change name (assuming UpdateName exists)
    err = found.UpdateName("Updated Name")
    s.Require().NoError(err)

    err = repo.Save(ctx, found)
    s.NoError(err)

    // Verify update
    updated, _ := repo.FindByID(ctx, newUser.ID())
    s.Equal("Updated Name", updated.Name())
}

func TestIntegrationSuite(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration tests")
    }
    suite.Run(t, new(IntegrationSuite))
}
```

## React Testing Patterns

### Pattern 1: Component Rendering Tests

```tsx
import { render, screen } from '@testing-library/react';
import { UserCard } from './user-card';

describe('UserCard', () => {
  const mockUser = {
    id: '1',
    name: 'John Doe',
    email: 'john@example.com',
    avatarUrl: '/avatar.jpg',
  };

  it('renders user information', () => {
    render(<UserCard user={mockUser} />);

    expect(screen.getByText('John Doe')).toBeInTheDocument();
    expect(screen.getByText('john@example.com')).toBeInTheDocument();
    expect(screen.getByRole('img')).toHaveAttribute('src', '/avatar.jpg');
  });

  it('renders fallback avatar when no URL provided', () => {
    render(<UserCard user={{ ...mockUser, avatarUrl: undefined }} />);

    expect(screen.getByText('JD')).toBeInTheDocument(); // Initials
  });
});
```

### Pattern 2: Interactive Component Tests

```tsx
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { vi } from 'vitest';
import { LoginForm } from './login-form';

describe('LoginForm', () => {
  const mockOnSubmit = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('submits with valid credentials', async () => {
    const user = userEvent.setup();
    mockOnSubmit.mockResolvedValue(undefined);

    render(<LoginForm onSubmit={mockOnSubmit} />);

    await user.type(screen.getByLabelText(/email/i), 'test@example.com');
    await user.type(screen.getByLabelText(/password/i), 'password123');
    await user.click(screen.getByRole('button', { name: /sign in/i }));

    await waitFor(() => {
      expect(mockOnSubmit).toHaveBeenCalledWith({
        email: 'test@example.com',
        password: 'password123',
      });
    });
  });

  it('displays validation errors', async () => {
    const user = userEvent.setup();

    render(<LoginForm onSubmit={mockOnSubmit} />);

    await user.click(screen.getByRole('button', { name: /sign in/i }));

    expect(await screen.findByText(/email is required/i)).toBeInTheDocument();
    expect(mockOnSubmit).not.toHaveBeenCalled();
  });

  it('displays server error', async () => {
    const user = userEvent.setup();
    mockOnSubmit.mockRejectedValue(new Error('Invalid credentials'));

    render(<LoginForm onSubmit={mockOnSubmit} />);

    await user.type(screen.getByLabelText(/email/i), 'test@example.com');
    await user.type(screen.getByLabelText(/password/i), 'password123');
    await user.click(screen.getByRole('button', { name: /sign in/i }));

    expect(await screen.findByText(/invalid credentials/i)).toBeInTheDocument();
  });
});
```

### Pattern 3: Hook Tests

```tsx
import { renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { vi } from 'vitest';
import { useUser } from './use-user';
import * as api from '@/lib/api';

vi.mock('@/lib/api');

function wrapper({ children }: { children: React.ReactNode }) {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return (
    <QueryClientProvider client={queryClient}>
      {children}
    </QueryClientProvider>
  );
}

describe('useUser', () => {
  it('fetches user successfully', async () => {
    const mockUser = { id: '1', name: 'Test User' };
    vi.mocked(api.get).mockResolvedValue(mockUser);

    const { result } = renderHook(() => useUser('1'), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(result.current.data).toEqual(mockUser);
  });

  it('handles error', async () => {
    vi.mocked(api.get).mockRejectedValue(new Error('Not found'));

    const { result } = renderHook(() => useUser('999'), { wrapper });

    await waitFor(() => expect(result.current.isError).toBe(true));

    expect(result.current.error?.message).toBe('Not found');
  });
});
```

### Pattern 4: Accessibility Tests

```tsx
import { render } from '@testing-library/react';
import { axe, toHaveNoViolations } from 'jest-axe';
import { Button } from './button';

expect.extend(toHaveNoViolations);

describe('Button accessibility', () => {
  it('has no violations', async () => {
    const { container } = render(<Button>Click me</Button>);
    const results = await axe(container);
    expect(results).toHaveNoViolations();
  });

  it('has no violations when disabled', async () => {
    const { container } = render(<Button disabled>Click me</Button>);
    const results = await axe(container);
    expect(results).toHaveNoViolations();
  });

  it('has no violations when loading', async () => {
    const { container } = render(<Button isLoading>Click me</Button>);
    const results = await axe(container);
    expect(results).toHaveNoViolations();
  });
});
```

## Test Utilities

### Go: Test Fixtures (DDD Pattern)

Create test fixtures using `Reconstitute` for domain entities with private fields.

```go
// testutil/fixtures.go
package testutil

import (
    "time"

    "github.com/google/uuid"
    "github.com/yourorg/app/internal/domain/user"
    "github.com/yourorg/app/internal/domain/product"
)

// UserFixture provides options for creating test users
type UserFixture struct {
    ID           uuid.UUID
    Email        string
    Name         string
    PasswordHash string
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

// DefaultUserFixture returns default values for a test user
func DefaultUserFixture() UserFixture {
    now := time.Now()
    return UserFixture{
        ID:           uuid.New(),
        Email:        "test@example.com",
        Name:         "Test User",
        PasswordHash: "hashed_password",
        CreatedAt:    now,
        UpdatedAt:    now,
    }
}

// NewTestUser creates a test user using Reconstitute (for entities with private fields)
func NewTestUser(opts ...func(*UserFixture)) *user.User {
    f := DefaultUserFixture()
    for _, opt := range opts {
        opt(&f)
    }
    return user.Reconstitute(f.ID, f.Email, f.Name, f.PasswordHash, f.CreatedAt, f.UpdatedAt)
}

// WithEmail sets a custom email
func WithEmail(email string) func(*UserFixture) {
    return func(f *UserFixture) {
        f.Email = email
    }
}

// WithName sets a custom name
func WithName(name string) func(*UserFixture) {
    return func(f *UserFixture) {
        f.Name = name
    }
}

// WithID sets a specific ID (useful for test assertions)
func WithID(id uuid.UUID) func(*UserFixture) {
    return func(f *UserFixture) {
        f.ID = id
    }
}

// Usage examples:
// Default user:
//   u := testutil.NewTestUser()
//
// Custom email:
//   u := testutil.NewTestUser(testutil.WithEmail("custom@example.com"))
//
// Multiple options:
//   u := testutil.NewTestUser(
//       testutil.WithEmail("admin@example.com"),
//       testutil.WithName("Admin User"),
//   )

// ProductFixture provides options for creating test products
type ProductFixture struct {
    ID          uuid.UUID
    Name        string
    Description string
    Category    string
    Price       int
    Stock       int
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

func DefaultProductFixture() ProductFixture {
    now := time.Now()
    return ProductFixture{
        ID:          uuid.New(),
        Name:        "Test Product",
        Description: "A test product description",
        Category:    "Electronics",
        Price:       1000,
        Stock:       10,
        CreatedAt:   now,
        UpdatedAt:   now,
    }
}

func NewTestProduct(opts ...func(*ProductFixture)) *product.Product {
    f := DefaultProductFixture()
    for _, opt := range opts {
        opt(&f)
    }
    return product.Reconstitute(f.ID, f.Name, f.Description, f.Category, f.Price, f.Stock, f.CreatedAt, f.UpdatedAt)
}
```

### React: Render with Providers

```tsx
// test-utils.tsx
import { render, RenderOptions } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

function createTestQueryClient() {
  return new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  });
}

interface ProvidersProps {
  children: React.ReactNode;
}

function Providers({ children }: ProvidersProps) {
  const queryClient = createTestQueryClient();
  return (
    <QueryClientProvider client={queryClient}>
      {children}
    </QueryClientProvider>
  );
}

export function renderWithProviders(
  ui: React.ReactElement,
  options?: RenderOptions
) {
  return render(ui, { wrapper: Providers, ...options });
}

// Usage
import { renderWithProviders } from '@/test-utils';

renderWithProviders(<MyComponent />);
```

## Testing Checklist

### Go Backend (CQRS)
- [ ] Command handler tests (Create, Update, Delete)
- [ ] Query handler tests (Get, List)
- [ ] Domain entity tests (validation, business rules)
- [ ] Repository adapter integration tests
- [ ] HTTP handler tests with mock handlers
- [ ] Domain errors properly tested
- [ ] Mock repositories use `Reconstitute` for test entities
- [ ] Test fixtures use functional options pattern

### React Frontend
- [ ] Component rendering tests
- [ ] User interaction tests (userEvent)
- [ ] Custom hook tests
- [ ] Error boundary tests
- [ ] Loading states tested
- [ ] Accessibility tests (jest-axe)

### General
- [ ] Edge cases covered
- [ ] Error scenarios tested
- [ ] Mocks properly configured
- [ ] Test data is realistic
- [ ] Tests are independent
- [ ] Integration tests use testcontainers
