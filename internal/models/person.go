package models

import "gorm.io/gorm"

// Person represents an individual who can be an author, narrator, or both.
// This is the base table for tracking people involved with books.
// A single person can have multiple roles - for example, an author who
// narrates their own memoir would appear in both BookAuthor and BookNarrator.
type Person struct {
	gorm.Model

	// Name is the display name of the person (e.g., "Stephen King")
	Name string `gorm:"not null;uniqueIndex"`

	// SortName is used for alphabetical sorting (e.g., "King, Stephen")
	SortName string

	// Description contains biographical information about the person
	Description string
}
