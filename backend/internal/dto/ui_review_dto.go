package dto

// UIReviewRequest is the input for UI Review.
type UIReviewRequest struct {
	Title       string `json:"title" binding:"required"`
	ReviewMode  string `json:"review_mode" binding:"required"`
	CodeSource  string `json:"code_source"`
	PageType    string `json:"page_type"`
	TargetStyle string `json:"target_style"`
	Description string `json:"description"`
	Code        string `json:"code"`
}

// UIReviewResult is the output for UI Review.
type UIReviewResult struct {
	Scores         []ScoreItem  `json:"scores"`
	Issues         []IssueItem  `json:"issues"`
	Recommendations []string    `json:"recommendations"`
	CodexPrompt    string       `json:"codex_prompt"`
}
