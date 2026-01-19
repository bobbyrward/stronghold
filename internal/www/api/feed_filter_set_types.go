package api

import (
	"context"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
)

type FeedFilterSetTypeResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type FeedFilterSetTypeHandler struct{}

func (h FeedFilterSetTypeHandler) ModelToResponse(c echo.Context, ctx context.Context, db *gorm.DB, row models.FeedFilterSetType) FeedFilterSetTypeResponse {
	return FeedFilterSetTypeResponse{
		ID:   row.ID,
		Name: row.Name,
	}
}

func (h FeedFilterSetTypeHandler) PreloadRelations(c echo.Context, ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	return db, nil
}

func (h FeedFilterSetTypeHandler) IDFromModel(row models.FeedFilterSetType) uint {
	return row.ID
}

// ListFeedFilterSetTypes returns all feed filter set types
func ListFeedFilterSetTypes(db *gorm.DB) echo.HandlerFunc {
	return readOnlyListHandler[models.FeedFilterSetType, FeedFilterSetTypeResponse](db, FeedFilterSetTypeHandler{})
}

// GetFeedFilterSetType returns a single feed filter set type by ID
func GetFeedFilterSetType(db *gorm.DB) echo.HandlerFunc {
	return readOnlyGetHandler[models.FeedFilterSetType, FeedFilterSetTypeResponse](db, FeedFilterSetTypeHandler{})
}
