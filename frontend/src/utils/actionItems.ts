import type { ActionEffort, ActionItem, ActionPriority } from '@/types/report'

export interface DisplayActionItem extends ActionItem {
  legacy?: boolean
}

const priorityRank: Record<ActionPriority, number> = {
  high: 0,
  medium: 1,
  low: 2,
}

const effortRank: Record<ActionEffort, number> = {
  small: 0,
  medium: 1,
  large: 2,
}

export function normalizeActionItems(
  actionItems: ActionItem[] | undefined,
  recommendations: string[] = [],
): DisplayActionItem[] {
  if (actionItems?.length) {
    return sortActionItems(actionItems.map((item) => ({ ...item })))
  }

  return sortActionItems(
    recommendations
      .filter((text) => text.trim().length > 0)
      .map((text, index) => ({
        id: `legacy-recommendation-${index + 1}`,
        title: text.trim(),
        priority: 'medium',
        effort: 'medium',
        category: 'recommendation',
        reason: '旧报告未包含统一 action_items 字段，此项由建议内容降级生成。',
        suggested_prompt: text.trim(),
        issue_title: text.trim(),
        issue_body: `## 背景\n旧报告未包含统一 action_items 字段。\n\n## 建议\n${text.trim()}\n\n## 验收\n- [ ] 已确认该建议是否仍适用\n- [ ] 已记录处理结果`,
        legacy: true,
      })),
  )
}

export function sortActionItems<T extends ActionItem>(items: T[]): T[] {
  return [...items].sort((a, b) => {
    const priorityDelta = priorityRank[a.priority] - priorityRank[b.priority]
    if (priorityDelta !== 0) return priorityDelta
    const effortDelta = effortRank[a.effort] - effortRank[b.effort]
    if (effortDelta !== 0) return effortDelta
    return a.title.localeCompare(b.title, 'zh-CN')
  })
}

export function groupActionItems(items: DisplayActionItem[]) {
  return {
    high: items.filter((item) => item.priority === 'high'),
    medium: items.filter((item) => item.priority === 'medium'),
    low: items.filter((item) => item.priority === 'low'),
  }
}

export function buildActionItemsMarkdown(
  reportId: string,
  reportTitle: string,
  items: DisplayActionItem[],
): string {
  const lines = [
    `# ${escapeMarkdown(reportTitle)} 行动计划`,
    '',
    `来源报告：${reportId} (/reports/${reportId})`,
    '',
  ]

  for (const item of sortActionItems(items)) {
    lines.push(`- [ ] **${escapeMarkdown(item.title)}**`)
    lines.push(`  - 优先级：${item.priority}`)
    lines.push(`  - 工作量：${item.effort}`)
    lines.push(`  - 分类：${escapeMarkdown(item.category)}`)
    lines.push(`  - 原因：${escapeMarkdown(item.reason)}`)
    lines.push(`  - Prompt：${escapeMarkdown(item.suggested_prompt)}`)
    if (item.issue_title) {
      lines.push(`  - Issue：${escapeMarkdown(item.issue_title)}`)
    }
    if (item.issue_body) {
      lines.push('  - Issue 草稿：')
      for (const line of item.issue_body.split('\n')) {
        lines.push(`    ${line}`)
      }
    }
    lines.push('')
  }

  return lines.join('\n').trimEnd() + '\n'
}

function escapeMarkdown(value: string): string {
  return value.replace(/\\/g, '\\\\').replace(/\|/g, '\\|')
}
