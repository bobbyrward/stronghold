package catalog

import (
	"context"
	"testing"

	"github.com/bobbyrward/stronghold/internal/hardcover"
	"github.com/bobbyrward/stronghold/internal/models"
)

func ptr(s string) *string { return &s }

func TestSyncAuthorBibliography(t *testing.T) {
	db, err := models.ConnectTestDB()
	if err != nil {
		t.Fatalf("ConnectTestDB: %v", err)
	}

	linked := models.Author{Name: "Brandon Sanderson", HardcoverRef: ptr("100")}
	other := models.Author{Name: "N.K. Jemisin", HardcoverRef: ptr("200")}
	unlinked := models.Author{Name: "Unlinked", HardcoverRef: nil}
	for _, a := range []*models.Author{&linked, &other, &unlinked} {
		if err := db.Create(a).Error; err != nil {
			t.Fatalf("create author %s: %v", a.Name, err)
		}
	}

	// Work "10" is co-authored: it appears in both authors' bibliographies.
	coauthored := hardcover.BookResult{HardcoverID: "10", Title: "Co-Authored Work", ReleaseDate: nil}
	client := hardcover.NewMockClient()
	client.Books["100"] = []hardcover.BookResult{
		{HardcoverID: "1", Title: "The Way of Kings", ReleaseDate: ptr("2010-08-31")},
		{HardcoverID: "2", Title: "Words of Radiance", ReleaseDate: nil},
		coauthored,
	}
	client.Books["200"] = []hardcover.BookResult{
		{HardcoverID: "3", Title: "The Fifth Season", ReleaseDate: ptr("2015-08-04")},
		coauthored,
	}

	ctx := context.Background()

	synced, err := SyncAuthorBibliography(ctx, db, client)
	if err != nil {
		t.Fatalf("SyncAuthorBibliography: %v", err)
	}
	// Four distinct works (1, 2, 3, 10) — the co-authored work counts once.
	if synced != 4 {
		t.Fatalf("expected 4 books synced, got %d", synced)
	}

	var count int64
	db.Model(&models.Book{}).Count(&count)
	if count != 4 {
		t.Fatalf("expected 4 book rows, got %d", count)
	}

	// The co-authored work is a single Book linked to both authors.
	var shared models.Book
	if err := db.Preload("Authors").Where("hardcover_ref = ?", "10").First(&shared).Error; err != nil {
		t.Fatalf("find co-authored book: %v", err)
	}
	if len(shared.Authors) != 2 {
		t.Fatalf("expected co-authored book to have 2 authors, got %d", len(shared.Authors))
	}

	// release_date parsed for the dated book, nil passes through.
	var wok models.Book
	if err := db.Where("hardcover_ref = ?", "1").First(&wok).Error; err != nil {
		t.Fatalf("find book 1: %v", err)
	}
	if wok.ReleaseDate == nil || wok.ReleaseDate.Year() != 2010 {
		t.Fatalf("expected 2010 release date, got %v", wok.ReleaseDate)
	}
	var wor models.Book
	if err := db.Where("hardcover_ref = ?", "2").First(&wor).Error; err != nil {
		t.Fatalf("find book 2: %v", err)
	}
	if wor.ReleaseDate != nil {
		t.Fatalf("expected nil release date, got %v", wor.ReleaseDate)
	}

	// Second run upserts: a changed title updates in place, no duplicate rows or
	// author links.
	client.Books["100"][0].Title = "The Way of Kings (Revised)"
	if _, err := SyncAuthorBibliography(ctx, db, client); err != nil {
		t.Fatalf("second SyncAuthorBibliography: %v", err)
	}
	db.Model(&models.Book{}).Count(&count)
	if count != 4 {
		t.Fatalf("expected 4 book rows after re-sync, got %d", count)
	}
	// Co-authored work still has exactly two author links — Append stayed idempotent.
	if err := db.Preload("Authors").Where("hardcover_ref = ?", "10").First(&shared).Error; err != nil {
		t.Fatalf("re-find co-authored book: %v", err)
	}
	if len(shared.Authors) != 2 {
		t.Fatalf("expected 2 authors after re-sync, got %d", len(shared.Authors))
	}
	if err := db.Where("hardcover_ref = ?", "1").First(&wok).Error; err != nil {
		t.Fatalf("re-find book 1: %v", err)
	}
	if wok.Title != "The Way of Kings (Revised)" {
		t.Fatalf("expected updated title, got %q", wok.Title)
	}
}
