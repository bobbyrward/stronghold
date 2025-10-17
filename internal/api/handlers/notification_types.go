package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
)

type NotificationTypeRequest struct {
	Name string `json:"name" validate:"required"`
}

type NotificationTypeResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// ListNotificationTypes returns all notification types
func ListNotificationTypes(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		slog.InfoContext(ctx, "Listing notification types")

		var notificationTypes []models.NotificationType
		result := db.Find(&notificationTypes)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to query notification types", slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query notification types"})
		}

		response := make([]NotificationTypeResponse, len(notificationTypes))
		for i, nt := range notificationTypes {
			response[i] = NotificationTypeResponse{
				ID:   nt.ID,
				Name: nt.Name,
			}
		}

		slog.InfoContext(ctx, "Successfully listed notification types", slog.Int("count", len(response)))
		return c.JSON(http.StatusOK, response)
	}
}

// CreateNotificationType creates a new notification type
func CreateNotificationType(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var req NotificationTypeRequest
		if err := c.Bind(&req); err != nil {
			slog.WarnContext(ctx, "Invalid request body", slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		}

		if req.Name == "" {
			slog.WarnContext(ctx, "Validation failed: name is required")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Name is required"})
		}

		slog.InfoContext(ctx, "Creating notification type", slog.String("name", req.Name))

		notificationType := models.NotificationType{
			Name: req.Name,
		}

		result := db.Create(&notificationType)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to create notification type", slog.String("name", req.Name), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create notification type"})
		}

		response := NotificationTypeResponse{
			ID:   notificationType.ID,
			Name: notificationType.Name,
		}

		slog.InfoContext(ctx, "Successfully created notification type", slog.Uint64("id", uint64(notificationType.ID)), slog.String("name", notificationType.Name))
		return c.JSON(http.StatusCreated, response)
	}
}

// GetNotificationType returns a single notification type by ID
func GetNotificationType(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		idStr := c.Param("id")

		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			slog.WarnContext(ctx, "Invalid ID parameter", slog.String("id", idStr), slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}

		slog.InfoContext(ctx, "Getting notification type", slog.Uint64("id", id))

		var notificationType models.NotificationType
		result := db.First(&notificationType, id)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				slog.WarnContext(ctx, "Notification type not found", slog.Uint64("id", id))
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Notification type not found"})
			}
			slog.ErrorContext(ctx, "Failed to query notification type", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query notification type"})
		}

		response := NotificationTypeResponse{
			ID:   notificationType.ID,
			Name: notificationType.Name,
		}

		slog.InfoContext(ctx, "Successfully retrieved notification type", slog.Uint64("id", id), slog.String("name", notificationType.Name))
		return c.JSON(http.StatusOK, response)
	}
}

// UpdateNotificationType updates an existing notification type
func UpdateNotificationType(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		idStr := c.Param("id")

		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			slog.WarnContext(ctx, "Invalid ID parameter", slog.String("id", idStr), slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}

		var req NotificationTypeRequest
		if err := c.Bind(&req); err != nil {
			slog.WarnContext(ctx, "Invalid request body", slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		}

		if req.Name == "" {
			slog.WarnContext(ctx, "Validation failed: name is required")
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Name is required"})
		}

		slog.InfoContext(ctx, "Updating notification type", slog.Uint64("id", id), slog.String("name", req.Name))

		var notificationType models.NotificationType
		result := db.First(&notificationType, id)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				slog.WarnContext(ctx, "Notification type not found", slog.Uint64("id", id))
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Notification type not found"})
			}
			slog.ErrorContext(ctx, "Failed to query notification type", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query notification type"})
		}

		notificationType.Name = req.Name
		result = db.Save(&notificationType)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to update notification type", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update notification type"})
		}

		response := NotificationTypeResponse{
			ID:   notificationType.ID,
			Name: notificationType.Name,
		}

		slog.InfoContext(ctx, "Successfully updated notification type", slog.Uint64("id", id), slog.String("name", notificationType.Name))
		return c.JSON(http.StatusOK, response)
	}
}

// DeleteNotificationType deletes a notification type
func DeleteNotificationType(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		idStr := c.Param("id")

		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			slog.WarnContext(ctx, "Invalid ID parameter", slog.String("id", idStr), slog.Any("error", err))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}

		slog.InfoContext(ctx, "Deleting notification type", slog.Uint64("id", id))

		result := db.Delete(&models.NotificationType{}, id)
		if result.Error != nil {
			slog.ErrorContext(ctx, "Failed to delete notification type", slog.Uint64("id", id), slog.Any("error", result.Error))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete notification type"})
		}

		if result.RowsAffected == 0 {
			slog.WarnContext(ctx, "Notification type not found", slog.Uint64("id", id))
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Notification type not found"})
		}

		slog.InfoContext(ctx, "Successfully deleted notification type", slog.Uint64("id", id))
		return c.NoContent(http.StatusNoContent)
	}
}
