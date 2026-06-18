/**
 * Validation Utilities
 * Form validation helpers for tool inputs
 */

import type {
  UIReviewInput,
  ProjectDoctorInput,
  AgentConfigInput,
  APIDocInput,
  DBSchemaInput,
} from '@/types/tool'

/**
 * Validate UI Review input
 * @param input UI Review input data
 * @returns Validation result with errors
 */
export function validateUIReviewInput(input: Partial<UIReviewInput>): {
  valid: boolean
  errors: string[]
} {
  const errors: string[] = []

  if (!input.title || input.title.trim().length === 0) {
    errors.push('请输入报告标题')
  }

  if (!input.review_mode) {
    errors.push('请选择审查模式')
  }

  // Mode-dependent validation
  if (input.review_mode === 'screenshot') {
    if (!input.screenshot) {
      errors.push('请上传 UI 截图')
    }
  } else if (input.review_mode === 'code') {
    if (!input.code || input.code.trim().length === 0) {
      errors.push('请输入前端代码')
    }
  } else if (input.review_mode === 'screenshot_code') {
    if (!input.screenshot) {
      errors.push('请上传 UI 截图')
    }
    if (!input.code || input.code.trim().length === 0) {
      errors.push('请输入前端代码')
    }
  }

  return {
    valid: errors.length === 0,
    errors,
  }
}

/**
 * Validate Project Doctor input
 * @param input Project Doctor input data
 * @returns Validation result with errors
 */
export function validateProjectDoctorInput(input: Partial<ProjectDoctorInput>): {
  valid: boolean
  errors: string[]
} {
  const errors: string[] = []

  if (!input.title || input.title.trim().length === 0) {
    errors.push('请输入报告标题')
  }

  if (!input.project_zip) {
    errors.push('请上传项目 ZIP 文件')
  }

  return {
    valid: errors.length === 0,
    errors,
  }
}

/**
 * Validate Agent Config input
 * @param input Agent Config input data
 * @returns Validation result with errors
 */
export function validateAgentConfigInput(input: Partial<AgentConfigInput>): {
  valid: boolean
  errors: string[]
} {
  const errors: string[] = []

  if (!input.title || input.title.trim().length === 0) {
    errors.push('请输入报告标题')
  }

  return {
    valid: errors.length === 0,
    errors,
  }
}

/**
 * Validate API Doc input
 * @param input API Doc input data
 * @returns Validation result with errors
 */
export function validateAPIDocInput(input: Partial<APIDocInput>): {
  valid: boolean
  errors: string[]
} {
  const errors: string[] = []

  if (!input.title || input.title.trim().length === 0) {
    errors.push('请输入报告标题')
  }

  if (!input.source_type) {
    errors.push('请选择输入源类型')
  }

  // Source type dependent validation
  if (input.source_type === 'code') {
    if (!input.code || input.code.trim().length === 0) {
      errors.push('请输入代码')
    }
  } else if (input.source_type === 'project_zip') {
    if (!input.project_zip) {
      errors.push('请上传项目 ZIP 文件')
    }
  } else if (input.source_type === 'manual') {
    if (!input.api_description || input.api_description.trim().length === 0) {
      errors.push('请输入 API 描述')
    }
  }

  return {
    valid: errors.length === 0,
    errors,
  }
}

/**
 * Validate DB Schema input
 * @param input DB Schema input data
 * @returns Validation result with errors
 */
export function validateDBSchemaInput(input: Partial<DBSchemaInput>): {
  valid: boolean
  errors: string[]
} {
  const errors: string[] = []

  if (!input.title || input.title.trim().length === 0) {
    errors.push('请输入报告标题')
  }

  if (!input.schema_content || input.schema_content.trim().length === 0) {
    errors.push('请输入数据库 Schema')
  }

  return {
    valid: errors.length === 0,
    errors,
  }
}

/**
 * Validate file type
 * @param file File to validate
 * @param acceptedTypes Accepted MIME types or extensions
 * @returns Whether file type is accepted
 */
export function validateFileType(
  file: File,
  acceptedTypes: string[]
): boolean {
  const fileType = file.type
  const fileName = file.name.toLowerCase()

  for (const accepted of acceptedTypes) {
    // Check MIME type
    if (accepted.includes('/') && fileType === accepted) {
      return true
    }

    // Check extension
    if (accepted.startsWith('.') && fileName.endsWith(accepted)) {
      return true
    }
  }

  return false
}

/**
 * Validate file size
 * @param file File to validate
 * @param maxSizeBytes Maximum size in bytes
 * @returns Whether file size is within limit
 */
export function validateFileSize(file: File, maxSizeBytes: number): boolean {
  return file.size <= maxSizeBytes
}