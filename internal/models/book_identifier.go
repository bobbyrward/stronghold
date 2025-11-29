package models

import "gorm.io/gorm"

// BookIdentifier stores external identifiers for books.
// Each book can have multiple identifiers from different sources.
// Common types include: isbn, isbn13, asin, goodreads, librarything, audible.
// The unique index on BookID+Type ensures only one identifier per type per book.
type BookIdentifier struct {
	gorm.Model

	// BookID references the book
	BookID uint `gorm:"not null;index"`

	// Type is the identifier system (e.g., "isbn", "asin", "goodreads")
	Type string `gorm:"not null;uniqueIndex:idx_book_identifier_type"`

	// Value is the actual identifier value
	Value string `gorm:"not null"`

	// Composite unique index to ensure one identifier per type per book
	// Note: BookID is included in the uniqueIndex via the Type field definition
	_ struct{} `gorm:"uniqueIndex:idx_book_identifier_type"`

	// Relationships
	Book Book `gorm:"foreignKey:BookID"`
}
