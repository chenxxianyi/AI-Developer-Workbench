/**
 * Axios API Client
 * Unified error handling and response parsing
 */

import axios, { AxiosError, type AxiosInstance } from 'axios'
import type { ApiResponse, ApiErrorResponse } from '@/types/api'

// Create Axios instance
const apiClient: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1',
  timeout: 100_000,
  headers: { 'Content-Type': 'application/json' },
})

// Request interceptor: attach JWT token
apiClient.interceptors.request.use((config) => {
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

    const body = response.data as ApiResponse<unknown>

    // Check if business code is not 0 (error)
    if (body.code !== 0) {
      const error: ApiErrorResponse = {
        code: body.code,
        message: body.message,
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
      return Promise.reject(error.response.data)
    }

    // Handle timeout
    if (error.code === 'ECONNABORTED') {
      return Promise.reject({
        code: -1,
        message: '请求超时，请稍后重试',
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
