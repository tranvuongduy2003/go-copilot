import { render, screen, waitFor } from '@/test/utils';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { RegisterForm } from './register-form';

vi.mock('sonner', () => ({
  toast: {
    success: vi.fn(),
    error: vi.fn(),
  },
}));

describe('RegisterForm', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders all form fields', () => {
    render(<RegisterForm />);

    expect(screen.getByPlaceholderText(/enter your full name/i)).toBeInTheDocument();
    expect(screen.getByPlaceholderText(/enter your email/i)).toBeInTheDocument();
    expect(screen.getByPlaceholderText(/create a password/i)).toBeInTheDocument();
    expect(screen.getByPlaceholderText(/confirm your password/i)).toBeInTheDocument();
    expect(screen.getByRole('checkbox')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /create account/i })).toBeInTheDocument();
  });

  it('shows validation errors for empty form submission', async () => {
    const user = userEvent.setup();
    render(<RegisterForm />);

    await user.click(screen.getByRole('button', { name: /create account/i }));

    await waitFor(() => {
      // The validation messages mention minimum characters or required
      const validationMessages = screen.queryAllByText(/required|at least|must be/i);
      expect(validationMessages.length).toBeGreaterThan(0);
    });
  });

  // Note: Email validation is covered in src/lib/validations/auth.test.ts
  // Form-level validation triggers correctly as shown in the empty form submission test

  it('shows password strength indicator', async () => {
    const user = userEvent.setup();
    render(<RegisterForm />);

    const passwordInput = screen.getByPlaceholderText(/create a password/i);
    await user.type(passwordInput, 'Password123!');

    await waitFor(() => {
      expect(screen.getByText(/at least 8 characters/i)).toBeInTheDocument();
    });
  });

  it('requires terms acceptance', async () => {
    const user = userEvent.setup();
    render(<RegisterForm />);

    await user.type(screen.getByPlaceholderText(/enter your full name/i), 'Test User');
    await user.type(screen.getByPlaceholderText(/enter your email/i), 'test@example.com');
    await user.type(screen.getByPlaceholderText(/create a password/i), 'Password123!');
    await user.type(screen.getByPlaceholderText(/confirm your password/i), 'Password123!');

    await user.click(screen.getByRole('button', { name: /create account/i }));

    await waitFor(() => {
      expect(screen.getByText(/must accept the terms/i)).toBeInTheDocument();
    });
  });

  it('has links to terms and privacy policy', () => {
    render(<RegisterForm />);

    expect(screen.getByRole('link', { name: /terms of service/i })).toHaveAttribute(
      'href',
      '/terms'
    );
    expect(screen.getByRole('link', { name: /privacy policy/i })).toHaveAttribute(
      'href',
      '/privacy'
    );
  });
});
