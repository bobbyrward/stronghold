package apiclient

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func createFilterKeysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "filter-keys",
		Short: "View filter keys (read-only)",
	}

	cmd.AddCommand(createFilterKeysListCmd())
	cmd.AddCommand(createFilterKeysGetCmd())

	return cmd
}

func createFilterKeysListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all filter keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := NewClient(apiURL)
			ctx := context.Background()

			var keys []map[string]interface{}
			err := client.Get(ctx, "/filter-keys", &keys)
			if err != nil {
				OutputError(err)
				return err
			}

			if format == "json" {
				return OutputJSON(keys)
			}

			// Table output
			headers := []string{"ID", "Name", "Description"}
			rows := [][]string{}
			for _, key := range keys {
				id := fmt.Sprintf("%.0f", key["id"].(float64))
				name := key["name"].(string)
				description := ""
				if desc, ok := key["description"].(string); ok {
					description = desc
				}
				rows = append(rows, []string{id, name, description})
			}

			return OutputTable(headers, rows)
		},
	}
}

func createFilterKeysGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a filter key by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := NewClient(apiURL)
			ctx := context.Background()

			id := args[0]
			var key map[string]interface{}
			err := client.Get(ctx, "/filter-keys/"+id, &key)
			if err != nil {
				OutputError(err)
				return err
			}

			if format == "json" {
				return OutputJSON(key)
			}

			// Table output
			description := ""
			if desc, ok := key["description"].(string); ok {
				description = desc
			}

			headers := []string{"Field", "Value"}
			rows := [][]string{
				{"ID", fmt.Sprintf("%.0f", key["id"].(float64))},
				{"Name", key["name"].(string)},
				{"Description", description},
			}

			return OutputTable(headers, rows)
		},
	}
}
