// Package booksearch implements the Book Search service.
package booksearch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/bobbyrward/stronghold/internal/config"
	"github.com/bobbyrward/stronghold/internal/cookies"
	"github.com/bobbyrward/stronghold/internal/models"
	"github.com/cappuccinotm/slogx/logger"
	"github.com/carlmjohnson/requests"
	"gorm.io/gorm"
)

type BookSearchService struct{}

func NewBookSearchService() *BookSearchService {
	return &BookSearchService{}
}

func noneOf(options []bool) bool {
	for _, v := range options {
		if v {
			return false
		}
	}

	return true
}

func (params *SearchParameters) Validate() error {
	searchOptions := []bool{params.Query != "", params.Hash != "", params.ID != nil}
	if noneOf(searchOptions) {
		return fmt.Errorf("one of Query or Hash or ID is required")
	}

	return nil
}

func createHTTPClient(enableLogging bool) *http.Client {
	searchConfig := &config.Config.BookSearch

	transport := &http.Transport{}

	if searchConfig.HttpsProxy != "" {
		transport.Proxy = http.ProxyURL(&url.URL{
			Scheme: "http",
			Host:   searchConfig.HttpsProxy,
		})
	}

	client := &http.Client{
		Timeout:   15 * time.Second,
		Transport: transport,
	}

	if enableLogging {
		l := logger.New(
			logger.WithLogger(slog.Default()),
			logger.WithBody(10240),
		)

		client.Transport = l.HTTPClientRoundTripper(transport)
	}

	return client
}

func (s *BookSearchService) Search(ctx context.Context, db *gorm.DB, params *SearchParameters) (*SearchResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params is required")
	}

	searchConfig := &config.Config.BookSearch

	if searchConfig.TokenCookieName == "" {
		return nil, errors.New("tokenCookieName not configured")
	}

	err := params.Validate()
	if err != nil {
		return nil, err
	}

	// Get API key from database
	credential, err := models.GetBookSearchCredential(db)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get book search credential", slog.Any("err", err))
		return nil, fmt.Errorf("book search credential not found in database")
	}

	baseURL := searchConfig.BaseURL
	searchEndpoint := searchConfig.SearchEndpoint

	url := fmt.Sprintf("%s%s", baseURL, searchEndpoint)

	request := SearchRequest{
		IncludeDlLink:      true,
		IncludeDescription: true,
		IncludeISBN:        true,
		PerPage:            params.MaxResults,
		Tor: SearchRequestTor{
			ID:             params.ID,
			BrowseLang:     []int{1},
			MainCategories: []int{13, 14},
			Categories:     []int{41, 47, 108, 63, 69, 0, 109, 108},
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
		Client(createHTTPClient(false)).
		Cookie(searchConfig.TokenCookieName, credential.APIKey).
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

type refreshTokenResponse struct {
	Success   bool   `json:"Success"`
	Message   string `json:"msg"`
	IPAddress string `json:"ip"`
	ASN       int    `json:"ASN"`
	AS        string `json:"AS"`
}

func (s *BookSearchService) RefreshToken(ctx context.Context, db *gorm.DB) error {
	searchConfig := &config.Config.BookSearch

	if searchConfig.TokenRefreshURL == "" {
		return errors.New("tokenRefreshUrl not configured")
	}

	if searchConfig.CookieDomain == "" {
		return errors.New("cookieDomain not configured")
	}

	if searchConfig.TokenCookieName == "" {
		return errors.New("tokenCookieName not configured")
	}

	// Get API key from database
	credential, err := models.GetBookSearchCredential(db)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get book search credential", slog.Any("err", err))
		return fmt.Errorf("book search credential not found in database")
	}

	slog.InfoContext(ctx, "Token refresh requested")

	client := createHTTPClient(false)
	client.Jar = requests.NewCookieJar()

	var response refreshTokenResponse

	err = requests.
		URL(searchConfig.TokenRefreshURL).
		Client(client).
		Cookie(searchConfig.TokenCookieName, credential.APIKey).
		ToJSON(&response).
		Fetch(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to refresh token", slog.Any("err", err))
		return errors.Join(err, fmt.Errorf("unable to search"))
	}

	responseCookies := client.Jar.Cookies(&url.URL{
		Scheme: "https",
		Host:   searchConfig.CookieDomain,
	})

	tokenCookie, found := cookies.FindCookieByName(responseCookies, searchConfig.TokenCookieName)
	if !found {
		slog.ErrorContext(ctx, "Token cookie not found in response", slog.String("cookieName", searchConfig.TokenCookieName))
		return fmt.Errorf("token cookie '%s' not found in response", searchConfig.TokenCookieName)
	}

	err = models.UpsertBookSearchCredential(db, tokenCookie.Value, response.IPAddress, fmt.Sprintf("%d", response.ASN))
	if err != nil {
		slog.ErrorContext(ctx, "Failed to upsert book search credential", slog.Any("err", err))
		return errors.Join(err, fmt.Errorf("failed to upsert book search credential"))
	}

	return nil
}
