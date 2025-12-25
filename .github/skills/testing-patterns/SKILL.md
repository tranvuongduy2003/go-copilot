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

### Pattern 2: Service Tests with Mocks

```go
func TestUserService_CreateUser(t *testing.T) {
    tests := []struct {
        name    string
        input   CreateUserInput
        setup   func(*MockUserRepository, *MockPasswordHasher, *MockEmailService)
        want    *User
        wantErr error
    }{
        {
            name: "success - creates user and sends welcome email",
            input: CreateUserInput{
                Email:    "test@example.com",
                Name:     "Test User",
                Password: "password123",
            },
            setup: func(repo *MockUserRepository, hasher *MockPasswordHasher, email *MockEmailService) {
                repo.On("FindByEmail", mock.Anything, "test@example.com").
                    Return(nil, ErrNotFound)
                hasher.On("Hash", "password123").
                    Return("hashed_password", nil)
                repo.On("Create", mock.Anything, mock.AnythingOfType("*User")).
                    Return(nil)
                email.On("SendWelcome", mock.Anything, "test@example.com", "Test User").
                    Return(nil)
            },
            want: &User{
                Email: "test@example.com",
                Name:  "Test User",
            },
            wantErr: nil,
        },
        {
            name: "error - email already exists",
            input: CreateUserInput{
                Email:    "exists@example.com",
                Name:     "Test User",
                Password: "password123",
            },
            setup: func(repo *MockUserRepository, hasher *MockPasswordHasher, email *MockEmailService) {
                repo.On("FindByEmail", mock.Anything, "exists@example.com").
                    Return(&User{Email: "exists@example.com"}, nil)
            },
            want:    nil,
            wantErr: ErrEmailExists,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            repo := new(MockUserRepository)
            hasher := new(MockPasswordHasher)
            emailSvc := new(MockEmailService)
            tt.setup(repo, hasher, emailSvc)

            svc := NewUserService(repo, hasher, emailSvc, slog.Default())

            // Act
            got, err := svc.CreateUser(context.Background(), tt.input)

            // Assert
            if tt.wantErr != nil {
                require.Error(t, err)
                assert.ErrorIs(t, err, tt.wantErr)
                return
            }

            require.NoError(t, err)
            assert.Equal(t, tt.want.Email, got.Email)
            assert.NotEmpty(t, got.ID)

            // Verify mock expectations
            repo.AssertExpectations(t)
            hasher.AssertExpectations(t)
            emailSvc.AssertExpectations(t)
        })
    }
}
```

### Pattern 3: HTTP Handler Tests

```go
func TestUserHandler_Create(t *testing.T) {
    tests := []struct {
        name       string
        body       interface{}
        setup      func(*MockUserService)
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
            setup: func(m *MockUserService) {
                m.On("CreateUser", mock.Anything, mock.AnythingOfType("CreateUserInput")).
                    Return(&User{
                        ID:    "user-123",
                        Email: "test@example.com",
                        Name:  "Test User",
                    }, nil)
            },
            wantStatus: http.StatusCreated,
            wantBody: func(t *testing.T, body []byte) {
                var resp Response
                require.NoError(t, json.Unmarshal(body, &resp))

                data := resp.Data.(map[string]interface{})
                assert.Equal(t, "user-123", data["id"])
                assert.Equal(t, "test@example.com", data["email"])
            },
        },
        {
            name: "validation error - invalid email",
            body: map[string]string{
                "email":    "invalid-email",
                "name":     "Test User",
                "password": "password123",
            },
            setup:      func(m *MockUserService) {},
            wantStatus: http.StatusBadRequest,
            wantBody: func(t *testing.T, body []byte) {
                var resp ErrorResponse
                require.NoError(t, json.Unmarshal(body, &resp))
                assert.Equal(t, "VALIDATION_ERROR", resp.Error.Code)
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            mockSvc := new(MockUserService)
            tt.setup(mockSvc)

            handler := NewUserHandler(mockSvc, slog.Default())
            router := chi.NewRouter()
            handler.RegisterRoutes(router)

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
        })
    }
}
```

### Pattern 4: Integration Tests with Test Containers

```go
//go:build integration

package integration_test

import (
    "context"
    "testing"

    "github.com/stretchr/testify/suite"
    "github.com/testcontainers/testcontainers-go/modules/postgres"
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

    s.runMigrations()
}

func (s *IntegrationSuite) TearDownSuite() {
    s.db.Close()
    s.container.Terminate(context.Background())
}

func (s *IntegrationSuite) SetupTest() {
    s.db.Exec(context.Background(), "TRUNCATE users, posts CASCADE")
}

func (s *IntegrationSuite) TestCreateAndFetchUser() {
    ctx := context.Background()
    repo := postgres.NewUserRepository(s.db)

    user := &User{
        ID:    "test-id",
        Email: "test@example.com",
        Name:  "Test User",
    }

    err := repo.Create(ctx, user)
    s.NoError(err)

    found, err := repo.FindByID(ctx, user.ID)
    s.NoError(err)
    s.Equal(user.Email, found.Email)
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

### Go: Test Fixtures

```go
// testutil/fixtures.go
package testutil

func NewTestUser(overrides ...func(*User)) *User {
    user := &User{
        ID:        uuid.New().String(),
        Email:     "test@example.com",
        Name:      "Test User",
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    for _, override := range overrides {
        override(user)
    }

    return user
}

// Usage
user := testutil.NewTestUser(func(u *User) {
    u.Email = "custom@example.com"
})
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

- [ ] Unit tests for business logic
- [ ] Integration tests for data layer
- [ ] Handler/component tests for API/UI
- [ ] Edge cases covered
- [ ] Error scenarios tested
- [ ] Loading states tested
- [ ] Accessibility tested
- [ ] Mocks properly configured
- [ ] Test data is realistic
- [ ] Tests are independent
