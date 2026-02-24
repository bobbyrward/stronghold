package api

import (
	"context"
	"log/slog"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/eventlog"
)

// EventLogger is an optional interface that ModelHandlers can implement
// to log mutation events (create, update, delete) to the event log.
type EventLogger[Model any] interface {
	LogEvent(db *gorm.DB, eventType string, model Model)
}

// logEventIfSupported checks if the handler implements EventLogger and calls it.
func logEventIfSupported[Model any, Request any, Response any](handler ModelHandler[Model, Request, Response], db *gorm.DB, eventType string, model Model) {
	if logger, ok := any(handler).(EventLogger[Model]); ok {
		logger.LogEvent(db, eventType, model)
	}
}

type ModelHandler[Model any, Request any, Response any] interface {
	ModelToResponse(echo.Context, context.Context, *gorm.DB, Model) Response
	RequestToModel(echo.Context, context.Context, *gorm.DB, Request) (Model, error)
	UpdateModel(echo.Context, context.Context, *gorm.DB, *Model, Request) error
	ParseQuery(echo.Context, context.Context, *gorm.DB) (*gorm.DB, error)
	PreloadRelations(echo.Context, context.Context, *gorm.DB) (*gorm.DB, error)
	IDFromModel(Model) uint
}

func genericListHandler[Model any, Request any, Response any](
	db *gorm.DB,
	handler ModelHandler[Model, Request, Response],
) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		typeName := reflect.TypeFor[Model]().Name()

		slog.InfoContext(ctx, "Listing", slog.String("type", typeName))

		db, err := handler.PreloadRelations(c, ctx, db)
		if err != nil {
			return err
		}

		db, err = handler.ParseQuery(c, ctx, db)
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

func genericCreateHandler[Model any, Request any, Response any](
	db *gorm.DB,
	handler ModelHandler[Model, Request, Response],
) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		typeName := reflect.TypeFor[Model]().Name()

		var req Request
		if err := BindRequest(c, ctx, &req); err != nil {
			return BadRequest(c, ctx, "Invalid request body")
		}

		if err := ValidateRequest(c, ctx, &req); err != nil {
			return BadRequest(c, ctx, "Invalid request body")
		}

		slog.InfoContext(ctx, "Creating row", slog.Any("request", req), slog.String("type", typeName))

		row, err := handler.RequestToModel(c, ctx, db, req)
		if err != nil {
			return err
		}

		if err := db.Create(&row).Error; err != nil {
			return InternalError(c, ctx, "Failed to create row", err)
		}

		slog.InfoContext(ctx, "Successfully created row", slog.Any("row", row), slog.String("type", typeName))

		id := handler.IDFromModel(row)

		// Preload relations before fetching the created row
		dbWithPreload, err := handler.PreloadRelations(c, ctx, db)
		if err != nil {
			return err
		}

		if err = GetByID(dbWithPreload, ctx, &row, id, typeName); err != nil {
			if err == gorm.ErrRecordNotFound {
				return NotFound(c, ctx, typeName, id)
			}
			return InternalError(c, ctx, "Failed to get row", err)
		}

		slog.InfoContext(ctx, "Successfully retrieved row", slog.Uint64("id", uint64(id)), slog.String("type", typeName), slog.Any("row", row))

		logEventIfSupported(handler, db, eventlog.EventCreated, row)

		return c.JSON(http.StatusCreated, handler.ModelToResponse(c, ctx, db, row))
	}
}

func genericGetHandler[Model any, Request any, Response any](
	db *gorm.DB,
	handler ModelHandler[Model, Request, Response],
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

func genericUpdateHandler[Model any, Request any, Response any](
	db *gorm.DB,
	handler ModelHandler[Model, Request, Response],
) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		typeName := reflect.TypeFor[Model]().Name()

		id, err := ParseIDParam(c, ctx)
		if err != nil {
			return BadRequest(c, ctx, "Invalid ID")
		}

		var req Request
		if err := BindRequest(c, ctx, &req); err != nil {
			return BadRequest(c, ctx, "Invalid request body")
		}

		if err := ValidateRequest(c, ctx, &req); err != nil {
			return BadRequest(c, ctx, "Invalid request body")
		}

		slog.InfoContext(ctx, "Updating row", slog.Uint64("id", uint64(id)), slog.String("type", typeName))

		var row Model
		if err := GetByID(db, ctx, &row, id, typeName); err != nil {
			if err == gorm.ErrRecordNotFound {
				return NotFound(c, ctx, typeName, id)
			}
			return InternalError(c, ctx, "Failed to query", err)
		}

		err = handler.UpdateModel(c, ctx, db, &row, req)
		if err != nil {
			return err
		}

		if err := db.Save(&row).Error; err != nil {
			return InternalError(c, ctx, "Failed to update feed", err)
		}

		// Preload relations before fetching the updated row
		dbWithPreload, err := handler.PreloadRelations(c, ctx, db)
		if err != nil {
			return err
		}

		if err := GetByID(dbWithPreload, ctx, &row, id, typeName); err != nil {
			if err == gorm.ErrRecordNotFound {
				return NotFound(c, ctx, typeName, id)
			}
			return InternalError(c, ctx, "Failed to query", err)
		}

		slog.InfoContext(ctx, "Successfully updated row", slog.Uint64("id", uint64(id)), slog.String("type", typeName))

		logEventIfSupported(handler, db, eventlog.EventUpdated, row)

		return c.JSON(http.StatusOK, handler.ModelToResponse(c, ctx, db, row))
	}
}

// genericDeleteHandler deletes a resource by ID
func genericDeleteHandler[Model any](db *gorm.DB) echo.HandlerFunc {
	return genericDeleteHandlerImpl[Model, any, any](db, nil)
}

// genericDeleteHandlerWithEventLog deletes a resource by ID and logs a mutation event
// via the handler's LogEvent method, which has access to the model's name and other fields.
func genericDeleteHandlerWithEventLog[Model any, Request any, Response any](
	db *gorm.DB,
	handler ModelHandler[Model, Request, Response],
) echo.HandlerFunc {
	return genericDeleteHandlerImpl[Model, Request, Response](db, handler)
}

func genericDeleteHandlerImpl[Model any, Request any, Response any](
	db *gorm.DB,
	handler ModelHandler[Model, Request, Response],
) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		typeName := reflect.TypeFor[Model]().Name()

		id, err := ParseIDParam(c, ctx)
		if err != nil {
			return BadRequest(c, ctx, "Invalid ID")
		}

		slog.InfoContext(ctx, "Deleting row", slog.Uint64("id", uint64(id)))

		// Fetch the row before deleting so LogEvent has access to model fields (name, etc.)
		var row Model
		if handler != nil {
			if err := GetByID(db, ctx, &row, id, typeName); err != nil {
				if err == gorm.ErrRecordNotFound {
					return NotFound(c, ctx, typeName, id)
				}
				return InternalError(c, ctx, "Failed to query", err)
			}
		}

		if err := DeleteByID(db, ctx, new(Model), id, typeName); err != nil {
			if err == gorm.ErrRecordNotFound {
				return NotFound(c, ctx, typeName, id)
			}
			return InternalError(c, ctx, "Failed to delete", err)
		}

		if handler != nil {
			logEventIfSupported(handler, db, eventlog.EventDeleted, row)
		}

		return c.NoContent(http.StatusNoContent)
	}
}
