package api

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/bobbyrward/stronghold/internal/config"
	"github.com/bobbyrward/stronghold/internal/importers/audiobooks"
	"github.com/bobbyrward/stronghold/internal/importers/audiobooks/audible"
	"github.com/bobbyrward/stronghold/internal/importers/audiobooks/metadata"
	"github.com/bobbyrward/stronghold/internal/importers/common"
	"github.com/bobbyrward/stronghold/internal/qbit"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// TorrentImportInfo contains information about a torrent being imported
type TorrentImportInfo struct {
	Hash             string `json:"hash"`
	Name             string `json:"name"`
	Category         string `json:"category"`
	Tags             string `json:"tags"`
	ASIN             string `json:"asin,omitempty"`
	Title            string `json:"title,omitempty"`
	Author           string `json:"author,omitempty"`
	SuggestedLibrary string `json:"suggested_library,omitempty"`
	LocalPath        string `json:"local_path"`
}

// SearchASINRequest contains the request body for ASIN search
type SearchASINRequest struct {
	Title  string `json:"title" validate:"required"`
	Author string `json:"author"`
}

// PreviewDirectoryRequest contains the request body for directory preview
type PreviewDirectoryRequest struct {
	Metadata metadata.BookMetadata `json:"metadata" validate:"required"`
}

// PreviewDirectoryResponse contains the previewed directory name
type PreviewDirectoryResponse struct {
	DirectoryName string `json:"directory_name"`
}

// ExecuteImportRequest contains the request body for executing import
type ExecuteImportRequest struct {
	Hash        string                `json:"hash" validate:"required"`
	Metadata    metadata.BookMetadata `json:"metadata" validate:"required"`
	LibraryName string                `json:"library_name" validate:"required"`
}

// ExecuteImportResponse contains the result of the import operation
type ExecuteImportResponse struct {
	Success         bool   `json:"success"`
	DestinationPath string `json:"destination_path"`
	Message         string `json:"message,omitempty"`
}

// GetTorrentInfo returns information about a torrent
func GetTorrentInfo(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		hash := c.Param("hash")

		slog.InfoContext(ctx, "Getting torrent import info", slog.String("hash", hash))

		// Create qBittorrent client
		qbitClient, err := qbit.CreateClient()
		if err != nil {
			return InternalError(c, ctx, "failed to create qBittorrent client", err)
		}

		// Get torrent by hash
		torrents, err := qbitClient.GetTorrentsCtx(ctx, qbit.TorrentFilterOptions{Hashes: []string{hash}})
		if err != nil || len(torrents) != 1 {
			return GenericNotFound(c, ctx, "torrent not found")
		}

		torrent := torrents[0]

		// Create metadata provider and audible client
		metadataProvider := metadata.NewFFProbeMetadataProvider()
		audibleClient := audible.NewAudibleApiClient()

		// Create importer system
		importer, err := audiobooks.NewAudiobookImporterSystem(
			qbitClient,
			config.Config.Importers,
			metadataProvider,
			audibleClient,
		)
		if err != nil {
			return InternalError(c, ctx, "failed to create audiobook importer", err)
		}

		// Extract metadata from torrent
		bookMetadata, err := importer.ExtractTorrentMetadata(ctx, torrent)

		// Create response with available information
		response := TorrentImportInfo{
			Hash:      torrent.Hash,
			Name:      torrent.Name,
			Category:  torrent.Category,
			Tags:      torrent.Tags,
			LocalPath: common.MapTorrentContentPathToLocalPath(torrent, config.Config.Qbit.DownloadPath, config.Config.Qbit.LocalDownloadPath),
		}

		// If metadata extraction succeeded, include it
		if err == nil {
			response.ASIN = bookMetadata.Asin
			response.Title = bookMetadata.Title
			if len(bookMetadata.Authors) > 0 {
				response.Author = bookMetadata.Authors[0].Name
			}
		}

		// Determine suggested library based on torrent category
		for _, importType := range config.Config.Importers.AudiobookImporter.ImportTypes {
			if importType.Category == torrent.Category {
				response.SuggestedLibrary = importType.Library
				break
			}
		}

		slog.InfoContext(ctx, "Successfully retrieved torrent info",
			slog.String("hash", hash),
			slog.String("name", torrent.Name))

		return c.JSON(http.StatusOK, response)
	}
}

// SearchASIN searches for audiobooks by title and author, returning ASINs and metadata
func SearchASIN(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var req SearchASINRequest
		if err := BindRequest(c, ctx, &req); err != nil {
			return BadRequest(c, ctx, "invalid request body")
		}

		if err := ValidateRequest(c, ctx, &req); err != nil {
			return BadRequest(c, ctx, "invalid request body")
		}

		slog.InfoContext(ctx, "Searching for ASINs",
			slog.String("title", req.Title),
			slog.String("author", req.Author))

		audibleClient := audible.NewAudibleApiClient()

		// Search by title and author
		asins, err := audibleClient.SearchByTitle(req.Title, req.Author)
		if err != nil {
			return InternalError(c, ctx, "failed to search for ASINs", err)
		}

		if len(asins) == 0 {
			slog.InfoContext(ctx, "No ASINs found",
				slog.String("title", req.Title),
				slog.String("author", req.Author))
			return c.JSON(http.StatusOK, []metadata.BookMetadata{})
		}

		// Get full metadata for each ASIN
		results := make([]metadata.BookMetadata, 0, len(asins))
		for _, asin := range asins {
			bookMetadata, err := audibleClient.GetMetadataFromAsin(asin)
			if err != nil {
				slog.WarnContext(ctx, "Failed to get metadata for ASIN",
					slog.String("asin", asin),
					slog.Any("error", err))
				continue
			}
			results = append(results, bookMetadata)
		}

		slog.InfoContext(ctx, "Successfully found ASINs",
			slog.String("title", req.Title),
			slog.Int("count", len(results)))

		return c.JSON(http.StatusOK, results)
	}
}

