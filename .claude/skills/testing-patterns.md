# Testing Patterns Skill

Generate comprehensive test suites for Go and React code.

## Usage

```
/project:skill:testing <file-path>
```

## Go Testing Patterns

### Table-Driven Unit Tests

```go
package user_test

import (
    "context"
    "testing"

    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "yourapp/internal/domain/user"
)

func TestNewUser(t *testing.T) {
    tests := []struct {
        name        string
        email       string
        userName    string
        role        user.Role
        expectedErr error
    }{
        {
            name:        "valid user",
            email:       "test@example.com",
            userName:    "Test User",
            role:        user.RoleUser,
            expectedErr: nil,
        },
        {
            name:        "empty name",
            email:       "test@example.com",
            userName:    "",
            role:        user.RoleUser,
            expectedErr: user.ErrInvalidName,
        },
        {
            name:        "invalid email",
            email:       "invalid-email",
            userName:    "Test User",
            role:        user.RoleUser,
            expectedErr: user.ErrInvalidEmail,
        },
        {
            name:        "invalid role",
            email:       "test@example.com",
            userName:    "Test User",
            role:        user.Role("invalid"),
            expectedErr: user.ErrInvalidRole,
        },
    }

    for _, testCase := range tests {
        t.Run(testCase.name, func(t *testing.T) {
            email, emailErr := user.NewEmail(testCase.email)
            if testCase.expectedErr == user.ErrInvalidEmail {
                assert.Error(t, emailErr)
                return
            }
            require.NoError(t, emailErr)

            createdUser, err := user.NewUser(email, testCase.userName, testCase.role)

            if testCase.expectedErr != nil {
                assert.ErrorIs(t, err, testCase.expectedErr)
                assert.Nil(t, createdUser)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, createdUser)
                assert.Equal(t, testCase.userName, createdUser.Name())
                assert.Equal(t, testCase.role, createdUser.Role())
            }
        })
    }
}
```

### Mock Repository Tests

```go
package command_test

import (
    "context"
    "testing"

    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "yourapp/internal/application/command"
    "yourapp/internal/domain/user"
)

type MockUserRepository struct {
    mock.Mock
}

func (mockRepository *MockUserRepository) FindByID(context context.Context, id uuid.UUID) (*user.User, error) {
    arguments := mockRepository.Called(context, id)
    if arguments.Get(0) == nil {
        return nil, arguments.Error(1)
    }
    return arguments.Get(0).(*user.User), arguments.Error(1)
}

func (mockRepository *MockUserRepository) Save(context context.Context, userEntity *user.User) error {
    arguments := mockRepository.Called(context, userEntity)
    return arguments.Error(0)
}

func (mockRepository *MockUserRepository) Update(context context.Context, userEntity *user.User) error {
    arguments := mockRepository.Called(context, userEntity)
    return arguments.Error(0)
}

func (mockRepository *MockUserRepository) Delete(context context.Context, id uuid.UUID) error {
    arguments := mockRepository.Called(context, id)
    return arguments.Error(0)
}

func TestCreateUserHandler_Handle(t *testing.T) {
    tests := []struct {
        name           string
        command        command.CreateUserCommand
        setupMock      func(*MockUserRepository)
        expectedError  bool
        expectedUserID bool
    }{
        {
            name: "successful creation",
            command: command.CreateUserCommand{
                Email: "test@example.com",
                Name:  "Test User",
            },
            setupMock: func(mockRepository *MockUserRepository) {
                mockRepository.On("Save", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)
            },
            expectedError:  false,
            expectedUserID: true,
        },
        {
            name: "repository error",
            command: command.CreateUserCommand{
                Email: "test@example.com",
                Name:  "Test User",
            },
            setupMock: func(mockRepository *MockUserRepository) {
                mockRepository.On("Save", mock.Anything, mock.AnythingOfType("*user.User")).
                    Return(errors.New("database error"))
            },
            expectedError:  true,
            expectedUserID: false,
        },
        {
            name: "validation error - empty name",
            command: command.CreateUserCommand{
                Email: "test@example.com",
                Name:  "",
            },
            setupMock:      func(mockRepository *MockUserRepository) {},
            expectedError:  true,
            expectedUserID: false,
        },
    }

    for _, testCase := range tests {
        t.Run(testCase.name, func(t *testing.T) {
            mockRepository := new(MockUserRepository)
            testCase.setupMock(mockRepository)

            handler := command.NewCreateUserHandler(mockRepository)

            result, err := handler.Handle(context.Background(), testCase.command)

            if testCase.expectedError {
                assert.Error(t, err)
                assert.Nil(t, result)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, result)
                if testCase.expectedUserID {
                    assert.NotEqual(t, uuid.Nil, result.ID)
                }
            }

            mockRepository.AssertExpectations(t)
        })
    }
}
```

