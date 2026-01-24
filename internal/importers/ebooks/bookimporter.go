package ebooks

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/autobrr/go-qbittorrent"
	"github.com/cappuccinotm/slogx"

	"github.com/bobbyrward/stronghold/internal/config"
	"github.com/bobbyrward/stronghold/internal/importers/common"
	"github.com/bobbyrward/stronghold/internal/notifications"
	"github.com/bobbyrward/stronghold/internal/qbit"
)

type BookImporterSystem struct {
	qbitClient qbit.QbitClient
}

func NewBookImporterSystem(qbitClient qbit.QbitClient) *BookImporterSystem {
	return &BookImporterSystem{
		qbitClient: qbitClient,
	}
}

func (bis *BookImporterSystem) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "Running book import process...")

	for _, importType := range config.Config.Importers.BookImporter.ImportTypes {
		slog.InfoContext(ctx, "Processing import type", slog.String("category", importType.Category))

		library, ok := config.FindLibraryByName(config.Config.Importers.BookImporter.Libraries, importType.Library)
		if !ok {
			return fmt.Errorf("unabled to find library: %s", importType.Library)
		}

		err := bis.ProcessImportType(ctx, importType, library)
		if err != nil {
			return errors.Join(err, fmt.Errorf("failed to process import type: %s", importType.Category))
		}
	}

	return nil
}

func (bis *BookImporterSystem) ProcessImportType(ctx context.Context, importType config.ImportType, library *config.ImportLibrary) error {
	torrents, err := qbit.GetUnimportedTorrentsByCategory(ctx, bis.qbitClient, importType.Category)
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to get unimported torrents for category: %s", importType.Category))
	}

	for _, torrent := range torrents {
		slog.InfoContext(ctx, "Found unimported torrent", slog.String("name", torrent.Name))
		bis.ImportTorrent(ctx, torrent, importType, library)
	}

	return nil
}

func (bis *BookImporterSystem) ImportTorrent(ctx context.Context, torrent qbittorrent.Torrent, importType config.ImportType, library *config.ImportLibrary) {
	files, err := common.MapTorrentFilesToLocalPaths(ctx, bis.qbitClient, torrent)
	if err != nil {
		slog.InfoContext(ctx, "Failed to map torrent files",
			slog.String("name", torrent.Name),
			slog.String("hash", torrent.Hash),
			slogx.Error(err),
		)

		bis.markForManualIntervention(ctx, torrent, importType.DiscordNotifier, "Failed to map torrent files: "+err.Error())
		return
	}

	books := make([]common.MappedTorrentFile, 0, len(files))

	for _, mappedFile := range files {
		switch filepath.Ext(mappedFile.BaseName) {
		case ".azw3":
			fallthrough
		case ".mobi":
			fallthrough
		case ".epub":
			books = append(books, mappedFile)
		}
	}

	if len(books) == 0 {
		slog.InfoContext(ctx, "Unable to find epubs in torrent", slog.String("name", torrent.Name))

		bis.markForManualIntervention(ctx, torrent, importType.DiscordNotifier, "No ebook files (.epub, .mobi, .azw3) found in torrent")
		return
	}

	slog.InfoContext(ctx, "Found epubs", slog.Int("count", len(books)))

	for _, mappedFile := range books {
		slog.InfoContext(ctx, "Copying file", slog.Any("mappedFile", mappedFile))

		// Use only the base filename to flatten directory structure
		destPath := filepath.Join(library.Path, filepath.Base(mappedFile.BaseName))
		err = copyFile(mappedFile.LocalPath, destPath)
		if err != nil {
			slog.InfoContext(ctx, "Unable to copy file", slog.Any("mappedFile", mappedFile), slog.String("name", torrent.Name), slog.Any("err", err))

			bis.markForManualIntervention(ctx, torrent, importType.DiscordNotifier, "Failed to copy file: "+err.Error())
			return
		}
	}

	err = qbit.TagTorrent(ctx, bis.qbitClient, torrent, config.Config.Importers.ImportedTag)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to add imported tag", slog.String("name", torrent.Name), slog.Any("err", err))
	}

	bis.sendDiscordNotification(ctx, torrent, library, books, importType)
}

