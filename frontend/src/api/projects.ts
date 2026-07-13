/**
 * Projects API
 * CRUD operations and stats for project profiles
 */

import apiClient from './client'
import type { PaginatedData } from '@/types/api'
import type {
  Project,
  ProjectSummary,
  ProjectStats,
  ProjectCreateInput,
  ProjectUpdateInput,
  ProjectListParams,
  ProjectDeleteResult,
} from '@/types/project'
import type { Report, ReportListParams } from '@/types/report'

export async function listProjects(params: ProjectListParams): Promise<PaginatedData<ProjectSummary>> {
  return apiClient.get('/projects', { params })
}

export async function getProject(id: string): Promise<Project> {
  return apiClient.get(`/projects/${id}`)
}

export async function createProject(input: ProjectCreateInput): Promise<Project> {
  return apiClient.post('/projects', input)
}

export async function updateProject(id: string, input: ProjectUpdateInput): Promise<Project> {
  return apiClient.patch(`/projects/${id}`, input)
}

export async function deleteProject(id: string): Promise<ProjectDeleteResult> {
  return apiClient.delete(`/projects/${id}`)
}

export async function getProjectStats(id: string): Promise<ProjectStats> {
  return apiClient.get(`/projects/${id}/stats`)
}

export async function listProjectReports(
  id: string,
  params: Pick<ReportListParams, 'page' | 'page_size'>,
): Promise<PaginatedData<Report>> {
  return apiClient.get(`/projects/${id}/reports`, { params })
}
