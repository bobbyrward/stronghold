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

// Helper to create a test author for subscription tests
func createTestAuthorForSubscription(t *testing.T, e *echo.Echo, name string) AuthorResponse {
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

// Helper to create test libraries for subscription tests
func createTestLibraries(t *testing.T, e *echo.Echo, suffix string) (ebookLib LibraryResponse, audiobookLib LibraryResponse) {
	// Create ebook library
	ebookBody := fmt.Sprintf(`{"name": "test-ebook-%s", "path": "/test/ebook/%s", "book_type_name": "ebook"}`, suffix, suffix)
	req := httptest.NewRequest(http.MethodPost, "/api/libraries", bytes.NewBufferString(ebookBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	err := json.Unmarshal(rec.Body.Bytes(), &ebookLib)
	require.NoError(t, err)

	// Create audiobook library
	audiobookBody := fmt.Sprintf(`{"name": "test-audiobook-%s", "path": "/test/audiobook/%s", "book_type_name": "audiobook"}`, suffix, suffix)
	req = httptest.NewRequest(http.MethodPost, "/api/libraries", bytes.NewBufferString(audiobookBody))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	err = json.Unmarshal(rec.Body.Bytes(), &audiobookLib)
	require.NoError(t, err)

	return ebookLib, audiobookLib
}

func TestAuthorSubscription_CRUD(t *testing.T) {
	e, cleanup := SetupTestServer(t)
	defer cleanup()

	author := createTestAuthorForSubscription(t, e, "Subscription Test Author")
	ebookLib, audiobookLib := createTestLibraries(t, e, "crud")

	t.Run("Get subscription returns 404 when none exists", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/authors/%d/subscription", author.ID), nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("Create subscription", func(t *testing.T) {
		body := fmt.Sprintf(`{"scope_name": "personal", "ebook_library_name": "%s", "audiobook_library_name": "%s"}`, ebookLib.Name, audiobookLib.Name)
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/authors/%d/subscription", author.ID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var sub AuthorSubscriptionResponse
		err := json.Unmarshal(rec.Body.Bytes(), &sub)
		require.NoError(t, err)

		assert.NotZero(t, sub.ID)
		assert.Equal(t, author.ID, sub.AuthorID)
		assert.Equal(t, "personal", sub.ScopeName)
		assert.Equal(t, ebookLib.Name, sub.EbookLibraryName)
		assert.Equal(t, audiobookLib.Name, sub.AudiobookLibraryName)
	})

	t.Run("Get subscription returns data", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/authors/%d/subscription", author.ID), nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var sub AuthorSubscriptionResponse
		err := json.Unmarshal(rec.Body.Bytes(), &sub)
		require.NoError(t, err)
		assert.Equal(t, "personal", sub.ScopeName)
		assert.Equal(t, ebookLib.Name, sub.EbookLibraryName)
		assert.Equal(t, audiobookLib.Name, sub.AudiobookLibraryName)
	})

	t.Run("Update subscription", func(t *testing.T) {
		body := fmt.Sprintf(`{"scope_name": "family", "ebook_library_name": "%s", "audiobook_library_name": "%s"}`, ebookLib.Name, audiobookLib.Name)
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/authors/%d/subscription", author.ID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var sub AuthorSubscriptionResponse
		err := json.Unmarshal(rec.Body.Bytes(), &sub)
		require.NoError(t, err)
		assert.Equal(t, "family", sub.ScopeName)
	})

	t.Run("Delete subscription", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/authors/%d/subscription", author.ID), nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("Get subscription returns 404 after delete", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/authors/%d/subscription", author.ID), nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestAuthorSubscription_Conflict(t *testing.T) {
	e, cleanup := SetupTestServer(t)
	defer cleanup()

	author := createTestAuthorForSubscription(t, e, "Conflict Test Author")
	ebookLib, audiobookLib := createTestLibraries(t, e, "conflict")

	// Create first subscription
	body := fmt.Sprintf(`{"scope_name": "personal", "ebook_library_name": "%s", "audiobook_library_name": "%s"}`, ebookLib.Name, audiobookLib.Name)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/authors/%d/subscription", author.ID), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	t.Run("Creating duplicate subscription returns 409", func(t *testing.T) {
		body := fmt.Sprintf(`{"scope_name": "family", "ebook_library_name": "%s", "audiobook_library_name": "%s"}`, ebookLib.Name, audiobookLib.Name)
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/authors/%d/subscription", author.ID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
	})
}

