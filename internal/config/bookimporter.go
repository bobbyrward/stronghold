package config

type BookImportConfig struct {
	ImportTypes []BookImportTypeConfig `yaml:"importTypes"`
}

type BookImportTypeConfig struct {
	Category              string `yaml:"category"`
	ImportedTag           string `yaml:"importedTag"`
	SourcePath            string `yaml:"sourcePath"`
	SourcePrefixPath      string `yaml:"sourcePrefixPath"`
	DestinationPath       string `yaml:"destinationPath"`
	CalibreDesktopURL     string `yaml:"calibreDesktopURL"`
	CalibreWebURL         string `yaml:"calibreWebURL"`
	ManualInterventionTag string `yaml:"manualInterventionTag"`
	DiscordNotifier       string `yaml:"notification"`
}
