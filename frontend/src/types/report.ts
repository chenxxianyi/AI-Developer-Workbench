/**
 * Report Types
 * Core report structure with nullable scores for Agent Config and API Doc
 */

import type { ToolType } from './tool'

export type ReportStatus = 'processing' | 'succeeded' | 'fallback' | 'failed'

export interface GeneratedFileMeta {
  id: string
  filename: string
  language?: string
  mime_type: string
  size_bytes: number
}

export interface ReportSummary {
  one_sentence: string
  strengths: string[]
  weaknesses: string[]
  top_priorities: string[]
}

export interface Issue {
  title: string
  severity: 'high' | 'medium' | 'low'
  category: string
  problem: string
  suggestion: string
  action: string
  viewport?: 'desktop' | 'mobile'
  region?: { x: number; y: number; width: number; height: number }
  contrast_suggestion?: string
  component_prompt?: string
}

export type ActionPriority = 'high' | 'medium' | 'low'
export type ActionEffort = 'small' | 'medium' | 'large'

export interface ActionItem {
  id: string
  title: string
  priority: ActionPriority
  effort: ActionEffort
  category: string
  reason: string
  suggested_prompt: string
  issue_title: string
  issue_body: string
}

/**
 * Generic Report type
 * T is the tool-specific report_data type
 *
 * IMPORTANT: total_score and grade can be NULL for Agent Config and API Doc
 * Frontend must NOT fabricate scores when null
 */
export interface Report<T = Record<string, unknown>> {
  id: string
  tool_type: ToolType
  title: string
  input_mode: string
  status: ReportStatus
  summary: string
  total_score: number | null  // NULL for Agent Config and API Doc
  grade: string | null        // NULL for Agent Config and API Doc
  input_data: Record<string, unknown>
  report_data: T
  generated_files: GeneratedFileMeta[]
  parent_report_id?: string | null  // lineage: re-run from this parent report
  project_id?: string | null
  created_at: string
  updated_at: string
}

// Tool-specific result types

export interface UIReviewResult {
  screenshot_contexts?: Array<{ kind: 'desktop' | 'mobile'; viewport: string }>
  scores: Array<{
    name: string
    score: number
    max_score: number
    comment: string
  }>
  issues: Issue[]
  recommendations: string[]
  action_items?: ActionItem[]
  codex_prompt: string
}

export interface ProjectDoctorResult {
  scores: Array<{
    name: string
    score: number
    max_score: number
    comment: string
  }>
  issues: Issue[]
  recommendations: string[]
  action_items?: ActionItem[]
  codex_prompt: string
}

export interface AgentConfigResult {
  generated_files_content: Record<string, string> // filename -> content
  recommendations: string[]
  action_items?: ActionItem[]
  codex_prompt: string
  // NOTE: no scores, total_score and grade are null at Report level
}

export interface APIDocResult {
  modules: Array<{
    name: string
    endpoints: unknown[]
  }>
  markdown_content?: string
  openapi_content?: string
  recommendations: string[]
  action_items?: ActionItem[]
  codex_prompt: string
  // NOTE: scores may be null
}

export interface DBSchemaResult {
  scores: Array<{
    name: string
    score: number
    max_score: number
    comment: string
  }>
  issues: Issue[]
  optimized_schema?: string
  migration_suggestions?: string[]
  recommendations: string[]
  action_items?: ActionItem[]
  codex_prompt: string
}

export interface ReportListParams {
  tool_type?: ToolType
  status?: ReportStatus
  sort?: 'newest' | 'oldest' | 'score_desc' | 'score_asc'
  page?: number
  page_size?: number
}

// ─── Report comparison (P1-06) ───

export interface ReportSummaryDTO {
  id: string
  title: string
  status: ReportStatus
  total_score: number | null
  grade: string | null
  created_at: string
  summary: string
  report_data: Record<string, unknown>
}

export interface IssueCountDelta {
  high: number
  medium: number
  low: number
  total: number
}

export interface IssueMatch {
  title: string
  category: string
  severity: 'high' | 'medium' | 'low' | string
  in_baseline: boolean
  in_target: boolean
}

export interface IssueComparison {
  resolved: IssueMatch[]
  new: IssueMatch[]
  persist: IssueMatch[]
}

export interface ActionItemDelta {
  resolved: ActionItem[]
  new: ActionItem[]
  persist: ActionItem[]
}

export interface ReportCompare {
  baseline_report: ReportSummaryDTO
  target_report: ReportSummaryDTO
  tool_type: ToolType
  score_delta?: number | null
  grade_delta?: string
  issue_count_delta: IssueCountDelta
  issues: IssueComparison
  action_items: ActionItemDelta
  warnings?: string[]
}
