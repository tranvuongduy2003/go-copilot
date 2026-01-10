import type { Permission, Role } from './user';

export type { Permission, Role } from './user';

export interface CreateRoleRequest {
  name: string;
  displayName: string;
  description?: string;
  permissionIds?: string[];
}

export interface UpdateRoleRequest {
  displayName?: string;
  description?: string;
  permissionIds?: string[];
}

export interface RoleWithPermissions extends Role {
  permissions: Permission[];
}

export interface SetRolePermissionsRequest {
  permissionIds: string[];
}
