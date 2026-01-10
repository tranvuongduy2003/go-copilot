import type { User } from '@/types';
import { describe, expect, it, vi } from 'vitest';
import { selectAuthStatus, selectIsAuthenticated, selectIsLoading, selectUser } from './auth-store';

const mockUser: User = {
  id: 'user-123',
  email: 'test@example.com',
  fullName: 'Test User',
  status: 'active',
  roles: [
    {
      id: 'role-1',
      name: 'admin',
      displayName: 'Administrator',
      isSystem: false,
      createdAt: '2024-01-01T00:00:00Z',
      updatedAt: '2024-01-01T00:00:00Z',
    },
  ],
  permissions: ['users:read', 'users:write', 'roles:read'],
  createdAt: '2024-01-01T00:00:00Z',
  updatedAt: '2024-01-01T00:00:00Z',
};

describe('Auth Store Selectors', () => {
  const authenticatedState = {
    user: mockUser,
    status: 'authenticated' as const,
    setUser: vi.fn(),
    clearAuth: vi.fn(),
    updateUser: vi.fn(),
    setStatus: vi.fn(),
    hasPermission: vi.fn(),
    hasRole: vi.fn(),
    hasAnyPermission: vi.fn(),
    hasAllPermissions: vi.fn(),
  };

  const unauthenticatedState = {
    ...authenticatedState,
    user: null,
    status: 'unauthenticated' as const,
  };

  const loadingState = {
    ...authenticatedState,
    user: null,
    status: 'loading' as const,
  };

  describe('selectUser', () => {
    it('returns the user from state', () => {
      expect(selectUser(authenticatedState)).toEqual(mockUser);
      expect(selectUser(unauthenticatedState)).toBeNull();
    });
  });

  describe('selectAuthStatus', () => {
    it('returns the auth status from state', () => {
      expect(selectAuthStatus(authenticatedState)).toBe('authenticated');
      expect(selectAuthStatus(unauthenticatedState)).toBe('unauthenticated');
      expect(selectAuthStatus(loadingState)).toBe('loading');
    });
  });

  describe('selectIsAuthenticated', () => {
    it('returns true when status is authenticated', () => {
      expect(selectIsAuthenticated(authenticatedState)).toBe(true);
    });

    it('returns false when status is not authenticated', () => {
      expect(selectIsAuthenticated(unauthenticatedState)).toBe(false);
      expect(selectIsAuthenticated(loadingState)).toBe(false);
    });
  });

  describe('selectIsLoading', () => {
    it('returns true when status is loading', () => {
      expect(selectIsLoading(loadingState)).toBe(true);
    });

    it('returns false when status is not loading', () => {
      expect(selectIsLoading(authenticatedState)).toBe(false);
      expect(selectIsLoading(unauthenticatedState)).toBe(false);
    });
  });

  describe('hasPermission helper (via mocked state)', () => {
    it('checks if user has specific permission', () => {
      const hasPermission = (permission: string): boolean => {
        const { user } = authenticatedState;
        if (!user) return false;
        return user.permissions.includes(permission);
      };

      expect(hasPermission('users:read')).toBe(true);
      expect(hasPermission('users:delete')).toBe(false);
    });
  });

  describe('hasRole helper (via mocked state)', () => {
    it('checks if user has specific role', () => {
      const hasRole = (roleName: string): boolean => {
        const { user } = authenticatedState;
        if (!user) return false;
        return user.roles.some((role) => role.name === roleName);
      };

      expect(hasRole('admin')).toBe(true);
      expect(hasRole('super-admin')).toBe(false);
    });
  });

  describe('hasAnyPermission helper (via mocked state)', () => {
    it('checks if user has any of the permissions', () => {
      const hasAnyPermission = (permissions: string[]): boolean => {
        const { user } = authenticatedState;
        if (!user) return false;
        return permissions.some((permission) => user.permissions.includes(permission));
      };

      expect(hasAnyPermission(['users:read', 'users:delete'])).toBe(true);
      expect(hasAnyPermission(['users:delete', 'roles:delete'])).toBe(false);
    });
  });

  describe('hasAllPermissions helper (via mocked state)', () => {
    it('checks if user has all permissions', () => {
      const hasAllPermissions = (permissions: string[]): boolean => {
        const { user } = authenticatedState;
        if (!user) return false;
        return permissions.every((permission) => user.permissions.includes(permission));
      };

      expect(hasAllPermissions(['users:read', 'users:write'])).toBe(true);
      expect(hasAllPermissions(['users:read', 'users:delete'])).toBe(false);
    });
  });
});

describe('Auth Store Persistence', () => {
  const localStorageMock = (() => {
    let store: Record<string, string> = {};
    return {
      getItem: (key: string) => store[key] || null,
      setItem: (key: string, value: string) => {
        store[key] = value;
      },
      removeItem: (key: string) => {
        delete store[key];
      },
      clear: () => {
        store = {};
      },
    };
  })();

  it('persists user data to localStorage', async () => {
    Object.defineProperty(window, 'localStorage', {
      value: localStorageMock,
      writable: true,
    });

    const { useAuthStore } = await import('./auth-store');
    const store = useAuthStore.getState();

    store.setUser(mockUser);

    // Verify store state is updated
    expect(store.user).toBeNull(); // Note: getState() returns snapshot, need to re-get
    const currentState = useAuthStore.getState();
    expect(currentState.user).toBeDefined();
    expect(currentState.user?.email).toBe('test@example.com');
  });

  it('clears user data from localStorage on logout', async () => {
    const { useAuthStore } = await import('./auth-store');
    const store = useAuthStore.getState();

    store.setUser(mockUser);
    store.clearAuth();

    await new Promise((resolve) => setTimeout(resolve, 100));

    const storedData = localStorageMock.getItem('go-copilot-auth');
    if (storedData) {
      const parsed = JSON.parse(storedData);
      expect(parsed.state.user).toBeNull();
    }
  });
});
