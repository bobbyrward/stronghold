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

// createTestAuthorForAlias creates an author for use in alias tests
func createTestAuthorForAlias(t *testing.T, e *echo.Echo, name string) AuthorResponse {
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

func TestAuthorAliases_CRUD(t *testing.T) {
	e, cleanup := SetupTestServer(t)
	defer cleanup()

	// Create author first
	author := createTestAuthorForAlias(t, e, "Test Author For Aliases")
	var aliasID uint

	t.Run("List aliases initially empty", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/authors/%d/aliases", author.ID), nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var aliases []AuthorAliasResponse
		err := json.Unmarshal(rec.Body.Bytes(), &aliases)
		require.NoError(t, err)
		assert.Empty(t, aliases)
	})

	t.Run("Create alias", func(t *testing.T) {
		body := `{"name": "Test Alias"}`
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/authors/%d/aliases", author.ID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var alias AuthorAliasResponse
		err := json.Unmarshal(rec.Body.Bytes(), &alias)
		require.NoError(t, err)

		assert.NotZero(t, alias.ID)
		assert.Equal(t, author.ID, alias.AuthorID)
		assert.Equal(t, "Test Alias", alias.Name)
		aliasID = alias.ID
	})

	t.Run("Get alias by ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/authors/%d/aliases/%d", author.ID, aliasID), nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var alias AuthorAliasResponse
		err := json.Unmarshal(rec.Body.Bytes(), &alias)
		require.NoError(t, err)
		assert.Equal(t, aliasID, alias.ID)
	})

	t.Run("Update alias", func(t *testing.T) {
		body := `{"name": "Updated Alias"}`
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/authors/%d/aliases/%d", author.ID, aliasID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var alias AuthorAliasResponse
		err := json.Unmarshal(rec.Body.Bytes(), &alias)
		require.NoError(t, err)
		assert.Equal(t, "Updated Alias", alias.Name)
	})

	t.Run("Delete alias", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/authors/%d/aliases/%d", author.ID, aliasID), nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("List shows alias removed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/authors/%d/aliases", author.ID), nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var aliases []AuthorAliasResponse
		err := json.Unmarshal(rec.Body.Bytes(), &aliases)
		require.NoError(t, err)
		assert.Empty(t, aliases)
	})
}

