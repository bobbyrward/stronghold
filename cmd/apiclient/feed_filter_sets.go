package apiclient

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func createFeedFilterSetsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "feed-filter-sets",
		Short: "Manage feed filter sets",
	}

	cmd.AddCommand(createFeedFilterSetsListCmd())
	cmd.AddCommand(createFeedFilterSetsGetCmd())
	cmd.AddCommand(createFeedFilterSetsCreateCmd())
	cmd.AddCommand(createFeedFilterSetsUpdateCmd())
	cmd.AddCommand(createFeedFilterSetsDeleteCmd())

	return cmd
}

func createFeedFilterSetsListCmd() *cobra.Command {
	var feedID string

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all feed filter sets",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := NewClient(apiURL)
			ctx := context.Background()

			path := "/feed-filter-sets"
			if feedID != "" {
				path += "?feed_filter_id=" + feedID
			}

			var sets []map[string]interface{}
			err := client.Get(ctx, path, &sets)
			if err != nil {
				OutputError(err)
				return err
			}

			if format == "json" {
				return OutputJSON(sets)
			}

			// Table output
			headers := []string{"ID", "Name", "Feed", "Type", "Category", "Notifier"}
			rows := [][]string{}
			for _, set := range sets {
				id := fmt.Sprintf("%.0f", set["id"].(float64))
				name := set["name"].(string)

				// Format feed as "feedname(id)"
				feedStr := "N/A"
				if feedData, ok := set["feed"].(map[string]interface{}); ok {
					feedName := feedData["name"].(string)
					feedID := feedData["id"].(float64)
					feedStr = FormatRelation(feedName, uint(feedID))
				}

				// Format type as "typename(id)"
				typeStr := "N/A"
				if typeData, ok := set["type"].(map[string]interface{}); ok {
					typeName := typeData["name"].(string)
					typeID := typeData["id"].(float64)
					typeStr = FormatRelation(typeName, uint(typeID))
				}

				// Format category as "categoryname(id)"
				categoryStr := "N/A"
				if categoryData, ok := set["category"].(map[string]interface{}); ok {
					categoryName := categoryData["name"].(string)
					categoryID := categoryData["id"].(float64)
					categoryStr = FormatRelation(categoryName, uint(categoryID))
				}

				// Format notifier as "notifiername(id)"
				notifierStr := "N/A"
				if notifierData, ok := set["notifier"].(map[string]interface{}); ok {
					notifierName := notifierData["name"].(string)
					notifierID := notifierData["id"].(float64)
					notifierStr = FormatRelation(notifierName, uint(notifierID))
				}

				rows = append(rows, []string{id, name, feedStr, typeStr, categoryStr, notifierStr})
			}

			return OutputTable(headers, rows)
		},
	}

	cmd.Flags().StringVar(&feedID, "feed-id", "", "Filter by feed ID")

	return cmd
}

func createFeedFilterSetsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a feed filter set by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := NewClient(apiURL)
			ctx := context.Background()

			id := args[0]
			var set map[string]interface{}
			err := client.Get(ctx, "/feed-filter-sets/"+id, &set)
			if err != nil {
				OutputError(err)
				return err
			}

			if format == "json" {
				return OutputJSON(set)
			}

			// Table output
			feedStr := "N/A"
			if feedData, ok := set["feed"].(map[string]interface{}); ok {
				feedName := feedData["name"].(string)
				feedID := feedData["id"].(float64)
				feedStr = FormatRelation(feedName, uint(feedID))
			}

			typeStr := "N/A"
			if typeData, ok := set["type"].(map[string]interface{}); ok {
				typeName := typeData["name"].(string)
				typeID := typeData["id"].(float64)
				typeStr = FormatRelation(typeName, uint(typeID))
			}

			categoryStr := "N/A"
			if categoryData, ok := set["category"].(map[string]interface{}); ok {
				categoryName := categoryData["name"].(string)
				categoryID := categoryData["id"].(float64)
				categoryStr = FormatRelation(categoryName, uint(categoryID))
			}

			notifierStr := "N/A"
			if notifierData, ok := set["notifier"].(map[string]interface{}); ok {
				notifierName := notifierData["name"].(string)
				notifierID := notifierData["id"].(float64)
				notifierStr = FormatRelation(notifierName, uint(notifierID))
			}

			headers := []string{"Field", "Value"}
			rows := [][]string{
				{"ID", fmt.Sprintf("%.0f", set["id"].(float64))},
				{"Name", set["name"].(string)},
				{"Feed", feedStr},
				{"Type", typeStr},
				{"Category", categoryStr},
				{"Notifier", notifierStr},
			}

			return OutputTable(headers, rows)
		},
	}
}

