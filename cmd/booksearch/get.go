package booksearch

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/bobbyrward/stronghold/internal/booksearch"
	"github.com/bobbyrward/stronghold/internal/models"
	"github.com/spf13/cobra"
)

func createGetByIDCommand() *cobra.Command {
	getByIDCmd := &cobra.Command{
		Use:   "get-by-id",
		Short: "Lookup a book by its ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGetByIDCommand(cmd, args[0])
		},
	}

	return getByIDCmd
}

func createGetByHashCommand() *cobra.Command {
	getByHashCmd := &cobra.Command{
		Use:   "get-by-hash",
		Short: "Lookup a book by its hash",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGetByHashCommand(cmd, args[0])
		},
	}

	return getByHashCmd
}

func runGetByIDCommand(cmd *cobra.Command, idString string) error {
	ctx := context.Background()

	id, err := strconv.Atoi(idString)
	if err != nil {
		return errors.Join(err, fmt.Errorf("invalid ID: %s", idString))
	}

	// Connect to database
	db, err := models.ConnectDB()
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to connect to database"))
	}

	searchService := booksearch.NewBookSearchService()
	params := booksearch.SearchParameters{ID: &id}

	results, err := searchService.Search(ctx, db, &params)
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to search for books"))
	}

	err = displaySingleResult(&params, results)
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to display results"))
	}

	return nil
}

func runGetByHashCommand(cmd *cobra.Command, hash string) error {
	ctx := context.Background()

	// Connect to database
	db, err := models.ConnectDB()
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to connect to database"))
	}

	searchService := booksearch.NewBookSearchService()
	params := booksearch.SearchParameters{Hash: hash}

	results, err := searchService.Search(ctx, db, &params)
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to search for books"))
	}

	err = displaySingleResult(&params, results)
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to display results"))
	}

	return nil
}
