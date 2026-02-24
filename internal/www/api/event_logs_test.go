package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
)

// seedEventLogs inserts test event log entries directly into the database.
func seedEventLogs(t *testing.T, db *gorm.DB, entries []models.EventLog) {
	t.Helper()
	for i := range entries {
		require.NoError(t, db.Create(&entries[i]).Error)
	}
}

func TestEventLogs_List_Defaults(t *testing.T) {
	db, err := models.ConnectTestDB()
	require.NoError(t, err)
	e := SetupTestServerWithDB(db)

	now := time.Now()
	seedEventLogs(t, db, []models.EventLog{
		{CreatedAt: now.Add(-1 * time.Hour), Category: "download", EventType: "torrent.added", Source: "feedwatcher2", EntityType: "torrent", EntityID: "abc123", Summary: "Downloaded torrent"},
		{CreatedAt: now.Add(-2 * time.Hour), Category: "import", EventType: "import.completed", Source: "ebook-importer", EntityType: "torrent", EntityID: "def456", Summary: "Imported ebook"},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/event-logs", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp PaginatedEventLogResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

	assert.Equal(t, int64(2), resp.Total)
	assert.Equal(t, 1, resp.Page)
	assert.Equal(t, 50, resp.PerPage)
	assert.Len(t, resp.Items, 2)

	// Results ordered by created_at DESC — most recent first
	assert.Equal(t, "Downloaded torrent", resp.Items[0].Summary)
	assert.Equal(t, "Imported ebook", resp.Items[1].Summary)

	// Facets should be populated
	assert.Contains(t, resp.Facets.Categories, "download")
	assert.Contains(t, resp.Facets.Categories, "import")
	assert.Contains(t, resp.Facets.Sources, "feedwatcher2")
	assert.Contains(t, resp.Facets.Sources, "ebook-importer")
	assert.Contains(t, resp.Facets.EventTypes, "torrent.added")
	assert.Contains(t, resp.Facets.EventTypes, "import.completed")
	assert.Contains(t, resp.Facets.EntityTypes, "torrent")
}

func TestEventLogs_List_EmptyResult(t *testing.T) {
	db, err := models.ConnectTestDB()
	require.NoError(t, err)
	e := SetupTestServerWithDB(db)

	req := httptest.NewRequest(http.MethodGet, "/api/event-logs", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp PaginatedEventLogResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

	assert.Equal(t, int64(0), resp.Total)
	assert.NotNil(t, resp.Items)
	assert.Empty(t, resp.Items)
	assert.NotNil(t, resp.Facets.Categories)
	assert.NotNil(t, resp.Facets.Sources)
	assert.NotNil(t, resp.Facets.EventTypes)
	assert.NotNil(t, resp.Facets.EntityTypes)
}

func TestEventLogs_List_Pagination(t *testing.T) {
	db, err := models.ConnectTestDB()
	require.NoError(t, err)
	e := SetupTestServerWithDB(db)

	now := time.Now()
	entries := make([]models.EventLog, 5)
	for i := range entries {
		entries[i] = models.EventLog{
			CreatedAt: now.Add(-time.Duration(i) * time.Hour),
			Category:  "download",
			EventType: "torrent.added",
			Source:    "feedwatcher2",
			EntityType: "torrent",
			EntityID:  fmt.Sprintf("hash-%d", i),
			Summary:   fmt.Sprintf("Event %d", i),
		}
	}
	seedEventLogs(t, db, entries)

	t.Run("page 1 with per_page 2", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/event-logs?page=1&per_page=2", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp PaginatedEventLogResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

		assert.Equal(t, int64(5), resp.Total)
		assert.Equal(t, 1, resp.Page)
		assert.Equal(t, 2, resp.PerPage)
		assert.Len(t, resp.Items, 2)
		assert.Equal(t, "Event 0", resp.Items[0].Summary) // most recent
		assert.Equal(t, "Event 1", resp.Items[1].Summary)
	})

	t.Run("page 2 with per_page 2", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/event-logs?page=2&per_page=2", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp PaginatedEventLogResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

		assert.Equal(t, int64(5), resp.Total)
		assert.Equal(t, 2, resp.Page)
		assert.Len(t, resp.Items, 2)
		assert.Equal(t, "Event 2", resp.Items[0].Summary)
		assert.Equal(t, "Event 3", resp.Items[1].Summary)
	})

	t.Run("last page has remaining items", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/event-logs?page=3&per_page=2", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp PaginatedEventLogResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

		assert.Equal(t, int64(5), resp.Total)
		assert.Len(t, resp.Items, 1)
		assert.Equal(t, "Event 4", resp.Items[0].Summary)
	})

	t.Run("page beyond results returns empty", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/event-logs?page=10&per_page=2", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp PaginatedEventLogResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

		assert.Equal(t, int64(5), resp.Total)
		assert.Empty(t, resp.Items)
	})

	t.Run("per_page clamped to 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/event-logs?per_page=999", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp PaginatedEventLogResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

		assert.Equal(t, 200, resp.PerPage)
	})

	t.Run("invalid page defaults to 1", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/event-logs?page=abc", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp PaginatedEventLogResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

		assert.Equal(t, 1, resp.Page)
	})
}

