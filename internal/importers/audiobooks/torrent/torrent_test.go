package torrent

import (
	"context"
	"errors"
	"testing"

	"github.com/autobrr/go-qbittorrent"
	"github.com/stretchr/testify/assert"

	"github.com/bobbyrward/stronghold/internal/importers/audiobooks/metadata"
	"github.com/bobbyrward/stronghold/internal/importers/audiobooks/source"
	"github.com/bobbyrward/stronghold/internal/importers/common"
	"github.com/bobbyrward/stronghold/internal/testutil"
)

// MockMetadataProvider is a mock implementation of metadata.MetadataProvider for testing
type MockMetadataProvider struct {
	GetMetadataFunc   func(ctx context.Context, path string) (metadata.MetadataTags, error)
	GetMetadataCalls  []string
	GetMetadataReturn struct {
		Tags metadata.MetadataTags
		Err  error
	}
}

func (m *MockMetadataProvider) GetMetadata(ctx context.Context, path string) (metadata.MetadataTags, error) {
	m.GetMetadataCalls = append(m.GetMetadataCalls, path)

	if m.GetMetadataFunc != nil {
		return m.GetMetadataFunc(ctx, path)
	}

	return m.GetMetadataReturn.Tags, m.GetMetadataReturn.Err
}

// MockMetadataTags is a mock implementation of metadata.MetadataTags for testing
type MockMetadataTags struct {
	ArtistValue       string
	ArtistOk          bool
	TitleValue        string
	TitleOk           bool
	AudibleASINValue  string
	AudibleASINOk     bool
}

func (m *MockMetadataTags) Artist() (string, bool) {
	return m.ArtistValue, m.ArtistOk
}

func (m *MockMetadataTags) Title() (string, bool) {
	return m.TitleValue, m.TitleOk
}

func (m *MockMetadataTags) AudibleASIN() (string, bool) {
	return m.AudibleASINValue, m.AudibleASINOk
}

func TestNewAudiobookFilesMetadata_M4B_Success(t *testing.T) {
	ctx := context.Background()

	torrent := qbittorrent.Torrent{
		Hash: "abc123",
		Name: "Test Audiobook",
	}

	files := []common.MappedTorrentFile{
		{
			BaseName:  "audiobook.m4b",
			LocalPath: "/path/to/audiobook.m4b",
		},
	}

	expectedTags := &MockMetadataTags{
		ArtistValue:      "Test Author",
		ArtistOk:         true,
		TitleValue:       "Test Title",
		TitleOk:          true,
		AudibleASINValue: "B01234567",
		AudibleASINOk:    true,
	}

	mockMetadata := &MockMetadataProvider{
		GetMetadataReturn: struct {
			Tags metadata.MetadataTags
			Err  error
		}{
			Tags: expectedTags,
			Err:  nil,
		},
	}

	mockQbit := &testutil.MockQbitClient{}

	result, err := NewAudiobookFilesMetadata(ctx, mockQbit, mockMetadata, torrent, files)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "abc123", result.Hash())
	assert.Equal(t, expectedTags, result.Tags())

	// Verify metadata provider was called with correct path
	assert.Len(t, mockMetadata.GetMetadataCalls, 1)
	assert.Equal(t, "/path/to/audiobook.m4b", mockMetadata.GetMetadataCalls[0])
}

