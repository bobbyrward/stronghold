package api

import (
	"context"
	"log/slog"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// ReadOnlyModelHandler is a simplified interface for read-only resources
// that only need List and Get operations (no Create, Update, Delete).
type ReadOnlyModelHandler[Model any, Response any] interface {
	ModelToResponse(echo.Context, context.Context, *gorm.DB, Model) Response
	PreloadRelations(echo.Context, context.Context, *gorm.DB) (*gorm.DB, error)
	IDFromModel(Model) uint
}

func readOnlyListHandler[Model any, Response any](
	db *gorm.DB,
	handler ReadOnlyModelHandler[Model, Response],
) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		typeName := reflect.TypeFor[Model]().Name()

		slog.InfoContext(ctx, "Listing", slog.String("type", typeName))

		db, err := handler.PreloadRelations(c, ctx, db)
		if err != nil {
			return err
		}

		rows, err := GetAll[Model](c, ctx, db)
		if err != nil {
			return err
		}

		response := make([]Response, len(rows))
		for i, f := range rows {
			response[i] = handler.ModelToResponse(c, ctx, db, f)
		}

		slog.InfoContext(
			ctx,
			"Successfully listed",
			slog.Int("count", len(response)),
			slog.String("type", typeName),
		)

		return c.JSON(http.StatusOK, response)
	}
}

func readOnlyGetHandler[Model any, Response any](
	db *gorm.DB,
	handler ReadOnlyModelHandler[Model, Response],
) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		typeName := reflect.TypeFor[Model]().Name()

		id, err := ParseIDParam(c, ctx)
		if err != nil {
			return BadRequest(c, ctx, "Invalid ID")
		}

		slog.InfoContext(ctx, "Getting row", slog.String("type", typeName), slog.Uint64("id", uint64(id)))

		db, err := handler.PreloadRelations(c, ctx, db)
		if err != nil {
			return err
		}

		var row Model
		if err = GetByID(db, ctx, &row, id, typeName); err != nil {
			if err == gorm.ErrRecordNotFound {
				return NotFound(c, ctx, typeName, id)
			}
			return InternalError(c, ctx, "Failed to get row", err)
		}

		slog.InfoContext(ctx, "Successfully retrieved row", slog.Uint64("id", uint64(id)), slog.String("type", typeName), slog.Any("row", row))
		return c.JSON(http.StatusOK, handler.ModelToResponse(c, ctx, db, row))
	}
}
