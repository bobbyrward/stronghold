package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
)

type FeedFilterRequest struct {
	Name         string `json:"name" validate:"required"`
	FeedName     string `json:"feed_name" validate:"required"`
	CategoryName string `json:"category_name" validate:"required"`
	NotifierName string `json:"notifier_name" validate:"required"`
}

type FeedFilterResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	FeedID   uint   `json:"feed_id"`
	FeedName string `json:"feed_name"`
	Category string `json:"category"`
	Notifier string `json:"notifier"`
}

// ListFeedFilters returns all feed filters
func ListFeedFilters(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		feedIDStr := c.QueryParam("feed_id")

		slog.InfoContext(ctx, "Listing feed filters", slog.String("feed_id_filter", feedIDStr))

		var feedFilters []models.FeedFilter
		query := db.Preload("Feed").Preload("TorrentCategory").Preload("Notifier")

		if feedIDStr != "" {
			feedID, err := strconv.ParseUint(feedIDStr, 10, 32)
			if err != nil {
				slog.WarnContext(ctx, "Invalid feed_id parameter", slog.String("feed_id", feedIDStr), slog.Any("error", err))
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid feed_id parameter"})
			}
			query = query.Where("feed_id = ?", feedID)
		}

		result := query.Find(&feedFilters)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to query feed filters", slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query feed filters"})
		}

		response := make([]FeedFilterResponse, len(feedFilters))
		for i, ff := range feedFilters {
			feedName := ""
			if ff.Feed.Name != "" {
				feedName = ff.Feed.Name
			}
			categoryName := ""
			if ff.TorrentCategory.Name != "" {
				categoryName = ff.TorrentCategory.Name
			}
			notifierName := ""
			if ff.Notifier.Name != "" {
				notifierName = ff.Notifier.Name
			}
			response[i] = FeedFilterResponse{
				ID:       ff.ID,
				Name:     ff.Name,
				FeedID:   ff.FeedID,
				FeedName: feedName,
				Category: categoryName,
				Notifier: notifierName,
			}
		}

		slog.InfoContext(ctx, "Successfully listed feed filters", slog.Int("count", len(response)))
		return c.JSON(http.StatusOK, response)
	}
}

// CreateFeedFilter creates a new feed filter
func CreateFeedFilter(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var req FeedFilterRequest
		if err := c.Bind(&req); err != nil {
			slog.WarnContext(ctx, "Invalid request body", slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		}

		if req.Name == "" {
			slog.WarnContext(ctx, "Validation failed: name is required")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Name is required"})
		}

		if req.FeedName == "" {
			slog.WarnContext(ctx, "Validation failed: feed_name is required")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Feed name is required"})
		}

		if req.CategoryName == "" {
			slog.WarnContext(ctx, "Validation failed: category_name is required")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Category name is required"})
		}

		if req.NotifierName == "" {
			slog.WarnContext(ctx, "Validation failed: notifier_name is required")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Notifier name is required"})
		}

		slog.InfoContext(ctx, "Creating feed filter",
			slog.String("name", req.Name),
			slog.String("feed_name", req.FeedName),
			slog.String("category_name", req.CategoryName),
			slog.String("notifier_name", req.NotifierName))

		// Look up feed by name
		var feed models.Feed
		result := db.Where("name = ?", req.FeedName).First(&feed)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				slog.WarnContext(ctx, "Feed not found", slog.String("feed_name", req.FeedName))
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Feed not found"})
			}
			slog.ErrorContext(ctx, "Failed to query feed", slog.String("feed_name", req.FeedName), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query feed"})
		}

		// Look up torrent category by name
		var torrentCategory models.TorrentCategory
		result = db.Where("name = ?", req.CategoryName).First(&torrentCategory)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				slog.WarnContext(ctx, "Torrent category not found", slog.String("category_name", req.CategoryName))
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Torrent category not found"})
			}
			slog.ErrorContext(ctx, "Failed to query torrent category", slog.String("category_name", req.CategoryName), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query torrent category"})
		}

		// Look up notifier by name
		var notifier models.Notifier
		result = db.Where("name = ?", req.NotifierName).First(&notifier)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				slog.WarnContext(ctx, "Notifier not found", slog.String("notifier_name", req.NotifierName))
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Notifier not found"})
			}
			slog.ErrorContext(ctx, "Failed to query notifier", slog.String("notifier_name", req.NotifierName), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query notifier"})
		}

		feedFilter := models.FeedFilter{
			Name:              req.Name,
			FeedID:            feed.ID,
			TorrentCategoryID: torrentCategory.ID,
			NotifierID:        notifier.ID,
		}

		result = db.Create(&feedFilter)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to create feed filter", slog.String("name", req.Name), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create feed filter"})
		}

		response := FeedFilterResponse{
			ID:       feedFilter.ID,
			Name:     feedFilter.Name,
			FeedID:   feedFilter.FeedID,
			FeedName: req.FeedName,
			Category: req.CategoryName,
			Notifier: req.NotifierName,
		}

		slog.InfoContext(ctx, "Successfully created feed filter", slog.Uint64("id", uint64(feedFilter.ID)), slog.String("name", feedFilter.Name))
		return c.JSON(http.StatusCreated, response)
	}
}

