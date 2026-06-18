package dto

// AgentConfigRequest is the input for Agent Config Studio.
type AgentConfigRequest struct {
	Title             string `json:"title" binding:"required"`
	ProjectName       string `json:"project_name" binding:"required"`
	ProjectType       string `json:"project_type"`
	FrontendStack     string `json:"frontend_stack"`
	BackendStack      string `json:"backend_stack"`
	Database          string `json:"database"`
	UIStyle           string `json:"ui_style"`
	CodingPreferences string `json:"coding_preferences"`
	StrictRules       string `json:"strict_rules"`
}

// AgentConfigResult is the output for Agent Config Studio.
type AgentConfigResult struct {
	GeneratedFilesContent map[string]string `json:"generated_files_content"`
	Recommendations       []string          `json:"recommendations"`
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
type IssueItem struct {
	Title      string `json:"title"`
	Severity   string `json:"severity"`
	Category   string `json:"category"`
	Problem    string `json:"problem"`
	Suggestion string `json:"suggestion"`
	Action     string `json:"action"`
}