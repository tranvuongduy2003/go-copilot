---
name: react-component-builder
description: Create React components with shadcn/ui, proper typing, and tests. Use when building new UI components.
---

# React Component Builder Skill

This skill guides you through building production-ready React components following the project's design system and best practices.

## When to Use This Skill

- Creating a new UI component
- Building a feature component
- Creating form components
- Implementing data display components

## CRITICAL: Design System Compliance

Every component MUST use the design system. No arbitrary colors, spacing, or typography.

### Colors (MUST USE)
```tsx
// Primary
className="bg-primary text-primary-foreground"
className="hover:bg-primary/90"

// Secondary
className="bg-secondary text-secondary-foreground"

// Destructive
className="bg-destructive text-destructive-foreground"

// Muted
className="bg-muted text-muted-foreground"

// Card/Background
className="bg-card bg-background"

// NEVER USE
className="bg-purple-500 bg-blue-600 bg-gray-100"
```

### Spacing (MUST USE)
```tsx
// Padding: p-1, p-2, p-3, p-4, p-6, p-8
// Margin: m-1, m-2, m-3, m-4, m-6, m-8
// Gap: gap-2, gap-4, gap-6

// NEVER USE
className="p-[13px] mt-[7px]"
```

## Component Templates

### Template 1: Basic Component

```tsx
// components/ui/status-badge.tsx
import { cva, type VariantProps } from 'class-variance-authority';
import { cn } from '@/lib/utils';

const statusBadgeVariants = cva(
  'inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium transition-colors',
  {
    variants: {
      status: {
        pending: 'bg-warning/10 text-warning border border-warning/20',
        active: 'bg-success/10 text-success border border-success/20',
        inactive: 'bg-muted text-muted-foreground border border-border',
        error: 'bg-destructive/10 text-destructive border border-destructive/20',
      },
    },
    defaultVariants: {
      status: 'pending',
    },
  }
);

interface StatusBadgeProps
  extends React.HTMLAttributes<HTMLSpanElement>,
    VariantProps<typeof statusBadgeVariants> {
  label?: string;
}

export function StatusBadge({
  status,
  label,
  className,
  ...props
}: StatusBadgeProps) {
  const displayLabel = label ?? status;

  return (
    <span
      className={cn(statusBadgeVariants({ status }), className)}
      {...props}
    >
      {displayLabel}
    </span>
  );
}
```

### Template 2: Card Component

```tsx
// components/features/product/product-card.tsx
import { Card, CardContent, CardFooter, CardHeader } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { ShoppingCart } from 'lucide-react';
import { cn } from '@/lib/utils';
import type { Product } from '@/types/product';

interface ProductCardProps {
  product: Product;
  onAddToCart?: (product: Product) => void;
  className?: string;
}

export function ProductCard({ product, onAddToCart, className }: ProductCardProps) {
  const formattedPrice = new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
  }).format(product.price / 100);

  return (
    <Card className={cn('group hover:shadow-lg transition-shadow', className)}>
      <CardHeader className="p-0">
        <div className="aspect-square overflow-hidden rounded-t-lg bg-muted">
          <img
            src={product.imageUrl}
            alt={product.name}
            className="h-full w-full object-cover transition-transform group-hover:scale-105"
          />
        </div>
      </CardHeader>
      <CardContent className="p-4">
        <div className="flex items-start justify-between gap-2">
          <div>
            <h3 className="font-semibold text-lg line-clamp-1">{product.name}</h3>
            <p className="text-sm text-muted-foreground line-clamp-2 mt-1">
              {product.description}
            </p>
          </div>
          <Badge variant="secondary">{product.category}</Badge>
        </div>
        <p className="text-2xl font-bold text-primary mt-4">{formattedPrice}</p>
      </CardContent>
      <CardFooter className="p-4 pt-0">
        <Button
          className="w-full"
          onClick={() => onAddToCart?.(product)}
          disabled={product.stock === 0}
        >
          <ShoppingCart className="mr-2 h-4 w-4" />
          {product.stock === 0 ? 'Out of Stock' : 'Add to Cart'}
        </Button>
      </CardFooter>
    </Card>
  );
}
```

### Template 3: Form Component

