package booksearch

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/charmbracelet/glamour"
	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/bobbyrward/stronghold/internal/booksearch"
	"github.com/bobbyrward/stronghold/internal/config"
)

func displaySearchResults(searchParams *booksearch.SearchParameters, result *booksearch.SearchResponse, format string) error {
	switch format {
	case "json":
		return displaySearchResultsJSON(searchParams, result)
	case "table":
		return displaySearchResultsTable(searchParams, result)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func displaySearchResultsJSON(searchParams *booksearch.SearchParameters, result *booksearch.SearchResponse) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

func displaySearchResultsTable(searchParams *booksearch.SearchParameters, result *booksearch.SearchResponse) error {
	if len(result.Data) == 0 {
		fmt.Printf("No books found for query: %s\n", searchParams.Query)
		return nil
	}

	header := table.Row{"Title", "Category", "Author(s)", "Series(s)", "Narrator(s)", "Filetype(s)"}

	tableWriter := table.NewWriter()
	tableWriter.AppendHeader(header)
	tableWriter.AppendFooter(table.Row{"", "", "", "", "Total", fmt.Sprintf("%d", len(result.Data))})

	for _, book := range result.Data {
		authors := strings.Join(slices.Collect(maps.Values(book.Authors)), ", ")
		if len(authors) > 30 {
			authors = authors[:27] + "..."
		}

		series := ""

		if len(book.Series) > 0 {
			seriesCollected := make([]string, 0)

			for _, value := range book.Series {
				seriesCollected = append(seriesCollected, fmt.Sprintf("%s(%d)", value.Name, value.Index))
			}

			series = strings.Join(seriesCollected, ", ")
			if len(series) > 30 {
				series = series[:27] + "..."
			}
		}

		narrators := ""

		if len(book.Narrators) > 0 {
			narrators = strings.Join(slices.Collect(maps.Values(book.Narrators)), ", ")
			if len(narrators) > 30 {
				narrators = narrators[:27] + "..."
			}
		}

		title := book.Title
		if len(title) > 40 {
			title = title[:37] + "..."
		}

		fileTypes := book.FileTypes
		if len(fileTypes) > 40 {
			fileTypes = fileTypes[:37] + "..."
		}

		categoryName := book.CategoryName
		if len(categoryName) > 40 {
			categoryName = categoryName[:37] + "..."
		}

		tableWriter.AppendRow(table.Row{title, categoryName, authors, series, narrators, fileTypes})
	}

	fmt.Println(tableWriter.Render())

	return nil
}

func displaySingleResult(searchParams *booksearch.SearchParameters, result *booksearch.SearchResponse) error {
	if len(result.Data) != 1 {
		fmt.Printf("No books found for query: %s\n", searchParams.Query)
		return nil
	}

	header := table.Row{"Field", "Value"}

	tableWriter := table.NewWriter()
	tableWriter.AppendHeader(header)

	book := result.Data[0]

	tableWriter.AppendRow(table.Row{"Title", book.Title})
	tableWriter.AppendRow(table.Row{"Category", book.CategoryName})
	tableWriter.AppendRow(table.Row{"Author(s)", strings.Join(slices.Collect(maps.Values(book.Authors)), ", ")})

	if len(book.Series) > 0 {
		series := ""

		seriesCollected := make([]string, 0)

		for _, value := range book.Series {
			seriesCollected = append(seriesCollected, fmt.Sprintf("%s(%d)", value.Name, value.Index))
		}

		series = strings.Join(seriesCollected, ", ")
		tableWriter.AppendRow(table.Row{"Series(s)", series})
	}

	if len(book.Narrators) > 0 {
		tableWriter.AppendRow(table.Row{"Narrator(s)", strings.Join(slices.Collect(maps.Values(book.Narrators)), ", ")})
	}

	tableWriter.AppendRow(table.Row{"Added", book.Added})
	tableWriter.AppendRow(table.Row{"FileTypes", book.FileTypes})
	tableWriter.AppendRow(table.Row{"ISBN", book.ISBN})
	tableWriter.AppendRow(table.Row{"Tags", book.Tags})

	markdown, err := htmltomarkdown.ConvertString(book.Description)
	if err != nil {
		markdown = book.Description
	}

	rendered, err := glamour.Render(markdown, "dark")
	if err != nil {
		rendered = book.Description
	}

	tableWriter.AppendRow(table.Row{"Description", rendered})

	link := fmt.Sprintf("%s%d", config.Config.BookSearch.TorrentUrlPrefix, book.ID)
	tableWriter.AppendRow(table.Row{"Link", fmt.Sprintf("\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\", link, link)})

	fmt.Println(tableWriter.Render())

	return nil
}
