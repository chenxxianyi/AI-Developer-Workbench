/**
 * Axios API Client
 * Unified error handling and response parsing
 */

import axios, { AxiosError, type AxiosInstance } from 'axios'
import type { ApiResponse, ApiErrorResponse } from '@/types/api'

// Create Axios instance
const apiClient: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api',
  timeout: 100_000, // Slightly above backend AI timeout
  headers: {
    'Content-Type': 'application/json',
  },
})

// Response interceptor: unwrap ApiResponse.data or throw normalized error
apiClient.interceptors.response.use(
  (response) => {
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
  (error: AxiosError<ApiErrorResponse>) => {
    // Handle backend error response
    if (error.response?.data) {
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