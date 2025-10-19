package handlers

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

	"github.com/bobbyrward/stronghold/internal/models"
)

// setupTestServer creates a test Echo server with a test database
func setupTestServer(t *testing.T) (*echo.Echo, func()) {
	db, err := models.ConnectTestDB()
	require.NoError(t, err, "Failed to connect to test database")

	e := echo.New()

	// Register filter key routes
	e.GET("/filter-keys", ListFilterKeys(db))
	e.POST("/filter-keys", CreateFilterKey(db))
	e.GET("/filter-keys/:id", GetFilterKey(db))
	e.PUT("/filter-keys/:id", UpdateFilterKey(db))
	e.DELETE("/filter-keys/:id", DeleteFilterKey(db))

	// Register filter operator routes
	e.GET("/filter-operators", ListFilterOperators(db))
	e.POST("/filter-operators", CreateFilterOperator(db))
	e.GET("/filter-operators/:id", GetFilterOperator(db))
	e.PUT("/filter-operators/:id", UpdateFilterOperator(db))
	e.DELETE("/filter-operators/:id", DeleteFilterOperator(db))

	// Register notification type routes
	e.GET("/notification-types", ListNotificationTypes(db))
	e.POST("/notification-types", CreateNotificationType(db))
	e.GET("/notification-types/:id", GetNotificationType(db))
	e.PUT("/notification-types/:id", UpdateNotificationType(db))
	e.DELETE("/notification-types/:id", DeleteNotificationType(db))

	// Register torrent category routes
	e.GET("/torrent-categories", ListTorrentCategories(db))
	e.POST("/torrent-categories", CreateTorrentCategory(db))
	e.GET("/torrent-categories/:id", GetTorrentCategory(db))
	e.PUT("/torrent-categories/:id", UpdateTorrentCategory(db))
	e.DELETE("/torrent-categories/:id", DeleteTorrentCategory(db))

	// Register notifier routes
	e.GET("/notifiers", ListNotifiers(db))
	e.POST("/notifiers", CreateNotifier(db))
	e.GET("/notifiers/:id", GetNotifier(db))
	e.PUT("/notifiers/:id", UpdateNotifier(db))
	e.DELETE("/notifiers/:id", DeleteNotifier(db))

	// Register feed routes
	e.GET("/feeds", ListFeeds(db))
	e.POST("/feeds", CreateFeed(db))
	e.GET("/feeds/:id", GetFeed(db))
	e.PUT("/feeds/:id", UpdateFeed(db))
	e.DELETE("/feeds/:id", DeleteFeed(db))

	// Register feed filter routes
	e.GET("/feed-filters", ListFeedFilters(db))
	e.POST("/feed-filters", CreateFeedFilter(db))
	e.GET("/feed-filters/:id", GetFeedFilter(db))
	e.PUT("/feed-filters/:id", UpdateFeedFilter(db))
	e.DELETE("/feed-filters/:id", DeleteFeedFilter(db))

	// Register feed author filter routes
	e.GET("/feed-author-filters", ListFeedAuthorFilters(db))
	e.POST("/feed-author-filters", CreateFeedAuthorFilter(db))
	e.GET("/feed-author-filters/:id", GetFeedAuthorFilter(db))
	e.PUT("/feed-author-filters/:id", UpdateFeedAuthorFilter(db))
	e.DELETE("/feed-author-filters/:id", DeleteFeedAuthorFilter(db))

	// Register feed filter set routes
	e.GET("/feed-filter-sets", ListFeedFilterSets(db))
	e.POST("/feed-filter-sets", CreateFeedFilterSet(db))
	e.GET("/feed-filter-sets/:id", GetFeedFilterSet(db))
	e.PUT("/feed-filter-sets/:id", UpdateFeedFilterSet(db))
	e.DELETE("/feed-filter-sets/:id", DeleteFeedFilterSet(db))

	// Register feed filter set entry routes
	e.GET("/feed-filter-set-entries", ListFeedFilterSetEntries(db))
	e.POST("/feed-filter-set-entries", CreateFeedFilterSetEntry(db))
	e.GET("/feed-filter-set-entries/:id", GetFeedFilterSetEntry(db))
	e.PUT("/feed-filter-set-entries/:id", UpdateFeedFilterSetEntry(db))
	e.DELETE("/feed-filter-set-entries/:id", DeleteFeedFilterSetEntry(db))

	cleanup := func() {
		// No cleanup needed for in-memory database
	}

	return e, cleanup
}

