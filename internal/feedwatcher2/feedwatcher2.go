package feedwatcher2

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
	"github.com/bobbyrward/stronghold/internal/qbit"
	"github.com/bobbyrward/stronghold/internal/torrentutil"
)

// Media type constants for torrent categorization.
const (
	MediaTypeEbook     = "ebook"
	MediaTypeAudiobook = "audiobook"
)

// extractIDFromGUID extracts the numeric ID from a GUID URL.
// Example: "https://www.myanonamouse.net/t/1213652" returns "1213652".
func extractIDFromGUID(guid string) string {
	lastSlash := strings.LastIndex(guid, "/")
	if lastSlash == -1 || lastSlash == len(guid)-1 {
		return guid
	}
	return guid[lastSlash+1:]
}

// FeedWatcher2 monitors RSS feeds and downloads torrents for subscribed authors.
type FeedWatcher2 struct {
	db                *gorm.DB
	qbitClient        qbit.QbitClient
	torrentDownloader *torrentutil.TorrentDownloader
	authorMatcher     *AuthorMatcher
}

// NewFeedWatcher2 creates a new FeedWatcher2 instance.
func NewFeedWatcher2(db *gorm.DB, qbitClient qbit.QbitClient, httpProxy, httpsProxy string) *FeedWatcher2 {
	return &FeedWatcher2{
		db:                db,
		qbitClient:        qbitClient,
		torrentDownloader: torrentutil.NewTorrentDownloader(httpProxy, httpsProxy),
		authorMatcher:     NewAuthorMatcher(db),
	}
}

