/**
 * Project Types
 */

import type { ToolType } from './tool'
import type { ReportStatus } from './report'
import type { QualityTrendPoint } from './system'

export interface Project {
  id: string
  name: string
  project_type?: ProjectType
  description: string
  repo_url: string
  frontend_stack: string
  backend_stack: string
  database: string
  ui_style: string
  coding_rules: string
  created_at: string
  updated_at: string
}

export interface ProjectSummary {
  id: string
  name: string
  project_type?: ProjectType
  description: string
  repo_url: string
  report_count: number
  average_score: number | null
  created_at: string
  updated_at: string
}

export interface ProjectStats {
  total_reports: number
  average_score: number | null
  tool_usage: Partial<Record<ToolType, number>>
  high_severity_count: number
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
  quality_trend: QualityTrendPoint[]
  latest_artifacts: ProjectArtifact[]
}

export interface ProjectArtifact {
  tool_type: ToolType
  report_id: string
  report_title: string
  filename: string
  mime_type: string
  created_at: string
}

export interface ProjectDeleteResult {
  detached_report_count: number
}

export interface ProjectCreateInput {
  name: string
  project_type?: ProjectType
  description?: string
  repo_url?: string
  frontend_stack?: string
  backend_stack?: string
  database?: string
  ui_style?: string
  coding_rules?: string
}

export interface ProjectUpdateInput {
  name?: string
  project_type?: ProjectType
  description?: string
  repo_url?: string
  frontend_stack?: string
  backend_stack?: string
  database?: string
  ui_style?: string
  coding_rules?: string
}

export type ProjectType =
  | 'interactive_app'
  | 'dashboard'
  | 'data_product'
  | 'content_site'
  | 'ecommerce'
  | 'utility_app'
  | 'landing_page'
  | 'analysis_existing'

export interface ProjectListParams {
  search?: string
  page?: number
  page_size?: number
}
