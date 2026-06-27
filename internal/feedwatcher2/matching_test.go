package feedwatcher2

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bobbyrward/stronghold/internal/models"
)

func TestNormalizeName_RemovesDots(t *testing.T) {
	result := normalizeName("J.F. Brink")
	assert.Equal(t, "jf brink", result)
}

func TestNormalizeName_MultipleDots(t *testing.T) {
	result := normalizeName("J.R.R. Tolkien")
	assert.Equal(t, "jrr tolkien", result)
}

func TestNormalizeName_Lowercase(t *testing.T) {
	result := normalizeName("BRANDON SANDERSON")
	assert.Equal(t, "brandon sanderson", result)
}

func TestNormalizeName_TrimWhitespace(t *testing.T) {
	result := normalizeName("  Brandon Sanderson  ")
	assert.Equal(t, "brandon sanderson", result)
}

func TestNormalizeName_Combined(t *testing.T) {
	result := normalizeName("  J.R.R. TOLKIEN  ")
	assert.Equal(t, "jrr tolkien", result)
}

func TestNormalizeName_NoChanges(t *testing.T) {
	result := normalizeName("brandon sanderson")
	assert.Equal(t, "brandon sanderson", result)
}

func TestNormalizeName_Empty(t *testing.T) {
	result := normalizeName("")
	assert.Equal(t, "", result)
}

func TestNormalizeName_OnlyDots(t *testing.T) {
	result := normalizeName("...")
	assert.Equal(t, "", result)
}

func TestLoadSubscriptions(t *testing.T) {
	db, err := models.ConnectTestDB()
	require.NoError(t, err)

	// Create test data
	scope := models.SubscriptionScope{Name: "personal"}
	err = db.FirstOrCreate(&scope, models.SubscriptionScope{Name: "personal"}).Error
	require.NoError(t, err)

	author := models.Author{Name: "Test Author"}
	err = db.Create(&author).Error
	require.NoError(t, err)

	subscription := models.AuthorSubscription{
		AuthorID: author.ID,
		ScopeID:  scope.ID,
	}
	err = db.Create(&subscription).Error
	require.NoError(t, err)

	// Create matcher and load
	am := NewAuthorMatcher(db)
	err = am.LoadSubscriptions(context.Background())
	require.NoError(t, err)

	// Verify cache is populated
	assert.NotEmpty(t, am.subscriptionCache)
	assert.Contains(t, am.subscriptionCache, "test author")
}

func TestLoadSubscriptions_WithAliases(t *testing.T) {
	db, err := models.ConnectTestDB()
	require.NoError(t, err)

	// Create test data
	scope := models.SubscriptionScope{Name: "family"}
	err = db.FirstOrCreate(&scope, models.SubscriptionScope{Name: "family"}).Error
	require.NoError(t, err)

	author := models.Author{Name: "Brandon Sanderson"}
	err = db.Create(&author).Error
	require.NoError(t, err)

	alias := models.AuthorAlias{
		AuthorID: author.ID,
		Name:     "Brando Sando",
	}
	err = db.Create(&alias).Error
	require.NoError(t, err)

	subscription := models.AuthorSubscription{
		AuthorID: author.ID,
		ScopeID:  scope.ID,
	}
	err = db.Create(&subscription).Error
	require.NoError(t, err)

	// Create matcher and load
	am := NewAuthorMatcher(db)
	err = am.LoadSubscriptions(context.Background())
	require.NoError(t, err)

	// Verify both author name and alias are in cache
	assert.Contains(t, am.subscriptionCache, "brandon sanderson")
	assert.Contains(t, am.subscriptionCache, "brando sando")

	// Both should point to the same subscription
	sub1 := am.subscriptionCache["brandon sanderson"]
	sub2 := am.subscriptionCache["brando sando"]
	assert.Equal(t, sub1.ID, sub2.ID)
}

func TestFindMatchingSubscription_DirectMatch(t *testing.T) {
	db, err := models.ConnectTestDB()
	require.NoError(t, err)

	scope := models.SubscriptionScope{Name: "personal"}
	err = db.FirstOrCreate(&scope, models.SubscriptionScope{Name: "personal"}).Error
	require.NoError(t, err)

	author := models.Author{Name: "Test Author"}
	err = db.Create(&author).Error
	require.NoError(t, err)

	subscription := models.AuthorSubscription{
		AuthorID: author.ID,
		ScopeID:  scope.ID,
	}
	err = db.Create(&subscription).Error
	require.NoError(t, err)

	am := NewAuthorMatcher(db)
	err = am.LoadSubscriptions(context.Background())
	require.NoError(t, err)

	result := am.FindMatchingSubscription([]string{"Test Author"})
	require.NotNil(t, result)
	assert.Equal(t, subscription.ID, result.ID)
}

