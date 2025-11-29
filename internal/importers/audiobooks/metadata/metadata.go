package metadata

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"strings"
	"text/template"
	"time"
)

// MetadataTags defines methods to retrieve metadata tags
type MetadataTags interface {
	Artist() (string, bool)
	Title() (string, bool)
	AudibleASIN() (string, bool)
}

// MetadataProvider defines method to get metadata for a given path
type MetadataProvider interface {
	GetMetadata(ctx context.Context, path string) (MetadataTags, error)
}

type Person struct {
	Name string  `json:"name"`
	Asin *string `json:"asin,omitempty"`
}

type Genre struct {
	Name string `json:"name"`
	Asin string `json:"asin"`
	Type string `json:"type"`
}

type Series struct {
	Name     string  `json:"name"`
	Asin     *string `json:"asin,omitempty"`
	Position *string `json:"position,omitempty"`
}

type BookMetadata struct {
	Asin            string    `json:"asin"`
	Authors         []Person  `json:"authors"`
	Copyright       *int      `json:"copyright,omitempty"`
	Description     string    `json:"description"`
	FormatType      string    `json:"formatType"`
	Genres          []Genre   `json:"genres"`
	Image           *string   `json:"image,omitempty"`
	IsAdult         *bool     `json:"isAdult,omitempty"`
	ISBN            *string   `json:"isbn,omitempty"`
	Language        string    `json:"language"`
	LiteratureType  *string   `json:"literatureType,omitempty"`
	Narrators       []Person  `json:"narrators"`
	PublisherName   string    `json:"publisherName"`
	Rating          string    `json:"rating"`
	Region          string    `json:"region"`
	ReleaseDate     time.Time `json:"releaseDate"`
	RuntimeLength   int       `json:"runtimeLengthMin"`
	PrimarySeries   *Series   `json:"seriesPrimary,omitempty"`
	SecondarySeries *Series   `json:"seriesSecondary,omitempty"`
	Subtitle        *string   `json:"subtitle,omitempty"`
	Summary         string    `json:"summary"`
	Title           string    `json:"title"`
}

const dirNameTemplate = `{{.Title}}{{if .PrimarySeries}} - {{.PrimarySeries.Name}}{{if .PrimarySeries.Position}} - Book {{.PrimarySeries.Position}}{{end}}{{end}}`

func (md *BookMetadata) GenerateDirectoryName() (string, error) {
	ctx := context.Background()

	slog.InfoContext(ctx, "Generating directory name from metadata",
		slog.String("title", md.Title),
		slog.String("asin", md.Asin))

	tmpl, err := template.New("opf").Parse(dirNameTemplate)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to parse directory name template",
			slog.String("title", md.Title), slog.Any("err", err))
		return "", err
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, &md)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to execute directory name template",
			slog.String("title", md.Title), slog.Any("err", err))
		return "", err
	}

	dirName := buf.String()
	slog.InfoContext(ctx, "Successfully generated directory name",
		slog.String("title", md.Title),
		slog.String("directoryName", dirName))

	return dirName, nil
}

func (md *BookMetadata) Summarize() string {
	var buffer strings.Builder

	truncatedTitle := md.Title

	if len(truncatedTitle) > 80 {
		truncatedTitle = truncatedTitle[:77] + "..."
	}

	if md.Authors == nil {
		return truncatedTitle
	}

	buffer.WriteString(truncatedTitle)
	buffer.WriteString(" by ")
	buffer.WriteString(md.Authors[0].Name)

	for i := 1; i < len(md.Authors); i++ {
		buffer.WriteString(" & ")
		buffer.WriteString(md.Authors[i].Name)
	}

	if md.PrimarySeries != nil {
		buffer.WriteString(" - ")
		buffer.WriteString(md.PrimarySeries.Name)

		if md.PrimarySeries.Position != nil {
			buffer.WriteString(" ")
			buffer.WriteString(*md.PrimarySeries.Position)
		}
	}

	buffer.WriteString(fmt.Sprintf(" (ASIN: %s)", md.Asin))

	return buffer.String()
}
