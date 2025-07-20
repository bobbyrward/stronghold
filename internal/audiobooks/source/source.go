package source

import (
	"io/fs"
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
	info := SourceInfo{
		Filename: filename,
	}

	fileInfo, err := os.Stat(filename)
	if err != nil {
		return info, err
	}

	info.IsDir = fileInfo.IsDir()

	if info.IsDir {
		_ = filepath.WalkDir(info.Filename, func(path string, d fs.DirEntry, err error) error {
			if filepath.Ext(path) == "mp3" {
				info.IsMp3 = true
			} else if filepath.Ext(path) == "m4b" {
				info.IsM4b = true
			}
			return nil
		})
	}

	return info, nil
}
