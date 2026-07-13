/**
 * System Types
 * System status and dashboard statistics
 */

import type { ToolType } from './tool'
import type { ReportStatus } from './report'

export interface SystemStatus {
  healthy: boolean
  provider: string
  text_model: string
  vision_model: string
  mock_mode: boolean
  upload_limits: {
    image_max_bytes: number
    zip_max_bytes: number
    zip_max_files: number
    zip_max_total_bytes: number
  }
}

export interface QualityTrendPoint {
  /** ISO date (YYYY-MM-DD) the bucket starts on (UTC). */
  bucket: string
  report_count: number
  /** null when no scored reports in this bucket. */
  average_score: number | null
  high_severity_count: number
}

export interface WeeklyStats {
  report_count_this_week: number
  /** null when no scored reports this week. */
  average_score_this_week: number | null
  high_severity_count_this_week: number
  most_used_tool_this_week: string
}

export interface DashboardStats {
  total_reports: number
  tool_usage: Record<ToolType, number>
  average_score: number | null // Can be null if no scored reports
  recent_reports: Array<{
    id: string
    tool_type: ToolType
    title: string
    status: ReportStatus
    total_score: number | null
    grade: string | null
    summary: string
    created_at: string
  }>
  weekly_stats?: WeeklyStats
  quality_trend?: QualityTrendPoint[]
}

export interface HealthStatus {
  status: string
  timestamp?: string
}