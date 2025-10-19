package apiclient

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestFeedsServer creates a mock server for feeds endpoints
func setupTestFeedsServer(t *testing.T) *httptest.Server {
	mux := http.NewServeMux()

	// List feeds
	mux.HandleFunc("/feeds", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			feeds := []map[string]interface{}{
				{"id": float64(1), "name": "Test Feed 1", "url": "http://example.com/feed1", "enabled": true},
				{"id": float64(2), "name": "Test Feed 2", "url": "http://example.com/feed2", "enabled": false},
			}
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(feeds); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		case "POST":
			var body map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			response := map[string]interface{}{
				"id":      float64(3),
				"name":    body["name"],
				"url":     body["url"],
				"enabled": body["enabled"],
			}
			w.WriteHeader(http.StatusCreated)
			if err := json.NewEncoder(w).Encode(response); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Get, update, delete feed
	mux.HandleFunc("/feeds/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/feeds/")

		switch r.Method {
		case "GET":
			if id == "999" {
				w.WriteHeader(http.StatusNotFound)
				if err := json.NewEncoder(w).Encode(map[string]string{"error": "Feed not found"}); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
			feed := map[string]interface{}{
				"id":      float64(1),
				"name":    "Test Feed",
				"url":     "http://example.com/feed",
				"enabled": true,
			}
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(feed); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		case "PUT":
			var body map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			response := map[string]interface{}{
				"id":      float64(1),
				"name":    body["name"],
				"url":     body["url"],
				"enabled": body["enabled"],
			}
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(response); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		case "DELETE":
			if id == "999" {
				w.WriteHeader(http.StatusNotFound)
				if err := json.NewEncoder(w).Encode(map[string]string{"error": "Feed not found"}); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return httptest.NewServer(mux)
}

func TestFeedsListCmd(t *testing.T) {
	server := setupTestFeedsServer(t)
	defer server.Close()

	tests := []struct {
		name           string
		format         string
		expectJSON     bool
		expectError    bool
		expectedOutput string
	}{
		{
			name:           "list feeds in table format",
			format:         "table",
			expectJSON:     false,
			expectError:    false,
			expectedOutput: "Test Feed 1",
		},
		{
			name:           "list feeds in JSON format",
			format:         "json",
			expectJSON:     true,
			expectError:    false,
			expectedOutput: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set global flags
			apiURL = server.URL
			format = tt.format

			cmd := createFeedsListCmd()
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)

			err := cmd.Execute()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFeedsGetCmd(t *testing.T) {
	server := setupTestFeedsServer(t)
	defer server.Close()

	tests := []struct {
		name        string
		args        []string
		format      string
		expectError bool
	}{
		{
			name:        "get feed by ID in table format",
			args:        []string{"1"},
			format:      "table",
			expectError: false,
		},
		{
			name:        "get feed by ID in JSON format",
			args:        []string{"1"},
			format:      "json",
			expectError: false,
		},
		{
			name:        "get feed with invalid ID",
			args:        []string{"999"},
			format:      "table",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiURL = server.URL
			format = tt.format

			cmd := createFeedsGetCmd()
			cmd.SetArgs(tt.args)
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)

			err := cmd.Execute()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFeedsCreateCmd(t *testing.T) {
	server := setupTestFeedsServer(t)
	defer server.Close()

	tests := []struct {
		name        string
		args        []string
		format      string
		expectError bool
	}{
		{
			name:        "create feed successfully",
			args:        []string{"--name", "New Feed", "--url", "http://example.com/newfeed"},
			format:      "table",
			expectError: false,
		},
		{
			name:        "create feed with JSON output",
			args:        []string{"--name", "New Feed", "--url", "http://example.com/newfeed"},
			format:      "json",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiURL = server.URL
			format = tt.format

			cmd := createFeedsCreateCmd()
			cmd.SetArgs(tt.args)
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)

			err := cmd.Execute()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFeedsUpdateCmd(t *testing.T) {
	server := setupTestFeedsServer(t)
	defer server.Close()

	tests := []struct {
		name        string
		args        []string
		format      string
		expectError bool
	}{
		{
			name:        "update feed successfully",
			args:        []string{"1", "--name", "Updated Feed"},
			format:      "table",
			expectError: false,
		},
		{
			name:        "update feed with JSON output",
			args:        []string{"1", "--url", "http://example.com/updated"},
			format:      "json",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiURL = server.URL
			format = tt.format

			cmd := createFeedsUpdateCmd()
			cmd.SetArgs(tt.args)
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)

			err := cmd.Execute()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFeedsDeleteCmd(t *testing.T) {
	server := setupTestFeedsServer(t)
	defer server.Close()

	tests := []struct {
		name        string
		args        []string
		format      string
		expectError bool
	}{
		{
			name:        "delete feed successfully",
			args:        []string{"1"},
			format:      "table",
			expectError: false,
		},
		{
			name:        "delete feed with JSON output",
			args:        []string{"1"},
			format:      "json",
			expectError: false,
		},
		{
			name:        "delete feed with invalid ID",
			args:        []string{"999"},
			format:      "table",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiURL = server.URL
			format = tt.format

			cmd := createFeedsDeleteCmd()
			cmd.SetArgs(tt.args)
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)

			err := cmd.Execute()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFeedsCmd(t *testing.T) {
	cmd := createFeedsCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "feeds", cmd.Use)
	assert.True(t, cmd.HasSubCommands())

	// Verify all subcommands are registered
	subcommands := []string{"list", "get", "create", "update", "delete"}
	for _, subcmd := range subcommands {
		found := false
		for _, c := range cmd.Commands() {
			if strings.HasPrefix(c.Use, subcmd) {
				found = true
				break
			}
		}
		assert.True(t, found, "Expected subcommand %s to be registered", subcmd)
	}
}


func TestFeedsCreateCmd_MissingRequiredFlags(t *testing.T) {
	server := setupTestFeedsServer(t)
	defer server.Close()

	apiURL = server.URL
	format = "table"

	cmd := createFeedsCreateCmd()
	cmd.SetArgs([]string{"--name", "Test"}) // Missing --url flag
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "required flag")
}

func TestFeedsGetCmd_NoArgs(t *testing.T) {
	server := setupTestFeedsServer(t)
	defer server.Close()

	apiURL = server.URL
	format = "table"

	cmd := createFeedsGetCmd()
	cmd.SetArgs([]string{}) // No ID provided
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.Execute()
	require.Error(t, err)
}

func TestFeedsCmd_Integration(t *testing.T) {
	server := setupTestFeedsServer(t)
	defer server.Close()

	apiURL = server.URL
	format = "json"

	// Create parent command
	parentCmd := &cobra.Command{Use: "test"}
	feedsCmd := createFeedsCmd()
	parentCmd.AddCommand(feedsCmd)

	// Test list command
	listCmd, _, err := parentCmd.Find([]string{"feeds", "list"})
	require.NoError(t, err)
	assert.Equal(t, "list", listCmd.Use)

	// Test get command
	getCmd, _, err := parentCmd.Find([]string{"feeds", "get"})
	require.NoError(t, err)
	assert.Equal(t, "get <id>", getCmd.Use)
}
