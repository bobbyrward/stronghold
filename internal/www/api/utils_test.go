package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/models"
)

// TestParseIDParam tests the ParseIDParam utility function
func TestParseIDParam(t *testing.T) {
	tests := []struct {
		name      string
		paramID   string
		wantID    uint
		wantError bool
	}{
		{
			name:      "valid ID",
			paramID:   "123",
			wantID:    123,
			wantError: false,
		},
		{
			name:      "valid ID zero",
			paramID:   "0",
			wantID:    0,
			wantError: false,
		},
		{
			name:      "valid large ID",
			paramID:   "4294967295",
			wantID:    4294967295,
			wantError: false,
		},
		{
			name:      "invalid ID - not a number",
			paramID:   "abc",
			wantID:    0,
			wantError: true,
		},
		{
			name:      "invalid ID - negative",
			paramID:   "-1",
			wantID:    0,
			wantError: true,
		},
		{
			name:      "invalid ID - empty",
			paramID:   "",
			wantID:    0,
			wantError: true,
		},
		{
			name:      "invalid ID - decimal",
			paramID:   "1.5",
			wantID:    0,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/test/:id", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.paramID)
			ctx := c.Request().Context()

			id, err := ParseIDParam(c, ctx)

			if tt.wantError {
				assert.Error(t, err)
				assert.Equal(t, uint(0), id)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantID, id)
			}
		})
	}
}

// TestParseQueryParamUint tests the ParseQueryParamUint utility function
func TestParseQueryParamUint(t *testing.T) {
	tests := []struct {
		name      string
		query     string
		param     string
		wantID    uint
		wantHas   bool
		wantError bool
	}{
		{
			name:      "valid param present",
			query:     "feed_id=42",
			param:     "feed_id",
			wantID:    42,
			wantHas:   true,
			wantError: false,
		},
		{
			name:      "param not present",
			query:     "",
			param:     "feed_id",
			wantID:    0,
			wantHas:   false,
			wantError: false,
		},
		{
			name:      "different param present",
			query:     "other_id=42",
			param:     "feed_id",
			wantID:    0,
			wantHas:   false,
			wantError: false,
		},
		{
			name:      "invalid param value",
			query:     "feed_id=abc",
			param:     "feed_id",
			wantID:    0,
			wantHas:   false,
			wantError: true,
		},
		{
			name:      "zero value",
			query:     "feed_id=0",
			param:     "feed_id",
			wantID:    0,
			wantHas:   true,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			url := "/test"
			if tt.query != "" {
				url += "?" + tt.query
			}
			req := httptest.NewRequest(http.MethodGet, url, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			ctx := c.Request().Context()

			id, has, err := ParseQueryParamUint(c, ctx, tt.param)

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantID, id)
				assert.Equal(t, tt.wantHas, has)
			}
		})
	}
}

// TestBindRequest tests the BindRequest utility function
func TestBindRequest(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		contentType string
		wantError   bool
	}{
		{
			name:        "valid JSON",
			body:        `{"name": "test"}`,
			contentType: "application/json",
			wantError:   false,
		},
		{
			name:        "invalid JSON",
			body:        `{"name": }`,
			contentType: "application/json",
			wantError:   true,
		},
		{
			name:        "empty body",
			body:        "",
			contentType: "application/json",
			wantError:   false, // Empty JSON is valid for binding
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", tt.contentType)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			ctx := c.Request().Context()

			var dest struct {
				Name string `json:"name"`
			}

			err := BindRequest(c, ctx, &dest)

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestBadRequest tests the BadRequest utility function
func TestBadRequest(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctx := c.Request().Context()

	err := BadRequest(c, ctx, "Test error message")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Test error message")
}

// TestNotFound tests the NotFound utility function
func TestNotFound(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctx := c.Request().Context()

	err := NotFound(c, ctx, "Resource", 42)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "Resource not found")
}

// TestInternalError tests the InternalError utility function
func TestInternalError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	ctx := c.Request().Context()

	testErr := errors.New("database connection failed")
	err := InternalError(c, ctx, "Failed to query", testErr)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Failed to query")
}

// TestLookupByName tests the LookupByName utility function
func TestLookupByName(t *testing.T) {
	db, err := models.ConnectTestDB()
	require.NoError(t, err)

	ctx := context.Background()

	t.Run("finds existing record", func(t *testing.T) {
		var notificationType models.NotificationType
		err := LookupByName(db, ctx, &notificationType, "discord", "Notification type")
		assert.NoError(t, err)
		assert.Equal(t, "discord", notificationType.Name)
		assert.NotZero(t, notificationType.ID)
	})

	t.Run("returns error for non-existent record", func(t *testing.T) {
		var notificationType models.NotificationType
		err := LookupByName(db, ctx, &notificationType, "nonexistent", "Notification type")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Notification type not found")
	})

	t.Run("finds filter key", func(t *testing.T) {
		var filterKey models.FilterKey
		err := LookupByName(db, ctx, &filterKey, "author", "Filter key")
		assert.NoError(t, err)
		assert.Equal(t, "author", filterKey.Name)
	})

	t.Run("finds filter operator", func(t *testing.T) {
		var filterOperator models.FilterOperator
		err := LookupByName(db, ctx, &filterOperator, "contains", "Filter operator")
		assert.NoError(t, err)
		assert.Equal(t, "contains", filterOperator.Name)
	})
}

