package api

import (
	"context"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
)

// TorrentCategoryRequest is unused but required to satisfy the ModelHandler interface
type TorrentCategoryRequest struct{}

type TorrentCategoryResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	ScopeID   uint   `json:"scope_id"`
	ScopeName string `json:"scope_name"`
	MediaType string `json:"media_type"`
}

type TorrentCategoryHandler struct{}

func (handler TorrentCategoryHandler) ModelToResponse(c echo.Context, ctx context.Context, db *gorm.DB, row models.TorrentCategory) TorrentCategoryResponse {
	return TorrentCategoryResponse{
		ID:        row.ID,
		Name:      row.Name,
		ScopeID:   row.ScopeID,
		ScopeName: row.Scope.Name,
		MediaType: row.MediaType,
	}
}

// RequestToModel is unused for read-only resources but required by interface
func (handler TorrentCategoryHandler) RequestToModel(c echo.Context, ctx context.Context, db *gorm.DB, req TorrentCategoryRequest) (models.TorrentCategory, error) {
	return models.TorrentCategory{}, nil
}

// UpdateModel is unused for read-only resources but required by interface
func (handler TorrentCategoryHandler) UpdateModel(c echo.Context, ctx context.Context, db *gorm.DB, row *models.TorrentCategory, req TorrentCategoryRequest) error {
	return nil
}

func (handler TorrentCategoryHandler) ParseQuery(c echo.Context, ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	return db, nil
}

func (handler TorrentCategoryHandler) PreloadRelations(c echo.Context, ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	return db.Preload("Scope"), nil
}

func (handler TorrentCategoryHandler) IDFromModel(row models.TorrentCategory) uint {
	return row.ID
}

// ListTorrentCategories returns all torrent categories
func ListTorrentCategories(db *gorm.DB) echo.HandlerFunc {
	return genericListHandler[models.TorrentCategory, TorrentCategoryRequest, TorrentCategoryResponse](db, TorrentCategoryHandler{})
}

// GetTorrentCategory returns a single torrent category by ID
func GetTorrentCategory(db *gorm.DB) echo.HandlerFunc {
	return genericGetHandler[models.TorrentCategory, TorrentCategoryRequest, TorrentCategoryResponse](db, TorrentCategoryHandler{})
}
