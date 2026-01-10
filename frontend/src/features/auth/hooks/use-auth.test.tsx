import * as authStore from '@/stores/auth-store';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { renderHook } from '@testing-library/react';
import type { ReactNode } from 'react';
import { MemoryRouter } from 'react-router';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import * as authQueries from '../api/auth.queries';
import { useAuth } from './use-auth';

vi.mock('@/stores/auth-store', () => ({
  useAuthStore: vi.fn(),
}));

vi.mock('../api/auth.queries', () => ({
  useCurrentUser: vi.fn(),
  useLogout: vi.fn(),
}));

const mockUser = {
  id: '123',
  email: 'test@example.com',
  fullName: 'Test User',
  status: 'active' as const,
  roles: [{ id: '1', name: 'admin', displayName: 'Administrator' }],
  permissions: ['users:read', 'users:write'],
  createdAt: '2024-01-01T00:00:00Z',
  updatedAt: '2024-01-01T00:00:00Z',
};

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
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

describe('useAuth', () => {
  const mockHasPermission = vi.fn();
  const mockHasRole = vi.fn();
  const mockHasAnyPermission = vi.fn();
  const mockHasAllPermissions = vi.fn();
  const mockLogout = vi.fn();
  const mockRefetchUser = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();

    vi.mocked(authStore.useAuthStore).mockReturnValue({
      user: null,
      status: 'unauthenticated',
      hasPermission: mockHasPermission,
      hasRole: mockHasRole,
      hasAnyPermission: mockHasAnyPermission,
      hasAllPermissions: mockHasAllPermissions,
    } as unknown as ReturnType<typeof authStore.useAuthStore>);

    vi.mocked(authQueries.useCurrentUser).mockReturnValue({
      isLoading: false,
      refetch: mockRefetchUser,
    } as unknown as ReturnType<typeof authQueries.useCurrentUser>);

    vi.mocked(authQueries.useLogout).mockReturnValue({
      mutate: mockLogout,
      isPending: false,
    } as unknown as ReturnType<typeof authQueries.useLogout>);
  });

  it('returns unauthenticated state when no user', () => {
    const { result } = renderHook(() => useAuth(), {
      wrapper: createWrapper(),
    });

    expect(result.current.user).toBeNull();
    expect(result.current.isAuthenticated).toBe(false);
    expect(result.current.status).toBe('unauthenticated');
  });

  it('returns authenticated state when user exists', () => {
    vi.mocked(authStore.useAuthStore).mockReturnValue({
      user: mockUser,
      status: 'authenticated',
      hasPermission: mockHasPermission,
      hasRole: mockHasRole,
      hasAnyPermission: mockHasAnyPermission,
      hasAllPermissions: mockHasAllPermissions,
    } as unknown as ReturnType<typeof authStore.useAuthStore>);

    const { result } = renderHook(() => useAuth(), {
      wrapper: createWrapper(),
    });

    expect(result.current.user).toEqual(mockUser);
    expect(result.current.isAuthenticated).toBe(true);
    expect(result.current.status).toBe('authenticated');
  });

  it('returns loading state when status is loading', () => {
    vi.mocked(authStore.useAuthStore).mockReturnValue({
      user: null,
      status: 'loading',
      hasPermission: mockHasPermission,
      hasRole: mockHasRole,
      hasAnyPermission: mockHasAnyPermission,
      hasAllPermissions: mockHasAllPermissions,
    } as unknown as ReturnType<typeof authStore.useAuthStore>);

    const { result } = renderHook(() => useAuth(), {
      wrapper: createWrapper(),
    });

    expect(result.current.isLoading).toBe(true);
  });

  it('returns loading state when useCurrentUser is loading', () => {
    vi.mocked(authQueries.useCurrentUser).mockReturnValue({
      isLoading: true,
      refetch: mockRefetchUser,
    } as unknown as ReturnType<typeof authQueries.useCurrentUser>);

    const { result } = renderHook(() => useAuth(), {
      wrapper: createWrapper(),
    });

    expect(result.current.isLoading).toBe(true);
  });

  it('exposes permission check functions', () => {
    const { result } = renderHook(() => useAuth(), {
      wrapper: createWrapper(),
    });

    expect(result.current.hasPermission).toBe(mockHasPermission);
    expect(result.current.hasRole).toBe(mockHasRole);
    expect(result.current.hasAnyPermission).toBe(mockHasAnyPermission);
    expect(result.current.hasAllPermissions).toBe(mockHasAllPermissions);
  });

  it('exposes logout function', () => {
    const { result } = renderHook(() => useAuth(), {
      wrapper: createWrapper(),
    });

    result.current.logout();

    expect(mockLogout).toHaveBeenCalled();
  });

  it('exposes refetchUser function', () => {
    const { result } = renderHook(() => useAuth(), {
      wrapper: createWrapper(),
    });

    result.current.refetchUser();

    expect(mockRefetchUser).toHaveBeenCalled();
  });

  it('returns isLoggingOut state', () => {
    vi.mocked(authQueries.useLogout).mockReturnValue({
      mutate: mockLogout,
      isPending: true,
    } as unknown as ReturnType<typeof authQueries.useLogout>);

    const { result } = renderHook(() => useAuth(), {
      wrapper: createWrapper(),
    });

    expect(result.current.isLoggingOut).toBe(true);
  });
});
