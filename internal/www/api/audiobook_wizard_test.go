package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bobbyrward/stronghold/internal/importers/audiobooks/metadata"
	"github.com/stretchr/testify/assert"
)

// TestGetLibraries tests the GET /audiobook-wizard/libraries endpoint
func TestGetLibraries(t *testing.T) {
	e, cleanup := setupTestServer(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/api/audiobook-wizard/libraries", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// The endpoint should exist and return 200
	// The libraries list may be empty in test environment
	assert.Equal(t, http.StatusOK, rec.Code)

	// Should return valid JSON array
	var libraries []any
	err := json.Unmarshal(rec.Body.Bytes(), &libraries)
	assert.NoError(t, err)
}

// TestSearchASIN tests the POST /audiobook-wizard/search-asin endpoint
func TestSearchASIN(t *testing.T) {
	e, cleanup := setupTestServer(t)
	defer cleanup()

	t.Run("missing_title", func(t *testing.T) {
		req := SearchASINRequest{
			Author: "Test Author",
		}

		body, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPost, "/api/audiobook-wizard/search-asin", bytes.NewReader(body))
		httpReq.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, httpReq)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid_json", func(t *testing.T) {
		httpReq := httptest.NewRequest(http.MethodPost, "/api/audiobook-wizard/search-asin", bytes.NewReader([]byte("invalid json")))
		httpReq.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, httpReq)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

// TestGetASINMetadata tests the GET /audiobook-wizard/asin/:asin/metadata endpoint
func TestGetASINMetadata(t *testing.T) {
	e, cleanup := setupTestServer(t)
	defer cleanup()

	t.Run("endpoint_exists", func(t *testing.T) {
		httpReq := httptest.NewRequest(http.MethodGet, "/api/audiobook-wizard/asin/B01234567/metadata", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, httpReq)

		// The endpoint should exist (not 404)
		// It will likely fail with an error trying to contact Audible API, which is expected
		assert.NotEqual(t, http.StatusNotFound, rec.Code)
	})
}

// TestPreviewDirectory tests the POST /audiobook-wizard/preview-directory endpoint
func TestPreviewDirectory(t *testing.T) {
	e, cleanup := setupTestServer(t)
	defer cleanup()

	t.Run("valid_metadata", func(t *testing.T) {
		req := PreviewDirectoryRequest{
			Metadata: metadata.BookMetadata{
				Title: "Test Book",
				Asin:  "B01234567",
			},
		}

		body, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPost, "/api/audiobook-wizard/preview-directory", bytes.NewReader(body))
		httpReq.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, httpReq)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response PreviewDirectoryResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response.DirectoryName)
	})

	t.Run("invalid_json", func(t *testing.T) {
		httpReq := httptest.NewRequest(http.MethodPost, "/api/audiobook-wizard/preview-directory", bytes.NewReader([]byte("invalid")))
		httpReq.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, httpReq)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

// TestGetTorrentInfo tests the GET /audiobook-wizard/torrent/:hash/info endpoint
func TestGetTorrentInfo(t *testing.T) {
	t.Skip("Skipping test that requires qBittorrent connection")

	e, cleanup := setupTestServer(t)
	defer cleanup()

	t.Run("endpoint_exists", func(t *testing.T) {
		httpReq := httptest.NewRequest(http.MethodGet, "/api/audiobook-wizard/torrent/abc123/info", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, httpReq)

		// The endpoint should exist (not 404)
		// It will likely fail trying to connect to qBittorrent, which is expected in test environment
		assert.NotEqual(t, http.StatusNotFound, rec.Code)
	})
}

// TestExecuteImport tests the POST /audiobook-wizard/execute-import endpoint
func TestExecuteImport(t *testing.T) {
	e, cleanup := setupTestServer(t)
	defer cleanup()

	t.Run("missing_hash", func(t *testing.T) {
		req := ExecuteImportRequest{
			Metadata: metadata.BookMetadata{
				Title: "Test Book",
				Asin:  "B01234567",
			},
			LibraryName: "test-library",
		}

		body, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPost, "/api/audiobook-wizard/execute-import", bytes.NewReader(body))
		httpReq.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, httpReq)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("missing_library_name", func(t *testing.T) {
		req := ExecuteImportRequest{
			Hash: "abc123",
			Metadata: metadata.BookMetadata{
				Title: "Test Book",
				Asin:  "B01234567",
			},
		}

		body, _ := json.Marshal(req)
		httpReq := httptest.NewRequest(http.MethodPost, "/api/audiobook-wizard/execute-import", bytes.NewReader(body))
		httpReq.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, httpReq)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid_json", func(t *testing.T) {
		httpReq := httptest.NewRequest(http.MethodPost, "/api/audiobook-wizard/execute-import", bytes.NewReader([]byte("invalid")))
		httpReq.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, httpReq)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

// TestSanitizeName tests the sanitizeName helper function
func TestSanitizeName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple_name",
			input:    "Book Title",
			expected: "Book Title",
		},
		{
			name:     "name_with_slash",
			input:    "Book/Title",
			expected: "Book-Title",
		},
		{
			name:     "multiple_slashes",
			input:    "Book/Sub/Title",
			expected: "Book-Sub-Title",
		},
		{
			name:     "empty_string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
