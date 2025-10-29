package qbit

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

type FilteredTorrents struct {
	Filtered  []Torrent
	Remaining []Torrent
}

func TagsFromTagList(tagList string) map[string]bool {
	tags := make(map[string]bool)

	for _, tag := range strings.Split(tagList, ",") {
		tags[tag] = true
	}

	return tags
}

func FilterTorrentsByTag(torrents []Torrent, tagFilter string) FilteredTorrents {
	result := FilteredTorrents{
		Filtered:  make([]Torrent, 0),
		Remaining: make([]Torrent, 0),
	}

	for _, torrent := range torrents {
		tags := TagsFromTagList(torrent.Tags)

		if _, ok := tags[tagFilter]; ok {
			result.Filtered = append(result.Filtered, torrent)
		} else {
			result.Remaining = append(result.Remaining, torrent)
		}
	}

	return result
}

func FilterUnimportedTorrents(torrents []Torrent) []Torrent {
	unimported := make([]Torrent, 0, len(torrents))

	for _, torrent := range torrents {
		if torrent.Tags == "" {
			unimported = append(unimported, torrent)
		}
	}

	return unimported
}

func TagTorrent(ctx context.Context, client QbitClient, torrent Torrent, tag string) error {
	return client.AddTagsCtx(ctx, []string{torrent.Hash}, tag)
}

func GetUnimportedTorrents(ctx context.Context, qbit QbitClient) ([]Torrent, error) {
	torrents, err := qbit.GetTorrentsCtx(ctx, TorrentFilterOptions{Tag: ""})
	if err != nil {
		return nil, errors.Join(err, fmt.Errorf("failed to get torrents"))
	}

	return FilterUnimportedTorrents(torrents), nil
}

func GetManualInterventionTorrents(ctx context.Context, qbit QbitClient, manualInterventionTag string) ([]Torrent, error) {
	torrents, err := qbit.GetTorrentsCtx(ctx, TorrentFilterOptions{Tag: ""})
	if err != nil {
		return nil, errors.Join(err, fmt.Errorf("failed to get torrents"))
	}

	return FilterTorrentsByTag(torrents, manualInterventionTag).Filtered, nil
}
