package audiobooks

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/autobrr/go-qbittorrent"
	"github.com/bobbyrward/stronghold/internal/config"
	"github.com/bobbyrward/stronghold/internal/importers/audiobooks/audible"
	"github.com/bobbyrward/stronghold/internal/importers/audiobooks/metadata"
	"github.com/bobbyrward/stronghold/internal/importers/audiobooks/torrent"
	"github.com/bobbyrward/stronghold/internal/importers/common"
	"github.com/bobbyrward/stronghold/internal/notifications"
	"github.com/bobbyrward/stronghold/internal/qbit"
	"github.com/cappuccinotm/slogx"
)

const (
	TITLE_SKIPPED = "__SKIPPED__"
)

type AudiobookImporterSystem struct {
	cfg              config.ImportersConfig
	qbitClient       qbit.QbitClient
	metadataProvider metadata.MetadataProvider
	audible          *audible.AudibleApiClient
}

func NewAudiobookImporterSystem(
	qbitClient qbit.QbitClient,
	cfg config.ImportersConfig,
	metadataProvider metadata.MetadataProvider,
	audibleApiClient *audible.AudibleApiClient,
) (*AudiobookImporterSystem, error) {
	importer := &AudiobookImporterSystem{
		cfg:              cfg,
		qbitClient:       qbitClient,
		metadataProvider: metadataProvider,
		audible:          audibleApiClient,
	}

	return importer, nil
}

func (abis *AudiobookImporterSystem) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "Running audiobook import process...")

	for _, importType := range abis.cfg.AudiobookImporter.ImportTypes {
		slog.InfoContext(ctx, "Processing import type", slog.String("category", importType.Category))

		library, ok := config.FindLibraryByName(abis.cfg.AudiobookImporter.Libraries, importType.Library)
		if !ok {
			return fmt.Errorf("unabled to find library: %s", importType.Library)
		}

		err := abis.ProcessImportType(ctx, importType, library)
		if err != nil {
			return errors.Join(err, fmt.Errorf("failed to process import type: %s", importType.Category))
		}
	}

	return nil
}

func (abis *AudiobookImporterSystem) ProcessImportType(ctx context.Context, importType config.ImportType, library *config.ImportLibrary) error {
	torrents, err := qbit.GetUnimportedTorrentsByCategory(
		ctx,
		abis.qbitClient,
		importType.Category,
	)
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to get unimported torrents for category: %s", importType.Category))
	}

	for _, torrent := range torrents {
		slog.InfoContext(ctx, "Found unimported torrent", slog.String("name", torrent.Name))
		abis.importTorrent(ctx, torrent, importType, library)
	}

	return nil
}

func (abis *AudiobookImporterSystem) MarkForManualIntervention(ctx context.Context, importTorrent qbittorrent.Torrent) {
	abis.MarkForManualInterventionWithNotification(ctx, importTorrent, "", "")
}

// MarkForManualInterventionWithNotification marks a torrent for manual intervention and sends a notification.
func (abis *AudiobookImporterSystem) MarkForManualInterventionWithNotification(ctx context.Context, importTorrent qbittorrent.Torrent, notifierName string, reason string) {
	err := qbit.TagTorrent(ctx, abis.qbitClient, importTorrent, config.Config.Importers.ManualInterventionTag)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to add manual intervention tag",
			slog.String("name", importTorrent.Name),
			slog.String("hash", importTorrent.Hash),
			slogx.Error(err),
		)
		return
	}

	slog.InfoContext(ctx, "marked torrent for manual intervention",
		slog.String("name", importTorrent.Name),
		slog.String("hash", importTorrent.Hash),
		slog.String("reason", reason),
	)

	// Send notification if notifier is configured
	if notifierName != "" {
		abis.sendManualInterventionNotification(ctx, importTorrent, notifierName, reason)
	}
}

func (abis *AudiobookImporterSystem) sendManualInterventionNotification(ctx context.Context, importTorrent qbittorrent.Torrent, notifierName string, reason string) {
	message := notifications.DiscordWebhookMessage{
		Username: "Stronghold Audiobook Importer",
		Embeds: []notifications.DiscordEmbed{
			{
				Title:       "⚠️ Manual Intervention Required",
				Description: fmt.Sprintf("Audiobook **%s** requires manual intervention", importTorrent.Name),
				Color:       0xFFA500, // Orange color
				Fields: []notifications.DiscordEmbedField{
					{
						Name:   "Reason",
						Value:  reason,
						Inline: false,
					},
					{
						Name:   "Torrent Hash",
						Value:  importTorrent.Hash,
						Inline: true,
					},
				},
			},
		},
	}

	err := notifications.SendNotification(ctx, notifierName, message)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to send manual intervention notification",
			slog.String("torrent", importTorrent.Name),
			slogx.Error(err))
	}
}

