package models

import "gorm.io/gorm"

// BookNarrator is a join table linking audiobooks to their narrators.
// This enables many-to-many relationships where an audiobook can have multiple
// narrators and a narrator (Person) can narrate multiple books.
// The composite unique index on BookID and PersonID prevents duplicate entries.
type BookNarrator struct {
	gorm.Model

	// BookID references the book
	BookID uint `gorm:"not null;uniqueIndex:idx_book_narrator"`

	// PersonID references the narrator (Person)
	PersonID uint `gorm:"not null;uniqueIndex:idx_book_narrator"`

	// Relationships
	Book   Book   `gorm:"foreignKey:BookID"`
	Person Person `gorm:"foreignKey:PersonID"`
}
