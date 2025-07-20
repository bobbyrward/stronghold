package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

/*
feedWatcher:
  feeds:
    - name: MAM
      url: https://04k6i.mrd.ninja/rss/52cf3a54
      filters:
        - name: Blaise Corvin Books
          category: personal-books
          match:
            - key: author
              value: Blaise Corvin
              operator: equals
*/

type FilterOperator int

const (
	FilterOperator_Equals FilterOperator = iota
	FilterOperator_Contains
	FilterOperator_Fnmatch
	FilterOperator_Regex
)

func (fo FilterOperator) String() string {
	str, ok := operatorToString[fo]
	if !ok {
		// TODO: Log this somehow
		return ""
	}

	return str
}

var operatorToString = map[FilterOperator]string{
	FilterOperator_Equals:   "equals",
	FilterOperator_Contains: "contains",
	FilterOperator_Fnmatch:  "fnmatch",
	FilterOperator_Regex:    "regex",
}

var stringToOperator = map[string]FilterOperator{
	"equals":   FilterOperator_Equals,
	"contains": FilterOperator_Contains,
	"fnmatch":  FilterOperator_Fnmatch,
	"regex":    FilterOperator_Regex,
}

type FilterKey int

const (
	FilterKey_Author FilterKey = iota
	FilterKey_Series
	FilterKey_Title
	FilterKey_Category
	FilterKey_Summary
)

func (fk FilterKey) String() string {
	str, ok := matchKeyToString[fk]
	if !ok {
		// TODO: Log this somehow
		return ""
	}

	return str
}

var matchKeyToString = map[FilterKey]string{
	FilterKey_Author:   "author",
	FilterKey_Series:   "series",
	FilterKey_Title:    "title",
	FilterKey_Category: "category",
	FilterKey_Summary:  "summary",
}

var stringToMatchKey = map[string]FilterKey{
	"author":   FilterKey_Author,
	"series":   FilterKey_Series,
	"title":    FilterKey_Title,
	"category": FilterKey_Category,
	"summary":  FilterKey_Summary,
}

type FeedWatcherConfig struct {
	Feeds []FeedWatcherConfigFeed `yaml:"feeds"`
}

type FeedWatcherConfigFeed struct {
	Name    string                        `yaml:"name"`
	URL     string                        `yaml:"url"`
	Filters []FeedWatcherConfigFeedFilter `yaml:"filters"`
}

type FeedWatcherConfigFeedFilter struct {
	Name         string                             `yaml:"name"`
	Category     string                             `yaml:"category"`
	Notification string                             `yaml:"notification,omitempty"`
	Matches      []FeedWatcherConfigFeedFilterMatch `yaml:"match"`
}

type FeedWatcherConfigFeedFilterMatch struct {
	Key      FilterKey      `yaml:"key"`
	Value    string         `yaml:"value"`
	Operator FilterOperator `yaml:"operator"`
}

func (fwcffm FeedWatcherConfigFeedFilterMatch) String() string {
	return fmt.Sprintf("FilterMatch{Key: %s, Value: %s, Operator: %s}", fwcffm.Key.String(), fwcffm.Value, fwcffm.Operator.String())
}

func (op *FilterOperator) UnmarshalYAML(value *yaml.Node) error {
	var strValue string

	err := value.Decode(&strValue)
	if err != nil {
		return err
	}

	enumValue, ok := stringToOperator[strValue]
	if !ok {
		return fmt.Errorf("invalid operator: %s", strValue)
	}

	*op = enumValue

	return nil
}

func (op *FilterOperator) MarshalYAML() (interface{}, error) {
	strValue, ok := operatorToString[*op]
	if !ok {
		return nil, fmt.Errorf("invalid operator: %d", *op)
	}

	return strValue, nil
}

func (key *FilterKey) UnmarshalYAML(value *yaml.Node) error {
	var strValue string

	err := value.Decode(&strValue)
	if err != nil {
		return err
	}

	enumValue, ok := stringToMatchKey[strValue]
	if !ok {
		return fmt.Errorf("invalid key: %s", strValue)
	}

	*key = enumValue

	return nil
}

func (key *FilterKey) MarshalYAML() (interface{}, error) {
	strValue, ok := matchKeyToString[*key]
	if !ok {
		return nil, fmt.Errorf("invalid key: %d", *key)
	}

	return strValue, nil
}
