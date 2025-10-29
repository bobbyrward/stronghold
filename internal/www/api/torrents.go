package api

import (
	"log/slog"
	"net/http"
	"slices"
	"strings"

	"github.com/bobbyrward/stronghold/internal/config"
	"github.com/bobbyrward/stronghold/internal/qbit"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type TorrentResponse struct {
	Hash     string            `json:"hash"`
	Name     string            `json:"name"`
	Category string            `json:"category"`
	State    qbit.TorrentState `json:"state"`
	Tags     string            `json:"tags"`
}

type TorrentChangeCategoryRequest struct {
	Category string `json:"category" validate:"required"`
}

type TorrentChangeTagsRequest struct {
	Tags string `json:"tags" validate:"required"`
}

func responseFromTorrent(torrent qbit.Torrent) TorrentResponse {
	return TorrentResponse{
		Hash:     torrent.Hash,
		Name:     torrent.Name,
		Category: torrent.Category,
		State:    torrent.State,
		Tags:     torrent.Tags,
	}
}

func sortTorrents(torrents []qbit.Torrent) {
	slices.SortStableFunc(torrents, func(a, b qbit.Torrent) int {
		cmp := strings.Compare(string(a.State), string(b.State))

		switch cmp {
		case 0:
			return strings.Compare(a.Category, b.Category)
		default:
			return cmp

		}
	})
}

func ListUnimportedTorrents(
	db *gorm.DB,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		slog.InfoContext(ctx, "Listing unimported torrents")

		qbitClient, err := qbit.CreateClient()
		if err != nil {
			return InternalError(c, ctx, "failed to create qBittorrent client", err)
		}

		torrents, err := qbit.GetUnimportedTorrents(ctx, qbitClient)
		if err != nil {
			return InternalError(c, ctx, "failed to get torrents from qBittorrent", err)
		}

		responses := make([]TorrentResponse, len(torrents))

		sortTorrents(torrents)

		for i, t := range torrents {
			responses[i] = responseFromTorrent(t)
		}

		return c.JSON(http.StatusOK, responses)
	}
}

func ListManualInterventionTorrents(
	db *gorm.DB,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		slog.InfoContext(ctx, "Listing manual intervention torrents")

		qbitClient, err := qbit.CreateClient()
		if err != nil {
			return InternalError(c, ctx, "failed to create qBittorrent client", err)
		}

		torrents, err := qbit.GetManualInterventionTorrents(
			ctx,
			qbitClient,
			config.Config.Importers.ManualInterventionTag,
		)
		if err != nil {
			return InternalError(c, ctx, "failed to get torrents from qBittorrent", err)
		}

		responses := make([]TorrentResponse, len(torrents))

		sortTorrents(torrents)

		for i, t := range torrents {
			responses[i] = responseFromTorrent(t)
		}

		return c.JSON(http.StatusOK, responses)
	}
}

func SetTorrentCategory(
	db *gorm.DB,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		hash := c.Param("hash")

		var req TorrentChangeCategoryRequest
		if err := BindRequest(c, ctx, &req); err != nil {
			return BadRequest(c, ctx, "Invalid request body")
		}

		if err := ValidateRequest(c, ctx, &req); err != nil {
			return BadRequest(c, ctx, "Invalid request body")
		}

		qbitClient, err := qbit.CreateClient()
		if err != nil {
			return InternalError(c, ctx, "failed to create qBittorrent client", err)
		}

		torrents, err := qbitClient.GetTorrentsCtx(ctx, qbit.TorrentFilterOptions{Hashes: []string{hash}})
		if len(torrents) != 1 || err != nil {
			return GenericNotFound(c, ctx, "torrent not found")
		}

		err = qbitClient.SetCategoryCtx(ctx, []string{hash}, req.Category)
		if err != nil {
			return InternalError(c, ctx, "failed to set torrent category", err)
		}

		return c.NoContent(http.StatusNoContent)
	}
}

func SetTorrentTags(
	db *gorm.DB,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		hash := c.Param("hash")

		var req TorrentChangeTagsRequest
		if err := BindRequest(c, ctx, &req); err != nil {
			return BadRequest(c, ctx, "Invalid request body")
		}

		if err := ValidateRequest(c, ctx, &req); err != nil {
			return BadRequest(c, ctx, "Invalid request body")
		}

		qbitClient, err := qbit.CreateClient()
		if err != nil {
			return InternalError(c, ctx, "failed to create qBittorrent client", err)
		}

		torrents, err := qbitClient.GetTorrentsCtx(ctx, qbit.TorrentFilterOptions{Hashes: []string{hash}})
		if len(torrents) != 1 || err != nil {
			return GenericNotFound(c, ctx, "torrent not found")
		}

		err = qbitClient.SetTags(ctx, []string{hash}, req.Tags)
		if err != nil {
			return InternalError(c, ctx, "failed to set torrent tags", err)
		}

		return c.NoContent(http.StatusNoContent)
	}
}
