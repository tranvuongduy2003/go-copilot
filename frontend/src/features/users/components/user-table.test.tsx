import { render, screen } from '@/test/utils';
import type { User } from '@/types/user';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import { UserTable } from './user-table';

const mockUsers: User[] = [
  {
    id: '1',
    email: 'admin@example.com',
    fullName: 'Admin User',
    status: 'active',
    roles: [
      {
        id: '1',
        name: 'admin',
        displayName: 'Administrator',
        description: '',
        isSystem: true,
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
      },
    ],
    permissions: ['users:read', 'users:write'],
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z',
  },
  {
    id: '2',
    email: 'user@example.com',
    fullName: 'Regular User',
    status: 'pending',
    roles: [
      {
        id: '2',
        name: 'user',
        displayName: 'User',
        description: '',
        isSystem: false,
        createdAt: '2024-01-02T00:00:00Z',
        updatedAt: '2024-01-02T00:00:00Z',
      },
    ],
    permissions: ['users:read'],
    createdAt: '2024-01-02T00:00:00Z',
    updatedAt: '2024-01-02T00:00:00Z',
  },
  {
    id: '3',
    email: 'banned@example.com',
    fullName: 'Banned User',
    status: 'banned',
    roles: [],
    permissions: [],
    createdAt: '2024-01-03T00:00:00Z',
    updatedAt: '2024-01-03T00:00:00Z',
  },
];

const defaultProps = {
  users: mockUsers,
  selectedUsers: [],
  onSelectUser: vi.fn(),
  onSelectAll: vi.fn(),
  onEdit: vi.fn(),
  onDelete: vi.fn(),
  onActivate: vi.fn(),
  onDeactivate: vi.fn(),
};

describe('UserTable', () => {
  it('renders table with user data', () => {
    render(<UserTable {...defaultProps} />);

    expect(screen.getByRole('table')).toBeInTheDocument();
    expect(screen.getByText('Admin User')).toBeInTheDocument();
    expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    expect(screen.getByText('Regular User')).toBeInTheDocument();
    expect(screen.getByText('user@example.com')).toBeInTheDocument();
  });

  it('renders column headers', () => {
    render(<UserTable {...defaultProps} />);

    expect(screen.getByText('User')).toBeInTheDocument();
    expect(screen.getByText('Status')).toBeInTheDocument();
    expect(screen.getByText('Roles')).toBeInTheDocument();
    expect(screen.getByText('Created')).toBeInTheDocument();
  });

  it('displays status badges with correct variants', () => {
    render(<UserTable {...defaultProps} />);

    const activeStatuses = screen.getAllByText('Active');
    const pendingStatus = screen.getByText('Pending');
    const bannedStatus = screen.getByText('Banned');

    expect(activeStatuses.length).toBeGreaterThan(0);
    expect(pendingStatus).toBeInTheDocument();
    expect(bannedStatus).toBeInTheDocument();
  });

  it('shows loading state', () => {
    render(<UserTable {...defaultProps} users={[]} isLoading={true} />);

    expect(screen.getByText(/loading users/i)).toBeInTheDocument();
  });

  it('shows empty state when no users', () => {
    render(<UserTable {...defaultProps} users={[]} />);

    expect(screen.getByText(/no users found/i)).toBeInTheDocument();
  });

  it('handles row selection', async () => {
    const user = userEvent.setup();
    const onSelectUser = vi.fn();
    render(<UserTable {...defaultProps} onSelectUser={onSelectUser} />);

    const checkboxes = screen.getAllByRole('checkbox');
    await user.click(checkboxes[1]);

    expect(onSelectUser).toHaveBeenCalledWith('1');
  });

  it('handles select all', async () => {
    const user = userEvent.setup();
    const onSelectAll = vi.fn();
    render(<UserTable {...defaultProps} onSelectAll={onSelectAll} />);

    const selectAllCheckbox = screen.getAllByRole('checkbox')[0];
    await user.click(selectAllCheckbox);

    expect(onSelectAll).toHaveBeenCalled();
  });

  it('shows all selected state when all users are selected', () => {
    render(<UserTable {...defaultProps} selectedUsers={['1', '2', '3']} />);

    const selectAllCheckbox = screen.getAllByRole('checkbox')[0];
    expect(selectAllCheckbox).toBeChecked();
  });

  it('renders user avatars with initials', () => {
    render(<UserTable {...defaultProps} />);

    expect(screen.getByText('AU')).toBeInTheDocument();
    expect(screen.getByText('RU')).toBeInTheDocument();
    expect(screen.getByText('BU')).toBeInTheDocument();
  });

  it('displays roles for each user', () => {
    render(<UserTable {...defaultProps} />);

    expect(screen.getByText('admin')).toBeInTheDocument();
    expect(screen.getByText('user')).toBeInTheDocument();
  });

  it('formats dates correctly', () => {
    render(<UserTable {...defaultProps} />);

    expect(screen.getByText('Jan 1, 2024')).toBeInTheDocument();
    expect(screen.getByText('Jan 2, 2024')).toBeInTheDocument();
    expect(screen.getByText('Jan 3, 2024')).toBeInTheDocument();
  });

  it('marks selected rows with data-state', () => {
    render(<UserTable {...defaultProps} selectedUsers={['1']} />);

    const rows = screen.getAllByRole('row');
    const dataRow = rows.find((row) => row.textContent?.includes('Admin User'));
    expect(dataRow).toHaveAttribute('data-state', 'selected');
  });
});
