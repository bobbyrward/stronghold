package qbit

import (
	"context"
	"errors"
	"fmt"
)

func GetTorrentsByCategory(ctx context.Context, qbit QbitClient, category string) ([]Torrent, error) {
	torrents, err := qbit.GetTorrentsCtx(ctx, TorrentFilterOptions{
		Category: category,
	})
	if err != nil {
		return nil, errors.Join(err, fmt.Errorf("failed to get torrents"))
	}

	return torrents, nil
}

func GetUnimportedTorrentsByCategory(ctx context.Context, qbit QbitClient, category string) ([]Torrent, error) {
	torrents, err := qbit.GetTorrentsCtx(ctx, TorrentFilterOptions{
		Category: category,
	})
	if err != nil {
		return nil, errors.Join(err, fmt.Errorf("failed to get torrents"))
	}

	return FilterUnimportedTorrents(torrents), nil
}
