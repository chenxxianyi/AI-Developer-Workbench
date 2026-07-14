package repository

import (
	"context"

	"ai-developer-workbench/internal/model"

	"gorm.io/gorm"
)

type AIRunRepository interface {
	Create(ctx context.Context, run *model.AIRun) error
	GetStats(ctx context.Context, toolType string, days int) (*AIObservabilityStats, error)
}

type AIObservabilityStats struct {
	TotalCalls    int64          `json:"total_calls"`
	SuccessRate   float64        `json:"success_rate"`
	ParseFailRate float64        `json:"parse_fail_rate"`
	AvgDurationMs float64        `json:"avg_duration_ms"`
	P50DurationMs float64        `json:"p50_duration_ms"`
	P95DurationMs float64        `json:"p95_duration_ms"`
	RetryRate     float64        `json:"retry_rate"`
	ByTool        []ToolAIStats  `json:"by_tool"`
	ByModel       []ModelAIStats `json:"by_model"`
}

type ToolAIStats struct {
	ToolType    string  `json:"tool_type"`
	TotalCalls  int64   `json:"total_calls"`
	SuccessRate float64 `json:"success_rate"`
}

type ModelAIStats struct {
	Model         string  `json:"model"`
	TotalCalls    int64   `json:"total_calls"`
	AvgDurationMs float64 `json:"avg_duration_ms"`
}

type aiRunRepository struct {
	db *gorm.DB
}

func NewAIRunRepository(db *gorm.DB) AIRunRepository {
	return &aiRunRepository{db: db}
}

func (r *aiRunRepository) Create(ctx context.Context, run *model.AIRun) error {
	return r.db.WithContext(ctx).Create(run).Error
}

func (r *aiRunRepository) GetStats(ctx context.Context, toolType string, days int) (*AIObservabilityStats, error) {
	since := r.db.NowFunc().AddDate(0, 0, -days)
	query := r.db.WithContext(ctx).Model(&model.AIRun{}).Where("created_at >= ?", since)
	if toolType != "" {
		query = query.Where("tool_type = ?", toolType)
	}

	type row struct {
		Count       int64
		ParseOk     int64
		Retried     int64
		SumDuration int64
	}

	var total row
	err := query.Select(
		"COUNT(*) as count",
		"SUM(CASE WHEN parse_success = 1 THEN 1 ELSE 0 END) as parse_ok",
		"SUM(CASE WHEN retry_count > 0 THEN 1 ELSE 0 END) as retried",
		"COALESCE(SUM(duration_ms), 0) as sum_duration",
	).Scan(&total).Error
	if err != nil {
		return nil, err
	}

	stats := &AIObservabilityStats{TotalCalls: total.Count}
	if total.Count > 0 {
		stats.ParseFailRate = float64(total.Count-total.ParseOk) / float64(total.Count) * 100
		stats.RetryRate = float64(total.Retried) / float64(total.Count) * 100
		stats.AvgDurationMs = float64(total.SumDuration) / float64(total.Count)
	}

	// By tool
	var byTool []ToolAIStats
	err = r.db.WithContext(ctx).Model(&model.AIRun{}).Where("created_at >= ?", since).
		Select("tool_type, COUNT(*) as total_calls, COALESCE(SUM(CASE WHEN parse_success = 1 THEN 1 ELSE 0 END),0)*100.0/COUNT(*) as success_rate").
		Group("tool_type").Order("total_calls DESC").Find(&byTool).Error
	if err == nil {
		stats.ByTool = byTool
	}

	// By model
	var byModel []ModelAIStats
	err = r.db.WithContext(ctx).Model(&model.AIRun{}).Where("created_at >= ?", since).
		Select("model, COUNT(*) as total_calls, COALESCE(AVG(duration_ms),0) as avg_duration_ms").
		Group("model").Order("total_calls DESC").Find(&byModel).Error
	if err == nil {
		stats.ByModel = byModel
	}

	return stats, nil
}
