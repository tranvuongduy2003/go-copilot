---
name: shadcn-theming
description: Configure and customize shadcn/ui themes. Use for theming and styling tasks.
---

# shadcn/ui Theming Skill

This skill guides you through configuring and customizing the shadcn/ui theme to match the project's design system.

## Project Design System

### Color Palette (OKLCH)

| Token | Light Mode | Dark Mode | Usage |
|-------|------------|-----------|-------|
| Primary | `oklch(0.7 0.15 290)` | `oklch(0.7 0.15 290)` | Main actions |
| Primary Dark | `oklch(0.6 0.2 280)` | `oklch(0.6 0.2 280)` | Hover states |
| Secondary | `oklch(0.75 0.15 220)` | `oklch(0.75 0.15 220)` | Secondary actions |
| Background | `oklch(0.98 0.01 260)` | `oklch(0.1 0.01 260)` | Page background |
| Card | `oklch(1 0 0)` | `oklch(0.15 0.01 260)` | Card surfaces |
| Success | `oklch(0.7 0.17 160)` | `oklch(0.7 0.17 160)` | Success states |
| Warning | `oklch(0.8 0.15 85)` | `oklch(0.8 0.15 85)` | Warning states |
| Error/Destructive | `oklch(0.65 0.2 15)` | `oklch(0.65 0.2 15)` | Error states |

## Complete globals.css Configuration

