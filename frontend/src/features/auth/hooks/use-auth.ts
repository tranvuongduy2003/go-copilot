import { useAuthStore } from '@/stores';
import { useCurrentUser, useLogout } from '../api';

export function useAuth() {
  const { user, status, hasPermission, hasRole, hasAnyPermission, hasAllPermissions } =
    useAuthStore();

  const { isLoading: isLoadingUser, refetch: refetchUser } = useCurrentUser();
  const { mutate: logout, isPending: isLoggingOut } = useLogout();

  const isAuthenticated = status === 'authenticated';
  const isLoading = status === 'loading' || isLoadingUser;

  return {
    user,
    status,
    isAuthenticated,
    isLoading,
    isLoggingOut,
    hasPermission,
    hasRole,
    hasAnyPermission,
    hasAllPermissions,
    logout,
    refetchUser,
  };
}
