package db

import (
	"errors"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	postgresURI, ok := os.LookupEnv("POSTGRES_URI")
	if !ok {
		return nil, fmt.Errorf("env var POSTGRES_URI not set")
	}

	db, err := gorm.Open(postgres.Open(postgresURI), &gorm.Config{})
	if err != nil {
		return nil, errors.Join(err, fmt.Errorf("failed to connect to db"))
	}

	return db, nil
}
