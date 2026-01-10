import { queryKeys } from '@/lib/api/query-client';
import type { CreateRoleRequest, UpdateRoleRequest } from '@/types/role';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  assignPermissionToRole,
  createRole,
  deleteRole,
  getPermission,
  getPermissions,
  getRole,
  getRolePermissions,
  getRoles,
  revokePermissionFromRole,
  updateRole,
  type RoleFilter,
} from './roles.api';

export function useRoles(filter?: RoleFilter) {
  return useQuery({
    queryKey: queryKeys.roles.list(filter),
    queryFn: () => getRoles(filter),
  });
}

export function useRole(id: string) {
  return useQuery({
    queryKey: queryKeys.roles.detail(id),
    queryFn: () => getRole(id),
    enabled: !!id,
  });
}

export function useRolePermissions(id: string) {
  return useQuery({
    queryKey: [...queryKeys.roles.detail(id), 'permissions'],
    queryFn: () => getRolePermissions(id),
    enabled: !!id,
  });
}

export function usePermissions() {
  return useQuery({
    queryKey: queryKeys.permissions.all,
    queryFn: () => getPermissions(),
  });
}

export function usePermission(id: string) {
  return useQuery({
    queryKey: queryKeys.permissions.detail(id),
    queryFn: () => getPermission(id),
    enabled: !!id,
  });
}

export function useCreateRole() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateRoleRequest) => createRole(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.roles.all });
    },
  });
}

export function useUpdateRole() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateRoleRequest }) => updateRole(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.roles.all });
      queryClient.invalidateQueries({ queryKey: queryKeys.roles.detail(variables.id) });
    },
  });
}

export function useDeleteRole() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => deleteRole(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.roles.all });
    },
  });
}

export function useAssignPermission() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ roleId, permissionId }: { roleId: string; permissionId: string }) =>
      assignPermissionToRole(roleId, permissionId),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.roles.detail(variables.roleId) });
      queryClient.invalidateQueries({
        queryKey: [...queryKeys.roles.detail(variables.roleId), 'permissions'],
      });
    },
  });
}

export function useRevokePermission() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ roleId, permissionId }: { roleId: string; permissionId: string }) =>
      revokePermissionFromRole(roleId, permissionId),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.roles.detail(variables.roleId) });
      queryClient.invalidateQueries({
        queryKey: [...queryKeys.roles.detail(variables.roleId), 'permissions'],
      });
    },
  });
}
