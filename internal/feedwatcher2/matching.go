package feedwatcher2

import (
	"context"
	"log/slog"
	"strings"

	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
)

// AuthorMatcher matches feed authors against subscribed authors and aliases.
type AuthorMatcher struct {
	db                *gorm.DB
	subscriptionCache map[string]*models.AuthorSubscription
}

// NewAuthorMatcher creates a new AuthorMatcher.
func NewAuthorMatcher(db *gorm.DB) *AuthorMatcher {
	return &AuthorMatcher{
		db:                db,
		subscriptionCache: make(map[string]*models.AuthorSubscription),
	}
}

// normalizeName normalizes an author name for matching:
// - Removes all periods (.)
// - Converts to lowercase
// - Trims whitespace
func normalizeName(name string) string {
	// Remove all periods
	name = strings.ReplaceAll(name, ".", "")
	// Lowercase
	name = strings.ToLower(name)
	// Trim whitespace
	name = strings.TrimSpace(name)
	return name
}

// LoadSubscriptions loads all author subscriptions and aliases into the cache.
func (am *AuthorMatcher) LoadSubscriptions(ctx context.Context) error {
	slog.InfoContext(ctx, "Loading author subscriptions")

	// Query all subscriptions with preloaded relationships
	var subscriptions []models.AuthorSubscription
	err := am.db.
		Preload("Author").
		Preload("Scope").
		Preload("Notifier").
		Find(&subscriptions).Error
	if err != nil {
		slog.ErrorContext(ctx, "Failed to load subscriptions", slog.Any("error", err))
		return err
	}

	// Build a map of author ID to subscription for alias lookup
	authorToSubscription := make(map[uint]*models.AuthorSubscription)
	for i := range subscriptions {
		sub := &subscriptions[i]
		authorToSubscription[sub.AuthorID] = sub

		// Add author name to cache
		normalizedName := normalizeName(sub.Author.Name)
		am.subscriptionCache[normalizedName] = sub
		slog.DebugContext(ctx, "Cached author subscription",
			slog.String("author", sub.Author.Name),
			slog.String("normalized", normalizedName))
	}

	// Query all aliases
	var aliases []models.AuthorAlias
	err = am.db.Find(&aliases).Error
	if err != nil {
		slog.ErrorContext(ctx, "Failed to load aliases", slog.Any("error", err))
		return err
	}

	// Add aliases to cache (only for authors with subscriptions)
	aliasCount := 0
	for _, alias := range aliases {
		if sub, ok := authorToSubscription[alias.AuthorID]; ok {
			normalizedAlias := normalizeName(alias.Name)
			am.subscriptionCache[normalizedAlias] = sub
			aliasCount++
			slog.DebugContext(ctx, "Cached alias",
				slog.String("alias", alias.Name),
				slog.String("normalized", normalizedAlias),
				slog.String("author", sub.Author.Name))
		}
	}

	slog.InfoContext(ctx, "Loaded subscriptions and aliases",
		slog.Int("subscriptions", len(subscriptions)),
		slog.Int("aliases", aliasCount),
		slog.Int("cache_entries", len(am.subscriptionCache)))

	return nil
}

// FindMatchingSubscription finds a subscription that matches any of the given feed authors.
// Returns nil if no match is found.
func (am *AuthorMatcher) FindMatchingSubscription(feedAuthors []string) *models.AuthorSubscription {
	for _, author := range feedAuthors {
		normalized := normalizeName(author)
		if sub, ok := am.subscriptionCache[normalized]; ok {
			return sub
		}
	}
	return nil
}
