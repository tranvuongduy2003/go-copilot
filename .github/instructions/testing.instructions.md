---
applyTo: "**/*_test.go,**/*.test.{ts,tsx}"
---

# Testing Instructions

These instructions apply to all test files in both Go and React.

## Go Testing

### Test File Organization

```
internal/
├── service/
│   ├── user_service.go
│   └── user_service_test.go      # Unit tests
├── repository/
│   └── postgres/
│       ├── user_repository.go
│       └── user_repository_test.go
└── handlers/
    ├── user_handler.go
    └── user_handler_test.go

tests/
└── integration/
    └── user_test.go              # Integration tests
```

### Table-Driven Tests

```go
func TestUserService_CreateUser(t *testing.T) {
    tests := []struct {
        name    string
        input   domain.CreateUserInput
        setup   func(*MockUserRepository, *MockPasswordHasher)
        want    *domain.User
        wantErr error
    }{
        {
            name: "success - creates user with hashed password",
            input: domain.CreateUserInput{
                Email:    "test@example.com",
                Name:     "Test User",
                Password: "password123",
            },
            setup: func(repo *MockUserRepository, hasher *MockPasswordHasher) {
                repo.On("FindByEmail", mock.Anything, "test@example.com").
                    Return(nil, domain.ErrNotFound)
                hasher.On("Hash", "password123").
                    Return("hashed_password", nil)
                repo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).
                    Return(nil)
            },
            want: &domain.User{
                Email: "test@example.com",
                Name:  "Test User",
            },
            wantErr: nil,
        },
        {
            name: "error - email already exists",
            input: domain.CreateUserInput{
                Email:    "existing@example.com",
                Name:     "Existing User",
                Password: "password123",
            },
            setup: func(repo *MockUserRepository, hasher *MockPasswordHasher) {
                repo.On("FindByEmail", mock.Anything, "existing@example.com").
                    Return(&domain.User{Email: "existing@example.com"}, nil)
            },
            want:    nil,
            wantErr: domain.ErrConflict,
        },
        {
            name: "error - password hashing fails",
            input: domain.CreateUserInput{
                Email:    "test@example.com",
                Name:     "Test User",
                Password: "password123",
            },
            setup: func(repo *MockUserRepository, hasher *MockPasswordHasher) {
                repo.On("FindByEmail", mock.Anything, "test@example.com").
                    Return(nil, domain.ErrNotFound)
                hasher.On("Hash", "password123").
                    Return("", errors.New("hash error"))
            },
            want:    nil,
            wantErr: errors.New("hash password"),
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            mockRepo := new(MockUserRepository)
            mockHasher := new(MockPasswordHasher)
            tt.setup(mockRepo, mockHasher)

            svc := service.NewUserService(mockRepo, mockHasher, slog.Default())

            // Act
            got, err := svc.CreateUser(context.Background(), tt.input)

            // Assert
            if tt.wantErr != nil {
                require.Error(t, err)
                assert.Contains(t, err.Error(), tt.wantErr.Error())
                return
            }

            require.NoError(t, err)
            assert.Equal(t, tt.want.Email, got.Email)
            assert.Equal(t, tt.want.Name, got.Name)
            assert.NotEmpty(t, got.ID)

            mockRepo.AssertExpectations(t)
            mockHasher.AssertExpectations(t)
        })
    }
}
```

### HTTP Handler Tests

```go
func TestUserHandler_Get(t *testing.T) {
    tests := []struct {
        name       string
        userID     string
        setup      func(*MockUserService)
        wantStatus int
        wantBody   map[string]interface{}
    }{
        {
            name:   "success - returns user",
            userID: "user-123",
            setup: func(m *MockUserService) {
                m.On("GetUser", mock.Anything, "user-123").
                    Return(&domain.User{
                        ID:    "user-123",
                        Email: "test@example.com",
                        Name:  "Test User",
                    }, nil)
            },
            wantStatus: http.StatusOK,
            wantBody: map[string]interface{}{
                "data": map[string]interface{}{
                    "id":    "user-123",
                    "email": "test@example.com",
                    "name":  "Test User",
                },
            },
        },
        {
            name:   "error - user not found",
            userID: "nonexistent",
            setup: func(m *MockUserService) {
                m.On("GetUser", mock.Anything, "nonexistent").
                    Return(nil, domain.ErrNotFound)
            },
            wantStatus: http.StatusNotFound,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockSvc := new(MockUserService)
            tt.setup(mockSvc)

            handler := NewUserHandler(mockSvc, slog.Default())
            router := chi.NewRouter()
            handler.RegisterRoutes(router)

            req := httptest.NewRequest(http.MethodGet, "/users/"+tt.userID, nil)
            rec := httptest.NewRecorder()

            router.ServeHTTP(rec, req)

            assert.Equal(t, tt.wantStatus, rec.Code)

            if tt.wantBody != nil {
                var got map[string]interface{}
                json.Unmarshal(rec.Body.Bytes(), &got)

                data := got["data"].(map[string]interface{})
                wantData := tt.wantBody["data"].(map[string]interface{})

                assert.Equal(t, wantData["id"], data["id"])
                assert.Equal(t, wantData["email"], data["email"])
            }
        })
    }
}
```

