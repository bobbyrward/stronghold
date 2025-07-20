package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestFilterOperator_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name        string
		yamlValue   string
		expected    FilterOperator
		expectError bool
	}{
		{
			name:      "equals operator",
			yamlValue: "equals",
			expected:  FilterOperator_Equals,
		},
		{
			name:      "contains operator",
			yamlValue: "contains",
			expected:  FilterOperator_Contains,
		},
		{
			name:      "fnmatch operator",
			yamlValue: "fnmatch",
			expected:  FilterOperator_Fnmatch,
		},
		{
			name:      "regex operator",
			yamlValue: "regex",
			expected:  FilterOperator_Regex,
		},
		{
			name:        "invalid operator",
			yamlValue:   "invalid",
			expectError: true,
		},
		{
			name:        "empty operator",
			yamlValue:   "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var op FilterOperator
			node := &yaml.Node{
				Kind:  yaml.ScalarNode,
				Value: tt.yamlValue,
			}

			err := op.UnmarshalYAML(node)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "invalid operator")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, op)
			}
		})
	}
}

func TestFilterOperator_MarshalYAML(t *testing.T) {
	tests := []struct {
		name        string
		operator    FilterOperator
		expected    string
		expectError bool
	}{
		{
			name:     "equals operator",
			operator: FilterOperator_Equals,
			expected: "equals",
		},
		{
			name:     "contains operator",
			operator: FilterOperator_Contains,
			expected: "contains",
		},
		{
			name:     "fnmatch operator",
			operator: FilterOperator_Fnmatch,
			expected: "fnmatch",
		},
		{
			name:     "regex operator",
			operator: FilterOperator_Regex,
			expected: "regex",
		},
		{
			name:        "invalid operator",
			operator:    FilterOperator(999),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.operator.MarshalYAML()

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "invalid operator")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestFilterKey_String(t *testing.T) {
	tests := []struct {
		name     string
		key      FilterKey
		expected string
	}{
		{
			name:     "author key",
			key:      FilterKey_Author,
			expected: "author",
		},
		{
			name:     "series key",
			key:      FilterKey_Series,
			expected: "series",
		},
		{
			name:     "title key",
			key:      FilterKey_Title,
			expected: "title",
		},
		{
			name:     "category key",
			key:      FilterKey_Category,
			expected: "category",
		},
		{
			name:     "summary key",
			key:      FilterKey_Summary,
			expected: "summary",
		},
		{
			name:     "invalid key returns empty string",
			key:      FilterKey(999),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.key.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFilterKey_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name        string
		yamlValue   string
		expected    FilterKey
		expectError bool
	}{
		{
			name:      "author key",
			yamlValue: "author",
			expected:  FilterKey_Author,
		},
		{
			name:      "series key",
			yamlValue: "series",
			expected:  FilterKey_Series,
		},
		{
			name:      "title key",
			yamlValue: "title",
			expected:  FilterKey_Title,
		},
		{
			name:      "category key",
			yamlValue: "category",
			expected:  FilterKey_Category,
		},
		{
			name:      "summary key",
			yamlValue: "summary",
			expected:  FilterKey_Summary,
		},
		{
			name:        "invalid key",
			yamlValue:   "invalid",
			expectError: true,
		},
		{
			name:        "empty key",
			yamlValue:   "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var key FilterKey
			node := &yaml.Node{
				Kind:  yaml.ScalarNode,
				Value: tt.yamlValue,
			}

			err := key.UnmarshalYAML(node)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "invalid key")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, key)
			}
		})
	}
}

