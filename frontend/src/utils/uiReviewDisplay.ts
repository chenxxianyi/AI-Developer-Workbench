const scoreNameMap: Record<string, string> = {
  'Visual Hierarchy': '视觉层级',
  Consistency: '一致性',
  Accessibility: '可访问性',
  'Color & Contrast': '色彩与对比度',
  'Responsive Design': '响应式设计',
  Overall: '总体评分',
}

const severityMap: Record<string, string> = {
  high: '高',
  medium: '中',
  low: '低',
}

export function getScoreDisplayName(name: string): string {
  return scoreNameMap[name] ?? name
}

export function getSeverityDisplayName(severity: string): string {
  return severityMap[severity.toLowerCase()] ?? severity
}
