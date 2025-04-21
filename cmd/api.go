package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bobbyrward/stronghold/internal/api"
)

func createApiCmd() *cobra.Command {
	apiCmd := &cobra.Command{
		Use:  "api",
		RunE: runApiCmd,
	}

	return apiCmd
}

func runApiCmd(cmd *cobra.Command, args []string) error {
	err := api.Run()
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to run api"))
	}

	return nil
}