```css
/* frontend/src/styles/globals.css */

@import "tailwindcss";

/* Custom fonts */
@font-face {
  font-family: 'Inter';
  font-style: normal;
  font-weight: 100 900;
  font-display: swap;
  src: url('/fonts/Inter-Variable.woff2') format('woff2');
}

@font-face {
  font-family: 'JetBrains Mono';
  font-style: normal;
  font-weight: 100 800;
  font-display: swap;
  src: url('/fonts/JetBrainsMono-Variable.woff2') format('woff2');
}

/* Design System Theme */
@theme inline {
  /* Typography */
  --font-sans: 'Inter', ui-sans-serif, system-ui, sans-serif;
  --font-mono: 'JetBrains Mono', ui-monospace, monospace;

  /* Spacing Scale */
  --spacing-0: 0px;
  --spacing-1: 4px;
  --spacing-2: 8px;
  --spacing-3: 12px;
  --spacing-4: 16px;
  --spacing-5: 20px;
  --spacing-6: 24px;
  --spacing-8: 32px;
  --spacing-10: 40px;
  --spacing-12: 48px;
  --spacing-16: 64px;
  --spacing-20: 80px;
  --spacing-24: 96px;
  --spacing-32: 128px;

  /* Border Radius */
  --radius-sm: 4px;
  --radius-md: 8px;
  --radius-lg: 12px;
  --radius-xl: 16px;
  --radius-2xl: 24px;
  --radius-full: 9999px;

  /* Shadows with violet tint */
  --shadow-sm: 0 1px 2px oklch(0.3 0.05 290 / 0.05);
  --shadow-md: 0 4px 6px -1px oklch(0.3 0.05 290 / 0.1), 0 2px 4px -2px oklch(0.3 0.05 290 / 0.1);
  --shadow-lg: 0 10px 15px -3px oklch(0.3 0.05 290 / 0.1), 0 4px 6px -4px oklch(0.3 0.05 290 / 0.1);
  --shadow-xl: 0 20px 25px -5px oklch(0.3 0.05 290 / 0.1), 0 8px 10px -6px oklch(0.3 0.05 290 / 0.1);

  /* Animations */
  --animate-accordion-down: accordion-down 0.2s ease-out;
  --animate-accordion-up: accordion-up 0.2s ease-out;
  --animate-fade-in: fade-in 0.2s ease-out;
  --animate-fade-out: fade-out 0.2s ease-out;
  --animate-slide-in-from-top: slide-in-from-top 0.2s ease-out;
  --animate-slide-in-from-bottom: slide-in-from-bottom 0.2s ease-out;
}

/* Light Mode (Default) */
:root {
  /* Primary - Violet/Purple */
  --primary: oklch(0.7 0.15 290);
  --primary-dark: oklch(0.6 0.2 280);
  --primary-foreground: oklch(0.98 0.01 290);

  /* Secondary - Cyan/Blue */
  --secondary: oklch(0.75 0.15 220);
  --secondary-foreground: oklch(0.15 0.02 220);

  /* Background & Surfaces */
  --background: oklch(0.98 0.01 260);
  --foreground: oklch(0.15 0.02 260);
  --card: oklch(1 0 0);
  --card-foreground: oklch(0.15 0.02 260);
  --popover: oklch(1 0 0);
  --popover-foreground: oklch(0.15 0.02 260);

  /* Muted */
  --muted: oklch(0.95 0.01 260);
  --muted-foreground: oklch(0.45 0.02 260);

  /* Accent */
  --accent: oklch(0.95 0.01 260);
  --accent-foreground: oklch(0.15 0.02 260);

  /* Semantic Colors */
  --success: oklch(0.7 0.17 160);
  --success-foreground: oklch(0.98 0.01 160);
  --warning: oklch(0.8 0.15 85);
  --warning-foreground: oklch(0.2 0.05 85);
  --destructive: oklch(0.65 0.2 15);
  --destructive-foreground: oklch(0.98 0.01 15);

  /* Borders & Input */
  --border: oklch(0.9 0.01 260);
  --input: oklch(0.9 0.01 260);
  --ring: oklch(0.7 0.15 290);

  /* Chart Colors */
  --chart-1: oklch(0.7 0.15 290);
  --chart-2: oklch(0.75 0.15 220);
  --chart-3: oklch(0.7 0.17 160);
  --chart-4: oklch(0.8 0.15 85);
  --chart-5: oklch(0.65 0.2 15);

  /* Radius */
  --radius: 8px;
}

/* Dark Mode */
.dark {
  /* Primary - Violet/Purple (same hue, adjusted lightness) */
  --primary: oklch(0.7 0.15 290);
  --primary-dark: oklch(0.6 0.2 280);
  --primary-foreground: oklch(0.15 0.02 290);

  /* Secondary - Cyan/Blue */
  --secondary: oklch(0.35 0.08 220);
  --secondary-foreground: oklch(0.9 0.02 220);

  /* Background & Surfaces */
  --background: oklch(0.1 0.01 260);
  --foreground: oklch(0.95 0.01 260);
  --card: oklch(0.15 0.01 260);
  --card-foreground: oklch(0.95 0.01 260);
  --popover: oklch(0.15 0.01 260);
  --popover-foreground: oklch(0.95 0.01 260);

  /* Muted */
  --muted: oklch(0.2 0.01 260);
  --muted-foreground: oklch(0.6 0.02 260);

  /* Accent */
  --accent: oklch(0.2 0.01 260);
  --accent-foreground: oklch(0.95 0.01 260);

  /* Borders & Input */
  --border: oklch(0.25 0.01 260);
  --input: oklch(0.25 0.01 260);
  --ring: oklch(0.7 0.15 290);
}

/* Keyframe Animations */
@keyframes accordion-down {
  from { height: 0; }
  to { height: var(--radix-accordion-content-height); }
}

@keyframes accordion-up {
  from { height: var(--radix-accordion-content-height); }
  to { height: 0; }
}

@keyframes fade-in {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes fade-out {
  from { opacity: 1; }
  to { opacity: 0; }
}

@keyframes slide-in-from-top {
  from { transform: translateY(-10px); opacity: 0; }
  to { transform: translateY(0); opacity: 1; }
}

@keyframes slide-in-from-bottom {
  from { transform: translateY(10px); opacity: 0; }
  to { transform: translateY(0); opacity: 1; }
}

/* Base Styles */
@layer base {
  * {
    @apply border-border;
  }

  body {
    @apply bg-background text-foreground;
    font-feature-settings: "rlig" 1, "calt" 1;
  }

  /* Focus styles */
  *:focus-visible {
    @apply outline-none ring-2 ring-ring ring-offset-2 ring-offset-background;
  }

  /* Selection */
  ::selection {
    @apply bg-primary/20 text-foreground;
  }

  /* Scrollbar */
  ::-webkit-scrollbar {
    @apply w-2 h-2;
  }

  ::-webkit-scrollbar-track {
    @apply bg-transparent;
  }

  ::-webkit-scrollbar-thumb {
    @apply bg-border rounded-full;
  }

  ::-webkit-scrollbar-thumb:hover {
    @apply bg-muted-foreground/50;
  }
}

/* Utility Classes */
@layer utilities {
  /* Gradient Text */
  .text-gradient {
    @apply bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent;
  }

  /* Gradient Background */
  .bg-gradient-primary {
    @apply bg-gradient-to-r from-primary to-primary-dark;
  }

  /* Glass Effect */
  .glass {
    @apply bg-background/80 backdrop-blur-lg border border-border/50;
  }

  /* Card Hover Effect */
  .card-hover {
    @apply transition-all duration-200 hover:shadow-lg hover:-translate-y-0.5;
  }

  /* No Scrollbar */
  .no-scrollbar::-webkit-scrollbar {
    display: none;
  }
  .no-scrollbar {
    -ms-overflow-style: none;
    scrollbar-width: none;
  }
}

/* Component Overrides */
@layer components {
  /* Button gradient variant */
  .btn-gradient {
    @apply bg-gradient-to-r from-primary to-primary-dark text-primary-foreground;
    @apply hover:opacity-90 transition-opacity;
  }

  /* Input focus enhancement */
  .input-enhanced {
    @apply transition-shadow duration-200;
    @apply focus:shadow-[0_0_0_3px] focus:shadow-primary/20;
  }

  /* Card with glow */
  .card-glow {
    @apply relative overflow-hidden;
  }
  .card-glow::before {
    content: '';
    @apply absolute -inset-px bg-gradient-to-r from-primary/20 to-secondary/20 rounded-[inherit] opacity-0 transition-opacity;
  }
  .card-glow:hover::before {
    @apply opacity-100;
  }
}
```

## Theme Toggle Component

