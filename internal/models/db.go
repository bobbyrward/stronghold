package models

import (
	"errors"
	"fmt"

	"github.com/bobbyrward/stronghold/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(config.Config.Postgres.URL), &gorm.Config{})
	if err != nil {
		return nil, errors.Join(err, fmt.Errorf("failed to connect to db"))
	}

	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&FeedItem{},
	)
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to auto migrate db"))
	}

	return nil
}
