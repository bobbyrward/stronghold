package booksearch

import (
	"github.com/spf13/cobra"
)

func CreateBookSearchCmd() *cobra.Command {
	bookSearchCmd := &cobra.Command{
		Use:   "book-search",
		Short: "Interact with the BookSearch service",
	}

	bookSearchCmd.AddCommand(createSearchCommand())
	bookSearchCmd.AddCommand(createGetByHashCommand())
	bookSearchCmd.AddCommand(createGetByIDCommand())

	return bookSearchCmd
}
