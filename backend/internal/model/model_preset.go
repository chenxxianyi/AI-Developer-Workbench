package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ModelPreset stores reusable AI model configurations.
type ModelPreset struct {
	ID             string    `gorm:"type:char(36);primaryKey" json:"id"`
	Name           string    `gorm:"type:varchar(128);not null;uniqueIndex" json:"name"`
	Provider       string    `gorm:"type:varchar(64);not null" json:"provider"`
	BaseURL        string    `gorm:"type:varchar(512);not null" json:"base_url"`
	Model          string    `gorm:"type:varchar(128);not null" json:"model"`
	VisionModel    string    `gorm:"type:varchar(128);not null" json:"vision_model"`
	TimeoutSeconds int       `gorm:"not null;default:90" json:"timeout_seconds"`
	MaxRetries     int       `gorm:"not null;default:1" json:"max_retries"`
	Status         string    `gorm:"type:varchar(20);not null;default:active" json:"status"`
	IsDefault      bool      `gorm:"not null;default:false" json:"is_default"`
	CreatedAt      time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt      time.Time `gorm:"not null" json:"updated_at"`
}

func (m *ModelPreset) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return nil
}
