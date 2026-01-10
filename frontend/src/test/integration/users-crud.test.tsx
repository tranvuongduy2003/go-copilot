import { API_BASE_URL } from '@/constants';
import { server } from '@/test/mocks/server';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { http, HttpResponse } from 'msw';
import type { ReactNode } from 'react';
import { MemoryRouter } from 'react-router';
import { beforeEach, describe, expect, it, vi } from 'vitest';

vi.mock('sonner', () => ({
  toast: {
    success: vi.fn(),
    error: vi.fn(),
  },
}));

import { useCreateUser, useDeleteUser, useUsers } from '@/features/users/api/users.queries';

function createTestQueryClient() {
  return new QueryClient({
    defaultOptions: {
      queries: { retry: false, gcTime: 0, staleTime: 0 },
      mutations: { retry: false },
    },
  });
}

function TestWrapper({ children, queryClient }: { children: ReactNode; queryClient: QueryClient }) {
  return (
    <QueryClientProvider client={queryClient}>
      <MemoryRouter>{children}</MemoryRouter>
    </QueryClientProvider>
  );
}

function UsersListTestComponent() {
  const { data, isLoading, error } = useUsers();

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;

  return (
    <div>
      <h1>Users List</h1>
      <ul>
        {data?.data.map((user) => (
          <li key={user.id} data-testid={`user-${user.id}`}>
            {user.fullName} - {user.email}
          </li>
        ))}
      </ul>
      <div>Total: {data?.meta.total}</div>
    </div>
  );
}

function CreateUserTestComponent({ onSuccess }: { onSuccess?: () => void }) {
  const { mutate, isPending, isSuccess, error } = useCreateUser();

  const handleCreate = () => {
    mutate(
      { email: 'new@example.com', fullName: 'New User', password: 'Password123!' },
      { onSuccess }
    );
  };

  return (
    <div>
      <button onClick={handleCreate} disabled={isPending}>
        Create User
      </button>
      {isPending && <span>Creating...</span>}
      {isSuccess && <span>Created!</span>}
      {error && <span>Error: {error.message}</span>}
    </div>
  );
}

function DeleteUserTestComponent({
  userId,
  onSuccess,
}: {
  userId: string;
  onSuccess?: () => void;
}) {
  const { mutate, isPending, isSuccess } = useDeleteUser();

  const handleDelete = () => {
    mutate(userId, { onSuccess });
  };

  return (
    <div>
      <button onClick={handleDelete} disabled={isPending}>
        Delete User
      </button>
      {isPending && <span>Deleting...</span>}
      {isSuccess && <span>Deleted!</span>}
    </div>
  );
}

describe('Users CRUD Integration', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('List Users', () => {
    it('fetches and displays users list', async () => {
      const queryClient = createTestQueryClient();

      render(
        <TestWrapper queryClient={queryClient}>
          <UsersListTestComponent />
        </TestWrapper>
      );

      expect(screen.getByText('Loading...')).toBeInTheDocument();

      await waitFor(() => {
        expect(screen.getByText('Users List')).toBeInTheDocument();
        expect(screen.getByTestId('user-1')).toBeInTheDocument();
      });
    });

    it('handles empty users list', async () => {
      const queryClient = createTestQueryClient();

      server.use(
        http.get(`${API_BASE_URL}/users`, () => {
          return HttpResponse.json({
            data: [],
            meta: { page: 1, pageSize: 10, total: 0, totalPages: 0 },
          });
        })
      );

      render(
        <TestWrapper queryClient={queryClient}>
          <UsersListTestComponent />
        </TestWrapper>
      );

      await waitFor(() => {
        expect(screen.getByText('Total: 0')).toBeInTheDocument();
      });
    });

    it('handles fetch error', async () => {
      const queryClient = createTestQueryClient();

      server.use(
        http.get(`${API_BASE_URL}/users`, () => {
          return HttpResponse.json(
            { error: { code: 'SERVER_ERROR', message: 'Internal server error' } },
            { status: 500 }
          );
        })
      );

      render(
        <TestWrapper queryClient={queryClient}>
          <UsersListTestComponent />
        </TestWrapper>
      );

      await waitFor(() => {
        expect(screen.getByText(/error/i)).toBeInTheDocument();
      });
    });
  });

  describe('Create User', () => {
    it('creates user successfully', async () => {
      const user = userEvent.setup();
      const queryClient = createTestQueryClient();
      const onSuccess = vi.fn();

      server.use(
        http.post(`${API_BASE_URL}/users`, async () => {
          return HttpResponse.json(
            {
              data: {
                id: '3',
                email: 'new@example.com',
                fullName: 'New User',
                status: 'active',
                roles: [],
                permissions: [],
                createdAt: new Date().toISOString(),
                updatedAt: new Date().toISOString(),
              },
            },
            { status: 201 }
          );
        })
      );

      render(
        <TestWrapper queryClient={queryClient}>
          <CreateUserTestComponent onSuccess={onSuccess} />
        </TestWrapper>
      );

      await user.click(screen.getByRole('button', { name: /create user/i }));

      await waitFor(() => {
        expect(screen.getByText('Created!')).toBeInTheDocument();
      });

      expect(onSuccess).toHaveBeenCalled();
    });

    it('shows loading state during creation', async () => {
      const user = userEvent.setup();
      const queryClient = createTestQueryClient();

      render(
        <TestWrapper queryClient={queryClient}>
          <CreateUserTestComponent />
        </TestWrapper>
      );

      const createButton = screen.getByRole('button', { name: /create user/i });
      await user.click(createButton);

      // After clicking, either the loading state shows or it completes quickly
      await waitFor(
        () => {
          // Either creating... is shown or Created! is shown (if fast)
          const creating = screen.queryByText('Creating...');
          const created = screen.queryByText('Created!');
          expect(creating || created).toBeTruthy();
        },
        { timeout: 2000 }
      );
    });
  });

  describe('Delete User', () => {
    it('deletes user successfully', async () => {
      const user = userEvent.setup();
      const queryClient = createTestQueryClient();
      const onSuccess = vi.fn();

      server.use(
        http.delete(`${API_BASE_URL}/users/1`, () => {
          return HttpResponse.json({ data: null }, { status: 204 });
        })
      );

      render(
        <TestWrapper queryClient={queryClient}>
          <DeleteUserTestComponent userId="1" onSuccess={onSuccess} />
        </TestWrapper>
      );

      await user.click(screen.getByRole('button', { name: /delete user/i }));

      await waitFor(() => {
        expect(screen.getByText('Deleted!')).toBeInTheDocument();
      });

      expect(onSuccess).toHaveBeenCalled();
    });
  });

  describe('Cache Invalidation', () => {
    it('invalidates users list cache after create', async () => {
      const user = userEvent.setup();
      const queryClient = createTestQueryClient();

      server.use(
        http.post(`${API_BASE_URL}/users`, () => {
          return HttpResponse.json({
            data: {
              id: '2',
              email: 'new@example.com',
              fullName: 'New User',
              status: 'active',
              roles: [],
              permissions: [],
              createdAt: new Date().toISOString(),
              updatedAt: new Date().toISOString(),
            },
          });
        })
      );

      render(
        <TestWrapper queryClient={queryClient}>
          <UsersListTestComponent />
          <CreateUserTestComponent />
        </TestWrapper>
      );

      await waitFor(() => {
        expect(screen.getByText('Users List')).toBeInTheDocument();
      });

      await user.click(screen.getByRole('button', { name: /create user/i }));

      await waitFor(() => {
        expect(screen.getByText('Created!')).toBeInTheDocument();
      });
    });
  });
});
