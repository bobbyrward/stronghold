package api

import (
	"context"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
)

type TorrentCategoryRequest struct {
	Name string `json:"name" validate:"required"`
}

type TorrentCategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type TorrentCategoryHandler struct{}

func (handler TorrentCategoryHandler) ModelToResponse(c echo.Context, ctx context.Context, db *gorm.DB, row models.TorrentCategory) TorrentCategoryResponse {
	return TorrentCategoryResponse{
		ID:   row.ID,
		Name: row.Name,
	}
}

func (handler TorrentCategoryHandler) RequestToModel(c echo.Context, ctx context.Context, db *gorm.DB, req TorrentCategoryRequest) (models.TorrentCategory, error) {
	return models.TorrentCategory{
		Name: req.Name,
	}, nil
}

func (handler TorrentCategoryHandler) UpdateModel(c echo.Context, ctx context.Context, db *gorm.DB, row *models.TorrentCategory, req TorrentCategoryRequest) error {
	row.Name = req.Name
	return nil
}

func (handler TorrentCategoryHandler) ParseQuery(c echo.Context, ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	return db, nil
}

func (handler TorrentCategoryHandler) PreloadRelations(c echo.Context, ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	return db, nil
}

func (handler TorrentCategoryHandler) IDFromModel(row models.TorrentCategory) uint {
	return row.ID
}

// ListTorrentCategories returns all torrent categories
func ListTorrentCategories(db *gorm.DB) echo.HandlerFunc {
	return genericListHandler[models.TorrentCategory, TorrentCategoryRequest, TorrentCategoryResponse](db, TorrentCategoryHandler{})
}

// CreateTorrentCategory creates a new torrent category
func CreateTorrentCategory(db *gorm.DB) echo.HandlerFunc {
	return genericCreateHandler[models.TorrentCategory, TorrentCategoryRequest, TorrentCategoryResponse](db, TorrentCategoryHandler{})
}

// GetTorrentCategory returns a single torrent category by ID
func GetTorrentCategory(db *gorm.DB) echo.HandlerFunc {
	return genericGetHandler[models.TorrentCategory, TorrentCategoryRequest, TorrentCategoryResponse](db, TorrentCategoryHandler{})
}

// UpdateTorrentCategory updates an existing torrent category
func UpdateTorrentCategory(db *gorm.DB) echo.HandlerFunc {
	return genericUpdateHandler[models.TorrentCategory, TorrentCategoryRequest, TorrentCategoryResponse](db, TorrentCategoryHandler{})
}

// DeleteTorrentCategory deletes a torrent category
func DeleteTorrentCategory(db *gorm.DB) echo.HandlerFunc {
	return genericDeleteHandler[models.TorrentCategory](db)
}
