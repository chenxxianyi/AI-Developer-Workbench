package prompts

// UIReviewPromptSchema describes the expected output format.
const UIReviewPromptSchema = `
Output a JSON object with this exact structure:
{
  "scores": [
    {"name": "string - Chinese dimension name", "score": number (0-100), "max_score": 100, "comment": "string - Simplified Chinese"}
  ],
  "screenshot_contexts": [{"kind":"desktop|mobile", "viewport":"WIDTHxHEIGHT"}],
  "issues": [
    {"title":"string", "severity":"high|medium|low", "category":"string", "problem":"string", "suggestion":"string", "action":"string", "viewport":"desktop|mobile", "region":{"x":0-100,"y":0-100,"width":0-100,"height":0-100}, "contrast_suggestion":"WCAG-aware fix", "component_prompt":"executable component-level fix prompt"}
  ],
  "recommendations": ["string - Simplified Chinese actionable recommendations"],
  "action_items": [
    {"id": "stable-kebab-case-id", "title": "string - Simplified Chinese", "priority": "high|medium|low", "effort": "small|medium|large", "category": "string", "reason": "string - Simplified Chinese", "suggested_prompt": "string - Simplified Chinese", "issue_title": "string - Simplified Chinese", "issue_body": "string - Markdown"}
  ],
  "codex_prompt": "string - Simplified Chinese prompt for generating CSS fixes"
}

Language requirements:
- 必须使用简体中文输出所有面向用户的内容，包括 scores.name、scores.comment、issues.title、issues.category、issues.problem、issues.suggestion、issues.action、recommendations 和 codex_prompt。
- action_items 中 title、reason、suggested_prompt、issue_title 和 issue_body 必须使用简体中文。
- 可以保留 UI、UX、WCAG、CSS、HTML、Vue、React、ARIA、CTA 等行业术语或专有名词。
- severity 字段必须继续使用 high、medium、low 三个枚举值，不要翻译。
` + ActionItemsPromptSchema

// BuildUIReviewPrompt creates the prompt for UI Review.
func BuildUIReviewPrompt(reviewMode, codeSource, pageType, targetStyle, description, code, projectSummary string) (string, string) {
	systemPrompt := `You are an expert UI/UX designer reviewing interfaces.
Your task is to analyze the provided UI screenshots or code for design quality.
For screenshot reviews, attribute every issue to desktop or mobile and provide percentage-based region coordinates. Every high-severity issue must include a non-empty component_prompt. Include WCAG contrast guidance when color is involved. Apply page-specific priorities for login, dashboard, form, and content pages.
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
- 代码来源: ` + codeSource + ` (paste or project_zip)
- 页面类型: ` + pageType + `
- 目标风格: ` + targetStyle + `
- 补充说明: ` + description + `

代码（如有）:
` + code + `

前端项目 ZIP 源码摘要（如有）:
` + projectSummary + `

请分析该 UI，并用简体中文提供评分、发现的问题和可执行的改进建议。`

	return BuildPrompt(systemPrompt, userPrompt)
}
