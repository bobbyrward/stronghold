package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
)

type FeedRequest struct {
	Name string `json:"name" validate:"required"`
	URL  string `json:"url" validate:"required"`
}

type FeedResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

// ListFeeds returns all feeds
func ListFeeds(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		slog.InfoContext(ctx, "Listing feeds")

		var feeds []models.Feed
		result := db.Find(&feeds)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to query feeds", slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query feeds"})
		}

		response := make([]FeedResponse, len(feeds))
		for i, f := range feeds {
			response[i] = FeedResponse{
				ID:   f.ID,
				Name: f.Name,
				URL:  f.URL,
			}
		}

		slog.InfoContext(ctx, "Successfully listed feeds", slog.Int("count", len(response)))
		return c.JSON(http.StatusOK, response)
	}
}

// CreateFeed creates a new feed
func CreateFeed(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var req FeedRequest
		if err := c.Bind(&req); err != nil {
			slog.WarnContext(ctx, "Invalid request body", slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		}

		if req.Name == "" {
			slog.WarnContext(ctx, "Validation failed: name is required")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Name is required"})
		}

		if req.URL == "" {
			slog.WarnContext(ctx, "Validation failed: url is required")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "URL is required"})
		}

		slog.InfoContext(ctx, "Creating feed", slog.String("name", req.Name), slog.String("url", req.URL))

		feed := models.Feed{
			Name: req.Name,
			URL:  req.URL,
		}

		result := db.Create(&feed)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to create feed", slog.String("name", req.Name), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create feed"})
		}

		response := FeedResponse{
			ID:   feed.ID,
			Name: feed.Name,
			URL:  feed.URL,
		}

		slog.InfoContext(ctx, "Successfully created feed", slog.Uint64("id", uint64(feed.ID)), slog.String("name", feed.Name))
		return c.JSON(http.StatusCreated, response)
	}
}

// GetFeed returns a single feed by ID
func GetFeed(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		idStr := c.Param("id")

		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			slog.WarnContext(ctx, "Invalid ID parameter", slog.String("id", idStr), slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}

		slog.InfoContext(ctx, "Getting feed", slog.Uint64("id", id))

		var feed models.Feed
		result := db.First(&feed, id)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				slog.WarnContext(ctx, "Feed not found", slog.Uint64("id", id))
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Feed not found"})
			}
			slog.ErrorContext(ctx, "Failed to query feed", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query feed"})
		}

		response := FeedResponse{
			ID:   feed.ID,
			Name: feed.Name,
			URL:  feed.URL,
		}

		slog.InfoContext(ctx, "Successfully retrieved feed", slog.Uint64("id", id), slog.String("name", feed.Name))
		return c.JSON(http.StatusOK, response)
	}
}

// UpdateFeed updates an existing feed
func UpdateFeed(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		idStr := c.Param("id")

		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			slog.WarnContext(ctx, "Invalid ID parameter", slog.String("id", idStr), slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}

		var req FeedRequest
		if err := c.Bind(&req); err != nil {
			slog.WarnContext(ctx, "Invalid request body", slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		}

		if req.Name == "" {
			slog.WarnContext(ctx, "Validation failed: name is required")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Name is required"})
		}

		if req.URL == "" {
			slog.WarnContext(ctx, "Validation failed: url is required")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "URL is required"})
		}

		slog.InfoContext(ctx, "Updating feed", slog.Uint64("id", id), slog.String("name", req.Name))

		var feed models.Feed
		result := db.First(&feed, id)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				slog.WarnContext(ctx, "Feed not found", slog.Uint64("id", id))
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Feed not found"})
			}
			slog.ErrorContext(ctx, "Failed to query feed", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query feed"})
		}

		feed.Name = req.Name
		feed.URL = req.URL

		result = db.Save(&feed)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to update feed", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update feed"})
		}

		response := FeedResponse{
			ID:   feed.ID,
			Name: feed.Name,
			URL:  feed.URL,
		}

		slog.InfoContext(ctx, "Successfully updated feed", slog.Uint64("id", id), slog.String("name", feed.Name))
		return c.JSON(http.StatusOK, response)
	}
}

// DeleteFeed deletes a feed
func DeleteFeed(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		idStr := c.Param("id")

		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			slog.WarnContext(ctx, "Invalid ID parameter", slog.String("id", idStr), slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}

		slog.InfoContext(ctx, "Deleting feed", slog.Uint64("id", id))

		result := db.Delete(&models.Feed{}, id)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to delete feed", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete feed"})
		}

		if result.RowsAffected == 0 {
			slog.WarnContext(ctx, "Feed not found", slog.Uint64("id", id))
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Feed not found"})
		}

		slog.InfoContext(ctx, "Successfully deleted feed", slog.Uint64("id", id))
		return c.NoContent(http.StatusNoContent)
	}
}
