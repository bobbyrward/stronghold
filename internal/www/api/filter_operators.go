package api

import (
	"context"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
)

// FilterOperatorRequest is unused but required to satisfy the ModelHandler interface
type FilterOperatorRequest struct{}

type FilterOperatorResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type FilterOperatorHandler struct{}

func (handler FilterOperatorHandler) ModelToResponse(c echo.Context, ctx context.Context, db *gorm.DB, row models.FilterOperator) FilterOperatorResponse {
	return FilterOperatorResponse{
		ID:   row.ID,
		Name: row.Name,
	}
}

// RequestToModel is unused for read-only resources but required by interface
func (handler FilterOperatorHandler) RequestToModel(c echo.Context, ctx context.Context, db *gorm.DB, req FilterOperatorRequest) (models.FilterOperator, error) {
	return models.FilterOperator{}, nil
}

// UpdateModel is unused for read-only resources but required by interface
func (handler FilterOperatorHandler) UpdateModel(c echo.Context, ctx context.Context, db *gorm.DB, row *models.FilterOperator, req FilterOperatorRequest) error {
	return nil
}

func (handler FilterOperatorHandler) ParseQuery(c echo.Context, ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	return db, nil
}

func (handler FilterOperatorHandler) PreloadRelations(c echo.Context, ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	return db, nil
}

func (handler FilterOperatorHandler) IDFromModel(row models.FilterOperator) uint {
	return row.ID
}

// ListFilterOperators returns all filter operators
func ListFilterOperators(db *gorm.DB) echo.HandlerFunc {
	return genericListHandler[models.FilterOperator, FilterOperatorRequest, FilterOperatorResponse](db, FilterOperatorHandler{})
}

// GetFilterOperator returns a single filter operator by ID
func GetFilterOperator(db *gorm.DB) echo.HandlerFunc {
	return genericGetHandler[models.FilterOperator, FilterOperatorRequest, FilterOperatorResponse](db, FilterOperatorHandler{})
}
