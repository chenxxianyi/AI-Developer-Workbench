package prompts

// UIReviewPromptSchema describes the expected output format.
const UIReviewPromptSchema = `
Output a JSON object with this exact structure:
{
  "scores": [
    {"name": "string - dimension name", "score": number (0-100), "max_score": 100, "comment": "string"}
  ],
  "issues": [
    {"title": "string", "severity": "high|medium|low", "category": "string", "problem": "string", "suggestion": "string", "action": "string"}
  ],
  "recommendations": ["string - actionable recommendations"],
  "codex_prompt": "string - prompt for generating CSS fixes"
}
`

// BuildUIReviewPrompt creates the prompt for UI Review.
func BuildUIReviewPrompt(reviewMode, pageType, targetStyle, description, code string) (string, string) {
	systemPrompt := `You are an expert UI/UX designer reviewing interfaces.
Your task is to analyze the provided UI screenshot or code for design quality.

Scoring dimensions:
1. Visual Hierarchy (0-100): Clear priority, size/color relationships
2. Consistency (0-100): Uniform spacing, typography, patterns
3. Accessibility (0-100): Focus states, contrast, labels
4. Color & Contrast (0-100): WCAG compliance, harmony
5. Responsive Design (0-100): Mobile adaptation, breakpoints
` + UIReviewPromptSchema

	userPrompt := `Review Information:
- Mode: ` + reviewMode + ` (screenshot, code, or screenshot_code)
- Page Type: ` + pageType + `
- Target Style: ` + targetStyle + `
- Description: ` + description + `

Code (if provided):
` + code + `

Analyze this UI and provide scores, issues, and actionable recommendations.`

	return BuildPrompt(systemPrompt, userPrompt)
}