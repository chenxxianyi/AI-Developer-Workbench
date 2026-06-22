import type { ToolType } from '@/types/tool'

interface ToolDisplayMeta {
  name: string
  shortDescription: string
  description: string
}

const toolDisplayMeta: Record<ToolType, ToolDisplayMeta> = {
  ui_review: {
    name: 'UI 质量审查',
    shortDescription: '截图或代码 UI/UX 审查',
    description: '基于截图或前端代码进行 UI/UX 质量审查，识别模板化痕迹、设计一致性和可用性问题。',
  },
  project_doctor: {
    name: '项目诊断',
    shortDescription: '项目健康度与工程质量检查',
    description: '静态分析项目结构、依赖管理、代码规范和潜在风险，生成工程健康度报告。',
  },
  agent_config: {
    name: 'Agent 配置生成',
    shortDescription: '生成 AI Agent 配置文件',
    description: '根据项目特征生成 AGENTS.md、TASK_PLAN.md 等配置文件，优化 AI Coding 工具协作效果。',
  },
  api_doc: {
    name: 'API 文档生成',
    shortDescription: '自动生成 API 文档',
    description: '从代码、项目文件或接口描述生成 Markdown/OpenAPI 文档，降低接口维护成本。',
  },
  db_schema: {
    name: '数据库结构审查',
    shortDescription: '分析并优化数据库结构',
    description: '审查 SQL、GORM、Prisma 等数据库定义，评估表结构、索引、性能和安全问题。',
  },
}

export function getToolDisplayMeta(toolType: ToolType): ToolDisplayMeta {
  return toolDisplayMeta[toolType]
}
