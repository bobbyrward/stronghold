package api

import (
	"context"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
)

type NotificationTypeResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type NotificationTypeHandler struct{}

func (h NotificationTypeHandler) ModelToResponse(c echo.Context, ctx context.Context, db *gorm.DB, row models.NotificationType) NotificationTypeResponse {
	return NotificationTypeResponse{
		ID:   row.ID,
		Name: row.Name,
	}
}

func (h NotificationTypeHandler) PreloadRelations(c echo.Context, ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	return db, nil
}

func (h NotificationTypeHandler) IDFromModel(row models.NotificationType) uint {
	return row.ID
}

// ListNotificationTypes returns all notification types
func ListNotificationTypes(db *gorm.DB) echo.HandlerFunc {
	return readOnlyListHandler[models.NotificationType, NotificationTypeResponse](db, NotificationTypeHandler{})
}

// GetNotificationType returns a single notification type by ID
func GetNotificationType(db *gorm.DB) echo.HandlerFunc {
	return readOnlyGetHandler[models.NotificationType, NotificationTypeResponse](db, NotificationTypeHandler{})
}