func TestAuthorAliases_AuthorNotFound(t *testing.T) {
	e, cleanup := SetupTestServer(t)
	defer cleanup()

	t.Run("List aliases for non-existent author returns 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/authors/99999/aliases", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("Create alias for non-existent author returns 404", func(t *testing.T) {
		body := `{"name": "Test Alias"}`
		req := httptest.NewRequest(http.MethodPost, "/api/authors/99999/aliases", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("Get alias for non-existent author returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/authors/99999/aliases/1", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		// GetAuthorAlias doesn't check if author exists first, it just checks if alias exists
		// and if alias.AuthorID matches. Since there's no alias with ID 1 belonging to author 99999,
		// it returns 404 for the alias not found
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("Update alias for non-existent author returns 404", func(t *testing.T) {
		body := `{"name": "Test Alias"}`
		req := httptest.NewRequest(http.MethodPut, "/api/authors/99999/aliases/1", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("Delete alias for non-existent author returns 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/authors/99999/aliases/1", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestAuthorAliases_UniqueConstraint(t *testing.T) {
	e, cleanup := SetupTestServer(t)
	defer cleanup()

	// Create two authors
	author1 := createTestAuthorForAlias(t, e, "Author One")
	author2 := createTestAuthorForAlias(t, e, "Author Two")

	// Create alias for author1
	body := `{"name": "Unique Alias Name"}`
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/authors/%d/aliases", author1.ID), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	t.Run("Creating duplicate alias name for different author fails", func(t *testing.T) {
		body := `{"name": "Unique Alias Name"}`
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/authors/%d/aliases", author2.ID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		// Should fail due to global uniqueness constraint
		assert.Equal(t, http.StatusInternalServerError, rec.Code) // Constraint violation becomes 500
	})

	t.Run("Creating duplicate alias name for same author fails", func(t *testing.T) {
		body := `{"name": "Unique Alias Name"}`
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/authors/%d/aliases", author1.ID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		// Should fail due to uniqueness constraint
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestAuthorAliases_AliasBelongsToAuthor(t *testing.T) {
	e, cleanup := SetupTestServer(t)
	defer cleanup()

	// Create two authors
	author1 := createTestAuthorForAlias(t, e, "Author With Alias")
	author2 := createTestAuthorForAlias(t, e, "Other Author")

	// Create alias for author1
	body := `{"name": "Author1 Alias"}`
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/authors/%d/aliases", author1.ID), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	var alias AuthorAliasResponse
	err := json.Unmarshal(rec.Body.Bytes(), &alias)
	require.NoError(t, err)

	t.Run("Get alias via wrong author returns 404", func(t *testing.T) {
		// Try to get author1's alias via author2's URL
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/authors/%d/aliases/%d", author2.ID, alias.ID), nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("Update alias via wrong author returns 404", func(t *testing.T) {
		body := `{"name": "Trying to Update"}`
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/authors/%d/aliases/%d", author2.ID, alias.ID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("Delete alias via wrong author returns 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/authors/%d/aliases/%d", author2.ID, alias.ID), nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("Alias still accessible via correct author", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/authors/%d/aliases/%d", author1.ID, alias.ID), nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var retrievedAlias AuthorAliasResponse
		err := json.Unmarshal(rec.Body.Bytes(), &retrievedAlias)
		require.NoError(t, err)
		assert.Equal(t, alias.ID, retrievedAlias.ID)
		assert.Equal(t, author1.ID, retrievedAlias.AuthorID)
	})
}

func TestAuthorAliases_ValidationErrors(t *testing.T) {
	e, cleanup := SetupTestServer(t)
	defer cleanup()

	author := createTestAuthorForAlias(t, e, "Author For Validation")

	t.Run("Create alias with empty name fails", func(t *testing.T) {
		body := `{"name": ""}`
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/authors/%d/aliases", author.ID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Create alias with missing name fails", func(t *testing.T) {
		body := `{}`
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/authors/%d/aliases", author.ID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Create alias with invalid JSON fails", func(t *testing.T) {
		body := `{invalid json`
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/authors/%d/aliases", author.ID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Invalid author_id format returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/authors/invalid/aliases", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Invalid alias id format returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/authors/%d/aliases/invalid", author.ID), nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestAuthorAliases_MultipleAliases(t *testing.T) {
	e, cleanup := SetupTestServer(t)
	defer cleanup()

	author := createTestAuthorForAlias(t, e, "Author With Multiple Aliases")

	// Create multiple aliases
	aliasNames := []string{"Alias One", "Alias Two", "Alias Three"}
	var createdAliases []AuthorAliasResponse

	for _, name := range aliasNames {
		body := fmt.Sprintf(`{"name": "%s"}`, name)
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/authors/%d/aliases", author.ID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		require.Equal(t, http.StatusCreated, rec.Code)

		var alias AuthorAliasResponse
		err := json.Unmarshal(rec.Body.Bytes(), &alias)
		require.NoError(t, err)
		createdAliases = append(createdAliases, alias)
	}

	t.Run("List returns all aliases", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/authors/%d/aliases", author.ID), nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var aliases []AuthorAliasResponse
		err := json.Unmarshal(rec.Body.Bytes(), &aliases)
		require.NoError(t, err)
		assert.Len(t, aliases, 3)

		// Verify all expected aliases are present
		foundNames := make(map[string]bool)
		for _, a := range aliases {
			foundNames[a.Name] = true
			assert.Equal(t, author.ID, a.AuthorID)
		}
		for _, name := range aliasNames {
			assert.True(t, foundNames[name], "Expected alias %q to be present", name)
		}
	})

	t.Run("Can get each alias individually", func(t *testing.T) {
		for _, created := range createdAliases {
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/authors/%d/aliases/%d", author.ID, created.ID), nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)

			var alias AuthorAliasResponse
			err := json.Unmarshal(rec.Body.Bytes(), &alias)
			require.NoError(t, err)
			assert.Equal(t, created.ID, alias.ID)
			assert.Equal(t, created.Name, alias.Name)
		}
	})
}
