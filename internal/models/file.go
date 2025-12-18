package models

import (
	"time"

	"gorm.io/gorm"
)

// File represents uploaded files in the system
type File struct {
	ID           int64          `json:"id" gorm:"type:bigint;primaryKey;autoIncrement"`
	OriginalName string         `json:"original_name" gorm:"type:varchar(255);not null"`
	FileName     string         `json:"file_name" gorm:"type:varchar(255);not null;uniqueIndex"`
	FilePath     string         `json:"file_path" gorm:"type:varchar(500);not null"`
	FileSize     int64          `json:"file_size" gorm:"not null"`
	MimeType     string         `json:"mime_type" gorm:"type:varchar(100);not null"`
	FileType     string         `json:"file_type" gorm:"type:varchar(50);not null"` // document, spreadsheet, image
	Extension    string         `json:"extension" gorm:"type:varchar(10);not null"`
	UploadedBy   int64          `json:"uploaded_by" gorm:"type:bigint;not null"`
	UploadedByAdmin *Admin      `json:"uploaded_by_admin,omitempty" gorm:"foreignKey:UploadedBy;references:ID"`
	IsActive     bool           `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
