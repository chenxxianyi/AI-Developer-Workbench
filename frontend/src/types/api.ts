/**
 * API Response Types
 * Unified response format from backend
 */

export interface ApiResponse<T> {
  code: 0
  message: string
  data: T
}

export interface ApiErrorResponse {
  code: number
  message: string
  error?: string
  request_id?: string
}

export interface PaginatedData<T> {
  items: T[]
  total: number
  page: number
  page_size: number
}