package api

import (
	"context"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
)

type FilterOperatorResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type FilterOperatorHandler struct{}

func (h FilterOperatorHandler) ModelToResponse(c echo.Context, ctx context.Context, db *gorm.DB, row models.FilterOperator) FilterOperatorResponse {
	return FilterOperatorResponse{
		ID:   row.ID,
		Name: row.Name,
	}
}

func (h FilterOperatorHandler) PreloadRelations(c echo.Context, ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	return db, nil
}

func (h FilterOperatorHandler) IDFromModel(row models.FilterOperator) uint {
	return row.ID
}

// ListFilterOperators returns all filter operators
func ListFilterOperators(db *gorm.DB) echo.HandlerFunc {
	return readOnlyListHandler[models.FilterOperator, FilterOperatorResponse](db, FilterOperatorHandler{})
}

// GetFilterOperator returns a single filter operator by ID
func GetFilterOperator(db *gorm.DB) echo.HandlerFunc {
	return readOnlyGetHandler[models.FilterOperator, FilterOperatorResponse](db, FilterOperatorHandler{})
}
