package audible

import (
	"encoding/json"
	"fmt"
	"io"
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
	aac := &AudibleApiClient{
		httpClient: &http.Client{},
	}

	return aac
}

func (aac *AudibleApiClient) SearchByTitle(title string) ([]string, error) {
	request, err := http.NewRequest("GET", asinSearchUrl, nil)
	if err != nil {
		return nil, err
	}

	query := request.URL.Query()
	query.Add("num_results", "10")
	query.Add("products_sort_by", "Relevance")
	query.Add("title", title)

	request.URL.RawQuery = query.Encode()

	response, err := aac.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			body = []byte{}
		}
		return nil, fmt.Errorf("unexpected response: status=%d, body=%s", response.StatusCode, body)
	}

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var parsedRepsonse asinSearchResponse

	err = json.Unmarshal(responseBytes, &parsedRepsonse)
	if err != nil {
		return nil, err
	}

	asins := make([]string, len(parsedRepsonse.Products))

	for idx, product := range parsedRepsonse.Products {
		asins[idx] = product.Asin
	}

	return asins, nil
}

func (aac *AudibleApiClient) GetMetadataFromAsin(asin string) (metadata.BookMetadata, error) {
	var md metadata.BookMetadata

	request, err := http.NewRequest("GET", fmt.Sprintf(asinMetadataUrl, asin), nil)
	if err != nil {
		return md, err
	}

	response, err := aac.httpClient.Do(request)
	if err != nil {
		return md, err
	}

	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			body = []byte{}
		}
		return md, fmt.Errorf("unexpected response: status=%d, body=%s", response.StatusCode, body)
	}

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return md, err
	}

	err = json.Unmarshal(responseBytes, &md)
	if err != nil {
		return md, err
	}

	return md, nil
}