func TestEventLogs_List_FilterByCategory(t *testing.T) {
	db, err := models.ConnectTestDB()
	require.NoError(t, err)
	e := SetupTestServerWithDB(db)

	now := time.Now()
	seedEventLogs(t, db, []models.EventLog{
		{CreatedAt: now, Category: "download", EventType: "torrent.added", Source: "feedwatcher2", Summary: "Download event"},
		{CreatedAt: now, Category: "import", EventType: "import.completed", Source: "ebook-importer", Summary: "Import event"},
		{CreatedAt: now, Category: "download", EventType: "torrent.added", Source: "feedwatcher2", Summary: "Another download"},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/event-logs?category=download", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp PaginatedEventLogResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

	assert.Equal(t, int64(2), resp.Total)
	assert.Len(t, resp.Items, 2)
	for _, item := range resp.Items {
		assert.Equal(t, "download", item.Category)
	}
}

func TestEventLogs_List_FilterBySource(t *testing.T) {
	db, err := models.ConnectTestDB()
	require.NoError(t, err)
	e := SetupTestServerWithDB(db)

	now := time.Now()
	seedEventLogs(t, db, []models.EventLog{
		{CreatedAt: now, Category: "download", EventType: "torrent.added", Source: "feedwatcher2", Summary: "FW event"},
		{CreatedAt: now, Category: "import", EventType: "import.completed", Source: "ebook-importer", Summary: "Importer event"},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/event-logs?source=ebook-importer", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp PaginatedEventLogResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

	assert.Equal(t, int64(1), resp.Total)
	assert.Len(t, resp.Items, 1)
	assert.Equal(t, "ebook-importer", resp.Items[0].Source)
}

func TestEventLogs_List_FilterByEventType(t *testing.T) {
	db, err := models.ConnectTestDB()
	require.NoError(t, err)
	e := SetupTestServerWithDB(db)

	now := time.Now()
	seedEventLogs(t, db, []models.EventLog{
		{CreatedAt: now, Category: "download", EventType: "torrent.added", Source: "feedwatcher2", Summary: "Added"},
		{CreatedAt: now, Category: "download", EventType: "torrent.duplicate_skipped", Source: "feedwatcher2", Summary: "Skipped"},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/event-logs?event_type=torrent.duplicate_skipped", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp PaginatedEventLogResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

	assert.Equal(t, int64(1), resp.Total)
	assert.Equal(t, "torrent.duplicate_skipped", resp.Items[0].EventType)
}

func TestEventLogs_List_FilterByEntityType(t *testing.T) {
	db, err := models.ConnectTestDB()
	require.NoError(t, err)
	e := SetupTestServerWithDB(db)

	now := time.Now()
	seedEventLogs(t, db, []models.EventLog{
		{CreatedAt: now, Category: "download", EventType: "torrent.added", Source: "feedwatcher2", EntityType: "torrent", Summary: "Torrent event"},
		{CreatedAt: now, Category: "mutation", EventType: "created", Source: "api", EntityType: "author", Summary: "Author event"},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/event-logs?entity_type=author", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp PaginatedEventLogResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

	assert.Equal(t, int64(1), resp.Total)
	assert.Equal(t, "author", resp.Items[0].EntityType)
}

func TestEventLogs_List_FilterByEntityID(t *testing.T) {
	db, err := models.ConnectTestDB()
	require.NoError(t, err)
	e := SetupTestServerWithDB(db)

	now := time.Now()
	seedEventLogs(t, db, []models.EventLog{
		{CreatedAt: now, Category: "download", EventType: "torrent.added", Source: "feedwatcher2", EntityType: "torrent", EntityID: "hash-aaa", Summary: "First"},
		{CreatedAt: now, Category: "download", EventType: "torrent.added", Source: "feedwatcher2", EntityType: "torrent", EntityID: "hash-bbb", Summary: "Second"},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/event-logs?entity_id=hash-aaa", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp PaginatedEventLogResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

	assert.Equal(t, int64(1), resp.Total)
	assert.Equal(t, "hash-aaa", resp.Items[0].EntityID)
}

func TestEventLogs_List_FilterBySummarySearch(t *testing.T) {
	db, err := models.ConnectTestDB()
	require.NoError(t, err)
	e := SetupTestServerWithDB(db)

	now := time.Now()
	seedEventLogs(t, db, []models.EventLog{
		{CreatedAt: now, Category: "download", EventType: "torrent.added", Source: "feedwatcher2", Summary: "Downloaded Brandon Sanderson book"},
		{CreatedAt: now, Category: "import", EventType: "import.completed", Source: "ebook-importer", Summary: "Imported Patrick Rothfuss novel"},
	})

	t.Run("partial match", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/event-logs?q=brandon", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp PaginatedEventLogResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

		assert.Equal(t, int64(1), resp.Total)
		assert.Contains(t, resp.Items[0].Summary, "Brandon")
	})

	t.Run("case insensitive", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/event-logs?q=PATRICK", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp PaginatedEventLogResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

		assert.Equal(t, int64(1), resp.Total)
		assert.Contains(t, resp.Items[0].Summary, "Patrick")
	})

	t.Run("no matches", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/event-logs?q=nonexistent999", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp PaginatedEventLogResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

		assert.Equal(t, int64(0), resp.Total)
		assert.Empty(t, resp.Items)
	})
}

func TestEventLogs_List_FilterByDateRange(t *testing.T) {
	db, err := models.ConnectTestDB()
	require.NoError(t, err)
	e := SetupTestServerWithDB(db)

	now := time.Now().UTC()
	seedEventLogs(t, db, []models.EventLog{
		{CreatedAt: now.Add(-48 * time.Hour), Category: "download", EventType: "torrent.added", Source: "feedwatcher2", Summary: "Old event"},
		{CreatedAt: now.Add(-1 * time.Hour), Category: "import", EventType: "import.completed", Source: "ebook-importer", Summary: "Recent event"},
	})

	t.Run("from filter", func(t *testing.T) {
		from := now.Add(-24 * time.Hour).Format(time.RFC3339)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/event-logs?from=%s", from), nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp PaginatedEventLogResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

		assert.Equal(t, int64(1), resp.Total)
		assert.Equal(t, "Recent event", resp.Items[0].Summary)
	})

	t.Run("to filter", func(t *testing.T) {
		to := now.Add(-24 * time.Hour).Format(time.RFC3339)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/event-logs?to=%s", to), nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp PaginatedEventLogResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

		assert.Equal(t, int64(1), resp.Total)
		assert.Equal(t, "Old event", resp.Items[0].Summary)
	})

	t.Run("from and to combined", func(t *testing.T) {
		from := now.Add(-2 * time.Hour).Format(time.RFC3339)
		to := now.Format(time.RFC3339)
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/event-logs?from=%s&to=%s", from, to), nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp PaginatedEventLogResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

		assert.Equal(t, int64(1), resp.Total)
		assert.Equal(t, "Recent event", resp.Items[0].Summary)
	})
}

func TestEventLogs_List_MultipleFilters(t *testing.T) {
	db, err := models.ConnectTestDB()
	require.NoError(t, err)
	e := SetupTestServerWithDB(db)

	now := time.Now()
	seedEventLogs(t, db, []models.EventLog{
		{CreatedAt: now, Category: "download", EventType: "torrent.added", Source: "feedwatcher2", EntityType: "torrent", Summary: "Match"},
		{CreatedAt: now, Category: "download", EventType: "torrent.added", Source: "discord-bot", EntityType: "torrent", Summary: "Wrong source"},
		{CreatedAt: now, Category: "import", EventType: "import.completed", Source: "feedwatcher2", EntityType: "torrent", Summary: "Wrong category"},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/event-logs?category=download&source=feedwatcher2", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp PaginatedEventLogResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

	assert.Equal(t, int64(1), resp.Total)
	assert.Equal(t, "Match", resp.Items[0].Summary)
}

func TestEventLogs_List_FacetsNotFilteredByNonDateFilters(t *testing.T) {
	db, err := models.ConnectTestDB()
	require.NoError(t, err)
	e := SetupTestServerWithDB(db)

	now := time.Now()
	seedEventLogs(t, db, []models.EventLog{
		{CreatedAt: now, Category: "download", EventType: "torrent.added", Source: "feedwatcher2", Summary: "DL"},
		{CreatedAt: now, Category: "import", EventType: "import.completed", Source: "ebook-importer", Summary: "Import"},
	})

	// Filter results to only download category, but facets should still show all categories
	req := httptest.NewRequest(http.MethodGet, "/api/event-logs?category=download", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp PaginatedEventLogResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

	// Only 1 result because of category filter
	assert.Equal(t, int64(1), resp.Total)

	// But facets should include both categories (facets use date range only)
	assert.Contains(t, resp.Facets.Categories, "download")
	assert.Contains(t, resp.Facets.Categories, "import")
	assert.Contains(t, resp.Facets.Sources, "feedwatcher2")
	assert.Contains(t, resp.Facets.Sources, "ebook-importer")
}

func TestEventLogs_Get(t *testing.T) {
	db, err := models.ConnectTestDB()
	require.NoError(t, err)
	e := SetupTestServerWithDB(db)

	now := time.Now()
	entry := models.EventLog{
		CreatedAt:  now,
		Category:   "download",
		EventType:  "torrent.added",
		Source:     "feedwatcher2",
		EntityType: "torrent",
		EntityID:   "abc123",
		Summary:    "Downloaded a book",
		Details:    `{"title":"Test Book","author":"Test Author"}`,
	}
	require.NoError(t, db.Create(&entry).Error)

	t.Run("get existing event", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/event-logs/%d", entry.ID), nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp EventLogResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

		assert.Equal(t, entry.ID, resp.ID)
		assert.Equal(t, "download", resp.Category)
		assert.Equal(t, "torrent.added", resp.EventType)
		assert.Equal(t, "feedwatcher2", resp.Source)
		assert.Equal(t, "torrent", resp.EntityType)
		assert.Equal(t, "abc123", resp.EntityID)
		assert.Equal(t, "Downloaded a book", resp.Summary)
		assert.Equal(t, `{"title":"Test Book","author":"Test Author"}`, resp.Details)
	})

	t.Run("get nonexistent event returns 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/event-logs/99999", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("invalid ID returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/event-logs/invalid", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
