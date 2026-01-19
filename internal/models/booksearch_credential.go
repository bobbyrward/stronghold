package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type BookSearchCredential struct {
	gorm.Model
	APIKey      string `gorm:"not null"`
	IPAddress   string
	ASN         string
	LastRefresh time.Time
}

// GetBookSearchCredential retrieves the current credential from database
// Returns the first (and should be only) credential record
func GetBookSearchCredential(db *gorm.DB) (*BookSearchCredential, error) {
	var credential BookSearchCredential
	result := db.First(&credential)
	if result.Error != nil {
		return nil, result.Error
	}
	return &credential, nil
}

// UpsertBookSearchCredential creates or updates the credential
func UpsertBookSearchCredential(db *gorm.DB, apiKey, ipAddress, asn string) error {
	var credential BookSearchCredential
	result := db.First(&credential)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return result.Error
		}

		// Create new record
		credential = BookSearchCredential{
			APIKey:      apiKey,
			IPAddress:   ipAddress,
			ASN:         asn,
			LastRefresh: time.Now(),
		}
		return db.Create(&credential).Error
	}

	// Update existing record
	return db.Model(&credential).Updates(map[string]any{
		"api_key":      apiKey,
		"ip_address":   ipAddress,
		"asn":          asn,
		"last_refresh": time.Now(),
	}).Error
}
