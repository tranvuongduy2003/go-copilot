import { Pagination } from '@/components/common/pagination';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import type { CreateUserRequest, UpdateUserRequest, User, UserFilter } from '@/types/user';
import { Plus } from 'lucide-react';
import { useCallback, useMemo, useState } from 'react';
import { toast } from 'sonner';
import {
  useActivateUser,
  useCreateUser,
  useDeactivateUser,
  useDeleteUser,
  useUpdateUser,
  useUsers,
} from '../api/users.queries';
import { DeleteUserDialog, UserFilters, UserFormDialog, UserTable } from '../components';

const DEFAULT_FILTERS: UserFilter = {
  page: 1,
  pageSize: 10,
};

export function UsersPage() {
  const [filters, setFilters] = useState<UserFilter>(DEFAULT_FILTERS);
  const [selectedUsers, setSelectedUsers] = useState<string[]>([]);
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [editingUser, setEditingUser] = useState<User | null>(null);
  const [deletingUser, setDeletingUser] = useState<User | null>(null);

  const { data: usersResponse, isLoading } = useUsers(filters);
  const createUser = useCreateUser();
  const updateUser = useUpdateUser();
  const deleteUser = useDeleteUser();
  const activateUser = useActivateUser();
  const deactivateUser = useDeactivateUser();

  const users = useMemo(() => usersResponse?.data ?? [], [usersResponse]);
  const pagination = useMemo(() => usersResponse?.meta, [usersResponse]);

  const handleFiltersChange = useCallback((newFilters: UserFilter) => {
    setFilters(newFilters);
    setSelectedUsers([]);
  }, []);

  const handleResetFilters = useCallback(() => {
    setFilters(DEFAULT_FILTERS);
    setSelectedUsers([]);
  }, []);

  const handleSelectUser = useCallback((userId: string) => {
    setSelectedUsers((prev) =>
      prev.includes(userId) ? prev.filter((id) => id !== userId) : [...prev, userId]
    );
  }, []);

  const handleSelectAll = useCallback(() => {
    setSelectedUsers((prev) => (prev.length === users.length ? [] : users.map((user) => user.id)));
  }, [users]);

  const handlePageChange = useCallback((page: number) => {
    setFilters((prev) => ({ ...prev, page }));
  }, []);

  const handlePageSizeChange = useCallback((pageSize: number) => {
    setFilters((prev) => ({ ...prev, pageSize, page: 1 }));
  }, []);

  const handleCreateUser = useCallback(() => {
    setEditingUser(null);
    setIsFormOpen(true);
  }, []);

  const handleEditUser = useCallback((user: User) => {
    setEditingUser(user);
    setIsFormOpen(true);
  }, []);

  const handleDeleteUser = useCallback((user: User) => {
    setDeletingUser(user);
  }, []);

  const handleFormSubmit = useCallback(
    async (data: CreateUserRequest | UpdateUserRequest) => {
      try {
        if (editingUser) {
          await updateUser.mutateAsync({ id: editingUser.id, data: data as UpdateUserRequest });
          toast.success('User updated successfully');
        } else {
          await createUser.mutateAsync(data as CreateUserRequest);
          toast.success('User created successfully');
        }
        setIsFormOpen(false);
        setEditingUser(null);
      } catch {
        toast.error(editingUser ? 'Failed to update user' : 'Failed to create user');
      }
    },
    [editingUser, createUser, updateUser]
  );

  const handleConfirmDelete = useCallback(async () => {
    if (!deletingUser) return;
    try {
      await deleteUser.mutateAsync(deletingUser.id);
      toast.success('User deleted successfully');
      setDeletingUser(null);
    } catch {
      toast.error('Failed to delete user');
    }
  }, [deletingUser, deleteUser]);

  const handleActivateUser = useCallback(
    async (user: User) => {
      try {
        await activateUser.mutateAsync(user.id);
        toast.success('User activated successfully');
      } catch {
        toast.error('Failed to activate user');
      }
    },
    [activateUser]
  );

  const handleDeactivateUser = useCallback(
    async (user: User) => {
      try {
        await deactivateUser.mutateAsync(user.id);
        toast.success('User deactivated successfully');
      } catch {
        toast.error('Failed to deactivate user');
      }
    },
    [deactivateUser]
  );

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Users</h1>
          <p className="text-muted-foreground">Manage user accounts and their roles</p>
        </div>
        <Button onClick={handleCreateUser}>
          <Plus className="mr-2 size-4" />
          Add User
        </Button>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>All Users</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <UserFilters
            filters={filters}
            onFiltersChange={handleFiltersChange}
            onReset={handleResetFilters}
          />
          <UserTable
            users={users}
            selectedUsers={selectedUsers}
            onSelectUser={handleSelectUser}
            onSelectAll={handleSelectAll}
            onEdit={handleEditUser}
            onDelete={handleDeleteUser}
            onActivate={handleActivateUser}
            onDeactivate={handleDeactivateUser}
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

      <UserFormDialog
        open={isFormOpen}
        onOpenChange={setIsFormOpen}
        user={editingUser}
        onSubmit={handleFormSubmit}
        isLoading={createUser.isPending || updateUser.isPending}
      />

      <DeleteUserDialog
        open={!!deletingUser}
        onOpenChange={(open) => !open && setDeletingUser(null)}
        user={deletingUser}
        onConfirm={handleConfirmDelete}
        isLoading={deleteUser.isPending}
      />
    </div>
  );
}
