package qbit

import (
	"context"
	"fmt"
)

func GetTorrentsByCategory(ctx context.Context, qbit QbitClient, category string) ([]Torrent, error) {
	torrents, err := qbit.GetTorrentsCtx(ctx, TorrentFilterOptions{
		Category: category,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get torrents: %w", err)
	}

	return torrents, nil
}

func GetUnimportedTorrentsByCategory(ctx context.Context, qbit QbitClient, category string) ([]Torrent, error) {
	torrents, err := qbit.GetTorrentsCtx(ctx, TorrentFilterOptions{
		Category: category,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get torrents: %w", err)
	}

	return FilterUnimportedTorrents(torrents), nil
}
