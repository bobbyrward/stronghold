package eventlog

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bobbyrward/stronghold/internal/models"
)

func TestLog(t *testing.T) {
	db, err := models.ConnectTestDB()
	require.NoError(t, err)

	t.Run("creates record with correct fields", func(t *testing.T) {
		Log(db, CategoryDownload, EventTorrentAdded, SourceFeedwatcher2,
			EntityTorrent, "abc123", "Downloaded torrent: Test Book", nil)

		var entry models.EventLog
		err := db.Last(&entry).Error
		require.NoError(t, err)

		assert.Equal(t, CategoryDownload, entry.Category)
		assert.Equal(t, EventTorrentAdded, entry.EventType)
		assert.Equal(t, SourceFeedwatcher2, entry.Source)
		assert.Equal(t, EntityTorrent, entry.EntityType)
		assert.Equal(t, "abc123", entry.EntityID)
		assert.Equal(t, "Downloaded torrent: Test Book", entry.Summary)
		assert.Empty(t, entry.Details)
		assert.False(t, entry.CreatedAt.IsZero())
	})

	t.Run("marshals details to JSON", func(t *testing.T) {
		details := map[string]any{
			"title":    "Test Book",
			"category": "ebooks",
			"hash":     "def456",
		}

		Log(db, CategoryDownload, EventTorrentAdded, SourceDiscordBot,
			EntityTorrent, "def456", "Downloaded via Discord", details)

		var entry models.EventLog
		err := db.Last(&entry).Error
		require.NoError(t, err)

		assert.Contains(t, entry.Details, `"title":"Test Book"`)
		assert.Contains(t, entry.Details, `"category":"ebooks"`)
		assert.Contains(t, entry.Details, `"hash":"def456"`)
	})

	t.Run("does not panic on unmarshalable details", func(t *testing.T) {
		// Channels cannot be marshaled to JSON
		ch := make(chan int)

		assert.NotPanics(t, func() {
			Log(db, CategoryDownload, EventTorrentAdded, SourceAPI,
				EntityTorrent, "ghi789", "Bad details", ch)
		})

		// Should still create the record with empty details
		var entry models.EventLog
		err := db.Last(&entry).Error
		require.NoError(t, err)
		assert.Equal(t, "{}", entry.Details)
	})
}

func TestCleanup(t *testing.T) {
	db, err := models.ConnectTestDB()
	require.NoError(t, err)

	t.Run("deletes old records and preserves recent ones", func(t *testing.T) {
		// Create an old record (100 days ago)
		oldEntry := models.EventLog{
			CreatedAt:  time.Now().AddDate(0, 0, -100),
			Category:   CategoryDownload,
			EventType:  EventTorrentAdded,
			Source:     SourceFeedwatcher2,
			EntityType: EntityTorrent,
			EntityID:   "old-hash",
			Summary:    "Old download",
		}
		require.NoError(t, db.Create(&oldEntry).Error)

		// Create a recent record (1 day ago)
		recentEntry := models.EventLog{
			CreatedAt:  time.Now().AddDate(0, 0, -1),
			Category:   CategoryDownload,
			EventType:  EventTorrentAdded,
			Source:     SourceFeedwatcher2,
			EntityType: EntityTorrent,
			EntityID:   "recent-hash",
			Summary:    "Recent download",
		}
		require.NoError(t, db.Create(&recentEntry).Error)

		// Run cleanup with 90-day retention
		Cleanup(context.Background(), db, 90)

		// Verify old record was deleted
		var count int64
		db.Model(&models.EventLog{}).Where("entity_id = ?", "old-hash").Count(&count)
		assert.Equal(t, int64(0), count)

		// Verify recent record was preserved
		db.Model(&models.EventLog{}).Where("entity_id = ?", "recent-hash").Count(&count)
		assert.Equal(t, int64(1), count)
	})
}
