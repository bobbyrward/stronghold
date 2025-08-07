package cmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/bobbyrward/stronghold/internal/feedwatcher"
	"github.com/bobbyrward/stronghold/internal/models"
)

func createFeedWatcherCmd() *cobra.Command {
	feedWatcherCmd := &cobra.Command{
		Use:  "feed-watcher",
		RunE: runFeedWatcher,
	}

	return feedWatcherCmd
}

func runFeedWatcher(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	
	slog.InfoContext(ctx, "Starting feed watcher command")

	db, err := models.ConnectDB()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to connect to database", slog.Any("err", err))
		return errors.Join(err, fmt.Errorf("failed to connect to database"))
	}

	err = models.AutoMigrate(db)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to auto-migrate database", slog.Any("err", err))
		return errors.Join(err, fmt.Errorf("failed to automigrate database"))
	}

	feedWatcher := feedwatcher.NewFeedWatcher()

	err = feedWatcher.Run(ctx, db)
	if err != nil {
		slog.ErrorContext(ctx, "Feed watcher failed", slog.Any("err", err))
		return errors.Join(err, fmt.Errorf("failed to run feed watcher"))
	}

	slog.InfoContext(ctx, "Feed watcher completed successfully")
	return nil
}
