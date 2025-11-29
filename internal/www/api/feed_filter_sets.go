package api

import (
	"context"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
)

type FeedFilterSetRequest struct {
	FeedFilterID        uint `json:"feed_filter_id" validate:"required"`
	FeedFilterSetTypeID uint `json:"type_id" validate:"required"`
}

type FeedFilterSetResponse struct {
	ID           uint   `json:"id"`
	FeedFilterID uint   `json:"feed_filter_id"`
	TypeName     string `json:"type_name"`
	TypeID       uint   `json:"type_id"`
}

type FeedFilterSetHandler struct{}

func (handler FeedFilterSetHandler) ModelToResponse(c echo.Context, ctx context.Context, db *gorm.DB, row models.FeedFilterSet) FeedFilterSetResponse {
	return FeedFilterSetResponse{
		ID:           row.ID,
		FeedFilterID: row.FeedFilterID,
		TypeName:     row.FeedFilterSetType.Name,
		TypeID:       row.FeedFilterSetType.ID,
	}
}

func (handler FeedFilterSetHandler) RequestToModel(c echo.Context, ctx context.Context, db *gorm.DB, req FeedFilterSetRequest) (models.FeedFilterSet, error) {
	return models.FeedFilterSet{
		FeedFilterID:        req.FeedFilterID,
		FeedFilterSetTypeID: req.FeedFilterSetTypeID,
	}, nil
}

func (handler FeedFilterSetHandler) UpdateModel(c echo.Context, ctx context.Context, db *gorm.DB, row *models.FeedFilterSet, req FeedFilterSetRequest) error {
	row.FeedFilterID = req.FeedFilterID
	row.FeedFilterSetTypeID = req.FeedFilterSetTypeID
	return nil
}

func (handler FeedFilterSetHandler) ParseQuery(c echo.Context, ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	feedFilterID, hasFeedFilterID, err := ParseQueryParamUint(c, ctx, "feed_filter_id")
	if err != nil {
		return db, BadRequest(c, ctx, "Invalid feed_filter_id parameter")
	}

	if hasFeedFilterID {
		db = db.Where("feed_filter_id = ?", feedFilterID)
	}

	return db, nil
}

func (handler FeedFilterSetHandler) PreloadRelations(c echo.Context, ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	return db.Preload("FeedFilterSetType"), nil
}

func (handler FeedFilterSetHandler) IDFromModel(row models.FeedFilterSet) uint {
	return row.ID
}

// ListFeedFilterSets returns all feed filter sets
func ListFeedFilterSets(db *gorm.DB) echo.HandlerFunc {
	return genericListHandler[models.FeedFilterSet, FeedFilterSetRequest, FeedFilterSetResponse](db, FeedFilterSetHandler{})
}

// CreateFeedFilterSet creates a new feed filter set
func CreateFeedFilterSet(db *gorm.DB) echo.HandlerFunc {
	return genericCreateHandler[models.FeedFilterSet, FeedFilterSetRequest, FeedFilterSetResponse](db, FeedFilterSetHandler{})
}

// GetFeedFilterSet returns a single feed filter set by ID
func GetFeedFilterSet(db *gorm.DB) echo.HandlerFunc {
	return genericGetHandler[models.FeedFilterSet, FeedFilterSetRequest, FeedFilterSetResponse](db, FeedFilterSetHandler{})
}

// UpdateFeedFilterSet updates an existing feed filter set
func UpdateFeedFilterSet(db *gorm.DB) echo.HandlerFunc {
	return genericUpdateHandler[models.FeedFilterSet, FeedFilterSetRequest, FeedFilterSetResponse](db, FeedFilterSetHandler{})
}

// DeleteFeedFilterSet deletes a feed filter set
func DeleteFeedFilterSet(db *gorm.DB) echo.HandlerFunc {
	return genericDeleteHandler[models.FeedFilterSet](db)
}