### HTTP Handler Tests

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
    "github.com/stretchr/testify/mock"
    "yourapp/internal/interfaces/http/handler"
)

func TestUserHandler_Create(t *testing.T) {
    tests := []struct {
        name           string
        requestBody    string
        setupMocks     func(*MockCreateUserHandler)
        expectedStatus int
        expectedBody   string
    }{
        {
            name:        "successful creation",
            requestBody: `{"email":"test@example.com","name":"Test User","password":"password123"}`,
            setupMocks: func(mockHandler *MockCreateUserHandler) {
                mockHandler.On("Handle", mock.Anything, mock.Anything).Return(&dto.UserDTO{
                    ID:    uuid.New(),
                    Email: "test@example.com",
                    Name:  "Test User",
                }, nil)
            },
            expectedStatus: http.StatusCreated,
        },
        {
            name:           "invalid json",
            requestBody:    `{invalid`,
            setupMocks:     func(mockHandler *MockCreateUserHandler) {},
            expectedStatus: http.StatusBadRequest,
        },
        {
            name:        "validation error",
            requestBody: `{"email":"invalid","name":"","password":"short"}`,
            setupMocks: func(mockHandler *MockCreateUserHandler) {
                mockHandler.On("Handle", mock.Anything, mock.Anything).
                    Return(nil, user.ErrInvalidEmail)
            },
            expectedStatus: http.StatusBadRequest,
        },
    }

    for _, testCase := range tests {
        t.Run(testCase.name, func(t *testing.T) {
            mockCreateHandler := new(MockCreateUserHandler)
            testCase.setupMocks(mockCreateHandler)

            userHandler := handler.NewUserHandler(mockCreateHandler, nil, nil)
            router := chi.NewRouter()
            userHandler.RegisterRoutes(router)

            request := httptest.NewRequest(
                http.MethodPost,
                "/users",
                bytes.NewBufferString(testCase.requestBody),
            )
            request.Header.Set("Content-Type", "application/json")
            recorder := httptest.NewRecorder()

            router.ServeHTTP(recorder, request)

            assert.Equal(t, testCase.expectedStatus, recorder.Code)
            mockCreateHandler.AssertExpectations(t)
        })
    }
}
```

## React Testing Patterns

### Component Tests

```tsx
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { UserCard } from './user-card';
import type { User } from '@/types/user';

const createTestQueryClient = () =>
    new QueryClient({
        defaultOptions: {
            queries: { retry: false },
            mutations: { retry: false },
        },
    });

const renderWithProviders = (ui: React.ReactElement) => {
    const testQueryClient = createTestQueryClient();
    return render(
        <QueryClientProvider client={testQueryClient}>{ui}</QueryClientProvider>
    );
};

const mockUser: User = {
    id: '1',
    name: 'John Doe',
    email: 'john@example.com',
    role: 'user',
    status: 'active',
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z',
};

describe('UserCard', () => {
    it('renders user information correctly', () => {
        renderWithProviders(<UserCard user={mockUser} />);

        expect(screen.getByText('John Doe')).toBeInTheDocument();
        expect(screen.getByText('john@example.com')).toBeInTheDocument();
    });

    it('displays avatar with initials when no image provided', () => {
        renderWithProviders(<UserCard user={mockUser} />);

        expect(screen.getByText('JD')).toBeInTheDocument();
    });

    it('calls onEdit callback when edit button is clicked', async () => {
        const user = userEvent.setup();
        const handleEdit = vi.fn();

        renderWithProviders(<UserCard user={mockUser} onEdit={handleEdit} />);

        await user.click(screen.getByRole('button', { name: /edit/i }));

        expect(handleEdit).toHaveBeenCalledTimes(1);
        expect(handleEdit).toHaveBeenCalledWith(mockUser);
    });

    it('calls onDelete callback when delete button is clicked', async () => {
        const user = userEvent.setup();
        const handleDelete = vi.fn();

        renderWithProviders(<UserCard user={mockUser} onDelete={handleDelete} />);

        await user.click(screen.getByRole('button', { name: /delete/i }));

        expect(handleDelete).toHaveBeenCalledTimes(1);
        expect(handleDelete).toHaveBeenCalledWith(mockUser);
    });

    it('does not render action buttons when callbacks not provided', () => {
        renderWithProviders(<UserCard user={mockUser} />);

        expect(screen.queryByRole('button', { name: /edit/i })).not.toBeInTheDocument();
        expect(screen.queryByRole('button', { name: /delete/i })).not.toBeInTheDocument();
    });

    it('applies correct variant styles', () => {
        const { rerender } = renderWithProviders(
            <UserCard user={mockUser} variant="compact" data-testid="card" />
        );

        expect(screen.getByTestId('card')).toHaveClass('p-2');

        rerender(
            <QueryClientProvider client={createTestQueryClient()}>
                <UserCard user={mockUser} variant="detailed" data-testid="card" />
            </QueryClientProvider>
        );

        expect(screen.getByTestId('card')).toHaveClass('p-6');
    });
});
```

### Hook Tests

```tsx
import { renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useUsers, useCreateUser } from '@/hooks/use-user';
import * as api from '@/lib/api';

