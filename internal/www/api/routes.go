package api

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/hardcover"
)

// RegisterRoutes registers all API routes with the Echo server
func RegisterRoutes(e *echo.Group, db *gorm.DB, hc hardcover.Client) {
	// Filter Keys (read-only reference data)
	e.GET("/filter-keys", ListFilterKeys(db))
	e.GET("/filter-keys/:id", GetFilterKey(db))

	// Filter Operators (read-only reference data)
	e.GET("/filter-operators", ListFilterOperators(db))
	e.GET("/filter-operators/:id", GetFilterOperator(db))

	// Notification Types (read-only reference data)
	e.GET("/notification-types", ListNotificationTypes(db))
	e.GET("/notification-types/:id", GetNotificationType(db))

	// Feed Filter Set Types (read-only reference data)
	e.GET("/feed-filter-set-types", ListFeedFilterSetTypes(db))
	e.GET("/feed-filter-set-types/:id", GetFeedFilterSetType(db))

	// Torrent Categories (read-only reference data)
	e.GET("/torrent-categories", ListTorrentCategories(db))
	e.GET("/torrent-categories/:id", GetTorrentCategory(db))

	// Subscription Scopes (read-only reference data)
	e.GET("/subscription-scopes", ListSubscriptionScopes(db))
	e.GET("/subscription-scopes/:id", GetSubscriptionScope(db))

	// Book Types (read-only reference data)
	e.GET("/book-types", ListBookTypes(db))
	e.GET("/book-types/:id", GetBookType(db))

	// Libraries
	e.GET("/libraries", ListLibraries(db))
	e.POST("/libraries", CreateLibrary(db))
	e.GET("/libraries/:id", GetLibrary(db))
	e.PUT("/libraries/:id", UpdateLibrary(db))
	e.DELETE("/libraries/:id", DeleteLibrary(db))

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
	e.POST("/book-torrent-dl", DownloadBookTorrent(db, nil))

	// Authors
	e.GET("/authors", ListAuthors(db, hc))
	e.POST("/authors", CreateAuthor(db, hc))
	e.GET("/authors/:id", GetAuthor(db, hc))
	e.PUT("/authors/:id", UpdateAuthor(db, hc))
	e.DELETE("/authors/:id", DeleteAuthor(db))

	// Author Aliases (nested under authors)
	e.GET("/authors/:author_id/aliases", ListAuthorAliases(db))
	e.POST("/authors/:author_id/aliases", CreateAuthorAlias(db))
	e.GET("/authors/:author_id/aliases/:id", GetAuthorAlias(db))
	e.PUT("/authors/:author_id/aliases/:id", UpdateAuthorAlias(db))
	e.DELETE("/authors/:author_id/aliases/:id", DeleteAuthorAlias(db))

	// Author Subscriptions (nested under authors, singleton per author)
	e.GET("/authors/:author_id/subscription", GetAuthorSubscription(db))
	e.POST("/authors/:author_id/subscription", CreateAuthorSubscription(db))
	e.PUT("/authors/:author_id/subscription", UpdateAuthorSubscription(db))
	e.DELETE("/authors/:author_id/subscription", DeleteAuthorSubscription(db))

	// Author Subscription Items (nested under subscription)
	e.GET("/authors/:author_id/subscription/items", ListAuthorSubscriptionItems(db))

	// Hardcover
	e.GET("/hardcover/authors/search", SearchHardcoverAuthors(hc))
}