func TestNewAudiobookFilesMetadata_MP3_Success(t *testing.T) {
	ctx := context.Background()

	torrent := qbittorrent.Torrent{
		Hash: "def456",
		Name: "Test MP3 Audiobook",
	}

	files := []common.MappedTorrentFile{
		{
			BaseName:  "chapter01.mp3",
			LocalPath: "/path/to/chapter01.mp3",
		},
		{
			BaseName:  "chapter02.mp3",
			LocalPath: "/path/to/chapter02.mp3",
		},
	}

	expectedTags := &MockMetadataTags{
		ArtistValue: "Another Author",
		ArtistOk:    true,
		TitleValue:  "Another Title",
		TitleOk:     true,
	}

	mockMetadata := &MockMetadataProvider{
		GetMetadataReturn: struct {
			Tags metadata.MetadataTags
			Err  error
		}{
			Tags: expectedTags,
			Err:  nil,
		},
	}

	mockQbit := &testutil.MockQbitClient{}

	result, err := NewAudiobookFilesMetadata(ctx, mockQbit, mockMetadata, torrent, files)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "def456", result.Hash())
	assert.Equal(t, expectedTags, result.Tags())

	// Verify metadata provider was called with first MP3 file
	assert.Len(t, mockMetadata.GetMetadataCalls, 1)
	assert.Equal(t, "/path/to/chapter01.mp3", mockMetadata.GetMetadataCalls[0])
}

func TestNewAudiobookFilesMetadata_M4B_and_MP3_Success(t *testing.T) {
	ctx := context.Background()

	torrent := qbittorrent.Torrent{
		Hash: "ghi789",
		Name: "Mixed Format Audiobook",
	}

	files := []common.MappedTorrentFile{
		{
			BaseName:  "audiobook.m4b",
			LocalPath: "/path/to/audiobook.m4b",
		},
		{
			BaseName:  "chapter01.mp3",
			LocalPath: "/path/to/chapter01.mp3",
		},
		{
			BaseName:  "chapter02.mp3",
			LocalPath: "/path/to/chapter02.mp3",
		},
	}

	expectedTags := &MockMetadataTags{
		ArtistValue: "Mixed Author",
		ArtistOk:    true,
		TitleValue:  "Mixed Title",
		TitleOk:     true,
	}

	mockMetadata := &MockMetadataProvider{
		GetMetadataReturn: struct {
			Tags metadata.MetadataTags
			Err  error
		}{
			Tags: expectedTags,
			Err:  nil,
		},
	}

	mockQbit := &testutil.MockQbitClient{}

	result, err := NewAudiobookFilesMetadata(ctx, mockQbit, mockMetadata, torrent, files)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "ghi789", result.Hash())
	assert.Equal(t, expectedTags, result.Tags())

	// When both M4B and MP3 present, should prefer M4B file
	assert.Len(t, mockMetadata.GetMetadataCalls, 1)
	assert.Equal(t, "/path/to/audiobook.m4b", mockMetadata.GetMetadataCalls[0])
}

func TestNewAudiobookFilesMetadata_MetadataProviderError(t *testing.T) {
	ctx := context.Background()

	torrent := qbittorrent.Torrent{
		Hash: "error123",
		Name: "Error Audiobook",
	}

	files := []common.MappedTorrentFile{
		{
			BaseName:  "audiobook.m4b",
			LocalPath: "/path/to/audiobook.m4b",
		},
	}

	expectedError := errors.New("failed to read metadata")

	mockMetadata := &MockMetadataProvider{
		GetMetadataReturn: struct {
			Tags metadata.MetadataTags
			Err  error
		}{
			Tags: nil,
			Err:  expectedError,
		},
	}

	mockQbit := &testutil.MockQbitClient{}

	result, err := NewAudiobookFilesMetadata(ctx, mockQbit, mockMetadata, torrent, files)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "unable to extract tag list from source files")
	assert.ErrorIs(t, err, expectedError)

	// Verify metadata provider was called
	assert.Len(t, mockMetadata.GetMetadataCalls, 1)
}

func TestNewAudiobookFilesMetadata_UnknownSourceType(t *testing.T) {
	ctx := context.Background()

	torrent := qbittorrent.Torrent{
		Hash: "unknown123",
		Name: "Unknown Format",
	}

	// Files with no recognized audio formats
	files := []common.MappedTorrentFile{
		{
			BaseName:  "cover.jpg",
			LocalPath: "/path/to/cover.jpg",
		},
		{
			BaseName:  "info.txt",
			LocalPath: "/path/to/info.txt",
		},
	}

	mockMetadata := &MockMetadataProvider{}
	mockQbit := &testutil.MockQbitClient{}

	result, err := NewAudiobookFilesMetadata(ctx, mockQbit, mockMetadata, torrent, files)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "unable to extract tag list from source files")

	// Metadata provider should not be called for unknown source types
	assert.Len(t, mockMetadata.GetMetadataCalls, 0)
}

