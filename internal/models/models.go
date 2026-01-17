package models

import (
	"time"

	"gorm.io/gorm"
)

type CommonFields struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type FeedItem struct {
	ID          uint      `gorm:"primaryKey"`
	Guid        string    `gorm:"not null;uniqueIndex"`
	Title       string    `gorm:"not null"`
	Link        string    `gorm:"not null"`
	Category    string    `gorm:"not null"`
	Description string    `gorm:"not null"`
	PubDate     time.Time `gorm:"not null"`
	CreatedAt   time.Time
}

type SearchResponseItem struct {
	ID           uint   `gorm:"primaryKey"`
	Title        string `gorm:"notnull"`
	Authors      string `gorm:"notnull"`
	MainCategory int    `gorm:"notnull"`
	Category     string `gorm:"notnull"`
	DlHash       string `gorm:"notnull"`
	FileTypes    string `gorm:"notnull"`
	TorrentID    uint   `gorm:"notnull"`
	Language     string `gorm:"notnull"`
	Series       string `gorm:"notnull"`
	Narrators    string `gorm:"notnull"`
	Size         string `gorm:"notnull"`
	Tags         string `gorm:"notnull"`
}

type NotificationType struct {
	gorm.Model
	Name string `gorm:"not null;uniqueIndex"`
}

type Notifier struct {
	gorm.Model
	Name               string `gorm:"not null;uniqueIndex"`
	NotificationTypeID uint
	NotificationType   *NotificationType
	URL                string
}

type Feed struct {
	gorm.Model
	Name string `gorm:"not null;uniqueIndex"`
	URL  string
}

type FeedAuthorFilter struct {
	gorm.Model
	FeedID            uint `gorm:"not null;uniqueIndex:idx_feed_author"`
	Feed              Feed
	TorrentCategoryID uint `gorm:"not null"`
	TorrentCategory   TorrentCategory
	NotifierID        uint `gorm:"not null"`
	Notifier          Notifier
	Author            string `gorm:"not null;uniqueIndex:idx_feed_author"`
}

type FeedFilter struct {
	gorm.Model
	Name              string
	FeedID            uint `gorm:"not null"`
	Feed              Feed
	TorrentCategoryID uint `gorm:"not null"`
	TorrentCategory   TorrentCategory
	NotifierID        uint `gorm:"not null"`
	Notifier          Notifier
}

type FilterKey struct {
	gorm.Model
	Name string `gorm:"not null;uniqueIndex"`
}

type FilterOperator struct {
	gorm.Model
	Name string `gorm:"not null;uniqueIndex"`
}

type FeedFilterSetType struct {
	gorm.Model
	Name string `gorm:"not null;uniqueIndex"`
}

type FeedFilterSet struct {
	gorm.Model
	FeedFilterID        uint `gorm:"not null"`
	FeedFilter          FeedFilter
	FeedFilterSetTypeID uint `gorm:"not null"`
	FeedFilterSetType   FeedFilterSetType
}

type FeedFilterSetEntry struct {
	gorm.Model
	FeedFilterSetID  uint `gorm:"not null"`
	FeedFilterSet    FeedFilterSet
	FilterKeyID      uint `gorm:"not null"`
	FilterKey        FilterKey
	FilterOperatorID uint `gorm:"not null"`
	FilterOperator   FilterOperator
	Value            string
}

// SubscriptionScope is a reference table for subscription scopes
type SubscriptionScope struct {
	CommonFields
	Name string `gorm:"not null;uniqueIndex"` // "personal" or "family"
}

// TorrentCategory (updated - replaces existing model)
type TorrentCategory struct {
	CommonFields
	Name      string `gorm:"not null;uniqueIndex"`
	ScopeID   uint   `gorm:"not null"`
	Scope     SubscriptionScope
	MediaType string `gorm:"not null"` // "audiobook" or "ebook"
}

// Author represents a writer of books
type Author struct {
	CommonFields
	Name         string  `gorm:"not null;uniqueIndex"`
	HardcoverRef *string // nullable until linked; slug/UUID TBD
}

// AuthorAlias represents an additional alias for an Author
type AuthorAlias struct {
	CommonFields
	AuthorID uint   `gorm:"not null"`
	Author   Author `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Name     string `gorm:"not null;uniqueIndex"` // globally unique
}

// AuthorSubscription represents a subscription to an Author
type AuthorSubscription struct {
	CommonFields
	AuthorID   uint   `gorm:"not null;uniqueIndex"` // one subscription per author
	Author     Author `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ScopeID    uint   `gorm:"not null"`
	Scope      SubscriptionScope
	NotifierID *uint
	Notifier   *Notifier
}

// AuthorSubscriptionItem represents a downloaded item from an AuthorSubscription
type AuthorSubscriptionItem struct {
	CommonFields
	AuthorSubscriptionID uint               `gorm:"not null"`
	AuthorSubscription   AuthorSubscription `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	TorrentHash          string             `gorm:"not null"`
	BooksearchID         string             `gorm:"not null;uniqueIndex"` // torrent ID extracted from feed GUID URL
	DownloadedAt         time.Time          `gorm:"not null"`
}
