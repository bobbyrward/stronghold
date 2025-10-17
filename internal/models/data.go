package models

import "gorm.io/gorm"

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

func populateTorrentCategories(db *gorm.DB) error {
	return populateTable(
		db,
		[]string{
			"audiobooks",
			"books",
			"personal-books",
			"personal-audiobooks",
		},
		func(name string) error {
			var record TorrentCategory
			return db.Where("name = ?", name).First(&record).Error
		},
		func(name string) error {
			record := TorrentCategory{Name: name}
			return db.Create(&record).Error
		},
	)
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
