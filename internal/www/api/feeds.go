package api

import (
	"context"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
)

type FeedRequest struct {
	Name string `json:"name" validate:"required"`
	URL  string `json:"url" validate:"required"`
}

type FeedResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type FeedHandler struct{}

func (handler FeedHandler) ModelToResponse(c echo.Context, ctx context.Context, db *gorm.DB, row models.Feed) FeedResponse {
	return FeedResponse{
		ID:   row.ID,
		Name: row.Name,
		URL:  row.URL,
	}
}

func (handler FeedHandler) RequestToModel(c echo.Context, ctx context.Context, db *gorm.DB, req FeedRequest) (models.Feed, error) {
	return models.Feed{
		Name: req.Name,
		URL:  req.URL,
	}, nil
}

func (handler FeedHandler) UpdateModel(c echo.Context, ctx context.Context, db *gorm.DB, row *models.Feed, req FeedRequest) error {
	row.Name = req.Name
	row.URL = req.URL
	return nil
}

func (handler FeedHandler) ParseQuery(c echo.Context, ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	return db, nil
}

func (handler FeedHandler) PreloadRelations(c echo.Context, ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	return db, nil
}

func (handler FeedHandler) IDFromModel(row models.Feed) uint {
	return row.ID
}

// ListFeeds returns all feeds
func ListFeeds(db *gorm.DB) echo.HandlerFunc {
	return genericListHandler[models.Feed, FeedRequest, FeedResponse](db, FeedHandler{})
}

// CreateFeed creates a new feed
func CreateFeed(db *gorm.DB) echo.HandlerFunc {
	return genericCreateHandler[models.Feed, FeedRequest, FeedResponse](db, FeedHandler{})
}

// GetFeed returns a single feed by ID
func GetFeed(db *gorm.DB) echo.HandlerFunc {
	return genericGetHandler[models.Feed, FeedRequest, FeedResponse](db, FeedHandler{})
}

// UpdateFeed updates an existing feed
func UpdateFeed(db *gorm.DB) echo.HandlerFunc {
	return genericUpdateHandler[models.Feed, FeedRequest, FeedResponse](db, FeedHandler{})
}

// DeleteFeed deletes a feed
func DeleteFeed(db *gorm.DB) echo.HandlerFunc {
	return genericDeleteHandler[models.Feed](db)
}
