package api

import (
	"context"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
)

// NotificationTypeRequest is unused but required to satisfy the ModelHandler interface
type NotificationTypeRequest struct{}

type NotificationTypeResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type NotificationTypeHandler struct{}

func (handler NotificationTypeHandler) ModelToResponse(c echo.Context, ctx context.Context, db *gorm.DB, row models.NotificationType) NotificationTypeResponse {
	return NotificationTypeResponse{
		ID:   row.ID,
		Name: row.Name,
	}
}

// RequestToModel is unused for read-only resources but required by interface
func (handler NotificationTypeHandler) RequestToModel(c echo.Context, ctx context.Context, db *gorm.DB, req NotificationTypeRequest) (models.NotificationType, error) {
	return models.NotificationType{}, nil
}

// UpdateModel is unused for read-only resources but required by interface
func (handler NotificationTypeHandler) UpdateModel(c echo.Context, ctx context.Context, db *gorm.DB, row *models.NotificationType, req NotificationTypeRequest) error {
	return nil
}

func (handler NotificationTypeHandler) ParseQuery(c echo.Context, ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	return db, nil
}

func (handler NotificationTypeHandler) PreloadRelations(c echo.Context, ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	return db, nil
}

func (handler NotificationTypeHandler) IDFromModel(row models.NotificationType) uint {
	return row.ID
}

// ListNotificationTypes returns all notification types
func ListNotificationTypes(db *gorm.DB) echo.HandlerFunc {
	return genericListHandler[models.NotificationType, NotificationTypeRequest, NotificationTypeResponse](db, NotificationTypeHandler{})
}

// GetNotificationType returns a single notification type by ID
func GetNotificationType(db *gorm.DB) echo.HandlerFunc {
	return genericGetHandler[models.NotificationType, NotificationTypeRequest, NotificationTypeResponse](db, NotificationTypeHandler{})
}
