import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import {
  capitalizeFirst,
  formatBytes,
  formatCompactNumber,
  formatCurrency,
  formatDate,
  formatDateTime,
  formatNumber,
  formatPercentage,
  formatRelativeTime,
  slugify,
  toTitleCase,
  truncateText,
} from './formatters';

describe('formatters', () => {
  describe('formatDate', () => {
    it('formats a Date object with default format', () => {
      const date = new Date('2024-03-15T10:30:00Z');
      expect(formatDate(date)).toBe('Mar 15, 2024');
    });

    it('formats an ISO string with default format', () => {
      expect(formatDate('2024-03-15T10:30:00Z')).toBe('Mar 15, 2024');
    });

    it('formats with custom format string', () => {
      const date = new Date('2024-03-15T10:30:00Z');
      expect(formatDate(date, 'yyyy-MM-dd')).toBe('2024-03-15');
    });

    it('returns empty string for null', () => {
      expect(formatDate(null)).toBe('');
    });

    it('returns empty string for undefined', () => {
      expect(formatDate(undefined)).toBe('');
    });

    it('returns empty string for invalid date string', () => {
      expect(formatDate('invalid-date')).toBe('');
    });

    it('handles different format patterns', () => {
      const date = new Date('2024-12-25T14:30:00Z');
      expect(formatDate(date, 'MMMM do, yyyy')).toBe('December 25th, 2024');
    });
  });

  describe('formatDateTime', () => {
    it('formats date with time using default format', () => {
      const date = new Date('2024-03-15T14:30:00');
      expect(formatDateTime(date)).toBe('Mar 15, 2024 2:30 PM');
    });

    it('formats with custom format', () => {
      const date = new Date('2024-03-15T14:30:00');
      expect(formatDateTime(date, 'yyyy-MM-dd HH:mm')).toBe('2024-03-15 14:30');
    });

    it('returns empty string for null', () => {
      expect(formatDateTime(null)).toBe('');
    });
  });

  describe('formatRelativeTime', () => {
    beforeEach(() => {
      vi.useFakeTimers();
      vi.setSystemTime(new Date('2024-03-15T12:00:00Z'));
    });

    afterEach(() => {
      vi.useRealTimers();
    });

    it('formats relative time with suffix', () => {
      const pastDate = new Date('2024-03-15T11:00:00Z');
      expect(formatRelativeTime(pastDate)).toBe('about 1 hour ago');
    });

    it('formats relative time without suffix', () => {
      const pastDate = new Date('2024-03-15T11:00:00Z');
      expect(formatRelativeTime(pastDate, { addSuffix: false })).toBe('about 1 hour');
    });

    it('formats ISO string', () => {
      const result = formatRelativeTime('2024-03-14T12:00:00Z');
      expect(result).toBe('1 day ago');
    });

    it('returns empty string for null', () => {
      expect(formatRelativeTime(null)).toBe('');
    });

    it('returns empty string for invalid date', () => {
      expect(formatRelativeTime('invalid')).toBe('');
    });
  });

  describe('formatCurrency', () => {
    it('formats USD by default', () => {
      expect(formatCurrency(1234.56)).toBe('$1,234.56');
    });

    it('formats EUR currency', () => {
      const result = formatCurrency(1234.56, 'EUR', 'de-DE');
      expect(result).toMatch(/1\.234,56.*€/);
    });

    it('formats JPY currency without decimals', () => {
      expect(formatCurrency(1234, 'JPY', 'ja-JP')).toBe('￥1,234');
    });

    it('returns empty string for null', () => {
      expect(formatCurrency(null)).toBe('');
    });

    it('returns empty string for undefined', () => {
      expect(formatCurrency(undefined)).toBe('');
    });

    it('handles zero', () => {
      expect(formatCurrency(0)).toBe('$0.00');
    });

    it('handles negative numbers', () => {
      expect(formatCurrency(-1234.56)).toBe('-$1,234.56');
    });
  });

  describe('formatNumber', () => {
    it('formats number with default options', () => {
      expect(formatNumber(1234567.89)).toBe('1,234,567.89');
    });

    it('formats with custom decimal places', () => {
      expect(formatNumber(1234.5, { minimumFractionDigits: 2, maximumFractionDigits: 2 })).toBe(
        '1,234.50'
      );
    });

    it('returns empty string for null', () => {
      expect(formatNumber(null)).toBe('');
    });

    it('returns empty string for undefined', () => {
      expect(formatNumber(undefined)).toBe('');
    });

    it('handles zero', () => {
      expect(formatNumber(0)).toBe('0');
    });
  });

  describe('formatPercentage', () => {
    it('formats decimal as percentage', () => {
      expect(formatPercentage(0.1234)).toBe('12%');
    });

    it('formats with decimal places', () => {
      expect(formatPercentage(0.1234, 2)).toBe('12.34%');
    });

    it('handles zero', () => {
      expect(formatPercentage(0)).toBe('0%');
    });

    it('handles 100%', () => {
      expect(formatPercentage(1)).toBe('100%');
    });

    it('returns empty string for null', () => {
      expect(formatPercentage(null)).toBe('');
    });

    it('returns empty string for undefined', () => {
      expect(formatPercentage(undefined)).toBe('');
    });
  });

  describe('formatCompactNumber', () => {
    it('formats thousands', () => {
      expect(formatCompactNumber(1500)).toBe('1.5K');
    });

    it('formats millions', () => {
      expect(formatCompactNumber(2500000)).toBe('2.5M');
    });

    it('formats billions', () => {
      expect(formatCompactNumber(1200000000)).toBe('1.2B');
    });

    it('does not compact small numbers', () => {
      expect(formatCompactNumber(999)).toBe('999');
    });

    it('returns empty string for null', () => {
      expect(formatCompactNumber(null)).toBe('');
    });

    it('returns empty string for undefined', () => {
      expect(formatCompactNumber(undefined)).toBe('');
    });
  });

  describe('formatBytes', () => {
    it('formats bytes', () => {
      expect(formatBytes(500)).toBe('500 Bytes');
    });

    it('formats kilobytes', () => {
      expect(formatBytes(1024)).toBe('1 KB');
    });

    it('formats megabytes', () => {
      expect(formatBytes(1048576)).toBe('1 MB');
    });

    it('formats gigabytes', () => {
      expect(formatBytes(1073741824)).toBe('1 GB');
    });

    it('formats with custom decimals', () => {
      expect(formatBytes(1536, 1)).toBe('1.5 KB');
    });

    it('handles zero bytes', () => {
      expect(formatBytes(0)).toBe('0 Bytes');
    });

    it('returns empty string for null', () => {
      expect(formatBytes(null)).toBe('');
    });

    it('returns empty string for undefined', () => {
      expect(formatBytes(undefined)).toBe('');
    });

    it('handles negative decimals as zero', () => {
      expect(formatBytes(1536, -1)).toBe('2 KB');
    });
  });

  describe('truncateText', () => {
    it('truncates long text', () => {
      expect(truncateText('Hello World!', 8)).toBe('Hello...');
    });

    it('does not truncate short text', () => {
      expect(truncateText('Hello', 10)).toBe('Hello');
    });

    it('uses custom suffix', () => {
      expect(truncateText('Hello World!', 8, '…')).toBe('Hello W…');
    });

    it('returns empty string for null', () => {
      expect(truncateText(null, 10)).toBe('');
    });

    it('returns empty string for undefined', () => {
      expect(truncateText(undefined, 10)).toBe('');
    });

    it('handles exact length', () => {
      expect(truncateText('Hello', 5)).toBe('Hello');
    });
  });

  describe('capitalizeFirst', () => {
    it('capitalizes first letter', () => {
      expect(capitalizeFirst('hello world')).toBe('Hello world');
    });

    it('lowercases rest of string', () => {
      expect(capitalizeFirst('HELLO WORLD')).toBe('Hello world');
    });

    it('handles single character', () => {
      expect(capitalizeFirst('a')).toBe('A');
    });

    it('returns empty string for null', () => {
      expect(capitalizeFirst(null)).toBe('');
    });

    it('returns empty string for undefined', () => {
      expect(capitalizeFirst(undefined)).toBe('');
    });

    it('returns empty string for empty string', () => {
      expect(capitalizeFirst('')).toBe('');
    });
  });

  describe('toTitleCase', () => {
    it('converts to title case', () => {
      expect(toTitleCase('hello world')).toBe('Hello World');
    });

    it('handles uppercase input', () => {
      expect(toTitleCase('HELLO WORLD')).toBe('Hello World');
    });

    it('handles mixed case input', () => {
      expect(toTitleCase('hElLo WoRlD')).toBe('Hello World');
    });

    it('returns empty string for null', () => {
      expect(toTitleCase(null)).toBe('');
    });

    it('returns empty string for undefined', () => {
      expect(toTitleCase(undefined)).toBe('');
    });

    it('handles single word', () => {
      expect(toTitleCase('hello')).toBe('Hello');
    });
  });

  describe('slugify', () => {
    it('converts text to slug', () => {
      expect(slugify('Hello World')).toBe('hello-world');
    });

    it('removes special characters', () => {
      expect(slugify('Hello, World!')).toBe('hello-world');
    });

    it('handles multiple spaces', () => {
      expect(slugify('Hello   World')).toBe('hello-world');
    });

    it('handles underscores', () => {
      expect(slugify('hello_world_test')).toBe('hello-world-test');
    });

    it('trims leading/trailing hyphens', () => {
      expect(slugify(' Hello World ')).toBe('hello-world');
    });

    it('returns empty string for null', () => {
      expect(slugify(null)).toBe('');
    });

    it('returns empty string for undefined', () => {
      expect(slugify(undefined)).toBe('');
    });

    it('handles numbers', () => {
      expect(slugify('Product 123')).toBe('product-123');
    });

    it('handles already slugified text', () => {
      expect(slugify('hello-world')).toBe('hello-world');
    });
  });
});
