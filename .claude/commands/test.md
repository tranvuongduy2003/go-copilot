# Testing Command

Create comprehensive test suites for Go backend and React frontend following best practices.

## Task: $ARGUMENTS

## Quick Commands

### Backend (Go)

```bash
# Run all tests
cd backend && go test ./...

# Run tests with coverage
cd backend && go test -cover -coverprofile=coverage.out ./...

# Run tests with verbose output
cd backend && go test -v ./...

# Run specific package tests
cd backend && go test ./internal/application/command/...

# Run tests matching pattern
cd backend && go test -run TestUserService ./...

# Generate HTML coverage report
cd backend && go tool cover -html=coverage.out -o coverage.html
```

### Frontend (React)

```bash
# Run all tests
cd frontend && npm test

# Run tests in watch mode
cd frontend && npm test -- --watch

# Run tests with coverage
cd frontend && npm test -- --coverage

# Run specific test file
cd frontend && npm test -- src/hooks/use-user.test.ts
```

## Backend Testing Patterns (Go + testify)

### Unit Test - Command Handler

```go
package command_test

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/yourorg/app/internal/application/command"
    "github.com/yourorg/app/internal/domain/user"
)

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

func (m *MockUserRepository) Save(ctx context.Context, u *user.User) error {
    return m.Called(ctx, u).Error(0)
}

func TestCreateUserHandler_Handle(t *testing.T) {
    tests := []struct {
        name    string
        cmd     command.CreateUserCommand
        setup   func(*MockUserRepository)
        wantErr bool
    }{
        {
            name: "success",
            cmd: command.CreateUserCommand{
                Email: "test@example.com",
                Name:  "Test User",
            },
            setup: func(m *MockUserRepository) {
                m.On("Save", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)
            },
            wantErr: false,
        },
        {
            name: "validation error - empty name",
            cmd: command.CreateUserCommand{
                Email: "test@example.com",
                Name:  "",
            },
            setup:   func(m *MockUserRepository) {},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            repo := new(MockUserRepository)
            tt.setup(repo)
            handler := command.NewCreateUserHandler(repo)

            result, err := handler.Handle(context.Background(), tt.cmd)

            if tt.wantErr {
                assert.Error(t, err)
                return
            }

            assert.NoError(t, err)
            assert.NotNil(t, result)
            repo.AssertExpectations(t)
        })
    }
}
```

### Unit Test - Domain Entity

```go
package user_test

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/yourorg/app/internal/domain/user"
)

func TestNewUser(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        uname   string
        wantErr error
    }{
        {
            name:    "valid user",
            email:   "test@example.com",
            uname:   "Test User",
            wantErr: nil,
        },
        {
            name:    "empty name",
            email:   "test@example.com",
            uname:   "",
            wantErr: user.ErrInvalidName,
        },
        {
            name:    "invalid email",
            email:   "invalid",
            uname:   "Test User",
            wantErr: user.ErrInvalidEmail,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            email, _ := user.NewEmail(tt.email)
            u, err := user.NewUser(email, tt.uname, user.RoleUser)

            if tt.wantErr != nil {
                assert.ErrorIs(t, err, tt.wantErr)
                return
            }

            assert.NoError(t, err)
            assert.NotNil(t, u)
            assert.Equal(t, tt.uname, u.Name())
        })
    }
}
```

### HTTP Handler Test

```go
package handler_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/go-chi/chi/v5"
    "github.com/stretchr/testify/assert"
)

func TestUserHandler_Create(t *testing.T) {
    tests := []struct {
        name       string
        body       string
        wantStatus int
    }{
        {
            name:       "success",
            body:       `{"email":"test@example.com","name":"Test","password":"password123"}`,
            wantStatus: http.StatusCreated,
        },
        {
            name:       "invalid json",
            body:       `{invalid}`,
            wantStatus: http.StatusBadRequest,
        },
        {
            name:       "missing required field",
            body:       `{"email":"test@example.com"}`,
            wantStatus: http.StatusBadRequest,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            handler := setupTestHandler(t)
            router := chi.NewRouter()
            handler.RegisterRoutes(router)

            req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(tt.body))
            req.Header.Set("Content-Type", "application/json")
            rec := httptest.NewRecorder()

            router.ServeHTTP(rec, req)

            assert.Equal(t, tt.wantStatus, rec.Code)
        })
    }
}
```

