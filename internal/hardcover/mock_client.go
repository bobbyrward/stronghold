package hardcover

import (
	"context"
	"strings"
)

// MockClient is a mock implementation of the Hardcover Client interface for testing.
type MockClient struct {
	Authors             map[string]AuthorSearchResult                                       // slug -> result
	SearchAuthorsFunc   func(ctx context.Context, query string) ([]AuthorSearchResult, error)
	GetAuthorBySlugFunc func(ctx context.Context, slug string) (*AuthorSearchResult, error)
}

// Compile-time check that MockClient implements Client interface.
var _ Client = (*MockClient)(nil)

// NewMockClient creates a new MockClient with an empty Authors map.
func NewMockClient() *MockClient {
	return &MockClient{
		Authors: make(map[string]AuthorSearchResult),
	}
}

// AddAuthor adds an author to the mock client's internal data store.
func (m *MockClient) AddAuthor(slug, name string) {
	m.Authors[slug] = AuthorSearchResult{Slug: slug, Name: name}
}

// SearchAuthors searches for authors by name query.
// If SearchAuthorsFunc is set, it delegates to that function.
// Otherwise, it performs a case-insensitive search of the Authors map.
func (m *MockClient) SearchAuthors(ctx context.Context, query string) ([]AuthorSearchResult, error) {
	if m.SearchAuthorsFunc != nil {
		return m.SearchAuthorsFunc(ctx, query)
	}
	// Default: search by name containing query (case-insensitive)
	var results []AuthorSearchResult
	queryLower := strings.ToLower(query)
	for _, author := range m.Authors {
		if strings.Contains(strings.ToLower(author.Name), queryLower) {
			results = append(results, author)
		}
	}
	return results, nil
}

// GetAuthorBySlug retrieves an author by their unique slug.
// If GetAuthorBySlugFunc is set, it delegates to that function.
// Otherwise, it looks up the author in the Authors map.
func (m *MockClient) GetAuthorBySlug(ctx context.Context, slug string) (*AuthorSearchResult, error) {
	if m.GetAuthorBySlugFunc != nil {
		return m.GetAuthorBySlugFunc(ctx, slug)
	}
	// Default: lookup in Authors map
	if author, ok := m.Authors[slug]; ok {
		return &author, nil
	}
	return nil, nil // Not found returns nil, not error
}
