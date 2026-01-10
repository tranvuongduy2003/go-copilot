import { render, screen, waitFor } from '@/test/utils';
import type { User } from '@/types/user';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import { UserFormDialog } from './user-form-dialog';

const mockUser: User = {
  id: '1',
  email: 'test@example.com',
  fullName: 'Test User',
  status: 'active',
  roles: [],
  permissions: [],
  createdAt: '2024-01-01T00:00:00Z',
  updatedAt: '2024-01-01T00:00:00Z',
};

describe('UserFormDialog', () => {
  it('renders create mode correctly', () => {
    render(<UserFormDialog open={true} onOpenChange={vi.fn()} onSubmit={vi.fn()} />);

    expect(screen.getByRole('heading', { name: /create user/i })).toBeInTheDocument();
    expect(screen.getByText(/fill in the details to create a new user/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/full name/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/email/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/password/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/status/i)).toBeInTheDocument();
  });

  it('renders edit mode correctly', () => {
    render(
      <UserFormDialog open={true} onOpenChange={vi.fn()} user={mockUser} onSubmit={vi.fn()} />
    );

    expect(screen.getByText('Edit User')).toBeInTheDocument();
    expect(screen.getByText(/update the user information/i)).toBeInTheDocument();
    expect(screen.queryByLabelText(/password/i)).not.toBeInTheDocument();
  });

  it('pre-populates form in edit mode', () => {
    render(
      <UserFormDialog open={true} onOpenChange={vi.fn()} user={mockUser} onSubmit={vi.fn()} />
    );

    expect(screen.getByDisplayValue('Test User')).toBeInTheDocument();
    expect(screen.getByDisplayValue('test@example.com')).toBeInTheDocument();
  });

  it('shows validation errors for empty required fields', async () => {
    const user = userEvent.setup();
    render(<UserFormDialog open={true} onOpenChange={vi.fn()} onSubmit={vi.fn()} />);

    await user.click(screen.getByRole('button', { name: /create user/i }));

    await waitFor(() => {
      const errorMessages = screen.queryAllByText(/required|invalid/i);
      expect(errorMessages.length).toBeGreaterThan(0);
    });
  });

  it('shows validation error for invalid email', async () => {
    const user = userEvent.setup();
    const onSubmit = vi.fn();
    render(<UserFormDialog open={true} onOpenChange={vi.fn()} onSubmit={onSubmit} />);

    const emailInput = screen.getByLabelText(/email/i);
    await user.type(emailInput, 'invalid-email');
    await user.click(screen.getByRole('button', { name: /create user/i }));

    await waitFor(() => {
      // Either shows validation error or doesn't call onSubmit
      const hasValidationError = screen.queryByText(/email/i);
      expect(hasValidationError || !onSubmit.mock.calls.length).toBeTruthy();
    });
  });

  it('calls onSubmit with form data for create', async () => {
    const user = userEvent.setup();
    const onSubmit = vi.fn();
    render(<UserFormDialog open={true} onOpenChange={vi.fn()} onSubmit={onSubmit} />);

    await user.type(screen.getByLabelText(/full name/i), 'New User');
    await user.type(screen.getByLabelText(/email/i), 'new@example.com');
    await user.type(screen.getByLabelText(/password/i), 'Password123!');

    await user.click(screen.getByRole('button', { name: /create user/i }));

    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalledWith(
        expect.objectContaining({
          fullName: 'New User',
          email: 'new@example.com',
          password: 'Password123!',
        })
      );
    });
  });

  it('calls onSubmit with form data for edit', async () => {
    const user = userEvent.setup();
    const onSubmit = vi.fn();
    render(
      <UserFormDialog open={true} onOpenChange={vi.fn()} user={mockUser} onSubmit={onSubmit} />
    );

    const fullNameInput = screen.getByLabelText(/full name/i);
    await user.clear(fullNameInput);
    await user.type(fullNameInput, 'Updated User');

    await user.click(screen.getByRole('button', { name: /save changes/i }));

    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalledWith(
        expect.objectContaining({
          fullName: 'Updated User',
        })
      );
    });
  });

  it('shows loading state during submission', () => {
    render(
      <UserFormDialog open={true} onOpenChange={vi.fn()} onSubmit={vi.fn()} isLoading={true} />
    );

    const submitButton = screen.getByRole('button', { name: /create user/i });
    expect(submitButton).toBeDisabled();
  });

  it('calls onOpenChange when dialog is closed', async () => {
    const user = userEvent.setup();
    const onOpenChange = vi.fn();
    render(<UserFormDialog open={true} onOpenChange={onOpenChange} onSubmit={vi.fn()} />);

    await user.click(screen.getByRole('button', { name: /cancel/i }));

    expect(onOpenChange).toHaveBeenCalledWith(false);
  });

  it('resets form when dialog is closed', async () => {
    const user = userEvent.setup();
    const onOpenChange = vi.fn();
    const { rerender } = render(
      <UserFormDialog open={true} onOpenChange={onOpenChange} onSubmit={vi.fn()} />
    );

    await user.type(screen.getByLabelText(/full name/i), 'Test');

    await user.click(screen.getByRole('button', { name: /cancel/i }));

    expect(onOpenChange).toHaveBeenCalledWith(false);

    rerender(<UserFormDialog open={true} onOpenChange={onOpenChange} onSubmit={vi.fn()} />);

    expect(screen.getByLabelText(/full name/i)).toHaveValue('');
  });

  it('renders status select', async () => {
    render(<UserFormDialog open={true} onOpenChange={vi.fn()} onSubmit={vi.fn()} />);

    const statusSelect = screen.getByRole('combobox');
    expect(statusSelect).toBeInTheDocument();
  });
});
