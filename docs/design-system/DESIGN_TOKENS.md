# Design Tokens

This document provides a comprehensive reference for all design tokens used in our design system. Design tokens are the visual design atoms of the design system — specifically, they are named entities that store visual design attributes.

## Token Categories

| Category | Description | Location |
|----------|-------------|----------|
| Colors | OKLCH color palette with semantic tokens | `design-system/tokens/colors.json` |
| Typography | Font families, sizes, weights, line heights | `design-system/tokens/typography.json` |
| Spacing | Spacing scale, border radius, shadows | `design-system/tokens/spacing.json` |

## Why OKLCH?

We use the OKLCH color space for several key advantages:

1. **Perceptual Uniformity**: Equal steps in lightness look equally different to human eyes
2. **Consistent Chroma**: Saturation remains visually consistent across hues
3. **Wide Gamut**: Supports P3 and future display technologies
4. **Predictable Mixing**: Color interpolations produce expected results

### OKLCH Format

```
oklch(L C H / A)
```

- **L (Lightness)**: 0 (black) to 1 (white)
- **C (Chroma)**: 0 (gray) to ~0.4 (most saturated)
- **H (Hue)**: 0-360 degrees
- **A (Alpha)**: Optional opacity value

## Color Tokens

### Primary Scale (Violet - Hue 290)

| Token | OKLCH Value | Usage |
|-------|-------------|-------|
| `primary-50` | `oklch(0.97 0.02 290)` | Subtle backgrounds |
| `primary-100` | `oklch(0.94 0.04 290)` | Hover states |
| `primary-200` | `oklch(0.88 0.08 290)` | Selection backgrounds |
| `primary-300` | `oklch(0.79 0.12 290)` | Borders |
| `primary-400` | `oklch(0.70 0.16 290)` | Icons |
| `primary-500` | `oklch(0.60 0.20 290)` | Primary actions (dark mode) |
| `primary-600` | `oklch(0.52 0.20 290)` | Primary actions (light mode) |
| `primary-700` | `oklch(0.44 0.18 290)` | Hover states |
| `primary-800` | `oklch(0.36 0.14 290)` | Active states |
| `primary-900` | `oklch(0.28 0.10 290)` | Text on light backgrounds |
| `primary-950` | `oklch(0.20 0.08 290)` | Darkest primary |

### Secondary Scale (Cyan - Hue 220)

| Token | OKLCH Value | Usage |
|-------|-------------|-------|
| `secondary-50` | `oklch(0.97 0.02 220)` | Subtle backgrounds |
| `secondary-100` | `oklch(0.94 0.04 220)` | Hover states |
| `secondary-200` | `oklch(0.88 0.08 220)` | Selection backgrounds |
| `secondary-300` | `oklch(0.79 0.12 220)` | Borders |
| `secondary-400` | `oklch(0.70 0.16 220)` | Icons |
| `secondary-500` | `oklch(0.60 0.18 220)` | Secondary actions |
| `secondary-600` | `oklch(0.52 0.18 220)` | Hover states |
| `secondary-700` | `oklch(0.44 0.16 220)` | Active states |
| `secondary-800` | `oklch(0.36 0.12 220)` | Text |
| `secondary-900` | `oklch(0.28 0.08 220)` | Dark text |
| `secondary-950` | `oklch(0.20 0.06 220)` | Darkest secondary |

### Semantic Color Tokens

```css
/* Light Mode */
--background: var(--color-neutral-50);
--foreground: var(--color-neutral-950);
--primary: var(--color-primary-600);
--primary-foreground: oklch(1 0 0);
--secondary: var(--color-secondary-100);
--secondary-foreground: var(--color-secondary-900);
--destructive: var(--color-error-500);
--success: var(--color-success-500);
--warning: var(--color-warning-500);

/* Dark Mode */
--background: var(--color-neutral-950);
--foreground: var(--color-neutral-50);
--primary: var(--color-primary-500);
```

## Spacing Tokens

Based on a 4px base unit for consistent rhythm.

