package qbit

import (
	"context"
	"log/slog"

	"github.com/autobrr/go-qbittorrent"

	"github.com/bobbyrward/stronghold/internal/config"
)

type QbitClient interface {
	GetTorrentsCtx(ctx context.Context, o TorrentFilterOptions) ([]Torrent, error)
	AddTagsCtx(ctx context.Context, hashes []string, tags string) error
	RemoveTagsCtx(ctx context.Context, hashes []string, tags string) error
	GetFilesInformationCtx(ctx context.Context, hash string) (*TorrentFiles, error)
	AddTorrentFromUrlCtx(ctx context.Context, url string, options map[string]string) error
	SetCategoryCtx(ctx context.Context, hashes []string, category string) error
	SetTags(ctx context.Context, hashes []string, tags string) error
}

func CreateClient() (QbitClient, error) {
	ctx := context.Background()

	slog.InfoContext(ctx, "Creating qBittorrent client",
		slog.String("host", config.Config.Qbit.URL),
		slog.String("username", config.Config.Qbit.Username))

	client := qbittorrent.NewClient(qbittorrent.Config{
		Host:     config.Config.Qbit.URL,
		Username: config.Config.Qbit.Username,
		Password: config.Config.Qbit.Password,
	})

	slog.InfoContext(ctx, "Successfully created qBittorrent client")

	return client, nil
}
