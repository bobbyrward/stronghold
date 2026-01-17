package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper to create a test author for subscription tests
func createTestAuthorForSubscription(t *testing.T, e *echo.Echo, name string) AuthorResponse {
	body := fmt.Sprintf(`{"name": "%s"}`, name)
	req := httptest.NewRequest(http.MethodPost, "/api/authors", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	var author AuthorResponse
	err := json.Unmarshal(rec.Body.Bytes(), &author)
	require.NoError(t, err)
	return author
}

func TestAuthorSubscription_CRUD(t *testing.T) {
	e, cleanup := SetupTestServer(t)
	defer cleanup()

	author := createTestAuthorForSubscription(t, e, "Subscription Test Author")

	t.Run("Get subscription returns 404 when none exists", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/authors/%d/subscription", author.ID), nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("Create subscription", func(t *testing.T) {
		body := `{"scope_name": "personal"}`
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/authors/%d/subscription", author.ID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var sub AuthorSubscriptionResponse
		err := json.Unmarshal(rec.Body.Bytes(), &sub)
		require.NoError(t, err)

		assert.NotZero(t, sub.ID)
		assert.Equal(t, author.ID, sub.AuthorID)
		assert.Equal(t, "personal", sub.ScopeName)
	})

	t.Run("Get subscription returns data", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/authors/%d/subscription", author.ID), nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var sub AuthorSubscriptionResponse
		err := json.Unmarshal(rec.Body.Bytes(), &sub)
		require.NoError(t, err)
		assert.Equal(t, "personal", sub.ScopeName)
	})

	t.Run("Update subscription", func(t *testing.T) {
		body := `{"scope_name": "family"}`
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/authors/%d/subscription", author.ID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var sub AuthorSubscriptionResponse
		err := json.Unmarshal(rec.Body.Bytes(), &sub)
		require.NoError(t, err)
		assert.Equal(t, "family", sub.ScopeName)
	})

	t.Run("Delete subscription", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/authors/%d/subscription", author.ID), nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("Get subscription returns 404 after delete", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/authors/%d/subscription", author.ID), nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestAuthorSubscription_Conflict(t *testing.T) {
	e, cleanup := SetupTestServer(t)
	defer cleanup()

	author := createTestAuthorForSubscription(t, e, "Conflict Test Author")

	// Create first subscription
	body := `{"scope_name": "personal"}`
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/authors/%d/subscription", author.ID), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	t.Run("Creating duplicate subscription returns 409", func(t *testing.T) {
		body := `{"scope_name": "family"}`
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/authors/%d/subscription", author.ID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
	})
}

