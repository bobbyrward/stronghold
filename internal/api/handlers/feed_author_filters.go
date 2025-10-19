package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
)

type FeedAuthorFilterRequest struct {
	Author       string `json:"author" validate:"required"`
	FeedName     string `json:"feed_name" validate:"required"`
	CategoryName string `json:"category_name" validate:"required"`
	NotifierName string `json:"notifier_name" validate:"required"`
}

type FeedAuthorFilterResponse struct {
	ID       uint   `json:"id"`
	Author   string `json:"author"`
	FeedID   uint   `json:"feed_id"`
	FeedName string `json:"feed_name"`
	Category string `json:"category"`
	Notifier string `json:"notifier"`
}

// ListFeedAuthorFilters returns all feed author filters
func ListFeedAuthorFilters(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		feedIDStr := c.QueryParam("feed_id")

		slog.InfoContext(ctx, "Listing feed author filters", slog.String("feed_id_filter", feedIDStr))

		var feedAuthorFilters []models.FeedAuthorFilter
		query := db.Preload("Feed").Preload("TorrentCategory").Preload("Notifier")

		if feedIDStr != "" {
			feedID, err := strconv.ParseUint(feedIDStr, 10, 32)
			if err != nil {
				slog.WarnContext(ctx, "Invalid feed_id parameter", slog.String("feed_id", feedIDStr), slog.Any("error", err))
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid feed_id parameter"})
			}
			query = query.Where("feed_id = ?", feedID)
		}

		result := query.Find(&feedAuthorFilters)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to query feed author filters", slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query feed author filters"})
		}

		response := make([]FeedAuthorFilterResponse, len(feedAuthorFilters))
		for i, faf := range feedAuthorFilters {
			feedName := ""
			if faf.Feed.Name != "" {
				feedName = faf.Feed.Name
			}
			categoryName := ""
			if faf.TorrentCategory.Name != "" {
				categoryName = faf.TorrentCategory.Name
			}
			notifierName := ""
			if faf.Notifier.Name != "" {
				notifierName = faf.Notifier.Name
			}
			response[i] = FeedAuthorFilterResponse{
				ID:       faf.ID,
				Author:   faf.Author,
				FeedID:   faf.FeedID,
				FeedName: feedName,
				Category: categoryName,
				Notifier: notifierName,
			}
		}

		slog.InfoContext(ctx, "Successfully listed feed author filters", slog.Int("count", len(response)))
		return c.JSON(http.StatusOK, response)
	}
}

// CreateFeedAuthorFilter creates a new feed author filter
func CreateFeedAuthorFilter(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var req FeedAuthorFilterRequest
		if err := c.Bind(&req); err != nil {
			slog.WarnContext(ctx, "Invalid request body", slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		}

		if req.Author == "" {
			slog.WarnContext(ctx, "Validation failed: author is required")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Author is required"})
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

		slog.InfoContext(ctx, "Creating feed author filter",
			slog.String("author", req.Author),
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

		feedAuthorFilter := models.FeedAuthorFilter{
			Author:            req.Author,
			FeedID:            feed.ID,
			TorrentCategoryID: torrentCategory.ID,
			NotifierID:        notifier.ID,
		}

		result = db.Create(&feedAuthorFilter)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to create feed author filter", slog.String("author", req.Author), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create feed author filter"})
		}

		response := FeedAuthorFilterResponse{
			ID:       feedAuthorFilter.ID,
			Author:   feedAuthorFilter.Author,
			FeedID:   feedAuthorFilter.FeedID,
			FeedName: req.FeedName,
			Category: req.CategoryName,
			Notifier: req.NotifierName,
		}

		slog.InfoContext(ctx, "Successfully created feed author filter", slog.Uint64("id", uint64(feedAuthorFilter.ID)), slog.String("author", feedAuthorFilter.Author))
		return c.JSON(http.StatusCreated, response)
	}
}

