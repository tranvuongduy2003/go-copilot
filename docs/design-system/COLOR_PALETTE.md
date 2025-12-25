# 🎨 Color Palette Documentation

> **For AI Agents**: This document defines the official color system. ALWAYS use these exact colors and semantic tokens. Never generate arbitrary hex values or create new colors.

---

## Overview

Our color system uses **OKLCH** (Oklch Lightness Chroma Hue) for perceptual uniformity. This means colors look equally vibrant across the palette and transitions are smooth.

```
oklch(L C H)
  │  │  └── Hue: 0-360 (color wheel position)
  │  └───── Chroma: 0-0.4 (saturation/vibrancy)
  └──────── Lightness: 0-1 (brightness)
```

---

## 🎯 Primary Brand Colors

**Hue: 285 (Purple/Violet)**

Use for: Primary actions, buttons, links, active states, brand elements

| Token | OKLCH Value | Usage |
|-------|-------------|-------|
| `brand-50` | `oklch(0.970 0.030 280)` | Subtle backgrounds |
| `brand-100` | `oklch(0.943 0.054 281)` | Light backgrounds, hover states |
| `brand-200` | `oklch(0.894 0.098 282)` | Borders, dividers |
| `brand-300` | `oklch(0.831 0.150 283)` | Icons, secondary elements |
| `brand-400` | `oklch(0.750 0.190 284)` | **Dark mode primary** |
| `brand-500` | `oklch(0.667 0.210 285)` | ⭐ **Light mode primary** |
| `brand-600` | `oklch(0.585 0.200 286)` | Hover state |
| `brand-700` | `oklch(0.510 0.180 287)` | Active/pressed state |
| `brand-800` | `oklch(0.432 0.150 288)` | Dark emphasis |
| `brand-900` | `oklch(0.365 0.120 289)` | Very dark |
| `brand-950` | `oklch(0.257 0.090 290)` | Dark mode backgrounds |

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

## 💎 Accent Colors

**Hue: 195 (Cyan/Teal)**

Use for: Secondary actions, highlights, links, decorative elements

| Token | OKLCH Value | Usage |
|-------|-------------|-------|
| `accent-50` | `oklch(0.972 0.025 195)` | Subtle backgrounds |
| `accent-100` | `oklch(0.948 0.050 195)` | Light backgrounds |
| `accent-200` | `oklch(0.901 0.095 195)` | Borders |
| `accent-300` | `oklch(0.837 0.140 195)` | Icons |
| `accent-400` | `oklch(0.752 0.175 195)` | **Dark mode accent** |
| `accent-500` | `oklch(0.665 0.190 195)` | ⭐ **Light mode accent** |
| `accent-600` | `oklch(0.580 0.175 195)` | Hover state |
| `accent-700` | `oklch(0.500 0.155 195)` | Active state |
| `accent-800` | `oklch(0.425 0.130 195)` | Dark emphasis |
| `accent-900` | `oklch(0.360 0.105 195)` | Very dark |
| `accent-950` | `oklch(0.250 0.075 195)` | Dark backgrounds |

---

## 🔘 Neutral Gray Scale

**Hue: 264 (Slight blue undertone for modern feel)**

Use for: Text, backgrounds, borders, shadows

| Token | OKLCH Value | Light Mode Usage | Dark Mode Usage |
|-------|-------------|------------------|-----------------|
| `gray-50` | `oklch(0.985 0.002 264)` | Background | — |
| `gray-100` | `oklch(0.967 0.003 264)` | Subtle bg | Foreground |
| `gray-200` | `oklch(0.928 0.006 264)` | Borders | — |
| `gray-300` | `oklch(0.872 0.010 264)` | Emphasis bg | — |
| `gray-400` | `oklch(0.707 0.015 264)` | Muted fg (dark) | Muted foreground |
| `gray-500` | `oklch(0.551 0.018 264)` | Placeholder | Subtle foreground |
| `gray-600` | `oklch(0.446 0.018 264)` | Muted foreground | — |
| `gray-700` | `oklch(0.372 0.016 264)` | — | Emphasis border |
| `gray-800` | `oklch(0.279 0.012 264)` | — | Borders |
| `gray-900` | `oklch(0.208 0.010 264)` | Foreground | Subtle bg |
| `gray-950` | `oklch(0.129 0.008 264)` | — | Background |

---

## ✅ Semantic Status Colors

### Success (Green)
**Hue: 155**

| Token | Value | Usage |
|-------|-------|-------|
| `success-500` | `oklch(0.680 0.195 155)` | Success states, checkmarks |
| `success-600` | `oklch(0.590 0.175 155)` | Hover state |

### Warning (Amber)
**Hue: 85**

| Token | Value | Usage |
|-------|-------|-------|
| `warning-500` | `oklch(0.770 0.195 85)` | Warnings, alerts |
| `warning-600` | `oklch(0.680 0.175 85)` | Hover state |

### Error (Red)
**Hue: 25**

| Token | Value | Usage |
|-------|-------|-------|
| `error-500` | `oklch(0.680 0.205 25)` | Errors, destructive |
| `error-600` | `oklch(0.590 0.190 25)` | Hover state |

---

## 🌓 Theme Mapping

### Semantic Variables

Always use semantic variables in components, not primitive colors:

```css
/* Light Mode */
--background: var(--color-gray-50);
--foreground: var(--color-gray-900);
--primary: var(--color-brand-500);
--border: var(--color-gray-200);

/* Dark Mode */
--background: var(--color-gray-950);
--foreground: var(--color-gray-50);
--primary: var(--color-brand-400);
--border: var(--color-gray-800);
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
<Card className="bg-gray-50 dark:bg-gray-950">...</Card>
```

---

## 🎨 Gradients

### Brand Gradient
```css
background: linear-gradient(135deg, var(--color-brand-400), var(--color-accent-400));
```

### Text Gradient
```tsx
<span className="text-gradient-brand">Gradient Text</span>
```

### Glow Effects
```tsx
<Button className="glow-brand">Glowing Button</Button>
```

---

## 🔒 Rules for AI Agents

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