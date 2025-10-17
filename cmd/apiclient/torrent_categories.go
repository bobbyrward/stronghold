package apiclient

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func createTorrentCategoriesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "torrent-categories",
		Short: "View torrent categories (read-only)",
	}

	cmd.AddCommand(createTorrentCategoriesListCmd())
	cmd.AddCommand(createTorrentCategoriesGetCmd())

	return cmd
}

func createTorrentCategoriesListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all torrent categories",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := NewClient(apiURL)
			ctx := context.Background()

			var categories []map[string]interface{}
			err := client.Get(ctx, "/torrent-categories", &categories)
			if err != nil {
				OutputError(err)
				return err
			}

			if format == "json" {
				return OutputJSON(categories)
			}

			// Table output
			headers := []string{"ID", "Name", "Description"}
			rows := [][]string{}
			for _, category := range categories {
				id := fmt.Sprintf("%.0f", category["id"].(float64))
				name := category["name"].(string)
				description := ""
				if desc, ok := category["description"].(string); ok {
					description = desc
				}
				rows = append(rows, []string{id, name, description})
			}

			return OutputTable(headers, rows)
		},
	}
}

func createTorrentCategoriesGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a torrent category by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := NewClient(apiURL)
			ctx := context.Background()

			id := args[0]
			var category map[string]interface{}
			err := client.Get(ctx, "/torrent-categories/"+id, &category)
			if err != nil {
				OutputError(err)
				return err
			}

			if format == "json" {
				return OutputJSON(category)
			}

			// Table output
			description := ""
			if desc, ok := category["description"].(string); ok {
				description = desc
			}

			headers := []string{"Field", "Value"}
			rows := [][]string{
				{"ID", fmt.Sprintf("%.0f", category["id"].(float64))},
				{"Name", category["name"].(string)},
				{"Description", description},
			}

			return OutputTable(headers, rows)
		},
	}
}
