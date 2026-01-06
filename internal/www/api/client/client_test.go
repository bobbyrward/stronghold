package client

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bobbyrward/stronghold/internal/www/api"
)

// setupTestServer creates a test HTTP server using the shared api.SetupTestServer
func setupTestServer(t *testing.T) (*httptest.Server, func()) {
	e, _ := api.SetupTestServer(t)
	server := httptest.NewServer(e)

	cleanup := func() {
		server.Close()
	}

	return server, cleanup
}

// getTorrentCategoryByName returns a pre-populated TorrentCategory by name
func getTorrentCategoryByName(t *testing.T, client *Client, ctx context.Context, name string) TorrentCategoryResponse {
	categories, err := client.TorrentCategories.List(ctx)
	require.NoError(t, err)
	for _, cat := range categories {
		if cat.Name == name {
			return cat
		}
	}
	t.Fatalf("TorrentCategory with name %q not found", name)
	return TorrentCategoryResponse{}
}

func TestNewClient(t *testing.T) {
	baseURL := "http://localhost:8000"
	client := NewClient(baseURL)

	assert.Equal(t, baseURL, client.BaseUrl)
	assert.NotNil(t, client.Feeds)
	assert.NotNil(t, client.FeedAuthorFilters)
	assert.NotNil(t, client.FeedFilterSetEntries)
	assert.NotNil(t, client.FeedFilterSetTypes)
	assert.NotNil(t, client.FeedFilterSets)
	assert.NotNil(t, client.FeedFilters)
	assert.NotNil(t, client.FilterKeys)
	assert.NotNil(t, client.FilterOperators)
	assert.NotNil(t, client.NotificationTypes)
	assert.NotNil(t, client.Notifiers)
	assert.NotNil(t, client.TorrentCategories)

	// Verify sub-clients have correct base URL and type names
	assert.Equal(t, baseURL, client.Feeds.BaseUrl)
	assert.Equal(t, "feeds", client.Feeds.TypeName)

	assert.Equal(t, baseURL, client.FeedAuthorFilters.BaseUrl)
	assert.Equal(t, "feed-author-filters", client.FeedAuthorFilters.TypeName)

	assert.Equal(t, baseURL, client.FeedFilterSetEntries.BaseUrl)
	assert.Equal(t, "feed-filter-set-entries", client.FeedFilterSetEntries.TypeName)

	assert.Equal(t, baseURL, client.FeedFilterSetTypes.BaseUrl)
	assert.Equal(t, "feed-filter-set-types", client.FeedFilterSetTypes.TypeName)

	assert.Equal(t, baseURL, client.FeedFilterSets.BaseUrl)
	assert.Equal(t, "feed-filter-sets", client.FeedFilterSets.TypeName)

	assert.Equal(t, baseURL, client.FeedFilters.BaseUrl)
	assert.Equal(t, "feed-filters", client.FeedFilters.TypeName)

	assert.Equal(t, baseURL, client.FilterKeys.BaseUrl)
	assert.Equal(t, "filter-keys", client.FilterKeys.TypeName)

	assert.Equal(t, baseURL, client.FilterOperators.BaseUrl)
	assert.Equal(t, "filter-operators", client.FilterOperators.TypeName)

	assert.Equal(t, baseURL, client.NotificationTypes.BaseUrl)
	assert.Equal(t, "notification-types", client.NotificationTypes.TypeName)

	assert.Equal(t, baseURL, client.Notifiers.BaseUrl)
	assert.Equal(t, "notifiers", client.Notifiers.TypeName)

	assert.Equal(t, baseURL, client.TorrentCategories.BaseUrl)
	assert.Equal(t, "torrent-categories", client.TorrentCategories.TypeName)
}