### Integration Tests

```go
//go:build integration

package integration_test

import (
    "context"
    "testing"

    "github.com/stretchr/testify/suite"
    "github.com/testcontainers/testcontainers-go/modules/postgres"
)

type UserIntegrationSuite struct {
    suite.Suite
    container *postgres.PostgresContainer
    db        *pgxpool.Pool
    svc       *service.UserService
}

func (s *UserIntegrationSuite) SetupSuite() {
    ctx := context.Background()

    // Start PostgreSQL container
    container, err := postgres.Run(ctx,
        "postgres:16-alpine",
        postgres.WithDatabase("testdb"),
        postgres.WithUsername("test"),
        postgres.WithPassword("test"),
    )
    s.Require().NoError(err)
    s.container = container

    // Connect to database
    connStr, err := container.ConnectionString(ctx, "sslmode=disable")
    s.Require().NoError(err)

    s.db, err = pgxpool.New(ctx, connStr)
    s.Require().NoError(err)

    // Run migrations
    s.runMigrations()

    // Create service with real repository
    repo := postgres.NewUserRepository(s.db)
    hasher := bcrypt.NewHasher()
    s.svc = service.NewUserService(repo, hasher, slog.Default())
}

func (s *UserIntegrationSuite) TearDownSuite() {
    s.db.Close()
    s.container.Terminate(context.Background())
}

func (s *UserIntegrationSuite) SetupTest() {
    // Clean database before each test
    s.db.Exec(context.Background(), "TRUNCATE users CASCADE")
}

func (s *UserIntegrationSuite) TestCreateAndGetUser() {
    ctx := context.Background()

    // Create user
    input := domain.CreateUserInput{
        Email:    "test@example.com",
        Name:     "Test User",
        Password: "password123",
    }

    created, err := s.svc.CreateUser(ctx, input)
    s.NoError(err)
    s.NotEmpty(created.ID)
    s.Equal(input.Email, created.Email)

    // Get user
    found, err := s.svc.GetUser(ctx, created.ID)
    s.NoError(err)
    s.Equal(created.ID, found.ID)
    s.Equal(created.Email, found.Email)
}

func (s *UserIntegrationSuite) TestCreateUser_DuplicateEmail() {
    ctx := context.Background()

    input := domain.CreateUserInput{
        Email:    "duplicate@example.com",
        Name:     "First User",
        Password: "password123",
    }

    // Create first user
    _, err := s.svc.CreateUser(ctx, input)
    s.NoError(err)

    // Try to create second user with same email
    input.Name = "Second User"
    _, err = s.svc.CreateUser(ctx, input)
    s.ErrorIs(err, domain.ErrConflict)
}

func TestUserIntegrationSuite(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration tests")
    }
    suite.Run(t, new(UserIntegrationSuite))
}
```

## React Testing

### Component Testing

```tsx
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { vi } from 'vitest';

// Test utilities
function createTestQueryClient() {
  return new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  });
}

function renderWithProviders(ui: React.ReactElement) {
  const queryClient = createTestQueryClient();
  return {
    ...render(
      <QueryClientProvider client={queryClient}>
        {ui}
      </QueryClientProvider>
    ),
    queryClient,
  };
}

describe('LoginForm', () => {
  const mockOnSubmit = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders all form fields', () => {
    renderWithProviders(<LoginForm onSubmit={mockOnSubmit} />);

    expect(screen.getByLabelText(/email/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/password/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /sign in/i })).toBeInTheDocument();
  });

  it('shows validation errors for empty submission', async () => {
    const user = userEvent.setup();
    renderWithProviders(<LoginForm onSubmit={mockOnSubmit} />);

    await user.click(screen.getByRole('button', { name: /sign in/i }));

    expect(await screen.findByText(/email is required/i)).toBeInTheDocument();
    expect(await screen.findByText(/password is required/i)).toBeInTheDocument();
    expect(mockOnSubmit).not.toHaveBeenCalled();
  });

  it('shows error for invalid email format', async () => {
    const user = userEvent.setup();
    renderWithProviders(<LoginForm onSubmit={mockOnSubmit} />);

    await user.type(screen.getByLabelText(/email/i), 'invalid-email');
    await user.type(screen.getByLabelText(/password/i), 'password123');
    await user.click(screen.getByRole('button', { name: /sign in/i }));

    expect(await screen.findByText(/invalid email/i)).toBeInTheDocument();
    expect(mockOnSubmit).not.toHaveBeenCalled();
  });

  it('submits form with valid data', async () => {
    const user = userEvent.setup();
    mockOnSubmit.mockResolvedValue(undefined);
    renderWithProviders(<LoginForm onSubmit={mockOnSubmit} />);

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

  it('shows loading state during submission', async () => {
    const user = userEvent.setup();
    mockOnSubmit.mockImplementation(() => new Promise(() => {})); // Never resolves
    renderWithProviders(<LoginForm onSubmit={mockOnSubmit} />);

    await user.type(screen.getByLabelText(/email/i), 'test@example.com');
    await user.type(screen.getByLabelText(/password/i), 'password123');
    await user.click(screen.getByRole('button', { name: /sign in/i }));

    expect(await screen.findByText(/loading/i)).toBeInTheDocument();
    expect(screen.getByRole('button')).toBeDisabled();
  });

  it('displays server error message', async () => {
    const user = userEvent.setup();
    mockOnSubmit.mockRejectedValue(new Error('Invalid credentials'));
    renderWithProviders(<LoginForm onSubmit={mockOnSubmit} />);

    await user.type(screen.getByLabelText(/email/i), 'test@example.com');
    await user.type(screen.getByLabelText(/password/i), 'password123');
    await user.click(screen.getByRole('button', { name: /sign in/i }));

    expect(await screen.findByText(/invalid credentials/i)).toBeInTheDocument();
  });
});
```

