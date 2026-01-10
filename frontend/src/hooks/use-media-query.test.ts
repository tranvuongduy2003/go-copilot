import { act, renderHook } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import {
  useIsDesktop,
  useIsMobile,
  useIsTablet,
  useMediaQuery,
  usePrefersDarkMode,
  usePrefersReducedMotion,
} from './use-media-query';

const createMatchMedia = (matches: boolean) => {
  const listeners: Array<(event: MediaQueryListEvent) => void> = [];

  const mediaQueryList = {
    matches,
    media: '',
    onchange: null,
    addListener: vi.fn((callback) => listeners.push(callback)),
    removeListener: vi.fn((callback) => {
      const index = listeners.indexOf(callback);
      if (index > -1) listeners.splice(index, 1);
    }),
    addEventListener: vi.fn((_, callback) => listeners.push(callback)),
    removeEventListener: vi.fn((_, callback) => {
      const index = listeners.indexOf(callback);
      if (index > -1) listeners.splice(index, 1);
    }),
    dispatchEvent: vi.fn((event: MediaQueryListEvent) => {
      listeners.forEach((listener) => listener(event));
      return true;
    }),
    triggerChange: (newMatches: boolean) => {
      mediaQueryList.matches = newMatches;
      listeners.forEach((listener) => listener({ matches: newMatches } as MediaQueryListEvent));
    },
  };

  return mediaQueryList;
};

describe('useMediaQuery', () => {
  let originalMatchMedia: typeof window.matchMedia;
  let mockMatchMedia: ReturnType<typeof createMatchMedia>;

  beforeEach(() => {
    originalMatchMedia = window.matchMedia;
    mockMatchMedia = createMatchMedia(false);
    window.matchMedia = vi.fn(() => mockMatchMedia) as typeof window.matchMedia;
  });

  afterEach(() => {
    window.matchMedia = originalMatchMedia;
  });

  it('returns false when media query does not match', () => {
    mockMatchMedia = createMatchMedia(false);
    window.matchMedia = vi.fn(() => mockMatchMedia) as typeof window.matchMedia;

    const { result } = renderHook(() => useMediaQuery('(min-width: 768px)'));

    expect(result.current).toBe(false);
  });

  it('returns true when media query matches', () => {
    mockMatchMedia = createMatchMedia(true);
    window.matchMedia = vi.fn(() => mockMatchMedia) as typeof window.matchMedia;

    const { result } = renderHook(() => useMediaQuery('(min-width: 768px)'));

    expect(result.current).toBe(true);
  });

  it('updates when media query changes', () => {
    mockMatchMedia = createMatchMedia(false);
    window.matchMedia = vi.fn(() => mockMatchMedia) as typeof window.matchMedia;

    const { result } = renderHook(() => useMediaQuery('(min-width: 768px)'));

    expect(result.current).toBe(false);

    act(() => {
      mockMatchMedia.triggerChange(true);
    });

    expect(result.current).toBe(true);
  });

  it('cleans up event listener on unmount', () => {
    mockMatchMedia = createMatchMedia(false);
    window.matchMedia = vi.fn(() => mockMatchMedia) as typeof window.matchMedia;

    const { unmount } = renderHook(() => useMediaQuery('(min-width: 768px)'));

    unmount();

    expect(mockMatchMedia.removeEventListener).toHaveBeenCalled();
  });

  it('updates when query prop changes', () => {
    const matchMediaImpl = vi.fn((query: string) => {
      const isMaxWidth = query.includes('max-width');
      return createMatchMedia(isMaxWidth);
    });
    window.matchMedia = matchMediaImpl as typeof window.matchMedia;

    const { result, rerender } = renderHook(({ query }) => useMediaQuery(query), {
      initialProps: { query: '(min-width: 768px)' },
    });

    expect(result.current).toBe(false);

    rerender({ query: '(max-width: 767px)' });

    expect(result.current).toBe(true);
  });
});