func TestFeedsClient_CRUD(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	client := NewClient(server.URL)
	ctx := context.Background()

	// Test List - should be empty initially
	feeds, err := client.Feeds.List(ctx)
	require.NoError(t, err)
	assert.Len(t, feeds, 0, "Should start with no feeds")

	// Test Create
	createReq := FeedRequest{
		Name: "Test Feed",
		URL:  "https://example.com/rss",
	}
	created, err := client.Feeds.Create(ctx, createReq)
	require.NoError(t, err)
	assert.NotZero(t, created.ID)
	assert.Equal(t, "Test Feed", created.Name)
	assert.Equal(t, "https://example.com/rss", created.URL)

	// Test Get
	retrieved, err := client.Feeds.Get(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, retrieved.ID)
	assert.Equal(t, created.Name, retrieved.Name)
	assert.Equal(t, created.URL, retrieved.URL)

	// Test List after create
	feeds, err = client.Feeds.List(ctx)
	require.NoError(t, err)
	assert.Len(t, feeds, 1)
	assert.Equal(t, created.ID, feeds[0].ID)

	// Test Delete
	err = client.Feeds.Delete(ctx, created.ID)
	require.NoError(t, err)

	// Verify deletion
	_, err = client.Feeds.Get(ctx, created.ID)
	assert.Error(t, err)

	// Test List after delete
	feeds, err = client.Feeds.List(ctx)
	require.NoError(t, err)
	assert.Len(t, feeds, 0)
}

func TestFilterKeysClient_List(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	client := NewClient(server.URL)
	ctx := context.Background()

	// Test List - should have seeded data
	keys, err := client.FilterKeys.List(ctx)
	require.NoError(t, err)
	assert.Greater(t, len(keys), 0, "Should have seeded filter keys")
}

func TestFilterOperatorsClient_List(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	client := NewClient(server.URL)
	ctx := context.Background()

	// Test List - should have seeded data
	operators, err := client.FilterOperators.List(ctx)
	require.NoError(t, err)
	assert.Greater(t, len(operators), 0, "Should have seeded filter operators")
}

func TestNotificationTypesClient_List(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	client := NewClient(server.URL)
	ctx := context.Background()

	// Test List - should have seeded data
	types, err := client.NotificationTypes.List(ctx)
	require.NoError(t, err)
	assert.Greater(t, len(types), 0, "Should have seeded notification types")
}

func TestFeedFilterSetTypesClient_List(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	client := NewClient(server.URL)
	ctx := context.Background()

	// Test List - should have seeded data
	types, err := client.FeedFilterSetTypes.List(ctx)
	require.NoError(t, err)
	assert.Greater(t, len(types), 0, "Should have seeded feed filter set types")
}

// TestTorrentCategoriesClient tests read-only operations for TorrentCategories
func TestTorrentCategoriesClient(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	client := NewClient(server.URL)
	ctx := context.Background()

	// Test List - should have pre-populated categories
	categories, err := client.TorrentCategories.List(ctx)
	require.NoError(t, err)
	assert.Len(t, categories, 4, "Should have 4 seeded torrent categories")

	// Find a known category for Get test
	var audiobooks TorrentCategoryResponse
	for _, cat := range categories {
		if cat.Name == "audiobooks" {
			audiobooks = cat
			break
		}
	}
	require.NotZero(t, audiobooks.ID, "audiobooks category should exist")
	assert.Equal(t, "family", audiobooks.ScopeName)
	assert.Equal(t, "audiobook", audiobooks.MediaType)

	// Test Get
	retrieved, err := client.TorrentCategories.Get(ctx, audiobooks.ID)
	require.NoError(t, err)
	assert.Equal(t, audiobooks.ID, retrieved.ID)
	assert.Equal(t, audiobooks.Name, retrieved.Name)
	assert.Equal(t, audiobooks.ScopeName, retrieved.ScopeName)
	assert.Equal(t, audiobooks.MediaType, retrieved.MediaType)
}

