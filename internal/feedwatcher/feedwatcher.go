package feedwatcher

import (
	"context"
	"errors"
	"fmt"

	"github.com/mmcdole/gofeed"
)

type FeedWatcher struct{}

type parsedEntry struct {
	Title    string
	Category string
	Series   string
	Author   string
	Narrator string
	Summary  string

	/*
	   for entry in parsed_feed["entries"]:
	       subtags = parse_summary(entry["summary"])

	       attributes = {}
	       attributes["title"] = entry["title"]
	       attributes["category"] = get_category_from_tags(entry["tags"])
	       attributes["series"] = subtags.get("series", "")
	       attributes["author"] = subtags.get("author", "")
	       attributes["narrator"] = subtags.get("narrator", "")
	       attributes["summary"] = entry["summary"]

	       for filter_name, filter in feed.filters.items():
	           if filter.has_matches(attributes):
	               if not state.has_download(entry["guid"]):
	                   print(f"Downloading {entry['title']}")

	                   state.add_download(
	                       feed_name,
	                       filter_name,
	                       entry["guid"],
	                       entry["link"],
	                   )

	                   download_torrent(entry["link"], filter.download_location)
	                   download_count += 1

	*/
}

func NewFeedWatcher() *FeedWatcher {
	return &FeedWatcher{}
}

func (fw *FeedWatcher) Run(ctx context.Context) error {
	parser := gofeed.NewParser()

	feed, err := parser.ParseURL("")
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to parse feed"))
	}

	for _, item := range feed.Items {
		fmt.Printf("Title: %s\n", item.Title)
		fmt.Printf("Description: %s\n", item.Description)
		fmt.Printf("Content: %s\n", item.Content)
		fmt.Printf("Link: %s\n", item.Link)

		for key, value := range item.Custom {
			fmt.Printf("Key=%s Value=%s\n", key, value)
		}
	}

	return nil
}

func parseDescription(description string) (parsedEntry, error) {
	parsed := parsedEntry{}

	return parsed, nil
}
