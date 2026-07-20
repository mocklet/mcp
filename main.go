package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"mocklet-mcp/client"
)

var (
	apiClient *client.ClientWithResponses
	mcpServer *server.MCPServer
)

func main() {
	apiUrl := os.Getenv("MOCKLET_API_URL")
	if apiUrl == "" {
		apiUrl = "http://localhost:8080"
	}
	token := os.Getenv("MOCKLET_SERVICE_TOKEN")
	if token == "" {
		fmt.Fprintf(os.Stderr, "Warning: MOCKLET_SERVICE_TOKEN is not set\n")
	}

	c, err := client.NewClientWithResponses(apiUrl, client.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		return nil
	}))
	if err != nil {
		log.Fatalf("Failed to create API client: %v", err)
	}
	apiClient = c

	mcpServer = server.NewMCPServer(
		"mocklet-mcp",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithPromptCapabilities(true),
	)

	registerTools()
	registerPrompts()

	fmt.Fprintf(os.Stderr, "Starting Mocklet MCP Server...\n")
	if err := server.ServeStdio(mcpServer); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
