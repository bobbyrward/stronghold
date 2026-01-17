package hardcover

import "context"

// Client defines the interface for interacting with the Hardcover API.
type Client interface {
	// SearchAuthors searches for authors by name query.
	SearchAuthors(ctx context.Context, query string) ([]AuthorSearchResult, error)

	// GetAuthorBySlug retrieves an author by their unique slug.
	GetAuthorBySlug(ctx context.Context, slug string) (*AuthorSearchResult, error)
}
