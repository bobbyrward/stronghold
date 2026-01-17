package feedwatcher2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jackpal/bencode-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
	"github.com/bobbyrward/stronghold/internal/notifications"
	"github.com/bobbyrward/stronghold/internal/qbit"
	"github.com/bobbyrward/stronghold/internal/testutil"
	"github.com/bobbyrward/stronghold/internal/torrentutil"
)

// createTestTorrent creates a valid torrent file bytes for testing.
func createTestTorrentBytes(t *testing.T, name string) []byte {
	t.Helper()

	torrent := map[string]interface{}{
		"info": map[string]interface{}{
			"name":         name,
			"piece length": 262144,
			"pieces":       "12345678901234567890",
			"length":       1024,
		},
		"announce": "http://tracker.example.com/announce",
	}

	var buf bytes.Buffer
	err := bencode.Marshal(&buf, torrent)
	require.NoError(t, err)

	return buf.Bytes()
}

// createTestFeedWatcher creates a FeedWatcher2 for testing with a no-proxy TorrentDownloader.
func createTestFeedWatcher(db *gorm.DB, qbitClient qbit.QbitClient) *FeedWatcher2 {
	return &FeedWatcher2{
		db:                db,
		qbitClient:        qbitClient,
		torrentDownloader: torrentutil.NewTestTorrentDownloader(),
		authorMatcher:     NewAuthorMatcher(db),
	}
}

