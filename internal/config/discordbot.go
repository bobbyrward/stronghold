package config

type DiscordBotConfig struct {
	Token         string           `yaml:"token"`
	GuildID       string           `yaml:"guildId"`
	BookSearchAPI BookSearchConfig `yaml:"bookSearchApi"`
}

