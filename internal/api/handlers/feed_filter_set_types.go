package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
)

type FeedFilterSetTypeRequest struct {
	Name string `json:"name" validate:"required"`
}

type FeedFilterSetTypeResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// ListFeedFilterSetTypes returns all feed filter set types
func ListFeedFilterSetTypes(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		slog.InfoContext(ctx, "Listing feed filter set types")

		var feedFilterSetTypes []models.FeedFilterSetType
		result := db.Find(&feedFilterSetTypes)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to query feed filter set types", slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query feed filter set types"})
		}

		response := make([]FeedFilterSetTypeResponse, len(feedFilterSetTypes))
		for i, ffst := range feedFilterSetTypes {
			response[i] = FeedFilterSetTypeResponse{
				ID:   ffst.ID,
				Name: ffst.Name,
			}
		}

		slog.InfoContext(ctx, "Successfully listed feed filter set types", slog.Int("count", len(response)))
		return c.JSON(http.StatusOK, response)
	}
}

// CreateFeedFilterSetType creates a new feed filter set type
func CreateFeedFilterSetType(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var req FeedFilterSetTypeRequest
		if err := c.Bind(&req); err != nil {
			slog.WarnContext(ctx, "Invalid request body", slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		}

		if req.Name == "" {
			slog.WarnContext(ctx, "Validation failed: name is required")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Name is required"})
		}

		slog.InfoContext(ctx, "Creating feed filter set type", slog.String("name", req.Name))

		feedFilterSetType := models.FeedFilterSetType{
			Name: req.Name,
		}

		result := db.Create(&feedFilterSetType)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to create feed filter set type", slog.String("name", req.Name), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create feed filter set type"})
		}

		response := FeedFilterSetTypeResponse{
			ID:   feedFilterSetType.ID,
			Name: feedFilterSetType.Name,
		}

		slog.InfoContext(ctx, "Successfully created feed filter set type", slog.Uint64("id", uint64(feedFilterSetType.ID)), slog.String("name", feedFilterSetType.Name))
		return c.JSON(http.StatusCreated, response)
	}
}

// GetFeedFilterSetType returns a single feed filter set type by ID
func GetFeedFilterSetType(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		idStr := c.Param("id")

		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			slog.WarnContext(ctx, "Invalid ID parameter", slog.String("id", idStr), slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}

		slog.InfoContext(ctx, "Getting feed filter set type", slog.Uint64("id", id))

		var feedFilterSetType models.FeedFilterSetType
		result := db.First(&feedFilterSetType, id)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				slog.WarnContext(ctx, "Feed filter set type not found", slog.Uint64("id", id))
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Feed filter set type not found"})
			}
			slog.ErrorContext(ctx, "Failed to query feed filter set type", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query feed filter set type"})
		}

		response := FeedFilterSetTypeResponse{
			ID:   feedFilterSetType.ID,
			Name: feedFilterSetType.Name,
		}

		slog.InfoContext(ctx, "Successfully retrieved feed filter set type", slog.Uint64("id", id), slog.String("name", feedFilterSetType.Name))
		return c.JSON(http.StatusOK, response)
	}
}

// UpdateFeedFilterSetType updates an existing feed filter set type
func UpdateFeedFilterSetType(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		idStr := c.Param("id")

		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			slog.WarnContext(ctx, "Invalid ID parameter", slog.String("id", idStr), slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}

		var req FeedFilterSetTypeRequest
		if err := c.Bind(&req); err != nil {
			slog.WarnContext(ctx, "Invalid request body", slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		}

		if req.Name == "" {
			slog.WarnContext(ctx, "Validation failed: name is required")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Name is required"})
		}

		slog.InfoContext(ctx, "Updating feed filter set type", slog.Uint64("id", id), slog.String("name", req.Name))

		var feedFilterSetType models.FeedFilterSetType
		result := db.First(&feedFilterSetType, id)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				slog.WarnContext(ctx, "Feed filter set type not found", slog.Uint64("id", id))
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Feed filter set type not found"})
			}
			slog.ErrorContext(ctx, "Failed to query feed filter set type", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query feed filter set type"})
		}

		feedFilterSetType.Name = req.Name
		result = db.Save(&feedFilterSetType)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to update feed filter set type", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update feed filter set type"})
		}

		response := FeedFilterSetTypeResponse{
			ID:   feedFilterSetType.ID,
			Name: feedFilterSetType.Name,
		}

		slog.InfoContext(ctx, "Successfully updated feed filter set type", slog.Uint64("id", id), slog.String("name", feedFilterSetType.Name))
		return c.JSON(http.StatusOK, response)
	}
}

// DeleteFeedFilterSetType deletes a feed filter set type
func DeleteFeedFilterSetType(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		idStr := c.Param("id")

		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			slog.WarnContext(ctx, "Invalid ID parameter", slog.String("id", idStr), slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}

		slog.InfoContext(ctx, "Deleting feed filter set type", slog.Uint64("id", id))

		result := db.Delete(&models.FeedFilterSetType{}, id)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to delete feed filter set type", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete feed filter set type"})
		}

		if result.RowsAffected == 0 {
			slog.WarnContext(ctx, "Feed filter set type not found", slog.Uint64("id", id))
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Feed filter set type not found"})
		}

		slog.InfoContext(ctx, "Successfully deleted feed filter set type", slog.Uint64("id", id))
		return c.NoContent(http.StatusNoContent)
	}
}
