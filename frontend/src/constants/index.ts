export const APP_NAME = import.meta.env.VITE_APP_NAME || 'Go Copilot';
export const APP_VERSION = import.meta.env.VITE_APP_VERSION || '1.0.0';
export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api/v1';
export const AUTH_STORAGE_KEY = import.meta.env.VITE_AUTH_STORAGE_KEY || 'auth_storage';
export const TOKEN_REFRESH_THRESHOLD = Number(import.meta.env.VITE_TOKEN_REFRESH_THRESHOLD) || 300;
export const ENABLE_MOCK = import.meta.env.VITE_ENABLE_MOCK === 'true';

export const QUERY_STALE_TIME = 5 * 60 * 1000;
export const QUERY_CACHE_TIME = 10 * 60 * 1000;

export const ROUTES = {
  HOME: '/',
  LOGIN: '/login',
  REGISTER: '/register',
  FORGOT_PASSWORD: '/forgot-password',
  RESET_PASSWORD: '/reset-password',
  DASHBOARD: '/dashboard',
  PROFILE: '/profile',
  SETTINGS: '/settings',
  USERS: '/users',
  USER_DETAIL: '/users/:id',
  USER_CREATE: '/users/create',
  USER_EDIT: '/users/:id/edit',
  ROLES: '/roles',
  ROLE_DETAIL: '/roles/:id',
  ROLE_CREATE: '/roles/create',
  AUDIT_LOGS: '/audit-logs',
  UNAUTHORIZED: '/unauthorized',
  NOT_FOUND: '/404',
} as const;

export const API_ENDPOINTS = {
  AUTH: {
    LOGIN: '/auth/login',
    REGISTER: '/auth/register',
    LOGOUT: '/auth/logout',
    REFRESH: '/auth/refresh',
    ME: '/auth/me',
    FORGOT_PASSWORD: '/auth/forgot-password',
    RESET_PASSWORD: '/auth/reset-password',
    CHANGE_PASSWORD: '/auth/change-password',
    SESSIONS: '/auth/sessions',
    REVOKE_SESSION: '/auth/sessions/:id/revoke',
  },
  USERS: {
    LIST: '/users',
    DETAIL: '/users/:id',
    CREATE: '/users',
    UPDATE: '/users/:id',
    DELETE: '/users/:id',
    ACTIVATE: '/users/:id/activate',
    DEACTIVATE: '/users/:id/deactivate',
    ROLES: '/users/:id/roles',
    ASSIGN_ROLE: '/users/:id/roles',
    REVOKE_ROLE: '/users/:id/roles/:roleId',
  },
  ROLES: {
    LIST: '/roles',
    DETAIL: '/roles/:id',
    CREATE: '/roles',
    UPDATE: '/roles/:id',
    DELETE: '/roles/:id',
    PERMISSIONS: '/roles/:id/permissions',
  },
  PERMISSIONS: {
    LIST: '/permissions',
    DETAIL: '/permissions/:id',
  },
} as const;

export const HTTP_STATUS = {
  OK: 200,
  CREATED: 201,
  NO_CONTENT: 204,
  BAD_REQUEST: 400,
  UNAUTHORIZED: 401,
  FORBIDDEN: 403,
  NOT_FOUND: 404,
  CONFLICT: 409,
  UNPROCESSABLE_ENTITY: 422,
  TOO_MANY_REQUESTS: 429,
  INTERNAL_SERVER_ERROR: 500,
  SERVICE_UNAVAILABLE: 503,
} as const;

export const ERROR_CODES = {
  NETWORK_ERROR: 'NETWORK_ERROR',
  TIMEOUT_ERROR: 'TIMEOUT_ERROR',
  VALIDATION_ERROR: 'VALIDATION_ERROR',
  UNAUTHORIZED: 'UNAUTHORIZED',
  FORBIDDEN: 'FORBIDDEN',
  NOT_FOUND: 'NOT_FOUND',
  CONFLICT: 'CONFLICT',
  SERVER_ERROR: 'SERVER_ERROR',
  UNKNOWN_ERROR: 'UNKNOWN_ERROR',
} as const;
