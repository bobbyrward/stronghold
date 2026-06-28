package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/bobbyrward/stronghold/internal/catalog"
	"github.com/bobbyrward/stronghold/internal/config"
	"github.com/bobbyrward/stronghold/internal/hardcover"
	"github.com/bobbyrward/stronghold/internal/models"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

func createDoctorCmd() *cobra.Command {
	doctorCmd := &cobra.Command{
		Use:   "doctor",
		Short: "System diagnostics and maintenance commands",
	}

	doctorCmd.AddCommand(createDoctorMigrateCmd())
	doctorCmd.AddCommand(createDoctorInitBookSearchCmd())
	doctorCmd.AddCommand(createDoctorBackfillHardcoverRefsCmd())
	doctorCmd.AddCommand(createDoctorSyncBibliographyCmd())

	return doctorCmd
}

func createDoctorSyncBibliographyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sync-bibliography",
		Short: "Fetch each Hardcover-linked author's works and upsert them as Book rows",
		Long: `For every author with a HardcoverRef, fetch their bibliography from Hardcover
and upsert each work into the Book catalog (keyed on the Hardcover work id, so
re-runs update rather than duplicate). Does not create acquisition targets.`,
		RunE: runDoctorSyncBibliographyCmd,
	}
}

func runDoctorSyncBibliographyCmd(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	slog.InfoContext(ctx, "Syncing author bibliographies")

	db, err := models.ConnectDB()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	client := hardcover.NewClient(config.Config.Hardcover.ApiToken)

	synced, err := catalog.SyncAuthorBibliography(ctx, db, client)
	if err != nil {
		return err
	}

	fmt.Printf("Bibliography sync complete: %d books upserted\n", synced)
	return nil
}

func createDoctorBackfillHardcoverRefsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "backfill-hardcover-refs",
		Short: "Rewrite slug-based Author.HardcoverRef values to canonical Hardcover ids",
		Long: `Existing authors may have HardcoverRef set to a Hardcover slug. Slugs change and
die when Hardcover merges duplicates, so refs should hold the canonical author id.
This command resolves each slug ref to its canonical id and rewrites it. Refs that
already parse as an integer are left untouched (idempotent), and slugs that cannot
be resolved are reported and left untouched.`,
		RunE: runDoctorBackfillHardcoverRefsCmd,
	}
}

func runDoctorBackfillHardcoverRefsCmd(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	slog.InfoContext(ctx, "Backfilling Hardcover refs")

	db, err := models.ConnectDB()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	client := hardcover.NewClient(config.Config.Hardcover.ApiToken)

	rewritten, skipped, unresolved, err := backfillHardcoverRefs(ctx, db, client)
	if err != nil {
		return err
	}

	fmt.Printf("Backfill complete: %d rewritten, %d already ids, %d unresolved\n", rewritten, skipped, unresolved)
	return nil
}

// backfillHardcoverRefs rewrites slug-based Author.HardcoverRef values to canonical
// Hardcover ids. Refs that already parse as an int are skipped (idempotent); slugs
// that cannot be resolved are reported and left untouched.
func backfillHardcoverRefs(ctx context.Context, db *gorm.DB, client hardcover.Client) (rewritten, skipped, unresolved int, err error) {
	var authors []models.Author
	if err := db.Where("hardcover_ref IS NOT NULL").Find(&authors).Error; err != nil {
		return 0, 0, 0, fmt.Errorf("failed to load authors: %w", err)
	}

	for _, author := range authors {
		ref := *author.HardcoverRef

		// Already an id — idempotent skip.
		if _, convErr := strconv.Atoi(ref); convErr == nil {
			skipped++
			continue
		}

		// Spacing is owned by the client's rate limiter (internal/hardcover).
		result, getErr := client.GetAuthorBySlug(ctx, ref)
		if getErr != nil {
			return rewritten, skipped, unresolved, fmt.Errorf("failed to resolve slug %q for author %d: %w", ref, author.ID, getErr)
		}
		if result == nil {
			slog.WarnContext(ctx, "Could not resolve slug; leaving untouched", slog.Uint64("author_id", uint64(author.ID)), slog.String("slug", ref))
			fmt.Printf("UNRESOLVED: author %d (%s) slug %q\n", author.ID, author.Name, ref)
			unresolved++
			continue
		}

		if updErr := db.Model(&author).Update("hardcover_ref", result.ID).Error; updErr != nil {
			return rewritten, skipped, unresolved, fmt.Errorf("failed to update author %d: %w", author.ID, updErr)
		}
		slog.InfoContext(ctx, "Rewrote hardcover_ref", slog.Uint64("author_id", uint64(author.ID)), slog.String("slug", ref), slog.String("id", result.ID))
		rewritten++
	}

	return rewritten, skipped, unresolved, nil
}

func createDoctorMigrateCmd() *cobra.Command {
	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
		RunE:  runDoctorMigrateCmd,
	}

	return migrateCmd
}

func runDoctorMigrateCmd(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	slog.InfoContext(ctx, "Starting database migration")

	_, err := models.ConnectAndMigrate(ctx)
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, "Database migration completed successfully")
	return nil
}

func createDoctorInitBookSearchCmd() *cobra.Command {
	var apiKey string

	initBookSearchCmd := &cobra.Command{
		Use:   "init-book-search",
		Short: "Initialize book search credentials in the database",
		Long: `Initialize book search credentials by providing an API key (token cookie value).
This creates or updates the BookSearchCredential record in the database.
You can obtain the token cookie value from your browser after logging into the book search site.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDoctorInitBookSearchCmd(cmd, args, apiKey)
		},
	}

	initBookSearchCmd.Flags().StringVar(&apiKey, "api-key", "", "Book search API key (token cookie value) - required")
	_ = initBookSearchCmd.MarkFlagRequired("api-key")

	return initBookSearchCmd
}

func runDoctorInitBookSearchCmd(cmd *cobra.Command, args []string, apiKey string) error {
	ctx := context.Background()

	slog.InfoContext(ctx, "Initializing book search credentials")

	if apiKey == "" {
		return fmt.Errorf("api-key is required")
	}

	db, err := models.ConnectDB()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to connect to database", slog.Any("err", err))
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Initialize with empty IP and ASN - these will be populated on first token refresh
	err = models.UpsertBookSearchCredential(db, apiKey, "", "")
	if err != nil {
		slog.ErrorContext(ctx, "Failed to initialize book search credentials", slog.Any("err", err))
		return fmt.Errorf("failed to initialize book search credentials: %w", err)
	}

	slog.InfoContext(ctx, "Book search credentials initialized successfully")
	fmt.Println("Book search credentials initialized successfully")
	fmt.Println("You can now run 'stronghold refresh-token' to update IP address and ASN information")

	return nil
}
