package hardcover

import "testing"

func TestAuthorRowToResult(t *testing.T) {
	t.Run("no canonical uses own id/slug/name", func(t *testing.T) {
		row := authorRow{ID: 42, Name: "Brandon Sanderson", Slug: "brandon-sanderson"}
		got := row.toResult()
		if got.ID != "42" || got.Slug != "brandon-sanderson" || got.Name != "Brandon Sanderson" {
			t.Fatalf("unexpected result: %+v", got)
		}
	})

	t.Run("canonical (merged) resolves id/slug/name as a unit from canonical", func(t *testing.T) {
		row := authorRow{ID: 7, Name: "Old Name", Slug: "old-slug"}
		row.Canonical = &struct {
			ID   int    `graphql:"id"`
			Name string `graphql:"name"`
			Slug string `graphql:"slug"`
		}{ID: 99, Name: "New Name", Slug: "new-slug"}

		got := row.toResult()
		if got.ID != "99" || got.Slug != "new-slug" || got.Name != "New Name" {
			t.Fatalf("expected canonical values, got: %+v", got)
		}
	})
}
