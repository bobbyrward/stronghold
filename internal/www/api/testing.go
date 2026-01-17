package api

import (
	"net/http"
	"testing"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/bobbyrward/stronghold/internal/hardcover"
	"github.com/bobbyrward/stronghold/internal/models"
)

// TestValidator is a validator for tests
type TestValidator struct {
	validator *validator.Validate
}

// NewTestValidator creates a new test validator
func NewTestValidator() *TestValidator {
	return &TestValidator{
		validator: validator.New(),
	}
}

// Validate validates the given struct
func (tv *TestValidator) Validate(i interface{}) error {
	if err := tv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

// SetupTestServer creates a test Echo server with a test database
func SetupTestServer(t *testing.T) (*echo.Echo, func()) {
	db, err := models.ConnectTestDB()
	require.NoError(t, err, "Failed to connect to test database")

	e := SetupTestServerWithDB(db)

	cleanup := func() {
		// No cleanup needed for in-memory database
	}

	return e, cleanup
}

// SetupTestServerWithDB creates a test Echo server with the provided database
func SetupTestServerWithDB(db *gorm.DB) *echo.Echo {
	e := echo.New()
	e.Validator = NewTestValidator()

	// Register all API routes under /api group (same as production)
	apiGroup := e.Group("/api")
	hc := hardcover.NewMockClient()

	// Add test authors to the mock Hardcover client
	hc.AddAuthor("brandon-sanderson", "Brandon Sanderson")
	hc.AddAuthor("brandon-mull", "Brandon Mull")
	hc.AddAuthor("patrick-rothfuss", "Patrick Rothfuss")
	hc.AddAuthor("joe-abercrombie", "Joe Abercrombie")

	RegisterRoutes(apiGroup, db, hc)

	return e
}
