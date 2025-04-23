package config

type NotificationsConfig struct {
	Notifiers []NotificationsConfigNotifier `yaml:"notifiers"`
}

type NotificationsConfigNotifier struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	Url  string `yaml:"url"`
}
