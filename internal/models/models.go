package models

import (
	"time"
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
