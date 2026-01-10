import { AUTH_STORAGE_KEY } from '@/constants';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const mockLocalStorage = (() => {
  let store: Record<string, string> = {};
  return {
    getItem: vi.fn((key: string) => store[key] || null),
    setItem: vi.fn((key: string, value: string) => {
      store[key] = value;
    }),
    removeItem: vi.fn((key: string) => {
      delete store[key];
    }),
    clear: vi.fn(() => {
      store = {};
    }),
    get length() {
      return Object.keys(store).length;
    },
    key: vi.fn((index: number) => Object.keys(store)[index] || null),
  };
})();

Object.defineProperty(global, 'localStorage', {
  value: mockLocalStorage,
  writable: true,
});

import { tokenService } from './token-service';

describe('TokenService', () => {
  beforeEach(() => {
    mockLocalStorage.clear();
    vi.clearAllMocks();
  });

  describe('setTokens', () => {
    it('stores tokens in localStorage', () => {
      tokenService.setTokens('access-token-123', 'refresh-token-456');

      expect(mockLocalStorage.setItem).toHaveBeenCalledWith(
        AUTH_STORAGE_KEY,
        JSON.stringify({
          accessToken: 'access-token-123',
          refreshToken: 'refresh-token-456',
        })
      );
    });
  });

  describe('getAccessToken', () => {
    it('retrieves access token from localStorage', () => {
      mockLocalStorage.setItem(
        AUTH_STORAGE_KEY,
        JSON.stringify({
          accessToken: 'access-token-123',
          refreshToken: 'refresh-token-456',
        })
      );

      expect(tokenService.getAccessToken()).toBe('access-token-123');
    });

    it('returns null when no tokens stored', () => {
      expect(tokenService.getAccessToken()).toBeNull();
    });

    it('returns null for invalid JSON in storage', () => {
      mockLocalStorage.setItem(AUTH_STORAGE_KEY, 'invalid-json');

      expect(tokenService.getAccessToken()).toBeNull();
    });
  });

  describe('getRefreshToken', () => {
    it('retrieves refresh token from localStorage', () => {
      mockLocalStorage.setItem(
        AUTH_STORAGE_KEY,
        JSON.stringify({
          accessToken: 'access-token-123',
          refreshToken: 'refresh-token-456',
        })
      );

      expect(tokenService.getRefreshToken()).toBe('refresh-token-456');
    });

    it('returns null when no tokens stored', () => {
      expect(tokenService.getRefreshToken()).toBeNull();
    });
  });

  describe('clearTokens', () => {
    it('removes tokens from localStorage', () => {
      mockLocalStorage.setItem(
        AUTH_STORAGE_KEY,
        JSON.stringify({
          accessToken: 'access-token-123',
          refreshToken: 'refresh-token-456',
        })
      );

      tokenService.clearTokens();

      expect(mockLocalStorage.removeItem).toHaveBeenCalledWith(AUTH_STORAGE_KEY);
    });
  });

  describe('hasTokens', () => {
    it('returns true when both tokens exist', () => {
      mockLocalStorage.setItem(
        AUTH_STORAGE_KEY,
        JSON.stringify({
          accessToken: 'access-token-123',
          refreshToken: 'refresh-token-456',
        })
      );

      expect(tokenService.hasTokens()).toBe(true);
    });

    it('returns false when no tokens exist', () => {
      expect(tokenService.hasTokens()).toBe(false);
    });

    it('returns false when only access token exists', () => {
      mockLocalStorage.setItem(
        AUTH_STORAGE_KEY,
        JSON.stringify({
          accessToken: 'access-token-123',
          refreshToken: '',
        })
      );

      expect(tokenService.hasTokens()).toBe(false);
    });

    it('returns false when only refresh token exists', () => {
      mockLocalStorage.setItem(
        AUTH_STORAGE_KEY,
        JSON.stringify({
          accessToken: '',
          refreshToken: 'refresh-token-456',
        })
      );

      expect(tokenService.hasTokens()).toBe(false);
    });
  });
});
