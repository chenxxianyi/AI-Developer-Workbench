package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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
	// GetRecentScoredReports returns reports created since `since` that have a
	// non-null total_score, ordered by created_at asc. Used by the service to
	// compute high-severity issue counts and the weekly summary.
	GetRecentScoredReports(ctx context.Context, since interface{}) ([]model.Report, error)
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

// GetDashboardStats retrieves aggregated dashboard statistics, including a
// 30-day quality trend (daily buckets) and a weekly summary.
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

	// 30-day quality trend (daily buckets, UTC). We fetch all reports in the
	// window and aggregate in Go because high-severity issue counts live
	// inside report_json, which is not cleanly aggregatable in SQL.
	now := time.Now().UTC()
	since30 := now.AddDate(0, 0, -30)
	windowReports, err := r.GetRecentScoredReports(ctx, since30)
	if err != nil {
		return nil, err
	}
	stats.QualityTrend = buildQualityTrend(windowReports, since30, now)

	// Weekly summary (last 7 days, UTC). Reuses the 30-day fetch to avoid a
	// second query; we filter to the last 7 days in Go.
	since7 := now.AddDate(0, 0, -7)
	weeklyReports := filterReportsSince(windowReports, since7)
	stats.WeeklyStats = buildWeeklyStats(weeklyReports)

	return stats, nil
}

// buildQualityTrend aggregates reports into daily buckets [since, now].
func buildQualityTrend(reports []model.Report, since, now time.Time) []dto.QualityTrendPointDTO {
	// Bucket by YYYY-MM-DD (UTC day start).
	buckets := make(map[string][]model.Report)
	for _, rp := range reports {
		day := rp.CreatedAt.UTC().Format("2006-01-02")
		buckets[day] = append(buckets[day], rp)
	}

	// Build a contiguous list of days from `since` to `now` (inclusive).
	var trend []dto.QualityTrendPointDTO
	for d := truncateToDay(since); !d.After(truncateToDay(now)); d = d.AddDate(0, 0, 1) {
		key := d.Format("2006-01-02")
		dayReports := buckets[key]
		point := dto.QualityTrendPointDTO{
			Bucket:       key,
			ReportCount:  int64(len(dayReports)),
			AverageScore: nil,
		}
		var sum int
		count := 0
		for _, rp := range dayReports {
			if rp.TotalScore != nil {
				sum += *rp.TotalScore
				count++
			}
		}
		if count > 0 {
			avg := float64(sum) / float64(count)
			point.AverageScore = &avg
		}
		point.HighSeverityCount = countHighSeverity(dayReports)
		trend = append(trend, point)
	}
	return trend
}

// buildWeeklyStats computes the 7-day summary from pre-fetched reports.
func buildWeeklyStats(reports []model.Report) *dto.WeeklyStatsDTO {
	w := &dto.WeeklyStatsDTO{
		ReportCountThisWeek: int64(len(reports)),
	}
	if len(reports) == 0 {
		return w
	}
	var sum int
	scoredCount := 0
	toolCounts := make(map[string]int64)
	for _, rp := range reports {
		if rp.TotalScore != nil {
			sum += *rp.TotalScore
			scoredCount++
		}
		toolCounts[rp.ToolType]++
	}
	if scoredCount > 0 {
		avg := float64(sum) / float64(scoredCount)
		w.AverageScoreThisWeek = &avg
	}
	w.HighSeverityCountThisWeek = countHighSeverity(reports)
	// Most used tool this week.
	var bestTool string
	var bestCount int64
	for tool, c := range toolCounts {
		if c > bestCount {
			bestTool = tool
			bestCount = c
		}
	}
	w.MostUsedToolThisWeek = bestTool
	return w
}

// countHighSeverity counts issues with severity "high" across reports' ReportJSON.
func countHighSeverity(reports []model.Report) int64 {
	var total int64
	for _, rp := range reports {
		var probe struct {
			Issues []struct {
				Severity string `json:"severity"`
			} `json:"issues"`
		}
		if err := json.Unmarshal(rp.ReportJSON, &probe); err != nil {
			continue
		}
		for _, it := range probe.Issues {
			if it.Severity == "high" {
				total++
			}
		}
	}
	return total
}

// filterReportsSince returns reports with created_at >= since.
func filterReportsSince(reports []model.Report, since time.Time) []model.Report {
	var out []model.Report
	for _, rp := range reports {
		if !rp.CreatedAt.Before(since) {
			out = append(out, rp)
		}
	}
	return out
}

// truncateToDay returns t truncated to the start of its UTC day.
func truncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

// GetRecentScoredReports returns reports created since `since` that have a
// non-null total_score, ordered by created_at asc.
func (r *reportRepository) GetRecentScoredReports(ctx context.Context, since interface{}) ([]model.Report, error) {
	var reports []model.Report
	if err := r.db.WithContext(ctx).
		Where("created_at >= ? AND total_score IS NOT NULL", since).
		Order("created_at ASC").
		Find(&reports).Error; err != nil {
		return nil, fmt.Errorf("failed to get recent scored reports: %w", err)
	}
	return reports, nil
}

// getDB returns the transaction DB if provided, otherwise the default DB.
func (r *reportRepository) getDB(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}