| Token | Value | Usage |
|-------|-------|-------|
| `0` | 0px | Reset |
| `px` | 1px | Hairline borders |
| `0.5` | 2px | Minimal spacing |
| `1` | 4px | Tight padding |
| `1.5` | 6px | - |
| `2` | 8px | Inline elements |
| `2.5` | 10px | - |
| `3` | 12px | Medium-small |
| `3.5` | 14px | - |
| `4` | 16px | Default component padding |
| `5` | 20px | Medium spacing |
| `6` | 24px | Card spacing |
| `7` | 28px | - |
| `8` | 32px | Section padding |
| `9` | 36px | - |
| `10` | 40px | Large sections |
| `11` | 44px | - |
| `12` | 48px | Page sections |
| `14` | 56px | - |
| `16` | 64px | 2XL spacing |
| `20` | 80px | 3XL spacing |
| `24` | 96px | Hero spacing |
| `28` | 112px | - |
| `32` | 128px | Maximum spacing |

## Border Radius Tokens

| Token | Value | Usage |
|-------|-------|-------|
| `none` | 0px | Sharp corners |
| `sm` | 4px | Badges, chips |
| `md` | 8px | Buttons, inputs (default) |
| `lg` | 12px | Cards, dialogs |
| `xl` | 16px | Large cards, modals |
| `2xl` | 24px | Feature sections |
| `3xl` | 32px | Large decorative |
| `full` | 9999px | Pills, avatars, circles |

## Shadow Tokens

All shadows use a violet tint (hue 290) for brand consistency.

| Token | Value | Usage |
|-------|-------|-------|
| `sm` | `0 1px 2px oklch(0.3 0.05 290 / 0.05)` | Subtle elevation |
| `md` | `0 4px 6px -1px oklch(0.3 0.05 290 / 0.1), ...` | Card elevation |
| `lg` | `0 10px 15px -3px oklch(0.3 0.05 290 / 0.1), ...` | Dropdown/popover |
| `xl` | `0 20px 25px -5px oklch(0.3 0.05 290 / 0.1), ...` | Modal |
| `2xl` | `0 25px 50px -12px oklch(0.3 0.05 290 / 0.25)` | High elevation |
| `inner` | `inset 0 2px 4px 0 oklch(0 0 0 / 0.05)` | Input inset |

## Z-Index Scale

| Token | Value | Usage |
|-------|-------|-------|
| `auto` | auto | Default stacking |
| `0` | 0 | Base |
| `10` | 10 | Base layer |
| `20` | 20 | Sticky elements |
| `30` | 30 | Fixed headers |
| `40` | 40 | Dropdowns |
| `50` | 50 | Modals |
| `60` | 60 | Popovers |
| `70` | 70 | Tooltips |
| `80` | 80 | Toasts |
| `90` | 90 | Maximum |

## Using Tokens in Code

### CSS Custom Properties

```css
.card {
  background: var(--card);
  color: var(--card-foreground);
  padding: var(--spacing-6);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
}
```

### Tailwind CSS

```tsx
<div className="bg-card text-card-foreground p-6 rounded-lg shadow-md">
  Card content
</div>
```

### Direct Token Reference

```tsx
import { tokens } from '@/design-system/tokens';

const styles = {
  padding: tokens.spacing[6],
  borderRadius: tokens.borderRadius.lg,
};
```

## Token Naming Convention

Tokens follow this naming pattern:

```
--{category}-{scale/variant}
```

Examples:
- `--color-primary-500` - Primary color at 500 scale
- `--spacing-4` - 4 unit spacing (16px)
- `--radius-md` - Medium border radius
- `--shadow-lg` - Large shadow

## Migration Guide

When updating tokens:

1. Update the JSON source file
2. Run the token build script: `npm run build:tokens`
3. Update any direct CSS property references
4. Test across light and dark modes
5. Verify WCAG contrast ratios

## Accessibility Considerations

- All color combinations meet WCAG 2.1 AA contrast requirements
- Interactive elements maintain 3:1 contrast against backgrounds
- Focus states use high-visibility ring color (primary-500)
- Text colors ensure minimum 4.5:1 contrast ratio
