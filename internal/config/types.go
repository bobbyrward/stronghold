package config

type ClusterConfig struct {
	Qbit         QbitConfig        `yaml:"qbit"`
	BookImporter BookImportConfig  `yaml:"bookImporter"`
	FeedWatcher  FeedWatcherConfig `yaml:"feedWatcher"`
}
