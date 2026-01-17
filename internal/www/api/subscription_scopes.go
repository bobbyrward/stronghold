package api

import (
	"context"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
)

// SubscriptionScopeRequest is unused but required to satisfy the ModelHandler interface
type SubscriptionScopeRequest struct{}

type SubscriptionScopeResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type SubscriptionScopeHandler struct{}

func (handler SubscriptionScopeHandler) ModelToResponse(c echo.Context, ctx context.Context, db *gorm.DB, row models.SubscriptionScope) SubscriptionScopeResponse {
	return SubscriptionScopeResponse{
		ID:   row.ID,
		Name: row.Name,
	}
}

// RequestToModel is unused for read-only resources but required by interface
func (handler SubscriptionScopeHandler) RequestToModel(c echo.Context, ctx context.Context, db *gorm.DB, req SubscriptionScopeRequest) (models.SubscriptionScope, error) {
	return models.SubscriptionScope{}, nil
}

// UpdateModel is unused for read-only resources but required by interface
func (handler SubscriptionScopeHandler) UpdateModel(c echo.Context, ctx context.Context, db *gorm.DB, row *models.SubscriptionScope, req SubscriptionScopeRequest) error {
	return nil
}

func (handler SubscriptionScopeHandler) ParseQuery(c echo.Context, ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	return db, nil
}

func (handler SubscriptionScopeHandler) PreloadRelations(c echo.Context, ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	return db, nil
}

func (handler SubscriptionScopeHandler) IDFromModel(row models.SubscriptionScope) uint {
	return row.ID
}

// ListSubscriptionScopes returns all subscription scopes
func ListSubscriptionScopes(db *gorm.DB) echo.HandlerFunc {
	return genericListHandler[models.SubscriptionScope, SubscriptionScopeRequest, SubscriptionScopeResponse](db, SubscriptionScopeHandler{})
}

// GetSubscriptionScope returns a single subscription scope by ID
func GetSubscriptionScope(db *gorm.DB) echo.HandlerFunc {
	return genericGetHandler[models.SubscriptionScope, SubscriptionScopeRequest, SubscriptionScopeResponse](db, SubscriptionScopeHandler{})
}