vi.mock('@/lib/api');

const createWrapper = () => {
    const queryClient = new QueryClient({
        defaultOptions: {
            queries: { retry: false },
            mutations: { retry: false },
        },
    });

    return ({ children }: { children: React.ReactNode }) => (
        <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    );
};

describe('useUsers', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('fetches users successfully', async () => {
        const mockUsers = [
            { id: '1', name: 'User 1', email: 'user1@example.com' },
            { id: '2', name: 'User 2', email: 'user2@example.com' },
        ];

        vi.mocked(api.get).mockResolvedValue(mockUsers);

        const { result } = renderHook(() => useUsers(), { wrapper: createWrapper() });

        await waitFor(() => expect(result.current.isSuccess).toBe(true));

        expect(result.current.data).toEqual(mockUsers);
        expect(api.get).toHaveBeenCalledWith('/users');
    });

    it('handles error state', async () => {
        vi.mocked(api.get).mockRejectedValue(new Error('Network error'));

        const { result } = renderHook(() => useUsers(), { wrapper: createWrapper() });

        await waitFor(() => expect(result.current.isError).toBe(true));

        expect(result.current.error).toBeDefined();
    });
});

describe('useCreateUser', () => {
    it('creates user and invalidates queries', async () => {
        const newUser = { id: '3', name: 'New User', email: 'new@example.com' };
        vi.mocked(api.post).mockResolvedValue(newUser);

        const { result } = renderHook(() => useCreateUser(), { wrapper: createWrapper() });

        await result.current.mutateAsync({
            name: 'New User',
            email: 'new@example.com',
            password: 'password123',
        });

        expect(api.post).toHaveBeenCalledWith('/users', {
            name: 'New User',
            email: 'new@example.com',
            password: 'password123',
        });
    });
});
```

### Form Tests

```tsx
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi } from 'vitest';
import { LoginForm } from './login-form';

describe('LoginForm', () => {
    it('submits form with valid data', async () => {
        const user = userEvent.setup();
        const handleSubmit = vi.fn().mockResolvedValue(undefined);

        render(<LoginForm onSubmit={handleSubmit} />);

        await user.type(screen.getByLabelText(/email/i), 'test@example.com');
        await user.type(screen.getByLabelText(/password/i), 'password123');
        await user.click(screen.getByRole('button', { name: /sign in/i }));

        await waitFor(() => {
            expect(handleSubmit).toHaveBeenCalledWith({
                email: 'test@example.com',
                password: 'password123',
            });
        });
    });

    it('displays validation errors for invalid input', async () => {
        const user = userEvent.setup();
        const handleSubmit = vi.fn();

        render(<LoginForm onSubmit={handleSubmit} />);

        await user.type(screen.getByLabelText(/email/i), 'invalid-email');
        await user.type(screen.getByLabelText(/password/i), 'short');
        await user.click(screen.getByRole('button', { name: /sign in/i }));

        await waitFor(() => {
            expect(screen.getByText(/invalid email/i)).toBeInTheDocument();
            expect(screen.getByText(/at least 8 characters/i)).toBeInTheDocument();
        });

        expect(handleSubmit).not.toHaveBeenCalled();
    });

    it('disables submit button while submitting', async () => {
        const user = userEvent.setup();
        const handleSubmit = vi.fn().mockImplementation(
            () => new Promise((resolve) => setTimeout(resolve, 100))
        );

        render(<LoginForm onSubmit={handleSubmit} />);

        await user.type(screen.getByLabelText(/email/i), 'test@example.com');
        await user.type(screen.getByLabelText(/password/i), 'password123');
        await user.click(screen.getByRole('button', { name: /sign in/i }));

        expect(screen.getByRole('button', { name: /loading/i })).toBeDisabled();
    });
});
```

## Checklist

- [ ] Table-driven tests for multiple scenarios
- [ ] Mock external dependencies
- [ ] Test error paths
- [ ] Test edge cases
- [ ] Test loading and error states (React)
- [ ] Test user interactions (React)
- [ ] Test accessibility
- [ ] Meaningful test names
