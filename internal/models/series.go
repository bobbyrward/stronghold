package models

import "gorm.io/gorm"

// Series represents a collection of related books.
// Examples include "The Stormlight Archive", "Harry Potter", or "The Expanse".
// Books are linked to series through the BookSeries join table which also
// tracks the position of each book within the series.
type Series struct {
	gorm.Model

	// Name is the series title (e.g., "The Stormlight Archive")
	Name string `gorm:"not null;uniqueIndex"`

	// Description provides information about the series
	Description string
}
