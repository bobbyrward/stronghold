package config

type ClusterConfig struct {
	Postgres      PostgresConfig      `yaml:"postgres"`
	Qbit          QbitConfig          `yaml:"qbit"`
	Notifications NotificationsConfig `yaml:"notifications"`
	BookImporter  BookImportConfig    `yaml:"bookImporter"`
	FeedWatcher   FeedWatcherConfig   `yaml:"feedWatcher"`
	DiscordBot    DiscordBotConfig    `yaml:"discordBot"`
	BookSearch    BookSearchConfig    `yaml:"bookSearch"`
	Logging       LoggingConfig       `yaml:"logging"`
}
