package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHardcoverSearch(t *testing.T) {
	e, cleanup := SetupTestServer(t)
	defer cleanup()

	t.Run("Search with results", func(t *testing.T) {
		// The mock client has pre-populated authors including "Brandon Sanderson" and "Brandon Mull"
		req := httptest.NewRequest(http.MethodGet, "/api/hardcover/authors/search?q=brandon", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var results []HardcoverAuthorSearchResponse
		err := json.Unmarshal(rec.Body.Bytes(), &results)
		require.NoError(t, err)

		// Should find the mock authors containing "brandon"
		assert.Len(t, results, 2, "Should find Brandon Sanderson and Brandon Mull")

		// Verify structure and content
		names := make(map[string]string)
		for _, r := range results {
			assert.NotEmpty(t, r.Slug)
			assert.NotEmpty(t, r.Name)
			names[r.Slug] = r.Name
		}
		assert.Equal(t, "Brandon Sanderson", names["brandon-sanderson"])
		assert.Equal(t, "Brandon Mull", names["brandon-mull"])
	})

	t.Run("Search with single result", func(t *testing.T) {
		// Search for "patrick" should only find Patrick Rothfuss
		req := httptest.NewRequest(http.MethodGet, "/api/hardcover/authors/search?q=patrick", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var results []HardcoverAuthorSearchResponse
		err := json.Unmarshal(rec.Body.Bytes(), &results)
		require.NoError(t, err)

		assert.Len(t, results, 1)
		assert.Equal(t, "patrick-rothfuss", results[0].Slug)
		assert.Equal(t, "Patrick Rothfuss", results[0].Name)
	})

	t.Run("Search is case insensitive", func(t *testing.T) {
		// Search with uppercase should still work
		req := httptest.NewRequest(http.MethodGet, "/api/hardcover/authors/search?q=BRANDON", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var results []HardcoverAuthorSearchResponse
		err := json.Unmarshal(rec.Body.Bytes(), &results)
		require.NoError(t, err)

		assert.Len(t, results, 2, "Case insensitive search should find both Brandon authors")
	})

	t.Run("Search with no results", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/hardcover/authors/search?q=nonexistentauthor12345", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var results []HardcoverAuthorSearchResponse
		err := json.Unmarshal(rec.Body.Bytes(), &results)
		require.NoError(t, err)

		assert.Empty(t, results)
	})

	t.Run("Search without query parameter returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/hardcover/authors/search", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Search with empty query parameter returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/hardcover/authors/search?q=", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Search by partial last name", func(t *testing.T) {
		// Search for "abercrombie" should find Joe Abercrombie
		req := httptest.NewRequest(http.MethodGet, "/api/hardcover/authors/search?q=abercrombie", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var results []HardcoverAuthorSearchResponse
		err := json.Unmarshal(rec.Body.Bytes(), &results)
		require.NoError(t, err)

		assert.Len(t, results, 1)
		assert.Equal(t, "joe-abercrombie", results[0].Slug)
		assert.Equal(t, "Joe Abercrombie", results[0].Name)
	})
}
