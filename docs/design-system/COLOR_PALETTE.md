# Color Palette Documentation

> **For AI Agents**: This document defines the official color system using W3C DTCG format. ALWAYS use these exact colors and semantic tokens. Never generate arbitrary hex values or create new colors.

---

## Overview

Our color system uses **OKLCH** (Oklch Lightness Chroma Hue) for perceptual uniformity. This means colors look equally vibrant across the palette and transitions are smooth.

```
oklch(L C H / A)
  │  │  │   └── Alpha: 0-1 (optional opacity)
  │  │  └───── Hue: 0-360 (color wheel position)
  │  └──────── Chroma: 0-0.4 (saturation/vibrancy)
  └─────────── Lightness: 0-1 (brightness)
```

### Token Format (W3C DTCG)

```json
{
  "color": {
    "$type": "color",
    "primary": {
      "500": {
        "$value": "oklch(0.60 0.20 290)",
        "$description": "Primary actions (dark mode)"
      }
    }
  }
}
```

---

## Primary Brand Colors

**Hue: 290 (Violet/Purple)**

Use for: Primary actions, buttons, links, active states, brand elements

| Token | OKLCH Value | Usage |
|-------|-------------|-------|
| `primary-50` | `oklch(0.97 0.02 290)` | Subtle backgrounds |
| `primary-100` | `oklch(0.94 0.04 290)` | Light backgrounds, hover states |
| `primary-200` | `oklch(0.88 0.08 290)` | Borders, dividers |
| `primary-300` | `oklch(0.79 0.12 290)` | Icons, secondary elements |
| `primary-400` | `oklch(0.70 0.16 290)` | Dark mode hover |
| `primary-500` | `oklch(0.60 0.20 290)` | **Dark mode primary** |
| `primary-600` | `oklch(0.52 0.20 290)` | **Light mode primary** |
| `primary-700` | `oklch(0.44 0.18 290)` | Hover state |
| `primary-800` | `oklch(0.36 0.14 290)` | Active/pressed state |
| `primary-900` | `oklch(0.28 0.10 290)` | Text on light backgrounds |
| `primary-950` | `oklch(0.20 0.08 290)` | Darkest primary |

### Usage Examples

```tsx
// ✅ CORRECT - Use semantic tokens
<Button className="bg-primary hover:bg-primary-hover">
  Primary Action
</Button>

// ✅ CORRECT - Use CSS variables
<div style={{ backgroundColor: 'var(--primary)' }}>
  Brand element
</div>

// ❌ WRONG - Never use arbitrary colors
<Button className="bg-[#7c3aed]">Wrong</Button>
```

---

## Secondary Colors

**Hue: 220 (Cyan/Blue)**

Use for: Secondary actions, highlights, links, decorative elements

| Token | OKLCH Value | Usage |
|-------|-------------|-------|
| `secondary-50` | `oklch(0.97 0.02 220)` | Subtle backgrounds |
| `secondary-100` | `oklch(0.94 0.04 220)` | Light backgrounds |
| `secondary-200` | `oklch(0.88 0.08 220)` | Borders |
| `secondary-300` | `oklch(0.79 0.12 220)` | Icons |
| `secondary-400` | `oklch(0.70 0.16 220)` | **Dark mode secondary** |
| `secondary-500` | `oklch(0.60 0.18 220)` | Secondary actions |
| `secondary-600` | `oklch(0.52 0.18 220)` | Hover state |
| `secondary-700` | `oklch(0.44 0.16 220)` | Active state |
| `secondary-800` | `oklch(0.36 0.12 220)` | Dark emphasis |
| `secondary-900` | `oklch(0.28 0.08 220)` | Text |
| `secondary-950` | `oklch(0.20 0.06 220)` | Darkest secondary |

---

## Neutral Gray Scale

**Hue: 290 (Slight violet undertone for brand consistency)**

Use for: Text, backgrounds, borders, shadows

| Token | OKLCH Value | Light Mode Usage | Dark Mode Usage |
|-------|-------------|------------------|-----------------|
| `neutral-50` | `oklch(0.98 0.005 290)` | Background | — |
| `neutral-100` | `oklch(0.96 0.008 290)` | Subtle bg | Foreground |
| `neutral-200` | `oklch(0.92 0.010 290)` | Borders | — |
| `neutral-300` | `oklch(0.87 0.012 290)` | Emphasis bg | — |
| `neutral-400` | `oklch(0.70 0.015 290)` | Muted fg (dark) | Muted foreground |
| `neutral-500` | `oklch(0.55 0.015 290)` | Placeholder | Subtle foreground |
| `neutral-600` | `oklch(0.45 0.015 290)` | Muted foreground | — |
| `neutral-700` | `oklch(0.35 0.012 290)` | — | Emphasis border |
| `neutral-800` | `oklch(0.25 0.010 290)` | — | Borders |
| `neutral-900` | `oklch(0.18 0.008 290)` | Foreground | Subtle bg |
| `neutral-950` | `oklch(0.12 0.005 290)` | — | Background |

