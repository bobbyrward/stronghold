package api

import (
	"context"

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

type NotificationTypeHandler struct{}

func (handler NotificationTypeHandler) ModelToResponse(c echo.Context, ctx context.Context, db *gorm.DB, row models.NotificationType) NotificationTypeResponse {
	return NotificationTypeResponse{
		ID:   row.ID,
		Name: row.Name,
	}
}

func (handler NotificationTypeHandler) RequestToModel(c echo.Context, ctx context.Context, db *gorm.DB, req NotificationTypeRequest) (models.NotificationType, error) {
	return models.NotificationType{
		Name: req.Name,
	}, nil
}

func (handler NotificationTypeHandler) UpdateModel(c echo.Context, ctx context.Context, db *gorm.DB, row *models.NotificationType, req NotificationTypeRequest) error {
	row.Name = req.Name
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

// CreateNotificationType creates a new notification type
func CreateNotificationType(db *gorm.DB) echo.HandlerFunc {
	return genericCreateHandler[models.NotificationType, NotificationTypeRequest, NotificationTypeResponse](db, NotificationTypeHandler{})
}

// GetNotificationType returns a single notification type by ID
func GetNotificationType(db *gorm.DB) echo.HandlerFunc {
	return genericGetHandler[models.NotificationType, NotificationTypeRequest, NotificationTypeResponse](db, NotificationTypeHandler{})
}

// UpdateNotificationType updates an existing notification type
func UpdateNotificationType(db *gorm.DB) echo.HandlerFunc {
	return genericUpdateHandler[models.NotificationType, NotificationTypeRequest, NotificationTypeResponse](db, NotificationTypeHandler{})
}

// DeleteNotificationType deletes a notification type
func DeleteNotificationType(db *gorm.DB) echo.HandlerFunc {
	return genericDeleteHandler[models.NotificationType](db)
}
