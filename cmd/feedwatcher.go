package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bobbyrward/stronghold/internal/feedwatcher"
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

	feedWatcher := feedwatcher.NewFeedWatcher()

	err := feedWatcher.Run(ctx)
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to run feed watcher"))
	}

	return nil
}
