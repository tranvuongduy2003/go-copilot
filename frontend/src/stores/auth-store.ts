import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';
import { immer } from 'zustand/middleware/immer';
import type { User, AuthStatus } from '@/types';
import { AUTH_STORAGE_KEY } from '@/constants';
import { tokenService } from '@/lib/api';

interface AuthState {
  user: User | null;
  status: AuthStatus;
}

interface AuthActions {
  setUser: (user: User) => void;
  clearAuth: () => void;
  updateUser: (updates: Partial<User>) => void;
  setStatus: (status: AuthStatus) => void;
  hasPermission: (permission: string) => boolean;
  hasRole: (roleName: string) => boolean;
  hasAnyPermission: (permissions: string[]) => boolean;
  hasAllPermissions: (permissions: string[]) => boolean;
}

type AuthStore = AuthState & AuthActions;

const initialState: AuthState = {
  user: null,
  status: 'loading',
};

export const useAuthStore = create<AuthStore>()(
  persist(
    immer((set, get) => ({
      ...initialState,

      setUser: (user: User) => {
        set((state) => {
          state.user = user;
          state.status = 'authenticated';
        });
      },

      clearAuth: () => {
        tokenService.clearTokens();
        set((state) => {
          state.user = null;
          state.status = 'unauthenticated';
        });
      },

      updateUser: (updates: Partial<User>) => {
        set((state) => {
          if (state.user) {
            state.user = { ...state.user, ...updates };
          }
        });
      },

      setStatus: (status: AuthStatus) => {
        set((state) => {
          state.status = status;
        });
      },

      hasPermission: (permission: string): boolean => {
        const { user } = get();
        if (!user) return false;
        return user.permissions.includes(permission);
      },

      hasRole: (roleName: string): boolean => {
        const { user } = get();
        if (!user) return false;
        return user.roles.some((role) => role.name === roleName);
      },

      hasAnyPermission: (permissions: string[]): boolean => {
        const { user } = get();
        if (!user) return false;
        return permissions.some((permission) => user.permissions.includes(permission));
      },

      hasAllPermissions: (permissions: string[]): boolean => {
        const { user } = get();
        if (!user) return false;
        return permissions.every((permission) => user.permissions.includes(permission));
      },
    })),
    {
      name: `${AUTH_STORAGE_KEY}_user`,
      storage: createJSONStorage(() => localStorage),
      partialize: (state) => ({
        user: state.user,
        status: state.user ? 'authenticated' : 'unauthenticated',
      }),
    }
  )
);

export const selectUser = (state: AuthStore) => state.user;
export const selectAuthStatus = (state: AuthStore) => state.status;
export const selectIsAuthenticated = (state: AuthStore) => state.status === 'authenticated';
export const selectIsLoading = (state: AuthStore) => state.status === 'loading';
