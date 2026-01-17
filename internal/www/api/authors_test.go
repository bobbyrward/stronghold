package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper to create a test author
func createTestAuthor(t *testing.T, e *echo.Echo, name string) AuthorResponse {
	body := fmt.Sprintf(`{"name": "%s"}`, name)
	req := httptest.NewRequest(http.MethodPost, "/api/authors", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	var author AuthorResponse
	err := json.Unmarshal(rec.Body.Bytes(), &author)
	require.NoError(t, err)
	return author
}

// Helper to create an alias for an author
func createTestAlias(t *testing.T, e *echo.Echo, authorID uint, name string) AuthorAliasResponse {
	body := fmt.Sprintf(`{"name": "%s"}`, name)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/authors/%d/aliases", authorID), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	var alias AuthorAliasResponse
	err := json.Unmarshal(rec.Body.Bytes(), &alias)
	require.NoError(t, err)
	return alias
}

func TestAuthors_CRUD(t *testing.T) {
	e, cleanup := SetupTestServer(t)
	defer cleanup()

	var authorID uint

	t.Run("Create author", func(t *testing.T) {
		body := `{"name": "Test Author"}`
		req := httptest.NewRequest(http.MethodPost, "/api/authors", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var author AuthorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &author)
		require.NoError(t, err)

		assert.NotZero(t, author.ID)
		assert.Equal(t, "Test Author", author.Name)
		authorID = author.ID
	})

	t.Run("List authors", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/authors", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var authors []AuthorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &authors)
		require.NoError(t, err)

		assert.NotEmpty(t, authors)
		found := false
		for _, a := range authors {
			if a.ID == authorID {
				found = true
				break
			}
		}
		assert.True(t, found, "Created author should be in list")
	})

	t.Run("Get author by ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/authors/%d", authorID), nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var author AuthorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &author)
		require.NoError(t, err)
		assert.Equal(t, authorID, author.ID)
		assert.Equal(t, "Test Author", author.Name)
	})

	t.Run("Update author", func(t *testing.T) {
		body := `{"name": "Updated Author"}`
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/authors/%d", authorID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var author AuthorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &author)
		require.NoError(t, err)
		assert.Equal(t, "Updated Author", author.Name)
	})

	t.Run("Delete author", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/authors/%d", authorID), nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("Get deleted author returns 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/authors/%d", authorID), nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestAuthors_HardcoverValidation(t *testing.T) {
	e, cleanup := SetupTestServer(t)
	defer cleanup()

	// The mock client in testing.go has these authors:
	// - "brandon-sanderson" -> "Brandon Sanderson"
	// - "brandon-mull" -> "Brandon Mull"
	// - "patrick-rothfuss" -> "Patrick Rothfuss"
	// - "joe-abercrombie" -> "Joe Abercrombie"

	t.Run("Create with valid hardcover_ref succeeds", func(t *testing.T) {
		body := `{"name": "Brandon Sanderson", "hardcover_ref": "brandon-sanderson"}`
		req := httptest.NewRequest(http.MethodPost, "/api/authors", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var author AuthorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &author)
		require.NoError(t, err)
		assert.NotNil(t, author.HardcoverRef)
		assert.Equal(t, "brandon-sanderson", *author.HardcoverRef)
	})

	t.Run("Create with invalid hardcover_ref returns 400", func(t *testing.T) {
		body := `{"name": "Invalid Author", "hardcover_ref": "nonexistent-author-slug"}`
		req := httptest.NewRequest(http.MethodPost, "/api/authors", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Create without hardcover_ref succeeds", func(t *testing.T) {
		body := `{"name": "No Hardcover Ref"}`
		req := httptest.NewRequest(http.MethodPost, "/api/authors", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var author AuthorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &author)
		require.NoError(t, err)
		assert.Nil(t, author.HardcoverRef)
	})

	t.Run("Create with empty hardcover_ref succeeds", func(t *testing.T) {
		body := `{"name": "Empty Hardcover Ref", "hardcover_ref": ""}`
		req := httptest.NewRequest(http.MethodPost, "/api/authors", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("Update with valid hardcover_ref succeeds", func(t *testing.T) {
		// First create an author without hardcover_ref
		author := createTestAuthor(t, e, "Author For Update")

		// Update with valid hardcover_ref
		body := `{"name": "Author For Update", "hardcover_ref": "patrick-rothfuss"}`
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/authors/%d", author.ID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var updated AuthorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &updated)
		require.NoError(t, err)
		assert.NotNil(t, updated.HardcoverRef)
		assert.Equal(t, "patrick-rothfuss", *updated.HardcoverRef)
	})

	t.Run("Update with invalid hardcover_ref returns 400", func(t *testing.T) {
		// First create an author
		author := createTestAuthor(t, e, "Author For Invalid Update")

		// Try to update with invalid hardcover_ref
		body := `{"name": "Author For Invalid Update", "hardcover_ref": "invalid-slug-here"}`
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/authors/%d", author.ID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestAuthors_FuzzySearch(t *testing.T) {
	e, cleanup := SetupTestServer(t)
	defer cleanup()

	// Create test author
	author := createTestAuthor(t, e, "Brandon Sanderson")

	// Create alias for the author
	createTestAlias(t, e, author.ID, "Brando Sando")

	t.Run("Search by name", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/authors?q=brandon", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var authors []AuthorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &authors)
		require.NoError(t, err)

		assert.NotEmpty(t, authors)
		found := false
		for _, a := range authors {
			if a.ID == author.ID {
				found = true
				break
			}
		}
		assert.True(t, found, "Should find author by name")
	})

	t.Run("Search is case-insensitive", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/authors?q=BRANDON", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var authors []AuthorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &authors)
		require.NoError(t, err)

		assert.NotEmpty(t, authors)
		found := false
		for _, a := range authors {
			if a.ID == author.ID {
				found = true
				break
			}
		}
		assert.True(t, found, "Should find author with case-insensitive search")
	})

	t.Run("Search by partial name", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/authors?q=sand", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var authors []AuthorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &authors)
		require.NoError(t, err)

		assert.NotEmpty(t, authors)
		found := false
		for _, a := range authors {
			if a.ID == author.ID {
				found = true
				break
			}
		}
		assert.True(t, found, "Should find author by partial name")
	})

	t.Run("Search by alias", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/authors?q=brando", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var authors []AuthorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &authors)
		require.NoError(t, err)

		found := false
		for _, a := range authors {
			if a.ID == author.ID {
				found = true
				break
			}
		}
		assert.True(t, found, "Should find author by alias")
	})

	t.Run("Search by alias case-insensitive", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/authors?q=SANDO", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var authors []AuthorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &authors)
		require.NoError(t, err)

		found := false
		for _, a := range authors {
			if a.ID == author.ID {
				found = true
				break
			}
		}
		assert.True(t, found, "Should find author by alias case-insensitive")
	})

	t.Run("Search without query returns all", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/authors", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var authors []AuthorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &authors)
		require.NoError(t, err)

		// Should have at least the author we created
		assert.NotEmpty(t, authors)
	})

	t.Run("Search with no matches returns empty list", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/authors?q=xyznonexistent123", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var authors []AuthorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &authors)
		require.NoError(t, err)

		assert.Empty(t, authors, "Should return empty list for no matches")
	})
}

