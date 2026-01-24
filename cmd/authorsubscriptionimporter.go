package cmd

import (
	"context"
	"errors"
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
		msg := "failed to connect to database"
		slog.ErrorContext(ctx, msg, slogx.Error(err))
		return errors.Join(errors.New(msg), err)
	}

	// Create qBittorrent client
	qbitClient, err := qbit.CreateClient()
	if err != nil {
		msg := "failed to create qBittorrent client"
		slog.ErrorContext(ctx, msg, slogx.Error(err))
		return errors.Join(errors.New(msg), err)
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
		msg := "failed to create audiobook importer system"
		slog.ErrorContext(ctx, msg, slogx.Error(err))
		return errors.Join(errors.New(msg), err)
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
		msg := "failed to run author subscription importer"
		slog.ErrorContext(ctx, msg, slogx.Error(err))
		return errors.Join(errors.New(msg), err)
	}

	slog.InfoContext(ctx, "Author subscription import completed successfully")

	return nil
}
