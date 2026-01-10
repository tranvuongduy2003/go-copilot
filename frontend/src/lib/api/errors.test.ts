import { describe, expect, it } from 'vitest';
import { ApiError, getErrorMessage, isApiError } from './errors';

describe('ApiError', () => {
  describe('constructor', () => {
    it('creates an error with all properties', () => {
      const error = new ApiError('Test error message', 400, 'VALIDATION_ERROR', [
        { field: 'email', message: 'Invalid email' },
      ]);

      expect(error.message).toBe('Test error message');
      expect(error.statusCode).toBe(400);
      expect(error.code).toBe('VALIDATION_ERROR');
      expect(error.details).toEqual([{ field: 'email', message: 'Invalid email' }]);
      expect(error.name).toBe('ApiError');
    });

    it('creates an error without optional properties', () => {
      const error = new ApiError('Server error', 500, 'INTERNAL_ERROR');

      expect(error.message).toBe('Server error');
      expect(error.statusCode).toBe(500);
      expect(error.code).toBe('INTERNAL_ERROR');
      expect(error.details).toBeUndefined();
    });
  });

  describe('status code helpers', () => {
    it('identifies unauthorized errors (401)', () => {
      const error = new ApiError('Unauthorized', 401, 'UNAUTHORIZED');
      expect(error.isUnauthorized()).toBe(true);
      expect(error.isForbidden()).toBe(false);
      expect(error.isNotFound()).toBe(false);
    });

    it('identifies forbidden errors (403)', () => {
      const error = new ApiError('Forbidden', 403, 'FORBIDDEN');
      expect(error.isUnauthorized()).toBe(false);
      expect(error.isForbidden()).toBe(true);
      expect(error.isNotFound()).toBe(false);
    });

    it('identifies not found errors (404)', () => {
      const error = new ApiError('Not found', 404, 'NOT_FOUND');
      expect(error.isUnauthorized()).toBe(false);
      expect(error.isForbidden()).toBe(false);
      expect(error.isNotFound()).toBe(true);
    });

    it('identifies validation errors (400)', () => {
      const error = new ApiError('Validation error', 400, 'VALIDATION_ERROR');
      expect(error.isValidationError()).toBe(true);
    });

    it('identifies validation errors (422)', () => {
      const error = new ApiError('Unprocessable entity', 422, 'VALIDATION_ERROR');
      expect(error.isValidationError()).toBe(true);
    });

    it('identifies server errors (500+)', () => {
      const error500 = new ApiError('Internal error', 500, 'INTERNAL_ERROR');
      const error502 = new ApiError('Bad gateway', 502, 'BAD_GATEWAY');
      const error503 = new ApiError('Service unavailable', 503, 'SERVICE_UNAVAILABLE');

      expect(error500.isServerError()).toBe(true);
      expect(error502.isServerError()).toBe(true);
      expect(error503.isServerError()).toBe(true);
    });

    it('does not identify client errors as server errors', () => {
      const error = new ApiError('Bad request', 400, 'BAD_REQUEST');
      expect(error.isServerError()).toBe(false);
    });
  });

  describe('getFieldError', () => {
    it('returns error message for a specific field', () => {
      const error = new ApiError('Validation failed', 400, 'VALIDATION_ERROR', [
        { field: 'email', message: 'Invalid email format' },
        { field: 'password', message: 'Password too short' },
      ]);

      expect(error.getFieldError('email')).toBe('Invalid email format');
      expect(error.getFieldError('password')).toBe('Password too short');
    });

    it('returns undefined for non-existent field', () => {
      const error = new ApiError('Validation failed', 400, 'VALIDATION_ERROR', [
        { field: 'email', message: 'Invalid email format' },
      ]);

      expect(error.getFieldError('name')).toBeUndefined();
    });

    it('returns undefined when no details exist', () => {
      const error = new ApiError('Validation failed', 400, 'VALIDATION_ERROR');
      expect(error.getFieldError('email')).toBeUndefined();
    });
  });

  describe('toJSON', () => {
    it('serializes error to JSON object', () => {
      const error = new ApiError('Test error', 400, 'TEST_ERROR', [
        { field: 'test', message: 'Test message' },
      ]);

      const json = error.toJSON();

      expect(json).toEqual({
        name: 'ApiError',
        message: 'Test error',
        statusCode: 400,
        code: 'TEST_ERROR',
        details: [{ field: 'test', message: 'Test message' }],
      });
    });
  });
});

describe('isApiError', () => {
  it('returns true for ApiError instances', () => {
    const error = new ApiError('Test', 400, 'TEST');
    expect(isApiError(error)).toBe(true);
  });

  it('returns false for regular Error instances', () => {
    const error = new Error('Test');
    expect(isApiError(error)).toBe(false);
  });

  it('returns false for plain objects', () => {
    const error = { message: 'Test', statusCode: 400 };
    expect(isApiError(error)).toBe(false);
  });

  it('returns false for null and undefined', () => {
    expect(isApiError(null)).toBe(false);
    expect(isApiError(undefined)).toBe(false);
  });
});

describe('getErrorMessage', () => {
  it('extracts message from ApiError', () => {
    const error = new ApiError('API error message', 400, 'TEST');
    expect(getErrorMessage(error)).toBe('API error message');
  });

  it('extracts message from regular Error', () => {
    const error = new Error('Regular error message');
    expect(getErrorMessage(error)).toBe('Regular error message');
  });

  it('returns string errors as-is', () => {
    expect(getErrorMessage('String error')).toBe('String error');
  });

  it('returns default message for unknown error types', () => {
    expect(getErrorMessage(null)).toBe('An unexpected error occurred');
    expect(getErrorMessage(undefined)).toBe('An unexpected error occurred');
    expect(getErrorMessage(123)).toBe('An unexpected error occurred');
    expect(getErrorMessage({})).toBe('An unexpected error occurred');
  });
});
