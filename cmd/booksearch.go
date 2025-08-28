package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bobbyrward/stronghold/internal/booksearch"
)

func createBookSearchCmd() *cobra.Command {
	var (
		query  string
		format string
		limit  int
		hash   string
	)

	bookSearchCmd := &cobra.Command{
		Use:   "book-search",
		Short: "Search for books using external APIs",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runBookSearch(cmd, args, query, hash, format, limit)
		},
	}

	bookSearchCmd.Flags().StringVarP(&query, "query", "q", "", "Search query")
	bookSearchCmd.Flags().StringVarP(&format, "format", "f", "table", "Output format (table, json)")
	bookSearchCmd.Flags().IntVarP(&limit, "limit", "l", 10, "Maximum number of results")
	bookSearchCmd.Flags().StringVarP(&hash, "hash", "x", "", "Search by hash")

	return bookSearchCmd
}

func runBookSearch(cmd *cobra.Command, args []string, query, hash, format string, limit int) error {
	ctx := context.Background()

	if (query != "") == (hash != "") {
		return fmt.Errorf("one of query or hash must be specified")
	}

	searchService := booksearch.NewBookSearchService()

	params := booksearch.SearchParameters{MaxResults: limit}

	if hash != "" {
		params.Hash = hash
	} else {
		params.Query = query
	}

	results, err := searchService.Search(ctx,
		&params,
	)
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to search for books"))
	}

	err = searchService.DisplayResults(&params, results, format)
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to display results"))
	}

	return nil
}
