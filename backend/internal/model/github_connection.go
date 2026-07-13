package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GitHubConnection stores an OAuth connection to a GitHub account.
type GitHubConnection struct {
	ID             string    `gorm:"type:char(36);primaryKey" json:"id"`
	ProjectID      *string   `gorm:"type:char(36);index" json:"project_id,omitempty"`
	GitHubUsername string    `gorm:"type:varchar(128);not null" json:"github_username"`
	AccessToken    string    `gorm:"type:text;not null" json:"-"` // never exposed
	Repository     string    `gorm:"type:varchar(256)" json:"repository"`
	Branch         string    `gorm:"type:varchar(128);not null;default:'main'" json:"branch"`
	Active         bool      `gorm:"not null;default:true" json:"active"`
	CreatedAt      time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt      time.Time `gorm:"not null" json:"updated_at"`
}

func (g *GitHubConnection) BeforeCreate(tx *gorm.DB) error {
	if g.ID == "" {
		g.ID = uuid.New().String()
	}
	return nil
}
