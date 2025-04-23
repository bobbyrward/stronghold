package audiobookshelf

import (
	"encoding/json"
	"fmt"
	"io"
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
	baseUrl = strings.TrimSuffix(baseUrl, "/")

	aac := &Client{
		httpClient: &http.Client{},
		token:      token,
		baseUrl:    baseUrl,
	}

	return aac
}

func (aac *Client) ListLibraries() ([]Library, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/api/libraries", aac.baseUrl), nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", aac.token))

	response, err := aac.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response: status=%d", response.StatusCode)
	}

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	parsedResponse := struct {
		Libraries []Library `json:"libraries"`
	}{}

	err = json.Unmarshal(responseBytes, &parsedResponse)
	if err != nil {
		return nil, err
	}

	return parsedResponse.Libraries, nil
}

func (aac *Client) ScanLibrary(libraryId string, force bool) error {
	request, err := http.NewRequest("POST", fmt.Sprintf("%s/api/libraries/%s/scan", aac.baseUrl, libraryId), nil)
	if err != nil {
		return err
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", aac.token))

	response, err := aac.httpClient.Do(request)
	if err != nil {
		return err
	}

	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response: status=%d", response.StatusCode)
	}

	return nil
}

func (aac *Client) GetLibrary(libraryId string, includeFilterData bool) (Library, error) {
	var library Library

	url := fmt.Sprintf("%s/api/libraries/%s", aac.baseUrl, libraryId)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return library, err
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", aac.token))

	response, err := aac.httpClient.Do(request)
	if err != nil {
		return library, err
	}

	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK {
		return library, fmt.Errorf("unexpected response: status=%d", response.StatusCode)
	}

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return library, err
	}

	err = json.Unmarshal(responseBytes, &library)
	if err != nil {
		return library, err
	}

	return library, nil
}
