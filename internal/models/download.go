package models

import (
	"time"

	"gorm.io/gorm"
)

// Download tracks active and completed torrent downloads.
// It connects the torrent source to the eventual book in the library.
// Status lifecycle: downloading -> completed -> importing -> imported
// Failed downloads can be marked as "failed" or "manual_intervention".
type Download struct {
	gorm.Model

	// TorrentHash is the unique hash of the torrent
	TorrentHash string `gorm:"uniqueIndex"`

	// TorrentName is the display name of the torrent
	TorrentName string

	// Category is the qBittorrent category (e.g., "audiobooks", "ebooks")
	Category string

	// BookID links to the book once identified (optional until import)
	BookID *uint

	// Status tracks the download lifecycle:
	// - downloading: torrent active in qBittorrent
	// - completed: download finished, awaiting import
	// - importing: import process running
	// - imported: successfully imported to library
	// - failed: import failed (check ErrorMessage)
	// - manual_intervention: needs manual review
	Status string `gorm:"not null;default:'downloading';index"`

	// ImportedAt is when the download was successfully imported
	ImportedAt *time.Time

	// ErrorMessage contains error details if Status is "failed"
	ErrorMessage string

	// RetryCount tracks how many import attempts have been made
	RetryCount int

	// Relationships
	Book *Book `gorm:"foreignKey:BookID"`
}
