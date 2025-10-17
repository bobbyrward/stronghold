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

func setupTestFeedFilterSetsServer(t *testing.T) *httptest.Server {
	mux := http.NewServeMux()

	// List feed filter sets
	mux.HandleFunc("/feed-filter-sets", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			feedFilterID := r.URL.Query().Get("feed_filter_id")
			sets := []map[string]interface{}{
				{
					"id":   float64(1),
					"name": "Set 1",
					"feed": map[string]interface{}{
						"id":   float64(1),
						"name": "Test Feed",
					},
					"type": map[string]interface{}{
						"id":   float64(1),
						"name": "include",
					},
					"category": map[string]interface{}{
						"id":   float64(1),
						"name": "movies",
					},
					"notifier": map[string]interface{}{
						"id":   float64(1),
						"name": "Discord",
					},
				},
			}

			// Filter by feed_filter_id if provided
			if feedFilterID != "" && feedFilterID != "1" {
				sets = []map[string]interface{}{}
			}

			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(sets); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		case "POST":
			var body map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			response := map[string]interface{}{
				"id":   float64(2),
				"name": body["name"],
				"feed": map[string]interface{}{
					"id":   body["feed_id"],
					"name": "Test Feed",
				},
				"type": map[string]interface{}{
					"id":   float64(1),
					"name": body["type_name"],
				},
				"category": map[string]interface{}{
					"id":   float64(1),
					"name": body["category_name"],
				},
				"notifier": map[string]interface{}{
					"id":   body["notifier_id"],
					"name": "Discord",
				},
			}
			w.WriteHeader(http.StatusCreated)
			if err := json.NewEncoder(w).Encode(response); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Get, update, delete feed filter set
	mux.HandleFunc("/feed-filter-sets/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/feed-filter-sets/")

		switch r.Method {
		case "GET":
			if id == "999" {
				w.WriteHeader(http.StatusNotFound)
				if err := json.NewEncoder(w).Encode(map[string]string{"error": "Set not found"}); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
			set := map[string]interface{}{
				"id":   float64(1),
				"name": "Set 1",
				"feed": map[string]interface{}{
					"id":   float64(1),
					"name": "Test Feed",
				},
				"type": map[string]interface{}{
					"id":   float64(1),
					"name": "include",
				},
				"category": map[string]interface{}{
					"id":   float64(1),
					"name": "movies",
				},
				"notifier": map[string]interface{}{
					"id":   float64(1),
					"name": "Discord",
				},
			}
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(set); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		case "PUT":
			var body map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			response := map[string]interface{}{
				"id":   float64(1),
				"name": body["name"],
				"feed": map[string]interface{}{
					"id":   body["feed_id"],
					"name": "Test Feed",
				},
				"type": map[string]interface{}{
					"id":   float64(1),
					"name": body["type_name"],
				},
				"category": map[string]interface{}{
					"id":   float64(1),
					"name": body["category_name"],
				},
				"notifier": map[string]interface{}{
					"id":   body["notifier_id"],
					"name": "Discord",
				},
			}
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(response); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		case "DELETE":
			if id == "999" {
				w.WriteHeader(http.StatusNotFound)
				if err := json.NewEncoder(w).Encode(map[string]string{"error": "Set not found"}); err != nil {
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

func TestFeedFilterSetsListCmd(t *testing.T) {
	server := setupTestFeedFilterSetsServer(t)
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
			name:           "list feed filter sets in table format",
			args:           []string{},
			format:         "table",
			expectJSON:     false,
			expectError:    false,
			expectedOutput: "Set 1",
		},
		{
			name:           "list feed filter sets with feed_id filter",
			args:           []string{"--feed-id", "1"},
			format:         "table",
			expectJSON:     false,
			expectError:    false,
			expectedOutput: "Set 1",
		},
		{
			name:           "list feed filter sets in JSON format",
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

			cmd := createFeedFilterSetsListCmd()
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

func TestFeedFilterSetsGetCmd(t *testing.T) {
	server := setupTestFeedFilterSetsServer(t)
	defer server.Close()

	tests := []struct {
		name        string
		args        []string
		format      string
		expectError bool
	}{
		{
			name:        "get feed filter set by ID in table format",
			args:        []string{"1"},
			format:      "table",
			expectError: false,
		},
		{
			name:        "get feed filter set by ID in JSON format",
			args:        []string{"1"},
			format:      "json",
			expectError: false,
		},
		{
			name:        "get feed filter set with invalid ID",
			args:        []string{"999"},
			format:      "table",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiURL = server.URL
			format = tt.format

			cmd := createFeedFilterSetsGetCmd()
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

func TestFeedFilterSetsCreateCmd(t *testing.T) {
	server := setupTestFeedFilterSetsServer(t)
	defer server.Close()

	tests := []struct {
		name        string
		args        []string
		format      string
		expectError bool
	}{
		{
			name: "create feed filter set successfully",
			args: []string{
				"--name", "New Set",
				"--feed-id", "1",
				"--type", "include",
				"--category", "movies",
				"--notifier-id", "1",
			},
			format:      "table",
			expectError: false,
		},
		{
			name: "create feed filter set with JSON output",
			args: []string{
				"--name", "New Set",
				"--feed-id", "2",
				"--type", "exclude",
				"--category", "tv",
				"--notifier-id", "2",
			},
			format:      "json",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiURL = server.URL
			format = tt.format

			cmd := createFeedFilterSetsCreateCmd()
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

func TestFeedFilterSetsUpdateCmd(t *testing.T) {
	server := setupTestFeedFilterSetsServer(t)
	defer server.Close()

	tests := []struct {
		name        string
		args        []string
		format      string
		expectError bool
	}{
		{
			name:        "update feed filter set successfully",
			args:        []string{"1", "--name", "Updated Set"},
			format:      "table",
			expectError: false,
		},
		{
			name:        "update feed filter set type with JSON output",
			args:        []string{"1", "--type", "exclude"},
			format:      "json",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiURL = server.URL
			format = tt.format

			cmd := createFeedFilterSetsUpdateCmd()
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

func TestFeedFilterSetsDeleteCmd(t *testing.T) {
	server := setupTestFeedFilterSetsServer(t)
	defer server.Close()

	tests := []struct {
		name        string
		args        []string
		format      string
		expectError bool
	}{
		{
			name:        "delete feed filter set successfully",
			args:        []string{"1"},
			format:      "table",
			expectError: false,
		},
		{
			name:        "delete feed filter set with JSON output",
			args:        []string{"1"},
			format:      "json",
			expectError: false,
		},
		{
			name:        "delete feed filter set with invalid ID",
			args:        []string{"999"},
			format:      "table",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiURL = server.URL
			format = tt.format

			cmd := createFeedFilterSetsDeleteCmd()
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

func TestFeedFilterSetsCmd(t *testing.T) {
	cmd := createFeedFilterSetsCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "feed-filter-sets", cmd.Use)
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

func TestFeedFilterSetsCreateCmd_MissingRequiredFlags(t *testing.T) {
	server := setupTestFeedFilterSetsServer(t)
	defer server.Close()

	apiURL = server.URL
	format = "table"

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "missing feed-id flag",
			args: []string{"--name", "Test", "--type", "include", "--category", "movies", "--notifier-id", "1"},
		},
		{
			name: "missing type flag",
			args: []string{"--name", "Test", "--feed-id", "1", "--category", "movies", "--notifier-id", "1"},
		},
		{
			name: "missing category flag",
			args: []string{"--name", "Test", "--feed-id", "1", "--type", "include", "--notifier-id", "1"},
		},
		{
			name: "missing notifier-id flag",
			args: []string{"--name", "Test", "--feed-id", "1", "--type", "include", "--category", "movies"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := createFeedFilterSetsCreateCmd()
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

func TestFeedFilterSetsCreateCmd_InvalidFeedID(t *testing.T) {
	server := setupTestFeedFilterSetsServer(t)
	defer server.Close()

	apiURL = server.URL
	format = "table"

	cmd := createFeedFilterSetsCreateCmd()
	cmd.SetArgs([]string{
		"--name", "Test",
		"--feed-id", "invalid",
		"--type", "include",
		"--category", "movies",
		"--notifier-id", "1",
	})
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.Execute()
	require.Error(t, err)
}

func TestFeedFilterSetsCreateCmd_InvalidNotifierID(t *testing.T) {
	server := setupTestFeedFilterSetsServer(t)
	defer server.Close()

	apiURL = server.URL
	format = "table"

	cmd := createFeedFilterSetsCreateCmd()
	cmd.SetArgs([]string{
		"--name", "Test",
		"--feed-id", "1",
		"--type", "include",
		"--category", "movies",
		"--notifier-id", "invalid",
	})
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.Execute()
	require.Error(t, err)
}

func TestFeedFilterSetsGetCmd_NoArgs(t *testing.T) {
	server := setupTestFeedFilterSetsServer(t)
	defer server.Close()

	apiURL = server.URL
	format = "table"

	cmd := createFeedFilterSetsGetCmd()
	cmd.SetArgs([]string{})
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.Execute()
	require.Error(t, err)
}
