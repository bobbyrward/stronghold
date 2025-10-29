package source

import (
	"context"
	"log/slog"
	"path/filepath"

	"github.com/bobbyrward/stronghold/internal/importers/common"
)

type SourceType int

const (
	SourceType_M4B SourceType = iota
	SourceType_MP3
	SourceType_Unknown
	SourceType_M4B_and_MP3
)

func (st SourceType) String() string {
	switch st {
	case SourceType_M4B:
		return "M4B"
	case SourceType_MP3:
		return "MP3"
	case SourceType_M4B_and_MP3:
		return "M4B and MP3"
	case SourceType_Unknown:
		return "Unknown"
	default:
		return "Unknown"
	}
}

type SourceInfo struct {
	M4bFiles   []common.MappedTorrentFile
	Mp3Files   []common.MappedTorrentFile
	SourceType SourceType
}

func AnalyzeSource(files []common.MappedTorrentFile) (SourceInfo, error) {
	ctx := context.Background()

	info := SourceInfo{}

	for _, file := range files {
		ext := filepath.Ext(file.BaseName)

		switch ext {
		case ".mp3":
			info.Mp3Files = append(info.Mp3Files, file)
		case ".m4b":
			info.M4bFiles = append(info.M4bFiles, file)
		default:
		}
	}

	hasMp3Files := len(info.Mp3Files) > 0
	hasM4bFiles := len(info.M4bFiles) > 0

	if !hasMp3Files && !hasM4bFiles {
		info.SourceType = SourceType_Unknown
	} else if hasMp3Files && hasM4bFiles {
		info.SourceType = SourceType_M4B_and_MP3
	} else if hasMp3Files && !hasM4bFiles {
		info.SourceType = SourceType_MP3
	} else if !hasMp3Files && hasM4bFiles {
		info.SourceType = SourceType_M4B
	}

	slog.InfoContext(ctx, "Source analysis complete",
		slog.Any("Mp3Files", info.Mp3Files),
		slog.Any("M4bFiles", info.M4bFiles),
		slog.String("SourceType", info.SourceType.String()),
	)

	return info, nil
}
