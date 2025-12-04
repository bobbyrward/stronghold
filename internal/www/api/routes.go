package api

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// RegisterRoutes registers all API routes with the Echo server
func RegisterRoutes(e *echo.Group, db *gorm.DB) {
	// Filter Keys
	e.GET("/filter-keys", ListFilterKeys(db))
	e.POST("/filter-keys", CreateFilterKey(db))
	e.GET("/filter-keys/:id", GetFilterKey(db))
	e.PUT("/filter-keys/:id", UpdateFilterKey(db))
	e.DELETE("/filter-keys/:id", DeleteFilterKey(db))

	// Filter Operators
	e.GET("/filter-operators", ListFilterOperators(db))
	e.POST("/filter-operators", CreateFilterOperator(db))
	e.GET("/filter-operators/:id", GetFilterOperator(db))
	e.PUT("/filter-operators/:id", UpdateFilterOperator(db))
	e.DELETE("/filter-operators/:id", DeleteFilterOperator(db))

	// Notification Types
	e.GET("/notification-types", ListNotificationTypes(db))
	e.POST("/notification-types", CreateNotificationType(db))
	e.GET("/notification-types/:id", GetNotificationType(db))
	e.PUT("/notification-types/:id", UpdateNotificationType(db))
	e.DELETE("/notification-types/:id", DeleteNotificationType(db))

	// Feed Filter Set Types
	e.GET("/feed-filter-set-types", ListFeedFilterSetTypes(db))
	e.POST("/feed-filter-set-types", CreateFeedFilterSetType(db))
	e.GET("/feed-filter-set-types/:id", GetFeedFilterSetType(db))
	e.PUT("/feed-filter-set-types/:id", UpdateFeedFilterSetType(db))
	e.DELETE("/feed-filter-set-types/:id", DeleteFeedFilterSetType(db))

	// Torrent Categories
	e.GET("/torrent-categories", ListTorrentCategories(db))
	e.POST("/torrent-categories", CreateTorrentCategory(db))
	e.GET("/torrent-categories/:id", GetTorrentCategory(db))
	e.PUT("/torrent-categories/:id", UpdateTorrentCategory(db))
	e.DELETE("/torrent-categories/:id", DeleteTorrentCategory(db))

	// Notifiers
	e.GET("/notifiers", ListNotifiers(db))
	e.POST("/notifiers", CreateNotifier(db))
	e.GET("/notifiers/:id", GetNotifier(db))
	e.PUT("/notifiers/:id", UpdateNotifier(db))
	e.DELETE("/notifiers/:id", DeleteNotifier(db))

	// Feeds
	e.GET("/feeds", ListFeeds(db))
	e.POST("/feeds", CreateFeed(db))
	e.GET("/feeds/:id", GetFeed(db))
	e.PUT("/feeds/:id", UpdateFeed(db))
	e.DELETE("/feeds/:id", DeleteFeed(db))

	// Feed Filters
	e.GET("/feed-filters", ListFeedFilters(db))
	e.POST("/feed-filters", CreateFeedFilter(db))
	e.GET("/feed-filters/:id", GetFeedFilter(db))
	e.PUT("/feed-filters/:id", UpdateFeedFilter(db))
	e.DELETE("/feed-filters/:id", DeleteFeedFilter(db))

	// Feed Author Filters
	e.GET("/feed-author-filters", ListFeedAuthorFilters(db))
	e.POST("/feed-author-filters", CreateFeedAuthorFilter(db))
	e.GET("/feed-author-filters/:id", GetFeedAuthorFilter(db))
	e.PUT("/feed-author-filters/:id", UpdateFeedAuthorFilter(db))
	e.DELETE("/feed-author-filters/:id", DeleteFeedAuthorFilter(db))

	// Feed Filter Sets
	e.GET("/feed-filter-sets", ListFeedFilterSets(db))
	e.POST("/feed-filter-sets", CreateFeedFilterSet(db))
	e.GET("/feed-filter-sets/:id", GetFeedFilterSet(db))
	e.PUT("/feed-filter-sets/:id", UpdateFeedFilterSet(db))
	e.DELETE("/feed-filter-sets/:id", DeleteFeedFilterSet(db))

	// Feed Filter Set Entries
	e.GET("/feed-filter-set-entries", ListFeedFilterSetEntries(db))
	e.POST("/feed-filter-set-entries", CreateFeedFilterSetEntry(db))
	e.GET("/feed-filter-set-entries/:id", GetFeedFilterSetEntry(db))
	e.PUT("/feed-filter-set-entries/:id", UpdateFeedFilterSetEntry(db))
	e.DELETE("/feed-filter-set-entries/:id", DeleteFeedFilterSetEntry(db))

	// Torrents
	e.GET("/torrents/unimported", ListUnimportedTorrents(db))
	e.GET("/torrents/manual", ListManualInterventionTorrents(db))
	e.POST("/torrents/:hash/category", SetTorrentCategory(db))
	e.POST("/torrents/:hash/tags", SetTorrentTags(db))

	// Audiobook Wizard
	e.GET("/audiobook-wizard/torrent/:hash/info", GetTorrentInfo(db))
	e.POST("/audiobook-wizard/search-asin", SearchASIN(db))
	e.GET("/audiobook-wizard/asin/:asin/metadata", GetASINMetadata(db))
	e.POST("/audiobook-wizard/preview-directory", PreviewDirectory(db))
	e.GET("/audiobook-wizard/libraries", GetLibraries(db))
	e.POST("/audiobook-wizard/execute-import", ExecuteImport(db))

	// Downloads
	e.POST("/book-torrent-dl", DownloadBookTorrent(db))
}
