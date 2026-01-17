package hardcover

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/hasura/go-graphql-client"
)

// RealClient is the production implementation of the Hardcover API client.
type RealClient struct {
	graphqlClient *graphql.Client
	token         string
}

// Compile-time check that RealClient implements Client interface.
var _ Client = (*RealClient)(nil)

// NewClient creates a new Hardcover API client with the given authentication token.
func NewClient(token string) *RealClient {
	return &RealClient{
		graphqlClient: graphql.NewClient("https://api.hardcover.app/v1/graphql", http.DefaultClient).
			WithRequestModifier(func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer "+token)
				r.Header.Set("User-Agent", "stronghold/1.0")
			}),
		token: token,
	}
}

// SearchAuthors searches for authors by name query using the Hardcover GraphQL API.
func (c *RealClient) SearchAuthors(ctx context.Context, query string) ([]AuthorSearchResult, error) {
	slog.InfoContext(ctx, "Searching Hardcover authors", slog.String("query", query))

	var q struct {
		Authors []struct {
			Name string `graphql:"name"`
			Slug string `graphql:"slug"`
		} `graphql:"authors(where: {name: {_eq: $name}})"`
	}

	/*
	   query MyQuery {
	       authors(where: {name: {_eq: "Greg Tolley"}}) {
	           slug
	           identifiers
	       }
	   }

	*/

	variables := map[string]interface{}{
		"name": query,
	}

	err := c.graphqlClient.Query(ctx, &q, variables)
	if err != nil {
		return nil, err
	}

	searchResults := make([]AuthorSearchResult, len(q.Authors))

	for idx, author := range q.Authors {
		searchResults[idx] = AuthorSearchResult{
			Slug: author.Slug,
			Name: author.Name,
		}
	}

	return searchResults, nil
}

// GetAuthorBySlug retrieves an author by their unique slug using the Hardcover GraphQL API.
func (c *RealClient) GetAuthorBySlug(ctx context.Context, slug string) (*AuthorSearchResult, error) {
	slog.InfoContext(ctx, "Getting Hardcover author by slug", slog.String("slug", slug))
	// TODO: Implement GraphQL query
	// POST to c.baseURL with Bearer token c.token
	// GraphQL query for getting author by slug
	return nil, nil
}
