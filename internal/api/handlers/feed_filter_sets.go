package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
)

type FeedFilterSetRequest struct {
	FeedFilterID uint   `json:"feed_filter_id" validate:"required"`
	TypeName     string `json:"type_name" validate:"required"`
}

type FeedFilterSetResponse struct {
	ID           uint   `json:"id"`
	FeedFilterID uint   `json:"feed_filter_id"`
	Type         string `json:"type"`
}

// ListFeedFilterSets returns all feed filter sets
func ListFeedFilterSets(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		feedFilterIDStr := c.QueryParam("feed_filter_id")

		slog.InfoContext(ctx, "Listing feed filter sets", slog.String("feed_filter_id_filter", feedFilterIDStr))

		var feedFilterSets []models.FeedFilterSet
		query := db.Preload("FeedFilterSetType")

		if feedFilterIDStr != "" {
			feedFilterID, err := strconv.ParseUint(feedFilterIDStr, 10, 32)
			if err != nil {
				slog.WarnContext(ctx, "Invalid feed_filter_id parameter", slog.String("feed_filter_id", feedFilterIDStr), slog.Any("error", err))
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid feed_filter_id parameter"})
			}
			query = query.Where("feed_filter_id = ?", feedFilterID)
		}

		result := query.Find(&feedFilterSets)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to query feed filter sets", slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query feed filter sets"})
		}

		response := make([]FeedFilterSetResponse, len(feedFilterSets))
		for i, ffs := range feedFilterSets {
			typeName := ""
			if ffs.FeedFilterSetType.Name != "" {
				typeName = ffs.FeedFilterSetType.Name
			}
			response[i] = FeedFilterSetResponse{
				ID:           ffs.ID,
				FeedFilterID: ffs.FeedFilterID,
				Type:         typeName,
			}
		}

		slog.InfoContext(ctx, "Successfully listed feed filter sets", slog.Int("count", len(response)))
		return c.JSON(http.StatusOK, response)
	}
}

// CreateFeedFilterSet creates a new feed filter set
func CreateFeedFilterSet(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var req FeedFilterSetRequest
		if err := c.Bind(&req); err != nil {
			slog.WarnContext(ctx, "Invalid request body", slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		}

		if req.FeedFilterID == 0 {
			slog.WarnContext(ctx, "Validation failed: feed_filter_id is required")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Feed filter ID is required"})
		}

		if req.TypeName == "" {
			slog.WarnContext(ctx, "Validation failed: type_name is required")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Type name is required"})
		}

		slog.InfoContext(ctx, "Creating feed filter set",
			slog.Uint64("feed_filter_id", uint64(req.FeedFilterID)),
			slog.String("type_name", req.TypeName))

		// Look up feed filter set type by name
		var feedFilterSetType models.FeedFilterSetType
		result := db.Where("name = ?", req.TypeName).First(&feedFilterSetType)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				slog.WarnContext(ctx, "Feed filter set type not found", slog.String("type_name", req.TypeName))
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Feed filter set type not found"})
			}
			slog.ErrorContext(ctx, "Failed to query feed filter set type", slog.String("type_name", req.TypeName), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query feed filter set type"})
		}

		feedFilterSet := models.FeedFilterSet{
			FeedFilterID:        req.FeedFilterID,
			FeedFilterSetTypeID: feedFilterSetType.ID,
		}

		result = db.Create(&feedFilterSet)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to create feed filter set", slog.Uint64("feed_filter_id", uint64(req.FeedFilterID)), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create feed filter set"})
		}

		response := FeedFilterSetResponse{
			ID:           feedFilterSet.ID,
			FeedFilterID: feedFilterSet.FeedFilterID,
			Type:         req.TypeName,
		}

		slog.InfoContext(ctx, "Successfully created feed filter set", slog.Uint64("id", uint64(feedFilterSet.ID)))
		return c.JSON(http.StatusCreated, response)
	}
}

