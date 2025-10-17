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

func setupTestFeedFilterSetEntriesServer(t *testing.T) *httptest.Server {
	mux := http.NewServeMux()

	// List feed filter set entries
	mux.HandleFunc("/feed-filter-set-entries", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			setID := r.URL.Query().Get("feed_filter_set_id")
			entries := []map[string]interface{}{
				{
					"id": float64(1),
					"set": map[string]interface{}{
						"id":   float64(1),
						"name": "Test Set",
					},
					"key": map[string]interface{}{
						"id":   float64(1),
						"name": "title",
					},
					"operator": map[string]interface{}{
						"id":   float64(1),
						"name": "contains",
					},
					"value": "test",
				},
			}

			// Filter by feed_filter_set_id if provided
			if setID != "" && setID != "1" {
				entries = []map[string]interface{}{}
			}

			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(entries); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		case "POST":
			var body map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			response := map[string]interface{}{
				"id": float64(2),
				"set": map[string]interface{}{
					"id":   body["feed_filter_set_id"],
					"name": "Test Set",
				},
				"key": map[string]interface{}{
					"id":   float64(1),
					"name": body["key_name"],
				},
				"operator": map[string]interface{}{
					"id":   float64(1),
					"name": body["operator_name"],
				},
				"value": body["value"],
			}
			w.WriteHeader(http.StatusCreated)
			if err := json.NewEncoder(w).Encode(response); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Get, update, delete feed filter set entry
	mux.HandleFunc("/feed-filter-set-entries/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/feed-filter-set-entries/")

		switch r.Method {
		case "GET":
			if id == "999" {
				w.WriteHeader(http.StatusNotFound)
				if err := json.NewEncoder(w).Encode(map[string]string{"error": "Entry not found"}); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
			entry := map[string]interface{}{
				"id": float64(1),
				"set": map[string]interface{}{
					"id":   float64(1),
					"name": "Test Set",
				},
				"key": map[string]interface{}{
					"id":   float64(1),
					"name": "title",
				},
				"operator": map[string]interface{}{
					"id":   float64(1),
					"name": "contains",
				},
				"value": "test",
			}
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(entry); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		case "PUT":
			var body map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			response := map[string]interface{}{
				"id": float64(1),
				"set": map[string]interface{}{
					"id":   body["feed_filter_set_id"],
					"name": "Test Set",
				},
				"key": map[string]interface{}{
					"id":   float64(1),
					"name": body["key_name"],
				},
				"operator": map[string]interface{}{
					"id":   float64(1),
					"name": body["operator_name"],
				},
				"value": body["value"],
			}
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(response); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		case "DELETE":
			if id == "999" {
				w.WriteHeader(http.StatusNotFound)
				if err := json.NewEncoder(w).Encode(map[string]string{"error": "Entry not found"}); err != nil {
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

func TestFeedFilterSetEntriesListCmd(t *testing.T) {
	server := setupTestFeedFilterSetEntriesServer(t)
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
			name:           "list feed filter set entries in table format",
			args:           []string{},
			format:         "table",
			expectJSON:     false,
			expectError:    false,
			expectedOutput: "Test Set",
		},
		{
			name:           "list feed filter set entries with set_id filter",
			args:           []string{"--set-id", "1"},
			format:         "table",
			expectJSON:     false,
			expectError:    false,
			expectedOutput: "Test Set",
		},
		{
			name:           "list feed filter set entries in JSON format",
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

			cmd := createFeedFilterSetEntriesListCmd()
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

func TestFeedFilterSetEntriesGetCmd(t *testing.T) {
	server := setupTestFeedFilterSetEntriesServer(t)
	defer server.Close()

	tests := []struct {
		name        string
		args        []string
		format      string
		expectError bool
	}{
		{
			name:        "get feed filter set entry by ID in table format",
			args:        []string{"1"},
			format:      "table",
			expectError: false,
		},
		{
			name:        "get feed filter set entry by ID in JSON format",
			args:        []string{"1"},
			format:      "json",
			expectError: false,
		},
		{
			name:        "get feed filter set entry with invalid ID",
			args:        []string{"999"},
			format:      "table",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiURL = server.URL
			format = tt.format

			cmd := createFeedFilterSetEntriesGetCmd()
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

func TestFeedFilterSetEntriesCreateCmd(t *testing.T) {
	server := setupTestFeedFilterSetEntriesServer(t)
	defer server.Close()

	tests := []struct {
		name        string
		args        []string
		format      string
		expectError bool
	}{
		{
			name: "create feed filter set entry successfully",
			args: []string{
				"--set-id", "1",
				"--key", "title",
				"--operator", "contains",
				"--value", "test",
			},
			format:      "table",
			expectError: false,
		},
		{
			name: "create feed filter set entry with JSON output",
			args: []string{
				"--set-id", "2",
				"--key", "description",
				"--operator", "matches",
				"--value", "pattern",
			},
			format:      "json",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiURL = server.URL
			format = tt.format

			cmd := createFeedFilterSetEntriesCreateCmd()
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

func TestFeedFilterSetEntriesUpdateCmd(t *testing.T) {
	server := setupTestFeedFilterSetEntriesServer(t)
	defer server.Close()

	tests := []struct {
		name        string
		args        []string
		format      string
		expectError bool
	}{
		{
			name:        "update feed filter set entry successfully",
			args:        []string{"1", "--value", "updated"},
			format:      "table",
			expectError: false,
		},
		{
			name:        "update feed filter set entry operator with JSON output",
			args:        []string{"1", "--operator", "equals"},
			format:      "json",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiURL = server.URL
			format = tt.format

			cmd := createFeedFilterSetEntriesUpdateCmd()
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

func TestFeedFilterSetEntriesDeleteCmd(t *testing.T) {
	server := setupTestFeedFilterSetEntriesServer(t)
	defer server.Close()

	tests := []struct {
		name        string
		args        []string
		format      string
		expectError bool
	}{
		{
			name:        "delete feed filter set entry successfully",
			args:        []string{"1"},
			format:      "table",
			expectError: false,
		},
		{
			name:        "delete feed filter set entry with JSON output",
			args:        []string{"1"},
			format:      "json",
			expectError: false,
		},
		{
			name:        "delete feed filter set entry with invalid ID",
			args:        []string{"999"},
			format:      "table",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiURL = server.URL
			format = tt.format

			cmd := createFeedFilterSetEntriesDeleteCmd()
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

func TestFeedFilterSetEntriesCmd(t *testing.T) {
	cmd := createFeedFilterSetEntriesCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "feed-filter-set-entries", cmd.Use)
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

func TestFeedFilterSetEntriesCreateCmd_MissingRequiredFlags(t *testing.T) {
	server := setupTestFeedFilterSetEntriesServer(t)
	defer server.Close()

	apiURL = server.URL
	format = "table"

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "missing set-id flag",
			args: []string{"--key", "title", "--operator", "contains", "--value", "test"},
		},
		{
			name: "missing key flag",
			args: []string{"--set-id", "1", "--operator", "contains", "--value", "test"},
		},
		{
			name: "missing operator flag",
			args: []string{"--set-id", "1", "--key", "title", "--value", "test"},
		},
		{
			name: "missing value flag",
			args: []string{"--set-id", "1", "--key", "title", "--operator", "contains"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := createFeedFilterSetEntriesCreateCmd()
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

func TestFeedFilterSetEntriesCreateCmd_InvalidSetID(t *testing.T) {
	server := setupTestFeedFilterSetEntriesServer(t)
	defer server.Close()

	apiURL = server.URL
	format = "table"

	cmd := createFeedFilterSetEntriesCreateCmd()
	cmd.SetArgs([]string{
		"--set-id", "invalid",
		"--key", "title",
		"--operator", "contains",
		"--value", "test",
	})
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.Execute()
	require.Error(t, err)
}

func TestFeedFilterSetEntriesGetCmd_NoArgs(t *testing.T) {
	server := setupTestFeedFilterSetEntriesServer(t)
	defer server.Close()

	apiURL = server.URL
	format = "table"

	cmd := createFeedFilterSetEntriesGetCmd()
	cmd.SetArgs([]string{})
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.Execute()
	require.Error(t, err)
}

func TestFeedFilterSetEntriesUpdateCmd_InvalidSetID(t *testing.T) {
	server := setupTestFeedFilterSetEntriesServer(t)
	defer server.Close()

	apiURL = server.URL
	format = "table"

	cmd := createFeedFilterSetEntriesUpdateCmd()
	cmd.SetArgs([]string{"1", "--set-id", "invalid"})
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.Execute()
	require.Error(t, err)
}