func TestNotifiersClient_CRUD(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	client := NewClient(server.URL)
	ctx := context.Background()

	// Get discord notification type ID
	types, err := client.NotificationTypes.List(ctx)
	require.NoError(t, err)
	var discordTypeID uint
	for _, nt := range types {
		if nt.Name == "discord" {
			discordTypeID = nt.ID
			break
		}
	}
	require.NotZero(t, discordTypeID, "Discord notification type not found")

	// Test Create
	createReq := NotifierRequest{
		Name:   "test-notifier",
		TypeID: discordTypeID,
		URL:    "https://discord.com/webhook/test",
	}
	created, err := client.Notifiers.Create(ctx, createReq)
	require.NoError(t, err)
	assert.NotZero(t, created.ID)
	assert.Equal(t, "test-notifier", created.Name)
	assert.Equal(t, "discord", created.TypeName)

	// Test Get
	retrieved, err := client.Notifiers.Get(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, retrieved.ID)
	assert.Equal(t, created.Name, retrieved.Name)

	// Test List
	notifiers, err := client.Notifiers.List(ctx)
	require.NoError(t, err)
	assert.Greater(t, len(notifiers), 0)

	// Test Delete
	err = client.Notifiers.Delete(ctx, created.ID)
	require.NoError(t, err)
}

func TestFeedAuthorFiltersClient_CRUD(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	client := NewClient(server.URL)
	ctx := context.Background()

	// Setup: Get discord notification type ID
	types, err := client.NotificationTypes.List(ctx)
	require.NoError(t, err)
	var discordTypeID uint
	for _, nt := range types {
		if nt.Name == "discord" {
			discordTypeID = nt.ID
			break
		}
	}
	require.NotZero(t, discordTypeID)

	// Setup: Get pre-populated category
	category := getTorrentCategoryByName(t, client, ctx, "personal-audiobooks")

	notifier, err := client.Notifiers.Create(ctx, NotifierRequest{
		Name:   "author-test-notifier",
		TypeID: discordTypeID,
		URL:    "https://discord.com/webhook/test",
	})
	require.NoError(t, err)

	feed, err := client.Feeds.Create(ctx, FeedRequest{
		Name: "Author Test Feed",
		URL:  "https://example.com/feed",
	})
	require.NoError(t, err)

	// Test List - should be empty initially
	filters, err := client.FeedAuthorFilters.List(ctx)
	require.NoError(t, err)
	assert.Len(t, filters, 0)

	// Test Create
	createReq := FeedAuthorFilterRequest{
		Author:     "John Doe",
		FeedID:     feed.ID,
		CategoryID: category.ID,
		NotifierID: notifier.ID,
	}
	created, err := client.FeedAuthorFilters.Create(ctx, createReq)
	require.NoError(t, err)
	assert.NotZero(t, created.ID)
	assert.Equal(t, "John Doe", created.Author)
	assert.Equal(t, "Author Test Feed", created.FeedName)
	assert.Equal(t, "personal-audiobooks", created.CategoryName)
	assert.Equal(t, "author-test-notifier", created.NotifierName)

	// Test Get
	retrieved, err := client.FeedAuthorFilters.Get(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, retrieved.ID)
	assert.Equal(t, "John Doe", retrieved.Author)

	// Test List after create
	filters, err = client.FeedAuthorFilters.List(ctx)
	require.NoError(t, err)
	assert.Len(t, filters, 1)

	// Test Delete
	err = client.FeedAuthorFilters.Delete(ctx, created.ID)
	require.NoError(t, err)

	// Verify deletion
	_, err = client.FeedAuthorFilters.Get(ctx, created.ID)
	assert.Error(t, err)
}

