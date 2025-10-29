package config

type ClusterConfig struct {
	Postgres      PostgresConfig      `yaml:"postgres"`
	Qbit          QbitConfig          `yaml:"qbit"`
	Notifications NotificationsConfig `yaml:"notifications"`
	FeedWatcher   FeedWatcherConfig   `yaml:"feedWatcher"`
	DiscordBot    DiscordBotConfig    `yaml:"discordBot"`
	BookSearch    BookSearchConfig    `yaml:"bookSearch"`
	Logging       LoggingConfig       `yaml:"logging"`
	APIClient     APIClientConfig     `yaml:"apiClient"`
	Importers     ImportersConfig     `yaml:"importers"`
}
