import { Checkbox } from '@/components/ui/checkbox';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Spinner } from '@/components/ui/spinner';
import type { Permission, Role } from '@/types/role';
import { Search } from 'lucide-react';
import { useMemo, useState } from 'react';
import { toast } from 'sonner';
import {
  useAssignPermission,
  usePermissions,
  useRevokePermission,
  useRolePermissions,
} from '../api/roles.queries';

interface PermissionsDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  role: Role | null;
}

function groupPermissionsByResource(permissions: Permission[]): Record<string, Permission[]> {
  return permissions.reduce(
    (groups, permission) => {
      const resource = permission.resource || 'other';
      if (!groups[resource]) {
        groups[resource] = [];
      }
      groups[resource].push(permission);
      return groups;
    },
    {} as Record<string, Permission[]>
  );
}

export function PermissionsDialog({ open, onOpenChange, role }: PermissionsDialogProps) {
  const [search, setSearch] = useState('');

  const { data: rolePermissionsResponse, isLoading: isLoadingRolePermissions } = useRolePermissions(
    role?.id ?? ''
  );
  const { data: allPermissionsResponse, isLoading: isLoadingAllPermissions } = usePermissions();
  const assignPermission = useAssignPermission();
  const revokePermission = useRevokePermission();

  const rolePermissionIds = useMemo(() => {
    const currentRolePermissions = rolePermissionsResponse?.data ?? [];
    return new Set(currentRolePermissions.map((permission) => permission.id));
  }, [rolePermissionsResponse?.data]);

  const filteredPermissions = useMemo(() => {
    const permissions = allPermissionsResponse?.data ?? [];
    if (!search) return permissions;
    const searchLower = search.toLowerCase();
    return permissions.filter(
      (permission) =>
        permission.name.toLowerCase().includes(searchLower) ||
        permission.description?.toLowerCase().includes(searchLower) ||
        permission.resource?.toLowerCase().includes(searchLower)
    );
  }, [allPermissionsResponse?.data, search]);

  const groupedPermissions = useMemo(
    () => groupPermissionsByResource(filteredPermissions),
    [filteredPermissions]
  );

  const handleTogglePermission = async (permission: Permission) => {
    if (!role) return;

    const hasPermission = rolePermissionIds.has(permission.id);

    try {
      if (hasPermission) {
        await revokePermission.mutateAsync({
          roleId: role.id,
          permissionId: permission.id,
        });
        toast.success(`Revoked ${permission.name} permission`);
      } else {
        await assignPermission.mutateAsync({
          roleId: role.id,
          permissionId: permission.id,
        });
        toast.success(`Assigned ${permission.name} permission`);
      }
    } catch {
      toast.error(`Failed to ${hasPermission ? 'revoke' : 'assign'} permission`);
    }
  };

  if (!role) return null;

  const isLoading = isLoadingRolePermissions || isLoadingAllPermissions;
  const isMutating = assignPermission.isPending || revokePermission.isPending;

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-h-[80vh] overflow-hidden sm:max-w-2xl">
        <DialogHeader>
          <DialogTitle>Manage Permissions for {role.name}</DialogTitle>
          <DialogDescription>
            Select the permissions you want to assign to this role.
          </DialogDescription>
        </DialogHeader>

        <div className="relative">
          <Search className="absolute left-2.5 top-2.5 size-4 text-muted-foreground" />
          <Input
            type="search"
            placeholder="Search permissions..."
            value={search}
            onChange={(event) => setSearch(event.target.value)}
            className="pl-8"
          />
        </div>

        <div className="max-h-96 space-y-6 overflow-y-auto pr-2">
          {isLoading ? (
            <div className="flex items-center justify-center py-8">
              <Spinner className="size-6" />
            </div>
          ) : Object.keys(groupedPermissions).length === 0 ? (
            <div className="py-8 text-center text-muted-foreground">No permissions found</div>
          ) : (
            Object.entries(groupedPermissions).map(([resource, permissions]) => (
              <div key={resource} className="space-y-3">
                <h4 className="text-sm font-semibold capitalize">{resource}</h4>
                <div className="grid gap-2">
                  {permissions.map((permission) => {
                    const isChecked = rolePermissionIds.has(permission.id);
                    return (
                      <div
                        key={permission.id}
                        className="flex items-start space-x-3 rounded-md border p-3"
                      >
                        <Checkbox
                          id={permission.id}
                          checked={isChecked}
                          onCheckedChange={() => handleTogglePermission(permission)}
                          disabled={isMutating}
                        />
                        <div className="space-y-1">
                          <Label
                            htmlFor={permission.id}
                            className="cursor-pointer font-medium leading-none"
                          >
                            {permission.name}
                          </Label>
                          {permission.description && (
                            <p className="text-sm text-muted-foreground">
                              {permission.description}
                            </p>
                          )}
                        </div>
                      </div>
                    );
                  })}
                </div>
              </div>
            ))
          )}
        </div>
      </DialogContent>
    </Dialog>
  );
}