func (bis *BookImporterSystem) markForManualIntervention(ctx context.Context, torrent qbittorrent.Torrent, notifierName string, reason string) {
	err := qbit.TagTorrent(ctx, bis.qbitClient, torrent, config.Config.Importers.ManualInterventionTag)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to add manual intervention tag",
			slog.String("name", torrent.Name),
			slog.String("hash", torrent.Hash),
			slogx.Error(err),
		)
		return
	}

	slog.InfoContext(ctx, "Marked torrent for manual intervention",
		slog.String("name", torrent.Name),
		slog.String("hash", torrent.Hash),
		slog.String("reason", reason),
	)

	// Send notification if notifier is configured
	if notifierName != "" {
		bis.sendManualInterventionNotification(ctx, torrent, notifierName, reason)
	}
}

func (bis *BookImporterSystem) sendManualInterventionNotification(ctx context.Context, torrent qbittorrent.Torrent, notifierName string, reason string) {
	message := notifications.DiscordWebhookMessage{
		Username: "Stronghold Book Importer",
		Embeds: []notifications.DiscordEmbed{
			{
				Title:       "⚠️ Manual Intervention Required",
				Description: fmt.Sprintf("Ebook **%s** requires manual intervention", torrent.Name),
				Color:       0xFFA500, // Orange color
				Fields: []notifications.DiscordEmbedField{
					{
						Name:   "Reason",
						Value:  reason,
						Inline: false,
					},
					{
						Name:   "Torrent Hash",
						Value:  torrent.Hash,
						Inline: true,
					},
				},
			},
		},
	}

	err := notifications.SendNotification(ctx, notifierName, message)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to send manual intervention notification",
			slog.String("torrent", torrent.Name),
			slogx.Error(err))
	}
}

func (bis *BookImporterSystem) sendDiscordNotification(ctx context.Context, torrent qbittorrent.Torrent, library *config.ImportLibrary, books []common.MappedTorrentFile, importType config.ImportType) {
	if importType.DiscordNotifier == "" {
		return
	}

	bookList := ""
	for i, book := range books {
		if i > 0 {
			bookList += "\n"
		}
		bookList += "• " + book.BaseName
	}

	message := notifications.DiscordWebhookMessage{
		Username: "Stronghold Book Importer",
		Embeds: []notifications.DiscordEmbed{
			{
				Title:       "📚 New Book(s) Imported",
				Description: fmt.Sprintf("Successfully imported %d book(s) from torrent **%s**", len(books), torrent.Name),
				Color:       0x00ff00,
				Fields: []notifications.DiscordEmbedField{
					{
						Name:   "Books",
						Value:  bookList,
						Inline: false,
					},
					{
						Name:   "Category",
						Value:  importType.Category,
						Inline: true,
					},
					{
						Name:   "Destination",
						Value:  library.Path,
						Inline: true,
					},
				},
			},
		},
	}

	err := notifications.SendNotification(ctx, importType.DiscordNotifier, message)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to send Discord notification", slog.String("torrent", torrent.Name), slog.Any("err", err))
	}
}

func copyFile(sourcePath string, destPath string) error {
	slog.Info("Copying file", slog.String("source", sourcePath), slog.String("dest", destPath))

	source, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer func() { _ = source.Close() }()

	// Create parent directory if it doesn't exist
	destDir := filepath.Dir(destPath)
	err = os.MkdirAll(destDir, 0755)
	if err != nil {
		return err
	}

	destination, err := os.Create(destPath)
	if err != nil {
		return err
	}

	_, err = io.Copy(destination, source)
	if err != nil {
		_ = destination.Close()
		return err
	}

	return destination.Close()
}