func TestNewAudiobookFilesMetadata_EmptyFiles(t *testing.T) {
	ctx := context.Background()

	torrent := qbittorrent.Torrent{
		Hash: "empty123",
		Name: "Empty Audiobook",
	}

	files := []common.MappedTorrentFile{}

	mockMetadata := &MockMetadataProvider{}
	mockQbit := &testutil.MockQbitClient{}

	result, err := NewAudiobookFilesMetadata(ctx, mockQbit, mockMetadata, torrent, files)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "unable to extract tag list from source files")
}

func TestNewAudiobookFilesMetadata_MultipleM4BFiles(t *testing.T) {
	ctx := context.Background()

	torrent := qbittorrent.Torrent{
		Hash: "multi123",
		Name: "Multi M4B Audiobook",
	}

	files := []common.MappedTorrentFile{
		{
			BaseName:  "part1.m4b",
			LocalPath: "/path/to/part1.m4b",
		},
		{
			BaseName:  "part2.m4b",
			LocalPath: "/path/to/part2.m4b",
		},
		{
			BaseName:  "part3.m4b",
			LocalPath: "/path/to/part3.m4b",
		},
	}

	expectedTags := &MockMetadataTags{
		ArtistValue: "Multi Author",
		ArtistOk:    true,
	}

	mockMetadata := &MockMetadataProvider{
		GetMetadataReturn: struct {
			Tags metadata.MetadataTags
			Err  error
		}{
			Tags: expectedTags,
			Err:  nil,
		},
	}

	mockQbit := &testutil.MockQbitClient{}

	result, err := NewAudiobookFilesMetadata(ctx, mockQbit, mockMetadata, torrent, files)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Should use first M4B file
	assert.Len(t, mockMetadata.GetMetadataCalls, 1)
	assert.Equal(t, "/path/to/part1.m4b", mockMetadata.GetMetadataCalls[0])
}

func TestAudiobookFilesMetadata_Hash(t *testing.T) {
	abfm := &AudiobookFilesMetadata{
		torrent: qbittorrent.Torrent{
			Hash: "test-hash-123",
			Name: "Test",
		},
	}

	assert.Equal(t, "test-hash-123", abfm.Hash())
}

func TestAudiobookFilesMetadata_Tags(t *testing.T) {
	expectedTags := &MockMetadataTags{
		ArtistValue: "Test Artist",
		ArtistOk:    true,
		TitleValue:  "Test Title",
		TitleOk:     true,
	}

	abfm := &AudiobookFilesMetadata{
		metadata: expectedTags,
	}

	tags := abfm.Tags()
	assert.Equal(t, expectedTags, tags)

	// Verify the tags work correctly
	artist, ok := tags.Artist()
	assert.True(t, ok)
	assert.Equal(t, "Test Artist", artist)

	title, ok := tags.Title()
	assert.True(t, ok)
	assert.Equal(t, "Test Title", title)
}

func TestAudiobookFilesMetadata_FindASIN(t *testing.T) {
	abfm := &AudiobookFilesMetadata{}

	asin, err := abfm.FindASIN()

	// Current implementation returns empty string and nil error
	assert.NoError(t, err)
	assert.Equal(t, "", asin)
}

