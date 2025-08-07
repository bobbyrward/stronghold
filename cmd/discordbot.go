package cmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/bobbyrward/stronghold/internal/discordbot"
)

func createDiscordBotCmd() *cobra.Command {
	discordBotCmd := &cobra.Command{
		Use:  "discord-bot",
		RunE: runDiscordBotCmd,
	}

	return discordBotCmd
}

func runDiscordBotCmd(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	
	slog.InfoContext(ctx, "Starting Discord bot command")
	
	err := discordbot.Run()
	if err != nil {
		slog.ErrorContext(ctx, "Discord bot failed", slog.Any("err", err))
		return errors.Join(err, fmt.Errorf("failed to run discord bot"))
	}

	slog.InfoContext(ctx, "Discord bot shut down gracefully")
	return nil
}