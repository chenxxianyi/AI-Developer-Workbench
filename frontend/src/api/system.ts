/**
 * System API
 * Health check, system status, and dashboard stats
 */

import apiClient from './client'
import type { SystemStatus, DashboardStats, HealthStatus } from '@/types/system'

/**
 * Check backend health status
 */
export async function getHealth(): Promise<HealthStatus> {
  return apiClient.get('/health')
}

/**
 * Get system status (Mock Mode, provider, upload limits)
 */
export async function getSystemStatus(): Promise<SystemStatus> {
  return apiClient.get('/system/status')
}

/**
 * Get dashboard statistics and recent reports
 */
export async function getDashboardStats(): Promise<DashboardStats> {
  return apiClient.get('/dashboard/stats')
}