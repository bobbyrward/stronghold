package eventlog

import (
	"context"
	"log/slog"
	"time"

	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
)

// Cleanup deletes event log entries older than retentionDays.
// It is fire-and-forget: errors are logged but never returned.
func Cleanup(ctx context.Context, db *gorm.DB, retentionDays int) {
	cutoff := time.Now().AddDate(0, 0, -retentionDays)

	result := db.Where("created_at < ?", cutoff).Delete(&models.EventLog{})
	if result.Error != nil {
		slog.ErrorContext(ctx, "Failed to cleanup event logs",
			slog.Int("retention_days", retentionDays),
			slog.Any("error", result.Error))
		return
	}

	if result.RowsAffected > 0 {
		slog.InfoContext(ctx, "Cleaned up old event logs",
			slog.Int64("deleted", result.RowsAffected),
			slog.Int("retention_days", retentionDays))
	}
}