func createFeedFilterSetsCreateCmd() *cobra.Command {
	var name, feedID, typeName, categoryName, notifierID string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new feed filter set",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := NewClient(apiURL)
			ctx := context.Background()

			// Convert IDs to int
			feedIDInt, err := strconv.Atoi(feedID)
			if err != nil {
				OutputError(fmt.Errorf("invalid feed ID: %s", feedID))
				return err
			}

			notifierIDInt, err := strconv.Atoi(notifierID)
			if err != nil {
				OutputError(fmt.Errorf("invalid notifier ID: %s", notifierID))
				return err
			}

			body := map[string]interface{}{
				"name":          name,
				"feed_id":       feedIDInt,
				"type_name":     typeName,
				"category_name": categoryName,
				"notifier_id":   notifierIDInt,
			}

			var set map[string]interface{}
			err = client.Post(ctx, "/feed-filter-sets", body, &set)
			if err != nil {
				OutputError(err)
				return err
			}

			if format == "json" {
				return OutputJSON(set)
			}

			fmt.Printf("Created feed filter set: %s (ID: %.0f)\n", set["name"].(string), set["id"].(float64))
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Filter set name (required)")
	cmd.Flags().StringVar(&feedID, "feed-id", "", "Feed ID (required)")
	cmd.Flags().StringVar(&typeName, "type", "", "Filter set type name (required)")
	cmd.Flags().StringVar(&categoryName, "category", "", "Torrent category name (required)")
	cmd.Flags().StringVar(&notifierID, "notifier-id", "", "Notifier ID (required)")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("feed-id")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("category")
	_ = cmd.MarkFlagRequired("notifier-id")

	return cmd
}

func createFeedFilterSetsUpdateCmd() *cobra.Command {
	var name, feedID, typeName, categoryName, notifierID string

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a feed filter set",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := NewClient(apiURL)
			ctx := context.Background()

			id := args[0]
			body := map[string]interface{}{}

			if cmd.Flags().Changed("name") {
				body["name"] = name
			}
			if cmd.Flags().Changed("feed-id") {
				feedIDInt, err := strconv.Atoi(feedID)
				if err != nil {
					OutputError(fmt.Errorf("invalid feed ID: %s", feedID))
					return err
				}
				body["feed_id"] = feedIDInt
			}
			if cmd.Flags().Changed("type") {
				body["type_name"] = typeName
			}
			if cmd.Flags().Changed("category") {
				body["category_name"] = categoryName
			}
			if cmd.Flags().Changed("notifier-id") {
				notifierIDInt, err := strconv.Atoi(notifierID)
				if err != nil {
					OutputError(fmt.Errorf("invalid notifier ID: %s", notifierID))
					return err
				}
				body["notifier_id"] = notifierIDInt
			}

			var set map[string]interface{}
			err := client.Put(ctx, "/feed-filter-sets/"+id, body, &set)
			if err != nil {
				OutputError(err)
				return err
			}

			if format == "json" {
				return OutputJSON(set)
			}

			fmt.Printf("Updated feed filter set: %s (ID: %.0f)\n", set["name"].(string), set["id"].(float64))
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Filter set name")
	cmd.Flags().StringVar(&feedID, "feed-id", "", "Feed ID")
	cmd.Flags().StringVar(&typeName, "type", "", "Filter set type name")
	cmd.Flags().StringVar(&categoryName, "category", "", "Torrent category name")
	cmd.Flags().StringVar(&notifierID, "notifier-id", "", "Notifier ID")

	return cmd
}

func createFeedFilterSetsDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "delete <id>",
		Aliases: []string{"rm"},
		Short:   "Delete a feed filter set",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := NewClient(apiURL)
			ctx := context.Background()

			id := args[0]
			err := client.Delete(ctx, "/feed-filter-sets/"+id)
			if err != nil {
				OutputError(err)
				return err
			}

			if format == "json" {
				return OutputJSON(map[string]string{"status": "deleted", "id": id})
			}

			idInt, _ := strconv.Atoi(id)
			fmt.Printf("Deleted feed filter set with ID: %d\n", idInt)
			return nil
		},
	}
}
