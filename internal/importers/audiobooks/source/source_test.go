package source

import (
	"testing"

	"github.com/bobbyrward/stronghold/internal/importers/common"
	"github.com/stretchr/testify/assert"
)

func TestSourceType_String(t *testing.T) {
	tests := []struct {
		name     string
		source   SourceType
		expected string
	}{
		{
			name:     "M4B source type",
			source:   SourceType_M4B,
			expected: "M4B",
		},
		{
			name:     "MP3 source type",
			source:   SourceType_MP3,
			expected: "MP3",
		},
		{
			name:     "M4B and MP3 source type",
			source:   SourceType_M4B_and_MP3,
			expected: "M4B and MP3",
		},
		{
			name:     "Unknown source type",
			source:   SourceType_Unknown,
			expected: "Unknown",
		},
		{
			name:     "Invalid source type",
			source:   SourceType(999),
			expected: "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.source.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAnalyzeSource_M4BOnly(t *testing.T) {
	files := []common.MappedTorrentFile{
		{
			BaseName:  "audiobook1.m4b",
			LocalPath: "/path/to/audiobook1.m4b",
		},
		{
			BaseName:  "audiobook2.m4b",
			LocalPath: "/path/to/audiobook2.m4b",
		},
	}

	result, err := AnalyzeSource(files)

	assert.NoError(t, err)
	assert.Equal(t, SourceType_M4B, result.SourceType)
	assert.Len(t, result.M4bFiles, 2)
	assert.Len(t, result.Mp3Files, 0)
	assert.Equal(t, "audiobook1.m4b", result.M4bFiles[0].BaseName)
	assert.Equal(t, "audiobook2.m4b", result.M4bFiles[1].BaseName)
}

func TestAnalyzeSource_MP3Only(t *testing.T) {
	files := []common.MappedTorrentFile{
		{
			BaseName:  "chapter01.mp3",
			LocalPath: "/path/to/chapter01.mp3",
		},
		{
			BaseName:  "chapter02.mp3",
			LocalPath: "/path/to/chapter02.mp3",
		},
		{
			BaseName:  "chapter03.mp3",
			LocalPath: "/path/to/chapter03.mp3",
		},
	}

	result, err := AnalyzeSource(files)

	assert.NoError(t, err)
	assert.Equal(t, SourceType_MP3, result.SourceType)
	assert.Len(t, result.Mp3Files, 3)
	assert.Len(t, result.M4bFiles, 0)
	assert.Equal(t, "chapter01.mp3", result.Mp3Files[0].BaseName)
	assert.Equal(t, "chapter02.mp3", result.Mp3Files[1].BaseName)
	assert.Equal(t, "chapter03.mp3", result.Mp3Files[2].BaseName)
}

func TestAnalyzeSource_BothM4BandMP3(t *testing.T) {
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

	result, err := AnalyzeSource(files)

	assert.NoError(t, err)
	assert.Equal(t, SourceType_M4B_and_MP3, result.SourceType)
	assert.Len(t, result.M4bFiles, 1)
	assert.Len(t, result.Mp3Files, 2)
	assert.Equal(t, "audiobook.m4b", result.M4bFiles[0].BaseName)
	assert.Equal(t, "chapter01.mp3", result.Mp3Files[0].BaseName)
}

func TestAnalyzeSource_UnknownOnly(t *testing.T) {
	files := []common.MappedTorrentFile{
		{
			BaseName:  "cover.jpg",
			LocalPath: "/path/to/cover.jpg",
		},
		{
			BaseName:  "info.txt",
			LocalPath: "/path/to/info.txt",
		},
		{
			BaseName:  "metadata.json",
			LocalPath: "/path/to/metadata.json",
		},
	}

	result, err := AnalyzeSource(files)

	assert.NoError(t, err)
	assert.Equal(t, SourceType_Unknown, result.SourceType)
	assert.Len(t, result.M4bFiles, 0)
	assert.Len(t, result.Mp3Files, 0)
}

func TestAnalyzeSource_EmptyFileList(t *testing.T) {
	files := []common.MappedTorrentFile{}

	result, err := AnalyzeSource(files)

	assert.NoError(t, err)
	assert.Equal(t, SourceType_Unknown, result.SourceType)
	assert.Len(t, result.M4bFiles, 0)
	assert.Len(t, result.Mp3Files, 0)
}

func TestAnalyzeSource_MixedWithNonAudioFiles(t *testing.T) {
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
			BaseName:  "info.txt",
			LocalPath: "/path/to/info.txt",
		},
		{
			BaseName:  "chapter01.mp3",
			LocalPath: "/path/to/chapter01.mp3",
		},
		{
			BaseName:  "README.md",
			LocalPath: "/path/to/README.md",
		},
	}

	result, err := AnalyzeSource(files)

	assert.NoError(t, err)
	assert.Equal(t, SourceType_M4B_and_MP3, result.SourceType)
	assert.Len(t, result.M4bFiles, 1)
	assert.Len(t, result.Mp3Files, 1)
	assert.Equal(t, "audiobook.m4b", result.M4bFiles[0].BaseName)
	assert.Equal(t, "chapter01.mp3", result.Mp3Files[0].BaseName)
}