func (abis *AudiobookImporterSystem) MarkAsImported(ctx context.Context, importTorrent qbittorrent.Torrent) {
	err := qbit.TagTorrent(ctx, abis.qbitClient, importTorrent, config.Config.Importers.ImportedTag)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to mark as imported",
			slog.String("name", importTorrent.Name),
			slog.String("hash", importTorrent.Hash),
			slogx.Error(err),
		)
	} else {
		slog.InfoContext(ctx, "marked torrent as imported",
			slog.String("name", importTorrent.Name),
			slog.String("hash", importTorrent.Hash),
		)
	}
}

func (abis *AudiobookImporterSystem) SendDiscordNotification(ctx context.Context, bookMetadata metadata.BookMetadata, importType config.ImportType) {
	if importType.DiscordNotifier == "" {
		return
	}

	// Build description with title and series info
	description := fmt.Sprintf("**%s**", bookMetadata.Title)
	if bookMetadata.PrimarySeries != nil {
		description += fmt.Sprintf(" - %s", bookMetadata.PrimarySeries.Name)
		if bookMetadata.PrimarySeries.Position != nil {
			description += fmt.Sprintf(" - Book %s", *bookMetadata.PrimarySeries.Position)
		}
	}

	// Build author list
	authorNames := make([]string, len(bookMetadata.Authors))
	for i, author := range bookMetadata.Authors {
		authorNames[i] = author.Name
	}
	authorsStr := strings.Join(authorNames, ", ")

	// Build fields
	fields := []notifications.DiscordEmbedField{
		{
			Name:   "Author(s)",
			Value:  authorsStr,
			Inline: false,
		},
	}

	// Add series field if applicable
	if bookMetadata.PrimarySeries != nil {
		seriesStr := bookMetadata.PrimarySeries.Name
		if bookMetadata.PrimarySeries.Position != nil {
			seriesStr += fmt.Sprintf(" - Book %s", *bookMetadata.PrimarySeries.Position)
		}
		fields = append(fields, notifications.DiscordEmbedField{
			Name:   "Series",
			Value:  seriesStr,
			Inline: true,
		})
	}

	// Add Audible link
	audibleURL := fmt.Sprintf("https://www.audible.com/pd/%s", bookMetadata.Asin)
	fields = append(fields, notifications.DiscordEmbedField{
		Name:   "Audible",
		Value:  fmt.Sprintf("[View on Audible](%s)", audibleURL),
		Inline: true,
	})

	message := notifications.DiscordWebhookMessage{
		Username: "Stronghold Audiobook Importer",
		Embeds: []notifications.DiscordEmbed{
			{
				Title:       "🎧 New Audiobook Imported",
				Description: description,
				Color:       0x00ff00,
				Fields:      fields,
			},
		},
	}

	err := notifications.SendNotification(ctx, importType.DiscordNotifier, message)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to send Discord notification", slog.String("title", bookMetadata.Title), slog.Any("err", err))
	}
}

