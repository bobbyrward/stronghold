package api

import (
	"context"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
)

type SubscriptionScopeResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type SubscriptionScopeHandler struct{}

func (h SubscriptionScopeHandler) ModelToResponse(c echo.Context, ctx context.Context, db *gorm.DB, row models.SubscriptionScope) SubscriptionScopeResponse {
	return SubscriptionScopeResponse{
		ID:   row.ID,
		Name: row.Name,
	}
}

func (h SubscriptionScopeHandler) PreloadRelations(c echo.Context, ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	return db, nil
}

func (h SubscriptionScopeHandler) IDFromModel(row models.SubscriptionScope) uint {
	return row.ID
}

// ListSubscriptionScopes returns all subscription scopes
func ListSubscriptionScopes(db *gorm.DB) echo.HandlerFunc {
	return readOnlyListHandler[models.SubscriptionScope, SubscriptionScopeResponse](db, SubscriptionScopeHandler{})
}

// GetSubscriptionScope returns a single subscription scope by ID
func GetSubscriptionScope(db *gorm.DB) echo.HandlerFunc {
	return readOnlyGetHandler[models.SubscriptionScope, SubscriptionScopeResponse](db, SubscriptionScopeHandler{})
}
