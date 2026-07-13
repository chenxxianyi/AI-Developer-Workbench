/**
 * Reports API
 * CRUD operations and file downloads
 */

import apiClient from './client'
import type { Report, ReportListParams, ReportCompare } from '@/types/report'
import type { ApiErrorResponse, PaginatedData } from '@/types/api'
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
  await throwIfBlobError(response.data, response.headers['content-type'])

  const filename = getDownloadFilename(
    response.headers['content-disposition'],
    `report-${id}.md`,
  )

  await downloadBlob(response.data as Blob, filename)
}

/**
 * Download report action items as GitHub Issue draft Markdown
 */
export async function downloadGitHubIssues(id: string): Promise<void> {
  const response = await apiClient.get(`/reports/${id}/export`, {
    params: { format: 'github-issues' },
    responseType: 'blob',
  })
  await throwIfBlobError(response.data, response.headers['content-type'])

  const filename = getDownloadFilename(
    response.headers['content-disposition'],
    `report-${id}-github-issues.md`,
  )

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

  await throwIfBlobError(response.data, response.headers['content-type'])

  await downloadBlob(response.data as Blob, filename)
}

/**
 * Compare two reports (baseline vs target)
 */
export async function compareReports(
  baselineId: string,
  targetId: string,
): Promise<ReportCompare> {
  return apiClient.get(`/reports/${baselineId}/compare/${targetId}`)
}

async function throwIfBlobError(data: unknown, contentType: unknown): Promise<void> {
  if (!(data instanceof Blob) || typeof contentType !== 'string' || !contentType.includes('application/json')) {
    return
  }

  try {
    const error = JSON.parse(await data.text()) as ApiErrorResponse
    throw new Error(error.error || error.message || '下载失败')
  } catch (error) {
    if (error instanceof SyntaxError) {
      throw new Error('下载失败')
    }
    throw error
  }
}

function getDownloadFilename(contentDisposition: unknown, fallback: string): string {
  if (typeof contentDisposition !== 'string' || !contentDisposition.trim()) {
    return fallback
  }

  const utf8Match = contentDisposition.match(/filename\*=UTF-8''([^;]+)/i)
  if (utf8Match?.[1]) {
    try {
      return decodeURIComponent(utf8Match[1].trim())
    } catch {
      return utf8Match[1].trim()
    }
  }

  const quotedMatch = contentDisposition.match(/filename="([^"]+)"/i)
  if (quotedMatch?.[1]) {
    return quotedMatch[1].trim()
  }

  const plainMatch = contentDisposition.match(/filename=([^;]+)/i)
  if (plainMatch?.[1]) {
    return plainMatch[1].trim()
  }

  return fallback
}