// TestGetByID tests the GetByID utility function
func TestGetByID(t *testing.T) {
	db, err := models.ConnectTestDB()
	require.NoError(t, err)

	ctx := context.Background()

	// Create a test record
	feed := models.Feed{
		Name: "Test Feed for GetByID",
		URL:  "https://example.com/getbyid",
	}
	result := db.Create(&feed)
	require.NoError(t, result.Error)
	require.NotZero(t, feed.ID)

	t.Run("finds existing record", func(t *testing.T) {
		var found models.Feed
		err := GetByID(db, ctx, &found, feed.ID, "feed")
		assert.NoError(t, err)
		assert.Equal(t, feed.ID, found.ID)
		assert.Equal(t, "Test Feed for GetByID", found.Name)
	})

	t.Run("returns ErrRecordNotFound for non-existent ID", func(t *testing.T) {
		var found models.Feed
		err := GetByID(db, ctx, &found, 99999, "feed")
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("returns ErrRecordNotFound for zero ID", func(t *testing.T) {
		var found models.Feed
		err := GetByID(db, ctx, &found, 0, "feed")
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

// TestDeleteByID tests the DeleteByID utility function
func TestDeleteByID(t *testing.T) {
	db, err := models.ConnectTestDB()
	require.NoError(t, err)

	ctx := context.Background()

	t.Run("deletes existing record", func(t *testing.T) {
		// Create a test record
		feed := models.Feed{
			Name: "Test Feed for Delete",
			URL:  "https://example.com/delete",
		}
		result := db.Create(&feed)
		require.NoError(t, result.Error)
		require.NotZero(t, feed.ID)

		// Delete it
		err := DeleteByID(db, ctx, &models.Feed{}, feed.ID, "feed")
		assert.NoError(t, err)

		// Verify deletion
		var found models.Feed
		result = db.First(&found, feed.ID)
		assert.Equal(t, gorm.ErrRecordNotFound, result.Error)
	})

	t.Run("returns ErrRecordNotFound for non-existent ID", func(t *testing.T) {
		err := DeleteByID(db, ctx, &models.Feed{}, 99999, "feed")
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

// TestErrorResponseConsistency tests that error responses are consistent
func TestErrorResponseConsistency(t *testing.T) {
	tests := []struct {
		name           string
		handler        func(c echo.Context, ctx context.Context) error
		expectedStatus int
		expectedKey    string
	}{
		{
			name: "BadRequest returns error key",
			handler: func(c echo.Context, ctx context.Context) error {
				return BadRequest(c, ctx, "bad request")
			},
			expectedStatus: http.StatusBadRequest,
			expectedKey:    "error",
		},
		{
			name: "NotFound returns error key",
			handler: func(c echo.Context, ctx context.Context) error {
				return NotFound(c, ctx, "Resource", 1)
			},
			expectedStatus: http.StatusNotFound,
			expectedKey:    "error",
		},
		{
			name: "InternalError returns error key",
			handler: func(c echo.Context, ctx context.Context) error {
				return InternalError(c, ctx, "internal error", errors.New("test"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedKey:    "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			ctx := c.Request().Context()

			err := tt.handler(c, ctx)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.Contains(t, rec.Body.String(), `"`+tt.expectedKey+`"`)
		})
	}
}

// TestLookupByNameWithDifferentModels tests LookupByName with various model types
func TestLookupByNameWithDifferentModels(t *testing.T) {
	db, err := models.ConnectTestDB()
	require.NoError(t, err)

	ctx := context.Background()

	// Test with seeded data
	tests := []struct {
		name         string
		model        interface{}
		lookupName   string
		resourceName string
		wantError    bool
	}{
		{
			name:         "NotificationType - discord",
			model:        &models.NotificationType{},
			lookupName:   "discord",
			resourceName: "Notification type",
			wantError:    false,
		},
		{
			name:         "FilterKey - author",
			model:        &models.FilterKey{},
			lookupName:   "author",
			resourceName: "Filter key",
			wantError:    false,
		},
		{
			name:         "FilterOperator - contains",
			model:        &models.FilterOperator{},
			lookupName:   "contains",
			resourceName: "Filter operator",
			wantError:    false,
		},
		{
			name:         "FeedFilterSetType - any",
			model:        &models.FeedFilterSetType{},
			lookupName:   "any",
			resourceName: "Feed filter set type",
			wantError:    false,
		},
		{
			name:         "Non-existent notification type",
			model:        &models.NotificationType{},
			lookupName:   "nonexistent",
			resourceName: "Notification type",
			wantError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := LookupByName(db, ctx, tt.model, tt.lookupName, tt.resourceName)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
