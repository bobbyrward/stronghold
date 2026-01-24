package logging

import (
	"log/slog"
	"os"

	"github.com/bobbyrward/stronghold/internal/config"
)

var logLevel = new(slog.LevelVar)

func SetLoggingLevel(level config.LoggingLevel) {
	logLevel.Set(slog.Level(level))
}

func SetupLogging(level config.LoggingLevel) {
	logLevel.Set(slog.Level(level))

	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	}))

	slog.SetDefault(logger)
}
