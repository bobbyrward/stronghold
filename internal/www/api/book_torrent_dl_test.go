package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bobbyrward/stronghold/internal/models"
	"github.com/bobbyrward/stronghold/internal/qbit"
	"github.com/bobbyrward/stronghold/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDownloadBookTorrent_Success tests successful torrent download
func TestDownloadBookTorrent_Success(t *testing.T) {
	db, err := models.ConnectTestDB()
	require.NoError(t, err)

	e := SetupTestServerWithDB(db)

	mockClient := &testutil.MockQbitClient{}
	var qbitClient qbit.QbitClient = mockClient

	// Register the route with mock client
	e.POST("/api/book-download", DownloadBookTorrent(db, &qbitClient))

	req := BookTorrentDLRequest{
		Category:   "audiobooks",
		TorrentID:  "12345",
		TorrentURL: "https://example.com/torrent/12345",
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/api/book-download", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httpReq)

	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify the mock was called with correct parameters
	require.Len(t, mockClient.AddTorrentFromUrlCtxCalls, 1)
	call := mockClient.AddTorrentFromUrlCtxCalls[0]
	assert.Equal(t, "https://example.com/torrent/12345", call.URL)
	assert.Equal(t, "true", call.Options["autoTMM"])
	assert.Equal(t, "audiobooks", call.Options["category"])

	// Verify response
	var response map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
}

// TestDownloadBookTorrent_MissingCategory tests validation error for missing category
func TestDownloadBookTorrent_MissingCategory(t *testing.T) {
	db, err := models.ConnectTestDB()
	require.NoError(t, err)

	e := SetupTestServerWithDB(db)

	mockClient := &testutil.MockQbitClient{}
	var qbitClient qbit.QbitClient = mockClient
	e.POST("/api/book-download", DownloadBookTorrent(db, &qbitClient))

	req := BookTorrentDLRequest{
		// Category missing
		TorrentID:  "12345",
		TorrentURL: "https://example.com/torrent/12345",
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/api/book-download", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httpReq)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Verify mock was not called
	assert.Len(t, mockClient.AddTorrentFromUrlCtxCalls, 0)
}

// TestDownloadBookTorrent_MissingTorrentID tests validation error for missing torrent_id
func TestDownloadBookTorrent_MissingTorrentID(t *testing.T) {
	db, err := models.ConnectTestDB()
	require.NoError(t, err)

	e := SetupTestServerWithDB(db)

	mockClient := &testutil.MockQbitClient{}
	var qbitClient qbit.QbitClient = mockClient
	e.POST("/api/book-download", DownloadBookTorrent(db, &qbitClient))

	req := BookTorrentDLRequest{
		Category: "audiobooks",
		// TorrentID missing
		TorrentURL: "https://example.com/torrent/12345",
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/api/book-download", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httpReq)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Verify mock was not called
	assert.Len(t, mockClient.AddTorrentFromUrlCtxCalls, 0)
}

// TestDownloadBookTorrent_MissingTorrentURL tests validation error for missing torrent_url
func TestDownloadBookTorrent_MissingTorrentURL(t *testing.T) {
	db, err := models.ConnectTestDB()
	require.NoError(t, err)

	e := SetupTestServerWithDB(db)

	mockClient := &testutil.MockQbitClient{}
	var qbitClient qbit.QbitClient = mockClient
	e.POST("/api/book-download", DownloadBookTorrent(db, &qbitClient))

	req := BookTorrentDLRequest{
		Category:  "audiobooks",
		TorrentID: "12345",
		// TorrentURL missing
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/api/book-download", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httpReq)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Verify mock was not called
	assert.Len(t, mockClient.AddTorrentFromUrlCtxCalls, 0)
}

// TestDownloadBookTorrent_InvalidJSON tests handling of invalid JSON
func TestDownloadBookTorrent_InvalidJSON(t *testing.T) {
	db, err := models.ConnectTestDB()
	require.NoError(t, err)

	e := SetupTestServerWithDB(db)

	mockClient := &testutil.MockQbitClient{}
	var qbitClient qbit.QbitClient = mockClient
	e.POST("/api/book-download", DownloadBookTorrent(db, &qbitClient))

	httpReq := httptest.NewRequest(http.MethodPost, "/api/book-download", bytes.NewReader([]byte("invalid json")))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httpReq)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Verify mock was not called
	assert.Len(t, mockClient.AddTorrentFromUrlCtxCalls, 0)
}

// TestDownloadBookTorrent_QbitClientError tests handling of qBittorrent client errors
func TestDownloadBookTorrent_QbitClientError(t *testing.T) {
	db, err := models.ConnectTestDB()
	require.NoError(t, err)

	e := SetupTestServerWithDB(db)

	mockClient := &testutil.MockQbitClient{
		AddTorrentFromUrlCtxReturn: errors.New("qbittorrent connection failed"),
	}
	var qbitClient qbit.QbitClient = mockClient
	e.POST("/api/book-download", DownloadBookTorrent(db, &qbitClient))

	req := BookTorrentDLRequest{
		Category:   "audiobooks",
		TorrentID:  "12345",
		TorrentURL: "https://example.com/torrent/12345",
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/api/book-download", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httpReq)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	// Verify mock was called
	require.Len(t, mockClient.AddTorrentFromUrlCtxCalls, 1)

	// Verify error response
	var response map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "failed to add torrent to qbit", response["error"])
}

// TestDownloadBookTorrent_NilQbitClient tests that handler creates client when nil
func TestDownloadBookTorrent_NilQbitClient(t *testing.T) {
	t.Skip("Skipping test that requires actual qBittorrent connection")

	db, err := models.ConnectTestDB()
	require.NoError(t, err)

	e := SetupTestServerWithDB(db)

	// Pass nil client - handler should create one
	e.POST("/api/book-download", DownloadBookTorrent(db, nil))

	req := BookTorrentDLRequest{
		Category:   "audiobooks",
		TorrentID:  "12345",
		TorrentURL: "https://example.com/torrent/12345",
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/api/book-download", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httpReq)

	// This will likely fail with internal error if qBittorrent is not configured
	// but the endpoint should exist and attempt to create a client
	assert.NotEqual(t, http.StatusNotFound, rec.Code)
}
