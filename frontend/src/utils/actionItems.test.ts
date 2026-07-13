import { describe, expect, it } from 'vitest'
import type { ActionItem } from '@/types/report'
import { buildActionItemsMarkdown, normalizeActionItems, sortActionItems } from './actionItems'

const baseItem: ActionItem = {
  id: 'base',
  title: 'Base',
  priority: 'medium',
  effort: 'medium',
  category: 'general',
  reason: 'Reason',
  suggested_prompt: 'Prompt',
  issue_title: 'Issue',
  issue_body: 'Body',
}

describe('action item utilities', () => {
  it('sorts by priority, effort, then title', () => {
    const items: ActionItem[] = [
      { ...baseItem, id: 'b', title: 'B', priority: 'medium', effort: 'small' },
      { ...baseItem, id: 'c', title: 'C', priority: 'high', effort: 'large' },
      { ...baseItem, id: 'a', title: 'A', priority: 'high', effort: 'small' },
    ]

    expect(sortActionItems(items).map((item) => item.id)).toEqual(['a', 'c', 'b'])
  })

  it('uses action_items when available', () => {
    const items = normalizeActionItems([{ ...baseItem, id: 'real' }], ['legacy recommendation'])

    expect(items).toHaveLength(1)
    expect(items[0].id).toBe('real')
    expect(items[0].legacy).toBeUndefined()
  })

  it('falls back to recommendations for legacy reports', () => {
    const items = normalizeActionItems(undefined, ['补齐测试'])

    expect(items).toHaveLength(1)
    expect(items[0].id).toBe('legacy-recommendation-1')
    expect(items[0].legacy).toBe(true)
    expect(items[0].suggested_prompt).toBe('补齐测试')
  })

  it('builds GitHub-readable markdown checklist', () => {
    const markdown = buildActionItemsMarkdown('report-1', 'UI 审查报告', [
      { ...baseItem, title: '修复上传键盘操作', priority: 'high', effort: 'small' },
    ])

    expect(markdown).toContain('# UI 审查报告 行动计划')
    expect(markdown).toContain('来源报告：report-1 (/reports/report-1)')
    expect(markdown).toContain('- [ ] **修复上传键盘操作**')
    expect(markdown).toContain('Prompt：Prompt')
    expect(markdown).toContain('Issue 草稿')
    expect(markdown).toContain('Body')
  })
})