---

## Semantic Status Colors

### Success (Green - Hue: 145)

| Token | Value | Usage |
|-------|-------|-------|
| `success-50` | `oklch(0.97 0.02 145)` | Success backgrounds |
| `success-500` | `oklch(0.60 0.18 145)` | Success states, checkmarks |
| `success-600` | `oklch(0.52 0.16 145)` | Hover state |

### Warning (Amber - Hue: 85)

| Token | Value | Usage |
|-------|-------|-------|
| `warning-50` | `oklch(0.97 0.02 85)` | Warning backgrounds |
| `warning-500` | `oklch(0.75 0.18 85)` | Warnings, alerts |
| `warning-600` | `oklch(0.65 0.16 70)` | Hover state |

### Error/Destructive (Red - Hue: 25)

| Token | Value | Usage |
|-------|-------|-------|
| `error-50` | `oklch(0.97 0.02 25)` | Error backgrounds |
| `error-500` | `oklch(0.60 0.20 25)` | Errors, destructive |
| `error-600` | `oklch(0.52 0.20 25)` | Hover state |

---

## Theme Mapping

### Semantic Variables

Always use semantic variables in components, not primitive colors:

```css
/* Light Mode */
:root {
  --background: var(--color-neutral-50);
  --foreground: var(--color-neutral-950);
  --primary: var(--color-primary-600);
  --primary-foreground: oklch(1 0 0);
  --border: var(--color-neutral-200);
}

/* Dark Mode */
.dark {
  --background: var(--color-neutral-950);
  --foreground: var(--color-neutral-50);
  --primary: var(--color-primary-500);
  --border: var(--color-neutral-800);
}
```

### Component Example

```tsx
// ✅ CORRECT - Uses semantic variables
<Card className="bg-card border-border text-card-foreground">
  <CardHeader>
    <CardTitle className="text-foreground">Title</CardTitle>
    <CardDescription className="text-muted-foreground">
      Description text
    </CardDescription>
  </CardHeader>
</Card>

// ❌ WRONG - Uses primitive colors directly
<Card className="bg-neutral-50 dark:bg-neutral-950">...</Card>
```

---

## Gradients

### Brand Gradient
```css
background: linear-gradient(135deg, var(--color-primary-400), var(--color-secondary-400));
```

### Text Gradient
```tsx
<span className="bg-linear-to-r from-primary-400 to-secondary-400 bg-clip-text text-transparent">
  Gradient Text
</span>
```

---

## Rules for AI Agents

### ✅ DO

1. **Always use CSS variables**: `var(--primary)`, `var(--background)`
2. **Use Tailwind semantic classes**: `bg-primary`, `text-foreground`
3. **Follow the 60-30-10 rule**:
   - 60% neutral (backgrounds, body text)
   - 30% secondary (cards, borders)
   - 10% accent (CTAs, highlights)
4. **Maintain contrast ratios**: 4.5:1 minimum for text

### ❌ DON'T

1. **Never use arbitrary hex colors**: `bg-[#7c3aed]`
2. **Never create new colors**: Stick to the defined palette
3. **Never skip dark mode**: Always consider both themes
4. **Never ignore accessibility**: Use proper contrast

---

## 📊 Color Combinations

### Recommended Pairings

| Context | Background | Text | Accent |
|---------|------------|------|--------|
| Page | `--background` | `--foreground` | `--primary` |
| Card | `--card` | `--card-foreground` | `--primary` |
| Muted | `--muted` | `--muted-foreground` | `--accent` |
| Alert | `--destructive` | `--destructive-foreground` | — |
| Success | `--success-subtle` | `--success-foreground` | — |

### Contrast Ratios

| Combination | Ratio | WCAG Level |
|-------------|-------|------------|
| `foreground` on `background` | 15.3:1 | AAA ✅ |
| `muted-foreground` on `background` | 4.8:1 | AA ✅ |
| `primary-foreground` on `primary` | 8.2:1 | AAA ✅ |
| `foreground` on `muted` | 11.2:1 | AAA ✅ |

---

## 🔄 Updating Colors

If you need to modify the color palette:

1. Update `design-system/tokens/colors.json`
2. Update `frontend/src/styles/globals.css`
3. Run design token build script
4. Update this documentation
5. Test in both light and dark modes
6. Verify accessibility compliance