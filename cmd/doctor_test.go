package cmd

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bobbyrward/stronghold/internal/hardcover"
	"github.com/bobbyrward/stronghold/internal/models"
)

func strptr(s string) *string { return &s }

func TestBackfillHardcoverRefs(t *testing.T) {
	db, err := models.ConnectTestDB()
	require.NoError(t, err)

	hc := hardcover.NewMockClient()
	hc.AddAuthor("1", "brandon-sanderson", "Brandon Sanderson")

	slugAuthor := models.Author{Name: "Brandon Sanderson", HardcoverRef: strptr("brandon-sanderson")}
	idAuthor := models.Author{Name: "Already Id", HardcoverRef: strptr("42")}
	unknownAuthor := models.Author{Name: "Dead Slug", HardcoverRef: strptr("gone-slug")}
	noRefAuthor := models.Author{Name: "No Ref"}
	require.NoError(t, db.Create(&slugAuthor).Error)
	require.NoError(t, db.Create(&idAuthor).Error)
	require.NoError(t, db.Create(&unknownAuthor).Error)
	require.NoError(t, db.Create(&noRefAuthor).Error)

	rewritten, skipped, unresolved, err := backfillHardcoverRefs(context.Background(), db, hc, 0)
	require.NoError(t, err)
	assert.Equal(t, 1, rewritten)
	assert.Equal(t, 1, skipped)
	assert.Equal(t, 1, unresolved)

	var gotSlug, gotID, gotUnknown models.Author
	require.NoError(t, db.First(&gotSlug, slugAuthor.ID).Error)
	assert.Equal(t, "1", *gotSlug.HardcoverRef, "slug ref should be rewritten to canonical id")

	require.NoError(t, db.First(&gotID, idAuthor.ID).Error)
	assert.Equal(t, "42", *gotID.HardcoverRef, "numeric ref should be untouched")

	require.NoError(t, db.First(&gotUnknown, unknownAuthor.ID).Error)
	assert.Equal(t, "gone-slug", *gotUnknown.HardcoverRef, "unresolved slug should be left untouched")
}
