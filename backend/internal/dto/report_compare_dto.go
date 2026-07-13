package dto

import "encoding/json"

// ReportCompareDTO is the response for GET /api/reports/:id/compare/:targetId.
// It describes the delta between two same-tool reports (baseline id → target targetId).
type ReportCompareDTO struct {
	BaselineReport  ReportSummaryDTO `json:"baseline_report"`
	TargetReport    ReportSummaryDTO `json:"target_report"`
	ToolType        string           `json:"tool_type"`
	ScoreDelta      *int             `json:"score_delta,omitempty"`
	GradeDelta      string           `json:"grade_delta,omitempty"`
	IssueCountDelta IssueCountDelta  `json:"issue_count_delta"`
	Issues          IssueComparison  `json:"issues"`
	ActionItems     ActionItemDelta  `json:"action_items"`
	Warnings        []string         `json:"warnings,omitempty"`
}

// ReportSummaryDTO is a minimal, comparison-safe report snapshot.
type ReportSummaryDTO struct {
	ID         string          `json:"id"`
	Title      string          `json:"title"`
	Status     string          `json:"status"`
	TotalScore *int            `json:"total_score"`
	Grade      *string         `json:"grade"`
	CreatedAt  string          `json:"created_at"`
	Summary    string          `json:"summary"`
	ReportData json.RawMessage `json:"report_data"`
}

// IssueCountDelta expresses the change in issue counts by severity.
type IssueCountDelta struct {
	High   int `json:"high"`
	Medium int `json:"medium"`
	Low    int `json:"low"`
	Total  int `json:"total"`
}

// IssueComparison groups issues by their resolution status.
type IssueComparison struct {
	Resolved []IssueMatch `json:"resolved"`
	New      []IssueMatch `json:"new"`
	Persist  []IssueMatch `json:"persist"`
}

// IssueMatch is an issue paired with the reports it appears in.
type IssueMatch struct {
	Title     string `json:"title"`
	Category  string `json:"category"`
	Severity  string `json:"severity"`
	InBaseline bool  `json:"in_baseline"`
	InTarget   bool  `json:"in_target"`
}

// ActionItemDelta summarizes action item changes.
type ActionItemDelta struct {
	Resolved []ActionItem `json:"resolved"`
	New      []ActionItem `json:"new"`
	Persist  []ActionItem `json:"persist"`
}
