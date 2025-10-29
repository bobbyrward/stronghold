package api

import (
	"context"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
)

type FilterKeyRequest struct {
	Name string `json:"name" validate:"required"`
}

type FilterKeyResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type FilterKeyHandler struct{}

func (handler FilterKeyHandler) ModelToResponse(c echo.Context, ctx context.Context, db *gorm.DB, row models.FilterKey) FilterKeyResponse {
	return FilterKeyResponse{
		ID:   row.ID,
		Name: row.Name,
	}
}

func (handler FilterKeyHandler) RequestToModel(c echo.Context, ctx context.Context, db *gorm.DB, req FilterKeyRequest) (models.FilterKey, error) {
	return models.FilterKey{
		Name: req.Name,
	}, nil
}

func (handler FilterKeyHandler) UpdateModel(c echo.Context, ctx context.Context, db *gorm.DB, row *models.FilterKey, req FilterKeyRequest) error {
	row.Name = req.Name
	return nil
}

func (handler FilterKeyHandler) ParseQuery(c echo.Context, ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	return db, nil
}

func (handler FilterKeyHandler) PreloadRelations(c echo.Context, ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	return db, nil
}

func (handler FilterKeyHandler) IDFromModel(row models.FilterKey) uint {
	return row.ID
}

// ListFilterKeys returns all filter keys
func ListFilterKeys(db *gorm.DB) echo.HandlerFunc {
	return genericListHandler[models.FilterKey, FilterKeyRequest, FilterKeyResponse](db, FilterKeyHandler{})
}

// CreateFilterKey creates a new filter key
func CreateFilterKey(db *gorm.DB) echo.HandlerFunc {
	return genericCreateHandler[models.FilterKey, FilterKeyRequest, FilterKeyResponse](db, FilterKeyHandler{})
}

// GetFilterKey returns a single filter key by ID
func GetFilterKey(db *gorm.DB) echo.HandlerFunc {
	return genericGetHandler[models.FilterKey, FilterKeyRequest, FilterKeyResponse](db, FilterKeyHandler{})
}

// UpdateFilterKey updates an existing filter key
func UpdateFilterKey(db *gorm.DB) echo.HandlerFunc {
	return genericUpdateHandler[models.FilterKey, FilterKeyRequest, FilterKeyResponse](db, FilterKeyHandler{})
}

// DeleteFilterKey deletes a filter key
func DeleteFilterKey(db *gorm.DB) echo.HandlerFunc {
	return genericDeleteHandler[models.FilterKey](db)
}
