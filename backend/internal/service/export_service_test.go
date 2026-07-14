package service

import (
	"context"
	"testing"

	"ai-developer-workbench/internal/dto"
	"ai-developer-workbench/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// stubReportRepo implements repository.ReportRepository for export tests.
type stubReportRepo struct {
	report *model.Report
}

func (s *stubReportRepo) Create(_ context.Context, _ *model.Report) error { return nil }
func (s *stubReportRepo) Update(_ context.Context, _ *model.Report) error { return nil }
func (s *stubReportRepo) GetByID(_ context.Context, _ string) (*model.Report, error) {
	return s.report, nil
}
func (s *stubReportRepo) List(_ context.Context, _ dto.ListReportsQuery) ([]model.Report, int64, error) {
	return nil, 0, nil
}
func (s *stubReportRepo) Delete(_ context.Context, _ *gorm.DB, _ string) error { return nil }
func (s *stubReportRepo) GetDashboardStats(_ context.Context) (*dto.DashboardStatsDTO, error) {
	return &dto.DashboardStatsDTO{}, nil
}
func (s *stubReportRepo) GetRecentScoredReports(_ context.Context, _ interface{}) ([]model.Report, error) {
	return nil, nil
}

// stubFileRepo implements repository.GeneratedFileRepository for export tests.
type stubFileRepo struct {
	files map[string]*model.GeneratedFile
}

func (s *stubFileRepo) Create(_ context.Context, _ *model.GeneratedFile) error       { return nil }
func (s *stubFileRepo) CreateBatch(_ context.Context, _ []model.GeneratedFile) error { return nil }
func (s *stubFileRepo) GetByReportID(_ context.Context, _ string) ([]model.GeneratedFile, error) {
	return nil, nil
}
func (s *stubFileRepo) GetByReportIDAndFilename(_ context.Context, _, filename string) (*model.GeneratedFile, error) {
	f, ok := s.files[filename]
	if !ok {
		return nil, assert.AnError
	}
	return f, nil
}

func TestExportMarkdown_ScoringReport(t *testing.T) {
	score := 85
	grade := "B"
	reportRepo := &stubReportRepo{
		report: &model.Report{
			ID:         "r1",
			Title:      "Test Report",
			ToolType:   model.ToolTypeUIReview,
			Status:     model.StatusSucceeded,
			Summary:    "All good",
			TotalScore: &score,
			Grade:      &grade,
		},
	}
	exportSvc := NewExportService(reportRepo, &stubFileRepo{files: map[string]*model.GeneratedFile{}})
	ctx := context.Background()

	content, filename, err := exportSvc.ExportMarkdown(ctx, "r1")
	require.NoError(t, err)
	assert.Contains(t, filename, "report.md")
	assert.Contains(t, string(content), "Test Report")
	assert.Contains(t, string(content), "ui_review")
	assert.Contains(t, string(content), "85/100")
	assert.Contains(t, string(content), "B")
}

func TestExportMarkdown_NonScoringReport(t *testing.T) {
	reportRepo := &stubReportRepo{
		report: &model.Report{
			ID:         "r2",
			Title:      "Agent Config Report",
			ToolType:   model.ToolTypeAgentConfig,
			Status:     model.StatusSucceeded,
			Summary:    "Generated config files",
			TotalScore: nil,
			Grade:      nil,
		},
	}
	exportSvc := NewExportService(reportRepo, &stubFileRepo{files: map[string]*model.GeneratedFile{}})
	ctx := context.Background()

	content, filename, err := exportSvc.ExportMarkdown(ctx, "r2")
	require.NoError(t, err)
	assert.Contains(t, filename, "report.md")
	assert.Contains(t, string(content), "Agent Config Report")
	assert.NotContains(t, string(content), "/100", "should not contain score for non-scoring report")
}

func TestExportMarkdown_WithGeneratedFiles(t *testing.T) {
	reportRepo := &stubReportRepo{
		report: &model.Report{
			ID:       "r4",
			Title:    "Report with files",
			ToolType: model.ToolTypeAgentConfig,
			Status:   model.StatusSucceeded,
			Summary:  "Has generated files",
			GeneratedFiles: []model.GeneratedFile{
				{Filename: "AGENTS.md", Content: "# AGENTS.md content"},
				{Filename: "TASK_PLAN.md", Content: "# Task Plan content"},
			},
		},
	}
	exportSvc := NewExportService(reportRepo, &stubFileRepo{files: map[string]*model.GeneratedFile{}})
	ctx := context.Background()

	content, _, err := exportSvc.ExportMarkdown(ctx, "r4")
	require.NoError(t, err)
	assert.Contains(t, string(content), "AGENTS.md content")
	assert.Contains(t, string(content), "Task Plan content")
}

func TestGetFileContent_Success(t *testing.T) {
	fileRepo := &stubFileRepo{
		files: map[string]*model.GeneratedFile{
			"API_DOCUMENTATION.md": {
				Filename: "API_DOCUMENTATION.md",
				Content:  "# API Docs",
				MimeType: "text/markdown",
			},
		},
	}
	exportSvc := NewExportService(&stubReportRepo{}, fileRepo)
	ctx := context.Background()

	content, filename, mimeType, err := exportSvc.GetFileContent(ctx, "r1", "API_DOCUMENTATION.md")
	require.NoError(t, err)
	assert.Equal(t, "API_DOCUMENTATION.md", filename)
	assert.Equal(t, "text/markdown", mimeType)
	assert.Equal(t, "# API Docs", string(content))
}

func TestGetFileContent_NotFound(t *testing.T) {
	exportSvc := NewExportService(
		&stubReportRepo{},
		&stubFileRepo{files: map[string]*model.GeneratedFile{}},
	)
	ctx := context.Background()

	_, _, _, err := exportSvc.GetFileContent(ctx, "r1", "nonexistent.md")
	assert.Error(t, err)
}
