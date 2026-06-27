package hardcover

import "context"

// Client defines the interface for interacting with the Hardcover API.
type Client interface {
	// SearchAuthors searches for authors by name query.
	SearchAuthors(ctx context.Context, query string) ([]AuthorSearchResult, error)

	// GetAuthorBySlug retrieves an author by their slug. Only used by the
	// slug→id backfill; new links validate by id.
	GetAuthorBySlug(ctx context.Context, slug string) (*AuthorSearchResult, error)

	// GetAuthorByID retrieves an author by their canonical id.
	GetAuthorByID(ctx context.Context, id string) (*AuthorSearchResult, error)

	// GetAuthorBooks fetches an author's bibliography (list of works) by their
	// canonical id, ordered by release date.
	GetAuthorBooks(ctx context.Context, id string) ([]BookResult, error)
}
