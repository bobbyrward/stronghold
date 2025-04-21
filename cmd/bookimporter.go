package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bobbyrward/stronghold/internal/bookimporter"
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

	bookImporterSystem := bookimporter.NewBookImporterSystem()

	err := bookImporterSystem.Run(ctx)
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to run book importer"))
	}

	return nil
}
