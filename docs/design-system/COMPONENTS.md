# Component Patterns

This document outlines the component patterns and best practices for building UI components using shadcn/ui with our design system.

## Tailwind CSS v4 + shadcn/ui

Our frontend uses Tailwind CSS v4 with CSS-first configuration and shadcn/ui components (new-york style).

### Key Dependencies

```json
{
  "dependencies": {
    "tailwindcss": "^4.0.0",
    "tw-animate-css": "^1.0.0",
    "@radix-ui/react-*": "latest"
  }
}
```

### CSS Configuration

```css
/* globals.css */
@import "tailwindcss";
@import "tw-animate-css";

@theme inline {
  --font-sans: 'Inter', ui-sans-serif, system-ui, sans-serif;
  --radius: 8px;
  /* ... other theme tokens */
}
```

## Component Architecture

### File Structure

```
src/components/
├── ui/                    # shadcn/ui base components
│   ├── button.tsx
│   ├── input.tsx
│   └── ...
├── features/              # Feature-specific components
│   ├── auth/
│   │   ├── login-form.tsx
│   │   └── signup-form.tsx
│   └── dashboard/
│       ├── stats-card.tsx
│       └── activity-feed.tsx
├── layouts/               # Layout components
│   ├── app-layout.tsx
│   ├── sidebar.tsx
│   └── header.tsx
└── shared/                # Shared composite components
    ├── data-table.tsx
    ├── page-header.tsx
    └── empty-state.tsx
```

### Naming Conventions

| Type | Convention | Example |
|------|------------|---------|
| Components | PascalCase | `UserProfile` |
| Files | kebab-case | `user-profile.tsx` |
| Props | PascalCase + Props | `UserProfileProps` |
| Hooks | camelCase with use | `useUserProfile` |
| Utilities | camelCase | `formatUserName` |

## Base Components (shadcn/ui)

### Button

```tsx
import { Button } from '@/components/ui/button';

// Variants
<Button variant="default">Primary</Button>
<Button variant="secondary">Secondary</Button>
<Button variant="outline">Outline</Button>
<Button variant="ghost">Ghost</Button>
<Button variant="link">Link</Button>
<Button variant="destructive">Destructive</Button>

// Sizes
<Button size="sm">Small</Button>
<Button size="default">Default</Button>
<Button size="lg">Large</Button>
<Button size="icon"><Icon /></Button>

// States
<Button disabled>Disabled</Button>
<Button loading>Loading</Button>

// With icons
<Button>
  <Plus className="mr-2 h-4 w-4" />
  Add Item
</Button>
```

### Input

```tsx
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';

<div className="grid gap-2">
  <Label htmlFor="email">Email</Label>
  <Input
    id="email"
    type="email"
    placeholder="name@example.com"
  />
</div>

// With error state
<div className="grid gap-2">
  <Label htmlFor="email">Email</Label>
  <Input
    id="email"
    type="email"
    aria-invalid="true"
    className="border-destructive"
  />
  <p className="text-xs text-destructive">
    Please enter a valid email address
  </p>
</div>
```

### Card

```tsx
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';

<Card>
  <CardHeader>
    <CardTitle>Card Title</CardTitle>
    <CardDescription>
      Card description text
    </CardDescription>
  </CardHeader>
  <CardContent>
    <p>Card content goes here</p>
  </CardContent>
  <CardFooter className="flex justify-between">
    <Button variant="outline">Cancel</Button>
    <Button>Save</Button>
  </CardFooter>
</Card>
```

### Dialog

```tsx
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';

<Dialog>
  <DialogTrigger asChild>
    <Button variant="outline">Open Dialog</Button>
  </DialogTrigger>
  <DialogContent className="sm:max-w-[425px]">
    <DialogHeader>
      <DialogTitle>Edit Profile</DialogTitle>
      <DialogDescription>
        Make changes to your profile here.
      </DialogDescription>
    </DialogHeader>
    <div className="grid gap-4 py-4">
      {/* Form fields */}
    </div>
    <DialogFooter>
      <Button type="submit">Save changes</Button>
    </DialogFooter>
  </DialogContent>
</Dialog>
```

