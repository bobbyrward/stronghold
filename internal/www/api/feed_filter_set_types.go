package api

import (
	"context"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
)

type FeedFilterSetTypeRequest struct {
	Name string `json:"name" validate:"required"`
}

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

func (handler FeedFilterSetTypeHandler) RequestToModel(c echo.Context, ctx context.Context, db *gorm.DB, req FeedFilterSetTypeRequest) (models.FeedFilterSetType, error) {
	return models.FeedFilterSetType{
		Name: req.Name,
	}, nil
}

func (handler FeedFilterSetTypeHandler) UpdateModel(c echo.Context, ctx context.Context, db *gorm.DB, row *models.FeedFilterSetType, req FeedFilterSetTypeRequest) error {
	row.Name = req.Name
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

// CreateFeedFilterSetType creates a new feed filter set type
func CreateFeedFilterSetType(db *gorm.DB) echo.HandlerFunc {
	return genericCreateHandler[models.FeedFilterSetType, FeedFilterSetTypeRequest, FeedFilterSetTypeResponse](db, FeedFilterSetTypeHandler{})
}

// GetFeedFilterSetType returns a single feed filter set type by ID
func GetFeedFilterSetType(db *gorm.DB) echo.HandlerFunc {
	return genericGetHandler[models.FeedFilterSetType, FeedFilterSetTypeRequest, FeedFilterSetTypeResponse](db, FeedFilterSetTypeHandler{})
}

// UpdateFeedFilterSetType updates an existing feed filter set type
func UpdateFeedFilterSetType(db *gorm.DB) echo.HandlerFunc {
	return genericUpdateHandler[models.FeedFilterSetType, FeedFilterSetTypeRequest, FeedFilterSetTypeResponse](db, FeedFilterSetTypeHandler{})
}

// DeleteFeedFilterSetType deletes a feed filter set type
func DeleteFeedFilterSetType(db *gorm.DB) echo.HandlerFunc {
	return genericDeleteHandler[models.FeedFilterSetType](db)
}
