package feedwatcher

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bobbyrward/stronghold/internal/config"
)

func TestParsedEntry_GetKeyValue(t *testing.T) {
	ctx := context.Background()
	entry := &parsedEntry{
		Title:       "Test Book",
		Category:    "Fiction",
		Series:      []string{"Series 1", "Series 2"},
		Authors:     []string{"Author 1", "Author 2"},
		Narrators:   []string{"Narrator 1"},
		Summary:     "Test summary",
		Tags:        "tag1,tag2,tag3",
		Description: "Test description",
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
			name:     "get tags",
			key:      config.FilterKey_Tags,
			expected: []string{"tag1,tag2,tag3"},
		},
		{
			name:     "get description",
			key:      config.FilterKey_Description,
			expected: []string{"Test description"},
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

func TestParsedEntry_HasMatch_WithPreprocessedAuthorFilters(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name         string
		entry        *parsedEntry
		feedConfig   config.FeedWatcherConfigFeed
		expectMatch  bool
		expectFilter string
	}{
		{
			name: "preprocessed author filter matches ebook",
			entry: &parsedEntry{
				Title:    "The Great Book",
				Category: "Ebooks - Fantasy",
				Authors:  []string{"John Smith"},
				Series:   []string{"Fantasy Series"},
			},
			feedConfig: config.FeedWatcherConfigFeed{
				Name: "TestFeed",
				URL:  "https://example.com/feed",
				AuthorFilters: []config.FeedWatcherConfigFilterByAuthor{
					{
						Author: "John Smith",
					},
				},
			},
			expectMatch:  true,
			expectFilter: "John Smith Books",
		},
		{
			name: "preprocessed author filter matches audiobook",
			entry: &parsedEntry{
				Title:    "The Great Audiobook",
				Category: "Audiobooks - Science Fiction",
				Authors:  []string{"Jane Doe"},
				Series:   []string{"Sci-Fi Series"},
			},
			feedConfig: config.FeedWatcherConfigFeed{
				Name: "TestFeed",
				URL:  "https://example.com/feed",
				AuthorFilters: []config.FeedWatcherConfigFilterByAuthor{
					{
						Author: "Jane Doe",
					},
				},
			},
			expectMatch:  true,
			expectFilter: "Jane Doe Audiobooks",
		},
		{
			name: "preprocessed author filter does not match wrong author",
			entry: &parsedEntry{
				Title:    "The Great Book",
				Category: "Ebooks - Fantasy",
				Authors:  []string{"Different Author"},
				Series:   []string{"Fantasy Series"},
			},
			feedConfig: config.FeedWatcherConfigFeed{
				Name: "TestFeed",
				URL:  "https://example.com/feed",
				AuthorFilters: []config.FeedWatcherConfigFilterByAuthor{
					{
						Author: "John Smith",
					},
				},
			},
			expectMatch:  false,
			expectFilter: "",
		},
		{
			name: "normal filter matches when author filters also present",
			entry: &parsedEntry{
				Title:    "Special Series Book",
				Category: "Ebooks - Mystery",
				Authors:  []string{"Unknown Author"},
				Series:   []string{"Mystery Series"},
			},
			feedConfig: config.FeedWatcherConfigFeed{
				Name: "TestFeed",
				URL:  "https://example.com/feed",
				Filters: []config.FeedWatcherConfigFeedFilter{
					{
						Name:     "Mystery Series Filter",
						Category: "mystery-books",
						Matches: []config.FeedWatcherConfigFeedFilterMatch{
							{
								Key:      config.FilterKey_Series,
								Operator: config.FilterOperator_Contains,
								Value:    "Mystery Series",
							},
						},
					},
				},
				AuthorFilters: []config.FeedWatcherConfigFilterByAuthor{
					{
						Author: "John Smith",
					},
				},
			},
			expectMatch:  true,
			expectFilter: "Mystery Series Filter",
		},
		{
			name: "normal filter takes precedence when it comes first",
			entry: &parsedEntry{
				Title:    "Premium Book",
				Category: "Ebooks - Fantasy",
				Authors:  []string{"John Smith"},
				Series:   []string{"Premium Series"},
			},
			feedConfig: config.FeedWatcherConfigFeed{
				Name: "TestFeed",
				URL:  "https://example.com/feed",
				Filters: []config.FeedWatcherConfigFeedFilter{
					{
						Name:     "Premium Series Filter",
						Category: "premium-books",
						Matches: []config.FeedWatcherConfigFeedFilterMatch{
							{
								Key:      config.FilterKey_Series,
								Operator: config.FilterOperator_Contains,
								Value:    "Premium",
							},
						},
					},
				},
				AuthorFilters: []config.FeedWatcherConfigFilterByAuthor{
					{
						Author: "John Smith",
					},
				},
			},
			expectMatch:  true,
			expectFilter: "Premium Series Filter",
		},
		{
			name: "author filter takes precedence when normal filter doesn't match",
			entry: &parsedEntry{
				Title:    "Regular Book",
				Category: "Ebooks - Fantasy",
				Authors:  []string{"John Smith"},
				Series:   []string{"Regular Series"},
			},
			feedConfig: config.FeedWatcherConfigFeed{
				Name: "TestFeed",
				URL:  "https://example.com/feed",
				Filters: []config.FeedWatcherConfigFeedFilter{
					{
						Name:     "Premium Series Filter",
						Category: "premium-books",
						Matches: []config.FeedWatcherConfigFeedFilterMatch{
							{
								Key:      config.FilterKey_Series,
								Operator: config.FilterOperator_Contains,
								Value:    "Premium",
							},
						},
					},
				},
				AuthorFilters: []config.FeedWatcherConfigFilterByAuthor{
					{
						Author: "John Smith",
					},
				},
			},
			expectMatch:  true,
			expectFilter: "John Smith Books",
		},
		{
			name: "multiple author filters create multiple preprocessed filters",
			entry: &parsedEntry{
				Title:    "The Book",
				Category: "Audiobooks - Fantasy",
				Authors:  []string{"Jane Doe"},
				Series:   []string{"Fantasy Series"},
			},
			feedConfig: config.FeedWatcherConfigFeed{
				Name: "TestFeed",
				URL:  "https://example.com/feed",
				AuthorFilters: []config.FeedWatcherConfigFilterByAuthor{
					{
						Author: "John Smith",
					},
					{
						Author: "Jane Doe",
					},
				},
			},
			expectMatch:  true,
			expectFilter: "Jane Doe Audiobooks",
		},
		{
			name: "author filter with notification preserves notification",
			entry: &parsedEntry{
				Title:    "The Book",
				Category: "Ebooks - Fantasy",
				Authors:  []string{"John Smith"},
				Series:   []string{"Fantasy Series"},
			},
			feedConfig: config.FeedWatcherConfigFeed{
				Name: "TestFeed",
				URL:  "https://example.com/feed",
				AuthorFilters: []config.FeedWatcherConfigFilterByAuthor{
					{
						Author:       "John Smith",
						Notification: "https://webhook.site/test",
					},
				},
			},
			expectMatch:  true,
			expectFilter: "John Smith Books",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a FeedWatcherConfig and preprocess it
			config := config.FeedWatcherConfig{
				Feeds: []config.FeedWatcherConfigFeed{tt.feedConfig},
			}
			config.Preprocess()

			// Now test matching against the preprocessed feed config
			matched, filter, err := tt.entry.HasMatch(ctx, &config.Feeds[0])
			require.NoError(t, err)
			assert.Equal(t, tt.expectMatch, matched, "Match result mismatch")

			if tt.expectMatch {
				require.NotNil(t, filter, "Expected filter to be non-nil")
				assert.Equal(t, tt.expectFilter, filter.Name, "Filter name mismatch")

				// Verify notification is preserved if it was set
				if len(tt.feedConfig.AuthorFilters) > 0 {
					for _, af := range tt.feedConfig.AuthorFilters {
						if af.Notification != "" && filter.Name == fmt.Sprintf("%s Books", af.Author) || filter.Name == fmt.Sprintf("%s Audiobooks", af.Author) {
							assert.Equal(t, af.Notification, filter.Notification, "Notification should be preserved")
						}
					}
				}
			} else {
				assert.Nil(t, filter, "Expected filter to be nil")
			}
		})
	}
}