## Composite Components

### Page Header

```tsx
interface PageHeaderProps {
  title: string;
  description?: string;
  actions?: React.ReactNode;
}

export function PageHeader({ title, description, actions }: PageHeaderProps) {
  return (
    <div className="flex items-center justify-between">
      <div className="space-y-1">
        <h1 className="text-2xl font-semibold tracking-tight">
          {title}
        </h1>
        {description && (
          <p className="text-sm text-muted-foreground">
            {description}
          </p>
        )}
      </div>
      {actions && (
        <div className="flex items-center gap-2">
          {actions}
        </div>
      )}
    </div>
  );
}

// Usage
<PageHeader
  title="Dashboard"
  description="Overview of your account"
  actions={
    <>
      <Button variant="outline">Export</Button>
      <Button>Add New</Button>
    </>
  }
/>
```

### Empty State

```tsx
interface EmptyStateProps {
  icon?: React.ReactNode;
  title: string;
  description: string;
  action?: React.ReactNode;
}

export function EmptyState({ icon, title, description, action }: EmptyStateProps) {
  return (
    <div className="flex flex-col items-center justify-center py-12 text-center">
      {icon && (
        <div className="mb-4 rounded-full bg-muted p-3">
          {icon}
        </div>
      )}
      <h3 className="text-lg font-medium">{title}</h3>
      <p className="mt-1 text-sm text-muted-foreground max-w-sm">
        {description}
      </p>
      {action && <div className="mt-4">{action}</div>}
    </div>
  );
}

// Usage
<EmptyState
  icon={<Inbox className="h-6 w-6 text-muted-foreground" />}
  title="No messages"
  description="You haven't received any messages yet. Start a conversation!"
  action={<Button>Compose Message</Button>}
/>
```

### Stats Card

```tsx
interface StatsCardProps {
  title: string;
  value: string | number;
  change?: {
    value: number;
    trend: 'up' | 'down';
  };
  icon?: React.ReactNode;
}

export function StatsCard({ title, value, change, icon }: StatsCardProps) {
  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between pb-2">
        <CardTitle className="text-sm font-medium text-muted-foreground">
          {title}
        </CardTitle>
        {icon}
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-bold">{value}</div>
        {change && (
          <p className={cn(
            "text-xs",
            change.trend === 'up' ? "text-success" : "text-destructive"
          )}>
            {change.trend === 'up' ? '+' : '-'}{Math.abs(change.value)}%
            <span className="text-muted-foreground ml-1">from last month</span>
          </p>
        )}
      </CardContent>
    </Card>
  );
}
```

## Form Patterns

### Form with React Hook Form + Zod

```tsx
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';

const formSchema = z.object({
  email: z.string().email('Invalid email address'),
  password: z.string().min(8, 'Password must be at least 8 characters'),
});

type FormValues = z.infer<typeof formSchema>;

export function LoginForm() {
  const form = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      email: '',
      password: '',
    },
  });

  const onSubmit = async (data: FormValues) => {
    // Handle form submission
  };

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
                <Input placeholder="name@example.com" {...field} />
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
                Must be at least 8 characters
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button type="submit" className="w-full">
          Sign In
        </Button>
      </form>
    </Form>
  );
}
```

## Data Display Patterns

### Data Table

```tsx
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';

interface DataTableProps<T> {
  columns: {
    key: keyof T;
    header: string;
    render?: (value: T[keyof T], row: T) => React.ReactNode;
  }[];
  data: T[];
  onRowClick?: (row: T) => void;
}

export function DataTable<T extends Record<string, unknown>>({
  columns,
  data,
  onRowClick,
}: DataTableProps<T>) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          {columns.map((column) => (
            <TableHead key={String(column.key)}>
              {column.header}
            </TableHead>
          ))}
        </TableRow>
      </TableHeader>
      <TableBody>
        {data.map((row, index) => (
          <TableRow
            key={index}
            onClick={() => onRowClick?.(row)}
            className={onRowClick ? 'cursor-pointer' : undefined}
          >
            {columns.map((column) => (
              <TableCell key={String(column.key)}>
                {column.render
                  ? column.render(row[column.key], row)
                  : String(row[column.key])}
              </TableCell>
            ))}
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
```

