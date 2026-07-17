/**
 * Axios API Client
 * Unified error handling and response parsing
 */

import axios, { AxiosError, type AxiosInstance } from 'axios'
import type { ApiErrorResponse } from '@/types/api'

// Create Axios instance
const apiClient: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api',
  timeout: 100_000,
  headers: { 'Content-Type': 'application/json' },
})

// Request interceptor: attach JWT token (skip for auth endpoints)
apiClient.interceptors.request.use((config) => {
  if (config.url && (config.url.endsWith('/auth/login') || config.url.endsWith('/auth/register'))) {
    return config
  }
  const token = localStorage.getItem('auth_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

// Response interceptor: unwrap ApiResponse.data or throw normalized error
apiClient.interceptors.response.use(
  (response) => {
    if (response.config.responseType === 'blob' || response.config.responseType === 'arraybuffer') {
      return response as any
    }

    const body = response.data as Record<string, unknown>

    // Backend error format (two variants):
    //   old: { error: { code: "UNAUTHORIZED", message: "..." } }
    //   new: { code: 50001, message: "internal error", error: "detail string" }
    if (body.error) {
      let message = '请求失败'
      if (typeof body.error === 'string') {
        message = body.error
      } else if (typeof body.error === 'object') {
        message = (body.error as { message?: string }).message || message
      }
      const error: ApiErrorResponse = {
        code: typeof body.code === 'number' ? body.code : 0,
        message,
      }
      return Promise.reject(error)
    }

    // Return just the data portion so callers get T directly
    return body.data as any
  },
  async (error: AxiosError<ApiErrorResponse | Blob>) => {
    // Handle backend error response
    if (error.response?.data) {
      const contentType = error.response.headers['content-type']
      if (error.response.data instanceof Blob && typeof contentType === 'string' && contentType.includes('application/json')) {
        try {
          const parsed = JSON.parse(await error.response.data.text()) as ApiErrorResponse
          return Promise.reject(parsed)
        } catch {
          return Promise.reject({
            code: error.response.status,
            message: `服务器错误 (${error.response.status})`,
          })
        }
      }
      return Promise.reject(normalizeError(error.response.data))
    }

    // Handle timeout
    if (error.code === 'ECONNABORTED') {
      const seconds = typeof error.config?.timeout === 'number' && error.config.timeout > 0
        ? Math.round(error.config.timeout / 1000)
        : 0
      return Promise.reject({
        code: -1,
        message: seconds ? `请求超过 ${seconds} 秒仍未完成，请稍后重试` : '请求超时，请稍后重试',
      })
    }

    // Handle network error
    if (!error.response) {
      return Promise.reject({
        code: -1,
        message: '网络连接失败，请检查网络',
      })
    }

    // Handle other HTTP errors
    return Promise.reject({
      code: error.response.status,
      message: `服务器错误 (${error.response.status})`,
    })
  }
)

export default apiClient

/** Normalize backend error responses to ApiErrorResponse format. */
function normalizeError(data: unknown): ApiErrorResponse {
  const obj = data as Record<string, unknown>
  // Old format: { error: { code: "UNAUTHORIZED", message: "..." } }
  if (obj.error && typeof obj.error === 'object') {
    const err = obj.error as { code?: string; message?: string }
    return { code: 0, message: err.message || '请求失败' }
  }
  // New format: { code: 50001, message: "internal error", error: "detail" }
  if (typeof obj.error === 'string') {
    return { code: typeof obj.code === 'number' ? obj.code : 0, message: obj.error }
  }
  return { code: typeof obj.code === 'number' ? obj.code : 0, message: (obj.message as string) || '请求失败' }
}

