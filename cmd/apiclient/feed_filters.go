package apiclient

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func createFeedFiltersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "feed-filters",
		Short: "Manage feed filters",
	}

	cmd.AddCommand(createFeedFiltersListCmd())
	cmd.AddCommand(createFeedFiltersGetCmd())
	cmd.AddCommand(createFeedFiltersCreateCmd())
	cmd.AddCommand(createFeedFiltersUpdateCmd())
	cmd.AddCommand(createFeedFiltersDeleteCmd())

	return cmd
}

func createFeedFiltersListCmd() *cobra.Command {
	var feedID string

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all feed filters",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := NewClient(apiURL)
			ctx := context.Background()

			path := "/feed-filters"
			if feedID != "" {
				path += "?feed_id=" + feedID
			}

			var filters []map[string]interface{}
			err := client.Get(ctx, path, &filters)
			if err != nil {
				OutputError(err)
				return err
			}

			if format == "json" {
				return OutputJSON(filters)
			}

			// Table output
			headers := []string{"ID", "Name", "Feed", "Category", "Notifier"}
			rows := [][]string{}
			for _, filter := range filters {
				id := fmt.Sprintf("%.0f", filter["id"].(float64))
				name := fmt.Sprintf("%v", filter["name"])

				// Format feed as "feedname(id)"
				feedStr := "N/A"
				if feedName, ok := filter["feed_name"].(string); ok {
					feedID := uint(filter["feed_id"].(float64))
					feedStr = FormatRelation(feedName, feedID)
				}

				// Format category as "categoryname(id)"
				categoryStr := fmt.Sprintf("%v", filter["category"])

				// Format notifier as "notifiername(id)"
				notifierStr := fmt.Sprintf("%v", filter["notifier"])

				rows = append(rows, []string{id, name, feedStr, categoryStr, notifierStr})
			}

			return OutputTable(headers, rows)
		},
	}

	cmd.Flags().StringVar(&feedID, "feed-id", "", "Filter by feed ID")

	return cmd
}

func createFeedFiltersGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a feed filter by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := NewClient(apiURL)
			ctx := context.Background()

			id := args[0]
			var filter map[string]interface{}
			err := client.Get(ctx, "/feed-filters/"+id, &filter)
			if err != nil {
				OutputError(err)
				return err
			}

			if format == "json" {
				return OutputJSON(filter)
			}

			// Table output
			feedStr := "N/A"
			if feedName, ok := filter["feed_name"].(string); ok {
				feedID := uint(filter["feed_id"].(float64))
				feedStr = FormatRelation(feedName, feedID)
			}

			headers := []string{"Field", "Value"}
			rows := [][]string{
				{"ID", fmt.Sprintf("%.0f", filter["id"].(float64))},
				{"Name", fmt.Sprintf("%v", filter["name"])},
				{"Feed", feedStr},
				{"Category", fmt.Sprintf("%v", filter["category"])},
				{"Notifier", fmt.Sprintf("%v", filter["notifier"])},
			}

			return OutputTable(headers, rows)
		},
	}
}

func createFeedFiltersCreateCmd() *cobra.Command {
	var name, feedName, categoryName, notifierName string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new feed filter",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := NewClient(apiURL)
			ctx := context.Background()

			body := map[string]interface{}{
				"name":          name,
				"feed_name":     feedName,
				"category_name": categoryName,
				"notifier_name": notifierName,
			}

			var filter map[string]interface{}
			err := client.Post(ctx, "/feed-filters", body, &filter)
			if err != nil {
				OutputError(err)
				return err
			}

			if format == "json" {
				return OutputJSON(filter)
			}

			fmt.Printf("Created feed filter with ID: %.0f\n", filter["id"].(float64))
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Filter name (required)")
	cmd.Flags().StringVar(&feedName, "feed-name", "", "Feed name (required)")
	cmd.Flags().StringVar(&categoryName, "category-name", "", "Category name (required)")
	cmd.Flags().StringVar(&notifierName, "notifier-name", "", "Notifier name (required)")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("feed-name")
	_ = cmd.MarkFlagRequired("category-name")
	_ = cmd.MarkFlagRequired("notifier-name")

	return cmd
}

func createFeedFiltersUpdateCmd() *cobra.Command {
	var name, feedName, categoryName, notifierName string

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a feed filter",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := NewClient(apiURL)
			ctx := context.Background()

			id := args[0]
			body := map[string]interface{}{}

			if cmd.Flags().Changed("name") {
				body["name"] = name
			}
			if cmd.Flags().Changed("feed-name") {
				body["feed_name"] = feedName
			}
			if cmd.Flags().Changed("category-name") {
				body["category_name"] = categoryName
			}
			if cmd.Flags().Changed("notifier-name") {
				body["notifier_name"] = notifierName
			}

			var filter map[string]interface{}
			err := client.Put(ctx, "/feed-filters/"+id, body, &filter)
			if err != nil {
				OutputError(err)
				return err
			}

			if format == "json" {
				return OutputJSON(filter)
			}

			fmt.Printf("Updated feed filter with ID: %.0f\n", filter["id"].(float64))
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Filter name")
	cmd.Flags().StringVar(&feedName, "feed-name", "", "Feed name")
	cmd.Flags().StringVar(&categoryName, "category-name", "", "Category name")
	cmd.Flags().StringVar(&notifierName, "notifier-name", "", "Notifier name")

	return cmd
}

func createFeedFiltersDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "delete <id>",
		Aliases: []string{"rm"},
		Short:   "Delete a feed filter",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := NewClient(apiURL)
			ctx := context.Background()

			id := args[0]
			err := client.Delete(ctx, "/feed-filters/"+id)
			if err != nil {
				OutputError(err)
				return err
			}

			if format == "json" {
				return OutputJSON(map[string]string{"status": "deleted", "id": id})
			}

			idInt, _ := strconv.Atoi(id)
			fmt.Printf("Deleted feed filter with ID: %d\n", idInt)
			return nil
		},
	}
}
