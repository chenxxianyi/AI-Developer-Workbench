package repository

import (
	"context"
	"fmt"

	"ai-developer-workbench/internal/dto"
	"ai-developer-workbench/internal/model"

	"gorm.io/gorm"
)

// ReportRepository handles database access for reports.
type ReportRepository interface {
	Create(ctx context.Context, report *model.Report) error
	Update(ctx context.Context, report *model.Report) error
	GetByID(ctx context.Context, id string) (*model.Report, error)
	List(ctx context.Context, query dto.ListReportsQuery) ([]model.Report, int64, error)
	Delete(ctx context.Context, tx *gorm.DB, id string) error
	GetDashboardStats(ctx context.Context) (*dto.DashboardStatsDTO, error)
}

type reportRepository struct {
	db *gorm.DB
}

// NewReportRepository creates a new report repository.
func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db: db}
}

// Create inserts a new report.
func (r *reportRepository) Create(ctx context.Context, report *model.Report) error {
	return r.db.WithContext(ctx).Create(report).Error
}

// Update updates an existing report.
func (r *reportRepository) Update(ctx context.Context, report *model.Report) error {
	return r.db.WithContext(ctx).Save(report).Error
}

// GetByID retrieves a report by ID with preloaded relations.
func (r *reportRepository) GetByID(ctx context.Context, id string) (*model.Report, error) {
	var report model.Report
	err := r.db.WithContext(ctx).
		Preload("GeneratedFiles").
		Preload("Assets").
		First(&report, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

// List retrieves a paginated list of reports with filtering and sorting.
func (r *reportRepository) List(ctx context.Context, query dto.ListReportsQuery) ([]model.Report, int64, error) {
	var reports []model.Report
	var total int64

	q := r.db.WithContext(ctx).Model(&model.Report{})

	// Apply tool_type filter with whitelist.
	if query.ToolType != "" && model.IsValidToolType(query.ToolType) {
		q = q.Where("tool_type = ?", query.ToolType)
	}

	// Apply status filter with whitelist.
	if query.Status != "" && dto.ValidStatusValues()[query.Status] {
		q = q.Where("status = ?", query.Status)
	}

	// Count total before pagination.
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sort with fixed mapping to prevent SQL injection.
	sortMapping := dto.ValidSortValues()
	if orderBy, ok := sortMapping[query.Sort]; ok {
		q = q.Order(orderBy)
	} else {
		q = q.Order("created_at DESC")
	}

	// Apply pagination.
	offset := (query.Page - 1) * query.PageSize
	if err := q.Offset(offset).Limit(query.PageSize).Find(&reports).Error; err != nil {
		return nil, 0, err
	}

	return reports, total, nil
}

// Delete removes a report by ID. Must be called within a transaction.
func (r *reportRepository) Delete(ctx context.Context, tx *gorm.DB, id string) error {
	db := r.getDB(tx)
	// CASCADE will handle generated_files and report_assets.
	return db.WithContext(ctx).Delete(&model.Report{}, "id = ?", id).Error
}

// GetDashboardStats retrieves aggregated dashboard statistics.
func (r *reportRepository) GetDashboardStats(ctx context.Context) (*dto.DashboardStatsDTO, error) {
	stats := &dto.DashboardStatsDTO{
		ToolUsage: make(map[string]int64),
	}

	// Total reports.
	if err := r.db.WithContext(ctx).Model(&model.Report{}).Count(&stats.TotalReports).Error; err != nil {
		return nil, fmt.Errorf("failed to count reports: %w", err)
	}

	// Per-tool usage counts.
	type toolCount struct {
		ToolType string
		Count    int64
	}
	var toolCounts []toolCount
	if err := r.db.WithContext(ctx).Model(&model.Report{}).
		Select("tool_type, COUNT(*) as count").
		Group("tool_type").
		Find(&toolCounts).Error; err != nil {
		return nil, fmt.Errorf("failed to count tool usage: %w", err)
	}
	// Ensure all tool types have entries.
	for _, t := range model.ValidToolTypes() {
		stats.ToolUsage[t] = 0
	}
	for _, tc := range toolCounts {
		stats.ToolUsage[tc.ToolType] = tc.Count
	}

	// Average score (excluding nulls).
	var avgResult struct {
		AvgScore *float64
	}
	if err := r.db.WithContext(ctx).Model(&model.Report{}).
		Select("AVG(total_score) as avg_score").
		Where("total_score IS NOT NULL").
		Scan(&avgResult).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate average score: %w", err)
	}
	stats.AverageScore = avgResult.AvgScore

	// Recent 5 reports.
	var recent []model.Report
	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(5).
		Find(&recent).Error; err != nil {
		return nil, fmt.Errorf("failed to get recent reports: %w", err)
	}
	stats.RecentReports = make([]dto.RecentReportDTO, 0, len(recent))
	for _, rp := range recent {
		stats.RecentReports = append(stats.RecentReports, dto.RecentReportDTO{
			ID:         rp.ID,
			ToolType:   rp.ToolType,
			Title:      rp.Title,
			Status:     rp.Status,
			TotalScore: rp.TotalScore,
			Grade:      rp.Grade,
			Summary:    rp.Summary,
			CreatedAt:  rp.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return stats, nil
}

// getDB returns the transaction DB if provided, otherwise the default DB.
func (r *reportRepository) getDB(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}