func TestCompleteWorkflow(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	client := NewClient(server.URL)
	ctx := context.Background()

	// Get reference data IDs
	types, err := client.NotificationTypes.List(ctx)
	require.NoError(t, err)
	var discordTypeID uint
	for _, nt := range types {
		if nt.Name == "discord" {
			discordTypeID = nt.ID
			break
		}
	}
	require.NotZero(t, discordTypeID)

	setTypes, err := client.FeedFilterSetTypes.List(ctx)
	require.NoError(t, err)
	var anyTypeID uint
	for _, st := range setTypes {
		if st.Name == "any" {
			anyTypeID = st.ID
			break
		}
	}
	require.NotZero(t, anyTypeID)

	keys, err := client.FilterKeys.List(ctx)
	require.NoError(t, err)
	var authorKeyID uint
	for _, k := range keys {
		if k.Name == "author" {
			authorKeyID = k.ID
			break
		}
	}
	require.NotZero(t, authorKeyID)

	operators, err := client.FilterOperators.List(ctx)
	require.NoError(t, err)
	var containsOpID uint
	for _, op := range operators {
		if op.Name == "contains" {
			containsOpID = op.ID
			break
		}
	}
	require.NotZero(t, containsOpID)

	// Step 1: Get a pre-populated TorrentCategory
	category := getTorrentCategoryByName(t, client, ctx, "books")

	// Step 2: Create a Notifier
	notifier, err := client.Notifiers.Create(ctx, NotifierRequest{
		Name:   "workflow-discord",
		TypeID: discordTypeID,
		URL:    "https://discord.com/webhook/test",
	})
	require.NoError(t, err)

	// Step 3: Create a Feed
	feed, err := client.Feeds.Create(ctx, FeedRequest{
		Name: "Workflow Feed",
		URL:  "https://example.com/rss",
	})
	require.NoError(t, err)

	// Step 4: Create a FeedFilter
	filter, err := client.FeedFilters.Create(ctx, FeedFilterRequest{
		Name:       "Workflow Filter",
		FeedID:     feed.ID,
		CategoryID: category.ID,
		NotifierID: notifier.ID,
	})
	require.NoError(t, err)
	assert.Equal(t, "Workflow Filter", filter.Name)
	assert.Equal(t, "Workflow Feed", filter.FeedName)

	// Step 5: Create a FeedFilterSet
	filterSet, err := client.FeedFilterSets.Create(ctx, FeedFilterSetRequest{
		FeedFilterID:        filter.ID,
		FeedFilterSetTypeID: anyTypeID,
	})
	require.NoError(t, err)
	assert.Equal(t, "any", filterSet.TypeName)

	// Step 6: Create FeedFilterSetEntry
	entry, err := client.FeedFilterSetEntries.Create(ctx, FeedFilterSetEntryRequest{
		FeedFilterSetID: filterSet.ID,
		KeyID:           authorKeyID,
		OperatorID:      containsOpID,
		Value:           "Test Author",
	})
	require.NoError(t, err)
	assert.Equal(t, "author", entry.KeyName)
	assert.Equal(t, "contains", entry.OperatorName)
	assert.Equal(t, "Test Author", entry.Value)

	// Verify all data was created
	feedFilters, err := client.FeedFilters.List(ctx)
	require.NoError(t, err)
	assert.Len(t, feedFilters, 1)

	filterSets, err := client.FeedFilterSets.List(ctx)
	require.NoError(t, err)
	assert.Len(t, filterSets, 1)

	entries, err := client.FeedFilterSetEntries.List(ctx)
	require.NoError(t, err)
	assert.Len(t, entries, 1)
}

func TestConnectionError(t *testing.T) {
	// Use an invalid URL to test connection errors
	client := NewClient("http://localhost:1")
	ctx := context.Background()

	_, err := client.Feeds.List(ctx)
	assert.Error(t, err)

	_, err = client.Feeds.Get(ctx, 1)
	assert.Error(t, err)

	_, err = client.Feeds.Create(ctx, FeedRequest{})
	assert.Error(t, err)

	err = client.Feeds.Delete(ctx, 1)
	assert.Error(t, err)
}

func TestContextCancellation(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	client := NewClient(server.URL)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.Feeds.List(ctx)
	assert.Error(t, err)
}
