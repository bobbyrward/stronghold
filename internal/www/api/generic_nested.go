package api

import (
	"context"
	"log/slog"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// NestedModelHandler handles child resources under a parent resource.
// It provides CRUD operations for resources that belong to a parent (e.g., AuthorAlias under Author).
type NestedModelHandler[Parent any, Model any, Request any, Response any] interface {
	// ParseParentID extracts the parent ID from the request path
	ParseParentID(echo.Context, context.Context) (uint, error)

	// ValidateParent validates that the parent exists, returns error if not found
	ValidateParent(*gorm.DB, context.Context, uint) error

	// GetParentID returns the parent ID from a model instance
	GetParentID(Model) uint

	// SetParentID sets the parent ID on a model instance
	SetParentID(*Model, uint)

	// ModelToResponse converts a model to its API response type
	ModelToResponse(echo.Context, context.Context, *gorm.DB, Model) Response

	// RequestToModel converts a request to a new model instance
	RequestToModel(echo.Context, context.Context, *gorm.DB, Request) (Model, error)

	// UpdateModel updates an existing model with request data
	UpdateModel(echo.Context, context.Context, *gorm.DB, *Model, Request) error

	// PreloadRelations adds any preloads needed for the model
	PreloadRelations(echo.Context, context.Context, *gorm.DB) (*gorm.DB, error)

	// IDFromModel returns the ID of a model instance
	IDFromModel(Model) uint

	// ParentForeignKey returns the column name for filtering by parent (e.g., "author_id")
	ParentForeignKey() string
}

func nestedListHandler[Parent any, Model any, Request any, Response any](
	db *gorm.DB,
	handler NestedModelHandler[Parent, Model, Request, Response],
) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		typeName := reflect.TypeFor[Model]().Name()

		parentID, err := handler.ParseParentID(c, ctx)
		if err != nil {
			return BadRequest(c, ctx, "Invalid parent ID")
		}

		if err := handler.ValidateParent(db, ctx, parentID); err != nil {
			if err == gorm.ErrRecordNotFound {
				parentTypeName := reflect.TypeFor[Parent]().Name()
				return NotFound(c, ctx, parentTypeName, parentID)
			}
			return InternalError(c, ctx, "Failed to query parent", err)
		}

		slog.InfoContext(ctx, "Listing nested resources",
			slog.String("type", typeName),
			slog.Uint64("parent_id", uint64(parentID)))

		dbWithPreload, err := handler.PreloadRelations(c, ctx, db)
		if err != nil {
			return err
		}

		var rows []Model
		if err := dbWithPreload.Where(handler.ParentForeignKey()+" = ?", parentID).Find(&rows).Error; err != nil {
			return InternalError(c, ctx, "Failed to list resources", err)
		}

		response := make([]Response, len(rows))
		for i, row := range rows {
			response[i] = handler.ModelToResponse(c, ctx, db, row)
		}

		slog.InfoContext(ctx, "Successfully listed nested resources",
			slog.Int("count", len(response)),
			slog.String("type", typeName))

		return c.JSON(http.StatusOK, response)
	}
}

func nestedCreateHandler[Parent any, Model any, Request any, Response any](
	db *gorm.DB,
	handler NestedModelHandler[Parent, Model, Request, Response],
) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		typeName := reflect.TypeFor[Model]().Name()

		parentID, err := handler.ParseParentID(c, ctx)
		if err != nil {
			return BadRequest(c, ctx, "Invalid parent ID")
		}

		if err := handler.ValidateParent(db, ctx, parentID); err != nil {
			if err == gorm.ErrRecordNotFound {
				parentTypeName := reflect.TypeFor[Parent]().Name()
				return NotFound(c, ctx, parentTypeName, parentID)
			}
			return InternalError(c, ctx, "Failed to query parent", err)
		}

		var req Request
		if err := BindRequest(c, ctx, &req); err != nil {
			return BadRequest(c, ctx, "Invalid request body")
		}
		if err := ValidateRequest(c, ctx, &req); err != nil {
			return BadRequest(c, ctx, "Validation failed")
		}

		slog.InfoContext(ctx, "Creating nested resource",
			slog.String("type", typeName),
			slog.Uint64("parent_id", uint64(parentID)),
			slog.Any("request", req))

		row, err := handler.RequestToModel(c, ctx, db, req)
		if err != nil {
			return err
		}

		handler.SetParentID(&row, parentID)

		if err := db.Create(&row).Error; err != nil {
			return InternalError(c, ctx, "Failed to create resource", err)
		}

		id := handler.IDFromModel(row)
		slog.InfoContext(ctx, "Successfully created nested resource",
			slog.Uint64("id", uint64(id)),
			slog.Uint64("parent_id", uint64(parentID)),
			slog.String("type", typeName))

		return c.JSON(http.StatusCreated, handler.ModelToResponse(c, ctx, db, row))
	}
}

