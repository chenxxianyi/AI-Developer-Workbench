package prompts

// UIReviewPromptSchema describes the expected output format.
const UIReviewPromptSchema = `
Output a JSON object with this exact structure:
{
  "scores": [
    {"name": "string - Chinese dimension name", "score": number (0-100), "max_score": 100, "comment": "string - Simplified Chinese"}
  ],
  "issues": [
    {"title": "string - Simplified Chinese", "severity": "high|medium|low", "category": "string - Simplified Chinese", "problem": "string - Simplified Chinese", "suggestion": "string - Simplified Chinese", "action": "string - Simplified Chinese"}
  ],
  "recommendations": ["string - Simplified Chinese actionable recommendations"],
  "codex_prompt": "string - Simplified Chinese prompt for generating CSS fixes"
}

Language requirements:
- 必须使用简体中文输出所有面向用户的内容，包括 scores.name、scores.comment、issues.title、issues.category、issues.problem、issues.suggestion、issues.action、recommendations 和 codex_prompt。
- 可以保留 UI、UX、WCAG、CSS、HTML、Vue、React、ARIA、CTA 等行业术语或专有名词。
- severity 字段必须继续使用 high、medium、low 三个枚举值，不要翻译。
`

// BuildUIReviewPrompt creates the prompt for UI Review.
func BuildUIReviewPrompt(reviewMode, pageType, targetStyle, description, code string) (string, string) {
	systemPrompt := `You are an expert UI/UX designer reviewing interfaces.
Your task is to analyze the provided UI screenshot or code for design quality.
All explanations, issue descriptions, comments, and recommendations must be written in Simplified Chinese unless they are proper nouns or technical terms.

Scoring dimensions, use these Chinese names exactly:
1. 视觉层级 (0-100): 信息优先级、字号、颜色和布局关系是否清晰
2. 一致性 (0-100): 间距、字体、组件模式和交互是否统一
3. 可访问性 (0-100): 焦点状态、标签、键盘可用性和辅助技术支持
4. 色彩与对比度 (0-100): WCAG 对比度、色彩协调性和可读性
5. 响应式设计 (0-100): 移动端适配、断点和不同屏幕宽度下的布局表现
` + UIReviewPromptSchema

	userPrompt := `评审信息:
- 模式: ` + reviewMode + ` (screenshot, code, or screenshot_code)
- 页面类型: ` + pageType + `
- 目标风格: ` + targetStyle + `
- 补充说明: ` + description + `

代码（如有）:
` + code + `

请分析该 UI，并用简体中文提供评分、发现的问题和可执行的改进建议。`

	return BuildPrompt(systemPrompt, userPrompt)
}
