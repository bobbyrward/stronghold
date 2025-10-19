package config

import (
	"log/slog"

	"gopkg.in/yaml.v3"
)

type LoggingLevel slog.Level

const (
	LoggingLevel_None  = LoggingLevel(1000)
	LoggingLevel_Debug = LoggingLevel(slog.LevelDebug)
	LoggingLevel_Info  = LoggingLevel(slog.LevelInfo)
	LoggingLevel_Warn  = LoggingLevel(slog.LevelWarn)
	LoggingLevel_Error = LoggingLevel(slog.LevelError)
)

type LoggingConfig struct {
	Level LoggingLevel `yaml:"level"`
}

func LoggingLevelFromString(levelStr string) LoggingLevel {
	switch levelStr {
	case "dbg":
		return LoggingLevel_Debug
	case "debug":
		return LoggingLevel_Debug
	case "info":
		return LoggingLevel_Info
	case "warn":
		return LoggingLevel_Warn
	case "warning":
		return LoggingLevel_Warn
	case "error":
		return LoggingLevel_Error
	case "":
		return LoggingLevel_Info
	case "none":
		return LoggingLevel_None
	default:
		return LoggingLevel_None
	}
}

func LoggingLevelToString(level LoggingLevel) string {
	switch level {
	case LoggingLevel_Debug:
		return "debug"
	case LoggingLevel_Info:
		return "info"
	case LoggingLevel_Warn:
		return "warn"
	case LoggingLevel_Error:
		return "error"
	case LoggingLevel_None:
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
