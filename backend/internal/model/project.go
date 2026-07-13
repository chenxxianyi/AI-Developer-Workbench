package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Project represents a project profile that reports can be associated with.
type Project struct {
	ID            string         `gorm:"type:char(36);primaryKey" json:"id"`
	Name          string         `gorm:"type:varchar(128);not null" json:"name"`
	Description   string         `gorm:"type:text" json:"description"`
	RepoURL       string         `gorm:"type:varchar(512)" json:"repo_url"`
	FrontendStack string         `gorm:"type:varchar(256)" json:"frontend_stack"`
	BackendStack  string         `gorm:"type:varchar(256)" json:"backend_stack"`
	Database      string         `gorm:"type:varchar(128)" json:"database"`
	UIStyle       string         `gorm:"type:varchar(256)" json:"ui_style"`
	CodingRules   string         `gorm:"type:text" json:"coding_rules"`
	CreatedAt     time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"not null" json:"updated_at"`

	// Reports associated with this project (0..N). On project delete the FK is
	// SET NULL so historical reports survive (product rule: no cascade delete).
	Reports []Report `gorm:"foreignKey:ProjectID;constraint:OnDelete:SET NULL" json:"reports,omitempty"`
}

// BeforeCreate sets UUID and timestamps before creating a project.
func (p *Project) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	p.CreatedAt = time.Now().UTC()
	p.UpdatedAt = time.Now().UTC()
	return nil
}

// BeforeUpdate sets the updated_at timestamp.
func (p *Project) BeforeUpdate(tx *gorm.DB) error {
	p.UpdatedAt = time.Now().UTC()
	return nil
}

// TableName returns the table name for Project.
func (Project) TableName() string {
	return "projects"
}