```tsx
// components/theme-toggle.tsx
import { Moon, Sun } from 'lucide-react';
import { useTheme } from '@/hooks/use-theme';
import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';

export function ThemeToggle() {
  const { theme, setTheme } = useTheme();

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" size="icon">
          <Sun className="h-5 w-5 rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0" />
          <Moon className="absolute h-5 w-5 rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100" />
          <span className="sr-only">Toggle theme</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        <DropdownMenuItem onClick={() => setTheme('light')}>
          Light
        </DropdownMenuItem>
        <DropdownMenuItem onClick={() => setTheme('dark')}>
          Dark
        </DropdownMenuItem>
        <DropdownMenuItem onClick={() => setTheme('system')}>
          System
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
```

## Theme Provider Hook

```tsx
// hooks/use-theme.tsx
import { createContext, useContext, useEffect, useState } from 'react';

type Theme = 'light' | 'dark' | 'system';

interface ThemeContextValue {
  theme: Theme;
  setTheme: (theme: Theme) => void;
  resolvedTheme: 'light' | 'dark';
}

const ThemeContext = createContext<ThemeContextValue | undefined>(undefined);

export function ThemeProvider({ children }: { children: React.ReactNode }) {
  const [theme, setTheme] = useState<Theme>(() => {
    if (typeof window !== 'undefined') {
      return (localStorage.getItem('theme') as Theme) || 'system';
    }
    return 'system';
  });

  const [resolvedTheme, setResolvedTheme] = useState<'light' | 'dark'>('light');

  useEffect(() => {
    const root = window.document.documentElement;
    const systemTheme = window.matchMedia('(prefers-color-scheme: dark)').matches
      ? 'dark'
      : 'light';

    const resolved = theme === 'system' ? systemTheme : theme;
    setResolvedTheme(resolved);

    root.classList.remove('light', 'dark');
    root.classList.add(resolved);

    localStorage.setItem('theme', theme);
  }, [theme]);

  // Listen for system theme changes
  useEffect(() => {
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
    const handleChange = () => {
      if (theme === 'system') {
        const resolved = mediaQuery.matches ? 'dark' : 'light';
        setResolvedTheme(resolved);
        document.documentElement.classList.remove('light', 'dark');
        document.documentElement.classList.add(resolved);
      }
    };

    mediaQuery.addEventListener('change', handleChange);
    return () => mediaQuery.removeEventListener('change', handleChange);
  }, [theme]);

  return (
    <ThemeContext.Provider value={{ theme, setTheme, resolvedTheme }}>
      {children}
    </ThemeContext.Provider>
  );
}

export function useTheme() {
  const context = useContext(ThemeContext);
  if (!context) {
    throw new Error('useTheme must be used within a ThemeProvider');
  }
  return context;
}
```

## Tailwind Config (v4)

With Tailwind CSS v4, configuration is primarily done in CSS. However, if you need JavaScript config:

```ts
// tailwind.config.ts
import type { Config } from 'tailwindcss';

export default {
  darkMode: 'class',
  content: ['./src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        border: 'var(--border)',
        input: 'var(--input)',
        ring: 'var(--ring)',
        background: 'var(--background)',
        foreground: 'var(--foreground)',
        primary: {
          DEFAULT: 'var(--primary)',
          foreground: 'var(--primary-foreground)',
        },
        secondary: {
          DEFAULT: 'var(--secondary)',
          foreground: 'var(--secondary-foreground)',
        },
        destructive: {
          DEFAULT: 'var(--destructive)',
          foreground: 'var(--destructive-foreground)',
        },
        success: {
          DEFAULT: 'var(--success)',
          foreground: 'var(--success-foreground)',
        },
        warning: {
          DEFAULT: 'var(--warning)',
          foreground: 'var(--warning-foreground)',
        },
        muted: {
          DEFAULT: 'var(--muted)',
          foreground: 'var(--muted-foreground)',
        },
        accent: {
          DEFAULT: 'var(--accent)',
          foreground: 'var(--accent-foreground)',
        },
        popover: {
          DEFAULT: 'var(--popover)',
          foreground: 'var(--popover-foreground)',
        },
        card: {
          DEFAULT: 'var(--card)',
          foreground: 'var(--card-foreground)',
        },
      },
      borderRadius: {
        lg: 'var(--radius)',
        md: 'calc(var(--radius) - 2px)',
        sm: 'calc(var(--radius) - 4px)',
      },
      fontFamily: {
        sans: ['var(--font-sans)'],
        mono: ['var(--font-mono)'],
      },
    },
  },
} satisfies Config;
```

## Usage Examples

```tsx
// Using theme colors
<div className="bg-primary text-primary-foreground">Primary</div>
<div className="bg-secondary text-secondary-foreground">Secondary</div>
<div className="bg-success text-success-foreground">Success</div>
<div className="bg-warning text-warning-foreground">Warning</div>
<div className="bg-destructive text-destructive-foreground">Error</div>

// Using gradients
<button className="btn-gradient px-4 py-2 rounded-md">
  Gradient Button
</button>

// Using glass effect
<div className="glass p-4 rounded-lg">
  Glassmorphism card
</div>

// Card with hover effect
<div className="card-hover bg-card p-4 rounded-lg border">
  Hover me
</div>
```
