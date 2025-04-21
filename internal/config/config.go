package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var Config ClusterConfig

func LoadConfig(configFilePath string) error {
	file, err := os.Open(configFilePath)
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to open config file"))
	}

	defer func() { _ = file.Close() }()

	decoder := yaml.NewDecoder(file)

	err = decoder.Decode(&Config)
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to decode config file"))
	}

	return nil
}
