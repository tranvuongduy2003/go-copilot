import { Loading } from '@/components/ui/spinner';
import { ROUTES } from '@/constants';
import type { ReactNode } from 'react';
import { Navigate, useLocation } from 'react-router';
import { useAuth } from '../hooks';

interface AuthGuardProps {
  children: ReactNode;
  requiredPermissions?: string[];
  requiredRoles?: string[];
  requireAll?: boolean;
}

export function AuthGuard({
  children,
  requiredPermissions = [],
  requiredRoles = [],
  requireAll = true,
}: AuthGuardProps) {
  const location = useLocation();
  const { isAuthenticated, isLoading, hasRole, hasAnyPermission, hasAllPermissions } = useAuth();

  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <Loading text="Checking authentication..." />
      </div>
    );
  }

  if (!isAuthenticated) {
    return <Navigate to={ROUTES.LOGIN} state={{ from: location }} replace />;
  }

  if (requiredPermissions.length > 0) {
    const hasRequiredPermissions = requireAll
      ? hasAllPermissions(requiredPermissions)
      : hasAnyPermission(requiredPermissions);

    if (!hasRequiredPermissions) {
      return <Navigate to={ROUTES.UNAUTHORIZED} replace />;
    }
  }

  if (requiredRoles.length > 0) {
    const hasRequiredRoles = requireAll
      ? requiredRoles.every((role) => hasRole(role))
      : requiredRoles.some((role) => hasRole(role));

    if (!hasRequiredRoles) {
      return <Navigate to={ROUTES.UNAUTHORIZED} replace />;
    }
  }

  return <>{children}</>;
}

interface GuestGuardProps {
  children: ReactNode;
}

export function GuestGuard({ children }: GuestGuardProps) {
  const location = useLocation();
  const { isAuthenticated, isLoading } = useAuth();
  const from = (location.state as { from?: Location })?.from?.pathname || ROUTES.DASHBOARD;

  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <Loading text="Loading..." />
      </div>
    );
  }

  if (isAuthenticated) {
    return <Navigate to={from} replace />;
  }

  return <>{children}</>;
}
