package common

import (
	"context"
	"errors"
	"testing"

	"github.com/autobrr/go-qbittorrent"
	"github.com/stretchr/testify/assert"

	"github.com/bobbyrward/stronghold/internal/config"
	"github.com/bobbyrward/stronghold/internal/testutil"
)

func TestMapTorrentSavePathToLocalPath(t *testing.T) {
	tests := []struct {
		name               string
		savePath           string
		remoteDownloadPath string
		localDownloadPath  string
		expectedLocalPath  string
	}{
		{
			name:               "Basic path mapping",
			savePath:           "/data/books",
			remoteDownloadPath: "/data",
			localDownloadPath:  "/mnt/downloads/",
			expectedLocalPath:  "/mnt/downloads/books",
		},
		{
			name:               "Remote path with trailing slash",
			savePath:           "/data/audiobooks",
			remoteDownloadPath: "/data/",
			localDownloadPath:  "/mnt/downloads",
			expectedLocalPath:  "/mnt/downloads/audiobooks",
		},
		{
			name:               "Nested subdirectory",
			savePath:           "/remote/media/audiobooks/fiction",
			remoteDownloadPath: "/remote/media",
			localDownloadPath:  "/local",
			expectedLocalPath:  "/local/audiobooks/fiction",
		},
		{
			name:               "Single level subdirectory",
			savePath:           "/downloads/books",
			remoteDownloadPath: "/downloads",
			localDownloadPath:  "/home/user/media",
			expectedLocalPath:  "/home/user/media/books",
		},
		{
			name:               "Same path returns local base",
			savePath:           "/data",
			remoteDownloadPath: "/data",
			localDownloadPath:  "/local",
			expectedLocalPath:  "/local",
		},
		{
			name:               "Multiple nested levels",
			savePath:           "/mnt/storage/downloads/2024/audiobooks/fiction",
			remoteDownloadPath: "/mnt/storage/downloads",
			localDownloadPath:  "/home/media",
			expectedLocalPath:  "/home/media/2024/audiobooks/fiction",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MapTorrentSavePathToLocalPath(
				tt.savePath,
				tt.remoteDownloadPath,
				tt.localDownloadPath,
			)
			assert.Equal(t, tt.expectedLocalPath, result)
		})
	}
}

