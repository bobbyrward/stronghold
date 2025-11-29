package feedwatcher

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/danwakefield/fnmatch"
	"github.com/mmcdole/gofeed"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/config"
	"github.com/bobbyrward/stronghold/internal/models"
	"github.com/bobbyrward/stronghold/internal/notifications"
	"github.com/bobbyrward/stronghold/internal/qbit"
)

type FeedWatcher struct{}

type parsedEntry struct {
	// Guid is the unique identifier for the entry.
	Guid string
	// Link is the URL associated with the entry.
	Link string
	// Title is the name or title of the entry.
	Title string
	// Category represents the category or genre of the entry.
	Category string
	// Series contains a list of series names associated with the entry.
	Series []string
	// Authors contains a list of authors associated with the entry.
	Authors []string
	// Narrators contains a list of narrators associated with the entry.
	Narrators []string
	// Summary provides a brief summary or description of the entry.
	Summary string
	// Leechers is the number of leechers for the entry.
	Leechers int
	// Seeders is the number of seeders for the entry.
	Seeders int
	// Added is the date and time when the entry was added.
	Added string
	// Tags contains a comma-separated list of tags associated with the entry.
	Tags string
	// Description provides a detailed description of the entry.
	Description string
}

func (pe *parsedEntry) GetKeyValue(ctx context.Context, key config.FilterKey) []string {
	switch key {
	case config.FilterKey_Author:
		return pe.Authors

	case config.FilterKey_Series:
		return pe.Series

	case config.FilterKey_Title:
		return []string{pe.Title}

	case config.FilterKey_Category:
		return []string{pe.Category}

	case config.FilterKey_Summary:
		return []string{pe.Summary}

	case config.FilterKey_Tags:
		return []string{pe.Tags}

	case config.FilterKey_Description:
		return []string{pe.Description}

	default:
		slog.WarnContext(ctx, "Unknown filter key", slog.String("key", key.String()))

	}

	return []string{}
}

func applyFilterOperator(ctx context.Context, operator config.FilterOperator, actualValues []string, filterValue string) bool {
	filterValue = strings.ToLower(filterValue)

	for _, actualValue := range actualValues {
		actualValue = strings.ToLower(actualValue)

		switch operator {
		case config.FilterOperator_Equals:
			if actualValue == filterValue {
				return true
			}

		case config.FilterOperator_Contains:
			if strings.Contains(actualValue, filterValue) {
				return true
			}

		case config.FilterOperator_Fnmatch:
			if fnmatch.Match(filterValue, actualValue, fnmatch.FNM_IGNORECASE) {
				return true
			}

		case config.FilterOperator_Regex:
			matched, err := regexp.Match(filterValue, []byte(actualValue))
			if err != nil {
				slog.WarnContext(ctx, "Invalid regex in filter", slog.String("regex", filterValue), slog.Any("err", err))
				continue
			}

			if matched {
				return true
			}
		}
	}

	return false
}

func (pe *parsedEntry) hasAllMatches(ctx context.Context, filter *config.FeedWatcherConfigFeedFilter) bool {
	for _, match := range filter.Matches {
		actualValues := pe.GetKeyValue(ctx, match.Key)

		if !applyFilterOperator(ctx, match.Operator, actualValues, match.Value) {
			return false
		}
	}

	return true
}

func (pe *parsedEntry) HasMatch(ctx context.Context, feedConfig *config.FeedWatcherConfigFeed) (bool, *config.FeedWatcherConfigFeedFilter, error) {
	for _, filter := range feedConfig.Filters {
		if pe.hasAllMatches(ctx, &filter) {
			return true, &filter, nil
		}
	}

	return false, nil, nil
}

func CreateFeedwatcherNotificationPayload(entry *parsedEntry, filter *config.FeedWatcherConfigFeedFilter) notifications.DiscordWebhookMessage {
	var embed notifications.DiscordEmbed

	embed.Author.Name = "Feedwatcher"
	embed.Url = entry.Link
	embed.Description = "Book Grabbed"
	embed.Title = entry.Title
	embed.Color = 16761392
	embed.Timestamp = time.Now().UTC().Format(time.RFC3339)

	addField := func(name string, value string, inline bool) {
		embed.Fields = append(embed.Fields, notifications.DiscordEmbedField{
			Name:   name,
			Value:  value,
			Inline: inline,
		})
	}

	addField("Category", entry.Category, false)
	addField("Series", strings.Join(entry.Series, ", "), false)
	addField("Authors", strings.Join(entry.Authors, ", "), false)

	if len(entry.Narrators) > 0 {
		addField("Narrators", strings.Join(entry.Narrators, ", "), false)
	}

	addField("Tags", entry.Tags, false)
	addField("Filter", filter.Name, false)
	addField("Description", entry.Description, false)

	req := notifications.DiscordWebhookMessage{
		Username: "Stronghold",
		Content:  "",
		Embeds:   []notifications.DiscordEmbed{embed},
	}

	return req
}

func NewFeedWatcher() *FeedWatcher {
	return &FeedWatcher{}
}

