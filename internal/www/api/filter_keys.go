package api

import (
	"context"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
)

type FilterKeyResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type FilterKeyHandler struct{}

func (h FilterKeyHandler) ModelToResponse(c echo.Context, ctx context.Context, db *gorm.DB, row models.FilterKey) FilterKeyResponse {
	return FilterKeyResponse{
		ID:   row.ID,
		Name: row.Name,
	}
}

func (h FilterKeyHandler) PreloadRelations(c echo.Context, ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	return db, nil
}

func (h FilterKeyHandler) IDFromModel(row models.FilterKey) uint {
	return row.ID
}

// ListFilterKeys returns all filter keys
func ListFilterKeys(db *gorm.DB) echo.HandlerFunc {
	return readOnlyListHandler[models.FilterKey, FilterKeyResponse](db, FilterKeyHandler{})
}

// GetFilterKey returns a single filter key by ID
func GetFilterKey(db *gorm.DB) echo.HandlerFunc {
	return readOnlyGetHandler[models.FilterKey, FilterKeyResponse](db, FilterKeyHandler{})
}
