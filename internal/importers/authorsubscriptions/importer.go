package authorsubscriptions

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/autobrr/go-qbittorrent"
	"github.com/cappuccinotm/slogx"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/config"
	"github.com/bobbyrward/stronghold/internal/feedwatcher2"
	"github.com/bobbyrward/stronghold/internal/importers/audiobooks"
	"github.com/bobbyrward/stronghold/internal/importers/ebooks"
	"github.com/bobbyrward/stronghold/internal/models"
	"github.com/bobbyrward/stronghold/internal/notifications"
	"github.com/bobbyrward/stronghold/internal/qbit"
)

// AuthorSubscriptionImporter handles importing torrents from the author-subscriptions category.
// It looks up the AuthorSubscriptionItem to determine book type and destination library.
type AuthorSubscriptionImporter struct {
	db              *gorm.DB
	qbitClient      qbit.QbitClient
	audiobookSystem *audiobooks.AudiobookImporterSystem
	ebookSystem     *ebooks.BookImporterSystem
}

// NewAuthorSubscriptionImporter creates a new AuthorSubscriptionImporter.
func NewAuthorSubscriptionImporter(
	db *gorm.DB,
	qbitClient qbit.QbitClient,
	audiobookSystem *audiobooks.AudiobookImporterSystem,
	ebookSystem *ebooks.BookImporterSystem,
) *AuthorSubscriptionImporter {
	return &AuthorSubscriptionImporter{
		db:              db,
		qbitClient:      qbitClient,
		audiobookSystem: audiobookSystem,
		ebookSystem:     ebookSystem,
	}
}

// Run processes all unimported torrents in the author-subscriptions category.
func (asi *AuthorSubscriptionImporter) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "Running author subscription import process...")

	torrents, err := qbit.GetUnimportedTorrentsByCategory(
		ctx,
		asi.qbitClient,
		feedwatcher2.AuthorSubscriptionCategory,
	)
	if err != nil {
		return fmt.Errorf("failed to get unimported torrents for author-subscriptions: %w", err)
	}

	slog.InfoContext(ctx, "Found unimported author subscription torrents",
		slog.Int("count", len(torrents)))

	for _, torrent := range torrents {
		err := asi.importTorrent(ctx, torrent)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to import author subscription torrent",
				slog.String("name", torrent.Name),
				slog.String("hash", torrent.Hash),
				slogx.Error(err))
			// Continue with other torrents
		}
	}

	return nil
}

// importTorrent imports a single torrent from the author-subscriptions category.
func (asi *AuthorSubscriptionImporter) importTorrent(ctx context.Context, torrent qbittorrent.Torrent) error {
	slog.InfoContext(ctx, "Processing author subscription torrent",
		slog.String("name", torrent.Name),
		slog.String("hash", torrent.Hash))

	// Look up the AuthorSubscriptionItem by torrent hash
	var item models.AuthorSubscriptionItem
	result := asi.db.
		Preload("BookType").
		Preload("AuthorSubscription").
		Preload("AuthorSubscription.EbookLibrary").
		Preload("AuthorSubscription.AudiobookLibrary").
		Preload("AuthorSubscription.Notifier").
		Where("torrent_hash = ?", torrent.Hash).
		First(&item)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			slog.WarnContext(ctx, "No AuthorSubscriptionItem found for torrent, marking for manual intervention",
				slog.String("name", torrent.Name),
				slog.String("hash", torrent.Hash))
			// No item found, so no notifier available - use empty string to skip notification
			return asi.markForManualIntervention(ctx, torrent, "", "No AuthorSubscriptionItem found for torrent hash")
		}
		return fmt.Errorf("failed to lookup AuthorSubscriptionItem: %w", result.Error)
	}

	slog.InfoContext(ctx, "Found AuthorSubscriptionItem",
		slog.String("title", item.Title),
		slog.String("book_type", item.BookType.Name),
		slog.Uint64("subscription_id", uint64(item.AuthorSubscriptionID)))

	// Determine the destination library based on book type
	var library *models.Library
	switch item.BookType.Name {
	case "audiobook":
		library = &item.AuthorSubscription.AudiobookLibrary
	case "ebook":
		library = &item.AuthorSubscription.EbookLibrary
	default:
		return fmt.Errorf("unknown book type: %s", item.BookType.Name)
	}

	slog.InfoContext(ctx, "Using library for import",
		slog.String("library_name", library.Name),
		slog.String("library_path", library.Path),
		slog.String("book_type", item.BookType.Name))

	// Create config adapters for the import systems
	importLibrary := &config.ImportLibrary{
		Name: library.Name,
		Path: library.Path,
	}

	// Get notifier name if set
	notifierName := ""
	if item.AuthorSubscription.Notifier != nil {
		notifierName = item.AuthorSubscription.Notifier.Name
	}

	importType := config.ImportType{
		Category:        feedwatcher2.AuthorSubscriptionCategory,
		Library:         library.Name,
		DiscordNotifier: notifierName,
	}

	// Route to appropriate importer based on book type
	switch item.BookType.Name {
	case "audiobook":
		asi.audiobookSystem.ImportTorrentWithLibrary(ctx, torrent, importType, importLibrary)
	case "ebook":
		asi.ebookSystem.ImportTorrent(ctx, torrent, importType, importLibrary)
	}

	return nil
}

// markForManualIntervention tags a torrent for manual handling and sends a notification.
func (asi *AuthorSubscriptionImporter) markForManualIntervention(ctx context.Context, torrent qbittorrent.Torrent, notifierName string, reason string) error {
	err := qbit.TagTorrent(ctx, asi.qbitClient, torrent, config.Config.Importers.ManualInterventionTag)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to add manual intervention tag",
			slog.String("name", torrent.Name),
			slog.String("hash", torrent.Hash),
			slogx.Error(err))
		return err
	}

	slog.InfoContext(ctx, "Marked torrent for manual intervention",
		slog.String("name", torrent.Name),
		slog.String("hash", torrent.Hash),
		slog.String("reason", reason))

	// Send Discord notification if notifier is configured
	asi.sendManualInterventionNotification(ctx, torrent, notifierName, reason)

	return nil
}

// sendManualInterventionNotification sends a Discord notification about a torrent requiring manual intervention.
func (asi *AuthorSubscriptionImporter) sendManualInterventionNotification(ctx context.Context, torrent qbittorrent.Torrent, notifierName string, reason string) {
	if notifierName == "" {
		return
	}

	message := notifications.DiscordWebhookMessage{
		Username: "Stronghold Author Subscription Importer",
		Embeds: []notifications.DiscordEmbed{
			{
				Title:       "⚠️ Manual Intervention Required",
				Description: fmt.Sprintf("Torrent **%s** requires manual intervention", torrent.Name),
				Color:       0xFFA500, // Orange color
				Fields: []notifications.DiscordEmbedField{
					{
						Name:   "Reason",
						Value:  reason,
						Inline: false,
					},
					{
						Name:   "Torrent Hash",
						Value:  torrent.Hash,
						Inline: true,
					},
					{
						Name:   "Category",
						Value:  feedwatcher2.AuthorSubscriptionCategory,
						Inline: true,
					},
				},
			},
		},
	}

	err := notifications.SendNotification(ctx, notifierName, message)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to send manual intervention notification",
			slog.String("torrent", torrent.Name),
			slogx.Error(err))
	}
}
