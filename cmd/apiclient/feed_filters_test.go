package apiclient

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestFeedFiltersServer(t *testing.T) *httptest.Server {
	mux := http.NewServeMux()

	// List feed filters
	mux.HandleFunc("/feed-filters", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			feedID := r.URL.Query().Get("feed_id")
			filters := []map[string]interface{}{
				{
					"id":        float64(1),
					"name":      "Filter 1",
					"feed_id":   float64(1),
					"feed_name": "Test Feed",
					"category":  "movies",
					"notifier":  "Discord",
				},
			}

			// Filter by feed_id if provided
			if feedID != "" && feedID != "1" {
				filters = []map[string]interface{}{}
			}

			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(filters); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		case "POST":
			var body map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			response := map[string]interface{}{
				"id":        float64(2),
				"name":      body["name"],
				"feed_id":   float64(1),
				"feed_name": body["feed_name"],
				"category":  body["category_name"],
				"notifier":  body["notifier_name"],
			}
			w.WriteHeader(http.StatusCreated)
			if err := json.NewEncoder(w).Encode(response); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Get, update, delete feed filter
	mux.HandleFunc("/feed-filters/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/feed-filters/")

		switch r.Method {
		case "GET":
			if id == "999" {
				w.WriteHeader(http.StatusNotFound)
				if err := json.NewEncoder(w).Encode(map[string]string{"error": "Filter not found"}); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
			filter := map[string]interface{}{
				"id":        float64(1),
				"name":      "Filter 1",
				"feed_id":   float64(1),
				"feed_name": "Test Feed",
				"category":  "movies",
				"notifier":  "Discord",
			}
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(filter); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		case "PUT":
			var body map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			response := map[string]interface{}{
				"id":        float64(1),
				"name":      body["name"],
				"feed_id":   float64(1),
				"feed_name": body["feed_name"],
				"category":  body["category_name"],
				"notifier":  body["notifier_name"],
			}
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(response); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		case "DELETE":
			if id == "999" {
				w.WriteHeader(http.StatusNotFound)
				if err := json.NewEncoder(w).Encode(map[string]string{"error": "Filter not found"}); err != nil {
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

func TestFeedFiltersListCmd(t *testing.T) {
	server := setupTestFeedFiltersServer(t)
	defer server.Close()

	tests := []struct {
		name           string
		args           []string
		format         string
		expectJSON     bool
		expectError    bool
		expectedOutput string
	}{
		{
			name:           "list feed filters in table format",
			args:           []string{},
			format:         "table",
			expectJSON:     false,
			expectError:    false,
			expectedOutput: "Filter 1",
		},
		{
			name:           "list feed filters with feed_id filter",
			args:           []string{"--feed-id", "1"},
			format:         "table",
			expectJSON:     false,
			expectError:    false,
			expectedOutput: "Filter 1",
		},
		{
			name:           "list feed filters in JSON format",
			args:           []string{},
			format:         "json",
			expectJSON:     true,
			expectError:    false,
			expectedOutput: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiURL = server.URL
			format = tt.format

			cmd := createFeedFiltersListCmd()
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

func TestFeedFiltersGetCmd(t *testing.T) {
	server := setupTestFeedFiltersServer(t)
	defer server.Close()

	tests := []struct {
		name        string
		args        []string
		format      string
		expectError bool
	}{
		{
			name:        "get feed filter by ID in table format",
			args:        []string{"1"},
			format:      "table",
			expectError: false,
		},
		{
			name:        "get feed filter by ID in JSON format",
			args:        []string{"1"},
			format:      "json",
			expectError: false,
		},
		{
			name:        "get feed filter with invalid ID",
			args:        []string{"999"},
			format:      "table",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiURL = server.URL
			format = tt.format

			cmd := createFeedFiltersGetCmd()
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

func TestFeedFiltersCreateCmd(t *testing.T) {
	server := setupTestFeedFiltersServer(t)
	defer server.Close()

	tests := []struct {
		name        string
		args        []string
		format      string
		expectError bool
	}{
		{
			name: "create feed filter successfully",
			args: []string{
				"--name", "New Filter",
				"--feed-name", "Test Feed",
				"--category-name", "movies",
				"--notifier-name", "Discord",
			},
			format:      "table",
			expectError: false,
		},
		{
			name: "create feed filter with JSON output",
			args: []string{
				"--name", "New Filter",
				"--feed-name", "Test Feed",
				"--category-name", "tv",
				"--notifier-name", "Slack",
			},
			format:      "json",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiURL = server.URL
			format = tt.format

			cmd := createFeedFiltersCreateCmd()
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

func TestFeedFiltersUpdateCmd(t *testing.T) {
	server := setupTestFeedFiltersServer(t)
	defer server.Close()

	tests := []struct {
		name        string
		args        []string
		format      string
		expectError bool
	}{
		{
			name:        "update feed filter successfully",
			args:        []string{"1", "--name", "Updated Filter"},
			format:      "table",
			expectError: false,
		},
		{
			name:        "update feed filter category with JSON output",
			args:        []string{"1", "--category-name", "tv"},
			format:      "json",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiURL = server.URL
			format = tt.format

			cmd := createFeedFiltersUpdateCmd()
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

func TestFeedFiltersDeleteCmd(t *testing.T) {
	server := setupTestFeedFiltersServer(t)
	defer server.Close()

	tests := []struct {
		name        string
		args        []string
		format      string
		expectError bool
	}{
		{
			name:        "delete feed filter successfully",
			args:        []string{"1"},
			format:      "table",
			expectError: false,
		},
		{
			name:        "delete feed filter with JSON output",
			args:        []string{"1"},
			format:      "json",
			expectError: false,
		},
		{
			name:        "delete feed filter with invalid ID",
			args:        []string{"999"},
			format:      "table",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiURL = server.URL
			format = tt.format

			cmd := createFeedFiltersDeleteCmd()
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

func TestFeedFiltersCmd(t *testing.T) {
	cmd := createFeedFiltersCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "feed-filters", cmd.Use)
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

func TestFeedFiltersCreateCmd_MissingRequiredFlags(t *testing.T) {
	server := setupTestFeedFiltersServer(t)
	defer server.Close()

	apiURL = server.URL
	format = "table"

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "missing feed-name flag",
			args: []string{"--name", "Test", "--category-name", "movies", "--notifier-name", "Discord"},
		},
		{
			name: "missing category-name flag",
			args: []string{"--name", "Test", "--feed-name", "Feed", "--notifier-name", "Discord"},
		},
		{
			name: "missing notifier-name flag",
			args: []string{"--name", "Test", "--feed-name", "Feed", "--category-name", "movies"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := createFeedFiltersCreateCmd()
			cmd.SetArgs(tt.args)
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)

			err := cmd.Execute()
			require.Error(t, err)
			assert.Contains(t, err.Error(), "required flag")
		})
	}
}

func TestFeedFiltersGetCmd_NoArgs(t *testing.T) {
	server := setupTestFeedFiltersServer(t)
	defer server.Close()

	apiURL = server.URL
	format = "table"

	cmd := createFeedFiltersGetCmd()
	cmd.SetArgs([]string{})
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.Execute()
	require.Error(t, err)
}
