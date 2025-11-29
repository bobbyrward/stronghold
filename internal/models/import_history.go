package models

import (
	"time"

	"gorm.io/gorm"
)

// ImportHistory provides an audit trail of all imports.
// Each record tracks a single file import from a download to the library.
// This allows tracking which downloads produced which files and when.
type ImportHistory struct {
	gorm.Model

	// DownloadID references the source download
	DownloadID uint `gorm:"not null"`

	// BookID references the book that was imported
	BookID uint `gorm:"not null"`

	// BookFileID references the file that was created
	BookFileID uint `gorm:"not null"`

	// SourcePath is the original location of the file
	SourcePath string

	// DestPath is the final location in the library
	DestPath string

	// ImportedAt is when this import occurred
	ImportedAt time.Time

	// Relationships
	Download Download `gorm:"foreignKey:DownloadID"`
	Book     Book     `gorm:"foreignKey:BookID"`
	BookFile BookFile `gorm:"foreignKey:BookFileID"`
}
