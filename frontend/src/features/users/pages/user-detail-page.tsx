import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Separator } from '@/components/ui/separator';
import { Skeleton } from '@/components/ui/skeleton';
import { ROUTES } from '@/constants';
import { ArrowLeft, Edit, Mail, Shield, User as UserIcon } from 'lucide-react';
import { Link, useNavigate, useParams } from 'react-router-dom';
import { useUser, useUserRoles } from '../api/users.queries';
import { UserStatusBadge } from '../components/user-table';

function UserDetailSkeleton() {
  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <Skeleton className="h-10 w-10" />
        <div className="space-y-2">
          <Skeleton className="h-8 w-48" />
          <Skeleton className="h-4 w-32" />
        </div>
      </div>
      <div className="grid gap-6 md:grid-cols-2">
        <Skeleton className="h-48" />
        <Skeleton className="h-48" />
      </div>
    </div>
  );
}

export function UserDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const { data: user, isLoading: isLoadingUser, error: userError } = useUser(id || '');
  const { data: roles, isLoading: isLoadingRoles } = useUserRoles(id || '');

  if (isLoadingUser) {
    return (
      <div className="container py-6">
        <UserDetailSkeleton />
      </div>
    );
  }

  if (userError || !user) {
    return (
      <div className="container py-6">
        <div className="flex flex-col items-center justify-center py-12 text-center">
          <UserIcon className="h-12 w-12 text-muted-foreground" />
          <h2 className="mt-4 text-xl font-semibold">User Not Found</h2>
          <p className="mt-2 text-muted-foreground">
            The user you're looking for doesn't exist or has been deleted.
          </p>
          <Button asChild className="mt-4">
            <Link to={ROUTES.USERS}>Back to Users</Link>
          </Button>
        </div>
      </div>
    );
  }

  const userRoles = roles || user.roles || [];

  return (
    <div className="container py-6">
      <div className="mb-6 flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Button variant="ghost" size="icon" onClick={() => navigate(-1)}>
            <ArrowLeft className="h-4 w-4" />
          </Button>
          <div>
            <h1 className="text-2xl font-bold">{user.fullName}</h1>
            <p className="text-sm text-muted-foreground">{user.email}</p>
          </div>
        </div>
        <Button asChild>
          <Link to={ROUTES.USER_EDIT.replace(':id', user.id)}>
            <Edit className="mr-2 h-4 w-4" />
            Edit User
          </Link>
        </Button>
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <UserIcon className="h-5 w-5" />
              User Information
            </CardTitle>
            <CardDescription>Basic user details and status</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center justify-between">
              <span className="text-sm font-medium text-muted-foreground">Status</span>
              <UserStatusBadge status={user.status} />
            </div>
            <Separator />
            <div className="flex items-center justify-between">
              <span className="text-sm font-medium text-muted-foreground">Email</span>
              <div className="flex items-center gap-2">
                <Mail className="h-4 w-4 text-muted-foreground" />
                <span className="text-sm">{user.email}</span>
              </div>
            </div>
            <Separator />
            <div className="flex items-center justify-between">
              <span className="text-sm font-medium text-muted-foreground">Full Name</span>
              <span className="text-sm">{user.fullName}</span>
            </div>
            <Separator />
            <div className="flex items-center justify-between">
              <span className="text-sm font-medium text-muted-foreground">Created</span>
              <span className="text-sm">
                {new Date(user.createdAt).toLocaleDateString('en-US', {
                  year: 'numeric',
                  month: 'long',
                  day: 'numeric',
                })}
              </span>
            </div>
            <Separator />
            <div className="flex items-center justify-between">
              <span className="text-sm font-medium text-muted-foreground">Last Updated</span>
              <span className="text-sm">
                {new Date(user.updatedAt).toLocaleDateString('en-US', {
                  year: 'numeric',
                  month: 'long',
                  day: 'numeric',
                })}
              </span>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Shield className="h-5 w-5" />
              Roles & Permissions
            </CardTitle>
            <CardDescription>Assigned roles and effective permissions</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div>
              <h4 className="mb-2 text-sm font-medium text-muted-foreground">Roles</h4>
              {isLoadingRoles ? (
                <div className="flex gap-2">
                  <Skeleton className="h-6 w-16" />
                  <Skeleton className="h-6 w-20" />
                </div>
              ) : userRoles.length > 0 ? (
                <div className="flex flex-wrap gap-2">
                  {userRoles.map((role) => (
                    <Badge key={role.id} variant={role.isSystem ? 'default' : 'secondary'}>
                      {role.displayName}
                    </Badge>
                  ))}
                </div>
              ) : (
                <p className="text-sm text-muted-foreground">No roles assigned</p>
              )}
            </div>
            <Separator />
            <div>
              <h4 className="mb-2 text-sm font-medium text-muted-foreground">Permissions</h4>
              {user.permissions && user.permissions.length > 0 ? (
                <div className="flex flex-wrap gap-1">
                  {user.permissions.slice(0, 10).map((permission) => (
                    <Badge key={permission} variant="outline" className="text-xs">
                      {permission}
                    </Badge>
                  ))}
                  {user.permissions.length > 10 && (
                    <Badge variant="outline" className="text-xs">
                      +{user.permissions.length - 10} more
                    </Badge>
                  )}
                </div>
              ) : (
                <p className="text-sm text-muted-foreground">No permissions</p>
              )}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
