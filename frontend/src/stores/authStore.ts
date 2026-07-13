import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import apiClient from '@/api/client'

export interface User {
  id: string
  username: string
  email: string
  role: 'user' | 'admin'
}

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const token = ref<string | null>(localStorage.getItem('auth_token'))

  const isLoggedIn = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.role === 'admin')

  async function login(username: string, password: string) {
    const res = await apiClient.post('/auth/login', { username, password }) as any
    token.value = res.token
    user.value = res.user
    localStorage.setItem('auth_token', res.token)
  }

  async function register(username: string, email: string, password: string) {
    await apiClient.post('/auth/register', { username, email, password })
  }

  async function fetchProfile() {
    if (!token.value) return
    try {
      user.value = await apiClient.get('/auth/profile') as any
    } catch {
      logout()
    }
  }

  async function updateProfile(data: { username?: string; email?: string }) {
    user.value = await apiClient.put('/auth/profile', data) as any
  }

  function logout() {
    token.value = null
    user.value = null
    localStorage.removeItem('auth_token')
  }

  // Auto-fetch profile on app start if token exists
  if (token.value) {
    fetchProfile()
  }

  return { user, token, isLoggedIn, isAdmin, login, register, fetchProfile, updateProfile, logout }
})
