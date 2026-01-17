package api

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
)

type AuthorAliasRequest struct {
	Name string `json:"name" validate:"required"`
}

type AuthorAliasResponse struct {
	ID       uint   `json:"id"`
	AuthorID uint   `json:"author_id"`
	Name     string `json:"name"`
}

func aliasToResponse(alias models.AuthorAlias) AuthorAliasResponse {
	return AuthorAliasResponse{
		ID:       alias.ID,
		AuthorID: alias.AuthorID,
		Name:     alias.Name,
	}
}

// ListAuthorAliases returns all aliases for a specific author
func ListAuthorAliases(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		authorID, err := ParseAuthorIDParam(c, ctx)
		if err != nil {
			return BadRequest(c, ctx, "Invalid author_id")
		}

		// Verify author exists
		var author models.Author
		if err := db.First(&author, authorID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return NotFound(c, ctx, "Author", authorID)
			}
			return InternalError(c, ctx, "Failed to query author", err)
		}

		var aliases []models.AuthorAlias
		if err := db.Where("author_id = ?", authorID).Find(&aliases).Error; err != nil {
			return InternalError(c, ctx, "Failed to list aliases", err)
		}

		response := make([]AuthorAliasResponse, len(aliases))
		for i, a := range aliases {
			response[i] = aliasToResponse(a)
		}

		slog.InfoContext(ctx, "Listed author aliases", slog.Uint64("author_id", uint64(authorID)), slog.Int("count", len(aliases)))
		return c.JSON(http.StatusOK, response)
	}
}

// CreateAuthorAlias creates a new alias for an author
func CreateAuthorAlias(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		authorID, err := ParseAuthorIDParam(c, ctx)
		if err != nil {
			return BadRequest(c, ctx, "Invalid author_id")
		}

		// Verify author exists
		var author models.Author
		if err := db.First(&author, authorID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return NotFound(c, ctx, "Author", authorID)
			}
			return InternalError(c, ctx, "Failed to query author", err)
		}

		var req AuthorAliasRequest
		if err := BindRequest(c, ctx, &req); err != nil {
			return BadRequest(c, ctx, "Invalid request body")
		}
		if err := ValidateRequest(c, ctx, &req); err != nil {
			return BadRequest(c, ctx, "Validation failed")
		}

		alias := models.AuthorAlias{
			AuthorID: authorID,
			Name:     req.Name,
		}

		if err := db.Create(&alias).Error; err != nil {
			slog.ErrorContext(ctx, "Failed to create alias", slog.Any("error", err))
			return InternalError(c, ctx, "Failed to create alias", err)
		}

		slog.InfoContext(ctx, "Created author alias", slog.Uint64("id", uint64(alias.ID)), slog.Uint64("author_id", uint64(authorID)))
		return c.JSON(http.StatusCreated, aliasToResponse(alias))
	}
}

// GetAuthorAlias returns a single alias by ID
func GetAuthorAlias(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		authorID, err := ParseAuthorIDParam(c, ctx)
		if err != nil {
			return BadRequest(c, ctx, "Invalid author_id")
		}

		id, err := ParseIDParam(c, ctx)
		if err != nil {
			return BadRequest(c, ctx, "Invalid id")
		}

		var alias models.AuthorAlias
		if err := db.First(&alias, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return NotFound(c, ctx, "Alias", id)
			}
			return InternalError(c, ctx, "Failed to query alias", err)
		}

		// Verify alias belongs to this author
		if alias.AuthorID != authorID {
			return NotFound(c, ctx, "Alias", id)
		}

		return c.JSON(http.StatusOK, aliasToResponse(alias))
	}
}

// UpdateAuthorAlias updates an existing alias
func UpdateAuthorAlias(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		authorID, err := ParseAuthorIDParam(c, ctx)
		if err != nil {
			return BadRequest(c, ctx, "Invalid author_id")
		}

		id, err := ParseIDParam(c, ctx)
		if err != nil {
			return BadRequest(c, ctx, "Invalid id")
		}

		var alias models.AuthorAlias
		if err := db.First(&alias, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return NotFound(c, ctx, "Alias", id)
			}
			return InternalError(c, ctx, "Failed to query alias", err)
		}

		// Verify alias belongs to this author
		if alias.AuthorID != authorID {
			return NotFound(c, ctx, "Alias", id)
		}

		var req AuthorAliasRequest
		if err := BindRequest(c, ctx, &req); err != nil {
			return BadRequest(c, ctx, "Invalid request body")
		}
		if err := ValidateRequest(c, ctx, &req); err != nil {
			return BadRequest(c, ctx, "Validation failed")
		}

		alias.Name = req.Name
		if err := db.Save(&alias).Error; err != nil {
			return InternalError(c, ctx, "Failed to update alias", err)
		}

		slog.InfoContext(ctx, "Updated author alias", slog.Uint64("id", uint64(id)))
		return c.JSON(http.StatusOK, aliasToResponse(alias))
	}
}

// DeleteAuthorAlias deletes an alias
func DeleteAuthorAlias(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		authorID, err := ParseAuthorIDParam(c, ctx)
		if err != nil {
			return BadRequest(c, ctx, "Invalid author_id")
		}

		id, err := ParseIDParam(c, ctx)
		if err != nil {
			return BadRequest(c, ctx, "Invalid id")
		}

		var alias models.AuthorAlias
		if err := db.First(&alias, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return NotFound(c, ctx, "Alias", id)
			}
			return InternalError(c, ctx, "Failed to query alias", err)
		}

		// Verify alias belongs to this author
		if alias.AuthorID != authorID {
			return NotFound(c, ctx, "Alias", id)
		}

		if err := db.Delete(&alias).Error; err != nil {
			return InternalError(c, ctx, "Failed to delete alias", err)
		}

		slog.InfoContext(ctx, "Deleted author alias", slog.Uint64("id", uint64(id)))
		return c.NoContent(http.StatusNoContent)
	}
}
