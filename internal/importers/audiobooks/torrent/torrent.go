package torrent

import (
	"context"
	"errors"
	"log/slog"

	"github.com/autobrr/go-qbittorrent"
	"github.com/bobbyrward/stronghold/internal/importers/audiobooks/metadata"
	"github.com/bobbyrward/stronghold/internal/importers/audiobooks/source"
	"github.com/bobbyrward/stronghold/internal/importers/common"
	"github.com/bobbyrward/stronghold/internal/qbit"
	"github.com/cappuccinotm/slogx"
)

// TorrentMetadata represents metadata extracted from a torrent's filesj
type TorrentMetadata interface {
	// Hash returns the torrent's hash
	Hash() string

	// Tags returns the extracted metadata tags
	Tags() metadata.MetadataTags
}

// AudiobookFilesMetadata implements TorrentMetadata for audiobook torrents
type AudiobookFilesMetadata struct {
	torrent    qbittorrent.Torrent
	files      []common.MappedTorrentFile
	sourceInfo source.SourceInfo
	metadata   metadata.MetadataTags
}

// NewAudiobookFilesMetadata creates a new AudiobookFilesMetadata instance
func NewAudiobookFilesMetadata(ctx context.Context, qbit qbit.QbitClient, metadataProvider metadata.MetadataProvider, torrent qbittorrent.Torrent, mappedFiles []common.MappedTorrentFile) (*AudiobookFilesMetadata, error) {
	sourceInfo, err := source.AnalyzeSource(mappedFiles)
	if err != nil {
		msg := "unable to analyze torrent source files"
		slog.ErrorContext(ctx, msg,
			slog.String("name", torrent.Name),
			slog.String("hash", torrent.Hash),
			slogx.Error(err),
		)
		return nil, errors.Join(errors.New(msg), err)
	}

	metadata, err := getTagList(ctx, metadataProvider, sourceInfo)
	if err != nil {
		msg := "unable to extract tag list from source files"
		slog.ErrorContext(ctx, msg,
			slog.String("name", torrent.Name),
			slog.String("hash", torrent.Hash),
			slog.String("sourceType", sourceInfo.SourceType.String()),
			slogx.Error(err),
		)
		return nil, errors.Join(errors.New(msg), err)
	}

	abTorrent := &AudiobookFilesMetadata{
		torrent:    torrent,
		files:      mappedFiles,
		sourceInfo: sourceInfo,
		metadata:   metadata,
	}

	return abTorrent, nil
}

func (abfm *AudiobookFilesMetadata) Hash() string {
	return abfm.torrent.Hash
}

func (abfm *AudiobookFilesMetadata) FindASIN() (string, error) {
	return "", nil
}

func (abfm *AudiobookFilesMetadata) Tags() metadata.MetadataTags {
	return abfm.metadata
}

// getTagList extracts metadata tags based on the source type
func getTagList(ctx context.Context, metadataProvider metadata.MetadataProvider, sourceInfo source.SourceInfo) (metadata.MetadataTags, error) {
	switch sourceInfo.SourceType {
	case source.SourceType_M4B:
		fallthrough
	case source.SourceType_M4B_and_MP3:
		return metadataProvider.GetMetadata(ctx, sourceInfo.M4bFiles[0].LocalPath)

	case source.SourceType_MP3:
		return metadataProvider.GetMetadata(ctx, sourceInfo.Mp3Files[0].LocalPath)
	}

	return nil, errors.New("unable to determine source type for tag extraction")
}
