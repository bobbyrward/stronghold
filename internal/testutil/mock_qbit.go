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

	// RemoveTagsCtx mocking
	RemoveTagsCtxFunc   func(ctx context.Context, hashes []string, tags string) error
	RemoveTagsCtxCalls  []RemoveTagsCall
	RemoveTagsCtxReturn error

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

	// SetCategoryCtx mocking
	SetCategoryCtxFunc   func(ctx context.Context, hashes []string, category string) error
	SetCategoryCtxCalls  []SetCategoryCtxCall
	SetCategoryCtxReturn error

	// SetTags mocking
	SetTagsCtxFunc   func(ctx context.Context, hashes []string, category string) error
	SetTagsCtxCalls  []SetTagsCtxCall
	SetTagsCtxReturn error
}

type AddTagsCall struct {
	Hashes []string
	Tags   string
}

type RemoveTagsCall struct {
	Hashes []string
	Tags   string
}

type AddTorrentFromUrlCall struct {
	URL     string
	Options map[string]string
}

type SetCategoryCtxCall struct {
	Hashes   []string
	Category string
}

type SetTagsCtxCall struct {
	Hashes []string
	Tags   string
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

func (m *MockQbitClient) RemoveTagsCtx(ctx context.Context, hashes []string, tags string) error {
	m.RemoveTagsCtxCalls = append(m.RemoveTagsCtxCalls, RemoveTagsCall{
		Hashes: hashes,
		Tags:   tags,
	})

	if m.RemoveTagsCtxFunc != nil {
		return m.RemoveTagsCtxFunc(ctx, hashes, tags)
	}

	return m.RemoveTagsCtxReturn
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

func (m *MockQbitClient) SetCategoryCtx(ctx context.Context, hashes []string, category string) error {
	m.SetCategoryCtxCalls = append(m.SetCategoryCtxCalls, SetCategoryCtxCall{
		Hashes:   hashes,
		Category: category,
	})

	if m.SetCategoryCtxFunc != nil {
		return m.SetCategoryCtxFunc(ctx, hashes, category)
	}

	return m.SetCategoryCtxReturn
}

func (m *MockQbitClient) SetTags(ctx context.Context, hashes []string, category string) error {
	m.SetTagsCtxCalls = append(m.SetTagsCtxCalls, SetTagsCtxCall{
		Hashes: hashes,
		Tags:   category,
	})

	if m.SetTagsCtxFunc != nil {
		return m.SetTagsCtxFunc(ctx, hashes, category)
	}

	return m.SetTagsCtxReturn
}
