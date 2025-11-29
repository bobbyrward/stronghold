package metadata

import (
	"html/template"
	"os"
)

//  title, author, narrator, publishYear, publisher, isbn, description, genres, language, series, volumeNumber

const opfTemplate = `<?xml version='1.0' encoding='utf-8'?>
<ns0:package xmlns:dc='http://purl.org/dc/elements/1.1/' xmlns:ns0='http://www.idpf.org/2007/opf' unique-identifier='BookId' version='2.0'>
  <ns0:metadata xmlns:dc="http://purl.org/dc/elements/1.1/"
          xmlns:opf="http://www.idpf.org/2007/opf">

<!-- Open Packaging Format (opf) File layout for audiobookshelf -->
<!-- author, narrator, genre series, {series/volume number}, and tag can be repeated as desired -->

	<dc:title>{{ .Title }}</dc:title>

	{{ if .Subtitle }}
	<dc:subtitle>{{ .Subtitle }}</dc:subtitle>
	{{ end }}

	<dc:description>{{ .Summary }}</dc:description>

	{{ range .Authors }}
	<dc:creator opf:role="aut">{{ .Name }}</dc:creator>      <!-- author -->
	{{ end }}

	{{ range .Narrators }}
	<dc:creator opf:role="nrt">{{ .Name }}</dc:creator>      <!-- narrator -->
	{{ end }}

	<dc:publisher>{{ .PublisherName }}</dc:publisher>

	<!-- <dc:date></dc:date> -->                           <!-- publish year --> 

	<dc:language>eng</dc:language>
	{{ range .Genres }}
	<dc:subject>{{ .Name }}</dc:subject>                     <!-- genre -->
	{{ end }}

	{{ if .ISBN }}
	<dc:identifier opf:scheme="ISBN">{{ .ISBN }}</dc:identifier>
	{{ end }}
	<dc:identifier opf:scheme="ASIN">{{ .Asin }}</dc:identifier>

	{{ if .PrimarySeries }}
	<ns0:meta name="calibre:series" content="{{ .PrimarySeries.Name }}" /> <!-- series -->
	{{ if .PrimarySeries.Position }}
	<ns0:meta name="calibre:series_index" content="{{ .PrimarySeries.Position }}" /> <!-- volumeNumber -->
	{{ end}}
	{{ end }}

	<dc:tag></dc:tag>
  </ns0:metadata>
</ns0:package>
`

func (md BookMetadata) WriteOpf(filename string) error {
	tmpl, err := template.New("opf").Parse(opfTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	err = tmpl.Execute(file, &md)
	if err != nil {
		_ = file.Close()
		return err
	}

	return file.Close()
}
