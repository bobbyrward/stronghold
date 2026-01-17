package api

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
)

// AuthorSubscriptionRequest is the request body for creating/updating author subscriptions
type AuthorSubscriptionRequest struct {
	ScopeName  string `json:"scope_name" validate:"required"`
	NotifierID *uint  `json:"notifier_id"`
}

// AuthorSubscriptionResponse is the response body for author subscriptions
type AuthorSubscriptionResponse struct {
	ID           uint    `json:"id"`
	AuthorID     uint    `json:"author_id"`
	AuthorName   string  `json:"author_name"`
	ScopeID      uint    `json:"scope_id"`
	ScopeName    string  `json:"scope_name"`
	NotifierID   *uint   `json:"notifier_id"`
	NotifierName *string `json:"notifier_name"`
}

// subscriptionToResponse converts an AuthorSubscription model to a response
func subscriptionToResponse(sub models.AuthorSubscription) AuthorSubscriptionResponse {
	resp := AuthorSubscriptionResponse{
		ID:         sub.ID,
		AuthorID:   sub.AuthorID,
		AuthorName: sub.Author.Name,
		ScopeID:    sub.ScopeID,
		ScopeName:  sub.Scope.Name,
		NotifierID: sub.NotifierID,
	}
	if sub.Notifier != nil {
		resp.NotifierName = &sub.Notifier.Name
	}
	return resp
}

// GetAuthorSubscription returns the subscription for an author
func GetAuthorSubscription(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		authorID, err := ParseAuthorIDParam(c, ctx)
		if err != nil {
			return BadRequest(c, ctx, "Invalid author_id")
		}

		slog.InfoContext(ctx, "Getting author subscription", slog.Uint64("author_id", uint64(authorID)))

		var sub models.AuthorSubscription
		err = db.Preload("Author").Preload("Scope").Preload("Notifier").
			Where("author_id = ?", authorID).First(&sub).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return GenericNotFound(c, ctx, "Subscription not found for this author")
			}
			return InternalError(c, ctx, "Failed to query subscription", err)
		}

		slog.InfoContext(ctx, "Successfully retrieved author subscription", slog.Uint64("id", uint64(sub.ID)), slog.Uint64("author_id", uint64(authorID)))
		return c.JSON(http.StatusOK, subscriptionToResponse(sub))
	}
}

// CreateAuthorSubscription creates a new subscription for an author
func CreateAuthorSubscription(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		authorID, err := ParseAuthorIDParam(c, ctx)
		if err != nil {
			return BadRequest(c, ctx, "Invalid author_id")
		}

		slog.InfoContext(ctx, "Creating author subscription", slog.Uint64("author_id", uint64(authorID)))

		// Verify author exists
		var author models.Author
		if err := db.First(&author, authorID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return NotFound(c, ctx, "Author", authorID)
			}
			return InternalError(c, ctx, "Failed to query author", err)
		}

		// Check if subscription already exists
		var existing models.AuthorSubscription
		if err := db.Where("author_id = ?", authorID).First(&existing).Error; err == nil {
			slog.WarnContext(ctx, "Subscription already exists", slog.Uint64("author_id", uint64(authorID)))
			return c.JSON(http.StatusConflict, map[string]string{"error": "Subscription already exists for this author"})
		}

		var req AuthorSubscriptionRequest
		if err := BindRequest(c, ctx, &req); err != nil {
			return BadRequest(c, ctx, "Invalid request body")
		}
		if err := ValidateRequest(c, ctx, &req); err != nil {
			return BadRequest(c, ctx, "Validation failed")
		}

		// Lookup scope by name
		var scope models.SubscriptionScope
		if err := LookupByName(db, ctx, &scope, req.ScopeName, "Subscription scope"); err != nil {
			return BadRequest(c, ctx, "Invalid scope_name: "+req.ScopeName)
		}

		sub := models.AuthorSubscription{
			AuthorID:   authorID,
			ScopeID:    scope.ID,
			NotifierID: req.NotifierID,
		}

		if err := db.Create(&sub).Error; err != nil {
			return InternalError(c, ctx, "Failed to create subscription", err)
		}

		// Reload with relations for response
		db.Preload("Author").Preload("Scope").Preload("Notifier").First(&sub, sub.ID)

		slog.InfoContext(ctx, "Created author subscription", slog.Uint64("id", uint64(sub.ID)), slog.Uint64("author_id", uint64(authorID)))
		return c.JSON(http.StatusCreated, subscriptionToResponse(sub))
	}
}

