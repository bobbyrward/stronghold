package models

import (
	"fmt"

	"gorm.io/gorm"
)

type (
	existsFunc func(name string) error
	createFunc func(name string) error
)

func populateTable(db *gorm.DB, names []string, exists existsFunc, create createFunc) error {
	for _, name := range names {
		err := exists(name)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := create(name); err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	return nil
}

func populateSubscriptionScopes(db *gorm.DB) error {
	return populateTable(
		db,
		[]string{
			"personal",
			"family",
			"kids",
			"general",
		},
		func(name string) error {
			var record SubscriptionScope
			return db.Where("name = ?", name).First(&record).Error
		},
		func(name string) error {
			record := SubscriptionScope{Name: name}
			return db.Create(&record).Error
		},
	)
}

func populateTorrentCategories(db *gorm.DB) error {
	// Must run after populateSubscriptionScopes
	type categoryDef struct {
		Name      string
		ScopeName string
		MediaType string
	}

	categories := []categoryDef{
		{"audiobooks", "family", "audiobook"},
		{"books", "family", "ebook"},
		{"personal-audiobooks", "personal", "audiobook"},
		{"personal-books", "personal", "ebook"},
		{"kids-audiobooks", "kids", "audiobook"},
		{"kids-books", "kids", "ebook"},
		{"general-audiobooks", "general", "audiobook"},
		{"general-books", "general", "ebook"},
	}

	// Build scope lookup map
	var scopes []SubscriptionScope
	if err := db.Find(&scopes).Error; err != nil {
		return err
	}

	scopeMap := make(map[string]uint)
	for _, s := range scopes {
		scopeMap[s.Name] = s.ID
	}

	for _, cat := range categories {
		scopeID, ok := scopeMap[cat.ScopeName]
		if !ok {
			return fmt.Errorf("scope not found: %s", cat.ScopeName)
		}

		var record TorrentCategory
		err := db.Where("name = ?", cat.Name).First(&record).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				record := TorrentCategory{
					Name:      cat.Name,
					ScopeID:   scopeID,
					MediaType: cat.MediaType,
				}
				if err := db.Create(&record).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	return nil
}

func populateNotificationType(db *gorm.DB) error {
	return populateTable(
		db,
		[]string{
			"discord",
		},
		func(name string) error {
			var record NotificationType
			return db.Where("name = ?", name).First(&record).Error
		},
		func(name string) error {
			record := NotificationType{Name: name}
			return db.Create(&record).Error
		},
	)
}

func populateFilterKeys(db *gorm.DB) error {
	return populateTable(
		db,
		[]string{
			"author",
			"series",
			"title",
			"category",
			"summary",
			"tags",
			"description",
		},
		func(name string) error {
			var record FilterKey
			return db.Where("name = ?", name).First(&record).Error
		},
		func(name string) error {
			record := FilterKey{Name: name}
			return db.Create(&record).Error
		},
	)
}

func populateFilterOperators(db *gorm.DB) error {
	return populateTable(
		db,
		[]string{
			"equals",
			"contains",
			"fnmatch",
			"regex",
		},
		func(name string) error {
			var record FilterOperator
			return db.Where("name = ?", name).First(&record).Error
		},
		func(name string) error {
			record := FilterOperator{Name: name}
			return db.Create(&record).Error
		},
	)
}

func populateFilterSetTypes(db *gorm.DB) error {
	return populateTable(
		db,
		[]string{
			"any",
			"all",
		},
		func(name string) error {
			var record FeedFilterSetType
			return db.Where("name = ?", name).First(&record).Error
		},
		func(name string) error {
			record := FeedFilterSetType{Name: name}
			return db.Create(&record).Error
		},
	)
}
