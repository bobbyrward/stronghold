package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bobbyrward/stronghold/internal/booksearch"
	"github.com/bobbyrward/stronghold/internal/models"
)

func createRefreshTokenCmd() *cobra.Command {
	refreshTokenCmd := &cobra.Command{
		Use:   "refresh-token",
		Short: "Refresh the book search API token",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRefreshToken(cmd, args)
		},
	}

	return refreshTokenCmd
}

func runRefreshToken(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Connect to database
	db, err := models.ConnectDB()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	searchService := booksearch.NewBookSearchService()

	err = searchService.RefreshToken(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to refresh token: %w", err)
	}

	fmt.Println("Token refresh completed successfully")

	return nil
}
