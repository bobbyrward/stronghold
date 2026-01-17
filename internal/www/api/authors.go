package api

import (
	"context"
	"log/slog"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/hardcover"
	"github.com/bobbyrward/stronghold/internal/models"
)

type AuthorRequest struct {
	Name         string  `json:"name" validate:"required"`
	HardcoverRef *string `json:"hardcover_ref"`
}

type AuthorResponse struct {
	ID           uint    `json:"id"`
	Name         string  `json:"name"`
	HardcoverRef *string `json:"hardcover_ref"`
}

type AuthorHandler struct {
	hardcoverClient hardcover.Client
}

func (handler AuthorHandler) ModelToResponse(c echo.Context, ctx context.Context, db *gorm.DB, row models.Author) AuthorResponse {
	return AuthorResponse{
		ID:           row.ID,
		Name:         row.Name,
		HardcoverRef: row.HardcoverRef,
	}
}

func (h AuthorHandler) RequestToModel(c echo.Context, ctx context.Context, db *gorm.DB, req AuthorRequest) (models.Author, error) {
	// Validate hardcover_ref if provided
	if req.HardcoverRef != nil && *req.HardcoverRef != "" {
		slog.DebugContext(ctx, "Validating hardcover_ref", slog.String("ref", *req.HardcoverRef))
		author, err := h.hardcoverClient.GetAuthorBySlug(ctx, *req.HardcoverRef)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to validate hardcover_ref", slog.Any("error", err))
			return models.Author{}, InternalError(c, ctx, "Failed to validate hardcover_ref", err)
		}
		if author == nil {
			slog.WarnContext(ctx, "Invalid hardcover_ref", slog.String("ref", *req.HardcoverRef))
			return models.Author{}, BadRequest(c, ctx, "invalid hardcover_ref: author not found on Hardcover")
		}
		slog.InfoContext(ctx, "Hardcover ref validated", slog.String("ref", *req.HardcoverRef), slog.String("name", author.Name))
	}

	return models.Author{
		Name:         req.Name,
		HardcoverRef: req.HardcoverRef,
	}, nil
}

func (h AuthorHandler) UpdateModel(c echo.Context, ctx context.Context, db *gorm.DB, row *models.Author, req AuthorRequest) error {
	// Validate hardcover_ref if provided and changed
	if req.HardcoverRef != nil && *req.HardcoverRef != "" {
		// Only validate if it's different from current
		if row.HardcoverRef == nil || *row.HardcoverRef != *req.HardcoverRef {
			slog.DebugContext(ctx, "Validating hardcover_ref on update", slog.String("ref", *req.HardcoverRef))
			author, err := h.hardcoverClient.GetAuthorBySlug(ctx, *req.HardcoverRef)
			if err != nil {
				slog.ErrorContext(ctx, "Failed to validate hardcover_ref", slog.Any("error", err))
				return InternalError(c, ctx, "Failed to validate hardcover_ref", err)
			}
			if author == nil {
				slog.WarnContext(ctx, "Invalid hardcover_ref", slog.String("ref", *req.HardcoverRef))
				return BadRequest(c, ctx, "invalid hardcover_ref: author not found on Hardcover")
			}
			slog.InfoContext(ctx, "Hardcover ref validated", slog.String("ref", *req.HardcoverRef), slog.String("name", author.Name))
		}
	}

	row.Name = req.Name
	row.HardcoverRef = req.HardcoverRef
	return nil
}

func (h AuthorHandler) ParseQuery(c echo.Context, ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	query := c.QueryParam("q")
	if query == "" {
		return db, nil
	}

	slog.DebugContext(ctx, "Fuzzy searching authors", slog.String("query", query))

	// Case-insensitive search using LOWER() for SQLite/PostgreSQL compatibility
	pattern := "%" + strings.ToLower(query) + "%"

	// Search author name OR any alias name via subquery
	aliasSubquery := db.Model(&models.AuthorAlias{}).
		Select("author_id").
		Where("LOWER(name) LIKE ?", pattern)

	db = db.Where("LOWER(name) LIKE ? OR id IN (?)", pattern, aliasSubquery)

	return db, nil
}

func (handler AuthorHandler) PreloadRelations(c echo.Context, ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	return db, nil
}

func (handler AuthorHandler) IDFromModel(row models.Author) uint {
	return row.ID
}

// ListAuthors returns all authors
func ListAuthors(db *gorm.DB, hc hardcover.Client) echo.HandlerFunc {
	return genericListHandler[models.Author, AuthorRequest, AuthorResponse](db, AuthorHandler{hardcoverClient: hc})
}

// CreateAuthor creates a new author
func CreateAuthor(db *gorm.DB, hc hardcover.Client) echo.HandlerFunc {
	return genericCreateHandler[models.Author, AuthorRequest, AuthorResponse](db, AuthorHandler{hardcoverClient: hc})
}

// GetAuthor returns a single author by ID
func GetAuthor(db *gorm.DB, hc hardcover.Client) echo.HandlerFunc {
	return genericGetHandler[models.Author, AuthorRequest, AuthorResponse](db, AuthorHandler{hardcoverClient: hc})
}

// UpdateAuthor updates an existing author
func UpdateAuthor(db *gorm.DB, hc hardcover.Client) echo.HandlerFunc {
	return genericUpdateHandler[models.Author, AuthorRequest, AuthorResponse](db, AuthorHandler{hardcoverClient: hc})
}

// DeleteAuthor deletes an author
func DeleteAuthor(db *gorm.DB) echo.HandlerFunc {
	return genericDeleteHandler[models.Author](db)
}
