package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
)

type NotifierRequest struct {
	Name     string `json:"name" validate:"required"`
	TypeName string `json:"type_name" validate:"required"`
	URL      string `json:"url"`
}

type NotifierResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	URL  string `json:"url"`
}

// ListNotifiers returns all notifiers
func ListNotifiers(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		slog.InfoContext(ctx, "Listing notifiers")

		var notifiers []models.Notifier
		result := db.Preload("NotificationType").Find(&notifiers)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to query notifiers", slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query notifiers"})
		}

		response := make([]NotifierResponse, len(notifiers))
		for i, n := range notifiers {
			typeName := ""
			if n.NotificationType != nil {
				typeName = n.NotificationType.Name
			}
			response[i] = NotifierResponse{
				ID:   n.ID,
				Name: n.Name,
				Type: typeName,
				URL:  n.URL,
			}
		}

		slog.InfoContext(ctx, "Successfully listed notifiers", slog.Int("count", len(response)))
		return c.JSON(http.StatusOK, response)
	}
}

// CreateNotifier creates a new notifier
func CreateNotifier(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var req NotifierRequest
		if err := c.Bind(&req); err != nil {
			slog.WarnContext(ctx, "Invalid request body", slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		}

		if req.Name == "" {
			slog.WarnContext(ctx, "Validation failed: name is required")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Name is required"})
		}

		if req.TypeName == "" {
			slog.WarnContext(ctx, "Validation failed: type_name is required")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Type name is required"})
		}

		slog.InfoContext(ctx, "Creating notifier", slog.String("name", req.Name), slog.String("type_name", req.TypeName))

		// Look up notification type by name
		var notificationType models.NotificationType
		result := db.Where("name = ?", req.TypeName).First(&notificationType)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				slog.WarnContext(ctx, "Notification type not found", slog.String("type_name", req.TypeName))
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Notification type not found"})
			}
			slog.ErrorContext(ctx, "Failed to query notification type", slog.String("type_name", req.TypeName), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query notification type"})
		}

		notifier := models.Notifier{
			Name:               req.Name,
			NotificationTypeID: notificationType.ID,
			URL:                req.URL,
		}

		result = db.Create(&notifier)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to create notifier", slog.String("name", req.Name), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create notifier"})
		}

		response := NotifierResponse{
			ID:   notifier.ID,
			Name: notifier.Name,
			Type: req.TypeName,
			URL:  notifier.URL,
		}

		slog.InfoContext(ctx, "Successfully created notifier", slog.Uint64("id", uint64(notifier.ID)), slog.String("name", notifier.Name))
		return c.JSON(http.StatusCreated, response)
	}
}

// GetNotifier returns a single notifier by ID
func GetNotifier(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		idStr := c.Param("id")

		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			slog.WarnContext(ctx, "Invalid ID parameter", slog.String("id", idStr), slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}

		slog.InfoContext(ctx, "Getting notifier", slog.Uint64("id", id))

		var notifier models.Notifier
		result := db.Preload("NotificationType").First(&notifier, id)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				slog.WarnContext(ctx, "Notifier not found", slog.Uint64("id", id))
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Notifier not found"})
			}
			slog.ErrorContext(ctx, "Failed to query notifier", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query notifier"})
		}

		typeName := ""
		if notifier.NotificationType != nil {
			typeName = notifier.NotificationType.Name
		}

		response := NotifierResponse{
			ID:   notifier.ID,
			Name: notifier.Name,
			Type: typeName,
			URL:  notifier.URL,
		}

		slog.InfoContext(ctx, "Successfully retrieved notifier", slog.Uint64("id", id), slog.String("name", notifier.Name))
		return c.JSON(http.StatusOK, response)
	}
}

// UpdateNotifier updates an existing notifier
func UpdateNotifier(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		idStr := c.Param("id")

		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			slog.WarnContext(ctx, "Invalid ID parameter", slog.String("id", idStr), slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}

		var req NotifierRequest
		if err := c.Bind(&req); err != nil {
			slog.WarnContext(ctx, "Invalid request body", slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		}

		if req.Name == "" {
			slog.WarnContext(ctx, "Validation failed: name is required")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Name is required"})
		}

		if req.TypeName == "" {
			slog.WarnContext(ctx, "Validation failed: type_name is required")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Type name is required"})
		}

		slog.InfoContext(ctx, "Updating notifier", slog.Uint64("id", id), slog.String("name", req.Name))

		var notifier models.Notifier
		result := db.First(&notifier, id)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				slog.WarnContext(ctx, "Notifier not found", slog.Uint64("id", id))
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Notifier not found"})
			}
			slog.ErrorContext(ctx, "Failed to query notifier", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query notifier"})
		}

		// Look up notification type by name
		var notificationType models.NotificationType
		result = db.Where("name = ?", req.TypeName).First(&notificationType)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				slog.WarnContext(ctx, "Notification type not found", slog.String("type_name", req.TypeName))
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Notification type not found"})
			}
			slog.ErrorContext(ctx, "Failed to query notification type", slog.String("type_name", req.TypeName), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query notification type"})
		}

		notifier.Name = req.Name
		notifier.NotificationTypeID = notificationType.ID
		notifier.URL = req.URL

		result = db.Save(&notifier)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to update notifier", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update notifier"})
		}

		response := NotifierResponse{
			ID:   notifier.ID,
			Name: notifier.Name,
			Type: req.TypeName,
			URL:  notifier.URL,
		}

		slog.InfoContext(ctx, "Successfully updated notifier", slog.Uint64("id", id), slog.String("name", notifier.Name))
		return c.JSON(http.StatusOK, response)
	}
}

// DeleteNotifier deletes a notifier
func DeleteNotifier(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		idStr := c.Param("id")

		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			slog.WarnContext(ctx, "Invalid ID parameter", slog.String("id", idStr), slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}

		slog.InfoContext(ctx, "Deleting notifier", slog.Uint64("id", id))

		result := db.Delete(&models.Notifier{}, id)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to delete notifier", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete notifier"})
		}

		if result.RowsAffected == 0 {
			slog.WarnContext(ctx, "Notifier not found", slog.Uint64("id", id))
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Notifier not found"})
		}

		slog.InfoContext(ctx, "Successfully deleted notifier", slog.Uint64("id", id))
		return c.NoContent(http.StatusNoContent)
	}
}
