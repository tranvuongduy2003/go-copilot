import { describe, expect, it } from 'vitest';
import { createRoleSchema, updateRoleSchema } from './role';

describe('Role Validations', () => {
  describe('createRoleSchema', () => {
    it('validates correct role data', () => {
      const result = createRoleSchema.safeParse({
        name: 'admin',
        displayName: 'Administrator',
      });
      expect(result.success).toBe(true);
    });

    it('validates role data with all fields', () => {
      const result = createRoleSchema.safeParse({
        name: 'content_moderator',
        displayName: 'Content Moderator',
        description: 'Can moderate user content',
        permissionIds: ['550e8400-e29b-41d4-a716-446655440000'],
      });
      expect(result.success).toBe(true);
    });

    describe('name validation', () => {
      it('rejects name shorter than 2 characters', () => {
        const result = createRoleSchema.safeParse({
          name: 'a',
          displayName: 'Single Char',
        });
        expect(result.success).toBe(false);
        if (!result.success) {
          expect(result.error.issues[0].message).toContain('at least 2 characters');
        }
      });

      it('rejects name longer than 50 characters', () => {
        const result = createRoleSchema.safeParse({
          name: 'a'.repeat(51),
          displayName: 'Long Name',
        });
        expect(result.success).toBe(false);
        if (!result.success) {
          expect(result.error.issues[0].message).toContain('at most 50 characters');
        }
      });

      it('rejects name starting with number', () => {
        const result = createRoleSchema.safeParse({
          name: '1admin',
          displayName: 'Admin',
        });
        expect(result.success).toBe(false);
        if (!result.success) {
          expect(result.error.issues[0].message).toContain('must start with a letter');
        }
      });

      it('rejects name with uppercase letters', () => {
        const result = createRoleSchema.safeParse({
          name: 'Admin',
          displayName: 'Admin',
        });
        expect(result.success).toBe(false);
      });

      it('rejects name with special characters', () => {
        const result = createRoleSchema.safeParse({
          name: 'admin-role',
          displayName: 'Admin',
        });
        expect(result.success).toBe(false);
      });

      it('allows underscores in name', () => {
        const result = createRoleSchema.safeParse({
          name: 'super_admin',
          displayName: 'Super Admin',
        });
        expect(result.success).toBe(true);
      });

      it('allows numbers in name (not at start)', () => {
        const result = createRoleSchema.safeParse({
          name: 'admin123',
          displayName: 'Admin 123',
        });
        expect(result.success).toBe(true);
      });
    });

    describe('displayName validation', () => {
      it('rejects displayName shorter than 2 characters', () => {
        const result = createRoleSchema.safeParse({
          name: 'admin',
          displayName: 'A',
        });
        expect(result.success).toBe(false);
        if (!result.success) {
          expect(result.error.issues[0].message).toContain('at least 2 characters');
        }
      });

      it('rejects displayName longer than 100 characters', () => {
        const result = createRoleSchema.safeParse({
          name: 'admin',
          displayName: 'A'.repeat(101),
        });
        expect(result.success).toBe(false);
        if (!result.success) {
          expect(result.error.issues[0].message).toContain('at most 100 characters');
        }
      });
    });

    describe('description validation', () => {
      it('allows empty description', () => {
        const result = createRoleSchema.safeParse({
          name: 'admin',
          displayName: 'Administrator',
        });
        expect(result.success).toBe(true);
      });

      it('rejects description longer than 500 characters', () => {
        const result = createRoleSchema.safeParse({
          name: 'admin',
          displayName: 'Administrator',
          description: 'A'.repeat(501),
        });
        expect(result.success).toBe(false);
        if (!result.success) {
          expect(result.error.issues[0].message).toContain('at most 500 characters');
        }
      });
    });

    describe('permissionIds validation', () => {
      it('allows empty permissionIds array', () => {
        const result = createRoleSchema.safeParse({
          name: 'admin',
          displayName: 'Administrator',
          permissionIds: [],
        });
        expect(result.success).toBe(true);
      });

      it('validates valid UUIDs', () => {
        const result = createRoleSchema.safeParse({
          name: 'admin',
          displayName: 'Administrator',
          permissionIds: [
            '550e8400-e29b-41d4-a716-446655440000',
            '6ba7b810-9dad-11d1-80b4-00c04fd430c8',
          ],
        });
        expect(result.success).toBe(true);
      });

      it('rejects invalid UUIDs', () => {
        const result = createRoleSchema.safeParse({
          name: 'admin',
          displayName: 'Administrator',
          permissionIds: ['not-a-uuid'],
        });
        expect(result.success).toBe(false);
      });
    });
  });

  describe('updateRoleSchema', () => {
    it('validates empty update (all optional)', () => {
      const result = updateRoleSchema.safeParse({});
      expect(result.success).toBe(true);
    });

    it('validates with only displayName', () => {
      const result = updateRoleSchema.safeParse({
        displayName: 'Updated Name',
      });
      expect(result.success).toBe(true);
    });

    it('validates with only description', () => {
      const result = updateRoleSchema.safeParse({
        description: 'Updated description',
      });
      expect(result.success).toBe(true);
    });

    it('validates with only permissionIds', () => {
      const result = updateRoleSchema.safeParse({
        permissionIds: ['550e8400-e29b-41d4-a716-446655440000'],
      });
      expect(result.success).toBe(true);
    });

    it('validates with all fields', () => {
      const result = updateRoleSchema.safeParse({
        displayName: 'Updated Name',
        description: 'Updated description',
        permissionIds: ['550e8400-e29b-41d4-a716-446655440000'],
      });
      expect(result.success).toBe(true);
    });

    it('does not allow name update (not in schema)', () => {
      const result = updateRoleSchema.safeParse({
        name: 'new_name',
        displayName: 'New Name',
      });
      expect(result.success).toBe(true);
      if (result.success) {
        expect(result.data).not.toHaveProperty('name');
      }
    });

    it('rejects invalid displayName when provided', () => {
      const result = updateRoleSchema.safeParse({
        displayName: 'A',
      });
      expect(result.success).toBe(false);
    });

    it('rejects invalid description when provided', () => {
      const result = updateRoleSchema.safeParse({
        description: 'A'.repeat(501),
      });
      expect(result.success).toBe(false);
    });

    it('rejects invalid permissionIds when provided', () => {
      const result = updateRoleSchema.safeParse({
        permissionIds: ['invalid'],
      });
      expect(result.success).toBe(false);
    });
  });
});