func TestFindMatchingSubscription_AliasMatch(t *testing.T) {
	db, err := models.ConnectTestDB()
	require.NoError(t, err)

	scope := models.SubscriptionScope{Name: "personal"}
	err = db.FirstOrCreate(&scope, models.SubscriptionScope{Name: "personal"}).Error
	require.NoError(t, err)

	author := models.Author{Name: "Real Name"}
	err = db.Create(&author).Error
	require.NoError(t, err)

	alias := models.AuthorAlias{
		AuthorID: author.ID,
		Name:     "Pen Name",
	}
	err = db.Create(&alias).Error
	require.NoError(t, err)

	subscription := models.AuthorSubscription{
		AuthorID: author.ID,
		ScopeID:  scope.ID,
	}
	err = db.Create(&subscription).Error
	require.NoError(t, err)

	am := NewAuthorMatcher(db)
	err = am.LoadSubscriptions(context.Background())
	require.NoError(t, err)

	result := am.FindMatchingSubscription([]string{"Pen Name"})
	require.NotNil(t, result)
	assert.Equal(t, subscription.ID, result.ID)
}

func TestFindMatchingSubscription_NoMatch(t *testing.T) {
	db, err := models.ConnectTestDB()
	require.NoError(t, err)

	am := NewAuthorMatcher(db)
	err = am.LoadSubscriptions(context.Background())
	require.NoError(t, err)

	result := am.FindMatchingSubscription([]string{"Unknown Author"})
	assert.Nil(t, result)
}

func TestFindMatchingSubscription_CaseInsensitive(t *testing.T) {
	db, err := models.ConnectTestDB()
	require.NoError(t, err)

	scope := models.SubscriptionScope{Name: "personal"}
	err = db.FirstOrCreate(&scope, models.SubscriptionScope{Name: "personal"}).Error
	require.NoError(t, err)

	author := models.Author{Name: "Test Author"}
	err = db.Create(&author).Error
	require.NoError(t, err)

	subscription := models.AuthorSubscription{
		AuthorID: author.ID,
		ScopeID:  scope.ID,
	}
	err = db.Create(&subscription).Error
	require.NoError(t, err)

	am := NewAuthorMatcher(db)
	err = am.LoadSubscriptions(context.Background())
	require.NoError(t, err)

	// Should match regardless of case
	result := am.FindMatchingSubscription([]string{"TEST AUTHOR"})
	require.NotNil(t, result)
	assert.Equal(t, subscription.ID, result.ID)

	result = am.FindMatchingSubscription([]string{"test author"})
	require.NotNil(t, result)
	assert.Equal(t, subscription.ID, result.ID)
}

func TestFindMatchingSubscription_DotInsensitive(t *testing.T) {
	db, err := models.ConnectTestDB()
	require.NoError(t, err)

	scope := models.SubscriptionScope{Name: "personal"}
	err = db.FirstOrCreate(&scope, models.SubscriptionScope{Name: "personal"}).Error
	require.NoError(t, err)

	// Author without dots
	author := models.Author{Name: "JF Brink"}
	err = db.Create(&author).Error
	require.NoError(t, err)

	subscription := models.AuthorSubscription{
		AuthorID: author.ID,
		ScopeID:  scope.ID,
	}
	err = db.Create(&subscription).Error
	require.NoError(t, err)

	am := NewAuthorMatcher(db)
	err = am.LoadSubscriptions(context.Background())
	require.NoError(t, err)

	// Should match even with dots in feed author name
	result := am.FindMatchingSubscription([]string{"J.F. Brink"})
	require.NotNil(t, result)
	assert.Equal(t, subscription.ID, result.ID)
}

func TestFindMatchingSubscription_MultipleAuthors(t *testing.T) {
	db, err := models.ConnectTestDB()
	require.NoError(t, err)

	scope := models.SubscriptionScope{Name: "personal"}
	err = db.FirstOrCreate(&scope, models.SubscriptionScope{Name: "personal"}).Error
	require.NoError(t, err)

	author := models.Author{Name: "Second Author"}
	err = db.Create(&author).Error
	require.NoError(t, err)

	subscription := models.AuthorSubscription{
		AuthorID: author.ID,
		ScopeID:  scope.ID,
	}
	err = db.Create(&subscription).Error
	require.NoError(t, err)

	am := NewAuthorMatcher(db)
	err = am.LoadSubscriptions(context.Background())
	require.NoError(t, err)

	// First author doesn't match, second does
	result := am.FindMatchingSubscription([]string{"Unknown Author", "Second Author"})
	require.NotNil(t, result)
	assert.Equal(t, subscription.ID, result.ID)
}
