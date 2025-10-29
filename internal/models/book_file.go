package models

import (
	"time"

	"gorm.io/gorm"
)

// BookFile represents a physical file in the library.
// Each book can have multiple files in different formats (epub, m4b, etc.).
// The Checksum field enables duplicate detection across imports.
type BookFile struct {
	gorm.Model

	// BookID references the book this file belongs to
	BookID uint `gorm:"not null;index"`

	// FilePath is the full absolute path to the file
	FilePath string `gorm:"not null"`

	// FileName is the base name of the file
	FileName string `gorm:"not null"`

	// FileType is the file format: epub, mobi, azw3, m4b, mp3
	FileType string `gorm:"not null"`

	// FileSize is the size in bytes
	FileSize int64

	// Checksum is the SHA256 hash for deduplication
	Checksum string

	// AddedAt is when the file was added to the library
	AddedAt time.Time

	// Relationships
	Book Book `gorm:"foreignKey:BookID"`
}