### Hook Testing

```tsx
import { renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { vi } from 'vitest';
import { useUser, useCreateUser } from './use-user';
import * as api from '@/lib/api';

vi.mock('@/lib/api');

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
}

describe('useUser', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('fetches user successfully', async () => {
    const mockUser = { id: '1', name: 'Test User', email: 'test@example.com' };
    vi.mocked(api.get).mockResolvedValue(mockUser);

    const { result } = renderHook(() => useUser('1'), {
      wrapper: createWrapper(),
    });

    expect(result.current.isLoading).toBe(true);

    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true);
    });

    expect(result.current.data).toEqual(mockUser);
    expect(api.get).toHaveBeenCalledWith('/users/1');
  });

  it('handles fetch error', async () => {
    vi.mocked(api.get).mockRejectedValue(new Error('User not found'));

    const { result } = renderHook(() => useUser('999'), {
      wrapper: createWrapper(),
    });

    await waitFor(() => {
      expect(result.current.isError).toBe(true);
    });

    expect(result.current.error?.message).toBe('User not found');
  });

  it('does not fetch when id is empty', () => {
    renderHook(() => useUser(''), {
      wrapper: createWrapper(),
    });

    expect(api.get).not.toHaveBeenCalled();
  });
});

describe('useCreateUser', () => {
  it('creates user and invalidates queries', async () => {
    const mockUser = { id: '1', name: 'New User', email: 'new@example.com' };
    vi.mocked(api.post).mockResolvedValue(mockUser);

    const queryClient = new QueryClient();
    const invalidateSpy = vi.spyOn(queryClient, 'invalidateQueries');

    const { result } = renderHook(() => useCreateUser(), {
      wrapper: ({ children }) => (
        <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
      ),
    });

    result.current.mutate({ name: 'New User', email: 'new@example.com' });

    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true);
    });

    expect(result.current.data).toEqual(mockUser);
    expect(invalidateSpy).toHaveBeenCalledWith({ queryKey: ['users'] });
  });
});
```

### Accessibility Testing

```tsx
import { render } from '@testing-library/react';
import { axe, toHaveNoViolations } from 'jest-axe';
import { Button } from './button';
import { Card } from './card';
import { LoginForm } from './login-form';

expect.extend(toHaveNoViolations);

describe('Accessibility', () => {
  it('Button has no accessibility violations', async () => {
    const { container } = render(<Button>Click me</Button>);
    const results = await axe(container);
    expect(results).toHaveNoViolations();
  });

  it('Card has no accessibility violations', async () => {
    const { container } = render(
      <Card>
        <h2>Card Title</h2>
        <p>Card content</p>
      </Card>
    );
    const results = await axe(container);
    expect(results).toHaveNoViolations();
  });

  it('LoginForm has no accessibility violations', async () => {
    const { container } = render(<LoginForm onSubmit={() => {}} />);
    const results = await axe(container);
    expect(results).toHaveNoViolations();
  });
});
```

## Testing Best Practices

### General
1. **Test behavior, not implementation** - Focus on what the code does
2. **Use descriptive test names** - Should describe the scenario
3. **Follow AAA pattern** - Arrange, Act, Assert
4. **Keep tests independent** - Each test is self-contained
5. **Don't test external libraries** - Mock them instead

### Go Specific
1. Use table-driven tests for multiple scenarios
2. Use `testify/assert` and `testify/require`
3. Use `testify/mock` for mocking
4. Tag integration tests with `//go:build integration`
5. Use test containers for database tests

### React Specific
1. Query by role/label, not test IDs
2. Use `userEvent` over `fireEvent`
3. Wrap in `waitFor` for async operations
4. Test accessibility with `jest-axe`
5. Mock API calls, not components

### Coverage Goals
- **Minimum**: 80% for new code
- **Critical paths**: 100% (auth, payments)
- **Focus on**: Edge cases, error handling
