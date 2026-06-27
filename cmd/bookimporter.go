package cmd

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/cappuccinotm/slogx"
	"github.com/spf13/cobra"

	"github.com/bobbyrward/stronghold/internal/eventlog"
	"github.com/bobbyrward/stronghold/internal/importers/ebooks"
	"github.com/bobbyrward/stronghold/internal/models"
	"github.com/bobbyrward/stronghold/internal/qbit"
)

func createBookImportCmd() *cobra.Command {
	bookImportCmd := &cobra.Command{
		Use:  "book-import",
		RunE: runBookImport,
	}

	return bookImportCmd
}

func runBookImport(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	slog.InfoContext(ctx, "Starting book import command")

	db, err := models.ConnectAndMigrate(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to connect to database", slogx.Error(err))
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	eventlog.Cleanup(ctx, db, 90)

	qbitClient, err := qbit.CreateClient()
	if err != nil {
		slog.ErrorContext(ctx, "failed to create qBittorrent client", slogx.Error(err))
		return fmt.Errorf("failed to create qBittorrent client: %w", err)
	}

	bookImporterSystem := ebooks.NewBookImporterSystem(qbitClient, db)

	err = bookImporterSystem.Run(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Book import failed", slog.Any("err", err))
		return fmt.Errorf("failed to run book importer: %w", err)
	}

	slog.InfoContext(ctx, "Book import completed successfully")
	return nil
}
