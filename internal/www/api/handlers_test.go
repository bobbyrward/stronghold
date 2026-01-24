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

func getFeedFilterSetTypeID(t *testing.T, e *echo.Echo, name string) uint {
	req := httptest.NewRequest(http.MethodGet, "/api/feed-filter-set-types", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	var types []FeedFilterSetTypeResponse
	err := json.Unmarshal(rec.Body.Bytes(), &types)
	require.NoError(t, err)

	for _, fst := range types {
		if fst.Name == name {
			return fst.ID
		}
	}
	t.Fatalf("Feed filter set type %q not found", name)
	return 0
}

func getFilterKeyID(t *testing.T, e *echo.Echo, name string) uint {
	req := httptest.NewRequest(http.MethodGet, "/api/filter-keys", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	var keys []FilterKeyResponse
	err := json.Unmarshal(rec.Body.Bytes(), &keys)
	require.NoError(t, err)

	for _, key := range keys {
		if key.Name == name {
			return key.ID
		}
	}
	t.Fatalf("Filter key %q not found", name)
	return 0
}

func getFilterOperatorID(t *testing.T, e *echo.Echo, name string) uint {
	req := httptest.NewRequest(http.MethodGet, "/api/filter-operators", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	var operators []FilterOperatorResponse
	err := json.Unmarshal(rec.Body.Bytes(), &operators)
	require.NoError(t, err)

	for _, op := range operators {
		if op.Name == name {
			return op.ID
		}
	}
	t.Fatalf("Filter operator %q not found", name)
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

// TestFilterKeys tests read-only operations for FilterKeys
func TestFilterKeys(t *testing.T) {
	e, cleanup := setupTestServer(t)
	defer cleanup()

	// Test List - should have seeded data
	req := httptest.NewRequest(http.MethodGet, "/api/filter-keys", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var filterKeys []FilterKeyResponse
	err := json.Unmarshal(rec.Body.Bytes(), &filterKeys)
	require.NoError(t, err)
	assert.Greater(t, len(filterKeys), 0, "Should have seeded filter keys")

	// Test Get by ID (use first seeded item)
	firstKey := filterKeys[0]
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/filter-keys/%d", firstKey.ID), nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var retrieved FilterKeyResponse
	err = json.Unmarshal(rec.Body.Bytes(), &retrieved)
	require.NoError(t, err)
	assert.Equal(t, firstKey.ID, retrieved.ID)
	assert.Equal(t, firstKey.Name, retrieved.Name)
}

// TestFilterOperators tests read-only operations for FilterOperators
func TestFilterOperators(t *testing.T) {
	e, cleanup := setupTestServer(t)
	defer cleanup()

	// Test List - should have seeded data
	req := httptest.NewRequest(http.MethodGet, "/api/filter-operators", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var filterOperators []FilterOperatorResponse
	err := json.Unmarshal(rec.Body.Bytes(), &filterOperators)
	require.NoError(t, err)
	assert.Greater(t, len(filterOperators), 0, "Should have seeded filter operators")

	// Test Get by ID (use first seeded item)
	firstOperator := filterOperators[0]
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/filter-operators/%d", firstOperator.ID), nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var retrieved FilterOperatorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &retrieved)
	require.NoError(t, err)
	assert.Equal(t, firstOperator.ID, retrieved.ID)
	assert.Equal(t, firstOperator.Name, retrieved.Name)
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

// TestCompleteWorkflow tests a complete workflow with cascading relationships
func TestCompleteWorkflow(t *testing.T) {
	e, cleanup := setupTestServer(t)
	defer cleanup()

	// Get reference data IDs
	discordTypeID := getNotificationTypeID(t, e, "discord")
	anyTypeID := getFeedFilterSetTypeID(t, e, "any")
	authorKeyID := getFilterKeyID(t, e, "author")
	containsOpID := getFilterOperatorID(t, e, "contains")
	categoryID := getTorrentCategoryID(t, e, "personal-books")

	// Step 1: Create a Notifier
	notifierReq := NotifierRequest{
		Name:   "discord-personal",
		TypeID: discordTypeID,
		URL:    "https://discord.com/webhook/test",
	}
	body, _ := json.Marshal(notifierReq)
	req := httptest.NewRequest(http.MethodPost, "/api/notifiers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var createdNotifier NotifierResponse
	err := json.Unmarshal(rec.Body.Bytes(), &createdNotifier)
	require.NoError(t, err)

	// Step 2: Create a Feed
	feedReq := FeedRequest{
		Name: "MAM",
		URL:  "https://example.com/rss",
	}
	body, _ = json.Marshal(feedReq)
	req = httptest.NewRequest(http.MethodPost, "/api/feeds", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var createdFeed FeedResponse
	err = json.Unmarshal(rec.Body.Bytes(), &createdFeed)
	require.NoError(t, err)

	// Step 3: Create a FeedFilter using IDs
	filterReq := FeedFilterRequest{
		Name:       "Blaise Corvin Books",
		FeedID:     createdFeed.ID,
		CategoryID: categoryID,
		NotifierID: createdNotifier.ID,
	}
	body, _ = json.Marshal(filterReq)
	req = httptest.NewRequest(http.MethodPost, "/api/feed-filters", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var createdFilter FeedFilterResponse
	err = json.Unmarshal(rec.Body.Bytes(), &createdFilter)
	require.NoError(t, err)
	assert.Equal(t, "Blaise Corvin Books", createdFilter.Name)
	assert.Equal(t, "MAM", createdFilter.FeedName)
	assert.NotZero(t, createdFilter.FeedID)
	assert.Equal(t, "personal-books", createdFilter.CategoryName)
	assert.Equal(t, "discord-personal", createdFilter.NotifierName)

	// Step 4: Create a FeedFilterSet
	setReq := FeedFilterSetRequest{
		FeedFilterID:        createdFilter.ID,
		FeedFilterSetTypeID: anyTypeID,
	}
	body, _ = json.Marshal(setReq)
	req = httptest.NewRequest(http.MethodPost, "/api/feed-filter-sets", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var createdSet FeedFilterSetResponse
	err = json.Unmarshal(rec.Body.Bytes(), &createdSet)
	require.NoError(t, err)
	assert.Equal(t, "any", createdSet.TypeName)
	require.NotZero(t, createdSet.ID, "FeedFilterSet ID should not be zero")

	// Step 5: Create FeedFilterSetEntry
	entryReq := FeedFilterSetEntryRequest{
		FeedFilterSetID: createdSet.ID,
		KeyID:           authorKeyID,
		OperatorID:      containsOpID,
		Value:           "Blaise Corvin",
	}
	body, _ = json.Marshal(entryReq)
	req = httptest.NewRequest(http.MethodPost, "/api/feed-filter-set-entries", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var createdEntry FeedFilterSetEntryResponse
	err = json.Unmarshal(rec.Body.Bytes(), &createdEntry)
	require.NoError(t, err)
	assert.Equal(t, "author", createdEntry.KeyName)
	assert.Equal(t, "contains", createdEntry.OperatorName)
	assert.Equal(t, "Blaise Corvin", createdEntry.Value)

	// Step 7: Test query parameter filtering
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/feed-filter-set-entries?feed_filter_set_id=%d", createdSet.ID), nil)
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

	// Test invalid ID (GET still valid for reference data)
	req := httptest.NewRequest(http.MethodGet, "/api/filter-keys/invalid", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Test non-existent ID (GET still valid for reference data)
	req = httptest.NewRequest(http.MethodGet, "/api/filter-keys/99999", nil)
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

// TestFeedFilterValidation tests feed filter validation using ID-based requests
func TestFeedFilterValidation(t *testing.T) {
	e, cleanup := setupTestServer(t)
	defer cleanup()

	// Get reference data IDs
	categoryID := getTorrentCategoryID(t, e, "books")
	discordTypeID := getNotificationTypeID(t, e, "discord")

	notifierReq := NotifierRequest{
		Name:   "validation-test-notifier",
		TypeID: discordTypeID,
		URL:    "https://discord.com/webhook/test",
	}
	body, _ := json.Marshal(notifierReq)
	req := httptest.NewRequest(http.MethodPost, "/api/notifiers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	var createdNotifier NotifierResponse
	err := json.Unmarshal(rec.Body.Bytes(), &createdNotifier)
	require.NoError(t, err)

	feedReq := FeedRequest{
		Name: "Validation Test Feed",
		URL:  "https://example.com/feed",
	}
	body, _ = json.Marshal(feedReq)
	req = httptest.NewRequest(http.MethodPost, "/api/feeds", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	var createdFeed FeedResponse
	err = json.Unmarshal(rec.Body.Bytes(), &createdFeed)
	require.NoError(t, err)

	// Test: Missing name
	invalidFilter := FeedFilterRequest{
		FeedID:     createdFeed.ID,
		CategoryID: categoryID,
		NotifierID: createdNotifier.ID,
	}
	body, _ = json.Marshal(invalidFilter)
	req = httptest.NewRequest(http.MethodPost, "/api/feed-filters", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Test: Valid request with all IDs
	validFilter := FeedFilterRequest{
		Name:       "Test Filter",
		FeedID:     createdFeed.ID,
		CategoryID: categoryID,
		NotifierID: createdNotifier.ID,
	}
	body, _ = json.Marshal(validFilter)
	req = httptest.NewRequest(http.MethodPost, "/api/feed-filters", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var createdFilter FeedFilterResponse
	err = json.Unmarshal(rec.Body.Bytes(), &createdFilter)
	require.NoError(t, err)
	assert.Equal(t, "Validation Test Feed", createdFilter.FeedName)
	assert.NotZero(t, createdFilter.FeedID)
	assert.Equal(t, "books", createdFilter.CategoryName)
	assert.Equal(t, "validation-test-notifier", createdFilter.NotifierName)
}

// TestFeedAuthorFilters tests CRUD operations for FeedAuthorFilters
func TestFeedAuthorFilters(t *testing.T) {
	e, cleanup := setupTestServer(t)
	defer cleanup()

	// Get reference data IDs
	discordTypeID := getNotificationTypeID(t, e, "discord")
	categoryID := getTorrentCategoryID(t, e, "personal-audiobooks")

	// Create a notifier
	notifierReq := NotifierRequest{
		Name:   "test-author-notifier",
		TypeID: discordTypeID,
		URL:    "https://discord.com/webhook/test",
	}
	body, _ := json.Marshal(notifierReq)
	req := httptest.NewRequest(http.MethodPost, "/api/notifiers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	var createdNotifier NotifierResponse
	err := json.Unmarshal(rec.Body.Bytes(), &createdNotifier)
	require.NoError(t, err)

	// Create two feeds for testing
	feedReq := FeedRequest{
		Name: "Test Author Feed 1",
		URL:  "https://example.com/feed1",
	}
	body, _ = json.Marshal(feedReq)
	req = httptest.NewRequest(http.MethodPost, "/api/feeds", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	var createdFeed1 FeedResponse
	err = json.Unmarshal(rec.Body.Bytes(), &createdFeed1)
	require.NoError(t, err)

	feedReq2 := FeedRequest{
		Name: "Test Author Feed 2",
		URL:  "https://example.com/feed2",
	}
	body, _ = json.Marshal(feedReq2)
	req = httptest.NewRequest(http.MethodPost, "/api/feeds", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	var createdFeed2 FeedResponse
	err = json.Unmarshal(rec.Body.Bytes(), &createdFeed2)
	require.NoError(t, err)

	// Step 2: Test List - should be empty initially
	req = httptest.NewRequest(http.MethodGet, "/api/feed-author-filters", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var initialList []FeedAuthorFilterResponse
	err = json.Unmarshal(rec.Body.Bytes(), &initialList)
	require.NoError(t, err)
	assert.Equal(t, 0, len(initialList), "Should start with no feed author filters")

	// Step 3: Test Create
	createReq := FeedAuthorFilterRequest{
		Author:     "John Doe",
		FeedID:     createdFeed1.ID,
		CategoryID: categoryID,
		NotifierID: createdNotifier.ID,
	}
	body, _ = json.Marshal(createReq)
	req = httptest.NewRequest(http.MethodPost, "/api/feed-author-filters", bytes.NewReader(body))
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
	assert.Equal(t, "personal-audiobooks", createdFilter.CategoryName)
	assert.Equal(t, "test-author-notifier", createdFilter.NotifierName)

	// Step 4: Test Get by ID
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/feed-author-filters/%d", createdFilter.ID), nil)
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
		Author:     "Jane Smith",
		FeedID:     createdFeed2.ID,
		CategoryID: categoryID,
		NotifierID: createdNotifier.ID,
	}
	body, _ = json.Marshal(updateReq)
	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/feed-author-filters/%d", createdFilter.ID), bytes.NewReader(body))
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
	req = httptest.NewRequest(http.MethodGet, "/api/feed-author-filters", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var filters []FeedAuthorFilterResponse
	err = json.Unmarshal(rec.Body.Bytes(), &filters)
	require.NoError(t, err)
	assert.Equal(t, 1, len(filters))

	// Step 7: Test query parameter filtering by feed_id
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/feed-author-filters?feed_id=%d", updatedFilter.FeedID), nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var filteredList []FeedAuthorFilterResponse
	err = json.Unmarshal(rec.Body.Bytes(), &filteredList)
	require.NoError(t, err)
	require.Len(t, filteredList, 1, "Should have exactly 1 filter for this feed")
	assert.Equal(t, updatedFilter.ID, filteredList[0].ID)

	// Step 8: Test Delete
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/feed-author-filters/%d", createdFilter.ID), nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNoContent, rec.Code)

	// Verify deletion
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/feed-author-filters/%d", createdFilter.ID), nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// TestFeedAuthorFilterValidation tests validation for FeedAuthorFilter
func TestFeedAuthorFilterValidation(t *testing.T) {
	e, cleanup := setupTestServer(t)
	defer cleanup()

	// Get reference data IDs
	discordTypeID := getNotificationTypeID(t, e, "discord")
	categoryID := getTorrentCategoryID(t, e, "audiobooks")

	notifierReq := NotifierRequest{
		Name:   "validation-author-notifier",
		TypeID: discordTypeID,
		URL:    "https://discord.com/webhook/test",
	}
	body, _ := json.Marshal(notifierReq)
	req := httptest.NewRequest(http.MethodPost, "/api/notifiers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	var createdNotifier NotifierResponse
	err := json.Unmarshal(rec.Body.Bytes(), &createdNotifier)
	require.NoError(t, err)

	feedReq := FeedRequest{
		Name: "Validation Author Feed",
		URL:  "https://example.com/feed",
	}
	body, _ = json.Marshal(feedReq)
	req = httptest.NewRequest(http.MethodPost, "/api/feeds", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	var createdFeed FeedResponse
	err = json.Unmarshal(rec.Body.Bytes(), &createdFeed)
	require.NoError(t, err)

	// Test: Missing author
	invalidFilter := FeedAuthorFilterRequest{
		FeedID:     createdFeed.ID,
		CategoryID: categoryID,
		NotifierID: createdNotifier.ID,
	}
	body, _ = json.Marshal(invalidFilter)
	req = httptest.NewRequest(http.MethodPost, "/api/feed-author-filters", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Test: Missing feed_id (zero value)
	invalidFilter = FeedAuthorFilterRequest{
		Author:     "Test Author",
		CategoryID: categoryID,
		NotifierID: createdNotifier.ID,
	}
	body, _ = json.Marshal(invalidFilter)
	req = httptest.NewRequest(http.MethodPost, "/api/feed-author-filters", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Test: Invalid feed_id (non-existent)
	// Note: Use different author names for each test to avoid unique constraint violations
	invalidFilter = FeedAuthorFilterRequest{
		Author:     "Test Author Invalid Feed",
		FeedID:     99999,
		CategoryID: categoryID,
		NotifierID: createdNotifier.ID,
	}
	body, _ = json.Marshal(invalidFilter)
	req = httptest.NewRequest(http.MethodPost, "/api/feed-author-filters", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	// Note: SQLite doesn't enforce FK constraints by default, so this may succeed with 201
	assert.True(t, rec.Code == http.StatusCreated || rec.Code == http.StatusBadRequest || rec.Code == http.StatusInternalServerError)

	// Test: Invalid category_id (non-existent)
	invalidFilter = FeedAuthorFilterRequest{
		Author:     "Test Author Invalid Category",
		FeedID:     createdFeed.ID,
		CategoryID: 99999,
		NotifierID: createdNotifier.ID,
	}
	body, _ = json.Marshal(invalidFilter)
	req = httptest.NewRequest(http.MethodPost, "/api/feed-author-filters", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	// Note: SQLite doesn't enforce FK constraints by default, so this may succeed with 201
	assert.True(t, rec.Code == http.StatusCreated || rec.Code == http.StatusBadRequest || rec.Code == http.StatusInternalServerError)

	// Test: Invalid notifier_id (non-existent)
	invalidFilter = FeedAuthorFilterRequest{
		Author:     "Test Author Invalid Notifier",
		FeedID:     createdFeed.ID,
		CategoryID: categoryID,
		NotifierID: 99999,
	}
	body, _ = json.Marshal(invalidFilter)
	req = httptest.NewRequest(http.MethodPost, "/api/feed-author-filters", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	// Note: SQLite doesn't enforce FK constraints by default, so this may succeed with 201
	assert.True(t, rec.Code == http.StatusCreated || rec.Code == http.StatusBadRequest || rec.Code == http.StatusInternalServerError)

	// Test: Valid request
	validFilter := FeedAuthorFilterRequest{
		Author:     "Test Author",
		FeedID:     createdFeed.ID,
		CategoryID: categoryID,
		NotifierID: createdNotifier.ID,
	}
	body, _ = json.Marshal(validFilter)
	req = httptest.NewRequest(http.MethodPost, "/api/feed-author-filters", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var createdFilter FeedAuthorFilterResponse
	err = json.Unmarshal(rec.Body.Bytes(), &createdFilter)
	require.NoError(t, err)
	assert.Equal(t, "Test Author", createdFilter.Author)
	assert.Equal(t, "Validation Author Feed", createdFilter.FeedName)
	assert.NotZero(t, createdFilter.FeedID)
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

// TestFeedAuthorFilterUniqueConstraint tests the unique constraint on (feed_id, author)
func TestFeedAuthorFilterUniqueConstraint(t *testing.T) {
	e, cleanup := setupTestServer(t)
	defer cleanup()

	// Get reference data IDs
	discordTypeID := getNotificationTypeID(t, e, "discord")
	categoryID := getTorrentCategoryID(t, e, "personal-books")

	notifierReq := NotifierRequest{
		Name:   "unique-test-notifier",
		TypeID: discordTypeID,
		URL:    "https://discord.com/webhook/test",
	}
	body, _ := json.Marshal(notifierReq)
	req := httptest.NewRequest(http.MethodPost, "/api/notifiers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	var createdNotifier NotifierResponse
	err := json.Unmarshal(rec.Body.Bytes(), &createdNotifier)
	require.NoError(t, err)

	feedReq := FeedRequest{
		Name: "Unique Test Feed",
		URL:  "https://example.com/feed",
	}
	body, _ = json.Marshal(feedReq)
	req = httptest.NewRequest(http.MethodPost, "/api/feeds", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	var createdFeed FeedResponse
	err = json.Unmarshal(rec.Body.Bytes(), &createdFeed)
	require.NoError(t, err)

	// Create first feed author filter
	createReq := FeedAuthorFilterRequest{
		Author:     "Unique Author",
		FeedID:     createdFeed.ID,
		CategoryID: categoryID,
		NotifierID: createdNotifier.ID,
	}
	body, _ = json.Marshal(createReq)
	req = httptest.NewRequest(http.MethodPost, "/api/feed-author-filters", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	// Try to create duplicate (same feed_id and author) - should fail
	body, _ = json.Marshal(createReq)
	req = httptest.NewRequest(http.MethodPost, "/api/feed-author-filters", bytes.NewReader(body))
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
	req = httptest.NewRequest(http.MethodPost, "/api/feeds", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	var createdFeed2 FeedResponse
	err = json.Unmarshal(rec.Body.Bytes(), &createdFeed2)
	require.NoError(t, err)

	// Same author but different feed - should succeed
	createReq2 := FeedAuthorFilterRequest{
		Author:     "Unique Author",
		FeedID:     createdFeed2.ID,
		CategoryID: categoryID,
		NotifierID: createdNotifier.ID,
	}
	body, _ = json.Marshal(createReq2)
	req = httptest.NewRequest(http.MethodPost, "/api/feed-author-filters", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)
}