func TestAuthorSubscription_InvalidScope(t *testing.T) {
	e, cleanup := SetupTestServer(t)
	defer cleanup()

	author := createTestAuthorForSubscription(t, e, "Invalid Scope Author")

	t.Run("Create with invalid scope returns 400", func(t *testing.T) {
		body := `{"scope_name": "invalid_scope"}`
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/authors/%d/subscription", author.ID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestAuthorSubscription_InvalidAuthor(t *testing.T) {
	e, cleanup := SetupTestServer(t)
	defer cleanup()

	t.Run("Create subscription for non-existent author returns 404", func(t *testing.T) {
		body := `{"scope_name": "personal"}`
		req := httptest.NewRequest(http.MethodPost, "/api/authors/99999/subscription", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("Get subscription for non-existent author returns 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/authors/99999/subscription", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("Update subscription for non-existent author returns 404", func(t *testing.T) {
		body := `{"scope_name": "family"}`
		req := httptest.NewRequest(http.MethodPut, "/api/authors/99999/subscription", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("Delete subscription for non-existent author returns 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/authors/99999/subscription", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestAuthorSubscription_InvalidAuthorID(t *testing.T) {
	e, cleanup := SetupTestServer(t)
	defer cleanup()

	t.Run("Invalid author_id returns 400", func(t *testing.T) {
		body := `{"scope_name": "personal"}`
		req := httptest.NewRequest(http.MethodPost, "/api/authors/invalid/subscription", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestAuthorSubscriptionItems_List(t *testing.T) {
	e, cleanup := SetupTestServer(t)
	defer cleanup()

	author := createTestAuthorForSubscription(t, e, "Items Test Author")

	t.Run("List items returns 404 when no subscription", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/authors/%d/subscription/items", author.ID), nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	// Create subscription
	body := `{"scope_name": "personal"}`
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/authors/%d/subscription", author.ID), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	t.Run("List items returns empty array when subscription exists", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/authors/%d/subscription/items", author.ID), nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var items []AuthorSubscriptionItemResponse
		err := json.Unmarshal(rec.Body.Bytes(), &items)
		require.NoError(t, err)
		assert.Empty(t, items)
	})
}

func TestAuthorSubscriptionItems_InvalidAuthor(t *testing.T) {
	e, cleanup := SetupTestServer(t)
	defer cleanup()

	t.Run("List items for non-existent author returns 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/authors/99999/subscription/items", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("Invalid author_id returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/authors/invalid/subscription/items", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestAuthorSubscription_UpdateNonExistent(t *testing.T) {
	e, cleanup := SetupTestServer(t)
	defer cleanup()

	author := createTestAuthorForSubscription(t, e, "Update Non-Existent Author")

	t.Run("Update non-existent subscription returns 404", func(t *testing.T) {
		body := `{"scope_name": "family"}`
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/authors/%d/subscription", author.ID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestAuthorSubscription_DeleteNonExistent(t *testing.T) {
	e, cleanup := SetupTestServer(t)
	defer cleanup()

	author := createTestAuthorForSubscription(t, e, "Delete Non-Existent Author")

	t.Run("Delete non-existent subscription returns 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/authors/%d/subscription", author.ID), nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestAuthorSubscription_WithNotifier(t *testing.T) {
	e, cleanup := SetupTestServer(t)
	defer cleanup()

	author := createTestAuthorForSubscription(t, e, "Notifier Test Author")

	// Create a notifier
	discordTypeID := getNotificationTypeID(t, e, "discord")
	notifierReq := NotifierRequest{
		Name:   "subscription-test-notifier",
		TypeID: discordTypeID,
		URL:    "https://discord.com/webhook/test",
	}
	notifierBody, _ := json.Marshal(notifierReq)
	req := httptest.NewRequest(http.MethodPost, "/api/notifiers", bytes.NewReader(notifierBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	var createdNotifier NotifierResponse
	err := json.Unmarshal(rec.Body.Bytes(), &createdNotifier)
	require.NoError(t, err)

	t.Run("Create subscription with notifier", func(t *testing.T) {
		body := fmt.Sprintf(`{"scope_name": "personal", "notifier_id": %d}`, createdNotifier.ID)
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/authors/%d/subscription", author.ID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var sub AuthorSubscriptionResponse
		err := json.Unmarshal(rec.Body.Bytes(), &sub)
		require.NoError(t, err)

		assert.NotZero(t, sub.ID)
		assert.Equal(t, author.ID, sub.AuthorID)
		assert.Equal(t, "personal", sub.ScopeName)
		require.NotNil(t, sub.NotifierID)
		assert.Equal(t, createdNotifier.ID, *sub.NotifierID)
		require.NotNil(t, sub.NotifierName)
		assert.Equal(t, "subscription-test-notifier", *sub.NotifierName)
	})

	t.Run("Get subscription with notifier", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/authors/%d/subscription", author.ID), nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var sub AuthorSubscriptionResponse
		err := json.Unmarshal(rec.Body.Bytes(), &sub)
		require.NoError(t, err)
		assert.Equal(t, "personal", sub.ScopeName)
		require.NotNil(t, sub.NotifierName)
		assert.Equal(t, "subscription-test-notifier", *sub.NotifierName)
	})
}

func TestAuthorSubscription_UpdateScope(t *testing.T) {
	e, cleanup := SetupTestServer(t)
	defer cleanup()

	author := createTestAuthorForSubscription(t, e, "Update Scope Author")

	// Create subscription
	body := `{"scope_name": "personal"}`
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/authors/%d/subscription", author.ID), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	t.Run("Update with invalid scope returns 400", func(t *testing.T) {
		body := `{"scope_name": "invalid_scope"}`
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/authors/%d/subscription", author.ID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Update to different valid scope succeeds", func(t *testing.T) {
		body := `{"scope_name": "kids"}`
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/authors/%d/subscription", author.ID), bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var sub AuthorSubscriptionResponse
		err := json.Unmarshal(rec.Body.Bytes(), &sub)
		require.NoError(t, err)
		assert.Equal(t, "kids", sub.ScopeName)
	})
}