func TestAnalyzeSource_CaseSensitiveExtensions(t *testing.T) {
	// Test that file extensions are case-sensitive
	files := []common.MappedTorrentFile{
		{
			BaseName:  "audiobook.M4B", // uppercase
			LocalPath: "/path/to/audiobook.M4B",
		},
		{
			BaseName:  "chapter01.MP3", // uppercase
			LocalPath: "/path/to/chapter01.MP3",
		},
		{
			BaseName:  "chapter02.Mp3", // mixed case
			LocalPath: "/path/to/chapter02.Mp3",
		},
	}

	result, err := AnalyzeSource(files)

	assert.NoError(t, err)
	// Since filepath.Ext is case-sensitive on most systems,
	// uppercase extensions should not match
	assert.Equal(t, SourceType_Unknown, result.SourceType)
	assert.Len(t, result.M4bFiles, 0)
	assert.Len(t, result.Mp3Files, 0)
}

func TestAnalyzeSource_NoExtension(t *testing.T) {
	files := []common.MappedTorrentFile{
		{
			BaseName:  "audiobook",
			LocalPath: "/path/to/audiobook",
		},
		{
			BaseName:  "chapter01",
			LocalPath: "/path/to/chapter01",
		},
	}

	result, err := AnalyzeSource(files)

	assert.NoError(t, err)
	assert.Equal(t, SourceType_Unknown, result.SourceType)
	assert.Len(t, result.M4bFiles, 0)
	assert.Len(t, result.Mp3Files, 0)
}

func TestAnalyzeSource_SingleM4BFile(t *testing.T) {
	files := []common.MappedTorrentFile{
		{
			BaseName:  "complete-audiobook.m4b",
			LocalPath: "/downloads/complete-audiobook.m4b",
		},
	}

	result, err := AnalyzeSource(files)

	assert.NoError(t, err)
	assert.Equal(t, SourceType_M4B, result.SourceType)
	assert.Len(t, result.M4bFiles, 1)
	assert.Len(t, result.Mp3Files, 0)
	assert.Equal(t, "complete-audiobook.m4b", result.M4bFiles[0].BaseName)
	assert.Equal(t, "/downloads/complete-audiobook.m4b", result.M4bFiles[0].LocalPath)
}

func TestAnalyzeSource_SingleMP3File(t *testing.T) {
	files := []common.MappedTorrentFile{
		{
			BaseName:  "complete-audiobook.mp3",
			LocalPath: "/downloads/complete-audiobook.mp3",
		},
	}

	result, err := AnalyzeSource(files)

	assert.NoError(t, err)
	assert.Equal(t, SourceType_MP3, result.SourceType)
	assert.Len(t, result.Mp3Files, 1)
	assert.Len(t, result.M4bFiles, 0)
	assert.Equal(t, "complete-audiobook.mp3", result.Mp3Files[0].BaseName)
	assert.Equal(t, "/downloads/complete-audiobook.mp3", result.Mp3Files[0].LocalPath)
}

func TestAnalyzeSource_MultipleM4BFiles(t *testing.T) {
	files := []common.MappedTorrentFile{
		{
			BaseName:  "book-part1.m4b",
			LocalPath: "/path/book-part1.m4b",
		},
		{
			BaseName:  "book-part2.m4b",
			LocalPath: "/path/book-part2.m4b",
		},
		{
			BaseName:  "book-part3.m4b",
			LocalPath: "/path/book-part3.m4b",
		},
		{
			BaseName:  "book-part4.m4b",
			LocalPath: "/path/book-part4.m4b",
		},
	}

	result, err := AnalyzeSource(files)

	assert.NoError(t, err)
	assert.Equal(t, SourceType_M4B, result.SourceType)
	assert.Len(t, result.M4bFiles, 4)
	assert.Len(t, result.Mp3Files, 0)

	for i, file := range result.M4bFiles {
		assert.Contains(t, file.BaseName, "book-part")
		assert.Contains(t, file.BaseName, ".m4b")
		assert.Equal(t, files[i].BaseName, file.BaseName)
	}
}

func TestAnalyzeSource_LargeMP3Collection(t *testing.T) {
	files := make([]common.MappedTorrentFile, 50)
	for i := 0; i < 50; i++ {
		files[i] = common.MappedTorrentFile{
			BaseName:  "chapter" + string(rune('0'+i)) + ".mp3",
			LocalPath: "/path/chapter" + string(rune('0'+i)) + ".mp3",
		}
	}

	result, err := AnalyzeSource(files)

	assert.NoError(t, err)
	assert.Equal(t, SourceType_MP3, result.SourceType)
	assert.Len(t, result.Mp3Files, 50)
	assert.Len(t, result.M4bFiles, 0)
}

func TestAnalyzeSource_OtherAudioFormats(t *testing.T) {
	// Test that other audio formats are not recognized
	files := []common.MappedTorrentFile{
		{
			BaseName:  "audiobook.flac",
			LocalPath: "/path/to/audiobook.flac",
		},
		{
			BaseName:  "audiobook.ogg",
			LocalPath: "/path/to/audiobook.ogg",
		},
		{
			BaseName:  "audiobook.wav",
			LocalPath: "/path/to/audiobook.wav",
		},
		{
			BaseName:  "audiobook.aac",
			LocalPath: "/path/to/audiobook.aac",
		},
	}

	result, err := AnalyzeSource(files)

	assert.NoError(t, err)
	assert.Equal(t, SourceType_Unknown, result.SourceType)
	assert.Len(t, result.M4bFiles, 0)
	assert.Len(t, result.Mp3Files, 0)
}
