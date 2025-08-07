package source

import (
	"context"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
)

type SourceInfo struct {
	Filename string
	IsDir    bool
	IsM4b    bool
	IsMp3    bool
}

func AnalyzeSource(filename string) (SourceInfo, error) {
	ctx := context.Background()
	
	slog.InfoContext(ctx, "Analyzing audiobook source", 
		slog.String("filename", filename))
	
	info := SourceInfo{
		Filename: filename,
	}

	fileInfo, err := os.Stat(filename)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to stat source file", 
			slog.String("filename", filename), slog.Any("err", err))
		return info, err
	}

	info.IsDir = fileInfo.IsDir()

	if info.IsDir {
		slog.InfoContext(ctx, "Source is directory, analyzing contents", 
			slog.String("filename", filename))
		
		fileCount := 0
		_ = filepath.WalkDir(info.Filename, func(path string, d fs.DirEntry, err error) error {
			if filepath.Ext(path) == ".mp3" {
				info.IsMp3 = true
				fileCount++
			} else if filepath.Ext(path) == ".m4b" {
				info.IsM4b = true
				fileCount++
			}
			return nil
		})
		
		slog.InfoContext(ctx, "Directory analysis complete", 
			slog.String("filename", filename),
			slog.Int("audioFiles", fileCount),
			slog.Bool("hasMp3", info.IsMp3),
			slog.Bool("hasM4b", info.IsM4b))
	} else {
		slog.InfoContext(ctx, "Source is single file", 
			slog.String("filename", filename),
			slog.Int64("size", fileInfo.Size()))
	}

	slog.InfoContext(ctx, "Successfully analyzed audiobook source", 
		slog.String("filename", filename),
		slog.Bool("isDir", info.IsDir),
		slog.Bool("hasMp3", info.IsMp3),
		slog.Bool("hasM4b", info.IsM4b))

	return info, nil
}
