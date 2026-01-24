package api

import (
	"context"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
)

type BookTypeResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type BookTypeHandler struct{}

func (h BookTypeHandler) ModelToResponse(c echo.Context, ctx context.Context, db *gorm.DB, row models.BookType) BookTypeResponse {
	return BookTypeResponse{
		ID:   row.ID,
		Name: row.Name,
	}
}

func (h BookTypeHandler) PreloadRelations(c echo.Context, ctx context.Context, db *gorm.DB) (*gorm.DB, error) {
	return db, nil
}

func (h BookTypeHandler) IDFromModel(row models.BookType) uint {
	return row.ID
}

// ListBookTypes returns all book types
func ListBookTypes(db *gorm.DB) echo.HandlerFunc {
	return readOnlyListHandler[models.BookType, BookTypeResponse](db, BookTypeHandler{})
}

// GetBookType returns a single book type by ID
func GetBookType(db *gorm.DB) echo.HandlerFunc {
	return readOnlyGetHandler[models.BookType, BookTypeResponse](db, BookTypeHandler{})
}
