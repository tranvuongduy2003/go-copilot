import { API_BASE_URL } from '@/constants';
import { http, HttpResponse } from 'msw';

const baseUrl = API_BASE_URL;

export const authHandlers = [
  http.post(`${baseUrl}/auth/login`, async ({ request }) => {
    const body = (await request.json()) as { email: string; password: string };

    if (body.email === 'test@example.com' && body.password === 'password123') {
      return HttpResponse.json({
        data: {
          accessToken: 'mock-access-token',
          refreshToken: 'mock-refresh-token',
          expiresIn: 3600,
          user: {
            id: '1',
            email: 'test@example.com',
            fullName: 'Test User',
            status: 'active',
            roles: [{ id: '1', name: 'admin', displayName: 'Administrator', description: '' }],
            permissions: ['users:read', 'users:write'],
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
          },
        },
      });
    }

    return HttpResponse.json(
      { error: { code: 'UNAUTHORIZED', message: 'Invalid credentials' } },
      { status: 401 }
    );
  }),

  http.post(`${baseUrl}/auth/register`, async ({ request }) => {
    const body = (await request.json()) as { email: string; fullName: string; password: string };

    return HttpResponse.json({
      data: {
        accessToken: 'mock-access-token',
        refreshToken: 'mock-refresh-token',
        expiresIn: 3600,
        user: {
          id: '2',
          email: body.email,
          fullName: body.fullName,
          status: 'pending',
          roles: [],
          permissions: [],
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString(),
        },
      },
    });
  }),

  http.post(`${baseUrl}/auth/logout`, () => {
    return HttpResponse.json({ data: null });
  }),

  http.post(`${baseUrl}/auth/refresh`, () => {
    return HttpResponse.json({
      data: {
        accessToken: 'new-mock-access-token',
        refreshToken: 'new-mock-refresh-token',
        expiresIn: 3600,
      },
    });
  }),

  http.get(`${baseUrl}/auth/me`, ({ request }) => {
    const authHeader = request.headers.get('Authorization');

    if (!authHeader || !authHeader.startsWith('Bearer ')) {
      return HttpResponse.json(
        { error: { code: 'UNAUTHORIZED', message: 'Not authenticated' } },
        { status: 401 }
      );
    }

    return HttpResponse.json({
      data: {
        id: '1',
        email: 'test@example.com',
        fullName: 'Test User',
        status: 'active',
        roles: [{ id: '1', name: 'admin', displayName: 'Administrator', description: '' }],
        permissions: ['users:read', 'users:write', 'roles:read'],
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      },
    });
  }),

  http.post(`${baseUrl}/auth/forgot-password`, () => {
    return HttpResponse.json({ data: null });
  }),

  http.post(`${baseUrl}/auth/reset-password`, () => {
    return HttpResponse.json({ data: null });
  }),

  http.post(`${baseUrl}/auth/change-password`, () => {
    return HttpResponse.json({ data: null });
  }),

  http.get(`${baseUrl}/auth/sessions`, () => {
    return HttpResponse.json({
      data: [
        {
          id: '1',
          userAgent: 'Mozilla/5.0',
          ipAddress: '127.0.0.1',
          createdAt: new Date().toISOString(),
          lastAccessAt: new Date().toISOString(),
          isCurrent: true,
        },
      ],
    });
  }),

  http.delete(`${baseUrl}/auth/sessions/:id/revoke`, () => {
    return HttpResponse.json({ data: null });
  }),
];
