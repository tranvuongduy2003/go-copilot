import { API_BASE_URL, ERROR_CODES, HTTP_STATUS } from '@/constants';
import axios, {
  type AxiosError,
  type AxiosInstance,
  type AxiosRequestConfig,
  type AxiosResponse,
  type InternalAxiosRequestConfig,
} from 'axios';
import { ApiError, type ApiErrorResponse } from './errors';
import { tokenService } from './token-service';

interface RefreshSubscriber {
  resolve: (token: string) => void;
  reject: (error: Error) => void;
}

let isRefreshing = false;
let refreshSubscribers: RefreshSubscriber[] = [];

const subscribeTokenRefresh = (callback: RefreshSubscriber): void => {
  refreshSubscribers.push(callback);
};

const onTokenRefreshed = (token: string): void => {
  refreshSubscribers.forEach((subscriber) => subscriber.resolve(token));
  refreshSubscribers = [];
};

const onTokenRefreshFailed = (error: Error): void => {
  refreshSubscribers.forEach((subscriber) => subscriber.reject(error));
  refreshSubscribers = [];
};

const PUBLIC_ENDPOINTS = [
  '/auth/login',
  '/auth/register',
  '/auth/forgot-password',
  '/auth/reset-password',
  '/auth/refresh',
];

const isPublicEndpoint = (url: string): boolean => {
  return PUBLIC_ENDPOINTS.some((endpoint) => url.includes(endpoint));
};

const createApiClient = (): AxiosInstance => {
  const instance = axios.create({
    baseURL: API_BASE_URL,
    timeout: 30000,
    headers: {
      'Content-Type': 'application/json',
      Accept: 'application/json',
    },
  });

  instance.interceptors.request.use(
    (config: InternalAxiosRequestConfig): InternalAxiosRequestConfig => {
      const url = config.url || '';

      if (!isPublicEndpoint(url)) {
        const token = tokenService.getAccessToken();
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }
      }

      return config;
    },
    (error: AxiosError) => {
      return Promise.reject(error);
    }
  );

  instance.interceptors.response.use(
    (response: AxiosResponse) => {
      return response;
    },
    async (error: AxiosError<ApiErrorResponse>) => {
      const originalRequest = error.config as InternalAxiosRequestConfig & {
        _retry?: boolean;
      };

      if (!originalRequest) {
        return Promise.reject(createApiError(error));
      }

      if (
        error.response?.status === HTTP_STATUS.UNAUTHORIZED &&
        !originalRequest._retry &&
        !isPublicEndpoint(originalRequest.url || '')
      ) {
        if (isRefreshing) {
          return new Promise((resolve, reject) => {
            subscribeTokenRefresh({
              resolve: (token: string) => {
                originalRequest.headers.Authorization = `Bearer ${token}`;
                resolve(instance(originalRequest));
              },
              reject: (refreshError: Error) => {
                reject(refreshError);
              },
            });
          });
        }

        originalRequest._retry = true;
        isRefreshing = true;

        try {
          const refreshToken = tokenService.getRefreshToken();
          if (!refreshToken) {
            throw new Error('No refresh token available');
          }

          const response = await axios.post<{
            data: { accessToken: string; refreshToken: string };
          }>(`${API_BASE_URL}/auth/refresh`, {
            refreshToken,
          });

          const { accessToken, refreshToken: newRefreshToken } = response.data.data;
          tokenService.setTokens(accessToken, newRefreshToken);

          onTokenRefreshed(accessToken);
          originalRequest.headers.Authorization = `Bearer ${accessToken}`;

          return instance(originalRequest);
        } catch (refreshError) {
          onTokenRefreshFailed(refreshError as Error);
          tokenService.clearTokens();
          window.location.href = '/login';
          return Promise.reject(refreshError);
        } finally {
          isRefreshing = false;
        }
      }

      return Promise.reject(createApiError(error));
    }
  );

  return instance;
};

const createApiError = (error: AxiosError<ApiErrorResponse>): ApiError => {
  if (error.code === 'ECONNABORTED') {
    return new ApiError(
      'Request timed out. Please try again.',
      0,
      ERROR_CODES.TIMEOUT_ERROR,
      undefined,
      error
    );
  }

  if (!error.response) {
    return new ApiError(
      'Network error. Please check your connection.',
      0,
      ERROR_CODES.NETWORK_ERROR,
      undefined,
      error
    );
  }

  const { status, data } = error.response;
  const message = data?.error?.message || getDefaultErrorMessage(status);
  const code = data?.error?.code || getErrorCode(status);
  const details = data?.error?.details;

  return new ApiError(message, status, code, details, error);
};

const getDefaultErrorMessage = (status: number): string => {
  switch (status) {
    case HTTP_STATUS.BAD_REQUEST:
      return 'Invalid request. Please check your input.';
    case HTTP_STATUS.UNAUTHORIZED:
      return 'Your session has expired. Please log in again.';
    case HTTP_STATUS.FORBIDDEN:
      return 'You do not have permission to perform this action.';
    case HTTP_STATUS.NOT_FOUND:
      return 'The requested resource was not found.';
    case HTTP_STATUS.CONFLICT:
      return 'A conflict occurred. The resource may already exist.';
    case HTTP_STATUS.UNPROCESSABLE_ENTITY:
      return 'Invalid data provided. Please check your input.';
    case HTTP_STATUS.TOO_MANY_REQUESTS:
      return 'Too many requests. Please try again later.';
    case HTTP_STATUS.INTERNAL_SERVER_ERROR:
      return 'An unexpected error occurred. Please try again.';
    case HTTP_STATUS.SERVICE_UNAVAILABLE:
      return 'Service temporarily unavailable. Please try again later.';
    default:
      return 'An error occurred. Please try again.';
  }
};

const getErrorCode = (status: number): string => {
  switch (status) {
    case HTTP_STATUS.BAD_REQUEST:
    case HTTP_STATUS.UNPROCESSABLE_ENTITY:
      return ERROR_CODES.VALIDATION_ERROR;
    case HTTP_STATUS.UNAUTHORIZED:
      return ERROR_CODES.UNAUTHORIZED;
    case HTTP_STATUS.FORBIDDEN:
      return ERROR_CODES.FORBIDDEN;
    case HTTP_STATUS.NOT_FOUND:
      return ERROR_CODES.NOT_FOUND;
    case HTTP_STATUS.CONFLICT:
      return ERROR_CODES.CONFLICT;
    case HTTP_STATUS.INTERNAL_SERVER_ERROR:
    case HTTP_STATUS.SERVICE_UNAVAILABLE:
      return ERROR_CODES.SERVER_ERROR;
    default:
      return ERROR_CODES.UNKNOWN_ERROR;
  }
};

export const apiClient = createApiClient();

export const api = {
  get: <TResponse>(url: string, config?: AxiosRequestConfig) =>
    apiClient.get<TResponse>(url, config).then((response) => response.data),

  post: <TResponse, TData = unknown>(url: string, data?: TData, config?: AxiosRequestConfig) =>
    apiClient.post<TResponse>(url, data, config).then((response) => response.data),

  put: <TResponse, TData = unknown>(url: string, data?: TData, config?: AxiosRequestConfig) =>
    apiClient.put<TResponse>(url, data, config).then((response) => response.data),

  patch: <TResponse, TData = unknown>(url: string, data?: TData, config?: AxiosRequestConfig) =>
    apiClient.patch<TResponse>(url, data, config).then((response) => response.data),

  delete: <TResponse>(url: string, config?: AxiosRequestConfig) =>
    apiClient.delete<TResponse>(url, config).then((response) => response.data),
};
