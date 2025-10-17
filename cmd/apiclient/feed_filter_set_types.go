package apiclient

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func createFeedFilterSetTypesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "feed-filter-set-types",
		Short: "View feed filter set types (read-only)",
	}

	cmd.AddCommand(createFeedFilterSetTypesListCmd())
	cmd.AddCommand(createFeedFilterSetTypesGetCmd())

	return cmd
}

func createFeedFilterSetTypesListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all feed filter set types",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := NewClient(apiURL)
			ctx := context.Background()

			var types []map[string]interface{}
			err := client.Get(ctx, "/feed-filter-set-types", &types)
			if err != nil {
				OutputError(err)
				return err
			}

			if format == "json" {
				return OutputJSON(types)
			}

			// Table output
			headers := []string{"ID", "Name", "Description"}
			rows := [][]string{}
			for _, setType := range types {
				id := fmt.Sprintf("%.0f", setType["id"].(float64))
				name := setType["name"].(string)
				description := ""
				if desc, ok := setType["description"].(string); ok {
					description = desc
				}
				rows = append(rows, []string{id, name, description})
			}

			return OutputTable(headers, rows)
		},
	}
}

func createFeedFilterSetTypesGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a feed filter set type by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := NewClient(apiURL)
			ctx := context.Background()

			id := args[0]
			var setType map[string]interface{}
			err := client.Get(ctx, "/feed-filter-set-types/"+id, &setType)
			if err != nil {
				OutputError(err)
				return err
			}

			if format == "json" {
				return OutputJSON(setType)
			}

			// Table output
			description := ""
			if desc, ok := setType["description"].(string); ok {
				description = desc
			}

			headers := []string{"Field", "Value"}
			rows := [][]string{
				{"ID", fmt.Sprintf("%.0f", setType["id"].(float64))},
				{"Name", setType["name"].(string)},
				{"Description", description},
			}

			return OutputTable(headers, rows)
		},
	}
}
