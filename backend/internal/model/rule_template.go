package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RuleTemplate stores reusable prompt templates and coding rules.
type RuleTemplate struct {
	ID        string    `gorm:"type:char(36);primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(128);not null;uniqueIndex" json:"name"`
	Category  string    `gorm:"type:varchar(64);not null;index" json:"category"` // ui_review, db_schema, coding_rules, agent_config, general
	Version   int       `gorm:"not null;default:1" json:"version"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	IsDefault bool      `gorm:"not null;default:false" json:"is_default"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
}

func (t *RuleTemplate) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}
