package ebooks

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/autobrr/go-qbittorrent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bobbyrward/stronghold/internal/config"
	"github.com/bobbyrward/stronghold/internal/testutil"
)

func TestImportTorrent_SingleFileInRoot(t *testing.T) {
	ctx := context.Background()

	// Create temporary directories
	tempDir := t.TempDir()
	sourceDir := filepath.Join(tempDir, "source")
	destDir := filepath.Join(tempDir, "dest")

	// Set up source file
	require.NoError(t, os.MkdirAll(sourceDir, 0755))
	require.NoError(t, os.MkdirAll(destDir, 0755))
	sourceFile := filepath.Join(sourceDir, "book.epub")
	require.NoError(t, os.WriteFile(sourceFile, []byte("test epub content"), 0644))

	// Set up config
	config.Config.Qbit = config.QbitConfig{
		DownloadPath:      "/remote",
		LocalDownloadPath: sourceDir,
	}
	config.Config.Importers.ImportedTag = "imported"
	config.Config.Importers.ManualInterventionTag = "manual"

	// Mock qBittorrent client
	torrent := qbittorrent.Torrent{
		Hash:     "test123",
		Name:     "Test Book",
		SavePath: "/remote",
	}

	mockClient := &testutil.MockQbitClient{
		GetFilesInformationCtxReturn: struct {
			Files *qbittorrent.TorrentFiles
			Err   error
		}{
			Files: &qbittorrent.TorrentFiles{
				{Name: "book.epub"},
			},
		},
	}

	library := &config.ImportLibrary{
		Name: "test-library",
		Path: destDir,
	}

	importType := config.ImportType{
		Category: "books",
		Library:  "test-library",
	}

	// Run import
	importer := NewBookImporterSystem(mockClient)
	importer.ImportTorrent(ctx, torrent, importType, library)

	// Verify file was copied to destination root
	destFile := filepath.Join(destDir, "book.epub")
	_, err := os.Stat(destFile)
	assert.NoError(t, err, "File should exist in destination root")

	// Verify content
	content, err := os.ReadFile(destFile)
	assert.NoError(t, err)
	assert.Equal(t, "test epub content", string(content))

	// Verify torrent was tagged as imported
	assert.Len(t, mockClient.AddTagsCtxCalls, 1)
	assert.Equal(t, "imported", mockClient.AddTagsCtxCalls[0].Tags)
}

func TestImportTorrent_SingleFileInNestedDirectory(t *testing.T) {
	ctx := context.Background()

	// Create temporary directories
	tempDir := t.TempDir()
	sourceDir := filepath.Join(tempDir, "source")
	destDir := filepath.Join(tempDir, "dest")

	// Set up source file in nested directory
	require.NoError(t, os.MkdirAll(filepath.Join(sourceDir, "Author/Book Title"), 0755))
	require.NoError(t, os.MkdirAll(destDir, 0755))
	sourceFile := filepath.Join(sourceDir, "Author/Book Title/book.epub")
	require.NoError(t, os.WriteFile(sourceFile, []byte("nested epub content"), 0644))

	// Set up config
	config.Config.Qbit = config.QbitConfig{
		DownloadPath:      "/remote",
		LocalDownloadPath: sourceDir,
	}
	config.Config.Importers.ImportedTag = "imported"
	config.Config.Importers.ManualInterventionTag = "manual"

	// Mock qBittorrent client
	torrent := qbittorrent.Torrent{
		Hash:     "test456",
		Name:     "Nested Book",
		SavePath: "/remote",
	}

	mockClient := &testutil.MockQbitClient{
		GetFilesInformationCtxReturn: struct {
			Files *qbittorrent.TorrentFiles
			Err   error
		}{
			Files: &qbittorrent.TorrentFiles{
				{Name: "Author/Book Title/book.epub"},
			},
		},
	}

	library := &config.ImportLibrary{
		Name: "test-library",
		Path: destDir,
	}

	importType := config.ImportType{
		Category: "books",
		Library:  "test-library",
	}

	// Run import
	importer := NewBookImporterSystem(mockClient)
	importer.ImportTorrent(ctx, torrent, importType, library)

	// Verify file was flattened to destination root
	destFile := filepath.Join(destDir, "book.epub")
	_, err := os.Stat(destFile)
	assert.NoError(t, err, "File should be flattened to destination root")

	// Verify no subdirectories were created
	entries, err := os.ReadDir(destDir)
	require.NoError(t, err)
	for _, entry := range entries {
		assert.False(t, entry.IsDir(), "No subdirectories should be created")
	}

	// Verify content
	content, err := os.ReadFile(destFile)
	assert.NoError(t, err)
	assert.Equal(t, "nested epub content", string(content))

	// Verify torrent was tagged as imported
	assert.Len(t, mockClient.AddTagsCtxCalls, 1)
	assert.Equal(t, "imported", mockClient.AddTagsCtxCalls[0].Tags)
}