func TestParsedEntry_HasMatch_NormalFilterBugWithAuthorFilters(t *testing.T) {
	ctx := context.Background()

	// This test specifically checks if normal filters still work when author filters are present
	entry := &parsedEntry{
		Title:       "The Special Book",
		Category:    "Ebooks - Mystery",
		Authors:     []string{"Random Author"},
		Series:      []string{"Detective Series"},
		Description: "A thrilling mystery",
	}

	feedConfig := config.FeedWatcherConfigFeed{
		Name: "TestFeed",
		URL:  "https://example.com/feed",
		Filters: []config.FeedWatcherConfigFeedFilter{
			{
				Name:     "Mystery Filter",
				Category: "mystery-category",
				Matches: []config.FeedWatcherConfigFeedFilterMatch{
					{
						Key:      config.FilterKey_Description,
						Operator: config.FilterOperator_Contains,
						Value:    "mystery",
					},
				},
			},
		},
		AuthorFilters: []config.FeedWatcherConfigFilterByAuthor{
			{
				Author: "John Smith",
			},
		},
	}

	// Preprocess the config
	config := config.FeedWatcherConfig{
		Feeds: []config.FeedWatcherConfigFeed{feedConfig},
	}
	config.Preprocess()

	// The normal filter should still match
	matched, filter, err := entry.HasMatch(ctx, &config.Feeds[0])
	require.NoError(t, err)
	assert.True(t, matched, "Normal filter should match even when author filters are present")
	require.NotNil(t, filter)
	assert.Equal(t, "Mystery Filter", filter.Name, "Should match the normal filter, not an author filter")
}

