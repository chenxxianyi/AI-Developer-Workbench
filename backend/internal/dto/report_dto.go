package dto

import "encoding/json"

// ReportDTO is the API response for a report, matching frontend Report<T> type.
type ReportDTO struct {
	ID             string             `json:"id"`
	ToolType       string             `json:"tool_type"`
	Title          string             `json:"title"`
	InputMode      string             `json:"input_mode"`
	Status         string             `json:"status"`
	Summary        string             `json:"summary"`
	TotalScore     *int               `json:"total_score"`
	Grade          *string            `json:"grade"`
	InputData      json.RawMessage    `json:"input_data"`
	ReportData     json.RawMessage    `json:"report_data"`
	GeneratedFiles []GeneratedFileDTO `json:"generated_files"`
	ParentReportID *string            `json:"parent_report_id,omitempty"`
	ProjectID      *string            `json:"project_id,omitempty"`
	CreatedAt      string             `json:"created_at"`
	UpdatedAt      string             `json:"updated_at"`
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
	TotalReports  int64                  `json:"total_reports"`
	ToolUsage     map[string]int64       `json:"tool_usage"`
	AverageScore  *float64               `json:"average_score"`
	RecentReports []RecentReportDTO      `json:"recent_reports"`
	WeeklyStats   *WeeklyStatsDTO        `json:"weekly_stats,omitempty"`
	QualityTrend  []QualityTrendPointDTO `json:"quality_trend,omitempty"`
}

// WeeklyStatsDTO holds the current-week dashboard summary.
type WeeklyStatsDTO struct {
	// ReportCountThisWeek is the number of reports created in the last 7 days.
	ReportCountThisWeek int64 `json:"report_count_this_week"`
	// AverageScoreThisWeek is the mean total_score of scored reports in the last 7 days.
	// Nil when no scored reports exist this week (not 0).
	AverageScoreThisWeek *float64 `json:"average_score_this_week"`
	// HighSeverityCountThisWeek counts issues with severity "high" across
	// this week's succeeded/fallback reports. Aggregated in Go from report_json.
	HighSeverityCountThisWeek int64 `json:"high_severity_count_this_week"`
	// MostUsedToolThisWeek is the tool_type with the most reports this week.
	// Empty when no reports exist this week.
	MostUsedToolThisWeek string `json:"most_used_tool_this_week"`
}

// QualityTrendPointDTO is one daily/weekly bucket of the quality trend.
type QualityTrendPointDTO struct {
	// Bucket is the ISO date (YYYY-MM-DD) the bucket starts on (UTC).
	Bucket string `json:"bucket"`
	// ReportCount is the number of reports in this bucket (all statuses).
	ReportCount int64 `json:"report_count"`
	// AverageScore is the mean total_score of scored reports in this bucket.
	// Nil when no scored reports exist in the bucket — UI shows "无数据".
	AverageScore *float64 `json:"average_score"`
	// HighSeverityCount counts high-severity issues in this bucket.
	HighSeverityCount int64 `json:"high_severity_count"`
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
	Healthy      bool            `json:"healthy"`
	Provider     string          `json:"provider"`
	TextModel    string          `json:"text_model"`
	VisionModel  string          `json:"vision_model"`
	MockMode     bool            `json:"mock_mode"`
	UploadLimits UploadLimitsDTO `json:"upload_limits"`
}

// UploadLimitsDTO matches frontend upload_limits.
type UploadLimitsDTO struct {
	ImageMaxBytes    int64 `json:"image_max_bytes"`
	ZipMaxBytes      int64 `json:"zip_max_bytes"`
	ZipMaxFiles      int   `json:"zip_max_files"`
	ZipMaxTotalBytes int64 `json:"zip_max_total_bytes"`
}
