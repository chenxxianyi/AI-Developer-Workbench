package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"strings"

	"ai-developer-workbench/internal/dto"
	"ai-developer-workbench/internal/repository"
)

var ErrNoActionItems = errors.New("report has no action items")

// ExportService handles report export operations.
type ExportService interface {
	ExportMarkdown(ctx context.Context, reportID string) ([]byte, string, error)
	ExportGitHubIssues(ctx context.Context, reportID string) ([]byte, string, error)
	GetFileContent(ctx context.Context, reportID, filename string) ([]byte, string, string, error)
}

type exportService struct {
	reportRepo repository.ReportRepository
	fileRepo   repository.GeneratedFileRepository
}

// NewExportService creates a new export service.
func NewExportService(reportRepo repository.ReportRepository, fileRepo repository.GeneratedFileRepository) ExportService {
	return &exportService{reportRepo: reportRepo, fileRepo: fileRepo}
}

// ExportMarkdown generates a Markdown export of a report.
func (s *exportService) ExportMarkdown(ctx context.Context, reportID string) ([]byte, string, error) {
	report, err := s.reportRepo.GetByID(ctx, reportID)
	if err != nil {
		return nil, "", fmt.Errorf("report not found")
	}

	var md string
	md = "# " + report.Title + "\n\n"
	md += fmt.Sprintf("- **Tool**: %s\n", report.ToolType)
	md += fmt.Sprintf("- **Status**: %s\n", report.Status)
	md += fmt.Sprintf("- **Created**: %s\n", report.CreatedAt.Format("2006-01-02 15:04:05"))

	if report.Summary != "" {
		md += fmt.Sprintf("\n## Summary\n\n%s\n", report.Summary)
	}
	if report.TotalScore != nil {
		md += fmt.Sprintf("\n**Score**: %d/100", *report.TotalScore)
		if report.Grade != nil {
			md += fmt.Sprintf(" (Grade: %s)", *report.Grade)
		}
		md += "\n"
	}
	if report.ErrorMessage != "" {
		md += fmt.Sprintf("\n## Error\n\n%s\n", report.ErrorMessage)
	}

	// Append generated files.
	if len(report.GeneratedFiles) > 0 {
		for _, f := range report.GeneratedFiles {
			md += fmt.Sprintf("\n---\n\n## %s\n\n", f.Filename)
			md += f.Content + "\n"
		}
	}

	filename := report.ToolType + "_report.md"
	return []byte(md), filename, nil
}

// ExportGitHubIssues generates Markdown issue drafts from report action items.
func (s *exportService) ExportGitHubIssues(ctx context.Context, reportID string) ([]byte, string, error) {
	report, err := s.reportRepo.GetByID(ctx, reportID)
	if err != nil {
		return nil, "", fmt.Errorf("report not found")
	}

	var payload struct {
		ActionItems []dto.ActionItem `json:"action_items"`
	}
	if err := json.Unmarshal(report.ReportJSON, &payload); err != nil {
		return nil, "", fmt.Errorf("parse report data: %w", err)
	}
	actionItems := dto.NormalizeActionItems(payload.ActionItems)
	if len(actionItems) == 0 {
		return nil, "", ErrNoActionItems
	}

	labels := defaultIssueLabels(report.ToolType)
	var md strings.Builder
	md.WriteString("# GitHub Issue Drafts\n\n")
	md.WriteString(fmt.Sprintf("Source report: `%s`\n\n", report.ID))

	for i, item := range actionItems {
		md.WriteString(fmt.Sprintf("## %d. %s\n\n", i+1, markdownSafe(item.IssueTitle)))
		md.WriteString(fmt.Sprintf("- Priority: `%s`\n", item.Priority))
		md.WriteString(fmt.Sprintf("- Effort: `%s`\n", item.Effort))
		md.WriteString(fmt.Sprintf("- Category: `%s`\n", markdownSafe(item.Category)))
		md.WriteString(fmt.Sprintf("- Suggested labels: `%s`\n\n", strings.Join(labels, "`, `")))
		md.WriteString("### Background\n\n")
		md.WriteString(markdownSafe(item.Reason))
		md.WriteString("\n\n### Fix Requirements\n\n")
		md.WriteString(markdownSafe(item.SuggestedPrompt))
		md.WriteString("\n\n### Draft Body\n\n")
		if item.IssueBody != "" {
			md.WriteString(markdownSafe(item.IssueBody))
		} else {
			md.WriteString(defaultIssueBody(item))
		}
		md.WriteString("\n\n### Acceptance Criteria\n\n")
		md.WriteString("- [ ] The fix requirements above are implemented\n")
		md.WriteString("- [ ] Relevant automated tests pass\n")
		md.WriteString("- [ ] The follow-up report or verification note is linked\n")
		md.WriteString("\n\n### Source\n\n")
		md.WriteString(fmt.Sprintf("- Report ID: `%s`\n- Report title: %s\n- Action item ID: `%s`\n\n", report.ID, markdownSafe(report.Title), markdownSafe(item.ID)))
	}

	return []byte(md.String()), report.ToolType + "_github_issues.md", nil
}

// GetFileContent retrieves the content of a specific generated file.
func (s *exportService) GetFileContent(ctx context.Context, reportID, filename string) ([]byte, string, string, error) {
	file, err := s.fileRepo.GetByReportIDAndFilename(ctx, reportID, filename)
	if err != nil {
		return nil, "", "", fmt.Errorf("file not found")
	}

	return []byte(file.Content), file.Filename, file.MimeType, nil
}

func defaultIssueLabels(toolType string) []string {
	switch toolType {
	case "ui_review":
		return []string{"ui", "accessibility", "quality"}
	case "project_doctor":
		return []string{"technical-debt", "project-health", "quality"}
	case "api_doc":
		return []string{"documentation", "api"}
	case "db_schema":
		return []string{"database", "schema"}
	case "agent_config":
		return []string{"agent-config", "documentation"}
	default:
		return []string{"quality"}
	}
}

func markdownSafe(value string) string {
	escaped := html.EscapeString(value)
	escaped = strings.ReplaceAll(escaped, "\r\n", "\n")
	return strings.ReplaceAll(escaped, "\r", "\n")
}

func defaultIssueBody(item dto.ActionItem) string {
	var b strings.Builder
	b.WriteString("## Problem\n\n")
	if item.Reason != "" {
		b.WriteString(markdownSafe(item.Reason))
	} else {
		b.WriteString("The source report identified this action item as needing follow-up.")
	}
	b.WriteString("\n\n## Fix\n\n")
	if item.SuggestedPrompt != "" {
		b.WriteString(markdownSafe(item.SuggestedPrompt))
	} else {
		b.WriteString("Apply the remediation described by the action item.")
	}
	return b.String()
}