func TestGetTagList_M4B(t *testing.T) {
	ctx := context.Background()

	files := []common.MappedTorrentFile{
		{
			BaseName:  "audiobook.m4b",
			LocalPath: "/path/to/audiobook.m4b",
		},
	}

	sourceInfo, err := analyzeSourceForTest(files)
	assert.NoError(t, err)

	expectedTags := &MockMetadataTags{
		ArtistValue: "M4B Artist",
		ArtistOk:    true,
	}

	mockMetadata := &MockMetadataProvider{
		GetMetadataReturn: struct {
			Tags metadata.MetadataTags
			Err  error
		}{
			Tags: expectedTags,
			Err:  nil,
		},
	}

	tags, err := getTagList(ctx, mockMetadata, sourceInfo)

	assert.NoError(t, err)
	assert.Equal(t, expectedTags, tags)
	assert.Len(t, mockMetadata.GetMetadataCalls, 1)
	assert.Equal(t, "/path/to/audiobook.m4b", mockMetadata.GetMetadataCalls[0])
}

func TestGetTagList_MP3(t *testing.T) {
	ctx := context.Background()

	files := []common.MappedTorrentFile{
		{
			BaseName:  "chapter01.mp3",
			LocalPath: "/path/to/chapter01.mp3",
		},
		{
			BaseName:  "chapter02.mp3",
			LocalPath: "/path/to/chapter02.mp3",
		},
	}

	sourceInfo, err := analyzeSourceForTest(files)
	assert.NoError(t, err)

	expectedTags := &MockMetadataTags{
		ArtistValue: "MP3 Artist",
		ArtistOk:    true,
	}

	mockMetadata := &MockMetadataProvider{
		GetMetadataReturn: struct {
			Tags metadata.MetadataTags
			Err  error
		}{
			Tags: expectedTags,
			Err:  nil,
		},
	}

	tags, err := getTagList(ctx, mockMetadata, sourceInfo)

	assert.NoError(t, err)
	assert.Equal(t, expectedTags, tags)
	assert.Len(t, mockMetadata.GetMetadataCalls, 1)
	assert.Equal(t, "/path/to/chapter01.mp3", mockMetadata.GetMetadataCalls[0])
}

func TestGetTagList_M4B_and_MP3(t *testing.T) {
	ctx := context.Background()

	files := []common.MappedTorrentFile{
		{
			BaseName:  "audiobook.m4b",
			LocalPath: "/path/to/audiobook.m4b",
		},
		{
			BaseName:  "chapter01.mp3",
			LocalPath: "/path/to/chapter01.mp3",
		},
	}

	sourceInfo, err := analyzeSourceForTest(files)
	assert.NoError(t, err)

	expectedTags := &MockMetadataTags{
		ArtistValue: "Mixed Artist",
		ArtistOk:    true,
	}

	mockMetadata := &MockMetadataProvider{
		GetMetadataReturn: struct {
			Tags metadata.MetadataTags
			Err  error
		}{
			Tags: expectedTags,
			Err:  nil,
		},
	}

	tags, err := getTagList(ctx, mockMetadata, sourceInfo)

	assert.NoError(t, err)
	assert.Equal(t, expectedTags, tags)
	// Should prefer M4B file
	assert.Len(t, mockMetadata.GetMetadataCalls, 1)
	assert.Equal(t, "/path/to/audiobook.m4b", mockMetadata.GetMetadataCalls[0])
}

func TestGetTagList_UnknownSourceType(t *testing.T) {
	ctx := context.Background()

	files := []common.MappedTorrentFile{
		{
			BaseName:  "cover.jpg",
			LocalPath: "/path/to/cover.jpg",
		},
	}

	sourceInfo, err := analyzeSourceForTest(files)
	assert.NoError(t, err)

	mockMetadata := &MockMetadataProvider{}

	tags, err := getTagList(ctx, mockMetadata, sourceInfo)

	assert.Error(t, err)
	assert.Nil(t, tags)
	assert.Contains(t, err.Error(), "unable to determine source type for tag extraction")
	// Should not call metadata provider
	assert.Len(t, mockMetadata.GetMetadataCalls, 0)
}

