package dto

import "encoding/json"

// ReportDTO is the API response for a report, matching frontend Report<T> type.
type ReportDTO struct {
	ID             string               `json:"id"`
	ToolType       string               `json:"tool_type"`
	Title          string               `json:"title"`
	InputMode      string               `json:"input_mode"`
	Status         string               `json:"status"`
	Summary        string               `json:"summary"`
	TotalScore     *int                 `json:"total_score"`
	Grade          *string              `json:"grade"`
	InputData      json.RawMessage      `json:"input_data"`
	ReportData     json.RawMessage      `json:"report_data"`
	GeneratedFiles []GeneratedFileDTO   `json:"generated_files"`
	CreatedAt      string               `json:"created_at"`
	UpdatedAt      string               `json:"updated_at"`
}

// GeneratedFileDTO is the metadata for a generated file.
type GeneratedFileDTO struct {
	ID        string `json:"id"`
	Filename  string `json:"filename"`
	Language  string `json:"language,omitempty"`
	MimeType  string `json:"mime_type"`
	SizeBytes uint64 `json:"size_bytes"`
}

// DashboardStatsDTO matches frontend DashboardStats type.
type DashboardStatsDTO struct {
	TotalReports  int64                       `json:"total_reports"`
	ToolUsage     map[string]int64            `json:"tool_usage"`
	AverageScore  *float64                    `json:"average_score"`
	RecentReports []RecentReportDTO           `json:"recent_reports"`
}

// RecentReportDTO is a summary of a recent report for the dashboard.
type RecentReportDTO struct {
	ID         string  `json:"id"`
	ToolType   string  `json:"tool_type"`
	Title      string  `json:"title"`
	Status     string  `json:"status"`
	TotalScore *int    `json:"total_score"`
	Grade      *string `json:"grade"`
	Summary    string  `json:"summary"`
	CreatedAt  string  `json:"created_at"`
}

// SystemStatusDTO matches frontend SystemStatus type.
type SystemStatusDTO struct {
	Healthy      bool             `json:"healthy"`
	Provider     string           `json:"provider"`
	TextModel    string           `json:"text_model"`
	VisionModel  string           `json:"vision_model"`
	MockMode     bool             `json:"mock_mode"`
	UploadLimits UploadLimitsDTO  `json:"upload_limits"`
}

// UploadLimitsDTO matches frontend upload_limits.
type UploadLimitsDTO struct {
	ImageMaxBytes    int64 `json:"image_max_bytes"`
	ZipMaxBytes      int64 `json:"zip_max_bytes"`
	ZipMaxFiles      int   `json:"zip_max_files"`
	ZipMaxTotalBytes int64 `json:"zip_max_total_bytes"`
}