func TestFilterKey_MarshalYAML(t *testing.T) {
	tests := []struct {
		name        string
		key         FilterKey
		expected    string
		expectError bool
	}{
		{
			name:     "author key",
			key:      FilterKey_Author,
			expected: "author",
		},
		{
			name:     "series key",
			key:      FilterKey_Series,
			expected: "series",
		},
		{
			name:     "title key",
			key:      FilterKey_Title,
			expected: "title",
		},
		{
			name:     "category key",
			key:      FilterKey_Category,
			expected: "category",
		},
		{
			name:     "summary key",
			key:      FilterKey_Summary,
			expected: "summary",
		},
		{
			name:        "invalid key",
			key:         FilterKey(999),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.key.MarshalYAML()

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "invalid key")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestFeedWatcherConfig_YAMLMarshaling(t *testing.T) {
	config := FeedWatcherConfig{
		Feeds: []FeedWatcherConfigFeed{
			{
				Name: "MAM",
				URL:  "https://feed.com/rss/feed",
				Filters: []FeedWatcherConfigFeedFilter{
					{
						Name:     "John Smith Books",
						Category: "personal-books",
						Matches: []FeedWatcherConfigFeedFilterMatch{
							{
								Key:      FilterKey_Author,
								Value:    "John Smith",
								Operator: FilterOperator_Equals,
							},
						},
					},
					{
						Name:         "John Smith Books",
						Category:     "personal-books",
						Notification: "https://webhook.com/webhook",
						Matches: []FeedWatcherConfigFeedFilterMatch{
							{
								Key:      FilterKey_Author,
								Value:    "John Smith",
								Operator: FilterOperator_Equals,
							},
						},
					},
				},
			},
		},
	}

	// Test marshaling
	// yamlData, err := yaml.Marshal(config)
	_, err := yaml.Marshal(config)
	require.NoError(t, err)

	// TODO: Fix these tests

	/*
		// Test unmarshaling
		var unmarshaledConfig FeedWatcherConfig
		err = yaml.Unmarshal(yamlData, &unmarshaledConfig)
		require.NoError(t, err)

		// Verify the unmarshaled config matches the original
		assert.Equal(t, config.Feeds[0].Name, unmarshaledConfig.Feeds[0].Name)
		assert.Equal(t, config.Feeds[0].URL, unmarshaledConfig.Feeds[0].URL)
		assert.Equal(t, config.Feeds[0].Filters[0].Name, unmarshaledConfig.Feeds[0].Filters[0].Name)
		assert.Equal(t, config.Feeds[0].Filters[0].Category, unmarshaledConfig.Feeds[0].Filters[0].Category)
		assert.Equal(t, config.Feeds[0].Filters[0].Matches[0].Key, unmarshaledConfig.Feeds[0].Filters[0].Matches[0].Key)
		assert.Equal(t, config.Feeds[0].Filters[0].Matches[0].Value, unmarshaledConfig.Feeds[0].Filters[0].Matches[0].Value)
		assert.Equal(t, config.Feeds[0].Filters[0].Matches[0].Operator, unmarshaledConfig.Feeds[0].Filters[0].Matches[0].Operator)
	*/
}

func TestComplexYAMLConfig(t *testing.T) {
	yamlData := `
feedWatcher:
  feeds:
    - name: MAM
      url: https://feed.com/rss/feed
      filters:
        - name: John Smith Books
          category: personal-books
          match:
            - key: author
              value: John Smith
              operator: equals
            - key: title
              value: "The*"
              operator: fnmatch
        - name: Sci-Fi Series
          category: sci-fi
          match:
            - key: category
              value: science
              operator: contains
            - key: summary
              value: "space|alien|robot"
              operator: regex
`

	var configWrapper struct {
		FeedWatcher FeedWatcherConfig `yaml:"feedWatcher"`
	}

	err := yaml.Unmarshal([]byte(yamlData), &configWrapper)
	require.NoError(t, err)

	config := configWrapper.FeedWatcher
	require.Len(t, config.Feeds, 1)

	feed := config.Feeds[0]
	assert.Equal(t, "MAM", feed.Name)
	assert.Equal(t, "https://feed.com/rss/feed", feed.URL)
	require.Len(t, feed.Filters, 2)

	// Test first filter
	filter1 := feed.Filters[0]
	assert.Equal(t, "John Smith Books", filter1.Name)
	assert.Equal(t, "personal-books", filter1.Category)
	require.Len(t, filter1.Matches, 2)

	match1 := filter1.Matches[0]
	assert.Equal(t, FilterKey_Author, match1.Key)
	assert.Equal(t, "John Smith", match1.Value)
	assert.Equal(t, FilterOperator_Equals, match1.Operator)

	match2 := filter1.Matches[1]
	assert.Equal(t, FilterKey_Title, match2.Key)
	assert.Equal(t, "The*", match2.Value)
	assert.Equal(t, FilterOperator_Fnmatch, match2.Operator)

	// Test second filter
	filter2 := feed.Filters[1]
	assert.Equal(t, "Sci-Fi Series", filter2.Name)
	assert.Equal(t, "sci-fi", filter2.Category)
	require.Len(t, filter2.Matches, 2)

	match3 := filter2.Matches[0]
	assert.Equal(t, FilterKey_Category, match3.Key)
	assert.Equal(t, "science", match3.Value)
	assert.Equal(t, FilterOperator_Contains, match3.Operator)

	match4 := filter2.Matches[1]
	assert.Equal(t, FilterKey_Summary, match4.Key)
	assert.Equal(t, "space|alien|robot", match4.Value)
	assert.Equal(t, FilterOperator_Regex, match4.Operator)
}

func TestInvalidYAMLUnmarshaling(t *testing.T) {
	tests := []struct {
		name     string
		yamlData string
		errMsg   string
	}{
		{
			name: "invalid operator",
			yamlData: `
feeds:
  - name: test
    url: test.com
    filters:
      - name: test-filter
        category: test
        match:
          - key: author
            value: test
            operator: invalid_operator
`,
			errMsg: "invalid operator",
		},
		{
			name: "invalid key",
			yamlData: `
feeds:
  - name: test
    url: test.com
    filters:
      - name: test-filter
        category: test
        match:
          - key: invalid_key
            value: test
            operator: equals
`,
			errMsg: "invalid key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var config FeedWatcherConfig
			err := yaml.Unmarshal([]byte(tt.yamlData), &config)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}

// Benchmark tests
func BenchmarkFilterOperator_UnmarshalYAML(b *testing.B) {
	node := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: "equals",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var op FilterOperator
		_ = op.UnmarshalYAML(node)
	}
}

func BenchmarkFilterKey_UnmarshalYAML(b *testing.B) {
	node := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: "author",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var key FilterKey
		_ = key.UnmarshalYAML(node)
	}
}

func BenchmarkCompleteConfigUnmarshal(b *testing.B) {
	yamlData := []byte(`
feeds:
  - name: MAM
    url: https://feed.com/rss/feed
    filters:
      - name: John Smith Books
        category: personal-books
        match:
          - key: author
            value: John Smith
            operator: equals
`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var config FeedWatcherConfig
		_ = yaml.Unmarshal(yamlData, &config)
	}
}
