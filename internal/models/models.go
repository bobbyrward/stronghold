package models

import (
	"time"

	"gorm.io/gorm"
)

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

type TorrentCategory struct {
	gorm.Model
	Name string `gorm:"not null;uniqueIndex"`
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
	FeedID            uint   `gorm:"not null;uniqueIndex:idx_feed_author"`
	Feed              Feed
	TorrentCategoryID uint   `gorm:"not null"`
	TorrentCategory   TorrentCategory
	NotifierID        uint   `gorm:"not null"`
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
