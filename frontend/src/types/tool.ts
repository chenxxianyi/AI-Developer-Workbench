/**
 * Tool Types
 * Five core tool types and their input structures
 */

export type ToolType =
  | 'ui_review'
  | 'project_doctor'
  | 'agent_config'
  | 'api_doc'
  | 'db_schema'

export interface ToolMeta {
  tool_type: ToolType
  name: string
  description: string
  icon: string // Lucide icon name
  color: string // Tailwind color class key
  usage_count: number
}

// Tool-specific input types

export interface UIReviewInput {
  title: string
  review_mode: 'screenshot' | 'code' | 'screenshot_code'
  code_source?: 'paste' | 'project_zip'
  page_type?: string
  target_style?: string
  description?: string
  code?: string
  screenshot?: File
  project_zip?: File
}

export interface ProjectDoctorInput {
  title: string
  project_name: string
  tech_stack?: string
  project_description?: string
  analysis_depth: 'basic' | 'standard' | 'deep'
  project_zip: File
}

export interface AgentConfigInput {
  title: string
  project_name: string
  project_type?: string
  frontend_stack?: string
  backend_stack?: string
  database?: string
  ui_style?: string
  coding_preferences?: string
  strict_rules?: string
}

export interface APIDocInput {
  title: string
  source_type: 'code' | 'project_zip' | 'manual'
  backend_stack?: string
  code?: string
  api_description?: string
  output_format: 'markdown' | 'openapi' | 'both'
  project_zip?: File
}

export interface DBSchemaInput {
  title: string
  schema_type: string
  database_type?: string
  business_context?: string
  schema_content: string
  target_goal?: string
}
