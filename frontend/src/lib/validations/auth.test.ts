import { describe, expect, it } from 'vitest';
import { forgotPasswordSchema, loginSchema, registerSchema, resetPasswordSchema } from './auth';

describe('Auth Validations', () => {
  describe('loginSchema', () => {
    it('validates correct login data', () => {
      const result = loginSchema.safeParse({
        email: 'test@example.com',
        password: 'password123',
        rememberMe: false,
      });
      expect(result.success).toBe(true);
    });

    it('rejects invalid email', () => {
      const result = loginSchema.safeParse({
        email: 'invalid-email',
        password: 'password123',
        rememberMe: false,
      });
      expect(result.success).toBe(false);
    });

    it('rejects empty email', () => {
      const result = loginSchema.safeParse({
        email: '',
        password: 'password123',
        rememberMe: false,
      });
      expect(result.success).toBe(false);
    });

    it('rejects empty password', () => {
      const result = loginSchema.safeParse({
        email: 'test@example.com',
        password: '',
        rememberMe: false,
      });
      expect(result.success).toBe(false);
    });
  });

  describe('registerSchema', () => {
    it('validates correct registration data', () => {
      const result = registerSchema.safeParse({
        fullName: 'Test User',
        email: 'test@example.com',
        password: 'StrongPass1!',
        confirmPassword: 'StrongPass1!',
        acceptTerms: true,
      });
      expect(result.success).toBe(true);
    });

    it('rejects mismatched passwords', () => {
      const result = registerSchema.safeParse({
        fullName: 'Test User',
        email: 'test@example.com',
        password: 'StrongPass1!',
        confirmPassword: 'DifferentPass1!',
        acceptTerms: true,
      });
      expect(result.success).toBe(false);
    });

    it('rejects short password', () => {
      const result = registerSchema.safeParse({
        fullName: 'Test User',
        email: 'test@example.com',
        password: '123',
        confirmPassword: '123',
        acceptTerms: true,
      });
      expect(result.success).toBe(false);
    });

    it('rejects empty name', () => {
      const result = registerSchema.safeParse({
        fullName: '',
        email: 'test@example.com',
        password: 'StrongPass1!',
        confirmPassword: 'StrongPass1!',
        acceptTerms: true,
      });
      expect(result.success).toBe(false);
    });
  });

  describe('forgotPasswordSchema', () => {
    it('validates correct email', () => {
      const result = forgotPasswordSchema.safeParse({
        email: 'test@example.com',
      });
      expect(result.success).toBe(true);
    });

    it('rejects invalid email', () => {
      const result = forgotPasswordSchema.safeParse({
        email: 'not-an-email',
      });
      expect(result.success).toBe(false);
    });
  });

  describe('resetPasswordSchema', () => {
    it('validates correct reset password data', () => {
      const result = resetPasswordSchema.safeParse({
        token: 'valid-token-123',
        password: 'NewStrongPass1!',
        confirmPassword: 'NewStrongPass1!',
      });
      expect(result.success).toBe(true);
    });

    it('rejects mismatched passwords', () => {
      const result = resetPasswordSchema.safeParse({
        token: 'valid-token-123',
        password: 'NewStrongPass1!',
        confirmPassword: 'DifferentPass1!',
      });
      expect(result.success).toBe(false);
    });

    it('rejects empty token', () => {
      const result = resetPasswordSchema.safeParse({
        token: '',
        password: 'NewStrongPass1!',
        confirmPassword: 'NewStrongPass1!',
      });
      expect(result.success).toBe(false);
    });
  });
});
