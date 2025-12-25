# Typography System

This document details our typography system, built around Inter for UI text and JetBrains Mono for code.

## Token Format (W3C DTCG)

Typography tokens follow the W3C Design Tokens Community Group specification:

```json
{
  "fontFamily": {
    "$type": "fontFamily",
    "sans": {
      "$value": ["Inter", "ui-sans-serif", "system-ui", "sans-serif"],
      "$description": "Primary font for body text and UI"
    }
  },
  "typography": {
    "$type": "typography",
    "heading": {
      "h1": {
        "$value": {
          "fontFamily": "{fontFamily.sans}",
          "fontSize": "{fontSize.4xl}",
          "fontWeight": "{fontWeight.bold}",
          "lineHeight": "{lineHeight.tight}"
        }
      }
    }
  }
}
```

## Font Families

### Inter (Sans-serif)

Primary font for all UI text, body content, and headings.

```css
--font-sans: 'Inter', ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
```

**Why Inter?**
- Designed specifically for computer screens
- Excellent legibility at small sizes
- Variable font with wide weight range
- Open source and free for commercial use
- Extensive language support

**Installation:**
```html
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap" rel="stylesheet">
```

### JetBrains Mono (Monospace)

Used for code blocks, inline code, and technical content.

```css
--font-mono: 'JetBrains Mono', ui-monospace, SFMono-Regular, 'SF Mono', Menlo, Consolas, 'Liberation Mono', monospace;
```

**Why JetBrains Mono?**
- Increased letter height for better readability
- Distinguished characters (0, O, l, 1)
- Code-specific ligatures
- Open source

**Installation:**
```html
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500&display=swap" rel="stylesheet">
```

## Type Scale

Our type scale is based on a 1.25 ratio (major third), providing clear visual hierarchy.

| Token | Size | Rem | Usage |
|-------|------|-----|-------|
| `xs` | 12px | 0.75rem | Labels, captions, badges |
| `sm` | 14px | 0.875rem | Secondary text, descriptions |
| `base` | 16px | 1rem | Body text (default) |
| `lg` | 18px | 1.125rem | Emphasized body, lead text |
| `xl` | 20px | 1.25rem | H4, card titles |
| `2xl` | 24px | 1.5rem | H3, section headers |
| `3xl` | 30px | 1.875rem | H2, page sections |
| `4xl` | 36px | 2.25rem | H1, page titles |
| `5xl` | 48px | 3rem | Display text |
| `6xl` | 60px | 3.75rem | Hero headlines |
| `7xl` | 72px | 4.5rem | Large hero text |

## Font Weights

| Token | Weight | Usage |
|-------|--------|-------|
| `normal` | 400 | Body text, descriptions |
| `medium` | 500 | Labels, emphasized text |
| `semibold` | 600 | Subheadings, buttons |
| `bold` | 700 | Headings, CTAs |

## Line Heights

| Token | Value | Usage |
|-------|-------|-------|
| `none` | 1 | Single line elements, icons |
| `tight` | 1.25 | Headings |
| `snug` | 1.375 | Subheadings |
| `normal` | 1.5 | Body text (default) |
| `relaxed` | 1.625 | Long-form content |
| `loose` | 2 | Spacious layouts |

## Letter Spacing

| Token | Value | Usage |
|-------|-------|-------|
| `tighter` | -0.05em | Display headings |
| `tight` | -0.025em | Large headings |
| `normal` | 0 | Body text (default) |
| `wide` | 0.025em | Small caps |
| `wider` | 0.05em | Labels |
| `widest` | 0.1em | Uppercase labels |

## Text Styles

Pre-composed text styles for consistent typography across the application.

### Headings

```tsx
// H1 - Page titles
<h1 className="text-4xl font-bold leading-tight tracking-tight">
  Page Title
</h1>

// H2 - Section headers
<h2 className="text-3xl font-semibold leading-tight tracking-tight">
  Section Header
</h2>

// H3 - Subsection headers
<h3 className="text-2xl font-semibold leading-snug">
  Subsection Header
</h3>

// H4 - Card titles
<h4 className="text-xl font-medium leading-snug">
  Card Title
</h4>
```

### Body Text

```tsx
// Large body - Lead paragraphs
<p className="text-lg leading-relaxed text-foreground-secondary">
  Lead paragraph text for introductions.
</p>

// Default body
<p className="text-base leading-normal text-foreground-secondary">
  Regular paragraph text for content.
</p>

// Small body - Supporting text
<p className="text-sm leading-normal text-muted-foreground">
  Supporting text and descriptions.
</p>
```

### Labels

```tsx
// Large label
<label className="text-base font-medium leading-none">
  Large Label
</label>

// Default label
<label className="text-sm font-medium leading-none">
  Default Label
</label>

// Small label
<label className="text-xs font-medium leading-none">
  Small Label
</label>
```

