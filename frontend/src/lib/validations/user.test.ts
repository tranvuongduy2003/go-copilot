import { describe, expect, it } from 'vitest';
import { createUserSchema, updateUserSchema, userFilterSchema, userStatusSchema } from './user';

describe('User Validations', () => {
  describe('userStatusSchema', () => {
    it('validates valid status values', () => {
      expect(userStatusSchema.safeParse('pending').success).toBe(true);
      expect(userStatusSchema.safeParse('active').success).toBe(true);
      expect(userStatusSchema.safeParse('inactive').success).toBe(true);
      expect(userStatusSchema.safeParse('banned').success).toBe(true);
    });

    it('rejects invalid status values', () => {
      expect(userStatusSchema.safeParse('unknown').success).toBe(false);
      expect(userStatusSchema.safeParse('').success).toBe(false);
      expect(userStatusSchema.safeParse(null).success).toBe(false);
    });
  });

  describe('createUserSchema', () => {
    it('validates correct user data', () => {
      const result = createUserSchema.safeParse({
        email: 'user@example.com',
        fullName: 'John Doe',
      });
      expect(result.success).toBe(true);
    });

    it('validates user data with optional fields', () => {
      const result = createUserSchema.safeParse({
        email: 'user@example.com',
        fullName: 'John Doe',
        password: 'SecurePass123!',
        roleIds: ['550e8400-e29b-41d4-a716-446655440000'],
      });
      expect(result.success).toBe(true);
    });

    it('rejects invalid email', () => {
      const result = createUserSchema.safeParse({
        email: 'invalid-email',
        fullName: 'John Doe',
      });
      expect(result.success).toBe(false);
    });

    it('rejects short full name', () => {
      const result = createUserSchema.safeParse({
        email: 'user@example.com',
        fullName: 'J',
      });
      expect(result.success).toBe(false);
      if (!result.success) {
        expect(result.error.issues[0].message).toContain('at least 2 characters');
      }
    });

    it('rejects long full name', () => {
      const result = createUserSchema.safeParse({
        email: 'user@example.com',
        fullName: 'A'.repeat(101),
      });
      expect(result.success).toBe(false);
      if (!result.success) {
        expect(result.error.issues[0].message).toContain('at most 100 characters');
      }
    });

    it('rejects short password when provided', () => {
      const result = createUserSchema.safeParse({
        email: 'user@example.com',
        fullName: 'John Doe',
        password: 'short',
      });
      expect(result.success).toBe(false);
    });

    it('rejects invalid UUID in roleIds', () => {
      const result = createUserSchema.safeParse({
        email: 'user@example.com',
        fullName: 'John Doe',
        roleIds: ['not-a-uuid'],
      });
      expect(result.success).toBe(false);
    });

    it('accepts empty roleIds array', () => {
      const result = createUserSchema.safeParse({
        email: 'user@example.com',
        fullName: 'John Doe',
        roleIds: [],
      });
      expect(result.success).toBe(true);
    });
  });

  describe('updateUserSchema', () => {
    it('validates with no fields (all optional)', () => {
      const result = updateUserSchema.safeParse({});
      expect(result.success).toBe(true);
    });

    it('validates with only email', () => {
      const result = updateUserSchema.safeParse({
        email: 'updated@example.com',
      });
      expect(result.success).toBe(true);
    });

    it('validates with only fullName', () => {
      const result = updateUserSchema.safeParse({
        fullName: 'Updated Name',
      });
      expect(result.success).toBe(true);
    });

    it('validates with only status', () => {
      const result = updateUserSchema.safeParse({
        status: 'active',
      });
      expect(result.success).toBe(true);
    });

    it('validates with all fields', () => {
      const result = updateUserSchema.safeParse({
        email: 'updated@example.com',
        fullName: 'Updated Name',
        status: 'inactive',
      });
      expect(result.success).toBe(true);
    });

    it('rejects invalid email when provided', () => {
      const result = updateUserSchema.safeParse({
        email: 'invalid-email',
      });
      expect(result.success).toBe(false);
    });

    it('rejects invalid status when provided', () => {
      const result = updateUserSchema.safeParse({
        status: 'unknown',
      });
      expect(result.success).toBe(false);
    });
  });

  describe('userFilterSchema', () => {
    it('validates empty filters (all optional)', () => {
      const result = userFilterSchema.safeParse({});
      expect(result.success).toBe(true);
    });

    it('validates pagination parameters', () => {
      const result = userFilterSchema.safeParse({
        page: 1,
        limit: 20,
      });
      expect(result.success).toBe(true);
    });

    it('validates search parameter', () => {
      const result = userFilterSchema.safeParse({
        search: 'john',
      });
      expect(result.success).toBe(true);
    });

    it('validates status filter', () => {
      const result = userFilterSchema.safeParse({
        status: 'active',
      });
      expect(result.success).toBe(true);
    });

    it('validates roleId filter', () => {
      const result = userFilterSchema.safeParse({
        roleId: '550e8400-e29b-41d4-a716-446655440000',
      });
      expect(result.success).toBe(true);
    });

    it('validates sortBy options', () => {
      const validSortFields = ['createdAt', 'updatedAt', 'email', 'fullName'];
      for (const sortBy of validSortFields) {
        const result = userFilterSchema.safeParse({ sortBy });
        expect(result.success).toBe(true);
      }
    });

    it('validates sortOrder options', () => {
      expect(userFilterSchema.safeParse({ sortOrder: 'asc' }).success).toBe(true);
      expect(userFilterSchema.safeParse({ sortOrder: 'desc' }).success).toBe(true);
    });

    it('validates all filters combined', () => {
      const result = userFilterSchema.safeParse({
        page: 2,
        limit: 50,
        search: 'admin',
        status: 'active',
        roleId: '550e8400-e29b-41d4-a716-446655440000',
        sortBy: 'email',
        sortOrder: 'asc',
      });
      expect(result.success).toBe(true);
    });

    it('rejects invalid status in filter', () => {
      const result = userFilterSchema.safeParse({
        status: 'invalid',
      });
      expect(result.success).toBe(false);
    });

    it('rejects invalid sortBy', () => {
      const result = userFilterSchema.safeParse({
        sortBy: 'invalidField',
      });
      expect(result.success).toBe(false);
    });

    it('rejects invalid sortOrder', () => {
      const result = userFilterSchema.safeParse({
        sortOrder: 'random',
      });
      expect(result.success).toBe(false);
    });

    it('rejects invalid roleId format', () => {
      const result = userFilterSchema.safeParse({
        roleId: 'not-a-uuid',
      });
      expect(result.success).toBe(false);
    });
  });
});