func TestImportTorrent_MultipleFilesInNestedDirectories(t *testing.T) {
	ctx := context.Background()

	// Create temporary directories
	tempDir := t.TempDir()
	sourceDir := filepath.Join(tempDir, "source")
	destDir := filepath.Join(tempDir, "dest")

	// Set up multiple source files in nested directories
	require.NoError(t, os.MkdirAll(filepath.Join(sourceDir, "Author1/Book1"), 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(sourceDir, "Author2/Series/Book2"), 0755))
	require.NoError(t, os.MkdirAll(destDir, 0755))

	file1 := filepath.Join(sourceDir, "Author1/Book1/book1.epub")
	file2 := filepath.Join(sourceDir, "Author2/Series/Book2/book2.mobi")
	require.NoError(t, os.WriteFile(file1, []byte("book1 content"), 0644))
	require.NoError(t, os.WriteFile(file2, []byte("book2 content"), 0644))

	// Set up config
	config.Config.Qbit = config.QbitConfig{
		DownloadPath:      "/remote",
		LocalDownloadPath: sourceDir,
	}
	config.Config.Importers.ImportedTag = "imported"
	config.Config.Importers.ManualInterventionTag = "manual"

	// Mock qBittorrent client
	torrent := qbittorrent.Torrent{
		Hash:     "test789",
		Name:     "Multiple Books",
		SavePath: "/remote",
	}

	mockClient := &testutil.MockQbitClient{
		GetFilesInformationCtxReturn: struct {
			Files *qbittorrent.TorrentFiles
			Err   error
		}{
			Files: &qbittorrent.TorrentFiles{
				{Name: "Author1/Book1/book1.epub"},
				{Name: "Author2/Series/Book2/book2.mobi"},
			},
		},
	}

	library := &config.ImportLibrary{
		Name: "test-library",
		Path: destDir,
	}

	importType := config.ImportType{
		Category: "books",
		Library:  "test-library",
	}

	// Run import
	importer := NewBookImporterSystem(mockClient)
	importer.ImportTorrent(ctx, torrent, importType, library)

	// Verify both files were flattened to destination root
	destFile1 := filepath.Join(destDir, "book1.epub")
	destFile2 := filepath.Join(destDir, "book2.mobi")

	_, err := os.Stat(destFile1)
	assert.NoError(t, err, "book1.epub should be in destination root")

	_, err = os.Stat(destFile2)
	assert.NoError(t, err, "book2.mobi should be in destination root")

	// Verify no subdirectories were created
	entries, err := os.ReadDir(destDir)
	require.NoError(t, err)
	for _, entry := range entries {
		assert.False(t, entry.IsDir(), "No subdirectories should be created")
	}

	// Verify torrent was tagged as imported
	assert.Len(t, mockClient.AddTagsCtxCalls, 1)
	assert.Equal(t, "imported", mockClient.AddTagsCtxCalls[0].Tags)
}