// GetFeedFilterSet returns a single feed filter set by ID
func GetFeedFilterSet(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		idStr := c.Param("id")

		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			slog.WarnContext(ctx, "Invalid ID parameter", slog.String("id", idStr), slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}

		slog.InfoContext(ctx, "Getting feed filter set", slog.Uint64("id", id))

		var feedFilterSet models.FeedFilterSet
		result := db.Preload("FeedFilterSetType").First(&feedFilterSet, id)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				slog.WarnContext(ctx, "Feed filter set not found", slog.Uint64("id", id))
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Feed filter set not found"})
			}
			slog.ErrorContext(ctx, "Failed to query feed filter set", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query feed filter set"})
		}

		typeName := ""
		if feedFilterSet.FeedFilterSetType.Name != "" {
			typeName = feedFilterSet.FeedFilterSetType.Name
		}

		response := FeedFilterSetResponse{
			ID:           feedFilterSet.ID,
			FeedFilterID: feedFilterSet.FeedFilterID,
			Type:         typeName,
		}

		slog.InfoContext(ctx, "Successfully retrieved feed filter set", slog.Uint64("id", id))
		return c.JSON(http.StatusOK, response)
	}
}

// UpdateFeedFilterSet updates an existing feed filter set
func UpdateFeedFilterSet(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		idStr := c.Param("id")

		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			slog.WarnContext(ctx, "Invalid ID parameter", slog.String("id", idStr), slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}

		var req FeedFilterSetRequest
		if err := c.Bind(&req); err != nil {
			slog.WarnContext(ctx, "Invalid request body", slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		}

		if req.FeedFilterID == 0 {
			slog.WarnContext(ctx, "Validation failed: feed_filter_id is required")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Feed filter ID is required"})
		}

		if req.TypeName == "" {
			slog.WarnContext(ctx, "Validation failed: type_name is required")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Type name is required"})
		}

		slog.InfoContext(ctx, "Updating feed filter set", slog.Uint64("id", id))

		var feedFilterSet models.FeedFilterSet
		result := db.First(&feedFilterSet, id)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				slog.WarnContext(ctx, "Feed filter set not found", slog.Uint64("id", id))
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Feed filter set not found"})
			}
			slog.ErrorContext(ctx, "Failed to query feed filter set", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query feed filter set"})
		}

		// Look up feed filter set type by name
		var feedFilterSetType models.FeedFilterSetType
		result = db.Where("name = ?", req.TypeName).First(&feedFilterSetType)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				slog.WarnContext(ctx, "Feed filter set type not found", slog.String("type_name", req.TypeName))
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Feed filter set type not found"})
			}
			slog.ErrorContext(ctx, "Failed to query feed filter set type", slog.String("type_name", req.TypeName), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query feed filter set type"})
		}

		feedFilterSet.FeedFilterID = req.FeedFilterID
		feedFilterSet.FeedFilterSetTypeID = feedFilterSetType.ID

		result = db.Save(&feedFilterSet)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to update feed filter set", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update feed filter set"})
		}

		response := FeedFilterSetResponse{
			ID:           feedFilterSet.ID,
			FeedFilterID: feedFilterSet.FeedFilterID,
			Type:         req.TypeName,
		}

		slog.InfoContext(ctx, "Successfully updated feed filter set", slog.Uint64("id", id))
		return c.JSON(http.StatusOK, response)
	}
}

// DeleteFeedFilterSet deletes a feed filter set
func DeleteFeedFilterSet(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		idStr := c.Param("id")

		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			slog.WarnContext(ctx, "Invalid ID parameter", slog.String("id", idStr), slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}

		slog.InfoContext(ctx, "Deleting feed filter set", slog.Uint64("id", id))

		result := db.Delete(&models.FeedFilterSet{}, id)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to delete feed filter set", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete feed filter set"})
		}

		if result.RowsAffected == 0 {
			slog.WarnContext(ctx, "Feed filter set not found", slog.Uint64("id", id))
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Feed filter set not found"})
		}

		slog.InfoContext(ctx, "Successfully deleted feed filter set", slog.Uint64("id", id))
		return c.NoContent(http.StatusNoContent)
	}
}
