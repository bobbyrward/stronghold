package hardcover

import (
	"context"
	"log/slog"
)

// RealClient is the production implementation of the Hardcover API client.
type RealClient struct {
	token   string
	baseURL string
}

// Compile-time check that RealClient implements Client interface.
var _ Client = (*RealClient)(nil)

// NewClient creates a new Hardcover API client with the given authentication token.
func NewClient(token string) *RealClient {
	return &RealClient{
		token:   token,
		baseURL: "https://api.hardcover.app/v1/graphql",
	}
}

// SearchAuthors searches for authors by name query using the Hardcover GraphQL API.
func (c *RealClient) SearchAuthors(ctx context.Context, query string) ([]AuthorSearchResult, error) {
	slog.InfoContext(ctx, "Searching Hardcover authors", slog.String("query", query))
	// TODO: Implement GraphQL query
	// POST to c.baseURL with Bearer token c.token
	// GraphQL query for searching authors by name
	return nil, nil
}

// GetAuthorBySlug retrieves an author by their unique slug using the Hardcover GraphQL API.
func (c *RealClient) GetAuthorBySlug(ctx context.Context, slug string) (*AuthorSearchResult, error) {
	slog.InfoContext(ctx, "Getting Hardcover author by slug", slog.String("slug", slug))
	// TODO: Implement GraphQL query
	// POST to c.baseURL with Bearer token c.token
	// GraphQL query for getting author by slug
	return nil, nil
}
