package cmd

import (
	"context"
	"errors"
	"log/slog"

	"github.com/bobbyrward/stronghold/internal/config"
	"github.com/bobbyrward/stronghold/internal/importers/audiobooks"
	"github.com/bobbyrward/stronghold/internal/importers/audiobooks/audible"
	"github.com/bobbyrward/stronghold/internal/importers/audiobooks/metadata"
	"github.com/bobbyrward/stronghold/internal/qbit"
	"github.com/cappuccinotm/slogx"
	"github.com/spf13/cobra"
)

func createAudiobookImporterCmd() *cobra.Command {
	bookImportCmd := &cobra.Command{
		Use:  "audiobook-importer",
		RunE: runAudiobookImporter,
	}

	return bookImportCmd
}

func runAudiobookImporter(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	slog.InfoContext(ctx, "Starting book import command")

	qbitClient, err := qbit.CreateClient()
	if err != nil {
		msg := "failed to create qBittorrent client"
		slog.ErrorContext(ctx, msg, slogx.Error(err))
		return errors.Join(errors.New(msg), err)
	}

	audibleApiClient := audible.NewAudibleApiClient()

	abookImporterSystem, err := audiobooks.NewAudiobookImporterSystem(
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

	err = abookImporterSystem.Run(ctx)
	if err != nil {
		msg := "failed to run audiobook importer system"
		slog.ErrorContext(ctx, msg, slogx.Error(err))
		return errors.Join(errors.New(msg), err)
	}

	slog.InfoContext(ctx, "Book import completed successfully")

	return nil
}
