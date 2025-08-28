package booksearch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/bobbyrward/stronghold/internal/config"
	"github.com/cappuccinotm/slogx/logger"
	"github.com/carlmjohnson/requests"
)

type BookSearchService struct{}

func NewBookSearchService() *BookSearchService {
	return &BookSearchService{}
}

func (params *SearchParameters) Validate() error {
	if (params.Hash != "") == (params.Query != "") {
		return fmt.Errorf("one of Query or Hash is required")
	}

	return nil
}

/*
func logRequest(req *http.Request, res *http.Response, err error, d time.Duration) {
	slog.Info(
		"REQUEST",
		slog.String("url", req.URL.String()),
		slog.String("method", req.Method),
		slog.String("status", res.Status),
		slog.Int("status_code", res.StatusCode),
		slogx.Error(err),
	)
}
*/

func createHttpClient() *http.Client {
	searchConfig := &config.Config.BookSearch

	transport := &http.Transport{}

	l := logger.New(
		logger.WithLogger(slog.Default()),
		logger.WithBody(10240),
	)

	if searchConfig.HttpsProxy != "" {
		transport.Proxy = http.ProxyURL(&url.URL{
			Scheme: "http",
			Host:   searchConfig.HttpsProxy,
		})
	}

	client := &http.Client{
		Timeout:   15 * time.Second,
		Transport: l.HTTPClientRoundTripper(transport),
	}

	return client
}

func (s *BookSearchService) Search(ctx context.Context, params *SearchParameters) (*SearchResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params is required")
	}

	err := params.Validate()
	if err != nil {
		return nil, err
	}

	searchConfig := &config.Config.BookSearch
	baseUrl := searchConfig.BaseURL
	searchEndpoint := searchConfig.SearchEndpoint

	url := fmt.Sprintf("%s%s", baseUrl, searchEndpoint)

	request := SearchRequest{
		IncludeDlLink:      true,
		IncludeDescription: true,
		IncludeISBN:        true,
		PerPage:            params.MaxResults,
		Tor: SearchRequestTor{
			BrowseLang:     []int{1},
			MainCategories: []int{13, 14},
			Categories:     []int{41, 47, 108, 63, 69, 0},
			SearchType:     "active",
			SortType:       "dateDesc",
			SearchIn: SearchRequestSrchIn{
				Author: true,
				Series: true,
				Title:  true,
			},
			Query: params.Query,
			Hash:  params.Hash,
		},
	}

	bytes, err := json.Marshal(request)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to marshal search request", slog.Any("err", err))
		return nil, err
	}
	slog.InfoContext(ctx, "Marshalled request", slog.String("request", string(bytes)))

	var response SearchResponse

	slog.InfoContext(ctx, "Making request", slog.String("url", url), slog.Any("request", request))

	err = requests.
		URL(url).
		Client(createHttpClient()).
		Cookie("mam_id", config.Config.BookSearch.APIKey).
		Method("POST").
		CheckStatus(200).
		BodyJSON(&request).
		ToJSON(&response).
		Fetch(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to search", slog.Any("err", err))
		return nil, errors.Join(err, fmt.Errorf("unable to search"))
	}

	slog.InfoContext(ctx, "made request", slog.Any("response", response))

	return &response, nil
}

func (s *BookSearchService) DisplayResults(searchParams *SearchParameters, result *SearchResponse, format string) error {
	switch format {
	case "json":
		return s.displayJSON(searchParams, result)
	case "table":
		return s.displayTable(searchParams, result)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func (s *BookSearchService) displayJSON(searchParams *SearchParameters, result *SearchResponse) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

func (s *BookSearchService) displayTable(searchParams *SearchParameters, result *SearchResponse) error {
	if len(result.Data) == 0 {
		fmt.Printf("No books found for query: %s\n", searchParams.Query)
		return nil
	}

	if searchParams.Query != "" {
		fmt.Printf("Found %d books for query: %s\n\n", result.TotalFound, searchParams.Query)
	} else {
		fmt.Printf("Found %d books for hash: %s\n\n", result.TotalFound, searchParams.Hash)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(w, "TITLE\tCATEGORY\tAUTHOR(S)\tSERIES(S)\tNARRATOR(S)\tFILETYPE(S)")
	_, _ = fmt.Fprintln(w, "-----\t--------\t---------\t----\t---------\t---------")

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

		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			title,
			categoryName,
			authors,
			series,
			narrators,
			fileTypes,
		)
	}

	return w.Flush()
}
