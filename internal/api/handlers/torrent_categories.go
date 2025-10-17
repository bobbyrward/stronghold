package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
)

type TorrentCategoryRequest struct {
	Name string `json:"name" validate:"required"`
}

type TorrentCategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// ListTorrentCategories returns all torrent categories
func ListTorrentCategories(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		slog.InfoContext(ctx, "Listing torrent categories")

		var torrentCategories []models.TorrentCategory
		result := db.Find(&torrentCategories)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to query torrent categories", slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query torrent categories"})
		}

		response := make([]TorrentCategoryResponse, len(torrentCategories))
		for i, tc := range torrentCategories {
			response[i] = TorrentCategoryResponse{
				ID:   tc.ID,
				Name: tc.Name,
			}
		}

		slog.InfoContext(ctx, "Successfully listed torrent categories", slog.Int("count", len(response)))
		return c.JSON(http.StatusOK, response)
	}
}

// CreateTorrentCategory creates a new torrent category
func CreateTorrentCategory(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var req TorrentCategoryRequest
		if err := c.Bind(&req); err != nil {
			slog.WarnContext(ctx, "Invalid request body", slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		}

		if req.Name == "" {
			slog.WarnContext(ctx, "Validation failed: name is required")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Name is required"})
		}

		slog.InfoContext(ctx, "Creating torrent category", slog.String("name", req.Name))

		torrentCategory := models.TorrentCategory{
			Name: req.Name,
		}

		result := db.Create(&torrentCategory)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to create torrent category", slog.String("name", req.Name), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create torrent category"})
		}

		response := TorrentCategoryResponse{
			ID:   torrentCategory.ID,
			Name: torrentCategory.Name,
		}

		slog.InfoContext(ctx, "Successfully created torrent category", slog.Uint64("id", uint64(torrentCategory.ID)), slog.String("name", torrentCategory.Name))
		return c.JSON(http.StatusCreated, response)
	}
}

// GetTorrentCategory returns a single torrent category by ID
func GetTorrentCategory(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		idStr := c.Param("id")

		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			slog.WarnContext(ctx, "Invalid ID parameter", slog.String("id", idStr), slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}

		slog.InfoContext(ctx, "Getting torrent category", slog.Uint64("id", id))

		var torrentCategory models.TorrentCategory
		result := db.First(&torrentCategory, id)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				slog.WarnContext(ctx, "Torrent category not found", slog.Uint64("id", id))
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Torrent category not found"})
			}
			slog.ErrorContext(ctx, "Failed to query torrent category", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query torrent category"})
		}

		response := TorrentCategoryResponse{
			ID:   torrentCategory.ID,
			Name: torrentCategory.Name,
		}

		slog.InfoContext(ctx, "Successfully retrieved torrent category", slog.Uint64("id", id), slog.String("name", torrentCategory.Name))
		return c.JSON(http.StatusOK, response)
	}
}

// UpdateTorrentCategory updates an existing torrent category
func UpdateTorrentCategory(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		idStr := c.Param("id")

		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			slog.WarnContext(ctx, "Invalid ID parameter", slog.String("id", idStr), slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}

		var req TorrentCategoryRequest
		if err := c.Bind(&req); err != nil {
			slog.WarnContext(ctx, "Invalid request body", slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		}

		if req.Name == "" {
			slog.WarnContext(ctx, "Validation failed: name is required")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Name is required"})
		}

		slog.InfoContext(ctx, "Updating torrent category", slog.Uint64("id", id), slog.String("name", req.Name))

		var torrentCategory models.TorrentCategory
		result := db.First(&torrentCategory, id)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				slog.WarnContext(ctx, "Torrent category not found", slog.Uint64("id", id))
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Torrent category not found"})
			}
			slog.ErrorContext(ctx, "Failed to query torrent category", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query torrent category"})
		}

		torrentCategory.Name = req.Name
		result = db.Save(&torrentCategory)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to update torrent category", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update torrent category"})
		}

		response := TorrentCategoryResponse{
			ID:   torrentCategory.ID,
			Name: torrentCategory.Name,
		}

		slog.InfoContext(ctx, "Successfully updated torrent category", slog.Uint64("id", id), slog.String("name", torrentCategory.Name))
		return c.JSON(http.StatusOK, response)
	}
}

// DeleteTorrentCategory deletes a torrent category
func DeleteTorrentCategory(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		idStr := c.Param("id")

		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			slog.WarnContext(ctx, "Invalid ID parameter", slog.String("id", idStr), slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}

		slog.InfoContext(ctx, "Deleting torrent category", slog.Uint64("id", id))

		result := db.Delete(&models.TorrentCategory{}, id)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to delete torrent category", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete torrent category"})
		}

		if result.RowsAffected == 0 {
			slog.WarnContext(ctx, "Torrent category not found", slog.Uint64("id", id))
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Torrent category not found"})
		}

		slog.InfoContext(ctx, "Successfully deleted torrent category", slog.Uint64("id", id))
		return c.NoContent(http.StatusNoContent)
	}
}
