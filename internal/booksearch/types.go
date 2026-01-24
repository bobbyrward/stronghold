package booksearch

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/bobbyrward/stronghold/internal/config"
	"github.com/bobbyrward/stronghold/internal/models"
)

const (
	MainCategoryEbooks     = 14
	MainCategoryAudiobooks = 13
)

type SearchRequest struct {
	IncludeDlLink      bool             `json:"dlLink,omitempty"`
	IncludeDescription bool             `json:"description,omitempty"`
	IncludeISBN        bool             `json:"isbn,omitempty"`
	PerPage            int              `json:"perpage,omitempty"`
	Tor                SearchRequestTor `json:"tor"`
}

type SearchRequestTor struct {
	ID             *int                `json:"id,omitempty"`
	BrowseLang     []int               `json:"browse_lang,omitempty"`
	MainCategories []int               `json:"main_cat,omitempty"`
	Categories     []int               `json:"cat,omitempty"`
	SearchType     string              `json:"searchType,omitempty"`
	SortType       string              `json:"sortType,omitempty"`
	SearchIn       SearchRequestSrchIn `json:"srchIn"`
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

	ID *int

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
	Description     string `json:"description"`
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

		AuthorInfo   string `json:"author_info"`
		NarratorInfo string `json:"narrator_info"`
		SeriesInfo   string `json:"series_info"`
		ISBNInner    any    `json:"isbn"`
	}{
		Alias: (*Alias)(sri),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	slog.Debug("Unmarshalling raw", "raw", aux)

	if aux.ISBNInner != nil {
		switch aux.ISBNInner.(type) {
		case string:
			sri.ISBN, _ = aux.ISBNInner.(string)
		case float64:
			v, _ := aux.ISBNInner.(float64)
			sri.ISBN = fmt.Sprintf("%v", int(v))
		default:
			return errors.New("unknown type for isbn field")
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
		seriesTemp := make(map[string][]any)

		if err := json.Unmarshal([]byte(aux.SeriesInfo), &seriesTemp); err != nil {
			return nil
		}

		sri.Series = make(map[string]SeriesEntry)

		for key, value := range seriesTemp {
			if len(value) != 3 {
				slog.Debug("Skipping series entry with unexpected length", "key", key, "value", value)
				continue
			}

			value0, value0ok := value[0].(string)
			value1, value1ok := value[1].(string)

			if !value0ok || !value1ok {
				slog.Debug("Skipping series entry with unexpected types", "key", key, "value", value)
				continue
			}

			index := 0
			index, _ = strconv.Atoi(value1)

			sri.Series[key] = SeriesEntry{Name: value0, Index: index}
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
