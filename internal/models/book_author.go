package models

import "gorm.io/gorm"

// BookAuthor is a join table linking books to their authors.
// This enables many-to-many relationships where a book can have multiple
// authors and an author (Person) can write multiple books.
// The composite unique index on BookID and PersonID prevents duplicate entries.
type BookAuthor struct {
	gorm.Model

	// BookID references the book
	BookID uint `gorm:"not null;uniqueIndex:idx_book_author"`

	// PersonID references the author (Person)
	PersonID uint `gorm:"not null;uniqueIndex:idx_book_author"`

	// Relationships
	Book   Book   `gorm:"foreignKey:BookID"`
	Person Person `gorm:"foreignKey:PersonID"`
}
