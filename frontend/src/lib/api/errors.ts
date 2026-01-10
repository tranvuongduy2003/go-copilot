import type { AxiosError } from 'axios';

export interface ApiErrorDetail {
  field: string;
  message: string;
}

export interface ApiErrorResponse {
  error: {
    code: string;
    message: string;
    details?: ApiErrorDetail[];
  };
}

export class ApiError extends Error {
  public readonly statusCode: number;
  public readonly code: string;
  public readonly details?: ApiErrorDetail[];
  public readonly originalError?: AxiosError;

  constructor(
    message: string,
    statusCode: number,
    code: string,
    details?: ApiErrorDetail[],
    originalError?: AxiosError
  ) {
    super(message);
    this.name = 'ApiError';
    this.statusCode = statusCode;
    this.code = code;
    this.details = details;
    this.originalError = originalError;

    Object.setPrototypeOf(this, ApiError.prototype);
  }

  public isUnauthorized(): boolean {
    return this.statusCode === 401;
  }

  public isForbidden(): boolean {
    return this.statusCode === 403;
  }

  public isNotFound(): boolean {
    return this.statusCode === 404;
  }

  public isValidationError(): boolean {
    return this.statusCode === 400 || this.statusCode === 422;
  }

  public isServerError(): boolean {
    return this.statusCode >= 500;
  }

  public getFieldError(field: string): string | undefined {
    return this.details?.find((detail) => detail.field === field)?.message;
  }

  public toJSON(): Record<string, unknown> {
    return {
      name: this.name,
      message: this.message,
      statusCode: this.statusCode,
      code: this.code,
      details: this.details,
    };
  }
}

export function isApiError(error: unknown): error is ApiError {
  return error instanceof ApiError;
}

export function getErrorMessage(error: unknown): string {
  if (isApiError(error)) {
    return error.message;
  }

  if (error instanceof Error) {
    return error.message;
  }

  if (typeof error === 'string') {
    return error;
  }

  return 'An unexpected error occurred';
}
