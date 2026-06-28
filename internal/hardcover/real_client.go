package hardcover

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/hasura/go-graphql-client"
	"golang.org/x/time/rate"
)

// Hardcover allows 60 req/min per account. requestsPerMinute keeps headroom under
// that cap; burst lets interactive single calls fire without waiting.
const (
	requestsPerMinute = 50
	burst             = 5
)

// rateLimitedTransport blocks each request until the shared limiter grants a
// token, so every Hardcover call site (sync, web author-add, doctor) shares one
// budget instead of pacing independently.
type rateLimitedTransport struct {
	limiter *rate.Limiter
	base    http.RoundTripper
}

func (t *rateLimitedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := t.limiter.Wait(req.Context()); err != nil {
		return nil, err
	}
	return t.base.RoundTrip(req)
}

// RealClient is the production implementation of the Hardcover API client.
type RealClient struct {
	graphqlClient *graphql.Client
	token         string
}

// Compile-time check that RealClient implements Client interface.
var _ Client = (*RealClient)(nil)

// NewClient creates a new Hardcover API client with the given authentication token.
func NewClient(token string) *RealClient {
	// ponytail: per-process limiter — a running web server and a separately
	// invoked doctor CLI each get their own and don't coordinate. Cross-process
	// (distributed) limiting is out of scope; not a normal workflow.
	httpClient := &http.Client{
		Transport: &rateLimitedTransport{
			limiter: rate.NewLimiter(rate.Every(time.Minute/requestsPerMinute), burst),
			base:    http.DefaultTransport,
		},
	}

	return &RealClient{
		graphqlClient: graphql.NewClient("https://api.hardcover.app/v1/graphql", httpClient).
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

// bookRow is the GraphQL selection for one work in an author's bibliography.
type bookRow struct {
	ID          int     `graphql:"id"`
	Title       string  `graphql:"title"`
	ReleaseDate *string `graphql:"release_date"`
}

func (b bookRow) toResult() BookResult {
	return BookResult{
		HardcoverID: strconv.Itoa(b.ID),
		Title:       b.Title,
		ReleaseDate: b.ReleaseDate,
	}
}

// GetAuthorBooks fetches an author's bibliography by canonical id, ordered by
// release date. See docs/hardcover-api.md query #2.
func (c *RealClient) GetAuthorBooks(ctx context.Context, id string) ([]BookResult, error) {
	slog.InfoContext(ctx, "Getting Hardcover author bibliography", slog.String("id", id))

	intID, err := strconv.Atoi(id)
	if err != nil {
		slog.WarnContext(ctx, "Invalid hardcover author id", slog.String("id", id))
		return nil, nil
	}

	var q struct {
		Books []bookRow `graphql:"books(where: {contributions: {author_id: {_eq: $authorId}}}, order_by: {release_date: asc})"`
	}

	variables := map[string]any{
		"authorId": intID,
	}

	if err := c.graphqlClient.Query(ctx, &q, variables); err != nil {
		return nil, err
	}

	results := make([]BookResult, len(q.Books))
	for idx, book := range q.Books {
		results[idx] = book.toResult()
	}

	slog.DebugContext(ctx, "Fetched author bibliography", slog.String("id", id), slog.Int("count", len(results)))
	return results, nil
}
