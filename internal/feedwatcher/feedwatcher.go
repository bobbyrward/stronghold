package feedwatcher

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/bobbyrward/stronghold/internal/config"
	"github.com/mmcdole/gofeed"
)

type FeedWatcher struct{}

type parsedEntry struct {
	Guid        string
	Link        string
	Title       string
	Category    string
	Series      string
	Authors     []string
	Narrators   []string
	Summary     string
	Leechers    int
	Seeders     int
	Added       string
	Tags        string
	Description string
	/*
		a   for entry in parsed_feed["entries"]:
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
	for _, feed := range config.Config.FeedWatcher.Feeds {
		err := fw.watchFeed(ctx, &feed)
		if err != nil {
			slog.WarnContext(ctx, "Unable to watch feed", slog.String("name", feed.Name))
		}
	}

	return nil
}

func (fw *FeedWatcher) watchFeed(ctx context.Context, feedConfig *config.FeedWatcherConfigFeed) error {
	parser := gofeed.NewParser()

	feed, err := parser.ParseURL(feedConfig.URL)
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to parse feed"))
	}

	for _, item := range feed.Items {
		for key, value := range item.Custom {
			fmt.Printf("Key=%s Value=%s\n", key, value)
		}

		entry, err := parseDescription(ctx, item.Description)
		if err != nil {
			slog.WarnContext(ctx, "Unable to parse feed item", slog.String("feedName", feedConfig.Name), slog.Any("err", err))
		}

		entry.Guid = item.GUID
		entry.Link = item.Link
		entry.Title = item.Title

		slog.InfoContext(ctx,
			"Parsed entry",
			slog.String("GUID", entry.Guid),
			slog.String("Link", entry.Link),
			slog.String("Title", entry.Title),
			slog.Any("Authors", entry.Authors),
			slog.Any("Narrators", entry.Narrators),
			slog.Int("Leechers", entry.Leechers),
			slog.Int("Seeders", entry.Seeders),
			slog.String("Category", entry.Category),
			slog.String("Series", entry.Series),
			slog.String("Summary", entry.Summary),
			slog.String("Added", entry.Added),
			slog.String("Tags", entry.Tags),
			slog.String("Description", entry.Description),
		)
	}

	return nil
}

func parseDescription(ctx context.Context, description string) (parsedEntry, error) {
	parsed := parsedEntry{}

	for _, part := range strings.Split(description, "<br/>") {
		if strings.TrimSpace(part) == "" {
			continue
		}

		label, value, ok := strings.Cut(part, ":")

		if !ok {
			fmt.Printf("!ok: %s\n", part)
			continue
		}

		label = strings.TrimSpace(label)
		value = strings.TrimSpace(value)

		switch label {
		case "Author(s)":
			parsed.Authors = strings.Split(value, ",")
		case "Narrator(s)":
			parsed.Narrators = strings.Split(value, ",")
		case "Series":
			parsed.Series = value
		case "Sumnmary":
			parsed.Summary = value
		case "Category":
			parsed.Category = value
		case "Leechers":
			intValue, err := strconv.Atoi(value)
			if err != nil {
				slog.WarnContext(ctx, "Unable to parse leechers", slog.String("leechers", value))
			}
			parsed.Leechers = intValue
		case "Seeders":
			intValue, err := strconv.Atoi(value)
			if err != nil {
				slog.WarnContext(ctx, "Unable to parse seeeders", slog.String("seeders", value))
			}
			parsed.Seeders = intValue
		case "Added":
			parsed.Added = value
		case "Tags":
			parsed.Tags = value
		case "Description":
			parsed.Description = value

		}
	}

	return parsed, nil
}
