package apiclient

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func createFilterOperatorsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "filter-operators",
		Short: "View filter operators (read-only)",
	}

	cmd.AddCommand(createFilterOperatorsListCmd())
	cmd.AddCommand(createFilterOperatorsGetCmd())

	return cmd
}

func createFilterOperatorsListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all filter operators",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := NewClient(apiURL)
			ctx := context.Background()

			var operators []map[string]interface{}
			err := client.Get(ctx, "/filter-operators", &operators)
			if err != nil {
				OutputError(err)
				return err
			}

			if format == "json" {
				return OutputJSON(operators)
			}

			// Table output
			headers := []string{"ID", "Name", "Description"}
			rows := [][]string{}
			for _, operator := range operators {
				id := fmt.Sprintf("%.0f", operator["id"].(float64))
				name := operator["name"].(string)
				description := ""
				if desc, ok := operator["description"].(string); ok {
					description = desc
				}
				rows = append(rows, []string{id, name, description})
			}

			return OutputTable(headers, rows)
		},
	}
}

func createFilterOperatorsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a filter operator by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := NewClient(apiURL)
			ctx := context.Background()

			id := args[0]
			var operator map[string]interface{}
			err := client.Get(ctx, "/filter-operators/"+id, &operator)
			if err != nil {
				OutputError(err)
				return err
			}

			if format == "json" {
				return OutputJSON(operator)
			}

			// Table output
			description := ""
			if desc, ok := operator["description"].(string); ok {
				description = desc
			}

			headers := []string{"Field", "Value"}
			rows := [][]string{
				{"ID", fmt.Sprintf("%.0f", operator["id"].(float64))},
				{"Name", operator["name"].(string)},
				{"Description", description},
			}

			return OutputTable(headers, rows)
		},
	}
}
