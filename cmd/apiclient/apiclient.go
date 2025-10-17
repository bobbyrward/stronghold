package apiclient

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var (
	apiURL   string
	format   string
	logLevel string
)

// CreateAPIClientCmd creates the apiclient command with all subcommands
func CreateAPIClientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apiclient",
		Short: "Interact with the Stronghold API",
		Long:  "Command-line client for managing Stronghold resources via the REST API",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Set up logging based on log-level flag
			setupLogging(logLevel)
		},
	}

	// Global flags
	cmd.PersistentFlags().StringVar(&apiURL, "api-url", "http://localhost:8000", "API server URL")
	cmd.PersistentFlags().StringVar(&format, "format", "table", "Output format: table or json")
	cmd.PersistentFlags().StringVar(&logLevel, "log-level", "none", "Log level: debug, info, warn, error, none")

	// Add subcommands for CRUD resources
	cmd.AddCommand(createFeedsCmd())
	cmd.AddCommand(createNotifiersCmd())
	cmd.AddCommand(createFeedFiltersCmd())
	cmd.AddCommand(createFeedFilterSetsCmd())
	cmd.AddCommand(createFeedFilterSetEntriesCmd())

	// Add subcommands for read-only reference data
	cmd.AddCommand(createFilterKeysCmd())
	cmd.AddCommand(createFilterOperatorsCmd())
	cmd.AddCommand(createNotificationTypesCmd())
	cmd.AddCommand(createFeedFilterSetTypesCmd())
	cmd.AddCommand(createTorrentCategoriesCmd())

	return cmd
}

func setupLogging(level string) {
	var slogLevel slog.Level

	switch level {
	case "debug":
		slogLevel = slog.LevelDebug
	case "info":
		slogLevel = slog.LevelInfo
	case "warn":
		slogLevel = slog.LevelWarn
	case "error":
		slogLevel = slog.LevelError
	case "none":
		// Set to a very high level to effectively disable logging
		slogLevel = slog.Level(1000)
	default:
		slogLevel = slog.Level(1000) // Default to no logging
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slogLevel,
	}))
	slog.SetDefault(logger)
}
