package booksearch

import (
	"context"
	"errors"
	"fmt"

	"github.com/bobbyrward/stronghold/internal/booksearch"
	"github.com/bobbyrward/stronghold/internal/models"
	"github.com/spf13/cobra"
)

func createSearchCommand() *cobra.Command {
	var (
		format string
		limit  int
	)

	searchCmd := &cobra.Command{
		Use:   "search",
		Short: "Search for books using external APIs",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := args[0]

			return runSearchCommand(cmd, query, format, limit)
		},
	}

	searchCmd.Flags().StringVarP(&format, "format", "f", "table", "Output format (table, json)")
	searchCmd.Flags().IntVar(&limit, "limit", 10, "Maximum number of results")

	return searchCmd
}

func runSearchCommand(cmd *cobra.Command, query, format string, limit int) error {
	ctx := context.Background()

	// Connect to database
	db, err := models.ConnectDB()
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to connect to database"))
	}

	searchService := booksearch.NewBookSearchService()

	params := booksearch.SearchParameters{MaxResults: limit, Query: query}

	results, err := searchService.Search(ctx, db, &params)
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to search for books"))
	}

	err = displaySearchResults(&params, results, format)
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to display results"))
	}

	return nil
}