## Loading States

### Skeleton Loading

```tsx
import { Skeleton } from '@/components/ui/skeleton';

// Card skeleton
export function CardSkeleton() {
  return (
    <Card>
      <CardHeader>
        <Skeleton className="h-5 w-[200px]" />
        <Skeleton className="h-4 w-[300px]" />
      </CardHeader>
      <CardContent>
        <Skeleton className="h-20 w-full" />
      </CardContent>
    </Card>
  );
}

// Table skeleton
export function TableSkeleton({ rows = 5 }: { rows?: number }) {
  return (
    <div className="space-y-2">
      <Skeleton className="h-10 w-full" />
      {Array.from({ length: rows }).map((_, i) => (
        <Skeleton key={i} className="h-12 w-full" />
      ))}
    </div>
  );
}
```

### Loading Button

```tsx
interface LoadingButtonProps extends ButtonProps {
  loading?: boolean;
}

export function LoadingButton({
  loading,
  children,
  disabled,
  ...props
}: LoadingButtonProps) {
  return (
    <Button disabled={disabled || loading} {...props}>
      {loading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
      {children}
    </Button>
  );
}
```

## Accessibility Patterns

### Focus Management

```tsx
// Auto-focus first input in dialog
<DialogContent onOpenAutoFocus={(e) => {
  const firstInput = e.currentTarget.querySelector('input');
  firstInput?.focus();
}}>
  {/* ... */}
</DialogContent>
```

### Screen Reader Text

```tsx
<Button>
  <Trash className="h-4 w-4" />
  <span className="sr-only">Delete item</span>
</Button>
```

### ARIA Labels

```tsx
<Button
  aria-label="Close dialog"
  aria-describedby="dialog-description"
>
  <X className="h-4 w-4" />
</Button>
```

## Animation Patterns

### Entrance Animations

```tsx
// Fade in
<div className="animate-fade-in">Content</div>

// Slide in from bottom
<div className="animate-slide-in-bottom">Content</div>

// Scale in
<div className="animate-scale-in">Content</div>
```

### Transition Classes

```tsx
// Hover transition
<div className="transition-colors hover:bg-accent">
  Hover me
</div>

// Transform transition
<div className="transition-transform hover:scale-105">
  Scale on hover
</div>

// All transitions
<div className="transition-all duration-200 ease-in-out">
  Animated
</div>
```

## Responsive Patterns

### Responsive Grid

```tsx
<div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
  {items.map((item) => (
    <Card key={item.id}>{/* ... */}</Card>
  ))}
</div>
```

### Responsive Stack

```tsx
<div className="flex flex-col sm:flex-row gap-4">
  <div className="flex-1">{/* ... */}</div>
  <div className="flex-1">{/* ... */}</div>
</div>
```

### Mobile-First Navigation

```tsx
export function MobileNav() {
  return (
    <Sheet>
      <SheetTrigger asChild className="lg:hidden">
        <Button variant="ghost" size="icon">
          <Menu className="h-5 w-5" />
        </Button>
      </SheetTrigger>
      <SheetContent side="left">
        <nav className="flex flex-col gap-4">
          {/* Navigation items */}
        </nav>
      </SheetContent>
    </Sheet>
  );
}
```

## Best Practices

### Do's

1. **Use semantic HTML elements**
2. **Maintain consistent spacing** using design tokens
3. **Include proper ARIA attributes** for accessibility
4. **Handle loading and error states** gracefully
5. **Use TypeScript** for type safety
6. **Follow the single responsibility principle**
7. **Extract reusable logic** into custom hooks

### Don'ts

1. **Don't use inline styles** - use Tailwind classes
2. **Don't create components for one-time use**
3. **Don't ignore accessibility requirements**
4. **Don't use arbitrary values** - stick to the design system
5. **Don't nest components too deeply**
6. **Don't forget to handle edge cases** (empty, error states)
