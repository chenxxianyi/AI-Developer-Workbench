package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"ai-developer-workbench/internal/dto"
	"ai-developer-workbench/internal/model"

	"gorm.io/gorm"
)

// ProjectRepository handles database access for projects.
type ProjectRepository interface {
	Create(ctx context.Context, project *model.Project) error
	GetByID(ctx context.Context, id string) (*model.Project, error)
	List(ctx context.Context, query dto.ListProjectsQuery) ([]ProjectListItem, int64, error)
	Update(ctx context.Context, project *model.Project) error
	Delete(ctx context.Context, id string) (int64, error)
	GetStats(ctx context.Context, id string) (*dto.ProjectStatsDTO, error)
	ListReports(ctx context.Context, id string, query dto.ListReportsQuery) ([]model.Report, int64, error)
}

// ProjectListItem is a project profile enriched with list-only report aggregates.
type ProjectListItem struct {
	Project      model.Project
	ReportCount  int64
	AverageScore *float64
}

type projectRepository struct {
	db *gorm.DB
}

// NewProjectRepository creates a new project repository.
func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) Create(ctx context.Context, project *model.Project) error {
	return r.db.WithContext(ctx).Create(project).Error
}

func (r *projectRepository) GetByID(ctx context.Context, id string) (*model.Project, error) {
	var project model.Project
	if err := r.db.WithContext(ctx).First(&project, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}
	return &project, nil
}

func (r *projectRepository) List(ctx context.Context, query dto.ListProjectsQuery) ([]ProjectListItem, int64, error) {
	base := r.db.WithContext(ctx).Model(&model.Project{})
	if query.Search != "" {
		like := "%" + query.Search + "%"
		base = base.Where("name LIKE ? OR description LIKE ?", like, like)
	}
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count projects: %w", err)
	}

	reportStats := r.db.WithContext(ctx).Model(&model.Report{}).
		Select("project_id, COUNT(*) AS report_count, AVG(total_score) AS average_score").
		Where("project_id IS NOT NULL").
		Group("project_id")

	type projectListRow struct {
		ID            string
		Name          string
		Description   string
		RepoURL       string
		FrontendStack string
		BackendStack  string
		Database      string
		UIStyle       string
		CodingRules   string
		CreatedAt     time.Time
		UpdatedAt     time.Time
		ReportCount   int64
		AverageScore  *float64
	}
	var rows []projectListRow
	offset := (query.Page - 1) * query.PageSize
	q := r.db.WithContext(ctx).Table("projects AS p").
		Select("p.*, COALESCE(rs.report_count, 0) AS report_count, rs.average_score").
		Joins("LEFT JOIN (?) AS rs ON rs.project_id = p.id", reportStats)
	if query.Search != "" {
		like := "%" + query.Search + "%"
		q = q.Where("p.name LIKE ? OR p.description LIKE ?", like, like)
	}
	if err := q.Order("p.updated_at DESC").Offset(offset).Limit(query.PageSize).Scan(&rows).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list projects: %w", err)
	}

	items := make([]ProjectListItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, ProjectListItem{
			Project: model.Project{
				ID:            row.ID,
				Name:          row.Name,
				Description:   row.Description,
				RepoURL:       row.RepoURL,
				FrontendStack: row.FrontendStack,
				BackendStack:  row.BackendStack,
				Database:      row.Database,
				UIStyle:       row.UIStyle,
				CodingRules:   row.CodingRules,
				CreatedAt:     row.CreatedAt,
				UpdatedAt:     row.UpdatedAt,
			},
			ReportCount:  row.ReportCount,
			AverageScore: row.AverageScore,
		})
	}
	return items, total, nil
}

func (r *projectRepository) Update(ctx context.Context, project *model.Project) error {
	return r.db.WithContext(ctx).Save(project).Error
}

func (r *projectRepository) Delete(ctx context.Context, id string) (int64, error) {
	var detached int64
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Report{}).Where("project_id = ?", id).Count(&detached).Error; err != nil {
			return err
		}
		// Clear report ownership explicitly as well as retaining the database
		// FK's ON DELETE SET NULL behavior. This keeps SQLite test databases and
		// MySQL deployments aligned.
		if err := tx.Model(&model.Report{}).Where("project_id = ?", id).Update("project_id", nil).Error; err != nil {
			return err
		}
		result := tx.Delete(&model.Project{}, "id = ?", id)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return detached, nil
}

