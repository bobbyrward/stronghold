package config

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var Config ClusterConfig

func generateDefaultConfig() string {
	return `# Stronghold Configuration File
# Logging Configuration
logging:
  level: ""

# API Configuration
apiClient:
  url: "http://localhost:8000"  # API server URL for CLI client

# PostgreSQL Database Configuration
postgres:
  url: ""  # Example: postgresql://user:password@localhost:5432/stronghold

# qBittorrent Configuration
qbit:
  url: ""       # Example: http://localhost:8080/
  username: ""
  password: ""
  downloadPath: ""  # Base download path for qBittorrent
  localDownloadPath: ""  # Local path where downloads are stored

importers:
  importedTag: imported
  manualInterventionTag: needs-manual

  ebooks:
    libraries: []
    # Example:
    # - name: personal-book
    #   path: /mnt/other/books/incoming
    importTypes: []
    # Example:
    # - category: books
    #   calibreDesktopURL: https://calibre-desktop.example.com/
    #   calibreWebURL: https://calibre.example.com
    #   notification: my-discord-notifier

  audiobooks:
    libraries: []
    # Example:
    # - name: audiobooks
    #   path: /audiobooks
    importTypes: []
    # Example:
    # - category: audiobooks
    #   library: audiobooks
    #   notification: my-discord-notifier



# Notification Configuration
notifications:
  notifiers: []
  # Example:
  # - name: my-discord-notifier
  #   type: discord
  #   url: https://discord.com/api/webhooks/...

# Book Importer Configuration
bookImporter:
  importTypes: []
  # Example:
  # - category: books
  #   importedTag: imported
  #   sourcePath: /path/to/source
  #   sourcePrefixPath: /data/
  #   destinationPath: /path/to/destination
  #   calibreDesktopURL: https://calibre-desktop.example.com/
  #   calibreWebURL: https://calibre.example.com
  #   manualInterventionTag: needs-manual

# Audiobook Import Configuration
audiobookImport:
  libraries: []
  # - name: ""
  #   path: ""

# Feed Watcher Configuration
feedWatcher:
  feeds: []
  # Example:
  # - name: My Feed
  #   url: https://example.com/rss
  #   filters:
  #     - name: Filter Name
  #       category: my-category
  #       notification: my-discord-notifier
  #       match:
  #         - key: author
  #           operator: contains
  #           value: Author Name

# Discord Bot Configuration
discordBot:
  token: ""
  guildId: ""
  bookSearchApi:
    baseUrl: ""
    searchEndpoint: ""
    httpProxy: ""
    httpsProxy: ""

# Book Search Configuration
bookSearch:
  baseUrl: ""
  searchEndpoint: ""
  cookieDomain: ""
  tokenRefreshUrl: ""
  httpProxy: ""
  httpsProxy: ""
`
}

func LoadConfig(configFilePath string) error {
	ctx := context.Background()

	// Determine config file path if not provided
	if configFilePath == "" {
		// Check environment variable
		configFilePath = os.Getenv("STRONGHOLD_CONFIG")

		if configFilePath == "" {
			// Use XDG_CONFIG_HOME or default to ~/.config
			xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
			if xdgConfigHome == "" {
				home, err := os.UserHomeDir()
				if err != nil {
					return fmt.Errorf("failed to get user home directory: %w", err)
				}
				xdgConfigHome = filepath.Join(home, ".config")
			}
			configFilePath = filepath.Join(xdgConfigHome, "stronghold", "config.yaml")
		}
	}

	// Check if config file exists
	_, err := os.Stat(configFilePath)
	if os.IsNotExist(err) {
		// Create parent directories
		configDir := filepath.Dir(configFilePath)
		err = os.MkdirAll(configDir, 0755)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to create config directory",
				slog.String("path", configDir),
				slog.Any("err", err))
			return fmt.Errorf("failed to create config directory: %w", err)
		}

		// Write default config
		defaultConfig := generateDefaultConfig()
		err = os.WriteFile(configFilePath, []byte(defaultConfig), 0644)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to write default config file",
				slog.String("path", configFilePath),
				slog.Any("err", err))
			return fmt.Errorf("failed to write default config file: %w", err)
		}

		slog.InfoContext(ctx, "Created default configuration file",
			slog.String("path", configFilePath))
	}

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
