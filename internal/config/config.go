package config

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

var Config ClusterConfig

func LoadConfig(configFilePath string) error {
	ctx := context.Background()
	
	slog.InfoContext(ctx, "Loading configuration", slog.String("path", configFilePath))
	
	file, err := os.Open(configFilePath)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to open config file", 
			slog.String("path", configFilePath),
			slog.Any("err", err))
		return errors.Join(err, fmt.Errorf("failed to open config file"))
	}

	defer func() { _ = file.Close() }()

	decoder := yaml.NewDecoder(file)

	err = decoder.Decode(&Config)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to decode config file", 
			slog.String("path", configFilePath),
			slog.Any("err", err))
		return errors.Join(err, fmt.Errorf("failed to decode config file"))
	}

	slog.InfoContext(ctx, "Successfully loaded configuration", 
		slog.String("path", configFilePath))
	
	return nil
}
