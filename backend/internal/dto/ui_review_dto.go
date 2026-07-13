package dto

// UIReviewRequest is the input for UI Review.
type UIReviewRequest struct {
	Title           string `json:"title" binding:"required"`
	ProjectID       string `json:"project_id,omitempty"`
	ReviewMode      string `json:"review_mode" binding:"required"`
	CodeSource      string `json:"code_source"`
	PageType        string `json:"page_type"`
	TargetStyle     string `json:"target_style"`
	Description     string `json:"description"`
	DesktopViewport string `json:"desktop_viewport,omitempty"`
	MobileViewport  string `json:"mobile_viewport,omitempty"`
	Code            string `json:"code"`
}

// ScreenshotContext identifies the source viewport used by an annotated issue.
type ScreenshotContext struct {
	Kind     string `json:"kind"`
	Viewport string `json:"viewport"`
}

// UIReviewResult is the output for UI Review.
type UIReviewResult struct {
	ScreenshotContexts []ScreenshotContext `json:"screenshot_contexts,omitempty"`
	Scores             []ScoreItem         `json:"scores"`
	Issues             []IssueItem         `json:"issues"`
	Recommendations    []string            `json:"recommendations"`
	ActionItems        []ActionItem        `json:"action_items,omitempty"`
	CodexPrompt        string              `json:"codex_prompt"`
}