func TestAuthors_Validation(t *testing.T) {
	e, cleanup := SetupTestServer(t)
	defer cleanup()

	t.Run("Create without name fails validation", func(t *testing.T) {
		body := `{}`
		req := httptest.NewRequest(http.MethodPost, "/api/authors", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Create with empty name fails validation", func(t *testing.T) {
		body := `{"name": ""}`
		req := httptest.NewRequest(http.MethodPost, "/api/authors", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Invalid JSON returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/authors", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Get nonexistent author returns 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/authors/99999", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("Update nonexistent author returns 404", func(t *testing.T) {
		body := `{"name": "Updated"}`
		req := httptest.NewRequest(http.MethodPut, "/api/authors/99999", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("Delete nonexistent author returns 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/authors/99999", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("Invalid ID format returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/authors/invalid", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestAuthors_MultipleAliasSearch(t *testing.T) {
	e, cleanup := SetupTestServer(t)
	defer cleanup()

	// Create author with multiple aliases
	author := createTestAuthor(t, e, "Robert Jordan")
	createTestAlias(t, e, author.ID, "Jim Rigney")
	createTestAlias(t, e, author.ID, "Reagan O'Neal")
	createTestAlias(t, e, author.ID, "Jackson O'Reilly")

	t.Run("Search finds author by first alias", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/authors?q=rigney", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var authors []AuthorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &authors)
		require.NoError(t, err)

		found := false
		for _, a := range authors {
			if a.ID == author.ID {
				found = true
				break
			}
		}
		assert.True(t, found, "Should find author by first alias")
	})

	t.Run("Search finds author by second alias", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/authors?q=reagan", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var authors []AuthorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &authors)
		require.NoError(t, err)

		found := false
		for _, a := range authors {
			if a.ID == author.ID {
				found = true
				break
			}
		}
		assert.True(t, found, "Should find author by second alias")
	})

	t.Run("Search finds author by third alias", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/authors?q=jackson", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var authors []AuthorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &authors)
		require.NoError(t, err)

		found := false
		for _, a := range authors {
			if a.ID == author.ID {
				found = true
				break
			}
		}
		assert.True(t, found, "Should find author by third alias")
	})

	t.Run("Search by real name still works with aliases", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/authors?q=jordan", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var authors []AuthorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &authors)
		require.NoError(t, err)

		found := false
		for _, a := range authors {
			if a.ID == author.ID {
				found = true
				break
			}
		}
		assert.True(t, found, "Should find author by real name even with aliases")
	})
}
