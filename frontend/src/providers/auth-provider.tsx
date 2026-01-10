import { Loading } from '@/components/ui/spinner';
import { useCurrentUser } from '@/features/auth/api/auth.queries';
import { tokenService } from '@/lib/api/token-service';
import { useAuthStore } from '@/stores/auth-store';
import { useEffect } from 'react';

interface AuthProviderProps {
  children: React.ReactNode;
}

export function AuthProvider({ children }: AuthProviderProps) {
  const { status, setUser, setStatus, clearAuth } = useAuthStore();
  const hasToken = !!tokenService.getAccessToken();

  const { data, isLoading, isError } = useCurrentUser();

  useEffect(() => {
    if (!hasToken) {
      setStatus('unauthenticated');
      return;
    }

    if (isLoading) {
      setStatus('loading');
      return;
    }

    if (isError) {
      clearAuth();
      return;
    }

    if (data) {
      setUser(data);
      setStatus('authenticated');
    }
  }, [hasToken, isLoading, isError, data, setUser, setStatus, clearAuth]);

  if (status === 'loading' && hasToken) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <Loading text="Initializing..." />
      </div>
    );
  }

  return <>{children}</>;
}
