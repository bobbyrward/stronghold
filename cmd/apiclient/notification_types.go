package apiclient

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func createNotificationTypesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "notification-types",
		Short: "View notification types (read-only)",
	}

	cmd.AddCommand(createNotificationTypesListCmd())
	cmd.AddCommand(createNotificationTypesGetCmd())

	return cmd
}

func createNotificationTypesListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all notification types",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := NewClient(apiURL)
			ctx := context.Background()

			var types []map[string]interface{}
			err := client.Get(ctx, "/notification-types", &types)
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
			for _, notifType := range types {
				id := fmt.Sprintf("%.0f", notifType["id"].(float64))
				name := notifType["name"].(string)
				description := ""
				if desc, ok := notifType["description"].(string); ok {
					description = desc
				}
				rows = append(rows, []string{id, name, description})
			}

			return OutputTable(headers, rows)
		},
	}
}

func createNotificationTypesGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a notification type by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := NewClient(apiURL)
			ctx := context.Background()

			id := args[0]
			var notifType map[string]interface{}
			err := client.Get(ctx, "/notification-types/"+id, &notifType)
			if err != nil {
				OutputError(err)
				return err
			}

			if format == "json" {
				return OutputJSON(notifType)
			}

			// Table output
			description := ""
			if desc, ok := notifType["description"].(string); ok {
				description = desc
			}

			headers := []string{"Field", "Value"}
			rows := [][]string{
				{"ID", fmt.Sprintf("%.0f", notifType["id"].(float64))},
				{"Name", notifType["name"].(string)},
				{"Description", description},
			}

			return OutputTable(headers, rows)
		},
	}
}
