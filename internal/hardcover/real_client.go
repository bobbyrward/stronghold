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

			Canonical *struct {
				Name string `graphql:"name"`
				Slug string `graphql:"slug"`
			}
		} `graphql:"authors(where: {name: {_eq: $name}})"`
	}

	variables := map[string]any{
		"name": query,
	}

	err := c.graphqlClient.Query(ctx, &q, variables)
	if err != nil {
		return nil, err
	}

	searchResults := make([]AuthorSearchResult, len(q.Authors))

	for idx, author := range q.Authors {
		slog.DebugContext(ctx, "Found author", slog.String("name", author.Name), slog.String("slug", author.Slug), slog.Any("canonical", author.Canonical))
		if author.Canonical != nil {
			searchResults[idx] = AuthorSearchResult{
				Slug: author.Canonical.Slug,
				Name: author.Canonical.Name,
			}
		} else {
			searchResults[idx] = AuthorSearchResult{
				Slug: author.Slug,
				Name: author.Name,
			}
		}
	}

	return searchResults, nil
}

// GetAuthorBySlug retrieves an author by their unique slug using the Hardcover GraphQL API.
func (c *RealClient) GetAuthorBySlug(ctx context.Context, slug string) (*AuthorSearchResult, error) {
	slog.InfoContext(ctx, "Getting Hardcover author by slug", slog.String("slug", slug))

	var q struct {
		Authors []struct {
			Name string `graphql:"name"`
			Slug string `graphql:"slug"`
		} `graphql:"authors(where: {slug: {_eq: $slug}})"`
	}

	variables := map[string]any{
		"slug": slug,
	}

	err := c.graphqlClient.Query(ctx, &q, variables)
	if err != nil {
		return nil, err
	}

	if len(q.Authors) == 0 {
		slog.WarnContext(ctx, "No author found with given slug", slog.String("slug", slug))
		return nil, nil
	}

	if len(q.Authors) > 1 {
		slog.WarnContext(ctx, "Multiple authors found with given slug", slog.String("slug", slug))
		return nil, nil
	}

	return &AuthorSearchResult{
		Slug: q.Authors[0].Slug,
		Name: q.Authors[0].Name,
	}, nil
}
