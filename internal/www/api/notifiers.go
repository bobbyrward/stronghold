package api

import (
	"context"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
)

type NotifierRequest struct {
	Name   string `json:"name" validate:"required"`
	TypeID uint   `json:"type_id" validate:"required"`
	URL    string `json:"url"`
}

type NotifierResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	TypeID   uint   `json:"type_id"`
	TypeName string `json:"type_name"`
	URL      string `json:"url"`
}

type NotifierHandler struct{}

func (handler NotifierHandler) ModelToResponse(c echo.Context, ctx context.Context, db *gorm.DB, row models.Notifier) NotifierResponse {
	typeName := ""
	if row.NotificationType != nil {
		typeName = row.NotificationType.Name
	}
	return NotifierResponse{
		ID:       row.ID,
		Name:     row.Name,
		TypeID:   row.NotificationTypeID,
		TypeName: typeName,
		URL:      row.URL,
	}
}

func (handler NotifierHandler) RequestToModel(c echo.Context, ctx context.Context, db *gorm.DB, req NotifierRequest) (models.Notifier, error) {
	return models.Notifier{
		Name:               req.Name,
		NotificationTypeID: req.TypeID,
		URL:                req.URL,
	}, nil
}

func (handler NotifierHandler) UpdateModel(c echo.Context, ctx context.Context, db *gorm.DB, row *models.Notifier, req NotifierRequest) error {
	row.Name = req.Name
	return nil
}

func (handler NotifierHandler) ParseQuery(c echo.Context, ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	return db, nil
}

func (handler NotifierHandler) PreloadRelations(c echo.Context, ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	return db.Preload("NotificationType"), nil
}

func (handler NotifierHandler) IDFromModel(row models.Notifier) uint {
	return row.ID
}

// ListNotifiers returns all notifiers
func ListNotifiers(db *gorm.DB) echo.HandlerFunc {
	return genericListHandler[models.Notifier, NotifierRequest, NotifierResponse](db, NotifierHandler{})
}

// CreateNotifier creates a new notifier
func CreateNotifier(db *gorm.DB) echo.HandlerFunc {
	return genericCreateHandler[models.Notifier, NotifierRequest, NotifierResponse](db, NotifierHandler{})
}

// GetNotifier returns a single notifier by ID
func GetNotifier(db *gorm.DB) echo.HandlerFunc {
	return genericGetHandler[models.Notifier, NotifierRequest, NotifierResponse](db, NotifierHandler{})
}

// UpdateNotifier updates an existing notifier
func UpdateNotifier(db *gorm.DB) echo.HandlerFunc {
	return genericUpdateHandler[models.Notifier, NotifierRequest, NotifierResponse](db, NotifierHandler{})
}

// DeleteNotifier deletes a notifier
func DeleteNotifier(db *gorm.DB) echo.HandlerFunc {
	return genericDeleteHandler[models.Notifier](db)
}
