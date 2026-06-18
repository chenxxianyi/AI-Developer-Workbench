/**
 * Format Utilities
 * File size, date, and text formatting
 */

/**
 * Format file size in human-readable format
 * @param bytes File size in bytes
 * @returns Formatted string (e.g., "2.4 KB")
 */
export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 Bytes'

  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))

  return `${Math.round((bytes / Math.pow(k, i)) * 100) / 100} ${sizes[i]}`
}

/**
 * Format date to locale string
 * @param dateString ISO date string
 * @returns Formatted date string
 */
export function formatDate(dateString: string): string {
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

/**
 * Format relative time (e.g., "2小时前")
 * @param dateString ISO date string
 * @returns Relative time string
 */
export function formatRelativeTime(dateString: string): string {
  const date = new Date(dateString)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()

  const minutes = Math.floor(diffMs / 60_000)
  const hours = Math.floor(diffMs / 3_600_000)
  const days = Math.floor(diffMs / 86_400_000)

  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  if (hours < 24) return `${hours}小时前`
  if (days < 7) return `${days}天前`

  return formatDate(dateString)
}

/**
 * Get score color class based on value
 * @param score Score value (0-100)
 * @returns Tailwind color class name
 */
export function getScoreColor(score: number | null): string {
  if (score === null) return 'text-text-muted'

  if (score >= 80) return 'text-success'
  if (score >= 60) return 'text-warning'
  return 'text-danger'
}

/**
 * Get severity color class
 * @param severity Issue severity level
 * @returns Tailwind color class name
 */
export function getSeverityColor(severity: 'high' | 'medium' | 'low'): string {
  switch (severity) {
    case 'high':
      return 'danger'
    case 'medium':
      return 'warning'
    case 'low':
      return 'text-muted'
  }
}