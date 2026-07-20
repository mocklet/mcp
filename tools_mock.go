package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"mocklet-mcp/client"
)

func createMockHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	harPath, err := request.RequireString("har_file_path")
	if err != nil {
		return mcp.NewToolResultError("har_file_path must be a string"), nil
	}

	ttlSeconds := int(request.GetFloat("ttl_seconds", 3600))

	contentType, body, err := createMultipartFileWithFields(harPath, nil)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to prepare request: %v", err)), nil
	}

	params := &client.CreateMockFileParams{
		TtlSeconds: ttlSeconds,
	}

	resp, err := apiClient.CreateMockFileWithBodyWithResponse(ctx, params, contentType, body)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Request failed: %v", err)), nil
	}

	if resp.StatusCode() != 201 && resp.StatusCode() != 200 {
		return formatError(resp.StatusCode(), resp.Body), nil
	}

	b, _ := json.MarshalIndent(resp.JSON201, "", "  ")
	return mcp.NewToolResultText(fmt.Sprintf("Mock created successfully:\n```json\n%s\n```", string(b))), nil
}

func listMocksHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	limit := int(request.GetFloat("limit", 5))

	params := &client.ListMocksParams{
		Limit: &limit,
	}

	resp, err := apiClient.ListMocksWithResponse(ctx, params)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Request failed: %v", err)), nil
	}

	if resp.StatusCode() != 200 {
		return formatError(resp.StatusCode(), resp.Body), nil
	}

	b, _ := json.MarshalIndent(resp.JSON200, "", "  ")
	return mcp.NewToolResultText(fmt.Sprintf("Mocks:\n```json\n%s\n```", string(b))), nil
}

func getMockStatsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	mockId, err := request.RequireString("mock_id")
	if err != nil {
		return mcp.NewToolResultError("mock_id must be a string"), nil
	}

	resp, err := apiClient.GetMockStatsWithResponse(ctx, mockId, &client.GetMockStatsParams{})
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Request failed: %v", err)), nil
	}

	if resp.StatusCode() != 200 {
		return formatError(resp.StatusCode(), resp.Body), nil
	}

	b, _ := json.MarshalIndent(resp.JSON200, "", "  ")
	return mcp.NewToolResultText(fmt.Sprintf("Stats:\n```json\n%s\n```", string(b))), nil
}

func deleteMockHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	mockId, err := request.RequireString("mock_id")
	if err != nil {
		return mcp.NewToolResultError("mock_id must be a string"), nil
	}

	resp, err := apiClient.DeleteMockWithResponse(ctx, mockId)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Request failed: %v", err)), nil
	}

	if resp.StatusCode() != 204 && resp.StatusCode() != 200 {
		return formatError(resp.StatusCode(), resp.Body), nil
	}

	return mcp.NewToolResultText("Mock deleted successfully"), nil
}
