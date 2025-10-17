package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
)

type FilterKeyRequest struct {
	Name string `json:"name" validate:"required"`
}

type FilterKeyResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// ListFilterKeys returns all filter keys
func ListFilterKeys(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		slog.InfoContext(ctx, "Listing filter keys")

		var filterKeys []models.FilterKey
		result := db.Find(&filterKeys)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to query filter keys", slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query filter keys"})
		}

		response := make([]FilterKeyResponse, len(filterKeys))
		for i, fk := range filterKeys {
			response[i] = FilterKeyResponse{
				ID:   fk.ID,
				Name: fk.Name,
			}
		}

		slog.InfoContext(ctx, "Successfully listed filter keys", slog.Int("count", len(response)))
		return c.JSON(http.StatusOK, response)
	}
}

// CreateFilterKey creates a new filter key
func CreateFilterKey(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var req FilterKeyRequest
		if err := c.Bind(&req); err != nil {
			slog.WarnContext(ctx, "Invalid request body", slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		}

		if req.Name == "" {
			slog.WarnContext(ctx, "Validation failed: name is required")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Name is required"})
		}

		slog.InfoContext(ctx, "Creating filter key", slog.String("name", req.Name))

		filterKey := models.FilterKey{
			Name: req.Name,
		}

		result := db.Create(&filterKey)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to create filter key", slog.String("name", req.Name), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create filter key"})
		}

		response := FilterKeyResponse{
			ID:   filterKey.ID,
			Name: filterKey.Name,
		}

		slog.InfoContext(ctx, "Successfully created filter key", slog.Uint64("id", uint64(filterKey.ID)), slog.String("name", filterKey.Name))
		return c.JSON(http.StatusCreated, response)
	}
}

// GetFilterKey returns a single filter key by ID
func GetFilterKey(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		idStr := c.Param("id")

		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			slog.WarnContext(ctx, "Invalid ID parameter", slog.String("id", idStr), slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}

		slog.InfoContext(ctx, "Getting filter key", slog.Uint64("id", id))

		var filterKey models.FilterKey
		result := db.First(&filterKey, id)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				slog.WarnContext(ctx, "Filter key not found", slog.Uint64("id", id))
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Filter key not found"})
			}
			slog.ErrorContext(ctx, "Failed to query filter key", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query filter key"})
		}

		response := FilterKeyResponse{
			ID:   filterKey.ID,
			Name: filterKey.Name,
		}

		slog.InfoContext(ctx, "Successfully retrieved filter key", slog.Uint64("id", id), slog.String("name", filterKey.Name))
		return c.JSON(http.StatusOK, response)
	}
}

// UpdateFilterKey updates an existing filter key
func UpdateFilterKey(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		idStr := c.Param("id")

		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			slog.WarnContext(ctx, "Invalid ID parameter", slog.String("id", idStr), slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}

		var req FilterKeyRequest
		if err := c.Bind(&req); err != nil {
			slog.WarnContext(ctx, "Invalid request body", slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		}

		if req.Name == "" {
			slog.WarnContext(ctx, "Validation failed: name is required")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Name is required"})
		}

		slog.InfoContext(ctx, "Updating filter key", slog.Uint64("id", id), slog.String("name", req.Name))

		var filterKey models.FilterKey
		result := db.First(&filterKey, id)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				slog.WarnContext(ctx, "Filter key not found", slog.Uint64("id", id))
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Filter key not found"})
			}
			slog.ErrorContext(ctx, "Failed to query filter key", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query filter key"})
		}

		filterKey.Name = req.Name
		result = db.Save(&filterKey)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to update filter key", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update filter key"})
		}

		response := FilterKeyResponse{
			ID:   filterKey.ID,
			Name: filterKey.Name,
		}

		slog.InfoContext(ctx, "Successfully updated filter key", slog.Uint64("id", id), slog.String("name", filterKey.Name))
		return c.JSON(http.StatusOK, response)
	}
}

// DeleteFilterKey deletes a filter key
func DeleteFilterKey(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		idStr := c.Param("id")

		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			slog.WarnContext(ctx, "Invalid ID parameter", slog.String("id", idStr), slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}

		slog.InfoContext(ctx, "Deleting filter key", slog.Uint64("id", id))

		result := db.Delete(&models.FilterKey{}, id)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to delete filter key", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete filter key"})
		}

		if result.RowsAffected == 0 {
			slog.WarnContext(ctx, "Filter key not found", slog.Uint64("id", id))
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Filter key not found"})
		}

		slog.InfoContext(ctx, "Successfully deleted filter key", slog.Uint64("id", id))
		return c.NoContent(http.StatusNoContent)
	}
}
