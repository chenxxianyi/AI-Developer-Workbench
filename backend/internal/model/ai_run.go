package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AIRun records a single AI service call for observability.
type AIRun struct {
	ID           string    `gorm:"type:char(36);primaryKey" json:"id"`
	ReportID     string    `gorm:"type:char(36);not null;index" json:"report_id"`
	JobID        *string   `gorm:"type:char(36);index" json:"job_id,omitempty"`
	ToolType     string    `gorm:"type:varchar(32);not null;index" json:"tool_type"`
	Provider     string    `gorm:"type:varchar(64);not null" json:"provider"`
	Model        string    `gorm:"type:varchar(128);not null" json:"model"`
	IsMock       bool      `gorm:"not null;default:false" json:"is_mock"`
	DurationMs   int64     `gorm:"not null;default:0" json:"duration_ms"`
	RetryCount   int       `gorm:"type:tinyint unsigned;not null;default:0" json:"retry_count"`
	ParseSuccess bool      `gorm:"not null;default:true" json:"parse_success"`
	FallbackUsed bool      `gorm:"not null;default:false" json:"fallback_used"`
	ErrorType    string    `gorm:"type:varchar(64);not null;default:''" json:"error_type"`
	CreatedAt    time.Time `gorm:"not null" json:"created_at"`

	Report *Report `gorm:"foreignKey:ReportID;constraint:OnDelete:CASCADE" json:"-"`
	Job    *Job    `gorm:"foreignKey:JobID;constraint:OnDelete:SET NULL" json:"-"`
}

// Error type constants.
const (
	AIErrorTimeout     = "timeout"
	AIErrorRateLimit   = "rate_limit"
	AIErrorProvider4xx = "provider_4xx"
	AIErrorProvider5xx = "provider_5xx"
	AIErrorNetwork     = "network"
	AIErrorInvalidJSON = "invalid_json"
	AIErrorValidation  = "validation"
	AIErrorCanceled    = "canceled"
	AIErrorInternal    = "internal"
)

func (a *AIRun) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return nil
}
