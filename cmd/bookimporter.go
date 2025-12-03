package cmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/cappuccinotm/slogx"
	"github.com/spf13/cobra"

	"github.com/bobbyrward/stronghold/internal/importers/ebooks"
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

	qbitClient, err := qbit.CreateClient()
	if err != nil {
		slog.ErrorContext(ctx, "failed to create qBittorrent client", slogx.Error(err))
		return errors.Join(err, fmt.Errorf("failed to create qBittorrent client"))
	}

	bookImporterSystem := ebooks.NewBookImporterSystem(qbitClient)

	err = bookImporterSystem.Run(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Book import failed", slog.Any("err", err))
		return errors.Join(err, fmt.Errorf("failed to run book importer"))
	}

	slog.InfoContext(ctx, "Book import completed successfully")
	return nil
}
