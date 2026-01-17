package cmd

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/bobbyrward/stronghold/internal/config"
	"github.com/bobbyrward/stronghold/internal/feedwatcher2"
	"github.com/bobbyrward/stronghold/internal/models"
	"github.com/bobbyrward/stronghold/internal/qbit"
)

func createFeedWatcher2Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "feedwatcher2",
		Short: "Run feedwatcher2 to monitor feeds for subscribed authors",
		Long: `Monitors RSS feeds from the database, matches items against author
subscriptions (including aliases), and sends matching torrents to qBittorrent.

This is the database-driven successor to feed-watcher, using Author and
AuthorSubscription models instead of config-based filters.`,
		RunE: runFeedWatcher2,
	}
}

func runFeedWatcher2(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	slog.InfoContext(ctx, "Starting feedwatcher2 command")

	db, err := models.ConnectDB()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to connect to database", slog.Any("error", err))
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	err = models.AutoMigrate(db)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to auto-migrate database", slog.Any("error", err))
		return fmt.Errorf("failed to automigrate database: %w", err)
	}

	qbitClient, err := qbit.CreateClient()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create qBittorrent client", slog.Any("error", err))
		return fmt.Errorf("failed to create qBittorrent client: %w", err)
	}

	fw := feedwatcher2.NewFeedWatcher2(
		db,
		qbitClient,
		config.Config.BookSearch.HttpProxy,
		config.Config.BookSearch.HttpsProxy,
	)

	if err := fw.Run(ctx); err != nil {
		slog.ErrorContext(ctx, "Feedwatcher2 failed", slog.Any("error", err))
		return fmt.Errorf("feedwatcher2 failed: %w", err)
	}

	slog.InfoContext(ctx, "Feedwatcher2 completed successfully")
	return nil
}
