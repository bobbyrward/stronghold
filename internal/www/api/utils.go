package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// ParseIDParam extracts and validates the "id" path parameter
func ParseIDParam(c echo.Context, ctx context.Context) (uint, error) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		slog.WarnContext(ctx, "Invalid ID parameter", slog.String("id", idStr), slog.Any("error", err))
		return 0, err
	}
	return uint(id), nil
}

// ParseAuthorIDParam extracts and validates the "author_id" path parameter
func ParseAuthorIDParam(c echo.Context, ctx context.Context) (uint, error) {
	idStr := c.Param("author_id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		slog.WarnContext(ctx, "Invalid author_id parameter", slog.String("author_id", idStr), slog.Any("error", err))
		return 0, err
	}
	return uint(id), nil
}

// ParseQueryParamUint parses an optional query parameter as uint
func ParseQueryParamUint(c echo.Context, ctx context.Context, param string) (uint, bool, error) {
	str := c.QueryParam(param)
	if str == "" {
		return 0, false, nil
	}
	id, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		slog.WarnContext(ctx, "Invalid "+param+" parameter", slog.String(param, str), slog.Any("error", err))
		return 0, false, err
	}
	return uint(id), true, nil
}

// BindRequest binds and logs request body errors
func BindRequest(c echo.Context, ctx context.Context, req interface{}) error {
	if err := c.Bind(req); err != nil {
		slog.WarnContext(ctx, "Invalid request body", slog.Any("error", err))
		return err
	}
	return nil
}

// ValidateRequest validated and logs request body errors
func ValidateRequest(c echo.Context, ctx context.Context, req interface{}) error {
	if err := c.Validate(req); err != nil {
		slog.WarnContext(ctx, "Invalid request body", slog.Any("error", err))
		return err
	}
	return nil
}

// BadRequest returns a 400 response with logging
func BadRequest(c echo.Context, ctx context.Context, message string) error {
	slog.WarnContext(ctx, message)
	return c.JSON(http.StatusBadRequest, map[string]string{"error": message})
}

// NotFound returns a 404 response with logging
func NotFound(c echo.Context, ctx context.Context, resource string, id uint) error {
	slog.WarnContext(ctx, resource+" not found", slog.Uint64("id", uint64(id)))
	return c.JSON(http.StatusNotFound, map[string]string{"error": resource + " not found"})
}

// NotFound returns a 404 response with logging
func GenericNotFound(c echo.Context, ctx context.Context, message string) error {
	slog.WarnContext(ctx, "not found", slog.String("msg", message))
	return c.JSON(http.StatusNotFound, map[string]string{"error": message})
}

// InternalError returns a 500 response with logging
func InternalError(c echo.Context, ctx context.Context, message string, err error) error {
	slog.ErrorContext(ctx, message, slog.Any("error", err))
	return c.JSON(http.StatusInternalServerError, map[string]string{"error": message})
}

// LookupByName finds a record by name field, returns user-friendly error
func LookupByName(db *gorm.DB, ctx context.Context, dest interface{}, name, resourceName string) error {
	result := db.Where("name = ?", name).First(dest)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			slog.WarnContext(ctx, resourceName+" not found", slog.String("name", name))
			return fmt.Errorf("%s not found", resourceName)
		}
		slog.ErrorContext(ctx, "Failed to query "+resourceName, slog.String("name", name), slog.Any("error", result.Error))
		return result.Error
	}
	return nil
}

// GetAll fetches all records of type T with standard error handling
func GetAll[T any](c echo.Context, ctx context.Context, db *gorm.DB) ([]T, error) {
	var rows []T
	if err := db.Find(&rows).Error; err != nil {
		return nil, InternalError(c, ctx, "Failed to GetAll", err)
	}

	return rows, nil
}

// GetByID fetches a record by ID with standard error handling
func GetByID(db *gorm.DB, ctx context.Context, dest interface{}, id uint, resourceName string) error {
	result := db.First(dest, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return gorm.ErrRecordNotFound
		}
		slog.ErrorContext(ctx, "Failed to query "+resourceName, slog.Uint64("id", uint64(id)), slog.Any("error", result.Error))
		return result.Error
	}
	return nil
}

// DeleteByID deletes a record and handles not-found
func DeleteByID(db *gorm.DB, ctx context.Context, model interface{}, id uint, resourceName string) error {
	result := db.Unscoped().Delete(model, id)
	if result.Error != nil {
		slog.ErrorContext(ctx, "Failed to delete "+resourceName, slog.Uint64("id", uint64(id)), slog.Any("error", result.Error))
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	slog.InfoContext(ctx, "Successfully deleted "+resourceName, slog.Uint64("id", uint64(id)))
	return nil
}

// EscapeLikePattern escapes special characters in LIKE patterns to prevent pattern injection
func EscapeLikePattern(pattern string) string {
	pattern = strings.ReplaceAll(pattern, "\\", "\\\\")
	pattern = strings.ReplaceAll(pattern, "%", "\\%")
	pattern = strings.ReplaceAll(pattern, "_", "\\_")
	return pattern
}
