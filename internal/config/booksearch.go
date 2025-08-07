package config

type BookSearchConfig struct {
	BaseURL        string `yaml:"baseUrl"`
	SearchEndpoint string `yaml:"searchEndpoint"`
	APIKey         string `yaml:"apiKey"`
	HttpProxy      string `yaml:"httpProxy"`
	HttpsProxy     string `yaml:"httpsProxy"`
}
