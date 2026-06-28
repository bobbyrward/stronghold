package hardcover

import (
	"context"
	"strings"
)

// MockClient is a mock implementation of the Hardcover Client interface for testing.
type MockClient struct {
	Authors             map[string]AuthorSearchResult // id -> result
	Books               map[string][]BookResult       // author id -> bibliography
	SearchAuthorsFunc   func(ctx context.Context, query string) ([]AuthorSearchResult, error)
	GetAuthorBySlugFunc func(ctx context.Context, slug string) (*AuthorSearchResult, error)
	GetAuthorByIDFunc   func(ctx context.Context, id string) (*AuthorSearchResult, error)
	GetAuthorBooksFunc  func(ctx context.Context, id string) ([]BookResult, error)
}

// Compile-time check that MockClient implements Client interface.
var _ Client = (*MockClient)(nil)

// NewMockClient creates a new MockClient with an empty Authors map.
func NewMockClient() *MockClient {
	return &MockClient{
		Authors: make(map[string]AuthorSearchResult),
		Books:   make(map[string][]BookResult),
	}
}

// AddAuthor adds an author to the mock client's internal data store.
func (m *MockClient) AddAuthor(id, slug, name string) {
	m.Authors[id] = AuthorSearchResult{ID: id, Slug: slug, Name: name}
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

// GetAuthorBySlug retrieves an author by their slug.
// If GetAuthorBySlugFunc is set, it delegates to that function.
// Otherwise, it scans the Authors map for a matching slug.
func (m *MockClient) GetAuthorBySlug(ctx context.Context, slug string) (*AuthorSearchResult, error) {
	if m.GetAuthorBySlugFunc != nil {
		return m.GetAuthorBySlugFunc(ctx, slug)
	}
	for _, author := range m.Authors {
		if author.Slug == slug {
			return &author, nil
		}
	}
	return nil, nil // Not found returns nil, not error
}

// GetAuthorByID retrieves an author by their canonical id.
// If GetAuthorByIDFunc is set, it delegates to that function.
// Otherwise, it looks up the author in the Authors map.
func (m *MockClient) GetAuthorByID(ctx context.Context, id string) (*AuthorSearchResult, error) {
	if m.GetAuthorByIDFunc != nil {
		return m.GetAuthorByIDFunc(ctx, id)
	}
	if author, ok := m.Authors[id]; ok {
		return &author, nil
	}
	return nil, nil
}

// GetAuthorBooks retrieves an author's bibliography by canonical id.
// If GetAuthorBooksFunc is set, it delegates to that function.
// Otherwise, it returns the books stored for that author id.
func (m *MockClient) GetAuthorBooks(ctx context.Context, id string) ([]BookResult, error) {
	if m.GetAuthorBooksFunc != nil {
		return m.GetAuthorBooksFunc(ctx, id)
	}
	return m.Books[id], nil
}