// Benchmark tests
func TestApplyFilterOperator_CaseInsensitive(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name         string
		operator     config.FilterOperator
		actualValues []string
		filterValue  string
		expected     bool
	}{
		// Equals operator - case insensitive tests
		{
			name:         "equals - lowercase matches uppercase",
			operator:     config.FilterOperator_Equals,
			actualValues: []string{"TEST", "VALUE"},
			filterValue:  "test",
			expected:     true,
		},
		{
			name:         "equals - uppercase matches lowercase",
			operator:     config.FilterOperator_Equals,
			actualValues: []string{"test", "value"},
			filterValue:  "TEST",
			expected:     true,
		},
		{
			name:         "equals - mixed case matches",
			operator:     config.FilterOperator_Equals,
			actualValues: []string{"TeSt", "VaLuE"},
			filterValue:  "tEsT",
			expected:     true,
		},
		{
			name:         "equals - case insensitive no match",
			operator:     config.FilterOperator_Equals,
			actualValues: []string{"TEST", "VALUE"},
			filterValue:  "nomatch",
			expected:     false,
		},

		// Contains operator - case insensitive tests
		{
			name:         "contains - lowercase matches uppercase",
			operator:     config.FilterOperator_Contains,
			actualValues: []string{"THIS IS A TEST"},
			filterValue:  "test",
			expected:     true,
		},
		{
			name:         "contains - uppercase matches lowercase",
			operator:     config.FilterOperator_Contains,
			actualValues: []string{"this is a test"},
			filterValue:  "TEST",
			expected:     true,
		},
		{
			name:         "contains - mixed case matches",
			operator:     config.FilterOperator_Contains,
			actualValues: []string{"ThIs Is A TeSt"},
			filterValue:  "iS a Te",
			expected:     true,
		},
		{
			name:         "contains - case insensitive no match",
			operator:     config.FilterOperator_Contains,
			actualValues: []string{"THIS IS A VALUE"},
			filterValue:  "test",
			expected:     false,
		},

		// Fnmatch operator - case insensitive tests (already implemented)
		{
			name:         "fnmatch - lowercase pattern matches uppercase",
			operator:     config.FilterOperator_Fnmatch,
			actualValues: []string{"TEST FILE.TXT"},
			filterValue:  "*.txt",
			expected:     true,
		},
		{
			name:         "fnmatch - uppercase pattern matches lowercase",
			operator:     config.FilterOperator_Fnmatch,
			actualValues: []string{"test file.txt"},
			filterValue:  "*.TXT",
			expected:     true,
		},
		{
			name:         "fnmatch - mixed case wildcard match",
			operator:     config.FilterOperator_Fnmatch,
			actualValues: []string{"TeSt FiLe.TxT"},
			filterValue:  "test *.txt",
			expected:     true,
		},
		{
			name:         "fnmatch - case insensitive no match",
			operator:     config.FilterOperator_Fnmatch,
			actualValues: []string{"TEST FILE.PDF"},
			filterValue:  "*.txt",
			expected:     false,
		},

		// Regex operator - case insensitive tests
		{
			name:         "regex - lowercase pattern matches uppercase",
			operator:     config.FilterOperator_Regex,
			actualValues: []string{"TEST123"},
			filterValue:  `test\d+`,
			expected:     true,
		},
		{
			name:         "regex - uppercase pattern matches lowercase",
			operator:     config.FilterOperator_Regex,
			actualValues: []string{"test123"},
			filterValue:  `TEST\d+`,
			expected:     true,
		},
		{
			name:         "regex - mixed case pattern matches",
			operator:     config.FilterOperator_Regex,
			actualValues: []string{"TeSt123"},
			filterValue:  `test\d+`,
			expected:     true,
		},
		{
			name:         "regex - case insensitive word boundary",
			operator:     config.FilterOperator_Regex,
			actualValues: []string{"The AUTHOR is John Doe"},
			filterValue:  `\bauthor\b`,
			expected:     true,
		},
		{
			name:         "regex - case insensitive character class",
			operator:     config.FilterOperator_Regex,
			actualValues: []string{"FANTASY"},
			filterValue:  `[a-z]+`,
			expected:     true,
		},
		{
			name:         "regex - case insensitive no match",
			operator:     config.FilterOperator_Regex,
			actualValues: []string{"TESTABC"},
			filterValue:  `test\d+`,
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
