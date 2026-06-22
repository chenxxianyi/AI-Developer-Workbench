import { describe, expect, it } from 'vitest'
import { getToolDisplayMeta } from './toolDisplay'

describe('tool display metadata', () => {
  it('uses Chinese copy for dashboard tool cards', () => {
    expect(getToolDisplayMeta('ui_review')).toEqual({
      name: 'UI 质量审查',
      shortDescription: '截图或代码 UI/UX 审查',
      description: '基于截图或前端代码进行 UI/UX 质量审查，识别模板化痕迹、设计一致性和可用性问题。',
    })

    expect(getToolDisplayMeta('project_doctor').name).toBe('项目诊断')
    expect(getToolDisplayMeta('agent_config').name).toBe('Agent 配置生成')
    expect(getToolDisplayMeta('api_doc').name).toBe('API 文档生成')
    expect(getToolDisplayMeta('db_schema').name).toBe('数据库结构审查')
  })
})
