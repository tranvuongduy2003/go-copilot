import { API_ENDPOINTS } from '@/constants';
import { apiClient } from '@/lib/api/client';
import type { ApiResponse, PaginatedResponse } from '@/types/api';
import type { CreateRoleRequest, Permission, Role, UpdateRoleRequest } from '@/types/role';

export interface RoleFilter {
  page?: number;
  pageSize?: number;
  search?: string;
  [key: string]: unknown;
}

export async function getRoles(filter?: RoleFilter): Promise<PaginatedResponse<Role>> {
  const response = await apiClient.get<PaginatedResponse<Role>>(API_ENDPOINTS.ROLES.LIST, {
    params: filter,
  });
  return response.data;
}

export async function getRole(id: string): Promise<ApiResponse<Role>> {
  const url = API_ENDPOINTS.ROLES.DETAIL.replace(':id', id);
  const response = await apiClient.get<ApiResponse<Role>>(url);
  return response.data;
}

export async function getRolePermissions(id: string): Promise<ApiResponse<Permission[]>> {
  const url = API_ENDPOINTS.ROLES.PERMISSIONS.replace(':id', id);
  const response = await apiClient.get<ApiResponse<Permission[]>>(url);
  return response.data;
}

export async function createRole(data: CreateRoleRequest): Promise<ApiResponse<Role>> {
  const response = await apiClient.post<ApiResponse<Role>>(API_ENDPOINTS.ROLES.CREATE, data);
  return response.data;
}

export async function updateRole(id: string, data: UpdateRoleRequest): Promise<ApiResponse<Role>> {
  const url = API_ENDPOINTS.ROLES.UPDATE.replace(':id', id);
  const response = await apiClient.put<ApiResponse<Role>>(url, data);
  return response.data;
}

export async function deleteRole(id: string): Promise<void> {
  const url = API_ENDPOINTS.ROLES.DELETE.replace(':id', id);
  await apiClient.delete(url);
}

export async function assignPermissionToRole(roleId: string, permissionId: string): Promise<void> {
  const url = API_ENDPOINTS.ROLES.PERMISSIONS.replace(':id', roleId);
  await apiClient.post(`${url}/${permissionId}`);
}

export async function revokePermissionFromRole(
  roleId: string,
  permissionId: string
): Promise<void> {
  const url = API_ENDPOINTS.ROLES.PERMISSIONS.replace(':id', roleId);
  await apiClient.delete(`${url}/${permissionId}`);
}

export async function getPermissions(): Promise<PaginatedResponse<Permission>> {
  const response = await apiClient.get<PaginatedResponse<Permission>>(
    API_ENDPOINTS.PERMISSIONS.LIST
  );
  return response.data;
}

export async function getPermission(id: string): Promise<ApiResponse<Permission>> {
  const url = API_ENDPOINTS.PERMISSIONS.DETAIL.replace(':id', id);
  const response = await apiClient.get<ApiResponse<Permission>>(url);
  return response.data;
}
