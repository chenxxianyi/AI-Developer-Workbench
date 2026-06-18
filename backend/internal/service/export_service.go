package service

import (
	"context"
	"fmt"

	"ai-developer-workbench/internal/repository"
)

// ExportService handles report export operations.
type ExportService interface {
	ExportMarkdown(ctx context.Context, reportID string) ([]byte, string, error)
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

// GetFileContent retrieves the content of a specific generated file.
func (s *exportService) GetFileContent(ctx context.Context, reportID, filename string) ([]byte, string, string, error) {
	file, err := s.fileRepo.GetByReportIDAndFilename(ctx, reportID, filename)
	if err != nil {
		return nil, "", "", fmt.Errorf("file not found")
	}

	return []byte(file.Content), file.Filename, file.MimeType, nil
}