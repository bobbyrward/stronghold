package config

type QbitConfig struct {
	URL               string `yaml:"url"`
	Username          string `yaml:"username"`
	Password          string `yaml:"password"`
	DownloadPath      string `yaml:"downloadPath"`
	LocalDownloadPath string `yaml:"localDownloadPath"`
}
