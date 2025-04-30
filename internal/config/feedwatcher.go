package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

/*
feedWatcher:
  url: asdasdasdasdasd
  filters:
    - category: personal-books
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
*/

type FilterOperator int

const (
	FilterOperator_Equals FilterOperator = iota
	FilterOperator_Contains
	FilterOperator_Fnmatch
)

var opratorToString = map[FilterOperator]string{
	FilterOperator_Equals:   "equals",
	FilterOperator_Contains: "contains",
	FilterOperator_Fnmatch:  "fnmatch",
}

var stringToOperator = map[string]FilterOperator{
	"equals":   FilterOperator_Equals,
	"contains": FilterOperator_Contains,
	"fnmatch":  FilterOperator_Fnmatch,
}

type FilterKey int

const (
	FilterKey_Author FilterKey = iota
	FilterKey_Series
	FilterKey_Title
	FilterKey_Category
	FilterKey_Summary
)

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
	Name    string                  `yaml:"name"`
	URL     string                  `yaml:"url"`
	Filters []FeedWatcherConfigFeed `yaml:"filters"`
}

type FeedWatcherConfigFeedFilters struct {
	Category string `yaml:"category"`
	Match    string `yaml:"match"`
}

type FeedWatcherConfigFeedFilterMatch struct {
	Key      FilterKey      `yaml:"key"`
	Value    string         `yaml:"value"`
	Operator FilterOperator `yaml:"operator"`
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
	strValue, ok := opratorToString[*op]
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
