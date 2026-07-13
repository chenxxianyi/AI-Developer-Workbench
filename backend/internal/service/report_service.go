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
	CreateProcessingReport(ctx context.Context, toolType, title, inputMode string, inputData json.RawMessage, parentReportID, projectID string) (*model.Report, error)
	SucceedReport(ctx context.Context, id string, reportJSON json.RawMessage, summary string, totalScore *int, grade *string, generatedFiles []model.GeneratedFile) (*dto.ReportDTO, error)
	FailReport(ctx context.Context, id string, errorMessage string) error
	FallbackReport(ctx context.Context, id string, reportJSON json.RawMessage, summary string) error
	GetReport(ctx context.Context, id string) (*dto.ReportDTO, error)
	ListReports(ctx context.Context, query dto.ListReportsQuery) (*dto.PaginatedResponse[dto.ReportDTO], error)
	DeleteReport(ctx context.Context, id string) error
	GetDashboardStats(ctx context.Context) (*dto.DashboardStatsDTO, error)
	// ValidateParentReport checks that parentReportID exists and has the same
	// toolType. Empty parentReportID is a no-op (returns nil). Returns the
	// parent report on success so callers can prefill inputs from it.
	ValidateParentReport(ctx context.Context, toolType, parentReportID string) (*model.Report, error)
	// ResolveProject returns a selected project profile for trusted prompt
	// context. An empty ID is intentionally valid and returns nil.
	ResolveProject(ctx context.Context, projectID string) (*model.Project, error)
	// CompareReports computes a delta between two same-tool reports.
	CompareReports(ctx context.Context, baselineID, targetID string) (*dto.ReportCompareDTO, error)
}

type reportService struct {
	cfg         *config.Config
	reportRepo  repository.ReportRepository
	fileRepo    repository.GeneratedFileRepository
	assetRepo   repository.ReportAssetRepository
	projectRepo repository.ProjectRepository
	db          *gorm.DB
}

// NewReportService creates a new report service.
func NewReportService(
	cfg *config.Config,
	reportRepo repository.ReportRepository,
	fileRepo repository.GeneratedFileRepository,
	assetRepo repository.ReportAssetRepository,
	projectRepo repository.ProjectRepository,
	db *gorm.DB,
) ReportService {
	return &reportService{
		cfg:         cfg,
		reportRepo:  reportRepo,
		fileRepo:    fileRepo,
		assetRepo:   assetRepo,
		projectRepo: projectRepo,
		db:          db,
	}
}

// CreateProcessingReport creates a new report with processing status.
func (s *reportService) CreateProcessingReport(ctx context.Context, toolType, title, inputMode string, inputData json.RawMessage, parentReportID, projectID string) (*model.Report, error) {
	if _, err := s.ResolveProject(ctx, projectID); err != nil {
		return nil, err
	}
	report := &model.Report{
		ToolType:   toolType,
		Title:      title,
		InputMode:  inputMode,
		Status:     model.StatusProcessing,
		InputJSON:  datatypes.JSON(inputData),
		ReportJSON: datatypes.JSON([]byte("{}")),
	}
	if parentReportID != "" {
		report.ParentReportID = &parentReportID
	}
	if projectID != "" {
		report.ProjectID = &projectID
	}
	if err := s.reportRepo.Create(ctx, report); err != nil {
		return nil, fmt.Errorf("failed to create processing report: %w", err)
	}
	return report, nil
}

// ResolveProject validates a selected project before a tool starts report
// creation. Empty IDs preserve the existing standalone workflow.
func (s *reportService) ResolveProject(ctx context.Context, projectID string) (*model.Project, error) {
	if projectID == "" {
		return nil, nil
	}
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}
	return project, nil
}

// ValidateParentReport checks that parentReportID exists and shares the same
// toolType. An empty parentReportID is treated as "no parent" (no-op).
func (s *reportService) ValidateParentReport(ctx context.Context, toolType, parentReportID string) (*model.Report, error) {
	if parentReportID == "" {
		return nil, nil
	}
	parent, err := s.reportRepo.GetByID(ctx, parentReportID)
	if err != nil {
		return nil, fmt.Errorf("parent report not found: %w", err)
	}
	if parent.ToolType != toolType {
		return nil, fmt.Errorf("parent report tool_type mismatch: expected %s, got %s", toolType, parent.ToolType)
	}
	return parent, nil
}