func TestAuthorSubscription_InvalidScope(t *testing.T) {
	e, cleanup := SetupTestServer(t)
	defer cleanup()

	author := createTestAuthorForSubscription(t, e, "Invalid Scope Author")
	ebookLib, audiobookLib := createTestLibraries(t, e, "invalidscope")

	t.Run("Create with invalid scope returns 400", func(t *testing.T) {
		body := fmt.Sprintf(`{"scope_name": "invalid_scope", "ebook_library_name": "%s", "audiobook_library_name": "%s"}`, ebookLib.Name, audiobookLib.Name)
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/authors/%d/subscription", author.ID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestAuthorSubscription_InvalidLibrary(t *testing.T) {
	e, cleanup := SetupTestServer(t)
	defer cleanup()

	author := createTestAuthorForSubscription(t, e, "Invalid Library Author")
	ebookLib, audiobookLib := createTestLibraries(t, e, "invalidlib")

	t.Run("Create with invalid ebook_library_name returns 400", func(t *testing.T) {
		body := fmt.Sprintf(`{"scope_name": "personal", "ebook_library_name": "nonexistent-library", "audiobook_library_name": "%s"}`, audiobookLib.Name)
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/authors/%d/subscription", author.ID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Create with invalid audiobook_library_name returns 400", func(t *testing.T) {
		body := fmt.Sprintf(`{"scope_name": "personal", "ebook_library_name": "%s", "audiobook_library_name": "nonexistent-library"}`, ebookLib.Name)
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/authors/%d/subscription", author.ID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Create with both invalid library names returns 400", func(t *testing.T) {
		body := `{"scope_name": "personal", "ebook_library_name": "nonexistent-ebook", "audiobook_library_name": "nonexistent-audiobook"}`
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/authors/%d/subscription", author.ID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestAuthorSubscription_UpdateInvalidLibrary(t *testing.T) {
	e, cleanup := SetupTestServer(t)
	defer cleanup()

	author := createTestAuthorForSubscription(t, e, "Update Invalid Library Author")
	ebookLib, audiobookLib := createTestLibraries(t, e, "updateinvalidlib")

	// Create subscription first
	body := fmt.Sprintf(`{"scope_name": "personal", "ebook_library_name": "%s", "audiobook_library_name": "%s"}`, ebookLib.Name, audiobookLib.Name)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/authors/%d/subscription", author.ID), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	t.Run("Update with invalid ebook_library_name returns 400", func(t *testing.T) {
		body := fmt.Sprintf(`{"scope_name": "personal", "ebook_library_name": "nonexistent-library", "audiobook_library_name": "%s"}`, audiobookLib.Name)
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/authors/%d/subscription", author.ID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Update with invalid audiobook_library_name returns 400", func(t *testing.T) {
		body := fmt.Sprintf(`{"scope_name": "personal", "ebook_library_name": "%s", "audiobook_library_name": "nonexistent-library"}`, ebookLib.Name)
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/authors/%d/subscription", author.ID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestAuthorSubscription_InvalidAuthor(t *testing.T) {
	e, cleanup := SetupTestServer(t)
	defer cleanup()

	ebookLib, audiobookLib := createTestLibraries(t, e, "invalidauthor")

	t.Run("Create subscription for non-existent author returns 404", func(t *testing.T) {
		body := fmt.Sprintf(`{"scope_name": "personal", "ebook_library_name": "%s", "audiobook_library_name": "%s"}`, ebookLib.Name, audiobookLib.Name)
		req := httptest.NewRequest(http.MethodPost, "/api/authors/99999/subscription", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("Get subscription for non-existent author returns 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/authors/99999/subscription", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("Update subscription for non-existent author returns 404", func(t *testing.T) {
		body := fmt.Sprintf(`{"scope_name": "family", "ebook_library_name": "%s", "audiobook_library_name": "%s"}`, ebookLib.Name, audiobookLib.Name)
		req := httptest.NewRequest(http.MethodPut, "/api/authors/99999/subscription", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("Delete subscription for non-existent author returns 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/authors/99999/subscription", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestAuthorSubscription_InvalidAuthorID(t *testing.T) {
	e, cleanup := SetupTestServer(t)
	defer cleanup()

	ebookLib, audiobookLib := createTestLibraries(t, e, "invalidid")

	t.Run("Invalid author_id returns 400", func(t *testing.T) {
		body := fmt.Sprintf(`{"scope_name": "personal", "ebook_library_name": "%s", "audiobook_library_name": "%s"}`, ebookLib.Name, audiobookLib.Name)
		req := httptest.NewRequest(http.MethodPost, "/api/authors/invalid/subscription", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestAuthorSubscriptionItems_List(t *testing.T) {
	e, cleanup := SetupTestServer(t)
	defer cleanup()

	author := createTestAuthorForSubscription(t, e, "Items Test Author")
	ebookLib, audiobookLib := createTestLibraries(t, e, "items")

	t.Run("List items returns 404 when no subscription", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/authors/%d/subscription/items", author.ID), nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	// Create subscription
	body := fmt.Sprintf(`{"scope_name": "personal", "ebook_library_name": "%s", "audiobook_library_name": "%s"}`, ebookLib.Name, audiobookLib.Name)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/authors/%d/subscription", author.ID), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	t.Run("List items returns empty array when subscription exists", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/authors/%d/subscription/items", author.ID), nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var items []AuthorSubscriptionItemResponse
		err := json.Unmarshal(rec.Body.Bytes(), &items)
		require.NoError(t, err)
		assert.Empty(t, items)
	})
}

func TestAuthorSubscriptionItems_InvalidAuthor(t *testing.T) {
	e, cleanup := SetupTestServer(t)
	defer cleanup()

	t.Run("List items for non-existent author returns 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/authors/99999/subscription/items", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("Invalid author_id returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/authors/invalid/subscription/items", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestAuthorSubscription_UpdateNonExistent(t *testing.T) {
	e, cleanup := SetupTestServer(t)
	defer cleanup()

	author := createTestAuthorForSubscription(t, e, "Update Non-Existent Author")
	ebookLib, audiobookLib := createTestLibraries(t, e, "updatenonexistent")

	t.Run("Update non-existent subscription returns 404", func(t *testing.T) {
		body := fmt.Sprintf(`{"scope_name": "family", "ebook_library_name": "%s", "audiobook_library_name": "%s"}`, ebookLib.Name, audiobookLib.Name)
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/authors/%d/subscription", author.ID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestAuthorSubscription_DeleteNonExistent(t *testing.T) {
	e, cleanup := SetupTestServer(t)
	defer cleanup()

	author := createTestAuthorForSubscription(t, e, "Delete Non-Existent Author")

	t.Run("Delete non-existent subscription returns 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/authors/%d/subscription", author.ID), nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestAuthorSubscription_WithNotifier(t *testing.T) {
	e, cleanup := SetupTestServer(t)
	defer cleanup()

	author := createTestAuthorForSubscription(t, e, "Notifier Test Author")
	ebookLib, audiobookLib := createTestLibraries(t, e, "notifier")

	// Create a notifier
	discordTypeID := getNotificationTypeID(t, e, "discord")
	notifierReq := NotifierRequest{
		Name:   "subscription-test-notifier",
		TypeID: discordTypeID,
		URL:    "https://discord.com/webhook/test",
	}
	notifierBody, _ := json.Marshal(notifierReq)
	req := httptest.NewRequest(http.MethodPost, "/api/notifiers", bytes.NewReader(notifierBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	var createdNotifier NotifierResponse
	err := json.Unmarshal(rec.Body.Bytes(), &createdNotifier)
	require.NoError(t, err)

	t.Run("Create subscription with notifier", func(t *testing.T) {
		body := fmt.Sprintf(`{"scope_name": "personal", "notifier_id": %d, "ebook_library_name": "%s", "audiobook_library_name": "%s"}`, createdNotifier.ID, ebookLib.Name, audiobookLib.Name)
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/authors/%d/subscription", author.ID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var sub AuthorSubscriptionResponse
		err := json.Unmarshal(rec.Body.Bytes(), &sub)
		require.NoError(t, err)

		assert.NotZero(t, sub.ID)
		assert.Equal(t, author.ID, sub.AuthorID)
		assert.Equal(t, "personal", sub.ScopeName)
		require.NotNil(t, sub.NotifierID)
		assert.Equal(t, createdNotifier.ID, *sub.NotifierID)
		require.NotNil(t, sub.NotifierName)
		assert.Equal(t, "subscription-test-notifier", *sub.NotifierName)
		assert.Equal(t, ebookLib.Name, sub.EbookLibraryName)
		assert.Equal(t, audiobookLib.Name, sub.AudiobookLibraryName)
	})

	t.Run("Get subscription with notifier", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/authors/%d/subscription", author.ID), nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var sub AuthorSubscriptionResponse
		err := json.Unmarshal(rec.Body.Bytes(), &sub)
		require.NoError(t, err)
		assert.Equal(t, "personal", sub.ScopeName)
		require.NotNil(t, sub.NotifierName)
		assert.Equal(t, "subscription-test-notifier", *sub.NotifierName)
	})
}

func TestAuthorSubscription_UpdateScope(t *testing.T) {
	e, cleanup := SetupTestServer(t)
	defer cleanup()

	author := createTestAuthorForSubscription(t, e, "Update Scope Author")
	ebookLib, audiobookLib := createTestLibraries(t, e, "updatescope")

	// Create subscription
	body := fmt.Sprintf(`{"scope_name": "personal", "ebook_library_name": "%s", "audiobook_library_name": "%s"}`, ebookLib.Name, audiobookLib.Name)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/authors/%d/subscription", author.ID), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	t.Run("Update with invalid scope returns 400", func(t *testing.T) {
		body := fmt.Sprintf(`{"scope_name": "invalid_scope", "ebook_library_name": "%s", "audiobook_library_name": "%s"}`, ebookLib.Name, audiobookLib.Name)
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/authors/%d/subscription", author.ID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Update to different valid scope succeeds", func(t *testing.T) {
		body := fmt.Sprintf(`{"scope_name": "kids", "ebook_library_name": "%s", "audiobook_library_name": "%s"}`, ebookLib.Name, audiobookLib.Name)
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/authors/%d/subscription", author.ID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var sub AuthorSubscriptionResponse
		err := json.Unmarshal(rec.Body.Bytes(), &sub)
		require.NoError(t, err)
		assert.Equal(t, "kids", sub.ScopeName)
	})
}
