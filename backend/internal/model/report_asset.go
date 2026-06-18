package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ReportAsset represents an uploaded file associated with a report.
type ReportAsset struct {
	ID           string    `gorm:"type:char(36);primaryKey" json:"id"`
	ReportID     string    `gorm:"type:char(36);not null;index:idx_report_assets_report" json:"report_id"`
	AssetType    string    `gorm:"type:varchar(32);not null" json:"asset_type"`
	OriginalName string    `gorm:"type:varchar(255);not null" json:"original_name"`
	StoredName   string    `gorm:"type:varchar(255);not null" json:"stored_name"`
	RelativePath string    `gorm:"type:varchar(1024);not null" json:"relative_path"`
	MimeType     string    `gorm:"type:varchar(100)" json:"mime_type"`
	SizeBytes    uint64    `gorm:"type:bigint unsigned;not null" json:"size_bytes"`
	SHA256       string    `gorm:"type:char(64)" json:"sha256"`
	CreatedAt    time.Time `gorm:"not null" json:"created_at"`
}

// BeforeCreate sets UUID and timestamp before creating a report asset.
func (r *ReportAsset) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	r.CreatedAt = time.Now().UTC()
	return nil
}

// TableName returns the table name for ReportAsset.
func (ReportAsset) TableName() string {
	return "report_assets"
}

// Asset type constants.
const (
	AssetTypeScreenshot = "screenshot"
	AssetTypeProjectZip  = "project_zip"
	AssetTypeSourceFile  = "source_file"
)