// ExtractTorrentMetadata extracts metadata from a torrent by analyzing its files
// and looking up book information via ASIN or title/author tags.
// Returns BookMetadata or an error if metadata cannot be extracted.
func (abis *AudiobookImporterSystem) ExtractTorrentMetadata(ctx context.Context, importTorrent qbittorrent.Torrent) (metadata.BookMetadata, error) {
	var bookMetadata metadata.BookMetadata

	files, err := common.MapTorrentFilesToLocalPaths(ctx, abis.qbitClient, importTorrent)
	if err != nil {
		slog.InfoContext(ctx, "Failed to map torrent files",
			slog.String("name", importTorrent.Name),
			slog.String("hash", importTorrent.Hash),
			slogx.Error(err),
		)
		return bookMetadata, err
	}

	torrentMetadata, err := torrent.NewAudiobookFilesMetadata(ctx, abis.qbitClient, abis.metadataProvider, importTorrent, files)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get audiobook files metadata",
			slog.String("name", importTorrent.Name),
			slog.String("hash", importTorrent.Hash),
			slog.Any("files", files),
			slogx.Error(err),
		)
		return bookMetadata, err
	}

	// Try ASIN lookup first
	asin, ok := torrentMetadata.Tags().AudibleASIN()
	if ok {
		bookMetadata, err = abis.lookupMetadataByAsin(ctx, asin)
		if err == nil {
			slog.InfoContext(ctx, "book metadata found by ASIN",
				slog.String("asin", asin),
				slog.String("title", bookMetadata.Title),
				slog.String("summary", bookMetadata.Summarize()),
			)
			return bookMetadata, nil
		}

		slog.ErrorContext(ctx, "Failed to lookup by ASIN, falling back to title lookup",
			slog.String("name", importTorrent.Name),
			slog.String("hash", importTorrent.Hash),
			slog.Any("files", files),
			slogx.Error(err),
		)
	}

	// Fall back to title/author lookup
	slog.InfoContext(ctx, "lookup by title and author tags", slog.Any("tags", torrentMetadata.Tags()))

	title, titleOk := torrentMetadata.Tags().Title()
	author, _ := torrentMetadata.Tags().Artist()

	if !titleOk {
		slog.ErrorContext(ctx, "Title tag not found in torrent metadata tags",
			slog.String("name", importTorrent.Name),
			slog.String("hash", importTorrent.Hash),
		)
		return bookMetadata, errors.New("title tag not found in torrent metadata tags")
	}

	bookMetadata, err = abis.lookupMetadataByTitle(ctx, title, author)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to lookup metadata by title",
			slog.String("name", importTorrent.Name),
			slog.String("hash", importTorrent.Hash),
			slog.String("title", title),
			slog.String("author", author),
		)
		return bookMetadata, err
	}

	slog.InfoContext(ctx, "Successfully extracted metadata", slog.Any("files", files), slog.Any("bookMetadata", bookMetadata))

	return bookMetadata, nil
}

// ExecuteImport performs the actual import operation: moving files to the library and writing metadata.
// Returns the destination path and any error encountered during the import process.
func (abis *AudiobookImporterSystem) ExecuteImport(ctx context.Context, importTorrent qbittorrent.Torrent, bookMetadata metadata.BookMetadata, library *config.ImportLibrary, localPath string) (string, error) {
	directoryName, err := bookMetadata.GenerateDirectoryName()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to generate directory name from metadata",
			slog.String("name", importTorrent.Name),
			slog.String("hash", importTorrent.Hash),
			slogx.Error(err),
		)
		return "", fmt.Errorf("failed to generate directory name: %w", err)
	}

	directoryName = sanitizeName(directoryName)
	fullDirName := path.Join(library.Path, directoryName)

	localPathInfo, err := os.Stat(localPath)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to stat local path",
			slog.String("name", importTorrent.Name),
			slog.String("hash", importTorrent.Hash),
			slog.String("localPath", localPath),
			slogx.Error(err),
		)
		return "", fmt.Errorf("failed to stat local path: %w", err)
	}

	if localPathInfo.IsDir() {
		err = moveFolderToLibrary(ctx, localPath, fullDirName)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to move folder to library",
				slog.String("name", importTorrent.Name),
				slog.String("hash", importTorrent.Hash),
				slog.String("localPath", localPath),
				slog.String("fullDirName", fullDirName),
				slogx.Error(err),
			)
			return "", fmt.Errorf("failed to move folder to library: %w", err)
		}
	} else {
		err = moveFileToLibrary(ctx, localPath, path.Base(localPath), fullDirName)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to move file to library",
				slog.String("name", importTorrent.Name),
				slog.String("hash", importTorrent.Hash),
				slog.String("localPath", localPath),
				slog.String("fullDirName", fullDirName),
				slogx.Error(err),
			)
			return "", fmt.Errorf("failed to move file to library: %w", err)
		}
	}

	err = bookMetadata.WriteOpf(path.Join(fullDirName, "metadata.opf"))
	if err != nil {
		slog.ErrorContext(ctx, "Failed to write opf",
			slog.String("name", importTorrent.Name),
			slog.String("hash", importTorrent.Hash),
			slogx.Error(err),
		)
		return "", fmt.Errorf("failed to write OPF metadata: %w", err)
	}

	slog.InfoContext(ctx, "Successfully imported audiobook",
		slog.String("name", importTorrent.Name),
		slog.String("destination", fullDirName),
	)

	return fullDirName, nil
}

func (abis *AudiobookImporterSystem) importTorrent(ctx context.Context, importTorrent qbittorrent.Torrent, importType config.ImportType, library *config.ImportLibrary) {
	abis.ImportTorrentWithLibrary(ctx, importTorrent, importType, library)
}

