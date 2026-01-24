package models

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/bobbyrward/stronghold/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SlogGormLogger integrates slog with GORM
type SlogGormLogger struct {
	SlowThreshold             time.Duration
	LogLevel                  logger.LogLevel
	IgnoreRecordNotFoundError bool
}

// NewSlogGormLogger creates a new GORM logger that uses slog
func NewSlogGormLogger() *SlogGormLogger {
	return &SlogGormLogger{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  logger.Info,
		IgnoreRecordNotFoundError: true,
	}
}

func (l *SlogGormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

func (l *SlogGormLogger) Info(ctx context.Context, msg string, data ...any) {
	if l.LogLevel >= logger.Info {
		slog.InfoContext(ctx, fmt.Sprintf(msg, data...))
	}
}

func (l *SlogGormLogger) Warn(ctx context.Context, msg string, data ...any) {
	if l.LogLevel >= logger.Warn {
		slog.WarnContext(ctx, fmt.Sprintf(msg, data...))
	}
}

func (l *SlogGormLogger) Error(ctx context.Context, msg string, data ...any) {
	if l.LogLevel >= logger.Error {
		slog.ErrorContext(ctx, fmt.Sprintf(msg, data...))
	}
}

func (l *SlogGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	if err != nil && (!errors.Is(err, logger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError) {
		slog.ErrorContext(ctx, "Database query failed",
			slog.String("sql", sql),
			slog.Duration("elapsed", elapsed),
			slog.Int64("rows", rows),
			slog.Any("error", err))
	} else if elapsed > l.SlowThreshold && l.SlowThreshold != 0 {
		slog.WarnContext(ctx, "Slow database query",
			slog.String("sql", sql),
			slog.Duration("elapsed", elapsed),
			slog.Duration("threshold", l.SlowThreshold),
			slog.Int64("rows", rows))
	} else if l.LogLevel == logger.Info {
		slog.InfoContext(ctx, "Database query",
			slog.String("sql", sql),
			slog.Duration("elapsed", elapsed),
			slog.Int64("rows", rows))
	}
}

// ConnectAndMigrate connects to the database and runs migrations.
// This is a convenience function that combines ConnectDB and AutoMigrate.
func ConnectAndMigrate(ctx context.Context) (*gorm.DB, error) {
	db, err := ConnectDB()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to connect to database", slog.Any("error", err))
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	err = AutoMigrate(db)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to auto-migrate database", slog.Any("error", err))
		return nil, fmt.Errorf("failed to automigrate database: %w", err)
	}

	return db, nil
}

func ConnectDB() (*gorm.DB, error) {
	ctx := context.Background()

	slog.InfoContext(ctx, "Connecting to database",
		slog.String("url", config.Config.Postgres.URL))

	gormConfig := &gorm.Config{
		Logger: NewSlogGormLogger(),
	}

	db, err := gorm.Open(postgres.Open(config.Config.Postgres.URL), gormConfig)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to connect to database",
			slog.String("url", config.Config.Postgres.URL),
			slog.Any("err", err))
		return nil, errors.Join(err, fmt.Errorf("failed to connect to db"))
	}

	slog.InfoContext(ctx, "Successfully connected to database")
	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	ctx := context.Background()

	slog.InfoContext(ctx, "Starting database auto-migration")

	err := db.AutoMigrate(
		// Existing models
		&FeedItem{},
		&SearchResponseItem{},
		&NotificationType{},
		&Notifier{},
		&Feed{},
		&FeedAuthorFilter{},
		&FeedFilter{},
		&FilterKey{},
		&FilterOperator{},
		&FeedFilterSetType{},
		&FeedFilterSet{},
		&FeedFilterSetEntry{},
		&BookSearchCredential{},
		// Feedwatcher2 models
		&SubscriptionScope{},
		&BookType{},
		&Library{},
		&TorrentCategory{},
		&Author{},
		&AuthorAlias{},
		&AuthorSubscription{},
		&AuthorSubscriptionItem{},
	)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to auto-migrate database", slog.Any("err", err))
		return errors.Join(err, fmt.Errorf("failed to auto migrate db"))
	}

	slog.InfoContext(ctx, "Successfully completed database auto-migration")

	return PopulateData(db)
}

func PopulateData(db *gorm.DB) error {
	err := populateSubscriptionScopes(db)
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to populate subscription scopes"))
	}

	err = populateBookTypes(db)
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to populate book types"))
	}

	err = populateTorrentCategories(db)
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to populate torrent categories"))
	}

	err = populateNotificationType(db)
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to populate notification types"))
	}

	err = populateFilterKeys(db)
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to populate filter keys"))
	}

	err = populateFilterOperators(db)
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to populate filter operators"))
	}

	err = populateFilterSetTypes(db)
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to populate filter set types"))
	}

	return nil
}

// ConnectTestDB creates an in-memory SQLite database for testing
func ConnectTestDB() (*gorm.DB, error) {
	ctx := context.Background()

	slog.InfoContext(ctx, "Connecting to test database (SQLite in-memory)")

	gormConfig := &gorm.Config{
		Logger: NewSlogGormLogger(),
	}

	db, err := gorm.Open(sqlite.Open(":memory:"), gormConfig)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to connect to test database", slog.Any("err", err))
		return nil, err
	}

	slog.InfoContext(ctx, "Successfully connected to test database")

	// Run auto migration
	err = AutoMigrate(db)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to auto-migrate test database", slog.Any("err", err))
		return nil, err
	}

	slog.InfoContext(ctx, "Test database ready with migrations and seed data")
	return db, nil
}
