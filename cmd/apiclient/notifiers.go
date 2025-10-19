package apiclient

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func createNotifiersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "notifiers",
		Short: "Manage notifiers",
	}

	cmd.AddCommand(createNotifiersListCmd())
	cmd.AddCommand(createNotifiersGetCmd())
	cmd.AddCommand(createNotifiersCreateCmd())
	cmd.AddCommand(createNotifiersUpdateCmd())
	cmd.AddCommand(createNotifiersDeleteCmd())

	return cmd
}

func createNotifiersListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all notifiers",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := NewClient(apiURL)
			ctx := context.Background()

			var notifiers []map[string]interface{}
			err := client.Get(ctx, "/notifiers", &notifiers)
			if err != nil {
				OutputError(err)
				return err
			}

			if format == "json" {
				return OutputJSON(notifiers)
			}

			// Table output
			headers := []string{"ID", "Name", "Type"}
			rows := [][]string{}
			for _, notifier := range notifiers {
				id := fmt.Sprintf("%.0f", notifier["id"].(float64))
				name := notifier["name"].(string)

				// Format type as "typename(id)"
				typeStr := "N/A"
				if typeData, ok := notifier["type"].(map[string]interface{}); ok {
					typeName := typeData["name"].(string)
					typeID := typeData["id"].(float64)
					typeStr = FormatRelation(typeName, uint(typeID))
				}

				rows = append(rows, []string{id, name, typeStr})
			}

			return OutputTable(headers, rows)
		},
	}
}

func createNotifiersGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a notifier by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := NewClient(apiURL)
			ctx := context.Background()

			id := args[0]
			var notifier map[string]interface{}
			err := client.Get(ctx, "/notifiers/"+id, &notifier)
			if err != nil {
				OutputError(err)
				return err
			}

			if format == "json" {
				return OutputJSON(notifier)
			}

			// Table output
			typeStr := "N/A"
			if typeData, ok := notifier["type"].(map[string]interface{}); ok {
				typeName := typeData["name"].(string)
				typeID := typeData["id"].(float64)
				typeStr = FormatRelation(typeName, uint(typeID))
			}

			headers := []string{"Field", "Value"}
			rows := [][]string{
				{"ID", fmt.Sprintf("%.0f", notifier["id"].(float64))},
				{"Name", notifier["name"].(string)},
				{"Type", typeStr},
				{"Webhook URL", fmt.Sprintf("%v", notifier["webhook_url"])},
			}

			return OutputTable(headers, rows)
		},
	}
}

func createNotifiersCreateCmd() *cobra.Command {
	var name, typeName, webhookURL string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new notifier",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := NewClient(apiURL)
			ctx := context.Background()

			body := map[string]interface{}{
				"name":        name,
				"type_name":   typeName,
				"webhook_url": webhookURL,
			}

			var notifier map[string]interface{}
			err := client.Post(ctx, "/notifiers", body, &notifier)
			if err != nil {
				OutputError(err)
				return err
			}

			if format == "json" {
				return OutputJSON(notifier)
			}

			fmt.Printf("Created notifier: %s (ID: %.0f)\n", notifier["name"].(string), notifier["id"].(float64))
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Notifier name (required)")
	cmd.Flags().StringVar(&typeName, "type", "", "Notification type name (required)")
	cmd.Flags().StringVar(&webhookURL, "webhook-url", "", "Webhook URL (required)")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("webhook-url")

	return cmd
}

func createNotifiersUpdateCmd() *cobra.Command {
	var name, typeName, webhookURL string

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a notifier",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := NewClient(apiURL)
			ctx := context.Background()

			id := args[0]
			body := map[string]interface{}{}

			if cmd.Flags().Changed("name") {
				body["name"] = name
			}
			if cmd.Flags().Changed("type") {
				body["type_name"] = typeName
			}
			if cmd.Flags().Changed("webhook-url") {
				body["webhook_url"] = webhookURL
			}

			var notifier map[string]interface{}
			err := client.Put(ctx, "/notifiers/"+id, body, &notifier)
			if err != nil {
				OutputError(err)
				return err
			}

			if format == "json" {
				return OutputJSON(notifier)
			}

			fmt.Printf("Updated notifier: %s (ID: %.0f)\n", notifier["name"].(string), notifier["id"].(float64))
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Notifier name")
	cmd.Flags().StringVar(&typeName, "type", "", "Notification type name")
	cmd.Flags().StringVar(&webhookURL, "webhook-url", "", "Webhook URL")

	return cmd
}

func createNotifiersDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "delete <id>",
		Aliases: []string{"rm"},
		Short:   "Delete a notifier",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := NewClient(apiURL)
			ctx := context.Background()

			id := args[0]
			err := client.Delete(ctx, "/notifiers/"+id)
			if err != nil {
				OutputError(err)
				return err
			}

			if format == "json" {
				return OutputJSON(map[string]string{"status": "deleted", "id": id})
			}

			idInt, _ := strconv.Atoi(id)
			fmt.Printf("Deleted notifier with ID: %d\n", idInt)
			return nil
		},
	}
}
