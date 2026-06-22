import { describe, expect, it } from 'vitest'
import { getScoreDisplayName, getSeverityDisplayName } from './uiReviewDisplay'

describe('UI review display helpers', () => {
  it('translates common score names to Chinese', () => {
    expect(getScoreDisplayName('Visual Hierarchy')).toBe('视觉层级')
    expect(getScoreDisplayName('Consistency')).toBe('一致性')
    expect(getScoreDisplayName('Accessibility')).toBe('可访问性')
    expect(getScoreDisplayName('Color & Contrast')).toBe('色彩与对比度')
    expect(getScoreDisplayName('Responsive Design')).toBe('响应式设计')
  })

  it('translates severity labels to Chinese', () => {
    expect(getSeverityDisplayName('high')).toBe('高')
    expect(getSeverityDisplayName('medium')).toBe('中')
    expect(getSeverityDisplayName('low')).toBe('低')
  })
})
