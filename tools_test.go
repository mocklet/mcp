package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"mocklet-mcp/client"
)

func TestHandlers(t *testing.T) {
	// Setup dummy HAR file for tests
	tmpDir := t.TempDir()
	harPath := filepath.Join(tmpDir, "test.har")
	err := os.WriteFile(harPath, []byte(`{"log": {"entries": []}}`), 0644)
	if err != nil {
		t.Fatalf("failed to write dummy har file: %v", err)
	}

	// Setup mock HTTP server
	mux := http.NewServeMux()
	
	// mocklet_validate_har endpoint
	mux.HandleFunc("/api/v1/har/validate-file", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"valid": true, "errors": []}`))
	})
	
	// mocklet_create_mock endpoint
	mux.HandleFunc("/api/v1/mocks/file", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"mock_id": "mock_123", "url": "http://mock_123.example.com"}`))
	})
	
	ts := httptest.NewServer(mux)
	defer ts.Close()

	// Initialize the global apiClient to point to the mock server
	c, err := client.NewClientWithResponses(ts.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	apiClient = c

	t.Run("mocklet_validate_har", func(t *testing.T) {
		req := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Arguments: map[string]interface{}{
					"har_file_path": harPath,
				},
			},
		}

		res, err := validateHarHandler(context.Background(), req)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if res.IsError {
			t.Fatalf("expected IsError to be false")
		}
		
		content := res.Content[0].(mcp.TextContent).Text
		if content == "" {
			t.Errorf("expected content in result")
		}
	})

	t.Run("mocklet_create_mock", func(t *testing.T) {
		req := mcp.CallToolRequest{
			Params: mcp.CallToolParams{
				Arguments: map[string]interface{}{
					"har_file_path": harPath,
					"ttl_seconds":   float64(3600),
				},
			},
		}

		res, err := createMockHandler(context.Background(), req)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if res.IsError {
			t.Fatalf("expected IsError to be false")
		}
		
		content := res.Content[0].(mcp.TextContent).Text
		if content == "" {
			t.Errorf("expected content in result")
		}
	})
}
