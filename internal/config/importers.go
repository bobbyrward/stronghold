package config

type ImportLibrary struct {
	Name string `yaml:"name" json:"name"`
	Path string `yaml:"path" json:"path"`
}

type ImportType struct {
	Category          string `yaml:"category"`
	Library           string `yaml:"library"`
	CalibreDesktopURL string `yaml:"calibreDesktopURL"`
	CalibreWebURL     string `yaml:"calibreWebURL"`
	DiscordNotifier   string `yaml:"notification"`
}

type AudiobookImporter struct {
	Libraries   []ImportLibrary `yaml:"libraries"`
	ImportTypes []ImportType    `yaml:"importTypes"`
}

type BookImporter struct {
	Libraries   []ImportLibrary `yaml:"libraries"`
	ImportTypes []ImportType    `yaml:"importTypes"`
}

type ImportersConfig struct {
	ImportedTag           string            `yaml:"importedTag"`
	ManualInterventionTag string            `yaml:"manualInterventionTag"`
	BookImporter          BookImporter      `yaml:"ebooks"`
	AudiobookImporter     AudiobookImporter `yaml:"audiobooks"`
}

func FindLibraryByName(libraries []ImportLibrary, libraryName string) (*ImportLibrary, bool) {
	for _, lib := range libraries {
		if lib.Name == libraryName {
			return &lib, true
		}
	}

	return nil, false
}
