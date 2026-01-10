import { format, formatDistanceToNow, isValid, parseISO } from 'date-fns';

export function formatDate(
  date: Date | string | null | undefined,
  formatString: string = 'MMM d, yyyy'
): string {
  if (!date) {
    return '';
  }

  const parsedDate = typeof date === 'string' ? parseISO(date) : date;

  if (!isValid(parsedDate)) {
    return '';
  }

  return format(parsedDate, formatString);
}

export function formatDateTime(
  date: Date | string | null | undefined,
  formatString: string = 'MMM d, yyyy h:mm a'
): string {
  return formatDate(date, formatString);
}

export function formatRelativeTime(
  date: Date | string | null | undefined,
  options: { addSuffix?: boolean } = { addSuffix: true }
): string {
  if (!date) {
    return '';
  }

  const parsedDate = typeof date === 'string' ? parseISO(date) : date;

  if (!isValid(parsedDate)) {
    return '';
  }

  return formatDistanceToNow(parsedDate, options);
}

export function formatCurrency(
  amount: number | null | undefined,
  currency: string = 'USD',
  locale: string = 'en-US'
): string {
  if (amount === null || amount === undefined) {
    return '';
  }

  return new Intl.NumberFormat(locale, {
    style: 'currency',
    currency,
  }).format(amount);
}

export function formatNumber(
  value: number | null | undefined,
  options: Intl.NumberFormatOptions = {},
  locale: string = 'en-US'
): string {
  if (value === null || value === undefined) {
    return '';
  }

  return new Intl.NumberFormat(locale, options).format(value);
}

export function formatPercentage(
  value: number | null | undefined,
  decimals: number = 0,
  locale: string = 'en-US'
): string {
  if (value === null || value === undefined) {
    return '';
  }

  return new Intl.NumberFormat(locale, {
    style: 'percent',
    minimumFractionDigits: decimals,
    maximumFractionDigits: decimals,
  }).format(value);
}

export function formatCompactNumber(
  value: number | null | undefined,
  locale: string = 'en-US'
): string {
  if (value === null || value === undefined) {
    return '';
  }

  return new Intl.NumberFormat(locale, {
    notation: 'compact',
    compactDisplay: 'short',
  }).format(value);
}

export function formatBytes(bytes: number | null | undefined, decimals: number = 2): string {
  if (bytes === null || bytes === undefined) {
    return '';
  }

  if (bytes === 0) {
    return '0 Bytes';
  }

  const kilobyte = 1024;
  const decimalPlaces = decimals < 0 ? 0 : decimals;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];

  const index = Math.floor(Math.log(bytes) / Math.log(kilobyte));
  const clampedIndex = Math.min(index, sizes.length - 1);

  return `${parseFloat((bytes / Math.pow(kilobyte, clampedIndex)).toFixed(decimalPlaces))} ${sizes[clampedIndex]}`;
}

export function truncateText(
  text: string | null | undefined,
  maxLength: number,
  suffix: string = '...'
): string {
  if (!text) {
    return '';
  }

  if (text.length <= maxLength) {
    return text;
  }

  return text.slice(0, maxLength - suffix.length) + suffix;
}

export function capitalizeFirst(text: string | null | undefined): string {
  if (!text) {
    return '';
  }

  return text.charAt(0).toUpperCase() + text.slice(1).toLowerCase();
}

export function toTitleCase(text: string | null | undefined): string {
  if (!text) {
    return '';
  }

  return text
    .toLowerCase()
    .split(' ')
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(' ');
}

export function slugify(text: string | null | undefined): string {
  if (!text) {
    return '';
  }

  return text
    .toLowerCase()
    .trim()
    .replace(/[^\w\s-]/g, '')
    .replace(/[\s_-]+/g, '-')
    .replace(/^-+|-+$/g, '');
}
