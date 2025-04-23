package cmd

import (
	"context"
	"errors"
	"fmt"

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

	db, err := models.ConnectDB() // Uncomment this line if you want to connect to the database
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to connect to database"))
	}

	err = models.AutoMigrate(db)
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to automigrate database"))
	}

	feedWatcher := feedwatcher.NewFeedWatcher()

	err = feedWatcher.Run(ctx, db)
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to run feed watcher"))
	}

	return nil
}
