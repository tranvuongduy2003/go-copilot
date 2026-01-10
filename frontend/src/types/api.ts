export interface ApiResponse<TData> {
  data: TData;
}

export interface ApiErrorResponse {
  error: {
    code: string;
    message: string;
    details?: ApiErrorDetail[];
  };
}

export interface ApiErrorDetail {
  field: string;
  message: string;
}

export interface PaginatedResponse<TData> {
  data: TData[];
  meta: PaginationMeta;
}

export interface PaginationMeta {
  page: number;
  limit: number;
  pageSize: number;
  total: number;
  totalPages: number;
  hasNext: boolean;
  hasPrevious: boolean;
}

export interface PaginationParams {
  page?: number;
  limit?: number;
}

export interface SortParams {
  sortBy?: string;
  sortOrder?: 'asc' | 'desc';
}

export type ListParams<TFilter = Record<string, unknown>> = PaginationParams & SortParams & TFilter;
