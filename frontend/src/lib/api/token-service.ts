import { AUTH_STORAGE_KEY } from '@/constants';

interface TokenStorage {
  accessToken: string;
  refreshToken: string;
}

class TokenService {
  private readonly storageKey: string;

  constructor(storageKey: string = AUTH_STORAGE_KEY) {
    this.storageKey = storageKey;
  }

  public getAccessToken(): string | null {
    try {
      const storage = this.getStorage();
      return storage?.accessToken || null;
    } catch {
      return null;
    }
  }

  public getRefreshToken(): string | null {
    try {
      const storage = this.getStorage();
      return storage?.refreshToken || null;
    } catch {
      return null;
    }
  }

  public setTokens(accessToken: string, refreshToken: string): void {
    try {
      const storage: TokenStorage = { accessToken, refreshToken };
      localStorage.setItem(this.storageKey, JSON.stringify(storage));
    } catch (error) {
      console.error('Failed to store tokens:', error);
    }
  }

  public clearTokens(): void {
    try {
      localStorage.removeItem(this.storageKey);
    } catch (error) {
      console.error('Failed to clear tokens:', error);
    }
  }

  public hasTokens(): boolean {
    return !!this.getAccessToken() && !!this.getRefreshToken();
  }

  private getStorage(): TokenStorage | null {
    try {
      const stored = localStorage.getItem(this.storageKey);
      if (!stored) {
        return null;
      }
      return JSON.parse(stored) as TokenStorage;
    } catch {
      return null;
    }
  }
}

export const tokenService = new TokenService();
