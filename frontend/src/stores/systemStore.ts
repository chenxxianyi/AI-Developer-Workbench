/**
 * System Store
 * System health and upload limits
 */

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { SystemStatus, DashboardStats } from '@/types/system'
import { getHealth, getSystemStatus, getDashboardStats } from '@/api/system'

export const useSystemStore = defineStore('system', () => {
  // State
  const health = ref<boolean | null>(null)
  const status = ref<SystemStatus | null>(null)
  const dashboardStats = ref<DashboardStats | null>(null)

  const loading = ref(false)
  const error = ref<string | null>(null)

  const statsLoading = ref(false)
  const statsError = ref<string | null>(null)

  // Actions
  async function fetchHealth() {
    loading.value = true
    error.value = null

    try {
      const response = await getHealth()
      health.value = response.status === 'ok'
    } catch (err: any) {
      health.value = false
      error.value = err.message || '健康检查失败'
    } finally {
      loading.value = false
    }
  }

  async function fetchStatus() {
    loading.value = true
    error.value = null

    try {
      status.value = await getSystemStatus()
      health.value = status.value.healthy
    } catch (err: any) {
      error.value = err.message || '获取系统状态失败'
      health.value = false
    } finally {
      loading.value = false
    }
  }

  async function fetchDashboardStats() {
    statsLoading.value = true
    statsError.value = null

    try {
      dashboardStats.value = await getDashboardStats()
    } catch (err: any) {
      statsError.value = err.message || '获取统计数据失败'
    } finally {
      statsLoading.value = false
    }
  }

  // Getters
  const uploadLimits = computed(() => status.value?.upload_limits ?? {
    image_max_bytes: 10_000_000,
    zip_max_bytes: 50_000_000,
    zip_max_files: 1000,
    zip_max_total_bytes: 100_000_000,
  })

  const providerInfo = computed(() => {
    if (!status.value) return 'Loading...'
    return `${status.value.provider} (${status.value.text_model})`
  })

  const isMockMode = computed(() => status.value?.mock_mode ?? false)

  return {
    // State
    health,
    status,
    dashboardStats,
    loading,
    error,
    statsLoading,
    statsError,

    // Actions
    fetchHealth,
    fetchStatus,
    fetchDashboardStats,

    // Getters
    uploadLimits,
    providerInfo,
    isMockMode,
  }
})