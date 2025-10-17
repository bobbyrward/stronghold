package apiclient

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func createFeedFilterSetEntriesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "feed-filter-set-entries",
		Short: "Manage feed filter set entries",
	}

	cmd.AddCommand(createFeedFilterSetEntriesListCmd())
	cmd.AddCommand(createFeedFilterSetEntriesGetCmd())
	cmd.AddCommand(createFeedFilterSetEntriesCreateCmd())
	cmd.AddCommand(createFeedFilterSetEntriesUpdateCmd())
	cmd.AddCommand(createFeedFilterSetEntriesDeleteCmd())

	return cmd
}

func createFeedFilterSetEntriesListCmd() *cobra.Command {
	var setID string

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all feed filter set entries",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := NewClient(apiURL)
			ctx := context.Background()

			path := "/feed-filter-set-entries"
			if setID != "" {
				path += "?feed_filter_set_id=" + setID
			}

			var entries []map[string]interface{}
			err := client.Get(ctx, path, &entries)
			if err != nil {
				OutputError(err)
				return err
			}

			if format == "json" {
				return OutputJSON(entries)
			}

			// Table output
			headers := []string{"ID", "Set", "Key", "Operator", "Value"}
			rows := [][]string{}
			for _, entry := range entries {
				id := fmt.Sprintf("%.0f", entry["id"].(float64))

				// Format set as "setname(id)"
				setStr := "N/A"
				if setData, ok := entry["set"].(map[string]interface{}); ok {
					setName := setData["name"].(string)
					setID := setData["id"].(float64)
					setStr = FormatRelation(setName, uint(setID))
				}

				// Format key as "keyname(id)"
				keyStr := "N/A"
				if keyData, ok := entry["key"].(map[string]interface{}); ok {
					keyName := keyData["name"].(string)
					keyID := keyData["id"].(float64)
					keyStr = FormatRelation(keyName, uint(keyID))
				}

				// Format operator as "opname(id)"
				operatorStr := "N/A"
				if operatorData, ok := entry["operator"].(map[string]interface{}); ok {
					operatorName := operatorData["name"].(string)
					operatorID := operatorData["id"].(float64)
					operatorStr = FormatRelation(operatorName, uint(operatorID))
				}

				value := fmt.Sprintf("%v", entry["value"])

				rows = append(rows, []string{id, setStr, keyStr, operatorStr, value})
			}

			return OutputTable(headers, rows)
		},
	}

	cmd.Flags().StringVar(&setID, "set-id", "", "Filter by set ID")

	return cmd
}

func createFeedFilterSetEntriesGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a feed filter set entry by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := NewClient(apiURL)
			ctx := context.Background()

			id := args[0]
			var entry map[string]interface{}
			err := client.Get(ctx, "/feed-filter-set-entries/"+id, &entry)
			if err != nil {
				OutputError(err)
				return err
			}

			if format == "json" {
				return OutputJSON(entry)
			}

			// Table output
			setStr := "N/A"
			if setData, ok := entry["set"].(map[string]interface{}); ok {
				setName := setData["name"].(string)
				setID := setData["id"].(float64)
				setStr = FormatRelation(setName, uint(setID))
			}

			keyStr := "N/A"
			if keyData, ok := entry["key"].(map[string]interface{}); ok {
				keyName := keyData["name"].(string)
				keyID := keyData["id"].(float64)
				keyStr = FormatRelation(keyName, uint(keyID))
			}

			operatorStr := "N/A"
			if operatorData, ok := entry["operator"].(map[string]interface{}); ok {
				operatorName := operatorData["name"].(string)
				operatorID := operatorData["id"].(float64)
				operatorStr = FormatRelation(operatorName, uint(operatorID))
			}

			headers := []string{"Field", "Value"}
			rows := [][]string{
				{"ID", fmt.Sprintf("%.0f", entry["id"].(float64))},
				{"Set", setStr},
				{"Key", keyStr},
				{"Operator", operatorStr},
				{"Value", fmt.Sprintf("%v", entry["value"])},
			}

			return OutputTable(headers, rows)
		},
	}
}

