export interface User {
  id: string;
  email: string;
  fullName: string;
  status: UserStatus;
  roles: Role[];
  permissions: string[];
  createdAt: string;
  updatedAt: string;
}

export type UserStatus = 'pending' | 'active' | 'inactive' | 'banned';

export interface Role {
  id: string;
  name: string;
  displayName: string;
  description?: string;
  isSystem: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface Permission {
  id: string;
  name: string;
  displayName: string;
  description?: string;
  resource: string;
  action: string;
}

export interface UserWithRoles extends User {
  roles: Role[];
}

export interface CreateUserRequest {
  email: string;
  fullName: string;
  password?: string;
  roleIds?: string[];
}

export interface UpdateUserRequest {
  email?: string;
  fullName?: string;
  status?: UserStatus;
}

export interface UserFilter {
  page?: number;
  pageSize?: number;
  limit?: number;
  search?: string;
  status?: UserStatus;
  roleId?: string;
  sortBy?: 'createdAt' | 'updatedAt' | 'email' | 'fullName';
  sortOrder?: 'asc' | 'desc';
  [key: string]: unknown;
}
