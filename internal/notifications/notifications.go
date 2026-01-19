package notifications

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/bobbyrward/stronghold/internal/config"
	"github.com/carlmjohnson/requests"
)

/*

{
  "name": "test webhook",
  "type": 1,
  "channel_id": "199737254929760256",
  "token": "3d89bb7572e0fb30d8128367b3b1b44fecd1726de135cbe28a41f8b2f777c372ba2939e72279b94526ff5d1bd4358d65cf11",
  "avatar": null,
  "guild_id": "199737254929760256",
  "id": "223704706495545344",
  "application_id": null,
  "user": {
    "username": "test",
    "discriminator": "7479",
    "id": "190320984123768832",
    "avatar": "b004ec1740a63ca06ae2e14c5cee11f3",
    "public_flags": 131328
  }
}



*/

type DiscordWebhookMessage struct {
	Username string         `json:"username,omitempty"`
	Content  string         `json:"content,omitempty"`
	Embeds   []DiscordEmbed `json:"embeds,omitempty"`
}

type DiscordEmbed struct {
	Author      DiscordEmbedAuthor  `json:"author,omitempty"`
	Url         string              `json:"url,omitempty"`
	Description string              `json:"description,omitempty"`
	Title       string              `json:"title,omitempty"`
	Color       int                 `json:"color,omitempty"`     // color code of the embed
	Timestamp   string              `json:"timestamp,omitempty"` // ISO8601 timestamp
	Fields      []DiscordEmbedField `json:"fields,omitempty"`    // max of 25 fields
}

type DiscordEmbedAuthor struct {
	Name    string `json:"name,omitempty"`
	Url     string `json:"url,omitempty"`
	IconUrl string `json:"icon_url,omitempty"`
}

type DiscordEmbedImage struct {
	Url    string `json:"url"`              // URL of the image
	Height int    `json:"height,omitempty"` // height of the image in pixels
	Width  int    `json:"width,omitempty"`  // width of the image in pixels
}

type DiscordEmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"` // whether the field should be displayed inline
}

func findNotifier(name string) (bool, *config.NotificationsConfigNotifier) {
	for _, notifier := range config.Config.Notifications.Notifiers {
		if notifier.Name == name {
			return true, &notifier
		}
	}

	return false, nil
}

func SendNotification(ctx context.Context, notificationName string, message any) error {
	if notificationName == "" {
		return nil
	}

	found, notifier := findNotifier(notificationName)
	if !found {
		slog.ErrorContext(ctx, "Notifier not found", slog.String("name", notificationName))
		return nil
	}

	switch notifier.Type {
	case "discord":
		discordMessage, ok := message.(DiscordWebhookMessage)

		if !ok {
			slog.ErrorContext(ctx, "Invalid message type for Discord notifier", slog.String("name", notificationName))
			return nil
		}

		return SendDiscordNotification(ctx, notificationName, discordMessage)
	default:
		slog.ErrorContext(ctx, "Unknown notifier type", slog.String("name", notificationName), slog.String("type", notifier.Type))

	}

	return nil
}

func SendDiscordNotification(ctx context.Context, notificationName string, message DiscordWebhookMessage) error {
	found, notifier := findNotifier(notificationName)
	if !found {
		slog.ErrorContext(ctx, "Notifier not found", slog.String("name", notificationName))
		return nil
	}

	bytes, err := json.Marshal(message)
	if err == nil {
		slog.InfoContext(ctx, "Message", slog.String("message", string(bytes)))
	}

	var responseString string

	err = requests.
		URL(notifier.Url).
		BodyJSON(&message).
		ToString(&responseString).
		Fetch(ctx)
	if err != nil {
		slog.ErrorContext(
			ctx,
			"Failed to send Discord notification",
			slog.String("name", notificationName),
			slog.Any("error", err),
			slog.String("response", responseString),
		)
	}

	return nil
}