// ImportTorrentWithLibrary imports a single audiobook torrent using the specified library.
// This is the public entry point for external callers like the AuthorSubscriptionImporter.
func (abis *AudiobookImporterSystem) ImportTorrentWithLibrary(ctx context.Context, importTorrent qbittorrent.Torrent, importType config.ImportType, library *config.ImportLibrary) {
	bookMetadata, err := abis.ExtractTorrentMetadata(ctx, importTorrent)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to extract metadata for torrent",
			slog.String("name", importTorrent.Name),
			slog.String("hash", importTorrent.Hash),
			slogx.Error(err),
		)

		abis.MarkForManualInterventionWithNotification(ctx, importTorrent, importType.DiscordNotifier, "Failed to extract metadata: "+err.Error())

		return
	}

	localPath := common.MapTorrentContentPathToLocalPath(importTorrent, config.Config.Qbit.DownloadPath, config.Config.Qbit.LocalDownloadPath)

	_, err = abis.ExecuteImport(ctx, importTorrent, bookMetadata, library, localPath)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to execute import for torrent",
			slog.String("name", importTorrent.Name),
			slog.String("hash", importTorrent.Hash),
			slogx.Error(err),
		)

		abis.MarkForManualInterventionWithNotification(ctx, importTorrent, importType.DiscordNotifier, "Failed to execute import: "+err.Error())

		return
	}

	abis.MarkAsImported(ctx, importTorrent)
	abis.SendDiscordNotification(ctx, bookMetadata, importType)
}

func (abis *AudiobookImporterSystem) lookupMetadataByAsin(ctx context.Context, asin string) (metadata.BookMetadata, error) {
	md := metadata.BookMetadata{}

	metadatas, err := abis.getBookMetadataFromASINs(ctx, []string{asin})
	if err != nil {
		return md, err
	}

	if len(metadatas) == 0 {
		return md, errors.New("no metadata found for ASIN")
	}

	return metadatas[0], nil
}

func (abis *AudiobookImporterSystem) getBookMetadataFromASINs(ctx context.Context, asins []string) ([]metadata.BookMetadata, error) {
	metadatas := make([]metadata.BookMetadata, 0, len(asins))

	for _, asin := range asins {
		md, err := abis.audible.GetMetadataFromAsin(asin)
		if err != nil {
			slog.WarnContext(ctx, "Failed to get metadata from ASIN", slog.String("asin", asin), slogx.Error(err))
			continue
		}

		metadatas = append(metadatas, md)
	}

	return metadatas, nil
}

func (abis *AudiobookImporterSystem) lookupMetadataByTitle(ctx context.Context, title string, author string) (metadata.BookMetadata, error) {
	md := metadata.BookMetadata{}

	asins, err := abis.audible.SearchByTitle(title, author)
	if err != nil {
		return md, err
	}

	switch len(asins) {
	case 0:
		return md, errors.New("no ASINs found for title")
	case 1:
		return abis.lookupMetadataByAsin(ctx, asins[0])
	default:
		asinMetadatas, _ := abis.getBookMetadataFromASINs(ctx, asins)

		var summaries []string

		if len(asinMetadatas) > 0 {
			summaries, _ = summarizeBookMetadatas(asinMetadatas)
		}

		slog.Info("Multiple ASINs found for title, manual selection required", slog.Any("asins", asins), slog.Any("summaries", summaries))

		return md, errors.New("multiple ASINs found for title, manual selection required")
	}
}

func summarizeBookMetadatas(metadatas []metadata.BookMetadata) ([]string, error) {
	bookChoices := make([]string, 0, len(metadatas))

	for _, md := range metadatas {
		bookChoices = append(bookChoices, md.Summarize())
	}

	return bookChoices, nil
}

func sanitizeName(name string) string {
	return strings.ReplaceAll(name, "/", "-")
}

func moveFileToLibrary(ctx context.Context, sourceFile string, baseName string, destDirectory string) error {
	err := os.MkdirAll(destDirectory, 0777)
	if err != nil {
		return err
	}

	fullDestination := path.Join(destDirectory, baseName)

	slog.InfoContext(ctx, "copying file to library", slog.String("source", sourceFile), slog.String("destination", fullDestination))

	err = os.Link(sourceFile, fullDestination)
	if err != nil {
		fmt.Printf("Hard link failed.  Falling back to copy\n")

		err = copyFile(ctx, sourceFile, fullDestination)
		if err != nil {
			return err
		}
	}

	return nil
}

