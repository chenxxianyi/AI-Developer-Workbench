/**
 * Tools API
 * Run five core analysis tools
 */

import apiClient from './client'
import type { Report } from '@/types/report'
import type {
  ToolMeta,
  AgentConfigInput,
  DBSchemaInput,
} from '@/types/tool'

/**
 * Get tool metadata
 */
export async function getTools(): Promise<ToolMeta[]> {
  return apiClient.get('/tools')
}

/**
 * Run UI Review tool
 * Uses FormData for file upload
 */
export async function runUIReview(formData: FormData): Promise<Report> {
  return apiClient.post('/tools/ui-review/run', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })
}

/**
 * Run Project Doctor tool
 * Uses FormData for ZIP file upload
 */
export async function runProjectDoctor(formData: FormData): Promise<Report> {
  return apiClient.post('/tools/project-doctor/run', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })
}

/**
 * Run Agent Config Studio tool
 * Uses JSON payload (no file upload)
 */
export async function runAgentConfig(payload: AgentConfigInput): Promise<Report> {
  return apiClient.post('/tools/agent-config/run', payload)
}

/**
 * Run API Doc Builder tool
 * Uses FormData if source_type is project_zip
 */
export async function runAPIDoc(formData: FormData): Promise<Report> {
  return apiClient.post('/tools/api-doc/run', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })
}

/**
 * Run DB Schema Review tool
 * Uses JSON payload (no file upload)
 */
export async function runDBSchema(payload: DBSchemaInput): Promise<Report> {
  return apiClient.post('/tools/db-schema/run', payload)
}