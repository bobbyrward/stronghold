package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/bobbyrward/stronghold/internal/config"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:          "stronghold",
		SilenceUsage: true,
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(onCobraInit)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "Path to the config file")
	rootCmd.AddCommand(createBookImportCmd())
	rootCmd.AddCommand(createFeedWatcherCmd())
	rootCmd.AddCommand(createApiCmd())
}

func internalCobraInit() error {
	err := config.LoadConfig(cfgFile)
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to load config"))
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	return nil
}

func onCobraInit() {
	err := internalCobraInit()
	if err != nil {
		fmt.Println("Error initializing Cobra:", err)
		os.Exit(1)
	}
}
