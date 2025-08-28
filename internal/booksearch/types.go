package booksearch

import (
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/bobbyrward/stronghold/internal/config"
	"github.com/bobbyrward/stronghold/internal/models"
)

const (
	MainCategory_Ebooks     = 14
	MainCategory_Audiobooks = 13
)

type SearchRequest struct {
	IncludeDlLink      bool             `json:"dlLink,omitempty"`
	IncludeDescription bool             `json:"description,omitempty"`
	IncludeISBN        bool             `json:"isbn,omitempty"`
	PerPage            int              `json:"perpage,omitempty"`
	Tor                SearchRequestTor `json:"tor"`
}

type SearchRequestTor struct {
	BrowseLang     []int               `json:"browse_lang,omitempty"`
	MainCategories []int               `json:"main_cat,omitempty"`
	Categories     []int               `json:"cat,omitempty"`
	SearchType     string              `json:"searchType,omitempty"`
	SortType       string              `json:"sortType,omitempty"`
	SearchIn       SearchRequestSrchIn `json:"srchIn,omitempty"`
	Query          string              `json:"text,omitempty"`
	Hash           string              `json:"hash,omitempty"`
}

type SearchRequestSrchIn struct {
	Author bool `json:"author,omitempty"`
	Series bool `json:"series,omitempty"`
	Title  bool `json:"title,omitempty"`
}

type SearchParameters struct {
	Query      string
	Hash       string
	MaxResults int

	// Offset             int
	// DateOffset         string

	// Languages    []int
	// Categories   []int
	// EndDate      string
	// Hash         string
	// MainCategory []int
}

type SearchResponse struct {
	Total      int                  `json:"total"`
	TotalFound int                  `json:"total_found"`
	Data       []SearchResponseItem `json:"data"`
}

type SeriesEntry struct {
	Name  string
	Index int
}

type SearchResponseItem struct {
	// Added           time.Time         `json:"added"`
	Added           string `json:"added"`
	Authors         map[string]string
	CategoryDisplay string `json:"cat"`
	Category        int    `json:"category"`
	CategoryName    string `json:"catname"`
	DlHash          string `json:"dl"`
	FileTypes       string `json:"filetype"`
	ID              int    `json:"id"`
	ISBN            string
	LanguageCode    string `json:"lang_code"`
	Language        int    `json:"language"`
	Title           string `json:"title"`
	Series          map[string]SeriesEntry
	Narrators       map[string]string
	MainCategory    int    `json:"main_cat"`
	Leechers        int    `json:"leechers"`
	Seeders         int    `json:"seeders"`
	NumFiles        int    `json:"numfiles"`
	Size            string `json:"size"`
	Tags            string `json:"tags"`
}

func (sri *SearchResponseItem) DownloadTorrentURL() string {
	return fmt.Sprintf("%s/tor/download.php/%s", config.Config.BookSearch.BaseURL, sri.DlHash)
}

func (sri *SearchResponseItem) UnmarshalJSON(data []byte) error {
	type Alias SearchResponseItem

	aux := &struct {
		*Alias

		AuthorInfo   string      `json:"author_info"`
		NarratorInfo string      `json:"narrator_info"`
		SeriesInfo   string      `json:"series_info"`
		ISBN_inner   interface{} `json:"isbn"`
	}{
		Alias: (*Alias)(sri),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	if aux.ISBN_inner != nil {
		switch aux.ISBN_inner.(type) {
		case string:
			sri.ISBN, _ = aux.ISBN_inner.(string)
		case float64:
			v, _ := aux.ISBN_inner.(float64)
			sri.ISBN = fmt.Sprintf("%v", int(v))
		default:
			panic(fmt.Sprintf("Unable to parse isbn: %v", aux.ISBN_inner))
		}
	}

	if aux.AuthorInfo != "" {
		if err := json.Unmarshal([]byte(aux.AuthorInfo), &sri.Authors); err != nil {
			return err
		}
	}

	if aux.NarratorInfo != "" {
		if err := json.Unmarshal([]byte(aux.NarratorInfo), &sri.Narrators); err != nil {
			return err
		}
	}

	if aux.SeriesInfo != "" {
		series_temp := make(map[string][]string)

		if err := json.Unmarshal([]byte(aux.SeriesInfo), &series_temp); err != nil {
			return nil
		}

		sri.Series = make(map[string]SeriesEntry)

		for key, value := range series_temp {
			index := 0
			index, _ = strconv.Atoi(value[1])

			sri.Series[key] = SeriesEntry{Name: value[0], Index: index}
		}
	}

	return nil
}

func (sri *SearchResponseItem) ToModel() models.SearchResponseItem {
	model := models.SearchResponseItem{
		Title:        sri.Title,
		Category:     sri.CategoryName,
		MainCategory: sri.MainCategory,
		DlHash:       sri.DlHash,
		FileTypes:    sri.FileTypes,
		TorrentID:    uint(sri.ID),
		Language:     sri.LanguageCode,
		Size:         sri.Size,
		Tags:         sri.Tags,
		Authors:      strings.Join(slices.Collect(maps.Values(sri.Authors)), ", "),
	}

	if len(sri.Series) > 0 {
		seriesCollected := make([]string, 0)

		for _, value := range sri.Series {
			seriesCollected = append(seriesCollected, fmt.Sprintf("%s(%d)", value.Name, value.Index))
		}

		model.Series = strings.Join(seriesCollected, ", ")
	}

	if len(sri.Narrators) > 0 {
		model.Narrators = strings.Join(slices.Collect(maps.Values(sri.Narrators)), ", ")
	}

	return model
}
