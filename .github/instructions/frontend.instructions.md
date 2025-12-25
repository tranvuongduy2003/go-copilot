---
applyTo: "frontend/src/**/*.{ts,tsx}"
---

# Frontend Development Instructions

These instructions apply to all TypeScript and React files in the frontend source directory.

## Project Structure

```
frontend/src/
├── components/
│   ├── ui/                    # shadcn/ui components
│   ├── features/              # Feature-specific components
│   └── layout/                # Layout components
├── hooks/                     # Custom React hooks
├── lib/
│   ├── api.ts                 # API client
│   ├── utils.ts               # Utility functions
│   └── validations.ts         # Zod schemas
├── pages/                     # Page components
├── stores/                    # Zustand stores
├── styles/
│   └── globals.css            # Global styles with Tailwind
└── types/                     # TypeScript type definitions
```

## TypeScript Standards

### Strict Mode Required
```json
{
  "compilerOptions": {
    "strict": true,
    "noImplicitAny": true,
    "strictNullChecks": true
  }
}
```

### Type Definitions
```typescript
// Prefer interface for object shapes
interface User {
  id: string;
  email: string;
  name: string;
  createdAt: string;
}

// Use type for unions, intersections, primitives
type Status = 'pending' | 'active' | 'completed';
type UserWithRole = User & { role: Role };

// Never use 'any' - use 'unknown' if type is truly unknown
function handleResponse(data: unknown): User {
  // Validate and narrow the type
  if (isUser(data)) {
    return data;
  }
  throw new Error('Invalid response');
}
```

### File Naming
```
components/
├── ui/
│   └── button.tsx           # kebab-case for components
├── features/
│   └── user-profile/
│       ├── user-profile.tsx
│       └── user-avatar.tsx
hooks/
├── use-user.ts              # camelCase with use- prefix
├── use-auth.ts
lib/
├── api.ts                   # camelCase for utilities
├── utils.ts
types/
├── user.ts                  # lowercase for type files
```

## React Patterns

### Component Structure
```tsx
import { forwardRef } from 'react';
import { cn } from '@/lib/utils';

// Props interface with explicit types
interface CardProps extends React.HTMLAttributes<HTMLDivElement> {
  variant?: 'default' | 'outlined';
  isLoading?: boolean;
}

// Named export (not default)
export const Card = forwardRef<HTMLDivElement, CardProps>(
  ({ className, variant = 'default', isLoading, children, ...props }, ref) => {
    if (isLoading) {
      return <CardSkeleton />;
    }

    return (
      <div
        ref={ref}
        className={cn(
          'rounded-lg border bg-card p-6',
          variant === 'outlined' && 'border-2',
          className
        )}
        {...props}
      >
        {children}
      </div>
    );
  }
);

Card.displayName = 'Card';
```

### Custom Hooks
```tsx
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '@/lib/api';
import type { User, CreateUserInput } from '@/types/user';

// Query keys as constants
export const userKeys = {
  all: ['users'] as const,
  detail: (id: string) => ['users', id] as const,
};

// Data fetching hook
export function useUser(id: string) {
  return useQuery({
    queryKey: userKeys.detail(id),
    queryFn: () => api.get<User>(`/users/${id}`),
    enabled: !!id,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}

// Mutation hook
export function useCreateUser() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: CreateUserInput) =>
      api.post<User>('/users', input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: userKeys.all });
    },
  });
}
```

### State Management (Zustand)
```tsx
import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  login: (user: User, token: string) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      token: null,
      isAuthenticated: false,
      login: (user, token) =>
        set({ user, token, isAuthenticated: true }),
      logout: () =>
        set({ user: null, token: null, isAuthenticated: false }),
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({ token: state.token }),
    }
  )
);
```

## Design System Compliance

### CRITICAL: Use Design System Colors Only

```tsx
// CORRECT: Design system colors
<div className="bg-primary text-primary-foreground">
<div className="bg-secondary">
<div className="text-muted-foreground">
<div className="border-border">
<div className="bg-destructive text-destructive-foreground">

// WRONG: Arbitrary colors - NEVER USE
<div className="bg-purple-500">
<div className="text-blue-600">
<div className="bg-[#7c3aed]">
```

### Use Spacing Scale
```tsx
// CORRECT: Standard spacing
<div className="p-4 mt-6 gap-4">
<div className="px-6 py-3">

// WRONG: Arbitrary spacing
<div className="p-[13px] mt-[7px]">
```