### Integration Test

```go
//go:build integration

package postgres_test

import (
    "context"
    "testing"

    "github.com/stretchr/testify/suite"
)

type UserRepositorySuite struct {
    suite.Suite
    repo *postgres.UserRepository
    db   *pgxpool.Pool
}

func (s *UserRepositorySuite) SetupSuite() {
    s.db = setupTestDB(s.T())
    s.repo = postgres.NewUserRepository(s.db)
}

func (s *UserRepositorySuite) TearDownSuite() {
    s.db.Close()
}

func (s *UserRepositorySuite) SetupTest() {
    s.db.Exec(context.Background(), "TRUNCATE users CASCADE")
}

func (s *UserRepositorySuite) TestSave_Success() {
    user, _ := user.NewUser(email, "Test User", user.RoleUser)

    err := s.repo.Save(context.Background(), user)
    s.NoError(err)

    found, err := s.repo.FindByID(context.Background(), user.ID())
    s.NoError(err)
    s.Equal(user.Name(), found.Name())
}

func TestUserRepositorySuite(t *testing.T) {
    suite.Run(t, new(UserRepositorySuite))
}
```

## Frontend Testing Patterns (Vitest + Testing Library)

### Component Test

```tsx
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { UserCard } from './user-card';

describe('UserCard', () => {
  const mockUser = {
    id: '1',
    name: 'John Doe',
    email: 'john@example.com',
  };

  it('renders user information', () => {
    render(<UserCard user={mockUser} />);

    expect(screen.getByText('John Doe')).toBeInTheDocument();
    expect(screen.getByText('john@example.com')).toBeInTheDocument();
  });

  it('calls onEdit when edit button clicked', () => {
    const onEdit = vi.fn();
    render(<UserCard user={mockUser} onEdit={onEdit} />);

    fireEvent.click(screen.getByRole('button', { name: /edit/i }));

    expect(onEdit).toHaveBeenCalledTimes(1);
  });
});
```

### Hook Test

```tsx
import { renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { describe, it, expect, vi } from 'vitest';
import { useUsers } from './use-users';
import * as api from '@/lib/api';

vi.mock('@/lib/api');

const createWrapper = () => {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
};

describe('useUsers', () => {
  it('fetches users successfully', async () => {
    const mockUsers = [{ id: '1', name: 'John' }];
    vi.mocked(api.get).mockResolvedValue(mockUsers);

    const { result } = renderHook(() => useUsers(), { wrapper: createWrapper() });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(result.current.data).toEqual(mockUsers);
  });

  it('handles error', async () => {
    vi.mocked(api.get).mockRejectedValue(new Error('Network error'));

    const { result } = renderHook(() => useUsers(), { wrapper: createWrapper() });

    await waitFor(() => expect(result.current.isError).toBe(true));
  });
});
```

### Form Test

```tsx
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi } from 'vitest';
import { LoginForm } from './login-form';

describe('LoginForm', () => {
  it('submits form with valid data', async () => {
    const onSubmit = vi.fn();
    render(<LoginForm onSubmit={onSubmit} />);

    await userEvent.type(screen.getByLabelText(/email/i), 'test@example.com');
    await userEvent.type(screen.getByLabelText(/password/i), 'password123');
    fireEvent.click(screen.getByRole('button', { name: /login/i }));

    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalledWith({
        email: 'test@example.com',
        password: 'password123',
      });
    });
  });

  it('shows validation errors', async () => {
    render(<LoginForm onSubmit={vi.fn()} />);

    fireEvent.click(screen.getByRole('button', { name: /login/i }));

    await waitFor(() => {
      expect(screen.getByText(/email is required/i)).toBeInTheDocument();
    });
  });
});
```

## Testing Best Practices

### Always Do

- Write tests BEFORE or alongside code (TDD/BDD)
- Use table-driven tests for multiple scenarios
- Mock external dependencies
- Test edge cases and error conditions
- Keep tests independent and isolated
- Use descriptive test names

### Test Coverage Goals

- Domain logic: 90%+
- Application layer: 80%+
- HTTP handlers: 70%+
- Frontend components: 70%+

### Never Do

- Never test implementation details
- Never write tests that depend on other tests
- Never test private methods directly
- Never ignore flaky tests
