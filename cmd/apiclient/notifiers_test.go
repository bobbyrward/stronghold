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

func setupTestNotifiersServer(t *testing.T) *httptest.Server {
	mux := http.NewServeMux()

	// List notifiers
	mux.HandleFunc("/notifiers", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			notifiers := []map[string]interface{}{
				{
					"id":      float64(1),
					"name":    "Discord Notifier",
					"enabled": true,
					"type": map[string]interface{}{
						"id":   float64(1),
						"name": "discord",
					},
					"webhook_url": "https://discord.com/webhook/123",
				},
				{
					"id":      float64(2),
					"name":    "Slack Notifier",
					"enabled": false,
					"type": map[string]interface{}{
						"id":   float64(2),
						"name": "slack",
					},
					"webhook_url": "https://slack.com/webhook/456",
				},
			}
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(notifiers); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		case "POST":
			var body map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			response := map[string]interface{}{
				"id":   float64(3),
				"name": body["name"],
				"type": map[string]interface{}{
					"id":   float64(1),
					"name": body["type_name"],
				},
				"webhook_url": body["webhook_url"],
				"enabled":     body["enabled"],
			}
			w.WriteHeader(http.StatusCreated)
			if err := json.NewEncoder(w).Encode(response); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Get, update, delete notifier
	mux.HandleFunc("/notifiers/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/notifiers/")

		switch r.Method {
		case "GET":
			if id == "999" {
				w.WriteHeader(http.StatusNotFound)
				if err := json.NewEncoder(w).Encode(map[string]string{"error": "Notifier not found"}); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
			notifier := map[string]interface{}{
				"id":   float64(1),
				"name": "Discord Notifier",
				"type": map[string]interface{}{
					"id":   float64(1),
					"name": "discord",
				},
				"webhook_url": "https://discord.com/webhook/123",
				"enabled":     true,
			}
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(notifier); err != nil {
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
				"type": map[string]interface{}{
					"id":   float64(1),
					"name": body["type_name"],
				},
				"webhook_url": body["webhook_url"],
				"enabled":     body["enabled"],
			}
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(response); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		case "DELETE":
			if id == "999" {
				w.WriteHeader(http.StatusNotFound)
				if err := json.NewEncoder(w).Encode(map[string]string{"error": "Notifier not found"}); err != nil {
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

func TestNotifiersListCmd(t *testing.T) {
	server := setupTestNotifiersServer(t)
	defer server.Close()

	tests := []struct {
		name           string
		format         string
		expectJSON     bool
		expectError    bool
		expectedOutput string
	}{
		{
			name:           "list notifiers in table format",
			format:         "table",
			expectJSON:     false,
			expectError:    false,
			expectedOutput: "Discord Notifier",
		},
		{
			name:           "list notifiers in JSON format",
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

			cmd := createNotifiersListCmd()
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

func TestNotifiersGetCmd(t *testing.T) {
	server := setupTestNotifiersServer(t)
	defer server.Close()

	tests := []struct {
		name        string
		args        []string
		format      string
		expectError bool
	}{
		{
			name:        "get notifier by ID in table format",
			args:        []string{"1"},
			format:      "table",
			expectError: false,
		},
		{
			name:        "get notifier by ID in JSON format",
			args:        []string{"1"},
			format:      "json",
			expectError: false,
		},
		{
			name:        "get notifier with invalid ID",
			args:        []string{"999"},
			format:      "table",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiURL = server.URL
			format = tt.format

			cmd := createNotifiersGetCmd()
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

func TestNotifiersCreateCmd(t *testing.T) {
	server := setupTestNotifiersServer(t)
	defer server.Close()

	tests := []struct {
		name        string
		args        []string
		format      string
		expectError bool
	}{
		{
			name: "create notifier successfully",
			args: []string{
				"--name", "New Notifier",
				"--type", "discord",
				"--webhook-url", "https://discord.com/webhook/new",
			},
			format:      "table",
			expectError: false,
		},
		{
			name: "create notifier with JSON output",
			args: []string{
				"--name", "New Notifier",
				"--type", "slack",
				"--webhook-url", "https://slack.com/webhook/new",
				"--enabled=false",
			},
			format:      "json",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiURL = server.URL
			format = tt.format

			cmd := createNotifiersCreateCmd()
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

func TestNotifiersUpdateCmd(t *testing.T) {
	server := setupTestNotifiersServer(t)
	defer server.Close()

	tests := []struct {
		name        string
		args        []string
		format      string
		expectError bool
	}{
		{
			name:        "update notifier successfully",
			args:        []string{"1", "--name", "Updated Notifier"},
			format:      "table",
			expectError: false,
		},
		{
			name:        "update notifier webhook with JSON output",
			args:        []string{"1", "--webhook-url", "https://discord.com/webhook/updated"},
			format:      "json",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiURL = server.URL
			format = tt.format

			cmd := createNotifiersUpdateCmd()
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

func TestNotifiersDeleteCmd(t *testing.T) {
	server := setupTestNotifiersServer(t)
	defer server.Close()

	tests := []struct {
		name        string
		args        []string
		format      string
		expectError bool
	}{
		{
			name:        "delete notifier successfully",
			args:        []string{"1"},
			format:      "table",
			expectError: false,
		},
		{
			name:        "delete notifier with JSON output",
			args:        []string{"1"},
			format:      "json",
			expectError: false,
		},
		{
			name:        "delete notifier with invalid ID",
			args:        []string{"999"},
			format:      "table",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiURL = server.URL
			format = tt.format

			cmd := createNotifiersDeleteCmd()
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

func TestNotifiersCmd(t *testing.T) {
	cmd := createNotifiersCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "notifiers", cmd.Use)
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

func TestNotifiersCreateCmd_MissingRequiredFlags(t *testing.T) {
	server := setupTestNotifiersServer(t)
	defer server.Close()

	apiURL = server.URL
	format = "table"

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "missing type flag",
			args: []string{"--name", "Test", "--webhook-url", "https://example.com"},
		},
		{
			name: "missing webhook-url flag",
			args: []string{"--name", "Test", "--type", "discord"},
		},
		{
			name: "missing name flag",
			args: []string{"--type", "discord", "--webhook-url", "https://example.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := createNotifiersCreateCmd()
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

func TestNotifiersGetCmd_NoArgs(t *testing.T) {
	server := setupTestNotifiersServer(t)
	defer server.Close()

	apiURL = server.URL
	format = "table"

	cmd := createNotifiersGetCmd()
	cmd.SetArgs([]string{})
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.Execute()
	require.Error(t, err)
}
