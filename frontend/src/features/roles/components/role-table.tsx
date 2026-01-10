import { Badge } from '@/components/ui/badge';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import type { Permission, Role } from '@/types/role';
import { RoleActions } from './role-actions';

interface RoleWithOptionalPermissions extends Role {
  permissions?: Permission[];
}

interface RoleTableProps {
  roles: RoleWithOptionalPermissions[];
  onEdit: (role: RoleWithOptionalPermissions) => void;
  onDelete: (role: RoleWithOptionalPermissions) => void;
  onManagePermissions: (role: RoleWithOptionalPermissions) => void;
  isLoading?: boolean;
}

function formatDate(dateString: string): string {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  });
}

export function RoleTable({
  roles,
  onEdit,
  onDelete,
  onManagePermissions,
  isLoading,
}: RoleTableProps) {
  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="text-muted-foreground">Loading roles...</div>
      </div>
    );
  }

  if (roles.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-12">
        <div className="text-muted-foreground">No roles found</div>
        <p className="mt-1 text-sm text-muted-foreground">Create a new role to get started</p>
      </div>
    );
  }

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Name</TableHead>
          <TableHead>Description</TableHead>
          <TableHead>Permissions</TableHead>
          <TableHead>Created</TableHead>
          <TableHead className="w-12">
            <span className="sr-only">Actions</span>
          </TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {roles.map((role) => (
          <TableRow key={role.id}>
            <TableCell>
              <div className="font-medium">{role.name}</div>
            </TableCell>
            <TableCell className="max-w-xs truncate text-muted-foreground">
              {role.description || '-'}
            </TableCell>
            <TableCell>
              <Badge variant="secondary">{role.permissions?.length ?? 0} permissions</Badge>
            </TableCell>
            <TableCell className="text-muted-foreground">{formatDate(role.createdAt)}</TableCell>
            <TableCell>
              <RoleActions
                role={role}
                onEdit={() => onEdit(role)}
                onDelete={() => onDelete(role)}
                onManagePermissions={() => onManagePermissions(role)}
              />
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
