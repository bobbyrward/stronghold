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

/*
   - name: MAM
     url: https://04k6i.mrd.ninja/rss/52cf3a54
     filters:
       - category: personal-books
         match:
           - key: author
             operator: equals
             value: Blaise Corvin
*/
