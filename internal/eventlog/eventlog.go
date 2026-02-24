package eventlog

import (
	"encoding/json"
	"log/slog"

	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
)

// Categories
const (
	CategoryDownload     = "download"
	CategoryImport       = "import"
	CategoryNotification = "notification"
	CategorySubscription = "subscription"
	CategorySearch       = "search"
	CategoryFeed         = "feed"
	CategoryMutation     = "mutation"
)

// Event types
const (
	// Download events
	EventTorrentAdded           = "torrent.added"
	EventTorrentDuplicateSkipped = "torrent.duplicate_skipped"

	// Import events
	EventImportStarted            = "import.started"
	EventImportCompleted          = "import.completed"
	EventImportFailed             = "import.failed"
	EventImportManualIntervention = "import.manual_intervention"

	// Notification events
	EventNotificationSent   = "notification.sent"
	EventNotificationFailed = "notification.failed"

	// Subscription events
	EventSubscriptionMatched = "subscription.matched"

	// Search events
	EventSearchRequested = "search.requested"
	EventSearchCompleted = "search.completed"
	EventSearchHardcover = "search.hardcover"

	// Feed events
	EventFeedPolled = "feed.polled"
	EventFeedError  = "feed.error"

	// Mutation events
	EventCreated = "created"
	EventUpdated = "updated"
	EventDeleted = "deleted"
)

// Sources
const (
	SourceFeedwatcher2              = "feedwatcher2"
	SourceDiscordBot                = "discord-bot"
	SourceAPI                       = "api"
	SourceEbookImporter             = "ebook-importer"
	SourceAudiobookImporter         = "audiobook-importer"
	SourceAuthorSubscriptionImporter = "author-subscription-importer"
)

// Entity types
const (
	EntityTorrent      = "torrent"
	EntityAuthor       = "author"
	EntityAuthorAlias  = "author_alias"
	EntitySubscription = "subscription"
	EntityFeed         = "feed"
	EntityNotifier     = "notifier"
	EntityLibrary      = "library"
	EntityFeedAuthorFilter = "feed_author_filter"
	EntitySearch       = "search"
)

// Log creates an event log entry. It is fire-and-forget: errors are logged but never returned.
// If db is nil, the event is silently skipped.
func Log(db *gorm.DB, category, eventType, source, entityType, entityID, summary string, details any) {
	if db == nil {
		return
	}

	var detailsJSON string
	if details != nil {
		b, err := json.Marshal(details)
		if err != nil {
			slog.Error("Failed to marshal event log details",
				slog.String("category", category),
				slog.String("event_type", eventType),
				slog.Any("error", err))
			detailsJSON = "{}"
		} else {
			detailsJSON = string(b)
		}
	}

	entry := models.EventLog{
		Category:   category,
		EventType:  eventType,
		Source:     source,
		EntityType: entityType,
		EntityID:   entityID,
		Summary:    summary,
		Details:    detailsJSON,
	}

	if err := db.Create(&entry).Error; err != nil {
		slog.Error("Failed to create event log entry",
			slog.String("category", category),
			slog.String("event_type", eventType),
			slog.String("source", source),
			slog.Any("error", err))
	}
}
