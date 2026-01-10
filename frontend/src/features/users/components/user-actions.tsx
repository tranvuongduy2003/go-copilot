import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import type { User } from '@/types/user';
import { Edit, MoreHorizontal, Shield, Trash2, UserCheck, UserX } from 'lucide-react';

interface UserActionsProps {
  user: User;
  onEdit: () => void;
  onDelete: () => void;
  onActivate: () => void;
  onDeactivate: () => void;
  onManageRoles?: () => void;
}

export function UserActions({
  user,
  onEdit,
  onDelete,
  onActivate,
  onDeactivate,
  onManageRoles,
}: UserActionsProps) {
  const canActivate = user.status === 'pending' || user.status === 'banned';
  const canDeactivate = user.status === 'active';

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
        {onManageRoles && (
          <DropdownMenuItem onClick={onManageRoles}>
            <Shield className="mr-2 size-4" />
            Manage Roles
          </DropdownMenuItem>
        )}
        <DropdownMenuSeparator />
        {canActivate && (
          <DropdownMenuItem onClick={onActivate}>
            <UserCheck className="mr-2 size-4" />
            Activate
          </DropdownMenuItem>
        )}
        {canDeactivate && (
          <DropdownMenuItem onClick={onDeactivate}>
            <UserX className="mr-2 size-4" />
            Deactivate
          </DropdownMenuItem>
        )}
        <DropdownMenuSeparator />
        <DropdownMenuItem onClick={onDelete} className="text-destructive focus:text-destructive">
          <Trash2 className="mr-2 size-4" />
          Delete
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
