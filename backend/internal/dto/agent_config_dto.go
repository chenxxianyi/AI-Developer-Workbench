package dto

// AgentConfigRequest is the input for Agent Config Studio.
type AgentConfigRequest struct {
	Title             string `json:"title" binding:"required"`
	ProjectID         string `json:"project_id,omitempty"`
	ProjectName       string `json:"project_name" binding:"required"`
	ProjectType       string `json:"project_type"`
	FrontendStack     string `json:"frontend_stack"`
	BackendStack      string `json:"backend_stack"`
	Database          string `json:"database"`
	UIStyle           string `json:"ui_style"`
	CodingPreferences string `json:"coding_preferences"`
	StrictRules       string `json:"strict_rules"`
	ParentReportID    string `json:"parent_report_id,omitempty"`
}

// AgentConfigResult is the output for Agent Config Studio.
type AgentConfigResult struct {
	GeneratedFilesContent map[string]string `json:"generated_files_content"`
	TargetFormat          string            `json:"target_format,omitempty"` // codex|copilot|cursor|windsurf
	MissingConfirmations  []string          `json:"missing_confirmations,omitempty"` // items needing user verification
	Recommendations       []string          `json:"recommendations"`
	ActionItems           []ActionItem      `json:"action_items,omitempty"`
	CodexPrompt           string            `json:"codex_prompt"`
}

// ScoreItem represents a score dimension.
type ScoreItem struct {
	Name     string `json:"name"`
	Score    int    `json:"score"`
	MaxScore int    `json:"max_score"`
	Comment  string `json:"comment"`
}

// IssueItem represents an identified issue.
type IssueRegion struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// IssueItem represents an identified issue. Region coordinates are percentages
// (0-100) so annotations remain stable across rendered screenshot sizes.
type IssueItem struct {
	Title              string       `json:"title"`
	Severity           string       `json:"severity"`
	Category           string       `json:"category"`
	Problem            string       `json:"problem"`
	Suggestion         string       `json:"suggestion"`
	Action             string       `json:"action"`
	Viewport           string       `json:"viewport,omitempty"`
	Region             *IssueRegion `json:"region,omitempty"`
	ContrastSuggestion string       `json:"contrast_suggestion,omitempty"`
	ComponentPrompt    string       `json:"component_prompt,omitempty"`
}
