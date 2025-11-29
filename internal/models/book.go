package models

import (
	"time"

	"gorm.io/gorm"
)

// Book represents a unique literary work, independent of format.
// The same book can exist as both an audiobook and ebook, tracked via BookType.
// Relationships to authors, narrators, series, and files are managed through
// separate join tables to support many-to-many relationships.
type Book struct {
	gorm.Model

	// Title is the main title of the book
	Title string `gorm:"not null;index"`

	// Subtitle is an optional secondary title
	Subtitle string

	// Description contains the book summary or blurb
	Description string

	// Language is the ISO 639-1 language code (e.g., "en", "es", "de")
	Language string `gorm:"default:'en'"`

	// PublishDate is the original publication date
	PublishDate *time.Time

	// Publisher is the name of the publishing company
	Publisher string

	// BookType indicates the available formats: "audiobook", "ebook", or "both"
	BookType string `gorm:"not null"`

	// Duration is the audiobook length in minutes (0 for ebooks)
	Duration int

	// Relationships
	Authors     []BookAuthor     `gorm:"foreignKey:BookID"`
	Narrators   []BookNarrator   `gorm:"foreignKey:BookID"`
	Series      []BookSeries     `gorm:"foreignKey:BookID"`
	Identifiers []BookIdentifier `gorm:"foreignKey:BookID"`
	Files       []BookFile       `gorm:"foreignKey:BookID"`
}