### Use Border Radius Scale
```tsx
// CORRECT: Standard radius
<div className="rounded-sm">   // 4px
<div className="rounded-md">   // 8px
<div className="rounded-lg">   // 12px
<div className="rounded-xl">   // 16px
<div className="rounded-full"> // pill

// WRONG: Arbitrary radius
<div className="rounded-[5px]">
```

### Typography
```tsx
// CORRECT: Use type scale
<h1 className="text-4xl font-bold tracking-tight">
<p className="text-base text-muted-foreground">
<span className="text-sm font-medium">

// Font families
<p className="font-sans">Inter text</p>
<code className="font-mono">JetBrains Mono</code>
```

## Form Handling

### React Hook Form with Zod
```tsx
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';

const formSchema = z.object({
  email: z.string().email('Invalid email address'),
  password: z.string().min(8, 'Password must be at least 8 characters'),
});

type FormValues = z.infer<typeof formSchema>;

export function LoginForm({ onSubmit }: { onSubmit: (values: FormValues) => Promise<void> }) {
  const form = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      email: '',
      password: '',
    },
  });

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
        <FormField
          control={form.control}
          name="email"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Email</FormLabel>
              <FormControl>
                <Input type="email" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button type="submit" isLoading={form.formState.isSubmitting}>
          Sign In
        </Button>
      </form>
    </Form>
  );
}
```

## Error Handling

### API Error Handling
```tsx
// API client with error handling
class ApiError extends Error {
  constructor(
    public status: number,
    public code: string,
    message: string,
    public details?: unknown
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

async function request<T>(endpoint: string, options?: RequestInit): Promise<T> {
  const response = await fetch(`${API_URL}${endpoint}`, options);

  if (!response.ok) {
    const error = await response.json();
    throw new ApiError(
      response.status,
      error.error?.code ?? 'UNKNOWN',
      error.error?.message ?? 'An error occurred',
      error.error?.details
    );
  }

  const data = await response.json();
  return data.data;
}
```

### Component Error States
```tsx
function UserProfile({ id }: { id: string }) {
  const { data: user, isLoading, error } = useUser(id);

  if (isLoading) {
    return <ProfileSkeleton />;
  }

  if (error) {
    return (
      <ErrorMessage
        title="Failed to load profile"
        message={error.message}
        onRetry={() => refetch()}
      />
    );
  }

  return <ProfileCard user={user} />;
}
```

## Accessibility

### Required Practices
```tsx
// All interactive elements need accessible names
<button aria-label="Close dialog">
  <X className="h-4 w-4" />
</button>

// Form inputs need labels
<label htmlFor="email">Email</label>
<input id="email" type="email" />

// Images need alt text
<img src={user.avatar} alt={`${user.name}'s avatar`} />

// Use semantic HTML
<nav aria-label="Main navigation">
<main>
<article>
<aside>

// Keyboard navigation
<div
  role="button"
  tabIndex={0}
  onClick={handleClick}
  onKeyDown={(e) => e.key === 'Enter' && handleClick()}
>
```

## Performance

### Memoization
```tsx
// Memoize expensive computations
const sortedItems = useMemo(
  () => [...items].sort((a, b) => a.name.localeCompare(b.name)),
  [items]
);

// Memoize callbacks passed to children
const handleClick = useCallback((id: string) => {
  setSelectedId(id);
}, []);

// Memoize components that receive objects/arrays
const MemoizedList = memo(function List({ items }: { items: Item[] }) {
  return items.map(item => <ListItem key={item.id} item={item} />);
});
```

### Code Splitting
```tsx
import { lazy, Suspense } from 'react';

// Lazy load routes
const Dashboard = lazy(() => import('@/pages/dashboard'));
const Settings = lazy(() => import('@/pages/settings'));

function App() {
  return (
    <Suspense fallback={<PageLoader />}>
      <Routes>
        <Route path="/dashboard" element={<Dashboard />} />
        <Route path="/settings" element={<Settings />} />
      </Routes>
    </Suspense>
  );
}
```

## Imports

### Use Absolute Imports
```tsx
// CORRECT: Absolute imports with @ alias
import { Button } from '@/components/ui/button';
import { useUser } from '@/hooks/use-user';
import { cn } from '@/lib/utils';

// WRONG: Relative imports (avoid deep paths)
import { Button } from '../../../components/ui/button';
```

### Import Order
```tsx
// 1. React
import { useState, useEffect } from 'react';

// 2. External libraries
import { useQuery } from '@tanstack/react-query';
import { z } from 'zod';

// 3. Internal aliases (@/)
import { Button } from '@/components/ui/button';
import { useAuth } from '@/hooks/use-auth';

// 4. Relative imports
import { ProfileCard } from './profile-card';

// 5. Types (if separate)
import type { User } from '@/types/user';
```