func createFeedFilterSetEntriesCreateCmd() *cobra.Command {
	var setID, keyName, operatorName, value string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new feed filter set entry",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := NewClient(apiURL)
			ctx := context.Background()

			// Convert setID to int
			setIDInt, err := strconv.Atoi(setID)
			if err != nil {
				OutputError(fmt.Errorf("invalid set ID: %s", setID))
				return err
			}

			body := map[string]interface{}{
				"feed_filter_set_id": setIDInt,
				"key_name":           keyName,
				"operator_name":      operatorName,
				"value":              value,
			}

			var entry map[string]interface{}
			err = client.Post(ctx, "/feed-filter-set-entries", body, &entry)
			if err != nil {
				OutputError(err)
				return err
			}

			if format == "json" {
				return OutputJSON(entry)
			}

			fmt.Printf("Created feed filter set entry with ID: %.0f\n", entry["id"].(float64))
			return nil
		},
	}

	cmd.Flags().StringVar(&setID, "set-id", "", "Filter set ID (required)")
	cmd.Flags().StringVar(&keyName, "key", "", "Filter key name (required)")
	cmd.Flags().StringVar(&operatorName, "operator", "", "Filter operator name (required)")
	cmd.Flags().StringVar(&value, "value", "", "Filter value (required)")
	_ = cmd.MarkFlagRequired("set-id")
	_ = cmd.MarkFlagRequired("key")
	_ = cmd.MarkFlagRequired("operator")
	_ = cmd.MarkFlagRequired("value")

	return cmd
}

func createFeedFilterSetEntriesUpdateCmd() *cobra.Command {
	var setID, keyName, operatorName, value string

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a feed filter set entry",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := NewClient(apiURL)
			ctx := context.Background()

			id := args[0]
			body := map[string]interface{}{}

			if cmd.Flags().Changed("set-id") {
				setIDInt, err := strconv.Atoi(setID)
				if err != nil {
					OutputError(fmt.Errorf("invalid set ID: %s", setID))
					return err
				}
				body["feed_filter_set_id"] = setIDInt
			}
			if cmd.Flags().Changed("key") {
				body["key_name"] = keyName
			}
			if cmd.Flags().Changed("operator") {
				body["operator_name"] = operatorName
			}
			if cmd.Flags().Changed("value") {
				body["value"] = value
			}

			var entry map[string]interface{}
			err := client.Put(ctx, "/feed-filter-set-entries/"+id, body, &entry)
			if err != nil {
				OutputError(err)
				return err
			}

			if format == "json" {
				return OutputJSON(entry)
			}

			fmt.Printf("Updated feed filter set entry with ID: %.0f\n", entry["id"].(float64))
			return nil
		},
	}

	cmd.Flags().StringVar(&setID, "set-id", "", "Filter set ID")
	cmd.Flags().StringVar(&keyName, "key", "", "Filter key name")
	cmd.Flags().StringVar(&operatorName, "operator", "", "Filter operator name")
	cmd.Flags().StringVar(&value, "value", "", "Filter value")

	return cmd
}

func createFeedFilterSetEntriesDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "delete <id>",
		Aliases: []string{"rm"},
		Short:   "Delete a feed filter set entry",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := NewClient(apiURL)
			ctx := context.Background()

			id := args[0]
			err := client.Delete(ctx, "/feed-filter-set-entries/"+id)
			if err != nil {
				OutputError(err)
				return err
			}

			if format == "json" {
				return OutputJSON(map[string]string{"status": "deleted", "id": id})
			}

			idInt, _ := strconv.Atoi(id)
			fmt.Printf("Deleted feed filter set entry with ID: %d\n", idInt)
			return nil
		},
	}
}
