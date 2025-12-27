# React Component Builder Skill

Generate React components following design system and best practices.

## Usage

```
/project:skill:react-component <component-name>
```

## Generated Files

For a component named `UserCard`, this skill generates:

### 1. Component File

**`src/components/features/user/user-card.tsx`**
```tsx
import { forwardRef } from 'react';
import { cn } from '@/lib/utils';
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import type { User } from '@/types/user';

interface UserCardProps extends React.HTMLAttributes<HTMLDivElement> {
    user: User;
    onEdit?: (user: User) => void;
    onDelete?: (user: User) => void;
    variant?: 'default' | 'compact' | 'detailed';
}

export const UserCard = forwardRef<HTMLDivElement, UserCardProps>(
    ({ user, onEdit, onDelete, variant = 'default', className, ...props }, ref) => {
        const initials = user.name
            .split(' ')
            .map((word) => word[0])
            .join('')
            .toUpperCase()
            .slice(0, 2);

        return (
            <Card
                ref={ref}
                className={cn(
                    'transition-shadow hover:shadow-md',
                    {
                        'p-4': variant === 'default',
                        'p-2': variant === 'compact',
                        'p-6': variant === 'detailed',
                    },
                    className
                )}
                {...props}
            >
                <CardHeader className="flex flex-row items-center gap-4">
                    <Avatar className="h-12 w-12">
                        <AvatarImage src={user.avatarUrl} alt={user.name} />
                        <AvatarFallback className="bg-primary text-primary-foreground">
                            {initials}
                        </AvatarFallback>
                    </Avatar>
                    <div className="flex-1">
                        <CardTitle className="text-lg">{user.name}</CardTitle>
                        <p className="text-sm text-muted-foreground">{user.email}</p>
                    </div>
                </CardHeader>

                {variant === 'detailed' && (
                    <CardContent>
                        <dl className="space-y-2 text-sm">
                            <div className="flex justify-between">
                                <dt className="text-muted-foreground">Role</dt>
                                <dd className="font-medium">{user.role}</dd>
                            </div>
                            <div className="flex justify-between">
                                <dt className="text-muted-foreground">Status</dt>
                                <dd>
                                    <span
                                        className={cn(
                                            'inline-flex items-center rounded-full px-2 py-1 text-xs font-medium',
                                            {
                                                'bg-success/10 text-success': user.status === 'active',
                                                'bg-muted text-muted-foreground': user.status === 'inactive',
                                            }
                                        )}
                                    >
                                        {user.status}
                                    </span>
                                </dd>
                            </div>
                            <div className="flex justify-between">
                                <dt className="text-muted-foreground">Joined</dt>
                                <dd className="font-medium">
                                    {new Date(user.createdAt).toLocaleDateString()}
                                </dd>
                            </div>
                        </dl>
                    </CardContent>
                )}

                {(onEdit || onDelete) && (
                    <CardFooter className="flex justify-end gap-2">
                        {onEdit && (
                            <Button variant="outline" size="sm" onClick={() => onEdit(user)}>
                                Edit
                            </Button>
                        )}
                        {onDelete && (
                            <Button
                                variant="destructive"
                                size="sm"
                                onClick={() => onDelete(user)}
                            >
                                Delete
                            </Button>
                        )}
                    </CardFooter>
                )}
            </Card>
        );
    }
);

UserCard.displayName = 'UserCard';
```

### 2. Type Definition

**`src/types/user.ts`**
```tsx
export interface User {
    id: string;
    name: string;
    email: string;
    avatarUrl?: string;
    role: 'admin' | 'user' | 'guest';
    status: 'active' | 'inactive';
    createdAt: string;
    updatedAt: string;
}

export interface CreateUserInput {
    name: string;
    email: string;
    password: string;
    role?: 'admin' | 'user' | 'guest';
}

export interface UpdateUserInput {
    name?: string;
    email?: string;
    role?: 'admin' | 'user' | 'guest';
    status?: 'active' | 'inactive';
}
```

### 3. Hook File

**`src/hooks/use-user.ts`**
```tsx
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '@/lib/api';
import type { User, CreateUserInput, UpdateUserInput } from '@/types/user';

export function useUser(id: string) {
    return useQuery({
        queryKey: ['user', id],
        queryFn: () => api.get<User>(`/users/${id}`),
        enabled: !!id,
    });
}

export function useUsers() {
    return useQuery({
        queryKey: ['users'],
        queryFn: () => api.get<User[]>('/users'),
    });
}

export function useCreateUser() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: (input: CreateUserInput) => api.post<User>('/users', input),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['users'] });
        },
    });
}

export function useUpdateUser() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: ({ id, ...data }: { id: string } & UpdateUserInput) =>
            api.put<User>(`/users/${id}`, data),
        onSuccess: (updatedUser) => {
            queryClient.setQueryData(['user', updatedUser.id], updatedUser);
            queryClient.invalidateQueries({ queryKey: ['users'] });
        },
    });
}

export function useDeleteUser() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: (id: string) => api.delete(`/users/${id}`),
        onSuccess: (_, deletedId) => {
            queryClient.removeQueries({ queryKey: ['user', deletedId] });
            queryClient.invalidateQueries({ queryKey: ['users'] });
        },
    });
}
```

