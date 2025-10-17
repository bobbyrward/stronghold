package apiclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIClient_Get(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse interface{}
		serverStatus   int
		expectedError  bool
		path           string
	}{
		{
			name:           "successful GET request",
			serverResponse: map[string]string{"message": "success"},
			serverStatus:   http.StatusOK,
			expectedError:  false,
			path:           "/test",
		},
		{
			name:           "GET request with 404",
			serverResponse: map[string]string{"error": "not found"},
			serverStatus:   http.StatusNotFound,
			expectedError:  true,
			path:           "/notfound",
		},
		{
			name:           "GET request with server error",
			serverResponse: map[string]string{"error": "internal error"},
			serverStatus:   http.StatusInternalServerError,
			expectedError:  true,
			path:           "/error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "GET", r.Method)
				assert.Equal(t, tt.path, r.URL.Path)
				w.WriteHeader(tt.serverStatus)
				if err := json.NewEncoder(w).Encode(tt.serverResponse); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}))
			defer server.Close()

			client := NewClient(server.URL)
			var response map[string]string
			err := client.Get(context.Background(), tt.path, &response)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.serverResponse, response)
			}
		})
	}
}

func TestAPIClient_Post(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		serverResponse interface{}
		serverStatus   int
		expectedError  bool
		path           string
	}{
		{
			name:           "successful POST request",
			requestBody:    map[string]string{"name": "test"},
			serverResponse: map[string]interface{}{"id": float64(1), "name": "test"},
			serverStatus:   http.StatusCreated,
			expectedError:  false,
			path:           "/test",
		},
		{
			name:           "POST request with validation error",
			requestBody:    map[string]string{},
			serverResponse: map[string]string{"error": "validation failed"},
			serverStatus:   http.StatusBadRequest,
			expectedError:  true,
			path:           "/test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "POST", r.Method)
				assert.Equal(t, tt.path, r.URL.Path)

				var receivedBody map[string]string
				err := json.NewDecoder(r.Body).Decode(&receivedBody)
				require.NoError(t, err)
				assert.Equal(t, tt.requestBody, receivedBody)

				w.WriteHeader(tt.serverStatus)
				if err := json.NewEncoder(w).Encode(tt.serverResponse); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}))
			defer server.Close()

			client := NewClient(server.URL)
			var response map[string]interface{}
			err := client.Post(context.Background(), tt.path, tt.requestBody, &response)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.serverResponse, response)
			}
		})
	}
}

func TestAPIClient_Put(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		serverResponse interface{}
		serverStatus   int
		expectedError  bool
		path           string
	}{
		{
			name:           "successful PUT request",
			requestBody:    map[string]string{"name": "updated"},
			serverResponse: map[string]interface{}{"id": float64(1), "name": "updated"},
			serverStatus:   http.StatusOK,
			expectedError:  false,
			path:           "/test/1",
		},
		{
			name:           "PUT request with not found error",
			requestBody:    map[string]string{"name": "updated"},
			serverResponse: map[string]string{"error": "not found"},
			serverStatus:   http.StatusNotFound,
			expectedError:  true,
			path:           "/test/999",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "PUT", r.Method)
				assert.Equal(t, tt.path, r.URL.Path)

				var receivedBody map[string]string
				err := json.NewDecoder(r.Body).Decode(&receivedBody)
				require.NoError(t, err)
				assert.Equal(t, tt.requestBody, receivedBody)

				w.WriteHeader(tt.serverStatus)
				if err := json.NewEncoder(w).Encode(tt.serverResponse); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}))
			defer server.Close()

			client := NewClient(server.URL)
			var response map[string]interface{}
			err := client.Put(context.Background(), tt.path, tt.requestBody, &response)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.serverResponse, response)
			}
		})
	}
}

func TestAPIClient_Delete(t *testing.T) {
	tests := []struct {
		name          string
		serverStatus  int
		expectedError bool
		path          string
	}{
		{
			name:          "successful DELETE request",
			serverStatus:  http.StatusNoContent,
			expectedError: false,
			path:          "/test/1",
		},
		{
			name:          "DELETE request with not found error",
			serverStatus:  http.StatusNotFound,
			expectedError: true,
			path:          "/test/999",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "DELETE", r.Method)
				assert.Equal(t, tt.path, r.URL.Path)
				w.WriteHeader(tt.serverStatus)
			}))
			defer server.Close()

			client := NewClient(server.URL)
			err := client.Delete(context.Background(), tt.path)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	baseURL := "http://localhost:8000"
	client := NewClient(baseURL)
	assert.NotNil(t, client)
	assert.Equal(t, baseURL, client.BaseURL)
}
