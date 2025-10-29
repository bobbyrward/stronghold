package metadata

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"gopkg.in/vansante/go-ffprobe.v2"
)

const mp3AudibleAsinPrefix = "http://www.audible.com/pd/"

// FFProbeTags implements MetadataTags using ffprobe tag list
type FFProbeTags struct {
	tagList ffprobe.Tags
}

// NewFFProbeTags creates a new FFProbeTags instance
func NewFFProbeTags(tagList ffprobe.Tags) MetadataTags {
	return &FFProbeTags{tagList: tagList}
}

// Artist retrieves the artist tag
func (fft *FFProbeTags) Artist() (string, bool) {
	value, err := fft.tagList.GetString("artist")
	if err != nil {
		return "", false
	}

	return value, true
}

// Title retrieves the title tag
func (fft *FFProbeTags) Title() (string, bool) {
	value, err := fft.tagList.GetString("title")
	if err != nil {
		return "", false
	}

	return value, true
}

// AudibleASIN retrieves the AUDIBLE_ASIN tag
func (fft *FFProbeTags) AudibleASIN() (string, bool) {
	value, err := fft.tagList.GetString("AUDIBLE_ASIN")
	if err != nil {
		return "", false
	}

	return strings.TrimPrefix(value, mp3AudibleAsinPrefix), true
}

// FFProbeMetadataProvider implements MetadataProvider using ffprobe
type FFProbeMetadataProvider struct{}

// NewFFProbeMetadataProvider creates a new FFProbeMetadataProvider instance
func NewFFProbeMetadataProvider() MetadataProvider {
	return &FFProbeMetadataProvider{}
}

// GetMetadata retrieves metadata tags for the given file path
func (ffmp *FFProbeMetadataProvider) GetMetadata(ctx context.Context, path string) (MetadataTags, error) {
	probeData, err := ffprobe.ProbeURL(ctx, path)
	if err != nil {
		return nil, errors.Join(
			fmt.Errorf("unable to probe metadata: %s", path),
			err,
		)
	}

	slog.InfoContext(ctx, "ffprobe format tags",
		slog.Any("tagList", probeData.Format.TagList))

	// Try format-level tags first
	if len(probeData.Format.TagList) > 0 {
		return NewFFProbeTags(probeData.Format.TagList), nil
	}

	// Fall back to first stream's tags if format tags are empty
	if len(probeData.Streams) > 0 && len(probeData.Streams[0].TagList) > 0 {
		slog.InfoContext(ctx, "Using stream-level tags as format tags were empty",
			slog.Any("streamTags", probeData.Streams[0].TagList))
		return NewFFProbeTags(probeData.Streams[0].TagList), nil
	}

	return NewFFProbeTags(ffprobe.Tags{}), nil
}
