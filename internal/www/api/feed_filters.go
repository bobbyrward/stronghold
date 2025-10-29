package api

import (
	"context"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
)

type FeedFilterRequest struct {
	Name       string `json:"name" validate:"required"`
	FeedID     uint   `json:"feed_id" validate:"required"`
	CategoryID uint   `json:"category_id" validate:"required"`
	NotifierID uint   `json:"notifier_id" validate:"required"`
}

type FeedFilterResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	FeedID       uint   `json:"feed_id"`
	FeedName     string `json:"feed_name"`
	CategoryID   uint   `json:"category_id"`
	CategoryName string `json:"category_name"`
	NotifierID   uint   `json:"notifier_id"`
	NotifierName string `json:"notifier_name"`
}

type FeedFilterHandler struct{}

func (handler FeedFilterHandler) ModelToResponse(c echo.Context, ctx context.Context, db *gorm.DB, row models.FeedFilter) FeedFilterResponse {
	return FeedFilterResponse{
		ID:           row.ID,
		Name:         row.Name,
		FeedID:       row.FeedID,
		FeedName:     row.Feed.Name,
		CategoryID:   row.TorrentCategoryID,
		CategoryName: row.TorrentCategory.Name,
		NotifierID:   row.NotifierID,
		NotifierName: row.Notifier.Name,
	}
}

func (handler FeedFilterHandler) RequestToModel(c echo.Context, ctx context.Context, db *gorm.DB, req FeedFilterRequest) (models.FeedFilter, error) {
	return models.FeedFilter{
		Name:              req.Name,
		FeedID:            req.FeedID,
		TorrentCategoryID: req.CategoryID,
		NotifierID:        req.NotifierID,
	}, nil
}

func (handler FeedFilterHandler) UpdateModel(c echo.Context, ctx context.Context, db *gorm.DB, row *models.FeedFilter, req FeedFilterRequest) error {
	row.Name = req.Name
	row.TorrentCategoryID = req.CategoryID
	row.NotifierID = req.NotifierID

	return nil
}

func (handler FeedFilterHandler) ParseQuery(c echo.Context, ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	feedID, hasFeedID, err := ParseQueryParamUint(c, ctx, "feed_id")
	if err != nil {
		return db, BadRequest(c, ctx, "Invalid feed_id parameter")
	}

	if hasFeedID {
		db = db.Where("feed_id = ?", feedID)
	}

	return db, nil
}

func (handler FeedFilterHandler) PreloadRelations(c echo.Context, ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	return db.Preload("Feed").Preload("TorrentCategory").Preload("Notifier"), nil
}

func (handler FeedFilterHandler) IDFromModel(row models.FeedFilter) uint {
	return row.ID
}

// ListFeedFilters returns all feed filters
func ListFeedFilters(db *gorm.DB) echo.HandlerFunc {
	return genericListHandler[models.FeedFilter, FeedFilterRequest, FeedFilterResponse](db, FeedFilterHandler{})
}

// CreateFeedFilter creates a new feed filter
func CreateFeedFilter(db *gorm.DB) echo.HandlerFunc {
	return genericCreateHandler[models.FeedFilter, FeedFilterRequest, FeedFilterResponse](db, FeedFilterHandler{})
}

// GetFeedFilter returns a single feed filter by ID
func GetFeedFilter(db *gorm.DB) echo.HandlerFunc {
	return genericGetHandler[models.FeedFilter, FeedFilterRequest, FeedFilterResponse](db, FeedFilterHandler{})
}

// UpdateFeedFilter updates an existing feed filter
func UpdateFeedFilter(db *gorm.DB) echo.HandlerFunc {
	return genericUpdateHandler[models.FeedFilter, FeedFilterRequest, FeedFilterResponse](db, FeedFilterHandler{})
}

// DeleteFeedFilter deletes a feed filter
func DeleteFeedFilter(db *gorm.DB) echo.HandlerFunc {
	return genericDeleteHandler[models.FeedFilter](db)
}
