package config

import (
	"log/slog"

	"gopkg.in/yaml.v3"
)

type LoggingLevel slog.Level

type LoggingConfig struct {
	Level LoggingLevel `yaml:"level"`
}

func LoggingLevelFromString(levelStr string) LoggingLevel {
	switch levelStr {
	case "dbg":
		return LoggingLevel(slog.LevelDebug)
	case "debug":
		return LoggingLevel(slog.LevelDebug)
	case "info":
		return LoggingLevel(slog.LevelInfo)
	case "warn":
		return LoggingLevel(slog.LevelWarn)
	case "warning":
		return LoggingLevel(slog.LevelWarn)
	case "error":
		return LoggingLevel(slog.LevelError)
	case "":
		return LoggingLevel(slog.LevelInfo)
	case "none":
		return LoggingLevel(1000)
	default:
		return LoggingLevel(1000)
	}
}

func LoggingLevelToString(level LoggingLevel) string {
	switch level {
	case LoggingLevel(slog.LevelDebug):
		return "debug"
	case LoggingLevel(slog.LevelInfo):
		return "info"
	case LoggingLevel(slog.LevelWarn):
		return "warn"
	case LoggingLevel(slog.LevelError):
		return "error"
	case LoggingLevel(1000):
		return "none"
	default:
		return "none"
	}
}

func (level *LoggingLevel) UnmarshalYAML(value *yaml.Node) error {
	var strValue string

	err := value.Decode(&strValue)
	if err != nil {
		return err
	}

	*level = LoggingLevelFromString(strValue)

	return nil
}

func (level *LoggingLevel) MarshalYAML() (interface{}, error) {
	strValue := LoggingLevelToString(*level)

	return strValue, nil
}
