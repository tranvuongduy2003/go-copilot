---
applyTo: "frontend/src/components/**/*.tsx"
---

# UI Component Development Instructions

These instructions apply to all React components in the components directory.

## CRITICAL: Design System Enforcement

Every component MUST use the defined design system. No exceptions.

### Color Palette (OKLCH)

| Token | CSS Variable | Tailwind Class | Usage |
|-------|--------------|----------------|-------|
| Primary | `--primary` | `bg-primary` | Main actions, links |
| Primary Dark | `--primary-dark` | Hover states |
| Secondary | `--secondary` | `bg-secondary` | Secondary actions |
| Destructive | `--destructive` | `bg-destructive` | Danger actions |
| Success | `--success` | `text-success` | Success states |
| Warning | `--warning` | `text-warning` | Warning states |
| Background | `--background` | `bg-background` | Page background |
| Card | `--card` | `bg-card` | Card surfaces |
| Muted | `--muted` | `bg-muted` | Subtle backgrounds |
| Border | `--border` | `border-border` | Borders |

```tsx
// CORRECT usage
<Button className="bg-primary text-primary-foreground hover:bg-primary/90">
<Card className="bg-card border-border">
<Badge className="bg-secondary text-secondary-foreground">

// WRONG - NEVER use arbitrary colors
<Button className="bg-purple-500 hover:bg-purple-600">
<Card className="bg-gray-100">
<Badge className="bg-blue-100 text-blue-800">
```

### Typography Scale

```tsx
// Headings
<h1 className="text-4xl font-bold tracking-tight">   // 36px
<h2 className="text-3xl font-semibold">               // 30px
<h3 className="text-2xl font-semibold">               // 24px
<h4 className="text-xl font-medium">                  // 20px

// Body text
<p className="text-base">                             // 16px
<p className="text-sm text-muted-foreground">         // 14px
<span className="text-xs">                            // 12px

// Code
<code className="font-mono text-sm">
```

### Spacing Scale

Always use standard spacing values:

```tsx
// Padding
<div className="p-1">   // 4px
<div className="p-2">   // 8px
<div className="p-3">   // 12px
<div className="p-4">   // 16px
<div className="p-6">   // 24px
<div className="p-8">   // 32px

// Gaps
<div className="gap-2"> // 8px
<div className="gap-4"> // 16px
<div className="gap-6"> // 24px

// Margins
<div className="mt-4 mb-6">
```

### Border Radius

```tsx
<div className="rounded-sm">   // 4px - small elements
<div className="rounded-md">   // 8px - buttons, inputs
<div className="rounded-lg">   // 12px - cards
<div className="rounded-xl">   // 16px - modals
<div className="rounded-2xl">  // 24px - large sections
<div className="rounded-full"> // pills, avatars
```

### Shadows

Use violet-tinted shadows for elevation:

```tsx
<div className="shadow-sm">   // Subtle elevation
<div className="shadow-md">   // Card elevation
<div className="shadow-lg">   // Modal/dropdown elevation
```

## Component Structure

### Basic Component Template

```tsx
import { forwardRef } from 'react';
import { cva, type VariantProps } from 'class-variance-authority';
import { cn } from '@/lib/utils';

const componentVariants = cva(
  // Base styles
  'inline-flex items-center justify-center rounded-md font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring',
  {
    variants: {
      variant: {
        default: 'bg-primary text-primary-foreground hover:bg-primary/90',
        secondary: 'bg-secondary text-secondary-foreground hover:bg-secondary/80',
        outline: 'border border-input bg-background hover:bg-accent',
        ghost: 'hover:bg-accent hover:text-accent-foreground',
        destructive: 'bg-destructive text-destructive-foreground hover:bg-destructive/90',
      },
      size: {
        sm: 'h-9 px-3 text-sm',
        default: 'h-10 px-4',
        lg: 'h-11 px-6 text-lg',
        icon: 'h-10 w-10',
      },
    },
    defaultVariants: {
      variant: 'default',
      size: 'default',
    },
  }
);

interface ComponentProps
  extends React.HTMLAttributes<HTMLDivElement>,
    VariantProps<typeof componentVariants> {
  isLoading?: boolean;
}

export const Component = forwardRef<HTMLDivElement, ComponentProps>(
  ({ className, variant, size, isLoading, children, ...props }, ref) => {
    return (
      <div
        ref={ref}
        className={cn(componentVariants({ variant, size }), className)}
        {...props}
      >
        {isLoading ? <Spinner /> : children}
      </div>
    );
  }
);

Component.displayName = 'Component';
```

## shadcn/ui Usage

### Installing Components

```bash
npx shadcn@latest add button
npx shadcn@latest add card
npx shadcn@latest add input
npx shadcn@latest add dialog
```

### Customizing shadcn/ui Components

```tsx
// Extend existing components with variants
import { Button, buttonVariants } from '@/components/ui/button';

// Add custom variant
const customButtonVariants = cva(
  buttonVariants(), // Include base styles
  {
    variants: {
      gradient: {
        true: 'bg-gradient-to-r from-primary to-primary-dark hover:opacity-90',
      },
    },
  }
);

// Or compose with existing
export function GradientButton({ className, ...props }: ButtonProps) {
  return (
    <Button
      className={cn(
        'bg-gradient-to-r from-primary to-primary-dark',
        className
      )}
      {...props}
    />
  );
}
```

