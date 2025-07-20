package feedwatcher

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bobbyrward/stronghold/internal/config"
)

func TestParsedEntry_GetKeyValue(t *testing.T) {
	ctx := context.Background()
	entry := &parsedEntry{
		Title:     "Test Book",
		Category:  "Fiction",
		Series:    []string{"Series 1", "Series 2"},
		Authors:   []string{"Author 1", "Author 2"},
		Narrators: []string{"Narrator 1"},
		Summary:   "Test summary",
	}

	tests := []struct {
		name     string
		key      config.FilterKey
		expected []string
	}{
		{
			name:     "get authors",
			key:      config.FilterKey_Author,
			expected: []string{"Author 1", "Author 2"},
		},
		{
			name:     "get series",
			key:      config.FilterKey_Series,
			expected: []string{"Series 1", "Series 2"},
		},
		{
			name:     "get title",
			key:      config.FilterKey_Title,
			expected: []string{"Test Book"},
		},
		{
			name:     "get category",
			key:      config.FilterKey_Category,
			expected: []string{"Fiction"},
		},
		{
			name:     "get summary",
			key:      config.FilterKey_Summary,
			expected: []string{"Test summary"},
		},
		{
			name:     "unknown key returns empty slice",
			key:      config.FilterKey(999),
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := entry.GetKeyValue(ctx, tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestApplyFilterOperator(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name         string
		operator     config.FilterOperator
		actualValues []string
		filterValue  string
		expected     bool
	}{
		{
			name:         "equals - match",
			operator:     config.FilterOperator_Equals,
			actualValues: []string{"test", "value"},
			filterValue:  "test",
			expected:     true,
		},
		{
			name:         "equals - no match",
			operator:     config.FilterOperator_Equals,
			actualValues: []string{"test", "value"},
			filterValue:  "nomatch",
			expected:     false,
		},
		{
			name:         "contains - match",
			operator:     config.FilterOperator_Contains,
			actualValues: []string{"this is a test"},
			filterValue:  "test",
			expected:     true,
		},
		{
			name:         "contains - no match",
			operator:     config.FilterOperator_Contains,
			actualValues: []string{"this is a value"},
			filterValue:  "test",
			expected:     false,
		},
		{
			name:         "fnmatch - wildcard match",
			operator:     config.FilterOperator_Fnmatch,
			actualValues: []string{"test file.txt"},
			filterValue:  "*.txt",
			expected:     true,
		},
		{
			name:         "fnmatch - no match",
			operator:     config.FilterOperator_Fnmatch,
			actualValues: []string{"test file.pdf"},
			filterValue:  "*.txt",
			expected:     false,
		},
		{
			name:         "regex - match",
			operator:     config.FilterOperator_Regex,
			actualValues: []string{"test123"},
			filterValue:  `test\d+`,
			expected:     true,
		},
		{
			name:         "regex - no match",
			operator:     config.FilterOperator_Regex,
			actualValues: []string{"testabc"},
			filterValue:  `test\d+`,
			expected:     false,
		},
		{
			name:         "regex - invalid pattern",
			operator:     config.FilterOperator_Regex,
			actualValues: []string{"test"},
			filterValue:  `[`,
			expected:     false,
		},
		{
			name:         "empty actual values",
			operator:     config.FilterOperator_Equals,
			actualValues: []string{},
			filterValue:  "test",
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := applyFilterOperator(ctx, tt.operator, tt.actualValues, tt.filterValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParsedEntry_HasMatch(t *testing.T) {
	ctx := context.Background()
	entry := &parsedEntry{
		Title:    "Test Book",
		Category: "Fiction",
		Authors:  []string{"John Doe"},
		Series:   []string{"Test Series"},
	}

	tests := []struct {
		name         string
		feedConfig   *config.FeedWatcherConfigFeed
		expectMatch  bool
		expectFilter string
	}{
		{
			name: "single match found",
			feedConfig: &config.FeedWatcherConfigFeed{
				Filters: []config.FeedWatcherConfigFeedFilter{
					{
						Name: "test-filter",
						Matches: []config.FeedWatcherConfigFeedFilterMatch{
							{
								Key:      config.FilterKey_Title,
								Operator: config.FilterOperator_Equals,
								Value:    "Test Book",
							},
						},
					},
				},
			},
			expectMatch:  true,
			expectFilter: "test-filter",
		},
		{
			name: "no match found",
			feedConfig: &config.FeedWatcherConfigFeed{
				Filters: []config.FeedWatcherConfigFeedFilter{
					{
						Name: "test-filter",
						Matches: []config.FeedWatcherConfigFeedFilterMatch{
							{
								Key:      config.FilterKey_Title,
								Operator: config.FilterOperator_Equals,
								Value:    "Different Book",
							},
						},
					},
				},
			},
			expectMatch:  false,
			expectFilter: "",
		},
		{
			name: "multiple filters, second matches",
			feedConfig: &config.FeedWatcherConfigFeed{
				Filters: []config.FeedWatcherConfigFeedFilter{
					{
						Name: "no-match-filter",
						Matches: []config.FeedWatcherConfigFeedFilterMatch{
							{
								Key:      config.FilterKey_Title,
								Operator: config.FilterOperator_Equals,
								Value:    "Different Book",
							},
						},
					},
					{
						Name: "matching-filter",
						Matches: []config.FeedWatcherConfigFeedFilterMatch{
							{
								Key:      config.FilterKey_Author,
								Operator: config.FilterOperator_Contains,
								Value:    "John",
							},
						},
					},
				},
			},
			expectMatch:  true,
			expectFilter: "matching-filter",
		},
		{
			name: "multiple filters, multiple matches, second matches",
			feedConfig: &config.FeedWatcherConfigFeed{
				Filters: []config.FeedWatcherConfigFeedFilter{
					{
						Name: "multiple-match-filter",
						Matches: []config.FeedWatcherConfigFeedFilterMatch{
							{
								Key:      config.FilterKey_Title,
								Operator: config.FilterOperator_Equals,
								Value:    "Different Book",
							},
							{
								Key:      config.FilterKey_Author,
								Operator: config.FilterOperator_Contains,
								Value:    "John",
							},
						},
					},
				},
			},
			expectMatch:  false,
			expectFilter: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matched, filter, err := entry.HasMatch(ctx, tt.feedConfig)
			require.NoError(t, err)
			assert.Equal(t, tt.expectMatch, matched)

			if tt.expectMatch {
				require.NotNil(t, filter)
				assert.Equal(t, tt.expectFilter, filter.Name)
			} else {
				assert.Nil(t, filter)
			}
		})
	}
}

func TestParseDescription(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		description string
		expected    parsedEntry
	}{
		{
			name: "complete description",
			description: `Author(s): John Doe, Jane Smith<br/>
Narrator(s): Bob Reader<br/>
Series: Test Series, Another Series<br/>
Category: Fiction<br/>
Leechers: 5<br/>
Seeders: 10<br/>
Added: 2023-01-01<br/>
Tags: tag1,tag2<br/>
Description: A great book`,
			expected: parsedEntry{
				Authors:     []string{"John Doe", " Jane Smith"},
				Narrators:   []string{"Bob Reader"},
				Series:      []string{"Test Series", " Another Series"},
				Category:    "Fiction",
				Leechers:    5,
				Seeders:     10,
				Added:       "2023-01-01",
				Tags:        "tag1,tag2",
				Description: "A great book",
			},
		},
		{
			name: "partial description",
			description: `Author(s): Single Author<br/>
Category: Non-Fiction<br/>
Seeders: 15`,
			expected: parsedEntry{
				Authors:  []string{"Single Author"},
				Category: "Non-Fiction",
				Seeders:  15,
			},
		},
		{
			name: "invalid numbers",
			description: `Leechers: not-a-number<br/>
Seeders: also-not-a-number`,
			expected: parsedEntry{
				Leechers: 0,
				Seeders:  0,
			},
		},
		{
			name: "empty parts and malformed lines",
			description: `Author(s): Test Author<br/>
<br/>
Invalid line without colon<br/>
Category: Test Category`,
			expected: parsedEntry{
				Authors:  []string{"Test Author"},
				Category: "Test Category",
			},
		},
		{
			name:        "empty description",
			description: "",
			expected:    parsedEntry{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseDescription(ctx, tt.description)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewFeedWatcher(t *testing.T) {
	fw := NewFeedWatcher()
	assert.NotNil(t, fw)
	assert.IsType(t, &FeedWatcher{}, fw)
}

// Benchmark tests
func BenchmarkApplyFilterOperator_Equals(b *testing.B) {
	ctx := context.Background()
	actualValues := []string{"test", "value", "benchmark"}
	filterValue := "benchmark"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		applyFilterOperator(ctx, config.FilterOperator_Equals, actualValues, filterValue)
	}
}

func BenchmarkApplyFilterOperator_Regex(b *testing.B) {
	ctx := context.Background()
	actualValues := []string{"test123", "value456", "benchmark789"}
	filterValue := `\d+`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		applyFilterOperator(ctx, config.FilterOperator_Regex, actualValues, filterValue)
	}
}

func BenchmarkParseDescription(b *testing.B) {
	ctx := context.Background()
	description := `Author(s): John Doe, Jane Smith<br/>
Narrator(s): Bob Reader<br/>
Series: Test Series, Another Series<br/>
Category: Fiction<br/>
Leechers: 5<br/>
Seeders: 10<br/>
Added: 2023-01-01<br/>
Tags: tag1,tag2<br/>
Description: A great book`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parseDescription(ctx, description)
	}
}