// GetFeedFilter returns a single feed filter by ID
func GetFeedFilter(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		idStr := c.Param("id")

		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			slog.WarnContext(ctx, "Invalid ID parameter", slog.String("id", idStr), slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}

		slog.InfoContext(ctx, "Getting feed filter", slog.Uint64("id", id))

		var feedFilter models.FeedFilter
		result := db.Preload("Feed").Preload("TorrentCategory").Preload("Notifier").First(&feedFilter, id)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				slog.WarnContext(ctx, "Feed filter not found", slog.Uint64("id", id))
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Feed filter not found"})
			}
			slog.ErrorContext(ctx, "Failed to query feed filter", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query feed filter"})
		}

		feedName := ""
		if feedFilter.Feed.Name != "" {
			feedName = feedFilter.Feed.Name
		}
		categoryName := ""
		if feedFilter.TorrentCategory.Name != "" {
			categoryName = feedFilter.TorrentCategory.Name
		}
		notifierName := ""
		if feedFilter.Notifier.Name != "" {
			notifierName = feedFilter.Notifier.Name
		}

		response := FeedFilterResponse{
			ID:       feedFilter.ID,
			Name:     feedFilter.Name,
			FeedID:   feedFilter.FeedID,
			FeedName: feedName,
			Category: categoryName,
			Notifier: notifierName,
		}

		slog.InfoContext(ctx, "Successfully retrieved feed filter", slog.Uint64("id", id), slog.String("name", feedFilter.Name))
		return c.JSON(http.StatusOK, response)
	}
}

// UpdateFeedFilter updates an existing feed filter
func UpdateFeedFilter(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		idStr := c.Param("id")

		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			slog.WarnContext(ctx, "Invalid ID parameter", slog.String("id", idStr), slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}

		var req FeedFilterRequest
		if err := c.Bind(&req); err != nil {
			slog.WarnContext(ctx, "Invalid request body", slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		}

		if req.Name == "" {
			slog.WarnContext(ctx, "Validation failed: name is required")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Name is required"})
		}

		if req.CategoryName == "" {
			slog.WarnContext(ctx, "Validation failed: category_name is required")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Category name is required"})
		}

		if req.NotifierName == "" {
			slog.WarnContext(ctx, "Validation failed: notifier_name is required")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Notifier name is required"})
		}

		slog.InfoContext(ctx, "Updating feed filter", slog.Uint64("id", id), slog.String("name", req.Name))

		var feedFilter models.FeedFilter
		result := db.First(&feedFilter, id)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				slog.WarnContext(ctx, "Feed filter not found", slog.Uint64("id", id))
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Feed filter not found"})
			}
			slog.ErrorContext(ctx, "Failed to query feed filter", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query feed filter"})
		}

		// Look up feed by name if provided
		if req.FeedName != "" {
			var feed models.Feed
			result = db.Where("name = ?", req.FeedName).First(&feed)
			if result.Error != nil {
				if result.Error == gorm.ErrRecordNotFound {
					slog.WarnContext(ctx, "Feed not found", slog.String("feed_name", req.FeedName))
					return c.JSON(http.StatusBadRequest, map[string]string{"error": "Feed not found"})
				}
				slog.ErrorContext(ctx, "Failed to query feed", slog.String("feed_name", req.FeedName), slog.Any("error", result.Error))
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query feed"})
			}
			feedFilter.FeedID = feed.ID
		}

		// Look up torrent category by name
		var torrentCategory models.TorrentCategory
		result = db.Where("name = ?", req.CategoryName).First(&torrentCategory)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				slog.WarnContext(ctx, "Torrent category not found", slog.String("category_name", req.CategoryName))
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Torrent category not found"})
			}
			slog.ErrorContext(ctx, "Failed to query torrent category", slog.String("category_name", req.CategoryName), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query torrent category"})
		}

		// Look up notifier by name
		var notifier models.Notifier
		result = db.Where("name = ?", req.NotifierName).First(&notifier)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				slog.WarnContext(ctx, "Notifier not found", slog.String("notifier_name", req.NotifierName))
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Notifier not found"})
			}
			slog.ErrorContext(ctx, "Failed to query notifier", slog.String("notifier_name", req.NotifierName), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query notifier"})
		}

		feedFilter.Name = req.Name
		feedFilter.TorrentCategoryID = torrentCategory.ID
		feedFilter.NotifierID = notifier.ID

		result = db.Save(&feedFilter)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to update feed filter", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update feed filter"})
		}

		feedName := req.FeedName
		if feedName == "" {
			// If feed_name wasn't provided in request, fetch it from DB
			result = db.Preload("Feed").First(&feedFilter, id)
			if result.Error == nil && feedFilter.Feed.Name != "" {
				feedName = feedFilter.Feed.Name
			}
		}

		response := FeedFilterResponse{
			ID:       feedFilter.ID,
			Name:     feedFilter.Name,
			FeedID:   feedFilter.FeedID,
			FeedName: feedName,
			Category: req.CategoryName,
			Notifier: req.NotifierName,
		}

		slog.InfoContext(ctx, "Successfully updated feed filter", slog.Uint64("id", id), slog.String("name", feedFilter.Name))
		return c.JSON(http.StatusOK, response)
	}
}

// DeleteFeedFilter deletes a feed filter
func DeleteFeedFilter(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		idStr := c.Param("id")

		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			slog.WarnContext(ctx, "Invalid ID parameter", slog.String("id", idStr), slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}

		slog.InfoContext(ctx, "Deleting feed filter", slog.Uint64("id", id))

		result := db.Delete(&models.FeedFilter{}, id)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to delete feed filter", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete feed filter"})
		}

		if result.RowsAffected == 0 {
			slog.WarnContext(ctx, "Feed filter not found", slog.Uint64("id", id))
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Feed filter not found"})
		}

		slog.InfoContext(ctx, "Successfully deleted feed filter", slog.Uint64("id", id))
		return c.NoContent(http.StatusNoContent)
	}
}
