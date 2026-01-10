import { API_BASE_URL } from '@/constants';
import { http, HttpResponse } from 'msw';

const baseUrl = API_BASE_URL;

const mockUsers = [
  {
    id: '1',
    email: 'admin@example.com',
    fullName: 'Admin User',
    status: 'active',
    roles: [{ id: '1', name: 'admin', displayName: 'Administrator', description: '' }],
    permissions: ['users:read', 'users:write'],
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  },
  {
    id: '2',
    email: 'user@example.com',
    fullName: 'Regular User',
    status: 'active',
    roles: [{ id: '2', name: 'user', displayName: 'User', description: '' }],
    permissions: ['users:read'],
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  },
];

export const userHandlers = [
  http.get(`${baseUrl}/users`, ({ request }) => {
    const url = new URL(request.url);
    const page = Number(url.searchParams.get('page')) || 1;
    const pageSize = Number(url.searchParams.get('pageSize')) || 10;

    return HttpResponse.json({
      data: mockUsers,
      meta: {
        page,
        pageSize,
        total: mockUsers.length,
        totalPages: 1,
      },
    });
  }),

  http.get(`${baseUrl}/users/:id`, ({ params }) => {
    const user = mockUsers.find((user) => user.id === params.id);

    if (!user) {
      return HttpResponse.json(
        { error: { code: 'NOT_FOUND', message: 'User not found' } },
        { status: 404 }
      );
    }

    return HttpResponse.json({ data: user });
  }),

  http.post(`${baseUrl}/users`, async ({ request }) => {
    const body = (await request.json()) as { email: string; fullName: string; password: string };

    const newUser = {
      id: String(mockUsers.length + 1),
      email: body.email,
      fullName: body.fullName,
      status: 'active',
      roles: [],
      permissions: [],
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    };

    return HttpResponse.json({ data: newUser }, { status: 201 });
  }),

  http.put(`${baseUrl}/users/:id`, async ({ params, request }) => {
    const body = (await request.json()) as { fullName?: string; email?: string };
    const user = mockUsers.find((user) => user.id === params.id);

    if (!user) {
      return HttpResponse.json(
        { error: { code: 'NOT_FOUND', message: 'User not found' } },
        { status: 404 }
      );
    }

    return HttpResponse.json({
      data: { ...user, ...body, updatedAt: new Date().toISOString() },
    });
  }),

  http.delete(`${baseUrl}/users/:id`, ({ params }) => {
    const user = mockUsers.find((user) => user.id === params.id);

    if (!user) {
      return HttpResponse.json(
        { error: { code: 'NOT_FOUND', message: 'User not found' } },
        { status: 404 }
      );
    }

    return HttpResponse.json({ data: null }, { status: 204 });
  }),

  http.post(`${baseUrl}/users/:id/activate`, ({ params }) => {
    const user = mockUsers.find((user) => user.id === params.id);

    if (!user) {
      return HttpResponse.json(
        { error: { code: 'NOT_FOUND', message: 'User not found' } },
        { status: 404 }
      );
    }

    return HttpResponse.json({
      data: { ...user, status: 'active', updatedAt: new Date().toISOString() },
    });
  }),

  http.post(`${baseUrl}/users/:id/deactivate`, ({ params }) => {
    const user = mockUsers.find((user) => user.id === params.id);

    if (!user) {
      return HttpResponse.json(
        { error: { code: 'NOT_FOUND', message: 'User not found' } },
        { status: 404 }
      );
    }

    return HttpResponse.json({
      data: { ...user, status: 'inactive', updatedAt: new Date().toISOString() },
    });
  }),
];
