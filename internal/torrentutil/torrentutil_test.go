package torrentutil

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackpal/bencode-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createTestTorrent creates a valid torrent file bytes for testing.
func createTestTorrent(t *testing.T, name string) []byte {
	t.Helper()

	torrent := map[string]interface{}{
		"info": map[string]interface{}{
			"name":         name,
			"piece length": 262144,
			"pieces":       "12345678901234567890", // 20 bytes (one piece hash)
			"length":       1024,
		},
		"announce": "http://tracker.example.com/announce",
	}

	var buf bytes.Buffer
	err := bencode.Marshal(&buf, torrent)
	require.NoError(t, err)

	return buf.Bytes()
}

func TestNewTorrentDownloader_NoProxy(t *testing.T) {
	td := NewTorrentDownloader("", "")
	assert.NotNil(t, td)
	assert.NotNil(t, td.httpClient)
}

func TestNewTorrentDownloader_WithProxy(t *testing.T) {
	td := NewTorrentDownloader("http://proxy.example.com:8080", "")
	assert.NotNil(t, td)
	assert.NotNil(t, td.httpClient)
}

func TestNewTorrentDownloader_InvalidProxy(t *testing.T) {
	// Invalid proxy URL should fall back to default client
	td := NewTorrentDownloader("://invalid", "")
	assert.NotNil(t, td)
	assert.NotNil(t, td.httpClient)
}

func TestDownloadAndHash(t *testing.T) {
	torrentData := createTestTorrent(t, "test-file.txt")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-bittorrent")
		_, _ = w.Write(torrentData)
	}))
	defer server.Close()

	td := NewTorrentDownloader("", "")
	hash, err := td.DownloadAndHash(context.Background(), server.URL+"/test.torrent")

	require.NoError(t, err)
	assert.Len(t, hash, 40) // SHA1 hex string is 40 characters
	assert.Regexp(t, "^[a-f0-9]{40}$", hash)
}

func TestDownloadAndHash_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	td := NewTorrentDownloader("", "")
	_, err := td.DownloadAndHash(context.Background(), server.URL+"/notfound.torrent")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected HTTP status: 404")
}

func TestDownloadAndHash_InvalidTorrent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-bittorrent")
		_, _ = w.Write([]byte("this is not a valid torrent file"))
	}))
	defer server.Close()

	td := NewTorrentDownloader("", "")
	_, err := td.DownloadAndHash(context.Background(), server.URL+"/invalid.torrent")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse torrent file")
}

func TestDownloadAndHash_ConnectionError(t *testing.T) {
	td := NewTorrentDownloader("", "")
	_, err := td.DownloadAndHash(context.Background(), "http://localhost:99999/test.torrent")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to download torrent")
}

func TestExtractInfoHash(t *testing.T) {
	torrentData := createTestTorrent(t, "test-file.txt")

	hash, err := ExtractInfoHash(torrentData)

	require.NoError(t, err)
	assert.Len(t, hash, 40)
	assert.Regexp(t, "^[a-f0-9]{40}$", hash)
}

func TestExtractInfoHash_MissingInfo(t *testing.T) {
	// Create torrent without info dict
	torrent := map[string]interface{}{
		"announce": "http://tracker.example.com/announce",
	}

	var buf bytes.Buffer
	err := bencode.Marshal(&buf, torrent)
	require.NoError(t, err)

	_, err = ExtractInfoHash(buf.Bytes())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing info dictionary")
}

func TestExtractInfoHash_InvalidBencode(t *testing.T) {
	_, err := ExtractInfoHash([]byte("not valid bencode"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse torrent file")
}

func TestExtractInfoHash_Deterministic(t *testing.T) {
	// Same torrent data should produce same hash
	torrentData := createTestTorrent(t, "deterministic-test.txt")

	hash1, err := ExtractInfoHash(torrentData)
	require.NoError(t, err)

	hash2, err := ExtractInfoHash(torrentData)
	require.NoError(t, err)

	assert.Equal(t, hash1, hash2)
}

func TestExtractInfoHash_DifferentTorrents(t *testing.T) {
	// Different torrents should produce different hashes
	torrent1 := createTestTorrent(t, "file1.txt")
	torrent2 := createTestTorrent(t, "file2.txt")

	hash1, err := ExtractInfoHash(torrent1)
	require.NoError(t, err)

	hash2, err := ExtractInfoHash(torrent2)
	require.NoError(t, err)

	assert.NotEqual(t, hash1, hash2)
}
