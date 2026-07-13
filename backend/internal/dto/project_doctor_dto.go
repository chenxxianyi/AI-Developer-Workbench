package dto

// ProjectDoctorRequest is the input for Project Doctor.
type ProjectDoctorRequest struct {
	Title              string `json:"title" binding:"required"`
	ProjectID          string `json:"project_id,omitempty"`
	ProjectName        string `json:"project_name" binding:"required"`
	TechStack          string `json:"tech_stack"`
	ProjectDescription string `json:"project_description"`
	AnalysisDepth      string `json:"analysis_depth" binding:"required"`
	ParentReportID     string `json:"parent_report_id,omitempty"`
}

// ProjectDoctorResult is the output for Project Doctor.
type ProjectDoctorResult struct {
	Scores          []ScoreItem  `json:"scores"`
	Issues          []IssueItem  `json:"issues"`
	Recommendations []string     `json:"recommendations"`
	ActionItems     []ActionItem `json:"action_items,omitempty"`
	CodexPrompt     string       `json:"codex_prompt"`
	// Enhanced fields
	EvidenceFiles   []EvidenceItem   `json:"evidence_files,omitempty"`
	TechDebtItems   []TechDebtItem   `json:"tech_debt,omitempty"`
}

// EvidenceItem represents a detected project artifact.
type EvidenceItem struct {
	Path    string `json:"path"`
	Type    string `json:"type"` // readme, agents_md, lockfile, dockerfile, ci, tests, docs
	Present bool   `json:"present"`
	Notes   string `json:"notes,omitempty"`
}

// TechDebtItem represents a prioritized technical debt item.
type TechDebtItem struct {
	Title       string `json:"title"`
	Impact      string `json:"impact"` // high|medium|low
	Cost        string `json:"cost"`   // high|medium|low
	Category    string `json:"category"`
	Description string `json:"description"`
	SuggestedFix string `json:"suggested_fix,omitempty"`
}
