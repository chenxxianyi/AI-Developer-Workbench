package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Report represents a tool execution report.
type Report struct {
	ID           string         `gorm:"type:char(36);primaryKey" json:"id"`
	ToolType     string         `gorm:"type:varchar(32);not null" json:"tool_type"`
	Title        string         `gorm:"type:varchar(255);not null" json:"title"`
	InputMode    string         `gorm:"type:varchar(32);not null;default:''" json:"input_mode"`
	Status       string         `gorm:"type:varchar(24);not null;default:'processing'" json:"status"`
	Summary      string         `gorm:"type:text" json:"summary"`
	TotalScore   *int           `gorm:"type:smallint unsigned" json:"total_score"`
	Grade        *string        `gorm:"type:varchar(64)" json:"grade"`
	InputJSON    datatypes.JSON `gorm:"type:json" json:"input_json"`
	ReportJSON   datatypes.JSON `gorm:"type:json;not null" json:"report_json"`
	FilePath     string         `gorm:"type:varchar(1024)" json:"file_path"`
	FileURL      string         `gorm:"type:varchar(1024)" json:"file_url"`
	ErrorMessage string         `gorm:"type:text" json:"error_message"`
	ParentReportID *string      `gorm:"type:char(36);index" json:"parent_report_id,omitempty"`
	ProjectID    *string        `gorm:"type:char(36);index" json:"project_id,omitempty"`
	CreatedAt    time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"not null" json:"updated_at"`

	// Relations
	GeneratedFiles []GeneratedFile `gorm:"foreignKey:ReportID;constraint:OnDelete:CASCADE" json:"generated_files"`
	Assets         []ReportAsset   `gorm:"foreignKey:ReportID;constraint:OnDelete:CASCADE" json:"assets"`
	// ParentReport is the report this one re-runs from (nullable). On parent
	// delete the FK is set NULL so the child report survives.
	ParentReport *Report `gorm:"foreignKey:ParentReportID;constraint:OnDelete:SET NULL" json:"parent_report,omitempty"`
	// Project is the project this report belongs to (nullable). On project
	// delete the FK is set NULL so the report survives.
	Project *Project `gorm:"foreignKey:ProjectID;constraint:OnDelete:SET NULL" json:"project,omitempty"`
}

// BeforeCreate sets UUID and timestamps before creating a report.
func (r *Report) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	r.CreatedAt = time.Now().UTC()
	r.UpdatedAt = time.Now().UTC()
	return nil
}

// BeforeUpdate sets the updated_at timestamp.
func (r *Report) BeforeUpdate(tx *gorm.DB) error {
	r.UpdatedAt = time.Now().UTC()
	return nil
}

// TableName returns the table name for Report.
func (Report) TableName() string {
	return "reports"
}

// Report status constants.
const (
	StatusProcessing = "processing"
	StatusSucceeded  = "succeeded"
	StatusFallback   = "fallback"
	StatusFailed     = "failed"
)

// Tool type constants.
const (
	ToolTypeUIReview      = "ui_review"
	ToolTypeProjectDoctor = "project_doctor"
	ToolTypeAgentConfig   = "agent_config"
	ToolTypeAPIDoc        = "api_doc"
	ToolTypeDBSchema      = "db_schema"
)

// ValidToolTypes returns the list of valid tool types.
func ValidToolTypes() []string {
	return []string{
		ToolTypeUIReview,
		ToolTypeProjectDoctor,
		ToolTypeAgentConfig,
		ToolTypeAPIDoc,
		ToolTypeDBSchema,
	}
}

// IsValidToolType checks if a tool type is valid.
func IsValidToolType(toolType string) bool {
	for _, t := range ValidToolTypes() {
		if t == toolType {
			return true
		}
	}
	return false
}
