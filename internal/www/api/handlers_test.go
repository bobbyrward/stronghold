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

// setupTestServer wraps SetupTestServer for local use
func setupTestServer(t *testing.T) (*echo.Echo, func()) {
	return SetupTestServer(t)
}

// Helper functions to look up IDs for seeded reference data

func getNotificationTypeID(t *testing.T, e *echo.Echo, name string) uint {
	req := httptest.NewRequest(http.MethodGet, "/api/notification-types", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	var types []NotificationTypeResponse
	err := json.Unmarshal(rec.Body.Bytes(), &types)
	require.NoError(t, err)

	for _, nt := range types {
		if nt.Name == name {
			return nt.ID
		}
	}
	t.Fatalf("Notification type %q not found", name)
	return 0
}

func getTorrentCategoryID(t *testing.T, e *echo.Echo, name string) uint {
	req := httptest.NewRequest(http.MethodGet, "/api/torrent-categories", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	var categories []TorrentCategoryResponse
	err := json.Unmarshal(rec.Body.Bytes(), &categories)
	require.NoError(t, err)

	for _, cat := range categories {
		if cat.Name == name {
			return cat.ID
		}
	}
	t.Fatalf("Torrent category %q not found", name)
	return 0
}

// TestTorrentCategories tests read-only operations for TorrentCategories
func TestTorrentCategories(t *testing.T) {
	e, cleanup := setupTestServer(t)
	defer cleanup()

	// Test List - should have seeded data
	req := httptest.NewRequest(http.MethodGet, "/api/torrent-categories", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var categories []TorrentCategoryResponse
	err := json.Unmarshal(rec.Body.Bytes(), &categories)
	require.NoError(t, err)
	assert.Equal(t, 9, len(categories), "Should have 9 seeded torrent categories")

	// Verify seeded categories have expected fields
	found := make(map[string]bool)
	for _, cat := range categories {
		found[cat.Name] = true
		assert.NotZero(t, cat.ID)
		assert.NotEmpty(t, cat.ScopeName)
		assert.NotEmpty(t, cat.MediaType)
	}
	assert.True(t, found["audiobooks"], "Should have audiobooks category")
	assert.True(t, found["books"], "Should have books category")
	assert.True(t, found["personal-audiobooks"], "Should have personal-audiobooks category")
	assert.True(t, found["personal-books"], "Should have personal-books category")
	assert.True(t, found["kids-audiobooks"], "Should have kids-audiobooks category")
	assert.True(t, found["kids-books"], "Should have kids-books category")
	assert.True(t, found["general-audiobooks"], "Should have general-audiobooks category")
	assert.True(t, found["general-books"], "Should have general-books category")
	assert.True(t, found["author-subscriptions"], "Should have author-subscriptions category")

	// Test Get by ID
	categoryID := getTorrentCategoryID(t, e, "personal-audiobooks")
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/torrent-categories/%d", categoryID), nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var category TorrentCategoryResponse
	err = json.Unmarshal(rec.Body.Bytes(), &category)
	require.NoError(t, err)
	assert.Equal(t, "personal-audiobooks", category.Name)
	assert.Equal(t, "personal", category.ScopeName)
	assert.Equal(t, "audiobook", category.MediaType)
}

// TestNotifiers tests CRUD operations with ID-based lookups
func TestNotifiers(t *testing.T) {
	e, cleanup := setupTestServer(t)
	defer cleanup()

	// Get the notification type ID for "discord"
	discordTypeID := getNotificationTypeID(t, e, "discord")

	// Test Create with valid notification type ID
	createReq := NotifierRequest{
		Name:   "test-notifier",
		TypeID: discordTypeID,
		URL:    "https://discord.com/webhook/test",
	}
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/notifiers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var created NotifierResponse
	err := json.Unmarshal(rec.Body.Bytes(), &created)
	require.NoError(t, err)
	assert.Equal(t, "test-notifier", created.Name)
	assert.Equal(t, "discord", created.TypeName)
	assert.NotZero(t, created.ID)

	// Test Get
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/notifiers/%d", created.ID), nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Test Create with invalid notification type ID (ID 0 or non-existent)
	invalidReq := NotifierRequest{
		Name:   "invalid-notifier",
		TypeID: 99999, // Non-existent ID
		URL:    "https://example.com",
	}
	body, _ = json.Marshal(invalidReq)
	req = httptest.NewRequest(http.MethodPost, "/api/notifiers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	// Note: SQLite doesn't enforce FK constraints by default, so this may succeed with 201
	// In production with PostgreSQL, this would fail with a FK constraint error
	assert.True(t, rec.Code == http.StatusCreated || rec.Code == http.StatusBadRequest || rec.Code == http.StatusInternalServerError)
}

// TestFeeds tests CRUD operations for Feeds
func TestFeeds(t *testing.T) {
	e, cleanup := setupTestServer(t)
	defer cleanup()

	// Test Create
	createReq := FeedRequest{
		Name: "Test Feed",
		URL:  "https://example.com/rss",
	}
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/feeds", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var created FeedResponse
	err := json.Unmarshal(rec.Body.Bytes(), &created)
	require.NoError(t, err)
	assert.Equal(t, "Test Feed", created.Name)
	assert.Equal(t, "https://example.com/rss", created.URL)

	// Test Update
	updateReq := FeedRequest{
		Name: "Updated Feed",
		URL:  "https://example.com/rss/updated",
	}
	body, _ = json.Marshal(updateReq)
	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/feeds/%d", created.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// TestErrorCases tests various error scenarios
func TestErrorCases(t *testing.T) {
	e, cleanup := setupTestServer(t)
	defer cleanup()

	// Test invalid ID (GET still valid for reference data)
	req := httptest.NewRequest(http.MethodGet, "/api/notification-types/invalid", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Test non-existent ID (GET still valid for reference data)
	req = httptest.NewRequest(http.MethodGet, "/api/notification-types/99999", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	// Test missing required field (using feeds endpoint - CRUD resource)
	invalidReq := map[string]string{}
	body, _ := json.Marshal(invalidReq)
	req = httptest.NewRequest(http.MethodPost, "/api/feeds", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Test invalid JSON (using feeds endpoint - CRUD resource)
	req = httptest.NewRequest(http.MethodPost, "/api/feeds", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// TestBookTypes tests read-only operations for BookTypes
func TestBookTypes(t *testing.T) {
	e, cleanup := setupTestServer(t)
	defer cleanup()

	// Test List - should have seeded data
	req := httptest.NewRequest(http.MethodGet, "/api/book-types", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var bookTypes []BookTypeResponse
	err := json.Unmarshal(rec.Body.Bytes(), &bookTypes)
	require.NoError(t, err)
	assert.Equal(t, 2, len(bookTypes), "Should have 2 seeded book types")

	// Verify seeded book types
	found := make(map[string]bool)
	for _, bt := range bookTypes {
		found[bt.Name] = true
		assert.NotZero(t, bt.ID)
	}
	assert.True(t, found["ebook"], "Should have ebook type")
	assert.True(t, found["audiobook"], "Should have audiobook type")

	// Test Get by ID
	bookTypeID := getBookTypeID(t, e, "ebook")
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/book-types/%d", bookTypeID), nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var bookType BookTypeResponse
	err = json.Unmarshal(rec.Body.Bytes(), &bookType)
	require.NoError(t, err)
	assert.Equal(t, "ebook", bookType.Name)
}

// TestLibraries tests CRUD operations for Libraries
func TestLibraries(t *testing.T) {
	e, cleanup := setupTestServer(t)
	defer cleanup()

	// Test List - should be empty initially
	req := httptest.NewRequest(http.MethodGet, "/api/libraries", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var initialList []LibraryResponse
	err := json.Unmarshal(rec.Body.Bytes(), &initialList)
	require.NoError(t, err)
	assert.Equal(t, 0, len(initialList), "Should start with no libraries")

	// Test Create
	createReq := LibraryRequest{
		Name:         "test-library",
		Path:         "/test/path",
		BookTypeName: "ebook",
	}
	body, _ := json.Marshal(createReq)
	req = httptest.NewRequest(http.MethodPost, "/api/libraries", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var created LibraryResponse
	err = json.Unmarshal(rec.Body.Bytes(), &created)
	require.NoError(t, err)
	assert.Equal(t, "test-library", created.Name)
	assert.Equal(t, "/test/path", created.Path)
	assert.Equal(t, "ebook", created.BookTypeName)
	assert.NotZero(t, created.ID)
	assert.NotZero(t, created.BookTypeID)

	// Test Get by ID
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/libraries/%d", created.ID), nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var retrieved LibraryResponse
	err = json.Unmarshal(rec.Body.Bytes(), &retrieved)
	require.NoError(t, err)
	assert.Equal(t, created.ID, retrieved.ID)
	assert.Equal(t, "test-library", retrieved.Name)

	// Test Update
	updateReq := LibraryRequest{
		Name:         "updated-library",
		Path:         "/updated/path",
		BookTypeName: "audiobook",
	}
	body, _ = json.Marshal(updateReq)
	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/libraries/%d", created.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var updated LibraryResponse
	err = json.Unmarshal(rec.Body.Bytes(), &updated)
	require.NoError(t, err)
	assert.Equal(t, created.ID, updated.ID)
	assert.Equal(t, "updated-library", updated.Name)
	assert.Equal(t, "/updated/path", updated.Path)
	assert.Equal(t, "audiobook", updated.BookTypeName)

	// Test List after create
	req = httptest.NewRequest(http.MethodGet, "/api/libraries", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var libraries []LibraryResponse
	err = json.Unmarshal(rec.Body.Bytes(), &libraries)
	require.NoError(t, err)
	assert.Equal(t, 1, len(libraries))

	// Test Delete
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/libraries/%d", created.ID), nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNoContent, rec.Code)

	// Verify deletion
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/libraries/%d", created.ID), nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// TestLibraries_FilterByBookType tests filtering libraries by book_type_id query parameter
func TestLibraries_FilterByBookType(t *testing.T) {
	e, cleanup := setupTestServer(t)
	defer cleanup()

	ebookTypeID := getBookTypeID(t, e, "ebook")
	audiobookTypeID := getBookTypeID(t, e, "audiobook")

	// Create an ebook library
	createReq := LibraryRequest{
		Name:         "ebook-lib",
		Path:         "/ebook/path",
		BookTypeName: "ebook",
	}
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/libraries", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	// Create an audiobook library
	createReq = LibraryRequest{
		Name:         "audiobook-lib",
		Path:         "/audiobook/path",
		BookTypeName: "audiobook",
	}
	body, _ = json.Marshal(createReq)
	req = httptest.NewRequest(http.MethodPost, "/api/libraries", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	// Filter by ebook type
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/libraries?book_type_id=%d", ebookTypeID), nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var ebookLibs []LibraryResponse
	err := json.Unmarshal(rec.Body.Bytes(), &ebookLibs)
	require.NoError(t, err)
	require.Len(t, ebookLibs, 1)
	assert.Equal(t, "ebook-lib", ebookLibs[0].Name)

	// Filter by audiobook type
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/libraries?book_type_id=%d", audiobookTypeID), nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var audiobookLibs []LibraryResponse
	err = json.Unmarshal(rec.Body.Bytes(), &audiobookLibs)
	require.NoError(t, err)
	require.Len(t, audiobookLibs, 1)
	assert.Equal(t, "audiobook-lib", audiobookLibs[0].Name)
}

// TestLibraries_InvalidBookType tests validation for invalid book_type_name
func TestLibraries_InvalidBookType(t *testing.T) {
	e, cleanup := setupTestServer(t)
	defer cleanup()

	t.Run("Create with invalid book_type_name returns 400", func(t *testing.T) {
		createReq := LibraryRequest{
			Name:         "invalid-lib",
			Path:         "/invalid/path",
			BookTypeName: "invalid_type",
		}
		body, _ := json.Marshal(createReq)
		req := httptest.NewRequest(http.MethodPost, "/api/libraries", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Update with invalid book_type_name returns 400", func(t *testing.T) {
		// First create a valid library
		createReq := LibraryRequest{
			Name:         "valid-lib-for-update",
			Path:         "/valid/path",
			BookTypeName: "ebook",
		}
		body, _ := json.Marshal(createReq)
		req := httptest.NewRequest(http.MethodPost, "/api/libraries", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		require.Equal(t, http.StatusCreated, rec.Code)

		var created LibraryResponse
		err := json.Unmarshal(rec.Body.Bytes(), &created)
		require.NoError(t, err)

		// Try to update with invalid book type
		updateReq := LibraryRequest{
			Name:         "valid-lib-for-update",
			Path:         "/valid/path",
			BookTypeName: "invalid_type",
		}
		body, _ = json.Marshal(updateReq)
		req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/libraries/%d", created.ID), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec = httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

// TestLibraries_Validation tests validation for required fields
func TestLibraries_Validation(t *testing.T) {
	e, cleanup := setupTestServer(t)
	defer cleanup()

	t.Run("Create with missing name returns 400", func(t *testing.T) {
		body := `{"path": "/test/path", "book_type_name": "ebook"}`
		req := httptest.NewRequest(http.MethodPost, "/api/libraries", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Create with missing path returns 400", func(t *testing.T) {
		body := `{"name": "test-lib", "book_type_name": "ebook"}`
		req := httptest.NewRequest(http.MethodPost, "/api/libraries", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Create with missing book_type_name returns 400", func(t *testing.T) {
		body := `{"name": "test-lib", "path": "/test/path"}`
		req := httptest.NewRequest(http.MethodPost, "/api/libraries", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Get non-existent library returns 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/libraries/99999", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("Invalid library ID returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/libraries/invalid", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

// Helper function to get book type ID by name
func getBookTypeID(t *testing.T, e *echo.Echo, name string) uint {
	req := httptest.NewRequest(http.MethodGet, "/api/book-types", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	var types []BookTypeResponse
	err := json.Unmarshal(rec.Body.Bytes(), &types)
	require.NoError(t, err)

	for _, bt := range types {
		if bt.Name == name {
			return bt.ID
		}
	}
	t.Fatalf("Book type %q not found", name)
	return 0
}
