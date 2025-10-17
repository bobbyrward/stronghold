package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
)

type FilterOperatorRequest struct {
	Name string `json:"name" validate:"required"`
}

type FilterOperatorResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// ListFilterOperators returns all filter operators
func ListFilterOperators(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		slog.InfoContext(ctx, "Listing filter operators")

		var filterOperators []models.FilterOperator
		result := db.Find(&filterOperators)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to query filter operators", slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query filter operators"})
		}

		response := make([]FilterOperatorResponse, len(filterOperators))
		for i, fo := range filterOperators {
			response[i] = FilterOperatorResponse{
				ID:   fo.ID,
				Name: fo.Name,
			}
		}

		slog.InfoContext(ctx, "Successfully listed filter operators", slog.Int("count", len(response)))
		return c.JSON(http.StatusOK, response)
	}
}

// CreateFilterOperator creates a new filter operator
func CreateFilterOperator(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var req FilterOperatorRequest
		if err := c.Bind(&req); err != nil {
			slog.WarnContext(ctx, "Invalid request body", slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		}

		if req.Name == "" {
			slog.WarnContext(ctx, "Validation failed: name is required")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Name is required"})
		}

		slog.InfoContext(ctx, "Creating filter operator", slog.String("name", req.Name))

		filterOperator := models.FilterOperator{
			Name: req.Name,
		}

		result := db.Create(&filterOperator)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to create filter operator", slog.String("name", req.Name), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create filter operator"})
		}

		response := FilterOperatorResponse{
			ID:   filterOperator.ID,
			Name: filterOperator.Name,
		}

		slog.InfoContext(ctx, "Successfully created filter operator", slog.Uint64("id", uint64(filterOperator.ID)), slog.String("name", filterOperator.Name))
		return c.JSON(http.StatusCreated, response)
	}
}

// GetFilterOperator returns a single filter operator by ID
func GetFilterOperator(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		idStr := c.Param("id")

		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			slog.WarnContext(ctx, "Invalid ID parameter", slog.String("id", idStr), slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}

		slog.InfoContext(ctx, "Getting filter operator", slog.Uint64("id", id))

		var filterOperator models.FilterOperator
		result := db.First(&filterOperator, id)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				slog.WarnContext(ctx, "Filter operator not found", slog.Uint64("id", id))
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Filter operator not found"})
			}
			slog.ErrorContext(ctx, "Failed to query filter operator", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query filter operator"})
		}

		response := FilterOperatorResponse{
			ID:   filterOperator.ID,
			Name: filterOperator.Name,
		}

		slog.InfoContext(ctx, "Successfully retrieved filter operator", slog.Uint64("id", id), slog.String("name", filterOperator.Name))
		return c.JSON(http.StatusOK, response)
	}
}

// UpdateFilterOperator updates an existing filter operator
func UpdateFilterOperator(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		idStr := c.Param("id")

		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			slog.WarnContext(ctx, "Invalid ID parameter", slog.String("id", idStr), slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}

		var req FilterOperatorRequest
		if err := c.Bind(&req); err != nil {
			slog.WarnContext(ctx, "Invalid request body", slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		}

		if req.Name == "" {
			slog.WarnContext(ctx, "Validation failed: name is required")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Name is required"})
		}

		slog.InfoContext(ctx, "Updating filter operator", slog.Uint64("id", id), slog.String("name", req.Name))

		var filterOperator models.FilterOperator
		result := db.First(&filterOperator, id)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				slog.WarnContext(ctx, "Filter operator not found", slog.Uint64("id", id))
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Filter operator not found"})
			}
			slog.ErrorContext(ctx, "Failed to query filter operator", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query filter operator"})
		}

		filterOperator.Name = req.Name
		result = db.Save(&filterOperator)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to update filter operator", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update filter operator"})
		}

		response := FilterOperatorResponse{
			ID:   filterOperator.ID,
			Name: filterOperator.Name,
		}

		slog.InfoContext(ctx, "Successfully updated filter operator", slog.Uint64("id", id), slog.String("name", filterOperator.Name))
		return c.JSON(http.StatusOK, response)
	}
}

// DeleteFilterOperator deletes a filter operator
func DeleteFilterOperator(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		idStr := c.Param("id")

		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			slog.WarnContext(ctx, "Invalid ID parameter", slog.String("id", idStr), slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}

		slog.InfoContext(ctx, "Deleting filter operator", slog.Uint64("id", id))

		result := db.Delete(&models.FilterOperator{}, id)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to delete filter operator", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete filter operator"})
		}

		if result.RowsAffected == 0 {
			slog.WarnContext(ctx, "Filter operator not found", slog.Uint64("id", id))
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Filter operator not found"})
		}

		slog.InfoContext(ctx, "Successfully deleted filter operator", slog.Uint64("id", id))
		return c.NoContent(http.StatusNoContent)
	}
}
