package common

import (
	"context"
	"errors"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/autobrr/go-qbittorrent"

	"github.com/bobbyrward/stronghold/internal/config"
	"github.com/bobbyrward/stronghold/internal/qbit"
)

type MappedTorrentFile struct {
	BaseName  string
	LocalPath string
}

func MapTorrentContentPathToLocalPath(torrent qbittorrent.Torrent, remoteDownloadPath string, localDownloadPath string) string {
	return filepath.Join(
		localDownloadPath,
		strings.TrimPrefix(
			torrent.ContentPath,
			strings.TrimSuffix(remoteDownloadPath, "/"),
		),
	)
}

func MapTorrentSavePathToLocalPath(savePath string, remoteDownloadPath string, localDownloadPath string) string {
	return filepath.Join(
		localDownloadPath,
		strings.TrimPrefix(
			savePath,
			strings.TrimSuffix(remoteDownloadPath, "/"),
		),
	)
}

func MapTorrentFilesToLocalPaths(ctx context.Context, qbit qbit.QbitClient, torrent qbittorrent.Torrent) ([]MappedTorrentFile, error) {
	localSavePath := MapTorrentSavePathToLocalPath(
		torrent.SavePath,
		config.Config.Qbit.DownloadPath,
		config.Config.Qbit.LocalDownloadPath,
	)

	torrentFiles, err := qbit.GetFilesInformationCtx(ctx, torrent.Hash)
	if err != nil {
		msg := "failed to get torrent files"
		slog.InfoContext(ctx, msg, slog.String("name", torrent.Name), slog.Any("err", err))

		return nil, errors.Join(errors.New(msg), err)
	}

	files := make([]MappedTorrentFile, 0, len(*torrentFiles))

	for _, torrentFile := range *torrentFiles {
		files = append(files, MappedTorrentFile{
			BaseName:  torrentFile.Name,
			LocalPath: filepath.Join(localSavePath, torrentFile.Name),
		})
	}

	return files, nil
}
