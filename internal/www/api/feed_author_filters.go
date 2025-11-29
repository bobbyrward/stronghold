package api

import (
	"context"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
)

type FeedAuthorFilterRequest struct {
	Author     string `json:"author" validate:"required"`
	FeedID     uint   `json:"feed_id" validate:"required"`
	CategoryID uint   `json:"category_id" validate:"required"`
	NotifierID uint   `json:"notifier_id" validate:"required"`
}

type FeedAuthorFilterResponse struct {
	ID           uint   `json:"id"`
	Author       string `json:"author"`
	FeedID       uint   `json:"feed_id"`
	FeedName     string `json:"feed_name"`
	CategoryID   uint   `json:"category_id"`
	CategoryName string `json:"category_name"`
	NotifierID   uint   `json:"notifier_id"`
	NotifierName string `json:"notifier_name"`
}

type FeedAuthorFilterHandler struct{}

func (handler FeedAuthorFilterHandler) ModelToResponse(c echo.Context, ctx context.Context, db *gorm.DB, row models.FeedAuthorFilter) FeedAuthorFilterResponse {
	return FeedAuthorFilterResponse{
		ID:           row.ID,
		Author:       row.Author,
		FeedID:       row.FeedID,
		FeedName:     row.Feed.Name,
		CategoryID:   row.TorrentCategoryID,
		CategoryName: row.TorrentCategory.Name,
		NotifierID:   row.NotifierID,
		NotifierName: row.Notifier.Name,
	}
}

func (handler FeedAuthorFilterHandler) RequestToModel(c echo.Context, ctx context.Context, db *gorm.DB, req FeedAuthorFilterRequest) (models.FeedAuthorFilter, error) {
	return models.FeedAuthorFilter{
		Author:            req.Author,
		FeedID:            req.FeedID,
		TorrentCategoryID: req.CategoryID,
		NotifierID:        req.NotifierID,
	}, nil
}

func (handler FeedAuthorFilterHandler) UpdateModel(c echo.Context, ctx context.Context, db *gorm.DB, row *models.FeedAuthorFilter, req FeedAuthorFilterRequest) error {
	row.Author = req.Author
	row.FeedID = req.FeedID
	row.TorrentCategoryID = req.CategoryID
	row.NotifierID = req.NotifierID
	return nil
}

func (handler FeedAuthorFilterHandler) ParseQuery(c echo.Context, ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	feedID, hasFeedID, err := ParseQueryParamUint(c, ctx, "feed_id")
	if err != nil {
		return db, BadRequest(c, ctx, "Invalid feed_id parameter")
	}

	if hasFeedID {
		db = db.Where("feed_id = ?", feedID)
	}

	return db, nil
}

func (handler FeedAuthorFilterHandler) PreloadRelations(c echo.Context, ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	return db.Preload("Feed").Preload("TorrentCategory").Preload("Notifier"), nil
}

func (handler FeedAuthorFilterHandler) IDFromModel(row models.FeedAuthorFilter) uint {
	return row.ID
}

// ListFeedAuthorFilters returns all feed author filters
func ListFeedAuthorFilters(db *gorm.DB) echo.HandlerFunc {
	return genericListHandler[models.FeedAuthorFilter, FeedAuthorFilterRequest, FeedAuthorFilterResponse](db, FeedAuthorFilterHandler{})
}

// CreateFeedAuthorFilter creates a new feed author filter
func CreateFeedAuthorFilter(db *gorm.DB) echo.HandlerFunc {
	return genericCreateHandler[models.FeedAuthorFilter, FeedAuthorFilterRequest, FeedAuthorFilterResponse](db, FeedAuthorFilterHandler{})
}

// GetFeedAuthorFilter returns a single feed author filter by ID
func GetFeedAuthorFilter(db *gorm.DB) echo.HandlerFunc {
	return genericGetHandler[models.FeedAuthorFilter, FeedAuthorFilterRequest, FeedAuthorFilterResponse](db, FeedAuthorFilterHandler{})
}

// UpdateFeedAuthorFilter updates an existing feed author filter
func UpdateFeedAuthorFilter(db *gorm.DB) echo.HandlerFunc {
	return genericUpdateHandler[models.FeedAuthorFilter, FeedAuthorFilterRequest, FeedAuthorFilterResponse](db, FeedAuthorFilterHandler{})
}

// DeleteFeedAuthorFilter deletes a feed author filter
func DeleteFeedAuthorFilter(db *gorm.DB) echo.HandlerFunc {
	return genericDeleteHandler[models.FeedAuthorFilter](db)
}
