package api

import (
	"context"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
)

type FeedFilterSetEntryRequest struct {
	FeedFilterSetID uint   `json:"feed_filter_set_id" validate:"required"`
	KeyID           uint   `json:"key_id" validate:"required"`
	OperatorID      uint   `json:"operator_id" validate:"required"`
	Value           string `json:"value" validate:"required"`
}

type FeedFilterSetEntryResponse struct {
	ID              uint   `json:"id"`
	FeedFilterSetID uint   `json:"feed_filter_set_id"`
	KeyID           uint   `json:"key_id"`
	KeyName         string `json:"key_name"`
	OperatorID      uint   `json:"operator_id"`
	OperatorName    string `json:"operator_name"`
	Value           string `json:"value"`
}

type FeedFilterSetEntryHandler struct{}

func (handler FeedFilterSetEntryHandler) ModelToResponse(c echo.Context, ctx context.Context, db *gorm.DB, row models.FeedFilterSetEntry) FeedFilterSetEntryResponse {
	return FeedFilterSetEntryResponse{
		ID:              row.ID,
		FeedFilterSetID: row.FeedFilterSetID,
		KeyID:           row.FilterKeyID,
		KeyName:         row.FilterKey.Name,
		OperatorID:      row.FilterOperatorID,
		OperatorName:    row.FilterOperator.Name,
		Value:           row.Value,
	}
}

func (handler FeedFilterSetEntryHandler) RequestToModel(c echo.Context, ctx context.Context, db *gorm.DB, req FeedFilterSetEntryRequest) (models.FeedFilterSetEntry, error) {
	return models.FeedFilterSetEntry{
		FeedFilterSetID:  req.FeedFilterSetID,
		FilterKeyID:      req.KeyID,
		FilterOperatorID: req.OperatorID,
		Value:            req.Value,
	}, nil
}

func (handler FeedFilterSetEntryHandler) UpdateModel(c echo.Context, ctx context.Context, db *gorm.DB, row *models.FeedFilterSetEntry, req FeedFilterSetEntryRequest) error {
	row.FeedFilterSetID = req.FeedFilterSetID
	row.FilterKeyID = req.KeyID
	row.FilterOperatorID = req.OperatorID
	row.Value = req.Value
	return nil
}

func (handler FeedFilterSetEntryHandler) ParseQuery(c echo.Context, ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	feedFilterSetID, hasFeedFilterSetID, err := ParseQueryParamUint(c, ctx, "feed_filter_set_id")
	if err != nil {
		return db, BadRequest(c, ctx, "Invalid feed_filter_set_id parameter")
	}

	if hasFeedFilterSetID {
		db = db.Where("feed_filter_set_id = ?", feedFilterSetID)
	}

	return db, nil
}

func (handler FeedFilterSetEntryHandler) PreloadRelations(c echo.Context, ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	return db.Preload("FilterKey").Preload("FilterOperator"), nil
}

func (handler FeedFilterSetEntryHandler) IDFromModel(row models.FeedFilterSetEntry) uint {
	return row.ID
}

// ListFeedFilterSetEntries returns all feed filter set entries
func ListFeedFilterSetEntries(db *gorm.DB) echo.HandlerFunc {
	return genericListHandler[models.FeedFilterSetEntry, FeedFilterSetEntryRequest, FeedFilterSetEntryResponse](db, FeedFilterSetEntryHandler{})
}

// CreateFeedFilterSetEntry creates a new feed filter set entry
func CreateFeedFilterSetEntry(db *gorm.DB) echo.HandlerFunc {
	return genericCreateHandler[models.FeedFilterSetEntry, FeedFilterSetEntryRequest, FeedFilterSetEntryResponse](db, FeedFilterSetEntryHandler{})
}

// GetFeedFilterSetEntry returns a single feed filter set entry by ID
func GetFeedFilterSetEntry(db *gorm.DB) echo.HandlerFunc {
	return genericGetHandler[models.FeedFilterSetEntry, FeedFilterSetEntryRequest, FeedFilterSetEntryResponse](db, FeedFilterSetEntryHandler{})
}

// UpdateFeedFilterSetEntry updates an existing feed filter set entry
func UpdateFeedFilterSetEntry(db *gorm.DB) echo.HandlerFunc {
	return genericUpdateHandler[models.FeedFilterSetEntry, FeedFilterSetEntryRequest, FeedFilterSetEntryResponse](db, FeedFilterSetEntryHandler{})
}

// DeleteFeedFilterSetEntry deletes a feed filter set entry
func DeleteFeedFilterSetEntry(db *gorm.DB) echo.HandlerFunc {
	return genericDeleteHandler[models.FeedFilterSetEntry](db)
}
