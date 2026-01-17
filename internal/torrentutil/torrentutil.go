package torrentutil

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/jackpal/bencode-go"
)

// TorrentDownloader downloads torrent files and extracts info hashes.
type TorrentDownloader struct {
	httpClient *http.Client
}

// NewTestTorrentDownloader creates a TorrentDownloader with no proxy for testing purposes.
func NewTestTorrentDownloader() *TorrentDownloader {
	return &TorrentDownloader{
		httpClient: &http.Client{},
	}
}

// NewTorrentDownloader creates a new TorrentDownloader with HTTP/HTTPS proxy support.
// Both httpProxy and httpsProxy are required and should be hostname:port (e.g., "myhost:8080").
func NewTorrentDownloader(httpProxy, httpsProxy string) *TorrentDownloader {
	return &TorrentDownloader{
		httpClient: &http.Client{
			Transport: &http.Transport{
				Proxy: func(req *http.Request) (*url.URL, error) {
					if req.URL.Scheme == "https" {
						return &url.URL{Scheme: "http", Host: httpsProxy}, nil
					}
					return &url.URL{Scheme: "http", Host: httpProxy}, nil
				},
			},
		},
	}
}

// DownloadAndHash downloads a torrent file from the given URL and returns its info hash.
// The info hash is the SHA1 hash of the bencoded "info" dictionary.
func (td *TorrentDownloader) DownloadAndHash(ctx context.Context, torrentURL string) (string, error) {
	slog.DebugContext(ctx, "Downloading torrent", slog.String("url", torrentURL))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, torrentURL, nil)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create HTTP request",
			slog.String("url", torrentURL),
			slog.Any("error", err))
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := td.httpClient.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to download torrent",
			slog.String("url", torrentURL),
			slog.Any("error", err))
		return "", fmt.Errorf("failed to download torrent: %w", err)
	}
	defer func() {
		// Drain body before closing for connection reuse
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		slog.ErrorContext(ctx, "Unexpected HTTP status",
			slog.String("url", torrentURL),
			slog.Int("status", resp.StatusCode))
		return "", fmt.Errorf("unexpected HTTP status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to read response body",
			slog.String("url", torrentURL),
			slog.Any("error", err))
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	slog.DebugContext(ctx, "Downloaded torrent", slog.Int("bytes", len(body)))

	hash, err := ExtractInfoHash(body)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to extract info hash",
			slog.String("url", torrentURL),
			slog.Any("error", err))
		return "", err
	}

	slog.DebugContext(ctx, "Extracted info hash",
		slog.String("url", torrentURL),
		slog.String("hash", hash))

	return hash, nil
}

// ExtractInfoHash extracts the info hash from raw torrent file bytes.
func ExtractInfoHash(torrentData []byte) (string, error) {
	reader := bytes.NewReader(torrentData)

	// Use bencode.Decode to parse into a generic map (Unmarshal with struct doesn't work for nested dicts)
	decoded, err := bencode.Decode(reader)
	if err != nil {
		return "", fmt.Errorf("failed to parse torrent file: %w", err)
	}

	torrent, ok := decoded.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("torrent file is not a dictionary")
	}

	info, ok := torrent["info"]
	if !ok || info == nil {
		return "", fmt.Errorf("torrent file missing info dictionary")
	}

	// Re-encode the info dictionary to get the exact bytes for hashing
	var infoBuf bytes.Buffer
	err = bencode.Marshal(&infoBuf, info)
	if err != nil {
		return "", fmt.Errorf("failed to re-encode info dictionary: %w", err)
	}

	// SHA1 hash the bencoded info dictionary
	hasher := sha1.New()
	hasher.Write(infoBuf.Bytes())
	hash := hasher.Sum(nil)

	return hex.EncodeToString(hash), nil
}
