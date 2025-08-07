package cmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

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
	ctx := context.Background()
	
	slog.InfoContext(ctx, "Starting API command")
	
	err := api.Run()
	if err != nil {
		slog.ErrorContext(ctx, "API server failed", slog.Any("err", err))
		return errors.Join(err, fmt.Errorf("failed to run api"))
	}

	slog.InfoContext(ctx, "API server shut down gracefully")
	return nil
}
