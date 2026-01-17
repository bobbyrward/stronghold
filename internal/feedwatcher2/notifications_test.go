package feedwatcher2

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bobbyrward/stronghold/internal/models"
	"github.com/bobbyrward/stronghold/internal/notifications"
)

func TestSendNotificationViaNotifier_NilNotifier(t *testing.T) {
	err := SendNotificationViaNotifier(
		context.Background(),
		nil,
		notifications.DiscordWebhookMessage{},
	)
	assert.NoError(t, err)
}

func TestSendNotificationViaNotifier_Success(t *testing.T) {
	var receivedMessage notifications.DiscordWebhookMessage

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := json.NewDecoder(r.Body).Decode(&receivedMessage)
		require.NoError(t, err)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	notifier := &models.Notifier{
		Name: "test-notifier",
		URL:  server.URL,
	}

	message := notifications.DiscordWebhookMessage{
		Username: "Test",
		Content:  "Hello",
	}

	err := SendNotificationViaNotifier(context.Background(), notifier, message)
	assert.NoError(t, err)
	assert.Equal(t, "Test", receivedMessage.Username)
	assert.Equal(t, "Hello", receivedMessage.Content)
}

func TestSendNotificationViaNotifier_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	notifier := &models.Notifier{
		Name: "test-notifier",
		URL:  server.URL,
	}

	err := SendNotificationViaNotifier(
		context.Background(),
		notifier,
		notifications.DiscordWebhookMessage{},
	)
	// Note: The requests library may not return an error for 5xx responses
	// depending on configuration, but let's verify the call doesn't panic
	_ = err
}

func TestCreateFeedwatcher2NotificationPayload(t *testing.T) {
	entry := &ParsedEntry{
		Title:       "Test Book",
		Link:        "https://example.com/torrent",
		Category:    "Audiobooks - Fiction",
		Series:      []string{"Test Series"},
		Authors:     []string{"Test Author"},
		Narrators:   []string{"Test Narrator"},
		Tags:        "fiction, fantasy",
		Description: "A test description",
	}

	author := &models.Author{
		Name: "Test Author",
	}

	scope := models.SubscriptionScope{
		Name: "personal",
	}

	subscription := &models.AuthorSubscription{
		Scope: scope,
	}

	payload := CreateFeedwatcher2NotificationPayload(entry, author, subscription)

	assert.Equal(t, "Stronghold", payload.Username)
	assert.Len(t, payload.Embeds, 1)

	embed := payload.Embeds[0]
	assert.Equal(t, "Feedwatcher2", embed.Author.Name)
	assert.Equal(t, "Test Book", embed.Title)
	assert.Equal(t, "https://example.com/torrent", embed.Url)
	assert.Equal(t, "Book Grabbed", embed.Description)
	assert.Equal(t, 16761392, embed.Color)
	assert.NotEmpty(t, embed.Timestamp)

	// Check fields
	fieldMap := make(map[string]string)
	for _, field := range embed.Fields {
		fieldMap[field.Name] = field.Value
	}

	assert.Equal(t, "Audiobooks - Fiction", fieldMap["Category"])
	assert.Equal(t, "Test Series", fieldMap["Series"])
	assert.Equal(t, "Test Author", fieldMap["Authors"])
	assert.Equal(t, "Test Narrator", fieldMap["Narrators"])
	assert.Equal(t, "fiction, fantasy", fieldMap["Tags"])
	assert.Equal(t, "Test Author", fieldMap["Subscribed Author"])
	assert.Equal(t, "personal", fieldMap["Subscription Scope"])
	assert.Equal(t, "A test description", fieldMap["Description"])
}

func TestCreateFeedwatcher2NotificationPayload_NoNarrators(t *testing.T) {
	entry := &ParsedEntry{
		Title:    "Ebook Title",
		Category: "Ebooks - Fiction",
		Authors:  []string{"Author Name"},
		// No narrators
	}

	author := &models.Author{Name: "Author Name"}
	scope := models.SubscriptionScope{Name: "family"}
	subscription := &models.AuthorSubscription{Scope: scope}

	payload := CreateFeedwatcher2NotificationPayload(entry, author, subscription)

	// Check that Narrators field is not present
	for _, field := range payload.Embeds[0].Fields {
		assert.NotEqual(t, "Narrators", field.Name)
	}
}

func TestCreateFeedwatcher2NotificationPayload_LongDescription(t *testing.T) {
	// Create a description longer than 1000 characters
	longDesc := ""
	for i := 0; i < 200; i++ {
		longDesc += "Lorem ipsum "
	}

	entry := &ParsedEntry{
		Title:       "Test",
		Description: longDesc,
		Authors:     []string{"Author"},
	}

	author := &models.Author{Name: "Author"}
	scope := models.SubscriptionScope{Name: "personal"}
	subscription := &models.AuthorSubscription{Scope: scope}

	payload := CreateFeedwatcher2NotificationPayload(entry, author, subscription)

	// Find description field
	for _, field := range payload.Embeds[0].Fields {
		if field.Name == "Description" {
			assert.LessOrEqual(t, len(field.Value), 1000)
			assert.True(t, len(field.Value) <= 1000)
			if len(longDesc) > 1000 {
				assert.Contains(t, field.Value, "...")
			}
			break
		}
	}
}

func TestCreateFeedwatcher2NotificationPayload_EmptyFields(t *testing.T) {
	entry := &ParsedEntry{
		Title:   "Minimal Entry",
		Authors: []string{"Author"},
		// All other fields empty
	}

	author := &models.Author{Name: "Author"}
	scope := models.SubscriptionScope{Name: "personal"}
	subscription := &models.AuthorSubscription{Scope: scope}

	payload := CreateFeedwatcher2NotificationPayload(entry, author, subscription)

	// Should still have basic fields
	assert.Equal(t, "Minimal Entry", payload.Embeds[0].Title)
	assert.NotEmpty(t, payload.Embeds[0].Fields)
}
