package cmd

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/cappuccinotm/slogx"
	"github.com/spf13/cobra"

	"github.com/bobbyrward/stronghold/internal/config"
	"github.com/bobbyrward/stronghold/internal/importers/audiobooks"
	"github.com/bobbyrward/stronghold/internal/importers/audiobooks/audible"
	"github.com/bobbyrward/stronghold/internal/importers/audiobooks/metadata"
	"github.com/bobbyrward/stronghold/internal/importers/authorsubscriptions"
	"github.com/bobbyrward/stronghold/internal/importers/ebooks"
	"github.com/bobbyrward/stronghold/internal/models"
	"github.com/bobbyrward/stronghold/internal/qbit"
)

func createAuthorSubscriptionImporterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "author-subscription-importer",
		Short: "Import torrents from author subscriptions",
		Long:  "Imports torrents in the author-subscriptions category, looking up the destination library from the AuthorSubscriptionItem record.",
		RunE:  runAuthorSubscriptionImporter,
	}

	return cmd
}

func runAuthorSubscriptionImporter(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	slog.InfoContext(ctx, "Starting author subscription import command")

	// Connect to database
	db, err := models.ConnectAndMigrate(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to connect to database", slogx.Error(err))
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Create qBittorrent client
	qbitClient, err := qbit.CreateClient()
	if err != nil {
		slog.ErrorContext(ctx, "failed to create qBittorrent client", slogx.Error(err))
		return fmt.Errorf("failed to create qBittorrent client: %w", err)
	}

	// Create audiobook importer system (for audiobook imports)
	audibleApiClient := audible.NewAudibleApiClient()
	audiobookSystem, err := audiobooks.NewAudiobookImporterSystem(
		qbitClient,
		config.Config.Importers,
		metadata.NewFFProbeMetadataProvider(),
		audibleApiClient,
	)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create audiobook importer system", slogx.Error(err))
		return fmt.Errorf("failed to create audiobook importer system: %w", err)
	}

	// Create ebook importer system (for ebook imports)
	ebookSystem := ebooks.NewBookImporterSystem(qbitClient)

	// Create and run the author subscription importer
	importer := authorsubscriptions.NewAuthorSubscriptionImporter(
		db,
		qbitClient,
		audiobookSystem,
		ebookSystem,
	)

	err = importer.Run(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to run author subscription importer", slogx.Error(err))
		return fmt.Errorf("failed to run author subscription importer: %w", err)
	}

	slog.InfoContext(ctx, "Author subscription import completed successfully")

	return nil
}
