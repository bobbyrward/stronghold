package api

import (
	"log/slog"
	"net/http"

	"github.com/bobbyrward/stronghold/internal/qbit"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// SearchASINRequest contains the request body for ASIN search
type BookTorrentDLRequest struct {
	Category   string `json:"category" validate:"required"`
	TorrentID  string `json:"torrent_id" validate:"required"`
	TorrentURL string `json:"torrent_url" validate:"required"`
}

// DownloadBookTorrent handles requests to download a book torrent using the provided category and torrent information.
func DownloadBookTorrent(db *gorm.DB, qbitClient *qbit.QbitClient) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var req BookTorrentDLRequest
		if err := BindRequest(c, ctx, &req); err != nil {
			return BadRequest(c, ctx, "invalid request body")
		}

		if err := ValidateRequest(c, ctx, &req); err != nil {
			return BadRequest(c, ctx, "invalid request body")
		}

		if qbitClient == nil {
			client, err := qbit.CreateClient()
			if err != nil {
				return InternalError(c, ctx, "failed to create qBittorrent client", err)
			}

			qbitClient = &client
		}

		slog.InfoContext(ctx, "Downloading book torrent",
			slog.String("category", req.Category),
			slog.Any("torrent_id", req.TorrentID),
		)

		err := (*qbitClient).AddTorrentFromUrlCtx(
			ctx,
			req.TorrentURL,
			map[string]string{
				"autoTMM":  "true",
				"category": req.Category,
			},
		)
		if err != nil {
			return InternalError(c, ctx, "failed to add torrent to qbit", err)
		}

		slog.InfoContext(ctx, "Successfully downloaded book torrent",
			slog.String("category", req.Category),
			slog.String("torrent_id", req.TorrentID),
		)

		return c.JSON(http.StatusOK, map[string]string{})
	}
}
