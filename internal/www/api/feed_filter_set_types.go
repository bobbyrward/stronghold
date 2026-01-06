package api

import (
	"context"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
)

// FeedFilterSetTypeRequest is unused but required to satisfy the ModelHandler interface
type FeedFilterSetTypeRequest struct{}

type FeedFilterSetTypeResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type FeedFilterSetTypeHandler struct{}

func (handler FeedFilterSetTypeHandler) ModelToResponse(c echo.Context, ctx context.Context, db *gorm.DB, row models.FeedFilterSetType) FeedFilterSetTypeResponse {
	return FeedFilterSetTypeResponse{
		ID:   row.ID,
		Name: row.Name,
	}
}

// RequestToModel is unused for read-only resources but required by interface
func (handler FeedFilterSetTypeHandler) RequestToModel(c echo.Context, ctx context.Context, db *gorm.DB, req FeedFilterSetTypeRequest) (models.FeedFilterSetType, error) {
	return models.FeedFilterSetType{}, nil
}

// UpdateModel is unused for read-only resources but required by interface
func (handler FeedFilterSetTypeHandler) UpdateModel(c echo.Context, ctx context.Context, db *gorm.DB, row *models.FeedFilterSetType, req FeedFilterSetTypeRequest) error {
	return nil
}

func (handler FeedFilterSetTypeHandler) ParseQuery(c echo.Context, ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	return db, nil
}

func (handler FeedFilterSetTypeHandler) PreloadRelations(c echo.Context, ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	return db, nil
}

func (handler FeedFilterSetTypeHandler) IDFromModel(row models.FeedFilterSetType) uint {
	return row.ID
}

// ListFeedFilterSetTypes returns all feed filter set types
func ListFeedFilterSetTypes(db *gorm.DB) echo.HandlerFunc {
	return genericListHandler[models.FeedFilterSetType, FeedFilterSetTypeRequest, FeedFilterSetTypeResponse](db, FeedFilterSetTypeHandler{})
}

// GetFeedFilterSetType returns a single feed filter set type by ID
func GetFeedFilterSetType(db *gorm.DB) echo.HandlerFunc {
	return genericGetHandler[models.FeedFilterSetType, FeedFilterSetTypeRequest, FeedFilterSetTypeResponse](db, FeedFilterSetTypeHandler{})
}