func copyFile(ctx context.Context, srcFilename string, destFilename string) error {
	sourceFile, err := os.Open(srcFilename)
	if err != nil {
		return err
	}
	defer func() { _ = sourceFile.Close() }()

	destinationFile, err := os.Create(destFilename)
	if err != nil {
		return err
	}
	defer func() { _ = destinationFile.Close() }()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

func moveFolderToLibrary(ctx context.Context, sourceDirectory string, destDirectory string) error {
	slog.InfoContext(ctx, "linking directory", slog.String("source", sourceDirectory), slog.String("destination", destDirectory))

	cmd := exec.Command("cp", "-al", sourceDirectory, destDirectory)

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Hard link failed.  Falling back to recursive copy\n")

		cmd = exec.Command("cp", "-r", sourceDirectory, destDirectory)

		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

/*
func autoSelectMetadata(ctx context.Context, sourceInfo SourceInfo) (metadata.BookMetadata, error) {
	var selectedMetadata metadata.BookMetadata

	if sourceInfo.IsM4b {
		slog.InfoContext(ctx, "Source file is M4B format", slog.String("path", sourceInfo.BookFiles[0].LocalPath), slog.Any("sourceInfo", sourceInfo))

		tagList, err := metadata.GetM4BTagList(context.Background(), sourceInfo.BookFiles[0].LocalPath)
		if err != nil {
			return selectedMetadata, err
		}

		asin, err := tagList.GetString("AUDIBLE_ASIN")
		if asin != "" && err == nil {
			slog.InfoContext(ctx, "Found AUDIBLE_ASIN tag from AUDIBLE_ASIN tag", slog.String("asin", asin))

			selectedMetadata, err = lookupMetadataByAsin(asin)
			if err == nil && selectedMetadata.Title != TITLE_SKIPPED {
				return selectedMetadata, nil
			}
		}

		title, err := tagList.GetString("title")
		if title != "" && err == nil {
			author, _ := tagList.GetString("artist")

			selectedMetadata, err = lookupMetadataByTitle(title, author)

			if err == nil && selectedMetadata.Title != TITLE_SKIPPED {
				return selectedMetadata, nil
			}
		}
	} else {
		slog.InfoContext(ctx, "Source file is not M4B format, assuming MP3", slog.String("path", sourceInfo.BookFiles[0].LocalPath), slog.Any("sourceInfo", sourceInfo))

		tagList, err := metadata.GetMp3TagList(context.Background(), sourceInfo.BookFiles[0].LocalPath)
		if err != nil {
			return selectedMetadata, err
		}

		woas, err := tagList.GetString("AUDIBLE_ASIN")
		if woas != "" && err == nil && strings.HasPrefix(woas, "http://www.audible.com/pd/") {
			woas = strings.TrimPrefix(woas, "http://www.audible.com/pd/")

			selectedMetadata, err = lookupMetadataByAsin(woas)
			if err == nil && selectedMetadata.Title != TITLE_SKIPPED {
				return selectedMetadata, nil
			}
		}

		title, err := tagList.GetString("title")
		if title != "" && err == nil {
			author, _ := tagList.GetString("artist")

			selectedMetadata, err = lookupMetadataByTitle(title, author)

			if err == nil && selectedMetadata.Title != TITLE_SKIPPED {
				return selectedMetadata, nil
			}
		}
	}

	prompt := promptui.Prompt{
		Label: "Book Title",
	}

	title, err := prompt.Run()
	if err != nil {
		return selectedMetadata, err
	}

	return LookupMetadataByTitle(title, "")
	return selectedMetadata, nil // DELETE ME
}

func lookupMetadataByAsin(asin string) (metadata.BookMetadata, error) {
	aac := audible.NewAudibleApiClient()
	md := metadata.BookMetadata{}

	metadatas, err := getBookMetadataFromASINs(aac, []string{asin})
	if err != nil {
		return md, err
	}

	return metadatas[0], nil
}


/*
func GetSourcePathFromArgs(args []string) string {
	return args[0]
}

func GetLibraryConfig(libraryName string) (*config.AbsImportConfigLibrary, error) {
	for _, library := range config.Config.Importers.AbsImportConfig.Libraries {
		if library.Name == libraryName {
			return &library, nil
		}
	}

	return nil, fmt.Errorf("library `%s` not found", libraryName)
}

func SummarizeBookMetadatas(metadatas []metadata.BookMetadata) ([]string, error) {
	bookChoices := make([]string, 0, len(metadatas))

	for _, md := range metadatas {
		bookChoices = append(bookChoices, md.Summarize())
	}

	return bookChoices, nil
}

func SelectBook(books []string) (int, error) {
	prompt := promptui.Select{
		Label: "Select a book",
		Items: books,
	}

	index, _, err := prompt.Run()
	if err != nil {
		return 0, err
	}

	return index, nil
}
*/
