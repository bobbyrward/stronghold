package api

import (
	"context"

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

type AuthorAliasHandler struct{}

func (h AuthorAliasHandler) ParseParentID(c echo.Context, ctx context.Context) (uint, error) {
	return ParseAuthorIDParam(c, ctx)
}

func (h AuthorAliasHandler) ValidateParent(db *gorm.DB, ctx context.Context, id uint) error {
	var author models.Author
	return db.First(&author, id).Error
}

func (h AuthorAliasHandler) GetParentID(row models.AuthorAlias) uint {
	return row.AuthorID
}

func (h AuthorAliasHandler) SetParentID(row *models.AuthorAlias, parentID uint) {
	row.AuthorID = parentID
}

func (h AuthorAliasHandler) ModelToResponse(c echo.Context, ctx context.Context, db *gorm.DB, row models.AuthorAlias) AuthorAliasResponse {
	return AuthorAliasResponse{
		ID:       row.ID,
		AuthorID: row.AuthorID,
		Name:     row.Name,
	}
}

func (h AuthorAliasHandler) RequestToModel(c echo.Context, ctx context.Context, db *gorm.DB, req AuthorAliasRequest) (models.AuthorAlias, error) {
	return models.AuthorAlias{
		Name: req.Name,
	}, nil
}

func (h AuthorAliasHandler) UpdateModel(c echo.Context, ctx context.Context, db *gorm.DB, row *models.AuthorAlias, req AuthorAliasRequest) error {
	row.Name = req.Name
	return nil
}

func (h AuthorAliasHandler) PreloadRelations(c echo.Context, ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	return db, nil
}

func (h AuthorAliasHandler) IDFromModel(row models.AuthorAlias) uint {
	return row.ID
}

func (h AuthorAliasHandler) ParentForeignKey() string {
	return "author_id"
}

// ListAuthorAliases returns all aliases for a specific author
func ListAuthorAliases(db *gorm.DB) echo.HandlerFunc {
	return nestedListHandler[models.Author, models.AuthorAlias, AuthorAliasRequest, AuthorAliasResponse](db, AuthorAliasHandler{})
}

// CreateAuthorAlias creates a new alias for an author
func CreateAuthorAlias(db *gorm.DB) echo.HandlerFunc {
	return nestedCreateHandler[models.Author, models.AuthorAlias, AuthorAliasRequest, AuthorAliasResponse](db, AuthorAliasHandler{})
}

// GetAuthorAlias returns a single alias by ID
func GetAuthorAlias(db *gorm.DB) echo.HandlerFunc {
	return nestedGetHandler[models.Author, models.AuthorAlias, AuthorAliasRequest, AuthorAliasResponse](db, AuthorAliasHandler{})
}

// UpdateAuthorAlias updates an existing alias
func UpdateAuthorAlias(db *gorm.DB) echo.HandlerFunc {
	return nestedUpdateHandler[models.Author, models.AuthorAlias, AuthorAliasRequest, AuthorAliasResponse](db, AuthorAliasHandler{})
}

// DeleteAuthorAlias deletes an alias
func DeleteAuthorAlias(db *gorm.DB) echo.HandlerFunc {
	return nestedDeleteHandler[models.Author, models.AuthorAlias, AuthorAliasRequest, AuthorAliasResponse](db, AuthorAliasHandler{})
}