describe('useIsMobile', () => {
  let originalMatchMedia: typeof window.matchMedia;

  beforeEach(() => {
    originalMatchMedia = window.matchMedia;
  });

  afterEach(() => {
    window.matchMedia = originalMatchMedia;
  });

  it('returns true when screen width is less than 768px', () => {
    const mockMatchMedia = createMatchMedia(true);
    window.matchMedia = vi.fn(() => mockMatchMedia) as typeof window.matchMedia;

    const { result } = renderHook(() => useIsMobile());

    expect(result.current).toBe(true);
  });

  it('returns false when screen width is 768px or more', () => {
    const mockMatchMedia = createMatchMedia(false);
    window.matchMedia = vi.fn(() => mockMatchMedia) as typeof window.matchMedia;

    const { result } = renderHook(() => useIsMobile());

    expect(result.current).toBe(false);
  });
});

describe('useIsTablet', () => {
  let originalMatchMedia: typeof window.matchMedia;

  beforeEach(() => {
    originalMatchMedia = window.matchMedia;
  });

  afterEach(() => {
    window.matchMedia = originalMatchMedia;
  });

  it('returns true when screen width is between 768px and 1023px', () => {
    const mockMatchMedia = createMatchMedia(true);
    window.matchMedia = vi.fn(() => mockMatchMedia) as typeof window.matchMedia;

    const { result } = renderHook(() => useIsTablet());

    expect(result.current).toBe(true);
  });

  it('returns false when screen width is outside tablet range', () => {
    const mockMatchMedia = createMatchMedia(false);
    window.matchMedia = vi.fn(() => mockMatchMedia) as typeof window.matchMedia;

    const { result } = renderHook(() => useIsTablet());

    expect(result.current).toBe(false);
  });
});

describe('useIsDesktop', () => {
  let originalMatchMedia: typeof window.matchMedia;

  beforeEach(() => {
    originalMatchMedia = window.matchMedia;
  });

  afterEach(() => {
    window.matchMedia = originalMatchMedia;
  });

  it('returns true when screen width is 1024px or more', () => {
    const mockMatchMedia = createMatchMedia(true);
    window.matchMedia = vi.fn(() => mockMatchMedia) as typeof window.matchMedia;

    const { result } = renderHook(() => useIsDesktop());

    expect(result.current).toBe(true);
  });

  it('returns false when screen width is less than 1024px', () => {
    const mockMatchMedia = createMatchMedia(false);
    window.matchMedia = vi.fn(() => mockMatchMedia) as typeof window.matchMedia;

    const { result } = renderHook(() => useIsDesktop());

    expect(result.current).toBe(false);
  });
});

describe('usePrefersDarkMode', () => {
  let originalMatchMedia: typeof window.matchMedia;

  beforeEach(() => {
    originalMatchMedia = window.matchMedia;
  });

  afterEach(() => {
    window.matchMedia = originalMatchMedia;
  });

  it('returns true when user prefers dark mode', () => {
    const mockMatchMedia = createMatchMedia(true);
    window.matchMedia = vi.fn(() => mockMatchMedia) as typeof window.matchMedia;

    const { result } = renderHook(() => usePrefersDarkMode());

    expect(result.current).toBe(true);
  });

  it('returns false when user prefers light mode', () => {
    const mockMatchMedia = createMatchMedia(false);
    window.matchMedia = vi.fn(() => mockMatchMedia) as typeof window.matchMedia;

    const { result } = renderHook(() => usePrefersDarkMode());

    expect(result.current).toBe(false);
  });
});

describe('usePrefersReducedMotion', () => {
  let originalMatchMedia: typeof window.matchMedia;

  beforeEach(() => {
    originalMatchMedia = window.matchMedia;
  });

  afterEach(() => {
    window.matchMedia = originalMatchMedia;
  });

  it('returns true when user prefers reduced motion', () => {
    const mockMatchMedia = createMatchMedia(true);
    window.matchMedia = vi.fn(() => mockMatchMedia) as typeof window.matchMedia;

    const { result } = renderHook(() => usePrefersReducedMotion());

    expect(result.current).toBe(true);
  });

  it('returns false when user does not prefer reduced motion', () => {
    const mockMatchMedia = createMatchMedia(false);
    window.matchMedia = vi.fn(() => mockMatchMedia) as typeof window.matchMedia;

    const { result } = renderHook(() => usePrefersReducedMotion());

    expect(result.current).toBe(false);
  });
});