// UpdateAuthorSubscription updates an author's subscription
func UpdateAuthorSubscription(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		authorID, err := ParseAuthorIDParam(c, ctx)
		if err != nil {
			return BadRequest(c, ctx, "Invalid author_id")
		}

		slog.InfoContext(ctx, "Updating author subscription", slog.Uint64("author_id", uint64(authorID)))

		var sub models.AuthorSubscription
		if err := db.Where("author_id = ?", authorID).First(&sub).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return GenericNotFound(c, ctx, "Subscription not found for this author")
			}
			return InternalError(c, ctx, "Failed to query subscription", err)
		}

		var req AuthorSubscriptionRequest
		if err := BindRequest(c, ctx, &req); err != nil {
			return BadRequest(c, ctx, "Invalid request body")
		}
		if err := ValidateRequest(c, ctx, &req); err != nil {
			return BadRequest(c, ctx, "Validation failed")
		}

		// Lookup scope by name
		var scope models.SubscriptionScope
		if err := LookupByName(db, ctx, &scope, req.ScopeName, "Subscription scope"); err != nil {
			return BadRequest(c, ctx, "Invalid scope_name: "+req.ScopeName)
		}

		sub.ScopeID = scope.ID
		sub.NotifierID = req.NotifierID

		if err := db.Save(&sub).Error; err != nil {
			return InternalError(c, ctx, "Failed to update subscription", err)
		}

		// Reload with relations for response
		db.Preload("Author").Preload("Scope").Preload("Notifier").First(&sub, sub.ID)

		slog.InfoContext(ctx, "Updated author subscription", slog.Uint64("id", uint64(sub.ID)), slog.Uint64("author_id", uint64(authorID)))
		return c.JSON(http.StatusOK, subscriptionToResponse(sub))
	}
}

// DeleteAuthorSubscription deletes an author's subscription
func DeleteAuthorSubscription(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		authorID, err := ParseAuthorIDParam(c, ctx)
		if err != nil {
			return BadRequest(c, ctx, "Invalid author_id")
		}

		slog.InfoContext(ctx, "Deleting author subscription", slog.Uint64("author_id", uint64(authorID)))

		var sub models.AuthorSubscription
		if err := db.Where("author_id = ?", authorID).First(&sub).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return GenericNotFound(c, ctx, "Subscription not found for this author")
			}
			return InternalError(c, ctx, "Failed to query subscription", err)
		}

		if err := db.Unscoped().Delete(&sub).Error; err != nil {
			return InternalError(c, ctx, "Failed to delete subscription", err)
		}

		slog.InfoContext(ctx, "Deleted author subscription", slog.Uint64("id", uint64(sub.ID)), slog.Uint64("author_id", uint64(authorID)))
		return c.NoContent(http.StatusNoContent)
	}
}

// AuthorSubscriptionItemResponse is the response body for subscription items
type AuthorSubscriptionItemResponse struct {
	ID                   uint      `json:"id"`
	AuthorSubscriptionID uint      `json:"author_subscription_id"`
	TorrentHash          string    `json:"torrent_hash"`
	BooksearchID         string    `json:"booksearch_id"`
	DownloadedAt         time.Time `json:"downloaded_at"`
}

// itemToResponse converts an AuthorSubscriptionItem model to a response
func itemToResponse(item models.AuthorSubscriptionItem) AuthorSubscriptionItemResponse {
	return AuthorSubscriptionItemResponse{
		ID:                   item.ID,
		AuthorSubscriptionID: item.AuthorSubscriptionID,
		TorrentHash:          item.TorrentHash,
		BooksearchID:         item.BooksearchID,
		DownloadedAt:         item.DownloadedAt,
	}
}

// ListAuthorSubscriptionItems returns subscription items for an author's subscription
func ListAuthorSubscriptionItems(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		authorID, err := ParseAuthorIDParam(c, ctx)
		if err != nil {
			return BadRequest(c, ctx, "Invalid author_id")
		}

		slog.InfoContext(ctx, "Listing author subscription items", slog.Uint64("author_id", uint64(authorID)))

		// Find subscription for this author
		var sub models.AuthorSubscription
		if err := db.Where("author_id = ?", authorID).First(&sub).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return GenericNotFound(c, ctx, "Subscription not found for this author")
			}
			return InternalError(c, ctx, "Failed to query subscription", err)
		}

		// Get items ordered by download time DESC
		var items []models.AuthorSubscriptionItem
		if err := db.Where("author_subscription_id = ?", sub.ID).
			Order("downloaded_at DESC").
			Find(&items).Error; err != nil {
			return InternalError(c, ctx, "Failed to list subscription items", err)
		}

		response := make([]AuthorSubscriptionItemResponse, len(items))
		for i, item := range items {
			response[i] = itemToResponse(item)
		}

		slog.InfoContext(ctx, "Listed subscription items",
			slog.Uint64("author_id", uint64(authorID)),
			slog.Int("count", len(items)))
		return c.JSON(http.StatusOK, response)
	}
}
