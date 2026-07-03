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
  created_at: string
  updated_at: string
}

// Tool-specific result types

export interface UIReviewResult {
  scores: Array<{
    name: string
    score: number
    max_score: number
    comment: string
  }>
  issues: Issue[]
  recommendations: string[]
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
  codex_prompt: string
}

export interface AgentConfigResult {
  generated_files_content: Record<string, string> // filename -> content
  recommendations: string[]
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
  codex_prompt: string
}

export interface ReportListParams {
  tool_type?: ToolType
  status?: ReportStatus
  sort?: 'newest' | 'oldest' | 'score_desc' | 'score_asc'
  page?: number
  page_size?: number
}