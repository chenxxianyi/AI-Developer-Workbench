package repository

import (
	"context"

	"ai-developer-workbench/internal/model"

	"gorm.io/gorm"
)

// ReportAssetRepository handles database access for report assets.
type ReportAssetRepository interface {
	Create(ctx context.Context, asset *model.ReportAsset) error
	GetByReportID(ctx context.Context, reportID string) ([]model.ReportAsset, error)
	GetRelativePathsByReportID(ctx context.Context, reportID string) ([]string, error)
}

type reportAssetRepository struct {
	db *gorm.DB
}

// NewReportAssetRepository creates a new report asset repository.
func NewReportAssetRepository(db *gorm.DB) ReportAssetRepository {
	return &reportAssetRepository{db: db}
}

// Create inserts a new report asset.
func (r *reportAssetRepository) Create(ctx context.Context, asset *model.ReportAsset) error {
	return r.db.WithContext(ctx).Create(asset).Error
}

// GetByReportID retrieves all assets for a report.
func (r *reportAssetRepository) GetByReportID(ctx context.Context, reportID string) ([]model.ReportAsset, error) {
	var assets []model.ReportAsset
	err := r.db.WithContext(ctx).Where("report_id = ?", reportID).Find(&assets).Error
	return assets, err
}

// GetRelativePathsByReportID retrieves the relative paths of all assets for a report.
// Used for disk cleanup during report deletion.
func (r *reportAssetRepository) GetRelativePathsByReportID(ctx context.Context, reportID string) ([]string, error) {
	var paths []string
	err := r.db.WithContext(ctx).
		Model(&model.ReportAsset{}).
		Where("report_id = ?", reportID).
		Pluck("relative_path", &paths).Error
	return paths, err
}
