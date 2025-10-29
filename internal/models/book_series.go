package models

import "gorm.io/gorm"

// BookSeries is a join table linking books to series with position information.
// This enables many-to-many relationships where a book can belong to multiple
// series (rare but possible) and a series contains multiple books.
// The Position field uses float64 to accommodate novellas (e.g., 1.5, 2.5).
type BookSeries struct {
	gorm.Model

	// BookID references the book
	BookID uint `gorm:"not null;uniqueIndex:idx_book_series"`

	// SeriesID references the series
	SeriesID uint `gorm:"not null;uniqueIndex:idx_book_series"`

	// Position is the book's position within the series.
	// Uses float64 to support novellas like "Book 1.5" or "Book 2.5".
	// Standard entries use whole numbers (1.0, 2.0, 3.0).
	Position float64

	// Relationships
	Book   Book   `gorm:"foreignKey:BookID"`
	Series Series `gorm:"foreignKey:SeriesID"`
}
