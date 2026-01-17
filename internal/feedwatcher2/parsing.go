package feedwatcher2

import (
	"context"
	"log/slog"
	"strconv"
	"strings"
)

// ParsedEntry represents a parsed feed item with extracted fields.
type ParsedEntry struct {
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

// parseDescription parses the HTML description from a feed item into a ParsedEntry.
func parseDescription(ctx context.Context, description string) (ParsedEntry, error) {
	parsed := ParsedEntry{}

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
			// Trim whitespace from each author
			for i, author := range parsed.Authors {
				parsed.Authors[i] = strings.TrimSpace(author)
			}
		case "Narrator(s)":
			parsed.Narrators = strings.Split(value, ",")
			for i, narrator := range parsed.Narrators {
				parsed.Narrators[i] = strings.TrimSpace(narrator)
			}
		case "Series":
			parsed.Series = strings.Split(value, ",")
			for i, series := range parsed.Series {
				parsed.Series[i] = strings.TrimSpace(series)
			}
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
