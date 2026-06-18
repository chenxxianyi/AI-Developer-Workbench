package repository

import (
	"context"

	"ai-developer-workbench/internal/model"

	"gorm.io/gorm"
)

// GeneratedFileRepository handles database access for generated files.
type GeneratedFileRepository interface {
	Create(ctx context.Context, file *model.GeneratedFile) error
	CreateBatch(ctx context.Context, files []model.GeneratedFile) error
	GetByReportID(ctx context.Context, reportID string) ([]model.GeneratedFile, error)
	GetByReportIDAndFilename(ctx context.Context, reportID string, filename string) (*model.GeneratedFile, error)
}

type generatedFileRepository struct {
	db *gorm.DB
}

// NewGeneratedFileRepository creates a new generated file repository.
func NewGeneratedFileRepository(db *gorm.DB) GeneratedFileRepository {
	return &generatedFileRepository{db: db}
}

// Create inserts a new generated file.
func (r *generatedFileRepository) Create(ctx context.Context, file *model.GeneratedFile) error {
	return r.db.WithContext(ctx).Create(file).Error
}

// CreateBatch inserts multiple generated files.
func (r *generatedFileRepository) CreateBatch(ctx context.Context, files []model.GeneratedFile) error {
	if len(files) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&files).Error
}

// GetByReportID retrieves all generated files for a report.
func (r *generatedFileRepository) GetByReportID(ctx context.Context, reportID string) ([]model.GeneratedFile, error) {
	var files []model.GeneratedFile
	err := r.db.WithContext(ctx).Where("report_id = ?", reportID).Find(&files).Error
	return files, err
}

// GetByReportIDAndFilename retrieves a specific generated file by report ID and filename.
func (r *generatedFileRepository) GetByReportIDAndFilename(ctx context.Context, reportID string, filename string) (*model.GeneratedFile, error) {
	var file model.GeneratedFile
	err := r.db.WithContext(ctx).Where("report_id = ? AND filename = ?", reportID, filename).First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}
