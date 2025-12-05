package config

type BookSearchConfig struct {
	BaseURL         string `yaml:"baseUrl"`
	SearchEndpoint  string `yaml:"searchEndpoint"`
	CookieDomain    string `yaml:"cookieDomain"`
	TokenRefreshURL string `yaml:"tokenRefreshUrl"`
	TokenCookieName string `yaml:"tokenCookieName"`
	HttpProxy       string `yaml:"httpProxy"`
	HttpsProxy      string `yaml:"httpsProxy"`
}
