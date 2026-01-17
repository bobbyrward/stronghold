package api

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/bobbyrward/stronghold/internal/hardcover"
)

// HardcoverAuthorSearchResponse represents an author result from Hardcover search.
type HardcoverAuthorSearchResponse struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

// SearchHardcoverAuthors handles GET /hardcover/authors/search
func SearchHardcoverAuthors(client hardcover.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		query := c.QueryParam("q")

		if query == "" {
			return BadRequest(c, ctx, "Query parameter 'q' is required")
		}

		slog.InfoContext(ctx, "Searching Hardcover authors", slog.String("query", query))

		results, err := client.SearchAuthors(ctx, query)
		if err != nil {
			return InternalError(c, ctx, "Failed to search Hardcover", err)
		}

		response := make([]HardcoverAuthorSearchResponse, len(results))
		for i, r := range results {
			response[i] = HardcoverAuthorSearchResponse{
				Slug: r.Slug,
				Name: r.Name,
			}
		}

		slog.InfoContext(ctx, "Hardcover search completed", slog.Int("results", len(response)))
		return c.JSON(http.StatusOK, response)
	}
}
