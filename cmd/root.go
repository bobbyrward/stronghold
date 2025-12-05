package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/bobbyrward/stronghold/internal/config"
	"github.com/bobbyrward/stronghold/internal/logging"
)

var (
	cfgFile  string
	logLevel string
	rootCmd  = &cobra.Command{
		Use:          "stronghold",
		SilenceUsage: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			level := config.Config.Logging.Level

			if logLevel != "" {
				level = config.LoggingLevelFromString(logLevel)
			}

			logging.SetupLogging(level)
		},
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

	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "", "Log level: debug, info, warn, error, none")
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "Path to the config file (default: $XDG_CONFIG_HOME/stronghold/config.yaml or $STRONGHOLD_CONFIG)")
	rootCmd.AddCommand(createBookImportCmd())
	rootCmd.AddCommand(createFeedWatcherCmd())
	rootCmd.AddCommand(createApiCmd())
	rootCmd.AddCommand(createDiscordBotCmd())
	rootCmd.AddCommand(createBookSearchCmd())
	rootCmd.AddCommand(createRefreshTokenCmd())
	rootCmd.AddCommand(createAudiobookImporterCmd())
	rootCmd.AddCommand(createDoctorCmd())
}

func internalCobraInit() error {
	logging.SetupLogging(config.LoggingLevel_Warn)

	err := config.LoadConfig(cfgFile)
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to load config"))
	}

	return nil
}

func onCobraInit() {
	err := internalCobraInit()
	if err != nil {
		// Use fmt.Println here since slog may not be initialized yet
		fmt.Println("Error initializing Cobra:", err)
		os.Exit(1)
	}
}
