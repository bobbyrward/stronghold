package api

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/api/handlers"
)

// RegisterRoutes registers all API routes with the Echo server
func RegisterRoutes(e *echo.Echo, db *gorm.DB) {
	// Filter Keys
	e.GET("/filter-keys", handlers.ListFilterKeys(db))
	e.POST("/filter-keys", handlers.CreateFilterKey(db))
	e.GET("/filter-keys/:id", handlers.GetFilterKey(db))
	e.PUT("/filter-keys/:id", handlers.UpdateFilterKey(db))
	e.DELETE("/filter-keys/:id", handlers.DeleteFilterKey(db))

	// Filter Operators
	e.GET("/filter-operators", handlers.ListFilterOperators(db))
	e.POST("/filter-operators", handlers.CreateFilterOperator(db))
	e.GET("/filter-operators/:id", handlers.GetFilterOperator(db))
	e.PUT("/filter-operators/:id", handlers.UpdateFilterOperator(db))
	e.DELETE("/filter-operators/:id", handlers.DeleteFilterOperator(db))

	// Notification Types
	e.GET("/notification-types", handlers.ListNotificationTypes(db))
	e.POST("/notification-types", handlers.CreateNotificationType(db))
	e.GET("/notification-types/:id", handlers.GetNotificationType(db))
	e.PUT("/notification-types/:id", handlers.UpdateNotificationType(db))
	e.DELETE("/notification-types/:id", handlers.DeleteNotificationType(db))

	// Feed Filter Set Types
	e.GET("/feed-filter-set-types", handlers.ListFeedFilterSetTypes(db))
	e.POST("/feed-filter-set-types", handlers.CreateFeedFilterSetType(db))
	e.GET("/feed-filter-set-types/:id", handlers.GetFeedFilterSetType(db))
	e.PUT("/feed-filter-set-types/:id", handlers.UpdateFeedFilterSetType(db))
	e.DELETE("/feed-filter-set-types/:id", handlers.DeleteFeedFilterSetType(db))

	// Torrent Categories
	e.GET("/torrent-categories", handlers.ListTorrentCategories(db))
	e.POST("/torrent-categories", handlers.CreateTorrentCategory(db))
	e.GET("/torrent-categories/:id", handlers.GetTorrentCategory(db))
	e.PUT("/torrent-categories/:id", handlers.UpdateTorrentCategory(db))
	e.DELETE("/torrent-categories/:id", handlers.DeleteTorrentCategory(db))

	// Notifiers
	e.GET("/notifiers", handlers.ListNotifiers(db))
	e.POST("/notifiers", handlers.CreateNotifier(db))
	e.GET("/notifiers/:id", handlers.GetNotifier(db))
	e.PUT("/notifiers/:id", handlers.UpdateNotifier(db))
	e.DELETE("/notifiers/:id", handlers.DeleteNotifier(db))

	// Feeds
	e.GET("/feeds", handlers.ListFeeds(db))
	e.POST("/feeds", handlers.CreateFeed(db))
	e.GET("/feeds/:id", handlers.GetFeed(db))
	e.PUT("/feeds/:id", handlers.UpdateFeed(db))
	e.DELETE("/feeds/:id", handlers.DeleteFeed(db))

	// Feed Filters
	e.GET("/feed-filters", handlers.ListFeedFilters(db))
	e.POST("/feed-filters", handlers.CreateFeedFilter(db))
	e.GET("/feed-filters/:id", handlers.GetFeedFilter(db))
	e.PUT("/feed-filters/:id", handlers.UpdateFeedFilter(db))
	e.DELETE("/feed-filters/:id", handlers.DeleteFeedFilter(db))

	// Feed Filter Sets
	e.GET("/feed-filter-sets", handlers.ListFeedFilterSets(db))
	e.POST("/feed-filter-sets", handlers.CreateFeedFilterSet(db))
	e.GET("/feed-filter-sets/:id", handlers.GetFeedFilterSet(db))
	e.PUT("/feed-filter-sets/:id", handlers.UpdateFeedFilterSet(db))
	e.DELETE("/feed-filter-sets/:id", handlers.DeleteFeedFilterSet(db))

	// Feed Filter Set Entries
	e.GET("/feed-filter-set-entries", handlers.ListFeedFilterSetEntries(db))
	e.POST("/feed-filter-set-entries", handlers.CreateFeedFilterSetEntry(db))
	e.GET("/feed-filter-set-entries/:id", handlers.GetFeedFilterSetEntry(db))
	e.PUT("/feed-filter-set-entries/:id", handlers.UpdateFeedFilterSetEntry(db))
	e.DELETE("/feed-filter-set-entries/:id", handlers.DeleteFeedFilterSetEntry(db))
}
