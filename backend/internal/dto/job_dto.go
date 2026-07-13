package dto

// JobDTO represents an async job in API responses.
type JobDTO struct {
	ID           string  `json:"id"`
	ToolType     string  `json:"tool_type"`
	ReportID     string  `json:"report_id"`
	ProjectID    *string `json:"project_id,omitempty"`
	Status       string  `json:"status"`
	Progress     int     `json:"progress"`
	Phase        string  `json:"phase"`
	ErrorMessage string  `json:"error_message,omitempty"`
	RetryOfJobID *string `json:"retry_of_job_id,omitempty"`
	RetryCount   int     `json:"retry_count"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

// CreateJobDTO is the request body for creating a new job.
type CreateJobDTO struct {
	ToolType string `json:"tool_type"`
}
