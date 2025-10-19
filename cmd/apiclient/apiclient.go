package apiclient

import (
	"github.com/spf13/cobra"
)

var (
	apiURL string
	format string
)

// CreateAPIClientCmd creates the apiclient command with all subcommands
func CreateAPIClientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apiclient",
		Short: "Interact with the Stronghold API",
		Long:  "Command-line client for managing Stronghold resources via the REST API",
	}

	// Global flags
	cmd.PersistentFlags().StringVarP(&apiURL, "api-url", "u", "https://stronghold.home.ohnozombi.es", "API server URL")
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
