import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from '@/stores/authStore'

// Mock API client
vi.mock('@/api/client', () => ({
  default: {
    post: vi.fn(),
    get: vi.fn(),
    put: vi.fn(),
  },
}))

describe('authStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
  })

  it('初始状态为未登录', () => {
    const auth = useAuthStore()
    expect(auth.isLoggedIn).toBe(false)
    expect(auth.user).toBeNull()
  })

  it('login 成功后设置 token 和 user', async () => {
    const auth = useAuthStore()
    const api = (await import('@/api/client')).default
    vi.mocked(api.post).mockResolvedValue({ token: 'test-token', user: { id: 'u1', username: 'test', email: 'test@test.com', role: 'user' } })

    await auth.login('test', 'password')
    expect(auth.isLoggedIn).toBe(true)
    expect(auth.user?.username).toBe('test')
    expect(localStorage.getItem('auth_token')).toBe('test-token')
  })

  it('register 成功后设置 token 和 user', async () => {
    const auth = useAuthStore()
    const api = (await import('@/api/client')).default
    vi.mocked(api.post).mockResolvedValue({ token: 'register-token', user: { id: 'u2', username: 'new', email: 'new@test.com', role: 'user' } })

    await auth.register('new', 'new@test.com', 'password')
    expect(auth.isLoggedIn).toBe(true)
    expect(auth.user?.username).toBe('new')
    expect(localStorage.getItem('auth_token')).toBe('register-token')
  })

  it('logout 清除状态', () => {
    const auth = useAuthStore()
    auth.token = 'token'
    auth.user = { id: 'u1', username: 'test', email: 'test@test.com', role: 'user' }
    auth.logout()
    expect(auth.isLoggedIn).toBe(false)
    expect(auth.user).toBeNull()
    expect(localStorage.getItem('auth_token')).toBeNull()
  })

  it('isAdmin 正确判断管理员', () => {
    const auth = useAuthStore()
    auth.user = { id: 'u1', username: 'admin', email: 'admin@test.com', role: 'admin' }
    expect(auth.isAdmin).toBe(true)
    auth.user!.role = 'user'
    expect(auth.isAdmin).toBe(false)
  })
})