// CompareReports computes a delta between two same-tool reports.
func (s *reportService) CompareReports(ctx context.Context, baselineID, targetID string) (*dto.ReportCompareDTO, error) {
	baseline, err := s.reportRepo.GetByID(ctx, baselineID)
	if err != nil {
		return nil, fmt.Errorf("baseline report not found: %w", err)
	}
	target, err := s.reportRepo.GetByID(ctx, targetID)
	if err != nil {
		return nil, fmt.Errorf("target report not found: %w", err)
	}
	if baseline.ToolType != target.ToolType {
		return nil, fmt.Errorf("cannot compare reports of different tool types: %s vs %s", baseline.ToolType, target.ToolType)
	}

	result := &dto.ReportCompareDTO{
		ToolType: baseline.ToolType,
		BaselineReport: dto.ReportSummaryDTO{
			ID:         baseline.ID,
			Title:      baseline.Title,
			Status:     baseline.Status,
			TotalScore: baseline.TotalScore,
			Grade:      baseline.Grade,
			CreatedAt:  baseline.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			Summary:    baseline.Summary,
			ReportData: json.RawMessage(baseline.ReportJSON),
		},
		TargetReport: dto.ReportSummaryDTO{
			ID:         target.ID,
			Title:      target.Title,
			Status:     target.Status,
			TotalScore: target.TotalScore,
			Grade:      target.Grade,
			CreatedAt:  target.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			Summary:    target.Summary,
			ReportData: json.RawMessage(target.ReportJSON),
		},
	}

	// Score delta: only when both reports have scores. Null is "not comparable".
	if baseline.TotalScore != nil && target.TotalScore != nil {
		delta := *target.TotalScore - *baseline.TotalScore
		result.ScoreDelta = &delta
	} else {
		result.Warnings = append(result.Warnings, "one or both reports have no score; score delta omitted")
	}

	if baseline.Grade != nil && target.Grade != nil {
		result.GradeDelta = gradeDelta(*baseline.Grade, *target.Grade)
	}

	// Parse issues + action items from report_data.
	baseIssues := extractIssues(json.RawMessage(baseline.ReportJSON))
	tgtIssues := extractIssues(json.RawMessage(target.ReportJSON))
	baseActions := extractActionItems(json.RawMessage(baseline.ReportJSON))
	tgtActions := extractActionItems(json.RawMessage(target.ReportJSON))

	// Issue counts by severity.
	result.IssueCountDelta = IssueCountDeltaFrom(baseIssues, tgtIssues)

	// Match issues by stable id (ActionItem has id, IssueItem does not), fallback category+title.
	result.Issues = compareIssues(baseIssues, tgtIssues)
	result.ActionItems = compareActionItems(baseActions, tgtActions)

	return result, nil
}

// SucceedReport updates a report to succeeded status, saves generated files in a transaction.
func (s *reportService) SucceedReport(ctx context.Context, id string, reportJSON json.RawMessage, summary string, totalScore *int, grade *string, generatedFiles []model.GeneratedFile) (*dto.ReportDTO, error) {
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Lock and read the report within the transaction.
		var report model.Report
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&report, "id = ?", id).Error; err != nil {
			return fmt.Errorf("report not found: %w", err)
		}

		report.Status = model.StatusSucceeded
		report.ReportJSON = datatypes.JSON(reportJSON)
		report.Summary = summary
		report.TotalScore = totalScore
		report.Grade = grade

		if err := tx.Save(&report).Error; err != nil {
			return fmt.Errorf("failed to update report: %w", err)
		}

		// Save generated files within the same transaction.
		for i := range generatedFiles {
			generatedFiles[i].ReportID = id
			generatedFiles[i].SizeBytes = uint64(len(generatedFiles[i].Content))
		}
		if len(generatedFiles) > 0 {
			if err := tx.Create(&generatedFiles).Error; err != nil {
				return fmt.Errorf("failed to save generated files: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
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
		ID:             report.ID,
		ToolType:       report.ToolType,
		Title:          report.Title,
		InputMode:      report.InputMode,
		Status:         report.Status,
		Summary:        report.Summary,
		TotalScore:     report.TotalScore,
		Grade:          report.Grade,
		InputData:      json.RawMessage(report.InputJSON),
		ReportData:     json.RawMessage(report.ReportJSON),
		ParentReportID: report.ParentReportID,
		ProjectID:      report.ProjectID,
		CreatedAt:      report.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      report.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
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
			ParentReportID: r.ParentReportID,
			ProjectID:      r.ProjectID,
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
