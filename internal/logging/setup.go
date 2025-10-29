package logging

import (
	"log/slog"
	"os"

	"github.com/bobbyrward/stronghold/internal/config"
)

func SetupLogging(level config.LoggingLevel) {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.Level(level),
	}))

	slog.SetDefault(logger)
}