```tsx
// components/features/auth/register-form.tsx
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
  FormDescription,
} from '@/components/ui/form';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { AlertCircle } from 'lucide-react';

const registerSchema = z.object({
  name: z.string().min(2, 'Name must be at least 2 characters'),
  email: z.string().email('Invalid email address'),
  password: z
    .string()
    .min(8, 'Password must be at least 8 characters')
    .regex(/[A-Z]/, 'Password must contain an uppercase letter')
    .regex(/[0-9]/, 'Password must contain a number'),
  confirmPassword: z.string(),
}).refine((data) => data.password === data.confirmPassword, {
  message: "Passwords don't match",
  path: ['confirmPassword'],
});

type RegisterFormValues = z.infer<typeof registerSchema>;

interface RegisterFormProps {
  onSubmit: (values: RegisterFormValues) => Promise<void>;
}

export function RegisterForm({ onSubmit }: RegisterFormProps) {
  const form = useForm<RegisterFormValues>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      name: '',
      email: '',
      password: '',
      confirmPassword: '',
    },
  });

  const handleSubmit = async (values: RegisterFormValues) => {
    try {
      await onSubmit(values);
    } catch (error) {
      form.setError('root', {
        message: error instanceof Error ? error.message : 'Registration failed',
      });
    }
  };

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-4">
        {form.formState.errors.root && (
          <Alert variant="destructive">
            <AlertCircle className="h-4 w-4" />
            <AlertDescription>
              {form.formState.errors.root.message}
            </AlertDescription>
          </Alert>
        )}

        <FormField
          control={form.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Name</FormLabel>
              <FormControl>
                <Input placeholder="John Doe" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="email"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Email</FormLabel>
              <FormControl>
                <Input type="email" placeholder="you@example.com" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="password"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Password</FormLabel>
              <FormControl>
                <Input type="password" {...field} />
              </FormControl>
              <FormDescription>
                At least 8 characters with uppercase and number
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="confirmPassword"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Confirm Password</FormLabel>
              <FormControl>
                <Input type="password" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <Button
          type="submit"
          className="w-full"
          disabled={form.formState.isSubmitting}
        >
          {form.formState.isSubmitting ? 'Creating account...' : 'Create Account'}
        </Button>
      </form>
    </Form>
  );
}
```

### Template 4: Data List Component

```tsx
// components/features/users/user-list.tsx
import { useUsers } from '@/hooks/use-users';
import { UserCard } from './user-card';
import { UserListSkeleton } from './user-list-skeleton';
import { EmptyState } from '@/components/ui/empty-state';
import { ErrorMessage } from '@/components/ui/error-message';
import { Users } from 'lucide-react';

interface UserListProps {
  onUserClick?: (userId: string) => void;
}

export function UserList({ onUserClick }: UserListProps) {
  const { data: users, isLoading, error, refetch } = useUsers();

  if (isLoading) {
    return <UserListSkeleton />;
  }

  if (error) {
    return (
      <ErrorMessage
        title="Failed to load users"
        message={error.message}
        onRetry={refetch}
      />
    );
  }

  if (!users?.length) {
    return (
      <EmptyState
        icon={Users}
        title="No users found"
        description="There are no users to display."
      />
    );
  }

  return (
    <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      {users.map((user) => (
        <UserCard
          key={user.id}
          user={user}
          onClick={() => onUserClick?.(user.id)}
        />
      ))}
    </div>
  );
}
```

### Template 5: Custom Hook

```tsx
// hooks/use-products.ts
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '@/lib/api';
import type { Product, CreateProductInput } from '@/types/product';

export const productKeys = {
  all: ['products'] as const,
  lists: () => [...productKeys.all, 'list'] as const,
  list: (filters: Record<string, unknown>) =>
    [...productKeys.lists(), filters] as const,
  details: () => [...productKeys.all, 'detail'] as const,
  detail: (id: string) => [...productKeys.details(), id] as const,
};

export function useProducts(filters?: { category?: string; page?: number }) {
  return useQuery({
    queryKey: productKeys.list(filters ?? {}),
    queryFn: () => api.get<Product[]>('/products', { params: filters }),
  });
}

export function useProduct(id: string) {
  return useQuery({
    queryKey: productKeys.detail(id),
    queryFn: () => api.get<Product>(`/products/${id}`),
    enabled: !!id,
  });
}

export function useCreateProduct() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: CreateProductInput) =>
      api.post<Product>('/products', input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: productKeys.lists() });
    },
  });
}
```

## Testing Component

```tsx
// components/features/product/product-card.test.tsx
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { vi } from 'vitest';
import { ProductCard } from './product-card';

const mockProduct = {
  id: '1',
  name: 'Test Product',
  description: 'A test product description',
  price: 2999,
  category: 'Electronics',
  stock: 10,
  imageUrl: '/test-image.jpg',
};

describe('ProductCard', () => {
  it('renders product information', () => {
    render(<ProductCard product={mockProduct} />);

    expect(screen.getByText('Test Product')).toBeInTheDocument();
    expect(screen.getByText('A test product description')).toBeInTheDocument();
    expect(screen.getByText('$29.99')).toBeInTheDocument();
    expect(screen.getByText('Electronics')).toBeInTheDocument();
  });

  it('calls onAddToCart when button is clicked', async () => {
    const user = userEvent.setup();
    const onAddToCart = vi.fn();

    render(<ProductCard product={mockProduct} onAddToCart={onAddToCart} />);

    await user.click(screen.getByRole('button', { name: /add to cart/i }));

    expect(onAddToCart).toHaveBeenCalledWith(mockProduct);
  });

  it('shows out of stock when stock is 0', () => {
    render(<ProductCard product={{ ...mockProduct, stock: 0 }} />);

    expect(screen.getByRole('button', { name: /out of stock/i })).toBeDisabled();
  });
});
```

## Component Checklist

- [ ] TypeScript interfaces defined for all props
- [ ] Uses design system colors (no arbitrary values)
- [ ] Uses design system spacing (no arbitrary values)
- [ ] Responsive design implemented
- [ ] Loading state handled
- [ ] Error state handled
- [ ] Empty state handled
- [ ] Accessible (proper ARIA labels, keyboard navigation)
- [ ] Tests written
- [ ] Uses forwardRef if needed for refs
- [ ] className prop passed through with cn()
