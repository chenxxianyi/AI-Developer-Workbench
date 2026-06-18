package dto

// ProjectDoctorRequest is the input for Project Doctor.
type ProjectDoctorRequest struct {
	Title               string `json:"title" binding:"required"`
	ProjectName         string `json:"project_name" binding:"required"`
	TechStack           string `json:"tech_stack"`
	ProjectDescription  string `json:"project_description"`
	AnalysisDepth       string `json:"analysis_depth" binding:"required"`
}

// ProjectDoctorResult is the output for Project Doctor.
type ProjectDoctorResult struct {
	Scores         []ScoreItem  `json:"scores"`
	Issues         []IssueItem  `json:"issues"`
	Recommendations []string    `json:"recommendations"`
	CodexPrompt    string       `json:"codex_prompt"`
}