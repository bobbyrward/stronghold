package api

import (
	"context"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
)

type LibraryRequest struct {
	Name         string `json:"name" validate:"required"`
	Path         string `json:"path" validate:"required"`
	BookTypeName string `json:"book_type_name" validate:"required"`
}

type LibraryResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Path         string `json:"path"`
	BookTypeID   uint   `json:"book_type_id"`
	BookTypeName string `json:"book_type_name"`
}

type LibraryHandler struct{}

func (h LibraryHandler) ModelToResponse(c echo.Context, ctx context.Context, db *gorm.DB, row models.Library) LibraryResponse {
	return LibraryResponse{
		ID:           row.ID,
		Name:         row.Name,
		Path:         row.Path,
		BookTypeID:   row.BookTypeID,
		BookTypeName: row.BookType.Name,
	}
}

func (h LibraryHandler) RequestToModel(c echo.Context, ctx context.Context, db *gorm.DB, req LibraryRequest) (models.Library, error) {
	var bookType models.BookType
	if err := LookupByName(db, ctx, &bookType, req.BookTypeName, "Book type"); err != nil {
		return models.Library{}, BadRequest(c, ctx, "Invalid book_type_name: "+req.BookTypeName)
	}

	return models.Library{
		Name:       req.Name,
		Path:       req.Path,
		BookTypeID: bookType.ID,
	}, nil
}

func (h LibraryHandler) UpdateModel(c echo.Context, ctx context.Context, db *gorm.DB, row *models.Library, req LibraryRequest) error {
	var bookType models.BookType
	if err := LookupByName(db, ctx, &bookType, req.BookTypeName, "Book type"); err != nil {
		return BadRequest(c, ctx, "Invalid book_type_name: "+req.BookTypeName)
	}

	row.Name = req.Name
	row.Path = req.Path
	row.BookTypeID = bookType.ID
	return nil
}

func (h LibraryHandler) ParseQuery(c echo.Context, ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	// Optional: filter by book_type_id query param
	db, err := ApplyUintFilter(c, ctx, db, "book_type_id", "book_type_id")
	if err != nil {
		return db, err
	}
	return db, nil
}

func (h LibraryHandler) PreloadRelations(c echo.Context, ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	return db.Preload("BookType"), nil
}

func (h LibraryHandler) IDFromModel(row models.Library) uint {
	return row.ID
}

// ListLibraries returns all libraries
func ListLibraries(db *gorm.DB) echo.HandlerFunc {
	return genericListHandler[models.Library, LibraryRequest, LibraryResponse](db, LibraryHandler{})
}

// CreateLibrary creates a new library
func CreateLibrary(db *gorm.DB) echo.HandlerFunc {
	return genericCreateHandler[models.Library, LibraryRequest, LibraryResponse](db, LibraryHandler{})
}

// GetLibrary returns a single library by ID
func GetLibrary(db *gorm.DB) echo.HandlerFunc {
	return genericGetHandler[models.Library, LibraryRequest, LibraryResponse](db, LibraryHandler{})
}

// UpdateLibrary updates an existing library
func UpdateLibrary(db *gorm.DB) echo.HandlerFunc {
	return genericUpdateHandler[models.Library, LibraryRequest, LibraryResponse](db, LibraryHandler{})
}

// DeleteLibrary deletes a library
func DeleteLibrary(db *gorm.DB) echo.HandlerFunc {
	return genericDeleteHandler[models.Library](db)
}
