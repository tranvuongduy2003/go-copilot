import { render, screen } from '@/test/utils';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import { Input } from './input';

describe('Input', () => {
  it('renders with default type', () => {
    render(<Input placeholder="Enter text" />);
    const input = screen.getByPlaceholderText('Enter text');
    expect(input).toBeInTheDocument();
    // Browser defaults to text type when not explicitly set
    expect(input.getAttribute('type') ?? 'text').toBe('text');
  });

  it('renders with email type', () => {
    render(<Input type="email" placeholder="Enter email" />);
    const input = screen.getByPlaceholderText('Enter email');
    expect(input).toHaveAttribute('type', 'email');
  });

  it('renders with password type', () => {
    render(<Input type="password" placeholder="Enter password" />);
    const input = screen.getByPlaceholderText('Enter password');
    expect(input).toHaveAttribute('type', 'password');
  });

  it('handles value changes', async () => {
    const user = userEvent.setup();
    const handleChange = vi.fn();
    render(<Input placeholder="Enter text" onChange={handleChange} />);

    const input = screen.getByPlaceholderText('Enter text');
    await user.type(input, 'Hello World');

    expect(handleChange).toHaveBeenCalled();
    expect(input).toHaveValue('Hello World');
  });

  it('renders as disabled', () => {
    render(<Input disabled placeholder="Disabled input" />);
    const input = screen.getByPlaceholderText('Disabled input');
    expect(input).toBeDisabled();
  });

  it('renders as required', () => {
    render(<Input required placeholder="Required input" />);
    const input = screen.getByPlaceholderText('Required input');
    expect(input).toBeRequired();
  });

  it('applies custom className', () => {
    render(<Input className="custom-class" placeholder="Custom input" />);
    const input = screen.getByPlaceholderText('Custom input');
    expect(input).toHaveClass('custom-class');
  });

  it('renders with aria-invalid for validation errors', () => {
    render(<Input aria-invalid="true" placeholder="Invalid input" />);
    const input = screen.getByPlaceholderText('Invalid input');
    expect(input).toHaveAttribute('aria-invalid', 'true');
  });

  it('supports controlled value', () => {
    const { rerender } = render(<Input value="initial" placeholder="Controlled" readOnly />);
    const input = screen.getByPlaceholderText('Controlled');
    expect(input).toHaveValue('initial');

    rerender(<Input value="updated" placeholder="Controlled" readOnly />);
    expect(input).toHaveValue('updated');
  });

  it('supports file input type', () => {
    render(<Input type="file" data-testid="file-input" />);
    const input = screen.getByTestId('file-input');
    expect(input).toHaveAttribute('type', 'file');
  });

  it('forwards ref correctly', () => {
    const ref = vi.fn();
    render(<Input ref={ref} placeholder="Ref input" />);
    expect(ref).toHaveBeenCalled();
  });
});
