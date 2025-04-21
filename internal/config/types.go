package config

type ClusterConfig struct {
	BookImporter BookImportConfig `yaml:"bookImporter"`
}

type BookImportConfig struct {
	QbitURL      string                 `yaml:"qbitURL"`
	QbitUsername string                 `yaml:"qbitUsername"`
	QbitPassword string                 `yaml:"qbitPassword"`
	ImportTypes  []BookImportTypeConfig `yaml:"importTypes"`
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
}