func nestedGetHandler[Parent any, Model any, Request any, Response any](
	db *gorm.DB,
	handler NestedModelHandler[Parent, Model, Request, Response],
) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		typeName := reflect.TypeFor[Model]().Name()

		parentID, err := handler.ParseParentID(c, ctx)
		if err != nil {
			return BadRequest(c, ctx, "Invalid parent ID")
		}

		id, err := ParseIDParam(c, ctx)
		if err != nil {
			return BadRequest(c, ctx, "Invalid ID")
		}

		slog.InfoContext(ctx, "Getting nested resource",
			slog.String("type", typeName),
			slog.Uint64("id", uint64(id)),
			slog.Uint64("parent_id", uint64(parentID)))

		dbWithPreload, err := handler.PreloadRelations(c, ctx, db)
		if err != nil {
			return err
		}

		var row Model
		if err := dbWithPreload.First(&row, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return NotFound(c, ctx, typeName, id)
			}
			return InternalError(c, ctx, "Failed to query resource", err)
		}

		// Verify resource belongs to this parent
		if handler.GetParentID(row) != parentID {
			return NotFound(c, ctx, typeName, id)
		}

		return c.JSON(http.StatusOK, handler.ModelToResponse(c, ctx, db, row))
	}
}

func nestedUpdateHandler[Parent any, Model any, Request any, Response any](
	db *gorm.DB,
	handler NestedModelHandler[Parent, Model, Request, Response],
) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		typeName := reflect.TypeFor[Model]().Name()

		parentID, err := handler.ParseParentID(c, ctx)
		if err != nil {
			return BadRequest(c, ctx, "Invalid parent ID")
		}

		id, err := ParseIDParam(c, ctx)
		if err != nil {
			return BadRequest(c, ctx, "Invalid ID")
		}

		var row Model
		if err := db.First(&row, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return NotFound(c, ctx, typeName, id)
			}
			return InternalError(c, ctx, "Failed to query resource", err)
		}

		// Verify resource belongs to this parent
		if handler.GetParentID(row) != parentID {
			return NotFound(c, ctx, typeName, id)
		}

		var req Request
		if err := BindRequest(c, ctx, &req); err != nil {
			return BadRequest(c, ctx, "Invalid request body")
		}
		if err := ValidateRequest(c, ctx, &req); err != nil {
			return BadRequest(c, ctx, "Validation failed")
		}

		slog.InfoContext(ctx, "Updating nested resource",
			slog.String("type", typeName),
			slog.Uint64("id", uint64(id)),
			slog.Uint64("parent_id", uint64(parentID)))

		if err := handler.UpdateModel(c, ctx, db, &row, req); err != nil {
			return err
		}

		if err := db.Save(&row).Error; err != nil {
			return InternalError(c, ctx, "Failed to update resource", err)
		}

		slog.InfoContext(ctx, "Successfully updated nested resource",
			slog.Uint64("id", uint64(id)),
			slog.String("type", typeName))

		return c.JSON(http.StatusOK, handler.ModelToResponse(c, ctx, db, row))
	}
}

func nestedDeleteHandler[Parent any, Model any, Request any, Response any](
	db *gorm.DB,
	handler NestedModelHandler[Parent, Model, Request, Response],
) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		typeName := reflect.TypeFor[Model]().Name()

		parentID, err := handler.ParseParentID(c, ctx)
		if err != nil {
			return BadRequest(c, ctx, "Invalid parent ID")
		}

		id, err := ParseIDParam(c, ctx)
		if err != nil {
			return BadRequest(c, ctx, "Invalid ID")
		}

		var row Model
		if err := db.First(&row, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return NotFound(c, ctx, typeName, id)
			}
			return InternalError(c, ctx, "Failed to query resource", err)
		}

		// Verify resource belongs to this parent
		if handler.GetParentID(row) != parentID {
			return NotFound(c, ctx, typeName, id)
		}

		slog.InfoContext(ctx, "Deleting nested resource",
			slog.String("type", typeName),
			slog.Uint64("id", uint64(id)),
			slog.Uint64("parent_id", uint64(parentID)))

		if err := db.Delete(&row).Error; err != nil {
			return InternalError(c, ctx, "Failed to delete resource", err)
		}

		slog.InfoContext(ctx, "Successfully deleted nested resource",
			slog.Uint64("id", uint64(id)),
			slog.String("type", typeName))

		return c.NoContent(http.StatusNoContent)
	}
}
