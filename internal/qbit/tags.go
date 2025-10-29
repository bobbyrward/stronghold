package qbit

import (
	"context"
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

func FilterUnimportedTorrents(torrents []Torrent, importedTag string, manualInterventionTag string) []Torrent {
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
