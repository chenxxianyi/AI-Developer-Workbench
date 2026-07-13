package repository

import (
	"context"

	"ai-developer-workbench/internal/model"

	"gorm.io/gorm"
)

// JobRepository handles database access for jobs.
type JobRepository interface {
	Create(ctx context.Context, job *model.Job) error
	GetByID(ctx context.Context, id string) (*model.Job, error)
	Update(ctx context.Context, job *model.Job) error
	GetRunningJobs(ctx context.Context) ([]model.Job, error)
}

type jobRepository struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) JobRepository {
	return &jobRepository{db: db}
}

func (r *jobRepository) Create(ctx context.Context, job *model.Job) error {
	return r.db.WithContext(ctx).Create(job).Error
}

func (r *jobRepository) GetByID(ctx context.Context, id string) (*model.Job, error) {
	var job model.Job
	err := r.db.WithContext(ctx).First(&job, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *jobRepository) Update(ctx context.Context, job *model.Job) error {
	return r.db.WithContext(ctx).Save(job).Error
}

func (r *jobRepository) GetRunningJobs(ctx context.Context) ([]model.Job, error) {
	var jobs []model.Job
	err := r.db.WithContext(ctx).Where("status IN ?", []string{model.JobStatusQueued, model.JobStatusRunning}).Find(&jobs).Error
	return jobs, err
}