func (fw *FeedWatcher) Run(ctx context.Context, db *gorm.DB) error {
	client, err := qbit.CreateClient()
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to create qBittorrent client"))
	}

	config.Config.FeedWatcher.Preprocess()

	for _, feed := range config.Config.FeedWatcher.Feeds {
		slog.DebugContext(ctx, "Watching feed", slog.String("name", feed.Name), slog.String("url", feed.URL))

		for _, filter := range feed.Filters {
			slog.DebugContext(ctx, "Filter", slog.String("name", filter.Name), slog.String("category", filter.Category), slog.String("notification", filter.Notification), slog.Any("matches", filter.Matches))
		}

		err := fw.watchFeed(ctx, &feed, client, db)
		if err != nil {
			slog.WarnContext(ctx, "Unable to watch feed", slog.String("name", feed.Name))
		}
	}

	return nil
}

func (fw *FeedWatcher) watchFeed(ctx context.Context, feedConfig *config.FeedWatcherConfigFeed, qbitClient qbit.QbitClient, db *gorm.DB) error {
	parser := gofeed.NewParser()

	feed, err := parser.ParseURL(feedConfig.URL)
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to parse feed"))
	}

	for _, item := range feed.Items {
		entry, err := parseDescription(ctx, item.Description)
		if err != nil {
			slog.WarnContext(ctx, "Unable to parse feed item", slog.String("feedName", feedConfig.Name), slog.Any("err", err))
		}

		entry.Guid = item.GUID
		entry.Link = item.Link
		entry.Title = item.Title

		matched, filter, err := entry.HasMatch(ctx, feedConfig)
		if err != nil {
			slog.ErrorContext(ctx, "Error checking for filter match", slog.String("feedName", feedConfig.Name), slog.Any("err", err))
			continue
		}
		if matched {
			var found models.FeedItem

			result := db.Where(&models.FeedItem{Guid: entry.Guid}, "Guid").First(&found)
			if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
				slog.InfoContext(
					ctx,
					"Error checking for existing feed item",
					slog.String("guid", entry.Guid),
					slog.String("feedName", feedConfig.Name),
					slog.Any("err", result.Error),
					slog.String("filterName", filter.Name),
				)
				continue
			}

			if result.Error == nil {
				slog.InfoContext(
					ctx,
					"Existing feed item already imported",
					slog.String("guid", entry.Guid),
					slog.String("feedName", feedConfig.Name),
					slog.String("filterName", filter.Name),
					slog.String("Title", item.Title),
				)
				continue
			}

			slog.InfoContext(ctx,
				"Matched Filter",
				slog.String("feedName", feedConfig.Name),
				slog.String("filterName", filter.Name),
				slog.String("GUID", entry.Guid),
				slog.String("Link", entry.Link),
				slog.String("Title", entry.Title),
				slog.Any("Authors", entry.Authors),
				slog.Any("Narrators", entry.Narrators),
				slog.Int("Leechers", entry.Leechers),
				slog.Int("Seeders", entry.Seeders),
				slog.String("Category", entry.Category),
				slog.Any("Series", entry.Series),
				slog.String("Summary", entry.Summary),
				slog.String("Added", entry.Added),
				slog.String("Tags", entry.Tags),
				slog.String("Description", entry.Description),
			)

			err = qbitClient.AddTorrentFromUrlCtx(
				ctx,
				entry.Link,
				map[string]string{
					"autoTMM":  "true",
					"category": filter.Category,
				},
			)
			if err != nil {
				slog.ErrorContext(ctx, "Failed to add torrent to qBittorrent",
					slog.String("link", entry.Link),
					slog.Any("err", err),
					slog.String("guid", entry.Guid),
					slog.Group("feedFilter",
						slog.String("feedName", feedConfig.Name),
						slog.String("filterName", filter.Name),
					),
				)
			}

			pubDate := time.Now()
			if item.PublishedParsed != nil {
				pubDate = *item.PublishedParsed
			}

			newFeedItem := models.FeedItem{
				Guid:        entry.Guid,
				Title:       entry.Title,
				Link:        entry.Link,
				Category:    entry.Category,
				Description: item.Description,
				PubDate:     pubDate,
			}

			result = db.Create(&newFeedItem)
			if result.Error != nil {
				slog.ErrorContext(ctx, "Failed to create feed item in database", slog.String("guid", entry.Guid), slog.String("feedName", feedConfig.Name), slog.Any("err", result.Error))
			}

			err := notifications.SendNotification(ctx, filter.Notification, CreateFeedwatcherNotificationPayload(&entry, filter))
			if err != nil {
				slog.ErrorContext(ctx, "Failed to send notification", slog.String("feedName", feedConfig.Name), slog.String("filterName", filter.Name), slog.Any("err", err))
			}
		}
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
			slog.WarnContext(ctx, "Unable to parse label and value from part", slog.String("part", part))
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
			parsed.Series = strings.Split(value, ",")
		case "Summary":
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
				slog.WarnContext(ctx, "Unable to parse seeders", slog.String("seeders", value))
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