func TestImportTorrent_MultipleFileFormats(t *testing.T) {
	ctx := context.Background()

	// Create temporary directories
	tempDir := t.TempDir()
	sourceDir := filepath.Join(tempDir, "source")
	destDir := filepath.Join(tempDir, "dest")

	// Set up source files with different formats
	require.NoError(t, os.MkdirAll(sourceDir, 0755))
	require.NoError(t, os.MkdirAll(destDir, 0755))

	epubFile := filepath.Join(sourceDir, "book.epub")
	mobiFile := filepath.Join(sourceDir, "book.mobi")
	azw3File := filepath.Join(sourceDir, "book.azw3")

	require.NoError(t, os.WriteFile(epubFile, []byte("epub content"), 0644))
	require.NoError(t, os.WriteFile(mobiFile, []byte("mobi content"), 0644))
	require.NoError(t, os.WriteFile(azw3File, []byte("azw3 content"), 0644))

	// Set up config
	config.Config.Qbit = config.QbitConfig{
		DownloadPath:      "/remote",
		LocalDownloadPath: sourceDir,
	}
	config.Config.Importers.ImportedTag = "imported"
	config.Config.Importers.ManualInterventionTag = "manual"

	// Mock qBittorrent client
	torrent := qbittorrent.Torrent{
		Hash:     "testformats",
		Name:     "Multiple Formats",
		SavePath: "/remote",
	}

	mockClient := &testutil.MockQbitClient{
		GetFilesInformationCtxReturn: struct {
			Files *qbittorrent.TorrentFiles
			Err   error
		}{
			Files: &qbittorrent.TorrentFiles{
				{Name: "book.epub"},
				{Name: "book.mobi"},
				{Name: "book.azw3"},
			},
		},
	}

	library := &config.ImportLibrary{
		Name: "test-library",
		Path: destDir,
	}

	importType := config.ImportType{
		Category: "books",
		Library:  "test-library",
	}

	// Run import
	importer := NewBookImporterSystem(mockClient)
	importer.ImportTorrent(ctx, torrent, importType, library)

	// Verify all three formats were imported
	_, err := os.Stat(filepath.Join(destDir, "book.epub"))
	assert.NoError(t, err, "epub should be imported")

	_, err = os.Stat(filepath.Join(destDir, "book.mobi"))
	assert.NoError(t, err, "mobi should be imported")

	_, err = os.Stat(filepath.Join(destDir, "book.azw3"))
	assert.NoError(t, err, "azw3 should be imported")

	// Verify torrent was tagged as imported
	assert.Len(t, mockClient.AddTagsCtxCalls, 1)
	assert.Equal(t, "imported", mockClient.AddTagsCtxCalls[0].Tags)
}

func TestImportTorrent_IgnoreNonBookFiles(t *testing.T) {
	ctx := context.Background()

	// Create temporary directories
	tempDir := t.TempDir()
	sourceDir := filepath.Join(tempDir, "source")
	destDir := filepath.Join(tempDir, "dest")

	// Set up source files - mix of book and non-book files
	require.NoError(t, os.MkdirAll(sourceDir, 0755))
	require.NoError(t, os.MkdirAll(destDir, 0755))

	epubFile := filepath.Join(sourceDir, "book.epub")
	txtFile := filepath.Join(sourceDir, "readme.txt")
	jpgFile := filepath.Join(sourceDir, "cover.jpg")

	require.NoError(t, os.WriteFile(epubFile, []byte("epub content"), 0644))
	require.NoError(t, os.WriteFile(txtFile, []byte("readme content"), 0644))
	require.NoError(t, os.WriteFile(jpgFile, []byte("image content"), 0644))

	// Set up config
	config.Config.Qbit = config.QbitConfig{
		DownloadPath:      "/remote",
		LocalDownloadPath: sourceDir,
	}
	config.Config.Importers.ImportedTag = "imported"
	config.Config.Importers.ManualInterventionTag = "manual"

	// Mock qBittorrent client
	torrent := qbittorrent.Torrent{
		Hash:     "testfilter",
		Name:     "Mixed Files",
		SavePath: "/remote",
	}

	mockClient := &testutil.MockQbitClient{
		GetFilesInformationCtxReturn: struct {
			Files *qbittorrent.TorrentFiles
			Err   error
		}{
			Files: &qbittorrent.TorrentFiles{
				{Name: "book.epub"},
				{Name: "readme.txt"},
				{Name: "cover.jpg"},
			},
		},
	}

	library := &config.ImportLibrary{
		Name: "test-library",
		Path: destDir,
	}

	importType := config.ImportType{
		Category: "books",
		Library:  "test-library",
	}

	// Run import
	importer := NewBookImporterSystem(mockClient)
	importer.ImportTorrent(ctx, torrent, importType, library)

	// Verify only epub was imported
	_, err := os.Stat(filepath.Join(destDir, "book.epub"))
	assert.NoError(t, err, "epub should be imported")

	_, err = os.Stat(filepath.Join(destDir, "readme.txt"))
	assert.True(t, os.IsNotExist(err), "txt file should not be imported")

	_, err = os.Stat(filepath.Join(destDir, "cover.jpg"))
	assert.True(t, os.IsNotExist(err), "jpg file should not be imported")

	// Verify torrent was tagged as imported
	assert.Len(t, mockClient.AddTagsCtxCalls, 1)
	assert.Equal(t, "imported", mockClient.AddTagsCtxCalls[0].Tags)
}

