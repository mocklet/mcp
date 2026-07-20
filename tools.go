package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

func registerTools() {
	mcpServer.AddTool(mcp.NewTool("mocklet_validate_har",
		mcp.WithDescription("Pre-flight validation of a generated HAR file before attempting to deploy it."),
		mcp.WithString("har_file_path", mcp.Required(), mcp.Description("Absolute path to the local HAR file")),
	), validateHarHandler)

	mcpServer.AddTool(mcp.NewTool("mocklet_create_mock",
		mcp.WithDescription("Create a one-off ephemeral mock server directly from a HAR file."),
		mcp.WithString("har_file_path", mcp.Required(), mcp.Description("Absolute path to the local HAR file")),
		mcp.WithNumber("ttl_seconds", mcp.Description("Time-to-live in seconds")),
	), createMockHandler)

	mcpServer.AddTool(mcp.NewTool("mocklet_list_mocks",
		mcp.WithDescription("Retrieve a paginated list of currently active ephemeral mocks."),
		mcp.WithString("search_query", mcp.Description("Search query")),
		mcp.WithNumber("limit", mcp.Description("Max results (default 5)")),
	), listMocksHandler)

	mcpServer.AddTool(mcp.NewTool("mocklet_get_mock_stats",
		mcp.WithDescription("Retrieve hit/miss statistics for a specific mock."),
		mcp.WithString("mock_id", mcp.Required(), mcp.Description("Mock ID")),
	), getMockStatsHandler)

	mcpServer.AddTool(mcp.NewTool("mocklet_delete_mock",
		mcp.WithDescription("Stop and delete an active mock server."),
		mcp.WithString("mock_id", mcp.Required(), mcp.Description("Mock ID")),
	), deleteMockHandler)

	mcpServer.AddTool(mcp.NewTool("mocklet_create_template",
		mcp.WithDescription("Upload a HAR file once to create a reusable persistent template."),
		mcp.WithString("name", mcp.Required(), mcp.Description("Template name")),
		mcp.WithString("description", mcp.Description("Description")),
		mcp.WithString("har_file_path", mcp.Required(), mcp.Description("Absolute path to the local HAR file")),
		mcp.WithNumber("default_ttl_seconds", mcp.Description("Default TTL")),
	), createTemplateHandler)

	mcpServer.AddTool(mcp.NewTool("mocklet_list_templates",
		mcp.WithDescription("Find existing templates available for reuse."),
		mcp.WithString("search_query", mcp.Description("Search query")),
		mcp.WithNumber("limit", mcp.Description("Max results (default 5)")),
	), listTemplatesHandler)

	mcpServer.AddTool(mcp.NewTool("mocklet_spawn_mock",
		mcp.WithDescription("Rapidly spin up an ephemeral mock from an existing template."),
		mcp.WithString("template_public_id", mcp.Required(), mcp.Description("Template Public ID")),
	), spawnMockHandler)

	mcpServer.AddTool(mcp.NewTool("mocklet_upload_template_revision",
		mcp.WithDescription("Update an existing template's behavior with a new HAR file."),
		mcp.WithString("template_public_id", mcp.Required(), mcp.Description("Template Public ID")),
		mcp.WithString("har_file_path", mcp.Required(), mcp.Description("Absolute path to the local HAR file")),
	), uploadTemplateRevisionHandler)
}

func formatError(status int, body []byte) *mcp.CallToolResult {
	return mcp.NewToolResultError(fmt.Sprintf("API Error %d: %s", status, string(body)))
}

func validateHarHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	harPath, err := request.RequireString("har_file_path")
	if err != nil {
		return mcp.NewToolResultError("har_file_path must be a string"), nil
	}

	contentType, body, err := createMultipartFileWithFields(harPath, nil)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to read file: %v", err)), nil
	}

	resp, err := apiClient.ValidateHarFileWithBodyWithResponse(ctx, contentType, body)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Request failed: %v", err)), nil
	}

	if resp.StatusCode() != 200 {
		return formatError(resp.StatusCode(), resp.Body), nil
	}

	b, _ := json.MarshalIndent(resp.JSON200, "", "  ")
	return mcp.NewToolResultText(fmt.Sprintf("HAR Validation Success:\n```json\n%s\n```", string(b))), nil
}
