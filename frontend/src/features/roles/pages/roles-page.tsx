import { Pagination } from '@/components/common/pagination';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import type { CreateRoleRequest, Role, UpdateRoleRequest } from '@/types/role';
import { Plus } from 'lucide-react';
import { useCallback, useState } from 'react';
import { toast } from 'sonner';
import { useCreateRole, useDeleteRole, useRoles, useUpdateRole } from '../api/roles.queries';
import { DeleteRoleDialog, PermissionsDialog, RoleFormDialog, RoleTable } from '../components';

interface RoleFilter {
  page: number;
  pageSize: number;
  search?: string;
  [key: string]: unknown;
}

const DEFAULT_FILTERS: RoleFilter = {
  page: 1,
  pageSize: 10,
};

export function RolesPage() {
  const [filters, setFilters] = useState<RoleFilter>(DEFAULT_FILTERS);
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [editingRole, setEditingRole] = useState<Role | null>(null);
  const [deletingRole, setDeletingRole] = useState<Role | null>(null);
  const [managingPermissionsRole, setManagingPermissionsRole] = useState<Role | null>(null);

  const { data: rolesResponse, isLoading } = useRoles(filters);
  const createRole = useCreateRole();
  const updateRole = useUpdateRole();
  const deleteRole = useDeleteRole();

  const roles = rolesResponse?.data ?? [];
  const pagination = rolesResponse?.meta;

  const handlePageChange = useCallback((page: number) => {
    setFilters((prev) => ({ ...prev, page }));
  }, []);

  const handlePageSizeChange = useCallback((pageSize: number) => {
    setFilters((prev) => ({ ...prev, pageSize, page: 1 }));
  }, []);

  const handleCreateRole = useCallback(() => {
    setEditingRole(null);
    setIsFormOpen(true);
  }, []);

  const handleEditRole = useCallback((role: Role) => {
    setEditingRole(role);
    setIsFormOpen(true);
  }, []);

  const handleDeleteRole = useCallback((role: Role) => {
    setDeletingRole(role);
  }, []);

  const handleManagePermissions = useCallback((role: Role) => {
    setManagingPermissionsRole(role);
  }, []);

  const handleFormSubmit = useCallback(
    async (data: CreateRoleRequest | UpdateRoleRequest) => {
      try {
        if (editingRole) {
          await updateRole.mutateAsync({ id: editingRole.id, data: data as UpdateRoleRequest });
          toast.success('Role updated successfully');
        } else {
          await createRole.mutateAsync(data as CreateRoleRequest);
          toast.success('Role created successfully');
        }
        setIsFormOpen(false);
        setEditingRole(null);
      } catch {
        toast.error(editingRole ? 'Failed to update role' : 'Failed to create role');
      }
    },
    [editingRole, createRole, updateRole]
  );

  const handleConfirmDelete = useCallback(async () => {
    if (!deletingRole) return;
    try {
      await deleteRole.mutateAsync(deletingRole.id);
      toast.success('Role deleted successfully');
      setDeletingRole(null);
    } catch {
      toast.error('Failed to delete role');
    }
  }, [deletingRole, deleteRole]);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Roles</h1>
          <p className="text-muted-foreground">Manage roles and their permissions</p>
        </div>
        <Button onClick={handleCreateRole}>
          <Plus className="mr-2 size-4" />
          Add Role
        </Button>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>All Roles</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <RoleTable
            roles={roles}
            onEdit={handleEditRole}
            onDelete={handleDeleteRole}
            onManagePermissions={handleManagePermissions}
            isLoading={isLoading}
          />
          {pagination && (
            <Pagination
              page={pagination.page}
              pageSize={pagination.pageSize}
              total={pagination.total}
              onPageChange={handlePageChange}
              onPageSizeChange={handlePageSizeChange}
            />
          )}
        </CardContent>
      </Card>

      <RoleFormDialog
        open={isFormOpen}
        onOpenChange={setIsFormOpen}
        role={editingRole}
        onSubmit={handleFormSubmit}
        isLoading={createRole.isPending || updateRole.isPending}
      />

      <DeleteRoleDialog
        open={!!deletingRole}
        onOpenChange={(open) => !open && setDeletingRole(null)}
        role={deletingRole}
        onConfirm={handleConfirmDelete}
        isLoading={deleteRole.isPending}
      />

      <PermissionsDialog
        open={!!managingPermissionsRole}
        onOpenChange={(open) => !open && setManagingPermissionsRole(null)}
        role={managingPermissionsRole}
      />
    </div>
  );
}
