// Package catalog materializes the book catalog from external sources. It sits
// above models and hardcover so the same logic is reusable by the doctor CLI, a
// future scheduler, and web refresh endpoints.
package catalog

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/bobbyrward/stronghold/internal/hardcover"
	"github.com/bobbyrward/stronghold/internal/models"
	"gorm.io/gorm"
)

// releaseDateLayout is Hardcover's release_date format (BookResult.ReleaseDate).
const releaseDateLayout = "2006-01-02"

// SyncAuthorBibliography fetches the bibliography of every Hardcover-linked author
// and upserts each work into the Book catalog, keyed on the Hardcover work id so
// re-runs update rather than duplicate. A co-authored work is a single Book linked
// to every tracked contributor (many-to-many). It does not create
// AcquisitionTargets. Returns the number of distinct works synced.
func SyncAuthorBibliography(ctx context.Context, db *gorm.DB, client hardcover.Client) (synced int, err error) {
	var authors []models.Author
	if err := db.Where("hardcover_ref IS NOT NULL").Find(&authors).Error; err != nil {
		return 0, fmt.Errorf("failed to load authors: %w", err)
	}

	slog.InfoContext(ctx, "Syncing author bibliographies", slog.Int("authors", len(authors)))

	// A work co-authored by two tracked authors is returned by both their
	// GetAuthorBooks calls; count it once.
	seen := make(map[string]struct{})

	for _, author := range authors {
		ref := *author.HardcoverRef

		books, getErr := client.GetAuthorBooks(ctx, ref)
		if getErr != nil {
			return synced, fmt.Errorf("failed to fetch bibliography for author %d (%s): %w", author.ID, ref, getErr)
		}

		for _, b := range books {
			book, upErr := upsertBook(ctx, db, b)
			if upErr != nil {
				return synced, upErr
			}

			// Link this author to the work. GORM upserts the join row, so the
			// append is idempotent across re-syncs and co-authors.
			if linkErr := db.Model(book).Association("Authors").Append(&author); linkErr != nil {
				return synced, fmt.Errorf("failed to link author %d to book %s: %w", author.ID, b.HardcoverID, linkErr)
			}

			if _, ok := seen[b.HardcoverID]; !ok {
				seen[b.HardcoverID] = struct{}{}
				synced++
			}
		}

		slog.InfoContext(ctx, "Synced author bibliography",
			slog.Uint64("author_id", uint64(author.ID)),
			slog.String("hardcover_ref", ref),
			slog.Int("books", len(books)))
	}

	slog.InfoContext(ctx, "Bibliography sync complete", slog.Int("books", synced))
	return synced, nil
}

// upsertBook finds the Book for a Hardcover work by its ref, updating title and
// release date, or creates it if absent. Returns the persisted row (with ID set)
// so the caller can attach author associations.
func upsertBook(ctx context.Context, db *gorm.DB, b hardcover.BookResult) (*models.Book, error) {
	var book models.Book
	res := db.Where("hardcover_ref = ?", b.HardcoverID).First(&book)

	switch {
	case errors.Is(res.Error, gorm.ErrRecordNotFound):
		book = models.Book{
			HardcoverRef: &b.HardcoverID,
			Title:        b.Title,
			ReleaseDate:  parseReleaseDate(ctx, b.ReleaseDate),
		}
		if err := db.Create(&book).Error; err != nil {
			return nil, fmt.Errorf("failed to create book %q (%s): %w", b.Title, b.HardcoverID, err)
		}
	case res.Error != nil:
		return nil, fmt.Errorf("failed to load book %s: %w", b.HardcoverID, res.Error)
	default:
		book.Title = b.Title
		book.ReleaseDate = parseReleaseDate(ctx, b.ReleaseDate)
		if err := db.Save(&book).Error; err != nil {
			return nil, fmt.Errorf("failed to update book %q (%s): %w", b.Title, b.HardcoverID, err)
		}
	}

	return &book, nil
}

// parseReleaseDate parses a Hardcover release_date; nil and unparseable values
// pass through as nil rather than failing the whole sync.
func parseReleaseDate(ctx context.Context, raw *string) *time.Time {
	if raw == nil || *raw == "" {
		return nil
	}
	t, err := time.Parse(releaseDateLayout, *raw)
	if err != nil {
		slog.WarnContext(ctx, "Unparseable release_date; storing nil", slog.String("value", *raw))
		return nil
	}
	return &t
}
