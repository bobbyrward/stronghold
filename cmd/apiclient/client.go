package apiclient

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/carlmjohnson/requests"
)

// APIClient handles HTTP requests to the API
type APIClient struct {
	BaseURL string
}

// NewClient creates a new API client
func NewClient(baseURL string) *APIClient {
	return &APIClient{BaseURL: baseURL}
}

// Get performs a GET request
func (c *APIClient) Get(ctx context.Context, path string, response interface{}) error {
	url := c.BaseURL + path
	slog.DebugContext(ctx, "GET request", slog.String("url", url))

	err := requests.
		URL(url).
		ToJSON(response).
		Fetch(ctx)

	if err != nil {
		slog.ErrorContext(ctx, "GET request failed", slog.String("url", url), slog.Any("error", err))
		return fmt.Errorf("GET %s failed: %w", path, err)
	}

	slog.DebugContext(ctx, "GET request successful", slog.String("url", url))
	return nil
}

// Post performs a POST request
func (c *APIClient) Post(ctx context.Context, path string, body interface{}, response interface{}) error {
	url := c.BaseURL + path
	slog.DebugContext(ctx, "POST request", slog.String("url", url))

	err := requests.
		URL(url).
		BodyJSON(body).
		ToJSON(response).
		Fetch(ctx)

	if err != nil {
		slog.ErrorContext(ctx, "POST request failed", slog.String("url", url), slog.Any("error", err))
		return fmt.Errorf("POST %s failed: %w", path, err)
	}

	slog.DebugContext(ctx, "POST request successful", slog.String("url", url))
	return nil
}

// Put performs a PUT request
func (c *APIClient) Put(ctx context.Context, path string, body interface{}, response interface{}) error {
	url := c.BaseURL + path
	slog.DebugContext(ctx, "PUT request", slog.String("url", url))

	err := requests.
		URL(url).
		Put().
		BodyJSON(body).
		ToJSON(response).
		Fetch(ctx)

	if err != nil {
		slog.ErrorContext(ctx, "PUT request failed", slog.String("url", url), slog.Any("error", err))
		return fmt.Errorf("PUT %s failed: %w", path, err)
	}

	slog.DebugContext(ctx, "PUT request successful", slog.String("url", url))
	return nil
}

// Delete performs a DELETE request
func (c *APIClient) Delete(ctx context.Context, path string) error {
	url := c.BaseURL + path
	slog.DebugContext(ctx, "DELETE request", slog.String("url", url))

	err := requests.
		URL(url).
		Delete().
		Fetch(ctx)

	if err != nil {
		slog.ErrorContext(ctx, "DELETE request failed", slog.String("url", url), slog.Any("error", err))
		return fmt.Errorf("DELETE %s failed: %w", path, err)
	}

	slog.DebugContext(ctx, "DELETE request successful", slog.String("url", url))
	return nil
}
