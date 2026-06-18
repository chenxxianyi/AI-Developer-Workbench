package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"ai-developer-workbench/internal/config"
	"ai-developer-workbench/internal/dto"
	"ai-developer-workbench/internal/model"
	"ai-developer-workbench/internal/repository"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ReportService handles report lifecycle operations.
type ReportService interface {
	CreateProcessingReport(ctx context.Context, toolType, title, inputMode string, inputData json.RawMessage) (*model.Report, error)
	SucceedReport(ctx context.Context, id string, reportJSON json.RawMessage, summary string, totalScore *int, grade *string, generatedFiles []model.GeneratedFile) (*dto.ReportDTO, error)
	FailReport(ctx context.Context, id string, errorMessage string) error
	FallbackReport(ctx context.Context, id string, reportJSON json.RawMessage, summary string) error
	GetReport(ctx context.Context, id string) (*dto.ReportDTO, error)
	ListReports(ctx context.Context, query dto.ListReportsQuery) (*dto.PaginatedResponse[dto.ReportDTO], error)
	DeleteReport(ctx context.Context, id string) error
	GetDashboardStats(ctx context.Context) (*dto.DashboardStatsDTO, error)
}

type reportService struct {
	cfg          *config.Config
	reportRepo   repository.ReportRepository
	fileRepo     repository.GeneratedFileRepository
	assetRepo    repository.ReportAssetRepository
	db           *gorm.DB
}

// NewReportService creates a new report service.
func NewReportService(
	cfg *config.Config,
	reportRepo repository.ReportRepository,
	fileRepo repository.GeneratedFileRepository,
	assetRepo repository.ReportAssetRepository,
	db *gorm.DB,
) ReportService {
	return &reportService{
		cfg:        cfg,
		reportRepo: reportRepo,
		fileRepo:   fileRepo,
		assetRepo:  assetRepo,
		db:         db,
	}
}

// CreateProcessingReport creates a new report with processing status.
func (s *reportService) CreateProcessingReport(ctx context.Context, toolType, title, inputMode string, inputData json.RawMessage) (*model.Report, error) {
	report := &model.Report{
		ToolType:   toolType,
		Title:      title,
		InputMode:  inputMode,
		Status:     model.StatusProcessing,
		InputJSON:  datatypes.JSON(inputData),
		ReportJSON: datatypes.JSON([]byte("{}")),
	}
	if err := s.reportRepo.Create(ctx, report); err != nil {
		return nil, fmt.Errorf("failed to create processing report: %w", err)
	}
	return report, nil
}

// SucceedReport updates a report to succeeded status with results.
func (s *reportService) SucceedReport(ctx context.Context, id string, reportJSON json.RawMessage, summary string, totalScore *int, grade *string, generatedFiles []model.GeneratedFile) (*dto.ReportDTO, error) {
	report, err := s.reportRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("report not found: %w", err)
	}

	report.Status = model.StatusSucceeded
	report.ReportJSON = datatypes.JSON(reportJSON)
	report.Summary = summary
	report.TotalScore = totalScore
	report.Grade = grade

	if err := s.reportRepo.Update(ctx, report); err != nil {
		return nil, fmt.Errorf("failed to update report: %w", err)
	}

	// Save generated files.
	for i := range generatedFiles {
		generatedFiles[i].ReportID = id
		generatedFiles[i].SizeBytes = uint64(len(generatedFiles[i].Content))
	}
	if len(generatedFiles) > 0 {
		if err := s.fileRepo.CreateBatch(ctx, generatedFiles); err != nil {
			slog.Warn("Failed to save generated files", "report_id", id, "error", err)
		}
	}

	return s.GetReport(ctx, id)
}

// FailReport updates a report to failed status with error message.
func (s *reportService) FailReport(ctx context.Context, id string, errorMessage string) error {
	report, err := s.reportRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("report not found: %w", err)
	}

	report.Status = model.StatusFailed
	report.ErrorMessage = errorMessage

	return s.reportRepo.Update(ctx, report)
}

// FallbackReport updates a report to fallback status with partial results.
func (s *reportService) FallbackReport(ctx context.Context, id string, reportJSON json.RawMessage, summary string) error {
	report, err := s.reportRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("report not found: %w", err)
	}

	report.Status = model.StatusFallback
	report.ReportJSON = datatypes.JSON(reportJSON)
	report.Summary = summary

	return s.reportRepo.Update(ctx, report)
}

