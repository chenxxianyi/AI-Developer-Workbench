package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"ai-developer-workbench/internal/dto"
	"ai-developer-workbench/internal/model"
	"ai-developer-workbench/internal/repository"

	"gorm.io/gorm"
)

// JobService manages async job lifecycle and state machine.
type JobService interface {
	CreateJob(ctx context.Context, toolType, reportID string, projectID *string) (*dto.JobDTO, error)
	GetJob(ctx context.Context, id string) (*dto.JobDTO, error)
	UpdateProgress(ctx context.Context, id string, progress int, phase string) error
	SucceedJob(ctx context.Context, id string) error
	FailJob(ctx context.Context, id string, errMsg string) error
	RequestCancel(ctx context.Context, id string) error
	CancelJob(ctx context.Context, id string) error
	RetryJob(ctx context.Context, id string) (*dto.JobDTO, error)
	RecoverOrphanedJobs(ctx context.Context) error
	JobToDTO(job *model.Job) *dto.JobDTO
}

type jobService struct {
	db      *gorm.DB
	jobRepo repository.JobRepository
}

func NewJobService(db *gorm.DB, jobRepo repository.JobRepository) JobService {
	return &jobService{db: db, jobRepo: jobRepo}
}

func (s *jobService) CreateJob(ctx context.Context, toolType, reportID string, projectID *string) (*dto.JobDTO, error) {
	job := &model.Job{
		ToolType:  toolType,
		ReportID:  reportID,
		ProjectID: projectID,
		Status:    model.JobStatusQueued,
		Phase:     "queued",
	}
	if err := s.jobRepo.Create(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to create job: %w", err)
	}
	return s.JobToDTO(job), nil
}

func (s *jobService) GetJob(ctx context.Context, id string) (*dto.JobDTO, error) {
	job, err := s.jobRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.JobToDTO(job), nil
}

func (s *jobService) UpdateProgress(ctx context.Context, id string, progress int, phase string) error {
	job, err := s.jobRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if progress < 0 || progress > 100 {
		return fmt.Errorf("progress must be 0-100, got %d", progress)
	}
	job.Progress = progress
	job.Phase = phase
	return s.jobRepo.Update(ctx, job)
}

func (s *jobService) transitionJob(ctx context.Context, id, toStatus string, setFields func(j *model.Job)) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var job model.Job
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&job, "id = ?", id).Error; err != nil {
			return err
		}
		if !model.CanTransition(job.Status, toStatus) {
			return fmt.Errorf("cannot transition job %s from %s to %s", id, job.Status, toStatus)
		}
		job.Status = toStatus
		if setFields != nil {
			setFields(&job)
		}
		return tx.Save(&job).Error
	})
}

func (s *jobService) SucceedJob(ctx context.Context, id string) error {
	return s.transitionJob(ctx, id, model.JobStatusSucceeded, func(j *model.Job) {
		j.Progress = 100
		j.Phase = "completed"
	})
}

func (s *jobService) FailJob(ctx context.Context, id string, errMsg string) error {
	return s.transitionJob(ctx, id, model.JobStatusFailed, func(j *model.Job) {
		j.ErrorMessage = errMsg
		j.Phase = "failed"
	})
}

func (s *jobService) RequestCancel(ctx context.Context, id string) error {
	return s.transitionJob(ctx, id, model.JobStatusCancelRequested, func(j *model.Job) {
		j.Phase = "canceling"
	})
}

func (s *jobService) CancelJob(ctx context.Context, id string) error {
	return s.transitionJob(ctx, id, model.JobStatusCanceled, func(j *model.Job) {
		j.Phase = "canceled"
	})
}

func (s *jobService) RetryJob(ctx context.Context, id string) (*dto.JobDTO, error) {
	original, err := s.jobRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !model.IsTerminalJobStatus(original.Status) {
		return nil, fmt.Errorf("job %s is not in a terminal state (%s), cannot retry", id, original.Status)
	}
	newJob := &model.Job{
		ToolType:     original.ToolType,
		ReportID:     original.ReportID,
		ProjectID:    original.ProjectID,
		Status:       model.JobStatusQueued,
		Phase:        "queued (retry)",
		RetryOfJobID: &original.ID,
		RetryCount:   original.RetryCount + 1,
	}
	if err := s.jobRepo.Create(ctx, newJob); err != nil {
		return nil, err
	}
	return s.JobToDTO(newJob), nil
}

func (s *jobService) RecoverOrphanedJobs(ctx context.Context) error {
	jobs, err := s.jobRepo.GetRunningJobs(ctx)
	if err != nil {
		return err
	}
	for _, job := range jobs {
		slog.Warn("recovering orphaned job", "job_id", job.ID, "status", job.Status)
		_ = s.FailJob(ctx, job.ID, "job orphaned after service restart")
	}
	return nil
}

func (s *jobService) JobToDTO(job *model.Job) *dto.JobDTO {
	return &dto.JobDTO{
		ID:           job.ID,
		ToolType:     job.ToolType,
		ReportID:     job.ReportID,
		ProjectID:    job.ProjectID,
		Status:       job.Status,
		Progress:     job.Progress,
		Phase:        job.Phase,
		ErrorMessage: job.ErrorMessage,
		RetryOfJobID: job.RetryOfJobID,
		RetryCount:   job.RetryCount,
		CreatedAt:    job.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    job.UpdatedAt.Format(time.RFC3339),
	}
}
