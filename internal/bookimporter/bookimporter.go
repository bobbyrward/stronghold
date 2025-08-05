package bookimporter

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/autobrr/go-qbittorrent"

	"github.com/bobbyrward/stronghold/internal/config"
	"github.com/bobbyrward/stronghold/internal/notifications"
	"github.com/bobbyrward/stronghold/internal/qbit"
)

type BookImporterSystem struct{}

func NewBookImporterSystem() *BookImporterSystem {
	return &BookImporterSystem{}
}

func (bis *BookImporterSystem) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "Running book import process...")

	for _, importType := range config.Config.BookImporter.ImportTypes {
		slog.InfoContext(ctx, "Processing import type", slog.String("category", importType.Category))

		err := bis.ProcessImportType(ctx, importType)
		if err != nil {
			return errors.Join(err, fmt.Errorf("failed to process import type: %s", importType.Category))
		}
	}

	return nil
}

func (bis *BookImporterSystem) ProcessImportType(ctx context.Context, importType config.BookImportTypeConfig) error {
	client, err := qbit.CreateClient()
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to create qBittorrent client"))
	}

	torrents, err := client.GetTorrentsCtx(ctx, qbittorrent.TorrentFilterOptions{
		Category: importType.Category,
	})
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to get torrents"))
	}

	for _, torrent := range torrents {
		tags := qbit.TagsFromTagList(torrent.Tags)

		if _, ok := tags[importType.ImportedTag]; ok {
			continue
		}

		if _, ok := tags[importType.ManualInterventionTag]; ok {
			continue
		}

		slog.InfoContext(ctx, "Found unimported torrent", slog.String("name", torrent.Name))
		bis.ImportTorrent(ctx, client, torrent, importType)
	}

	return nil
}

func (bis *BookImporterSystem) ImportTorrent(ctx context.Context, client *qbittorrent.Client, torrent qbittorrent.Torrent, importType config.BookImportTypeConfig) {
	unprefixedSavePath := strings.TrimPrefix(torrent.SavePath, importType.SourcePrefixPath)
	localSavePath := filepath.Join(importType.SourcePath, unprefixedSavePath)

	torrentFiles, err := client.GetFilesInformationCtx(ctx, torrent.Hash)
	if err != nil {
		slog.InfoContext(ctx, "Failed to get torrent files", slog.String("name", torrent.Name), slog.Any("err", err))

		err = client.AddTagsCtx(ctx, []string{torrent.Hash}, importType.ManualInterventionTag)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to add manual intervention tag", slog.String("name", torrent.Name), slog.Any("err", err))
		}

		return
	}

	books := make([]string, 0)

	for _, torrentFile := range *torrentFiles {
		if strings.HasSuffix(torrentFile.Name, ".epub") {
			books = append(books, torrentFile.Name)
			continue
		}

		if strings.HasSuffix(torrentFile.Name, ".mobi") {
			books = append(books, torrentFile.Name)
			continue
		}

		if strings.HasSuffix(torrentFile.Name, ".azw3") {
			books = append(books, torrentFile.Name)
			continue
		}
	}

	if len(books) == 0 {
		slog.InfoContext(ctx, "Unable to find epubs in torrent", slog.String("name", torrent.Name), slog.Any("err", err))

		err = client.AddTagsCtx(ctx, []string{torrent.Hash}, importType.ManualInterventionTag)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to add manual intervention tag", slog.String("name", torrent.Name), slog.Any("err", err))
		}

		return
	}

	slog.InfoContext(ctx, "Found epubs", slog.Int("count", len(books)))

	for _, filename := range books {
		slog.InfoContext(ctx, "Copying file", slog.String("filename", filename))

		err = copyFile(filepath.Join(localSavePath, filename), filepath.Join(importType.DestinationPath, filepath.Base(filename)))
		if err != nil {
			slog.InfoContext(ctx, "Unable to copy file", slog.String("filename", filename), slog.String("name", torrent.Name), slog.Any("err", err))

			err = client.AddTagsCtx(ctx, []string{torrent.Hash}, importType.ManualInterventionTag)
			if err != nil {
				slog.ErrorContext(ctx, "Failed to add manual intervention tag", slog.String("name", torrent.Name), slog.Any("err", err))
			}
			return
		}
	}

	err = client.AddTagsCtx(ctx, []string{torrent.Hash}, importType.ImportedTag)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to add imported tag", slog.String("name", torrent.Name), slog.Any("err", err))
	}

	bis.sendDiscordNotification(ctx, torrent, books, importType)
}

func (bis *BookImporterSystem) sendDiscordNotification(ctx context.Context, torrent qbittorrent.Torrent, books []string, importType config.BookImportTypeConfig) {
	if importType.DiscordNotifier == "" {
		return
	}

	bookList := ""
	for i, book := range books {
		if i > 0 {
			bookList += "\n"
		}
		bookList += "• " + filepath.Base(book)
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
						Value:  importType.DestinationPath,
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