// GetReport retrieves a report by ID and converts to DTO.
func (s *reportService) GetReport(ctx context.Context, id string) (*dto.ReportDTO, error) {
	report, err := s.reportRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Convert model to DTO.
	dtoReport := &dto.ReportDTO{
		ID:         report.ID,
		ToolType:   report.ToolType,
		Title:      report.Title,
		InputMode:  report.InputMode,
		Status:     report.Status,
		Summary:    report.Summary,
		TotalScore: report.TotalScore,
		Grade:      report.Grade,
		InputData:  json.RawMessage(report.InputJSON),
		ReportData: json.RawMessage(report.ReportJSON),
		CreatedAt:  report.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  report.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Convert generated files.
	dtoReport.GeneratedFiles = make([]dto.GeneratedFileDTO, 0, len(report.GeneratedFiles))
	for _, f := range report.GeneratedFiles {
		dtoReport.GeneratedFiles = append(dtoReport.GeneratedFiles, dto.GeneratedFileDTO{
			ID:        f.ID,
			Filename:  f.Filename,
			Language:  f.Language,
			MimeType:  f.MimeType,
			SizeBytes: f.SizeBytes,
		})
	}

	return dtoReport, nil
}

// ListReports retrieves a paginated list of reports.
func (s *reportService) ListReports(ctx context.Context, query dto.ListReportsQuery) (*dto.PaginatedResponse[dto.ReportDTO], error) {
	query.SetDefaults()

	reports, total, err := s.reportRepo.List(ctx, query)
	if err != nil {
		return nil, err
	}

	items := make([]dto.ReportDTO, 0, len(reports))
	for _, r := range reports {
		items = append(items, dto.ReportDTO{
			ID:             r.ID,
			ToolType:       r.ToolType,
			Title:          r.Title,
			InputMode:      r.InputMode,
			Status:         r.Status,
			Summary:        r.Summary,
			TotalScore:     r.TotalScore,
			Grade:          r.Grade,
			InputData:      json.RawMessage([]byte("{}")),
			ReportData:     json.RawMessage([]byte("{}")),
			GeneratedFiles: []dto.GeneratedFileDTO{},
			CreatedAt:      r.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:      r.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return &dto.PaginatedResponse[dto.ReportDTO]{
		Items:    items,
		Total:    total,
		Page:     query.Page,
		PageSize: query.PageSize,
	}, nil
}

// DeleteReport deletes a report and its associated files.
// Database deletion is done in a transaction, then disk cleanup is attempted.
func (s *reportService) DeleteReport(ctx context.Context, id string) error {
	// Get relative paths for disk cleanup before deletion.
	relativePaths, err := s.assetRepo.GetRelativePathsByReportID(ctx, id)
	if err != nil {
		slog.Warn("Failed to get asset paths for cleanup", "report_id", id, "error", err)
	}

	// Delete in transaction.
	err = s.db.Transaction(func(tx *gorm.DB) error {
		return s.reportRepo.Delete(ctx, tx, id)
	})
	if err != nil {
		return fmt.Errorf("failed to delete report: %w", err)
	}

	// Disk cleanup (non-transactional, failures are logged).
	s.cleanupDiskFiles(id, relativePaths)

	return nil
}

// cleanupDiskFiles removes uploaded files and directories for a deleted report.
func (s *reportService) cleanupDiskFiles(reportID string, relativePaths []string) {
	// Clean up individual asset files.
	for _, relPath := range relativePaths {
		fullPath := filepath.Join(s.cfg.Upload.Dir, relPath)
		if err := os.Remove(fullPath); err != nil {
			slog.Warn("Failed to remove asset file", "path", fullPath, "error", err)
		}
	}

	// Clean up report directories.
	dirsToClean := []string{
		filepath.Join(s.cfg.Upload.Dir, reportID),
		filepath.Join(s.cfg.Upload.TempDir, reportID),
	}
	for _, dir := range dirsToClean {
		if err := os.RemoveAll(dir); err != nil {
			slog.Warn("Failed to remove report directory", "dir", dir, "error", err)
		}
	}
}

// GetDashboardStats retrieves dashboard statistics.
func (s *reportService) GetDashboardStats(ctx context.Context) (*dto.DashboardStatsDTO, error) {
	return s.reportRepo.GetDashboardStats(ctx)
}