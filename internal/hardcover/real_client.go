package hardcover

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

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

// authorRow is the shared GraphQL selection for an author plus its canonical
// (merge target). Hardcover merges duplicate authors; when canonical is set it
// is the live author, so id/slug/name must all be resolved from it as a unit.
type authorRow struct {
	ID   int    `graphql:"id"`
	Name string `graphql:"name"`
	Slug string `graphql:"slug"`

	Canonical *struct {
		ID   int    `graphql:"id"`
		Name string `graphql:"name"`
		Slug string `graphql:"slug"`
	}
}

func (a authorRow) toResult() AuthorSearchResult {
	if a.Canonical != nil {
		return AuthorSearchResult{
			ID:   strconv.Itoa(a.Canonical.ID),
			Slug: a.Canonical.Slug,
			Name: a.Canonical.Name,
		}
	}
	return AuthorSearchResult{
		ID:   strconv.Itoa(a.ID),
		Slug: a.Slug,
		Name: a.Name,
	}
}

// SearchAuthors searches for authors by name query using the Hardcover GraphQL API.
func (c *RealClient) SearchAuthors(ctx context.Context, query string) ([]AuthorSearchResult, error) {
	slog.InfoContext(ctx, "Searching Hardcover authors", slog.String("query", query))

	var q struct {
		Authors []authorRow `graphql:"authors(where: {name: {_eq: $name}})"`
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
		searchResults[idx] = author.toResult()
	}

	return searchResults, nil
}

// GetAuthorBySlug retrieves an author by their slug using the Hardcover GraphQL API.
func (c *RealClient) GetAuthorBySlug(ctx context.Context, slug string) (*AuthorSearchResult, error) {
	slog.InfoContext(ctx, "Getting Hardcover author by slug", slog.String("slug", slug))

	var q struct {
		Authors []authorRow `graphql:"authors(where: {slug: {_eq: $slug}})"`
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

	result := q.Authors[0].toResult()
	return &result, nil
}

// GetAuthorByID retrieves an author by their canonical id using the Hardcover GraphQL API.
func (c *RealClient) GetAuthorByID(ctx context.Context, id string) (*AuthorSearchResult, error) {
	slog.InfoContext(ctx, "Getting Hardcover author by id", slog.String("id", id))

	intID, err := strconv.Atoi(id)
	if err != nil {
		slog.WarnContext(ctx, "Invalid hardcover author id", slog.String("id", id))
		return nil, nil
	}

	var q struct {
		Authors []authorRow `graphql:"authors(where: {id: {_eq: $id}})"`
	}

	variables := map[string]any{
		"id": intID,
	}

	if err := c.graphqlClient.Query(ctx, &q, variables); err != nil {
		return nil, err
	}

	if len(q.Authors) == 0 {
		slog.WarnContext(ctx, "No author found with given id", slog.String("id", id))
		return nil, nil
	}

	result := q.Authors[0].toResult()
	return &result, nil
}
