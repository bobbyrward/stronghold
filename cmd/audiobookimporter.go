package cmd

import (
	"context"
	"fmt"
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
		slog.ErrorContext(ctx, "failed to create qBittorrent client", slogx.Error(err))
		return fmt.Errorf("failed to create qBittorrent client: %w", err)
	}

	audibleApiClient := audible.NewAudibleApiClient()

	abookImporterSystem, err := audiobooks.NewAudiobookImporterSystem(
		qbitClient,
		config.Config.Importers,
		metadata.NewFFProbeMetadataProvider(),
		audibleApiClient,
	)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create audiobook importer system", slogx.Error(err))
		return fmt.Errorf("failed to create audiobook importer system: %w", err)
	}

	err = abookImporterSystem.Run(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to run audiobook importer system", slogx.Error(err))
		return fmt.Errorf("failed to run audiobook importer system: %w", err)
	}

	slog.InfoContext(ctx, "Book import completed successfully")

	return nil
}
