package hardcover

// AuthorSearchResult represents an author returned from Hardcover API searches.
type AuthorSearchResult struct {
	ID   string // canonical author id (stable); used as Author.HardcoverRef
	Slug string // for UI display + hardcover.app link only
	Name string
}
