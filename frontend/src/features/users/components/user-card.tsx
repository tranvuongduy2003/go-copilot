import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardFooter, CardHeader } from '@/components/ui/card';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { ROUTES } from '@/constants';
import type { User, UserStatus } from '@/types';
import {
  Ban,
  CheckCircle,
  Edit,
  Eye,
  Mail,
  MoreVertical,
  Trash2,
  UserCheck,
  UserX,
} from 'lucide-react';
import { Link } from 'react-router-dom';

interface UserCardProps {
  user: User;
  onEdit?: (user: User) => void;
  onDelete?: (user: User) => void;
  onActivate?: (user: User) => void;
  onDeactivate?: (user: User) => void;
}

function getStatusConfig(status: UserStatus) {
  const configs: Record<UserStatus, { color: string; icon: React.ReactNode }> = {
    active: { color: 'bg-green-500', icon: <CheckCircle className="h-3 w-3" /> },
    inactive: { color: 'bg-gray-500', icon: <UserX className="h-3 w-3" /> },
    pending: { color: 'bg-yellow-500', icon: <UserCheck className="h-3 w-3" /> },
    banned: { color: 'bg-red-500', icon: <Ban className="h-3 w-3" /> },
  };
  return configs[status];
}

function getInitials(name: string): string {
  return name
    .split(' ')
    .map((part) => part[0])
    .join('')
    .toUpperCase()
    .slice(0, 2);
}

export function UserCard({ user, onEdit, onDelete, onActivate, onDeactivate }: UserCardProps) {
  const statusConfig = getStatusConfig(user.status);

  return (
    <Card className="relative overflow-hidden transition-shadow hover:shadow-md">
      <CardHeader className="pb-2">
        <div className="flex items-start justify-between">
          <div className="flex items-center gap-3">
            <div className="relative">
              <Avatar className="h-12 w-12">
                <AvatarFallback className="bg-primary/10 text-primary">
                  {getInitials(user.fullName)}
                </AvatarFallback>
              </Avatar>
              <span
                className={`absolute bottom-0 right-0 h-3 w-3 rounded-full border-2 border-background ${statusConfig.color}`}
              />
            </div>
            <div className="min-w-0 flex-1">
              <h3 className="truncate font-semibold">{user.fullName}</h3>
              <div className="flex items-center gap-1 text-sm text-muted-foreground">
                <Mail className="h-3 w-3" />
                <span className="truncate">{user.email}</span>
              </div>
            </div>
          </div>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" className="h-8 w-8">
                <MoreVertical className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem asChild>
                <Link to={ROUTES.USER_DETAIL.replace(':id', user.id)}>
                  <Eye className="mr-2 h-4 w-4" />
                  View Details
                </Link>
              </DropdownMenuItem>
              {onEdit && (
                <DropdownMenuItem onClick={() => onEdit(user)}>
                  <Edit className="mr-2 h-4 w-4" />
                  Edit
                </DropdownMenuItem>
              )}
              <DropdownMenuSeparator />
              {user.status === 'active' && onDeactivate && (
                <DropdownMenuItem onClick={() => onDeactivate(user)}>
                  <UserX className="mr-2 h-4 w-4" />
                  Deactivate
                </DropdownMenuItem>
              )}
              {user.status !== 'active' && onActivate && (
                <DropdownMenuItem onClick={() => onActivate(user)}>
                  <UserCheck className="mr-2 h-4 w-4" />
                  Activate
                </DropdownMenuItem>
              )}
              {onDelete && (
                <>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem
                    onClick={() => onDelete(user)}
                    className="text-destructive focus:text-destructive"
                  >
                    <Trash2 className="mr-2 h-4 w-4" />
                    Delete
                  </DropdownMenuItem>
                </>
              )}
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </CardHeader>
      <CardContent className="pb-2">
        <div className="flex flex-wrap gap-1">
          {user.roles.slice(0, 3).map((role) => (
            <Badge key={role.id} variant="secondary" className="text-xs">
              {role.displayName}
            </Badge>
          ))}
          {user.roles.length > 3 && (
            <Badge variant="outline" className="text-xs">
              +{user.roles.length - 3}
            </Badge>
          )}
        </div>
      </CardContent>
      <CardFooter className="border-t pt-3 text-xs text-muted-foreground">
        <span>
          Joined{' '}
          {new Date(user.createdAt).toLocaleDateString('en-US', {
            month: 'short',
            year: 'numeric',
          })}
        </span>
      </CardFooter>
    </Card>
  );
}

interface UserCardGridProps {
  users: User[];
  onEdit?: (user: User) => void;
  onDelete?: (user: User) => void;
  onActivate?: (user: User) => void;
  onDeactivate?: (user: User) => void;
}

export function UserCardGrid({
  users,
  onEdit,
  onDelete,
  onActivate,
  onDeactivate,
}: UserCardGridProps) {
  if (users.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-12 text-center">
        <p className="text-muted-foreground">No users found</p>
      </div>
    );
  }

  return (
    <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
      {users.map((user) => (
        <UserCard
          key={user.id}
          user={user}
          onEdit={onEdit}
          onDelete={onDelete}
          onActivate={onActivate}
          onDeactivate={onDeactivate}
        />
      ))}
    </div>
  );
}
