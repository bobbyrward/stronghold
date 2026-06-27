package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// getSubscriptionScopeID is a helper function to get a scope ID by name
func getSubscriptionScopeID(t *testing.T, e *echo.Echo, name string) uint {
	req := httptest.NewRequest(http.MethodGet, "/api/subscription-scopes", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	var scopes []SubscriptionScopeResponse
	err := json.Unmarshal(rec.Body.Bytes(), &scopes)
	require.NoError(t, err)

	for _, s := range scopes {
		if s.Name == name {
			return s.ID
		}
	}
	t.Fatalf("Subscription scope %q not found", name)
	return 0
}

func TestSubscriptionScopes(t *testing.T) {
	e, cleanup := SetupTestServer(t)
	defer cleanup()

	t.Run("List returns all scopes", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/subscription-scopes", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var scopes []SubscriptionScopeResponse
		err := json.Unmarshal(rec.Body.Bytes(), &scopes)
		require.NoError(t, err)

		// Should have 4 scopes: personal, family, kids, general
		assert.Len(t, scopes, 4)

		// Verify expected scope names exist
		names := make(map[string]bool)
		for _, s := range scopes {
			names[s.Name] = true
		}
		assert.True(t, names["personal"], "should have 'personal' scope")
		assert.True(t, names["family"], "should have 'family' scope")
		assert.True(t, names["kids"], "should have 'kids' scope")
		assert.True(t, names["general"], "should have 'general' scope")
	})

	t.Run("Get by ID returns scope", func(t *testing.T) {
		// First get a valid ID
		scopeID := getSubscriptionScopeID(t, e, "personal")

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/subscription-scopes/%d", scopeID), nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var scope SubscriptionScopeResponse
		err := json.Unmarshal(rec.Body.Bytes(), &scope)
		require.NoError(t, err)

		assert.Equal(t, scopeID, scope.ID)
		assert.Equal(t, "personal", scope.Name)
	})

	t.Run("Get non-existent ID returns 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/subscription-scopes/99999", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("Get invalid ID returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/subscription-scopes/invalid", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
