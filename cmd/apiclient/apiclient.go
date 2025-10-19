package apiclient

import (
	"github.com/bobbyrward/stronghold/internal/config"
	"github.com/spf13/cobra"
)

var (
	apiURL string
	format string
)

const defaultAPIURL = "http://localhost:8000"

// CreateAPIClientCmd creates the apiclient command with all subcommands
func CreateAPIClientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apiclient",
		Short: "Interact with the Stronghold API",
		Long:  "Command-line client for managing Stronghold resources via the REST API",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// If the user didn't explicitly set the api-url flag, try to load from config
			if !cmd.Flags().Changed("api-url") {
				// Try to load config (non-fatal if it fails)
				if config.Config.APIClient.URL != "" {
					apiURL = config.Config.APIClient.URL
				}
			}
			return nil
		},
	}

	// Global flags
	cmd.PersistentFlags().StringVarP(&apiURL, "api-url", "u", defaultAPIURL, "API server URL")
	cmd.PersistentFlags().StringVarP(&format, "format", "f", "table", "Output format: table or json")

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