// Run executes the feed watcher, processing all feeds from the database.
func (fw *FeedWatcher2) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "Starting feedwatcher2")

	// Load subscriptions into the matcher
	err := fw.authorMatcher.LoadSubscriptions(ctx)
	if err != nil {
		return fmt.Errorf("failed to load subscriptions: %w", err)
	}

	// Query all feeds from database
	var feeds []models.Feed
	result := fw.db.Find(&feeds)
	if result.Error != nil {
		slog.ErrorContext(ctx, "Failed to query feeds", slog.Any("error", result.Error))
		return fmt.Errorf("failed to query feeds: %w", result.Error)
	}

	slog.InfoContext(ctx, "Found feeds to process", slog.Int("count", len(feeds)))

	// Process each feed
	var errs []error
	for _, feed := range feeds {
		err := fw.watchFeed(ctx, &feed)
		if err != nil {
			slog.WarnContext(ctx, "Error processing feed",
				slog.String("feed_name", feed.Name),
				slog.Any("error", err))
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	slog.InfoContext(ctx, "Feedwatcher2 completed successfully")
	return nil
}

// watchFeed processes a single RSS feed.
func (fw *FeedWatcher2) watchFeed(ctx context.Context, feed *models.Feed) error {
	slog.InfoContext(ctx, "Processing feed",
		slog.String("name", feed.Name),
		slog.String("url", feed.URL))

	parser := gofeed.NewParser()
	parsedFeed, err := parser.ParseURL(feed.URL)
	if err != nil {
		return fmt.Errorf("failed to parse feed %s: %w", feed.Name, err)
	}

	slog.InfoContext(ctx, "Parsed feed",
		slog.String("name", feed.Name),
		slog.Int("items", len(parsedFeed.Items)))

	for _, item := range parsedFeed.Items {
		err := fw.processItem(ctx, feed, item)
		if err != nil {
			slog.WarnContext(ctx, "Error processing feed item",
				slog.String("feed_name", feed.Name),
				slog.String("item_title", item.Title),
				slog.Any("error", err))
			// Continue processing other items
		}
	}

	return nil
}

// processItem processes a single feed item.
func (fw *FeedWatcher2) processItem(ctx context.Context, feed *models.Feed, item *gofeed.Item) error {
	// Parse the description to extract metadata
	entry, err := parseDescription(ctx, item.Description)
	if err != nil {
		return fmt.Errorf("failed to parse description: %w", err)
	}

	// Set fields from the feed item
	entry.Guid = item.GUID
	entry.Link = item.Link
	entry.Title = item.Title

	// Extract the ID from the GUID URL for deduplication and storage
	booksearchID := extractIDFromGUID(item.GUID)

	// Find matching subscription
	subscription := fw.authorMatcher.FindMatchingSubscription(entry.Authors)
	if subscription == nil {
		// No match, skip this item
		return nil
	}

	slog.InfoContext(ctx, "Found matching subscription",
		slog.String("title", entry.Title),
		slog.String("author", subscription.Author.Name),
		slog.String("scope", subscription.Scope.Name),
		slog.Any("feed_authors", entry.Authors))

	// Determine the torrent category based on scope and media type
	category, err := fw.determineTorrentCategory(ctx, &subscription.Scope, entry.Category)
	if err != nil {
		return fmt.Errorf("failed to determine torrent category: %w", err)
	}

	slog.InfoContext(ctx, "Determined torrent category",
		slog.String("category", category.Name),
		slog.String("feed_category", entry.Category))

	// Check for duplicate by booksearch ID before downloading
	var existingItem models.AuthorSubscriptionItem
	result := fw.db.Where("booksearch_id = ?", booksearchID).First(&existingItem)
	if result.Error == nil {
		// Already exists, skip
		slog.InfoContext(ctx, "Item already downloaded, skipping",
			slog.String("booksearch_id", booksearchID),
			slog.String("title", entry.Title))
		return nil
	}
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check for existing item: %w", result.Error)
	}

	// Download torrent and extract hash
	hash, err := fw.torrentDownloader.DownloadAndHash(ctx, entry.Link)
	if err != nil {
		return fmt.Errorf("failed to download torrent: %w", err)
	}

	slog.InfoContext(ctx, "Downloaded torrent",
		slog.String("hash", hash),
		slog.String("title", entry.Title))

	// Add torrent to qBittorrent
	err = fw.qbitClient.AddTorrentFromUrlCtx(
		ctx,
		entry.Link,
		map[string]string{
			"autoTMM":  "true",
			"category": category.Name,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to add torrent to qBittorrent: %w", err)
	}

	slog.InfoContext(ctx, "Added torrent to qBittorrent",
		slog.String("title", entry.Title),
		slog.String("category", category.Name),
		slog.String("hash", hash))

	// Create AuthorSubscriptionItem record
	subscriptionItem := models.AuthorSubscriptionItem{
		AuthorSubscriptionID: subscription.ID,
		TorrentHash:          hash,
		BooksearchID:         booksearchID,
		DownloadedAt:         time.Now(),
	}

	result = fw.db.Create(&subscriptionItem)
	if result.Error != nil {
		slog.ErrorContext(ctx, "Failed to create subscription item record",
			slog.String("hash", hash),
			slog.Any("error", result.Error))
		// Don't return error - torrent was already added to qBittorrent
	}

	// Send notification
	payload := CreateFeedwatcher2NotificationPayload(&entry, &subscription.Author, subscription)
	err = SendNotificationViaNotifier(ctx, subscription.Notifier, payload)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to send notification",
			slog.String("title", entry.Title),
			slog.Any("error", err))
		// Don't return error - torrent was already added
	}

	slog.InfoContext(ctx, "Successfully processed feed item",
		slog.String("title", entry.Title),
		slog.String("author", subscription.Author.Name),
		slog.String("category", category.Name),
		slog.String("hash", hash))

	return nil
}

// determineTorrentCategory determines the appropriate TorrentCategory based on scope and media type.
func (fw *FeedWatcher2) determineTorrentCategory(ctx context.Context, scope *models.SubscriptionScope, feedCategory string) (*models.TorrentCategory, error) {
	// Determine media type from feed category
	mediaType := MediaTypeEbook
	if strings.HasPrefix(feedCategory, "Audiobooks") {
		mediaType = MediaTypeAudiobook
	}

	slog.DebugContext(ctx, "Determining torrent category",
		slog.String("scope", scope.Name),
		slog.String("feed_category", feedCategory),
		slog.String("media_type", mediaType))

	// Query TorrentCategory
	var category models.TorrentCategory
	result := fw.db.Where("scope_id = ? AND media_type = ?", scope.ID, mediaType).First(&category)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no torrent category found for scope=%s, media_type=%s", scope.Name, mediaType)
		}
		return nil, result.Error
	}

	return &category, nil
}