// TestFilterKeys tests CRUD operations for FilterKeys
func TestFilterKeys(t *testing.T) {
	e, cleanup := setupTestServer(t)
	defer cleanup()

	// Test List - should have seeded data
	req := httptest.NewRequest(http.MethodGet, "/filter-keys", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var filterKeys []FilterKeyResponse
	err := json.Unmarshal(rec.Body.Bytes(), &filterKeys)
	require.NoError(t, err)
	assert.Greater(t, len(filterKeys), 0, "Should have seeded filter keys")

	// Test Create
	createReq := FilterKeyRequest{Name: "test-key"}
	body, _ := json.Marshal(createReq)
	req = httptest.NewRequest(http.MethodPost, "/filter-keys", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var created FilterKeyResponse
	err = json.Unmarshal(rec.Body.Bytes(), &created)
	require.NoError(t, err)
	assert.Equal(t, "test-key", created.Name)
	assert.NotZero(t, created.ID)

	// Test Get by ID
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/filter-keys/%d", created.ID), nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var retrieved FilterKeyResponse
	err = json.Unmarshal(rec.Body.Bytes(), &retrieved)
	require.NoError(t, err)
	assert.Equal(t, created.ID, retrieved.ID)
	assert.Equal(t, created.Name, retrieved.Name)

	// Test Update
	updateReq := FilterKeyRequest{Name: "updated-key"}
	body, _ = json.Marshal(updateReq)
	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/filter-keys/%d", created.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var updated FilterKeyResponse
	err = json.Unmarshal(rec.Body.Bytes(), &updated)
	require.NoError(t, err)
	assert.Equal(t, "updated-key", updated.Name)

	// Test Delete
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/filter-keys/%d", created.ID), nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNoContent, rec.Code)

	// Verify deletion
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/filter-keys/%d", created.ID), nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// TestFilterOperators tests CRUD operations for FilterOperators
func TestFilterOperators(t *testing.T) {
	e, cleanup := setupTestServer(t)
	defer cleanup()

	// Test List - should have seeded data
	req := httptest.NewRequest(http.MethodGet, "/filter-operators", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var filterOperators []FilterOperatorResponse
	err := json.Unmarshal(rec.Body.Bytes(), &filterOperators)
	require.NoError(t, err)
	assert.Greater(t, len(filterOperators), 0, "Should have seeded filter operators")

	// Test Create
	createReq := FilterOperatorRequest{Name: "test-operator"}
	body, _ := json.Marshal(createReq)
	req = httptest.NewRequest(http.MethodPost, "/filter-operators", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var created FilterOperatorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &created)
	require.NoError(t, err)
	assert.Equal(t, "test-operator", created.Name)
	assert.NotZero(t, created.ID)
}

// TestTorrentCategories tests CRUD operations for TorrentCategories
func TestTorrentCategories(t *testing.T) {
	e, cleanup := setupTestServer(t)
	defer cleanup()

	// Test Create
	createReq := TorrentCategoryRequest{Name: "test-category"}
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/torrent-categories", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var created TorrentCategoryResponse
	err := json.Unmarshal(rec.Body.Bytes(), &created)
	require.NoError(t, err)
	assert.Equal(t, "test-category", created.Name)

	// Test List
	req = httptest.NewRequest(http.MethodGet, "/torrent-categories", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var categories []TorrentCategoryResponse
	err = json.Unmarshal(rec.Body.Bytes(), &categories)
	require.NoError(t, err)
	assert.Greater(t, len(categories), 0)
}

// TestNotifiers tests CRUD operations with name-based lookups
func TestNotifiers(t *testing.T) {
	e, cleanup := setupTestServer(t)
	defer cleanup()

	// Test Create with valid notification type name
	createReq := NotifierRequest{
		Name:     "test-notifier",
		TypeName: "discord", // From seeded data
		URL:      "https://discord.com/webhook/test",
	}
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/notifiers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var created NotifierResponse
	err := json.Unmarshal(rec.Body.Bytes(), &created)
	require.NoError(t, err)
	assert.Equal(t, "test-notifier", created.Name)
	assert.Equal(t, "discord", created.Type)
	assert.NotZero(t, created.ID)

	// Test Get
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/notifiers/%d", created.ID), nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Test Create with invalid notification type name
	invalidReq := NotifierRequest{
		Name:     "invalid-notifier",
		TypeName: "nonexistent",
		URL:      "https://example.com",
	}
	body, _ = json.Marshal(invalidReq)
	req = httptest.NewRequest(http.MethodPost, "/notifiers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
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
	req := httptest.NewRequest(http.MethodPost, "/feeds", bytes.NewReader(body))
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
	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/feeds/%d", created.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// TestCompleteWorkflow tests a complete workflow with cascading relationships
func TestCompleteWorkflow(t *testing.T) {
	e, cleanup := setupTestServer(t)
	defer cleanup()

	// Step 1: Create a TorrentCategory
	categoryReq := TorrentCategoryRequest{Name: "test-books"}
	body, _ := json.Marshal(categoryReq)
	req := httptest.NewRequest(http.MethodPost, "/torrent-categories", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	// Step 2: Create a Notifier
	notifierReq := NotifierRequest{
		Name:     "discord-personal",
		TypeName: "discord",
		URL:      "https://discord.com/webhook/test",
	}
	body, _ = json.Marshal(notifierReq)
	req = httptest.NewRequest(http.MethodPost, "/notifiers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	// Step 3: Create a Feed
	feedReq := FeedRequest{
		Name: "MAM",
		URL:  "https://example.com/rss",
	}
	body, _ = json.Marshal(feedReq)
	req = httptest.NewRequest(http.MethodPost, "/feeds", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	// Step 4: Create a FeedFilter using name lookups
	filterReq := FeedFilterRequest{
		Name:         "Blaise Corvin Books",
		FeedName:     "MAM",
		CategoryName: "test-books",
		NotifierName: "discord-personal",
	}
	body, _ = json.Marshal(filterReq)
	req = httptest.NewRequest(http.MethodPost, "/feed-filters", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var createdFilter FeedFilterResponse
	err := json.Unmarshal(rec.Body.Bytes(), &createdFilter)
	require.NoError(t, err)
	assert.Equal(t, "Blaise Corvin Books", createdFilter.Name)
	assert.Equal(t, "MAM", createdFilter.FeedName)
	assert.NotZero(t, createdFilter.FeedID)
	assert.Equal(t, "test-books", createdFilter.Category)
	assert.Equal(t, "discord-personal", createdFilter.Notifier)

	// Step 5: Create a FeedFilterSet
	setReq := FeedFilterSetRequest{
		FeedFilterID: createdFilter.ID,
		TypeName:     "any", // Using seeded type
	}
	body, _ = json.Marshal(setReq)
	req = httptest.NewRequest(http.MethodPost, "/feed-filter-sets", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var createdSet FeedFilterSetResponse
	err = json.Unmarshal(rec.Body.Bytes(), &createdSet)
	require.NoError(t, err)
	assert.Equal(t, "any", createdSet.Type)
	require.NotZero(t, createdSet.ID, "FeedFilterSet ID should not be zero")

	// Step 6: Create FeedFilterSetEntry
	entryReq := FeedFilterSetEntryRequest{
		FeedFilterSetID: createdSet.ID,
		KeyName:         "author",
		OperatorName:    "contains",
		Value:           "Blaise Corvin",
	}
	body, _ = json.Marshal(entryReq)
	req = httptest.NewRequest(http.MethodPost, "/feed-filter-set-entries", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var createdEntry FeedFilterSetEntryResponse
	err = json.Unmarshal(rec.Body.Bytes(), &createdEntry)
	require.NoError(t, err)
	assert.Equal(t, "author", createdEntry.Key)
	assert.Equal(t, "contains", createdEntry.Operator)
	assert.Equal(t, "Blaise Corvin", createdEntry.Value)

	// Step 7: Test query parameter filtering
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/feed-filter-set-entries?feed_filter_set_id=%d", createdSet.ID), nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var entries []FeedFilterSetEntryResponse
	err = json.Unmarshal(rec.Body.Bytes(), &entries)
	require.NoError(t, err)
	require.Len(t, entries, 1, "Should have exactly 1 entry")
	assert.Equal(t, createdEntry.ID, entries[0].ID)
}

// TestErrorCases tests various error scenarios
func TestErrorCases(t *testing.T) {
	e, cleanup := setupTestServer(t)
	defer cleanup()

	// Test invalid ID
	req := httptest.NewRequest(http.MethodGet, "/filter-keys/invalid", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Test non-existent ID
	req = httptest.NewRequest(http.MethodGet, "/filter-keys/99999", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	// Test missing required field
	invalidReq := map[string]string{}
	body, _ := json.Marshal(invalidReq)
	req = httptest.NewRequest(http.MethodPost, "/filter-keys", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Test invalid JSON
	req = httptest.NewRequest(http.MethodPost, "/filter-keys", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// TestFeedFilterValidation tests feed filter validation including feed_name
func TestFeedFilterValidation(t *testing.T) {
	e, cleanup := setupTestServer(t)
	defer cleanup()

	// Setup: Create required entities
	categoryReq := TorrentCategoryRequest{Name: "validation-test-category"}
	body, _ := json.Marshal(categoryReq)
	req := httptest.NewRequest(http.MethodPost, "/torrent-categories", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	notifierReq := NotifierRequest{
		Name:     "validation-test-notifier",
		TypeName: "discord",
		URL:      "https://discord.com/webhook/test",
	}
	body, _ = json.Marshal(notifierReq)
	req = httptest.NewRequest(http.MethodPost, "/notifiers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	feedReq := FeedRequest{
		Name: "Validation Test Feed",
		URL:  "https://example.com/feed",
	}
	body, _ = json.Marshal(feedReq)
	req = httptest.NewRequest(http.MethodPost, "/feeds", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	// Test: Missing feed_name
	invalidFilter := FeedFilterRequest{
		Name:         "Test Filter",
		CategoryName: "validation-test-category",
		NotifierName: "validation-test-notifier",
	}
	body, _ = json.Marshal(invalidFilter)
	req = httptest.NewRequest(http.MethodPost, "/feed-filters", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Test: Invalid feed_name
	invalidFilter = FeedFilterRequest{
		Name:         "Test Filter",
		FeedName:     "nonexistent-feed",
		CategoryName: "validation-test-category",
		NotifierName: "validation-test-notifier",
	}
	body, _ = json.Marshal(invalidFilter)
	req = httptest.NewRequest(http.MethodPost, "/feed-filters", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Test: Valid feed_name
	validFilter := FeedFilterRequest{
		Name:         "Test Filter",
		FeedName:     "Validation Test Feed",
		CategoryName: "validation-test-category",
		NotifierName: "validation-test-notifier",
	}
	body, _ = json.Marshal(validFilter)
	req = httptest.NewRequest(http.MethodPost, "/feed-filters", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var createdFilter FeedFilterResponse
	err := json.Unmarshal(rec.Body.Bytes(), &createdFilter)
	require.NoError(t, err)
	assert.Equal(t, "Validation Test Feed", createdFilter.FeedName)
	assert.NotZero(t, createdFilter.FeedID)
}

// TestFeedAuthorFilters tests CRUD operations for FeedAuthorFilters
func TestFeedAuthorFilters(t *testing.T) {
	e, cleanup := setupTestServer(t)
	defer cleanup()

	// Step 1: Create required entities for testing
	// Create a torrent category
	categoryReq := TorrentCategoryRequest{Name: "test-author-category"}
	body, _ := json.Marshal(categoryReq)
	req := httptest.NewRequest(http.MethodPost, "/torrent-categories", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	// Create a notifier
	notifierReq := NotifierRequest{
		Name:     "test-author-notifier",
		TypeName: "discord",
		URL:      "https://discord.com/webhook/test",
	}
	body, _ = json.Marshal(notifierReq)
	req = httptest.NewRequest(http.MethodPost, "/notifiers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	// Create two feeds for testing
	feedReq := FeedRequest{
		Name: "Test Author Feed 1",
		URL:  "https://example.com/feed1",
	}
	body, _ = json.Marshal(feedReq)
	req = httptest.NewRequest(http.MethodPost, "/feeds", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	var createdFeed1 FeedResponse
	err := json.Unmarshal(rec.Body.Bytes(), &createdFeed1)
	require.NoError(t, err)

	feedReq2 := FeedRequest{
		Name: "Test Author Feed 2",
		URL:  "https://example.com/feed2",
	}
	body, _ = json.Marshal(feedReq2)
	req = httptest.NewRequest(http.MethodPost, "/feeds", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	// Step 2: Test List - should be empty initially
	req = httptest.NewRequest(http.MethodGet, "/feed-author-filters", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var initialList []FeedAuthorFilterResponse
	err = json.Unmarshal(rec.Body.Bytes(), &initialList)
	require.NoError(t, err)
	assert.Equal(t, 0, len(initialList), "Should start with no feed author filters")

	// Step 3: Test Create
	createReq := FeedAuthorFilterRequest{
		Author:       "John Doe",
		FeedName:     "Test Author Feed 1",
		CategoryName: "test-author-category",
		NotifierName: "test-author-notifier",
	}
	body, _ = json.Marshal(createReq)
	req = httptest.NewRequest(http.MethodPost, "/feed-author-filters", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var createdFilter FeedAuthorFilterResponse
	err = json.Unmarshal(rec.Body.Bytes(), &createdFilter)
	require.NoError(t, err)
	assert.NotZero(t, createdFilter.ID)
	assert.Equal(t, "John Doe", createdFilter.Author)
	assert.Equal(t, "Test Author Feed 1", createdFilter.FeedName)
	assert.Equal(t, createdFeed1.ID, createdFilter.FeedID)
	assert.Equal(t, "test-author-category", createdFilter.Category)
	assert.Equal(t, "test-author-notifier", createdFilter.Notifier)

	// Step 4: Test Get by ID
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/feed-author-filters/%d", createdFilter.ID), nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var retrievedFilter FeedAuthorFilterResponse
	err = json.Unmarshal(rec.Body.Bytes(), &retrievedFilter)
	require.NoError(t, err)
	assert.Equal(t, createdFilter.ID, retrievedFilter.ID)
	assert.Equal(t, "John Doe", retrievedFilter.Author)

	// Step 5: Test Update
	updateReq := FeedAuthorFilterRequest{
		Author:       "Jane Smith",
		FeedName:     "Test Author Feed 2",
		CategoryName: "test-author-category",
		NotifierName: "test-author-notifier",
	}
	body, _ = json.Marshal(updateReq)
	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/feed-author-filters/%d", createdFilter.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var updatedFilter FeedAuthorFilterResponse
	err = json.Unmarshal(rec.Body.Bytes(), &updatedFilter)
	require.NoError(t, err)
	assert.Equal(t, createdFilter.ID, updatedFilter.ID)
	assert.Equal(t, "Jane Smith", updatedFilter.Author)
	assert.Equal(t, "Test Author Feed 2", updatedFilter.FeedName)

	// Step 6: Test List after create
	req = httptest.NewRequest(http.MethodGet, "/feed-author-filters", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var filters []FeedAuthorFilterResponse
	err = json.Unmarshal(rec.Body.Bytes(), &filters)
	require.NoError(t, err)
	assert.Equal(t, 1, len(filters))

	// Step 7: Test query parameter filtering by feed_id
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/feed-author-filters?feed_id=%d", updatedFilter.FeedID), nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var filteredList []FeedAuthorFilterResponse
	err = json.Unmarshal(rec.Body.Bytes(), &filteredList)
	require.NoError(t, err)
	require.Len(t, filteredList, 1, "Should have exactly 1 filter for this feed")
	assert.Equal(t, updatedFilter.ID, filteredList[0].ID)

	// Step 8: Test Delete
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/feed-author-filters/%d", createdFilter.ID), nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNoContent, rec.Code)

	// Verify deletion
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/feed-author-filters/%d", createdFilter.ID), nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// TestFeedAuthorFilterValidation tests validation for FeedAuthorFilter
func TestFeedAuthorFilterValidation(t *testing.T) {
	e, cleanup := setupTestServer(t)
	defer cleanup()

	// Setup: Create required entities
	categoryReq := TorrentCategoryRequest{Name: "validation-author-category"}
	body, _ := json.Marshal(categoryReq)
	req := httptest.NewRequest(http.MethodPost, "/torrent-categories", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	notifierReq := NotifierRequest{
		Name:     "validation-author-notifier",
		TypeName: "discord",
		URL:      "https://discord.com/webhook/test",
	}
	body, _ = json.Marshal(notifierReq)
	req = httptest.NewRequest(http.MethodPost, "/notifiers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	feedReq := FeedRequest{
		Name: "Validation Author Feed",
		URL:  "https://example.com/feed",
	}
	body, _ = json.Marshal(feedReq)
	req = httptest.NewRequest(http.MethodPost, "/feeds", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	// Test: Missing author
	invalidFilter := FeedAuthorFilterRequest{
		FeedName:     "Validation Author Feed",
		CategoryName: "validation-author-category",
		NotifierName: "validation-author-notifier",
	}
	body, _ = json.Marshal(invalidFilter)
	req = httptest.NewRequest(http.MethodPost, "/feed-author-filters", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Test: Missing feed_name
	invalidFilter = FeedAuthorFilterRequest{
		Author:       "Test Author",
		CategoryName: "validation-author-category",
		NotifierName: "validation-author-notifier",
	}
	body, _ = json.Marshal(invalidFilter)
	req = httptest.NewRequest(http.MethodPost, "/feed-author-filters", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Test: Invalid feed_name
	invalidFilter = FeedAuthorFilterRequest{
		Author:       "Test Author",
		FeedName:     "nonexistent-feed",
		CategoryName: "validation-author-category",
		NotifierName: "validation-author-notifier",
	}
	body, _ = json.Marshal(invalidFilter)
	req = httptest.NewRequest(http.MethodPost, "/feed-author-filters", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Test: Invalid category_name
	invalidFilter = FeedAuthorFilterRequest{
		Author:       "Test Author",
		FeedName:     "Validation Author Feed",
		CategoryName: "nonexistent-category",
		NotifierName: "validation-author-notifier",
	}
	body, _ = json.Marshal(invalidFilter)
	req = httptest.NewRequest(http.MethodPost, "/feed-author-filters", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Test: Invalid notifier_name
	invalidFilter = FeedAuthorFilterRequest{
		Author:       "Test Author",
		FeedName:     "Validation Author Feed",
		CategoryName: "validation-author-category",
		NotifierName: "nonexistent-notifier",
	}
	body, _ = json.Marshal(invalidFilter)
	req = httptest.NewRequest(http.MethodPost, "/feed-author-filters", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Test: Valid request
	validFilter := FeedAuthorFilterRequest{
		Author:       "Test Author",
		FeedName:     "Validation Author Feed",
		CategoryName: "validation-author-category",
		NotifierName: "validation-author-notifier",
	}
	body, _ = json.Marshal(validFilter)
	req = httptest.NewRequest(http.MethodPost, "/feed-author-filters", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var createdFilter FeedAuthorFilterResponse
	err := json.Unmarshal(rec.Body.Bytes(), &createdFilter)
	require.NoError(t, err)
	assert.Equal(t, "Test Author", createdFilter.Author)
	assert.Equal(t, "Validation Author Feed", createdFilter.FeedName)
	assert.NotZero(t, createdFilter.FeedID)
}

// TestFeedAuthorFilterUniqueConstraint tests the unique constraint on (feed_id, author)
func TestFeedAuthorFilterUniqueConstraint(t *testing.T) {
	e, cleanup := setupTestServer(t)
	defer cleanup()

	// Setup: Create required entities
	categoryReq := TorrentCategoryRequest{Name: "unique-test-category"}
	body, _ := json.Marshal(categoryReq)
	req := httptest.NewRequest(http.MethodPost, "/torrent-categories", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	notifierReq := NotifierRequest{
		Name:     "unique-test-notifier",
		TypeName: "discord",
		URL:      "https://discord.com/webhook/test",
	}
	body, _ = json.Marshal(notifierReq)
	req = httptest.NewRequest(http.MethodPost, "/notifiers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	feedReq := FeedRequest{
		Name: "Unique Test Feed",
		URL:  "https://example.com/feed",
	}
	body, _ = json.Marshal(feedReq)
	req = httptest.NewRequest(http.MethodPost, "/feeds", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	// Create first feed author filter
	createReq := FeedAuthorFilterRequest{
		Author:       "Unique Author",
		FeedName:     "Unique Test Feed",
		CategoryName: "unique-test-category",
		NotifierName: "unique-test-notifier",
	}
	body, _ = json.Marshal(createReq)
	req = httptest.NewRequest(http.MethodPost, "/feed-author-filters", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	// Try to create duplicate (same feed_id and author) - should fail
	body, _ = json.Marshal(createReq)
	req = httptest.NewRequest(http.MethodPost, "/feed-author-filters", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	// Create another feed
	feedReq2 := FeedRequest{
		Name: "Unique Test Feed 2",
		URL:  "https://example.com/feed2",
	}
	body, _ = json.Marshal(feedReq2)
	req = httptest.NewRequest(http.MethodPost, "/feeds", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	// Same author but different feed - should succeed
	createReq2 := FeedAuthorFilterRequest{
		Author:       "Unique Author",
		FeedName:     "Unique Test Feed 2",
		CategoryName: "unique-test-category",
		NotifierName: "unique-test-notifier",
	}
	body, _ = json.Marshal(createReq2)
	req = httptest.NewRequest(http.MethodPost, "/feed-author-filters", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)
}
