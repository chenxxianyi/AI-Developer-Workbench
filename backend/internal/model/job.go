package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Job represents an async tool execution task.
type Job struct {
	ID            string     `gorm:"type:char(36);primaryKey" json:"id"`
	ToolType      string     `gorm:"type:varchar(32);not null;index" json:"tool_type"`
	ReportID      string     `gorm:"type:char(36);not null;index" json:"report_id"`
	ProjectID     *string    `gorm:"type:char(36);index" json:"project_id,omitempty"`
	Status        string     `gorm:"type:varchar(24);not null;default:'queued';index" json:"status"`
	Progress      int        `gorm:"type:tinyint unsigned;not null;default:0" json:"progress"`
	Phase         string     `gorm:"type:varchar(128);not null;default:''" json:"phase"`
	ErrorMessage  string     `gorm:"type:text" json:"error_message,omitempty"`
	RetryOfJobID  *string    `gorm:"type:char(36);index" json:"retry_of_job_id,omitempty"`
	RetryCount    int        `gorm:"type:tinyint unsigned;not null;default:0" json:"retry_count"`
	CreatedAt     time.Time  `gorm:"not null" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"not null" json:"updated_at"`

	// Relations
	Report    *Report `gorm:"foreignKey:ReportID;constraint:OnDelete:CASCADE" json:"-"`
}

// Job status constants.
const (
	JobStatusQueued          = "queued"
	JobStatusRunning         = "running"
	JobStatusSucceeded       = "succeeded"
	JobStatusFailed          = "failed"
	JobStatusCancelRequested = "cancel_requested"
	JobStatusCanceled        = "canceled"
)

// ValidJobTransitions maps current status to allowed next statuses.
var ValidJobTransitions = map[string][]string{
	JobStatusQueued:          {JobStatusRunning, JobStatusCancelRequested},
	JobStatusRunning:         {JobStatusSucceeded, JobStatusFailed, JobStatusCancelRequested},
	JobStatusCancelRequested: {JobStatusCanceled, JobStatusFailed},
	JobStatusFailed:          {}, // terminal
	JobStatusCanceled:        {}, // terminal
	JobStatusSucceeded:       {}, // terminal
}

// IsTerminal returns true if the job is in a terminal state.
func IsTerminalJobStatus(status string) bool {
	return status == JobStatusSucceeded || status == JobStatusFailed || status == JobStatusCanceled
}

// CanTransition checks if a state transition is allowed.
func CanTransition(from, to string) bool {
	allowed, ok := ValidJobTransitions[from]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}

// BeforeCreate sets the Job ID before creation.
func (j *Job) BeforeCreate(tx *gorm.DB) error {
	if j.ID == "" {
		j.ID = uuid.New().String()
	}
	return nil
}
