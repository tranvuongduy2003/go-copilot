---
description: Create a new React component with shadcn/ui following design system
---

# Create React Component

Create a new React component following the project's design system and best practices.

## Component Details

**Component Name**: {{componentName}}

**Description**: {{description}}

**Component Type**:
- [ ] UI Component (reusable, goes in `components/ui/`)
- [ ] Feature Component (feature-specific, goes in `components/features/`)
- [ ] Layout Component (goes in `components/layout/`)
- [ ] Page Component (goes in `pages/`)

## CRITICAL: Design System Rules

You MUST use the design system. No arbitrary colors, spacing, or typography.

### Colors (USE THESE)
```tsx
// Primary actions
className="bg-primary text-primary-foreground hover:bg-primary/90"

// Secondary actions
className="bg-secondary text-secondary-foreground"

// Destructive/danger
className="bg-destructive text-destructive-foreground"

// Muted/subtle
className="bg-muted text-muted-foreground"

// Success/warning states
className="text-success" // oklch(0.7 0.17 160)
className="text-warning" // oklch(0.8 0.15 85)
```

### Spacing (USE THESE)
```tsx
// Padding: p-1 (4px), p-2 (8px), p-3 (12px), p-4 (16px), p-6 (24px), p-8 (32px)
// Gaps: gap-2, gap-4, gap-6
// Margins: mt-2, mb-4, etc.
```

### NEVER USE
```tsx
// NO arbitrary colors
className="bg-purple-500 text-blue-600"
// NO arbitrary spacing
className="p-[13px] mt-[7px]"
```

## Implementation Steps

### 1. Create Component File

Location: `frontend/src/components/{{path}}/{{componentName}}.tsx`

Structure:
```tsx
import { forwardRef } from 'react';
import { cva, type VariantProps } from 'class-variance-authority';
import { cn } from '@/lib/utils';

const componentVariants = cva(
  'base-classes-here',
  {
    variants: {
      variant: {
        default: 'bg-primary text-primary-foreground',
        secondary: 'bg-secondary text-secondary-foreground',
      },
      size: {
        sm: 'text-sm px-2 py-1',
        default: 'px-4 py-2',
        lg: 'text-lg px-6 py-3',
      },
    },
    defaultVariants: {
      variant: 'default',
      size: 'default',
    },
  }
);

interface ComponentProps
  extends React.HTMLAttributes<HTMLElement>,
    VariantProps<typeof componentVariants> {
  // Additional props
}

export const Component = forwardRef<HTMLElement, ComponentProps>(
  ({ className, variant, size, ...props }, ref) => {
    return (
      <element
        ref={ref}
        className={cn(componentVariants({ variant, size }), className)}
        {...props}
      />
    );
  }
);

Component.displayName = 'Component';
```

### 2. Add Types (if needed)

Location: `frontend/src/types/{{typeName}}.ts`

### 3. Create Tests

Location: `frontend/src/components/{{path}}/{{componentName}}.test.tsx`

```tsx
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Component } from './component';

describe('Component', () => {
  it('renders correctly', () => {
    render(<Component>Content</Component>);
    expect(screen.getByText('Content')).toBeInTheDocument();
  });

  // Add more tests
});
```

### 4. Handle States

Implement all necessary states:
- Default/idle state
- Hover state
- Focus state (keyboard accessible)
- Active/pressed state
- Disabled state
- Loading state (if applicable)
- Error state (if applicable)

### 5. Accessibility

Ensure:
- Proper ARIA labels
- Keyboard navigation
- Focus management
- Screen reader support

## Using shadcn/ui

If extending a shadcn/ui component:
```bash
npx shadcn@latest add {{componentName}}
```

Then customize in `frontend/src/components/ui/{{componentName}}.tsx`

## Output

Provide:
1. The component file
2. Types file (if needed)
3. Test file
4. Usage examples
5. Props documentation