### Common Component Patterns

#### Card Component

```tsx
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';

export function FeatureCard({ title, description, children, actions }) {
  return (
    <Card className="hover:shadow-md transition-shadow">
      <CardHeader>
        <CardTitle className="text-xl">{title}</CardTitle>
        <CardDescription>{description}</CardDescription>
      </CardHeader>
      <CardContent>{children}</CardContent>
      {actions && (
        <CardFooter className="flex justify-end gap-2">
          {actions}
        </CardFooter>
      )}
    </Card>
  );
}
```

#### Form Field

```tsx
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';

<FormField
  control={form.control}
  name="email"
  render={({ field }) => (
    <FormItem>
      <FormLabel>Email</FormLabel>
      <FormControl>
        <Input
          type="email"
          placeholder="you@example.com"
          {...field}
        />
      </FormControl>
      <FormMessage />
    </FormItem>
  )}
/>
```

#### Dialog/Modal

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

export function ConfirmDialog({
  trigger,
  title,
  description,
  onConfirm,
  confirmText = 'Confirm',
  destructive = false,
}) {
  const [open, setOpen] = useState(false);

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>{trigger}</DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{title}</DialogTitle>
          <DialogDescription>{description}</DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="outline" onClick={() => setOpen(false)}>
            Cancel
          </Button>
          <Button
            variant={destructive ? 'destructive' : 'default'}
            onClick={() => {
              onConfirm();
              setOpen(false);
            }}
          >
            {confirmText}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
```

## Accessibility Requirements

### Focus States

```tsx
// All interactive elements need visible focus states
<button className="focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2">

// Use focus-visible (not focus) to avoid focus on click
```

### ARIA Labels

```tsx
// Icon buttons need labels
<Button variant="ghost" size="icon" aria-label="Close menu">
  <X className="h-4 w-4" />
</Button>

// Loading states
<Button disabled aria-busy={isLoading}>
  {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
  {isLoading ? 'Loading...' : 'Submit'}
</Button>

// Expandable sections
<Button
  aria-expanded={isOpen}
  aria-controls="menu-content"
  onClick={() => setIsOpen(!isOpen)}
>
  Menu
</Button>
<div id="menu-content" hidden={!isOpen}>
  {/* content */}
</div>
```

### Keyboard Navigation

```tsx
// Ensure keyboard operability
<div
  role="button"
  tabIndex={0}
  onClick={handleClick}
  onKeyDown={(e) => {
    if (e.key === 'Enter' || e.key === ' ') {
      e.preventDefault();
      handleClick();
    }
  }}
>
  Clickable div (prefer <button> when possible)
</div>
```

## Responsive Design

### Breakpoint Usage

```tsx
// Mobile-first approach
<div className="
  p-4
  md:p-6
  lg:p-8
">

// Grid layouts
<div className="
  grid
  grid-cols-1
  sm:grid-cols-2
  lg:grid-cols-3
  xl:grid-cols-4
  gap-4
">

// Hidden/shown at breakpoints
<div className="hidden md:block">Desktop only</div>
<div className="md:hidden">Mobile only</div>
```

### Container Patterns

```tsx
// Page container
<div className="container mx-auto px-4 py-8">

// Max width constraints
<div className="max-w-md mx-auto">   // Forms
<div className="max-w-2xl mx-auto">  // Content
<div className="max-w-7xl mx-auto">  // Wide layouts
```

## Animation

### Use Tailwind Animations

```tsx
// Spin (loading indicators)
<Loader2 className="h-4 w-4 animate-spin" />

// Pulse (skeleton loaders)
<div className="animate-pulse bg-muted h-4 w-full rounded" />

// Transition for state changes
<div className="transition-colors hover:bg-accent">
<div className="transition-transform hover:scale-105">
<div className="transition-opacity opacity-0 data-[state=open]:opacity-100">
```

### Custom Animations (in globals.css)

```css
@keyframes fade-in {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes slide-up {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
```

## Dark Mode Support

### Theme-Aware Classes

```tsx
// Colors automatically switch in dark mode via CSS variables
<div className="bg-background text-foreground">
<div className="bg-card border-border">

// If needed, explicit dark mode overrides
<div className="bg-white dark:bg-zinc-900">
```

## Testing Components

```tsx
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Button } from './button';

describe('Button', () => {
  it('renders children', () => {
    render(<Button>Click me</Button>);
    expect(screen.getByRole('button')).toHaveTextContent('Click me');
  });

  it('shows loading state', () => {
    render(<Button isLoading>Submit</Button>);
    expect(screen.getByRole('button')).toBeDisabled();
    expect(screen.getByText(/loading/i)).toBeInTheDocument();
  });

  it('calls onClick when clicked', async () => {
    const onClick = vi.fn();
    render(<Button onClick={onClick}>Click</Button>);

    await userEvent.click(screen.getByRole('button'));

    expect(onClick).toHaveBeenCalledTimes(1);
  });

  it('is accessible', async () => {
    const { container } = render(<Button>Accessible Button</Button>);
    const results = await axe(container);
    expect(results).toHaveNoViolations();
  });
});
```
