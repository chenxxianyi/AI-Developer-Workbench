/**
 * Reports API
 * CRUD operations and file downloads
 */

import apiClient from './client'
import type { Report, ReportListParams } from '@/types/report'
import type { PaginatedData } from '@/types/api'
import { downloadBlob } from '@/utils/download'

/**
 * List reports with pagination, filter, and sort
 */
export async function listReports(params: ReportListParams): Promise<PaginatedData<Report>> {
  return apiClient.get('/reports', { params })
}

/**
 * Get report detail by ID
 */
export async function getReport(id: string): Promise<Report> {
  return apiClient.get(`/reports/${id}`)
}

/**
 * Delete report by ID
 */
export async function deleteReport(id: string): Promise<void> {
  return apiClient.delete(`/reports/${id}`)
}

/**
 * Download report as Markdown
 */
export async function downloadReport(id: string): Promise<void> {
  const response = await apiClient.get(`/reports/${id}/export`, {
    params: { format: 'markdown' },
    responseType: 'blob',
  })

  // Extract filename from Content-Disposition header
  const contentDisposition = response.headers['content-disposition']
  let filename = `report-${id}.md`

  if (contentDisposition) {
    const filenameMatch = contentDisposition.match(/filename="?(.+)"?/)
    if (filenameMatch && filenameMatch[1]) {
      filename = filenameMatch[1]
    }
  }

  await downloadBlob(response.data as Blob, filename)
}

/**
 * Download generated file from report
 */
export async function downloadGeneratedFile(
  reportId: string,
  filename: string
): Promise<void> {
  const response = await apiClient.get(
    `/reports/${reportId}/files/${encodeURIComponent(filename)}`,
    {
      responseType: 'blob',
    }
  )

  // Check if response is JSON error (not a file)
  const contentType = response.headers['content-type']
  if (contentType && typeof contentType === 'string' && contentType.includes('application/json')) {
    // Parse error message
    const text = await (response.data as Blob).text()
    const error = JSON.parse(text)
    throw error
  }

  await downloadBlob(response.data as Blob, filename)
}