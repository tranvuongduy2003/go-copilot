import { API_ENDPOINTS } from '@/constants';
import { api } from '@/lib/api';
import type {
  ApiResponse,
  CreateUserRequest,
  PaginatedResponse,
  Role,
  UpdateUserRequest,
  User,
  UserFilter,
} from '@/types';

export const usersApi = {
  getUsers: async (filters?: UserFilter): Promise<PaginatedResponse<User>> => {
    const params = new URLSearchParams();

    if (filters) {
      Object.entries(filters).forEach(([key, value]) => {
        if (value !== undefined && value !== null && value !== '') {
          params.append(key, String(value));
        }
      });
    }

    const queryString = params.toString();
    const url = queryString
      ? `${API_ENDPOINTS.USERS.LIST}?${queryString}`
      : API_ENDPOINTS.USERS.LIST;

    return api.get<PaginatedResponse<User>>(url);
  },

  getUser: async (id: string): Promise<User> => {
    const url = API_ENDPOINTS.USERS.DETAIL.replace(':id', id);
    const response = await api.get<ApiResponse<User>>(url);
    return response.data;
  },

  createUser: async (data: CreateUserRequest): Promise<User> => {
    const response = await api.post<ApiResponse<User>>(API_ENDPOINTS.USERS.CREATE, data);
    return response.data;
  },

  updateUser: async (id: string, data: UpdateUserRequest): Promise<User> => {
    const url = API_ENDPOINTS.USERS.UPDATE.replace(':id', id);
    const response = await api.put<ApiResponse<User>>(url, data);
    return response.data;
  },

  deleteUser: async (id: string): Promise<void> => {
    const url = API_ENDPOINTS.USERS.DELETE.replace(':id', id);
    await api.delete(url);
  },

  activateUser: async (id: string): Promise<User> => {
    const url = API_ENDPOINTS.USERS.ACTIVATE.replace(':id', id);
    const response = await api.post<ApiResponse<User>>(url);
    return response.data;
  },

  deactivateUser: async (id: string): Promise<User> => {
    const url = API_ENDPOINTS.USERS.DEACTIVATE.replace(':id', id);
    const response = await api.post<ApiResponse<User>>(url);
    return response.data;
  },

  getUserRoles: async (id: string): Promise<Role[]> => {
    const url = API_ENDPOINTS.USERS.ROLES.replace(':id', id);
    const response = await api.get<ApiResponse<Role[]>>(url);
    return response.data;
  },

  assignRole: async (userId: string, roleId: string): Promise<void> => {
    const url = API_ENDPOINTS.USERS.ASSIGN_ROLE.replace(':id', userId);
    await api.post(url, { roleId });
  },

  revokeRole: async (userId: string, roleId: string): Promise<void> => {
    const url = API_ENDPOINTS.USERS.REVOKE_ROLE.replace(':id', userId).replace(':roleId', roleId);
    await api.delete(url);
  },
};