### Code

```tsx
// Inline code
<code className="font-mono text-sm bg-muted px-1.5 py-0.5 rounded">
  inlineCode()
</code>

// Code block
<pre className="font-mono text-sm bg-muted p-4 rounded-lg overflow-x-auto">
  <code>
    {`function example() {
  return 'Hello, world!';
}`}
  </code>
</pre>
```

## Component Typography

### Buttons

| Size | Font Size | Font Weight | Line Height |
|------|-----------|-------------|-------------|
| Small | 14px (sm) | 500 (medium) | 1 (none) |
| Default | 14px (sm) | 500 (medium) | 1 (none) |
| Large | 16px (base) | 500 (medium) | 1 (none) |

```tsx
<Button size="sm">Small Button</Button>
<Button>Default Button</Button>
<Button size="lg">Large Button</Button>
```

### Form Inputs

| Element | Font Size | Font Weight |
|---------|-----------|-------------|
| Input text | 14px (sm) | 400 (normal) |
| Placeholder | 14px (sm) | 400 (normal) |
| Label | 14px (sm) | 500 (medium) |
| Helper text | 12px (xs) | 400 (normal) |
| Error text | 12px (xs) | 400 (normal) |

### Cards

| Element | Style |
|---------|-------|
| Title | text-lg font-semibold |
| Description | text-sm text-muted-foreground |
| Content | text-base |

### Tables

| Element | Style |
|---------|-------|
| Header | text-sm font-medium |
| Cell | text-sm |
| Caption | text-sm text-muted-foreground |

## Responsive Typography

For optimal readability across devices:

```css
/* Mobile-first responsive headings */
h1 {
  font-size: 1.875rem; /* 30px */
}

@media (min-width: 768px) {
  h1 {
    font-size: 2.25rem; /* 36px */
  }
}

@media (min-width: 1024px) {
  h1 {
    font-size: 2.5rem; /* 40px */
  }
}
```

With Tailwind:

```tsx
<h1 className="text-3xl md:text-4xl lg:text-5xl font-bold">
  Responsive Heading
</h1>
```

## Text Colors

| Token | Usage |
|-------|-------|
| `--foreground` | Primary text, headings |
| `--foreground-secondary` | Body text |
| `--foreground-muted` | Placeholder, disabled |
| `--muted-foreground` | Captions, helper text |
| `--primary` | Links, emphasized |
| `--destructive` | Error messages |
| `--success` | Success messages |

## Best Practices

### Hierarchy

1. Use a maximum of 3-4 text sizes per page
2. Maintain clear contrast between heading levels
3. Use font weight to create emphasis, not size alone

### Readability

1. Body text: 16px minimum for comfortable reading
2. Line length: 60-80 characters optimal
3. Line height: 1.5 for body, 1.25 for headings
4. Contrast: 4.5:1 minimum for normal text

### Alignment

1. Left-align body text (LTR languages)
2. Center-align only for short text (headings, CTAs)
3. Never justify body text on the web

### Consistency

```tsx
// DO: Use text style utilities
<h2 className="text-2xl font-semibold leading-snug">Heading</h2>

// DON'T: Use arbitrary values
<h2 style={{ fontSize: '23px', fontWeight: 590 }}>Heading</h2>
```

## CSS Custom Properties

```css
/* Font families */
--font-sans: ...;
--font-mono: ...;

/* Font sizes */
--font-size-xs: 0.75rem;
--font-size-sm: 0.875rem;
--font-size-base: 1rem;
/* ... */

/* Font weights */
--font-weight-normal: 400;
--font-weight-medium: 500;
--font-weight-semibold: 600;
--font-weight-bold: 700;

/* Line heights */
--line-height-tight: 1.25;
--line-height-normal: 1.5;
/* ... */

/* Letter spacing */
--letter-spacing-tight: -0.025em;
--letter-spacing-normal: 0;
/* ... */
```

## Tailwind Classes Quick Reference

```tsx
// Font family
font-sans  // Inter
font-mono  // JetBrains Mono

// Font size
text-xs    // 12px
text-sm    // 14px
text-base  // 16px
text-lg    // 18px
text-xl    // 20px
text-2xl   // 24px
text-3xl   // 30px
text-4xl   // 36px

// Font weight
font-normal    // 400
font-medium    // 500
font-semibold  // 600
font-bold      // 700

// Line height
leading-none     // 1
leading-tight    // 1.25
leading-snug     // 1.375
leading-normal   // 1.5
leading-relaxed  // 1.625
leading-loose    // 2

// Letter spacing
tracking-tighter // -0.05em
tracking-tight   // -0.025em
tracking-normal  // 0
tracking-wide    // 0.025em
tracking-wider   // 0.05em
tracking-widest  // 0.1em

// Text color
text-foreground
text-foreground-secondary
text-muted-foreground
text-primary
```
