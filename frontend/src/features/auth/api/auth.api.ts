import { API_ENDPOINTS } from '@/constants';
import { api, tokenService } from '@/lib/api';
import type {
  ApiResponse,
  ChangePasswordRequest,
  ForgotPasswordRequest,
  LoginRequest,
  LoginResponse,
  RefreshTokenResponse,
  RegisterRequest,
  RegisterResponse,
  ResetPasswordRequest,
  Session,
  UpdateProfileRequest,
  User,
} from '@/types';

export const authApi = {
  login: async (data: LoginRequest): Promise<LoginResponse> => {
    const response = await api.post<ApiResponse<LoginResponse>>(API_ENDPOINTS.AUTH.LOGIN, data);
    const { accessToken, refreshToken } = response.data;
    tokenService.setTokens(accessToken, refreshToken);
    return response.data;
  },

  register: async (data: RegisterRequest): Promise<RegisterResponse> => {
    const response = await api.post<ApiResponse<RegisterResponse>>(
      API_ENDPOINTS.AUTH.REGISTER,
      data
    );
    const { accessToken, refreshToken } = response.data;
    tokenService.setTokens(accessToken, refreshToken);
    return response.data;
  },

  logout: async (): Promise<void> => {
    try {
      await api.post(API_ENDPOINTS.AUTH.LOGOUT);
    } finally {
      tokenService.clearTokens();
    }
  },

  refreshToken: async (refreshToken: string): Promise<RefreshTokenResponse> => {
    const response = await api.post<ApiResponse<RefreshTokenResponse>>(API_ENDPOINTS.AUTH.REFRESH, {
      refreshToken,
    });
    const { accessToken, refreshToken: newRefreshToken } = response.data;
    tokenService.setTokens(accessToken, newRefreshToken);
    return response.data;
  },

  getCurrentUser: async (): Promise<User> => {
    const response = await api.get<ApiResponse<User>>(API_ENDPOINTS.AUTH.ME);
    return response.data;
  },

  updateProfile: async (data: UpdateProfileRequest): Promise<User> => {
    const response = await api.patch<ApiResponse<User>>(API_ENDPOINTS.AUTH.ME, data);
    return response.data;
  },

  forgotPassword: async (data: ForgotPasswordRequest): Promise<void> => {
    await api.post(API_ENDPOINTS.AUTH.FORGOT_PASSWORD, data);
  },

  resetPassword: async (data: ResetPasswordRequest): Promise<void> => {
    await api.post(API_ENDPOINTS.AUTH.RESET_PASSWORD, data);
  },

  changePassword: async (data: ChangePasswordRequest): Promise<void> => {
    await api.post(API_ENDPOINTS.AUTH.CHANGE_PASSWORD, data);
  },

  getSessions: async (): Promise<Session[]> => {
    const response = await api.get<ApiResponse<Session[]>>(API_ENDPOINTS.AUTH.SESSIONS);
    return response.data;
  },

  revokeSession: async (sessionId: string): Promise<void> => {
    const url = API_ENDPOINTS.AUTH.REVOKE_SESSION.replace(':id', sessionId);
    await api.post(url);
  },
};
