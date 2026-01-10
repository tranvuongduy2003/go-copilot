import { AuthLayout, MainLayout } from '@/components/layout';
import { ROUTES } from '@/constants';
import { AuthGuard, GuestGuard } from '@/features/auth/components';
import { lazy, Suspense } from 'react';
import { createBrowserRouter, Navigate, type RouteObject } from 'react-router-dom';
import { PageLoader } from './page-loader';

const DashboardPage = lazy(() =>
  import('@/pages/dashboard-page').then((module) => ({ default: module.DashboardPage }))
);
const NotFoundPage = lazy(() =>
  import('@/pages/not-found-page').then((module) => ({ default: module.NotFoundPage }))
);
const ProfilePage = lazy(() =>
  import('@/pages/profile-page').then((module) => ({ default: module.ProfilePage }))
);
const SettingsPage = lazy(() =>
  import('@/pages/settings-page').then((module) => ({ default: module.SettingsPage }))
);
const UnauthorizedPage = lazy(() =>
  import('@/pages/unauthorized-page').then((module) => ({ default: module.UnauthorizedPage }))
);
const ServerErrorPage = lazy(() =>
  import('@/pages/server-error-page').then((module) => ({ default: module.ServerErrorPage }))
);

const LoginPage = lazy(() =>
  import('@/features/auth/pages/login.page').then((module) => ({ default: module.LoginPage }))
);
const RegisterPage = lazy(() =>
  import('@/features/auth/pages/register.page').then((module) => ({ default: module.RegisterPage }))
);
const ForgotPasswordPage = lazy(() =>
  import('@/features/auth/pages/forgot-password.page').then((module) => ({
    default: module.ForgotPasswordPage,
  }))
);
const ResetPasswordPage = lazy(() =>
  import('@/features/auth/pages/reset-password.page').then((module) => ({
    default: module.ResetPasswordPage,
  }))
);

const UsersPage = lazy(() =>
  import('@/features/users/pages/users-page').then((module) => ({ default: module.UsersPage }))
);
const UserDetailPage = lazy(() =>
  import('@/features/users/pages/user-detail-page').then((module) => ({
    default: module.UserDetailPage,
  }))
);

const RolesPage = lazy(() =>
  import('@/features/roles/pages/roles-page').then((module) => ({ default: module.RolesPage }))
);

function withSuspense(Component: React.ComponentType) {
  return (
    <Suspense fallback={<PageLoader />}>
      <Component />
    </Suspense>
  );
}

const authRoutes: RouteObject[] = [
  {
    element: (
      <GuestGuard>
        <AuthLayout />
      </GuestGuard>
    ),
    children: [
      { path: ROUTES.LOGIN, element: withSuspense(LoginPage) },
      { path: ROUTES.REGISTER, element: withSuspense(RegisterPage) },
      { path: ROUTES.FORGOT_PASSWORD, element: withSuspense(ForgotPasswordPage) },
      { path: ROUTES.RESET_PASSWORD, element: withSuspense(ResetPasswordPage) },
    ],
  },
];

const protectedRoutes: RouteObject[] = [
  {
    element: (
      <AuthGuard>
        <MainLayout />
      </AuthGuard>
    ),
    children: [
      { index: true, element: <Navigate to={ROUTES.DASHBOARD} replace /> },
      { path: ROUTES.DASHBOARD, element: withSuspense(DashboardPage) },
      { path: ROUTES.USERS, element: withSuspense(UsersPage) },
      { path: ROUTES.USER_DETAIL, element: withSuspense(UserDetailPage) },
      { path: ROUTES.ROLES, element: withSuspense(RolesPage) },
      { path: ROUTES.PROFILE, element: withSuspense(ProfilePage) },
      { path: ROUTES.SETTINGS, element: withSuspense(SettingsPage) },
    ],
  },
];

const publicRoutes: RouteObject[] = [
  { path: ROUTES.UNAUTHORIZED, element: withSuspense(UnauthorizedPage) },
  { path: '/500', element: withSuspense(ServerErrorPage) },
  { path: '*', element: withSuspense(NotFoundPage) },
];

export const router = createBrowserRouter([...authRoutes, ...protectedRoutes, ...publicRoutes]);
