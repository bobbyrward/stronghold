package cmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/bobbyrward/stronghold/internal/models"
	"github.com/spf13/cobra"
)

func createDoctorCmd() *cobra.Command {
	doctorCmd := &cobra.Command{
		Use:   "doctor",
		Short: "System diagnostics and maintenance commands",
	}

	doctorCmd.AddCommand(createDoctorMigrateCmd())
	doctorCmd.AddCommand(createDoctorInitBookSearchCmd())

	return doctorCmd
}

func createDoctorMigrateCmd() *cobra.Command {
	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
		RunE:  runDoctorMigrateCmd,
	}

	return migrateCmd
}

func runDoctorMigrateCmd(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	slog.InfoContext(ctx, "Starting database migration")

	_, err := models.ConnectAndMigrate(ctx)
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, "Database migration completed successfully")
	return nil
}

func createDoctorInitBookSearchCmd() *cobra.Command {
	var apiKey string

	initBookSearchCmd := &cobra.Command{
		Use:   "init-book-search",
		Short: "Initialize book search credentials in the database",
		Long: `Initialize book search credentials by providing an API key (token cookie value).
This creates or updates the BookSearchCredential record in the database.
You can obtain the token cookie value from your browser after logging into the book search site.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDoctorInitBookSearchCmd(cmd, args, apiKey)
		},
	}

	initBookSearchCmd.Flags().StringVar(&apiKey, "api-key", "", "Book search API key (token cookie value) - required")
	_ = initBookSearchCmd.MarkFlagRequired("api-key")

	return initBookSearchCmd
}

func runDoctorInitBookSearchCmd(cmd *cobra.Command, args []string, apiKey string) error {
	ctx := context.Background()

	slog.InfoContext(ctx, "Initializing book search credentials")

	if apiKey == "" {
		return fmt.Errorf("api-key is required")
	}

	db, err := models.ConnectDB()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to connect to database", slog.Any("err", err))
		return errors.Join(err, fmt.Errorf("failed to connect to database"))
	}

	// Initialize with empty IP and ASN - these will be populated on first token refresh
	err = models.UpsertBookSearchCredential(db, apiKey, "", "")
	if err != nil {
		slog.ErrorContext(ctx, "Failed to initialize book search credentials", slog.Any("err", err))
		return errors.Join(err, fmt.Errorf("failed to initialize book search credentials"))
	}

	slog.InfoContext(ctx, "Book search credentials initialized successfully")
	fmt.Println("Book search credentials initialized successfully")
	fmt.Println("You can now run 'stronghold refresh-token' to update IP address and ASN information")

	return nil
}
