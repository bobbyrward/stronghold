package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
)

type FeedFilterSetEntryRequest struct {
	FeedFilterSetID uint   `json:"feed_filter_set_id" validate:"required"`
	KeyName         string `json:"key_name" validate:"required"`
	OperatorName    string `json:"operator_name" validate:"required"`
	Value           string `json:"value" validate:"required"`
}

type FeedFilterSetEntryResponse struct {
	ID              uint   `json:"id"`
	FeedFilterSetID uint   `json:"feed_filter_set_id"`
	Key             string `json:"key"`
	Operator        string `json:"operator"`
	Value           string `json:"value"`
}

// ListFeedFilterSetEntries returns all feed filter set entries
func ListFeedFilterSetEntries(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		feedFilterSetIDStr := c.QueryParam("feed_filter_set_id")

		slog.InfoContext(ctx, "Listing feed filter set entries", slog.String("feed_filter_set_id_filter", feedFilterSetIDStr))

		var feedFilterSetEntries []models.FeedFilterSetEntry
		query := db.Preload("FilterKey").Preload("FilterOperator")

		if feedFilterSetIDStr != "" {
			feedFilterSetID, err := strconv.ParseUint(feedFilterSetIDStr, 10, 32)
			if err != nil {
				slog.WarnContext(ctx, "Invalid feed_filter_set_id parameter", slog.String("feed_filter_set_id", feedFilterSetIDStr), slog.Any("error", err))
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid feed_filter_set_id parameter"})
			}
			query = query.Where("feed_filter_set_id = ?", feedFilterSetID)
		}

		result := query.Find(&feedFilterSetEntries)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to query feed filter set entries", slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query feed filter set entries"})
		}

		response := make([]FeedFilterSetEntryResponse, len(feedFilterSetEntries))
		for i, ffse := range feedFilterSetEntries {
			keyName := ""
			if ffse.FilterKey.Name != "" {
				keyName = ffse.FilterKey.Name
			}
			operatorName := ""
			if ffse.FilterOperator.Name != "" {
				operatorName = ffse.FilterOperator.Name
			}
			response[i] = FeedFilterSetEntryResponse{
				ID:              ffse.ID,
				FeedFilterSetID: ffse.FeedFilterSetID,
				Key:             keyName,
				Operator:        operatorName,
				Value:           ffse.Value,
			}
		}

		slog.InfoContext(ctx, "Successfully listed feed filter set entries", slog.Int("count", len(response)))
		return c.JSON(http.StatusOK, response)
	}
}

// CreateFeedFilterSetEntry creates a new feed filter set entry
func CreateFeedFilterSetEntry(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var req FeedFilterSetEntryRequest
		if err := c.Bind(&req); err != nil {
			slog.WarnContext(ctx, "Invalid request body", slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		}

		if req.FeedFilterSetID == 0 {
			slog.WarnContext(ctx, "Validation failed: feed_filter_set_id is required")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Feed filter set ID is required"})
		}

		if req.KeyName == "" {
			slog.WarnContext(ctx, "Validation failed: key_name is required")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Key name is required"})
		}

		if req.OperatorName == "" {
			slog.WarnContext(ctx, "Validation failed: operator_name is required")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Operator name is required"})
		}

		if req.Value == "" {
			slog.WarnContext(ctx, "Validation failed: value is required")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Value is required"})
		}

		slog.InfoContext(ctx, "Creating feed filter set entry",
			slog.Uint64("feed_filter_set_id", uint64(req.FeedFilterSetID)),
			slog.String("key_name", req.KeyName),
			slog.String("operator_name", req.OperatorName),
			slog.String("value", req.Value))

		// Look up filter key by name
		var filterKey models.FilterKey
		result := db.Where("name = ?", req.KeyName).First(&filterKey)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				slog.WarnContext(ctx, "Filter key not found", slog.String("key_name", req.KeyName))
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Filter key not found"})
			}
			slog.ErrorContext(ctx, "Failed to query filter key", slog.String("key_name", req.KeyName), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query filter key"})
		}

		// Look up filter operator by name
		var filterOperator models.FilterOperator
		result = db.Where("name = ?", req.OperatorName).First(&filterOperator)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				slog.WarnContext(ctx, "Filter operator not found", slog.String("operator_name", req.OperatorName))
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Filter operator not found"})
			}
			slog.ErrorContext(ctx, "Failed to query filter operator", slog.String("operator_name", req.OperatorName), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query filter operator"})
		}

		feedFilterSetEntry := models.FeedFilterSetEntry{
			FeedFilterSetID:  req.FeedFilterSetID,
			FilterKeyID:      filterKey.ID,
			FilterOperatorID: filterOperator.ID,
			Value:            req.Value,
		}

		result = db.Create(&feedFilterSetEntry)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to create feed filter set entry", slog.Uint64("feed_filter_set_id", uint64(req.FeedFilterSetID)), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create feed filter set entry"})
		}

		response := FeedFilterSetEntryResponse{
			ID:              feedFilterSetEntry.ID,
			FeedFilterSetID: feedFilterSetEntry.FeedFilterSetID,
			Key:             req.KeyName,
			Operator:        req.OperatorName,
			Value:           feedFilterSetEntry.Value,
		}

		slog.InfoContext(ctx, "Successfully created feed filter set entry", slog.Uint64("id", uint64(feedFilterSetEntry.ID)))
		return c.JSON(http.StatusCreated, response)
	}
}