func (r *projectRepository) GetStats(ctx context.Context, id string) (*dto.ProjectStatsDTO, error) {
	stats := &dto.ProjectStatsDTO{
		ToolUsage:       make(map[string]int64),
		RecentReports:   []dto.RecentReportDTO{},
		QualityTrend:    []dto.QualityTrendPointDTO{},
		LatestArtifacts: []dto.ProjectArtifactDTO{},
	}

	var reports []model.Report
	if err := r.db.WithContext(ctx).
		Select("id, tool_type, title, status, summary, total_score, grade, report_json, created_at").
		Where("project_id = ?", id).
		Order("created_at DESC").
		Find(&reports).Error; err != nil {
		return nil, fmt.Errorf("failed to get project reports: %w", err)
	}

	for _, t := range model.ValidToolTypes() {
		stats.ToolUsage[t] = 0
	}
	var scoreSum float64
	var scoreCount int64
	for i, rp := range reports {
		stats.TotalReports++
		stats.ToolUsage[rp.ToolType]++
		if rp.TotalScore != nil {
			scoreSum += float64(*rp.TotalScore)
			scoreCount++
		}
		high := countHighSeverityInJSON(rp.ReportJSON)
		stats.HighSeverityCount += high
		if i < 5 {
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
	}
	if scoreCount > 0 {
		average := scoreSum / float64(scoreCount)
		stats.AverageScore = &average
	}
	stats.QualityTrend = buildProjectTrend(reports)

	var artifactReports []model.Report
	artifactTools := []string{model.ToolTypeAgentConfig, model.ToolTypeAPIDoc, model.ToolTypeDBSchema}
	if err := r.db.WithContext(ctx).
		Preload("GeneratedFiles").
		Where("project_id = ? AND tool_type IN ? AND status IN ?", id, artifactTools, []string{model.StatusSucceeded, model.StatusFallback}).
		Order("created_at DESC").
		Find(&artifactReports).Error; err != nil {
		return nil, fmt.Errorf("failed to get latest project artifacts: %w", err)
	}
	seenTool := make(map[string]bool, len(artifactTools))
	for _, report := range artifactReports {
		if seenTool[report.ToolType] || len(report.GeneratedFiles) == 0 {
			continue
		}
		seenTool[report.ToolType] = true
		for _, file := range report.GeneratedFiles {
			stats.LatestArtifacts = append(stats.LatestArtifacts, dto.ProjectArtifactDTO{
				ToolType:    report.ToolType,
				ReportID:    report.ID,
				ReportTitle: report.Title,
				Filename:    file.Filename,
				MimeType:    file.MimeType,
				CreatedAt:   report.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			})
		}
	}

	return stats, nil
}

// ListReports returns a project's report history without leaking other projects'
// report data.
func (r *projectRepository) ListReports(ctx context.Context, id string, query dto.ListReportsQuery) ([]model.Report, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.Report{}).Where("project_id = ?", id)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count project reports: %w", err)
	}
	var reports []model.Report
	offset := (query.Page - 1) * query.PageSize
	if err := q.Order("created_at DESC").Offset(offset).Limit(query.PageSize).Find(&reports).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list project reports: %w", err)
	}
	return reports, total, nil
}

func countHighSeverityInJSON(raw []byte) int64 {
	var payload struct {
		Issues []dto.IssueItem `json:"issues"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return 0
	}
	var count int64
	for _, issue := range payload.Issues {
		if issue.Severity == "high" {
			count++
		}
	}
	return count
}

func buildProjectTrend(reports []model.Report) []dto.QualityTrendPointDTO {
	cutoff := time.Now().UTC().AddDate(0, 0, -30)
	type bucketTotals struct {
		reportCount int64
		scoreSum    float64
		scoreCount  int64
		highRisk    int64
	}
	buckets := make(map[string]*bucketTotals)
	for _, report := range reports {
		if report.CreatedAt.Before(cutoff) {
			continue
		}
		key := report.CreatedAt.UTC().Format("2006-01-02")
		bucket := buckets[key]
		if bucket == nil {
			bucket = &bucketTotals{}
			buckets[key] = bucket
		}
		bucket.reportCount++
		bucket.highRisk += countHighSeverityInJSON(report.ReportJSON)
		if report.TotalScore != nil {
			bucket.scoreSum += float64(*report.TotalScore)
			bucket.scoreCount++
		}
	}
	keys := make([]string, 0, len(buckets))
	for key := range buckets {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	result := make([]dto.QualityTrendPointDTO, 0, len(keys))
	for _, key := range keys {
		bucket := buckets[key]
		point := dto.QualityTrendPointDTO{
			Bucket:            key,
			ReportCount:       bucket.reportCount,
			HighSeverityCount: bucket.highRisk,
		}
		if bucket.scoreCount > 0 {
			average := bucket.scoreSum / float64(bucket.scoreCount)
			point.AverageScore = &average
		}
		result = append(result, point)
	}
	return result
}