func TestGetTagList_MetadataError(t *testing.T) {
	ctx := context.Background()

	files := []common.MappedTorrentFile{
		{
			BaseName:  "audiobook.m4b",
			LocalPath: "/path/to/audiobook.m4b",
		},
	}

	sourceInfo, err := analyzeSourceForTest(files)
	assert.NoError(t, err)

	expectedError := errors.New("metadata extraction failed")

	mockMetadata := &MockMetadataProvider{
		GetMetadataReturn: struct {
			Tags metadata.MetadataTags
			Err  error
		}{
			Tags: nil,
			Err:  expectedError,
		},
	}

	tags, err := getTagList(ctx, mockMetadata, sourceInfo)

	assert.Error(t, err)
	assert.Nil(t, tags)
	assert.ErrorIs(t, err, expectedError)
	assert.Len(t, mockMetadata.GetMetadataCalls, 1)
}

// Helper function to analyze source for tests
func analyzeSourceForTest(files []common.MappedTorrentFile) (source.SourceInfo, error) {
	// Import the actual source package function
	return source.AnalyzeSource(files)
}

func TestNewAudiobookFilesMetadata_WithMixedNonAudioFiles(t *testing.T) {
	ctx := context.Background()

	torrent := qbittorrent.Torrent{
		Hash: "mixed123",
		Name: "Mixed Files",
	}

	files := []common.MappedTorrentFile{
		{
			BaseName:  "audiobook.m4b",
			LocalPath: "/path/to/audiobook.m4b",
		},
		{
			BaseName:  "cover.jpg",
			LocalPath: "/path/to/cover.jpg",
		},
		{
			BaseName:  "info.nfo",
			LocalPath: "/path/to/info.nfo",
		},
		{
			BaseName:  "metadata.xml",
			LocalPath: "/path/to/metadata.xml",
		},
	}

	expectedTags := &MockMetadataTags{
		ArtistValue: "Test Artist",
		ArtistOk:    true,
	}

	mockMetadata := &MockMetadataProvider{
		GetMetadataReturn: struct {
			Tags metadata.MetadataTags
			Err  error
		}{
			Tags: expectedTags,
			Err:  nil,
		},
	}

	mockQbit := &testutil.MockQbitClient{}

	result, err := NewAudiobookFilesMetadata(ctx, mockQbit, mockMetadata, torrent, files)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	// Should correctly identify and use the M4B file
	assert.Len(t, mockMetadata.GetMetadataCalls, 1)
	assert.Equal(t, "/path/to/audiobook.m4b", mockMetadata.GetMetadataCalls[0])
}

func TestNewAudiobookFilesMetadata_CustomMetadataFunc(t *testing.T) {
	// This test demonstrates using a custom function for more complex metadata scenarios
	ctx := context.Background()

	torrent := qbittorrent.Torrent{
		Hash: "custom456",
		Name: "Custom Metadata Test",
	}

	files := []common.MappedTorrentFile{
		{
			BaseName:  "book.m4b",
			LocalPath: "/path/to/book.m4b",
		},
	}

	callCount := 0
	mockMetadata := &MockMetadataProvider{}
	mockMetadata.GetMetadataFunc = func(ctx context.Context, path string) (metadata.MetadataTags, error) {
		callCount++
		// Custom logic - could do different things based on path
		if path == "/path/to/book.m4b" {
			return &MockMetadataTags{
				ArtistValue:      "Custom Artist",
				ArtistOk:         true,
				TitleValue:       "Custom Title",
				TitleOk:          true,
				AudibleASINValue: "B09876543",
				AudibleASINOk:    true,
			}, nil
		}
		return nil, errors.New("unexpected path")
	}

	mockQbit := &testutil.MockQbitClient{}

	result, err := NewAudiobookFilesMetadata(ctx, mockQbit, mockMetadata, torrent, files)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, callCount)

	// Verify the custom tags
	tags := result.Tags()
	artist, ok := tags.Artist()
	assert.True(t, ok)
	assert.Equal(t, "Custom Artist", artist)

	title, ok := tags.Title()
	assert.True(t, ok)
	assert.Equal(t, "Custom Title", title)

	asin, ok := tags.AudibleASIN()
	assert.True(t, ok)
	assert.Equal(t, "B09876543", asin)
}
