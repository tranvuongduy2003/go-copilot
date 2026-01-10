import { LoginForm } from '@/features/auth/components/login-form';
import { useAuthStore } from '@/stores';
import { server } from '@/test/mocks/server';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { http, HttpResponse } from 'msw';
import type { ReactNode } from 'react';
import { MemoryRouter, Route, Routes } from 'react-router';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1';

vi.mock('sonner', () => ({
  toast: {
    success: vi.fn(),
    error: vi.fn(),
  },
}));

function DashboardPage() {
  return <div>Dashboard</div>;
}

function createTestQueryClient() {
  return new QueryClient({
    defaultOptions: {
      queries: { retry: false, gcTime: 0, staleTime: 0 },
      mutations: { retry: false },
    },
  });
}

function renderWithProviders(
  ui: ReactNode,
  { initialEntries = ['/login'] }: { initialEntries?: string[] } = {}
) {
  const queryClient = createTestQueryClient();

  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter initialEntries={initialEntries}>
        <Routes>
          <Route path="/login" element={ui} />
          <Route path="/" element={<DashboardPage />} />
          <Route path="/dashboard" element={<DashboardPage />} />
        </Routes>
      </MemoryRouter>
    </QueryClientProvider>
  );
}

describe('Auth Integration - Login Flow', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    useAuthStore.getState().clearAuth();
  });

  it('completes full login flow successfully', async () => {
    const user = userEvent.setup();

    server.use(
      http.post(`${API_BASE_URL}/auth/login`, async () => {
        return HttpResponse.json({
          data: {
            accessToken: 'test-access-token',
            refreshToken: 'test-refresh-token',
            expiresIn: 3600,
            user: {
              id: '1',
              email: 'test@example.com',
              fullName: 'Test User',
              status: 'active',
              roles: [{ id: '1', name: 'admin', displayName: 'Admin', description: '' }],
              permissions: ['users:read', 'users:write'],
              createdAt: new Date().toISOString(),
              updatedAt: new Date().toISOString(),
            },
          },
        });
      })
    );

    renderWithProviders(<LoginForm />);

    await user.type(screen.getByPlaceholderText(/enter your email/i), 'test@example.com');
    await user.type(screen.getByPlaceholderText(/enter your password/i), 'password123');
    await user.click(screen.getByRole('button', { name: /sign in/i }));

    await waitFor(() => {
      expect(screen.getByText('Dashboard')).toBeInTheDocument();
    });

    const authState = useAuthStore.getState();
    expect(authState.user).toBeTruthy();
    expect(authState.user?.email).toBe('test@example.com');
  });

  it('shows error message on invalid credentials', async () => {
    const user = userEvent.setup();

    server.use(
      http.post(`${API_BASE_URL}/auth/login`, () => {
        return HttpResponse.json(
          { error: { code: 'UNAUTHORIZED', message: 'Invalid credentials' } },
          { status: 401 }
        );
      })
    );

    renderWithProviders(<LoginForm />);

    await user.type(screen.getByPlaceholderText(/enter your email/i), 'wrong@example.com');
    await user.type(screen.getByPlaceholderText(/enter your password/i), 'wrongpassword');
    await user.click(screen.getByRole('button', { name: /sign in/i }));

    await waitFor(() => {
      expect(screen.getByRole('alert')).toBeInTheDocument();
    });

    const authState = useAuthStore.getState();
    expect(authState.user).toBeNull();
  });

  it('validates form before submission', async () => {
    const user = userEvent.setup();
    renderWithProviders(<LoginForm />);

    await user.click(screen.getByRole('button', { name: /sign in/i }));

    await waitFor(() => {
      expect(screen.getByText(/email is required/i)).toBeInTheDocument();
    });
  });
});

describe('Auth Integration - Logout Flow', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    useAuthStore.getState().setUser({
      id: '1',
      email: 'test@example.com',
      fullName: 'Test User',
      status: 'active',
      roles: [],
      permissions: [],
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    });
  });

  it('clears auth state on logout', () => {
    // Verify user is set initially
    expect(useAuthStore.getState().user).toBeTruthy();

    useAuthStore.getState().clearAuth();

    const authState = useAuthStore.getState();
    expect(authState.user).toBeNull();
    expect(authState.status).toBe('unauthenticated');
  });
});
