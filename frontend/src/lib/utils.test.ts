import { describe, expect, it } from 'vitest';
import { cn } from './utils';

describe('cn (className utility)', () => {
  it('merges class names', () => {
    const result = cn('foo', 'bar');
    expect(result).toBe('foo bar');
  });

  it('handles conditional classes', () => {
    const condition = false;
    const result = cn('foo', condition && 'bar', 'baz');
    expect(result).toBe('foo baz');
  });

  it('handles undefined and null values', () => {
    const result = cn('foo', undefined, null, 'bar');
    expect(result).toBe('foo bar');
  });

  it('handles arrays of classes', () => {
    const result = cn(['foo', 'bar'], 'baz');
    expect(result).toBe('foo bar baz');
  });

  it('handles objects with boolean values', () => {
    const result = cn('foo', { bar: true, baz: false });
    expect(result).toBe('foo bar');
  });

  it('merges Tailwind classes correctly', () => {
    const result = cn('px-4 py-2', 'px-6');
    expect(result).toBe('py-2 px-6');
  });

  it('merges conflicting Tailwind utilities', () => {
    const result = cn('text-red-500', 'text-blue-500');
    expect(result).toBe('text-blue-500');
  });

  it('handles mixed bg classes', () => {
    const result = cn('bg-primary', 'bg-secondary');
    expect(result).toBe('bg-secondary');
  });

  it('preserves non-conflicting utilities', () => {
    const result = cn('p-4', 'mx-auto', 'text-center');
    expect(result).toBe('p-4 mx-auto text-center');
  });

  it('handles responsive variants', () => {
    const result = cn('text-sm', 'md:text-base', 'lg:text-lg');
    expect(result).toBe('text-sm md:text-base lg:text-lg');
  });

  it('handles hover and focus variants', () => {
    const result = cn('bg-blue-500', 'hover:bg-blue-600', 'focus:bg-blue-700');
    expect(result).toBe('bg-blue-500 hover:bg-blue-600 focus:bg-blue-700');
  });

  it('handles empty inputs', () => {
    const result = cn();
    expect(result).toBe('');
  });

  it('handles only falsy values', () => {
    const result = cn(false, null, undefined, '');
    expect(result).toBe('');
  });

  it('handles deeply nested arrays', () => {
    const result = cn(['foo', ['bar', 'baz']]);
    expect(result).toBe('foo bar baz');
  });

  it('combines shadcn component classes correctly', () => {
    const baseClasses = 'inline-flex items-center justify-center rounded-md text-sm font-medium';
    const variantClasses = 'bg-primary text-primary-foreground';
    const sizeClasses = 'h-10 px-4 py-2';
    const customClasses = 'w-full';

    const result = cn(baseClasses, variantClasses, sizeClasses, customClasses);
    expect(result).toContain('inline-flex');
    expect(result).toContain('bg-primary');
    expect(result).toContain('h-10');
    expect(result).toContain('w-full');
  });
});
