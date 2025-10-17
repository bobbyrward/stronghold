package apiclient

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func createFeedsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "feeds",
		Short: "Manage feeds",
	}

	cmd.AddCommand(createFeedsListCmd())
	cmd.AddCommand(createFeedsGetCmd())
	cmd.AddCommand(createFeedsCreateCmd())
	cmd.AddCommand(createFeedsUpdateCmd())
	cmd.AddCommand(createFeedsDeleteCmd())

	return cmd
}

func createFeedsListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all feeds",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := NewClient(apiURL)
			ctx := context.Background()

			var feeds []map[string]interface{}
			err := client.Get(ctx, "/feeds", &feeds)
			if err != nil {
				OutputError(err)
				return err
			}

			if format == "json" {
				return OutputJSON(feeds)
			}

			// Table output
			headers := []string{"ID", "Name", "URL", "Enabled"}
			rows := [][]string{}
			for _, feed := range feeds {
				id := fmt.Sprintf("%.0f", feed["id"].(float64))
				name := feed["name"].(string)
				url := feed["url"].(string)
				enabled := "false"
				if e, ok := feed["enabled"].(bool); ok && e {
					enabled = "true"
				}
				rows = append(rows, []string{id, name, url, enabled})
			}

			return OutputTable(headers, rows)
		},
	}
}

func createFeedsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a feed by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := NewClient(apiURL)
			ctx := context.Background()

			id := args[0]
			var feed map[string]interface{}
			err := client.Get(ctx, "/feeds/"+id, &feed)
			if err != nil {
				OutputError(err)
				return err
			}

			if format == "json" {
				return OutputJSON(feed)
			}

			// Table output
			headers := []string{"Field", "Value"}
			rows := [][]string{
				{"ID", fmt.Sprintf("%.0f", feed["id"].(float64))},
				{"Name", feed["name"].(string)},
				{"URL", feed["url"].(string)},
				{"Enabled", fmt.Sprintf("%v", feed["enabled"])},
			}

			return OutputTable(headers, rows)
		},
	}
}

func createFeedsCreateCmd() *cobra.Command {
	var name, url string
	var enabled bool

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new feed",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := NewClient(apiURL)
			ctx := context.Background()

			body := map[string]interface{}{
				"name":    name,
				"url":     url,
				"enabled": enabled,
			}

			var feed map[string]interface{}
			err := client.Post(ctx, "/feeds", body, &feed)
			if err != nil {
				OutputError(err)
				return err
			}

			if format == "json" {
				return OutputJSON(feed)
			}

			fmt.Printf("Created feed: %s (ID: %.0f)\n", feed["name"].(string), feed["id"].(float64))
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Feed name (required)")
	cmd.Flags().StringVar(&url, "url", "", "Feed URL (required)")
	cmd.Flags().BoolVar(&enabled, "enabled", true, "Whether the feed is enabled")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("url")

	return cmd
}

func createFeedsUpdateCmd() *cobra.Command {
	var name, url string
	var enabled bool

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a feed",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := NewClient(apiURL)
			ctx := context.Background()

			id := args[0]
			body := map[string]interface{}{}

			if cmd.Flags().Changed("name") {
				body["name"] = name
			}
			if cmd.Flags().Changed("url") {
				body["url"] = url
			}
			if cmd.Flags().Changed("enabled") {
				body["enabled"] = enabled
			}

			var feed map[string]interface{}
			err := client.Put(ctx, "/feeds/"+id, body, &feed)
			if err != nil {
				OutputError(err)
				return err
			}

			if format == "json" {
				return OutputJSON(feed)
			}

			fmt.Printf("Updated feed: %s (ID: %.0f)\n", feed["name"].(string), feed["id"].(float64))
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Feed name")
	cmd.Flags().StringVar(&url, "url", "", "Feed URL")
	cmd.Flags().BoolVar(&enabled, "enabled", true, "Whether the feed is enabled")

	return cmd
}

func createFeedsDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "delete <id>",
		Aliases: []string{"rm"},
		Short:   "Delete a feed",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := NewClient(apiURL)
			ctx := context.Background()

			id := args[0]
			err := client.Delete(ctx, "/feeds/"+id)
			if err != nil {
				OutputError(err)
				return err
			}

			if format == "json" {
				return OutputJSON(map[string]string{"status": "deleted", "id": id})
			}

			// Convert id to int for display
			idInt, _ := strconv.Atoi(id)
			fmt.Printf("Deleted feed with ID: %d\n", idInt)
			return nil
		},
	}
}
