package config

type FeedWatcherConfig struct{}

type FeedWatcherConfigFeed struct {
	URL     string                  `yaml:"url"`
	Filters []FeedWatcherConfigFeed `yaml:"filters"`
}

type FeedWatcherConfigFeedFilters struct {
	Match    string `yaml:"match"`
	Category string `yaml:"category"`
}

type FeedWatcherConfigFeedFilterMatchAuthor struct {
	Value    string `yaml:"value"`
	Operator string `yaml:"value"`
}

/*

feedWatcher:
  url: asdasdasdasdasd
  filters:
    match:
      - key: author
        operator: equals
        value: Blah H. Blahher
      - key: series
        operator: contains
        value: asdasdasd
      - key: titles
        operator: fnmatch
        value: def*ijkl
      - key: category
        operator: regex
        value: ^def.+ijkl$
      - key: summary
        operator: equals
        value: qwe
      category: personal-books
*/
