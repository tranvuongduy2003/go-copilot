import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import type { Role } from '@/types/role';
import { Edit, MoreHorizontal, Shield, Trash2 } from 'lucide-react';

interface RoleActionsProps {
  role: Role;
  onEdit: () => void;
  onDelete: () => void;
  onManagePermissions: () => void;
}

export function RoleActions({ role, onEdit, onDelete, onManagePermissions }: RoleActionsProps) {
  const isSystemRole = role.name === 'admin' || role.name === 'user';

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" size="icon" className="size-8">
          <MoreHorizontal className="size-4" />
          <span className="sr-only">Open menu</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        <DropdownMenuItem onClick={onEdit}>
          <Edit className="mr-2 size-4" />
          Edit
        </DropdownMenuItem>
        <DropdownMenuItem onClick={onManagePermissions}>
          <Shield className="mr-2 size-4" />
          Manage Permissions
        </DropdownMenuItem>
        {!isSystemRole && (
          <>
            <DropdownMenuSeparator />
            <DropdownMenuItem
              onClick={onDelete}
              className="text-destructive focus:text-destructive"
            >
              <Trash2 className="mr-2 size-4" />
              Delete
            </DropdownMenuItem>
          </>
        )}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
