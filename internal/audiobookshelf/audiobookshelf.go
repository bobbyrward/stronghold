package audiobookshelf

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	httpClient *http.Client
	token      string
	baseUrl    string
}

type UnixTime struct {
	time.Time
}

func (u *UnixTime) UnmarshalJSON(b []byte) error {
	var timestamp int64

	err := json.Unmarshal(b, &timestamp)
	if err != nil {
		return err
	}

	u.Time = time.UnixMicro(timestamp)

	return nil
}

func (u UnixTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", u.UnixMicro())), nil
}

type Folder struct {
	Id        string   `json:"id"`
	FullPath  string   `json:"fullPath"`
	LibraryId string   `json:"libraryId"`
	AddedAt   UnixTime `json:"addedAt"`
}

type LibrarySettings struct {
	CoverAspectRation         int    `json:"coverAspectRatio "`
	DisableWatcher            bool   `json:"disableWatcher "`
	SkipMatchingMediaWithAsin bool   `json:"skipMatchingMediaWithAsin "`
	SkipMatchingMediaWithIsbn bool   `json:"skipMatchingMediaWithIsbn "`
	AutoScanCronExpression    string `json:"autoScanCronExpression "`
}

type Library struct {
	Id           string          `json:"id"`
	Name         string          `json:"name"`
	Folders      []Folder        `json:"folders"`
	DisplayOrder int             `json:"displayOrder"`
	Icon         string          `json:"icon"`
	MediaType    string          `json:"mediaType"`
	Provider     string          `json:"provider"`
	Settings     LibrarySettings `json:"settings"`
	CreatedAt    UnixTime        `json:"createdAt"`
	LastUpdate   UnixTime        `json:"lastUpdate"`
}

type AuthorMinified struct {
	Id          string   `json:"id"`
	Asin        string   `json:"asin"`
	Name        string   `json:"name"`
	Description *string  `json:"description,omitempty"`
	ImagePath   *string  `json:"imagePath,omitempty"`
	AddedAt     UnixTime `json:"addedAt"`
	UpdatedAt   UnixTime `json:"updatedAt"`
}

type Series struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Description *string  `json:"desciption,omitempty"`
	AddedAt     UnixTime `json:"addedAt"`
	UpdatedAt   UnixTime `json:"updatedAt"`
}

type LibraryFilterData struct {
	Authors   []AuthorMinified `json:"authors"`
	Genres    []string         `json:"genres"`
	Tags      []string         `json:"tags"`
	Series    []Series         `json:"series"`
	Narrators []string         `json:"narrators"`
	Languages []string         `json:"languages"`
}

type LibraryDetails struct {
	FilterData       LibraryFilterData `json:"filterData"`
	Issues           int               `json:"issues"`
	NumUserPlaylists int               `json:"numUserPlaylists"`
	Library          Library           `json:"library"`
}

func NewClient(baseUrl string, token string) *Client {
	ctx := context.Background()
	baseUrl = strings.TrimSuffix(baseUrl, "/")

	slog.InfoContext(ctx, "Creating Audiobookshelf client", 
		slog.String("baseUrl", baseUrl))

	aac := &Client{
		httpClient: &http.Client{},
		token:      token,
		baseUrl:    baseUrl,
	}

	return aac
}

func (aac *Client) ListLibraries() ([]Library, error) {
	ctx := context.Background()
	url := fmt.Sprintf("%s/api/libraries", aac.baseUrl)
	
	slog.InfoContext(ctx, "Listing Audiobookshelf libraries", slog.String("url", url))
	
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create request", slog.String("url", url), slog.Any("err", err))
		return nil, err
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", aac.token))

	response, err := aac.httpClient.Do(request)
	if err != nil {
		slog.ErrorContext(ctx, "HTTP request failed", slog.String("url", url), slog.Any("err", err))
		return nil, err
	}

	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK {
		slog.ErrorContext(ctx, "Unexpected HTTP status", 
			slog.String("url", url),
			slog.Int("statusCode", response.StatusCode))
		return nil, fmt.Errorf("unexpected response: status=%d", response.StatusCode)
	}

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to read response body", slog.String("url", url), slog.Any("err", err))
		return nil, err
	}

	parsedResponse := struct {
		Libraries []Library `json:"libraries"`
	}{}

	err = json.Unmarshal(responseBytes, &parsedResponse)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to unmarshal response", slog.String("url", url), slog.Any("err", err))
		return nil, err
	}

	slog.InfoContext(ctx, "Successfully listed libraries", 
		slog.String("url", url),
		slog.Int("libraryCount", len(parsedResponse.Libraries)))

	return parsedResponse.Libraries, nil
}

func (aac *Client) ScanLibrary(libraryId string, force bool) error {
	ctx := context.Background()
	url := fmt.Sprintf("%s/api/libraries/%s/scan", aac.baseUrl, libraryId)
	
	slog.InfoContext(ctx, "Scanning Audiobookshelf library", 
		slog.String("url", url),
		slog.String("libraryId", libraryId),
		slog.Bool("force", force))
	
	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create scan request", slog.String("url", url), slog.Any("err", err))
		return err
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", aac.token))

	response, err := aac.httpClient.Do(request)
	if err != nil {
		slog.ErrorContext(ctx, "HTTP scan request failed", slog.String("url", url), slog.Any("err", err))
		return err
	}

	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK {
		slog.ErrorContext(ctx, "Scan request failed", 
			slog.String("url", url),
			slog.Int("statusCode", response.StatusCode))
		return fmt.Errorf("unexpected response: status=%d", response.StatusCode)
	}

	slog.InfoContext(ctx, "Successfully triggered library scan", 
		slog.String("libraryId", libraryId))
	return nil
}

func (aac *Client) GetLibrary(libraryId string, includeFilterData bool) (Library, error) {
	var library Library
	ctx := context.Background()

	url := fmt.Sprintf("%s/api/libraries/%s", aac.baseUrl, libraryId)
	
	slog.InfoContext(ctx, "Getting Audiobookshelf library details", 
		slog.String("url", url),
		slog.String("libraryId", libraryId),
		slog.Bool("includeFilterData", includeFilterData))

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create get library request", slog.String("url", url), slog.Any("err", err))
		return library, err
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", aac.token))

	response, err := aac.httpClient.Do(request)
	if err != nil {
		slog.ErrorContext(ctx, "HTTP get library request failed", slog.String("url", url), slog.Any("err", err))
		return library, err
	}

	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK {
		slog.ErrorContext(ctx, "Get library request failed", 
			slog.String("url", url),
			slog.Int("statusCode", response.StatusCode))
		return library, fmt.Errorf("unexpected response: status=%d", response.StatusCode)
	}

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to read get library response", slog.String("url", url), slog.Any("err", err))
		return library, err
	}

	err = json.Unmarshal(responseBytes, &library)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to unmarshal library response", slog.String("url", url), slog.Any("err", err))
		return library, err
	}

	slog.InfoContext(ctx, "Successfully retrieved library details", 
		slog.String("libraryId", libraryId),
		slog.String("libraryName", library.Name))
	return library, nil
}
