import { ROUTES } from '@/constants';
import { queryKeys } from '@/lib/api';
import { useAuthStore } from '@/stores';
import type {
  ChangePasswordRequest,
  ForgotPasswordRequest,
  LoginRequest,
  RegisterRequest,
  ResetPasswordRequest,
  UpdateProfileRequest,
} from '@/types';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useNavigate } from 'react-router';
import { toast } from 'sonner';
import { authApi } from './auth.api';

export function useLogin() {
  const navigate = useNavigate();
  const { setUser } = useAuthStore();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: LoginRequest) => authApi.login(data),
    onSuccess: (response) => {
      setUser(response.user);
      queryClient.setQueryData(queryKeys.auth.currentUser(), response.user);
      toast.success('Welcome back!');
      navigate(ROUTES.DASHBOARD);
    },
    onError: () => {
      toast.error('Invalid email or password');
    },
  });
}

export function useRegister() {
  const navigate = useNavigate();
  const { setUser } = useAuthStore();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: RegisterRequest) => authApi.register(data),
    onSuccess: (response) => {
      setUser(response.user);
      queryClient.setQueryData(queryKeys.auth.currentUser(), response.user);
      toast.success('Account created successfully!');
      navigate(ROUTES.DASHBOARD);
    },
    onError: () => {
      toast.error('Failed to create account. Please try again.');
    },
  });
}

export function useLogout() {
  const navigate = useNavigate();
  const { clearAuth } = useAuthStore();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => authApi.logout(),
    onSuccess: () => {
      clearAuth();
      queryClient.clear();
      toast.success('Logged out successfully');
      navigate(ROUTES.LOGIN);
    },
    onError: () => {
      clearAuth();
      queryClient.clear();
      navigate(ROUTES.LOGIN);
    },
  });
}

export function useForgotPassword() {
  return useMutation({
    mutationFn: (data: ForgotPasswordRequest) => authApi.forgotPassword(data),
    onSuccess: () => {
      toast.success('Password reset instructions sent to your email');
    },
    onError: () => {
      toast.error('Failed to send reset instructions. Please try again.');
    },
  });
}

export function useResetPassword() {
  const navigate = useNavigate();

  return useMutation({
    mutationFn: (data: ResetPasswordRequest) => authApi.resetPassword(data),
    onSuccess: () => {
      toast.success('Password reset successfully');
      navigate(ROUTES.LOGIN);
    },
    onError: () => {
      toast.error('Failed to reset password. The link may have expired.');
    },
  });
}

export function useChangePassword() {
  return useMutation({
    mutationFn: (data: ChangePasswordRequest) => authApi.changePassword(data),
    onSuccess: () => {
      toast.success('Password changed successfully');
    },
    onError: () => {
      toast.error('Failed to change password. Please check your current password.');
    },
  });
}

export function useCurrentUser() {
  const { setUser, setStatus, clearAuth } = useAuthStore();

  return useQuery({
    queryKey: queryKeys.auth.currentUser(),
    queryFn: async () => {
      const user = await authApi.getCurrentUser();
      setUser(user);
      return user;
    },
    retry: false,
    staleTime: 5 * 60 * 1000,
    meta: {
      onError: () => {
        clearAuth();
        setStatus('unauthenticated');
      },
    },
  });
}

export function useSessions() {
  return useQuery({
    queryKey: queryKeys.auth.sessions(),
    queryFn: () => authApi.getSessions(),
  });
}

export function useRevokeSession() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (sessionId: string) => authApi.revokeSession(sessionId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.auth.sessions() });
      toast.success('Session revoked successfully');
    },
    onError: () => {
      toast.error('Failed to revoke session');
    },
  });
}

export function useUpdateProfile() {
  const queryClient = useQueryClient();
  const { updateUser } = useAuthStore();

  return useMutation({
    mutationFn: (data: UpdateProfileRequest) => authApi.updateProfile(data),
    onSuccess: (user) => {
      updateUser(user);
      queryClient.setQueryData(queryKeys.auth.currentUser(), user);
      toast.success('Profile updated successfully');
    },
    onError: () => {
      toast.error('Failed to update profile');
    },
  });
}
