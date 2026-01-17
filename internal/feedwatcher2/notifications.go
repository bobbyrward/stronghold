package feedwatcher2

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/carlmjohnson/requests"

	"github.com/bobbyrward/stronghold/internal/models"
	"github.com/bobbyrward/stronghold/internal/notifications"
)

// SendNotificationViaNotifier sends a Discord notification using the database Notifier's URL.
// Returns nil if notifier is nil (no notification configured).
func SendNotificationViaNotifier(ctx context.Context, notifier *models.Notifier, message notifications.DiscordWebhookMessage) error {
	if notifier == nil {
		slog.DebugContext(ctx, "No notifier configured, skipping notification")
		return nil
	}

	slog.InfoContext(ctx, "Sending notification via database notifier",
		slog.String("notifier_name", notifier.Name),
		slog.String("url", notifier.URL))

	err := requests.
		URL(notifier.URL).
		BodyJSON(&message).
		Fetch(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to send notification",
			slog.String("notifier_name", notifier.Name),
			slog.Any("error", err))
		return err
	}

	slog.InfoContext(ctx, "Notification sent successfully",
		slog.String("notifier_name", notifier.Name))
	return nil
}

// CreateFeedwatcher2NotificationPayload creates a Discord webhook message for a feed match.
func CreateFeedwatcher2NotificationPayload(entry *ParsedEntry, author *models.Author, subscription *models.AuthorSubscription) notifications.DiscordWebhookMessage {
	var embed notifications.DiscordEmbed

	embed.Author.Name = "Feedwatcher2"
	embed.Url = entry.Link
	embed.Description = "Book Grabbed"
	embed.Title = entry.Title
	embed.Color = 16761392
	embed.Timestamp = time.Now().UTC().Format(time.RFC3339)

	addField := func(name string, value string, inline bool) {
		if value == "" {
			return
		}
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
	addField("Subscribed Author", author.Name, true)
	addField("Subscription Scope", subscription.Scope.Name, true)

	// Truncate description if too long for Discord
	description := entry.Description
	if len(description) > 1000 {
		description = description[:997] + "..."
	}
	addField("Description", description, false)

	return notifications.DiscordWebhookMessage{
		Username: "Stronghold",
		Content:  "",
		Embeds:   []notifications.DiscordEmbed{embed},
	}
}