### 4. Test File

**`src/components/features/user/user-card.test.tsx`**
```tsx
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { UserCard } from './user-card';
import type { User } from '@/types/user';

const mockUser: User = {
    id: '1',
    name: 'John Doe',
    email: 'john@example.com',
    role: 'user',
    status: 'active',
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z',
};

describe('UserCard', () => {
    it('renders user name and email', () => {
        render(<UserCard user={mockUser} />);

        expect(screen.getByText('John Doe')).toBeInTheDocument();
        expect(screen.getByText('john@example.com')).toBeInTheDocument();
    });

    it('displays user initials in avatar fallback', () => {
        render(<UserCard user={mockUser} />);

        expect(screen.getByText('JD')).toBeInTheDocument();
    });

    it('calls onEdit when edit button is clicked', () => {
        const handleEdit = vi.fn();
        render(<UserCard user={mockUser} onEdit={handleEdit} />);

        fireEvent.click(screen.getByRole('button', { name: /edit/i }));

        expect(handleEdit).toHaveBeenCalledWith(mockUser);
    });

    it('calls onDelete when delete button is clicked', () => {
        const handleDelete = vi.fn();
        render(<UserCard user={mockUser} onDelete={handleDelete} />);

        fireEvent.click(screen.getByRole('button', { name: /delete/i }));

        expect(handleDelete).toHaveBeenCalledWith(mockUser);
    });

    it('shows detailed information in detailed variant', () => {
        render(<UserCard user={mockUser} variant="detailed" />);

        expect(screen.getByText('Role')).toBeInTheDocument();
        expect(screen.getByText('Status')).toBeInTheDocument();
        expect(screen.getByText('Joined')).toBeInTheDocument();
    });

    it('hides action buttons when callbacks not provided', () => {
        render(<UserCard user={mockUser} />);

        expect(screen.queryByRole('button', { name: /edit/i })).not.toBeInTheDocument();
        expect(screen.queryByRole('button', { name: /delete/i })).not.toBeInTheDocument();
    });
});
```

### 5. List Component

**`src/components/features/user/user-list.tsx`**
```tsx
import { UserCard } from './user-card';
import { Skeleton } from '@/components/ui/skeleton';
import { useUsers } from '@/hooks/use-user';
import type { User } from '@/types/user';

interface UserListProps {
    onEdit?: (user: User) => void;
    onDelete?: (user: User) => void;
}

export function UserList({ onEdit, onDelete }: UserListProps) {
    const { data: users, isLoading, error } = useUsers();

    if (isLoading) {
        return (
            <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
                {Array.from({ length: 6 }).map((_, index) => (
                    <Skeleton key={index} className="h-32 rounded-lg" />
                ))}
            </div>
        );
    }

    if (error) {
        return (
            <div className="rounded-lg border border-destructive/20 bg-destructive/5 p-4 text-center">
                <p className="text-destructive">Failed to load users. Please try again.</p>
            </div>
        );
    }

    if (!users || users.length === 0) {
        return (
            <div className="rounded-lg border border-dashed p-8 text-center">
                <p className="text-muted-foreground">No users found.</p>
            </div>
        );
    }

    return (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {users.map((user) => (
                <UserCard
                    key={user.id}
                    user={user}
                    onEdit={onEdit}
                    onDelete={onDelete}
                />
            ))}
        </div>
    );
}
```

## Component Patterns

### With Variants (cva)

```tsx
import { cva, type VariantProps } from 'class-variance-authority';

const userCardVariants = cva(
    'rounded-lg border bg-card text-card-foreground shadow-sm',
    {
        variants: {
            variant: {
                default: 'p-4',
                compact: 'p-2',
                detailed: 'p-6',
            },
            elevated: {
                true: 'shadow-lg',
                false: '',
            },
        },
        defaultVariants: {
            variant: 'default',
            elevated: false,
        },
    }
);

interface UserCardProps
    extends React.HTMLAttributes<HTMLDivElement>,
        VariantProps<typeof userCardVariants> {
    user: User;
}
```

### With Form

```tsx
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';

const userFormSchema = z.object({
    name: z.string().min(1, 'Name is required').max(100),
    email: z.string().email('Invalid email address'),
    role: z.enum(['admin', 'user', 'guest']),
});

type UserFormValues = z.infer<typeof userFormSchema>;

export function UserForm({ onSubmit }: { onSubmit: (values: UserFormValues) => void }) {
    const form = useForm<UserFormValues>({
        resolver: zodResolver(userFormSchema),
        defaultValues: {
            name: '',
            email: '',
            role: 'user',
        },
    });

    return (
        <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                {/* Form fields */}
            </form>
        </Form>
    );
}
```

## Checklist

- [ ] Component with TypeScript props interface
- [ ] forwardRef for ref forwarding
- [ ] Variants using cva or conditional classes
- [ ] Design system colors and spacing only
- [ ] Loading, error, and empty states
- [ ] Accessibility (ARIA labels, keyboard navigation)
- [ ] Unit tests with Testing Library
- [ ] Associated hooks for data fetching
- [ ] Type definitions
