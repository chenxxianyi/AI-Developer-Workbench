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
}