// GetASINMetadata retrieves full metadata for a specific ASIN
func GetASINMetadata(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		asin := c.Param("asin")

		slog.InfoContext(ctx, "Getting ASIN metadata", slog.String("asin", asin))

		audibleClient := audible.NewAudibleApiClient()

		// Get metadata from Audible
		bookMetadata, err := audibleClient.GetMetadataFromAsin(asin)
		if err != nil {
			return InternalError(c, ctx, "failed to get metadata for ASIN", err)
		}

		slog.InfoContext(ctx, "Successfully retrieved ASIN metadata",
			slog.String("asin", asin),
			slog.String("title", bookMetadata.Title))

		return c.JSON(http.StatusOK, bookMetadata)
	}
}

// PreviewDirectory generates a preview of the directory name for given metadata
func PreviewDirectory(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var req PreviewDirectoryRequest
		if err := BindRequest(c, ctx, &req); err != nil {
			return BadRequest(c, ctx, "invalid request body")
		}

		if err := ValidateRequest(c, ctx, &req); err != nil {
			return BadRequest(c, ctx, "invalid request body")
		}

		slog.InfoContext(ctx, "Previewing directory name",
			slog.String("title", req.Metadata.Title),
			slog.String("asin", req.Metadata.Asin))

		// Generate directory name
		directoryName, err := req.Metadata.GenerateDirectoryName()
		if err != nil {
			return InternalError(c, ctx, "failed to generate directory name", err)
		}

		// Sanitize the directory name (same logic as in importer)
		directoryName = sanitizeName(directoryName)

		response := PreviewDirectoryResponse{
			DirectoryName: directoryName,
		}

		slog.InfoContext(ctx, "Successfully generated directory preview",
			slog.String("directory", directoryName))

		return c.JSON(http.StatusOK, response)
	}
}

// GetLibraries returns the list of available audiobook libraries
func GetLibraries(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		slog.InfoContext(ctx, "Getting audiobook libraries")

		libraries := config.Config.Importers.AudiobookImporter.Libraries

		slog.InfoContext(ctx, "Successfully retrieved libraries",
			slog.Int("count", len(libraries)))

		return c.JSON(http.StatusOK, libraries)
	}
}

// ExecuteImport performs the actual import operation
func ExecuteImport(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var req ExecuteImportRequest
		if err := BindRequest(c, ctx, &req); err != nil {
			return BadRequest(c, ctx, "invalid request body")
		}

		if err := ValidateRequest(c, ctx, &req); err != nil {
			return BadRequest(c, ctx, "invalid request body")
		}

		slog.InfoContext(ctx, "Executing audiobook import",
			slog.String("hash", req.Hash),
			slog.String("library", req.LibraryName),
			slog.String("title", req.Metadata.Title))

		// Create qBittorrent client
		qbitClient, err := qbit.CreateClient()
		if err != nil {
			return InternalError(c, ctx, "failed to create qBittorrent client", err)
		}

		// Get torrent by hash
		torrents, err := qbitClient.GetTorrentsCtx(ctx, qbit.TorrentFilterOptions{Hashes: []string{req.Hash}})
		if err != nil || len(torrents) != 1 {
			return GenericNotFound(c, ctx, "torrent not found")
		}

		torrent := torrents[0]

		// Find the library
		library, ok := config.FindLibraryByName(config.Config.Importers.AudiobookImporter.Libraries, req.LibraryName)
		if !ok {
			return BadRequest(c, ctx, "library not found")
		}

		// Create metadata provider and audible client
		metadataProvider := metadata.NewFFProbeMetadataProvider()
		audibleClient := audible.NewAudibleApiClient()

		// Create importer system
		importer, err := audiobooks.NewAudiobookImporterSystem(
			qbitClient,
			config.Config.Importers,
			metadataProvider,
			audibleClient,
		)
		if err != nil {
			return InternalError(c, ctx, "failed to create audiobook importer", err)
		}

		// Get local path
		localPath := common.MapTorrentContentPathToLocalPath(torrent, config.Config.Qbit.DownloadPath, config.Config.Qbit.LocalDownloadPath)

		// Execute the import
		destinationPath, err := importer.ExecuteImport(ctx, torrent, req.Metadata, library, localPath)
		if err != nil {
			return InternalError(c, ctx, "failed to execute import", err)
		}

		// Add imported tag
		err = qbitClient.AddTagsCtx(ctx, []string{req.Hash}, config.Config.Importers.ImportedTag)
		if err != nil {
			slog.WarnContext(ctx, "Failed to add imported tag",
				slog.String("hash", req.Hash),
				slog.Any("error", err))
		}

		// Remove manual_intervention tag
		err = qbitClient.RemoveTagsCtx(ctx, []string{req.Hash}, config.Config.Importers.ManualInterventionTag)
		if err != nil {
			slog.WarnContext(ctx, "Failed to remove manual_intervention tag",
				slog.String("hash", req.Hash),
				slog.Any("error", err))
		}

		response := ExecuteImportResponse{
			Success:         true,
			DestinationPath: destinationPath,
			Message:         "Import completed successfully",
		}

		slog.InfoContext(ctx, "Successfully executed import",
			slog.String("hash", req.Hash),
			slog.String("destination", destinationPath))

		return c.JSON(http.StatusOK, response)
	}
}

// sanitizeName sanitizes a directory name by replacing invalid characters
func sanitizeName(name string) string {
	return strings.ReplaceAll(name, "/", "-")
}
