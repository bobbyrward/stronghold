package audible

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/bobbyrward/stronghold/internal/audiobooks/metadata"
)

const (
	asinSearchUrl   = "https://api.audible.com/1.0/catalog/products"
	asinMetadataUrl = "https://api.audnex.us/books/%s"
)

type asinItem struct {
	Asin string `json:"asin"`
}

type asinSearchResponse struct {
	Products     []asinItem `json:"products"`
	TotalResults int        `json:"total_results"`
}

type AudibleApiClient struct {
	httpClient *http.Client
}

func NewAudibleApiClient() *AudibleApiClient {
	ctx := context.Background()
	
	slog.InfoContext(ctx, "Creating Audible API client")
	
	aac := &AudibleApiClient{
		httpClient: &http.Client{},
	}

	return aac
}

func (aac *AudibleApiClient) SearchByTitle(title string) ([]string, error) {
	ctx := context.Background()
	
	slog.InfoContext(ctx, "Searching Audible by title", 
		slog.String("title", title),
		slog.String("url", asinSearchUrl))
	
	request, err := http.NewRequest("GET", asinSearchUrl, nil)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create Audible search request", 
			slog.String("title", title), slog.Any("err", err))
		return nil, err
	}

	query := request.URL.Query()
	query.Add("num_results", "10")
	query.Add("products_sort_by", "Relevance")
	query.Add("title", title)

	request.URL.RawQuery = query.Encode()

	response, err := aac.httpClient.Do(request)
	if err != nil {
		slog.ErrorContext(ctx, "Audible search request failed", 
			slog.String("title", title), slog.Any("err", err))
		return nil, err
	}

	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			body = []byte{}
		}
		slog.ErrorContext(ctx, "Audible search returned error status", 
			slog.String("title", title),
			slog.Int("statusCode", response.StatusCode),
			slog.String("responseBody", string(body)))
		return nil, fmt.Errorf("unexpected response: status=%d, body=%s", response.StatusCode, body)
	}

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to read Audible search response", 
			slog.String("title", title), slog.Any("err", err))
		return nil, err
	}

	var parsedRepsonse asinSearchResponse

	err = json.Unmarshal(responseBytes, &parsedRepsonse)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to unmarshal Audible search response", 
			slog.String("title", title), slog.Any("err", err))
		return nil, err
	}

	asins := make([]string, len(parsedRepsonse.Products))

	for idx, product := range parsedRepsonse.Products {
		asins[idx] = product.Asin
	}

	slog.InfoContext(ctx, "Successfully searched Audible by title", 
		slog.String("title", title),
		slog.Int("resultCount", len(asins)),
		slog.Int("totalResults", parsedRepsonse.TotalResults))

	return asins, nil
}

func (aac *AudibleApiClient) GetMetadataFromAsin(asin string) (metadata.BookMetadata, error) {
	var md metadata.BookMetadata
	ctx := context.Background()
	url := fmt.Sprintf(asinMetadataUrl, asin)
	
	slog.InfoContext(ctx, "Getting Audible metadata from ASIN", 
		slog.String("asin", asin),
		slog.String("url", url))

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create metadata request", 
			slog.String("asin", asin), slog.Any("err", err))
		return md, err
	}

	response, err := aac.httpClient.Do(request)
	if err != nil {
		slog.ErrorContext(ctx, "Audible metadata request failed", 
			slog.String("asin", asin), slog.Any("err", err))
		return md, err
	}

	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			body = []byte{}
		}
		slog.ErrorContext(ctx, "Audible metadata returned error status", 
			slog.String("asin", asin),
			slog.Int("statusCode", response.StatusCode),
			slog.String("responseBody", string(body)))
		return md, fmt.Errorf("unexpected response: status=%d, body=%s", response.StatusCode, body)
	}

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to read metadata response", 
			slog.String("asin", asin), slog.Any("err", err))
		return md, err
	}

	err = json.Unmarshal(responseBytes, &md)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to unmarshal metadata response", 
			slog.String("asin", asin), slog.Any("err", err))
		return md, err
	}

	slog.InfoContext(ctx, "Successfully retrieved Audible metadata", 
		slog.String("asin", asin),
		slog.String("title", md.Title))

	return md, nil
}
