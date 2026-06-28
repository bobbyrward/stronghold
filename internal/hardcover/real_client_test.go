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

func TestBookRowToResult(t *testing.T) {
	t.Run("maps id, title and release date", func(t *testing.T) {
		date := "2017-11-14"
		got := bookRow{ID: 123, Title: "Oathbringer", ReleaseDate: &date}.toResult()
		if got.HardcoverID != "123" || got.Title != "Oathbringer" || got.ReleaseDate == nil || *got.ReleaseDate != date {
			t.Fatalf("unexpected result: %+v", got)
		}
	})

	t.Run("nil release date passes through", func(t *testing.T) {
		got := bookRow{ID: 1, Title: "Untitled"}.toResult()
		if got.ReleaseDate != nil {
			t.Fatalf("expected nil release date, got: %v", *got.ReleaseDate)
		}
	})
}
