package apiclient

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func createFeedAuthorFiltersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "feed-author-filters",
		Short: "Manage feed author filters",
	}

	cmd.AddCommand(createFeedAuthorFiltersListCmd())
	cmd.AddCommand(createFeedAuthorFiltersGetCmd())
	cmd.AddCommand(createFeedAuthorFiltersCreateCmd())
	cmd.AddCommand(createFeedAuthorFiltersUpdateCmd())
	cmd.AddCommand(createFeedAuthorFiltersDeleteCmd())

	return cmd
}

func createFeedAuthorFiltersListCmd() *cobra.Command {
	var feedID string

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all feed author filters",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := NewClient(apiURL)
			ctx := context.Background()

			path := "/feed-author-filters"
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
			headers := []string{"ID", "Author", "Feed", "Category", "Notifier"}
			rows := [][]string{}
			for _, filter := range filters {
				id := fmt.Sprintf("%.0f", filter["id"].(float64))
				author := fmt.Sprintf("%v", filter["author"])

				// Format feed as "feedname(id)"
				feedStr := "N/A"
				if feedName, ok := filter["feed_name"].(string); ok {
					feedID := uint(filter["feed_id"].(float64))
					feedStr = FormatRelation(feedName, feedID)
				}

				// Format category
				categoryStr := fmt.Sprintf("%v", filter["category"])

				// Format notifier
				notifierStr := fmt.Sprintf("%v", filter["notifier"])

				rows = append(rows, []string{id, author, feedStr, categoryStr, notifierStr})
			}

			return OutputTable(headers, rows)
		},
	}

	cmd.Flags().StringVar(&feedID, "feed-id", "", "Filter by feed ID")

	return cmd
}

func createFeedAuthorFiltersGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a feed author filter by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := NewClient(apiURL)
			ctx := context.Background()

			id := args[0]
			var filter map[string]interface{}
			err := client.Get(ctx, "/feed-author-filters/"+id, &filter)
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
				{"Author", fmt.Sprintf("%v", filter["author"])},
				{"Feed", feedStr},
				{"Category", fmt.Sprintf("%v", filter["category"])},
				{"Notifier", fmt.Sprintf("%v", filter["notifier"])},
			}

			return OutputTable(headers, rows)
		},
	}
}

func createFeedAuthorFiltersCreateCmd() *cobra.Command {
	var author, feedName, categoryName, notifierName string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new feed author filter",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := NewClient(apiURL)
			ctx := context.Background()

			body := map[string]interface{}{
				"author":        author,
				"feed_name":     feedName,
				"category_name": categoryName,
				"notifier_name": notifierName,
			}

			var filter map[string]interface{}
			err := client.Post(ctx, "/feed-author-filters", body, &filter)
			if err != nil {
				OutputError(err)
				return err
			}

			if format == "json" {
				return OutputJSON(filter)
			}

			fmt.Printf("Created feed author filter with ID: %.0f\n", filter["id"].(float64))
			return nil
		},
	}

	cmd.Flags().StringVar(&author, "author", "", "Author name (required)")
	cmd.Flags().StringVar(&feedName, "feed-name", "", "Feed name (required)")
	cmd.Flags().StringVar(&categoryName, "category-name", "", "Category name (required)")
	cmd.Flags().StringVar(&notifierName, "notifier-name", "", "Notifier name (required)")
	_ = cmd.MarkFlagRequired("author")
	_ = cmd.MarkFlagRequired("feed-name")
	_ = cmd.MarkFlagRequired("category-name")
	_ = cmd.MarkFlagRequired("notifier-name")

	return cmd
}

func createFeedAuthorFiltersUpdateCmd() *cobra.Command {
	var author, feedName, categoryName, notifierName string

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update an existing feed author filter",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := NewClient(apiURL)
			ctx := context.Background()

			id := args[0]
			body := map[string]interface{}{
				"author":        author,
				"feed_name":     feedName,
				"category_name": categoryName,
				"notifier_name": notifierName,
			}

			var filter map[string]interface{}
			err := client.Put(ctx, "/feed-author-filters/"+id, body, &filter)
			if err != nil {
				OutputError(err)
				return err
			}

			if format == "json" {
				return OutputJSON(filter)
			}

			fmt.Printf("Updated feed author filter with ID: %s\n", id)
			return nil
		},
	}

	cmd.Flags().StringVar(&author, "author", "", "Author name (required)")
	cmd.Flags().StringVar(&feedName, "feed-name", "", "Feed name (optional, omit to keep current)")
	cmd.Flags().StringVar(&categoryName, "category-name", "", "Category name (required)")
	cmd.Flags().StringVar(&notifierName, "notifier-name", "", "Notifier name (required)")
	_ = cmd.MarkFlagRequired("author")
	_ = cmd.MarkFlagRequired("category-name")
	_ = cmd.MarkFlagRequired("notifier-name")

	return cmd
}

func createFeedAuthorFiltersDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a feed author filter",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := NewClient(apiURL)
			ctx := context.Background()

			id := args[0]
			err := client.Delete(ctx, "/feed-author-filters/"+id)
			if err != nil {
				OutputError(err)
				return err
			}

			fmt.Printf("Deleted feed author filter with ID: %s\n", id)
			return nil
		},
	}
}