func TestMapTorrentFilesToLocalPaths(t *testing.T) {
	tests := []struct {
		name          string
		torrent       qbittorrent.Torrent
		mockFiles     *qbittorrent.TorrentFiles
		mockErr       error
		qbitConfig    config.QbitConfig
		expectedFiles []MappedTorrentFile
		expectedError bool
	}{
		{
			name: "Single file torrent",
			torrent: qbittorrent.Torrent{
				Hash:     "abc123",
				Name:     "Test Book",
				SavePath: "/data/books",
			},
			mockFiles: &qbittorrent.TorrentFiles{
				{Name: "Book.epub"},
			},
			mockErr: nil,
			qbitConfig: config.QbitConfig{
				DownloadPath:      "/data",
				LocalDownloadPath: "/mnt/downloads/",
			},
			expectedFiles: []MappedTorrentFile{
				{
					BaseName:  "Book.epub",
					LocalPath: "/mnt/downloads/books/Book.epub",
				},
			},
			expectedError: false,
		},
		{
			name: "Multiple files in torrent",
			torrent: qbittorrent.Torrent{
				Hash:     "def456",
				Name:     "Audio Series",
				SavePath: "/remote/audiobooks/series",
			},
			mockFiles: &qbittorrent.TorrentFiles{
				{Name: "book1.m4b"},
				{Name: "book2.m4b"},
				{Name: "book3.m4b"},
			},
			mockErr: nil,
			qbitConfig: config.QbitConfig{
				DownloadPath:      "/remote/audiobooks",
				LocalDownloadPath: "/home/media",
			},
			expectedFiles: []MappedTorrentFile{
				{
					BaseName:  "book1.m4b",
					LocalPath: "/home/media/series/book1.m4b",
				},
				{
					BaseName:  "book2.m4b",
					LocalPath: "/home/media/series/book2.m4b",
				},
				{
					BaseName:  "book3.m4b",
					LocalPath: "/home/media/series/book3.m4b",
				},
			},
			expectedError: false,
		},
		{
			name: "Files with subdirectories in torrent",
			torrent: qbittorrent.Torrent{
				Hash:     "ghi789",
				Name:     "Collection",
				SavePath: "/data/downloads/collection",
			},
			mockFiles: &qbittorrent.TorrentFiles{
				{Name: "subfolder/file1.m4b"},
				{Name: "subfolder/file2.mp3"},
				{Name: "other/file3.m4b"},
			},
			mockErr: nil,
			qbitConfig: config.QbitConfig{
				DownloadPath:      "/data/downloads",
				LocalDownloadPath: "/mnt/local",
			},
			expectedFiles: []MappedTorrentFile{
				{
					BaseName:  "subfolder/file1.m4b",
					LocalPath: "/mnt/local/collection/subfolder/file1.m4b",
				},
				{
					BaseName:  "subfolder/file2.mp3",
					LocalPath: "/mnt/local/collection/subfolder/file2.mp3",
				},
				{
					BaseName:  "other/file3.m4b",
					LocalPath: "/mnt/local/collection/other/file3.m4b",
				},
			},
			expectedError: false,
		},
		{
			name: "Torrent saved at remote base path",
			torrent: qbittorrent.Torrent{
				Hash:     "jkl012",
				Name:     "Root Book",
				SavePath: "/downloads",
			},
			mockFiles: &qbittorrent.TorrentFiles{
				{Name: "root-book.epub"},
			},
			mockErr: nil,
			qbitConfig: config.QbitConfig{
				DownloadPath:      "/downloads",
				LocalDownloadPath: "/local/media",
			},
			expectedFiles: []MappedTorrentFile{
				{
					BaseName:  "root-book.epub",
					LocalPath: "/local/media/root-book.epub",
				},
			},
			expectedError: false,
		},
		{
			name: "Empty torrent (no files)",
			torrent: qbittorrent.Torrent{
				Hash:     "empty123",
				Name:     "Empty Torrent",
				SavePath: "/data/empty",
			},
			mockFiles: &qbittorrent.TorrentFiles{},
			mockErr:   nil,
			qbitConfig: config.QbitConfig{
				DownloadPath:      "/data",
				LocalDownloadPath: "/local",
			},
			expectedFiles: []MappedTorrentFile{},
			expectedError: false,
		},
		{
			name: "Error getting files information",
			torrent: qbittorrent.Torrent{
				Hash:     "error123",
				Name:     "Failed Torrent",
				SavePath: "/data/failed",
			},
			mockFiles: nil,
			mockErr:   errors.New("connection timeout"),
			qbitConfig: config.QbitConfig{
				DownloadPath:      "/data",
				LocalDownloadPath: "/local",
			},
			expectedFiles: nil,
			expectedError: true,
		},
		{
			name: "Deeply nested torrent path",
			torrent: qbittorrent.Torrent{
				Hash:     "nested456",
				Name:     "Nested Book",
				SavePath: "/mnt/storage/torrents/2024/fiction/mystery",
			},
			mockFiles: &qbittorrent.TorrentFiles{
				{Name: "book.m4b"},
			},
			mockErr: nil,
			qbitConfig: config.QbitConfig{
				DownloadPath:      "/mnt/storage/torrents",
				LocalDownloadPath: "/home/user/audiobooks",
			},
			expectedFiles: []MappedTorrentFile{
				{
					BaseName:  "book.m4b",
					LocalPath: "/home/user/audiobooks/2024/fiction/mystery/book.m4b",
				},
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up config for this test
			config.Config.Qbit = tt.qbitConfig

			// Create mock client
			mockClient := &testutil.MockQbitClient{
				GetFilesInformationCtxReturn: struct {
					Files *qbittorrent.TorrentFiles
					Err   error
				}{
					Files: tt.mockFiles,
					Err:   tt.mockErr,
				},
			}

			// Call the function
			ctx := context.Background()
			files, err := MapTorrentFilesToLocalPaths(ctx, mockClient, tt.torrent)

			// Verify expectations
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, files)
				assert.Contains(t, err.Error(), "failed to get torrent files")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedFiles, files)
			}

			// Verify mock was called correctly
			assert.Len(t, mockClient.GetFilesInformationCtxCalls, 1)
			assert.Equal(t, tt.torrent.Hash, mockClient.GetFilesInformationCtxCalls[0])
		})
	}
}

func TestMapTorrentFilesToLocalPaths_WithCustomMockFunction(t *testing.T) {
	// This test demonstrates using a custom mock function for more complex scenarios
	torrent := qbittorrent.Torrent{
		Hash:     "custom123",
		Name:     "Custom Test",
		SavePath: "/remote/books/custom",
	}

	config.Config.Qbit = config.QbitConfig{
		DownloadPath:      "/remote/books",
		LocalDownloadPath: "/local",
	}

	mockClient := &testutil.MockQbitClient{}
	mockClient.GetFilesInformationCtxFunc = func(ctx context.Context, hash string) (*qbittorrent.TorrentFiles, error) {
		// Custom logic - could do different things based on hash
		if hash == "custom123" {
			return &qbittorrent.TorrentFiles{
				{Name: "custom-file.m4b"},
			}, nil
		}
		return nil, errors.New("unexpected hash")
	}

	ctx := context.Background()
	files, err := MapTorrentFilesToLocalPaths(ctx, mockClient, torrent)

	assert.NoError(t, err)
	assert.Len(t, files, 1)
	assert.Equal(t, "custom-file.m4b", files[0].BaseName)
	assert.Equal(t, "/local/custom/custom-file.m4b", files[0].LocalPath)
}