// GetFeedFilterSetEntry returns a single feed filter set entry by ID
func GetFeedFilterSetEntry(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		idStr := c.Param("id")

		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			slog.WarnContext(ctx, "Invalid ID parameter", slog.String("id", idStr), slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}

		slog.InfoContext(ctx, "Getting feed filter set entry", slog.Uint64("id", id))

		var feedFilterSetEntry models.FeedFilterSetEntry
		result := db.Preload("FilterKey").Preload("FilterOperator").First(&feedFilterSetEntry, id)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				slog.WarnContext(ctx, "Feed filter set entry not found", slog.Uint64("id", id))
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Feed filter set entry not found"})
			}
			slog.ErrorContext(ctx, "Failed to query feed filter set entry", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query feed filter set entry"})
		}

		keyName := ""
		if feedFilterSetEntry.FilterKey.Name != "" {
			keyName = feedFilterSetEntry.FilterKey.Name
		}
		operatorName := ""
		if feedFilterSetEntry.FilterOperator.Name != "" {
			operatorName = feedFilterSetEntry.FilterOperator.Name
		}

		response := FeedFilterSetEntryResponse{
			ID:              feedFilterSetEntry.ID,
			FeedFilterSetID: feedFilterSetEntry.FeedFilterSetID,
			Key:             keyName,
			Operator:        operatorName,
			Value:           feedFilterSetEntry.Value,
		}

		slog.InfoContext(ctx, "Successfully retrieved feed filter set entry", slog.Uint64("id", id))
		return c.JSON(http.StatusOK, response)
	}
}

// UpdateFeedFilterSetEntry updates an existing feed filter set entry
func UpdateFeedFilterSetEntry(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		idStr := c.Param("id")

		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			slog.WarnContext(ctx, "Invalid ID parameter", slog.String("id", idStr), slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}

		var req FeedFilterSetEntryRequest
		if err := c.Bind(&req); err != nil {
			slog.WarnContext(ctx, "Invalid request body", slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		}

		if req.FeedFilterSetID == 0 {
			slog.WarnContext(ctx, "Validation failed: feed_filter_set_id is required")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Feed filter set ID is required"})
		}

		if req.KeyName == "" {
			slog.WarnContext(ctx, "Validation failed: key_name is required")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Key name is required"})
		}

		if req.OperatorName == "" {
			slog.WarnContext(ctx, "Validation failed: operator_name is required")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Operator name is required"})
		}

		if req.Value == "" {
			slog.WarnContext(ctx, "Validation failed: value is required")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Value is required"})
		}

		slog.InfoContext(ctx, "Updating feed filter set entry", slog.Uint64("id", id))

		var feedFilterSetEntry models.FeedFilterSetEntry
		result := db.First(&feedFilterSetEntry, id)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				slog.WarnContext(ctx, "Feed filter set entry not found", slog.Uint64("id", id))
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Feed filter set entry not found"})
			}
			slog.ErrorContext(ctx, "Failed to query feed filter set entry", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query feed filter set entry"})
		}

		// Look up filter key by name
		var filterKey models.FilterKey
		result = db.Where("name = ?", req.KeyName).First(&filterKey)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				slog.WarnContext(ctx, "Filter key not found", slog.String("key_name", req.KeyName))
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Filter key not found"})
			}
			slog.ErrorContext(ctx, "Failed to query filter key", slog.String("key_name", req.KeyName), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query filter key"})
		}

		// Look up filter operator by name
		var filterOperator models.FilterOperator
		result = db.Where("name = ?", req.OperatorName).First(&filterOperator)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				slog.WarnContext(ctx, "Filter operator not found", slog.String("operator_name", req.OperatorName))
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Filter operator not found"})
			}
			slog.ErrorContext(ctx, "Failed to query filter operator", slog.String("operator_name", req.OperatorName), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query filter operator"})
		}

		feedFilterSetEntry.FeedFilterSetID = req.FeedFilterSetID
		feedFilterSetEntry.FilterKeyID = filterKey.ID
		feedFilterSetEntry.FilterOperatorID = filterOperator.ID
		feedFilterSetEntry.Value = req.Value

		result = db.Save(&feedFilterSetEntry)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to update feed filter set entry", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update feed filter set entry"})
		}

		response := FeedFilterSetEntryResponse{
			ID:              feedFilterSetEntry.ID,
			FeedFilterSetID: feedFilterSetEntry.FeedFilterSetID,
			Key:             req.KeyName,
			Operator:        req.OperatorName,
			Value:           feedFilterSetEntry.Value,
		}

		slog.InfoContext(ctx, "Successfully updated feed filter set entry", slog.Uint64("id", id))
		return c.JSON(http.StatusOK, response)
	}
}

// DeleteFeedFilterSetEntry deletes a feed filter set entry
func DeleteFeedFilterSetEntry(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		idStr := c.Param("id")

		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			slog.WarnContext(ctx, "Invalid ID parameter", slog.String("id", idStr), slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}

		slog.InfoContext(ctx, "Deleting feed filter set entry", slog.Uint64("id", id))

		result := db.Delete(&models.FeedFilterSetEntry{}, id)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to delete feed filter set entry", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete feed filter set entry"})
		}

		if result.RowsAffected == 0 {
			slog.WarnContext(ctx, "Feed filter set entry not found", slog.Uint64("id", id))
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Feed filter set entry not found"})
		}

		slog.InfoContext(ctx, "Successfully deleted feed filter set entry", slog.Uint64("id", id))
		return c.NoContent(http.StatusNoContent)
	}
}
