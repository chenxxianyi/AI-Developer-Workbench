package model

import (
	"time"
)

// User represents a user account (unified from Builder).
type User struct {
	ID               string     `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Username         string     `gorm:"uniqueIndex;not null;size:100" json:"username"`
	Email            string     `gorm:"uniqueIndex;not null;size:255" json:"email"`
	PasswordHash     string     `gorm:"not null;size:255" json:"-"`
	Role             string     `gorm:"not null;default:user;size:20" json:"role"`    // user | admin
	Status           string     `gorm:"not null;default:active;size:20" json:"status"` // active | disabled
	LegacyBuilderID  *uint      `json:"legacy_builder_id,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// Requirement stores the structured project requirements.
type Requirement struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	ProjectID string    `gorm:"not null;index" json:"project_id"`
	Content   string    `gorm:"type:json;not null" json:"content"` // JSON
	Version   int       `gorm:"not null;default:1" json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Blueprint stores the AI-generated project blueprint.
type Blueprint struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	ProjectID string    `gorm:"not null;index" json:"project_id"`
	Content   string    `gorm:"type:json;not null" json:"content"` // JSON: product_positioning, pages, components, db, api, ui, tech_stack
	Status    string    `gorm:"not null;default:draft;size:20" json:"status"` // draft | generated | confirmed
	Version   int       `gorm:"not null;default:1" json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Task represents a generation/build/tool/report task (unified from Builder tasks + Workbench jobs).
type Task struct {
	ID          string     `gorm:"primaryKey;type:varchar(36)" json:"id"`
	ProjectID   string     `gorm:"index" json:"project_id"`
	UserID      string     `gorm:"index" json:"user_id"`
	Type        string     `gorm:"not null;size:30" json:"type"`       // generation | build | tool_run | report
	Status      string     `gorm:"not null;default:pending;size:20" json:"status"` // pending | running | success | failed | cancelled
	Progress    int        `gorm:"default:0" json:"progress"`
	Stage       string     `gorm:"size:100" json:"stage,omitempty"`
	Message     string     `gorm:"size:500" json:"message,omitempty"`
	Result      string     `gorm:"type:text" json:"result,omitempty"`
	ErrorCode   string     `gorm:"size:50" json:"error_code,omitempty"`
	ErrorDetail string     `gorm:"type:text" json:"error_detail,omitempty"`
	RetryCount  int        `gorm:"default:0" json:"retry_count"`
	MaxRetries  int        `gorm:"default:3" json:"max_retries"`
	Retryable   bool       `gorm:"default:true" json:"retryable"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	FinishedAt  *time.Time `json:"finished_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// ProjectFile represents a file in the project workspace.
type ProjectFile struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	ProjectID string    `gorm:"not null;index" json:"project_id"`
	Name      string    `gorm:"not null;size:500" json:"name"`
	Path      string    `gorm:"not null;size:1000" json:"path"`
	Size      int64     `gorm:"not null;default:0" json:"size"`
	MimeType  string    `gorm:"size:100" json:"mime_type,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
