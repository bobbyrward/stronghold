package api

import (
	"context"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
)

type TorrentCategoryResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	ScopeID   uint   `json:"scope_id"`
	ScopeName string `json:"scope_name"`
	MediaType string `json:"media_type"`
}

type TorrentCategoryHandler struct{}

func (h TorrentCategoryHandler) ModelToResponse(c echo.Context, ctx context.Context, db *gorm.DB, row models.TorrentCategory) TorrentCategoryResponse {
	return TorrentCategoryResponse{
		ID:        row.ID,
		Name:      row.Name,
		ScopeID:   row.ScopeID,
		ScopeName: row.Scope.Name,
		MediaType: row.MediaType,
	}
}

func (h TorrentCategoryHandler) PreloadRelations(c echo.Context, ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	return db.Preload("Scope"), nil
}

func (h TorrentCategoryHandler) IDFromModel(row models.TorrentCategory) uint {
	return row.ID
}

// ListTorrentCategories returns all torrent categories
func ListTorrentCategories(db *gorm.DB) echo.HandlerFunc {
	return readOnlyListHandler[models.TorrentCategory, TorrentCategoryResponse](db, TorrentCategoryHandler{})
}

// GetTorrentCategory returns a single torrent category by ID
func GetTorrentCategory(db *gorm.DB) echo.HandlerFunc {
	return readOnlyGetHandler[models.TorrentCategory, TorrentCategoryResponse](db, TorrentCategoryHandler{})
}
