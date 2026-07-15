package service

import (
	"errors"
	"fmt"
	"time"

	"ai-developer-workbench/internal/model"
	"ai-developer-workbench/pkg/sse"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TaskService manages task lifecycle, state transitions, and SSE publishing.
type TaskService struct {
	db     *gorm.DB
	broker *sse.Broker
}

// NewTaskService creates a new TaskService.
func NewTaskService(db *gorm.DB, broker *sse.Broker) *TaskService {
	return &TaskService{db: db, broker: broker}
}

// ── Task CRUD ──

// Create creates a new task and publishes a "created" event.
func (s *TaskService) Create(projectID, userID, taskType string) (*model.Task, error) {
	task := &model.Task{
		ID:         uuid.New().String(),
		ProjectID:  projectID,
		UserID:     userID,
		Type:       taskType,
		Status:     "pending",
		MaxRetries: 3,
		Retryable:  true,
	}
	if err := s.db.Create(task).Error; err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}
	return task, nil
}

// Get retrieves a task by ID.
func (s *TaskService) Get(taskID string) (*model.Task, error) {
	var task model.Task
	if err := s.db.First(&task, "id = ?", taskID).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

// ── State Machine ──

var validTransitions = map[string][]string{
	"pending":   {"running"},
	"running":   {"success", "failed", "cancelled"},
	"failed":    {"pending"}, // retry
	"cancelled": {},
	"success":   {},
}

// Start transitions a task from pending → running.
func (s *TaskService) Start(taskID string) error {
	task, err := s.Get(taskID)
	if err != nil {
		return err
	}
	if task.Status != "pending" {
		return errors.New("only pending tasks can be started")
	}
	now := time.Now()
	return s.db.Model(task).Updates(map[string]interface{}{
		"status":     "running",
		"started_at": &now,
	}).Error
}

// UpdateProgress updates task progress and optional stage/message.
func (s *TaskService) UpdateProgress(taskID string, progress int, stage, message string) error {
	return s.db.Model(&model.Task{}).Where("id = ?", taskID).Updates(map[string]interface{}{
		"progress": progress,
		"stage":    stage,
		"message":  message,
	}).Error
}

// Complete transitions a task to success.
func (s *TaskService) Complete(taskID string, result string) error {
	now := time.Now()
	return s.db.Model(&model.Task{}).Where("id = ? AND status = ?", taskID, "running").Updates(map[string]interface{}{
		"status":      "success",
		"progress":    100,
		"result":      result,
		"finished_at": &now,
	}).Error
}

// Fail transitions a task to failed with error details.
func (s *TaskService) Fail(taskID string, errorCode, errorDetail string) error {
	now := time.Now()
	return s.db.Model(&model.Task{}).Where("id = ? AND status = ?", taskID, "running").Updates(map[string]interface{}{
		"status":       "failed",
		"error_code":   errorCode,
		"error_detail": errorDetail,
		"finished_at":  &now,
	}).Error
}

// Cancel transitions a running task to cancelled.
func (s *TaskService) Cancel(taskID string) error {
	return s.db.Model(&model.Task{}).Where("id = ? AND status = ?", taskID, "running").Update("status", "cancelled").Error
}

// Retry resets a failed task back to pending (if retryable and within limit).
func (s *TaskService) Retry(taskID string) (*model.Task, error) {
	task, err := s.Get(taskID)
	if err != nil {
		return nil, err
	}
	if task.Status != "failed" {
		return nil, errors.New("only failed tasks can be retried")
	}
	if !task.Retryable || task.RetryCount >= task.MaxRetries {
		return nil, errors.New("task is not retryable")
	}

	if err := s.db.Model(task).Updates(map[string]interface{}{
		"status":      "pending",
		"retry_count": task.RetryCount + 1,
		"progress":    0,
		"error_code":  "",
	}).Error; err != nil {
		return nil, err
	}

	// Refresh
	return s.Get(taskID)
}

// ── SSE Publishing ──

// PublishEvent sends an SSE event to all subscribers of the task.
func (s *TaskService) PublishEvent(taskID string, event sse.Event) {
	s.broker.Publish(taskID, event)
}

// PublishProgress persists progress and sends an SSE event.
func (s *TaskService) PublishProgress(taskID string, stage string, progress int, message string) {
	_ = s.UpdateProgress(taskID, progress, stage, message)
	s.PublishEvent(taskID, sse.Event{
		Type:      "stage_progress",
		TaskID:    taskID,
		Status:    "running",
		Stage:     stage,
		Progress:  progress,
		Message:   message,
		Timestamp: time.Now().UnixMilli(),
	})
}

// PublishCompleted sends a task completion event.
func (s *TaskService) PublishCompleted(taskID string) {
	s.PublishEvent(taskID, sse.Event{
		Type:      "task_completed",
		TaskID:    taskID,
		Status:    "success",
		Progress:  100,
		Timestamp: time.Now().UnixMilli(),
	})
}

// PublishFailed sends a task failure event.
func (s *TaskService) PublishFailed(taskID string, errMsg string) {
	s.PublishEvent(taskID, sse.Event{
		Type:      "task_failed",
		TaskID:    taskID,
		Status:    "failed",
		Message:   errMsg,
		Timestamp: time.Now().UnixMilli(),
	})
}

// ── Cleanup ──

// CleanupStaleTasks marks tasks running longer than timeout as failed.
func (s *TaskService) CleanupStaleTasks(timeout time.Duration) error {
	cutoff := time.Now().Add(-timeout)
	return s.db.Model(&model.Task{}).
		Where("status = ? AND started_at < ?", "running", cutoff).
		Updates(map[string]interface{}{
			"status":       "failed",
			"error_code":   "TASK_TIMEOUT",
			"error_detail": fmt.Sprintf("任务运行超过 %v 自动失败", timeout),
			"finished_at":  time.Now(),
		}).Error
}