// createMockRSSFeed creates an RSS feed XML string with the given items.
func createMockRSSFeed(items []mockFeedItem) string {
	itemsXML := ""
	for _, item := range items {
		itemsXML += fmt.Sprintf(`
		<item>
			<guid>%s</guid>
			<title>%s</title>
			<link>%s</link>
			<description>Author(s): %s&lt;br/&gt;Category: %s&lt;br/&gt;Description: Test description</description>
		</item>`, item.GUID, item.Title, item.Link, item.Author, item.Category)
	}

	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
	<channel>
		<title>Test Feed</title>
		<link>http://example.com</link>
		<description>Test RSS Feed</description>
		%s
	</channel>
</rss>`, itemsXML)
}

type mockFeedItem struct {
	GUID     string
	Title    string
	Link     string
	Author   string
	Category string
}

func setupIntegrationTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := models.ConnectTestDB()
	require.NoError(t, err)

	return db
}

func createTestScope(t *testing.T, db *gorm.DB, name string) *models.SubscriptionScope {
	t.Helper()

	var scope models.SubscriptionScope
	err := db.FirstOrCreate(&scope, models.SubscriptionScope{Name: name}).Error
	require.NoError(t, err)
	return &scope
}

func createTestTorrentCategory(t *testing.T, db *gorm.DB, name string, scopeID uint, mediaType string) *models.TorrentCategory {
	t.Helper()

	var category models.TorrentCategory
	err := db.FirstOrCreate(&category, models.TorrentCategory{
		Name:      name,
		ScopeID:   scopeID,
		MediaType: mediaType,
	}).Error
	require.NoError(t, err)
	return &category
}

func TestWatchFeed_NoSubscriptions(t *testing.T) {
	db := setupIntegrationTestDB(t)

	// Create RSS server
	rssServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		feed := createMockRSSFeed([]mockFeedItem{
			{GUID: "https://www.myanonamouse.net/t/1000", Title: "Test Book", Link: "http://torrent.example.com/1", Author: "Unknown Author", Category: "Audiobooks"},
		})
		w.Header().Set("Content-Type", "application/rss+xml")
		_, _ = w.Write([]byte(feed))
	}))
	defer rssServer.Close()

	// Create feed in DB
	feed := models.Feed{Name: "Test Feed", URL: rssServer.URL}
	err := db.Create(&feed).Error
	require.NoError(t, err)

	// Create mock qBit client
	mockQbit := &testutil.MockQbitClient{}

	fw := createTestFeedWatcher(db, mockQbit)
	err = fw.Run(context.Background())
	require.NoError(t, err)

	// No torrents should be added since no subscriptions match
	assert.Empty(t, mockQbit.AddTorrentFromUrlCtxCalls)
}

func TestWatchFeed_MatchByAuthorName(t *testing.T) {
	db := setupIntegrationTestDB(t)

	// Create scope and category
	scope := createTestScope(t, db, "personal")
	createTestTorrentCategory(t, db, "personal-audiobooks", scope.ID, "audiobook")

	// Create author and subscription
	author := models.Author{Name: "Brandon Sanderson"}
	err := db.Create(&author).Error
	require.NoError(t, err)

	subscription := models.AuthorSubscription{AuthorID: author.ID, ScopeID: scope.ID}
	err = db.Create(&subscription).Error
	require.NoError(t, err)

	// Create torrent server
	torrentData := createTestTorrentBytes(t, "test-book.mp3")
	torrentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-bittorrent")
		_, _ = w.Write(torrentData)
	}))
	defer torrentServer.Close()

	// Create RSS server
	rssServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		feed := createMockRSSFeed([]mockFeedItem{
			{GUID: "https://www.myanonamouse.net/t/1001", Title: "Mistborn", Link: torrentServer.URL + "/mistborn.torrent", Author: "Brandon Sanderson", Category: "Audiobooks - Fantasy"},
		})
		w.Header().Set("Content-Type", "application/rss+xml")
		_, _ = w.Write([]byte(feed))
	}))
	defer rssServer.Close()

	// Create feed in DB
	feed := models.Feed{Name: "Test Feed", URL: rssServer.URL}
	err = db.Create(&feed).Error
	require.NoError(t, err)

	// Create mock qBit client
	mockQbit := &testutil.MockQbitClient{}

	fw := createTestFeedWatcher(db, mockQbit)
	err = fw.Run(context.Background())
	require.NoError(t, err)

	// Should have added one torrent
	require.Len(t, mockQbit.AddTorrentFromUrlCtxCalls, 1)
	assert.Equal(t, "personal-audiobooks", mockQbit.AddTorrentFromUrlCtxCalls[0].Options["category"])

	// Should have created subscription item with extracted ID
	var items []models.AuthorSubscriptionItem
	err = db.Find(&items).Error
	require.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Equal(t, "1001", items[0].BooksearchID)
}

func TestWatchFeed_MatchByAlias(t *testing.T) {
	db := setupIntegrationTestDB(t)

	// Create scope and category
	scope := createTestScope(t, db, "family")
	createTestTorrentCategory(t, db, "books", scope.ID, "ebook")

	// Create author with alias
	author := models.Author{Name: "JF Brink"}
	err := db.Create(&author).Error
	require.NoError(t, err)

	alias := models.AuthorAlias{AuthorID: author.ID, Name: "J.F. Brink"}
	err = db.Create(&alias).Error
	require.NoError(t, err)

	subscription := models.AuthorSubscription{AuthorID: author.ID, ScopeID: scope.ID}
	err = db.Create(&subscription).Error
	require.NoError(t, err)

	// Create torrent server
	torrentData := createTestTorrentBytes(t, "test-ebook.epub")
	torrentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-bittorrent")
		_, _ = w.Write(torrentData)
	}))
	defer torrentServer.Close()

	// Create RSS server - note author has dots
	rssServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		feed := createMockRSSFeed([]mockFeedItem{
			{GUID: "https://www.myanonamouse.net/t/1002", Title: "Some Ebook", Link: torrentServer.URL + "/ebook.torrent", Author: "J.F. Brink", Category: "Ebooks - Fiction"},
		})
		w.Header().Set("Content-Type", "application/rss+xml")
		_, _ = w.Write([]byte(feed))
	}))
	defer rssServer.Close()

	// Create feed
	feed := models.Feed{Name: "Ebook Feed", URL: rssServer.URL}
	err = db.Create(&feed).Error
	require.NoError(t, err)

	mockQbit := &testutil.MockQbitClient{}

	fw := createTestFeedWatcher(db, mockQbit)
	err = fw.Run(context.Background())
	require.NoError(t, err)

	// Should match via alias
	require.Len(t, mockQbit.AddTorrentFromUrlCtxCalls, 1)
	assert.Equal(t, "books", mockQbit.AddTorrentFromUrlCtxCalls[0].Options["category"])
}

func TestWatchFeed_Deduplication(t *testing.T) {
	db := setupIntegrationTestDB(t)

	// Create scope and category
	scope := createTestScope(t, db, "personal")
	createTestTorrentCategory(t, db, "personal-audiobooks", scope.ID, "audiobook")

	// Create author and subscription
	author := models.Author{Name: "Test Author"}
	err := db.Create(&author).Error
	require.NoError(t, err)

	subscription := models.AuthorSubscription{AuthorID: author.ID, ScopeID: scope.ID}
	err = db.Create(&subscription).Error
	require.NoError(t, err)

	// Pre-create an existing subscription item with the same ID (extracted from GUID URL)
	// Deduplication happens by ID before downloading the torrent
	existingItem := models.AuthorSubscriptionItem{
		AuthorSubscriptionID: subscription.ID,
		TorrentHash:          "previoushash1234567890abcdef12345678",
		BooksearchID:         "1003",
		DownloadedAt:         time.Now().Add(-time.Hour),
	}
	err = db.Create(&existingItem).Error
	require.NoError(t, err)

	// Create RSS server - no torrent server needed since dedup happens before download
	rssServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		feed := createMockRSSFeed([]mockFeedItem{
			{GUID: "https://www.myanonamouse.net/t/1003", Title: "Same Item", Link: "http://torrent.example.com/test.torrent", Author: "Test Author", Category: "Audiobooks"},
		})
		w.Header().Set("Content-Type", "application/rss+xml")
		_, _ = w.Write([]byte(feed))
	}))
	defer rssServer.Close()

	feed := models.Feed{Name: "Test Feed", URL: rssServer.URL}
	err = db.Create(&feed).Error
	require.NoError(t, err)

	mockQbit := &testutil.MockQbitClient{}

	fw := createTestFeedWatcher(db, mockQbit)
	err = fw.Run(context.Background())
	require.NoError(t, err)

	// Should NOT add torrent since ID already exists
	assert.Empty(t, mockQbit.AddTorrentFromUrlCtxCalls)

	// Should still only have one item
	var items []models.AuthorSubscriptionItem
	err = db.Find(&items).Error
	require.NoError(t, err)
	assert.Len(t, items, 1)
}

func TestWatchFeed_CorrectCategory_PersonalAudiobook(t *testing.T) {
	db := setupIntegrationTestDB(t)

	scope := createTestScope(t, db, "personal")
	createTestTorrentCategory(t, db, "personal-audiobooks", scope.ID, "audiobook")

	author := models.Author{Name: "Audio Author"}
	err := db.Create(&author).Error
	require.NoError(t, err)

	subscription := models.AuthorSubscription{AuthorID: author.ID, ScopeID: scope.ID}
	err = db.Create(&subscription).Error
	require.NoError(t, err)

	torrentData := createTestTorrentBytes(t, "audiobook.mp3")
	torrentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(torrentData)
	}))
	defer torrentServer.Close()

	rssServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		feed := createMockRSSFeed([]mockFeedItem{
			{GUID: "https://www.myanonamouse.net/t/1004", Title: "Audiobook", Link: torrentServer.URL, Author: "Audio Author", Category: "Audiobooks - Fantasy"},
		})
		_, _ = w.Write([]byte(feed))
	}))
	defer rssServer.Close()

	feed := models.Feed{Name: "Feed", URL: rssServer.URL}
	err = db.Create(&feed).Error
	require.NoError(t, err)

	mockQbit := &testutil.MockQbitClient{}
	fw := createTestFeedWatcher(db, mockQbit)
	err = fw.Run(context.Background())
	require.NoError(t, err)

	require.Len(t, mockQbit.AddTorrentFromUrlCtxCalls, 1)
	assert.Equal(t, "personal-audiobooks", mockQbit.AddTorrentFromUrlCtxCalls[0].Options["category"])
}

func TestWatchFeed_CorrectCategory_FamilyEbook(t *testing.T) {
	db := setupIntegrationTestDB(t)

	scope := createTestScope(t, db, "family")
	createTestTorrentCategory(t, db, "books", scope.ID, "ebook")

	author := models.Author{Name: "Ebook Author"}
	err := db.Create(&author).Error
	require.NoError(t, err)

	subscription := models.AuthorSubscription{AuthorID: author.ID, ScopeID: scope.ID}
	err = db.Create(&subscription).Error
	require.NoError(t, err)

	torrentData := createTestTorrentBytes(t, "ebook.epub")
	torrentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(torrentData)
	}))
	defer torrentServer.Close()

	rssServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		feed := createMockRSSFeed([]mockFeedItem{
			{GUID: "https://www.myanonamouse.net/t/1005", Title: "Ebook", Link: torrentServer.URL, Author: "Ebook Author", Category: "Ebooks - Romance"},
		})
		_, _ = w.Write([]byte(feed))
	}))
	defer rssServer.Close()

	feed := models.Feed{Name: "Feed", URL: rssServer.URL}
	err = db.Create(&feed).Error
	require.NoError(t, err)

	mockQbit := &testutil.MockQbitClient{}
	fw := createTestFeedWatcher(db, mockQbit)
	err = fw.Run(context.Background())
	require.NoError(t, err)

	require.Len(t, mockQbit.AddTorrentFromUrlCtxCalls, 1)
	assert.Equal(t, "books", mockQbit.AddTorrentFromUrlCtxCalls[0].Options["category"])
}

// E2E Test

func TestFeedWatcher2_E2E(t *testing.T) {
	db := setupIntegrationTestDB(t)

	// Setup reference data
	personalScope := createTestScope(t, db, "personal")
	createTestTorrentCategory(t, db, "personal-audiobooks", personalScope.ID, "audiobook")

	// Create author with alias
	author := models.Author{Name: "JF Brink"}
	err := db.Create(&author).Error
	require.NoError(t, err)

	alias := models.AuthorAlias{AuthorID: author.ID, Name: "J.F. Brink"}
	err = db.Create(&alias).Error
	require.NoError(t, err)

	// Create notification server
	var receivedNotification bool
	var receivedMessage notifications.DiscordWebhookMessage
	notifyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedNotification = true
		_ = json.NewDecoder(r.Body).Decode(&receivedMessage)
		w.WriteHeader(http.StatusOK)
	}))
	defer notifyServer.Close()

	// Create notifier
	notificationType := models.NotificationType{Name: "discord"}
	err = db.FirstOrCreate(&notificationType, models.NotificationType{Name: "discord"}).Error
	require.NoError(t, err)

	notifier := models.Notifier{
		Name:               "test-notifier",
		NotificationTypeID: notificationType.ID,
		URL:                notifyServer.URL,
	}
	err = db.Create(&notifier).Error
	require.NoError(t, err)

	// Create subscription with notifier
	subscription := models.AuthorSubscription{
		AuthorID:   author.ID,
		ScopeID:    personalScope.ID,
		NotifierID: &notifier.ID,
	}
	err = db.Create(&subscription).Error
	require.NoError(t, err)

	// Create torrent server
	torrentData := createTestTorrentBytes(t, "e2e-test-book.mp3")
	torrentServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-bittorrent")
		_, _ = w.Write(torrentData)
	}))
	defer torrentServer.Close()

	// Create RSS server with item matching the alias
	rssServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		feed := createMockRSSFeed([]mockFeedItem{
			{
				GUID:     "https://www.myanonamouse.net/t/1006",
				Title:    "E2E Test Book",
				Link:     torrentServer.URL + "/book.torrent",
				Author:   "J.F. Brink", // Uses dotted alias
				Category: "Audiobooks - Fiction",
			},
		})
		w.Header().Set("Content-Type", "application/rss+xml")
		_, _ = w.Write([]byte(feed))
	}))
	defer rssServer.Close()

	// Create feed pointing to RSS server
	feed := models.Feed{Name: "E2E Test Feed", URL: rssServer.URL}
	err = db.Create(&feed).Error
	require.NoError(t, err)

	// Create mock qBit client
	mockQbit := &testutil.MockQbitClient{}

	// Run feedwatcher2
	fw := createTestFeedWatcher(db, mockQbit)
	err = fw.Run(context.Background())
	require.NoError(t, err)

	// Verify torrent added to qBittorrent with correct category
	require.Len(t, mockQbit.AddTorrentFromUrlCtxCalls, 1)
	assert.Equal(t, "personal-audiobooks", mockQbit.AddTorrentFromUrlCtxCalls[0].Options["category"])
	assert.Contains(t, mockQbit.AddTorrentFromUrlCtxCalls[0].URL, torrentServer.URL)

	// Verify AuthorSubscriptionItem created with extracted ID
	var items []models.AuthorSubscriptionItem
	err = db.Find(&items).Error
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.NotEmpty(t, items[0].TorrentHash)
	assert.Equal(t, "1006", items[0].BooksearchID)
	assert.Equal(t, subscription.ID, items[0].AuthorSubscriptionID)

	// Verify notification sent
	assert.True(t, receivedNotification, "Notification should have been sent")
	assert.Equal(t, "Stronghold", receivedMessage.Username)
	assert.Len(t, receivedMessage.Embeds, 1)
	assert.Equal(t, "E2E Test Book", receivedMessage.Embeds[0].Title)
	assert.Equal(t, "Feedwatcher2", receivedMessage.Embeds[0].Author.Name)
}
