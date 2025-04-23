package metadata

import (
	"bytes"
	"text/template"
	"time"
)

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
	tmpl, err := template.New("opf").Parse(dirNameTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, &md)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
