package testutil

import (
	"context"

	"github.com/autobrr/go-qbittorrent"
)

// MockQbitClient is a reusable mock implementation of qbit.QbitClient for testing
type MockQbitClient struct {
	// GetTorrentsCtx mocking
	GetTorrentsCtxFunc   func(ctx context.Context, o qbittorrent.TorrentFilterOptions) ([]qbittorrent.Torrent, error)
	GetTorrentsCtxCalls  []qbittorrent.TorrentFilterOptions
	GetTorrentsCtxReturn struct {
		Torrents []qbittorrent.Torrent
		Err      error
	}

	// AddTagsCtx mocking
	AddTagsCtxFunc   func(ctx context.Context, hashes []string, tags string) error
	AddTagsCtxCalls  []AddTagsCall
	AddTagsCtxReturn error

	// GetFilesInformationCtx mocking
	GetFilesInformationCtxFunc   func(ctx context.Context, hash string) (*qbittorrent.TorrentFiles, error)
	GetFilesInformationCtxCalls  []string
	GetFilesInformationCtxReturn struct {
		Files *qbittorrent.TorrentFiles
		Err   error
	}

	// AddTorrentFromUrlCtx mocking
	AddTorrentFromUrlCtxFunc   func(ctx context.Context, url string, options map[string]string) error
	AddTorrentFromUrlCtxCalls  []AddTorrentFromUrlCall
	AddTorrentFromUrlCtxReturn error
}

type AddTagsCall struct {
	Hashes []string
	Tags   string
}

type AddTorrentFromUrlCall struct {
	URL     string
	Options map[string]string
}

func (m *MockQbitClient) GetTorrentsCtx(ctx context.Context, o qbittorrent.TorrentFilterOptions) ([]qbittorrent.Torrent, error) {
	m.GetTorrentsCtxCalls = append(m.GetTorrentsCtxCalls, o)

	if m.GetTorrentsCtxFunc != nil {
		return m.GetTorrentsCtxFunc(ctx, o)
	}

	return m.GetTorrentsCtxReturn.Torrents, m.GetTorrentsCtxReturn.Err
}

func (m *MockQbitClient) AddTagsCtx(ctx context.Context, hashes []string, tags string) error {
	m.AddTagsCtxCalls = append(m.AddTagsCtxCalls, AddTagsCall{
		Hashes: hashes,
		Tags:   tags,
	})

	if m.AddTagsCtxFunc != nil {
		return m.AddTagsCtxFunc(ctx, hashes, tags)
	}

	return m.AddTagsCtxReturn
}

func (m *MockQbitClient) GetFilesInformationCtx(ctx context.Context, hash string) (*qbittorrent.TorrentFiles, error) {
	m.GetFilesInformationCtxCalls = append(m.GetFilesInformationCtxCalls, hash)

	if m.GetFilesInformationCtxFunc != nil {
		return m.GetFilesInformationCtxFunc(ctx, hash)
	}

	return m.GetFilesInformationCtxReturn.Files, m.GetFilesInformationCtxReturn.Err
}

func (m *MockQbitClient) AddTorrentFromUrlCtx(ctx context.Context, url string, options map[string]string) error {
	m.AddTorrentFromUrlCtxCalls = append(m.AddTorrentFromUrlCtxCalls, AddTorrentFromUrlCall{
		URL:     url,
		Options: options,
	})

	if m.AddTorrentFromUrlCtxFunc != nil {
		return m.AddTorrentFromUrlCtxFunc(ctx, url, options)
	}

	return m.AddTorrentFromUrlCtxReturn
}
