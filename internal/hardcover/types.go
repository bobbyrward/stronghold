package hardcover

// AuthorSearchResult represents an author returned from Hardcover API searches.
type AuthorSearchResult struct {
	ID   string // canonical author id (stable); used as Author.HardcoverRef
	Slug string // for UI display + hardcover.app link only
	Name string
}

// BookResult is one work from an author's Hardcover bibliography. Maps to a
// catalog Book; editions are deliberately not pulled (see docs/hardcover-api.md).
type BookResult struct {
	HardcoverID string  // Hardcover books (work) id, stored as decimal string
	Title       string
	ReleaseDate *string // ISO date from Hardcover; nil when unknown
}