// GetFeedAuthorFilter returns a single feed author filter by ID
func GetFeedAuthorFilter(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		idStr := c.Param("id")

		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			slog.WarnContext(ctx, "Invalid ID parameter", slog.String("id", idStr), slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}

		slog.InfoContext(ctx, "Getting feed author filter", slog.Uint64("id", id))

		var feedAuthorFilter models.FeedAuthorFilter
		result := db.Preload("Feed").Preload("TorrentCategory").Preload("Notifier").First(&feedAuthorFilter, id)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				slog.WarnContext(ctx, "Feed author filter not found", slog.Uint64("id", id))
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Feed author filter not found"})
			}
			slog.ErrorContext(ctx, "Failed to query feed author filter", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query feed author filter"})
		}

		feedName := ""
		if feedAuthorFilter.Feed.Name != "" {
			feedName = feedAuthorFilter.Feed.Name
		}
		categoryName := ""
		if feedAuthorFilter.TorrentCategory.Name != "" {
			categoryName = feedAuthorFilter.TorrentCategory.Name
		}
		notifierName := ""
		if feedAuthorFilter.Notifier.Name != "" {
			notifierName = feedAuthorFilter.Notifier.Name
		}

		response := FeedAuthorFilterResponse{
			ID:       feedAuthorFilter.ID,
			Author:   feedAuthorFilter.Author,
			FeedID:   feedAuthorFilter.FeedID,
			FeedName: feedName,
			Category: categoryName,
			Notifier: notifierName,
		}

		slog.InfoContext(ctx, "Successfully retrieved feed author filter", slog.Uint64("id", id), slog.String("author", feedAuthorFilter.Author))
		return c.JSON(http.StatusOK, response)
	}
}

// UpdateFeedAuthorFilter updates an existing feed author filter
func UpdateFeedAuthorFilter(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		idStr := c.Param("id")

		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			slog.WarnContext(ctx, "Invalid ID parameter", slog.String("id", idStr), slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}

		var req FeedAuthorFilterRequest
		if err := c.Bind(&req); err != nil {
			slog.WarnContext(ctx, "Invalid request body", slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		}

		if req.Author == "" {
			slog.WarnContext(ctx, "Validation failed: author is required")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Author is required"})
		}

		if req.CategoryName == "" {
			slog.WarnContext(ctx, "Validation failed: category_name is required")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Category name is required"})
		}

		if req.NotifierName == "" {
			slog.WarnContext(ctx, "Validation failed: notifier_name is required")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Notifier name is required"})
		}

		slog.InfoContext(ctx, "Updating feed author filter", slog.Uint64("id", id), slog.String("author", req.Author))

		var feedAuthorFilter models.FeedAuthorFilter
		result := db.First(&feedAuthorFilter, id)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				slog.WarnContext(ctx, "Feed author filter not found", slog.Uint64("id", id))
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Feed author filter not found"})
			}
			slog.ErrorContext(ctx, "Failed to query feed author filter", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query feed author filter"})
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
			feedAuthorFilter.FeedID = feed.ID
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

		feedAuthorFilter.Author = req.Author
		feedAuthorFilter.TorrentCategoryID = torrentCategory.ID
		feedAuthorFilter.NotifierID = notifier.ID

		result = db.Save(&feedAuthorFilter)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to update feed author filter", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update feed author filter"})
		}

		feedName := req.FeedName
		if feedName == "" {
			// If feed_name wasn't provided in request, fetch it from DB
			result = db.Preload("Feed").First(&feedAuthorFilter, id)
			if result.Error == nil && feedAuthorFilter.Feed.Name != "" {
				feedName = feedAuthorFilter.Feed.Name
			}
		}

		response := FeedAuthorFilterResponse{
			ID:       feedAuthorFilter.ID,
			Author:   feedAuthorFilter.Author,
			FeedID:   feedAuthorFilter.FeedID,
			FeedName: feedName,
			Category: req.CategoryName,
			Notifier: req.NotifierName,
		}

		slog.InfoContext(ctx, "Successfully updated feed author filter", slog.Uint64("id", id), slog.String("author", feedAuthorFilter.Author))
		return c.JSON(http.StatusOK, response)
	}
}

// DeleteFeedAuthorFilter deletes a feed author filter
func DeleteFeedAuthorFilter(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		idStr := c.Param("id")

		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			slog.WarnContext(ctx, "Invalid ID parameter", slog.String("id", idStr), slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}

		slog.InfoContext(ctx, "Deleting feed author filter", slog.Uint64("id", id))

		result := db.Delete(&models.FeedAuthorFilter{}, id)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to delete feed author filter", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete feed author filter"})
		}

		if result.RowsAffected == 0 {
			slog.WarnContext(ctx, "Feed author filter not found", slog.Uint64("id", id))
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Feed author filter not found"})
		}

		slog.InfoContext(ctx, "Successfully deleted feed author filter", slog.Uint64("id", id))
		return c.NoContent(http.StatusNoContent)
	}
}
