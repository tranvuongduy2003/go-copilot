import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { renderHook, waitFor } from '@testing-library/react';
import type { ReactNode } from 'react';
import { MemoryRouter } from 'react-router';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { useUsers, useUser, useCreateUser, useDeleteUser } from './users.queries';

vi.mock('./users.api', () => ({
  usersApi: {
    getUsers: vi.fn(),
    getUser: vi.fn(),
    createUser: vi.fn(),
    deleteUser: vi.fn(),
  },
}));

import { usersApi } from './users.api';

const mockedUsersApi = vi.mocked(usersApi);

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
        gcTime: 0,
        staleTime: 0,
      },
      mutations: {
        retry: false,
      },
    },
  });

  return function Wrapper({ children }: { children: ReactNode }) {
    return (
      <QueryClientProvider client={queryClient}>
        <MemoryRouter>{children}</MemoryRouter>
      </QueryClientProvider>
    );
  };
}

describe('useUsers', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('fetches users list successfully', async () => {
    const mockUsersResponse = {
      data: [
        {
          id: '1',
          email: 'user1@test.com',
          fullName: 'User One',
          status: 'active' as const,
          roles: [],
          permissions: [],
          createdAt: '',
          updatedAt: '',
        },
        {
          id: '2',
          email: 'user2@test.com',
          fullName: 'User Two',
          status: 'active' as const,
          roles: [],
          permissions: [],
          createdAt: '',
          updatedAt: '',
        },
      ],
      meta: { page: 1, pageSize: 10, total: 2, totalPages: 1 },
    };

    mockedUsersApi.getUsers.mockResolvedValue(mockUsersResponse);

    const { result } = renderHook(() => useUsers(), { wrapper: createWrapper() });

    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true);
    });

    expect(result.current.data).toEqual(mockUsersResponse);
    expect(mockedUsersApi.getUsers).toHaveBeenCalledWith(undefined);
  });

  it('fetches users with filters', async () => {
    const filters = { status: 'active' as const, page: 1, pageSize: 10 };
    const mockUsersResponse = {
      data: [
        {
          id: '1',
          email: 'user1@test.com',
          fullName: 'User One',
          status: 'active' as const,
          roles: [],
          permissions: [],
          createdAt: '',
          updatedAt: '',
        },
      ],
      meta: { page: 1, pageSize: 10, total: 1, totalPages: 1 },
    };

    mockedUsersApi.getUsers.mockResolvedValue(mockUsersResponse);

    const { result } = renderHook(() => useUsers(filters), { wrapper: createWrapper() });

    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true);
    });

    expect(mockedUsersApi.getUsers).toHaveBeenCalledWith(filters);
  });

  it('handles error state', async () => {
    mockedUsersApi.getUsers.mockRejectedValue(new Error('API Error'));

    const { result } = renderHook(() => useUsers(), { wrapper: createWrapper() });

    await waitFor(() => {
      expect(result.current.isError).toBe(true);
    });

    expect(result.current.error).toBeInstanceOf(Error);
  });
});

describe('useUser', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('fetches single user by ID', async () => {
    const mockUser = {
      id: '1',
      email: 'user@test.com',
      fullName: 'Test User',
      status: 'active' as const,
      roles: [],
      permissions: [],
      createdAt: '',
      updatedAt: '',
    };

    mockedUsersApi.getUser.mockResolvedValue(mockUser);

    const { result } = renderHook(() => useUser('1'), { wrapper: createWrapper() });

    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true);
    });

    expect(result.current.data).toEqual(mockUser);
    expect(mockedUsersApi.getUser).toHaveBeenCalledWith('1');
  });

  it('does not fetch when ID is empty', () => {
    const { result } = renderHook(() => useUser(''), { wrapper: createWrapper() });

    expect(result.current.isFetching).toBe(false);
    expect(mockedUsersApi.getUser).not.toHaveBeenCalled();
  });
});

describe('useCreateUser', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('creates user successfully', async () => {
    const newUser = {
      id: '3',
      email: 'new@test.com',
      fullName: 'New User',
      status: 'active' as const,
      roles: [],
      permissions: [],
      createdAt: '',
      updatedAt: '',
    };

    mockedUsersApi.createUser.mockResolvedValue(newUser);

    const { result } = renderHook(() => useCreateUser(), { wrapper: createWrapper() });

    result.current.mutate({
      email: 'new@test.com',
      fullName: 'New User',
      password: 'password123',
    });

    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true);
    });

    expect(mockedUsersApi.createUser).toHaveBeenCalledWith({
      email: 'new@test.com',
      fullName: 'New User',
      password: 'password123',
    });
  });
});

describe('useDeleteUser', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('deletes user successfully', async () => {
    mockedUsersApi.deleteUser.mockResolvedValue(undefined);

    const { result } = renderHook(() => useDeleteUser(), { wrapper: createWrapper() });

    result.current.mutate('1');

    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true);
    });

    expect(mockedUsersApi.deleteUser).toHaveBeenCalledWith('1');
  });
});
