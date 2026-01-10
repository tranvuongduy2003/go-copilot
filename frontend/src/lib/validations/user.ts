import { z } from 'zod';
import { emailSchema, paginationSchema } from './auth';

export const userStatusSchema = z.enum(['pending', 'active', 'inactive', 'banned']);

export const createUserSchema = z.object({
  email: emailSchema,
  fullName: z
    .string()
    .min(2, 'Full name must be at least 2 characters')
    .max(100, 'Full name must be at most 100 characters'),
  password: z.string().min(8, 'Password must be at least 8 characters').optional(),
  roleIds: z.array(z.string().uuid()).optional(),
});

export const updateUserSchema = z.object({
  email: emailSchema.optional(),
  fullName: z
    .string()
    .min(2, 'Full name must be at least 2 characters')
    .max(100, 'Full name must be at most 100 characters')
    .optional(),
  status: userStatusSchema.optional(),
});

export const userFilterSchema = paginationSchema.extend({
  search: z.string().optional(),
  status: userStatusSchema.optional(),
  roleId: z.string().uuid().optional(),
  sortBy: z.enum(['createdAt', 'updatedAt', 'email', 'fullName']).optional(),
  sortOrder: z.enum(['asc', 'desc']).optional(),
});

export type UserStatus = z.infer<typeof userStatusSchema>;
export type CreateUserInput = z.infer<typeof createUserSchema>;
export type UpdateUserInput = z.infer<typeof updateUserSchema>;
export type UserFilterInput = z.infer<typeof userFilterSchema>;