func TestImportTorrent_NoBookFilesFound(t *testing.T) {
	ctx := context.Background()

	// Create temporary directories
	tempDir := t.TempDir()
	sourceDir := filepath.Join(tempDir, "source")
	destDir := filepath.Join(tempDir, "dest")

	require.NoError(t, os.MkdirAll(sourceDir, 0755))
	require.NoError(t, os.MkdirAll(destDir, 0755))

	// Set up config
	config.Config.Qbit = config.QbitConfig{
		DownloadPath:      "/remote",
		LocalDownloadPath: sourceDir,
	}
	config.Config.Importers.ImportedTag = "imported"
	config.Config.Importers.ManualInterventionTag = "manual"

	// Mock qBittorrent client - only non-book files
	torrent := qbittorrent.Torrent{
		Hash:     "testnobooks",
		Name:     "No Books",
		SavePath: "/remote",
	}

	mockClient := &testutil.MockQbitClient{
		GetFilesInformationCtxReturn: struct {
			Files *qbittorrent.TorrentFiles
			Err   error
		}{
			Files: &qbittorrent.TorrentFiles{
				{Name: "readme.txt"},
				{Name: "cover.jpg"},
			},
		},
	}

	library := &config.ImportLibrary{
		Name: "test-library",
		Path: destDir,
	}

	importType := config.ImportType{
		Category: "books",
		Library:  "test-library",
	}

	// Run import
	importer := NewBookImporterSystem(mockClient)
	importer.ImportTorrent(ctx, torrent, importType, library)

	// Verify no files were imported
	entries, err := os.ReadDir(destDir)
	require.NoError(t, err)
	assert.Empty(t, entries, "No files should be imported")

	// Verify torrent was tagged for manual intervention
	assert.Len(t, mockClient.AddTagsCtxCalls, 1)
	assert.Equal(t, "manual", mockClient.AddTagsCtxCalls[0].Tags)
}

func TestImportTorrent_DeeplyNestedStructure(t *testing.T) {
	ctx := context.Background()

	// Create temporary directories
	tempDir := t.TempDir()
	sourceDir := filepath.Join(tempDir, "source")
	destDir := filepath.Join(tempDir, "dest")

	// Set up deeply nested source file
	nestedPath := filepath.Join(sourceDir, "Level1/Level2/Level3/Level4/Level5")
	require.NoError(t, os.MkdirAll(nestedPath, 0755))
	require.NoError(t, os.MkdirAll(destDir, 0755))

	sourceFile := filepath.Join(nestedPath, "deep-book.epub")
	require.NoError(t, os.WriteFile(sourceFile, []byte("deeply nested content"), 0644))

	// Set up config
	config.Config.Qbit = config.QbitConfig{
		DownloadPath:      "/remote",
		LocalDownloadPath: sourceDir,
	}
	config.Config.Importers.ImportedTag = "imported"
	config.Config.Importers.ManualInterventionTag = "manual"

	// Mock qBittorrent client
	torrent := qbittorrent.Torrent{
		Hash:     "testdeep",
		Name:     "Deeply Nested",
		SavePath: "/remote",
	}

	mockClient := &testutil.MockQbitClient{
		GetFilesInformationCtxReturn: struct {
			Files *qbittorrent.TorrentFiles
			Err   error
		}{
			Files: &qbittorrent.TorrentFiles{
				{Name: "Level1/Level2/Level3/Level4/Level5/deep-book.epub"},
			},
		},
	}

	library := &config.ImportLibrary{
		Name: "test-library",
		Path: destDir,
	}

	importType := config.ImportType{
		Category: "books",
		Library:  "test-library",
	}

	// Run import
	importer := NewBookImporterSystem(mockClient)
	importer.ImportTorrent(ctx, torrent, importType, library)

	// Verify file was flattened to destination root
	destFile := filepath.Join(destDir, "deep-book.epub")
	_, err := os.Stat(destFile)
	assert.NoError(t, err, "File should be flattened to destination root")

	// Verify no subdirectories were created
	entries, err := os.ReadDir(destDir)
	require.NoError(t, err)
	assert.Len(t, entries, 1, "Only one file should exist")
	assert.False(t, entries[0].IsDir(), "Entry should be a file, not a directory")

	// Verify content
	content, err := os.ReadFile(destFile)
	assert.NoError(t, err)
	assert.Equal(t, "deeply nested content", string(content))

	// Verify torrent was tagged as imported
	assert.Len(t, mockClient.AddTagsCtxCalls, 1)
	assert.Equal(t, "imported", mockClient.AddTagsCtxCalls[0].Tags)
}
